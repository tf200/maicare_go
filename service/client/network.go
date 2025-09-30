package clientp

import (
	"context"
	"encoding/json"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"

	"go.uber.org/zap"
)

func (s *clientService) GetClientSender(ctx context.Context, clientID int64) (*GetClientSenderResponse, error) {
	sender, err := s.Store.GetClientSender(ctx, clientID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetClientSender", "Failed to get client sender", zap.Int64("client_id", clientID), zap.Error(err))
		return nil, err
	}

	var contacts []SenderContact
	if err := json.Unmarshal(sender.Contacts, &contacts); err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetClientSender", "Failed to unmarshal contacts", zap.Int64("client_id", clientID), zap.Error(err))
		return nil, err
	}
	response := &GetClientSenderResponse{
		ID:           sender.ID,
		Types:        sender.Types,
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

func (s *clientService) CreateClientEmergencyContact(ctx context.Context, req CreateClientEmergencyContactParams, clientID int64) (*CreateClientEmergencyContactResponse, error) {
	arg := db.CreateEmemrgencyContactParams{
		ClientID:         clientID,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Email:            req.Email,
		PhoneNumber:      req.PhoneNumber,
		Address:          req.Address,
		Relationship:     req.Relationship,
		RelationStatus:   req.RelationStatus,
		MedicalReports:   req.MedicalReports,
		IncidentsReports: req.IncidentsReports,
		GoalsReports:     req.GoalsReports,
	}
	clientEmergencyContact, err := s.Store.CreateEmemrgencyContact(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateClientEmergencyContact", "Failed to create client emergency contact", zap.Int64("client_id", clientID), zap.Error(err))
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
		RelationStatus:   clientEmergencyContact.RelationStatus,
		CreatedAt:        clientEmergencyContact.CreatedAt.Time,
		IsVerified:       clientEmergencyContact.IsVerified,
		MedicalReports:   clientEmergencyContact.MedicalReports,
		IncidentsReports: clientEmergencyContact.IncidentsReports,
		GoalsReports:     clientEmergencyContact.GoalsReports,
	}
	return response, nil
}
