package service

import (
	"maicare_go/async/aclient"
	"maicare_go/bucket"
	db "maicare_go/db/sqlc"
	grpclient "maicare_go/grpclient/proto"
	"maicare_go/hub"
	"maicare_go/logger"
	"maicare_go/service/ai"
	"maicare_go/service/appointment"
	"maicare_go/service/attachment"
	"maicare_go/service/audit"
	"maicare_go/service/auth"
	"maicare_go/service/care"
	clientp "maicare_go/service/client"
	contractp "maicare_go/service/contract"
	"maicare_go/service/deps"
	"maicare_go/service/ecr"
	"maicare_go/service/employees"
	"maicare_go/service/handbook"
	"maicare_go/service/invoice"
	"maicare_go/service/leave"
	"maicare_go/service/notification"
	"maicare_go/service/organization"
	"maicare_go/service/schedule"
	"maicare_go/service/sender"
	"maicare_go/service/settings"
	"maicare_go/token"
	"maicare_go/util"
)

type BusinessService struct {
	*deps.ServiceDependencies
	AuthService         auth.AuthService
	ClientService       clientp.ClientService
	EmployeeService     employees.EmployeeService
	HandbookService     handbook.HandbookService
	InvoiceService      invoice.InvoiceService
	LeaveService        leave.LeaveService
	AppointmentService  appointment.AppointmentService
	AttachmentService   attachment.AttachmentService
	ContractService     contractp.ContractService
	ECRService          ecr.ECRService
	OrganizationService organization.OrganizationService
	CarePlanService     care.CarePlanService
	NotificationService notification.NotificationService
	ScheduleService     schedule.ScheduleService
	SenderService       sender.SenderService
	SettingsService     settings.SettingsService
	AuditService        *audit.AuditService
}

func NewBusinessService(store *db.Store, tokenMaker token.Maker, logger logger.Logger, config *util.Config, b2Client bucket.ObjectStorageInterface, grpcClient grpclient.GrpcClientInterface, wsHub *hub.Hub, asynqClient aclient.AsynqClientInterface, aiService ai.AIService) *BusinessService {
	deps := deps.NewServiceDependencies(store, tokenMaker, logger, config, b2Client, grpcClient, wsHub, aiService)
	authService := auth.NewAuthService(deps)
	clientService := clientp.NewClientService(deps, asynqClient)
	employeeService := employees.NewEmployeeService(deps, asynqClient)
	handbookService := handbook.NewHandbookService(deps)
	invoiceService := invoice.NewInvoiceService(deps)
	leaveService := leave.NewLeaveService(deps)
	appointmentService := appointment.NewAppointmentService(deps, asynqClient)
	attachmentService := attachment.NewAttachmentService(deps)
	contractService := contractp.NewContractService(deps)
	ecrService := ecr.NewECRService(deps)
	organizationService := organization.NewOrganizationService(deps)
	carePlanService := care.NewCarePlanService(deps)
	notificationService := notification.NewNotificationService(deps)
	scheduleService := schedule.NewScheduleService(deps, asynqClient)
	senderService := sender.NewSenderService(deps)
	settingsService := settings.NewSettingsService(deps)
	auditService := audit.NewAuditService(store)
	return &BusinessService{
		ServiceDependencies: deps,
		AuthService:         authService,
		ClientService:       clientService,
		EmployeeService:     employeeService,
		HandbookService:     handbookService,
		InvoiceService:      invoiceService,
		LeaveService:        leaveService,
		AppointmentService:  appointmentService,
		AttachmentService:   attachmentService,
		ContractService:     contractService,
		ECRService:          ecrService,
		OrganizationService: organizationService,
		CarePlanService:     carePlanService,
		NotificationService: notificationService,
		ScheduleService:     scheduleService,
		SenderService:       senderService,
		SettingsService:     settingsService,
		AuditService:        auditService,
	}
}

// func NewMockBusinessService(ctrl *gomock.Controller) *BusinessService {
// 	authService := mocks.NewMockAuthService(ctrl)

// 	employeeService := mocks.NewMockEmployeeService(ctrl)
// 	invoiceService := mocks.NewMockInvoiceService(ctrl)

// 	return &BusinessService{
// 		AuthService:     authService,
// 		EmployeeService: employeeService,
// 		InvoiceService:  invoiceService,
// 	}
// }
