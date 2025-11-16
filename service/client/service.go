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
	GetClientsCount(ctx context.Context) (*GetClientsCountResponse, error)
	GetClientDetails(ctx context.Context, clientID uuid.UUID) (*GetClientApiResponse, error)
	GetClientAddresses(ctx context.Context, clientID uuid.UUID) (*GetClientAddressesApiResponse, error)
	UpdateClientDetails(ctx context.Context, req UpdateClientDetailsRequest, clientID uuid.UUID) (*UpdateClientDetailsResponse, error)
	UpdateClientStatus(ctx context.Context, req UpdateClientStatusRequest, clientID uuid.UUID) (*UpdateClientStatusResponse, error)
	ListStatusHistory(ctx context.Context, clientID uuid.UUID) ([]ListStatusHistoryApiResponse, error)
	SetClientProfilePicture(ctx context.Context, req SetClientProfilePictureRequest, clientID uuid.UUID) (*SetClientProfilePictureResponse, error)
	// Client Documents
	AddClientDocument(ctx context.Context, req AddClientDocumentApiRequest, clientID uuid.UUID) (*AddClientDocumentApiResponse, error)
	ListClientDocuments(ctx *gin.Context, req ListClientDocumentsApiRequest, clientID uuid.UUID) (*pagination.Response[ListClientDocumentsApiResponse], error)
	DeleteClientDocument(ctx context.Context, clientID uuid.UUID, attachmentID uuid.UUID) (*DeleteClientDocumentApiResponse, error)
	GetMissingClientDocuments(ctx context.Context, clientID uuid.UUID) (*GetMissingClientDocumentsApiResponse, error)

	// Client Appointment Card
	CreateAppointmentCard(req CreateAppointmentCardRequest, clientID uuid.UUID, ctx context.Context) (*CreateAppointmentCardResponse, error)
	GetAppointmentCard(ctx context.Context, clientID uuid.UUID) (*GetAppointmentCardResponse, error)
	UpdateAppointmentCard(req UpdateAppointmentCardRequest, clientID uuid.UUID, ctx context.Context) (*UpdateAppointmentCardResponse, error)
	GenerateAppointmentCardDocumentApi(ctx context.Context, clientID uuid.UUID) (*GenerateAppointmentCardDocumentApiResponse, error)

	// Client Incidents
	CreateIncident(ctx context.Context, req CreateIncidentRequest, clientID uuid.UUID) (*CreateIncidentResponse, error)
	ListIncidents(ctx *gin.Context, req ListIncidentsRequest, clientID uuid.UUID) (*pagination.Response[ListIncidentsResponse], error)
	GetIncident(ctx context.Context, incidentID int64) (*GetIncidentResponse, error)
	UpdateIncident(ctx context.Context, req UpdateIncidentRequest, incidentID int64) (*UpdateIncidentResponse, error)
	DeleteIncident(ctx context.Context, incidentID int64) error
	GenerateIncidentFile(ctx context.Context, incidentID int64) (*GenerateIncidentFileResponse, error)
	ConfirmIncident(ctx context.Context, incidentID int64) (*ConfirmIncidentResponse, error)
	ListAllIncidents(ctx *gin.Context, req *ListAllIncidentsRequest) (*pagination.Response[ListAllIncidentsResponse], error)

	// Client Diagnoses
	CreateClientDiagnosis(ctx context.Context, req CreateClientDiagnosisRequest, clientID uuid.UUID) (*CreateClientDiagnosisResponse, error)
	ListClientDiagnoses(ctx *gin.Context, req ListClientDiagnosesRequest, clientID uuid.UUID) (*pagination.Response[ListClientDiagnosesResponse], error)
	GetClientDiagnosis(ctx context.Context, diagnosisID int64) (*GetClientDiagnosisResponse, error)
	UpdateClientDiagnosis(ctx context.Context, req UpdateClientDiagnosisRequest, diagnosisID int64) (*UpdateClientDiagnosisResponse, error)
	DeleteClientDiagnosis(ctx context.Context, diagnosisID int64) (*DeleteClientDiagnosisResponse, error)
	// Client Medications
	CreateClientMedication(ctx context.Context, req CreateClientMedicationRequest, diagnosisID *int64) (*CreateClientMedicationResponse, error)
	ListMedicationsByDiagnosisID(ctx *gin.Context, req ListClientMedicationsRequest, diagnosisID *int64) (*pagination.Response[ListClientMedicationsResponse], error)
	GetClientMedication(ctx context.Context, medicationID int64) (*GetClientMedicationResponse, error)
	UpdateClientMedication(ctx context.Context, req UpdateClientMedicationRequest, medicationID int64) (*UpdateClientMedicationResponse, error)
	DeleteClientMedication(ctx context.Context, medicationID int64) error

	// Client Sender
	GetClientSender(ctx context.Context, clientID uuid.UUID) (*GetClientSenderResponse, error)

	// Client Emergency Contacts
	CreateClientEmergencyContact(ctx context.Context, req CreateClientEmergencyContactParams, clientID uuid.UUID) (*CreateClientEmergencyContactResponse, error)
	ListClientEmergencyContacts(ctx *gin.Context, req ListClientEmergencyContactsRequest, clientID uuid.UUID) (*pagination.Response[ListClientEmergencyContactsResponse], error)
	GetClientEmergencyContact(ctx context.Context, contactID int64) (*GetClientEmergencyContactResponse, error)
	UpdateClientEmergencyContact(ctx context.Context, req UpdateClientEmergencyContactParams, contactID int64) (*UpdateClientEmergencyContactResponse, error)
	DeleteClientEmergencyContact(ctx context.Context, contactID int64) (*DeleteClientEmergencyContactResponse, error)

	// Client Involved employees
	AssignEmployeeToClient(ctx context.Context, req AssignEmployeeRequest, clientID uuid.UUID) (*AssignEmployeeResponse, error)
	ListAssignedEmployees(ctx *gin.Context, req ListAssignedEmployeesRequest, clientID uuid.UUID) (*pagination.Response[ListAssignedEmployeesResponse], error)
	GetAssignedEmployee(ctx context.Context, assignmentID int64) (*GetAssignedEmployeeResponse, error)
	UpdateAssignedEmployee(ctx context.Context, req UpdateAssignedEmployeeRequest, assignmentID int64) (*UpdateAssignedEmployeeResponse, error)
	DeleteAssignedEmployee(ctx context.Context, assignmentID int64) (*DeleteAssignedEmployeeResponse, error)

	// Client Network Emails
	GetClientRelatedEmail(ctx context.Context, clientID uuid.UUID) (*GetClientRelatedEmailsResponse, error)

	// Client Progress Reports
	CreateProgressReport(ctx context.Context, req *CreateProgressReportRequest, clientID uuid.UUID) (*CreateProgressReportResponse, error)
	ListProgressReports(ctx *gin.Context, req *ListProgressReportsRequest, clientID uuid.UUID) (*pagination.Response[ListProgressReportsResponse], error)
	GetProgressReport(ctx context.Context, reportID int64) (*GetProgressReportResponse, error)
	UpdateProgressReport(ctx context.Context, req *UpdateProgressReportRequest, reportID int64) (*GetProgressReportResponse, error)
	DeleteProgressReport(ctx context.Context, reportID int64) error
	GenerateAutoReports(ctx context.Context, req *GenerateAutoReportsRequest, clientID uuid.UUID) (*GenerateAutoReportsResponse, error)
	ConfirmAiProgressReport(ctx context.Context, clientID uuid.UUID, req *ConfirmProgressReportRequest, reportID int64) (*ConfirmProgressReportResponse, error)
	ListAiGeneratedReports(ctx *gin.Context, req *ListAiGeneratedReportsRequest, clientID uuid.UUID) (*pagination.Response[ListAiGeneratedReportsResponse], error)

	// Registration Form
	CreateRegistrationForm(ctx context.Context, req *CreateRegistrationFormRequest) (*CreateRegistrationFormResponse, error)
	ListRegistrationForms(ctx *gin.Context, req *ListRegistrationFormsRequest) (*pagination.Response[ListRegistrationFormsResponse], error)
	GetRegistrationFormB(tx context.Context, formID int64) (*GetRegistrationFormResponse, error)
	UpdateRegistrationForm(ctx context.Context, req *UpdateRegistrationFormRequest, formID int64) (*UpdateRegistrationFormResponse, error)
	DeleteRegistrationForm(ctx context.Context, formID int64) error
	UpdateRegistrationFormStatus(ctx context.Context, req *UpdateRegistrationFormStatusRequest, formID int64, employeeID uuid.UUID) error

	// Location Transfer
	RequestLocationTransfer(ctx context.Context, clientID uuid.UUID, req LocationTransferRequest) error
	ApproveLocationTransfer(ctx context.Context, employeeID uuid.UUID, req ApproveOrRejectLocationTransferRequest) error
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
