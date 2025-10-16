package clientp

import (
	"context"
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
	GetClientDetails(ctx context.Context, clientID int64) (*GetClientApiResponse, error)
	GetClientAddresses(ctx context.Context, clientID int64) (*GetClientAddressesApiResponse, error)
	UpdateClientDetails(ctx context.Context, req UpdateClientDetailsRequest, clientID int64) (*UpdateClientDetailsResponse, error)
	UpdateClientStatus(ctx context.Context, req UpdateClientStatusRequest, clientID int64) (*UpdateClientStatusResponse, error)
	ListStatusHistory(ctx context.Context, clientID int64) ([]ListStatusHistoryApiResponse, error)
	SetClientProfilePicture(ctx context.Context, req SetClientProfilePictureRequest, clientID int64) (*SetClientProfilePictureResponse, error)
	// Client Documents
	AddClientDocument(ctx context.Context, req AddClientDocumentApiRequest, clientID int64) (*AddClientDocumentApiResponse, error)
	ListClientDocuments(ctx *gin.Context, req ListClientDocumentsApiRequest, clientID int64) (*pagination.Response[ListClientDocumentsApiResponse], error)
	DeleteClientDocument(ctx context.Context, clientID int64, attachmentID uuid.UUID) (*DeleteClientDocumentApiResponse, error)
	GetMissingClientDocuments(ctx context.Context, clientID int64) (*GetMissingClientDocumentsApiResponse, error)

	// Client Appointment Card
	CreateAppointmentCard(req CreateAppointmentCardRequest, clientID int64, ctx context.Context) (*CreateAppointmentCardResponse, error)
	GetAppointmentCard(ctx context.Context, clientID int64) (*GetAppointmentCardResponse, error)
	UpdateAppointmentCard(req UpdateAppointmentCardRequest, clientID int64, ctx context.Context) (*UpdateAppointmentCardResponse, error)
	GenerateAppointmentCardDocumentApi(ctx context.Context, clientID int64) (*GenerateAppointmentCardDocumentApiResponse, error)

	// Client Incidents
	CreateIncident(ctx context.Context, req CreateIncidentRequest, clientID int64) (*CreateIncidentResponse, error)
	ListIncidents(ctx *gin.Context, req ListIncidentsRequest, clientID int64) (*pagination.Response[ListIncidentsResponse], error)
	GetIncident(ctx context.Context, incidentID int64) (*GetIncidentResponse, error)
	UpdateIncident(ctx context.Context, req UpdateIncidentRequest, incidentID int64) (*UpdateIncidentResponse, error)
	DeleteIncident(ctx context.Context, incidentID int64) error
	GenerateIncidentFile(ctx context.Context, incidentID int64) (*GenerateIncidentFileResponse, error)
	ConfirmIncident(ctx context.Context, incidentID int64) (*ConfirmIncidentResponse, error)

	// Client Diagnoses
	CreateClientDiagnosis(ctx context.Context, req CreateClientDiagnosisRequest, clientID int64) (*CreateClientDiagnosisResponse, error)
	ListClientDiagnoses(ctx *gin.Context, req ListClientDiagnosesRequest, clientID int64) (*pagination.Response[ListClientDiagnosesResponse], error)
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
	GetClientSender(ctx context.Context, clientID int64) (*GetClientSenderResponse, error)

	// Client Emergency Contacts
	CreateClientEmergencyContact(ctx context.Context, req CreateClientEmergencyContactParams, clientID int64) (*CreateClientEmergencyContactResponse, error)
	ListClientEmergencyContacts(ctx *gin.Context, req ListClientEmergencyContactsRequest, clientID int64) (*pagination.Response[ListClientEmergencyContactsResponse], error)
	GetClientEmergencyContact(ctx context.Context, contactID int64) (*GetClientEmergencyContactResponse, error)
	UpdateClientEmergencyContact(ctx context.Context, req UpdateClientEmergencyContactParams, contactID int64) (*UpdateClientEmergencyContactResponse, error)
	DeleteClientEmergencyContact(ctx context.Context, contactID int64) (*DeleteClientEmergencyContactResponse, error)

	// Client Involved employees
	AssignEmployeeToClient(ctx context.Context, req AssignEmployeeRequest, clientID int64) (*AssignEmployeeResponse, error)
	ListAssignedEmployees(ctx *gin.Context, req ListAssignedEmployeesRequest, clientID int64) (*pagination.Response[ListAssignedEmployeesResponse], error)
	GetAssignedEmployee(ctx context.Context, assignmentID int64) (*GetAssignedEmployeeResponse, error)
	UpdateAssignedEmployee(ctx context.Context, req UpdateAssignedEmployeeRequest, assignmentID int64) (*UpdateAssignedEmployeeResponse, error)
	DeleteAssignedEmployee(ctx context.Context, assignmentID int64) (*DeleteAssignedEmployeeResponse, error)

	// Client Network Emails
	GetClientRelatedEmail(ctx context.Context, clientID int64) (*GetClientRelatedEmailsResponse, error)

	// Client Progress Reports
	CreateProgressReport(ctx context.Context, req *CreateProgressReportRequest, clientID int64) (*CreateProgressReportResponse, error)
	ListProgressReports(ctx *gin.Context, req *ListProgressReportsRequest, clientID int64) (*pagination.Response[ListProgressReportsResponse], error)
	GetProgressReport(ctx context.Context, reportID int64) (*GetProgressReportResponse, error)
	UpdateProgressReport(ctx context.Context, req *UpdateProgressReportRequest, reportID int64) (*GetProgressReportResponse, error)
	DeleteProgressReport(ctx context.Context, reportID int64) error
}

type clientService struct {
	*deps.ServiceDependencies
}

func NewClientService(deps *deps.ServiceDependencies) ClientService {
	return &clientService{
		ServiceDependencies: deps,
	}
}
