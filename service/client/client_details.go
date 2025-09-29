package clientp

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *clientService) CreateClientDetails(req CreateClientDetailsRequest, ctx context.Context) (*CreateClientDetailsResponse, error) {
	parsedDateOfBirth, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateClientDetails",
			"Failed to parse date of birth", zap.Error(err))
		return nil, fmt.Errorf("failed to parse date of birth")
	}
	AddressesJSON, err := json.Marshal(req.Addresses)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateClientDetails",
			"Failed to marshal addresses", zap.Error(err))
		return nil, fmt.Errorf("failed to marshal addresses")
	}

	client, err := s.Store.CreateClientDetails(ctx, db.CreateClientDetailsParams{
		FirstName:                  req.FirstName,
		LastName:                   req.LastName,
		DateOfBirth:                pgtype.Date{Time: parsedDateOfBirth, Valid: true},
		Identity:                   true,
		Bsn:                        req.Bsn,
		BsnVerifiedBy:              req.BsnVerifiedBy,
		Source:                     req.Source,
		Birthplace:                 req.Birthplace,
		Email:                      req.Email,
		PhoneNumber:                req.PhoneNumber,
		Organisation:               req.Organisation,
		Departement:                req.Departement,
		Infix:                      req.Infix,
		Gender:                     req.Gender,
		Filenumber:                 req.Filenumber,
		SenderID:                   req.SenderID,
		LocationID:                 req.LocationID,
		Addresses:                  AddressesJSON,
		LegalMeasure:               req.LegalMeasure,
		EducationCurrentlyEnrolled: req.EducationCurrentlyEnrolled,
		EducationInstitution:       req.EducationInstitution,
		EducationMentorName:        req.EducationMentorName,
		EducationMentorPhone:       req.EducationMentorPhone,
		EducationMentorEmail:       req.EducationMentorEmail,
		EducationAdditionalNotes:   req.EducationAdditionalNotes,
		EducationLevel:             req.EducationLevel,
		WorkCurrentlyEmployed:      req.WorkCurrentlyEmployed,
		WorkCurrentEmployer:        req.WorkCurrentEmployer,
		WorkCurrentEmployerPhone:   req.WorkCurrentEmployerPhone,
		WorkCurrentEmployerEmail:   req.WorkCurrentEmployerEmail,
		WorkCurrentPosition:        req.WorkCurrentPosition,
		WorkStartDate:              pgtype.Date{Time: req.WorkStartDate, Valid: true},
		WorkAdditionalNotes:        req.WorkAdditionalNotes,
		LivingSituation:            req.LivingSituation,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateClientDetails",
			"Failed to create client details", zap.Error(err))
		return nil, fmt.Errorf("failed to create client details")
	}

	var addresses []Address
	err = json.Unmarshal(client.Addresses, &addresses)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateClientDetails",
			"Failed to unmarshal addresses", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal addresses")
	}

	result := &CreateClientDetailsResponse{
		ID:                         client.ID,
		FirstName:                  client.FirstName,
		LastName:                   client.LastName,
		DateOfBirth:                client.DateOfBirth.Time,
		Identity:                   client.Identity,
		Status:                     client.Status,
		Bsn:                        client.Bsn,
		Source:                     client.Source,
		Birthplace:                 client.Birthplace,
		Email:                      client.Email,
		PhoneNumber:                client.PhoneNumber,
		Organisation:               client.Organisation,
		Departement:                client.Departement,
		Gender:                     client.Gender,
		Filenumber:                 client.Filenumber,
		ProfilePicture:             client.ProfilePicture,
		Infix:                      client.Infix,
		Created:                    client.CreatedAt.Time,
		SenderID:                   client.SenderID,
		LocationID:                 client.LocationID,
		DepartureReason:            client.DepartureReason,
		DepartureReport:            client.DepartureReport,
		Addresses:                  addresses,
		LegalMeasure:               client.LegalMeasure,
		EducationCurrentlyEnrolled: client.EducationCurrentlyEnrolled,
		EducationInstitution:       client.EducationInstitution,
		EducationMentorName:        client.EducationMentorName,
		EducationMentorPhone:       client.EducationMentorPhone,
		EducationMentorEmail:       client.EducationMentorEmail,
		EducationAdditionalNotes:   client.EducationAdditionalNotes,
		EducationLevel:             client.EducationLevel,
		WorkCurrentlyEmployed:      client.WorkCurrentlyEmployed,
		WorkCurrentEmployer:        client.WorkCurrentEmployer,
		WorkCurrentEmployerPhone:   client.WorkCurrentEmployerPhone,
		WorkCurrentEmployerEmail:   client.WorkCurrentEmployerEmail,
		WorkCurrentPosition:        client.WorkCurrentPosition,
		WorkStartDate:              client.WorkStartDate.Time,
		WorkAdditionalNotes:        client.WorkAdditionalNotes,
		LivingSituation:            client.LivingSituation,
		LivingSituationNotes:       client.LivingSituationNotes,
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "CreateClientDetails",
		"Successfully created client details", zap.Int64("ClientID", client.ID))
	return result, nil
}

func (s *clientService) ListClientDetails(ctx *gin.Context, req ListClientsApiParams) (*pagination.Response[ListClientsApiResponse], error) {
	params := req.GetParams()

	clients, err := s.Store.ListClientDetails(ctx, db.ListClientDetailsParams{
		Limit:      params.Limit,
		Offset:     params.Offset,
		Status:     req.Status,
		LocationID: req.LocationID,
		Search:     req.Search,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListClientDetails",
			"Failed to list client details", zap.Error(err))
		return nil, fmt.Errorf("failed to list client details")
	}
	if len(clients) == 0 {
		s.Logger.LogBusinessEvent(logger.LogLevelInfo, "ListClientDetails",
			"No clients found")
		pagObj := pagination.NewResponse(ctx, req.Request, []ListClientsApiResponse{}, 0)
		return &pagObj, nil
	}

	totalCount := clients[0].TotalCount

	clientList := make([]ListClientsApiResponse, len(clients))
	for i, client := range clients {
		var addresses []Address
		err = json.Unmarshal(client.Addresses, &addresses)
		if err != nil {
			s.Logger.LogBusinessEvent(logger.LogLevelError, "ListClientDetails",
				"Failed to unmarshal addresses", zap.Error(err))
			return nil, fmt.Errorf("failed to unmarshal addresses")
		}
		clientList[i] = ListClientsApiResponse{
			ID:                    client.ID,
			FirstName:             client.FirstName,
			LastName:              client.LastName,
			DateOfBirth:           client.DateOfBirth.Time,
			Identity:              client.Identity,
			Status:                client.Status,
			Bsn:                   client.Bsn,
			Source:                client.Source,
			Birthplace:            client.Birthplace,
			Email:                 client.Email,
			PhoneNumber:           client.PhoneNumber,
			Organisation:          client.Organisation,
			Departement:           client.Departement,
			Gender:                client.Gender,
			Filenumber:            client.Filenumber,
			ProfilePicture:        s.GenerateResponsePresignedURL(client.ProfilePicture, ctx),
			Infix:                 client.Infix,
			CreatedAt:             client.CreatedAt.Time,
			SenderID:              client.SenderID,
			LocationID:            client.LocationID,
			DepartureReason:       client.DepartureReason,
			DepartureReport:       client.DepartureReport,
			Addresses:             addresses,
			LegalMeasure:          client.LegalMeasure,
			HasUntakenMedications: client.HasUntakenMedications,
		}
	}
	pagObj := pagination.NewResponse(ctx, req.Request, clientList, totalCount)
	return &pagObj, nil
}

func (s *clientService) GetClientsCount(ctx context.Context) (*GetClientsCountResponse, error) {
	count, err := s.Store.GetClientCounts(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetClientsCount",
			"Failed to get clients count", zap.Error(err))
		return nil, fmt.Errorf("failed to get clients count")
	}
	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "GetClientsCount",
		"Successfully retrieved clients count")
	return &GetClientsCountResponse{
		TotalClients:         count.TotalClients,
		ClientsInCare:        count.ClientsInCare,
		ClientsOnWaitingList: count.ClientsOnWaitingList,
		ClientsOutOfCare:     count.ClientsOutOfCare,
	}, nil
}

func (s *clientService) GetClientDetails(ctx context.Context, clientID int64) (*GetClientApiResponse, error) {
	client, err := s.Store.GetClientDetails(ctx, clientID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetClientDetails",
			"Failed to get client details", zap.Error(err), zap.Int64("ClientID", clientID))
		return nil, fmt.Errorf("failed to get client details")
	}

	return &GetClientApiResponse{
		ID:                         client.ID,
		FirstName:                  client.FirstName,
		LastName:                   client.LastName,
		DateOfBirth:                client.DateOfBirth.Time,
		Identity:                   client.Identity,
		Status:                     client.Status,
		Bsn:                        client.Bsn,
		BsnVerifiedBy:              client.BsnVerifiedBy,
		BsnVerifiedByFirstName:     client.BsnVerifiedByFirstName,
		BsnVerifiedByLastName:      client.BsnVerifiedByLastName,
		Source:                     client.Source,
		Birthplace:                 client.Birthplace,
		Email:                      client.Email,
		PhoneNumber:                client.PhoneNumber,
		Organisation:               client.Organisation,
		Departement:                client.Departement,
		Gender:                     client.Gender,
		Filenumber:                 client.Filenumber,
		ProfilePicture:             s.GenerateResponsePresignedURL(client.ProfilePicture, ctx),
		Infix:                      client.Infix,
		CreatedAt:                  client.CreatedAt.Time,
		SenderID:                   client.SenderID,
		LocationID:                 client.LocationID,
		DepartureReason:            client.DepartureReason,
		DepartureReport:            client.DepartureReport,
		LegalMeasure:               client.LegalMeasure,
		HasUntakenMedications:      client.HasUntakenMedications,
		EducationCurrentlyEnrolled: client.EducationCurrentlyEnrolled,
		EducationInstitution:       client.EducationInstitution,
		EducationMentorName:        client.EducationMentorName,
		EducationMentorEmail:       client.EducationMentorEmail,
		EducationMentorPhone:       client.EducationMentorPhone,
		EducationAdditionalNotes:   client.EducationAdditionalNotes,
		EducationLevel:             client.EducationLevel,
		WorkCurrentlyEmployed:      client.WorkCurrentlyEmployed,
		WorkCurrentEmployer:        client.WorkCurrentEmployer,
		WorkCurrentEmployerPhone:   client.WorkCurrentEmployerPhone,
		WorkCurrentEmployerEmail:   client.WorkCurrentEmployerEmail,
		WorkCurrentPosition:        client.WorkCurrentPosition,
		WorkStartDate:              client.WorkStartDate.Time,
		WorkAdditionalNotes:        client.WorkAdditionalNotes,
		LivingSituation:            client.LivingSituation,
		LivingSituationNotes:       client.LivingSituationNotes,
	}, nil

}

func (s *clientService) GetClientAddresses(ctx context.Context, clientID int64) (*GetClientAddressesApiResponse, error) {
	address, err := s.Store.GetClientAddresses(ctx, clientID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetClientAddresses",
			"Failed to get client addresses", zap.Error(err), zap.Int64("ClientID", clientID))
		return nil, fmt.Errorf("failed to get client addresses")
	}

	var addresses []Address
	err = json.Unmarshal(address, &addresses)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetClientAddresses",
			"Failed to unmarshal addresses", zap.Error(err), zap.Int64("ClientID", clientID))
		return nil, fmt.Errorf("failed to unmarshal addresses")
	}

	return &GetClientAddressesApiResponse{
		Addresses: addresses,
	}, nil
}

func (s *clientService) UpdateClientDetails(ctx context.Context, req UpdateClientDetailsRequest, clientID int64) (*UpdateClientDetailsResponse, error) {
	client, err := s.Store.UpdateClientDetails(ctx, db.UpdateClientDetailsParams{
		ID:                         clientID,
		FirstName:                  req.FirstName,
		LastName:                   req.LastName,
		DateOfBirth:                pgtype.Date{Time: req.DateOfBirth, Valid: true},
		Identity:                   req.Identity,
		Bsn:                        req.Bsn,
		BsnVerifiedBy:              req.BsnVerifiedBy,
		Source:                     req.Source,
		Birthplace:                 req.Birthplace,
		Email:                      req.Email,
		PhoneNumber:                req.PhoneNumber,
		Organisation:               req.Organisation,
		Departement:                req.Departement,
		Gender:                     req.Gender,
		Filenumber:                 req.Filenumber,
		ProfilePicture:             req.ProfilePicture,
		Infix:                      req.Infix,
		SenderID:                   req.SenderID,
		LocationID:                 req.LocationID,
		DepartureReason:            req.DepartureReason,
		DepartureReport:            req.DepartureReport,
		LegalMeasure:               req.LegalMeasure,
		EducationCurrentlyEnrolled: req.EducationCurrentlyEnrolled,
		EducationInstitution:       req.EducationInstitution,
		EducationMentorName:        req.EducationMentorName,
		EducationMentorPhone:       req.EducationMentorPhone,
		EducationMentorEmail:       req.EducationMentorEmail,
		EducationAdditionalNotes:   req.EducationAdditionalNotes,
		EducationLevel:             req.EducationLevel,
		WorkCurrentlyEmployed:      req.WorkCurrentlyEmployed,
		WorkCurrentEmployer:        req.WorkCurrentEmployer,
		WorkCurrentEmployerPhone:   req.WorkCurrentEmployerPhone,
		WorkCurrentEmployerEmail:   req.WorkCurrentEmployerEmail,
		WorkCurrentPosition:        req.WorkCurrentPosition,
		WorkStartDate:              pgtype.Date{Time: req.WorkStartDate, Valid: true},
		WorkAdditionalNotes:        req.WorkAdditionalNotes,
		LivingSituation:            req.LivingSituation,
		LivingSituationNotes:       req.LivingSituationNotes,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateClientDetails",
			"Failed to update client details", zap.Error(err), zap.Int64("ClientID", clientID))
		return nil, fmt.Errorf("failed to update client details")
	}

	var addresses []Address
	err = json.Unmarshal(client.Addresses, &addresses)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateClientDetails",
			"Failed to unmarshal addresses", zap.Error(err), zap.Int64("ClientID", clientID))
		return nil, fmt.Errorf("failed to unmarshal addresses")
	}

	result := &UpdateClientDetailsResponse{
		ID:                    client.ID,
		FirstName:             client.FirstName,
		LastName:              client.LastName,
		DateOfBirth:           client.DateOfBirth.Time,
		Identity:              client.Identity,
		Status:                client.Status,
		Bsn:                   client.Bsn,
		BsnVerifiedBy:         client.BsnVerifiedBy,
		Source:                client.Source,
		Birthplace:            client.Birthplace,
		Email:                 client.Email,
		PhoneNumber:           client.PhoneNumber,
		Organisation:          client.Organisation,
		Departement:           client.Departement,
		Gender:                client.Gender,
		Filenumber:            client.Filenumber,
		ProfilePicture:        client.ProfilePicture,
		Infix:                 client.Infix,
		SenderID:              client.SenderID,
		LocationID:            client.LocationID,
		DepartureReason:       client.DepartureReason,
		DepartureReport:       client.DepartureReport,
		Addresses:             addresses,
		LegalMeasure:          client.LegalMeasure,
		HasUntakenMedications: client.HasUntakenMedications,
	}
	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "UpdateClientDetails",
		"Successfully updated client details", zap.Int64("ClientID", client.ID))
	return result, nil
}

func (s *clientService) UpdateClientStatus(ctx context.Context, req UpdateClientStatusRequest, clientID int64) (*UpdateClientStatusResponse, error) {
	switch req.IsSchedueled {
	case true:
		return s.handleSchedueledStatusUpdates(ctx, req, clientID)
	case false:
		return s.handleNormalStatusUpdates(ctx, req, clientID)
	default:
		return nil, fmt.Errorf("invalid is_schedueled value")
	}
}

func (s *clientService) handleSchedueledStatusUpdates(ctx context.Context, req UpdateClientStatusRequest, clientID int64) (*UpdateClientStatusResponse, error) {
	if req.SchedueledFor.Before(time.Now()) {
		return nil, fmt.Errorf("scheduled time must be in the future")
	}

	schedueledChange, err := s.Store.CreateSchedueledClientStatusChange(ctx, db.CreateSchedueledClientStatusChangeParams{
		ClientID:      clientID,
		NewStatus:     &req.Status,
		Reason:        &req.Reason,
		ScheduledDate: pgtype.Date{Time: req.SchedueledFor, Valid: true},
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateClientStatus",
			"Failed to create scheduled status change", zap.Error(err), zap.Int64("ClientID", clientID))
		return nil, fmt.Errorf("failed to create scheduled status change")
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "UpdateClientStatus",
		"Successfully created scheduled status change", zap.Int64("ClientID", clientID),
		zap.String("NewStatus", req.Status), zap.Time("ScheduledFor", req.SchedueledFor))

	return &UpdateClientStatusResponse{
		ID:     clientID,
		Status: schedueledChange.NewStatus,
	}, nil

}

func (s *clientService) ListStatusHistory(ctx context.Context, clientID int64) ([]ListStatusHistoryApiResponse, error) {
	arg := db.ListClientStatusHistoryParams{
		ClientID: clientID,
		Limit:    10,
		Offset:   0,
	}
	histories, err := s.Store.ListClientStatusHistory(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListStatusHistory",
			"Failed to list status history", zap.Error(err), zap.Int64("ClientID", clientID))
		return nil, fmt.Errorf("failed to list status history")
	}
	if len(histories) == 0 {
		s.Logger.LogBusinessEvent(logger.LogLevelInfo, "ListStatusHistory",
			"No status history found", zap.Int64("ClientID", clientID))
		return []ListStatusHistoryApiResponse{}, nil
	}
	var historyList []ListStatusHistoryApiResponse
	for _, history := range histories {
		historyList = append(historyList, ListStatusHistoryApiResponse{
			ID:        history.ID,
			ClientID:  history.ClientID,
			OldStatus: history.OldStatus,
			NewStatus: history.NewStatus,
			Reason:    history.Reason,
			ChangedAt: history.ChangedAt.Time,
			ChangedBy: history.ChangedBy,
		})
	}
	return historyList, nil
}

func (s *clientService) handleNormalStatusUpdates(ctx context.Context, req UpdateClientStatusRequest, clientID int64) (*UpdateClientStatusResponse, error) {
	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateClientStatus",
			"Failed to begin transaction", zap.Error(err), zap.Int64("ClientID", clientID))
		return nil, fmt.Errorf("failed to begin transaction")
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
			s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateClientStatus",
				"Failed to rollback transaction", zap.Error(rollbackErr), zap.Int64("ClientID", clientID))
		}
	}()

	qtx := s.Store.WithTx(tx)

	oldClient, err := qtx.GetClientDetails(ctx, clientID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateClientStatus",
			"Failed to get client details", zap.Error(err), zap.Int64("ClientID", clientID))
		return nil, fmt.Errorf("failed to get client details")
	}

	client, err := qtx.UpdateClientStatus(ctx, db.UpdateClientStatusParams{
		ID:     clientID,
		Status: &req.Status,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateClientStatus",
			"Failed to update client status", zap.Error(err), zap.Int64("ClientID", clientID))
		return nil, fmt.Errorf("failed to update client status")
	}

	_, err = qtx.CreateClientStatusHistory(ctx, db.CreateClientStatusHistoryParams{
		ClientID:  clientID,
		OldStatus: oldClient.Status,
		NewStatus: req.Status,
		Reason:    &req.Reason,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateClientStatus",
			"Failed to create client status history", zap.Error(err), zap.Int64("ClientID", clientID))
		return nil, fmt.Errorf("failed to create client status history")
	}

	if err = tx.Commit(ctx); err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateClientStatus",
			"Failed to commit transaction", zap.Error(err), zap.Int64("ClientID", clientID))
		return nil, fmt.Errorf("failed to commit transaction")
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "UpdateClientStatus",
		"Successfully updated client status", zap.Int64("ClientID", client.ID),
		zap.String("NewStatus", req.Status))

	return &UpdateClientStatusResponse{
		ID:     client.ID,
		Status: client.Status,
	}, nil
}

func (s *clientService) SetClientProfilePicture(ctx context.Context, req SetClientProfilePictureRequest, clientID int64) (*SetClientProfilePictureResponse, error) {
	arg := db.SetClientProfilePictureTxParams{
		ClientID:     clientID,
		AttachmentID: req.AttachmentID,
	}
	client, err := s.Store.SetClientProfilePictureTx(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "SetClientProfilePicture",
			"Failed to set client profile picture", zap.Error(err), zap.Int64("ClientID", clientID))
		return nil, fmt.Errorf("failed to set client profile picture")
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "SetClientProfilePicture",
		"Successfully set client profile picture", zap.Int64("ClientID", client.User.ID))

	return &SetClientProfilePictureResponse{
		ID:             client.User.ID,
		ProfilePicture: client.User.ProfilePicture,
	}, nil
}

func (s *clientService) AddClientDocument(ctx context.Context, req AddClientDocumentApiRequest, clientID int64) (*AddClientDocumentApiResponse, error) {
	arg := db.AddClientDocumentTxParams{
		ClientID:     clientID,
		AttachmentID: req.AttachmentID,
		Label:        req.Label,
	}

	clientDoc, err := s.Store.AddClientDocumentTx(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "AddClientDocument",
			"Failed to add client document", zap.Error(err), zap.Int64("ClientID", clientID))
		return nil, fmt.Errorf("failed to add client document")
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "AddClientDocument",
		"Successfully added client document", zap.Int64("ClientID", clientID),
		zap.Int64("DocumentID", clientDoc.ClientDocument.ID))
	return &AddClientDocumentApiResponse{
		ID:           clientDoc.ClientDocument.ID,
		AttachmentID: clientDoc.ClientDocument.AttachmentUuid,
		ClientID:     clientDoc.ClientDocument.ClientID,
		Label:        clientDoc.ClientDocument.Label,
		Name:         clientDoc.Attachment.Name,
		File:         clientDoc.Attachment.File,
		Size:         clientDoc.Attachment.Size,
		IsUsed:       clientDoc.Attachment.IsUsed,
		Tag:          clientDoc.Attachment.Tag,
		UpdatedAt:    clientDoc.Attachment.Updated.Time,
		CreatedAt:    clientDoc.Attachment.Created.Time,
	}, nil
}

func (s *clientService) ListClientDocuments(ctx *gin.Context, req ListClientDocumentsApiRequest, clientID int64) (*pagination.Response[ListClientDocumentsApiResponse], error) {
	params := req.GetParams()
	clientDocs, err := s.Store.ListClientDocuments(ctx, db.ListClientDocumentsParams{
		ClientID: clientID,
		Offset:   params.Offset,
		Limit:    params.Limit,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListClientDocuments",
			"Failed to list client documents", zap.Error(err), zap.Int64("ClientID", clientID))
		return nil, fmt.Errorf("failed to list client documents")
	}
	if len(clientDocs) == 0 {
		s.Logger.LogBusinessEvent(logger.LogLevelInfo, "ListClientDocuments",
			"No client documents found", zap.Int64("ClientID", clientID))
		pag := pagination.NewResponse(ctx, req.Request, []ListClientDocumentsApiResponse{}, 0)
		return &pag, nil
	}

	totalCount := clientDocs[0].TotalCount

	var docList []ListClientDocumentsApiResponse
	for _, doc := range clientDocs {
		docList = append(docList, ListClientDocumentsApiResponse{
			ID:             doc.ID,
			AttachmentUuid: doc.AttachmentUuid,
			ClientID:       doc.ClientID,
			Label:          doc.Label,
			Name:           doc.Name,
			File:           s.GenerateResponsePresignedURL(&doc.File, ctx),
			Size:           doc.Size,
			IsUsed:         doc.IsUsed,
			Tag:            doc.Tag,
			UpdatedAt:      doc.Updated.Time,
			CreatedAt:      doc.Created.Time,
		})
	}
	pag := pagination.NewResponse(ctx, req.Request, docList, totalCount)
	return &pag, nil
}

func (s *clientService) DeleteClientDocument(ctx context.Context, clientID int64, documentID uuid.UUID) (*DeleteClientDocumentApiResponse, error) {
	clientDoc, err := s.Store.DeleteClientDocumentTx(ctx, db.DeleteClientDocumentParams{
		AttachmentID: documentID,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "DeleteClientDocument",
			"Failed to delete client document", zap.Error(err), zap.Int64("ClientID", clientID))
		return nil, fmt.Errorf("failed to delete client document")
	}
	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "DeleteClientDocument",
		"Successfully deleted client document", zap.Int64("ClientID", clientID),
		zap.String("DocumentID", documentID.String()))
	return &DeleteClientDocumentApiResponse{
		ID:           clientDoc.ClientDocument.ID,
		AttachmentID: clientDoc.ClientDocument.AttachmentUuid,
	}, nil
}

func (s *clientService) GetMissingClientDocuments(ctx context.Context, clientID int64) (*GetMissingClientDocumentsApiResponse, error) {
	missingDocs, err := s.Store.GetMissingClientDocuments(ctx, clientID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetMissingClientDocuments",
			"Failed to get missing client documents", zap.Error(err), zap.Int64("ClientID", clientID))
		return nil, fmt.Errorf("failed to get missing client documents")
	}
	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "GetMissingClientDocuments",
		"Successfully retrieved missing client documents", zap.Int64("ClientID", clientID))
	return &GetMissingClientDocumentsApiResponse{
		MissingDocs: missingDocs,
	}, nil
}
