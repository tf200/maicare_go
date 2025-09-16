package client

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"

	"go.uber.org/zap"
)

type CreateAppointmentCardRequest struct {
	ClientID               int64
	GeneralInformation     []string
	ImportantContacts      []string
	HouseholdInfo          []string
	OrganizationAgreements []string
	YouthOfficerAgreements []string
	TreatmentAgreements    []string
	SmokingRules           []string
	Work                   []string
	SchoolInternship       []string
	Travel                 []string
	Leave                  []string
}

func (s *clientService) CreateAppointmentCard(req CreateAppointmentCardRequest, ctx context.Context) (*db.AppointmentCard, error) {
	appointmentCard, err := s.Store.CreateAppointmentCard(ctx, db.CreateAppointmentCardParams{
		ClientID:               req.ClientID,
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
	return &appointmentCard, nil
}

func (s *clientService) GetAppointmentCard(ctx context.Context, clientID int64) (*db.GetAppointmentCardRow, error) {
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
	return &appointmentCard, nil
}
