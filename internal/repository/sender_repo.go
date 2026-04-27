package repository

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/google/uuid"

	db "maicare_go/db/sqlc"
	"maicare_go/internal/domain"
)

type SenderRepository struct {
	store *db.Store
}

func NewSenderRepository(store *db.Store) *SenderRepository {
	return &SenderRepository{store: store}
}

func (r *SenderRepository) Create(ctx context.Context, params domain.CreateSenderParams) (*domain.SenderEntity, error) {
	contactsJSON, err := json.Marshal(params.Contacts)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal contacts: %w", err)
	}

	sender, err := r.store.CreateSender(ctx, db.CreateSenderParams{
		Types:               db.SenderTypesEnum(params.Types),
		Name:                params.Name,
		Street:              params.Street,
		HouseNumber:         params.HouseNumber,
		HouseNumberAddition: params.HouseNumberAddition,
		PostalCode:          params.PostalCode,
		City:                params.City,
		Land:                params.Land,
		Kvknumber:           params.KVKNumber,
		Btwnumber:           params.BTWNumber,
		PhoneNumber:         params.PhoneNumber,
		ClientNumber:        params.ClientNumber,
		Contacts:            contactsJSON,
	})
	if err != nil {
		return nil, err
	}

	return mapSenderDBToDomain(sender), nil
}

func (r *SenderRepository) List(ctx context.Context, params domain.ListSendersParams) ([]domain.SenderEntity, int64, error) {
	senders, err := r.store.ListSenders(ctx, db.ListSendersParams{
		Limit:           params.Limit,
		Offset:          params.Offset,
		IncludeArchived: params.IncludeArchived,
		Search:          params.Search,
	})
	if err != nil {
		return nil, 0, err
	}

	totalCount, err := r.store.CountSenders(ctx, params.IncludeArchived)
	if err != nil {
		return nil, 0, err
	}

	result := make([]domain.SenderEntity, len(senders))
	for i, s := range senders {
		result[i] = *mapSenderDBToDomain(s)
	}

	return result, totalCount, nil
}

func (r *SenderRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.SenderEntity, error) {
	sender, err := r.store.GetSenderById(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapSenderDBToDomain(sender), nil
}

func (r *SenderRepository) Update(ctx context.Context, params domain.UpdateSenderParams) (*domain.SenderEntity, error) {
	arg := db.UpdateSenderParams{
		ID:                  params.ID,
		Name:                params.Name,
		Street:              params.Street,
		HouseNumber:         params.HouseNumber,
		HouseNumberAddition: params.HouseNumberAddition,
		PostalCode:          params.PostalCode,
		City:                params.City,
		Land:                params.Land,
		Kvknumber:           params.KVKNumber,
		Btwnumber:           params.BTWNumber,
		PhoneNumber:         params.PhoneNumber,
		ClientNumber:        params.ClientNumber,
		EmailAddress:        params.EmailAddress,
		IsArchived:          params.IsArchived,
		Types:               db.NullSenderTypesFromPtr(params.Types),
	}

	if params.Contacts != nil {
		contactsJSON, err := json.Marshal(params.Contacts)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal contacts: %w", err)
		}
		arg.Contacts = contactsJSON
	}

	sender, err := r.store.UpdateSender(ctx, arg)
	if err != nil {
		return nil, err
	}

	return mapSenderDBToDomain(sender), nil
}

func (r *SenderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.store.DeleteSender(ctx, id)
}

func (r *SenderRepository) CreateInvoiceTemplate(ctx context.Context, senderID uuid.UUID, templateIDs []uuid.UUID) error {
	_, err := r.store.CreateSenderInvoiceTemplate(ctx, db.CreateSenderInvoiceTemplateParams{
		ID:              senderID,
		InvoiceTemplate: templateIDs,
	})
	return err
}

func (r *SenderRepository) GetTemplateItemsBySourceTable(ctx context.Context, ids []uuid.UUID) ([]domain.InvoiceTemplateItem, error) {
	items, err := r.store.GetTemplateItemsBySourceTable(ctx, ids)
	if err != nil {
		return nil, err
	}

	result := make([]domain.InvoiceTemplateItem, len(items))
	for i, item := range items {
		result[i] = domain.InvoiceTemplateItem{
			ID:           item.ID,
			ItemTag:      item.ItemTag,
			Description:  item.Description,
			SourceTable:  item.SourceTable,
			SourceColumn: item.SourceColumn,
		}
	}

	return result, nil
}

func (r *SenderRepository) GetTemplateItemsByIDs(ctx context.Context, ids []uuid.UUID) ([]uuid.UUID, error) {
	return r.store.GetTemplateItemsByIds(ctx, ids)
}

// ==================== Helpers ====================

func mapSenderDBToDomain(s db.Sender) *domain.SenderEntity {
	contacts := make([]domain.SenderContactInfo, 0)
	if len(s.Contacts) > 0 {
		_ = json.Unmarshal(s.Contacts, &contacts)
	}

	return &domain.SenderEntity{
		ID:                  s.ID,
		Types:               string(s.Types),
		Name:                s.Name,
		Street:              s.Street,
		HouseNumber:         s.HouseNumber,
		HouseNumberAddition: s.HouseNumberAddition,
		PostalCode:          s.PostalCode,
		City:                s.City,
		Land:                s.Land,
		KVKNumber:           s.Kvknumber,
		BTWNumber:           s.Btwnumber,
		PhoneNumber:         s.PhoneNumber,
		ClientNumber:        s.ClientNumber,
		EmailAddress:        s.EmailAddress,
		Contacts:            contacts,
		InvoiceTemplate:     s.InvoiceTemplate,
		IsArchived:          s.IsArchived,
		CreatedAt:           s.CreatedAt.Time,
		UpdatedAt:           s.UpdatedAt.Time,
	}
}
