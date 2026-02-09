package contract

import (
	"context"
	"fmt"
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
	CreateContract(ctx context.Context, req CreateContractRequest, clientID uuid.UUID) (*CreateContractResponse, error)
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

func (s *contractService) CreateContract(ctx context.Context, req CreateContractRequest, clientID uuid.UUID) (*CreateContractResponse, error) {
	contract, err := s.Store.CreateContract(ctx, db.CreateContractParams{
		TypeID:          req.TypeID,
		StartDate:       pgtype.Timestamptz{Time: req.StartDate, Valid: true},
		EndDate:         pgtype.Timestamptz{Time: req.EndDate, Valid: true},
		ReminderPeriod:  req.ReminderPeriod,
		Vat:             req.Vat,
		Price:           req.Price,
		PriceTimeUnit:   db.PriceTimeUnitEnum(req.PriceTimeUnit),
		Hours:           req.Hours,
		HoursType:       db.NullHoursTypeFromPtr(req.HoursType),
		CareName:        req.CareName,
		CareType:        db.CareTypeEnum(req.CareType),
		ClientID:        clientID,
		SenderID:        req.SenderID,
		Status:          "draft",
		AttachmentIds:   req.AttachmentIds,
		FinancingAct:    db.FinancingActEnum(req.FinancingAct),
		FinancingOption: db.FinancingOptionEnum(req.FinancingOption),
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateContract", "Failed to create contract", zap.String("client_id", clientID.String()), zap.Error(err))
		return nil, err
	}

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

	contract, err := qtx.UpdateContract(ctx, db.UpdateContractParams{
		ID:              contractID,
		TypeID:          req.TypeID,
		StartDate:       pgtype.Timestamptz{Time: req.StartDate, Valid: true},
		EndDate:         pgtype.Timestamptz{Time: req.EndDate, Valid: true},
		ReminderPeriod:  req.ReminderPeriod,
		VAT:             req.Vat,
		Price:           req.Price,
		PriceTimeUnit:   db.NullPriceTimeUnitFromPtr(req.PriceTimeUnit),
		Hours:           req.Hours,
		HoursType:       db.NullHoursTypeFromPtr(req.HoursType),
		CareName:        req.CareName,
		CareType:        db.NullCareTypeFromPtr(req.CareType),
		SenderID:        req.SenderID,
		AttachmentIds:   req.AttachmentIds,
		FinancingAct:    db.NullFinancingActFromPtr(req.FinancingAct),
		FinancingOption: db.NullFinancingOptionFromPtr(req.FinancingOption),
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

	response := &GetClientContractResponse{
		ID:              contract.ID,
		TypeID:          contract.TypeID,
		TypeName:        contract.ContractTypeName,
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
		FinancingOption: func() []db.FinancingOptionEnum {
			if req.FinancingOption == nil {
				return nil
			}
			options := make([]db.FinancingOptionEnum, len(req.FinancingOption))
			for i, option := range req.FinancingOption {
				options[i] = db.FinancingOptionEnum(option)
			}
			return options
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
		contractsRes[i] = ListContractsResponse{
			ID:              contract.ID,
			ClientID:        contract.ClientID,
			Status:          string(contract.Status),
			StartDate:       contract.StartDate.Time,
			EndDate:         contract.EndDate.Time,
			Price:           contract.Price,
			PriceTimeUnit:   string(contract.PriceTimeUnit),
			CareName:        contract.CareName,
			CareType:        string(contract.CareType),
			FinancingAct:    string(contract.FinancingAct),
			FinancingOption: string(contract.FinancingOption),
			SenderID:        contract.SenderID,
			SenderName:      contract.SenderName,
			ClientFirstName: contract.ClientFirstName,
			ClientLastName:  contract.ClientLastName,
			CreatedAt:       contract.CreatedAt.Time,
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
