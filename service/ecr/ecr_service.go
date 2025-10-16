package ecr

import (
	"context"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"
	"maicare_go/service/deps"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ecrService struct {
	*deps.ServiceDependencies
}

func NewECRService(deps *deps.ServiceDependencies) ECRService {
	return &ecrService{
		ServiceDependencies: deps,
	}
}

func (s *ecrService) DischargeOverview(ctx *gin.Context, req DischargeOverviewRequest) (*pagination.Response[DischargeOverviewResponse], error) {
	params := req.GetParams()

	overview, err := s.Store.DischargeOverview(ctx, db.DischargeOverviewParams{
		Limit:      params.Limit,
		Offset:     params.Offset,
		FilterType: req.FilterType,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "DischargeOverview", "Failed to get discharge overview", zap.Error(err))
		return nil, err
	}

	overviewRes := []DischargeOverviewResponse{}
	for _, item := range overview {
		overviewRes = append(overviewRes, DischargeOverviewResponse{
			ID:                 item.ID,
			FirstName:          item.FirstName,
			LastName:           item.LastName,
			CurrentStatus:      item.CurrentStatus,
			ScheduledStatus:    item.ScheduledStatus,
			StatusChangeReason: item.StatusChangeReason,
			StatusChangeDate:   item.StatusChangeDate.Time,
			ContractEndDate:    item.ContractEndDate.Time,
			ContractStatus:     item.ContractStatus,
			DepartureReason:    item.DepartureReason,
			FollowUpPlan:       item.FollowUpPlan,
			DischargeType:      item.DischargeType,
		})
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "DischargeOverview", fmt.Sprintf("Retrieved %d discharge overview records", len(overviewRes)), zap.Int("record_count", len(overviewRes)))

	count, err := s.Store.TotalDischargeCount(context.Background())
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "DischargeOverview", "Failed to get total discharge count", zap.Error(err))
		return nil, err
	}

	pag := pagination.NewResponse(ctx, req.Request, overviewRes, count)
	return &pag, nil
}

func (s *ecrService) TotalDischargeCount(ctx *gin.Context) (*TotalDischargeCountResponse, error) {
	// Initialize WaitGroup to track 4 Go routines
	var wg sync.WaitGroup
	wg.Add(4)

	// Variables to hold results and errors for each query
	var totalDischargeCount int64
	var totalDischargeErr error
	var urgentCasesCount int64
	var urgentCasesErr error
	var statusChangeCount int64
	var statusChangeErr error
	var contractEndingCount int64
	var contractEndingErr error

	// Get the request context to respect timeouts/cancellations
	reqCtx := ctx.Request.Context()

	// Go routine for TotalDischargeCount
	go func() {
		defer wg.Done() // Signal completion when done
		count, err := s.Store.TotalDischargeCount(reqCtx)
		totalDischargeCount = count
		totalDischargeErr = err
	}()

	// Go routine for UrgentCasesCount
	go func() {
		defer wg.Done()
		count, err := s.Store.UrgentCasesCount(reqCtx)
		urgentCasesCount = count
		urgentCasesErr = err
	}()

	// Go routine for StatusChangeCount
	go func() {
		defer wg.Done()
		count, err := s.Store.StatusChangeCount(reqCtx)
		statusChangeCount = count
		statusChangeErr = err
	}()

	// Go routine for ContractEndCount
	go func() {
		defer wg.Done()
		count, err := s.Store.ContractEndCount(reqCtx)
		contractEndingCount = count
		contractEndingErr = err
	}()

	// Wait for all Go routines to complete
	wg.Wait()

	// Check errors and return the first one encountered
	if totalDischargeErr != nil {
		return nil, totalDischargeErr
	}
	if urgentCasesErr != nil {
		return nil, urgentCasesErr
	}
	if statusChangeErr != nil {
		return nil, statusChangeErr
	}
	if contractEndingErr != nil {
		return nil, contractEndingErr
	}

	// All queries succeeded, construct and send the response
	response := &TotalDischargeCountResponse{
		TotalDischargeCount: totalDischargeCount,
		UrgentCasesCount:    urgentCasesCount,
		StatusChangeCount:   statusChangeCount,
		ContractEndCount:    contractEndingCount,
	}

	return response, nil
}

func (s *ecrService) ListEmployeesByContractEndDate(ctx context.Context) ([]ListEmployeesByContractEndDateResponse, error) {
	employees, err := s.Store.ListEmployeesByContractEndDate(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListEmployeesByContractEndDate", "Failed to list employees by contract end date", zap.Error(err))
		return nil, err
	}

	response := []ListEmployeesByContractEndDateResponse{}
	for _, emp := range employees {
		response = append(response, ListEmployeesByContractEndDateResponse{
			ID:                emp.ID,
			FirstName:         emp.FirstName,
			LastName:          emp.LastName,
			Position:          emp.Position,
			Department:        emp.Department,
			EmployeeNumber:    emp.EmployeeNumber,
			EmploymentNumber:  emp.EmploymentNumber,
			Email:             emp.Email,
			ContractStartDate: emp.ContractStartDate.Time,
			ContractEndDate:   emp.ContractEndDate.Time,
			ContractType:      emp.ContractType,
		})
	}

	return response, nil
}

func (s *ecrService) ListLatestPayments(ctx context.Context) ([]ListLatestPaymentsResponse, error) {
	latestPayments, err := s.Store.ListLatestPayments(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListLatestPayments", "Failed to list latest payments", zap.Error(err))
		return nil, err
	}

	response := []ListLatestPaymentsResponse{}
	for _, payment := range latestPayments {
		response = append(response, ListLatestPaymentsResponse{
			InvoiceID:     payment.InvoiceID,
			InvoiceNumber: payment.InvoiceNumber,
			PaymentMethod: payment.PaymentMethod,
			PaymentStatus: payment.PaymentStatus,
			Amount:        payment.Amount,
			PaymentDate:   payment.PaymentDate.Time,
			UpdatedAt:     payment.UpdatedAt.Time,
		})
	}

	return response, nil
}

func (s *ecrService) ListUpcomingAppointments(ctx context.Context, employeeID int64) ([]ListUpcomingAppointmentsResponse, error) {
	appointments, err := s.Store.ListUpcomingAppointments(ctx, &employeeID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListUpcomingAppointments", "Failed to list upcoming appointments", zap.Error(err))
		return nil, err
	}

	response := []ListUpcomingAppointmentsResponse{}
	for _, appt := range appointments {
		response = append(response, ListUpcomingAppointmentsResponse{
			ID:          appt.ID,
			StartTime:   appt.StartTime.Time,
			EndTime:     appt.EndTime.Time,
			Location:    appt.Location,
			Description: appt.Description,
		})
	}

	return response, nil
}