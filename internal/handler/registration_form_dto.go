package handler

import (
	"time"

	"maicare_go/internal/domain"
	"maicare_go/internal/httpapi"

	"github.com/google/uuid"
)

// ========== Request DTOs ==========

type createRegistrationFormRequest struct {
	ClientFirstName               string     `json:"client_first_name" binding:"required"`
	ClientLastName                string     `json:"client_last_name" binding:"required"`
	ClientDateOfBirth             *time.Time `json:"client_date_of_birth"`
	ClientBsnNumber               string     `json:"client_bsn_number" binding:"required"`
	ClientGender                  string     `json:"client_gender" binding:"required,oneof=male female other unknown"`
	ClientNationality             string     `json:"client_nationality" binding:"required"`
	ClientPhoneNumber             string     `json:"client_phone_number" binding:"required"`
	ClientEmail                   string     `json:"client_email" binding:"required,email"`
	ClientStreet                  string     `json:"client_street" binding:"required"`
	ClientHouseNumber             string     `json:"client_house_number" binding:"required"`
	ClientHouseNumberAddition     *string    `json:"client_house_number_addition"`
	ClientPostalCode              string     `json:"client_postal_code" binding:"required"`
	ClientCity                    string     `json:"client_city" binding:"required"`
	ReferrerFirstName             string     `json:"referrer_first_name" binding:"required"`
	ReferrerLastName              string     `json:"referrer_last_name" binding:"required"`
	ReferrerOrganization          string     `json:"referrer_organization" binding:"required"`
	ReferrerJobTitle              string     `json:"referrer_job_title" binding:"required"`
	ReferrerPhoneNumber           string     `json:"referrer_phone_number" binding:"required"`
	ReferrerEmail                 string     `json:"referrer_email" binding:"required,email"`
	Guardian1FirstName            string     `json:"guardian1_first_name" binding:"required"`
	Guardian1LastName             string     `json:"guardian1_last_name" binding:"required"`
	Guardian1Relationship         string     `json:"guardian1_relationship" binding:"required"`
	Guardian1PhoneNumber          string     `json:"guardian1_phone_number" binding:"required"`
	Guardian1Email                string     `json:"guardian1_email" binding:"required,email"`
	Guardian2FirstName            string     `json:"guardian2_first_name" binding:"required"`
	Guardian2LastName             string     `json:"guardian2_last_name" binding:"required"`
	Guardian2Relationship         string     `json:"guardian2_relationship" binding:"required"`
	Guardian2PhoneNumber          string     `json:"guardian2_phone_number" binding:"required"`
	Guardian2Email                string     `json:"guardian2_email" binding:"required,email"`
	EducationInstitution          *string    `json:"education_institution"`
	EducationMentorName           *string    `json:"education_mentor_name"`
	EducationMentorPhone          *string    `json:"education_mentor_phone"`
	EducationMentorEmail          *string    `json:"education_mentor_email"`
	EducationCurrentlyEnrolled    bool       `json:"education_currently_enrolled"`
	EducationAdditionalNotes      *string    `json:"education_additional_notes"`
	EducationLevel                string     `json:"education_level" binding:"omitempty,oneof=primary secondary higher none"`
	WorkCurrentEmployer           *string    `json:"work_current_employer"`
	WorkEmployerPhone             *string    `json:"work_employer_phone"`
	WorkEmployerEmail             *string    `json:"work_employer_email"`
	WorkCurrentPosition           *string    `json:"work_current_position"`
	WorkCurrentlyEmployed         bool       `json:"work_currently_employed"`
	WorkStartDate                 *time.Time `json:"work_start_date"`
	WorkAdditionalNotes           *string    `json:"work_additional_notes"`
	CareProtectedLiving           *bool      `json:"care_protected_living"`
	CareAssistedIndependentLiving *bool      `json:"care_assisted_independent_living"`
	CareRoomTrainingCenter        *bool      `json:"care_room_training_center"`
	CareAmbulatoryGuidance        *bool      `json:"care_ambulatory_guidance"`
	ClientGoals                   []string   `json:"client_goals"`
	ApplicationReason             *string    `json:"application_reason"`
	RiskAggressiveBehavior        *bool      `json:"risk_aggressive_behavior"`
	RiskSuicidalSelfharm          *bool      `json:"risk_suicidal_selfharm"`
	RiskSubstanceAbuse            *bool      `json:"risk_substance_abuse"`
	RiskPsychiatricIssues         *bool      `json:"risk_psychiatric_issues"`
	RiskCriminalHistory           *bool      `json:"risk_criminal_history"`
	RiskFlightBehavior            *bool      `json:"risk_flight_behavior"`
	RiskWeaponPossession          *bool      `json:"risk_weapon_possession"`
	RiskSexualBehavior            *bool      `json:"risk_sexual_behavior"`
	RiskDayNightRhythm            *bool      `json:"risk_day_night_rhythm"`
	RiskOther                     *bool      `json:"risk_other"`
	RiskOtherDescription          *string    `json:"risk_other_description"`
	RiskAdditionalNotes           *string    `json:"risk_additional_notes"`
	DocumentReferral              *uuid.UUID `json:"document_referral"`
	DocumentEducationReport       *uuid.UUID `json:"document_education_report"`
	DocumentPsychiatricReport     *uuid.UUID `json:"document_psychiatric_report"`
	DocumentDiagnosis             *uuid.UUID `json:"document_diagnosis"`
	DocumentSafetyPlan            *uuid.UUID `json:"document_safety_plan"`
	DocumentIDCopy                *uuid.UUID `json:"document_id_copy"`
	ApplicationDate               time.Time  `json:"application_date" binding:"required"`
	ReferrerSignature             *bool      `json:"referrer_signature"`
}

type listRegistrationFormsRequest struct {
	httpapi.PageRequest
	Status                 *string `form:"status" binding:"omitempty,oneof=pending processed rejected"`
	RiskAggressiveBehavior *bool   `form:"risk_aggressive_behavior"`
	RiskSuicidalSelfharm   *bool   `form:"risk_suicidal_selfharm"`
	RiskSubstanceAbuse     *bool   `form:"risk_substance_abuse"`
	RiskPsychiatricIssues  *bool   `form:"risk_psychiatric_issues"`
	RiskCriminalHistory    *bool   `form:"risk_criminal_history"`
	RiskFlightBehavior     *bool   `form:"risk_flight_behavior"`
	RiskWeaponPossession   *bool   `form:"risk_weapon_possession"`
	RiskSexualBehavior     *bool   `form:"risk_sexual_behavior"`
	RiskDayNightRhythm     *bool   `form:"risk_day_night_rhythm"`
}

type updateRegistrationFormRequest struct {
	ClientFirstName               *string    `json:"client_first_name"`
	ClientLastName                *string    `json:"client_last_name"`
	ClientDateOfBirth             *time.Time `json:"client_date_of_birth"`
	ClientBsnNumber               *string    `json:"client_bsn_number"`
	ClientGender                  *string    `json:"client_gender" binding:"omitempty,oneof=male female other unknown"`
	ClientNationality             *string    `json:"client_nationality"`
	ClientPhoneNumber             *string    `json:"client_phone_number"`
	ClientEmail                   *string    `json:"client_email"`
	ClientStreet                  *string    `json:"client_street"`
	ClientHouseNumber             *string    `json:"client_house_number"`
	ClientHouseNumberAddition     *string    `json:"client_house_number_addition"`
	ClientPostalCode              *string    `json:"client_postal_code"`
	ClientCity                    *string    `json:"client_city"`
	ReferrerFirstName             *string    `json:"referrer_first_name"`
	ReferrerLastName              *string    `json:"referrer_last_name"`
	ReferrerOrganization          *string    `json:"referrer_organization"`
	ReferrerJobTitle              *string    `json:"referrer_job_title"`
	ReferrerPhoneNumber           *string    `json:"referrer_phone_number"`
	ReferrerEmail                 *string    `json:"referrer_email"`
	Guardian1FirstName            *string    `json:"guardian1_first_name"`
	Guardian1LastName             *string    `json:"guardian1_last_name"`
	Guardian1Relationship         *string    `json:"guardian1_relationship"`
	Guardian1PhoneNumber          *string    `json:"guardian1_phone_number"`
	Guardian1Email                *string    `json:"guardian1_email"`
	Guardian2FirstName            *string    `json:"guardian2_first_name"`
	Guardian2LastName             *string    `json:"guardian2_last_name"`
	Guardian2Relationship         *string    `json:"guardian2_relationship"`
	Guardian2PhoneNumber          *string    `json:"guardian2_phone_number"`
	Guardian2Email                *string    `json:"guardian2_email"`
	EducationInstitution          *string    `json:"education_institution"`
	EducationMentorName           *string    `json:"education_mentor_name"`
	EducationMentorPhone          *string    `json:"education_mentor_phone"`
	EducationMentorEmail          *string    `json:"education_mentor_email"`
	EducationCurrentlyEnrolled    *bool      `json:"education_currently_enrolled"`
	EducationAdditionalNotes      *string    `json:"education_additional_notes"`
	EducationLevel                *string    `json:"education_level" binding:"omitempty,oneof=primary secondary higher none"`
	WorkCurrentEmployer           *string    `json:"work_current_employer"`
	WorkEmployerPhone             *string    `json:"work_employer_phone"`
	WorkEmployerEmail             *string    `json:"work_employer_email"`
	WorkCurrentPosition           *string    `json:"work_current_position"`
	WorkCurrentlyEmployed         *bool      `json:"work_currently_employed"`
	WorkStartDate                 *time.Time `json:"work_start_date"`
	WorkAdditionalNotes           *string    `json:"work_additional_notes"`
	CareProtectedLiving           *bool      `json:"care_protected_living"`
	CareAssistedIndependentLiving *bool      `json:"care_assisted_independent_living"`
	CareRoomTrainingCenter        *bool      `json:"care_room_training_center"`
	CareAmbulatoryGuidance        *bool      `json:"care_ambulatory_guidance"`
	ApplicationReason             *string    `json:"application_reason"`
	ClientGoals                   []string   `json:"client_goals"`
	RiskAggressiveBehavior        *bool      `json:"risk_aggressive_behavior"`
	RiskSuicidalSelfharm          *bool      `json:"risk_suicidal_selfharm"`
	RiskSubstanceAbuse            *bool      `json:"risk_substance_abuse"`
	RiskPsychiatricIssues         *bool      `json:"risk_psychiatric_issues"`
	RiskCriminalHistory           *bool      `json:"risk_criminal_history"`
	RiskFlightBehavior            *bool      `json:"risk_flight_behavior"`
	RiskWeaponPossession          *bool      `json:"risk_weapon_possession"`
	RiskSexualBehavior            *bool      `json:"risk_sexual_behavior"`
	RiskDayNightRhythm            *bool      `json:"risk_day_night_rhythm"`
	RiskOther                     *bool      `json:"risk_other"`
	RiskOtherDescription          *string    `json:"risk_other_description"`
	RiskAdditionalNotes           *string    `json:"risk_additional_notes"`
	ApplicationDate               *time.Time `json:"application_date"`
	ReferrerSignature             *bool      `json:"referrer_signature"`
}

type updateRegistrationFormStatusRequest struct {
	Status                    string    `json:"status" binding:"required,oneof=pending processed rejected"`
	IntakeAppointmentDate     time.Time `json:"intake_appointment_date"`
	IntakeAppointmentLocation *string   `json:"intake_appointment_location"`
	AddmissionType            *string   `json:"admission_type" binding:"omitempty,oneof=crisis_admission regular_placement"`
	RejectionReason           *string   `json:"rejection_reason"`
}

type processRegistrationFormRequest struct {
	IntakeAppointmentLocation string   `json:"intake_appointment_location" binding:"required"`
	AddmissionType            string   `json:"admission_type" binding:"required,oneof=crisis_admission regular_placement"`
	ProposedDates             []string `json:"proposed_dates" binding:"required,min=1"`
}

type selectIntakeDateRequest struct {
	SelectedDate time.Time `json:"selected_date" binding:"required"`
}

// ========== Response DTOs ==========

type registrationFormResponse struct {
	ID                            uuid.UUID         `json:"id"`
	ClientFirstName               string            `json:"client_first_name"`
	ClientLastName                string            `json:"client_last_name"`
	ClientDateOfBirth             *time.Time        `json:"client_date_of_birth"`
	ClientBsnNumber               string            `json:"client_bsn_number"`
	ClientGender                  string            `json:"client_gender"`
	ClientNationality             string            `json:"client_nationality"`
	ClientPhoneNumber             string            `json:"client_phone_number"`
	ClientEmail                   string            `json:"client_email"`
	ClientStreet                  string            `json:"client_street"`
	ClientHouseNumber             string            `json:"client_house_number"`
	ClientHouseNumberAddition     *string           `json:"client_house_number_addition"`
	ClientPostalCode              string            `json:"client_postal_code"`
	ClientCity                    string            `json:"client_city"`
	ReferrerFirstName             string            `json:"referrer_first_name"`
	ReferrerLastName              string            `json:"referrer_last_name"`
	ReferrerOrganization          string            `json:"referrer_organization"`
	ReferrerJobTitle              string            `json:"referrer_job_title"`
	ReferrerPhoneNumber           string            `json:"referrer_phone_number"`
	ReferrerEmail                 string            `json:"referrer_email"`
	Guardian1FirstName            string            `json:"guardian1_first_name"`
	Guardian1LastName             string            `json:"guardian1_last_name"`
	Guardian1Relationship         string            `json:"guardian1_relationship"`
	Guardian1PhoneNumber          string            `json:"guardian1_phone_number"`
	Guardian1Email                string            `json:"guardian1_email"`
	Guardian2FirstName            string            `json:"guardian2_first_name"`
	Guardian2LastName             string            `json:"guardian2_last_name"`
	Guardian2Relationship         string            `json:"guardian2_relationship"`
	Guardian2PhoneNumber          string            `json:"guardian2_phone_number"`
	Guardian2Email                string            `json:"guardian2_email"`
	EducationInstitution          *string           `json:"education_institution"`
	EducationMentorName           *string           `json:"education_mentor_name"`
	EducationMentorPhone          *string           `json:"education_mentor_phone"`
	EducationMentorEmail          *string           `json:"education_mentor_email"`
	EducationCurrentlyEnrolled    bool              `json:"education_currently_enrolled"`
	EducationAdditionalNotes      *string           `json:"education_additional_notes"`
	EducationLevel                string            `json:"education_level"`
	WorkCurrentEmployer           *string           `json:"work_current_employer"`
	WorkEmployerPhone             *string           `json:"work_employer_phone"`
	WorkEmployerEmail             *string           `json:"work_employer_email"`
	WorkCurrentPosition           *string           `json:"work_current_position"`
	WorkCurrentlyEmployed         bool              `json:"work_currently_employed"`
	WorkStartDate                 *time.Time        `json:"work_start_date"`
	WorkAdditionalNotes           *string           `json:"work_additional_notes"`
	CareProtectedLiving           *bool             `json:"care_protected_living"`
	CareAssistedIndependentLiving *bool             `json:"care_assisted_independent_living"`
	CareRoomTrainingCenter        *bool             `json:"care_room_training_center"`
	CareAmbulatoryGuidance        *bool             `json:"care_ambulatory_guidance"`
	ApplicationReason             *string           `json:"application_reason"`
	ClientGoals                   []string          `json:"client_goals"`
	RiskAggressiveBehavior        *bool             `json:"risk_aggressive_behavior"`
	RiskSuicidalSelfharm          *bool             `json:"risk_suicidal_selfharm"`
	RiskSubstanceAbuse            *bool             `json:"risk_substance_abuse"`
	RiskPsychiatricIssues         *bool             `json:"risk_psychiatric_issues"`
	RiskCriminalHistory           *bool             `json:"risk_criminal_history"`
	RiskFlightBehavior            *bool             `json:"risk_flight_behavior"`
	RiskWeaponPossession          *bool             `json:"risk_weapon_possession"`
	RiskSexualBehavior            *bool             `json:"risk_sexual_behavior"`
	RiskDayNightRhythm            *bool             `json:"risk_day_night_rhythm"`
	RiskOther                     *bool             `json:"risk_other"`
	RiskOtherDescription          *string           `json:"risk_other_description"`
	RiskAdditionalNotes           *string           `json:"risk_additional_notes"`
	DocumentReferral              *documentResponse `json:"document_referral"`
	DocumentEducationReport       *documentResponse `json:"document_education_report"`
	DocumentActionPlan            *documentResponse `json:"document_action_plan"`
	DocumentPsychiatricReport     *documentResponse `json:"document_psychiatric_report"`
	DocumentDiagnosis             *documentResponse `json:"document_diagnosis"`
	DocumentSafetyPlan            *documentResponse `json:"document_safety_plan"`
	DocumentIDCopy                *documentResponse `json:"document_id_copy"`
	ApplicationDate               time.Time         `json:"application_date"`
	ReferrerSignature             *bool             `json:"referrer_signature"`
	FormStatus                    string            `json:"form_status"`
	CreatedAt                     time.Time         `json:"created_at"`
	UpdatedAt                     time.Time         `json:"updated_at"`
	SubmittedAt                   time.Time         `json:"submitted_at"`
	ProcessedAt                   time.Time         `json:"processed_at"`
	ProcessedByEmployeeID         *uuid.UUID        `json:"processed_by_employee_id"`
	ProcessedByEmployeeName       *string           `json:"processed_by_employee_name"`
	IntakeAppointmentDate         time.Time         `json:"intake_appointment_date"`
	IntakeAppointmentLocation     *string           `json:"intake_appointment_location"`
	AddmissionType                string            `json:"admission_type"`
	IntakeOptions                 []string          `json:"intake_options"`
	IntakeFormID                  *uuid.UUID        `json:"intake_form_id"`
	RejectionReason               *string           `json:"rejection_reason"`
}

type registrationFormListItemResponse struct {
	ID                            uuid.UUID  `json:"id"`
	ClientFirstName               string     `json:"client_first_name"`
	ClientLastName                string     `json:"client_last_name"`
	ClientBsnNumber               string     `json:"client_bsn_number"`
	ReferrerFirstName             string     `json:"referrer_first_name"`
	ReferrerLastName              string     `json:"referrer_last_name"`
	CareProtectedLiving           *bool      `json:"care_protected_living"`
	CareAssistedIndependentLiving *bool      `json:"care_assisted_independent_living"`
	CareRoomTrainingCenter        *bool      `json:"care_room_training_center"`
	CareAmbulatoryGuidance        *bool      `json:"care_ambulatory_guidance"`
	RiskCount                     int        `json:"risk_count"`
	FormStatus                    string     `json:"form_status"`
	SubmittedAt                   time.Time  `json:"submitted_at"`
	IntakeFormID                  *uuid.UUID `json:"intake_form_id"`
}

type documentResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	File string    `json:"file"`
	Size int32     `json:"size"`
}

type publicIntakeOptionsResponse struct {
	ClientFirstName string   `json:"client_first_name"`
	IntakeLocation  string   `json:"intake_location"`
	ProposedDates   []string `json:"proposed_dates"`
}

// ========== Mappers ==========

func toRegistrationFormResponse(f domain.RegistrationForm) registrationFormResponse {
	return registrationFormResponse{
		ID:                            f.ID,
		ClientFirstName:               f.ClientFirstName,
		ClientLastName:                f.ClientLastName,
		ClientDateOfBirth:             f.ClientDateOfBirth,
		ClientBsnNumber:               f.ClientBsnNumber,
		ClientGender:                  f.ClientGender,
		ClientNationality:             f.ClientNationality,
		ClientPhoneNumber:             f.ClientPhoneNumber,
		ClientEmail:                   f.ClientEmail,
		ClientStreet:                  f.ClientStreet,
		ClientHouseNumber:             f.ClientHouseNumber,
		ClientHouseNumberAddition:     f.ClientHouseNumberAddition,
		ClientPostalCode:              f.ClientPostalCode,
		ClientCity:                    f.ClientCity,
		ReferrerFirstName:             f.ReferrerFirstName,
		ReferrerLastName:              f.ReferrerLastName,
		ReferrerOrganization:          f.ReferrerOrganization,
		ReferrerJobTitle:              f.ReferrerJobTitle,
		ReferrerPhoneNumber:           f.ReferrerPhoneNumber,
		ReferrerEmail:                 f.ReferrerEmail,
		Guardian1FirstName:            f.Guardian1FirstName,
		Guardian1LastName:             f.Guardian1LastName,
		Guardian1Relationship:         f.Guardian1Relationship,
		Guardian1PhoneNumber:          f.Guardian1PhoneNumber,
		Guardian1Email:                f.Guardian1Email,
		Guardian2FirstName:            f.Guardian2FirstName,
		Guardian2LastName:             f.Guardian2LastName,
		Guardian2Relationship:         f.Guardian2Relationship,
		Guardian2PhoneNumber:          f.Guardian2PhoneNumber,
		Guardian2Email:                f.Guardian2Email,
		EducationInstitution:          f.EducationInstitution,
		EducationMentorName:           f.EducationMentorName,
		EducationMentorPhone:          f.EducationMentorPhone,
		EducationMentorEmail:          f.EducationMentorEmail,
		EducationCurrentlyEnrolled:    f.EducationCurrentlyEnrolled,
		EducationAdditionalNotes:      f.EducationAdditionalNotes,
		EducationLevel:                f.EducationLevel,
		WorkCurrentEmployer:           f.WorkCurrentEmployer,
		WorkEmployerPhone:             f.WorkEmployerPhone,
		WorkEmployerEmail:             f.WorkEmployerEmail,
		WorkCurrentPosition:           f.WorkCurrentPosition,
		WorkCurrentlyEmployed:         f.WorkCurrentlyEmployed,
		WorkStartDate:                 f.WorkStartDate,
		WorkAdditionalNotes:           f.WorkAdditionalNotes,
		CareProtectedLiving:           f.CareProtectedLiving,
		CareAssistedIndependentLiving: f.CareAssistedIndependentLiving,
		CareRoomTrainingCenter:        f.CareRoomTrainingCenter,
		CareAmbulatoryGuidance:        f.CareAmbulatoryGuidance,
		ApplicationReason:             f.ApplicationReason,
		ClientGoals:                   f.ClientGoals,
		RiskAggressiveBehavior:        f.RiskAggressiveBehavior,
		RiskSuicidalSelfharm:          f.RiskSuicidalSelfharm,
		RiskSubstanceAbuse:            f.RiskSubstanceAbuse,
		RiskPsychiatricIssues:         f.RiskPsychiatricIssues,
		RiskCriminalHistory:           f.RiskCriminalHistory,
		RiskFlightBehavior:            f.RiskFlightBehavior,
		RiskWeaponPossession:          f.RiskWeaponPossession,
		RiskSexualBehavior:            f.RiskSexualBehavior,
		RiskDayNightRhythm:            f.RiskDayNightRhythm,
		RiskOther:                     f.RiskOther,
		RiskOtherDescription:          f.RiskOtherDescription,
		RiskAdditionalNotes:           f.RiskAdditionalNotes,
		DocumentReferral:              toDocumentResponse(f.DocumentReferral),
		DocumentEducationReport:       toDocumentResponse(f.DocumentEducationReport),
		DocumentActionPlan:            toDocumentResponse(f.DocumentActionPlan),
		DocumentPsychiatricReport:     toDocumentResponse(f.DocumentPsychiatricReport),
		DocumentDiagnosis:             toDocumentResponse(f.DocumentDiagnosis),
		DocumentSafetyPlan:            toDocumentResponse(f.DocumentSafetyPlan),
		DocumentIDCopy:                toDocumentResponse(f.DocumentIDCopy),
		ApplicationDate:               f.ApplicationDate,
		ReferrerSignature:             f.ReferrerSignature,
		FormStatus:                    f.FormStatus,
		CreatedAt:                     f.CreatedAt,
		UpdatedAt:                     f.UpdatedAt,
		SubmittedAt:                   f.SubmittedAt,
		ProcessedAt:                   f.ProcessedAt,
		ProcessedByEmployeeID:         f.ProcessedByEmployeeID,
		ProcessedByEmployeeName:       f.ProcessedByEmployeeName,
		IntakeAppointmentDate:         f.IntakeAppointmentDate,
		IntakeAppointmentLocation:     f.IntakeAppointmentLocation,
		AddmissionType:                f.AddmissionType,
		IntakeOptions:                 f.IntakeOptions,
		IntakeFormID:                  f.IntakeFormID,
		RejectionReason:               f.RejectionReason,
	}
}

func toDocumentResponse(d *domain.Document) *documentResponse {
	if d == nil {
		return nil
	}
	return &documentResponse{
		ID:   d.ID,
		Name: d.Name,
		File: d.File,
		Size: d.Size,
	}
}

func toRegistrationFormListItemResponse(item domain.RegistrationFormListItem) registrationFormListItemResponse {
	return registrationFormListItemResponse{
		ID:                            item.ID,
		ClientFirstName:               item.ClientFirstName,
		ClientLastName:                item.ClientLastName,
		ClientBsnNumber:               item.ClientBsnNumber,
		ReferrerFirstName:             item.ReferrerFirstName,
		ReferrerLastName:              item.ReferrerLastName,
		CareProtectedLiving:           item.CareProtectedLiving,
		CareAssistedIndependentLiving: item.CareAssistedIndependentLiving,
		CareRoomTrainingCenter:        item.CareRoomTrainingCenter,
		CareAmbulatoryGuidance:        item.CareAmbulatoryGuidance,
		RiskCount:                     item.RiskCount,
		FormStatus:                    item.FormStatus,
		SubmittedAt:                   item.SubmittedAt,
		IntakeFormID:                  item.IntakeFormID,
	}
}

func toPublicIntakeOptionsResponse(opts domain.PublicIntakeOptions) publicIntakeOptionsResponse {
	return publicIntakeOptionsResponse{
		ClientFirstName: opts.ClientFirstName,
		IntakeLocation:  opts.IntakeLocation,
		ProposedDates:   opts.ProposedDates,
	}
}

func toCreateRegistrationFormParams(req createRegistrationFormRequest) domain.CreateRegistrationFormParams {
	return domain.CreateRegistrationFormParams{
		ClientFirstName:               req.ClientFirstName,
		ClientLastName:                req.ClientLastName,
		ClientDateOfBirth:             req.ClientDateOfBirth,
		ClientBsnNumber:               req.ClientBsnNumber,
		ClientGender:                  req.ClientGender,
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
		EducationLevel:                req.EducationLevel,
		WorkCurrentEmployer:           req.WorkCurrentEmployer,
		WorkEmployerPhone:             req.WorkEmployerPhone,
		WorkEmployerEmail:             req.WorkEmployerEmail,
		WorkCurrentPosition:           req.WorkCurrentPosition,
		WorkCurrentlyEmployed:         req.WorkCurrentlyEmployed,
		WorkStartDate:                 req.WorkStartDate,
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
		ApplicationDate:               req.ApplicationDate,
		ReferrerSignature:             req.ReferrerSignature,
	}
}

func toUpdateRegistrationFormParams(id uuid.UUID, req updateRegistrationFormRequest) domain.UpdateRegistrationFormParams {
	return domain.UpdateRegistrationFormParams{
		ID:                            id,
		ClientFirstName:               req.ClientFirstName,
		ClientLastName:                req.ClientLastName,
		ClientDateOfBirth:             req.ClientDateOfBirth,
		ClientBsnNumber:               req.ClientBsnNumber,
		ClientGender:                  req.ClientGender,
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
		EducationLevel:                req.EducationLevel,
		WorkCurrentEmployer:           req.WorkCurrentEmployer,
		WorkEmployerPhone:             req.WorkEmployerPhone,
		WorkEmployerEmail:             req.WorkEmployerEmail,
		WorkCurrentPosition:           req.WorkCurrentPosition,
		WorkCurrentlyEmployed:         req.WorkCurrentlyEmployed,
		WorkStartDate:                 req.WorkStartDate,
		WorkAdditionalNotes:           req.WorkAdditionalNotes,
		CareProtectedLiving:           req.CareProtectedLiving,
		CareAssistedIndependentLiving: req.CareAssistedIndependentLiving,
		CareRoomTrainingCenter:        req.CareRoomTrainingCenter,
		CareAmbulatoryGuidance:        req.CareAmbulatoryGuidance,
		ApplicationReason:             req.ApplicationReason,
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
		ApplicationDate:               req.ApplicationDate,
		ReferrerSignature:             req.ReferrerSignature,
	}
}

func toListRegistrationFormsParams(req listRegistrationFormsRequest) domain.ListRegistrationFormsParams {
	page := req.PageRequest.Params()
	return domain.ListRegistrationFormsParams{
		Limit:                  page.Limit,
		Offset:                 page.Offset,
		Status:                 req.Status,
		RiskAggressiveBehavior: req.RiskAggressiveBehavior,
		RiskSuicidalSelfharm:   req.RiskSuicidalSelfharm,
		RiskSubstanceAbuse:     req.RiskSubstanceAbuse,
		RiskPsychiatricIssues:  req.RiskPsychiatricIssues,
		RiskCriminalHistory:    req.RiskCriminalHistory,
		RiskFlightBehavior:     req.RiskFlightBehavior,
		RiskWeaponPossession:   req.RiskWeaponPossession,
		RiskSexualBehavior:     req.RiskSexualBehavior,
		RiskDayNightRhythm:     req.RiskDayNightRhythm,
	}
}

func toUpdateRegistrationFormStatusParams(id uuid.UUID, req updateRegistrationFormStatusRequest, employeeID uuid.UUID) domain.UpdateRegistrationFormStatusParams {
	return domain.UpdateRegistrationFormStatusParams{
		ID:                        id,
		Status:                    req.Status,
		ProcessedByEmployeeID:     &employeeID,
		IntakeAppointmentLocation: req.IntakeAppointmentLocation,
		AddmissionType:            req.AddmissionType,
		RejectionReason:           req.RejectionReason,
	}
}
