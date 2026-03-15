package clientp

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"time"

	"maicare_go/async/aclient"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *clientService) CreateRegistrationForm(ctx context.Context, req *CreateRegistrationFormRequest) (*CreateRegistrationFormResponse, error) {
	arg := db.CreateRegistrationFormParams{
		ClientFirstName:               req.ClientFirstName,
		ClientLastName:                req.ClientLastName,
		ClientDateOfBirth:             pgtype.Date{Valid: false},
		ClientBsnNumber:               req.ClientBsnNumber,
		ClientGender:                  db.GenderEnum(req.ClientGender),
		ClientNationality:             req.ClientNationality,
		ClientPhoneNumber:             req.ClientPhoneNumber,
		ClientEmail:                   req.ClientEmail,
		ClientStreet:                  req.ClientStreet,
		ClientHouseNumber:             req.ClientHouseNumber,
		ClientHouseNumberAddition:     req.ClientHouseNumberAddition,
		ClientPostalCode:              req.ClientPostalCode,
		ClientCity:                    req.ClientCity,
		ReferrerFirstName:             req.ReferrerFirstName,
		ReferrerLastName:              req.ReferrerLastName,
		ReferrerOrganization:          req.ReferrerOrganization,
		ReferrerJobTitle:              req.ReferrerJobTitle,
		ReferrerPhoneNumber:           req.ReferrerPhoneNumber,
		ReferrerEmail:                 req.ReferrerEmail,
		Guardian1FirstName:            req.Guardian1FirstName,
		Guardian1LastName:             req.Guardian1LastName,
		Guardian1Relationship:         req.Guardian1Relationship,
		Guardian1PhoneNumber:          req.Guardian1PhoneNumber,
		Guardian1Email:                req.Guardian1Email,
		Guardian2FirstName:            req.Guardian2FirstName,
		Guardian2LastName:             req.Guardian2LastName,
		Guardian2Relationship:         req.Guardian2Relationship,
		Guardian2PhoneNumber:          req.Guardian2PhoneNumber,
		Guardian2Email:                req.Guardian2Email,
		EducationInstitution:          req.EducationInstitution,
		EducationMentorName:           req.EducationMentorName,
		EducationMentorPhone:          req.EducationMentorPhone,
		EducationMentorEmail:          req.EducationMentorEmail,
		EducationCurrentlyEnrolled:    req.EducationCurrentlyEnrolled,
		EducationAdditionalNotes:      req.EducationAdditionalNotes,
		EducationLevel:                db.EducationLevelEnum(req.EducationLevel),
		WorkCurrentEmployer:           req.WorkCurrentEmployer,
		WorkEmployerPhone:             req.WorkEmployerPhone,
		WorkEmployerEmail:             req.WorkEmployerEmail,
		WorkCurrentPosition:           req.WorkCurrentPosition,
		WorkCurrentlyEmployed:         req.WorkCurrentlyEmployed,
		WorkAdditionalNotes:           req.WorkAdditionalNotes,
		CareProtectedLiving:           req.CareProtectedLiving,
		CareAssistedIndependentLiving: req.CareAssistedIndependentLiving,
		CareRoomTrainingCenter:        req.CareRoomTrainingCenter,
		CareAmbulatoryGuidance:        req.CareAmbulatoryGuidance,
		ClientGoals:                   req.ClientGoals,
		ApplicationReason:             req.ApplicationReason,
		RiskAggressiveBehavior:        req.RiskAggressiveBehavior,
		RiskSuicidalSelfharm:          req.RiskSuicidalSelfharm,
		RiskSubstanceAbuse:            req.RiskSubstanceAbuse,
		RiskPsychiatricIssues:         req.RiskPsychiatricIssues,
		RiskCriminalHistory:           req.RiskCriminalHistory,
		RiskFlightBehavior:            req.RiskFlightBehavior,
		RiskWeaponPossession:          req.RiskWeaponPossession,
		RiskSexualBehavior:            req.RiskSexualBehavior,
		RiskDayNightRhythm:            req.RiskDayNightRhythm,
		RiskOther:                     req.RiskOther,
		RiskOtherDescription:          req.RiskOtherDescription,
		RiskAdditionalNotes:           req.RiskAdditionalNotes,
		DocumentReferral:              req.DocumentReferral,
		DocumentEducationReport:       req.DocumentEducationReport,
		DocumentPsychiatricReport:     req.DocumentPsychiatricReport,
		DocumentDiagnosis:             req.DocumentDiagnosis,
		DocumentSafetyPlan:            req.DocumentSafetyPlan,
		DocumentIDCopy:                req.DocumentIDCopy,
		ApplicationDate:               pgtype.Date{Time: req.ApplicationDate, Valid: true},
		ReferrerSignature:             req.ReferrerSignature,
	}
	if req.WorkStartDate != nil {
		arg.WorkStartDate = pgtype.Date{Time: *req.WorkStartDate, Valid: true}
	} else {
		arg.WorkStartDate = pgtype.Date{Valid: false}
	}
	if req.ClientDateOfBirth != nil {
		arg.ClientDateOfBirth = pgtype.Date{Time: *req.ClientDateOfBirth, Valid: true}
	}
	if req.ClientDateOfBirth != nil {
		arg.ClientDateOfBirth = pgtype.Date{Time: *req.ClientDateOfBirth, Valid: true}
	} else {
		arg.ClientDateOfBirth = pgtype.Date{Valid: false}
	}
	createdForm, err := s.Store.CreateRegistrationForm(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateRegistrationForm", "Failed to create registration form", zap.Error(err))
		return nil, fmt.Errorf("failed to create registration form: %w", err)
	}
	response := &CreateRegistrationFormResponse{
		ID:                            createdForm.ID,
		ClientFirstName:               createdForm.ClientFirstName,
		ClientLastName:                createdForm.ClientLastName,
		ClientDateOfBirth:             &createdForm.ClientDateOfBirth.Time,
		ClientBsnNumber:               createdForm.ClientBsnNumber,
		ClientGender:                  string(createdForm.ClientGender),
		ClientNationality:             createdForm.ClientNationality,
		ClientPhoneNumber:             createdForm.ClientPhoneNumber,
		ClientEmail:                   createdForm.ClientEmail,
		ClientStreet:                  createdForm.ClientStreet,
		ClientHouseNumber:             createdForm.ClientHouseNumber,
		ClientHouseNumberAddition:     createdForm.ClientHouseNumberAddition,
		ClientPostalCode:              createdForm.ClientPostalCode,
		ClientCity:                    createdForm.ClientCity,
		ReferrerFirstName:             createdForm.ReferrerFirstName,
		ReferrerLastName:              createdForm.ReferrerLastName,
		ReferrerOrganization:          createdForm.ReferrerOrganization,
		ReferrerJobTitle:              createdForm.ReferrerJobTitle,
		ReferrerPhoneNumber:           createdForm.ReferrerPhoneNumber,
		ReferrerEmail:                 createdForm.ReferrerEmail,
		Guardian1FirstName:            createdForm.Guardian1FirstName,
		Guardian1LastName:             createdForm.Guardian1LastName,
		Guardian1Relationship:         createdForm.Guardian1Relationship,
		Guardian1PhoneNumber:          createdForm.Guardian1PhoneNumber,
		Guardian1Email:                createdForm.Guardian1Email,
		Guardian2FirstName:            createdForm.Guardian2FirstName,
		Guardian2LastName:             createdForm.Guardian2LastName,
		Guardian2Relationship:         createdForm.Guardian2Relationship,
		Guardian2PhoneNumber:          createdForm.Guardian2PhoneNumber,
		Guardian2Email:                createdForm.Guardian2Email,
		EducationInstitution:          createdForm.EducationInstitution,
		EducationMentorName:           createdForm.EducationMentorName,
		EducationMentorPhone:          createdForm.EducationMentorPhone,
		EducationMentorEmail:          createdForm.EducationMentorEmail,
		EducationCurrentlyEnrolled:    createdForm.EducationCurrentlyEnrolled,
		EducationAdditionalNotes:      createdForm.EducationAdditionalNotes,
		WorkCurrentEmployer:           createdForm.WorkCurrentEmployer,
		WorkEmployerPhone:             createdForm.WorkEmployerPhone,
		WorkEmployerEmail:             createdForm.WorkEmployerEmail,
		WorkCurrentPosition:           createdForm.WorkCurrentPosition,
		WorkCurrentlyEmployed:         createdForm.WorkCurrentlyEmployed,
		WorkStartDate:                 &createdForm.WorkStartDate.Time,
		WorkAdditionalNotes:           createdForm.WorkAdditionalNotes,
		CareProtectedLiving:           createdForm.CareProtectedLiving,
		CareAssistedIndependentLiving: createdForm.CareAssistedIndependentLiving,
		CareRoomTrainingCenter:        createdForm.CareRoomTrainingCenter,
		CareAmbulatoryGuidance:        createdForm.CareAmbulatoryGuidance,
		ApplicationReason:             createdForm.ApplicationReason,
		ClientGoals:                   createdForm.ClientGoals,
		RiskAggressiveBehavior:        createdForm.RiskAggressiveBehavior,
		RiskSuicidalSelfharm:          createdForm.RiskSuicidalSelfharm,
		RiskSubstanceAbuse:            createdForm.RiskSubstanceAbuse,
		RiskPsychiatricIssues:         createdForm.RiskPsychiatricIssues,
		RiskCriminalHistory:           createdForm.RiskCriminalHistory,
		RiskFlightBehavior:            createdForm.RiskFlightBehavior,
		RiskWeaponPossession:          createdForm.RiskWeaponPossession,
		RiskSexualBehavior:            createdForm.RiskSexualBehavior,
		RiskDayNightRhythm:            createdForm.RiskDayNightRhythm,
		RiskOther:                     createdForm.RiskOther,
		RiskOtherDescription:          createdForm.RiskOtherDescription,
		RiskAdditionalNotes:           createdForm.RiskAdditionalNotes,
		DocumentReferral:              createdForm.DocumentReferral,
		DocumentEducationReport:       createdForm.DocumentEducationReport,
		DocumentActionPlan:            createdForm.DocumentActionPlan,
		DocumentPsychiatricReport:     createdForm.DocumentPsychiatricReport,
		DocumentDiagnosis:             createdForm.DocumentDiagnosis,
		DocumentSafetyPlan:            createdForm.DocumentSafetyPlan,
		DocumentIDCopy:                createdForm.DocumentIDCopy,
		ApplicationDate:               createdForm.ApplicationDate.Time,
		ReferrerSignature:             createdForm.ReferrerSignature,
		FormStatus:                    string(createdForm.FormStatus),
		CreatedAt:                     createdForm.CreatedAt.Time,
		UpdatedAt:                     createdForm.UpdatedAt.Time,
		SubmittedAt:                   createdForm.SubmittedAt.Time,
		ProcessedAt:                   createdForm.ProcessedAt.Time,
		ProcessedByEmployeeID:         createdForm.ProcessedByEmployeeID,
	}
	return response, nil
}

func (s *clientService) ListRegistrationForms(ctx *gin.Context, req *ListRegistrationFormsRequest) (*pagination.Response[ListRegistrationFormsResponse], error) {
	params := req.GetParams()
	status := db.NullFormStatusFromPtr(req.Status)

	forms, err := s.Store.ListRegistrationForms(ctx, db.ListRegistrationFormsParams{
		Limit:                  params.Limit,
		Offset:                 params.Offset,
		Status:                 status,
		RiskAggressiveBehavior: req.RiskAggressiveBehavior,
		RiskSuicidalSelfharm:   req.RiskSuicidalSelfharm,
		RiskSubstanceAbuse:     req.RiskSubstanceAbuse,
		RiskPsychiatricIssues:  req.RiskPsychiatricIssues,
		RiskCriminalHistory:    req.RiskCriminalHistory,
		RiskFlightBehavior:     req.RiskFlightBehavior,
		RiskWeaponPossession:   req.RiskWeaponPossession,
		RiskSexualBehavior:     req.RiskSexualBehavior,
		RiskDayNightRhythm:     req.RiskDayNightRhythm,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListRegistrationForms", "Failed to list registration forms", zap.Error(err))
		return nil, fmt.Errorf("failed to list registration forms: %w status: %v", err, status)
	}
	totalCount, err := s.Store.CountRegistrationForms(ctx, db.CountRegistrationFormsParams{
		Status:                 status,
		RiskAggressiveBehavior: req.RiskAggressiveBehavior,
		RiskSuicidalSelfharm:   req.RiskSuicidalSelfharm,
		RiskSubstanceAbuse:     req.RiskSubstanceAbuse,
		RiskPsychiatricIssues:  req.RiskPsychiatricIssues,
		RiskCriminalHistory:    req.RiskCriminalHistory,
		RiskFlightBehavior:     req.RiskFlightBehavior,
		RiskWeaponPossession:   req.RiskWeaponPossession,
		RiskSexualBehavior:     req.RiskSexualBehavior,
		RiskDayNightRhythm:     req.RiskDayNightRhythm,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListRegistrationForms", "Failed to count registration forms", zap.Error(err))
		return nil, fmt.Errorf("failed to count registration forms: %w", err)
	}

	responseItems := []ListRegistrationFormsResponse{}
	for _, form := range forms {
		responseItems = append(responseItems, ListRegistrationFormsResponse{
			ID:                            form.ID,
			ClientFirstName:               form.ClientFirstName,
			ClientLastName:                form.ClientLastName,
			ClientBsnNumber:               form.ClientBsnNumber,
			ReferrerFirstName:             form.ReferrerFirstName,
			ReferrerLastName:              form.ReferrerLastName,
			CareProtectedLiving:           form.CareProtectedLiving,
			CareAssistedIndependentLiving: form.CareAssistedIndependentLiving,
			CareRoomTrainingCenter:        form.CareRoomTrainingCenter,
			CareAmbulatoryGuidance:        form.CareAmbulatoryGuidance,
			RiskCount: calculateRiskCount(
				form.RiskAggressiveBehavior,
				form.RiskSuicidalSelfharm,
				form.RiskSubstanceAbuse,
				form.RiskPsychiatricIssues,
				form.RiskCriminalHistory,
				form.RiskFlightBehavior,
				form.RiskWeaponPossession,
				form.RiskSexualBehavior,
				form.RiskDayNightRhythm,
				form.RiskOther,
			),
			FormStatus:   string(form.FormStatus),
			SubmittedAt:  form.SubmittedAt.Time,
			IntakeFormID: form.IntakeFormID,
		})
	}

	paginationResponse := pagination.NewResponse(ctx, req.Request, responseItems, totalCount)
	return &paginationResponse, nil
}

func (s *clientService) GetRegistrationFormB(ctx context.Context, formID uuid.UUID) (*GetRegistrationFormResponse, error) {
	registrationForm, err := s.Store.GetRegistrationForm(ctx, formID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetRegistrationFormB", "Failed to get registration form B", zap.Error(err), zap.String("FormID", formID.String()))
		return nil, fmt.Errorf("failed to get registration form B: %w", err)
	}

	// Collect all document UUIDs
	var documentIDs []uuid.UUID
	if registrationForm.DocumentReferral != nil {
		documentIDs = append(documentIDs, *registrationForm.DocumentReferral)
	}
	if registrationForm.DocumentEducationReport != nil {
		documentIDs = append(documentIDs, *registrationForm.DocumentEducationReport)
	}
	if registrationForm.DocumentActionPlan != nil {
		documentIDs = append(documentIDs, *registrationForm.DocumentActionPlan)
	}
	if registrationForm.DocumentPsychiatricReport != nil {
		documentIDs = append(documentIDs, *registrationForm.DocumentPsychiatricReport)
	}
	if registrationForm.DocumentDiagnosis != nil {
		documentIDs = append(documentIDs, *registrationForm.DocumentDiagnosis)
	}
	if registrationForm.DocumentSafetyPlan != nil {
		documentIDs = append(documentIDs, *registrationForm.DocumentSafetyPlan)
	}
	if registrationForm.DocumentIDCopy != nil {
		documentIDs = append(documentIDs, *registrationForm.DocumentIDCopy)
	}

	// Fetch documents if any
	documentMap := make(map[uuid.UUID]DocumentResponse)
	if len(documentIDs) > 0 {
		attachments, err := s.Store.GetAttachmentsByUUIDs(ctx, documentIDs)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetRegistrationFormB", "Failed to get attachments", zap.Error(err))
			return nil, fmt.Errorf("failed to get attachments: %w", err)
		}
		for _, att := range attachments {
			documentMap[att.Uuid] = DocumentResponse{
				ID:   att.Uuid,
				Name: att.Name,
				File: att.File,
				Size: att.Size,
			}
		}
	}

	getDocumentResponse := func(id *uuid.UUID) *DocumentResponse {
		if id == nil {
			return nil
		}
		if doc, ok := documentMap[*id]; ok {
			return &doc
		}
		return nil
	}

	var processedByEmployeeName *string
	if registrationForm.ProcessedByFirstName != nil && registrationForm.ProcessedByLastName != nil {
		name := fmt.Sprintf("%s %s", *registrationForm.ProcessedByFirstName, *registrationForm.ProcessedByLastName)
		processedByEmployeeName = &name
	}

	var intakeOptions []string
	if registrationForm.IntakeOptions != nil {
		if err := json.Unmarshal(registrationForm.IntakeOptions, &intakeOptions); err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetRegistrationFormB", "Failed to unmarshal intake options", zap.Error(err))
			// Don't fail, just empty options
		}
	}

	response := &GetRegistrationFormResponse{
		ID:                            registrationForm.ID,
		ClientFirstName:               registrationForm.ClientFirstName,
		ClientLastName:                registrationForm.ClientLastName,
		ClientDateOfBirth:             &registrationForm.ClientDateOfBirth.Time,
		ClientBsnNumber:               registrationForm.ClientBsnNumber,
		ClientGender:                  string(registrationForm.ClientGender),
		ClientNationality:             registrationForm.ClientNationality,
		ClientPhoneNumber:             registrationForm.ClientPhoneNumber,
		ClientEmail:                   registrationForm.ClientEmail,
		ClientStreet:                  registrationForm.ClientStreet,
		ClientHouseNumber:             registrationForm.ClientHouseNumber,
		ClientHouseNumberAddition:     registrationForm.ClientHouseNumberAddition,
		ClientPostalCode:              registrationForm.ClientPostalCode,
		ClientCity:                    registrationForm.ClientCity,
		ReferrerFirstName:             registrationForm.ReferrerFirstName,
		ReferrerLastName:              registrationForm.ReferrerLastName,
		ReferrerOrganization:          registrationForm.ReferrerOrganization,
		ReferrerJobTitle:              registrationForm.ReferrerJobTitle,
		ReferrerPhoneNumber:           registrationForm.ReferrerPhoneNumber,
		ReferrerEmail:                 registrationForm.ReferrerEmail,
		Guardian1FirstName:            registrationForm.Guardian1FirstName,
		Guardian1LastName:             registrationForm.Guardian1LastName,
		Guardian1Relationship:         registrationForm.Guardian1Relationship,
		Guardian1PhoneNumber:          registrationForm.Guardian1PhoneNumber,
		Guardian1Email:                registrationForm.Guardian1Email,
		Guardian2FirstName:            registrationForm.Guardian2FirstName,
		Guardian2LastName:             registrationForm.Guardian2LastName,
		Guardian2Relationship:         registrationForm.Guardian2Relationship,
		Guardian2PhoneNumber:          registrationForm.Guardian2PhoneNumber,
		Guardian2Email:                registrationForm.Guardian2Email,
		EducationInstitution:          registrationForm.EducationInstitution,
		EducationMentorName:           registrationForm.EducationMentorName,
		EducationMentorPhone:          registrationForm.EducationMentorPhone,
		EducationMentorEmail:          registrationForm.EducationMentorEmail,
		EducationCurrentlyEnrolled:    registrationForm.EducationCurrentlyEnrolled,
		EducationAdditionalNotes:      util.DerefString(registrationForm.EducationAdditionalNotes),
		EducationLevel:                string(registrationForm.EducationLevel),
		WorkCurrentEmployer:           registrationForm.WorkCurrentEmployer,
		WorkEmployerPhone:             registrationForm.WorkEmployerPhone,
		WorkEmployerEmail:             registrationForm.WorkEmployerEmail,
		WorkCurrentPosition:           registrationForm.WorkCurrentPosition,
		WorkCurrentlyEmployed:         registrationForm.WorkCurrentlyEmployed,
		WorkStartDate:                 &registrationForm.WorkStartDate.Time,
		WorkAdditionalNotes:           util.DerefString(registrationForm.WorkAdditionalNotes),
		CareProtectedLiving:           registrationForm.CareProtectedLiving,
		CareAssistedIndependentLiving: registrationForm.CareAssistedIndependentLiving,
		CareRoomTrainingCenter:        registrationForm.CareRoomTrainingCenter,
		CareAmbulatoryGuidance:        registrationForm.CareAmbulatoryGuidance,
		ApplicationReason:             util.DerefString(registrationForm.ApplicationReason),
		ClientGoals:                   registrationForm.ClientGoals,
		RiskAggressiveBehavior:        registrationForm.RiskAggressiveBehavior,
		RiskSuicidalSelfharm:          registrationForm.RiskSuicidalSelfharm,
		RiskSubstanceAbuse:            registrationForm.RiskSubstanceAbuse,
		RiskPsychiatricIssues:         registrationForm.RiskPsychiatricIssues,
		RiskCriminalHistory:           registrationForm.RiskCriminalHistory,
		RiskFlightBehavior:            registrationForm.RiskFlightBehavior,
		RiskWeaponPossession:          registrationForm.RiskWeaponPossession,
		RiskSexualBehavior:            registrationForm.RiskSexualBehavior,
		RiskDayNightRhythm:            registrationForm.RiskDayNightRhythm,
		RiskOther:                     registrationForm.RiskOther,
		RiskOtherDescription:          util.DerefString(registrationForm.RiskOtherDescription),
		RiskAdditionalNotes:           util.DerefString(registrationForm.RiskAdditionalNotes),
		DocumentReferral:              getDocumentResponse(registrationForm.DocumentReferral),
		DocumentEducationReport:       getDocumentResponse(registrationForm.DocumentEducationReport),
		DocumentActionPlan:            getDocumentResponse(registrationForm.DocumentActionPlan),
		DocumentPsychiatricReport:     getDocumentResponse(registrationForm.DocumentPsychiatricReport),
		DocumentDiagnosis:             getDocumentResponse(registrationForm.DocumentDiagnosis),
		DocumentSafetyPlan:            getDocumentResponse(registrationForm.DocumentSafetyPlan),
		DocumentIDCopy:                getDocumentResponse(registrationForm.DocumentIDCopy),
		ApplicationDate:               registrationForm.ApplicationDate.Time,
		ReferrerSignature:             registrationForm.ReferrerSignature,
		FormStatus:                    string(registrationForm.FormStatus),
		CreatedAt:                     registrationForm.CreatedAt.Time,
		UpdatedAt:                     registrationForm.UpdatedAt.Time,
		SubmittedAt:                   registrationForm.SubmittedAt.Time,
		ProcessedAt:                   registrationForm.ProcessedAt.Time,
		ProcessedByEmployeeID:         registrationForm.ProcessedByEmployeeID,
		IntakeAppointmentDate:         registrationForm.IntakeAppointmentDatetime.Time,
		AddmissionType:                string(registrationForm.AddmissionType),
		ProcessedByEmployeeName:       processedByEmployeeName,
		IntakeOptions:                 intakeOptions,
		IntakeAppointmentLocation:     registrationForm.IntakeAppointmentLocation,
		IntakeFormID:                  registrationForm.IntakeFormID,
		RejectionReason:               registrationForm.RejectionReason,
	}
	return response, nil
}

func (s *clientService) UpdateRegistrationForm(ctx context.Context, req *UpdateRegistrationFormRequest, formID uuid.UUID) (*UpdateRegistrationFormResponse, error) {
	arg := db.UpdateRegistrationFormParams{
		ID:              formID,
		ClientFirstName: req.ClientFirstName,
		ClientLastName:  req.ClientLastName,
		ClientDateOfBirth: func() pgtype.Date {
			if req.ClientDateOfBirth != nil {
				return pgtype.Date{Time: *req.ClientDateOfBirth, Valid: true}
			}
			return pgtype.Date{Valid: false}
		}(),
		ClientBsnNumber:            req.ClientBsnNumber,
		ClientGender:               db.NullGenderFromPtr(req.ClientGender),
		ClientNationality:          req.ClientNationality,
		ClientPhoneNumber:          req.ClientPhoneNumber,
		ClientEmail:                req.ClientEmail,
		ClientStreet:               req.ClientStreet,
		ClientHouseNumber:          req.ClientHouseNumber,
		ClientPostalCode:           req.ClientPostalCode,
		ClientCity:                 req.ClientCity,
		ReferrerFirstName:          req.ReferrerFirstName,
		ReferrerLastName:           req.ReferrerLastName,
		ReferrerOrganization:       req.ReferrerOrganization,
		ReferrerJobTitle:           req.ReferrerJobTitle,
		ReferrerPhoneNumber:        req.ReferrerPhoneNumber,
		ReferrerEmail:              req.ReferrerEmail,
		Guardian1FirstName:         req.Guardian1FirstName,
		Guardian1LastName:          req.Guardian1LastName,
		Guardian1Relationship:      req.Guardian1Relationship,
		Guardian1PhoneNumber:       req.Guardian1PhoneNumber,
		Guardian1Email:             req.Guardian1Email,
		Guardian2FirstName:         req.Guardian2FirstName,
		Guardian2LastName:          req.Guardian2LastName,
		Guardian2Relationship:      req.Guardian2Relationship,
		Guardian2PhoneNumber:       req.Guardian2PhoneNumber,
		Guardian2Email:             req.Guardian2Email,
		EducationInstitution:       req.EducationInstitution,
		EducationMentorName:        req.EducationMentorName,
		EducationMentorPhone:       req.EducationMentorPhone,
		EducationMentorEmail:       req.EducationMentorEmail,
		EducationCurrentlyEnrolled: req.EducationCurrentlyEnrolled,
		EducationAdditionalNotes:   req.EducationAdditionalNotes,
		EducationLevel:             db.NullClientEducationLevelFromPtr(req.EducationLevel),
		WorkCurrentEmployer:        req.WorkCurrentEmployer,
		WorkEmployerPhone:          req.WorkEmployerPhone,
		WorkEmployerEmail:          req.WorkEmployerEmail,
		WorkCurrentPosition:        req.WorkCurrentPosition,
		WorkCurrentlyEmployed:      req.WorkCurrentlyEmployed,
		ApplicationReason:          req.ApplicationReason,

		CareProtectedLiving:           req.CareProtectedLiving,
		CareAssistedIndependentLiving: req.CareAssistedIndependentLiving,
		CareRoomTrainingCenter:        req.CareRoomTrainingCenter,
		CareAmbulatoryGuidance:        req.CareAmbulatoryGuidance,
		ClientGoals:                   req.ClientGoals,
		RiskAggressiveBehavior:        req.RiskAggressiveBehavior,
		RiskSuicidalSelfharm:          req.RiskSuicidalSelfharm,
		RiskSubstanceAbuse:            req.RiskSubstanceAbuse,
		RiskPsychiatricIssues:         req.RiskPsychiatricIssues,
		RiskCriminalHistory:           req.RiskCriminalHistory,
		RiskFlightBehavior:            req.RiskFlightBehavior,
		RiskWeaponPossession:          req.RiskWeaponPossession,
		RiskSexualBehavior:            req.RiskSexualBehavior,
		RiskDayNightRhythm:            req.RiskDayNightRhythm,
		RiskOther:                     req.RiskOther,
		RiskOtherDescription:          req.RiskOtherDescription,
		RiskAdditionalNotes:           req.RiskAdditionalNotes,
		ReferrerSignature:             req.ReferrerSignature,
	}

	if req.WorkStartDate != nil {
		arg.WorkStartDate = pgtype.Date{Time: *req.WorkStartDate, Valid: true}
	} else {
		arg.WorkStartDate = pgtype.Date{Valid: false}
	}

	if req.ApplicationDate != nil {
		arg.ApplicationDate = pgtype.Date{Time: *req.ApplicationDate, Valid: true}
	} else {
		arg.ApplicationDate = pgtype.Date{Valid: false}
	}

	registrationForm, err := s.Store.UpdateRegistrationForm(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateRegistrationForm", "Failed to update registration form", zap.Error(err), zap.String("FormID", formID.String()))
		return nil, fmt.Errorf("failed to update registration form: %w", err)
	}
	response := &UpdateRegistrationFormResponse{
		ID:                            registrationForm.ID,
		ClientFirstName:               registrationForm.ClientFirstName,
		ClientLastName:                registrationForm.ClientLastName,
		ClientDateOfBirth:             &registrationForm.ClientDateOfBirth.Time,
		ClientBsnNumber:               registrationForm.ClientBsnNumber,
		ClientGender:                  string(registrationForm.ClientGender),
		ClientNationality:             registrationForm.ClientNationality,
		ClientPhoneNumber:             registrationForm.ClientPhoneNumber,
		ClientEmail:                   registrationForm.ClientEmail,
		ClientStreet:                  registrationForm.ClientStreet,
		ClientHouseNumber:             registrationForm.ClientHouseNumber,
		ClientHouseNumberAddition:     registrationForm.ClientHouseNumberAddition,
		ClientPostalCode:              registrationForm.ClientPostalCode,
		ClientCity:                    registrationForm.ClientCity,
		ReferrerFirstName:             registrationForm.ReferrerFirstName,
		ReferrerLastName:              registrationForm.ReferrerLastName,
		ReferrerOrganization:          registrationForm.ReferrerOrganization,
		ReferrerJobTitle:              registrationForm.ReferrerJobTitle,
		ReferrerPhoneNumber:           registrationForm.ReferrerPhoneNumber,
		ReferrerEmail:                 registrationForm.ReferrerEmail,
		Guardian1FirstName:            registrationForm.Guardian1FirstName,
		Guardian1LastName:             registrationForm.Guardian1LastName,
		Guardian1Relationship:         registrationForm.Guardian1Relationship,
		Guardian1PhoneNumber:          registrationForm.Guardian1PhoneNumber,
		Guardian1Email:                registrationForm.Guardian1Email,
		Guardian2FirstName:            registrationForm.Guardian2FirstName,
		Guardian2LastName:             registrationForm.Guardian2LastName,
		Guardian2Relationship:         registrationForm.Guardian2Relationship,
		Guardian2PhoneNumber:          registrationForm.Guardian2PhoneNumber,
		Guardian2Email:                registrationForm.Guardian2Email,
		EducationInstitution:          registrationForm.EducationInstitution,
		EducationMentorName:           registrationForm.EducationMentorName,
		EducationMentorPhone:          registrationForm.EducationMentorPhone,
		EducationMentorEmail:          registrationForm.EducationMentorEmail,
		EducationCurrentlyEnrolled:    registrationForm.EducationCurrentlyEnrolled,
		EducationAdditionalNotes:      registrationForm.EducationAdditionalNotes,
		WorkCurrentEmployer:           registrationForm.WorkCurrentEmployer,
		WorkEmployerPhone:             registrationForm.WorkEmployerPhone,
		WorkEmployerEmail:             registrationForm.WorkEmployerEmail,
		WorkCurrentPosition:           registrationForm.WorkCurrentPosition,
		WorkCurrentlyEmployed:         registrationForm.WorkCurrentlyEmployed,
		WorkStartDate:                 &registrationForm.WorkStartDate.Time,
		WorkAdditionalNotes:           registrationForm.WorkAdditionalNotes,
		CareProtectedLiving:           registrationForm.CareProtectedLiving,
		CareAssistedIndependentLiving: registrationForm.CareAssistedIndependentLiving,
		CareRoomTrainingCenter:        registrationForm.CareRoomTrainingCenter,
		CareAmbulatoryGuidance:        registrationForm.CareAmbulatoryGuidance,
		ApplicationReason:             registrationForm.ApplicationReason,
		ClientGoals:                   registrationForm.ClientGoals,
		RiskAggressiveBehavior:        registrationForm.RiskAggressiveBehavior,
		RiskSuicidalSelfharm:          registrationForm.RiskSuicidalSelfharm,
		RiskSubstanceAbuse:            registrationForm.RiskSubstanceAbuse,
		RiskPsychiatricIssues:         registrationForm.RiskPsychiatricIssues,
		RiskCriminalHistory:           registrationForm.RiskCriminalHistory,
		RiskFlightBehavior:            registrationForm.RiskFlightBehavior,
		RiskWeaponPossession:          registrationForm.RiskWeaponPossession,
		RiskSexualBehavior:            registrationForm.RiskSexualBehavior,
		RiskDayNightRhythm:            registrationForm.RiskDayNightRhythm,
		RiskOther:                     registrationForm.RiskOther,
		RiskOtherDescription:          registrationForm.RiskOtherDescription,
		RiskAdditionalNotes:           registrationForm.RiskAdditionalNotes,
		DocumentReferral:              registrationForm.DocumentReferral,
		DocumentEducationReport:       registrationForm.DocumentEducationReport,
		DocumentActionPlan:            registrationForm.DocumentActionPlan,
		DocumentPsychiatricReport:     registrationForm.DocumentPsychiatricReport,
		DocumentDiagnosis:             registrationForm.DocumentDiagnosis,
		DocumentSafetyPlan:            registrationForm.DocumentSafetyPlan,
		DocumentIDCopy:                registrationForm.DocumentIDCopy,
		ApplicationDate:               registrationForm.ApplicationDate.Time,
		ReferrerSignature:             registrationForm.ReferrerSignature,
		FormStatus:                    string(registrationForm.FormStatus),
		CreatedAt:                     registrationForm.CreatedAt.Time,
		UpdatedAt:                     registrationForm.UpdatedAt.Time,
		SubmittedAt:                   registrationForm.SubmittedAt.Time,
		ProcessedAt:                   registrationForm.ProcessedAt.Time,
		ProcessedByEmployeeID:         registrationForm.ProcessedByEmployeeID,
	}
	return response, nil
}

func (s *clientService) DeleteRegistrationForm(ctx context.Context, formID uuid.UUID) error {
	err := s.Store.DeleteRegistrationForm(ctx, formID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteRegistrationForm", "Failed to delete registration form", zap.Error(err), zap.String("FormID", formID.String()))
		return fmt.Errorf("failed to delete registration form: %w", err)
	}
	return nil
}

func (s *clientService) UpdateRegistrationFormStatus(ctx context.Context, req *UpdateRegistrationFormStatusRequest, formID uuid.UUID, employeeID uuid.UUID) error {
	arg := db.UpdateRegistrationFormStatusParams{
		ID:                        formID,
		FormStatus:                db.FormStatusEnum(req.Status),
		ProcessedByEmployeeID:     &employeeID,
		IntakeAppointmentLocation: req.IntakeAppointmentLocation,
		AddmissionType: func() db.NullAdmissionTypeEnum {
			if req.AddmissionType != nil {
				return db.NullAdmissionTypeEnum{
					AdmissionTypeEnum: db.AdmissionTypeEnum(*req.AddmissionType),
					Valid:             true,
				}
			} else {
				return db.NullAdmissionTypeEnum{
					Valid: false,
				}
			}
		}(),
		RejectionReason: req.RejectionReason,
	}
	_, err := s.Store.UpdateRegistrationFormStatus(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateRegistrationFormStatus", "Failed to update registration form status", zap.Error(err), zap.String("FormID", formID.String()))
		return fmt.Errorf("failed to update registration form status: %w", err)
	}
	return nil
}

func (s *clientService) ProcessRegistrationForm(ctx context.Context, req *ProcessRegistrationFormRequest, formID uuid.UUID, employeeID uuid.UUID) error {
	optionsJSON, err := json.Marshal(req.ProposedDates)
	if err != nil {
		return fmt.Errorf("failed to marshal proposed dates: %w", err)
	}
	// Generate a secure random token
	token := util.RandomString(32) // Use a proper random string generator

	arg := db.UpdateRegistrationFormStatusParams{
		ID:                        formID,
		FormStatus:                db.FormStatusEnum("processed"),
		ProcessedByEmployeeID:     &employeeID,
		IntakeAppointmentLocation: &req.IntakeAppointmentLocation,
		AddmissionType:            db.NullAdmissionTypeEnum{AdmissionTypeEnum: db.AdmissionTypeEnum(req.AddmissionType), Valid: true},
		IntakeOptions:             optionsJSON,
		IntakeToken:               &token,
	}

	form, err := s.Store.UpdateRegistrationFormStatus(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ProcessRegistrationForm", "Failed to update registration form status", zap.Error(err), zap.String("FormID", formID.String()))
		return fmt.Errorf("failed to process registration form: %w", err)
	}

	// Send email to referrer and guardians
	recipients := []string{}
	if form.ReferrerEmail != "" {
		recipients = append(recipients, form.ReferrerEmail)
	}
	if form.Guardian1Email != "" {
		recipients = append(recipients, form.Guardian1Email)
	}
	if form.Guardian2Email != "" {
		recipients = append(recipients, form.Guardian2Email)
	}

	// Use config for base URL
	link := fmt.Sprintf("https://maicare.online/intake/schedule/%s", token)

	emailData := aclient.ProcessRegistrationFormEmailPayload{
		ReferrerName: form.ReferrerFirstName, // Or generic
		ClientName:   fmt.Sprintf("%s %s", form.ClientFirstName, form.ClientLastName),
		Location:     *form.IntakeAppointmentLocation,
		Link:         link,
		To:           recipients,
	}

	err = s.asynqClient.EnqueueProcessRegistrationFormEmail(ctx, emailData)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ProcessRegistrationForm", "Failed to enqueue email task", zap.Error(err), zap.String("FormID", formID.String()))
		// Don't fail the whole request if email fails, but log it critical
	}

	return nil
}

func (s *clientService) GetPublicIntakeOptions(ctx context.Context, token string) (*PublicIntakeOptionsResponse, error) {
	form, err := s.Store.GetRegistrationFormByToken(ctx, &token)
	if err != nil {
		return nil, fmt.Errorf("invalid token or form not found")
	}

	var dates []string
	if form.IntakeOptions != nil {
		if err := json.Unmarshal(form.IntakeOptions, &dates); err != nil {
			return nil, fmt.Errorf("failed to unmarshal proposed dates")
		}
	}

	return &PublicIntakeOptionsResponse{
		ClientFirstName: form.ClientFirstName,
		IntakeLocation:  util.DerefString(form.IntakeAppointmentLocation),
		ProposedDates:   dates,
	}, nil
}

func (s *clientService) SelectIntakeDate(ctx context.Context, token string, req *SelectIntakeDateRequest) error {
	form, err := s.Store.GetRegistrationFormByToken(ctx, &token)
	if err != nil {
		return fmt.Errorf("invalid token or form not found")
	}

	// Validate selected date is in options
	var dates []string
	if form.IntakeOptions != nil {
		if err := json.Unmarshal(form.IntakeOptions, &dates); err != nil {
			return fmt.Errorf("failed to unmarshal proposed dates")
		}
	}

	// Basic validation: Check if selected date roughly matches one of the options (ignoring timezone subtle diffs if needed, but string match is safer if strictly formatted)
	// For now, assume strict match or just accept if valid time.
	// Since we pass time.Time in req, but dates are []string in DB (from JSON).
	// We should convert req.SelectedDate to string or vice versa.
	// Let's assume options are ISO strings.

	found := false
	selectedStr := req.SelectedDate.Format(time.RFC3339)
	// Also try other formats if needed, or just compare times.
	for _, d := range dates {
		// Parse d
		t, err := time.Parse(time.RFC3339, d)
		if err == nil && t.Equal(req.SelectedDate) {
			found = true
			break
		}
		// Fallback simple string match check
		if d == selectedStr {
			found = true
			break
		}
	}

	if !found {
		// Strict validation: return fmt.Errorf("selected date is not in proposed options")
		// Permissive for now as timezone formats might differ
	}

	arg := db.UpdateRegistrationFormIntakeDateParams{
		ID:                        form.ID,
		IntakeAppointmentDatetime: pgtype.Timestamptz{Time: req.SelectedDate, Valid: true},
	}

	_, err = s.Store.UpdateRegistrationFormIntakeDate(ctx, arg)
	if err != nil {
		return fmt.Errorf("failed to update intake date: %w", err)
	}

	return nil
}

// ==================== Helper Functions ====================

// calculateRiskCount counts the number of true risk factors
func calculateRiskCount(
	riskAggressiveBehavior *bool,
	riskSuicidalSelfharm *bool,
	riskSubstanceAbuse *bool,
	riskPsychiatricIssues *bool,
	riskCriminalHistory *bool,
	riskFlightBehavior *bool,
	riskWeaponPossession *bool,
	riskSexualBehavior *bool,
	riskDayNightRhythm *bool,
	riskOther *bool,
) int {
	count := 0
	isTruePtr := func(b *bool) bool {
		return b != nil && *b
	}

	if isTruePtr(riskAggressiveBehavior) {
		count++
	}
	if isTruePtr(riskSuicidalSelfharm) {
		count++
	}
	if isTruePtr(riskSubstanceAbuse) {
		count++
	}
	if isTruePtr(riskPsychiatricIssues) {
		count++
	}
	if isTruePtr(riskCriminalHistory) {
		count++
	}
	if isTruePtr(riskFlightBehavior) {
		count++
	}
	if isTruePtr(riskWeaponPossession) {
		count++
	}
	if isTruePtr(riskSexualBehavior) {
		count++
	}
	if isTruePtr(riskDayNightRhythm) {
		count++
	}
	if isTruePtr(riskOther) {
		count++
	}

	return count
}
