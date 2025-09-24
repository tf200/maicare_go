package clientp

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"

	"go.uber.org/zap"
)

func (s *clientService) CreateAppointmentCard(req CreateAppointmentCardRequest, clientID int64, ctx context.Context) (*CreateAppointmentCardResponse, error) {
	appointmentCard, err := s.Store.CreateAppointmentCard(ctx, db.CreateAppointmentCardParams{
		ClientID:               clientID,
		GeneralInformation:     req.GeneralInformation,
		ImportantContacts:      req.ImportantContacts,
		HouseholdInfo:          req.HouseholdInfo,
		OrganizationAgreements: req.OrganizationAgreements,
		YouthOfficerAgreements: req.YouthOfficerAgreements,
		TreatmentAgreements:    req.TreatmentAgreements,
		SmokingRules:           req.SmokingRules,
		Work:                   req.Work,
		SchoolInternship:       req.SchoolInternship,
		Travel:                 req.Travel,
		Leave:                  req.Leave,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateAppointmentCard",
			"Failed to create appointment card", zap.Error(err))
		return nil, fmt.Errorf("failed to create appointment card")
	}
	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "CreateAppointmentCard",
		"Successfully created appointment card", zap.Int64("AppointmentCardID", appointmentCard.ID))
	return &CreateAppointmentCardResponse{
		ID:                     appointmentCard.ID,
		ClientID:               appointmentCard.ClientID,
		GeneralInformation:     appointmentCard.GeneralInformation,
		ImportantContacts:      appointmentCard.ImportantContacts,
		HouseholdInfo:          appointmentCard.HouseholdInfo,
		OrganizationAgreements: appointmentCard.OrganizationAgreements,
		YouthOfficerAgreements: appointmentCard.YouthOfficerAgreements,
		TreatmentAgreements:    appointmentCard.TreatmentAgreements,
		SmokingRules:           appointmentCard.SmokingRules,
		Work:                   appointmentCard.Work,
		SchoolInternship:       appointmentCard.SchoolInternship,
		Travel:                 appointmentCard.Travel,
		Leave:                  appointmentCard.Leave,
		CreatedAt:              appointmentCard.CreatedAt.Time,
		UpdatedAt:              appointmentCard.UpdatedAt.Time,
		FileUrl:                appointmentCard.FileUrl,
	}, nil
}

func (s *clientService) GetAppointmentCard(ctx context.Context, clientID int64) (*GetAppointmentCardResponse, error) {
	appointmentCard, err := s.Store.GetAppointmentCard(ctx, clientID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.Logger.LogBusinessEvent(logger.LogLevelWarn, "GetAppointmentCard",
				"Appointment card not found", zap.Int64("ClientID", clientID))
			return nil, fmt.Errorf("appointment card not found")
		}
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetAppointmentCard",
			"Failed to get appointment card", zap.Error(err))
		return nil, fmt.Errorf("failed to get appointment card")
	}
	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "GetAppointmentCard",
		"Successfully retrieved appointment card", zap.Int64("AppointmentCardID", appointmentCard.ID))
	return &GetAppointmentCardResponse{
		ID:                     appointmentCard.ID,
		ClientID:               appointmentCard.ClientID,
		GeneralInformation:     appointmentCard.GeneralInformation,
		ImportantContacts:      appointmentCard.ImportantContacts,
		HouseholdInfo:          appointmentCard.HouseholdInfo,
		OrganizationAgreements: appointmentCard.OrganizationAgreements,
		YouthOfficerAgreements: appointmentCard.YouthOfficerAgreements,
		TreatmentAgreements:    appointmentCard.TreatmentAgreements,
		SmokingRules:           appointmentCard.SmokingRules,
		Work:                   appointmentCard.Work,
		SchoolInternship:       appointmentCard.SchoolInternship,
		Travel:                 appointmentCard.Travel,
		Leave:                  appointmentCard.Leave,
		CreatedAt:              appointmentCard.CreatedAt.Time,
		UpdatedAt:              appointmentCard.UpdatedAt.Time,
		FileUrl:                appointmentCard.FileUrl,
	}, nil
}
