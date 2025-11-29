package clientp

import (
	"context"
	"fmt"

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
		ClientBsnNumber:               req.ClientBsnNumber,
		ClientGender:                  db.ClientGenderEnum(req.ClientGender),
		ClientNationality:             req.ClientNationality,
		ClientPhoneNumber:             req.ClientPhoneNumber,
		ClientEmail:                   req.ClientEmail,
		ClientStreet:                  req.ClientStreet,
		ClientHouseNumber:             req.ClientHouseNumber,
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
		EducationLevel:                db.NullClientEducationLevelFromPtr(req.EducationLevel),
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
	createdForm, err := s.Store.CreateRegistrationForm(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateRegistrationForm", "Failed to create registration form", zap.Error(err))
		return nil, fmt.Errorf("failed to create registration form: %w", err)
	}
	response := &CreateRegistrationFormResponse{
		ID:                            createdForm.ID,
		ClientFirstName:               createdForm.ClientFirstName,
		ClientLastName:                createdForm.ClientLastName,
		ClientBsnNumber:               createdForm.ClientBsnNumber,
		ClientGender:                  string(createdForm.ClientGender),
		ClientNationality:             createdForm.ClientNationality,
		ClientPhoneNumber:             createdForm.ClientPhoneNumber,
		ClientEmail:                   createdForm.ClientEmail,
		ClientStreet:                  createdForm.ClientStreet,
		ClientHouseNumber:             createdForm.ClientHouseNumber,
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
			ClientGender:                  string(form.ClientGender),
			ClientNationality:             form.ClientNationality,
			ClientPhoneNumber:             form.ClientPhoneNumber,
			ClientEmail:                   form.ClientEmail,
			ClientStreet:                  form.ClientStreet,
			ClientHouseNumber:             form.ClientHouseNumber,
			ClientPostalCode:              form.ClientPostalCode,
			ClientCity:                    form.ClientCity,
			ReferrerFirstName:             form.ReferrerFirstName,
			ReferrerLastName:              form.ReferrerLastName,
			ReferrerOrganization:          form.ReferrerOrganization,
			ReferrerJobTitle:              form.ReferrerJobTitle,
			ReferrerPhoneNumber:           form.ReferrerPhoneNumber,
			ReferrerEmail:                 form.ReferrerEmail,
			Guardian1FirstName:            form.Guardian1FirstName,
			Guardian1LastName:             form.Guardian1LastName,
			Guardian1Relationship:         form.Guardian1Relationship,
			Guardian1PhoneNumber:          form.Guardian1PhoneNumber,
			Guardian1Email:                form.Guardian1Email,
			Guardian2FirstName:            form.Guardian2FirstName,
			Guardian2LastName:             form.Guardian2LastName,
			Guardian2Relationship:         form.Guardian2Relationship,
			Guardian2PhoneNumber:          form.Guardian2PhoneNumber,
			Guardian2Email:                form.Guardian2Email,
			EducationInstitution:          form.EducationInstitution,
			EducationMentorName:           form.EducationMentorName,
			EducationMentorPhone:          form.EducationMentorPhone,
			EducationMentorEmail:          form.EducationMentorEmail,
			EducationCurrentlyEnrolled:    form.EducationCurrentlyEnrolled,
			EducationAdditionalNotes:      form.EducationAdditionalNotes,
			WorkCurrentEmployer:           form.WorkCurrentEmployer,
			WorkEmployerPhone:             form.WorkEmployerPhone,
			WorkEmployerEmail:             form.WorkEmployerEmail,
			WorkCurrentPosition:           form.WorkCurrentPosition,
			WorkCurrentlyEmployed:         form.WorkCurrentlyEmployed,
			WorkStartDate:                 &form.WorkStartDate.Time,
			WorkAdditionalNotes:           form.WorkAdditionalNotes,
			CareProtectedLiving:           form.CareProtectedLiving,
			CareAssistedIndependentLiving: form.CareAssistedIndependentLiving,
			CareRoomTrainingCenter:        form.CareRoomTrainingCenter,
			CareAmbulatoryGuidance:        form.CareAmbulatoryGuidance,
			ApplicationReason:             form.ApplicationReason,
			ClientGoals:                   form.ClientGoals,
			RiskAggressiveBehavior:        form.RiskAggressiveBehavior,
			RiskSuicidalSelfharm:          form.RiskSuicidalSelfharm,
			RiskSubstanceAbuse:            form.RiskSubstanceAbuse,
			RiskPsychiatricIssues:         form.RiskPsychiatricIssues,
			RiskCriminalHistory:           form.RiskCriminalHistory,
			RiskFlightBehavior:            form.RiskFlightBehavior,
			RiskWeaponPossession:          form.RiskWeaponPossession,
			RiskSexualBehavior:            form.RiskSexualBehavior,
			RiskDayNightRhythm:            form.RiskDayNightRhythm,
			RiskOther:                     form.RiskOther,
			RiskOtherDescription:          form.RiskOtherDescription,
			RiskAdditionalNotes:           form.RiskAdditionalNotes,
			DocumentReferral:              form.DocumentReferral,
			DocumentEducationReport:       form.DocumentEducationReport,
			DocumentActionPlan:            form.DocumentActionPlan,
			DocumentPsychiatricReport:     form.DocumentPsychiatricReport,
			DocumentDiagnosis:             form.DocumentDiagnosis,
			DocumentSafetyPlan:            form.DocumentSafetyPlan,
			DocumentIDCopy:                form.DocumentIDCopy,
			ApplicationDate:               form.ApplicationDate.Time,
			ReferrerSignature:             form.ReferrerSignature,
			FormStatus:                    string(form.FormStatus),
			CreatedAt:                     form.CreatedAt.Time,
			UpdatedAt:                     form.UpdatedAt.Time,
			SubmittedAt:                   form.SubmittedAt.Time,
			ProcessedAt:                   form.ProcessedAt.Time,
			ProcessedByEmployeeID:         form.ProcessedByEmployeeID,
			RiskCount:                     calculateRiskCount(form),
			IntakeAppointmentDate:         form.IntakeAppointmentDatetime.Time,
			AddmissionType:                form.AddmissionType,
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
	response := &GetRegistrationFormResponse{
		ID:                            registrationForm.ID,
		ClientFirstName:               registrationForm.ClientFirstName,
		ClientLastName:                registrationForm.ClientLastName,
		ClientBsnNumber:               registrationForm.ClientBsnNumber,
		ClientGender:                  string(registrationForm.ClientGender),
		ClientNationality:             registrationForm.ClientNationality,
		ClientPhoneNumber:             registrationForm.ClientPhoneNumber,
		ClientEmail:                   registrationForm.ClientEmail,
		ClientStreet:                  registrationForm.ClientStreet,
		ClientHouseNumber:             registrationForm.ClientHouseNumber,
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
		EducationLevel:                db.ClientEducationLevelPtrFromEnum(registrationForm.EducationLevel),
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
		ClientGoals:                   util.DerefString(registrationForm.ClientGoals),
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
		IntakeAppointmentDate:         registrationForm.IntakeAppointmentDatetime.Time,
		AddmissionType:                registrationForm.AddmissionType,
	}
	return response, nil
}

func (s *clientService) UpdateRegistrationForm(ctx context.Context, req *UpdateRegistrationFormRequest, formID uuid.UUID) (*UpdateRegistrationFormResponse, error) {
	arg := db.UpdateRegistrationFormParams{
		ID:                         formID,
		ClientFirstName:            req.ClientFirstName,
		ClientLastName:             req.ClientLastName,
		ClientBsnNumber:            req.ClientBsnNumber,
		ClientGender:               db.NullClientGenderFromPtr(req.ClientGender),
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

		CareProtectedLiving:           req.CareProtectedLiving,
		CareAssistedIndependentLiving: req.CareAssistedIndependentLiving,
		CareRoomTrainingCenter:        req.CareRoomTrainingCenter,
		CareAmbulatoryGuidance:        req.CareAmbulatoryGuidance,
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
		ApplicationDate:               req.ApplicationDate,
		ReferrerSignature:             req.ReferrerSignature,
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
		ClientBsnNumber:               registrationForm.ClientBsnNumber,
		ClientGender:                  string(registrationForm.ClientGender),
		ClientNationality:             registrationForm.ClientNationality,
		ClientPhoneNumber:             registrationForm.ClientPhoneNumber,
		ClientEmail:                   registrationForm.ClientEmail,
		ClientStreet:                  registrationForm.ClientStreet,
		ClientHouseNumber:             registrationForm.ClientHouseNumber,
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
		AddmissionType:            req.AddmissionType,
	}
	updatedForm, err := s.Store.UpdateRegistrationFormStatus(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateRegistrationFormStatus", "Failed to update registration form status", zap.Error(err), zap.String("FormID", formID.String()))
		return fmt.Errorf("failed to update registration form status: %w", err)
	}
	if req.Status == "approved" {
		err = s.asynqClient.EnqueueAcceptedRegistration(ctx, aclient.AcceptedRegistrationFormPayload{
			ReferrerName:        updatedForm.ReferrerFirstName + " " + updatedForm.ReferrerLastName,
			ChildName:           updatedForm.ClientFirstName + " " + updatedForm.ClientLastName,
			ChildBSN:            updatedForm.ClientBsnNumber,
			AppointmentDate:     updatedForm.IntakeAppointmentDatetime.Time.Format("2006-01-02 15:04:05"),
			AppointmentLocation: *updatedForm.IntakeAppointmentLocation,
			To:                  updatedForm.ReferrerEmail,
		})
	}
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateRegistrationFormStatus", "Failed to enqueue accepted registration email", zap.Error(err), zap.String("FormID", formID.String()))
	}
	return nil
}

// ==================== Helper Functions ====================

// calculateRiskCount counts the number of true risk factors for a registration form
func calculateRiskCount(rf db.RegistrationForm) int {
	count := 0

	// Helper function to check if a *bool is true
	isTruePtr := func(b *bool) bool {
		return b != nil && *b
	}

	if isTruePtr(rf.RiskAggressiveBehavior) {
		count++
	}
	if isTruePtr(rf.RiskSuicidalSelfharm) {
		count++
	}
	if isTruePtr(rf.RiskSubstanceAbuse) {
		count++
	}
	if isTruePtr(rf.RiskPsychiatricIssues) {
		count++
	}
	if isTruePtr(rf.RiskCriminalHistory) {
		count++
	}
	if isTruePtr(rf.RiskFlightBehavior) {
		count++
	}
	if isTruePtr(rf.RiskWeaponPossession) {
		count++
	}
	if isTruePtr(rf.RiskSexualBehavior) {
		count++
	}
	if isTruePtr(rf.RiskDayNightRhythm) {
		count++
	}
	if isTruePtr(rf.RiskOther) {
		count++
	}

	return count
}
