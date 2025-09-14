package service

import (
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/service/auth"
	"maicare_go/service/deps"
	"maicare_go/service/employees"
	"maicare_go/service/invoice"
	"maicare_go/token"
	"maicare_go/util"
)

type BusinessService struct {
	*deps.ServiceDependencies
	AuthService     auth.AuthService
	EmployeeService employees.EmployeeService
	InvoiceService  invoice.InvoiceService
}

func NewBusinessService(store *db.Store, tokenMaker token.Maker, logger logger.Logger, config *util.Config) *BusinessService {
	deps := deps.NewServiceDependencies(store, tokenMaker, logger, config)
	authService := auth.NewAuthService(deps)
	employeeService := employees.NewEmployeeService(deps)
	invoiceService := invoice.NewInvoiceService(deps)
	return &BusinessService{
		ServiceDependencies: deps,
		AuthService:         authService,
		EmployeeService:     employeeService,
		InvoiceService:      invoiceService,
	}
}
