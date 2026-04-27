package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	db "maicare_go/db/sqlc"
)

var ErrClientNotFound = errors.New("client not found")

type Client struct {
	ID                         uuid.UUID
	FirstName                  string
	LastName                   string
	DateOfBirth                time.Time
	Identity                   bool
	Status                     string
	Bsn                        *string
	BsnVerifiedBy              *uuid.UUID
	Email                      string
	PhoneNumber                *string
	Gender                     string
	Filenumber                 string
	CreatedAt                  time.Time
	SenderID                   *uuid.UUID
	LocationID                 *uuid.UUID
	EducationCurrentlyEnrolled bool
	EducationInstitution       *string
	EducationMentorName        *string
	EducationMentorPhone       *string
	EducationMentorEmail       *string
	EducationAdditionalNotes   *string
	EducationLevel             string
	WorkCurrentlyEmployed      bool
	WorkCurrentEmployer        *string
	WorkCurrentEmployerPhone   *string
	WorkCurrentEmployerEmail   *string
	WorkCurrentPosition        *string
	WorkStartDate              time.Time
	WorkAdditionalNotes        *string
}

type ClientListItem struct {
	ID           uuid.UUID
	FirstName    string
	LastName     string
	Bsn          *string
	Filenumber   string
	LocationName *string
	CareType     *string
	Status       string
	GoalsCount   int64
	RiskCount    int64
	CreatedAt    time.Time
}

type WaitingListClient struct {
	ID             uuid.UUID
	FirstName      string
	LastName       string
	CareType       *string
	Bsn            *string
	SenderName     *string
	DaysInWaitlist int32
	AdmissionType  *string
}

type InCareClient struct {
	ID                uuid.UUID
	Bsn               *string
	FirstName         string
	LastName          string
	CoordinatorName   *string
	LocationName      *string
	Status            string
	CareStartDate     *time.Time
	DaysInCare        int32
	HasActiveContract bool
}

type ClientCounts struct {
	TotalClients         int64
	ClientsInCare        int64
	ClientsOnWaitingList int64
	ClientsOutOfCare     int64
}

type ClientStatusCounts struct {
	ClientsInOrScheduledInCare     int64
	ClientsOnWaitingList           int64
	ClientsOutOrScheduledOutOfCare int64
}

type ClientPage struct {
	Items      []ClientListItem
	TotalCount int64
}

type WaitingListClientPage struct {
	Items      []WaitingListClient
	TotalCount int64
}

type InCareClientPage struct {
	Items      []InCareClient
	TotalCount int64
}

type CreateClientParams struct {
	FirstName                  string
	LastName                   string
	DateOfBirth                time.Time
	Bsn                        *string
	BsnVerifiedBy              *uuid.UUID
	Email                      string
	PhoneNumber                *string
	CareType                   *string
	SenderID                   *uuid.UUID
	LocationID                 *uuid.UUID
	EducationCurrentlyEnrolled bool
	EducationInstitution       *string
	EducationMentorName        *string
	EducationMentorPhone       *string
	EducationMentorEmail       *string
	EducationAdditionalNotes   *string
	WorkCurrentlyEmployed      bool
	WorkCurrentEmployer        *string
	WorkCurrentEmployerPhone   *string
	WorkCurrentEmployerEmail   *string
	WorkCurrentPosition        *string
	WorkStartDate              time.Time
	WorkAdditionalNotes        *string
}

type ListClientsParams struct {
	Limit      int32
	Offset     int32
	Status     *string
	LocationID *uuid.UUID
	Search     *string
}

type ListWaitingListClientsParams struct {
	Limit     int32
	Offset    int32
	Search    *string
	Placement *string
	SortDays  string
}

type ListInCareClientsParams struct {
	Limit          int32
	Offset         int32
	Search         *string
	Status         []string
	SortDaysInCare string
}

type UpdateClientParams struct {
	FirstName                  *string
	LastName                   *string
	DateOfBirth                time.Time
	Identity                   *bool
	Bsn                        *string
	BsnVerifiedBy              *uuid.UUID
	Source                     *string
	Birthplace                 *string
	Nationality                *string
	Email                      *string
	PhoneNumber                *string
	OrganizationID             *uuid.UUID
	Departement                *string
	Gender                     *string
	Filenumber                 *string
	ProfilePicture             *string
	Infix                      *string
	SenderID                   *uuid.UUID
	LocationID                 *uuid.UUID
	DepartureReason            *string
	DepartureReport            *string
	LegalMeasure               *string
	EducationCurrentlyEnrolled *bool
	EducationInstitution       *string
	EducationMentorName        *string
	EducationMentorPhone       *string
	EducationMentorEmail       *string
	EducationAdditionalNotes   *string
	EducationLevel             *string
	WorkCurrentlyEmployed      *bool
	WorkCurrentEmployer        *string
	WorkCurrentEmployerPhone   *string
	WorkCurrentEmployerEmail   *string
	WorkCurrentPosition        *string
	WorkStartDate              time.Time
	WorkAdditionalNotes        *string
	LivingSituation            *string
	LivingSituationNotes       *string
}

type ClientAddress struct {
	BelongsTo           *string
	Street              string
	HouseNumber         string
	HouseNumberAddition *string
	PostalCode          string
	City                string
	PhoneNumber         *string
}

type ClientPageClient struct {
	ID          uuid.UUID
	FirstName   string
	LastName    string
	Bsn         *string
	FileNumber  string
	Gender      string
	DateOfBirth *time.Time
	Age         *int32
	CareType    *string
	Address     ClientAddress
	Location    *ClientLocation
}

type ClientLocation struct {
	ID   uuid.UUID
	Name string
}

type ClientInCare struct {
	CareStartDate            *time.Time
	PlacedInCareAt           *time.Time
	DaysInCare               *int32
	EvaluationIntervalsWeeks int32
	LastEvaluationAnchorDate *time.Time
	NextEvaluationDate       *time.Time
}

type ClientCareSchedule struct {
	CareStartDate      *time.Time
	PlacedInCareAt     *time.Time
	DaysUntilStart     *int32
	ShouldBeActiveNow  bool
	NextEvaluationDate *time.Time
}

type ClientDischargeSchedule struct {
	DischargeDate          *time.Time
	DischargeReason        *string
	FinalEvaluation        *string
	DaysUntilDischarge     *int32
	IsDue                  bool
	MissingFinalEvaluation bool
}

type ClientDischargeSummary struct {
	DischargeDate   *time.Time
	DischargeReason *string
	FinalEvaluation *string
}

type ClientSenderMinimal struct {
	Name         string
	EmailAddress *string
	PhoneNumber  *string
}

type ClientCoordinator struct {
	EmployeeID *uuid.UUID
	FirstName  *string
	LastName   *string
	StartDate  *time.Time
}

type ClientContractSummary struct {
	HasActiveApprovedContract bool
	ActiveContract            *ClientActiveContract
	DaysUntilContractEnd      *int32
}

type ClientActiveContract struct {
	ID              uuid.UUID
	Status          *string
	StartDate       *time.Time
	EndDate         *time.Time
	FinancingAct    *string
	FinancingOption *string
	CareType        *string
}

type ClientEvaluationSummary struct {
	NextEvaluationDate *time.Time
	DaysLeft           *int32
	Priority           *string
	Draft              *ClientEvaluationDraft
	LastCompleted      *ClientEvaluationLastCompleted
}

type ClientEvaluationDraft struct {
	ID        uuid.UUID
	UpdatedAt *time.Time
}

type ClientEvaluationLastCompleted struct {
	ID                  uuid.UUID
	SubmittedAt         *time.Time
	CreatedByEmployeeID *uuid.UUID
	CreatorName         *string
}

type ClientGoalSummary struct {
	Title     string
	Priority  string
	TopicName *string
}

type ClientDocuments struct {
	Existing []string
	Missing  []string
}

type ClientIntake struct {
	SelfSufficiencyScore *int32
	Conclusion           *string
	ConclusionNotes      *string
}

type ClientRiskSummary struct {
	Flags []string
	Notes *string
}

type ClientPageAlert struct {
	Code     string
	Severity string
	Message  string
}

type ClientPageMeta struct {
	WaitlistSince time.Time
	LastUpdatedAt time.Time
}

type ClientPageCounts struct {
	Contracts    int64
	Incidents    int64
	Reports      int64
	Evaluations  int64
	Documents    int64
	Appointments int64
}

type ClientStatusTimeline struct {
	LastChangeReason *string
	LastChangedAt    *time.Time
	LastStatus       *string
}

type ClientPageDetail struct {
	SchemaVersion     int32
	Status            string
	Client            ClientPageClient
	Care              *ClientInCare
	CareSchedule      *ClientCareSchedule
	DischargeSchedule *ClientDischargeSchedule
	DischargeSummary  *ClientDischargeSummary
	Sender            *ClientSenderMinimal
	Coordinator       *ClientCoordinator
	ContractSummary   *ClientContractSummary
	EvaluationSummary *ClientEvaluationSummary
	EmergencyContacts []ClientEmergencyContact
	Documents         ClientDocuments
	Goals             []ClientGoalSummary
	Intake            ClientIntake
	Risks             ClientRiskSummary
	Counts            ClientPageCounts
	Alerts            []ClientPageAlert
	Meta              ClientPageMeta
	StatusTimeline    *ClientStatusTimeline
}

type ClientDocument struct {
	ID           uuid.UUID
	AttachmentID *uuid.UUID
	Uuid         uuid.UUID
	ClientID     uuid.UUID
	Label        string
	Name         string
	File         string
	FileURL      string
	Size         int32
	IsUsed       bool
	Tag          *string
	UpdatedAt    time.Time
	CreatedAt    time.Time
}

type AttachmentFile struct {
	UUID      uuid.UUID
	Name      string
	File      string
	Size      int32
	IsUsed    bool
	Tag       *string
	UpdatedAt time.Time
	CreatedAt time.Time
}

type AddClientDocumentItem struct {
	AttachmentID uuid.UUID
	Label        string
}

type AddClientDocumentParams struct {
	Documents []AddClientDocumentItem
}

type AddClientDocumentResult struct {
	ID           uuid.UUID
	AttachmentID *uuid.UUID
	ClientID     uuid.UUID
	Label        string
	Name         string
	File         string
	Size         int32
	IsUsed       bool
	Tag          *string
	UpdatedAt    time.Time
	CreatedAt    time.Time
}

type ListClientDocumentsParams struct {
	ClientID uuid.UUID
	Limit    int32
	Offset   int32
}

type ListClientDocumentsResult struct {
	Documents  []ClientDocument
	TotalCount int64
}

type DeleteClientDocumentResult struct {
	ID           uuid.UUID
	AttachmentID *uuid.UUID
}

type ClientRepository interface {
	CreateClient(ctx context.Context, params CreateClientParams) (*Client, error)
	ListClients(ctx context.Context, params ListClientsParams) (*ClientPage, error)
	ListWaitingListClients(ctx context.Context, params ListWaitingListClientsParams) (*WaitingListClientPage, error)
	ListInCareClients(ctx context.Context, params ListInCareClientsParams) (*InCareClientPage, error)
	GetClientCounts(ctx context.Context) (*ClientCounts, error)
	GetClientStatusCounts(ctx context.Context) (*ClientStatusCounts, error)
	GetClientByID(ctx context.Context, id uuid.UUID) (*ClientPageDetail, error)
	UpdateClient(ctx context.Context, id uuid.UUID, params UpdateClientParams) (*Client, error)
	GetClientAddresses(ctx context.Context, id uuid.UUID) ([]ClientAddress, error)
	UpdateClientStatus(ctx context.Context, clientID uuid.UUID, params UpdateClientStatusParams) (*UpdateClientStatusResult, error)
	PutClientInCare(ctx context.Context, clientID uuid.UUID, targetStatus db.ClientStatusEnum, careStartDay time.Time, params PutClientInCareParams) (*PutClientInCareResult, error)
	PutClientOutOfCare(ctx context.Context, clientID uuid.UUID, targetStatus db.ClientStatusEnum, dischargeDay time.Time, dischargeReason db.DischargeReasonEnum, finalEvaluation *string, params PutClientOutOfCareParams) (*PutClientOutOfCareResult, error)
	ListStatusHistory(ctx context.Context, params ListStatusHistoryParams) ([]ClientStatusHistory, error)
	AddClientDocuments(ctx context.Context, clientID uuid.UUID, params AddClientDocumentParams) ([]AddClientDocumentResult, error)
	ListClientDocuments(ctx context.Context, params ListClientDocumentsParams) (*ListClientDocumentsResult, error)
	DeleteClientDocument(ctx context.Context, clientID uuid.UUID, documentID uuid.UUID) (*DeleteClientDocumentResult, error)
	GetMissingClientDocuments(ctx context.Context, clientID uuid.UUID) ([]string, error)
	GetAttachmentsByUUIDs(ctx context.Context, ids []uuid.UUID) ([]AttachmentFile, error)
	CreateClientGoal(ctx context.Context, clientID uuid.UUID, params CreateClientGoalParams) (*ClientGoal, error)
	UpdateClientGoal(ctx context.Context, clientID uuid.UUID, goalID uuid.UUID, params UpdateClientGoalParams) (*UpdateClientGoalResult, error)
	GetClientGoalsForEvaluationPage(ctx context.Context, clientID uuid.UUID, employeeID uuid.UUID) (*ClientGoalsForEvaluationPage, error)
	CreateGoalEvaluation(ctx context.Context, clientID uuid.UUID, employeeID uuid.UUID, params CreateGoalEvaluationParams) (*GoalEvaluation, error)
	GetGoalEvaluationBootstrap(ctx context.Context, clientID uuid.UUID) (*GoalEvaluationBootstrap, error)
	ListClientSubmittedEvaluations(ctx context.Context, params ListClientSubmittedEvaluationsParams) (*ListClientSubmittedEvaluationsResult, error)
	ListGoalEvaluationHistory(ctx context.Context, params ListGoalEvaluationHistoryParams) (*ListGoalEvaluationHistoryResult, error)
	CreateLocationTransfer(ctx context.Context, clientID uuid.UUID, params CreateLocationTransferParams) error
	ApproveLocationTransfer(ctx context.Context, employeeID uuid.UUID, params ApproveLocationTransferParams) error
	ListLocationTransferRequests(ctx context.Context, params ListLocationTransferParams) (*ListLocationTransferResult, error)
	GetGoalEvaluation(ctx context.Context, evaluationID uuid.UUID) (*GoalEvaluation, error)
	ListUpcomingEvaluations(ctx context.Context, params ListUpcomingEvaluationsParams) (*ListUpcomingEvaluationsResult, error)
	ListRecentSubmittedEvaluations(ctx context.Context, params ListRecentSubmittedEvaluationsParams) (*ListRecentSubmittedEvaluationsResult, error)
	ListRecentDraftEvaluations(ctx context.Context, params ListRecentDraftEvaluationsParams) (*ListRecentDraftEvaluationsResult, error)
	CreateClientDiagnosis(ctx context.Context, params CreateClientDiagnosisParams) (*ClientDiagnosis, error)
	ListClientDiagnoses(ctx context.Context, params ListClientDiagnosesParams) (*ListClientDiagnosesResult, error)
	GetClientDiagnosis(ctx context.Context, clientID, diagnosisID uuid.UUID) (*ClientDiagnosis, error)
	UpdateClientDiagnosis(ctx context.Context, params UpdateClientDiagnosisParams) (*ClientDiagnosis, error)
	DeleteClientDiagnosis(ctx context.Context, clientID, diagnosisID uuid.UUID) (*DeleteClientDiagnosisResult, error)
	CreateClientMedicationOrder(ctx context.Context, params CreateClientMedicationOrderParams) (*ClientMedicationOrder, error)
	ListClientMedicationOrders(ctx context.Context, params ListClientMedicationOrdersParams) (*ListClientMedicationOrdersResult, error)
	GetClientMedicationOrder(ctx context.Context, clientID, orderID uuid.UUID) (*ClientMedicationOrder, error)
	UpdateClientMedicationOrder(ctx context.Context, params UpdateClientMedicationOrderParams) (*ClientMedicationOrder, error)
	DeleteClientMedicationOrder(ctx context.Context, clientID, orderID uuid.UUID) (*DeleteClientMedicationOrderResult, error)
	GetClientMedicalOverview(ctx context.Context, clientID uuid.UUID) (*ClientMedicalOverview, error)
	GetClientSender(ctx context.Context, clientID uuid.UUID) (*Sender, error)
	CreateClientEmergencyContact(ctx context.Context, params CreateClientEmergencyContactParams) (*ClientEmergencyContact, error)
	ListClientEmergencyContacts(ctx context.Context, params ListClientEmergencyContactsParams) (*ListClientEmergencyContactsResult, error)
	GetClientEmergencyContact(ctx context.Context, contactID uuid.UUID) (*ClientEmergencyContact, error)
	UpdateClientEmergencyContact(ctx context.Context, params UpdateClientEmergencyContactParams) (*ClientEmergencyContact, error)
	DeleteClientEmergencyContact(ctx context.Context, contactID uuid.UUID) (*DeleteClientEmergencyContactResult, error)
	CreateAssignedEmployee(ctx context.Context, params CreateAssignedEmployeeParams) (*AssignedEmployee, error)
	ListAssignedEmployees(ctx context.Context, params ListAssignedEmployeesParams) (*ListAssignedEmployeesResult, error)
	GetAssignedEmployee(ctx context.Context, assignmentID uuid.UUID) (*AssignedEmployee, error)
	UpdateAssignedEmployee(ctx context.Context, params UpdateAssignedEmployeeParams) (*AssignedEmployee, error)
	DeleteAssignedEmployee(ctx context.Context, assignmentID uuid.UUID) (*DeleteAssignedEmployeeResult, error)
	GetClientRelatedEmails(ctx context.Context, clientID uuid.UUID) (*ClientRelatedEmails, error)
	CreateProgressReport(ctx context.Context, params CreateProgressReportParams) (*ProgressReport, error)
	ListProgressReports(ctx context.Context, params ListProgressReportsParams) (*ListProgressReportsResult, error)
	GetProgressReport(ctx context.Context, reportID uuid.UUID) (*ProgressReport, error)
	UpdateProgressReport(ctx context.Context, params UpdateProgressReportParams) (*ProgressReport, error)
	DeleteProgressReport(ctx context.Context, reportID uuid.UUID) error
	GetProgressReportsByDateRange(ctx context.Context, params GetProgressReportsByDateRangeParams) ([]ProgressReport, error)
	CreateAiGeneratedReport(ctx context.Context, params CreateAiGeneratedReportParams) (*AiGeneratedReport, error)
	ListAiGeneratedReports(ctx context.Context, params ListAiGeneratedReportsParams) (*ListAiGeneratedReportsResult, error)
}

type UpdateClientStatusParams struct {
	Status       string
	Reason       string
	IsScheduled  bool
	ScheduledFor time.Time
}

type PutClientInCareParams struct {
	CareStartDate         string
	CoordinatorEmployeeID uuid.UUID
	PlacedInCareAt        *time.Time
	Reason                *string
}

type PutClientOutOfCareParams struct {
	DischargeDate   string
	DischargeReason string
	FinalEvaluation *string
	Reason          *string
}

type UpdateClientStatusResult struct {
	ID     uuid.UUID
	Status string
}

type PutClientInCareResult struct {
	ID                      uuid.UUID
	Status                  string
	CareStartDate           time.Time
	PlacedInCareAt          time.Time
	NextEvaluationDate      *time.Time
	CoordinatorAssignmentID uuid.UUID
	Warning                 *string
}

type PutClientOutOfCareResult struct {
	ID              uuid.UUID
	Status          string
	DischargeDate   time.Time
	DischargeReason string
	FinalEvaluation *string
}

type ClientStatusHistory struct {
	ID        uuid.UUID
	ClientID  uuid.UUID
	OldStatus *string
	NewStatus string
	ChangedAt time.Time
	ChangedBy *uuid.UUID
	Reason    *string
}

type ListStatusHistoryParams struct {
	ClientID uuid.UUID
	Limit    int32
	Offset   int32
}

type ClientService interface {
	CreateClient(ctx context.Context, params CreateClientParams) (*Client, error)
	ListClients(ctx context.Context, params ListClientsParams) (*ClientPage, error)
	ListWaitingListClients(ctx context.Context, params ListWaitingListClientsParams) (*WaitingListClientPage, error)
	ListInCareClients(ctx context.Context, params ListInCareClientsParams) (*InCareClientPage, error)
	GetClientCounts(ctx context.Context) (*ClientCounts, error)
	GetClientStatusCounts(ctx context.Context) (*ClientStatusCounts, error)
	GetClientByID(ctx context.Context, id uuid.UUID) (*ClientPageDetail, error)
	UpdateClient(ctx context.Context, id uuid.UUID, params UpdateClientParams) (*Client, error)
	GetClientAddresses(ctx context.Context, id uuid.UUID) ([]ClientAddress, error)
	UpdateClientStatus(ctx context.Context, clientID uuid.UUID, params UpdateClientStatusParams) (*UpdateClientStatusResult, error)
	PutClientInCare(ctx context.Context, clientID uuid.UUID, params PutClientInCareParams) (*PutClientInCareResult, error)
	PutClientOutOfCare(ctx context.Context, clientID uuid.UUID, params PutClientOutOfCareParams) (*PutClientOutOfCareResult, error)
	ListStatusHistory(ctx context.Context, clientID uuid.UUID) ([]ClientStatusHistory, error)
	AddClientDocument(ctx context.Context, clientID uuid.UUID, params AddClientDocumentParams) ([]AddClientDocumentResult, error)
	ListClientDocuments(ctx context.Context, params ListClientDocumentsParams) (*ListClientDocumentsResult, error)
	DeleteClientDocument(ctx context.Context, clientID uuid.UUID, documentID uuid.UUID) (*DeleteClientDocumentResult, error)
	GetMissingClientDocuments(ctx context.Context, clientID uuid.UUID) ([]string, error)
	CreateClientGoal(ctx context.Context, clientID uuid.UUID, params CreateClientGoalParams) (*ClientGoal, error)
	UpdateClientGoal(ctx context.Context, clientID uuid.UUID, goalID uuid.UUID, params UpdateClientGoalParams) (*UpdateClientGoalResult, error)
	GetClientGoalsForEvaluationPage(ctx context.Context, clientID uuid.UUID, employeeID uuid.UUID) (*ClientGoalsForEvaluationPage, error)
	CreateGoalEvaluation(ctx context.Context, clientID uuid.UUID, employeeID uuid.UUID, params CreateGoalEvaluationParams) (*GoalEvaluation, error)
	GetGoalEvaluationBootstrap(ctx context.Context, clientID uuid.UUID) (*GoalEvaluationBootstrap, error)
	ListClientSubmittedEvaluations(ctx context.Context, params ListClientSubmittedEvaluationsParams) (*ListClientSubmittedEvaluationsResult, error)
	ListGoalEvaluationHistory(ctx context.Context, params ListGoalEvaluationHistoryParams) (*ListGoalEvaluationHistoryResult, error)
	RequestLocationTransfer(ctx context.Context, clientID uuid.UUID, params CreateLocationTransferParams) error
	ApproveLocationTransfer(ctx context.Context, employeeID uuid.UUID, params ApproveLocationTransferParams) error
	ListLocationTransferRequests(ctx context.Context, params ListLocationTransferParams) (*ListLocationTransferResult, error)
	GetGoalEvaluation(ctx context.Context, evaluationID uuid.UUID) (*GoalEvaluation, error)
	ListUpcomingEvaluations(ctx context.Context, params ListUpcomingEvaluationsParams) (*ListUpcomingEvaluationsResult, error)
	ListRecentSubmittedEvaluations(ctx context.Context, params ListRecentSubmittedEvaluationsParams) (*ListRecentSubmittedEvaluationsResult, error)
	ListRecentDraftEvaluations(ctx context.Context, params ListRecentDraftEvaluationsParams) (*ListRecentDraftEvaluationsResult, error)
	CreateClientDiagnosis(ctx context.Context, clientID uuid.UUID, employeeID uuid.UUID, params CreateClientDiagnosisParams) (*ClientDiagnosis, error)
	ListClientDiagnoses(ctx context.Context, params ListClientDiagnosesParams) (*ListClientDiagnosesResult, error)
	GetClientDiagnosis(ctx context.Context, clientID, diagnosisID uuid.UUID) (*ClientDiagnosis, error)
	UpdateClientDiagnosis(ctx context.Context, clientID, diagnosisID, employeeID uuid.UUID, params UpdateClientDiagnosisParams) (*ClientDiagnosis, error)
	DeleteClientDiagnosis(ctx context.Context, clientID, diagnosisID uuid.UUID) (*DeleteClientDiagnosisResult, error)
	CreateClientMedicationOrder(ctx context.Context, clientID uuid.UUID, employeeID uuid.UUID, params CreateClientMedicationOrderParams) (*ClientMedicationOrder, error)
	ListClientMedicationOrders(ctx context.Context, params ListClientMedicationOrdersParams) (*ListClientMedicationOrdersResult, error)
	GetClientMedicationOrder(ctx context.Context, clientID, orderID uuid.UUID) (*ClientMedicationOrder, error)
	UpdateClientMedicationOrder(ctx context.Context, clientID, orderID, employeeID uuid.UUID, params UpdateClientMedicationOrderParams) (*ClientMedicationOrder, error)
	DeleteClientMedicationOrder(ctx context.Context, clientID, orderID uuid.UUID) (*DeleteClientMedicationOrderResult, error)
	GetClientMedicalOverview(ctx context.Context, clientID uuid.UUID) (*ClientMedicalOverview, error)
	GetClientSender(ctx context.Context, clientID uuid.UUID) (*Sender, error)
	CreateClientEmergencyContact(ctx context.Context, clientID uuid.UUID, params CreateClientEmergencyContactParams) (*ClientEmergencyContact, error)
	ListClientEmergencyContacts(ctx context.Context, params ListClientEmergencyContactsParams) (*ListClientEmergencyContactsResult, error)
	GetClientEmergencyContact(ctx context.Context, contactID uuid.UUID) (*ClientEmergencyContact, error)
	UpdateClientEmergencyContact(ctx context.Context, contactID uuid.UUID, params UpdateClientEmergencyContactParams) (*ClientEmergencyContact, error)
	DeleteClientEmergencyContact(ctx context.Context, contactID uuid.UUID) (*DeleteClientEmergencyContactResult, error)
	CreateAssignedEmployee(ctx context.Context, clientID uuid.UUID, params CreateAssignedEmployeeParams) (*AssignedEmployee, error)
	ListAssignedEmployees(ctx context.Context, params ListAssignedEmployeesParams) (*ListAssignedEmployeesResult, error)
	GetAssignedEmployee(ctx context.Context, assignmentID uuid.UUID) (*AssignedEmployee, error)
	UpdateAssignedEmployee(ctx context.Context, assignmentID uuid.UUID, params UpdateAssignedEmployeeParams) (*AssignedEmployee, error)
	DeleteAssignedEmployee(ctx context.Context, assignmentID uuid.UUID) (*DeleteAssignedEmployeeResult, error)
	GetClientRelatedEmails(ctx context.Context, clientID uuid.UUID) (*ClientRelatedEmails, error)
	CreateProgressReport(ctx context.Context, clientID uuid.UUID, params CreateProgressReportParams) (*ProgressReport, error)
	ListProgressReports(ctx context.Context, params ListProgressReportsParams) (*ListProgressReportsResult, error)
	GetProgressReport(ctx context.Context, reportID uuid.UUID) (*ProgressReport, error)
	UpdateProgressReport(ctx context.Context, reportID uuid.UUID, params UpdateProgressReportParams) (*ProgressReport, error)
	DeleteProgressReport(ctx context.Context, reportID uuid.UUID) error
	GenerateAutoReports(ctx context.Context, clientID uuid.UUID, startDate, endDate time.Time) (string, error)
	ConfirmAiProgressReport(ctx context.Context, clientID uuid.UUID, reportText string, startDate, endDate time.Time) (*AiGeneratedReport, error)
	ListAiGeneratedReports(ctx context.Context, params ListAiGeneratedReportsParams) (*ListAiGeneratedReportsResult, error)
}

var (
	ErrClientGoalClientNotFound        = errors.New("client for goal creation not found")
	ErrClientGoalGoalNotFound          = errors.New("client goal not found")
	ErrClientGoalTopicNotFound         = errors.New("topic for goal creation not found")
	ErrClientGoalDraftEvaluationExists = errors.New("client goal updates are blocked while a draft evaluation exists")
	ErrClientGoalEmptyPatch            = errors.New("at least one goal field must be provided")
	ErrClientGoalTitleRequired         = errors.New("title is required")
)

type ClientGoal struct {
	ID                uuid.UUID
	ClientID          uuid.UUID
	Title             string
	Description       *string
	Priority          string
	Status            string
	TopicID           *uuid.UUID
	TopicNameSnapshot *string
	Source            string
	SortOrder         int32
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type CreateClientGoalParams struct {
	Title       string
	Description *string
	Priority    *string
	TopicID     uuid.UUID
	SortOrder   *int32
}

type UpdateClientGoalParams struct {
	Title       *string
	Description *string
	Priority    *string
	TopicID     *uuid.UUID
	SortOrder   *int32
}

type UpdateClientGoalResult struct {
	MutationType      string
	GoalID            uuid.UUID
	ReplacementGoalID *uuid.UUID
	Goal              ClientGoal
}

type ClientGoalsForEvaluationPage struct {
	NextEvaluationDate    *time.Time
	MyDraftEvaluationID   *uuid.UUID
	IsResponsibleEmployee bool
	CanUpdateGoals        bool
	GoalUpdateBlockReason *string
	Goals                 []ClientGoalForEvaluationPage
}

type ClientGoalForEvaluationPage struct {
	ID                     uuid.UUID
	TopicName              *string
	Title                  string
	Priority               string
	LastEvaluationProgress *string
}

type GoalEvaluation struct {
	ID                      uuid.UUID
	ClientID                uuid.UUID
	EvaluationDate          time.Time
	PeriodStart             *time.Time
	PeriodEnd               *time.Time
	EvaluationIntervalWeeks int32
	Status                  string
	OverallNotes            *string
	CreatedByEmployeeID     *uuid.UUID
	CreatorName             *string
	CreatedAt               time.Time
	UpdatedAt               time.Time
	SubmitError             *string
	Items                   []GoalEvaluationItem
}

type GoalEvaluationItem struct {
	ID                uuid.UUID
	EvaluationID      uuid.UUID
	GoalID            uuid.UUID
	GoalTitle         string
	GoalDescription   *string
	TopicNameSnapshot *string
	Progress          string
	Notes             *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type CreateGoalEvaluationParams struct {
	OverallNotes *string
	Submit       bool
	Items        []GoalEvaluationItemParams
}

type GoalEvaluationItemParams struct {
	GoalID   uuid.UUID
	Progress string
	Notes    *string
}

type GoalEvaluationBootstrap struct {
	ClientID                uuid.UUID
	ClientFirstName         string
	ClientLastName          string
	NextEvaluationDate      *time.Time
	DaysLeft                *int32
	Priority                *string
	ExistingDraft           *GoalEvaluationBootstrapDraft
	LastCompletedEvaluation *GoalEvaluationBootstrapCompleted
	ActiveGoals             []GoalEvaluationBootstrapActiveGoal
}

type GoalEvaluationBootstrapDraft struct {
	ID             uuid.UUID
	EvaluationDate time.Time
	UpdatedAt      time.Time
}

type GoalEvaluationBootstrapCompleted struct {
	ID                  uuid.UUID
	EvaluationDate      time.Time
	SubmittedAt         time.Time
	OverallNotes        *string
	CreatedByEmployeeID *uuid.UUID
	CreatorName         *string
}

type GoalEvaluationBootstrapActiveGoal struct {
	GoalID            uuid.UUID
	Title             string
	TopicNameSnapshot *string
	Priority          string
	SortOrder         int32
	LastProgress      *string
	LastNotes         *string
}

type ListClientSubmittedEvaluationsParams struct {
	ClientID uuid.UUID
	Limit    int32
	Offset   int32
}

type ListClientSubmittedEvaluationsResult struct {
	Items      []ListClientSubmittedEvaluationsItem
	TotalCount int64
}

type ListClientSubmittedEvaluationsItem struct {
	EvaluationID        uuid.UUID
	EvaluationDate      time.Time
	SubmittedAt         time.Time
	FilledGoalsCount    int32
	TotalGoalsCount     int32
	CreatedByEmployeeID *uuid.UUID
	CreatorName         *string
}

type ListGoalEvaluationHistoryParams struct {
	ClientID uuid.UUID
	GoalID   uuid.UUID
	Limit    int32
	Offset   int32
}

type ListGoalEvaluationHistoryResult struct {
	Items      []ListGoalEvaluationHistoryItem
	TotalCount int64
}

type ListGoalEvaluationHistoryItem struct {
	EvaluationID        uuid.UUID
	EvaluationDate      time.Time
	SubmittedAt         time.Time
	Progress            string
	Notes               *string
	CreatedByEmployeeID *uuid.UUID
	CreatorName         *string
	PeriodStart         *time.Time
	PeriodEnd           *time.Time
}

type LocationTransfer struct {
	ID                 uuid.UUID
	ClientID           uuid.UUID
	FromLocationID     *uuid.UUID
	ToLocationID       *uuid.UUID
	NewMentorID        *uuid.UUID
	RequestDate        time.Time
	Status             string
	ApprovedRejectedBy *uuid.UUID
	ApprovedRejectedAt *time.Time
	Reason             *string
	MentorFirstName    *string
	MentorLastName     *string
}

type CreateLocationTransferParams struct {
	FromLocationID uuid.UUID
	ToLocationID   uuid.UUID
	Reason         string
	NewMentorID    *uuid.UUID
}

type ApproveLocationTransferParams struct {
	TransferID uuid.UUID
	Status     string
}

type ListLocationTransferParams struct {
	Limit  int32
	Offset int32
}

type ListLocationTransferResult struct {
	Items      []LocationTransfer
	TotalCount int64
}

type ListUpcomingEvaluationsParams struct {
	EmployeeID uuid.UUID
	Limit      int32
	Offset     int32
}

type ListUpcomingEvaluationsResult struct {
	Items      []UpcomingEvaluation
	TotalCount int64
}

type UpcomingEvaluation struct {
	ClientID         uuid.UUID
	ClientFirstName  string
	ClientLastName   string
	DueDate          time.Time
	DaysLeft         int32
	Priority         string
	HasDraft         bool
	FilledGoalsCount int32
	TotalGoalsCount  int32
}

type ListRecentSubmittedEvaluationsParams struct {
	EmployeeID uuid.UUID
	Limit      int32
	Offset     int32
}

type ListRecentSubmittedEvaluationsResult struct {
	Items      []RecentSubmittedEvaluation
	TotalCount int64
}

type RecentSubmittedEvaluation struct {
	EvaluationID       uuid.UUID
	ClientID           uuid.UUID
	ClientFirstName    string
	ClientLastName     string
	EvaluationDate     time.Time
	SubmittedAt        time.Time
	NextEvaluationDate *time.Time
	FilledGoalsCount   int32
	TotalGoalsCount    int32
}

type ListRecentDraftEvaluationsParams struct {
	EmployeeID uuid.UUID
	Limit      int32
	Offset     int32
}

type ListRecentDraftEvaluationsResult struct {
	Items      []RecentDraftEvaluation
	TotalCount int64
}

type RecentDraftEvaluation struct {
	EvaluationID     uuid.UUID
	ClientID         uuid.UUID
	ClientFirstName  string
	ClientLastName   string
	DueDate          time.Time
	UpdatedAt        time.Time
	DaysLeft         int32
	Priority         string
	FilledGoalsCount int32
	TotalGoalsCount  int32
}

// =====================
// Medical - Diagnoses
// =====================

type ClientDiagnosis struct {
	ID                  uuid.UUID
	ClientID            uuid.UUID
	CodeSystem          string
	Code                string
	Title               *string
	Description         *string
	Status              string
	Severity            string
	DiagnosedOn         *time.Time
	ResolvedOn          *time.Time
	DiagnosingClinician *string
	Notes               *string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type CreateClientDiagnosisParams struct {
	ClientID            uuid.UUID
	CodeSystem          string
	Code                string
	Title               *string
	Description         *string
	Status              *string
	Severity            *string
	DiagnosedOn         *time.Time
	ResolvedOn          *time.Time
	DiagnosingClinician *string
	Notes               *string
	CreatedByEmployeeID *uuid.UUID
	UpdatedByEmployeeID *uuid.UUID
}

type UpdateClientDiagnosisParams struct {
	ClientID            uuid.UUID
	ID                  uuid.UUID
	CodeSystem          *string
	Code                *string
	Title               *string
	Description         *string
	Status              *string
	Severity            *string
	DiagnosedOn         *time.Time
	ResolvedOn          *time.Time
	DiagnosingClinician *string
	Notes               *string
	UpdatedByEmployeeID *uuid.UUID
}

type ListClientDiagnosesParams struct {
	ClientID uuid.UUID
	Limit    int32
	Offset   int32
}

type ListClientDiagnosesResult struct {
	Items      []ClientDiagnosis
	TotalCount int64
}

type DeleteClientDiagnosisResult struct {
	ID uuid.UUID
}

// =====================
// Medical - Medication Orders
// =====================

type ClientMedicationOrder struct {
	ID                           uuid.UUID
	ClientID                     uuid.UUID
	DiagnosisID                  *uuid.UUID
	MedicationName               string
	DosageText                   string
	DoseAmount                   *float64
	DoseUnit                     *string
	Route                        *string
	FrequencyText                *string
	Schedule                     []byte
	IsPrn                        bool
	PrnIndication                *string
	MaxDosesPer24h               *int32
	StartDate                    time.Time
	EndDate                      *time.Time
	Status                       string
	AdminMode                    string
	ResponsibleEmployeeID        *uuid.UUID
	ResponsibleEmployeeFirstName *string
	ResponsibleEmployeeLastName  *string
	IsCritical                   bool
	Notes                        *string
	SourceAttachmentUUID         *uuid.UUID
	DiagnosisTitle               *string
	DiagnosisCodeSystem          *string
	DiagnosisCode                *string
	CreatedAt                    time.Time
	UpdatedAt                    time.Time
}

type CreateClientMedicationOrderParams struct {
	ClientID              uuid.UUID
	DiagnosisID           *uuid.UUID
	MedicationName        string
	DosageText            string
	DoseAmount            *float64
	DoseUnit              *string
	Route                 *string
	FrequencyText         *string
	Schedule              []byte
	IsPrn                 bool
	PrnIndication         *string
	MaxDosesPer24h        *int32
	StartDate             time.Time
	EndDate               *time.Time
	Status                *string
	AdminMode             *string
	ResponsibleEmployeeID *uuid.UUID
	IsCritical            bool
	Notes                 *string
	SourceAttachmentUUID  *uuid.UUID
	CreatedByEmployeeID   *uuid.UUID
	UpdatedByEmployeeID   *uuid.UUID
}

type UpdateClientMedicationOrderParams struct {
	ClientID              uuid.UUID
	ID                    uuid.UUID
	DiagnosisID           *uuid.UUID
	MedicationName        *string
	DosageText            *string
	DoseAmount            *float64
	DoseUnit              *string
	Route                 *string
	FrequencyText         *string
	Schedule              []byte
	IsPrn                 *bool
	PrnIndication         *string
	MaxDosesPer24h        *int32
	StartDate             *time.Time
	EndDate               *time.Time
	Status                *string
	AdminMode             *string
	ResponsibleEmployeeID *uuid.UUID
	IsCritical            *bool
	Notes                 *string
	SourceAttachmentUUID  *uuid.UUID
	UpdatedByEmployeeID   *uuid.UUID
}

type ListClientMedicationOrdersParams struct {
	ClientID    uuid.UUID
	Status      *string
	AdminMode   *string
	DiagnosisID *uuid.UUID
	Search      *string
	Limit       int32
	Offset      int32
}

type ListClientMedicationOrdersResult struct {
	Items      []ClientMedicationOrder
	TotalCount int64
}

type DeleteClientMedicationOrderResult struct {
	ID uuid.UUID
}

// =====================
// Medical - Overview
// =====================

type ClientMedicalOverview struct {
	Diagnoses        []ClientDiagnosis
	MedicationOrders []ClientMedicationOrder
}

// =====================
// Network - Sender
// =====================

type Sender struct {
	ID                  uuid.UUID
	Types               string
	Name                string
	Street              *string
	HouseNumber         *string
	HouseNumberAddition *string
	PostalCode          *string
	City                *string
	Land                *string
	Kvknumber           *string
	Btwnumber           *string
	PhoneNumber         *string
	ClientNumber        *string
	EmailAddress        *string
	Contacts            []byte
	IsArchived          bool
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type SenderContact struct {
	Name        *string
	Email       *string
	PhoneNumber *string
}

// =====================
// Network - Emergency Contacts
// =====================

type ClientEmergencyContact struct {
	ID               uuid.UUID
	ClientID         uuid.UUID
	FirstName        *string
	LastName         *string
	Email            *string
	PhoneNumber      *string
	Address          *string
	Relationship     *string
	RelationStatus   *string
	CreatedAt        time.Time
	IsVerified       bool
	MedicalReports   bool
	IncidentsReports bool
	GoalsReports     bool
}

type CreateClientEmergencyContactParams struct {
	ClientID         uuid.UUID
	FirstName        *string
	LastName         *string
	Email            *string
	PhoneNumber      *string
	Address          *string
	Relationship     *string
	RelationStatus   *string
	MedicalReports   bool
	IncidentsReports bool
	GoalsReports     bool
}

type UpdateClientEmergencyContactParams struct {
	ID               uuid.UUID
	FirstName        *string
	LastName         *string
	Email            *string
	PhoneNumber      *string
	Address          *string
	Relationship     *string
	RelationStatus   *string
	MedicalReports   *bool
	IncidentsReports *bool
	GoalsReports     *bool
}

type ListClientEmergencyContactsParams struct {
	ClientID uuid.UUID
	Search   string
	Limit    int32
	Offset   int32
}

type ListClientEmergencyContactsResult struct {
	Items      []ClientEmergencyContact
	TotalCount int64
}

type DeleteClientEmergencyContactResult struct {
	ID uuid.UUID
}

// =====================
// Network - Assigned Employees
// =====================

type AssignedEmployee struct {
	ID                 uuid.UUID
	ClientID           uuid.UUID
	EmployeeID         uuid.UUID
	StartDate          time.Time
	Role               string
	CreatedAt          time.Time
	EmployeeFirstName  string
	EmployeeLastName   string
	UserID             uuid.UUID
	ClientFirstName    string
	ClientLastName     string
	ClientLocationName *string
}

type CreateAssignedEmployeeParams struct {
	ClientID   uuid.UUID
	EmployeeID uuid.UUID
	StartDate  time.Time
	Role       string
}

type UpdateAssignedEmployeeParams struct {
	ID         uuid.UUID
	EmployeeID *uuid.UUID
	StartDate  time.Time
	Role       *string
}

type ListAssignedEmployeesParams struct {
	ClientID uuid.UUID
	Limit    int32
	Offset   int32
}

type ListAssignedEmployeesResult struct {
	Items      []AssignedEmployee
	TotalCount int64
}

type DeleteAssignedEmployeeResult struct {
	ID uuid.UUID
}

// =====================
// Network - Related Emails
// =====================

type ClientRelatedEmails struct {
	Emails []*string
}

// =====================
// Progress Reports
// =====================

type ProgressReport struct {
	ID                     uuid.UUID
	ClientID               uuid.UUID
	Date                   time.Time
	Title                  *string
	ReportText             string
	EmployeeID             *uuid.UUID
	Type                   string
	EmotionalState         string
	CreatedAt              time.Time
	EmployeeFirstName      string
	EmployeeLastName       string
	EmployeeProfilePicture *string
}

type CreateProgressReportParams struct {
	ClientID       uuid.UUID
	EmployeeID     *uuid.UUID
	Title          *string
	Date           time.Time
	ReportText     string
	Type           string
	EmotionalState string
}

type UpdateProgressReportParams struct {
	ID             uuid.UUID
	EmployeeID     *uuid.UUID
	Title          *string
	Date           time.Time
	ReportText     *string
	Type           *string
	EmotionalState *string
}

type ListProgressReportsParams struct {
	ClientID uuid.UUID
	Type     *string
	Limit    int32
	Offset   int32
}

type ListProgressReportsResult struct {
	Items      []ProgressReport
	TotalCount int64
}

type GetProgressReportsByDateRangeParams struct {
	ClientID  uuid.UUID
	StartDate time.Time
	EndDate   time.Time
}

// =====================
// AI Generated Reports
// =====================

type AiGeneratedReport struct {
	ID         uuid.UUID
	ClientID   uuid.UUID
	ReportText string
	StartDate  time.Time
	EndDate    time.Time
	CreatedAt  time.Time
}

type CreateAiGeneratedReportParams struct {
	ClientID   uuid.UUID
	ReportText string
	StartDate  time.Time
	EndDate    time.Time
}

type ListAiGeneratedReportsParams struct {
	ClientID uuid.UUID
	Limit    int32
	Offset   int32
}

type ListAiGeneratedReportsResult struct {
	Items      []AiGeneratedReport
	TotalCount int64
}

// =====================
// Auto Report Generator
// =====================

type AutoReportGenerator interface {
	GenerateAutoReports(ctx context.Context, pastReportsText string) (string, error)
}
