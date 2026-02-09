package clientp

import (
	"context"
	"fmt"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *clientService) CreateClientDetails(req CreateClientDetailsRequest, ctx context.Context) (*CreateClientDetailsResponse, error) {
	parsedDateOfBirth, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientDetails",
			"Failed to parse date of birth", zap.Error(err))
		return nil, fmt.Errorf("failed to parse date of birth")
	}

	client, err := s.Store.CreateClientDetails(ctx, db.CreateClientDetailsParams{
		FirstName:                  req.FirstName,
		LastName:                   req.LastName,
		DateOfBirth:                pgtype.Date{Time: parsedDateOfBirth, Valid: true},
		Identity:                   true,
		Bsn:                        req.Bsn,
		BsnVerifiedBy:              req.BsnVerifiedBy,
		Email:                      req.Email,
		PhoneNumber:                req.PhoneNumber,
		CareType:                   db.NullIntakeCareTypeFromPtr(req.CareType),
		SenderID:                   req.SenderID,
		LocationID:                 req.LocationID,
		EducationCurrentlyEnrolled: req.EducationCurrentlyEnrolled,
		EducationInstitution:       req.EducationInstitution,
		EducationMentorName:        req.EducationMentorName,
		EducationMentorPhone:       req.EducationMentorPhone,
		EducationMentorEmail:       req.EducationMentorEmail,
		EducationAdditionalNotes:   req.EducationAdditionalNotes,
		WorkCurrentlyEmployed:      req.WorkCurrentlyEmployed,
		WorkCurrentEmployer:        req.WorkCurrentEmployer,
		WorkCurrentEmployerPhone:   req.WorkCurrentEmployerPhone,
		WorkCurrentEmployerEmail:   req.WorkCurrentEmployerEmail,
		WorkCurrentPosition:        req.WorkCurrentPosition,
		WorkStartDate:              pgtype.Date{Time: req.WorkStartDate, Valid: true},
		WorkAdditionalNotes:        req.WorkAdditionalNotes,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientDetails",
			"Failed to create client details", zap.Error(err))
		return nil, fmt.Errorf("failed to create client details")
	}

	result := &CreateClientDetailsResponse{
		ID:                         client.ID,
		FirstName:                  client.FirstName,
		LastName:                   client.LastName,
		DateOfBirth:                client.DateOfBirth.Time,
		Identity:                   client.Identity,
		Status:                     string(client.Status),
		Bsn:                        client.Bsn,
		Email:                      client.Email,
		PhoneNumber:                client.PhoneNumber,
		Gender:                     string(client.Gender),
		Filenumber:                 client.Filenumber,
		Created:                    client.CreatedAt.Time,
		SenderID:                   client.SenderID,
		LocationID:                 client.LocationID,
		EducationCurrentlyEnrolled: client.EducationCurrentlyEnrolled,
		EducationInstitution:       client.EducationInstitution,
		EducationMentorName:        client.EducationMentorName,
		EducationMentorPhone:       client.EducationMentorPhone,
		EducationMentorEmail:       client.EducationMentorEmail,
		EducationAdditionalNotes:   client.EducationAdditionalNotes,
		EducationLevel:             string(client.EducationLevel),
		WorkCurrentlyEmployed:      client.WorkCurrentlyEmployed,
		WorkCurrentEmployer:        client.WorkCurrentEmployer,
		WorkCurrentEmployerPhone:   client.WorkCurrentEmployerPhone,
		WorkCurrentEmployerEmail:   client.WorkCurrentEmployerEmail,
		WorkCurrentPosition:        client.WorkCurrentPosition,
		WorkStartDate:              client.WorkStartDate.Time,
		WorkAdditionalNotes:        client.WorkAdditionalNotes,
	}

	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "CreateClientDetails",
		"Successfully created client details", zap.String("ClientID", client.ID.String()))
	return result, nil
}

func (s *clientService) ListClientDetails(ctx *gin.Context, req ListClientsApiParams) (*pagination.Response[ListClientsApiResponse], error) {
	params := req.GetParams()

	var clients []db.ListClientDetailsRow
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		clients, err = q.ListClientDetails(ctx, db.ListClientDetailsParams{
			Limit:      params.Limit,
			Offset:     params.Offset,
			Status:     db.NullClientStatusFromPtr(req.Status),
			LocationID: req.LocationID,
			Search:     req.Search,
		})
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListClientDetails",
			"Failed to list client details", zap.Error(err))
		return nil, fmt.Errorf("failed to list client details")
	}
	if len(clients) == 0 {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "ListClientDetails",
			"No clients found")
		pagObj := pagination.NewResponse(ctx, req.Request, []ListClientsApiResponse{}, 0)
		return &pagObj, nil
	}

	totalCount := clients[0].TotalCount

	clientList := make([]ListClientsApiResponse, len(clients))
	for i, client := range clients {
		clientList[i] = ListClientsApiResponse{
			ID:           client.ID,
			FirstName:    client.FirstName,
			LastName:     client.LastName,
			DateOfBirth:  client.DateOfBirth.Time,
			Identity:     client.Identity,
			Status:       string(client.Status),
			Bsn:          client.Bsn,
			Email:        client.Email,
			PhoneNumber:  client.PhoneNumber,
			Gender:       string(client.Gender),
			Filenumber:   client.Filenumber,
			CreatedAt:    client.CreatedAt.Time,
			SenderID:     client.SenderID,
			LocationID:   client.LocationID,
			LocationName: client.LocationName,
		}
	}
	pagObj := pagination.NewResponse(ctx, req.Request, clientList, totalCount)
	return &pagObj, nil
}

func (s *clientService) ListWaitingListClients(ctx *gin.Context, req ListWaitingListClientsParams) (*pagination.Response[ListWaitingListClientsResponse], error) {
	params := req.GetParams()
	sortDays := "desc"
	if req.SortDays != nil {
		sortDays = *req.SortDays
	}

	var clients []db.ListWaitingListClientsRow
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		clients, err = q.ListWaitingListClients(ctx, db.ListWaitingListClientsParams{
			Search:    req.Search,
			Placement: db.NullIntakeCareTypeFromPtr(req.Placement),
			SortDays:  sortDays,
			Offset:    params.Offset,
			Limit:     params.Limit,
		})
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListWaitingListClients",
			"Failed to list waiting list clients", zap.Error(err))
		return nil, fmt.Errorf("failed to list waiting list clients")
	}

	if len(clients) == 0 {
		pagObj := pagination.NewResponse(ctx, req.Request, []ListWaitingListClientsResponse{}, 0)
		return &pagObj, nil
	}

	totalCount := clients[0].TotalCount
	response := make([]ListWaitingListClientsResponse, len(clients))
	for i, client := range clients {
		var admissionType *string
		if client.AdmissionType.Valid {
			value := string(client.AdmissionType.AdmissionTypeEnum)
			admissionType = &value
		}

		response[i] = ListWaitingListClientsResponse{
			ID:             client.ID,
			FirstName:      client.FirstName,
			LastName:       client.LastName,
			CareType:       db.IntakeCareTypePtrFromEnum(client.CareType),
			SenderName:     client.SenderName,
			DaysInWaitlist: client.DaysInWaitlist,
			AdmissionType:  admissionType,
			Bsn:            client.Bsn,
		}
	}

	pagObj := pagination.NewResponse(ctx, req.Request, response, totalCount)
	return &pagObj, nil
}

func (s *clientService) GetClientsCount(ctx context.Context) (*GetClientsCountResponse, error) {
	var count db.GetClientCountsRow
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		count, err = q.GetClientCounts(ctx)
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetClientsCount",
			"Failed to get clients count", zap.Error(err))
		return nil, fmt.Errorf("failed to get clients count")
	}
	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "GetClientsCount",
		"Successfully retrieved clients count")
	return &GetClientsCountResponse{
		TotalClients:         count.TotalClients,
		ClientsInCare:        count.ClientsInCare,
		ClientsOnWaitingList: count.ClientsOnWaitingList,
		ClientsOutOfCare:     count.ClientsOutOfCare,
	}, nil
}

func (s *clientService) GetClientDetails(ctx context.Context, clientID uuid.UUID) (*GetClientApiResponse, error) {
	var client db.GetClientDetailsRow
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		client, err = q.GetClientDetails(ctx, clientID)
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetClientDetails",
			"Failed to get client details", zap.Error(err), zap.String("ClientID", clientID.String()))
		return nil, fmt.Errorf("failed to get client details")
	}

	return &GetClientApiResponse{
		ID:                         client.ID,
		FirstName:                  client.FirstName,
		LastName:                   client.LastName,
		DateOfBirth:                client.DateOfBirth.Time,
		Identity:                   client.Identity,
		Status:                     string(client.Status),
		Bsn:                        client.Bsn,
		BsnVerifiedBy:              client.BsnVerifiedBy,
		BsnVerifiedByFirstName:     client.BsnVerifiedByFirstName,
		BsnVerifiedByLastName:      client.BsnVerifiedByLastName,
		Email:                      client.Email,
		PhoneNumber:                client.PhoneNumber,
		Gender:                     string(client.Gender),
		Filenumber:                 client.Filenumber,
		CreatedAt:                  client.CreatedAt.Time,
		SenderID:                   client.SenderID,
		LocationID:                 client.LocationID,
		EducationCurrentlyEnrolled: client.EducationCurrentlyEnrolled,
		EducationInstitution:       client.EducationInstitution,
		EducationMentorName:        client.EducationMentorName,
		EducationMentorEmail:       client.EducationMentorEmail,
		EducationMentorPhone:       client.EducationMentorPhone,
		EducationAdditionalNotes:   client.EducationAdditionalNotes,
		EducationLevel:             string(client.EducationLevel),
		WorkCurrentlyEmployed:      client.WorkCurrentlyEmployed,
		WorkCurrentEmployer:        client.WorkCurrentEmployer,
		WorkCurrentEmployerPhone:   client.WorkCurrentEmployerPhone,
		WorkCurrentEmployerEmail:   client.WorkCurrentEmployerEmail,
		WorkCurrentPosition:        client.WorkCurrentPosition,
		WorkStartDate:              client.WorkStartDate.Time,
		WorkAdditionalNotes:        client.WorkAdditionalNotes,
	}, nil
}

func (s *clientService) GetClientAddresses(ctx context.Context, clientID uuid.UUID) (*GetClientAddressesApiResponse, error) {

	return &GetClientAddressesApiResponse{}, nil
}

func (s *clientService) UpdateClientDetails(ctx context.Context, req UpdateClientDetailsRequest, clientID uuid.UUID) (*UpdateClientDetailsResponse, error) {

	result := &UpdateClientDetailsResponse{}
	return result, nil
}

func (s *clientService) UpdateClientStatus(ctx context.Context, req UpdateClientStatusRequest, clientID uuid.UUID) (*UpdateClientStatusResponse, error) {
	switch req.IsSchedueled {
	case true:
		return s.handleSchedueledStatusUpdates(ctx, req, clientID)
	case false:
		return s.handleNormalStatusUpdates(ctx, req, clientID)
	default:
		return nil, fmt.Errorf("invalid is_schedueled value")
	}
}

func (s *clientService) handleSchedueledStatusUpdates(ctx context.Context, req UpdateClientStatusRequest, clientID uuid.UUID) (*UpdateClientStatusResponse, error) {
	if req.SchedueledFor.Before(time.Now()) {
		return nil, fmt.Errorf("scheduled time must be in the future")
	}

	var schedueledChange db.ScheduledStatusChange
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		schedueledChange, err = q.CreateSchedueledClientStatusChange(ctx, db.CreateSchedueledClientStatusChangeParams{
			ClientID:      clientID,
			NewStatus:     db.NullClientStatusFromPtr(&req.Status),
			Reason:        &req.Reason,
			ScheduledDate: pgtype.Date{Time: req.SchedueledFor, Valid: true},
		})
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateClientStatus",
			"Failed to create scheduled status change", zap.Error(err), zap.String("ClientID", clientID.String()))
		return nil, fmt.Errorf("failed to create scheduled status change")
	}

	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "UpdateClientStatus",
		"Successfully created scheduled status change", zap.String("ClientID", clientID.String()),
		zap.String("NewStatus", req.Status), zap.Time("ScheduledFor", req.SchedueledFor))

	if !schedueledChange.NewStatus.Valid {
		return nil, fmt.Errorf("scheduled status change not created properly")
	}

	return &UpdateClientStatusResponse{
		ID:     clientID,
		Status: string(schedueledChange.NewStatus.ClientStatusEnum),
	}, nil
}

func (s *clientService) ListStatusHistory(ctx context.Context, clientID uuid.UUID) ([]ListStatusHistoryApiResponse, error) {
	arg := db.ListClientStatusHistoryParams{
		ClientID: clientID,
		Limit:    10,
		Offset:   0,
	}
	var histories []db.ClientStatusHistory
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		histories, err = q.ListClientStatusHistory(ctx, arg)
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListStatusHistory",
			"Failed to list status history", zap.Error(err), zap.String("ClientID", clientID.String()))
		return nil, fmt.Errorf("failed to list status history")
	}
	if len(histories) == 0 {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "ListStatusHistory",
			"No status history found", zap.String("ClientID", clientID.String()))
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

func (s *clientService) handleNormalStatusUpdates(ctx context.Context, req UpdateClientStatusRequest, clientID uuid.UUID) (*UpdateClientStatusResponse, error) {
	var result *UpdateClientStatusResponse
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		// 1. Fetch old details (automatically subject to RLS visibility)
		oldClient, err := q.GetClientDetails(ctx, clientID)
		if err != nil {
			return fmt.Errorf("failed to get client details: %w", err)
		}
		// 2. Perform the update (protected by RLS UPDATE policy)
		client, err := q.UpdateClientStatus(ctx, db.UpdateClientStatusParams{
			ID:     clientID,
			Status: db.ClientStatusEnum(req.Status),
		})
		if err != nil {
			return fmt.Errorf("failed to update client status: %w", err)
		}
		// 3. Create history record (protected by RLS INSERT policy)
		_, err = q.CreateClientStatusHistory(ctx, db.CreateClientStatusHistoryParams{
			ClientID:  clientID,
			OldStatus: util.StringPtr(string(oldClient.Status)),
			NewStatus: req.Status,
			Reason:    &req.Reason,
		})
		if err != nil {
			return fmt.Errorf("failed to create client status history: %w", err)
		}
		// Prepare response object
		result = &UpdateClientStatusResponse{
			ID:     client.ID,
			Status: string(client.Status),
		}
		return nil
	})
	if err != nil {
		// Log the error returned by the transaction
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateClientStatus",
			"Transaction failed", zap.Error(err), zap.String("ClientID", clientID.String()))
		return nil, fmt.Errorf("failed to update client status: %v", err)
	}
	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "UpdateClientStatus",
		"Successfully updated client status", zap.String("ClientID", clientID.String()),
		zap.String("NewStatus", req.Status))
	return result, nil
}

func (s *clientService) AddClientDocument(ctx context.Context, req AddClientDocumentApiRequest, clientID uuid.UUID) (*AddClientDocumentApiResponse, error) {
	arg := db.AddClientDocumentTxParams{
		ClientID:     clientID,
		AttachmentID: req.AttachmentID,
		Label:        req.Label,
	}

	clientDoc, err := s.Store.AddClientDocumentTx(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "AddClientDocument",
			"Failed to add client document", zap.Error(err), zap.String("ClientID", clientID.String()))
		return nil, fmt.Errorf("failed to add client document")
	}

	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "AddClientDocument",
		"Successfully added client document", zap.String("ClientID", clientID.String()),
		zap.String("DocumentID", clientDoc.ClientDocument.ID.String()))
	return &AddClientDocumentApiResponse{
		ID:           clientDoc.ClientDocument.ID,
		AttachmentID: clientDoc.ClientDocument.AttachmentUuid,
		ClientID:     clientDoc.ClientDocument.ClientID,
		Label:        string(clientDoc.ClientDocument.Label),
		Name:         clientDoc.Attachment.Name,
		File:         clientDoc.Attachment.File,
		Size:         clientDoc.Attachment.Size,
		IsUsed:       clientDoc.Attachment.IsUsed,
		Tag:          clientDoc.Attachment.Tag,
		UpdatedAt:    clientDoc.Attachment.Updated.Time,
		CreatedAt:    clientDoc.Attachment.Created.Time,
	}, nil
}

func (s *clientService) ListClientDocuments(ctx *gin.Context, req ListClientDocumentsApiRequest, clientID uuid.UUID) (*pagination.Response[ListClientDocumentsApiResponse], error) {
	params := req.GetParams()
	var clientDocs []db.ListClientDocumentsRow
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		clientDocs, err = q.ListClientDocuments(ctx, db.ListClientDocumentsParams{
			ClientID: clientID,
			Offset:   params.Offset,
			Limit:    params.Limit,
		})
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListClientDocuments",
			"Failed to list client documents", zap.Error(err), zap.String("ClientID", clientID.String()))
		return nil, fmt.Errorf("failed to list client documents")
	}
	if len(clientDocs) == 0 {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "ListClientDocuments",
			"No client documents found", zap.String("ClientID", clientID.String()))
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
			Label:          string(doc.Label),
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

func (s *clientService) DeleteClientDocument(ctx context.Context, clientID uuid.UUID, documentID uuid.UUID) (*DeleteClientDocumentApiResponse, error) {
	clientDoc, err := s.Store.DeleteClientDocumentTx(ctx, db.DeleteClientDocumentParams{
		AttachmentID: documentID,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteClientDocument",
			"Failed to delete client document", zap.Error(err), zap.String("ClientID", clientID.String()))
		return nil, fmt.Errorf("failed to delete client document")
	}
	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "DeleteClientDocument",
		"Successfully deleted client document", zap.String("ClientID", clientID.String()),
		zap.String("DocumentID", documentID.String()))
	return &DeleteClientDocumentApiResponse{
		ID:           clientDoc.ClientDocument.ID,
		AttachmentID: clientDoc.ClientDocument.AttachmentUuid,
	}, nil
}

func (s *clientService) GetMissingClientDocuments(ctx context.Context, clientID uuid.UUID) (*GetMissingClientDocumentsApiResponse, error) {
	var missingDocs []string
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		missingDocs, err = q.GetMissingClientDocuments(ctx, clientID)
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetMissingClientDocuments",
			"Failed to get missing client documents", zap.Error(err), zap.String("ClientID", clientID.String()))
		return nil, fmt.Errorf("failed to get missing client documents")
	}
	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "GetMissingClientDocuments",
		"Successfully retrieved missing client documents", zap.String("ClientID", clientID.String()))
	return &GetMissingClientDocumentsApiResponse{
		MissingDocs: missingDocs,
	}, nil
}
