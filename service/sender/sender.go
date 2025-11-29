package sender

import (
	"context"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *senderService) CreateSender(ctx context.Context, req *CreateSenderRequest) (*CreateSenderResponse, error) {
	contactsParam, err := json.Marshal(req.Contacts)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateSender", "Failed to marshal contacts", zap.Error(err))
		return nil, err
	}

	sender, err := s.Store.CreateSender(ctx, db.CreateSenderParams{
		Types:        db.SenderTypesEnum(req.Types),
		Name:         req.Name,
		Address:      req.Address,
		PostalCode:   req.PostalCode,
		Place:        req.Place,
		Land:         req.Land,
		Kvknumber:    req.KVKNumber,
		Btwnumber:    req.BTWNumber,
		PhoneNumber:  req.PhoneNumber,
		ClientNumber: req.ClientNumber,
		Contacts:     contactsParam,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateSender", "Failed to create sender", zap.Error(err))
		return nil, err
	}
	contactsResp := make([]SenderContact, 0)
	if err := json.Unmarshal(sender.Contacts, &contactsResp); err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateSender", "Failed to unmarshal contacts", zap.Error(err))
		return nil, err
	}

	response := &CreateSenderResponse{
		ID:           sender.ID,
		Types:        string(sender.Types),
		Name:         sender.Name,
		Address:      sender.Address,
		PostalCode:   sender.PostalCode,
		Place:        sender.Place,
		Land:         sender.Land,
		KVKNumber:    sender.Kvknumber,
		BTWNumber:    sender.Btwnumber,
		PhoneNumber:  sender.PhoneNumber,
		ClientNumber: sender.ClientNumber,
		Contacts:     contactsResp,
		CreatedAt:    sender.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:    sender.UpdatedAt.Time.Format(time.RFC3339),
	}

	return response, nil
}

func (s *senderService) ListSenders(ctx *gin.Context, req *ListSendersRequest) (*pagination.Response[ListSendersResponse], error) {
	params := req.GetParams()

	arg := db.ListSendersParams{
		Limit:           params.Limit,
		Offset:          params.Offset,
		Search:          req.Search,
		IncludeArchived: req.IncludeArchived,
	}

	senders, err := s.Store.ListSenders(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListSenders", "Failed to list senders", zap.Error(err))
		return nil, err
	}

	totalCount, err := s.Store.CountSenders(ctx, req.IncludeArchived)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListSenders", "Failed to count senders", zap.Error(err))
		return nil, err
	}

	list := []ListSendersResponse{}
	for _, sender := range senders {
		contactsResp := make([]SenderContact, 0)
		if err := json.Unmarshal(sender.Contacts, &contactsResp); err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListSenders", "Failed to unmarshal contacts", zap.Error(err))
			return nil, err
		}

		response := ListSendersResponse{
			ID:           sender.ID,
			Types:        string(sender.Types),
			Name:         sender.Name,
			Address:      sender.Address,
			PostalCode:   sender.PostalCode,
			Place:        sender.Place,
			Land:         sender.Land,
			KVKNumber:    sender.Kvknumber,
			BTWNumber:    sender.Btwnumber,
			PhoneNumber:  sender.PhoneNumber,
			ClientNumber: sender.ClientNumber,
			Contacts:     contactsResp,
			CreatedAt:    sender.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:    sender.UpdatedAt.Time.Format(time.RFC3339),
		}
		list = append(list, response)
	}

	res := pagination.NewResponse(ctx, req.Request, list, totalCount)
	return &res, nil
}

func (s *senderService) GetSenderByID(ctx context.Context, senderID uuid.UUID) (*GetSenderByIdResponse, error) {
	sender, err := s.Store.GetSenderById(ctx, senderID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetSenderByID", "Failed to get sender by ID", zap.Error(err))
		return nil, err
	}

	contactsResp := []SenderContact{}
	if err := json.Unmarshal(sender.Contacts, &contactsResp); err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetSenderByID", "Failed to unmarshal contacts", zap.Error(err))
		return nil, err
	}

	tempIDs, err := s.Store.GetTemplateItemsBySourceTable(ctx, sender.InvoiceTemplate)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetSenderByID", "Failed to get template items", zap.Error(err))
		return nil, err
	}

	templateItems := []TemplateItem{}
	if len(tempIDs) != 0 {
		for _, tmpl := range tempIDs {
			templateItems = append(templateItems, TemplateItem{
				ID:           tmpl.ID,
				ItemTag:      tmpl.ItemTag,
				Description:  tmpl.Description,
				SourceTable:  tmpl.SourceTable,
				SourceColumn: tmpl.SourceColumn,
			})
		}
	}

	response := &GetSenderByIdResponse{
		ID:                   sender.ID,
		Types:                string(sender.Types),
		Name:                 sender.Name,
		Address:              sender.Address,
		PostalCode:           sender.PostalCode,
		Place:                sender.Place,
		Land:                 sender.Land,
		Kvknumber:            sender.Kvknumber,
		Btwnumber:            sender.Btwnumber,
		PhoneNumber:          sender.PhoneNumber,
		ClientNumber:         sender.ClientNumber,
		EmailAddress:         sender.EmailAddress,
		Contacts:             contactsResp,
		IsArchived:           sender.IsArchived,
		InvoiceTemplateItems: templateItems,
		CreatedAt:            sender.CreatedAt.Time,
		UpdatedAt:            sender.UpdatedAt.Time,
	}

	return response, nil
}

func (s *senderService) UpdateSender(ctx context.Context, senderID uuid.UUID, req *UpdateSenderRequest) (*UpdateSenderResponse, error) {
	arg := db.UpdateSenderParams{
		ID:           senderID,
		Name:         req.Name,
		Address:      req.Address,
		PostalCode:   req.PostalCode,
		Place:        req.Place,
		Land:         req.Land,
		Kvknumber:    req.Kvknumber,
		Btwnumber:    req.Btwnumber,
		PhoneNumber:  req.PhoneNumber,
		ClientNumber: req.ClientNumber,
		EmailAddress: req.EmailAddress,
		IsArchived:   req.IsArchived,
		Types:        db.NullSenderTypesFromPtr(req.Types),
	}

	if req.Contacts != nil {
		contactsParam, err := json.Marshal(req.Contacts)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateSender", "Failed to marshal contacts", zap.Error(err))
			return nil, err
		}
		arg.Contacts = contactsParam
	}
	updatedSender, err := s.Store.UpdateSender(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateSender", "Failed to update sender", zap.Error(err))
		return nil, err
	}

	contactsResp := []SenderContact{}
	if err := json.Unmarshal(updatedSender.Contacts, &contactsResp); err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateSender", "Failed to unmarshal contacts", zap.Error(err))
		return nil, err
	}

	response := &UpdateSenderResponse{
		ID:           updatedSender.ID,
		Types:        string(updatedSender.Types),
		Name:         updatedSender.Name,
		Address:      updatedSender.Address,
		PostalCode:   updatedSender.PostalCode,
		Place:        updatedSender.Place,
		Land:         updatedSender.Land,
		Kvknumber:    updatedSender.Kvknumber,
		Btwnumber:    updatedSender.Btwnumber,
		PhoneNumber:  updatedSender.PhoneNumber,
		ClientNumber: updatedSender.ClientNumber,
		EmailAddress: updatedSender.EmailAddress,
		Contacts:     contactsResp,
		IsArchived:   updatedSender.IsArchived,
		CreatedAt:    updatedSender.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:    updatedSender.UpdatedAt.Time.Format(time.RFC3339),
	}

	return response, nil
}

func (s *senderService) DeleteSender(ctx context.Context, senderID uuid.UUID) error {
	err := s.Store.DeleteSender(ctx, senderID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteSender", "Failed to delete sender", zap.Error(err))
		return err
	}
	return nil
}

func (s *senderService) CreateSenderInvoiceTemplate(ctx context.Context, senderID uuid.UUID, req *CreateSenderInvoiceTemplateRequest) error {
	tmplids, err := s.Store.GetTemplateItemsByIds(ctx, req.InvoiceTemplateIDs)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateSenderInvoiceTemplate", "Failed to get template items by IDs", zap.Error(err))
		return err
	}
	_, err = s.Store.CreateSenderInvoiceTemplate(ctx, db.CreateSenderInvoiceTemplateParams{
		ID:              senderID,
		InvoiceTemplate: tmplids,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateSenderInvoiceTemplate", "Failed to create sender invoice template", zap.Error(err))
		return err
	}
	return nil
}
