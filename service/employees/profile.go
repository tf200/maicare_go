package employees

import (
	"context"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/util"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type CreateEmployeeRequest struct {
	EmployeeNumber            *string
	EmploymentNumber          *string
	LocationID                *int64
	IsSubcontractor           *bool
	FirstName                 string
	LastName                  string
	DateOfBirth               *string
	Gender                    *string
	Email                     string
	PrivateEmailAddress       *string
	AuthenticationPhoneNumber *string
	WorkPhoneNumber           *string
	PrivatePhoneNumber        *string
	HomeTelephoneNumber       *string
	RoleID                    int32
	Position                  *string
	Department                *string
}

type CreateEmployeeResult struct {
	ID                        int64
	UserID                    int64
	FirstName                 string
	LastName                  string
	Position                  *string
	Department                *string
	EmployeeNumber            *string
	EmploymentNumber          *string
	PrivateEmailAddress       *string
	Email                     string
	AuthenticationPhoneNumber *string
	PrivatePhoneNumber        *string
	WorkPhoneNumber           *string
	DateOfBirth               time.Time
	HomeTelephoneNumber       *string
	CreatedAt                 time.Time
	IsSubcontractor           *bool
	Gender                    *string
	LocationID                *int64
	HasBorrowed               bool
	OutOfService              *bool
	IsArchived                bool
}

func (s *employeeService) CreateEmployee(req CreateEmployeeRequest, ctx context.Context) (*CreateEmployeeResult, error) {
	password := util.RandomString(12)
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateEmployee", "Failed to hash password", zap.Error(err))
		return nil, fmt.Errorf("failed to hash password")
	}

	var parsedDateOfBirth time.Time
	if req.DateOfBirth != nil {
		parsedDateOfBirth, err = time.Parse("2006-01-02", *req.DateOfBirth)
		if err != nil {
			s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateEmployee", "Failed to parse date of birth", zap.Error(err))
			return nil, fmt.Errorf("invalid date of birth format")
		}
	}

	contractType := "loondienst"
	if req.IsSubcontractor != nil && *req.IsSubcontractor {
		contractType = "ZZP"
	}

	employee, err := s.Store.CreateEmployeeWithAccountTx(
		ctx,
		db.CreateEmployeeWithAccountTxParams{
			CreateUserParams: db.CreateUserParams{
				Password: hashedPassword,
				Email:    req.Email,
				IsActive: true,
			},

			CreateEmployeeParams: db.CreateEmployeeProfileParams{
				FirstName:                 req.FirstName,
				LastName:                  req.LastName,
				EmployeeNumber:            req.EmployeeNumber,
				EmploymentNumber:          req.EmploymentNumber,
				LocationID:                req.LocationID,
				IsSubcontractor:           req.IsSubcontractor,
				DateOfBirth:               pgtype.Date{Time: parsedDateOfBirth, Valid: true},
				Gender:                    req.Gender,
				Email:                     req.Email,
				PrivateEmailAddress:       req.PrivateEmailAddress,
				AuthenticationPhoneNumber: req.AuthenticationPhoneNumber,
				WorkPhoneNumber:           req.WorkPhoneNumber,
				PrivatePhoneNumber:        req.PrivatePhoneNumber,
				HomeTelephoneNumber:       req.HomeTelephoneNumber,
				Position:                  req.Position,
				Department:                req.Department,
				ContractType:              &contractType,
			},
			RoleID: req.RoleID,
		},
	)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateEmployee", "Failed to create employee with account", zap.Error(err))
		return nil, fmt.Errorf("failed to create employee with account: %v", err)
	}

	res := &CreateEmployeeResult{
		ID:                        employee.Employee.ID,
		EmployeeNumber:            employee.Employee.EmployeeNumber,
		EmploymentNumber:          employee.Employee.EmploymentNumber,
		FirstName:                 employee.Employee.FirstName,
		LastName:                  employee.Employee.LastName,
		IsSubcontractor:           employee.Employee.IsSubcontractor,
		DateOfBirth:               employee.Employee.DateOfBirth.Time,
		Gender:                    employee.Employee.Gender,
		Email:                     employee.Employee.Email,
		PrivateEmailAddress:       employee.Employee.PrivateEmailAddress,
		AuthenticationPhoneNumber: employee.Employee.AuthenticationPhoneNumber,
		WorkPhoneNumber:           employee.Employee.WorkPhoneNumber,
		PrivatePhoneNumber:        employee.Employee.PrivatePhoneNumber,
		HomeTelephoneNumber:       employee.Employee.HomeTelephoneNumber,
		OutOfService:              employee.Employee.OutOfService,
		HasBorrowed:               employee.Employee.HasBorrowed,
		UserID:                    employee.User.ID,
		CreatedAt:                 employee.Employee.CreatedAt.Time,
		IsArchived:                employee.Employee.IsArchived,
		LocationID:                employee.Employee.LocationID,
	}
	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "CreateEmployee", "Successfully created employee with account", zap.Int64("EmployeeID", res.ID), zap.Int64("UserID", res.UserID))
	return res, nil
}

type ListEmployeesRequest struct {
	Limit               int32
	Offset              int32
	IncludeArchived     *bool
	IncludeOutOfService *bool
	Department          *string
	Position            *string
	LocationID          *int64
	Search              *string
}

func (s *employeeService) ListEmployees(req ListEmployeesRequest, ctx context.Context) ([]db.ListEmployeeProfileRow, *int64, error) {
	employees, err := s.Store.ListEmployeeProfile(ctx, db.ListEmployeeProfileParams{
		Limit:               req.Limit,
		Offset:              req.Offset,
		IncludeArchived:     req.IncludeArchived,
		IncludeOutOfService: req.IncludeOutOfService,
		Department:          req.Department,
		Position:            req.Position,
		LocationID:          req.LocationID,
		Search:              req.Search,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListEmployees",
			"Failed to list employees", zap.Error(err))
		return nil, nil, fmt.Errorf("failed to list employees")
	}
	totalCount, err := s.Store.CountEmployeeProfile(ctx, db.CountEmployeeProfileParams{
		IncludeArchived:     req.IncludeArchived,
		IncludeOutOfService: req.IncludeOutOfService,
		Department:          req.Department,
		Position:            req.Position,
		LocationID:          req.LocationID,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListEmployees",
			"Failed to count employees", zap.Error(err))
		return nil, nil, fmt.Errorf("failed to count employees")
	}
	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "ListEmployees", "Successfully listed employees",
		zap.Int("Count", len(employees)), zap.Int64("TotalCount", totalCount))
	return employees, &totalCount, nil
}

type UpdateEmployeeIsSubcontractorRequest struct {
	EmployeeID      int64
	IsSubcontractor *bool
}

type UpdateEmployeeIsSubcontractorResult struct {
	ID                int64
	IsSubcontractor   *bool
	ContractType      *string
	ContractHours     *float64
	ContractStartDate time.Time
	ContractEndDate   time.Time
	ContractRate      *float64
}

func (s *employeeService) UpdateEmployeeIsSubcontractor(
	req UpdateEmployeeIsSubcontractorRequest,
	ctx context.Context) (*UpdateEmployeeIsSubcontractorResult, error) {
	contractType := "loondienst"
	if req.IsSubcontractor != nil && *req.IsSubcontractor {
		contractType = "ZZP"
	}

	emp, err := s.Store.UpdateEmployeeIsSubcontractor(ctx, db.UpdateEmployeeIsSubcontractorParams{
		ID:              req.EmployeeID,
		IsSubcontractor: req.IsSubcontractor,
		ContractType:    &contractType,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateEmployeeIsSubcontractor", "Failed to update employee subcontractor status", zap.Error(err), zap.Int64("EmployeeID", req.EmployeeID))
		return nil, fmt.Errorf("failed to update employee subcontractor status")
	}
	res := &UpdateEmployeeIsSubcontractorResult{
		ID:                emp.ID,
		IsSubcontractor:   emp.IsSubcontractor,
		ContractType:      emp.ContractType,
		ContractHours:     emp.ContractHours,
		ContractStartDate: emp.ContractStartDate.Time,
		ContractEndDate:   emp.ContractEndDate.Time,
		ContractRate:      emp.ContractRate,
	}
	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "UpdateEmployeeIsSubcontractor", "Successfully updated employee subcontractor status", zap.Int64("EmployeeID", req.EmployeeID), zap.Bool("IsSubcontractor", *req.IsSubcontractor))
	return res, nil
}
