package clientp

import (
	"time"

	"maicare_go/pagination"

	"github.com/google/uuid"
)

// Address represents a client address - matches location table structure
type Address struct {
	BelongsTo           *string `json:"belongs_to"`
	Street              *string `json:"street"`
	HouseNumber         *string `json:"house_number"`
	HouseNumberAddition *string `json:"house_number_addition"`
	PostalCode          *string `json:"postal_code"`
	City                *string `json:"city"`
	PhoneNumber         *string `json:"phone_number"`
}

// CreateClientDetailsRequest represents a request to create a new client
type CreateClientDetailsRequest struct {
	FirstName                  string     `json:"first_name" binding:"required"`
	LastName                   string     `json:"last_name" binding:"required"`
	Email                      string     `json:"email" binding:"required,email"`
	OrganizationID             *uuid.UUID `json:"organization_id"`
	LocationID                 *uuid.UUID `json:"location_id"`
	LegalMeasure               *string    `json:"legal_measure"`
	Birthplace                 *string    `json:"birthplace"`
	Departement                *string    `json:"departement"`
	Gender                     string     `json:"gender" binding:"required,oneof=male female other"`
	CareType                   *string    `json:"care_type" binding:"omitempty,oneof=protected_living training_center supported_independent_living ambulatory_support other"`
	Filenumber                 string     `json:"filenumber"`
	DateOfBirth                string     `json:"date_of_birth" binding:"required" time_format:"2006-01-02"`
	PhoneNumber                *string    `json:"phone_number" binding:"required"`
	SenderID                   *uuid.UUID `json:"sender_id" binding:"required"`
	Infix                      *string    `json:"infix"`
	Source                     *string    `json:"source" binding:"required"`
	Nationality                *string    `json:"nationality"`
	Bsn                        *string    `json:"bsn"`
	BsnVerifiedBy              *uuid.UUID `json:"bsn_verified_by"` // needs to be checked
	Addresses                  []Address  `json:"addresses"`
	EducationCurrentlyEnrolled bool       `json:"education_currently_enrolled"`
	EducationInstitution       *string    `json:"education_institution"`
	EducationMentorName        *string    `json:"education_mentor_name"`
	EducationMentorPhone       *string    `json:"education_mentor_phone"`
	EducationMentorEmail       *string    `json:"education_mentor_email"`
	EducationAdditionalNotes   *string    `json:"education_additional_notes"`
	EducationLevel             *string    `json:"education_level" binding:"oneof=primary secondary higher none"`
	WorkCurrentlyEmployed      bool       `json:"work_currently_employed"`
	WorkCurrentEmployer        *string    `json:"work_current_employer"`
	WorkCurrentEmployerPhone   *string    `json:"work_employer_phone"`
	WorkCurrentEmployerEmail   *string    `json:"work_employer_email"`
	WorkCurrentPosition        *string    `json:"work_current_position"`
	WorkStartDate              time.Time  `json:"work_start_date"`
	WorkAdditionalNotes        *string    `json:"work_additional_notes"`
	LivingSituation            *string    `json:"living_situation" binding:"oneof=home foster_care youth_care_institution other"`
	LivingSituationNotes       *string    `json:"living_situation_notes"`
}

// CreateClientDetailsResponse represents a response to a create client request
type CreateClientDetailsResponse struct {
	ID                         uuid.UUID  `json:"id"`
	FirstName                  string     `json:"first_name"`
	LastName                   string     `json:"last_name"`
	DateOfBirth                time.Time  `json:"date_of_birth"`
	Identity                   bool       `json:"identity"`
	Status                     string     `json:"status"`
	Bsn                        *string    `json:"bsn"`
	BsnVerifiedBy              *uuid.UUID `json:"bsn_verified_by"` // needs to be checked
	Source                     *string    `json:"source"`
	Birthplace                 *string    `json:"birthplace"`
	Nationality                *string    `json:"nationality"`
	Email                      string     `json:"email"`
	PhoneNumber                *string    `json:"phone_number"`
	OrganizationID             *uuid.UUID `json:"organization_id"`
	Departement                *string    `json:"departement"`
	Gender                     string     `json:"gender"`
	Filenumber                 string     `json:"filenumber"`
	ProfilePicture             *string    `json:"profile_picture"`
	Infix                      *string    `json:"infix"`
	Created                    time.Time  `json:"created"`
	SenderID                   *uuid.UUID `json:"sender_id"`
	LocationID                 *uuid.UUID `json:"location_id"`
	DepartureReason            *string    `json:"departure_reason"`
	DepartureReport            *string    `json:"departure_report"`
	Addresses                  []Address  `json:"addresses"`
	LegalMeasure               *string    `json:"legal_measure"`
	EducationCurrentlyEnrolled bool       `json:"education_currently_enrolled"`
	EducationInstitution       *string    `json:"education_institution"`
	EducationMentorName        *string    `json:"education_mentor_name"`
	EducationMentorEmail       *string    `json:"education_mentor_email"`
	EducationMentorPhone       *string    `json:"education_mentor_phone"`
	EducationAdditionalNotes   *string    `json:"education_additional_notes"`
	EducationLevel             string     `json:"education_level"`
	WorkCurrentlyEmployed      bool       `json:"work_currently_employed"`
	WorkCurrentEmployer        *string    `json:"work_current_employer"`
	WorkCurrentEmployerPhone   *string    `json:"work_employer_phone"`
	WorkCurrentEmployerEmail   *string    `json:"work_employer_email"`
	WorkCurrentPosition        *string    `json:"work_current_position"`
	WorkStartDate              time.Time  `json:"work_start_date"`
	WorkAdditionalNotes        *string    `json:"work_additional_notes"`
	LivingSituation            *string    `json:"living_situation"`
	LivingSituationNotes       *string    `json:"living_situation_notes"`
}

// ListClientsApiParams represents a request to list clients
type ListClientsApiParams struct {
	pagination.Request
	Status     *string    `form:"status"`
	LocationID *uuid.UUID `form:"location_id"`
	Search     *string    `form:"search"`
}

// ListClientsApiResponse represents a response to a list clients request
type ListClientsApiResponse struct {
	ID           uuid.UUID `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Bsn          *string   `json:"bsn"`
	Filenumber   string    `json:"filenumber"`
	LocationName *string   `json:"location_name"`
	CareType     *string   `json:"care_type"`
	Status       string    `json:"status"`
	GoalsCount   int64     `json:"goals_count"`
	RiskCount    int64     `json:"risk_count"`
	CreatedAt    time.Time `json:"created_at"`
}

type ListWaitingListClientsParams struct {
	pagination.Request
	Search    *string `form:"search" binding:"omitempty,max=120"`
	Placement *string `form:"placement" binding:"omitempty,oneof=protected_living training_center supported_independent_living ambulatory_support other"`
	SortDays  *string `form:"sort_days" binding:"omitempty,oneof=asc desc"`
}

type ListWaitingListClientsResponse struct {
	ID             uuid.UUID `json:"id"`
	FirstName      string    `json:"first_name"`
	Bsn            *string   `json:"bsn"`
	LastName       string    `json:"last_name"`
	CareType       *string   `json:"care_type"`
	SenderName     *string   `json:"sender_name"`
	DaysInWaitlist int32     `json:"days_in_waitlist"`
	AdmissionType  *string   `json:"admission_type"`
}

type ListInCareClientsParams struct {
	pagination.Request
	Search         *string  `form:"search" binding:"omitempty,max=120"`
	Status         []string `form:"status" binding:"omitempty,dive,oneof=in_care scheduled_in_care"`
	SortDaysInCare *string  `form:"sort_days_in_care" binding:"omitempty,oneof=asc desc"`
}

type ListInCareClientsResponse struct {
	ID                uuid.UUID  `json:"id"`
	Bsn               *string    `json:"bsn"`
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name"`
	CoordinatorName   *string    `json:"coordinator_name"`
	LocationName      *string    `json:"location_name"`
	Status            string     `json:"status"`
	CareStartDate     *time.Time `json:"care_start_date"`
	DaysInCare        int32      `json:"days_in_care"`
	HasActiveContract bool       `json:"has_active_contract"`
}

// GetClientsCountApi gets the count of clients
type GetClientsCountResponse struct {
	TotalClients         int64 `json:"total_clients"`
	ClientsInCare        int64 `json:"clients_in_care"`
	ClientsOnWaitingList int64 `json:"clients_on_waiting_list"`
	ClientsOutOfCare     int64 `json:"clients_out_of_care"`
}

// GetClientApiResponse represents the client GET page response.
// This is a stable envelope + nested objects; consumers should not rely on old flat fields.
type GetClientApiResponse struct {
	SchemaVersion     int32                            `json:"schema_version"`
	Status            string                           `json:"status"`
	Client            ClientPageClientResponse         `json:"client"`
	Care              *ClientInCareResponse            `json:"care,omitempty"`
	CareSchedule      *ClientCareScheduleResponse      `json:"care_schedule,omitempty"`
	DischargeSchedule *ClientDischargeScheduleResponse `json:"discharge_schedule,omitempty"`
	DischargeSummary  *ClientDischargeSummaryResponse  `json:"discharge_summary,omitempty"`
	Sender            *ClientSenderMinimalResponse     `json:"sender"`
	Coordinator       *ClientCoordinatorResponse       `json:"coordinator,omitempty"`
	ContractSummary   *ClientContractSummaryResponse   `json:"contract_summary,omitempty"`
	EvaluationSummary *ClientEvaluationSummaryResponse `json:"evaluation_summary,omitempty"`
	EmergencyContacts []ClientEmergencySummary         `json:"emergency_contacts"`
	Documents         ClientDocumentsSummary           `json:"documents"`
	Goals             []ClientGoalSummaryResponse      `json:"goals"`
	Intake            ClientIntakeResponse             `json:"intake"`
	Risks             ClientRiskSummary                `json:"risks"`
	Counts            ClientPageCounts                 `json:"counts"`
	Alerts            []ClientPageAlert                `json:"alerts"`
	Meta              ClientPageMetaResponse           `json:"meta"`
	StatusTimeline    *ClientStatusTimelineResponse    `json:"status_timeline,omitempty"`
}

type ClientPageClientResponse struct {
	ID          uuid.UUID               `json:"id"`
	FirstName   string                  `json:"first_name"`
	LastName    string                  `json:"last_name"`
	Bsn         *string                 `json:"bsn"`
	FileNumber  string                  `json:"file_number"`
	Gender      string                  `json:"gender"`
	DateOfBirth *time.Time              `json:"date_of_birth"`
	Age         *int32                  `json:"age"`
	CareType    *string                 `json:"care_type"`
	Address     ClientAddressResponse   `json:"address"`
	Location    *ClientLocationResponse `json:"location"`
}

type ClientLocationResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ClientPageMetaResponse struct {
	WaitlistSince time.Time `json:"waitlist_since"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
}

type ClientCareScheduleResponse struct {
	CareStartDate      *time.Time `json:"care_start_date"`
	PlacedInCareAt     *time.Time `json:"placed_in_care_at"`
	DaysUntilStart     *int32     `json:"days_until_start"`
	ShouldBeActiveNow  bool       `json:"should_be_active_now"`
	NextEvaluationDate *time.Time `json:"next_evaluation_date"`
}

type ClientDischargeScheduleResponse struct {
	DischargeDate          *time.Time `json:"discharge_date"`
	DischargeReason        *string    `json:"discharge_reason"`
	FinalEvaluation        *string    `json:"final_evaluation"`
	DaysUntilDischarge     *int32     `json:"days_until_discharge"`
	IsDue                  bool       `json:"is_due"`
	MissingFinalEvaluation bool       `json:"missing_final_evaluation"`
}

type ClientDischargeSummaryResponse struct {
	DischargeDate   *time.Time `json:"discharge_date"`
	DischargeReason *string    `json:"discharge_reason"`
	FinalEvaluation *string    `json:"final_evaluation"`
}

type ClientCoordinatorResponse struct {
	EmployeeID *uuid.UUID `json:"employee_id"`
	FirstName  *string    `json:"first_name"`
	LastName   *string    `json:"last_name"`
	StartDate  *time.Time `json:"start_date"`
}

type ClientStatusTimelineResponse struct {
	LastChangeReason *string    `json:"last_change_reason"`
	LastChangedAt    *time.Time `json:"last_changed_at"`
	LastStatus       *string    `json:"last_status"`
}

type ClientInCareResponse struct {
	CareStartDate            *time.Time `json:"care_start_date"`
	PlacedInCareAt           *time.Time `json:"placed_in_care_at"`
	DaysInCare               *int32     `json:"days_in_care"`
	EvaluationIntervalsWeeks int32      `json:"evaluation_intervals_weeks"`
	LastEvaluationAnchorDate *time.Time `json:"last_evaluation_anchor_date"`
	NextEvaluationDate       *time.Time `json:"next_evaluation_date"`
}

type ClientContractSummaryResponse struct {
	HasActiveApprovedContract bool                          `json:"has_active_approved_contract"`
	ActiveContract            *ClientActiveContractResponse `json:"active_contract,omitempty"`
	DaysUntilContractEnd      *int32                        `json:"days_until_contract_end"`
}

type ClientActiveContractResponse struct {
	ID              uuid.UUID  `json:"id"`
	Status          *string    `json:"status"`
	StartDate       *time.Time `json:"start_date"`
	EndDate         *time.Time `json:"end_date"`
	FinancingAct    *string    `json:"financing_act"`
	FinancingOption *string    `json:"financing_option"`
	CareType        *string    `json:"care_type"`
}

type ClientEvaluationSummaryResponse struct {
	NextEvaluationDate *time.Time                             `json:"next_evaluation_date"`
	DaysLeft           *int32                                 `json:"days_left"`
	Priority           *string                                `json:"priority"`
	Draft              *ClientEvaluationDraftSummaryResponse  `json:"draft,omitempty"`
	LastCompleted      *ClientEvaluationLastCompletedResponse `json:"last_completed,omitempty"`
}

type ClientEvaluationDraftSummaryResponse struct {
	ID        uuid.UUID  `json:"id"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type ClientEvaluationLastCompletedResponse struct {
	ID                  uuid.UUID  `json:"id"`
	SubmittedAt         *time.Time `json:"submitted_at"`
	CreatedByEmployeeID *uuid.UUID `json:"created_by_employee_id"`
	CreatorName         *string    `json:"creator_name"`
}

type ClientAddressResponse struct {
	Street              string  `json:"street"`
	HouseNumber         string  `json:"house_number"`
	HouseNumberAddition *string `json:"house_number_addition"`
	PostalCode          string  `json:"postal_code"`
	City                string  `json:"city"`
}

type ClientSenderMinimalResponse struct {
	Name         string  `json:"name"`
	EmailAddress *string `json:"email_address"`
	PhoneNumber  *string `json:"phone_number"`
}

type ClientIntakeResponse struct {
	SelfSufficiencyScore *int32  `json:"self_sufficiency_score"`
	Conclusion           *string `json:"conclusion"`
	ConclusionNotes      *string `json:"conclusion_notes"`
}

type ClientGoalSummaryResponse struct {
	Title     string  `json:"title"`
	Priority  string  `json:"priority"`
	TopicName *string `json:"topic_name"`
}

type ClientEmergencySummary struct {
	ID           uuid.UUID `json:"id"`
	FirstName    *string   `json:"first_name"`
	LastName     *string   `json:"last_name"`
	Relationship *string   `json:"relationship"`
	PhoneNumber  *string   `json:"phone_number"`
	Email        *string   `json:"email"`
}

type ClientDocumentsSummary struct {
	Existing []string `json:"existing"`
	Missing  []string `json:"missing"`
}

type ClientRiskSummary struct {
	Flags []string `json:"flags"`
	Notes *string  `json:"notes"`
}

type ClientPageCounts struct {
	Contracts    int64 `json:"contracts"`
	Incidents    int64 `json:"incidents"`
	Reports      int64 `json:"reports"`
	Evaluations  int64 `json:"evaluations"`
	Documents    int64 `json:"documents"`
	Appointments int64 `json:"appointments"`
}

type ClientPageAlert struct {
	Code     string `json:"code"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

// GetClientAddressesApiResponse represents a response to a get client addresses request
type GetClientAddressesApiResponse struct {
	Addresses []Address `json:"addresses"`
}

// UpdateClientDetailsRequest represents a request to update a client
type UpdateClientDetailsRequest struct {
	FirstName                  *string    `json:"first_name"`
	LastName                   *string    `json:"last_name"`
	DateOfBirth                time.Time  `json:"date_of_birth"`
	Identity                   *bool      `json:"identity"`
	Bsn                        *string    `json:"bsn"`
	BsnVerifiedBy              *uuid.UUID `json:"bsn_verified_by"`
	Source                     *string    `json:"source"`
	Birthplace                 *string    `json:"birthplace"`
	Nationality                *string    `json:"nationality"`
	Email                      *string    `json:"email"`
	PhoneNumber                *string    `json:"phone_number"`
	OrganizationID             *uuid.UUID `json:"organization_id"`
	Departement                *string    `json:"departement"`
	Gender                     *string    `json:"gender"`
	Filenumber                 *string    `json:"filenumber"`
	ProfilePicture             *string    `json:"profile_picture"`
	Infix                      *string    `json:"infix"`
	SenderID                   *uuid.UUID `json:"sender_id"`
	LocationID                 *uuid.UUID `json:"location_id"`
	DepartureReason            *string    `json:"departure_reason"`
	DepartureReport            *string    `json:"departure_report"`
	LegalMeasure               *string    `json:"legal_measure"`
	EducationCurrentlyEnrolled *bool      `json:"education_currently_enrolled"`
	EducationInstitution       *string    `json:"education_institution"`
	EducationMentorName        *string    `json:"education_mentor_name"`
	EducationMentorPhone       *string    `json:"education_mentor_phone"`
	EducationMentorEmail       *string    `json:"education_mentor_email"`
	EducationAdditionalNotes   *string    `json:"education_additional_notes"`
	EducationLevel             *string    `json:"education_level"`
	WorkCurrentlyEmployed      *bool      `json:"work_currently_employed"`
	WorkCurrentEmployer        *string    `json:"work_current_employer"`
	WorkCurrentEmployerPhone   *string    `json:"work_employer_phone"`
	WorkCurrentEmployerEmail   *string    `json:"work_employer_email"`
	WorkCurrentPosition        *string    `json:"work_current_position"`
	WorkStartDate              time.Time  `json:"work_start_date"`
	WorkAdditionalNotes        *string    `json:"work_additional_notes"`
	LivingSituation            *string    `json:"living_situation"`
	LivingSituationNotes       *string    `json:"living_situation_notes"`
}

// UpdateClientDetailsResponse represents a response to an update client request
type UpdateClientDetailsResponse struct {
	ID                    uuid.UUID  `json:"id"`
	FirstName             string     `json:"first_name"`
	LastName              string     `json:"last_name"`
	DateOfBirth           time.Time  `json:"date_of_birth"`
	Identity              bool       `json:"identity"`
	Status                string     `json:"status"`
	Bsn                   *string    `json:"bsn"`
	BsnVerifiedBy         *uuid.UUID `json:"bsn_verified_by"`
	Source                *string    `json:"source"`
	Birthplace            *string    `json:"birthplace"`
	Email                 string     `json:"email"`
	PhoneNumber           *string    `json:"phone_number"`
	OrganizationID        *uuid.UUID `json:"organization_id"`
	Departement           *string    `json:"departement"`
	Gender                string     `json:"gender"`
	Filenumber            string     `json:"filenumber"`
	ProfilePicture        *string    `json:"profile_picture"`
	Infix                 *string    `json:"infix"`
	Created               time.Time  `json:"created"`
	SenderID              *uuid.UUID `json:"sender_id"`
	LocationID            *uuid.UUID `json:"location_id"`
	DepartureReason       *string    `json:"departure_reason"`
	DepartureReport       *string    `json:"departure_report"`
	Addresses             []Address  `json:"addresses"`
	LegalMeasure          *string    `json:"legal_measure"`
	HasUntakenMedications bool       `json:"has_untaken_medications"`
}

// UpdateClientStatusRequest represents a request to update a client status
type UpdateClientStatusRequest struct {
	Status        string    `json:"status" binding:"required"`
	Reason        string    `json:"reason"`
	IsSchedueled  bool      `json:"schedueled"`
	SchedueledFor time.Time `json:"schedueled_for"`
}

// UpdateClientStatusResponse represents a response to an update client request
type UpdateClientStatusResponse struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

type PutClientInCareRequest struct {
	CareStartDate         string     `json:"care_start_date" binding:"required"`
	CoordinatorEmployeeID uuid.UUID  `json:"coordinator_employee_id" binding:"required"`
	PlacedInCareAt        *time.Time `json:"placed_in_care_at"`
	Reason                *string    `json:"reason"`
}

type PutClientInCareResponse struct {
	ID                  uuid.UUID  `json:"id"`
	Status              string     `json:"status"`
	CareStartDate       time.Time  `json:"care_start_date"`
	PlacedInCareAt      time.Time  `json:"placed_in_care_at"`
	NextEvaluationDate  *time.Time `json:"next_evaluation_date"`
	CoordinatorAssignID uuid.UUID  `json:"coordinator_assignment_id"`
	Warning             *string    `json:"warning,omitempty"`
}

type PutClientOutOfCareRequest struct {
	DischargeDate   string  `json:"discharge_date" binding:"required"`
	DischargeReason string  `json:"discharge_reason" binding:"required,oneof=treatment_completed terminated_by_mutual_agreement terminated_by_client terminated_by_provider terminated_due_to_external_factors other"`
	FinalEvaluation *string `json:"final_evaluation"`
	Reason          *string `json:"reason"`
}

type PutClientOutOfCareResponse struct {
	ID              uuid.UUID `json:"id"`
	Status          string    `json:"status"`
	DischargeDate   time.Time `json:"discharge_date"`
	DischargeReason string    `json:"discharge_reason"`
	FinalEvaluation *string   `json:"final_evaluation"`
}

// ListStatusHistoryApiResponse represents a response to a list status history request
type ListStatusHistoryApiResponse struct {
	ID        uuid.UUID  `json:"id"`
	ClientID  uuid.UUID  `json:"client_id"`
	OldStatus *string    `json:"old_status"`
	NewStatus string     `json:"new_status"`
	ChangedAt time.Time  `json:"changed_at"`
	ChangedBy *uuid.UUID `json:"changed_by"`
	Reason    *string    `json:"reason"`
}

// SetClientProfilePictureRequest represents a request to update a client
type SetClientProfilePictureRequest struct {
	AttachmentID uuid.UUID `json:"attachement_id" binding:"required"`
}

// SetClientProfilePictureResponse represents a response to a set client profile picture request
type SetClientProfilePictureResponse struct {
	ID             uuid.UUID `json:"id"`
	ProfilePicture *string   `json:"profile_picture"`
}

// AddClientDocumentApiRequest represents a request to add a document to a client
type AddClientDocumentApiRequest struct {
	AttachmentID uuid.UUID `json:"attachment_id" binding:"required"`
	Label        string    `json:"label" binding:"required"`
}

// AddClientDocumentApiResponse represents a response to an add client document request
type AddClientDocumentApiResponse struct {
	ID           uuid.UUID  `json:"id"`
	AttachmentID *uuid.UUID `json:"attachment_id"`
	ClientID     uuid.UUID  `json:"client_id"`
	Label        string     `json:"label"`
	Name         string     `json:"name"`
	File         string     `json:"file"`
	Size         int32      `json:"size"`
	IsUsed       bool       `json:"is_used"`
	Tag          *string    `json:"tag"`
	UpdatedAt    time.Time  `json:"updated"`
	CreatedAt    time.Time  `json:"created"`
}

// ListClientDocumentsApiRequest represents a request to list client documents
type ListClientDocumentsApiRequest struct {
	pagination.Request
}

// ListClientDocumentsApiResponse represents a response to a list client documents request
type ListClientDocumentsApiResponse struct {
	ID             uuid.UUID  `json:"id"`
	AttachmentUuid *uuid.UUID `json:"attachment_uuid"`
	ClientID       uuid.UUID  `json:"client_id"`
	Label          string     `json:"label"`
	Uuid           uuid.UUID  `json:"uuid"`
	Name           string     `json:"name"`
	File           *string    `json:"file"`
	Size           int32      `json:"size"`
	IsUsed         bool       `json:"is_used"`
	Tag            *string    `json:"tag"`
	UpdatedAt      time.Time  `json:"updated_at"`
	CreatedAt      time.Time  `json:"created_at"`
}

// DeleteClientDocumentApiRequest represents a request to delete a client document
type DeleteClientDocumentApiRequest struct {
	AttachmentID uuid.UUID `json:"attachement_id" binding:"required"`
}

// DeleteClientDocumentApiResponse represents a response to a delete client document request
type DeleteClientDocumentApiResponse struct {
	ID           uuid.UUID  `json:"id"`
	AttachmentID *uuid.UUID `json:"attachment_id"`
}

// GetMissingClientDocumentsApiResponse represents a response to a get missing client documents request
type GetMissingClientDocumentsApiResponse struct {
	MissingDocs []string `json:"missing_docs"`
}
