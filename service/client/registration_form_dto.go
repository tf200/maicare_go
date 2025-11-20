package clientp

import (
	"time"

	"maicare_go/pagination"

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
	EducationLevel                *string    `json:"education_level"`
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
	ProcessedByEmployeeID         *uuid.UUID `json:"processed_by_employee_id"`
}

// ListRegistrationFormsRequest represents the request body for listing registration forms
type ListRegistrationFormsRequest struct {
	pagination.Request
	Status                 *string `form:"status" json:"status" binding:"omitempty,oneof=pending approved rejected"`
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
	ProcessedByEmployeeID         *uuid.UUID `json:"processed_by_employee_id"`
	IntakeAppointmentDate         time.Time  `json:"intake_appointment_date,omitempty"`
	AddmissionType                *string    `json:"admission_type"` // "crisis_admission" or "regular_placement"
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
	EducationLevel                *string    `json:"education_level"`
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
	ProcessedByEmployeeID         *uuid.UUID `json:"processed_by_employee_id"`
	IntakeAppointmentDate         time.Time  `json:"intake_appointment_date,omitempty"`
	AddmissionType                *string    `json:"admission_type"` // "crisis_admission" or "regular_placement"
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
	EducationLevel                *string     `json:"education_level"`
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
	ProcessedByEmployeeID         *uuid.UUID `json:"processed_by_employee_id"`
}

// UpdateRegistrationFormStatusRequest represents the response body for updating a registration form status
type UpdateRegistrationFormStatusRequest struct {
	Status                    string    `json:"status" binding:"required,oneof=approved rejected" example:"approved"`
	IntakeAppointmentDate     time.Time `json:"intake_appointment_date" binding:"required_if=Status approved" example:"2023-10-01T10:00:00Z"`
	IntakeAppointmentLocation *string   `json:"intake_appointment_location" binding:"required_if=Status approved" example:"Amsterdam Central Station"`
	AddmissionType            *string   `json:"admission_type" binding:"required_if=Status approved,oneof=crisis_admission regular_placement" example:"regular_placement"`
}
