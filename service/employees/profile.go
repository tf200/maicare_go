package employees

import (
	"context"
	"fmt"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *employeeService) CreateEmployee(
	req CreateEmployeeProfileRequest,
	ctx context.Context,
) (*CreateEmployeeProfileResponse, error) {
	password := util.RandomString(12)
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"CreateEmployee",
			"Failed to hash password",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to hash password")
	}

	var parsedDateOfBirth pgtype.Date
	if req.DateOfBirth != nil {
		parsedDateOfBirth, err = util.StringToPgDate(*req.DateOfBirth)
		if err != nil {
			s.Logger.LogBusinessEvent(
				ctx,
				logger.LogLevelError,
				"CreateEmployee",
				"Failed to parse date of birth",
				zap.Error(err),
			)
			return nil, fmt.Errorf("invalid date of birth format")
		}
	}
	var parsedContractStartDate, parsedContractEndDate pgtype.Date
	if req.ContractStartDate != nil {
		parsedContractStartDate, err = util.StringToPgDate(*req.ContractStartDate)
		if err != nil {
			s.Logger.LogBusinessEvent(
				ctx,
				logger.LogLevelError,
				"CreateEmployee",
				"Failed to parse contract start date",
				zap.Error(err),
			)
			return nil, fmt.Errorf("invalid contract start date format")
		}
	}
	if req.ContractEndDate != nil {
		parsedContractEndDate, err = util.StringToPgDate(*req.ContractEndDate)
		if err != nil {
			s.Logger.LogBusinessEvent(
				ctx,
				logger.LogLevelError,
				"CreateEmployee",
				"Failed to parse contract end date",
				zap.Error(err),
			)
			return nil, fmt.Errorf("invalid contract end date format")
		}
	}

	contractType := db.EmployeeContractTypeEnum(req.ContractType)

	employee, err := s.Store.CreateEmployeeWithAccountTx(
		ctx,
		db.CreateEmployeeWithAccountTxParams{
			CreateUserParams: db.CreateUserParams{
				Password: hashedPassword,
				Email:    req.WorkEmailAddress,
				IsActive: true,
			},

			CreateEmployeeParams: db.CreateEmployeeProfileParams{
				FirstName:           req.FirstName,
				LastName:            req.LastName,
				Bsn:                 req.Bsn,
				Street:              req.Street,
				HouseNumber:         req.HouseNumber,
				HouseNumberAddition: req.HouseNumberAddition,
				PostalCode:          req.PostalCode,
				City:                req.City,
				Position:            req.Position,
				Department:          req.Department,
				EmployeeNumber:      req.EmployeeNumber,
				EmploymentNumber:    req.EmploymentNumber,
				PrivateEmailAddress: req.PrivateEmailAddress,
				WorkEmailAddress:    &req.WorkEmailAddress,
				WorkPhoneNumber:     req.WorkPhoneNumber,
				PrivatePhoneNumber:  req.PrivatePhoneNumber,
				DateOfBirth:         parsedDateOfBirth,
				HomeTelephoneNumber: req.HomeTelephoneNumber,
				Gender:              db.GenderEnum(req.Gender),
				LocationID:          req.LocationID,
				ContractHours:       req.ContractHours,
				ContractType:        contractType,
				ContractStartDate:   parsedContractStartDate,
				ContractEndDate:     parsedContractEndDate,
				ContractRate:        req.ContractRate,
			},
			RoleID: req.RoleID,
		},
	)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"CreateEmployee",
			"Failed to create employee with account",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create employee with account: %v", err)
	}

	res := &CreateEmployeeProfileResponse{
		ID:                  employee.Employee.ID,
		EmployeeNumber:      employee.Employee.EmployeeNumber,
		EmploymentNumber:    employee.Employee.EmploymentNumber,
		FirstName:           employee.Employee.FirstName,
		LastName:            employee.Employee.LastName,
		DateOfBirth:         employee.Employee.DateOfBirth.Time,
		Gender:              string(employee.Employee.Gender),
		Email:               *employee.Employee.WorkEmailAddress,
		PrivateEmailAddress: employee.Employee.PrivateEmailAddress,
		WorkPhoneNumber:     employee.Employee.WorkPhoneNumber,
		PrivatePhoneNumber:  employee.Employee.PrivatePhoneNumber,
		HomeTelephoneNumber: employee.Employee.HomeTelephoneNumber,
		OutOfService:        employee.Employee.OutOfService,
		HasBorrowed:         employee.Employee.HasBorrowed,
		UserID:              employee.User.ID,
		CreatedAt:           employee.Employee.CreatedAt.Time,
		IsArchived:          employee.Employee.IsArchived,
		LocationID:          employee.Employee.LocationID,
	}
	s.Logger.LogBusinessEvent(
		ctx,
		logger.LogLevelInfo,
		"CreateEmployee",
		"Successfully created employee with account",
		zap.String("EmployeeID", res.ID.String()),
		zap.String("UserID", res.UserID.String()),
	)
	return res, nil
}

func (s *employeeService) ListEmployees(
	req ListEmployeeRequest,
	ctx *gin.Context,
) (*pagination.Response[ListEmployeeResponse], error) {
	params := req.GetParams()
	contractTypeFilter := db.NullEmployeeContractTypeEnum{}
	if req.ContractType != nil {
		contractTypeFilter = db.NullEmployeeContractTypeEnum{
			EmployeeContractTypeEnum: db.EmployeeContractTypeEnum(*req.ContractType),
			Valid:                    true,
		}
	}

	employees, err := s.Store.ListEmployeeProfile(ctx, db.ListEmployeeProfileParams{
		Limit:               params.Limit,
		Offset:              params.Offset,
		IncludeArchived:     req.IncludeArchived,
		IncludeOutOfService: req.IncludeOutOfService,
		LocationID:          req.LocationID,
		ContractType:        contractTypeFilter,
		Search:              req.Search,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListEmployees",
			"Failed to list employees", zap.Error(err))
		return nil, fmt.Errorf("failed to list employees")
	}
	totalCount, err := s.Store.CountEmployeeProfile(ctx, db.CountEmployeeProfileParams{
		IncludeArchived:     req.IncludeArchived,
		IncludeOutOfService: req.IncludeOutOfService,
		LocationID:          req.LocationID,
		ContractType:        contractTypeFilter,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListEmployees",
			"Failed to count employees", zap.Error(err))
		return nil, fmt.Errorf("failed to count employees")
	}

	responseEmployees := make([]ListEmployeeResponse, len(employees))

	for i, employee := range employees {
		var contractEndDate *time.Time
		if employee.ContractEndDate.Valid {
			parsedDate := employee.ContractEndDate.Time
			contractEndDate = &parsedDate
		}
		responseEmployees[i] = ListEmployeeResponse{
			ID:              employee.ID,
			FirstName:       employee.FirstName,
			LastName:        employee.LastName,
			Bsn:             employee.Bsn,
			ContractType:    string(employee.ContractType),
			Department:      employee.Department,
			LocationAddress: employee.LocationAddress,
			ContractEndDate: contractEndDate,
		}
	}

	response := pagination.NewResponse(ctx, req.Request, responseEmployees, totalCount)

	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "ListEmployees", "Successfully listed employees",
		zap.Int("Count", len(employees)), zap.Int64("TotalCount", totalCount))
	return &response, nil
}

func (s *employeeService) UpdateEmployeeIsSubcontractor(
	req UpdateEmployeeIsSubcontractorRequest,
	employeeID uuid.UUID,
	ctx context.Context,
) (*UpdateEmployeeIsSubcontractorResponse, error) {
	contractType := db.EmployeeContractTypeEnumLoondienst
	if req.IsSubcontractor != nil && *req.IsSubcontractor {
		contractType = db.EmployeeContractTypeEnumZZP
	}

	emp, err := s.Store.UpdateEmployeeIsSubcontractor(ctx, db.UpdateEmployeeIsSubcontractorParams{
		ID:           employeeID,
		ContractType: contractType,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"UpdateEmployeeIsSubcontractor",
			"Failed to update employee subcontractor status",
			zap.Error(err),
			zap.String("EmployeeID", employeeID.String()),
			zap.Bool("IsSubcontractor", *req.IsSubcontractor),
		)
		return nil, fmt.Errorf("failed to update employee subcontractor status")
	}
	res := &UpdateEmployeeIsSubcontractorResponse{
		ID:                emp.ID,
		ContractType:      string(emp.ContractType),
		ContractHours:     emp.ContractHours,
		ContractStartDate: emp.ContractStartDate.Time,
		ContractEndDate:   emp.ContractEndDate.Time,
		ContractRate:      emp.ContractRate,
	}
	s.Logger.LogBusinessEvent(
		ctx,
		logger.LogLevelInfo,
		"UpdateEmployeeIsSubcontractor",
		"Successfully updated employee subcontractor status",
		zap.String("EmployeeID", employeeID.String()),
		zap.Bool("IsSubcontractor", *req.IsSubcontractor),
	)
	return res, nil
}

func (s *employeeService) GetEmployeeProfile(
	userID uuid.UUID,
	ctx context.Context,
) (*GetEmployeeProfileResponse, error) {
	profile, err := s.Store.GetEmployeeProfileByUserID(ctx, userID)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"GetEmployeeProfile",
			"Failed to get employee profile by user ID",
			zap.Error(err),
			zap.String("UserID", userID.String()),
		)
		return nil, fmt.Errorf("failed to get employee profile: %w", err)
	}

	var permissions []Permission
	if err := json.Unmarshal(profile.Permissions, &permissions); err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"GetEmployeeProfile",
			"Failed to unmarshal permissions",
			zap.Error(err),
			zap.String("UserID", userID.String()),
		)
		return nil, fmt.Errorf("failed to parse permissions: %w", err)
	}

	res := &GetEmployeeProfileResponse{
		UserID:      profile.UserID,
		EmployeeID:  profile.EmployeeID,
		FirstName:   profile.FirstName,
		LastName:    profile.LastName,
		Email:       profile.Email,
		TwoFactor:   profile.TwoFactorEnabled,
		LastLogin:   profile.LastLogin.Time,
		Permissions: permissions,
	}

	s.Logger.LogBusinessEvent(
		ctx,
		logger.LogLevelInfo,
		"GetEmployeeProfile",
		"Successfully retrieved employee profile",
		zap.String("UserID", userID.String()),
		zap.String("EmployeeID", profile.EmployeeID.String()),
	)
	return res, nil
}

func (s *employeeService) GetEmployeeProfileByID(
	employeeID, currentUserID uuid.UUID,
	ctx context.Context,
) (*GetEmployeeProfileByIDResponse, error) {
	employee, err := s.Store.GetEmployeeProfileByID(ctx, employeeID)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"GetEmployeeProfileByID",
			"Failed to get employee profile by ID",
			zap.Error(err),
			zap.String("EmployeeID", employeeID.String()),
		)
		return nil, fmt.Errorf("failed to get employee profile: %w", err)
	}

	res := &GetEmployeeProfileByIDResponse{
		ID:                  employee.ID,
		UserID:              employee.UserID,
		FirstName:           employee.FirstName,
		LastName:            employee.LastName,
		Position:            employee.Position,
		Department:          employee.Department,
		EmployeeNumber:      employee.EmployeeNumber,
		EmploymentNumber:    employee.EmploymentNumber,
		PrivateEmailAddress: employee.PrivateEmailAddress,
		PrivatePhoneNumber:  employee.PrivatePhoneNumber,
		WorkPhoneNumber:     employee.WorkPhoneNumber,
		DateOfBirth:         employee.DateOfBirth.Time,
		HomeTelephoneNumber: employee.HomeTelephoneNumber,
		CreatedAt:           employee.CreatedAt.Time,
		Gender:              string(employee.Gender),
		LocationID:          employee.LocationID,
		HasBorrowed:         employee.HasBorrowed,
		OutOfService:        employee.OutOfService,
		IsArchived:          employee.IsArchived,
		ProfilePicture:      s.GenerateResponsePresignedURL(employee.ProfilePicture, ctx),
		IsLoggedInUser:      employee.UserID == currentUserID,
	}

	s.Logger.LogBusinessEvent(
		ctx,
		logger.LogLevelInfo,
		"GetEmployeeProfileByID",
		"Successfully retrieved employee profile by ID",
		zap.String("EmployeeID", employeeID.String()),
		zap.String("UserID", employee.UserID.String()),
	)
	return res, nil
}

func (s *employeeService) UpdateEmployeeProfile(
	req UpdateEmployeeProfileRequest,
	employeeID uuid.UUID,
	ctx context.Context,
) (*UpdateEmployeeProfileResponse, error) {
	var parsedDate time.Time
	var err error
	if req.DateOfBirth != nil {
		parsedDate, err = time.Parse("2006-01-02", *req.DateOfBirth)
		if err != nil {
			s.Logger.LogBusinessEvent(
				ctx,
				logger.LogLevelError,
				"UpdateEmployeeProfile",
				"Failed to parse date of birth",
				zap.Error(err),
				zap.String("EmployeeID", employeeID.String()),
			)
			return nil, fmt.Errorf("invalid date of birth format: %w", err)
		}
	}

	employee, err := s.Store.UpdateEmployeeProfile(ctx, db.UpdateEmployeeProfileParams{
		ID:                  employeeID,
		FirstName:           req.FirstName,
		LastName:            req.LastName,
		Position:            req.Position,
		Department:          req.Department,
		EmployeeNumber:      req.EmployeeNumber,
		EmploymentNumber:    req.EmploymentNumber,
		PrivateEmailAddress: req.PrivateEmailAddress,
		PrivatePhoneNumber:  req.PrivatePhoneNumber,
		WorkPhoneNumber:     req.WorkPhoneNumber,
		DateOfBirth:         pgtype.Date{Time: parsedDate, Valid: true},
		HomeTelephoneNumber: req.HomeTelephoneNumber,
		Gender:              db.NullGenderFromPtr(req.Gender),
		LocationID:          req.LocationID,
		HasBorrowed:         req.HasBorrowed,
		OutOfService:        req.OutOfService,
		IsArchived:          req.IsArchived,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"UpdateEmployeeProfile",
			"Failed to update employee profile",
			zap.Error(err),
			zap.String("EmployeeID", employeeID.String()),
		)
		return nil, fmt.Errorf("failed to update employee profile: %w", err)
	}

	res := &UpdateEmployeeProfileResponse{
		ID:                  employee.ID,
		UserID:              employee.UserID,
		FirstName:           employee.FirstName,
		LastName:            employee.LastName,
		Position:            employee.Position,
		Department:          employee.Department,
		EmployeeNumber:      employee.EmployeeNumber,
		EmploymentNumber:    employee.EmploymentNumber,
		PrivateEmailAddress: employee.PrivateEmailAddress,
		PrivatePhoneNumber:  employee.PrivatePhoneNumber,
		WorkPhoneNumber:     employee.WorkPhoneNumber,
		DateOfBirth:         employee.DateOfBirth.Time,
		HomeTelephoneNumber: employee.HomeTelephoneNumber,
		CreatedAt:           employee.CreatedAt.Time,
		Gender:              string(employee.Gender),
		LocationID:          employee.LocationID,
		HasBorrowed:         employee.HasBorrowed,
		OutOfService:        employee.OutOfService,
		IsArchived:          employee.IsArchived,
	}

	s.Logger.LogBusinessEvent(
		ctx,
		logger.LogLevelInfo,
		"UpdateEmployeeProfile",
		"Successfully updated employee profile",
		zap.String("EmployeeID", employeeID.String()),
	)
	return res, nil
}

func (s *employeeService) SetEmployeeProfilePicture(
	req SetEmployeeProfilePictureRequest,
	employeeID uuid.UUID,
	ctx context.Context,
) (*SetEmployeeProfilePictureResponse, error) {
	attachmentID, err := uuid.Parse(req.AttachmentID)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"SetEmployeeProfilePicture",
			"Failed to parse attachment ID",
			zap.Error(err),
			zap.String("EmployeeID", employeeID.String()),
		)
		return nil, fmt.Errorf("invalid attachment ID format: %w", err)
	}

	arg := db.SetEmployeeProfilePictureTxParams{
		EmployeeID:    employeeID,
		AttachementID: attachmentID,
	}
	user, err := s.Store.SetEmployeeProfilePictureTx(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"SetEmployeeProfilePicture",
			"Failed to set employee profile picture",
			zap.Error(err),
			zap.String("EmployeeID", employeeID.String()),
		)
		return nil, fmt.Errorf("failed to set employee profile picture: %w", err)
	}

	res := &SetEmployeeProfilePictureResponse{
		ID:             user.User.ID,
		Email:          user.User.Email,
		ProfilePicture: user.User.ProfilePicture,
	}

	s.Logger.LogBusinessEvent(
		ctx,
		logger.LogLevelInfo,
		"SetEmployeeProfilePicture",
		"Successfully set employee profile picture",
		zap.String("EmployeeID", employeeID.String()),
	)
	return res, nil
}

func (s *employeeService) GetEmployeeCounts(
	ctx context.Context,
) (*GetEmployeeCountsResponse, error) {
	counts, err := s.Store.GetEmployeeCounts(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"GetEmployeeCounts",
			"Failed to get employee counts",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get employee counts: %w", err)
	}

	res := &GetEmployeeCountsResponse{
		TotalEmployees:      counts.TotalEmployees,
		TotalSubcontractors: counts.TotalSubcontractors,
		TotalArchived:       counts.TotalArchived,
		TotalOutOfService:   counts.TotalOutOfService,
	}

	s.Logger.LogBusinessEvent(
		ctx,
		logger.LogLevelInfo,
		"GetEmployeeCounts",
		"Successfully retrieved employee counts",
	)
	return res, nil
}

func (s *employeeService) SearchEmployeesByNameOrEmail(
	req SearchEmployeesByNameOrEmailRequest,
	ctx context.Context,
) ([]SearchEmployeesByNameOrEmailResponse, error) {
	employees, err := s.Store.SearchEmployeesByNameOrEmail(ctx, req.Search)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"SearchEmployeesByNameOrEmail",
			"Failed to search employees",
			zap.Error(err),
			zap.String("Search", *req.Search),
		)
		return nil, fmt.Errorf("failed to search employees: %w", err)
	}

	responseEmployees := make([]SearchEmployeesByNameOrEmailResponse, len(employees))
	for i, employee := range employees {
		responseEmployees[i] = SearchEmployeesByNameOrEmailResponse{
			ID:        employee.ID,
			FirstName: employee.FirstName,
			LastName:  employee.LastName,
		}
	}

	s.Logger.LogBusinessEvent(
		ctx,
		logger.LogLevelInfo,
		"SearchEmployeesByNameOrEmail",
		"Successfully searched employees",
		zap.Int("Count", len(employees)),
	)
	return responseEmployees, nil
}
