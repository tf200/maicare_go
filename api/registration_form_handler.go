package api

import (
	db "maicare_go/db/sqlc"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateRegistrationFormRequest struct {
	ClientFirstName               string    `json:"client_first_name"`
	ClientLastName                string    `json:"client_last_name"`
	ClientBsnNumber               string    `json:"client_bsn_number"`
	ClientGender                  string    `json:"client_gender"`
	ClientNationality             string    `json:"client_nationality"`
	ClientPhoneNumber             string    `json:"client_phone_number"`
	ClientEmail                   string    `json:"client_email"`
	ClientStreet                  string    `json:"client_street"`
	ClientHouseNumber             string    `json:"client_house_number"`
	ClientPostalCode              string    `json:"client_postal_code"`
	ClientCity                    string    `json:"client_city"`
	ReferrerFirstName             string    `json:"referrer_first_name"`
	ReferrerLastName              string    `json:"referrer_last_name"`
	ReferrerOrganization          string    `json:"referrer_organization"`
	ReferrerJobTitle              string    `json:"referrer_job_title"`
	ReferrerPhoneNumber           string    `json:"referrer_phone_number"`
	ReferrerEmail                 string    `json:"referrer_email"`
	Guardian1FirstName            string    `json:"guardian1_first_name"`
	Guardian1LastName             string    `json:"guardian1_last_name"`
	Guardian1Relationship         string    `json:"guardian1_relationship"`
	Guardian1PhoneNumber          string    `json:"guardian1_phone_number"`
	Guardian1Email                string    `json:"guardian1_email"`
	Guardian2FirstName            string    `json:"guardian2_first_name"`
	Guardian2LastName             string    `json:"guardian2_last_name"`
	Guardian2Relationship         string    `json:"guardian2_relationship"`
	Guardian2PhoneNumber          string    `json:"guardian2_phone_number"`
	Guardian2Email                string    `json:"guardian2_email"`
	EducationInstitution          string    `json:"education_institution"`
	EducationMentorName           string    `json:"education_mentor_name"`
	EducationMentorPhone          string    `json:"education_mentor_phone"`
	EducationMentorEmail          string    `json:"education_mentor_email"`
	EducationCurrentlyEnrolled    bool      `json:"education_currently_enrolled"`
	EducationAdditionalNotes      *string   `json:"education_additional_notes"`
	CareProtectedLiving           *bool     `json:"care_protected_living"`
	CareAssistedIndependentLiving *bool     `json:"care_assisted_independent_living"`
	CareRoomTrainingCenter        *bool     `json:"care_room_training_center"`
	CareAmbulatoryGuidance        *bool     `json:"care_ambulatory_guidance"`
	RiskAggressiveBehavior        *bool     `json:"risk_aggressive_behavior"`
	RiskSuicidalSelfharm          *bool     `json:"risk_suicidal_selfharm"`
	RiskSubstanceAbuse            *bool     `json:"risk_substance_abuse"`
	RiskPsychiatricIssues         *bool     `json:"risk_psychiatric_issues"`
	RiskCriminalHistory           *bool     `json:"risk_criminal_history"`
	RiskFlightBehavior            *bool     `json:"risk_flight_behavior"`
	RiskWeaponPossession          *bool     `json:"risk_weapon_possession"`
	RiskSexualBehavior            *bool     `json:"risk_sexual_behavior"`
	RiskDayNightRhythm            *bool     `json:"risk_day_night_rhythm"`
	RiskOther                     *bool     `json:"risk_other"`
	RiskOtherDescription          *string   `json:"risk_other_description"`
	RiskAdditionalNotes           *string   `json:"risk_additional_notes"`
	DocumentReferral              uuid.UUID `json:"document_referral"`
	DocumentEducationReport       uuid.UUID `json:"document_education_report"`
	DocumentPsychiatricReport     uuid.UUID `json:"document_psychiatric_report"`
	DocumentDiagnosis             uuid.UUID `json:"document_diagnosis"`
	DocumentSafetyPlan            uuid.UUID `json:"document_safety_plan"`
	DocumentIDCopy                uuid.UUID `json:"document_id_copy"`
	ApplicationDate               time.Time `json:"application_date"`
	ReferrerSignature             *bool     `json:"referrer_signature"`
}

type CreateRegistrationFormResponse struct {
	ID                            int64     `json:"id"`
	ClientFirstName               string    `json:"client_first_name"`
	ClientLastName                string    `json:"client_last_name"`
	ClientBsnNumber               string    `json:"client_bsn_number"`
	ClientGender                  string    `json:"client_gender"`
	ClientNationality             string    `json:"client_nationality"`
	ClientPhoneNumber             string    `json:"client_phone_number"`
	ClientEmail                   string    `json:"client_email"`
	ClientStreet                  string    `json:"client_street"`
	ClientHouseNumber             string    `json:"client_house_number"`
	ClientPostalCode              string    `json:"client_postal_code"`
	ClientCity                    string    `json:"client_city"`
	ReferrerFirstName             string    `json:"referrer_first_name"`
	ReferrerLastName              string    `json:"referrer_last_name"`
	ReferrerOrganization          string    `json:"referrer_organization"`
	ReferrerJobTitle              string    `json:"referrer_job_title"`
	ReferrerPhoneNumber           string    `json:"referrer_phone_number"`
	ReferrerEmail                 string    `json:"referrer_email"`
	Guardian1FirstName            string    `json:"guardian1_first_name"`
	Guardian1LastName             string    `json:"guardian1_last_name"`
	Guardian1Relationship         string    `json:"guardian1_relationship"`
	Guardian1PhoneNumber          string    `json:"guardian1_phone_number"`
	Guardian1Email                string    `json:"guardian1_email"`
	Guardian2FirstName            string    `json:"guardian2_first_name"`
	Guardian2LastName             string    `json:"guardian2_last_name"`
	Guardian2Relationship         string    `json:"guardian2_relationship"`
	Guardian2PhoneNumber          string    `json:"guardian2_phone_number"`
	Guardian2Email                string    `json:"guardian2_email"`
	EducationInstitution          string    `json:"education_institution"`
	EducationMentorName           string    `json:"education_mentor_name"`
	EducationMentorPhone          string    `json:"education_mentor_phone"`
	EducationMentorEmail          string    `json:"education_mentor_email"`
	EducationCurrentlyEnrolled    bool      `json:"education_currently_enrolled"`
	EducationAdditionalNotes      *string   `json:"education_additional_notes"`
	CareProtectedLiving           *bool     `json:"care_protected_living"`
	CareAssistedIndependentLiving *bool     `json:"care_assisted_independent_living"`
	CareRoomTrainingCenter        *bool     `json:"care_room_training_center"`
	CareAmbulatoryGuidance        *bool     `json:"care_ambulatory_guidance"`
	ApplicationReason             *string   `json:"application_reason"`
	ClientGoals                   *string   `json:"client_goals"`
	RiskAggressiveBehavior        *bool     `json:"risk_aggressive_behavior"`
	RiskSuicidalSelfharm          *bool     `json:"risk_suicidal_selfharm"`
	RiskSubstanceAbuse            *bool     `json:"risk_substance_abuse"`
	RiskPsychiatricIssues         *bool     `json:"risk_psychiatric_issues"`
	RiskCriminalHistory           *bool     `json:"risk_criminal_history"`
	RiskFlightBehavior            *bool     `json:"risk_flight_behavior"`
	RiskWeaponPossession          *bool     `json:"risk_weapon_possession"`
	RiskSexualBehavior            *bool     `json:"risk_sexual_behavior"`
	RiskDayNightRhythm            *bool     `json:"risk_day_night_rhythm"`
	RiskOther                     *bool     `json:"risk_other"`
	RiskOtherDescription          *string   `json:"risk_other_description"`
	RiskAdditionalNotes           *string   `json:"risk_additional_notes"`
	DocumentReferral              uuid.UUID `json:"document_referral"`
	DocumentEducationReport       uuid.UUID `json:"document_education_report"`
	DocumentActionPlan            uuid.UUID `json:"document_action_plan"`
	DocumentPsychiatricReport     uuid.UUID `json:"document_psychiatric_report"`
	DocumentDiagnosis             uuid.UUID `json:"document_diagnosis"`
	DocumentSafetyPlan            uuid.UUID `json:"document_safety_plan"`
	DocumentIDCopy                uuid.UUID `json:"document_id_copy"`
	ApplicationDate               time.Time `json:"application_date"`
	ReferrerSignature             *bool     `json:"referrer_signature"`
	FormStatus                    string    `json:"form_status"`
	CreatedAt                     time.Time `json:"created_at"`
	UpdatedAt                     time.Time `json:"updated_at"`
	SubmittedAt                   time.Time `json:"submitted_at"`
	ProcessedAt                   time.Time `json:"processed_at"`
	ProcessedByEmployeeID         *int64    `json:"processed_by_employee_id"`
}

// @Summary Create Registration Form
// @Description Create a new registration form
// @Tags Registration Form
// @Accept json
// @Produce json
// @Param request body CreateRegistrationFormRequest true "Create Registration Form Request"
// @Success 200 {object} CreateRegistrationFormResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /registration_form [post]
func (server *Server) CreateRegistrationFormApi(ctx *gin.Context) {
	var req CreateRegistrationFormRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		DocumentReferral:              pgtype.UUID{Bytes: req.DocumentReferral, Valid: true},
		DocumentEducationReport:       pgtype.UUID{Bytes: req.DocumentEducationReport, Valid: true},
		DocumentPsychiatricReport:     pgtype.UUID{Bytes: req.DocumentPsychiatricReport, Valid: true},
		DocumentDiagnosis:             pgtype.UUID{Bytes: req.DocumentDiagnosis, Valid: true},
		DocumentSafetyPlan:            pgtype.UUID{Bytes: req.DocumentSafetyPlan, Valid: true},
		DocumentIDCopy:                pgtype.UUID{Bytes: req.DocumentIDCopy, Valid: true},
		ApplicationDate:               pgtype.Date{Time: req.ApplicationDate, Valid: true},
		ReferrerSignature:             req.ReferrerSignature,
	}
	createdForm, err := server.store.CreateRegistrationForm(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		DocumentReferral:              createdForm.DocumentReferral.Bytes,
		DocumentEducationReport:       createdForm.DocumentEducationReport.Bytes,
		DocumentActionPlan:            createdForm.DocumentActionPlan.Bytes,
		DocumentPsychiatricReport:     createdForm.DocumentPsychiatricReport.Bytes,
		DocumentDiagnosis:             createdForm.DocumentDiagnosis.Bytes,
		DocumentSafetyPlan:            createdForm.DocumentSafetyPlan.Bytes,
		DocumentIDCopy:                createdForm.DocumentIDCopy.Bytes,
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
	ctx.JSON(http.StatusOK, res)

}
