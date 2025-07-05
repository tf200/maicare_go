package api

import (
	"database/sql"
	"fmt"
	"maicare_go/async"
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"maicare_go/util"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateRegistrationFormRequest represents the request body for creating a registration form
type CreateRegistrationFormRequest struct {
	ClientFirstName               string     `json:"client_first_name"`
	ClientLastName                string     `json:"client_last_name"`
	ClientBsnNumber               string     `json:"client_bsn_number"`
	ClientGender                  string     `json:"client_gender"`
	ClientNationality             string     `json:"client_nationality"`
	ClientPhoneNumber             string     `json:"client_phone_number"`
	ClientEmail                   string     `json:"client_email"`
	ClientStreet                  string     `json:"client_street"`
	ClientHouseNumber             string     `json:"client_house_number"`
	ClientPostalCode              string     `json:"client_postal_code"`
	ClientCity                    string     `json:"client_city"`
	ReferrerFirstName             string     `json:"referrer_first_name"`
	ReferrerLastName              string     `json:"referrer_last_name"`
	ReferrerOrganization          string     `json:"referrer_organization"`
	ReferrerJobTitle              string     `json:"referrer_job_title"`
	ReferrerPhoneNumber           string     `json:"referrer_phone_number"`
	ReferrerEmail                 string     `json:"referrer_email"`
	Guardian1FirstName            string     `json:"guardian1_first_name"`
	Guardian1LastName             string     `json:"guardian1_last_name"`
	Guardian1Relationship         string     `json:"guardian1_relationship"`
	Guardian1PhoneNumber          string     `json:"guardian1_phone_number"`
	Guardian1Email                string     `json:"guardian1_email"`
	Guardian2FirstName            string     `json:"guardian2_first_name"`
	Guardian2LastName             string     `json:"guardian2_last_name"`
	Guardian2Relationship         string     `json:"guardian2_relationship"`
	Guardian2PhoneNumber          string     `json:"guardian2_phone_number"`
	Guardian2Email                string     `json:"guardian2_email"`
	EducationInstitution          *string    `json:"education_institution"`
	EducationMentorName           *string    `json:"education_mentor_name"`
	EducationMentorPhone          *string    `json:"education_mentor_phone"`
	EducationMentorEmail          *string    `json:"education_mentor_email"`
	EducationCurrentlyEnrolled    bool       `json:"education_currently_enrolled"`
	EducationAdditionalNotes      *string    `json:"education_additional_notes"`
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
	ApplicationDate               time.Time  `json:"application_date"`
	ReferrerSignature             *bool      `json:"referrer_signature"`
}

// CreateRegistrationFormResponse represents the response body for creating a registration form
type CreateRegistrationFormResponse struct {
	ID                            int64      `json:"id"`
	ClientFirstName               string     `json:"client_first_name"`
	ClientLastName                string     `json:"client_last_name"`
	ClientBsnNumber               string     `json:"client_bsn_number"`
	ClientGender                  string     `json:"client_gender"`
	ClientNationality             string     `json:"client_nationality"`
	ClientPhoneNumber             string     `json:"client_phone_number"`
	ClientEmail                   string     `json:"client_email"`
	ClientStreet                  string     `json:"client_street"`
	ClientHouseNumber             string     `json:"client_house_number"`
	ClientPostalCode              string     `json:"client_postal_code"`
	ClientCity                    string     `json:"client_city"`
	ReferrerFirstName             string     `json:"referrer_first_name"`
	ReferrerLastName              string     `json:"referrer_last_name"`
	ReferrerOrganization          string     `json:"referrer_organization"`
	ReferrerJobTitle              string     `json:"referrer_job_title"`
	ReferrerPhoneNumber           string     `json:"referrer_phone_number"`
	ReferrerEmail                 string     `json:"referrer_email"`
	Guardian1FirstName            string     `json:"guardian1_first_name"`
	Guardian1LastName             string     `json:"guardian1_last_name"`
	Guardian1Relationship         string     `json:"guardian1_relationship"`
	Guardian1PhoneNumber          string     `json:"guardian1_phone_number"`
	Guardian1Email                string     `json:"guardian1_email"`
	Guardian2FirstName            string     `json:"guardian2_first_name"`
	Guardian2LastName             string     `json:"guardian2_last_name"`
	Guardian2Relationship         string     `json:"guardian2_relationship"`
	Guardian2PhoneNumber          string     `json:"guardian2_phone_number"`
	Guardian2Email                string     `json:"guardian2_email"`
	EducationInstitution          *string    `json:"education_institution"`
	EducationMentorName           *string    `json:"education_mentor_name"`
	EducationMentorPhone          *string    `json:"education_mentor_phone"`
	EducationMentorEmail          *string    `json:"education_mentor_email"`
	EducationCurrentlyEnrolled    bool       `json:"education_currently_enrolled"`
	EducationAdditionalNotes      *string    `json:"education_additional_notes"`
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
	ApplicationReason             *string    `json:"application_reason"`
	ClientGoals                   *string    `json:"client_goals"`
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
	DocumentActionPlan            *uuid.UUID `json:"document_action_plan"`
	DocumentPsychiatricReport     *uuid.UUID `json:"document_psychiatric_report"`
	DocumentDiagnosis             *uuid.UUID `json:"document_diagnosis"`
	DocumentSafetyPlan            *uuid.UUID `json:"document_safety_plan"`
	DocumentIDCopy                *uuid.UUID `json:"document_id_copy"`
	ApplicationDate               time.Time  `json:"application_date"`
	ReferrerSignature             *bool      `json:"referrer_signature"`
	FormStatus                    string     `json:"form_status"`
	CreatedAt                     time.Time  `json:"created_at"`
	UpdatedAt                     time.Time  `json:"updated_at"`
	SubmittedAt                   time.Time  `json:"submitted_at"`
	ProcessedAt                   time.Time  `json:"processed_at"`
	ProcessedByEmployeeID         *int64     `json:"processed_by_employee_id"`
}

// @Summary Create Registration Form
// @Description Create a new registration form
// @Tags Registration Form
// @Accept json
// @Produce json
// @Param request body CreateRegistrationFormRequest true "Create Registration Form Request"
// @Success 200 {object} CreateRegistrationFormResponse
// @Failure 400 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /registration_form [post]
func (server *Server) CreateRegistrationFormApi(ctx *gin.Context) {
	var req CreateRegistrationFormRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateRegistrationFormParams{
		ClientFirstName:               req.ClientFirstName,
		ClientLastName:                req.ClientLastName,
		ClientBsnNumber:               req.ClientBsnNumber,
		ClientGender:                  req.ClientGender,
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

	createdForm, err := server.store.CreateRegistrationForm(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := CreateRegistrationFormResponse{
		ID:                            createdForm.ID,
		ClientFirstName:               createdForm.ClientFirstName,
		ClientLastName:                createdForm.ClientLastName,
		ClientBsnNumber:               createdForm.ClientBsnNumber,
		ClientGender:                  createdForm.ClientGender,
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
		FormStatus:                    createdForm.FormStatus,
		CreatedAt:                     createdForm.CreatedAt.Time,
		UpdatedAt:                     createdForm.UpdatedAt.Time,
		SubmittedAt:                   createdForm.SubmittedAt.Time,
		ProcessedAt:                   createdForm.ProcessedAt.Time,
		ProcessedByEmployeeID:         createdForm.ProcessedByEmployeeID,
	}
	res := SuccessResponse(response, "Registration form created successfully")
	ctx.JSON(http.StatusCreated, res)

}

// ListRegistrationFormsRequest represents the request body for listing registration forms
type ListRegistrationFormsRequest struct {
	pagination.Request
	Status                 *string `form:"status" json:"status" binding:"oneof=pending approved rejected all"`
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

// ListRegistrationFormsResponse represents the response body for listing registration forms
type ListRegistrationFormsResponse struct {
	ID                            int64      `json:"id"`
	ClientFirstName               string     `json:"client_first_name"`
	ClientLastName                string     `json:"client_last_name"`
	ClientBsnNumber               string     `json:"client_bsn_number"`
	ClientGender                  string     `json:"client_gender"`
	ClientNationality             string     `json:"client_nationality"`
	ClientPhoneNumber             string     `json:"client_phone_number"`
	ClientEmail                   string     `json:"client_email"`
	ClientStreet                  string     `json:"client_street"`
	ClientHouseNumber             string     `json:"client_house_number"`
	ClientPostalCode              string     `json:"client_postal_code"`
	ClientCity                    string     `json:"client_city"`
	ReferrerFirstName             string     `json:"referrer_first_name"`
	ReferrerLastName              string     `json:"referrer_last_name"`
	ReferrerOrganization          string     `json:"referrer_organization"`
	ReferrerJobTitle              string     `json:"referrer_job_title"`
	ReferrerPhoneNumber           string     `json:"referrer_phone_number"`
	ReferrerEmail                 string     `json:"referrer_email"`
	Guardian1FirstName            string     `json:"guardian1_first_name"`
	Guardian1LastName             string     `json:"guardian1_last_name"`
	Guardian1Relationship         string     `json:"guardian1_relationship"`
	Guardian1PhoneNumber          string     `json:"guardian1_phone_number"`
	Guardian1Email                string     `json:"guardian1_email"`
	Guardian2FirstName            string     `json:"guardian2_first_name"`
	Guardian2LastName             string     `json:"guardian2_last_name"`
	Guardian2Relationship         string     `json:"guardian2_relationship"`
	Guardian2PhoneNumber          string     `json:"guardian2_phone_number"`
	Guardian2Email                string     `json:"guardian2_email"`
	EducationInstitution          *string    `json:"education_institution"`
	EducationMentorName           *string    `json:"education_mentor_name"`
	EducationMentorPhone          *string    `json:"education_mentor_phone"`
	EducationMentorEmail          *string    `json:"education_mentor_email"`
	EducationCurrentlyEnrolled    bool       `json:"education_currently_enrolled"`
	WorkCurrentEmployer           *string    `json:"work_current_employer"`
	WorkEmployerPhone             *string    `json:"work_employer_phone"`
	WorkEmployerEmail             *string    `json:"work_employer_email"`
	WorkCurrentPosition           *string    `json:"work_current_position"`
	WorkCurrentlyEmployed         bool       `json:"work_currently_employed"`
	WorkStartDate                 *time.Time `json:"work_start_date"`
	WorkAdditionalNotes           *string    `json:"work_additional_notes"`
	EducationAdditionalNotes      *string    `json:"education_additional_notes"`
	CareProtectedLiving           *bool      `json:"care_protected_living"`
	CareAssistedIndependentLiving *bool      `json:"care_assisted_independent_living"`
	CareRoomTrainingCenter        *bool      `json:"care_room_training_center"`
	CareAmbulatoryGuidance        *bool      `json:"care_ambulatory_guidance"`
	ApplicationReason             *string    `json:"application_reason"`
	ClientGoals                   *string    `json:"client_goals"`
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
	RiskCount                     int        `json:"risk_count"`
	DocumentReferral              *uuid.UUID `json:"document_referral"`
	DocumentEducationReport       *uuid.UUID `json:"document_education_report"`
	DocumentActionPlan            *uuid.UUID `json:"document_action_plan"`
	DocumentPsychiatricReport     *uuid.UUID `json:"document_psychiatric_report"`
	DocumentDiagnosis             *uuid.UUID `json:"document_diagnosis"`
	DocumentSafetyPlan            *uuid.UUID `json:"document_safety_plan"`
	DocumentIDCopy                *uuid.UUID `json:"document_id_copy"`
	ApplicationDate               time.Time  `json:"application_date"`
	ReferrerSignature             *bool      `json:"referrer_signature"`
	FormStatus                    string     `json:"form_status"`
	CreatedAt                     time.Time  `json:"created_at"`
	UpdatedAt                     time.Time  `json:"updated_at"`
	SubmittedAt                   time.Time  `json:"submitted_at"`
	ProcessedAt                   time.Time  `json:"processed_at"`
	ProcessedByEmployeeID         *int64     `json:"processed_by_employee_id"`
	IntakeAppointmentDate         time.Time  `json:"intake_appointment_date,omitempty"`
	AddmissionType                *string    `json:"admission_type"` // "crisis_admission" or "regular_placement"
}

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

// @Summary List Registration Forms
// @Description List all registration forms
// @Tags Registration Form
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param status query string false "Form status" Enums(pending, approved, rejected)
// @Param risk_aggressive_behavior query bool false "Risk aggressive behavior"
// @Param risk_suicidal_selfharm query bool false "Risk suicidal self-harm"
// @Param risk_substance_abuse query bool false "Risk substance abuse"
// @Param risk_psychiatric_issues query bool false "Risk psychiatric issues"
// @Param risk_criminal_history query bool false "Risk criminal history"
// @Param risk_flight_behavior query bool false "Risk flight behavior"
// @Param risk_weapon_possession query bool false "Risk weapon possession"
// @Param risk_sexual_behavior query bool false "Risk sexual behavior"
// @Param risk_day_night_rhythm query bool false "Risk day-night rhythm"
// @Success 200 {object} Response[pagination.Response[ListRegistrationFormsResponse]]
// @Failure 400 {object}  Response[any]
// @Failure 500 {object}  Response[any]
// @Router /registration_form [get]
func (server *Server) ListRegistrationFormsApi(ctx *gin.Context) {
	var req ListRegistrationFormsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := req.GetParams()

	arg := db.ListRegistrationFormsParams{
		Limit:                  params.Limit,
		Offset:                 params.Offset,
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
	registrationForms, err := server.store.ListRegistrationForms(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if len(registrationForms) == 0 {
		ctx.JSON(http.StatusOK, SuccessResponse([]ListRegistrationFormsResponse{}, "No registration forms found"))
		return
	}

	contArg := db.CountRegistrationFormsParams{
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
	totalCount, err := server.store.CountRegistrationForms(ctx, contArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseRf := make([]ListRegistrationFormsResponse, len(registrationForms))
	for i, rf := range registrationForms {
		responseRf[i] = ListRegistrationFormsResponse{
			ID:                            rf.ID,
			ClientFirstName:               rf.ClientFirstName,
			ClientLastName:                rf.ClientLastName,
			ClientBsnNumber:               rf.ClientBsnNumber,
			ClientGender:                  rf.ClientGender,
			ClientNationality:             rf.ClientNationality,
			ClientPhoneNumber:             rf.ClientPhoneNumber,
			ClientEmail:                   rf.ClientEmail,
			ClientStreet:                  rf.ClientStreet,
			ClientHouseNumber:             rf.ClientHouseNumber,
			ClientPostalCode:              rf.ClientPostalCode,
			ClientCity:                    rf.ClientCity,
			ReferrerFirstName:             rf.ReferrerFirstName,
			ReferrerLastName:              rf.ReferrerLastName,
			ReferrerOrganization:          rf.ReferrerOrganization,
			ReferrerJobTitle:              rf.ReferrerJobTitle,
			ReferrerPhoneNumber:           rf.ReferrerPhoneNumber,
			ReferrerEmail:                 rf.ReferrerEmail,
			Guardian1FirstName:            rf.Guardian1FirstName,
			Guardian1LastName:             rf.Guardian1LastName,
			Guardian1Relationship:         rf.Guardian1Relationship,
			Guardian1PhoneNumber:          rf.Guardian1PhoneNumber,
			Guardian1Email:                rf.Guardian1Email,
			Guardian2FirstName:            rf.Guardian2FirstName,
			Guardian2LastName:             rf.Guardian2LastName,
			Guardian2Relationship:         rf.Guardian2Relationship,
			Guardian2PhoneNumber:          rf.Guardian2PhoneNumber,
			Guardian2Email:                rf.Guardian2Email,
			EducationInstitution:          rf.EducationInstitution,
			EducationMentorName:           rf.EducationMentorName,
			EducationMentorPhone:          rf.EducationMentorPhone,
			EducationMentorEmail:          rf.EducationMentorEmail,
			EducationCurrentlyEnrolled:    rf.EducationCurrentlyEnrolled,
			EducationAdditionalNotes:      rf.EducationAdditionalNotes,
			WorkCurrentEmployer:           rf.WorkCurrentEmployer,
			WorkEmployerPhone:             rf.WorkEmployerPhone,
			WorkEmployerEmail:             rf.WorkEmployerEmail,
			WorkCurrentPosition:           rf.WorkCurrentPosition,
			WorkCurrentlyEmployed:         rf.WorkCurrentlyEmployed,
			WorkStartDate:                 &rf.WorkStartDate.Time,
			WorkAdditionalNotes:           rf.WorkAdditionalNotes,
			CareProtectedLiving:           rf.CareProtectedLiving,
			CareAssistedIndependentLiving: rf.CareAssistedIndependentLiving,
			CareRoomTrainingCenter:        rf.CareRoomTrainingCenter,
			CareAmbulatoryGuidance:        rf.CareAmbulatoryGuidance,
			ApplicationReason:             rf.ApplicationReason,
			ClientGoals:                   rf.ClientGoals,
			RiskAggressiveBehavior:        rf.RiskAggressiveBehavior,
			RiskSuicidalSelfharm:          rf.RiskSuicidalSelfharm,
			RiskSubstanceAbuse:            rf.RiskSubstanceAbuse,
			RiskPsychiatricIssues:         rf.RiskPsychiatricIssues,
			RiskCriminalHistory:           rf.RiskCriminalHistory,
			RiskFlightBehavior:            rf.RiskFlightBehavior,
			RiskWeaponPossession:          rf.RiskWeaponPossession,
			RiskSexualBehavior:            rf.RiskSexualBehavior,
			RiskDayNightRhythm:            rf.RiskDayNightRhythm,
			RiskOther:                     rf.RiskOther,
			RiskOtherDescription:          rf.RiskOtherDescription,
			RiskAdditionalNotes:           rf.RiskAdditionalNotes,
			DocumentReferral:              rf.DocumentReferral,
			DocumentEducationReport:       rf.DocumentEducationReport,
			DocumentActionPlan:            rf.DocumentActionPlan,
			DocumentPsychiatricReport:     rf.DocumentPsychiatricReport,
			DocumentDiagnosis:             rf.DocumentDiagnosis,
			DocumentSafetyPlan:            rf.DocumentSafetyPlan,
			DocumentIDCopy:                rf.DocumentIDCopy,
			ApplicationDate:               rf.ApplicationDate.Time,
			ReferrerSignature:             rf.ReferrerSignature,
			FormStatus:                    rf.FormStatus,
			CreatedAt:                     rf.CreatedAt.Time,
			UpdatedAt:                     rf.UpdatedAt.Time,
			SubmittedAt:                   rf.SubmittedAt.Time,
			ProcessedAt:                   rf.ProcessedAt.Time,
			ProcessedByEmployeeID:         rf.ProcessedByEmployeeID,
			RiskCount:                     calculateRiskCount(rf),
			IntakeAppointmentDate:         rf.IntakeAppointmentDatetime.Time,
			AddmissionType:                rf.AddmissionType,
		}
	}

	response := pagination.NewResponse(ctx, req.Request, responseRf, totalCount)
	res := SuccessResponse(response, "Registration forms retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

type GetRegistrationFormResponse struct {
	ID                            int64      `json:"id"`
	ClientFirstName               string     `json:"client_first_name"`
	ClientLastName                string     `json:"client_last_name"`
	ClientBsnNumber               string     `json:"client_bsn_number"`
	ClientGender                  string     `json:"client_gender"`
	ClientNationality             string     `json:"client_nationality"`
	ClientPhoneNumber             string     `json:"client_phone_number"`
	ClientEmail                   string     `json:"client_email"`
	ClientStreet                  string     `json:"client_street"`
	ClientHouseNumber             string     `json:"client_house_number"`
	ClientPostalCode              string     `json:"client_postal_code"`
	ClientCity                    string     `json:"client_city"`
	ReferrerFirstName             string     `json:"referrer_first_name"`
	ReferrerLastName              string     `json:"referrer_last_name"`
	ReferrerOrganization          string     `json:"referrer_organization"`
	ReferrerJobTitle              string     `json:"referrer_job_title"`
	ReferrerPhoneNumber           string     `json:"referrer_phone_number"`
	ReferrerEmail                 string     `json:"referrer_email"`
	Guardian1FirstName            string     `json:"guardian1_first_name"`
	Guardian1LastName             string     `json:"guardian1_last_name"`
	Guardian1Relationship         string     `json:"guardian1_relationship"`
	Guardian1PhoneNumber          string     `json:"guardian1_phone_number"`
	Guardian1Email                string     `json:"guardian1_email"`
	Guardian2FirstName            string     `json:"guardian2_first_name"`
	Guardian2LastName             string     `json:"guardian2_last_name"`
	Guardian2Relationship         string     `json:"guardian2_relationship"`
	Guardian2PhoneNumber          string     `json:"guardian2_phone_number"`
	Guardian2Email                string     `json:"guardian2_email"`
	EducationInstitution          *string    `json:"education_institution"`
	EducationMentorName           *string    `json:"education_mentor_name"`
	EducationMentorPhone          *string    `json:"education_mentor_phone"`
	EducationMentorEmail          *string    `json:"education_mentor_email"`
	EducationCurrentlyEnrolled    bool       `json:"education_currently_enrolled"`
	EducationAdditionalNotes      string     `json:"education_additional_notes"`
	WorkCurrentEmployer           *string    `json:"work_current_employer"`
	WorkEmployerPhone             *string    `json:"work_employer_phone"`
	WorkEmployerEmail             *string    `json:"work_employer_email"`
	WorkCurrentPosition           *string    `json:"work_current_position"`
	WorkCurrentlyEmployed         bool       `json:"work_currently_employed"`
	WorkStartDate                 *time.Time `json:"work_start_date"`
	WorkAdditionalNotes           string     `json:"work_additional_notes"`
	CareProtectedLiving           *bool      `json:"care_protected_living"`
	CareAssistedIndependentLiving *bool      `json:"care_assisted_independent_living"`
	CareRoomTrainingCenter        *bool      `json:"care_room_training_center"`
	CareAmbulatoryGuidance        *bool      `json:"care_ambulatory_guidance"`
	ApplicationReason             string     `json:"application_reason"`
	ClientGoals                   string     `json:"client_goals"`
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
	RiskOtherDescription          string     `json:"risk_other_description"`
	RiskAdditionalNotes           string     `json:"risk_additional_notes"`
	DocumentReferral              *uuid.UUID `json:"document_referral"`
	DocumentEducationReport       *uuid.UUID `json:"document_education_report"`
	DocumentActionPlan            *uuid.UUID `json:"document_action_plan"`
	DocumentPsychiatricReport     *uuid.UUID `json:"document_psychiatric_report"`
	DocumentDiagnosis             *uuid.UUID `json:"document_diagnosis"`
	DocumentSafetyPlan            *uuid.UUID `json:"document_safety_plan"`
	DocumentIDCopy                *uuid.UUID `json:"document_id_copy"`
	ApplicationDate               time.Time  `json:"application_date"`
	ReferrerSignature             *bool      `json:"referrer_signature"`
	FormStatus                    string     `json:"form_status"`
	CreatedAt                     time.Time  `json:"created_at"`
	UpdatedAt                     time.Time  `json:"updated_at"`
	SubmittedAt                   time.Time  `json:"submitted_at"`
	ProcessedAt                   time.Time  `json:"processed_at"`
	ProcessedByEmployeeID         *int64     `json:"processed_by_employee_id"`
	IntakeAppointmentDate         time.Time  `json:"intake_appointment_date,omitempty"`
	AddmissionType                *string    `json:"admission_type"` // "crisis_admission" or "regular_placement"
}

// @Summary Get Registration Form
// @Description Get a registration form by ID
// @Tags Registration Form
// @Produce json
// @Param id path int true "Registration Form ID"
// @Success 200 {object} Response[GetRegistrationFormResponse]
// @Failure 400 {object} Response[any]
// @Failure 404 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /registration_form/{id} [get]
func (server *Server) GetRegistrationFormApi(ctx *gin.Context) {
	rfId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	registrationForm, err := server.store.GetRegistrationForm(ctx, rfId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := GetRegistrationFormResponse{
		ID:                            registrationForm.ID,
		ClientFirstName:               registrationForm.ClientFirstName,
		ClientLastName:                registrationForm.ClientLastName,
		ClientBsnNumber:               registrationForm.ClientBsnNumber,
		ClientGender:                  registrationForm.ClientGender,
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
		FormStatus:                    registrationForm.FormStatus,
		CreatedAt:                     registrationForm.CreatedAt.Time,
		UpdatedAt:                     registrationForm.UpdatedAt.Time,
		SubmittedAt:                   registrationForm.SubmittedAt.Time,
		ProcessedAt:                   registrationForm.ProcessedAt.Time,
		ProcessedByEmployeeID:         registrationForm.ProcessedByEmployeeID,
		IntakeAppointmentDate:         registrationForm.IntakeAppointmentDatetime.Time,
		AddmissionType:                registrationForm.AddmissionType,
	}
	res := SuccessResponse(response, "Registration form retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateRegistrationFormRequest represents the request body for updating a registration form
type UpdateRegistrationFormRequest struct {
	ClientFirstName               *string     `json:"client_first_name"`
	ClientLastName                *string     `json:"client_last_name"`
	ClientBsnNumber               *string     `json:"client_bsn_number"`
	ClientGender                  *string     `json:"client_gender"`
	ClientNationality             *string     `json:"client_nationality"`
	ClientPhoneNumber             *string     `json:"client_phone_number"`
	ClientEmail                   *string     `json:"client_email"`
	ClientStreet                  *string     `json:"client_street"`
	ClientHouseNumber             *string     `json:"client_house_number"`
	ClientPostalCode              *string     `json:"client_postal_code"`
	ClientCity                    *string     `json:"client_city"`
	ReferrerFirstName             *string     `json:"referrer_first_name"`
	ReferrerLastName              *string     `json:"referrer_last_name"`
	ReferrerOrganization          *string     `json:"referrer_organization"`
	ReferrerJobTitle              *string     `json:"referrer_job_title"`
	ReferrerPhoneNumber           *string     `json:"referrer_phone_number"`
	ReferrerEmail                 *string     `json:"referrer_email"`
	Guardian1FirstName            *string     `json:"guardian1_first_name"`
	Guardian1LastName             *string     `json:"guardian1_last_name"`
	Guardian1Relationship         *string     `json:"guardian1_relationship"`
	Guardian1PhoneNumber          *string     `json:"guardian1_phone_number"`
	Guardian1Email                *string     `json:"guardian1_email"`
	Guardian2FirstName            *string     `json:"guardian2_first_name"`
	Guardian2LastName             *string     `json:"guardian2_last_name"`
	Guardian2Relationship         *string     `json:"guardian2_relationship"`
	Guardian2PhoneNumber          *string     `json:"guardian2_phone_number"`
	Guardian2Email                *string     `json:"guardian2_email"`
	EducationInstitution          *string     `json:"education_institution"`
	EducationMentorName           *string     `json:"education_mentor_name"`
	EducationMentorPhone          *string     `json:"education_mentor_phone"`
	EducationMentorEmail          *string     `json:"education_mentor_email"`
	EducationCurrentlyEnrolled    *bool       `json:"education_currently_enrolled"`
	EducationAdditionalNotes      *string     `json:"education_additional_notes"`
	WorkCurrentEmployer           *string     `json:"work_current_employer"`
	WorkEmployerPhone             *string     `json:"work_employer_phone"`
	WorkEmployerEmail             *string     `json:"work_employer_email"`
	WorkCurrentPosition           *string     `json:"work_current_position"`
	WorkCurrentlyEmployed         *bool       `json:"work_currently_employed"`
	WorkStartDate                 *time.Time  `json:"work_start_date"`
	WorkAdditionalNotes           *string     `json:"work_additional_notes"`
	CareProtectedLiving           *bool       `json:"care_protected_living"`
	CareAssistedIndependentLiving *bool       `json:"care_assisted_independent_living"`
	CareRoomTrainingCenter        *bool       `json:"care_room_training_center"`
	CareAmbulatoryGuidance        *bool       `json:"care_ambulatory_guidance"`
	RiskAggressiveBehavior        *bool       `json:"risk_aggressive_behavior"`
	RiskSuicidalSelfharm          *bool       `json:"risk_suicidal_selfharm"`
	RiskSubstanceAbuse            *bool       `json:"risk_substance_abuse"`
	RiskPsychiatricIssues         *bool       `json:"risk_psychiatric_issues"`
	RiskCriminalHistory           *bool       `json:"risk_criminal_history"`
	RiskFlightBehavior            *bool       `json:"risk_flight_behavior"`
	RiskWeaponPossession          *bool       `json:"risk_weapon_possession"`
	RiskSexualBehavior            *bool       `json:"risk_sexual_behavior"`
	RiskDayNightRhythm            *bool       `json:"risk_day_night_rhythm"`
	RiskOther                     *bool       `json:"risk_other"`
	RiskOtherDescription          *string     `json:"risk_other_description"`
	RiskAdditionalNotes           *string     `json:"risk_additional_notes"`
	DocumentReferral              *uuid.UUID  `json:"document_referral"`
	DocumentEducationReport       *uuid.UUID  `json:"document_education_report"`
	DocumentPsychiatricReport     *uuid.UUID  `json:"document_psychiatric_report"`
	DocumentDiagnosis             *uuid.UUID  `json:"document_diagnosis"`
	DocumentSafetyPlan            *uuid.UUID  `json:"document_safety_plan"`
	DocumentIDCopy                *uuid.UUID  `json:"document_id_copy"`
	ApplicationDate               pgtype.Date `json:"application_date"`
	ReferrerSignature             *bool       `json:"referrer_signature"`
}

// UpdateRegistrationFormResponse represents the response body for updating a registration form
type UpdateRegistrationFormResponse struct {
	ID                            int64      `json:"id"`
	ClientFirstName               string     `json:"client_first_name"`
	ClientLastName                string     `json:"client_last_name"`
	ClientBsnNumber               string     `json:"client_bsn_number"`
	ClientGender                  string     `json:"client_gender"`
	ClientNationality             string     `json:"client_nationality"`
	ClientPhoneNumber             string     `json:"client_phone_number"`
	ClientEmail                   string     `json:"client_email"`
	ClientStreet                  string     `json:"client_street"`
	ClientHouseNumber             string     `json:"client_house_number"`
	ClientPostalCode              string     `json:"client_postal_code"`
	ClientCity                    string     `json:"client_city"`
	ReferrerFirstName             string     `json:"referrer_first_name"`
	ReferrerLastName              string     `json:"referrer_last_name"`
	ReferrerOrganization          string     `json:"referrer_organization"`
	ReferrerJobTitle              string     `json:"referrer_job_title"`
	ReferrerPhoneNumber           string     `json:"referrer_phone_number"`
	ReferrerEmail                 string     `json:"referrer_email"`
	Guardian1FirstName            string     `json:"guardian1_first_name"`
	Guardian1LastName             string     `json:"guardian1_last_name"`
	Guardian1Relationship         string     `json:"guardian1_relationship"`
	Guardian1PhoneNumber          string     `json:"guardian1_phone_number"`
	Guardian1Email                string     `json:"guardian1_email"`
	Guardian2FirstName            string     `json:"guardian2_first_name"`
	Guardian2LastName             string     `json:"guardian2_last_name"`
	Guardian2Relationship         string     `json:"guardian2_relationship"`
	Guardian2PhoneNumber          string     `json:"guardian2_phone_number"`
	Guardian2Email                string     `json:"guardian2_email"`
	EducationInstitution          *string    `json:"education_institution"`
	EducationMentorName           *string    `json:"education_mentor_name"`
	EducationMentorPhone          *string    `json:"education_mentor_phone"`
	EducationMentorEmail          *string    `json:"education_mentor_email"`
	EducationCurrentlyEnrolled    bool       `json:"education_currently_enrolled"`
	EducationAdditionalNotes      *string    `json:"education_additional_notes"`
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
	ApplicationReason             *string    `json:"application_reason"`
	ClientGoals                   *string    `json:"client_goals"`
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
	DocumentActionPlan            *uuid.UUID `json:"document_action_plan"`
	DocumentPsychiatricReport     *uuid.UUID `json:"document_psychiatric_report"`
	DocumentDiagnosis             *uuid.UUID `json:"document_diagnosis"`
	DocumentSafetyPlan            *uuid.UUID `json:"document_safety_plan"`
	DocumentIDCopy                *uuid.UUID `json:"document_id_copy"`
	ApplicationDate               time.Time  `json:"application_date"`
	ReferrerSignature             *bool      `json:"referrer_signature"`
	FormStatus                    string     `json:"form_status"`
	CreatedAt                     time.Time  `json:"created_at"`
	UpdatedAt                     time.Time  `json:"updated_at"`
	SubmittedAt                   time.Time  `json:"submitted_at"`
	ProcessedAt                   time.Time  `json:"processed_at"`
	ProcessedByEmployeeID         *int64     `json:"processed_by_employee_id"`
}

// @Summary Update Registration Form
// @Description Update a registration form by ID
// @Tags Registration Form
// @Accept json
// @Produce json
// @Param id path int true "Registration Form ID"
// @Param request body UpdateRegistrationFormRequest true "Update Registration Form Request"
// @Success 200 {object} Response[UpdateRegistrationFormResponse]
// @Failure 400 {object} Response[any]
// @Failure 404 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /registration_form/{id} [put]
func (server *Server) UpdateRegistrationFormApi(ctx *gin.Context) {
	rfId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req UpdateRegistrationFormRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateRegistrationFormParams{
		ID:                         rfId,
		ClientFirstName:            req.ClientFirstName,
		ClientLastName:             req.ClientLastName,
		ClientBsnNumber:            req.ClientBsnNumber,
		ClientGender:               req.ClientGender,
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
	registrationForm, err := server.store.UpdateRegistrationForm(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := UpdateRegistrationFormResponse{
		ID:                            registrationForm.ID,
		ClientFirstName:               registrationForm.ClientFirstName,
		ClientLastName:                registrationForm.ClientLastName,
		ClientBsnNumber:               registrationForm.ClientBsnNumber,
		ClientGender:                  registrationForm.ClientGender,
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
		FormStatus:                    registrationForm.FormStatus,
		CreatedAt:                     registrationForm.CreatedAt.Time,
		UpdatedAt:                     registrationForm.UpdatedAt.Time,
		SubmittedAt:                   registrationForm.SubmittedAt.Time,
		ProcessedAt:                   registrationForm.ProcessedAt.Time,
		ProcessedByEmployeeID:         registrationForm.ProcessedByEmployeeID,
	}
	res := SuccessResponse(response, "Registration form updated successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Delete Registration Form
// @Description Delete a registration form by ID
// @Tags Registration Form
// @Produce json
// @Param id path int true "Registration Form ID"
// @Success 200 {object} Response[any]
// @Failure 400 {object} Response[any]
// @Failure 404 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /registration_form/{id} [delete]
func (server *Server) DeleteRegistrationFormApi(ctx *gin.Context) {
	rfId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.DeleteRegistrationForm(ctx, rfId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := SuccessResponse[any](nil, "Registration form deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateRegistrationFormStatusRequest represents the response body for updating a registration form status
type UpdateRegistrationFormStatusRequest struct {
	Status                    string    `json:"status" binding:"required,oneof=approved rejected" example:"approved"`
	IntakeAppointmentDate     time.Time `json:"intake_appointment_date" binding:"required_if=Status approved" example:"2023-10-01T10:00:00Z"`
	IntakeAppointmentLocation *string   `json:"intake_appointment_location" binding:"required_if=Status approved" example:"Amsterdam Central Station"`
	AddmissionType            *string   `json:"admission_type" binding:"required_if=Status approved,oneof=crisis_admission regular_placement" example:"regular_placement"`
}

// @Summary Update Registration Form Status
// @Description Update the status of a registration form by ID
// @Tags Registration Form
// @Produce json
// @Param id path int true "Registration Form ID"
// @Param request body UpdateRegistrationFormStatusRequest true "Update Registration Form Status Request"
// @Success 200 {object} Response[any]
// @Failure 400 {object} Response[any]
// @Failure 404 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /registration_form/{id}/status [post]
func (server *Server) UpdateRegistrationFormStatusApi(ctx *gin.Context) {
	var req UpdateRegistrationFormStatusRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	employeeID, err := server.store.GetEmployeeIDByUserID(ctx, payload.UserId)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	rfId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateRegistrationFormStatusParams{
		ID:                        rfId,
		FormStatus:                req.Status,
		ProcessedByEmployeeID:     &employeeID,
		IntakeAppointmentLocation: req.IntakeAppointmentLocation,
		AddmissionType:            req.AddmissionType,
	}

	registrationForm, err := server.store.UpdateRegistrationFormStatus(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if req.Status == "approved" {
		// Enqueue the task to create an intake appointment
		err = server.asynqClient.EnqueueAcceptedRegistration(ctx, async.AcceptedRegistrationFormPayload{
			ReferrerName:        registrationForm.ReferrerFirstName + " " + registrationForm.ReferrerLastName,
			ChildName:           registrationForm.ClientFirstName + " " + registrationForm.ClientLastName,
			ChildBSN:            registrationForm.ClientBsnNumber,
			AppointmentDate:     registrationForm.IntakeAppointmentDatetime.Time.Format("2006-01-02 15:04:05"),
			AppointmentLocation: *registrationForm.IntakeAppointmentLocation,
			To:                  registrationForm.ReferrerEmail,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to enqueue intake appointment creation task: %w", err)))
			return
		}

		res := SuccessResponse[any](nil, "Registration Form Status updated")
		ctx.JSON(http.StatusOK, res)

	}
}
