package clientp

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/service/pdf"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *clientService) GetAppointmentCard(ctx context.Context, clientID uuid.UUID) (*GetAppointmentCardResponse, error) {
	appointmentCard, err := s.Store.GetAppointmentCard(ctx, clientID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "GetAppointmentCard",
				"Appointment card not found", zap.String("ClientID", clientID.String()))
			return nil, nil
		}
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetAppointmentCard",
			"Failed to get appointment card", zap.Error(err))

		return nil, fmt.Errorf("failed to get appointment card")

	}
	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "GetAppointmentCard",
		"Successfully retrieved appointment card", zap.String("AppointmentCardID", appointmentCard.ID.String()))
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
	}, nil
}

func (s *clientService) UpdateAppointmentCard(req UpdateAppointmentCardRequest, clientID uuid.UUID, ctx context.Context) (*UpdateAppointmentCardResponse, error) {
	arg := db.UpdateAppointmentCardParams{
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
	}
	appointmentCard, err := s.Store.UpdateAppointmentCard(ctx, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			createdAppointmentCard, createErr := s.Store.CreateAppointmentCard(ctx, db.CreateAppointmentCardParams{
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
			if createErr != nil {
				s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateAppointmentCard",
					"Failed to create appointment card during upsert", zap.Error(createErr))
				return nil, fmt.Errorf("failed to update appointment card")
			}
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "UpdateAppointmentCard",
				"Successfully created appointment card during upsert", zap.String("AppointmentCardID", createdAppointmentCard.ID.String()))
			return &UpdateAppointmentCardResponse{
				ID:                     createdAppointmentCard.ID,
				ClientID:               createdAppointmentCard.ClientID,
				GeneralInformation:     createdAppointmentCard.GeneralInformation,
				ImportantContacts:      createdAppointmentCard.ImportantContacts,
				HouseholdInfo:          createdAppointmentCard.HouseholdInfo,
				OrganizationAgreements: createdAppointmentCard.OrganizationAgreements,
				YouthOfficerAgreements: createdAppointmentCard.YouthOfficerAgreements,
				TreatmentAgreements:    createdAppointmentCard.TreatmentAgreements,
				SmokingRules:           createdAppointmentCard.SmokingRules,
				Work:                   createdAppointmentCard.Work,
				SchoolInternship:       createdAppointmentCard.SchoolInternship,
				Travel:                 createdAppointmentCard.Travel,
				Leave:                  createdAppointmentCard.Leave,
				CreatedAt:              createdAppointmentCard.CreatedAt.Time,
				UpdatedAt:              createdAppointmentCard.UpdatedAt.Time,
			}, nil
		}
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateAppointmentCard",
			"Failed to update appointment card", zap.Error(err))
		return nil, fmt.Errorf("failed to update appointment card")
	}
	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "UpdateAppointmentCard",
		"Successfully updated appointment card", zap.String("AppointmentCardID", appointmentCard.ID.String()))
	return &UpdateAppointmentCardResponse{
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
	}, nil
}

func (s *clientService) GenerateAppointmentCardDocumentApi(ctx context.Context, clientID uuid.UUID) ([]byte, string, error) {
	appointmentCard, err := s.Store.GetAppointmentCard(ctx, clientID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelWarn, "GenerateAppointmentCardDocumentApi",
				"Appointment card not found", zap.String("ClientID", clientID.String()))
			return nil, "", fmt.Errorf("appointment card not found")
		}
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GenerateAppointmentCardDocumentApi",
			"Failed to retrieve appointment card", zap.Error(err))
		return nil, "", fmt.Errorf("failed to retrieve appointment card")
	}

	pdfArg := pdf.AppointmentCard{
		ID:                     appointmentCard.ID,
		ClientName:             appointmentCard.FirstName + " " + appointmentCard.LastName,
		Date:                   appointmentCard.CreatedAt.Time.Format("02-01-2006"),
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
	}

	pdfBytes, err := s.PDFService.GenerateAppointmentCardPDF(ctx, pdfArg)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GenerateAppointmentCardDocumentApi",
			"Failed to generate appointment card PDF", zap.Error(err))
		return nil, "", fmt.Errorf("failed to generate appointment card PDF")
	}

	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "GenerateAppointmentCardDocumentApi",
		"Successfully generated appointment card document", zap.String("ClientID", clientID.String()))

	fileName := fmt.Sprintf("appointment_card_%s.pdf", appointmentCard.ID.String())
	return pdfBytes, fileName, nil
}
