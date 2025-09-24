package employees

import (
	"context"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"
	"maicare_go/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *employeeService) CreateEmployee(req CreateEmployeeProfileRequest, ctx context.Context) (*CreateEmployeeProfileResponse, error) {
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

	res := &CreateEmployeeProfileResponse{
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

func (s *employeeService) ListEmployees(req ListEmployeeRequest, ctx *gin.Context) (*pagination.Response[ListEmployeeResponse], error) {
	params := req.GetParams()
	employees, err := s.Store.ListEmployeeProfile(ctx, db.ListEmployeeProfileParams{
		Limit:               params.Limit,
		Offset:              params.Offset,
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
		return nil, fmt.Errorf("failed to list employees")
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
		return nil, fmt.Errorf("failed to count employees")
	}

	responseEmployees := make([]ListEmployeeResponse, len(employees))

	for i, employee := range employees {
		responseEmployees[i] = ListEmployeeResponse{
			ID:                        employee.ID,
			UserID:                    employee.UserID,
			FirstName:                 employee.FirstName,
			LastName:                  employee.LastName,
			Position:                  employee.Position,
			Department:                employee.Department,
			EmployeeNumber:            employee.EmployeeNumber,
			EmploymentNumber:          employee.EmploymentNumber,
			PrivateEmailAddress:       employee.PrivateEmailAddress,
			Email:                     employee.Email,
			AuthenticationPhoneNumber: employee.AuthenticationPhoneNumber,
			PrivatePhoneNumber:        employee.PrivatePhoneNumber,
			WorkPhoneNumber:           employee.WorkPhoneNumber,
			DateOfBirth:               employee.DateOfBirth.Time,
			HomeTelephoneNumber:       employee.HomeTelephoneNumber,
			CreatedAt:                 employee.CreatedAt.Time,
			IsSubcontractor:           employee.IsSubcontractor,
			Gender:                    employee.Gender,
			LocationID:                employee.LocationID,
			HasBorrowed:               employee.HasBorrowed,
			OutOfService:              employee.OutOfService,
			IsArchived:                employee.IsArchived,
			ProfilePicture:            s.GenerateResponsePresignedURL(employee.ProfilePicture, ctx),
			Age:                       int64(time.Since(employee.DateOfBirth.Time).Hours() / 24 / 365),
			RoleID:                    employee.RoleID,
			RoleName:                  employee.RoleName,
		}
	}

	response := pagination.NewResponse(ctx, req.Request, responseEmployees, totalCount)

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "ListEmployees", "Successfully listed employees",
		zap.Int("Count", len(employees)), zap.Int64("TotalCount", totalCount))
	return &response, nil
}

func (s *employeeService) UpdateEmployeeIsSubcontractor(
	req UpdateEmployeeIsSubcontractorRequest,
	employeeID int64,
	ctx context.Context) (*UpdateEmployeeIsSubcontractorResponse, error) {
	contractType := "loondienst"
	if req.IsSubcontractor != nil && *req.IsSubcontractor {
		contractType = "ZZP"
	}

	emp, err := s.Store.UpdateEmployeeIsSubcontractor(ctx, db.UpdateEmployeeIsSubcontractorParams{
		ID:              employeeID,
		IsSubcontractor: req.IsSubcontractor,
		ContractType:    &contractType,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateEmployeeIsSubcontractor", "Failed to update employee subcontractor status", zap.Error(err), zap.Int64("EmployeeID", employeeID), zap.Bool("IsSubcontractor", *req.IsSubcontractor))
		return nil, fmt.Errorf("failed to update employee subcontractor status")
	}
	res := &UpdateEmployeeIsSubcontractorResponse{
		ID:                emp.ID,
		IsSubcontractor:   emp.IsSubcontractor,
		ContractType:      emp.ContractType,
		ContractHours:     emp.ContractHours,
		ContractStartDate: emp.ContractStartDate.Time,
		ContractEndDate:   emp.ContractEndDate.Time,
		ContractRate:      emp.ContractRate,
	}
	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "UpdateEmployeeIsSubcontractor", "Successfully updated employee subcontractor status", zap.Int64("EmployeeID", employeeID), zap.Bool("IsSubcontractor", *req.IsSubcontractor))
	return res, nil
}
