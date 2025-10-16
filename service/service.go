package service

import (
	"maicare_go/bucket"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/service/appointment"
	"maicare_go/service/attachment"
	"maicare_go/service/auth"
	clientp "maicare_go/service/client"
	contractp "maicare_go/service/contract"
	"maicare_go/service/deps"
	"maicare_go/service/ecr"
	"maicare_go/service/employees"
	"maicare_go/service/invoice"
	"maicare_go/token"
	"maicare_go/util"
)

type BusinessService struct {
	*deps.ServiceDependencies
	AuthService        auth.AuthService
	ClientService      clientp.ClientService
	EmployeeService    employees.EmployeeService
	InvoiceService     invoice.InvoiceService
	AppointmentService appointment.AppointmentService
	AttachmentService  attachment.AttachmentService
	ContractService    contractp.ContractService
	ECRService         ecr.ECRService
}

func NewBusinessService(store *db.Store, tokenMaker token.Maker, logger logger.Logger, config *util.Config, b2Client bucket.ObjectStorageInterface) *BusinessService {
	deps := deps.NewServiceDependencies(store, tokenMaker, logger, config, b2Client)
	authService := auth.NewAuthService(deps)
	clientService := clientp.NewClientService(deps)
	employeeService := employees.NewEmployeeService(deps)
	invoiceService := invoice.NewInvoiceService(deps)
	appointmentService := appointment.NewAppointmentService(deps)
	attachmentService := attachment.NewAttachmentService(deps)
	contractService := contractp.NewContractService(deps)
	ecrService := ecr.NewECRService(deps)
	return &BusinessService{
		ServiceDependencies: deps,
		AuthService:         authService,
		ClientService:       clientService,
		EmployeeService:     employeeService,
		InvoiceService:      invoiceService,
		AppointmentService:  appointmentService,
		AttachmentService:   attachmentService,
		ContractService:     contractService,
		ECRService:          ecrService,
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
