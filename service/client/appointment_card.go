package clientp

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pdf"

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

func (s *clientService) UpdateAppointmentCard(req UpdateAppointmentCardRequest, clientID int64, ctx context.Context) (*UpdateAppointmentCardResponse, error) {
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
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateAppointmentCard",
			"Failed to update appointment card", zap.Error(err))
		return nil, fmt.Errorf("failed to update appointment card")
	}
	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "UpdateAppointmentCard",
		"Successfully updated appointment card", zap.Int64("AppointmentCardID", appointmentCard.ID))
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

func (s *clientService) GenerateAppointmentCardDocumentApi(ctx context.Context, clientID int64) (*GenerateAppointmentCardDocumentApiResponse, error) {
	appointmentCard, err := s.Store.GetAppointmentCard(ctx, clientID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.Logger.LogBusinessEvent(logger.LogLevelWarn, "GenerateAppointmentCardDocumentApi",
				"Appointment card not found", zap.Int64("ClientID", clientID))
			return nil, fmt.Errorf("appointment card not found")
		}
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GenerateAppointmentCardDocumentApi",
			"Failed to retrieve appointment card", zap.Error(err))
		return nil, fmt.Errorf("failed to retrieve appointment card")
	}

	if appointmentCard.FileUrl != nil && *appointmentCard.FileUrl != "" {
		err = s.B2Client.Delete(ctx, *appointmentCard.FileUrl)
		if err != nil {
			s.Logger.LogBusinessEvent(logger.LogLevelError, "GenerateAppointmentCardDocumentApi",
				"Failed to delete existing appointment card document", zap.Error(err))
			return nil, fmt.Errorf("failed to delete existing appointment card document")
		}
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

	pdfUrl, err := pdf.GenerateAndUploadAppointmentCardPDF(ctx, pdfArg, s.B2Client)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GenerateAppointmentCardDocumentApi",
			"Failed to generate and upload appointment card PDF", zap.Error(err))
		return nil, fmt.Errorf("failed to generate and upload appointment card PDF")
	}

	fileKey, err := s.Store.UpdateAppointmentCardUrl(ctx, db.UpdateAppointmentCardUrlParams{
		ClientID: clientID,
		FileUrl:  &pdfUrl,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GenerateAppointmentCardDocumentApi",
			"Failed to update appointment card with file URL", zap.Error(err))
		return nil, fmt.Errorf("failed to update appointment card with file URL")
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "GenerateAppointmentCardDocumentApi",
		"Successfully generated appointment card document", zap.Int64("ClientID", clientID))
	return &GenerateAppointmentCardDocumentApiResponse{
		FileUrl:  s.GenerateResponsePresignedURL(fileKey, ctx),
		ClientID: clientID,
	}, nil
}
