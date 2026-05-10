package repository

import (
	"context"
	"fmt"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/internal/domain"
	"maicare_go/pkg/conv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ContractRepository struct {
	store *db.Store
}

func NewContractRepository(store *db.Store) domain.ContractRepository {
	return &ContractRepository{store: store}
}

// ==================== ContractType ====================

func (r *ContractRepository) CreateContractType(ctx context.Context, params domain.CreateContractTypeParams) (*domain.ContractType, error) {
	ct, err := r.store.CreateContractType(ctx, params.Name)
	if err != nil {
		return nil, err
	}
	return &domain.ContractType{ID: ct.ID, Name: ct.Name}, nil
}

func (r *ContractRepository) ListContractTypes(ctx context.Context) ([]domain.ContractType, error) {
	rows, err := r.store.ListContractTypes(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.ContractType, len(rows))
	for i, row := range rows {
		result[i] = domain.ContractType{ID: row.ID, Name: row.Name}
	}
	return result, nil
}

func (r *ContractRepository) DeleteContractType(ctx context.Context, contractTypeID uuid.UUID) error {
	return r.store.DeleteContractType(ctx, contractTypeID)
}

// ==================== Contract CRUD ====================

func (r *ContractRepository) CreateContract(ctx context.Context, params domain.CreateContractParams) (*domain.Contract, error) {
	reminderPeriod := int32(90)
	if params.ReminderPeriod != nil {
		reminderPeriod = *params.ReminderPeriod
	}

	vat := params.Vat
	if vat == nil {
		defaultVAT := int32(20)
		vat = &defaultVAT
	}

	contract, err := r.store.CreateContract(ctx, db.CreateContractParams{
		TypeID:          params.TypeID,
		Status:          db.ContractStatusEnumDraft,
		StartDate:       conv.PgTimestamptzFromTime(params.StartDate),
		EndDate:         conv.PgTimestamptzFromTime(params.EndDate),
		ReminderPeriod:  reminderPeriod,
		Vat:             vat,
		Price:           params.Price,
		PriceTimeUnit:   db.PriceTimeUnitEnum(params.PriceTimeUnit),
		Hours:           params.Hours,
		HoursType:       db.NullHoursTypeFromPtr(params.HoursType),
		CareName:        params.CareName,
		CareType:        db.CareTypeEnum(params.CareType),
		ClientID:        params.ClientID,
		SenderID:        params.SenderID,
		AttachmentIds:   params.AttachmentIds,
		FinancingAct:    db.FinancingActEnum(params.FinancingAct),
		FinancingOption: db.FinancingOptionEnum(params.FinancingOption),
	})
	if err != nil {
		return nil, err
	}
	return toDomainContract(contract), nil
}

func (r *ContractRepository) GetContractByID(ctx context.Context, contractID uuid.UUID) (*domain.ContractDetail, error) {
	row, err := r.store.GetClientContract(ctx, contractID)
	if err != nil {
		return nil, err
	}
	return toDomainContractDetail(row), nil
}

func (r *ContractRepository) UpdateContract(ctx context.Context, params domain.UpdateContractParams, employeeID uuid.UUID) (*domain.Contract, error) {
	tx, err := r.store.ConnPool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, "SET LOCAL myapp.current_employee_id = $1", employeeID); err != nil {
		return nil, fmt.Errorf("failed to set current employee ID: %w", err)
	}

	qtx := r.store.WithTx(tx)

	existing, err := qtx.GetClientContract(ctx, params.ContractID)
	if err != nil {
		return nil, err
	}

	startDate := existing.StartDate.Time
	if params.StartDate != nil {
		startDate = *params.StartDate
	}
	endDate := existing.EndDate.Time
	if params.EndDate != nil {
		endDate = *params.EndDate
	}

	reminderPeriod := existing.ReminderPeriod
	if params.ReminderPeriod != nil {
		reminderPeriod = *params.ReminderPeriod
	}

	vat := existing.Vat
	if params.Vat != nil {
		vat = params.Vat
	}

	price := existing.Price
	if params.Price != nil {
		price = *params.Price
	}

	priceTimeUnit := string(existing.PriceTimeUnit)
	if params.PriceTimeUnit != nil {
		priceTimeUnit = *params.PriceTimeUnit
	}

	hours := existing.Hours
	if params.Hours != nil {
		hours = params.Hours
	}

	hoursType := db.HoursTypePtrFromEnum(existing.HoursType)
	if params.HoursType != nil {
		hoursType = params.HoursType
	}

	careName := existing.CareName
	if params.CareName != nil {
		careName = *params.CareName
	}

	careType := string(existing.CareType)

	typeID := existing.TypeID
	if params.TypeID != nil {
		typeID = params.TypeID
	}

	senderID := existing.SenderID
	if params.SenderID != nil {
		senderID = *params.SenderID
	}

	attachmentIDs := existing.AttachmentIds
	if params.AttachmentIds != nil {
		attachmentIDs = params.AttachmentIds
	}

	financingAct := string(existing.FinancingAct)
	if params.FinancingAct != nil {
		financingAct = *params.FinancingAct
	}

	financingOption := string(existing.FinancingOption)
	if params.FinancingOption != nil {
		financingOption = *params.FinancingOption
	}

	contract, err := qtx.UpdateContract(ctx, db.UpdateContractParams{
		ID:              params.ContractID,
		TypeID:          typeID,
		StartDate:       conv.PgTimestamptzFromTime(startDate),
		EndDate:         conv.PgTimestamptzFromTime(endDate),
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
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return toDomainContract(contract), nil
}

func (r *ContractRepository) UpdateContractStatus(ctx context.Context, params domain.UpdateContractStatusParams, employeeID uuid.UUID) (*domain.Contract, error) {
	tx, err := r.store.ConnPool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, "SET LOCAL myapp.current_employee_id = $1", employeeID); err != nil {
		return nil, fmt.Errorf("failed to set current employee ID: %w", err)
	}

	qtx := r.store.WithTx(tx)

	contract, err := qtx.UpdateContractStatus(ctx, db.UpdateContractStatusParams{
		Status:     db.ContractStatusEnum(params.Status),
		ContractID: params.ContractID,
	})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return toDomainContract(contract), nil
}

func (r *ContractRepository) ListClientContracts(ctx context.Context, params domain.ListClientContractsParams) ([]domain.ClientContractListItem, int64, error) {
	rows, err := r.store.ListClientContracts(ctx, db.ListClientContractsParams{
		ClientID: params.ClientID,
		Limit:    params.Limit,
		Offset:   params.Offset,
	})
	if err != nil {
		return nil, 0, err
	}

	if len(rows) == 0 {
		return []domain.ClientContractListItem{}, 0, nil
	}

	totalCount := rows[0].TotalCount
	result := make([]domain.ClientContractListItem, len(rows))
	for i, row := range rows {
		result[i] = domain.ClientContractListItem{
			StartDate:       row.StartDate.Time,
			EndDate:         row.EndDate.Time,
			DaysLeft:        row.DaysLeft,
			CareName:        row.CareName,
			CareType:        string(row.CareType),
			FinancingAct:    string(row.FinancingAct),
			FinancingOption: string(row.FinancingOption),
		}
	}
	return result, totalCount, nil
}

func (r *ContractRepository) ListContracts(ctx context.Context, params domain.ListContractsParams) ([]domain.ContractListItem, int64, error) {
	rows, err := r.store.ListContracts(ctx, db.ListContractsParams{
		Limit:  params.Limit,
		Offset: params.Offset,
		Search: params.Search,
		EndDateFrom: func() pgtype.Timestamptz {
			if params.EndDateFrom == nil {
				return pgtype.Timestamptz{Valid: false}
			}
			return conv.PgTimestamptzFromTime(*params.EndDateFrom)
		}(),
		EndDateTo: func() pgtype.Timestamptz {
			if params.EndDateTo == nil {
				return pgtype.Timestamptz{Valid: false}
			}
			return conv.PgTimestamptzFromTime(*params.EndDateTo)
		}(),
		Status: func() []db.ContractStatusEnum {
			if len(params.Status) == 0 {
				return nil
			}
			result := make([]db.ContractStatusEnum, len(params.Status))
			for i, s := range params.Status {
				result[i] = db.ContractStatusEnum(s)
			}
			return result
		}(),
		CareType: func() []db.CareTypeEnum {
			if len(params.CareType) == 0 {
				return nil
			}
			result := make([]db.CareTypeEnum, len(params.CareType))
			for i, ct := range params.CareType {
				result[i] = db.CareTypeEnum(ct)
			}
			return result
		}(),
		FinancingAct: func() []db.FinancingActEnum {
			if len(params.FinancingAct) == 0 {
				return nil
			}
			result := make([]db.FinancingActEnum, len(params.FinancingAct))
			for i, fa := range params.FinancingAct {
				result[i] = db.FinancingActEnum(fa)
			}
			return result
		}(),
		FinancingOption: func() *db.FinancingOptionEnum {
			if params.FinancingOption == nil {
				return nil
			}
			value := db.FinancingOptionEnum(*params.FinancingOption)
			return &value
		}(),
	})
	if err != nil {
		return nil, 0, err
	}

	if len(rows) == 0 {
		return []domain.ContractListItem{}, 0, nil
	}

	totalCount := rows[0].TotalCount
	result := make([]domain.ContractListItem, len(rows))
	for i, row := range rows {
		var approvedAt *time.Time
		if row.ApprovedAt.Valid {
			approvedAt = &row.ApprovedAt.Time
		}
		result[i] = domain.ContractListItem{
			ID:               row.ID,
			ClientID:         row.ClientID,
			ClientFirstName:  row.ClientFirstName,
			ClientLastName:   row.ClientLastName,
			ClientFilenumber: row.ClientFilenumber,
			SenderID:         row.SenderID,
			SenderName:       row.SenderName,
			CareName:         row.CareName,
			CareType:         string(row.CareType),
			Price:            row.Price,
			PriceTimeUnit:    string(row.PriceTimeUnit),
			Hours:            row.Hours,
			HoursType:        db.HoursTypePtrFromEnum(row.HoursType),
			FinancingAct:     string(row.FinancingAct),
			FinancingOption:  string(row.FinancingOption),
			StartDate:        row.StartDate.Time,
			EndDate:          row.EndDate.Time,
			DaysLeft:         row.DaysLeft,
			Status:           string(row.Status),
			ApprovedAt:       approvedAt,
			UpdatedAt:        row.UpdatedAt.Time,
		}
	}
	return result, totalCount, nil
}

func (r *ContractRepository) GetAttachmentFiles(ctx context.Context, ids []uuid.UUID) ([]domain.AttachmentFile, error) {
	rows, err := r.store.GetAttachmentsByUUIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	result := make([]domain.AttachmentFile, len(rows))
	for i, row := range rows {
		result[i] = domain.AttachmentFile{
			UUID:      row.Uuid,
			Name:      row.Name,
			File:      row.File,
			Size:      row.Size,
			IsUsed:    row.IsUsed,
			Tag:       row.Tag,
			UpdatedAt: row.Updated.Time,
			CreatedAt: row.Created.Time,
		}
	}
	return result, nil
}

func (r *ContractRepository) GetContractAuditLog(ctx context.Context, contractID uuid.UUID) ([]domain.ContractAuditLog, error) {
	rows, err := r.store.GetContractAudit(ctx, contractID)
	if err != nil {
		return nil, err
	}

	result := make([]domain.ContractAuditLog, len(rows))
	for i, row := range rows {
		result[i] = domain.ContractAuditLog{
			AuditID:            row.AuditID,
			ContractID:         row.ContractID,
			Operation:          string(row.Operation),
			ChangedBy:          row.ChangedBy,
			ChangedAt:          row.ChangedAt.Time,
			OldValues:          row.OldValues,
			NewValues:          row.NewValues,
			ChangedFields:      row.ChangedFields,
			ChangedByFirstName: row.ChangedByFirstName,
			ChangedByLastName:  row.ChangedByLastName,
		}
	}
	return result, nil
}

// ==================== Mappers ====================

func toDomainContract(c db.Contract) *domain.Contract {
	return &domain.Contract{
		ID:              c.ID,
		TypeID:          c.TypeID,
		Status:          string(c.Status),
		StartDate:       c.StartDate.Time,
		EndDate:         c.EndDate.Time,
		ReminderPeriod:  c.ReminderPeriod,
		Vat:             c.Vat,
		Price:           c.Price,
		PriceTimeUnit:   string(c.PriceTimeUnit),
		Hours:           c.Hours,
		HoursType:       db.HoursTypePtrFromEnum(c.HoursType),
		CareName:        c.CareName,
		CareType:        string(c.CareType),
		ClientID:        c.ClientID,
		SenderID:        c.SenderID,
		AttachmentIds:   c.AttachmentIds,
		FinancingAct:    string(c.FinancingAct),
		FinancingOption: string(c.FinancingOption),
		DepartureReason: c.DepartureReason,
		DepartureReport: c.DepartureReport,
		UpdatedAt:       c.UpdatedAt.Time,
		CreatedAt:       c.CreatedAt.Time,
	}
}

func toDomainContractDetail(row db.GetClientContractRow) *domain.ContractDetail {
	var approvedAt *time.Time
	if row.ApprovedAt.Valid {
		approvedAt = &row.ApprovedAt.Time
	}
	return &domain.ContractDetail{
		Contract: domain.Contract{
			ID:              row.ID,
			TypeID:          row.TypeID,
			Status:          string(row.Status),
			ApprovedAt:      approvedAt,
			StartDate:       row.StartDate.Time,
			EndDate:         row.EndDate.Time,
			ReminderPeriod:  row.ReminderPeriod,
			Vat:             row.Vat,
			Price:           row.Price,
			PriceTimeUnit:   string(row.PriceTimeUnit),
			Hours:           row.Hours,
			HoursType:       db.HoursTypePtrFromEnum(row.HoursType),
			CareName:        row.CareName,
			CareType:        string(row.CareType),
			ClientID:        row.ClientID,
			SenderID:        row.SenderID,
			AttachmentIds:   row.AttachmentIds,
			FinancingAct:    string(row.FinancingAct),
			FinancingOption: string(row.FinancingOption),
			DepartureReason: row.DepartureReason,
			DepartureReport: row.DepartureReport,
			UpdatedAt:       row.UpdatedAt.Time,
			CreatedAt:       row.CreatedAt.Time,
		},
		TypeName:                  row.ContractTypeName,
		ClientFirstName:           row.ClientFirstName,
		ClientLastName:            row.ClientLastName,
		ClientFilenumber:          row.ClientFilenumber,
		ClientBsn:                 row.ClientBsn,
		SenderName:                row.SenderName,
		SenderType:                string(row.SenderType),
		SenderStreet:              row.SenderStreet,
		SenderHouseNumber:         row.SenderHouseNumber,
		SenderHouseNumberAddition: row.SenderHouseNumberAddition,
		SenderPostalCode:          row.SenderPostalCode,
		SenderCity:                row.SenderCity,
		SenderLand:                row.SenderLand,
		SenderKvknumber:           row.SenderKvknumber,
		SenderBtwnumber:           row.SenderBtwnumber,
		SenderPhoneNumber:         row.SenderPhoneNumber,
		SenderClientNumber:        row.SenderClientNumber,
		SenderEmailAddress:        row.SenderEmailAddress,
	}
}
