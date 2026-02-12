package contract

import (
	"context"
	"fmt"
	"strings"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"
	"maicare_go/service/deps"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type ContractService interface {
	CreateContractType(ctx context.Context, req CreateContractTypeRequest) (*CreateContractTypeResponse, error)
	ListContractTypes(ctx context.Context) ([]ListContractTypesResponse, error)
	DeleteContractType(ctx context.Context, contractTypeID uuid.UUID) (*DeleteContractTypeResponse, error)
	CreateContract(ctx context.Context, req CreateContractRequest) (*CreateContractResponse, error)
	ListClientContracts(ctx *gin.Context, req ListClientContractsRequest, clientID uuid.UUID) (*pagination.Response[ListClientContractsResponse], error)
	UpdateContract(ctx context.Context, req UpdateContractRequest, contractID uuid.UUID, employeeID uuid.UUID) (*UpdateContractResponse, error)
	UpdateContractStatus(ctx context.Context, req UpdateContractStatusRequest, contractID uuid.UUID, employeeID uuid.UUID) (*UpdateContractStatusResponse, error)
	GetClientContract(ctx context.Context, contractID uuid.UUID) (*GetClientContractResponse, error)
	ListContracts(ctx *gin.Context, req ListContractsRequest) (*pagination.Response[ListContractsResponse], error)
	GetContractAuditLog(ctx context.Context, contractID uuid.UUID) ([]GetContractAuditLogResponse, error)
}

type contractService struct {
	*deps.ServiceDependencies
}

const (
	defaultContractReminderPeriod int32 = 90
	defaultContractVAT            int32 = 20
)

func normalizeAttachmentIDs(attachmentIDs []uuid.UUID) []uuid.UUID {
	if len(attachmentIDs) == 0 {
		return []uuid.UUID{}
	}

	normalized := make([]uuid.UUID, 0, len(attachmentIDs))
	seen := make(map[uuid.UUID]struct{}, len(attachmentIDs))
	for _, id := range attachmentIDs {
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		normalized = append(normalized, id)
	}

	return normalized
}

func (s *contractService) validateAttachmentIds(ctx context.Context, attachmentIds []uuid.UUID) error {
	if len(attachmentIds) == 0 {
		return nil
	}

	attachments, err := s.Store.GetAttachmentsByUUIDs(ctx, attachmentIds)
	if err != nil {
		return fmt.Errorf("failed to validate attachment IDs: %w", err)
	}

	found := make(map[uuid.UUID]bool, len(attachments))
	for _, a := range attachments {
		found[a.Uuid] = true
	}
	var missing []string
	for _, id := range attachmentIds {
		if !found[id] {
			missing = append(missing, id.String())
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("attachment IDs not found: %v", missing)
	}

	for _, a := range attachments {
		if !a.IsUsed {
			return fmt.Errorf("attachment %s has not been confirmed (upload not completed)", a.Uuid.String())
		}
	}

	return nil
}

func (s *contractService) fetchAttachmentDetails(ctx context.Context, attachmentIds []uuid.UUID) []AttachmentDetail {
	if len(attachmentIds) == 0 {
		return []AttachmentDetail{}
	}

	attachments, err := s.Store.GetAttachmentsByUUIDs(ctx, attachmentIds)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "fetchAttachmentDetails", "Failed to fetch attachments", zap.Error(err))
		return []AttachmentDetail{}
	}

	details := make([]AttachmentDetail, 0, len(attachments))
	for _, a := range attachments {
		detail := AttachmentDetail{
			ID:   a.Uuid,
			Name: a.Name,
			Size: int64(a.Size),
		}
		url := s.GenerateResponsePresignedURL(&a.File, ctx)
		if url != nil {
			detail.DownloadURL = *url
		}
		details = append(details, detail)
	}

	return details
}

func validateContractCarePricing(careType string, priceTimeUnit string, hours *float64, hoursType *string) error {
	switch careType {
	case string(db.CareTypeEnumAmbulante):
		if priceTimeUnit != string(db.PriceTimeUnitEnumMinute) && priceTimeUnit != string(db.PriceTimeUnitEnumHourly) {
			return fmt.Errorf("ambulante contracts require price_time_unit to be minute or hourly")
		}
		if hours == nil || *hours <= 0 {
			return fmt.Errorf("ambulante contracts require hours to be greater than 0")
		}
		if hoursType == nil {
			return fmt.Errorf("ambulante contracts require hours_type")
		}
	case string(db.CareTypeEnumAccommodation):
		if priceTimeUnit != string(db.PriceTimeUnitEnumDaily) && priceTimeUnit != string(db.PriceTimeUnitEnumWeekly) {
			return fmt.Errorf("accommodation contracts require price_time_unit to be daily or weekly")
		}
		if hours != nil || hoursType != nil {
			return fmt.Errorf("accommodation contracts require hours and hours_type to be null")
		}
	default:
		return fmt.Errorf("invalid care_type")
	}

	return nil
}

func NewContractService(deps *deps.ServiceDependencies) ContractService {
	return &contractService{
		ServiceDependencies: deps,
	}
}

func (s *contractService) CreateContractType(ctx context.Context, req CreateContractTypeRequest) (*CreateContractTypeResponse, error) {
	contractType, err := s.Store.CreateContractType(ctx, req.Name)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateContractType", "Failed to create contract type", zap.String("name", req.Name), zap.Error(err))
		return nil, err
	}

	response := &CreateContractTypeResponse{
		ID:   contractType.ID,
		Name: contractType.Name,
	}
	return response, nil
}

func (s *contractService) ListContractTypes(ctx context.Context) ([]ListContractTypesResponse, error) {
	contractTypes, err := s.Store.ListContractTypes(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListContractTypes", "Failed to list contract types", zap.Error(err))
		return nil, err
	}

	if len(contractTypes) == 0 {
		return []ListContractTypesResponse{}, nil
	}

	contractTypesRes := make([]ListContractTypesResponse, len(contractTypes))
	for i, contractType := range contractTypes {
		contractTypesRes[i] = ListContractTypesResponse{
			ID:   contractType.ID,
			Name: contractType.Name,
		}
	}

	return contractTypesRes, nil
}

func (s *contractService) DeleteContractType(ctx context.Context, contractTypeID uuid.UUID) (*DeleteContractTypeResponse, error) {
	err := s.Store.DeleteContractType(ctx, contractTypeID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteContractType", "Failed to delete contract type", zap.String("contract_type_id", contractTypeID.String()), zap.Error(err))
		return nil, err
	}

	return &DeleteContractTypeResponse{ID: contractTypeID}, nil
}

func (s *contractService) CreateContract(ctx context.Context, req CreateContractRequest) (*CreateContractResponse, error) {
	careName := strings.TrimSpace(req.CareName)
	if careName == "" {
		return nil, fmt.Errorf("care_name is required")
	}

	if req.Price <= 0 {
		return nil, fmt.Errorf("price must be greater than 0")
	}

	if req.EndDate.Before(req.StartDate) || req.EndDate.Equal(req.StartDate) {
		return nil, fmt.Errorf("end_date must be after start_date")
	}

	if req.ReminderPeriod != nil && *req.ReminderPeriod < 0 {
		return nil, fmt.Errorf("reminder_period must be greater than or equal to 0")
	}

	if err := validateContractCarePricing(req.CareType, req.PriceTimeUnit, req.Hours, req.HoursType); err != nil {
		return nil, err
	}

	attachmentIDs := normalizeAttachmentIDs(req.AttachmentIds)
	if err := s.validateAttachmentIds(ctx, attachmentIDs); err != nil {
		return nil, err
	}

	reminderPeriod := defaultContractReminderPeriod
	if req.ReminderPeriod != nil {
		reminderPeriod = *req.ReminderPeriod
	}

	vat := req.Vat
	if vat == nil {
		defaultVAT := defaultContractVAT
		vat = &defaultVAT
	}
	if *vat < 0 || *vat > 100 {
		return nil, fmt.Errorf("VAT must be between 0 and 100")
	}

	var contract db.Contract
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		createdContract, createErr := q.CreateContract(ctx, db.CreateContractParams{
			TypeID:          req.TypeID,
			StartDate:       pgtype.Timestamptz{Time: req.StartDate, Valid: true},
			EndDate:         pgtype.Timestamptz{Time: req.EndDate, Valid: true},
			ReminderPeriod:  reminderPeriod,
			Vat:             vat,
			Price:           req.Price,
			PriceTimeUnit:   db.PriceTimeUnitEnum(req.PriceTimeUnit),
			Hours:           req.Hours,
			HoursType:       db.NullHoursTypeFromPtr(req.HoursType),
			CareName:        careName,
			CareType:        db.CareTypeEnum(req.CareType),
			ClientID:        req.ClientID,
			SenderID:        req.SenderID,
			Status:          "draft",
			AttachmentIds:   attachmentIDs,
			FinancingAct:    db.FinancingActEnum(req.FinancingAct),
			FinancingOption: db.FinancingOptionEnum(req.FinancingOption),
		})
		if createErr != nil {
			return createErr
		}

		contract = createdContract
		return nil
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateContract", "Failed to create contract", zap.String("client_id", req.ClientID.String()), zap.Error(err))
		return nil, err
	}

	attachmentDetails := s.fetchAttachmentDetails(ctx, contract.AttachmentIds)

	response := &CreateContractResponse{
		ID:              contract.ID,
		TypeID:          contract.TypeID,
		Status:          string(contract.Status),
		StartDate:       contract.StartDate.Time,
		EndDate:         contract.EndDate.Time,
		ReminderPeriod:  contract.ReminderPeriod,
		Vat:             contract.Vat,
		Price:           contract.Price,
		PriceTimeUnit:   string(contract.PriceTimeUnit),
		Hours:           contract.Hours,
		HoursType:       db.HoursTypePtrFromEnum(contract.HoursType),
		CareName:        contract.CareName,
		CareType:        string(contract.CareType),
		ClientID:        contract.ClientID,
		SenderID:        contract.SenderID,
		AttachmentIds:   contract.AttachmentIds,
		Attachments:     attachmentDetails,
		FinancingAct:    string(contract.FinancingAct),
		FinancingOption: string(contract.FinancingOption),
		DepartureReason: contract.DepartureReason,
		DepartureReport: contract.DepartureReport,
		UpdatedAt:       contract.UpdatedAt,
		CreatedAt:       contract.CreatedAt,
	}
	return response, nil
}

func (s *contractService) ListClientContracts(ctx *gin.Context, req ListClientContractsRequest, clientID uuid.UUID) (*pagination.Response[ListClientContractsResponse], error) {
	params := req.GetParams()

	contracts, err := s.Store.ListClientContracts(ctx, db.ListClientContractsParams{
		ClientID: clientID,
		Limit:    params.Limit,
		Offset:   params.Offset,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListClientContracts", "Failed to list client contracts", zap.String("client_id", clientID.String()), zap.Error(err))
		return nil, err
	}

	if len(contracts) == 0 {
		pag := pagination.NewResponse(ctx, req.Request, []ListClientContractsResponse{}, 0)
		return &pag, nil
	}

	totalCount := contracts[0].TotalCount

	contractsRes := make([]ListClientContractsResponse, len(contracts))
	for i, contract := range contracts {
		contractsRes[i] = ListClientContractsResponse{
			ID:              contract.ID,
			TypeID:          contract.TypeID,
			Status:          string(contract.Status),
			StartDate:       contract.StartDate.Time,
			EndDate:         contract.EndDate.Time,
			ReminderPeriod:  contract.ReminderPeriod,
			Vat:             contract.Vat,
			Price:           contract.Price,
			PriceTimeUnit:   string(contract.PriceTimeUnit),
			Hours:           contract.Hours,
			HoursType:       db.HoursTypePtrFromEnum(contract.HoursType),
			CareName:        contract.CareName,
			CareType:        string(contract.CareType),
			ClientID:        contract.ClientID,
			ClientFirstName: contract.ClientFirstName,
			ClientLastName:  contract.ClientLastName,
			SenderID:        contract.SenderID,
			SenderName:      contract.SenderName,
			AttachmentIds:   contract.AttachmentIds,
			FinancingAct:    string(contract.FinancingAct),
			FinancingOption: string(contract.FinancingOption),
			DepartureReason: contract.DepartureReason,
			DepartureReport: contract.DepartureReport,
			UpdatedAt:       contract.UpdatedAt.Time,
			CreatedAt:       contract.CreatedAt.Time,
		}
	}

	pag := pagination.NewResponse(ctx, req.Request, contractsRes, totalCount)
	return &pag, nil
}

func (s *contractService) UpdateContract(ctx context.Context, req UpdateContractRequest, contractID uuid.UUID, employeeID uuid.UUID) (*UpdateContractResponse, error) {
	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateContract", "Failed to begin transaction", zap.String("contract_id", contractID.String()), zap.Error(err))
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := s.Store.WithTx(tx)

	_, err = tx.Exec(ctx, "SET LOCAL myapp.current_employee_id = $1", employeeID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateContract", "Failed to set current employee ID", zap.String("contract_id", contractID.String()), zap.String("employee_id", employeeID.String()), zap.Error(err))
		return nil, err
	}

	existingContract, err := qtx.GetClientContract(ctx, contractID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateContract", "Failed to get contract", zap.String("contract_id", contractID.String()), zap.Error(err))
		return nil, err
	}

	startDate := existingContract.StartDate.Time
	if req.StartDate != nil {
		startDate = *req.StartDate
	}

	endDate := existingContract.EndDate.Time
	if req.EndDate != nil {
		endDate = *req.EndDate
	}

	if endDate.Before(startDate) || endDate.Equal(startDate) {
		return nil, fmt.Errorf("end_date must be after start_date")
	}

	reminderPeriod := existingContract.ReminderPeriod
	if req.ReminderPeriod != nil {
		reminderPeriod = *req.ReminderPeriod
	}

	vat := existingContract.Vat
	if req.Vat != nil {
		vat = req.Vat
	}

	price := existingContract.Price
	if req.Price != nil {
		price = *req.Price
	}

	priceTimeUnit := string(existingContract.PriceTimeUnit)
	if req.PriceTimeUnit != nil {
		priceTimeUnit = *req.PriceTimeUnit
	}

	hours := existingContract.Hours
	if req.Hours != nil {
		hours = req.Hours
	}

	hoursType := db.HoursTypePtrFromEnum(existingContract.HoursType)
	if req.HoursType != nil {
		hoursType = req.HoursType
	}

	careName := existingContract.CareName
	if req.CareName != nil {
		careName = *req.CareName
	}

	careType := string(existingContract.CareType)
	if req.CareType != nil {
		careType = *req.CareType
	}

	if careType == string(db.CareTypeEnumAccommodation) {
		hours = nil
		hoursType = nil
	}

	if err := validateContractCarePricing(careType, priceTimeUnit, hours, hoursType); err != nil {
		return nil, err
	}

	typeID := existingContract.TypeID
	if req.TypeID != nil {
		typeID = req.TypeID
	}

	senderID := existingContract.SenderID
	if req.SenderID != nil {
		senderID = *req.SenderID
	}

	attachmentIDs := existingContract.AttachmentIds
	if req.AttachmentIds != nil {
		attachmentIDs = req.AttachmentIds
	}

	if err := s.validateAttachmentIds(ctx, attachmentIDs); err != nil {
		return nil, err
	}

	financingAct := string(existingContract.FinancingAct)
	if req.FinancingAct != nil {
		financingAct = *req.FinancingAct
	}

	financingOption := string(existingContract.FinancingOption)
	if req.FinancingOption != nil {
		financingOption = *req.FinancingOption
	}

	contract, err := qtx.UpdateContract(ctx, db.UpdateContractParams{
		ID:              contractID,
		TypeID:          typeID,
		StartDate:       pgtype.Timestamptz{Time: startDate, Valid: true},
		EndDate:         pgtype.Timestamptz{Time: endDate, Valid: true},
		ReminderPeriod:  reminderPeriod,
		Vat:             vat,
		Price:           price,
		PriceTimeUnit:   db.PriceTimeUnitEnum(priceTimeUnit),
		Hours:           hours,
		HoursType:       db.NullHoursTypeFromPtr(hoursType),
		CareName:        careName,
		CareType:        db.CareTypeEnum(careType),
		SenderID:        senderID,
		AttachmentIds:   attachmentIDs,
		FinancingAct:    db.FinancingActEnum(financingAct),
		FinancingOption: db.FinancingOptionEnum(financingOption),
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateContract", "Failed to update contract", zap.String("contract_id", contractID.String()), zap.Error(err))
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateContract", "Failed to commit transaction", zap.String("contract_id", contractID.String()), zap.Error(err))
		return nil, err
	}

	response := &UpdateContractResponse{
		ID:              contract.ID,
		TypeID:          contract.TypeID,
		Status:          string(contract.Status),
		StartDate:       contract.StartDate.Time,
		EndDate:         contract.EndDate.Time,
		ReminderPeriod:  contract.ReminderPeriod,
		Vat:             contract.Vat,
		Price:           contract.Price,
		PriceFrequency:  string(contract.PriceTimeUnit),
		Hours:           contract.Hours,
		HoursType:       db.HoursTypePtrFromEnum(contract.HoursType),
		CareName:        contract.CareName,
		CareType:        string(contract.CareType),
		ClientID:        contract.ClientID,
		SenderID:        contract.SenderID,
		AttachmentIds:   contract.AttachmentIds,
		FinancingAct:    string(contract.FinancingAct),
		FinancingOption: string(contract.FinancingOption),
		DepartureReason: contract.DepartureReason,
		DepartureReport: contract.DepartureReport,
		UpdatedAt:       contract.UpdatedAt.Time,
		CreatedAt:       contract.CreatedAt.Time,
	}
	return response, nil
}

func (s *contractService) UpdateContractStatus(ctx context.Context, req UpdateContractStatusRequest, contractID uuid.UUID, employeeID uuid.UUID) (*UpdateContractStatusResponse, error) {
	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateContractStatus", "Failed to begin transaction", zap.String("contract_id", contractID.String()), zap.Error(err))
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := s.Store.WithTx(tx)

	_, err = tx.Exec(ctx, "SET LOCAL myapp.current_employee_id = $1", employeeID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateContractStatus", "Failed to set current employee ID", zap.String("contract_id", contractID.String()), zap.String("employee_id", employeeID.String()), zap.Error(err))
		return nil, err
	}

	contract, err := qtx.GetClientContract(ctx, contractID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateContractStatus", "Failed to get client contract", zap.String("contract_id", contractID.String()), zap.Error(err))
		return nil, err
	}

	if req.Status == "approved" && contract.EndDate.Time.Before(time.Now()) {
		err := fmt.Errorf("cannot approve contract that has already ended")
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateContractStatus", "Cannot approve ended contract", zap.String("contract_id", contractID.String()), zap.String("status", req.Status), zap.Error(err))
		return nil, err
	}

	updatedContract, err := qtx.UpdateContractStatus(ctx, db.UpdateContractStatusParams{
		ContractID: contractID,
		Status:     db.ContractStatusEnum(req.Status),
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateContractStatus", "Failed to update contract status", zap.String("contract_id", contractID.String()), zap.String("status", req.Status), zap.Error(err))
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateContractStatus", "Failed to commit transaction", zap.String("contract_id", contractID.String()), zap.Error(err))
		return nil, err
	}

	response := &UpdateContractStatusResponse{
		ID:     updatedContract.ID,
		Status: string(updatedContract.Status),
	}
	return response, nil
}

func (s *contractService) GetClientContract(ctx context.Context, contractID uuid.UUID) (*GetClientContractResponse, error) {
	contract, err := s.Store.GetClientContract(ctx, contractID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetClientContract", "Failed to get client contract", zap.String("contract_id", contractID.String()), zap.Error(err))
		return nil, err
	}

	var approvedAt *time.Time
	if contract.ApprovedAt.Valid {
		approvedAt = &contract.ApprovedAt.Time
	}

	response := &GetClientContractResponse{
		ID:              contract.ID,
		TypeID:          contract.TypeID,
		TypeName:        contract.ContractTypeName,
		Status:          string(contract.Status),
		ApprovedAt:      approvedAt,
		StartDate:       contract.StartDate.Time,
		EndDate:         contract.EndDate.Time,
		ReminderPeriod:  contract.ReminderPeriod,
		Vat:             contract.Vat,
		Price:           contract.Price,
		PriceTimeUnit:   string(contract.PriceTimeUnit),
		Hours:           contract.Hours,
		HoursType:       db.HoursTypePtrFromEnum(contract.HoursType),
		CareName:        contract.CareName,
		CareType:        string(contract.CareType),
		AttachmentIds:   contract.AttachmentIds,
		FinancingAct:    string(contract.FinancingAct),
		FinancingOption: string(contract.FinancingOption),
		DepartureReason: contract.DepartureReason,
		DepartureReport: contract.DepartureReport,
		UpdatedAt:       contract.UpdatedAt.Time,
		CreatedAt:       contract.CreatedAt.Time,

		ClientID:         contract.ClientID,
		ClientFirstName:  contract.ClientFirstName,
		ClientLastName:   contract.ClientLastName,
		ClientFilenumber: contract.ClientFilenumber,
		ClientBsn:        contract.ClientBsn,

		SenderID:                  contract.SenderID,
		SenderName:                contract.SenderName,
		SenderType:                string(contract.SenderType),
		SenderStreet:              contract.SenderStreet,
		SenderHouseNumber:         contract.SenderHouseNumber,
		SenderHouseNumberAddition: contract.SenderHouseNumberAddition,
		SenderPostalCode:          contract.SenderPostalCode,
		SenderCity:                contract.SenderCity,
		SenderLand:                contract.SenderLand,
		SenderKvknumber:           contract.SenderKvknumber,
		SenderBtwnumber:           contract.SenderBtwnumber,
		SenderPhoneNumber:         contract.SenderPhoneNumber,
		SenderClientNumber:        contract.SenderClientNumber,
		SenderEmailAddress:        contract.SenderEmailAddress,
	}
	return response, nil
}

func (s *contractService) ListContracts(ctx *gin.Context, req ListContractsRequest) (*pagination.Response[ListContractsResponse], error) {
	params := req.GetParams()

	contracts, err := s.Store.ListContracts(ctx, db.ListContractsParams{
		Limit:  params.Limit,
		Offset: params.Offset,
		Search: req.Search,
		Status: func() []db.ContractStatusEnum {
			if req.Status == nil {
				return nil
			}
			statuses := make([]db.ContractStatusEnum, len(req.Status))
			for i, status := range req.Status {
				statuses[i] = db.ContractStatusEnum(status)
			}
			return statuses
		}(),
		CareType: func() []db.CareTypeEnum {
			if req.CareType == nil {
				return nil
			}
			types := make([]db.CareTypeEnum, len(req.CareType))
			for i, careType := range req.CareType {
				types[i] = db.CareTypeEnum(careType)
			}
			return types
		}(),
		FinancingAct: func() []db.FinancingActEnum {
			if req.FinancingAct == nil {
				return nil
			}
			acts := make([]db.FinancingActEnum, len(req.FinancingAct))
			for i, act := range req.FinancingAct {
				acts[i] = db.FinancingActEnum(act)
			}
			return acts
		}(),
		FinancingOption: func() db.NullFinancingOptionEnum {
			if req.FinancingOption == nil {
				return db.NullFinancingOptionEnum{Valid: false}
			}
			return db.NullFinancingOptionEnum{
				FinancingOptionEnum: db.FinancingOptionEnum(*req.FinancingOption),
				Valid:               true,
			}
		}(),
		EndDateFrom: func() pgtype.Timestamptz {
			if req.EndDateFrom == nil {
				return pgtype.Timestamptz{Valid: false}
			}
			return pgtype.Timestamptz{Time: *req.EndDateFrom, Valid: true}
		}(),
		EndDateTo: func() pgtype.Timestamptz {
			if req.EndDateTo == nil {
				return pgtype.Timestamptz{Valid: false}
			}
			return pgtype.Timestamptz{Time: *req.EndDateTo, Valid: true}
		}(),
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListContracts", "Failed to list contracts", zap.Error(err))
		return nil, err
	}

	if len(contracts) == 0 {
		pag := pagination.NewResponse(ctx, req.Request, []ListContractsResponse{}, 0)
		return &pag, nil
	}

	totalCount := contracts[0].TotalCount

	contractsRes := make([]ListContractsResponse, len(contracts))
	for i, contract := range contracts {
		var approvedAt *time.Time
		if contract.ApprovedAt.Valid {
			approvedAt = &contract.ApprovedAt.Time
		}

		contractsRes[i] = ListContractsResponse{
			ID:               contract.ID,
			ClientID:         contract.ClientID,
			ClientFirstName:  contract.ClientFirstName,
			ClientLastName:   contract.ClientLastName,
			ClientFilenumber: contract.ClientFilenumber,
			SenderID:         contract.SenderID,
			SenderName:       contract.SenderName,
			CareName:         contract.CareName,
			CareType:         string(contract.CareType),
			Price:            contract.Price,
			PriceTimeUnit:    string(contract.PriceTimeUnit),
			Hours:            contract.Hours,
			HoursType:        db.HoursTypePtrFromEnum(contract.HoursType),
			FinancingAct:     string(contract.FinancingAct),
			FinancingOption:  string(contract.FinancingOption),
			StartDate:        contract.StartDate.Time,
			EndDate:          contract.EndDate.Time,
			DaysLeft:         contract.DaysLeft,
			Status:           string(contract.Status),
			ApprovedAt:       approvedAt,
			UpdatedAt:        contract.UpdatedAt.Time,
		}
	}

	pag := pagination.NewResponse(ctx, req.Request, contractsRes, totalCount)
	return &pag, nil
}

func (s *contractService) GetContractAuditLog(ctx context.Context, contractID uuid.UUID) ([]GetContractAuditLogResponse, error) {
	auditLogs, err := s.Store.GetContractAudit(ctx, contractID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetContractAuditLog", "Failed to get contract audit logs", zap.String("contract_id", contractID.String()), zap.Error(err))
		return nil, err
	}

	if len(auditLogs) == 0 {
		return []GetContractAuditLogResponse{}, nil
	}

	auditLogsRes := make([]GetContractAuditLogResponse, len(auditLogs))
	for i, log := range auditLogs {
		auditLogsRes[i] = GetContractAuditLogResponse{
			AuditID:            log.AuditID,
			ContractID:         log.ContractID,
			Operation:          string(log.Operation),
			ChangedBy:          log.ChangedBy,
			ChangedAt:          log.ChangedAt,
			OldValues:          util.ParseJSONToObject(log.OldValues),
			NewValues:          util.ParseJSONToObject(log.NewValues),
			ChangedFields:      log.ChangedFields,
			ChangedByFirstName: log.ChangedByFirstName,
			ChangedByLastName:  log.ChangedByLastName,
		}
	}

	return auditLogsRes, nil
}
