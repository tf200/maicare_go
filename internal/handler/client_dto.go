package handler

import (
	"time"

	"github.com/goccy/go-json"
	"maicare_go/internal/domain"
	"maicare_go/internal/httpapi"

	"github.com/google/uuid"
)

// --- Create Client ---

type createClientRequest struct {
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
	DateOfBirth                string     `json:"date_of_birth" binding:"required"`
	PhoneNumber                *string    `json:"phone_number" binding:"required"`
	SenderID                   *uuid.UUID `json:"sender_id" binding:"required"`
	Infix                      *string    `json:"infix"`
	Source                     *string    `json:"source" binding:"required"`
	Nationality                *string    `json:"nationality"`
	Bsn                        *string    `json:"bsn"`
	BsnVerifiedBy              *uuid.UUID `json:"bsn_verified_by"`
	Addresses                  []address  `json:"addresses"`
	EducationCurrentlyEnrolled bool       `json:"education_currently_enrolled"`
	EducationInstitution       *string    `json:"education_institution"`
	EducationMentorName        *string    `json:"education_mentor_name"`
	EducationMentorPhone       *string    `json:"education_mentor_phone"`
	EducationMentorEmail       *string    `json:"education_mentor_email"`
	EducationAdditionalNotes   *string    `json:"education_additional_notes"`
	EducationLevel             *string    `json:"education_level" binding:"omitempty,oneof=primary secondary higher none"`
	WorkCurrentlyEmployed      bool       `json:"work_currently_employed"`
	WorkCurrentEmployer        *string    `json:"work_current_employer"`
	WorkCurrentEmployerPhone   *string    `json:"work_employer_phone"`
	WorkCurrentEmployerEmail   *string    `json:"work_employer_email"`
	WorkCurrentPosition        *string    `json:"work_current_position"`
	WorkStartDate              time.Time  `json:"work_start_date"`
	WorkAdditionalNotes        *string    `json:"work_additional_notes"`
	LivingSituation            *string    `json:"living_situation" binding:"omitempty,oneof=home foster_care youth_care_institution other"`
	LivingSituationNotes       *string    `json:"living_situation_notes"`
}

type address struct {
	BelongsTo           *string `json:"belongs_to"`
	Street              *string `json:"street"`
	HouseNumber         *string `json:"house_number"`
	HouseNumberAddition *string `json:"house_number_addition"`
	PostalCode          *string `json:"postal_code"`
	City                *string `json:"city"`
	PhoneNumber         *string `json:"phone_number"`
}

type createClientResponse struct {
	ID                         uuid.UUID  `json:"id"`
	FirstName                  string     `json:"first_name"`
	LastName                   string     `json:"last_name"`
	DateOfBirth                time.Time  `json:"date_of_birth"`
	Identity                   bool       `json:"identity"`
	Status                     string     `json:"status"`
	Bsn                        *string    `json:"bsn"`
	BsnVerifiedBy              *uuid.UUID `json:"bsn_verified_by"`
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
	Addresses                  []address  `json:"addresses"`
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

// --- List Clients ---

type listClientsRequest struct {
	httpapi.PageRequest
	Status     *string    `form:"status"`
	LocationID *uuid.UUID `form:"location_id"`
	Search     *string    `form:"search"`
}

type listClientsResponse struct {
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

// --- List Waiting List Clients ---

type listWaitingListClientsRequest struct {
	httpapi.PageRequest
	Search    *string `form:"search" binding:"omitempty,max=120"`
	Placement *string `form:"placement" binding:"omitempty,oneof=protected_living training_center supported_independent_living ambulatory_support other"`
	SortDays  *string `form:"sort_days" binding:"omitempty,oneof=asc desc"`
}

type listWaitingListClientsResponse struct {
	ID             uuid.UUID `json:"id"`
	FirstName      string    `json:"first_name"`
	Bsn            *string   `json:"bsn"`
	LastName       string    `json:"last_name"`
	CareType       *string   `json:"care_type"`
	SenderName     *string   `json:"sender_name"`
	DaysInWaitlist int32     `json:"days_in_waitlist"`
	AdmissionType  *string   `json:"admission_type"`
}

// --- List In-Care Clients ---

type listInCareClientsRequest struct {
	httpapi.PageRequest
	Search         *string  `form:"search" binding:"omitempty,max=120"`
	Status         []string `form:"status" binding:"omitempty,dive,oneof=in_care scheduled_in_care"`
	SortDaysInCare *string  `form:"sort_days_in_care" binding:"omitempty,oneof=asc desc"`
}

type listInCareClientsResponse struct {
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

// --- Counts ---

type getClientsCountResponse struct {
	TotalClients         int64 `json:"total_clients"`
	ClientsInCare        int64 `json:"clients_in_care"`
	ClientsOnWaitingList int64 `json:"clients_on_waiting_list"`
	ClientsOutOfCare     int64 `json:"clients_out_of_care"`
}

type getClientStatusCountsResponse struct {
	ClientsInOrScheduledInCare     int64 `json:"clients_in_or_scheduled_in_care"`
	ClientsOnWaitingList           int64 `json:"clients_on_waiting_list"`
	ClientsOutOrScheduledOutOfCare int64 `json:"clients_out_or_scheduled_out_of_care"`
}

// --- Mappers ---

func toCreateClientParams(req createClientRequest) (domain.CreateClientParams, error) {
	parsedDateOfBirth, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return domain.CreateClientParams{}, err
	}

	return domain.CreateClientParams{
		FirstName:                  req.FirstName,
		LastName:                   req.LastName,
		DateOfBirth:                parsedDateOfBirth,
		Bsn:                        req.Bsn,
		BsnVerifiedBy:              req.BsnVerifiedBy,
		Email:                      req.Email,
		PhoneNumber:                req.PhoneNumber,
		CareType:                   req.CareType,
		SenderID:                   req.SenderID,
		LocationID:                 req.LocationID,
		EducationCurrentlyEnrolled: req.EducationCurrentlyEnrolled,
		EducationInstitution:       req.EducationInstitution,
		EducationMentorName:        req.EducationMentorName,
		EducationMentorPhone:       req.EducationMentorPhone,
		EducationMentorEmail:       req.EducationMentorEmail,
		EducationAdditionalNotes:   req.EducationAdditionalNotes,
		WorkCurrentlyEmployed:      req.WorkCurrentlyEmployed,
		WorkCurrentEmployer:        req.WorkCurrentEmployer,
		WorkCurrentEmployerPhone:   req.WorkCurrentEmployerPhone,
		WorkCurrentEmployerEmail:   req.WorkCurrentEmployerEmail,
		WorkCurrentPosition:        req.WorkCurrentPosition,
		WorkStartDate:              req.WorkStartDate,
		WorkAdditionalNotes:        req.WorkAdditionalNotes,
	}, nil
}

func toCreateClientResponse(client *domain.Client) createClientResponse {
	return createClientResponse{
		ID:                         client.ID,
		FirstName:                  client.FirstName,
		LastName:                   client.LastName,
		DateOfBirth:                client.DateOfBirth,
		Identity:                   client.Identity,
		Status:                     client.Status,
		Bsn:                        client.Bsn,
		BsnVerifiedBy:              client.BsnVerifiedBy,
		Email:                      client.Email,
		PhoneNumber:                client.PhoneNumber,
		Gender:                     client.Gender,
		Filenumber:                 client.Filenumber,
		Created:                    client.CreatedAt,
		SenderID:                   client.SenderID,
		LocationID:                 client.LocationID,
		EducationCurrentlyEnrolled: client.EducationCurrentlyEnrolled,
		EducationInstitution:       client.EducationInstitution,
		EducationMentorName:        client.EducationMentorName,
		EducationMentorPhone:       client.EducationMentorPhone,
		EducationMentorEmail:       client.EducationMentorEmail,
		EducationAdditionalNotes:   client.EducationAdditionalNotes,
		EducationLevel:             client.EducationLevel,
		WorkCurrentlyEmployed:      client.WorkCurrentlyEmployed,
		WorkCurrentEmployer:        client.WorkCurrentEmployer,
		WorkCurrentEmployerPhone:   client.WorkCurrentEmployerPhone,
		WorkCurrentEmployerEmail:   client.WorkCurrentEmployerEmail,
		WorkCurrentPosition:        client.WorkCurrentPosition,
		WorkStartDate:              client.WorkStartDate,
		WorkAdditionalNotes:        client.WorkAdditionalNotes,
	}
}

func toListClientsParams(req listClientsRequest) domain.ListClientsParams {
	params := req.PageRequest.Params()
	return domain.ListClientsParams{
		Limit:      params.Limit,
		Offset:     params.Offset,
		Status:     req.Status,
		LocationID: req.LocationID,
		Search:     req.Search,
	}
}

func toListClientsResponse(item domain.ClientListItem) listClientsResponse {
	return listClientsResponse{
		ID:           item.ID,
		FirstName:    item.FirstName,
		LastName:     item.LastName,
		Bsn:          item.Bsn,
		Filenumber:   item.Filenumber,
		LocationName: item.LocationName,
		CareType:     item.CareType,
		Status:       item.Status,
		GoalsCount:   item.GoalsCount,
		RiskCount:    item.RiskCount,
		CreatedAt:    item.CreatedAt,
	}
}

func toListWaitingListClientsParams(req listWaitingListClientsRequest) domain.ListWaitingListClientsParams {
	params := req.PageRequest.Params()
	sortDays := "desc"
	if req.SortDays != nil {
		sortDays = *req.SortDays
	}

	return domain.ListWaitingListClientsParams{
		Limit:     params.Limit,
		Offset:    params.Offset,
		Search:    req.Search,
		Placement: req.Placement,
		SortDays:  sortDays,
	}
}

func toListWaitingListClientsResponse(item domain.WaitingListClient) listWaitingListClientsResponse {
	return listWaitingListClientsResponse{
		ID:             item.ID,
		FirstName:      item.FirstName,
		Bsn:            item.Bsn,
		LastName:       item.LastName,
		CareType:       item.CareType,
		SenderName:     item.SenderName,
		DaysInWaitlist: item.DaysInWaitlist,
		AdmissionType:  item.AdmissionType,
	}
}

func toListInCareClientsParams(req listInCareClientsRequest) domain.ListInCareClientsParams {
	params := req.PageRequest.Params()
	sortDaysInCare := "desc"
	if req.SortDaysInCare != nil {
		sortDaysInCare = *req.SortDaysInCare
	}

	return domain.ListInCareClientsParams{
		Limit:          params.Limit,
		Offset:         params.Offset,
		Search:         req.Search,
		Status:         req.Status,
		SortDaysInCare: sortDaysInCare,
	}
}

func toListInCareClientsResponse(item domain.InCareClient) listInCareClientsResponse {
	return listInCareClientsResponse{
		ID:                item.ID,
		Bsn:               item.Bsn,
		FirstName:         item.FirstName,
		LastName:          item.LastName,
		CoordinatorName:   item.CoordinatorName,
		LocationName:      item.LocationName,
		Status:            item.Status,
		CareStartDate:     item.CareStartDate,
		DaysInCare:        item.DaysInCare,
		HasActiveContract: item.HasActiveContract,
	}
}

func toGetClientsCountResponse(counts *domain.ClientCounts) getClientsCountResponse {
	return getClientsCountResponse{
		TotalClients:         counts.TotalClients,
		ClientsInCare:        counts.ClientsInCare,
		ClientsOnWaitingList: counts.ClientsOnWaitingList,
		ClientsOutOfCare:     counts.ClientsOutOfCare,
	}
}

func toGetClientStatusCountsResponse(counts *domain.ClientStatusCounts) getClientStatusCountsResponse {
	return getClientStatusCountsResponse{
		ClientsInOrScheduledInCare:     counts.ClientsInOrScheduledInCare,
		ClientsOnWaitingList:           counts.ClientsOnWaitingList,
		ClientsOutOrScheduledOutOfCare: counts.ClientsOutOrScheduledOutOfCare,
	}
}

// --- Get Client By ID ---

type getClientResponse struct {
	SchemaVersion     int32                            `json:"schema_version"`
	Status            string                           `json:"status"`
	Client            clientPageClientResponse         `json:"client"`
	Care              *clientInCareResponse            `json:"care,omitempty"`
	CareSchedule      *clientCareScheduleResponse      `json:"care_schedule,omitempty"`
	DischargeSchedule *clientDischargeScheduleResponse `json:"discharge_schedule,omitempty"`
	DischargeSummary  *clientDischargeSummaryResponse  `json:"discharge_summary,omitempty"`
	Sender            *clientSenderMinimalResponse     `json:"sender"`
	Coordinator       *clientCoordinatorResponse       `json:"coordinator,omitempty"`
	ContractSummary   *clientContractSummaryResponse   `json:"contract_summary,omitempty"`
	EvaluationSummary *clientEvaluationSummaryResponse `json:"evaluation_summary,omitempty"`
	EmergencyContacts []clientEmergencySummary         `json:"emergency_contacts"`
	Documents         clientDocumentsSummary           `json:"documents"`
	Goals             []clientGoalSummaryResponse      `json:"goals"`
	Intake            clientIntakeResponse             `json:"intake"`
	Risks             clientRiskSummary                `json:"risks"`
	Counts            clientPageCounts                 `json:"counts"`
	Alerts            []clientPageAlert                `json:"alerts"`
	Meta              clientPageMetaResponse           `json:"meta"`
	StatusTimeline    *clientStatusTimelineResponse    `json:"status_timeline,omitempty"`
}

type clientPageClientResponse struct {
	ID          uuid.UUID               `json:"id"`
	FirstName   string                  `json:"first_name"`
	LastName    string                  `json:"last_name"`
	Bsn         *string                 `json:"bsn"`
	FileNumber  string                  `json:"file_number"`
	Gender      string                  `json:"gender"`
	DateOfBirth *time.Time              `json:"date_of_birth"`
	Age         *int32                  `json:"age"`
	CareType    *string                 `json:"care_type"`
	Address     clientAddressResponse   `json:"address"`
	Location    *clientLocationResponse `json:"location"`
}

type clientLocationResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type clientAddressResponse struct {
	Street              string  `json:"street"`
	HouseNumber         string  `json:"house_number"`
	HouseNumberAddition *string `json:"house_number_addition"`
	PostalCode          string  `json:"postal_code"`
	City                string  `json:"city"`
}

type clientInCareResponse struct {
	CareStartDate            *time.Time `json:"care_start_date"`
	PlacedInCareAt           *time.Time `json:"placed_in_care_at"`
	DaysInCare               *int32     `json:"days_in_care"`
	EvaluationIntervalsWeeks int32      `json:"evaluation_intervals_weeks"`
	LastEvaluationAnchorDate *time.Time `json:"last_evaluation_anchor_date"`
	NextEvaluationDate       *time.Time `json:"next_evaluation_date"`
}

type clientCareScheduleResponse struct {
	CareStartDate      *time.Time `json:"care_start_date"`
	PlacedInCareAt     *time.Time `json:"placed_in_care_at"`
	DaysUntilStart     *int32     `json:"days_until_start"`
	ShouldBeActiveNow  bool       `json:"should_be_active_now"`
	NextEvaluationDate *time.Time `json:"next_evaluation_date"`
}

type clientDischargeScheduleResponse struct {
	DischargeDate          *time.Time `json:"discharge_date"`
	DischargeReason        *string    `json:"discharge_reason"`
	FinalEvaluation        *string    `json:"final_evaluation"`
	DaysUntilDischarge     *int32     `json:"days_until_discharge"`
	IsDue                  bool       `json:"is_due"`
	MissingFinalEvaluation bool       `json:"missing_final_evaluation"`
}

type clientDischargeSummaryResponse struct {
	DischargeDate   *time.Time `json:"discharge_date"`
	DischargeReason *string    `json:"discharge_reason"`
	FinalEvaluation *string    `json:"final_evaluation"`
}

type clientSenderMinimalResponse struct {
	Name         string  `json:"name"`
	EmailAddress *string `json:"email_address"`
	PhoneNumber  *string `json:"phone_number"`
}

type clientCoordinatorResponse struct {
	EmployeeID *uuid.UUID `json:"employee_id"`
	FirstName  *string    `json:"first_name"`
	LastName   *string    `json:"last_name"`
	StartDate  *time.Time `json:"start_date"`
}

type clientContractSummaryResponse struct {
	HasActiveApprovedContract bool                          `json:"has_active_approved_contract"`
	ActiveContract            *clientActiveContractResponse `json:"active_contract,omitempty"`
	DaysUntilContractEnd      *int32                        `json:"days_until_contract_end"`
}

type clientActiveContractResponse struct {
	ID              uuid.UUID  `json:"id"`
	Status          *string    `json:"status"`
	StartDate       *time.Time `json:"start_date"`
	EndDate         *time.Time `json:"end_date"`
	FinancingAct    *string    `json:"financing_act"`
	FinancingOption *string    `json:"financing_option"`
	CareType        *string    `json:"care_type"`
}

type clientEvaluationSummaryResponse struct {
	NextEvaluationDate *time.Time                             `json:"next_evaluation_date"`
	DaysLeft           *int32                                 `json:"days_left"`
	Priority           *string                                `json:"priority"`
	Draft              *clientEvaluationDraftSummaryResponse  `json:"draft,omitempty"`
	LastCompleted      *clientEvaluationLastCompletedResponse `json:"last_completed,omitempty"`
}

type clientEvaluationDraftSummaryResponse struct {
	ID        uuid.UUID  `json:"id"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type clientEvaluationLastCompletedResponse struct {
	ID                  uuid.UUID  `json:"id"`
	SubmittedAt         *time.Time `json:"submitted_at"`
	CreatedByEmployeeID *uuid.UUID `json:"created_by_employee_id"`
	CreatorName         *string    `json:"creator_name"`
}

type clientGoalSummaryResponse struct {
	Title     string  `json:"title"`
	Priority  string  `json:"priority"`
	TopicName *string `json:"topic_name"`
}

type clientEmergencySummary struct {
	ID           uuid.UUID `json:"id"`
	FirstName    *string   `json:"first_name"`
	LastName     *string   `json:"last_name"`
	Relationship *string   `json:"relationship"`
	PhoneNumber  *string   `json:"phone_number"`
	Email        *string   `json:"email"`
}

type clientDocumentsSummary struct {
	Existing []string `json:"existing"`
	Missing  []string `json:"missing"`
}

type clientIntakeResponse struct {
	SelfSufficiencyScore *int32  `json:"self_sufficiency_score"`
	Conclusion           *string `json:"conclusion"`
	ConclusionNotes      *string `json:"conclusion_notes"`
}

type clientRiskSummary struct {
	Flags []string `json:"flags"`
	Notes *string  `json:"notes"`
}

type clientPageCounts struct {
	Contracts    int64 `json:"contracts"`
	Incidents    int64 `json:"incidents"`
	Reports      int64 `json:"reports"`
	Evaluations  int64 `json:"evaluations"`
	Documents    int64 `json:"documents"`
	Appointments int64 `json:"appointments"`
}

type clientPageAlert struct {
	Code     string `json:"code"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

type clientPageMetaResponse struct {
	WaitlistSince time.Time `json:"waitlist_since"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
}

type clientStatusTimelineResponse struct {
	LastChangeReason *string    `json:"last_change_reason"`
	LastChangedAt    *time.Time `json:"last_changed_at"`
	LastStatus       *string    `json:"last_status"`
}

// --- Update Client ---

type updateClientRequest struct {
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

type updateClientResponse struct {
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
	Addresses             []address  `json:"addresses"`
	LegalMeasure          *string    `json:"legal_measure"`
	HasUntakenMedications bool       `json:"has_untaken_medications"`
}

// --- Get Client Addresses ---

type getClientAddressesResponse struct {
	Addresses []address `json:"addresses"`
}

// --- Mappers (Fragment 2) ---

func toGetClientResponse(detail *domain.ClientPageDetail) getClientResponse {
	var sender *clientSenderMinimalResponse
	if detail.Sender != nil {
		sender = &clientSenderMinimalResponse{
			Name:         detail.Sender.Name,
			EmailAddress: detail.Sender.EmailAddress,
			PhoneNumber:  detail.Sender.PhoneNumber,
		}
	}

	var location *clientLocationResponse
	if detail.Client.Location != nil {
		location = &clientLocationResponse{
			ID:   detail.Client.Location.ID,
			Name: detail.Client.Location.Name,
		}
	}

	var care *clientInCareResponse
	if detail.Care != nil {
		care = &clientInCareResponse{
			CareStartDate:            detail.Care.CareStartDate,
			PlacedInCareAt:           detail.Care.PlacedInCareAt,
			DaysInCare:               detail.Care.DaysInCare,
			EvaluationIntervalsWeeks: detail.Care.EvaluationIntervalsWeeks,
			LastEvaluationAnchorDate: detail.Care.LastEvaluationAnchorDate,
			NextEvaluationDate:       detail.Care.NextEvaluationDate,
		}
	}

	var careSchedule *clientCareScheduleResponse
	if detail.CareSchedule != nil {
		careSchedule = &clientCareScheduleResponse{
			CareStartDate:      detail.CareSchedule.CareStartDate,
			PlacedInCareAt:     detail.CareSchedule.PlacedInCareAt,
			DaysUntilStart:     detail.CareSchedule.DaysUntilStart,
			ShouldBeActiveNow:  detail.CareSchedule.ShouldBeActiveNow,
			NextEvaluationDate: detail.CareSchedule.NextEvaluationDate,
		}
	}

	var dischargeSchedule *clientDischargeScheduleResponse
	if detail.DischargeSchedule != nil {
		dischargeSchedule = &clientDischargeScheduleResponse{
			DischargeDate:          detail.DischargeSchedule.DischargeDate,
			DischargeReason:        detail.DischargeSchedule.DischargeReason,
			FinalEvaluation:        detail.DischargeSchedule.FinalEvaluation,
			DaysUntilDischarge:     detail.DischargeSchedule.DaysUntilDischarge,
			IsDue:                  detail.DischargeSchedule.IsDue,
			MissingFinalEvaluation: detail.DischargeSchedule.MissingFinalEvaluation,
		}
	}

	var dischargeSummary *clientDischargeSummaryResponse
	if detail.DischargeSummary != nil {
		dischargeSummary = &clientDischargeSummaryResponse{
			DischargeDate:   detail.DischargeSummary.DischargeDate,
			DischargeReason: detail.DischargeSummary.DischargeReason,
			FinalEvaluation: detail.DischargeSummary.FinalEvaluation,
		}
	}

	var coordinator *clientCoordinatorResponse
	if detail.Coordinator != nil {
		coordinator = &clientCoordinatorResponse{
			EmployeeID: detail.Coordinator.EmployeeID,
			FirstName:  detail.Coordinator.FirstName,
			LastName:   detail.Coordinator.LastName,
			StartDate:  detail.Coordinator.StartDate,
		}
	}

	var contractSummary *clientContractSummaryResponse
	if detail.ContractSummary != nil {
		var activeContract *clientActiveContractResponse
		if detail.ContractSummary.ActiveContract != nil {
			activeContract = &clientActiveContractResponse{
				ID:              detail.ContractSummary.ActiveContract.ID,
				Status:          detail.ContractSummary.ActiveContract.Status,
				StartDate:       detail.ContractSummary.ActiveContract.StartDate,
				EndDate:         detail.ContractSummary.ActiveContract.EndDate,
				FinancingAct:    detail.ContractSummary.ActiveContract.FinancingAct,
				FinancingOption: detail.ContractSummary.ActiveContract.FinancingOption,
				CareType:        detail.ContractSummary.ActiveContract.CareType,
			}
		}
		contractSummary = &clientContractSummaryResponse{
			HasActiveApprovedContract: detail.ContractSummary.HasActiveApprovedContract,
			ActiveContract:            activeContract,
			DaysUntilContractEnd:      detail.ContractSummary.DaysUntilContractEnd,
		}
	}

	var evaluationSummary *clientEvaluationSummaryResponse
	if detail.EvaluationSummary != nil {
		var draft *clientEvaluationDraftSummaryResponse
		if detail.EvaluationSummary.Draft != nil {
			draft = &clientEvaluationDraftSummaryResponse{
				ID:        detail.EvaluationSummary.Draft.ID,
				UpdatedAt: detail.EvaluationSummary.Draft.UpdatedAt,
			}
		}
		var lastCompleted *clientEvaluationLastCompletedResponse
		if detail.EvaluationSummary.LastCompleted != nil {
			lastCompleted = &clientEvaluationLastCompletedResponse{
				ID:                  detail.EvaluationSummary.LastCompleted.ID,
				SubmittedAt:         detail.EvaluationSummary.LastCompleted.SubmittedAt,
				CreatedByEmployeeID: detail.EvaluationSummary.LastCompleted.CreatedByEmployeeID,
				CreatorName:         detail.EvaluationSummary.LastCompleted.CreatorName,
			}
		}
		evaluationSummary = &clientEvaluationSummaryResponse{
			NextEvaluationDate: detail.EvaluationSummary.NextEvaluationDate,
			DaysLeft:           detail.EvaluationSummary.DaysLeft,
			Priority:           detail.EvaluationSummary.Priority,
			Draft:              draft,
			LastCompleted:      lastCompleted,
		}
	}

	emergencyContacts := make([]clientEmergencySummary, len(detail.EmergencyContacts))
	for i, ec := range detail.EmergencyContacts {
		emergencyContacts[i] = clientEmergencySummary{
			ID:           ec.ID,
			FirstName:    ec.FirstName,
			LastName:     ec.LastName,
			Relationship: ec.Relationship,
			PhoneNumber:  ec.PhoneNumber,
			Email:        ec.Email,
		}
	}

	goals := make([]clientGoalSummaryResponse, len(detail.Goals))
	for i, g := range detail.Goals {
		goals[i] = clientGoalSummaryResponse{
			Title:     g.Title,
			Priority:  g.Priority,
			TopicName: g.TopicName,
		}
	}

	alerts := make([]clientPageAlert, len(detail.Alerts))
	for i, a := range detail.Alerts {
		alerts[i] = clientPageAlert{
			Code:     a.Code,
			Severity: a.Severity,
			Message:  a.Message,
		}
	}

	var statusTimeline *clientStatusTimelineResponse
	if detail.StatusTimeline != nil {
		statusTimeline = &clientStatusTimelineResponse{
			LastChangeReason: detail.StatusTimeline.LastChangeReason,
			LastChangedAt:    detail.StatusTimeline.LastChangedAt,
			LastStatus:       detail.StatusTimeline.LastStatus,
		}
	}

	return getClientResponse{
		SchemaVersion: detail.SchemaVersion,
		Status:        detail.Status,
		Client: clientPageClientResponse{
			ID:          detail.Client.ID,
			FirstName:   detail.Client.FirstName,
			LastName:    detail.Client.LastName,
			Bsn:         detail.Client.Bsn,
			FileNumber:  detail.Client.FileNumber,
			Gender:      detail.Client.Gender,
			DateOfBirth: detail.Client.DateOfBirth,
			Age:         detail.Client.Age,
			CareType:    detail.Client.CareType,
			Address: clientAddressResponse{
				Street:              detail.Client.Address.Street,
				HouseNumber:         detail.Client.Address.HouseNumber,
				HouseNumberAddition: detail.Client.Address.HouseNumberAddition,
				PostalCode:          detail.Client.Address.PostalCode,
				City:                detail.Client.Address.City,
			},
			Location: location,
		},
		Care:              care,
		CareSchedule:      careSchedule,
		DischargeSchedule: dischargeSchedule,
		DischargeSummary:  dischargeSummary,
		Sender:            sender,
		Coordinator:       coordinator,
		ContractSummary:   contractSummary,
		EvaluationSummary: evaluationSummary,
		EmergencyContacts: emergencyContacts,
		Documents: clientDocumentsSummary{
			Existing: detail.Documents.Existing,
			Missing:  detail.Documents.Missing,
		},
		Goals: goals,
		Intake: clientIntakeResponse{
			SelfSufficiencyScore: detail.Intake.SelfSufficiencyScore,
			Conclusion:           detail.Intake.Conclusion,
			ConclusionNotes:      detail.Intake.ConclusionNotes,
		},
		Risks: clientRiskSummary{
			Flags: detail.Risks.Flags,
			Notes: detail.Risks.Notes,
		},
		Counts: clientPageCounts{
			Contracts:    detail.Counts.Contracts,
			Incidents:    detail.Counts.Incidents,
			Reports:      detail.Counts.Reports,
			Evaluations:  detail.Counts.Evaluations,
			Documents:    detail.Counts.Documents,
			Appointments: detail.Counts.Appointments,
		},
		Alerts: alerts,
		Meta: clientPageMetaResponse{
			WaitlistSince: detail.Meta.WaitlistSince,
			LastUpdatedAt: detail.Meta.LastUpdatedAt,
		},
		StatusTimeline: statusTimeline,
	}
}

func toUpdateClientParams(req updateClientRequest) domain.UpdateClientParams {
	return domain.UpdateClientParams{
		FirstName:                  req.FirstName,
		LastName:                   req.LastName,
		DateOfBirth:                req.DateOfBirth,
		Identity:                   req.Identity,
		Bsn:                        req.Bsn,
		BsnVerifiedBy:              req.BsnVerifiedBy,
		Source:                     req.Source,
		Birthplace:                 req.Birthplace,
		Nationality:                req.Nationality,
		Email:                      req.Email,
		PhoneNumber:                req.PhoneNumber,
		OrganizationID:             req.OrganizationID,
		Departement:                req.Departement,
		Gender:                     req.Gender,
		Filenumber:                 req.Filenumber,
		ProfilePicture:             req.ProfilePicture,
		Infix:                      req.Infix,
		SenderID:                   req.SenderID,
		LocationID:                 req.LocationID,
		DepartureReason:            req.DepartureReason,
		DepartureReport:            req.DepartureReport,
		LegalMeasure:               req.LegalMeasure,
		EducationCurrentlyEnrolled: req.EducationCurrentlyEnrolled,
		EducationInstitution:       req.EducationInstitution,
		EducationMentorName:        req.EducationMentorName,
		EducationMentorPhone:       req.EducationMentorPhone,
		EducationMentorEmail:       req.EducationMentorEmail,
		EducationAdditionalNotes:   req.EducationAdditionalNotes,
		EducationLevel:             req.EducationLevel,
		WorkCurrentlyEmployed:      req.WorkCurrentlyEmployed,
		WorkCurrentEmployer:        req.WorkCurrentEmployer,
		WorkCurrentEmployerPhone:   req.WorkCurrentEmployerPhone,
		WorkCurrentEmployerEmail:   req.WorkCurrentEmployerEmail,
		WorkCurrentPosition:        req.WorkCurrentPosition,
		WorkStartDate:              req.WorkStartDate,
		WorkAdditionalNotes:        req.WorkAdditionalNotes,
		LivingSituation:            req.LivingSituation,
		LivingSituationNotes:       req.LivingSituationNotes,
	}
}

func toUpdateClientResponse(client *domain.Client) updateClientResponse {
	return updateClientResponse{
		ID:            client.ID,
		FirstName:     client.FirstName,
		LastName:      client.LastName,
		DateOfBirth:   client.DateOfBirth,
		Identity:      client.Identity,
		Status:        client.Status,
		Bsn:           client.Bsn,
		BsnVerifiedBy: client.BsnVerifiedBy,
		Email:         client.Email,
		PhoneNumber:   client.PhoneNumber,
		Gender:        client.Gender,
		Filenumber:    client.Filenumber,
		Created:       client.CreatedAt,
		SenderID:      client.SenderID,
		LocationID:    client.LocationID,
	}
}

func toGetClientAddressesResponse(addresses []domain.ClientAddress) getClientAddressesResponse {
	resp := make([]address, len(addresses))
	for i, a := range addresses {
		resp[i] = address{
			BelongsTo:           a.BelongsTo,
			Street:              &a.Street,
			HouseNumber:         &a.HouseNumber,
			HouseNumberAddition: a.HouseNumberAddition,
			PostalCode:          &a.PostalCode,
			City:                &a.City,
			PhoneNumber:         a.PhoneNumber,
		}
	}
	return getClientAddressesResponse{Addresses: resp}
}

// --- Update Client Status ---

type updateClientStatusRequest struct {
	Status       string    `json:"status" binding:"required"`
	Reason       string    `json:"reason"`
	IsScheduled  bool      `json:"schedueled"`
	ScheduledFor time.Time `json:"schedueled_for"`
}

type updateClientStatusResponse struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

func toUpdateClientStatusParams(req updateClientStatusRequest) domain.UpdateClientStatusParams {
	return domain.UpdateClientStatusParams{
		Status:       req.Status,
		Reason:       req.Reason,
		IsScheduled:  req.IsScheduled,
		ScheduledFor: req.ScheduledFor,
	}
}

func toUpdateClientStatusResponse(result *domain.UpdateClientStatusResult) updateClientStatusResponse {
	return updateClientStatusResponse{
		ID:     result.ID,
		Status: result.Status,
	}
}

// --- Put Client In Care ---

type putClientInCareRequest struct {
	CareStartDate         string     `json:"care_start_date" binding:"required"`
	CoordinatorEmployeeID uuid.UUID  `json:"coordinator_employee_id" binding:"required"`
	PlacedInCareAt        *time.Time `json:"placed_in_care_at"`
	Reason                *string    `json:"reason"`
}

type putClientInCareResponse struct {
	ID                      uuid.UUID  `json:"id"`
	Status                  string     `json:"status"`
	CareStartDate           time.Time  `json:"care_start_date"`
	PlacedInCareAt          time.Time  `json:"placed_in_care_at"`
	NextEvaluationDate      *time.Time `json:"next_evaluation_date"`
	CoordinatorAssignmentID uuid.UUID  `json:"coordinator_assignment_id"`
	Warning                 *string    `json:"warning,omitempty"`
}

func toPutClientInCareParams(req putClientInCareRequest) domain.PutClientInCareParams {
	return domain.PutClientInCareParams{
		CareStartDate:         req.CareStartDate,
		CoordinatorEmployeeID: req.CoordinatorEmployeeID,
		PlacedInCareAt:        req.PlacedInCareAt,
		Reason:                req.Reason,
	}
}

func toPutClientInCareResponse(result *domain.PutClientInCareResult) putClientInCareResponse {
	return putClientInCareResponse{
		ID:                      result.ID,
		Status:                  result.Status,
		CareStartDate:           result.CareStartDate,
		PlacedInCareAt:          result.PlacedInCareAt,
		NextEvaluationDate:      result.NextEvaluationDate,
		CoordinatorAssignmentID: result.CoordinatorAssignmentID,
		Warning:                 result.Warning,
	}
}

// --- Put Client Out Of Care ---

type putClientOutOfCareRequest struct {
	DischargeDate   string  `json:"discharge_date" binding:"required"`
	DischargeReason string  `json:"discharge_reason" binding:"required,oneof=treatment_completed terminated_by_mutual_agreement terminated_by_client terminated_by_provider terminated_due_to_external_factors other"`
	FinalEvaluation *string `json:"final_evaluation"`
	Reason          *string `json:"reason"`
}

type putClientOutOfCareResponse struct {
	ID              uuid.UUID `json:"id"`
	Status          string    `json:"status"`
	DischargeDate   time.Time `json:"discharge_date"`
	DischargeReason string    `json:"discharge_reason"`
	FinalEvaluation *string   `json:"final_evaluation"`
}

func toPutClientOutOfCareParams(req putClientOutOfCareRequest) domain.PutClientOutOfCareParams {
	return domain.PutClientOutOfCareParams{
		DischargeDate:   req.DischargeDate,
		DischargeReason: req.DischargeReason,
		FinalEvaluation: req.FinalEvaluation,
		Reason:          req.Reason,
	}
}

func toPutClientOutOfCareResponse(result *domain.PutClientOutOfCareResult) putClientOutOfCareResponse {
	return putClientOutOfCareResponse{
		ID:              result.ID,
		Status:          result.Status,
		DischargeDate:   result.DischargeDate,
		DischargeReason: result.DischargeReason,
		FinalEvaluation: result.FinalEvaluation,
	}
}

// --- List Status History ---

type listStatusHistoryResponse struct {
	ID        uuid.UUID  `json:"id"`
	ClientID  uuid.UUID  `json:"client_id"`
	OldStatus *string    `json:"old_status"`
	NewStatus string     `json:"new_status"`
	ChangedAt time.Time  `json:"changed_at"`
	ChangedBy *uuid.UUID `json:"changed_by"`
	Reason    *string    `json:"reason"`
}

func toListStatusHistoryResponse(history domain.ClientStatusHistory) listStatusHistoryResponse {
	return listStatusHistoryResponse{
		ID:        history.ID,
		ClientID:  history.ClientID,
		OldStatus: history.OldStatus,
		NewStatus: history.NewStatus,
		ChangedAt: history.ChangedAt,
		ChangedBy: history.ChangedBy,
		Reason:    history.Reason,
	}
}

// --- Add Client Document ---

type addClientDocumentItem struct {
	AttachmentID uuid.UUID `json:"attachment_id" binding:"required"`
	Label        string    `json:"label" binding:"required"`
}

type addClientDocumentRequest struct {
	Documents []addClientDocumentItem `json:"documents" binding:"required"`
}

type addClientDocumentResult struct {
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

type addClientDocumentResponse struct {
	Documents []addClientDocumentResult `json:"documents"`
}

func toAddClientDocumentParams(req addClientDocumentRequest) domain.AddClientDocumentParams {
	docs := make([]domain.AddClientDocumentItem, 0, len(req.Documents))
	for _, d := range req.Documents {
		docs = append(docs, domain.AddClientDocumentItem{
			AttachmentID: d.AttachmentID,
			Label:        d.Label,
		})
	}
	return domain.AddClientDocumentParams{Documents: docs}
}

func toAddClientDocumentResponse(results []domain.AddClientDocumentResult) addClientDocumentResponse {
	docs := make([]addClientDocumentResult, 0, len(results))
	for _, r := range results {
		docs = append(docs, addClientDocumentResult{
			ID:           r.ID,
			AttachmentID: r.AttachmentID,
			ClientID:     r.ClientID,
			Label:        r.Label,
			Name:         r.Name,
			File:         r.File,
			Size:         r.Size,
			IsUsed:       r.IsUsed,
			Tag:          r.Tag,
			UpdatedAt:    r.UpdatedAt,
			CreatedAt:    r.CreatedAt,
		})
	}
	return addClientDocumentResponse{Documents: docs}
}

// --- List Client Documents ---

type listClientDocumentsRequest struct {
	httpapi.PageRequest
}

type listClientDocumentsResponse struct {
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

func toListClientDocumentsResponse(doc domain.ClientDocument) listClientDocumentsResponse {
	var filePtr *string
	if doc.FileURL != "" {
		filePtr = &doc.FileURL
	}
	return listClientDocumentsResponse{
		ID:             doc.ID,
		AttachmentUuid: doc.AttachmentID,
		ClientID:       doc.ClientID,
		Label:          doc.Label,
		Uuid:           doc.Uuid,
		Name:           doc.Name,
		File:           filePtr,
		Size:           doc.Size,
		IsUsed:         doc.IsUsed,
		Tag:            doc.Tag,
		UpdatedAt:      doc.UpdatedAt,
		CreatedAt:      doc.CreatedAt,
	}
}

// --- Delete Client Document ---

type deleteClientDocumentResponse struct {
	ID           uuid.UUID  `json:"id"`
	AttachmentID *uuid.UUID `json:"attachment_id"`
}

func toDeleteClientDocumentResponse(result *domain.DeleteClientDocumentResult) deleteClientDocumentResponse {
	return deleteClientDocumentResponse{
		ID:           result.ID,
		AttachmentID: result.AttachmentID,
	}
}

// --- Get Missing Client Documents ---

type getMissingClientDocumentsResponse struct {
	MissingDocs []string `json:"missing_docs"`
}

// --- Create Client Goal ---

type createClientGoalRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description *string   `json:"description"`
	Priority    *string   `json:"priority" binding:"omitempty,oneof=low medium high"`
	TopicID     uuid.UUID `json:"topic_id" binding:"required"`
	SortOrder   *int32    `json:"sort_order"`
}

type createClientGoalResponse struct {
	ID                uuid.UUID  `json:"id"`
	ClientID          uuid.UUID  `json:"client_id"`
	Title             string     `json:"title"`
	Description       *string    `json:"description"`
	Priority          string     `json:"priority"`
	Status            string     `json:"status"`
	TopicID           *uuid.UUID `json:"topic_id"`
	TopicNameSnapshot *string    `json:"topic_name_snapshot"`
	Source            string     `json:"source"`
	SortOrder         int32      `json:"sort_order"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

func toCreateClientGoalParams(req createClientGoalRequest) domain.CreateClientGoalParams {
	return domain.CreateClientGoalParams{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		TopicID:     req.TopicID,
		SortOrder:   req.SortOrder,
	}
}

func toCreateClientGoalResponse(goal domain.ClientGoal) createClientGoalResponse {
	return createClientGoalResponse{
		ID:                goal.ID,
		ClientID:          goal.ClientID,
		Title:             goal.Title,
		Description:       goal.Description,
		Priority:          goal.Priority,
		Status:            goal.Status,
		TopicID:           goal.TopicID,
		TopicNameSnapshot: goal.TopicNameSnapshot,
		Source:            goal.Source,
		SortOrder:         goal.SortOrder,
		CreatedAt:         goal.CreatedAt,
		UpdatedAt:         goal.UpdatedAt,
	}
}

// --- Update Client Goal ---

type updateClientGoalRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Priority    *string    `json:"priority" binding:"omitempty,oneof=low medium high"`
	TopicID     *uuid.UUID `json:"topic_id"`
	SortOrder   *int32     `json:"sort_order"`
}

type updateClientGoalResponse struct {
	MutationType      string                   `json:"mutation_type"`
	GoalID            uuid.UUID                `json:"goal_id"`
	ReplacementGoalID *uuid.UUID               `json:"replacement_goal_id,omitempty"`
	Goal              createClientGoalResponse `json:"goal"`
}

func toUpdateClientGoalParams(req updateClientGoalRequest) domain.UpdateClientGoalParams {
	return domain.UpdateClientGoalParams{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		TopicID:     req.TopicID,
		SortOrder:   req.SortOrder,
	}
}

func toUpdateClientGoalResponse(result *domain.UpdateClientGoalResult) updateClientGoalResponse {
	return updateClientGoalResponse{
		MutationType:      result.MutationType,
		GoalID:            result.GoalID,
		ReplacementGoalID: result.ReplacementGoalID,
		Goal:              toCreateClientGoalResponse(result.Goal),
	}
}

// --- Get Client Goals For Evaluation Page ---

type clientGoalForEvaluationPageResponse struct {
	ID                     uuid.UUID `json:"id"`
	TopicName              *string   `json:"topic_name"`
	Title                  string    `json:"title"`
	Priority               string    `json:"priority"`
	LastEvaluationProgress *string   `json:"last_evaluation_progress"`
}

type getClientGoalsForEvaluationPageResponse struct {
	NextEvaluationDate    *time.Time                            `json:"next_evaluation_date"`
	MyDraftEvaluationID   *uuid.UUID                            `json:"my_draft_evaluation_id"`
	IsResponsibleEmployee bool                                  `json:"is_responsible_employee"`
	CanUpdateGoals        bool                                  `json:"can_update_goals"`
	GoalUpdateBlockReason *string                               `json:"goal_update_block_reason"`
	Goals                 []clientGoalForEvaluationPageResponse `json:"goals"`
}

func toGetClientGoalsForEvaluationPageResponse(result *domain.ClientGoalsForEvaluationPage) getClientGoalsForEvaluationPageResponse {
	goals := make([]clientGoalForEvaluationPageResponse, 0, len(result.Goals))
	for _, g := range result.Goals {
		goals = append(goals, clientGoalForEvaluationPageResponse{
			ID:                     g.ID,
			TopicName:              g.TopicName,
			Title:                  g.Title,
			Priority:               g.Priority,
			LastEvaluationProgress: g.LastEvaluationProgress,
		})
	}
	return getClientGoalsForEvaluationPageResponse{
		NextEvaluationDate:    result.NextEvaluationDate,
		MyDraftEvaluationID:   result.MyDraftEvaluationID,
		IsResponsibleEmployee: result.IsResponsibleEmployee,
		CanUpdateGoals:        result.CanUpdateGoals,
		GoalUpdateBlockReason: result.GoalUpdateBlockReason,
		Goals:                 goals,
	}
}

// --- Create Goal Evaluation ---

type createGoalEvaluationItemRequest struct {
	GoalID   uuid.UUID `json:"goal_id" binding:"required"`
	Progress string    `json:"progress"`
	Notes    *string   `json:"notes"`
}

type createGoalEvaluationRequest struct {
	OverallNotes *string                           `json:"overall_notes"`
	Submit       bool                              `json:"submit"`
	Items        []createGoalEvaluationItemRequest `json:"items"`
}

type goalEvaluationItemResponse struct {
	ID                uuid.UUID `json:"id"`
	EvaluationID      uuid.UUID `json:"evaluation_id"`
	GoalID            uuid.UUID `json:"goal_id"`
	GoalTitle         string    `json:"goal_title"`
	GoalDescription   *string   `json:"goal_description"`
	TopicNameSnapshot *string   `json:"topic_name_snapshot"`
	Progress          string    `json:"progress"`
	Notes             *string   `json:"notes"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type goalEvaluationResponse struct {
	ID                      uuid.UUID                    `json:"id"`
	ClientID                uuid.UUID                    `json:"client_id"`
	EvaluationDate          time.Time                    `json:"evaluation_date"`
	PeriodStart             *time.Time                   `json:"period_start"`
	PeriodEnd               *time.Time                   `json:"period_end"`
	EvaluationIntervalWeeks int32                        `json:"evaluation_interval_weeks"`
	Status                  string                       `json:"status"`
	OverallNotes            *string                      `json:"overall_notes"`
	CreatedByEmployeeID     *uuid.UUID                   `json:"created_by_employee_id"`
	CreatorName             *string                      `json:"creator_name"`
	CreatedAt               time.Time                    `json:"created_at"`
	UpdatedAt               time.Time                    `json:"updated_at"`
	SubmitError             *string                      `json:"submit_error,omitempty"`
	Items                   []goalEvaluationItemResponse `json:"items"`
}

func toCreateGoalEvaluationParams(req createGoalEvaluationRequest) domain.CreateGoalEvaluationParams {
	items := make([]domain.GoalEvaluationItemParams, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, domain.GoalEvaluationItemParams{
			GoalID:   item.GoalID,
			Progress: item.Progress,
			Notes:    item.Notes,
		})
	}
	return domain.CreateGoalEvaluationParams{
		OverallNotes: req.OverallNotes,
		Submit:       req.Submit,
		Items:        items,
	}
}

func toGoalEvaluationResponse(eval domain.GoalEvaluation) goalEvaluationResponse {
	items := make([]goalEvaluationItemResponse, 0, len(eval.Items))
	for _, item := range eval.Items {
		items = append(items, goalEvaluationItemResponse{
			ID:                item.ID,
			EvaluationID:      item.EvaluationID,
			GoalID:            item.GoalID,
			GoalTitle:         item.GoalTitle,
			GoalDescription:   item.GoalDescription,
			TopicNameSnapshot: item.TopicNameSnapshot,
			Progress:          item.Progress,
			Notes:             item.Notes,
			CreatedAt:         item.CreatedAt,
			UpdatedAt:         item.UpdatedAt,
		})
	}
	return goalEvaluationResponse{
		ID:                      eval.ID,
		ClientID:                eval.ClientID,
		EvaluationDate:          eval.EvaluationDate,
		PeriodStart:             eval.PeriodStart,
		PeriodEnd:               eval.PeriodEnd,
		EvaluationIntervalWeeks: eval.EvaluationIntervalWeeks,
		Status:                  eval.Status,
		OverallNotes:            eval.OverallNotes,
		CreatedByEmployeeID:     eval.CreatedByEmployeeID,
		CreatorName:             eval.CreatorName,
		CreatedAt:               eval.CreatedAt,
		UpdatedAt:               eval.UpdatedAt,
		SubmitError:             eval.SubmitError,
		Items:                   items,
	}
}

// --- Get Goal Evaluation Bootstrap ---

type goalEvaluationBootstrapDraftResponse struct {
	ID             uuid.UUID `json:"id"`
	EvaluationDate time.Time `json:"evaluation_date"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type goalEvaluationBootstrapCompletedResponse struct {
	ID                  uuid.UUID  `json:"id"`
	EvaluationDate      time.Time  `json:"evaluation_date"`
	SubmittedAt         time.Time  `json:"submitted_at"`
	OverallNotes        *string    `json:"overall_notes"`
	CreatedByEmployeeID *uuid.UUID `json:"created_by_employee_id"`
	CreatorName         *string    `json:"creator_name"`
}

type goalEvaluationBootstrapActiveGoalResponse struct {
	GoalID            uuid.UUID `json:"goal_id"`
	Title             string    `json:"title"`
	TopicNameSnapshot *string   `json:"topic_name_snapshot"`
	Priority          string    `json:"priority"`
	SortOrder         int32     `json:"sort_order"`
	LastProgress      *string   `json:"last_progress"`
	LastNotes         *string   `json:"last_notes"`
}

type getGoalEvaluationBootstrapResponse struct {
	ClientID                uuid.UUID                                   `json:"client_id"`
	ClientFirstName         string                                      `json:"client_first_name"`
	ClientLastName          string                                      `json:"client_last_name"`
	NextEvaluationDate      *time.Time                                  `json:"next_evaluation_date"`
	DaysLeft                *int32                                      `json:"days_left"`
	Priority                *string                                     `json:"priority"`
	ExistingDraft           *goalEvaluationBootstrapDraftResponse       `json:"existing_draft"`
	LastCompletedEvaluation *goalEvaluationBootstrapCompletedResponse   `json:"last_completed_evaluation"`
	ActiveGoals             []goalEvaluationBootstrapActiveGoalResponse `json:"active_goals"`
}

func toGetGoalEvaluationBootstrapResponse(result *domain.GoalEvaluationBootstrap) getGoalEvaluationBootstrapResponse {
	var draft *goalEvaluationBootstrapDraftResponse
	if result.ExistingDraft != nil {
		draft = &goalEvaluationBootstrapDraftResponse{
			ID:             result.ExistingDraft.ID,
			EvaluationDate: result.ExistingDraft.EvaluationDate,
			UpdatedAt:      result.ExistingDraft.UpdatedAt,
		}
	}

	var completed *goalEvaluationBootstrapCompletedResponse
	if result.LastCompletedEvaluation != nil {
		completed = &goalEvaluationBootstrapCompletedResponse{
			ID:                  result.LastCompletedEvaluation.ID,
			EvaluationDate:      result.LastCompletedEvaluation.EvaluationDate,
			SubmittedAt:         result.LastCompletedEvaluation.SubmittedAt,
			OverallNotes:        result.LastCompletedEvaluation.OverallNotes,
			CreatedByEmployeeID: result.LastCompletedEvaluation.CreatedByEmployeeID,
			CreatorName:         result.LastCompletedEvaluation.CreatorName,
		}
	}

	goals := make([]goalEvaluationBootstrapActiveGoalResponse, 0, len(result.ActiveGoals))
	for _, g := range result.ActiveGoals {
		goals = append(goals, goalEvaluationBootstrapActiveGoalResponse{
			GoalID:            g.GoalID,
			Title:             g.Title,
			TopicNameSnapshot: g.TopicNameSnapshot,
			Priority:          g.Priority,
			SortOrder:         g.SortOrder,
			LastProgress:      g.LastProgress,
			LastNotes:         g.LastNotes,
		})
	}

	return getGoalEvaluationBootstrapResponse{
		ClientID:                result.ClientID,
		ClientFirstName:         result.ClientFirstName,
		ClientLastName:          result.ClientLastName,
		NextEvaluationDate:      result.NextEvaluationDate,
		DaysLeft:                result.DaysLeft,
		Priority:                result.Priority,
		ExistingDraft:           draft,
		LastCompletedEvaluation: completed,
		ActiveGoals:             goals,
	}
}

// --- List Client Submitted Evaluations ---

type listClientSubmittedEvaluationsRequest struct {
	httpapi.PageRequest
}

type listClientSubmittedEvaluationsResponse struct {
	EvaluationID        uuid.UUID  `json:"evaluation_id"`
	EvaluationDate      time.Time  `json:"evaluation_date"`
	SubmittedAt         time.Time  `json:"submitted_at"`
	FilledGoalsCount    int32      `json:"filled_goals_count"`
	TotalGoalsCount     int32      `json:"total_goals_count"`
	CreatedByEmployeeID *uuid.UUID `json:"created_by_employee_id"`
	CreatorName         *string    `json:"creator_name"`
}

func toListClientSubmittedEvaluationsResponse(item domain.ListClientSubmittedEvaluationsItem) listClientSubmittedEvaluationsResponse {
	return listClientSubmittedEvaluationsResponse{
		EvaluationID:        item.EvaluationID,
		EvaluationDate:      item.EvaluationDate,
		SubmittedAt:         item.SubmittedAt,
		FilledGoalsCount:    item.FilledGoalsCount,
		TotalGoalsCount:     item.TotalGoalsCount,
		CreatedByEmployeeID: item.CreatedByEmployeeID,
		CreatorName:         item.CreatorName,
	}
}

// --- List Goal Evaluation History ---

type listGoalEvaluationHistoryRequest struct {
	httpapi.PageRequest
}

type listGoalEvaluationHistoryResponse struct {
	EvaluationID        uuid.UUID  `json:"evaluation_id"`
	EvaluationDate      time.Time  `json:"evaluation_date"`
	SubmittedAt         time.Time  `json:"submitted_at"`
	Progress            string     `json:"progress"`
	Notes               *string    `json:"notes"`
	CreatedByEmployeeID *uuid.UUID `json:"created_by_employee_id"`
	CreatorName         *string    `json:"creator_name"`
	PeriodStart         *time.Time `json:"period_start"`
	PeriodEnd           *time.Time `json:"period_end"`
}

func toListGoalEvaluationHistoryResponse(item domain.ListGoalEvaluationHistoryItem) listGoalEvaluationHistoryResponse {
	return listGoalEvaluationHistoryResponse{
		EvaluationID:        item.EvaluationID,
		EvaluationDate:      item.EvaluationDate,
		SubmittedAt:         item.SubmittedAt,
		Progress:            item.Progress,
		Notes:               item.Notes,
		CreatedByEmployeeID: item.CreatedByEmployeeID,
		CreatorName:         item.CreatorName,
		PeriodStart:         item.PeriodStart,
		PeriodEnd:           item.PeriodEnd,
	}
}

// --- Location Transfer ---

type requestLocationTransferRequest struct {
	FromLocationID uuid.UUID  `json:"from_location_id" binding:"required"`
	ToLocationID   uuid.UUID  `json:"to_location_id" binding:"required"`
	Reason         string     `json:"reason" binding:"required"`
	NewMentorID    *uuid.UUID `json:"new_mentor_id"`
}

type approveOrRejectLocationTransferRequest struct {
	TransferID uuid.UUID `json:"transfer_id" binding:"required"`
	Status     string    `json:"status" binding:"required,oneof=approved rejected"`
}

type listLocationTransferRequestsRequest struct {
	httpapi.PageRequest
}

type locationTransferResponse struct {
	ID                 uuid.UUID  `json:"id"`
	ClientID           uuid.UUID  `json:"client_id"`
	FromLocationID     *uuid.UUID `json:"from_location_id"`
	ToLocationID       *uuid.UUID `json:"to_location_id"`
	NewMentorID        *uuid.UUID `json:"new_mentor_id"`
	RequestDate        time.Time  `json:"request_date"`
	Status             string     `json:"status"`
	ApprovedRejectedBy *uuid.UUID `json:"approved_rejected_by"`
	ApprovedRejectedAt *time.Time `json:"approved_rejected_at"`
	Reason             *string    `json:"reason"`
	MentorFirstName    *string    `json:"mentor_first_name"`
	MentorLastName     *string    `json:"mentor_last_name"`
}

func toCreateLocationTransferParams(req requestLocationTransferRequest) domain.CreateLocationTransferParams {
	return domain.CreateLocationTransferParams{
		FromLocationID: req.FromLocationID,
		ToLocationID:   req.ToLocationID,
		Reason:         req.Reason,
		NewMentorID:    req.NewMentorID,
	}
}

func toApproveLocationTransferParams(req approveOrRejectLocationTransferRequest) domain.ApproveLocationTransferParams {
	return domain.ApproveLocationTransferParams{
		TransferID: req.TransferID,
		Status:     req.Status,
	}
}

func toLocationTransferResponse(t domain.LocationTransfer) locationTransferResponse {
	return locationTransferResponse{
		ID:                 t.ID,
		ClientID:           t.ClientID,
		FromLocationID:     t.FromLocationID,
		ToLocationID:       t.ToLocationID,
		NewMentorID:        t.NewMentorID,
		RequestDate:        t.RequestDate,
		Status:             t.Status,
		ApprovedRejectedBy: t.ApprovedRejectedBy,
		ApprovedRejectedAt: t.ApprovedRejectedAt,
		Reason:             t.Reason,
		MentorFirstName:    t.MentorFirstName,
		MentorLastName:     t.MentorLastName,
	}
}

// --- Top-level Evaluations ---

type listUpcomingEvaluationsRequest struct {
	httpapi.PageRequest
}

type upcomingEvaluationResponse struct {
	ClientID         uuid.UUID `json:"client_id"`
	ClientFirstName  string    `json:"client_first_name"`
	ClientLastName   string    `json:"client_last_name"`
	DueDate          time.Time `json:"due_date"`
	DaysLeft         int32     `json:"days_left"`
	Priority         string    `json:"priority"`
	HasDraft         bool      `json:"has_draft"`
	FilledGoalsCount int32     `json:"filled_goals_count"`
	TotalGoalsCount  int32     `json:"total_goals_count"`
}

type listRecentSubmittedEvaluationsRequest struct {
	httpapi.PageRequest
}

type recentSubmittedEvaluationResponse struct {
	EvaluationID       uuid.UUID  `json:"evaluation_id"`
	ClientID           uuid.UUID  `json:"client_id"`
	ClientFirstName    string     `json:"client_first_name"`
	ClientLastName     string     `json:"client_last_name"`
	EvaluationDate     time.Time  `json:"evaluation_date"`
	SubmittedAt        time.Time  `json:"submitted_at"`
	NextEvaluationDate *time.Time `json:"next_evaluation_date"`
	FilledGoalsCount   int32      `json:"filled_goals_count"`
	TotalGoalsCount    int32      `json:"total_goals_count"`
}

type listRecentDraftEvaluationsRequest struct {
	httpapi.PageRequest
}

type recentDraftEvaluationResponse struct {
	EvaluationID     uuid.UUID `json:"evaluation_id"`
	ClientID         uuid.UUID `json:"client_id"`
	ClientFirstName  string    `json:"client_first_name"`
	ClientLastName   string    `json:"client_last_name"`
	DueDate          time.Time `json:"due_date"`
	UpdatedAt        time.Time `json:"updated_at"`
	DaysLeft         int32     `json:"days_left"`
	Priority         string    `json:"priority"`
	FilledGoalsCount int32     `json:"filled_goals_count"`
	TotalGoalsCount  int32     `json:"total_goals_count"`
}

func toUpcomingEvaluationResponse(e domain.UpcomingEvaluation) upcomingEvaluationResponse {
	return upcomingEvaluationResponse{
		ClientID:         e.ClientID,
		ClientFirstName:  e.ClientFirstName,
		ClientLastName:   e.ClientLastName,
		DueDate:          e.DueDate,
		DaysLeft:         e.DaysLeft,
		Priority:         e.Priority,
		HasDraft:         e.HasDraft,
		FilledGoalsCount: e.FilledGoalsCount,
		TotalGoalsCount:  e.TotalGoalsCount,
	}
}

func toRecentSubmittedEvaluationResponse(e domain.RecentSubmittedEvaluation) recentSubmittedEvaluationResponse {
	return recentSubmittedEvaluationResponse{
		EvaluationID:       e.EvaluationID,
		ClientID:           e.ClientID,
		ClientFirstName:    e.ClientFirstName,
		ClientLastName:     e.ClientLastName,
		EvaluationDate:     e.EvaluationDate,
		SubmittedAt:        e.SubmittedAt,
		NextEvaluationDate: e.NextEvaluationDate,
		FilledGoalsCount:   e.FilledGoalsCount,
		TotalGoalsCount:    e.TotalGoalsCount,
	}
}

func toRecentDraftEvaluationResponse(e domain.RecentDraftEvaluation) recentDraftEvaluationResponse {
	return recentDraftEvaluationResponse{
		EvaluationID:     e.EvaluationID,
		ClientID:         e.ClientID,
		ClientFirstName:  e.ClientFirstName,
		ClientLastName:   e.ClientLastName,
		DueDate:          e.DueDate,
		UpdatedAt:        e.UpdatedAt,
		DaysLeft:         e.DaysLeft,
		Priority:         e.Priority,
		FilledGoalsCount: e.FilledGoalsCount,
		TotalGoalsCount:  e.TotalGoalsCount,
	}
}

// --- Medical - Diagnoses ---

type createClientDiagnosisRequest struct {
	CodeSystem          string     `json:"code_system"`
	Code                string     `json:"code"`
	Title               *string    `json:"title"`
	Description         *string    `json:"description"`
	Status              *string    `json:"status"`
	Severity            *string    `json:"severity"`
	DiagnosedOn         *time.Time `json:"diagnosed_on"`
	ResolvedOn          *time.Time `json:"resolved_on"`
	DiagnosingClinician *string    `json:"diagnosing_clinician"`
	Notes               *string    `json:"notes"`
}

type clientDiagnosisResponse struct {
	ID                  uuid.UUID  `json:"id"`
	ClientID            uuid.UUID  `json:"client_id"`
	CodeSystem          string     `json:"code_system"`
	Code                string     `json:"code"`
	Title               *string    `json:"title"`
	Description         *string    `json:"description"`
	Status              string     `json:"status"`
	Severity            string     `json:"severity"`
	DiagnosedOn         *time.Time `json:"diagnosed_on"`
	ResolvedOn          *time.Time `json:"resolved_on"`
	DiagnosingClinician *string    `json:"diagnosing_clinician"`
	Notes               *string    `json:"notes"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type listClientDiagnosesRequest struct {
	httpapi.PageRequest
}

type updateClientDiagnosisRequest struct {
	CodeSystem          *string    `json:"code_system"`
	Code                *string    `json:"code"`
	Title               *string    `json:"title"`
	Description         *string    `json:"description"`
	Status              *string    `json:"status"`
	Severity            *string    `json:"severity"`
	DiagnosedOn         *time.Time `json:"diagnosed_on"`
	ResolvedOn          *time.Time `json:"resolved_on"`
	DiagnosingClinician *string    `json:"diagnosing_clinician"`
	Notes               *string    `json:"notes"`
}

type deleteClientDiagnosisResponse struct {
	ID uuid.UUID `json:"id"`
}

func toDomainCreateClientDiagnosisParams(req createClientDiagnosisRequest, clientID uuid.UUID) domain.CreateClientDiagnosisParams {
	return domain.CreateClientDiagnosisParams{
		ClientID:            clientID,
		CodeSystem:          req.CodeSystem,
		Code:                req.Code,
		Title:               req.Title,
		Description:         req.Description,
		Status:              req.Status,
		Severity:            req.Severity,
		DiagnosedOn:         req.DiagnosedOn,
		ResolvedOn:          req.ResolvedOn,
		DiagnosingClinician: req.DiagnosingClinician,
		Notes:               req.Notes,
	}
}

func toDomainUpdateClientDiagnosisParams(req updateClientDiagnosisRequest, clientID, diagnosisID uuid.UUID) domain.UpdateClientDiagnosisParams {
	return domain.UpdateClientDiagnosisParams{
		ClientID:            clientID,
		ID:                  diagnosisID,
		CodeSystem:          req.CodeSystem,
		Code:                req.Code,
		Title:               req.Title,
		Description:         req.Description,
		Status:              req.Status,
		Severity:            req.Severity,
		DiagnosedOn:         req.DiagnosedOn,
		ResolvedOn:          req.ResolvedOn,
		DiagnosingClinician: req.DiagnosingClinician,
		Notes:               req.Notes,
	}
}

func toClientDiagnosisResponse(d domain.ClientDiagnosis) clientDiagnosisResponse {
	return clientDiagnosisResponse{
		ID:                  d.ID,
		ClientID:            d.ClientID,
		CodeSystem:          d.CodeSystem,
		Code:                d.Code,
		Title:               d.Title,
		Description:         d.Description,
		Status:              d.Status,
		Severity:            d.Severity,
		DiagnosedOn:         d.DiagnosedOn,
		ResolvedOn:          d.ResolvedOn,
		DiagnosingClinician: d.DiagnosingClinician,
		Notes:               d.Notes,
		CreatedAt:           d.CreatedAt,
		UpdatedAt:           d.UpdatedAt,
	}
}

// --- Medical - Medication Orders ---

type createClientMedicationOrderRequest struct {
	DiagnosisID           *uuid.UUID      `json:"diagnosis_id"`
	MedicationName        string          `json:"medication_name"`
	DosageText            string          `json:"dosage_text"`
	DoseAmount            *float64        `json:"dose_amount"`
	DoseUnit              *string         `json:"dose_unit"`
	Route                 *string         `json:"route"`
	FrequencyText         *string         `json:"frequency_text"`
	Schedule              json.RawMessage `json:"schedule"`
	IsPrn                 bool            `json:"is_prn"`
	PrnIndication         *string         `json:"prn_indication"`
	MaxDosesPer24h        *int32          `json:"max_doses_per_24h"`
	StartDate             time.Time       `json:"start_date"`
	EndDate               *time.Time      `json:"end_date"`
	Status                *string         `json:"status"`
	AdminMode             *string         `json:"admin_mode"`
	ResponsibleEmployeeID *uuid.UUID      `json:"responsible_employee_id"`
	IsCritical            bool            `json:"is_critical"`
	Notes                 *string         `json:"notes"`
	SourceAttachmentUUID  *uuid.UUID      `json:"source_attachment_uuid"`
}

type clientMedicationOrderResponse struct {
	ID                           uuid.UUID       `json:"id"`
	ClientID                     uuid.UUID       `json:"client_id"`
	DiagnosisID                  *uuid.UUID      `json:"diagnosis_id"`
	MedicationName               string          `json:"medication_name"`
	DosageText                   string          `json:"dosage_text"`
	DoseAmount                   *float64        `json:"dose_amount"`
	DoseUnit                     *string         `json:"dose_unit"`
	Route                        *string         `json:"route"`
	FrequencyText                *string         `json:"frequency_text"`
	Schedule                     json.RawMessage `json:"schedule"`
	IsPrn                        bool            `json:"is_prn"`
	PrnIndication                *string         `json:"prn_indication"`
	MaxDosesPer24h               *int32          `json:"max_doses_per_24h"`
	StartDate                    time.Time       `json:"start_date"`
	EndDate                      *time.Time      `json:"end_date"`
	Status                       string          `json:"status"`
	AdminMode                    string          `json:"admin_mode"`
	ResponsibleEmployeeID        *uuid.UUID      `json:"responsible_employee_id"`
	ResponsibleEmployeeFirstName *string         `json:"responsible_employee_first_name"`
	ResponsibleEmployeeLastName  *string         `json:"responsible_employee_last_name"`
	IsCritical                   bool            `json:"is_critical"`
	Notes                        *string         `json:"notes"`
	SourceAttachmentUUID         *uuid.UUID      `json:"source_attachment_uuid"`
	DiagnosisTitle               *string         `json:"diagnosis_title"`
	DiagnosisCodeSystem          *string         `json:"diagnosis_code_system"`
	DiagnosisCode                *string         `json:"diagnosis_code"`
	CreatedAt                    time.Time       `json:"created_at"`
	UpdatedAt                    time.Time       `json:"updated_at"`
}

type listClientMedicationOrdersRequest struct {
	httpapi.PageRequest
	Status      *string    `json:"status" form:"status"`
	AdminMode   *string    `json:"admin_mode" form:"admin_mode"`
	DiagnosisID *uuid.UUID `json:"diagnosis_id" form:"diagnosis_id"`
	Search      *string    `json:"search" form:"search"`
}

type updateClientMedicationOrderRequest struct {
	DiagnosisID           *uuid.UUID      `json:"diagnosis_id"`
	MedicationName        *string         `json:"medication_name"`
	DosageText            *string         `json:"dosage_text"`
	DoseAmount            *float64        `json:"dose_amount"`
	DoseUnit              *string         `json:"dose_unit"`
	Route                 *string         `json:"route"`
	FrequencyText         *string         `json:"frequency_text"`
	Schedule              json.RawMessage `json:"schedule"`
	IsPrn                 *bool           `json:"is_prn"`
	PrnIndication         *string         `json:"prn_indication"`
	MaxDosesPer24h        *int32          `json:"max_doses_per_24h"`
	StartDate             *time.Time      `json:"start_date"`
	EndDate               *time.Time      `json:"end_date"`
	Status                *string         `json:"status"`
	AdminMode             *string         `json:"admin_mode"`
	ResponsibleEmployeeID *uuid.UUID      `json:"responsible_employee_id"`
	IsCritical            *bool           `json:"is_critical"`
	Notes                 *string         `json:"notes"`
	SourceAttachmentUUID  *uuid.UUID      `json:"source_attachment_uuid"`
}

type deleteClientMedicationOrderResponse struct {
	ID uuid.UUID `json:"id"`
}

func toDomainCreateClientMedicationOrderParams(req createClientMedicationOrderRequest, clientID uuid.UUID) domain.CreateClientMedicationOrderParams {
	return domain.CreateClientMedicationOrderParams{
		ClientID:              clientID,
		DiagnosisID:           req.DiagnosisID,
		MedicationName:        req.MedicationName,
		DosageText:            req.DosageText,
		DoseAmount:            req.DoseAmount,
		DoseUnit:              req.DoseUnit,
		Route:                 req.Route,
		FrequencyText:         req.FrequencyText,
		Schedule:              req.Schedule,
		IsPrn:                 req.IsPrn,
		PrnIndication:         req.PrnIndication,
		MaxDosesPer24h:        req.MaxDosesPer24h,
		StartDate:             req.StartDate,
		EndDate:               req.EndDate,
		Status:                req.Status,
		AdminMode:             req.AdminMode,
		ResponsibleEmployeeID: req.ResponsibleEmployeeID,
		IsCritical:            req.IsCritical,
		Notes:                 req.Notes,
		SourceAttachmentUUID:  req.SourceAttachmentUUID,
	}
}

func toDomainUpdateClientMedicationOrderParams(req updateClientMedicationOrderRequest, clientID, orderID uuid.UUID) domain.UpdateClientMedicationOrderParams {
	return domain.UpdateClientMedicationOrderParams{
		ClientID:              clientID,
		ID:                    orderID,
		DiagnosisID:           req.DiagnosisID,
		MedicationName:        req.MedicationName,
		DosageText:            req.DosageText,
		DoseAmount:            req.DoseAmount,
		DoseUnit:              req.DoseUnit,
		Route:                 req.Route,
		FrequencyText:         req.FrequencyText,
		Schedule:              req.Schedule,
		IsPrn:                 req.IsPrn,
		PrnIndication:         req.PrnIndication,
		MaxDosesPer24h:        req.MaxDosesPer24h,
		StartDate:             req.StartDate,
		EndDate:               req.EndDate,
		Status:                req.Status,
		AdminMode:             req.AdminMode,
		ResponsibleEmployeeID: req.ResponsibleEmployeeID,
		IsCritical:            req.IsCritical,
		Notes:                 req.Notes,
		SourceAttachmentUUID:  req.SourceAttachmentUUID,
	}
}

func toClientMedicationOrderResponse(o domain.ClientMedicationOrder) clientMedicationOrderResponse {
	return clientMedicationOrderResponse{
		ID:                           o.ID,
		ClientID:                     o.ClientID,
		DiagnosisID:                  o.DiagnosisID,
		MedicationName:               o.MedicationName,
		DosageText:                   o.DosageText,
		DoseAmount:                   o.DoseAmount,
		DoseUnit:                     o.DoseUnit,
		Route:                        o.Route,
		FrequencyText:                o.FrequencyText,
		Schedule:                     json.RawMessage(o.Schedule),
		IsPrn:                        o.IsPrn,
		PrnIndication:                o.PrnIndication,
		MaxDosesPer24h:               o.MaxDosesPer24h,
		StartDate:                    o.StartDate,
		EndDate:                      o.EndDate,
		Status:                       o.Status,
		AdminMode:                    o.AdminMode,
		ResponsibleEmployeeID:        o.ResponsibleEmployeeID,
		ResponsibleEmployeeFirstName: o.ResponsibleEmployeeFirstName,
		ResponsibleEmployeeLastName:  o.ResponsibleEmployeeLastName,
		IsCritical:                   o.IsCritical,
		Notes:                        o.Notes,
		SourceAttachmentUUID:         o.SourceAttachmentUUID,
		DiagnosisTitle:               o.DiagnosisTitle,
		DiagnosisCodeSystem:          o.DiagnosisCodeSystem,
		DiagnosisCode:                o.DiagnosisCode,
		CreatedAt:                    o.CreatedAt,
		UpdatedAt:                    o.UpdatedAt,
	}
}

// --- Medical - Overview ---

type clientMedicalOverviewResponse struct {
	Diagnoses        []clientDiagnosisResponse       `json:"diagnoses"`
	MedicationOrders []clientMedicationOrderResponse `json:"medication_orders"`
}

// =====================
// Network - Sender
// =====================

type senderContact struct {
	Name        *string `json:"name"`
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phone_number"`
}

type getClientSenderResponse struct {
	ID                  uuid.UUID       `json:"id"`
	Types               string          `json:"types"`
	Name                string          `json:"name"`
	Street              *string         `json:"street"`
	HouseNumber         *string         `json:"house_number"`
	HouseNumberAddition *string         `json:"house_number_addition"`
	PostalCode          *string         `json:"postal_code"`
	City                *string         `json:"city"`
	Land                *string         `json:"land"`
	Kvknumber           *string         `json:"kvknumber"`
	Btwnumber           *string         `json:"btwnumber"`
	PhoneNumber         *string         `json:"phone_number"`
	ClientNumber        *string         `json:"client_number"`
	EmailAddress        *string         `json:"email_address"`
	Contacts            []senderContact `json:"contacts"`
	IsArchived          bool            `json:"is_archived"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
}

func toGetClientSenderResponse(s domain.Sender) getClientSenderResponse {
	var contacts []senderContact
	if len(s.Contacts) > 0 {
		_ = json.Unmarshal(s.Contacts, &contacts)
	}
	return getClientSenderResponse{
		ID:                  s.ID,
		Types:               s.Types,
		Name:                s.Name,
		Street:              s.Street,
		HouseNumber:         s.HouseNumber,
		HouseNumberAddition: s.HouseNumberAddition,
		PostalCode:          s.PostalCode,
		City:                s.City,
		Land:                s.Land,
		Kvknumber:           s.Kvknumber,
		Btwnumber:           s.Btwnumber,
		PhoneNumber:         s.PhoneNumber,
		ClientNumber:        s.ClientNumber,
		EmailAddress:        s.EmailAddress,
		Contacts:            contacts,
		IsArchived:          s.IsArchived,
		CreatedAt:           s.CreatedAt,
		UpdatedAt:           s.UpdatedAt,
	}
}

// =====================
// Network - Emergency Contacts
// =====================

type createClientEmergencyContactRequest struct {
	FirstName        *string `json:"first_name"`
	LastName         *string `json:"last_name"`
	Email            *string `json:"email"`
	PhoneNumber      *string `json:"phone_number"`
	Address          *string `json:"address"`
	Relationship     *string `json:"relationship"`
	RelationStatus   *string `json:"relation_status"`
	MedicalReports   bool    `json:"medical_reports"`
	IncidentsReports bool    `json:"incidents_reports"`
	GoalsReports     bool    `json:"goals_reports"`
}

type clientEmergencyContactResponse struct {
	ID               uuid.UUID `json:"id"`
	ClientID         uuid.UUID `json:"client_id"`
	FirstName        *string   `json:"first_name"`
	LastName         *string   `json:"last_name"`
	Email            *string   `json:"email"`
	PhoneNumber      *string   `json:"phone_number"`
	Address          *string   `json:"address"`
	Relationship     *string   `json:"relationship"`
	RelationStatus   *string   `json:"relation_status"`
	CreatedAt        time.Time `json:"created_at"`
	IsVerified       bool      `json:"is_verified"`
	MedicalReports   bool      `json:"medical_reports"`
	IncidentsReports bool      `json:"incidents_reports"`
	GoalsReports     bool      `json:"goals_reports"`
}

type listClientEmergencyContactsRequest struct {
	httpapi.PageRequest
	Search string `form:"search"`
}

type updateClientEmergencyContactRequest struct {
	FirstName        *string `json:"first_name"`
	LastName         *string `json:"last_name"`
	Email            *string `json:"email"`
	PhoneNumber      *string `json:"phone_number"`
	Address          *string `json:"address"`
	Relationship     *string `json:"relationship"`
	RelationStatus   *string `json:"relation_status"`
	MedicalReports   *bool   `json:"medical_reports"`
	IncidentsReports *bool   `json:"incidents_reports"`
	GoalsReports     *bool   `json:"goals_reports"`
}

type deleteClientEmergencyContactResponse struct {
	ID uuid.UUID `json:"id"`
}

func toDomainCreateClientEmergencyContactParams(req createClientEmergencyContactRequest, clientID uuid.UUID) domain.CreateClientEmergencyContactParams {
	return domain.CreateClientEmergencyContactParams{
		ClientID:         clientID,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Email:            req.Email,
		PhoneNumber:      req.PhoneNumber,
		Address:          req.Address,
		Relationship:     req.Relationship,
		RelationStatus:   req.RelationStatus,
		MedicalReports:   req.MedicalReports,
		IncidentsReports: req.IncidentsReports,
		GoalsReports:     req.GoalsReports,
	}
}

func toDomainUpdateClientEmergencyContactParams(req updateClientEmergencyContactRequest, contactID uuid.UUID) domain.UpdateClientEmergencyContactParams {
	return domain.UpdateClientEmergencyContactParams{
		ID:               contactID,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Email:            req.Email,
		PhoneNumber:      req.PhoneNumber,
		Address:          req.Address,
		Relationship:     req.Relationship,
		RelationStatus:   req.RelationStatus,
		MedicalReports:   req.MedicalReports,
		IncidentsReports: req.IncidentsReports,
		GoalsReports:     req.GoalsReports,
	}
}

func toClientEmergencyContactResponse(c domain.ClientEmergencyContact) clientEmergencyContactResponse {
	return clientEmergencyContactResponse{
		ID:               c.ID,
		ClientID:         c.ClientID,
		FirstName:        c.FirstName,
		LastName:         c.LastName,
		Email:            c.Email,
		PhoneNumber:      c.PhoneNumber,
		Address:          c.Address,
		Relationship:     c.Relationship,
		RelationStatus:   c.RelationStatus,
		CreatedAt:        c.CreatedAt,
		IsVerified:       c.IsVerified,
		MedicalReports:   c.MedicalReports,
		IncidentsReports: c.IncidentsReports,
		GoalsReports:     c.GoalsReports,
	}
}

// =====================
// Network - Assigned Employees
// =====================

type createAssignedEmployeeRequest struct {
	EmployeeID uuid.UUID `json:"employee_id" binding:"required"`
	StartDate  time.Time `json:"start_date" binding:"required"`
	Role       string    `json:"role" binding:"required"`
}

type assignedEmployeeResponse struct {
	ID           uuid.UUID `json:"id"`
	ClientID     uuid.UUID `json:"client_id"`
	EmployeeID   uuid.UUID `json:"employee_id"`
	StartDate    time.Time `json:"start_date"`
	Role         string    `json:"role"`
	EmployeeName string    `json:"employee_name"`
	CreatedAt    time.Time `json:"created_at"`
}

type listAssignedEmployeesRequest struct {
	httpapi.PageRequest
}

type updateAssignedEmployeeRequest struct {
	EmployeeID *uuid.UUID `json:"employee_id"`
	StartDate  time.Time  `json:"start_date"`
	Role       *string    `json:"role"`
}

type deleteAssignedEmployeeResponse struct {
	ID uuid.UUID `json:"id"`
}

func toDomainCreateAssignedEmployeeParams(req createAssignedEmployeeRequest, clientID uuid.UUID) domain.CreateAssignedEmployeeParams {
	return domain.CreateAssignedEmployeeParams{
		ClientID:   clientID,
		EmployeeID: req.EmployeeID,
		StartDate:  req.StartDate,
		Role:       req.Role,
	}
}

func toDomainUpdateAssignedEmployeeParams(req updateAssignedEmployeeRequest, assignmentID uuid.UUID) domain.UpdateAssignedEmployeeParams {
	return domain.UpdateAssignedEmployeeParams{
		ID:         assignmentID,
		EmployeeID: req.EmployeeID,
		StartDate:  req.StartDate,
		Role:       req.Role,
	}
}

func toAssignedEmployeeResponse(e domain.AssignedEmployee) assignedEmployeeResponse {
	employeeName := e.EmployeeFirstName + " " + e.EmployeeLastName
	return assignedEmployeeResponse{
		ID:           e.ID,
		ClientID:     e.ClientID,
		EmployeeID:   e.EmployeeID,
		StartDate:    e.StartDate,
		Role:         e.Role,
		EmployeeName: employeeName,
		CreatedAt:    e.CreatedAt,
	}
}

// =====================
// Network - Related Emails
// =====================

type getClientRelatedEmailsResponse struct {
	Emails []*string `json:"emails"`
}

// --- Progress Reports ---

type createProgressReportRequest struct {
	EmployeeID     *uuid.UUID `json:"employee_id"`
	Title          *string    `json:"title"`
	Date           time.Time  `json:"date"`
	ReportText     string     `json:"report_text" binding:"required"`
	Type           string     `json:"type" binding:"required,oneof=morning_report evening_report night_report shift_report one_to_one_report process_report contact_journal other"`
	EmotionalState string     `json:"emotional_state" binding:"required,oneof=normal excited happy sad angry anxious depressed"`
}

type progressReportResponse struct {
	ID                     uuid.UUID  `json:"id"`
	ClientID               uuid.UUID  `json:"client_id"`
	Date                   time.Time  `json:"date"`
	Title                  *string    `json:"title"`
	ReportText             string     `json:"report_text"`
	EmployeeID             *uuid.UUID `json:"employee_id"`
	Type                   string     `json:"type"`
	EmotionalState         string     `json:"emotional_state"`
	CreatedAt              time.Time  `json:"created_at"`
	EmployeeFirstName      string     `json:"employee_first_name"`
	EmployeeLastName       string     `json:"employee_last_name"`
	EmployeeProfilePicture *string    `json:"employee_profile_picture"`
}

type listProgressReportsRequest struct {
	httpapi.PageRequest
	Type *string `form:"type" binding:"omitempty,oneof=morning_report evening_report night_report shift_report one_to_one_report process_report contact_journal other"`
}

type updateProgressReportRequest struct {
	EmployeeID     *uuid.UUID `json:"employee_id"`
	Title          *string    `json:"title"`
	Date           time.Time  `json:"date"`
	ReportText     *string    `json:"report_text"`
	Type           *string    `json:"type"`
	EmotionalState *string    `json:"emotional_state"`
}

func toDomainCreateProgressReportParams(req createProgressReportRequest, clientID uuid.UUID) domain.CreateProgressReportParams {
	return domain.CreateProgressReportParams{
		ClientID:       clientID,
		EmployeeID:     req.EmployeeID,
		Title:          req.Title,
		Date:           req.Date,
		ReportText:     req.ReportText,
		Type:           req.Type,
		EmotionalState: req.EmotionalState,
	}
}

func toDomainUpdateProgressReportParams(req updateProgressReportRequest, reportID uuid.UUID) domain.UpdateProgressReportParams {
	return domain.UpdateProgressReportParams{
		ID:             reportID,
		EmployeeID:     req.EmployeeID,
		Title:          req.Title,
		Date:           req.Date,
		ReportText:     req.ReportText,
		Type:           req.Type,
		EmotionalState: req.EmotionalState,
	}
}

func toProgressReportResponse(r domain.ProgressReport) progressReportResponse {
	return progressReportResponse{
		ID:                     r.ID,
		ClientID:               r.ClientID,
		Date:                   r.Date,
		Title:                  r.Title,
		ReportText:             r.ReportText,
		EmployeeID:             r.EmployeeID,
		Type:                   r.Type,
		EmotionalState:         r.EmotionalState,
		CreatedAt:              r.CreatedAt,
		EmployeeFirstName:      r.EmployeeFirstName,
		EmployeeLastName:       r.EmployeeLastName,
		EmployeeProfilePicture: r.EmployeeProfilePicture,
	}
}

// --- AI Progress Reports ---

type generateAutoReportsRequest struct {
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

type generateAutoReportsResponse struct {
	Report string `json:"report"`
}

type confirmProgressReportRequest struct {
	ReportText string    `json:"report_text" binding:"required"`
	StartDate  time.Time `json:"start_date" binding:"required"`
	EndDate    time.Time `json:"end_date" binding:"required"`
}

type aiGeneratedReportResponse struct {
	ID         uuid.UUID `json:"id"`
	ClientID   uuid.UUID `json:"client_id"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	ReportText string    `json:"report_text"`
	CreatedAt  time.Time `json:"created_at"`
}

type listAiGeneratedReportsRequest struct {
	httpapi.PageRequest
}

func toAiGeneratedReportResponse(r domain.AiGeneratedReport) aiGeneratedReportResponse {
	return aiGeneratedReportResponse{
		ID:         r.ID,
		ClientID:   r.ClientID,
		StartDate:  r.StartDate,
		EndDate:    r.EndDate,
		ReportText: r.ReportText,
		CreatedAt:  r.CreatedAt,
	}
}
