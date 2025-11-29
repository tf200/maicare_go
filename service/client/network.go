package clientp

import (
	"context"
	"encoding/json"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"
	"maicare_go/service/notification"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *clientService) GetClientSender(ctx context.Context, clientID uuid.UUID) (*GetClientSenderResponse, error) {
	sender, err := s.Store.GetClientSender(ctx, clientID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetClientSender", "Failed to get client sender", zap.String("client_id", clientID.String()), zap.Error(err))
		return nil, err
	}

	var contacts []SenderContact
	if err := json.Unmarshal(sender.Contacts, &contacts); err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetClientSender", "Failed to unmarshal contacts", zap.String("client_id", clientID.String()), zap.Error(err))
		return nil, err
	}
	response := &GetClientSenderResponse{
		ID:           sender.ID,
		Types:        string(sender.Types),
		Name:         sender.Name,
		Address:      sender.Address,
		PostalCode:   sender.PostalCode,
		Place:        sender.Place,
		Land:         sender.Land,
		Kvknumber:    sender.Kvknumber,
		Btwnumber:    sender.Btwnumber,
		PhoneNumber:  sender.PhoneNumber,
		ClientNumber: sender.ClientNumber,
		IsArchived:   sender.IsArchived,
		Contacts:     contacts,
	}
	return response, nil
}

func (s *clientService) CreateClientEmergencyContact(ctx context.Context, req CreateClientEmergencyContactParams, clientID uuid.UUID) (*CreateClientEmergencyContactResponse, error) {
	arg := db.CreateEmemrgencyContactParams{
		ClientID:         clientID,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Email:            req.Email,
		PhoneNumber:      req.PhoneNumber,
		Address:          req.Address,
		Relationship:     req.Relationship,
		RelationStatus:   db.NullRelationStatusFromPtr(req.RelationStatus),
		MedicalReports:   req.MedicalReports,
		IncidentsReports: req.IncidentsReports,
		GoalsReports:     req.GoalsReports,
	}
	clientEmergencyContact, err := s.Store.CreateEmemrgencyContact(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientEmergencyContact", "Failed to create client emergency contact", zap.String("client_id", clientID.String()), zap.Error(err))
		return nil, err
	}

	response := &CreateClientEmergencyContactResponse{
		ID:               clientEmergencyContact.ID,
		ClientID:         clientEmergencyContact.ClientID,
		FirstName:        clientEmergencyContact.FirstName,
		LastName:         clientEmergencyContact.LastName,
		Email:            clientEmergencyContact.Email,
		PhoneNumber:      clientEmergencyContact.PhoneNumber,
		Address:          clientEmergencyContact.Address,
		Relationship:     clientEmergencyContact.Relationship,
		RelationStatus:   db.RelationStatusPtrFromEnum(clientEmergencyContact.RelationStatus),
		CreatedAt:        clientEmergencyContact.CreatedAt.Time,
		IsVerified:       clientEmergencyContact.IsVerified,
		MedicalReports:   clientEmergencyContact.MedicalReports,
		IncidentsReports: clientEmergencyContact.IncidentsReports,
		GoalsReports:     clientEmergencyContact.GoalsReports,
	}
	return response, nil
}

func (s *clientService) ListClientEmergencyContacts(ctx *gin.Context, req ListClientEmergencyContactsRequest, clientID uuid.UUID) (*pagination.Response[ListClientEmergencyContactsResponse], error) {
	params := req.GetParams()
	contacts, err := s.Store.ListEmergencyContacts(ctx, db.ListEmergencyContactsParams{
		ClientID: clientID,
		Limit:    params.Limit,
		Offset:   params.Offset,
		Search:   req.Search,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListClientEmergencyContacts", "Failed to list client emergency contacts", zap.String("client_id", clientID.String()), zap.Error(err))
		return nil, err
	}

	if len(contacts) == 0 {
		pag := pagination.NewResponse(ctx, req.Request, []ListClientEmergencyContactsResponse{}, 0)
		return &pag, nil
	}

	totalCount := contacts[0].TotalCount
	var results []ListClientEmergencyContactsResponse
	for _, contact := range contacts {
		results = append(results, ListClientEmergencyContactsResponse{
			ID:               contact.ID,
			ClientID:         contact.ClientID,
			FirstName:        contact.FirstName,
			LastName:         contact.LastName,
			Email:            contact.Email,
			PhoneNumber:      contact.PhoneNumber,
			Address:          contact.Address,
			Relationship:     contact.Relationship,
			RelationStatus:   db.RelationStatusPtrFromEnum(contact.RelationStatus),
			CreatedAt:        contact.CreatedAt.Time,
			IsVerified:       contact.IsVerified,
			MedicalReports:   contact.MedicalReports,
			IncidentsReports: contact.IncidentsReports,
			GoalsReports:     contact.GoalsReports,
		})
	}

	pag := pagination.NewResponse(ctx, req.Request, results, totalCount)
	return &pag, nil
}

func (s *clientService) GetClientEmergencyContact(ctx context.Context, contactID uuid.UUID) (*GetClientEmergencyContactResponse, error) {
	contact, err := s.Store.GetEmergencyContact(ctx, contactID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetClientEmergencyContact", "Failed to get client emergency contact", zap.String("contact_id", contactID.String()), zap.Error(err))
		return nil, err
	}

	response := &GetClientEmergencyContactResponse{
		ID:               contact.ID,
		ClientID:         contact.ClientID,
		FirstName:        contact.FirstName,
		LastName:         contact.LastName,
		Email:            contact.Email,
		PhoneNumber:      contact.PhoneNumber,
		Address:          contact.Address,
		Relationship:     contact.Relationship,
		RelationStatus:   db.RelationStatusPtrFromEnum(contact.RelationStatus),
		CreatedAt:        contact.CreatedAt.Time,
		IsVerified:       contact.IsVerified,
		MedicalReports:   contact.MedicalReports,
		IncidentsReports: contact.IncidentsReports,
		GoalsReports:     contact.GoalsReports,
	}
	return response, nil
}

func (s *clientService) UpdateClientEmergencyContact(ctx context.Context, req UpdateClientEmergencyContactParams, contactID uuid.UUID) (*UpdateClientEmergencyContactResponse, error) {
	arg := db.UpdateEmergencyContactParams{
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Email:            req.Email,
		PhoneNumber:      req.PhoneNumber,
		Address:          req.Address,
		Relationship:     req.Relationship,
		RelationStatus:   db.NullRelationStatusFromPtr(req.RelationStatus),
		MedicalReports:   req.MedicalReports,
		IncidentsReports: req.IncidentsReports,
		GoalsReports:     req.GoalsReports,
		ID:               contactID,
	}
	clientEmergencyContact, err := s.Store.UpdateEmergencyContact(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateClientEmergencyContact", "Failed to update client emergency contact", zap.String("contact_id", contactID.String()), zap.Error(err))
		return nil, err
	}

	response := &UpdateClientEmergencyContactResponse{
		ID:               clientEmergencyContact.ID,
		ClientID:         clientEmergencyContact.ClientID,
		FirstName:        clientEmergencyContact.FirstName,
		LastName:         clientEmergencyContact.LastName,
		Email:            clientEmergencyContact.Email,
		PhoneNumber:      clientEmergencyContact.PhoneNumber,
		Address:          clientEmergencyContact.Address,
		Relationship:     clientEmergencyContact.Relationship,
		RelationStatus:   db.RelationStatusPtrFromEnum(clientEmergencyContact.RelationStatus),
		CreatedAt:        clientEmergencyContact.CreatedAt.Time,
		IsVerified:       clientEmergencyContact.IsVerified,
		MedicalReports:   clientEmergencyContact.MedicalReports,
		IncidentsReports: clientEmergencyContact.IncidentsReports,
		GoalsReports:     clientEmergencyContact.GoalsReports,
	}
	return response, nil
}

func (s *clientService) DeleteClientEmergencyContact(ctx context.Context, contactID uuid.UUID) (*DeleteClientEmergencyContactResponse, error) {
	contact, err := s.Store.DeleteEmergencyContact(ctx, contactID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteClientEmergencyContact", "Failed to delete client emergency contact", zap.String("contact_id", contactID.String()), zap.Error(err))
		return nil, err
	}
	return &DeleteClientEmergencyContactResponse{ID: contact.ID}, nil
}

func (s *clientService) AssignEmployeeToClient(ctx context.Context, req AssignEmployeeRequest, clientID uuid.UUID) (*AssignEmployeeResponse, error) {
	arg := db.AssignEmployeeParams{
		ClientID:   clientID,
		EmployeeID: req.EmployeeID,
		StartDate:  pgtype.Date{Time: req.StartDate, Valid: true},
		Role:       req.Role,
	}
	assign, err := s.Store.AssignEmployee(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "AssignEmployeeToClient", "Failed to assign employee to client", zap.String("client_id", clientID.String()), zap.String("employee_id", req.EmployeeID.String()), zap.Error(err))
		return nil, err
	}

	notificationData := notification.NewClientAssignmentData{
		ClientID:        assign.ClientID,
		ClientFirstName: assign.ClientFirstName,
		ClientLastName:  assign.ClientLastName,
		ClientLocation:  assign.ClientLocationName,
	}

	err = s.asynqClient.EnqueueNotificationTask(ctx, notification.NotificationPayload{
		RecipientUserIDs: []uuid.UUID{assign.UserID},
		Type:             notification.TypeNewClientAssignment,
		Data: notification.NotificationData{
			NewClientAssignment: &notificationData,
		},
		CreatedAt: time.Now(),
		Message:   notificationData.NewClientAssignmentMessage(),
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "AssignEmployeeToClient", "Failed to enqueue notification task", zap.String("client_id", clientID.String()), zap.String("employee_id", req.EmployeeID.String()), zap.Error(err))
	}

	response := &AssignEmployeeResponse{
		ID:         assign.ID,
		ClientID:   assign.ClientID,
		EmployeeID: assign.EmployeeID,
		StartDate:  assign.StartDate.Time,
		Role:       assign.Role,
		CreatedAt:  assign.CreatedAt.Time,
	}
	return response, nil
}

func (s *clientService) ListAssignedEmployees(ctx *gin.Context, req ListAssignedEmployeesRequest, clientID uuid.UUID) (*pagination.Response[ListAssignedEmployeesResponse], error) {
	params := req.GetParams()
	assignedEmployees, err := s.Store.ListAssignedEmployees(ctx, db.ListAssignedEmployeesParams{
		ClientID: clientID,
		Limit:    params.Limit,
		Offset:   params.Offset,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListAssignedEmployees", "Failed to list assigned employees", zap.String("client_id", clientID.String()), zap.Error(err))
		return nil, err
	}

	if len(assignedEmployees) == 0 {
		pag := pagination.NewResponse(ctx, req.Request, []ListAssignedEmployeesResponse{}, 0)
		return &pag, nil
	}
	totalCount := assignedEmployees[0].TotalCount
	var results []ListAssignedEmployeesResponse
	for _, emp := range assignedEmployees {
		results = append(results, ListAssignedEmployeesResponse{
			ID:           emp.ID,
			ClientID:     emp.ClientID,
			EmployeeID:   emp.EmployeeID,
			StartDate:    emp.StartDate.Time,
			Role:         emp.Role,
			EmployeeName: emp.EmployeeFirstName + " " + emp.EmployeeLastName,
			CreatedAt:    emp.CreatedAt.Time,
		})
	}
	pag := pagination.NewResponse(ctx, req.Request, results, totalCount)
	return &pag, nil
}

func (s *clientService) GetAssignedEmployee(ctx context.Context, assignmentID uuid.UUID) (*GetAssignedEmployeeResponse, error) {
	assignment, err := s.Store.GetAssignedEmployee(ctx, assignmentID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetAssignedEmployee", "Failed to get assigned employee", zap.String("assignment_id", assignmentID.String()), zap.Error(err))
		return nil, err
	}

	response := &GetAssignedEmployeeResponse{
		ID:           assignment.ID,
		ClientID:     assignment.ClientID,
		EmployeeID:   assignment.EmployeeID,
		StartDate:    assignment.StartDate.Time,
		Role:         assignment.Role,
		EmployeeName: assignment.EmployeeFirstName + " " + assignment.EmployeeLastName,
		CreatedAt:    assignment.CreatedAt.Time,
	}
	return response, nil
}

func (s *clientService) UpdateAssignedEmployee(ctx context.Context, req UpdateAssignedEmployeeRequest, assignmentID uuid.UUID) (*UpdateAssignedEmployeeResponse, error) {
	arg := db.UpdateAssignedEmployeeParams{
		ID:         assignmentID,
		EmployeeID: req.EmployeeID,
		StartDate:  pgtype.Date{Time: req.StartDate, Valid: true},
		Role:       req.Role,
	}
	assign, err := s.Store.UpdateAssignedEmployee(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateAssignedEmployee", "Failed to update assigned employee", zap.String("assignment_id", assignmentID.String()), zap.Error(err))
		return nil, err
	}

	response := &UpdateAssignedEmployeeResponse{
		ID:         assign.ID,
		ClientID:   assign.ClientID,
		EmployeeID: assign.EmployeeID,
		StartDate:  assign.StartDate.Time,
		Role:       assign.Role,
		CreatedAt:  assign.CreatedAt.Time,
	}
	return response, nil
}

func (s *clientService) DeleteAssignedEmployee(ctx context.Context, assignmentID uuid.UUID) (*DeleteAssignedEmployeeResponse, error) {
	assignment, err := s.Store.DeleteAssignedEmployee(ctx, assignmentID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteAssignedEmployee", "Failed to delete assigned employee", zap.String("assignment_id", assignmentID.String()), zap.Error(err))
		return nil, err
	}
	return &DeleteAssignedEmployeeResponse{ID: assignment.ID}, nil
}

func (s *clientService) GetClientRelatedEmail(ctx context.Context, clientID uuid.UUID) (*GetClientRelatedEmailsResponse, error) {
	email, err := s.Store.GetClientRelatedEmails(ctx, clientID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetClientRelatedEmail", "Failed to get client related email", zap.String("client_id", clientID.String()), zap.Error(err))
		return nil, err
	}

	response := &GetClientRelatedEmailsResponse{
		Emails: email,
	}
	return response, nil
}
