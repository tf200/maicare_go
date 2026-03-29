package clientp

import (
	"context"

	"maicare_go/async/aclient"
	"maicare_go/pagination"
	"maicare_go/service/deps"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ClientService interface {
	// Client Details
	CreateClientDetails(req CreateClientDetailsRequest, ctx context.Context) (*CreateClientDetailsResponse, error)
	ListClientDetails(ctx *gin.Context, req ListClientsApiParams) (*pagination.Response[ListClientsApiResponse], error)
	ListWaitingListClients(ctx *gin.Context, req ListWaitingListClientsParams) (*pagination.Response[ListWaitingListClientsResponse], error)
	ListInCareClients(ctx *gin.Context, req ListInCareClientsParams) (*pagination.Response[ListInCareClientsResponse], error)
	GetClientsCount(ctx context.Context) (*GetClientsCountResponse, error)
	GetClientStatusCounts(ctx context.Context) (*GetClientStatusCountsResponse, error)
	GetClientDetails(ctx context.Context, clientID uuid.UUID) (*GetClientApiResponse, error)
	GetClientAddresses(ctx context.Context, clientID uuid.UUID) (*GetClientAddressesApiResponse, error)
	UpdateClientDetails(ctx context.Context, req UpdateClientDetailsRequest, clientID uuid.UUID) (*UpdateClientDetailsResponse, error)
	UpdateClientStatus(ctx context.Context, req UpdateClientStatusRequest, clientID uuid.UUID) (*UpdateClientStatusResponse, error)
	PutClientInCare(ctx context.Context, req PutClientInCareRequest, clientID uuid.UUID) (*PutClientInCareResponse, error)
	PutClientOutOfCare(ctx context.Context, req PutClientOutOfCareRequest, clientID uuid.UUID) (*PutClientOutOfCareResponse, error)
	ListStatusHistory(ctx context.Context, clientID uuid.UUID) ([]ListStatusHistoryApiResponse, error)
	// Client Documents
	AddClientDocument(ctx context.Context, req AddClientDocumentApiRequest, clientID uuid.UUID) (*AddClientDocumentApiResponse, error)
	ListClientDocuments(ctx *gin.Context, req ListClientDocumentsApiRequest, clientID uuid.UUID) (*pagination.Response[ListClientDocumentsApiResponse], error)
	DeleteClientDocument(ctx context.Context, clientID uuid.UUID, attachmentID uuid.UUID) (*DeleteClientDocumentApiResponse, error)
	GetMissingClientDocuments(ctx context.Context, clientID uuid.UUID) (*GetMissingClientDocumentsApiResponse, error)

	// Client Appointment Card
	GetAppointmentCard(ctx context.Context, clientID uuid.UUID) (*GetAppointmentCardResponse, error)
	UpdateAppointmentCard(req UpdateAppointmentCardRequest, clientID uuid.UUID, ctx context.Context) (*UpdateAppointmentCardResponse, error)
	GenerateAppointmentCardDocumentApi(ctx context.Context, clientID uuid.UUID) ([]byte, string, error)

	// Client Incidents
	CreateIncident(ctx context.Context, req CreateIncidentRequest) (*CreateIncidentResponse, error)
	ListIncidents(ctx *gin.Context, req ListIncidentsRequest, clientID uuid.UUID) (*pagination.Response[ListIncidentsResponse], error)
	GetIncident(ctx context.Context, incidentID uuid.UUID) (*GetIncidentResponse, error)
	UpdateIncident(ctx context.Context, req UpdateIncidentRequest, incidentID uuid.UUID) (*UpdateIncidentResponse, error)
	DeleteIncident(ctx context.Context, incidentID uuid.UUID) error
	GenerateIncidentFile(ctx context.Context, incidentID uuid.UUID) ([]byte, string, error)
	ConfirmIncident(ctx context.Context, incidentID uuid.UUID, confirmedByUserID uuid.UUID) (*ConfirmIncidentResponse, error)
	ListAllIncidents(ctx *gin.Context, req *ListAllIncidentsRequest) (*pagination.Response[ListAllIncidentsResponse], error)
	GetIncidentCounts(ctx context.Context) (*GetIncidentCountsResponse, error)

	// Client Diagnoses
	CreateClientDiagnosis(ctx context.Context, req CreateClientDiagnosisRequest, clientID uuid.UUID) (*ClientDiagnosisResponse, error)
	ListClientDiagnoses(ctx *gin.Context, req ListClientDiagnosesRequest, clientID uuid.UUID) (*pagination.Response[ClientDiagnosisResponse], error)
	GetClientDiagnosis(ctx context.Context, clientID uuid.UUID, diagnosisID uuid.UUID) (*ClientDiagnosisResponse, error)
	UpdateClientDiagnosis(ctx context.Context, req UpdateClientDiagnosisRequest, clientID uuid.UUID, diagnosisID uuid.UUID) (*ClientDiagnosisResponse, error)
	DeleteClientDiagnosis(ctx context.Context, clientID uuid.UUID, diagnosisID uuid.UUID) (*DeleteClientDiagnosisResponse, error)

	// Client Medication Orders
	CreateClientMedicationOrder(ctx context.Context, req CreateClientMedicationOrderRequest, clientID uuid.UUID) (*ClientMedicationOrderResponse, error)
	ListClientMedicationOrders(ctx *gin.Context, req ListClientMedicationOrdersRequest, clientID uuid.UUID) (*pagination.Response[ClientMedicationOrderResponse], error)
	GetClientMedicationOrder(ctx context.Context, clientID uuid.UUID, orderID uuid.UUID) (*ClientMedicationOrderResponse, error)
	UpdateClientMedicationOrder(ctx context.Context, req UpdateClientMedicationOrderRequest, clientID uuid.UUID, orderID uuid.UUID) (*ClientMedicationOrderResponse, error)
	DeleteClientMedicationOrder(ctx context.Context, clientID uuid.UUID, orderID uuid.UUID) (*DeleteClientMedicationOrderResponse, error)

	// Client Medical Overview
	GetClientMedicalOverview(ctx context.Context, clientID uuid.UUID) (*ClientMedicalOverviewResponse, error)

	// Client Sender
	GetClientSender(ctx context.Context, clientID uuid.UUID) (*GetClientSenderResponse, error)

	// Client Emergency Contacts
	CreateClientEmergencyContact(ctx context.Context, req CreateClientEmergencyContactParams, clientID uuid.UUID) (*CreateClientEmergencyContactResponse, error)
	ListClientEmergencyContacts(ctx *gin.Context, req ListClientEmergencyContactsRequest, clientID uuid.UUID) (*pagination.Response[ListClientEmergencyContactsResponse], error)
	GetClientEmergencyContact(ctx context.Context, contactID uuid.UUID) (*GetClientEmergencyContactResponse, error)
	UpdateClientEmergencyContact(ctx context.Context, req UpdateClientEmergencyContactParams, contactID uuid.UUID) (*UpdateClientEmergencyContactResponse, error)
	DeleteClientEmergencyContact(ctx context.Context, contactID uuid.UUID) (*DeleteClientEmergencyContactResponse, error)

	// Client Involved employees
	AssignEmployeeToClient(ctx context.Context, req AssignEmployeeRequest, clientID uuid.UUID) (*AssignEmployeeResponse, error)
	ListAssignedEmployees(ctx *gin.Context, req ListAssignedEmployeesRequest, clientID uuid.UUID) (*pagination.Response[ListAssignedEmployeesResponse], error)
	GetAssignedEmployee(ctx context.Context, assignmentID uuid.UUID) (*GetAssignedEmployeeResponse, error)
	UpdateAssignedEmployee(ctx context.Context, req UpdateAssignedEmployeeRequest, assignmentID uuid.UUID) (*UpdateAssignedEmployeeResponse, error)
	DeleteAssignedEmployee(ctx context.Context, assignmentID uuid.UUID) (*DeleteAssignedEmployeeResponse, error)

	// Client Network Emails
	GetClientRelatedEmail(ctx context.Context, clientID uuid.UUID) (*GetClientRelatedEmailsResponse, error)

	// Client Progress Reports
	CreateProgressReport(ctx context.Context, req *CreateProgressReportRequest, clientID uuid.UUID) (*CreateProgressReportResponse, error)
	ListProgressReports(ctx *gin.Context, req *ListProgressReportsRequest, clientID uuid.UUID) (*pagination.Response[ListProgressReportsResponse], error)
	GetProgressReport(ctx context.Context, reportID uuid.UUID) (*GetProgressReportResponse, error)
	UpdateProgressReport(ctx context.Context, req *UpdateProgressReportRequest, reportID uuid.UUID) (*GetProgressReportResponse, error)
	DeleteProgressReport(ctx context.Context, reportID uuid.UUID) error
	GenerateAutoReports(ctx context.Context, req *GenerateAutoReportsRequest, clientID uuid.UUID) (*GenerateAutoReportsResponse, error)
	ConfirmAiProgressReport(ctx context.Context, clientID uuid.UUID, req *ConfirmProgressReportRequest) (*ConfirmProgressReportResponse, error)
	ListAiGeneratedReports(ctx *gin.Context, req *ListAiGeneratedReportsRequest, clientID uuid.UUID) (*pagination.Response[ListAiGeneratedReportsResponse], error)

	// Registration Form
	CreateRegistrationForm(ctx context.Context, req *CreateRegistrationFormRequest) (*CreateRegistrationFormResponse, error)
	ListRegistrationForms(ctx *gin.Context, req *ListRegistrationFormsRequest) (*pagination.Response[ListRegistrationFormsResponse], error)
	GetRegistrationFormB(tx context.Context, formID uuid.UUID) (*GetRegistrationFormResponse, error)
	UpdateRegistrationForm(ctx context.Context, req *UpdateRegistrationFormRequest, formID uuid.UUID) (*UpdateRegistrationFormResponse, error)
	DeleteRegistrationForm(ctx context.Context, formID uuid.UUID) error
	UpdateRegistrationFormStatus(ctx context.Context, req *UpdateRegistrationFormStatusRequest, formID uuid.UUID, employeeID uuid.UUID) error

	// Process Registration Form
	ProcessRegistrationForm(ctx context.Context, req *ProcessRegistrationFormRequest, formID uuid.UUID, employeeID uuid.UUID) error

	// Public Intake endpoints
	GetPublicIntakeOptions(ctx context.Context, token string) (*PublicIntakeOptionsResponse, error)
	SelectIntakeDate(ctx context.Context, token string, req *SelectIntakeDateRequest) error

	// Intake Form
	CreateIntakeForm(ctx context.Context, req *CreateIntakeFormRequest) (*CreateIntakeFormResponse, error)
	ListIntakeForms(ctx *gin.Context, req *ListIntakeFormsRequest) (*pagination.Response[ListIntakeFormsResponse], error)
	GetIntakeFormTotals(ctx context.Context) (*GetIntakeFormTotalsResponse, error)
	GetIntakeForm(ctx context.Context, intakeFormID uuid.UUID) (*GetIntakeFormResponse, error)
	UpdateIntakeForm(ctx context.Context, intakeFormID uuid.UUID, req *UpdateIntakeFormRequest) (*UpdateIntakeFormResponse, error)
	CreateIntakeFormGoals(ctx context.Context, intakeFormID uuid.UUID, req *CreateIntakeFormGoalsRequest) (*CreateIntakeFormGoalsResponse, error)
	UpdateIntakeConclusion(ctx context.Context, intakeFormID uuid.UUID, req *UpdateIntakeConclusionRequest) (*UpdateIntakeConclusionResponse, error)

	// Intake Maturity Assessments

	GenerateIntakeGoals(ctx context.Context, req *GenerateIntakeGoalsRequest) (*GenerateIntakeGoalsResponse, error)

	// Promote Intake to Client
	PromoteIntakeToClient(ctx context.Context, req *PromoteIntakeToClientRequest) (*PromoteIntakeToClientResponse, error)

	// Goal Evaluations
	CreateClientGoal(ctx context.Context, clientID uuid.UUID, req CreateClientGoalRequest) (*CreateClientGoalResponse, error)
	UpdateClientGoal(ctx context.Context, clientID uuid.UUID, goalID uuid.UUID, req UpdateClientGoalRequest) (*UpdateClientGoalResponse, error)
	CreateGoalEvaluation(ctx context.Context, clientID uuid.UUID, employeeID uuid.UUID, req CreateGoalEvaluationRequest) (*GoalEvaluationResponse, error)
	GetGoalEvaluation(ctx context.Context, evaluationID uuid.UUID) (*GoalEvaluationResponse, error)
	GetGoalEvaluationBootstrap(ctx context.Context, clientID uuid.UUID) (*GoalEvaluationBootstrapResponse, error)
	GetClientGoalsForEvaluationPage(ctx context.Context, clientID uuid.UUID, employeeID uuid.UUID) (*GetClientGoalsForEvaluationPageResponse, error)
	ListClientSubmittedEvaluations(ctx *gin.Context, clientID uuid.UUID, req ListClientSubmittedEvaluationsRequest) (*pagination.Response[ListClientSubmittedEvaluationsResponse], error)
	ListGoalEvaluationHistory(ctx *gin.Context, clientID uuid.UUID, goalID uuid.UUID, req ListGoalEvaluationHistoryRequest) (*pagination.Response[ListGoalEvaluationHistoryResponse], error)
	ListUpcomingEvaluations(ctx *gin.Context, coordinatorID uuid.UUID, req ListUpcomingEvaluationsRequest) (*pagination.Response[ListUpcomingEvaluationsResponse], error)
	ListRecentSubmittedEvaluations(ctx *gin.Context, employeeID uuid.UUID, req ListRecentSubmittedEvaluationsRequest) (*pagination.Response[ListRecentSubmittedEvaluationsResponse], error)
	ListRecentDraftEvaluations(ctx *gin.Context, employeeID uuid.UUID, req ListRecentDraftEvaluationsRequest) (*pagination.Response[ListRecentDraftEvaluationsResponse], error)

	// Location Transfer
	RequestLocationTransfer(ctx context.Context, clientID uuid.UUID, req LocationTransferRequest) error
	ApproveLocationTransfer(ctx context.Context, employeeID uuid.UUID, req ApproveOrRejectLocationTransferRequest) error
	ListLocationTransferRequests(ctx *gin.Context, req ListLocationTransferRequestsRequest) (*pagination.Response[ListLocationTransferRequestsResponse], error)
}

type clientService struct {
	*deps.ServiceDependencies
	asynqClient aclient.AsynqClientInterface
}

func NewClientService(deps *deps.ServiceDependencies, asynqClient aclient.AsynqClientInterface) ClientService {
	return &clientService{
		ServiceDependencies: deps,
		asynqClient:         asynqClient,
	}
}
