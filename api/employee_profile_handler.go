package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type Permission struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Resource string `json:"resource"`
	Method   string `json:"method"`
}

// GetEmployeeProfileResponse represents the response for GetEmployeeProfileApi
type GetEmployeeProfileResponse struct {
	UserID      int64        `json:"user_id"`
	Email       string       `json:"email"`
	EmployeeID  int64        `json:"employee_id"`
	FirstName   string       `json:"first_name"`
	LastName    string       `json:"last_name"`
	RoleID      int32        `json:"role_id"`
	Permissions []Permission `json:"permissions"`
}

// @Summary Get employee profile by user ID
// @Description Get employee profile by user ID
// @Tags employees
// @Produce json
// @Success 200 {object} Response[GetEmployeeProfileResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/profile [get]
func (server *Server) GetEmployeeProfileApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	log.Printf("Payload: %v", payload)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err)) // comment gere
		return
	}
	log.Print("here: ")
	profile, err := server.store.GetEmployeeProfileByUserID(ctx, payload.UserId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	var permissions []Permission
	if err := json.Unmarshal(profile.Permissions, &permissions); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(GetEmployeeProfileResponse{
		UserID:      profile.UserID,
		EmployeeID:  profile.EmployeeID,
		FirstName:   profile.FirstName,
		LastName:    profile.LastName,
		Email:       profile.Email,
		RoleID:      profile.RoleID,
		Permissions: permissions,
	}, "Employee profile retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// CreateEmployeeProfileRequest represents the request for CreateEmployeeProfileApi
type CreateEmployeeProfileRequest struct {
	EmployeeNumber            *string `json:"employee_number" example:"123456"`
	EmploymentNumber          *string `json:"employment_number" example:"123456"`
	Location                  *int64  `json:"location" example:"1"`
	IsSubcontractor           *bool   `json:"is_subcontractor" example:"false"`
	FirstName                 string  `json:"first_name" binding:"required" example:"fara"`
	LastName                  string  `json:"last_name" binding:"required" example:"joe"`
	DateOfBirth               *string `json:"date_of_birth" example:"2000-01-01"`
	Gender                    *string `json:"gender" exmple:"man"`
	Email                     string  `json:"email" binding:"required,email" example:"emai@exe.com"`
	PrivateEmailAddress       *string `json:"private_email_address" binding:"email" example:"joe@ex.com"`
	AuthenticationPhoneNumber *string `json:"authentication_phone_number" example:"1234567890"`
	WorkPhoneNumber           *string `json:"work_phone_number" example:"1234567890"`
	PrivatePhoneNumber        *string `json:"private_phone_number" example:"1234567890"`
	HomeTelephoneNumber       *string `json:"home_telephone_number" example:"1234567890"`
	OutOfService              *bool   `json:"out_of_service" example:"false"`
	RoleID                    int32   `json:"role_id" binding:"required" example:"1"`
	Position                  *string `json:"position" example:"developer"`
	Department                *string `json:"department" example:"IT"`
}

// CreateEmployeeProfileResponse represents the response for CreateEmployeeProfileApi
type CreateEmployeeProfileResponse struct {
	ID                        int64              `json:"id"`
	UserID                    int64              `json:"user_id"`
	FirstName                 string             `json:"first_name"`
	LastName                  string             `json:"last_name"`
	Position                  *string            `json:"position"`
	Department                *string            `json:"department"`
	EmployeeNumber            *string            `json:"employee_number"`
	EmploymentNumber          *string            `json:"employment_number"`
	PrivateEmailAddress       *string            `json:"private_email_address"`
	Email                     string             `json:"email"`
	AuthenticationPhoneNumber *string            `json:"authentication_phone_number"`
	PrivatePhoneNumber        *string            `json:"private_phone_number"`
	WorkPhoneNumber           *string            `json:"work_phone_number"`
	DateOfBirth               pgtype.Date        `json:"date_of_birth"`
	HomeTelephoneNumber       *string            `json:"home_telephone_number"`
	CreatedAt                 pgtype.Timestamptz `json:"created_at"`
	IsSubcontractor           *bool              `json:"is_subcontractor"`
	Gender                    *string            `json:"gender"`
	LocationID                *int64             `json:"location_id"`
	HasBorrowed               bool               `json:"has_borrowed"`
	OutOfService              *bool              `json:"out_of_service"`
	IsArchived                bool               `json:"is_archived"`
}

// @Summary Create employee profile
// @Description Create a new employee profile with associated user account
// @Tags employees
// @Accept json
// @Produce json
// @Param request body CreateEmployeeProfileRequest true "Employee profile details"
// @Success 201 {object} Response[CreateEmployeeProfileResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees [post]
func (server *Server) CreateEmployeeProfileApi(ctx *gin.Context) {
	var req CreateEmployeeProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(util.RandomString(6))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var parsedDate time.Time
	if req.DateOfBirth != nil {
		parsedDate, err = time.Parse("2006-01-02", *req.DateOfBirth)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}

	employee, err := server.store.CreateEmployeeWithAccountTx(
		ctx,
		db.CreateEmployeeWithAccountTxParams{
			CreateUserParams: db.CreateUserParams{
				Password: hashedPassword,
				Email:    req.Email,
				IsActive: true,
				RoleID:   req.RoleID,
			},

			CreateEmployeeParams: db.CreateEmployeeProfileParams{
				FirstName:                 req.FirstName,
				LastName:                  req.LastName,
				EmployeeNumber:            req.EmployeeNumber,
				EmploymentNumber:          req.EmploymentNumber,
				LocationID:                req.Location,
				IsSubcontractor:           req.IsSubcontractor,
				DateOfBirth:               pgtype.Date{Time: parsedDate, Valid: true},
				Gender:                    req.Gender,
				Email:                     req.Email,
				PrivateEmailAddress:       req.PrivateEmailAddress,
				AuthenticationPhoneNumber: req.AuthenticationPhoneNumber,
				WorkPhoneNumber:           req.WorkPhoneNumber,
				PrivatePhoneNumber:        req.PrivatePhoneNumber,
				HomeTelephoneNumber:       req.HomeTelephoneNumber,
				OutOfService:              req.OutOfService,
				Position:                  req.Position,
				Department:                req.Department,
			},
		},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(CreateEmployeeProfileResponse{
		ID:                        employee.Employee.ID,
		EmployeeNumber:            employee.Employee.EmployeeNumber,
		EmploymentNumber:          employee.Employee.EmploymentNumber,
		FirstName:                 employee.Employee.FirstName,
		LastName:                  employee.Employee.LastName,
		IsSubcontractor:           employee.Employee.IsSubcontractor,
		DateOfBirth:               employee.Employee.DateOfBirth,
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
		CreatedAt:                 employee.Employee.CreatedAt,
		IsArchived:                employee.Employee.IsArchived,
		LocationID:                employee.Employee.LocationID,
	}, "Employee profile created successfully")

	ctx.JSON(http.StatusCreated, res)
}

// ListEmployeeRequest represents the request for ListEmployeeProfileApi
type ListEmployeeRequest struct {
	pagination.Request
	IncludeArchived     *bool   `form:"is_archived"`
	IncludeOutOfService *bool   `form:"out_of_service"`
	Department          *string `form:"department"`
	Position            *string `form:"position"`
	LocationID          *int64  `form:"location_id"`
	Search              *string `form:"search"`
}

// ListEmployeeResponse represents the response for ListEmployeeProfileApi
type ListEmployeeResponse struct {
	ID                        int64              `json:"id"`
	UserID                    int64              `json:"user_id"`
	FirstName                 string             `json:"first_name"`
	LastName                  string             `json:"last_name"`
	Position                  *string            `json:"position"`
	Department                *string            `json:"department"`
	EmployeeNumber            *string            `json:"employee_number"`
	EmploymentNumber          *string            `json:"employment_number"`
	PrivateEmailAddress       *string            `json:"private_email_address"`
	Email                     string             `json:"email"`
	AuthenticationPhoneNumber *string            `json:"authentication_phone_number"`
	PrivatePhoneNumber        *string            `json:"private_phone_number"`
	WorkPhoneNumber           *string            `json:"work_phone_number"`
	DateOfBirth               pgtype.Date        `json:"date_of_birth"`
	HomeTelephoneNumber       *string            `json:"home_telephone_number"`
	CreatedAt                 pgtype.Timestamptz `json:"created_at"`
	IsSubcontractor           *bool              `json:"is_subcontractor"`
	Gender                    *string            `json:"gender"`
	LocationID                *int64             `json:"location_id"`
	HasBorrowed               bool               `json:"has_borrowed"`
	OutOfService              *bool              `json:"out_of_service"`
	IsArchived                bool               `json:"is_archived"`
	ProfilePicture            *string            `json:"profile_picture"`
}

// @Summary List employee profiles
// @Description Get a paginated list of employee profiles with optional filters
// @Tags employees
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param is_archived query bool false "Include archived employees"
// @Param out_of_service query bool false "Include out of service employees"
// @Param department query string false "Filter by department"
// @Param position query string false "Filter by position"
// @Param location_id query integer false "Filter by location ID"
// @Param search query string false "Search term for employee name or number"
// @Success 200 {object} Response[pagination.Response[ListEmployeeResponse]]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees [get]
func (server *Server) ListEmployeeProfileApi(ctx *gin.Context) {
	var req ListEmployeeRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := req.GetParams()

	arg := db.ListEmployeeProfileParams{
		Limit:               params.Limit,
		Offset:              params.Offset,
		IncludeArchived:     req.IncludeArchived,
		IncludeOutOfService: req.IncludeOutOfService,
		Search:              req.Search,
	}

	employees, err := server.store.ListEmployeeProfile(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	countArg := db.CountEmployeeProfileParams{
		IncludeArchived:     arg.IncludeArchived,
		IncludeOutOfService: arg.IncludeOutOfService,
		Department:          arg.Department,
		Position:            arg.Position,
		LocationID:          arg.LocationID,
	}

	// Get total count
	totalCount, err := server.store.CountEmployeeProfile(ctx, countArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
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
			DateOfBirth:               employee.DateOfBirth,
			HomeTelephoneNumber:       employee.HomeTelephoneNumber,
			CreatedAt:                 employee.CreatedAt,
			IsSubcontractor:           employee.IsSubcontractor,
			Gender:                    employee.Gender,
			LocationID:                employee.LocationID,
			HasBorrowed:               employee.HasBorrowed,
			OutOfService:              employee.OutOfService,
			IsArchived:                employee.IsArchived,
			ProfilePicture:            employee.ProfilePicture,
		}
	}

	response := pagination.NewResponse(ctx, req.Request, responseEmployees, totalCount)
	res := SuccessResponse(response, "Employee profiles retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// GetEmployeeProfileByIDApiResponse represents the response for GetEmployeeProfileByIDApi
type GetEmployeeProfileByIDApiResponse struct {
	ID                        int64              `json:"id"`
	UserID                    int64              `json:"user_id"`
	FirstName                 string             `json:"first_name"`
	LastName                  string             `json:"last_name"`
	Position                  *string            `json:"position"`
	Department                *string            `json:"department"`
	EmployeeNumber            *string            `json:"employee_number"`
	EmploymentNumber          *string            `json:"employment_number"`
	PrivateEmailAddress       *string            `json:"private_email_address"`
	Email                     string             `json:"email"`
	AuthenticationPhoneNumber *string            `json:"authentication_phone_number"`
	PrivatePhoneNumber        *string            `json:"private_phone_number"`
	WorkPhoneNumber           *string            `json:"work_phone_number"`
	DateOfBirth               pgtype.Date        `json:"date_of_birth"`
	HomeTelephoneNumber       *string            `json:"home_telephone_number"`
	CreatedAt                 pgtype.Timestamptz `json:"created_at"`
	IsSubcontractor           *bool              `json:"is_subcontractor"`
	Gender                    *string            `json:"gender"`
	LocationID                *int64             `json:"location_id"`
	HasBorrowed               bool               `json:"has_borrowed"`
	OutOfService              *bool              `json:"out_of_service"`
	IsArchived                bool               `json:"is_archived"`
	ProfilePicture            *string            `json:"profile_picture"`
	RoleID                    int32              `json:"role_id"`
}

// @Summary Get employee profile by  ID
// @Description Get employee profile by ID
// @Tags employees
// @Produce json
// @Success 200 {object} Response[GetEmployeeProfileByIDApiResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id} [get]
func (server *Server) GetEmployeeProfileByIDApi(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	employee, err := server.store.GetEmployeeProfileByID(ctx, employeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(GetEmployeeProfileByIDApiResponse{
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
		DateOfBirth:               employee.DateOfBirth,
		HomeTelephoneNumber:       employee.HomeTelephoneNumber,
		CreatedAt:                 employee.CreatedAt,
		IsSubcontractor:           employee.IsSubcontractor,
		Gender:                    employee.Gender,
		LocationID:                employee.LocationID,
		HasBorrowed:               employee.HasBorrowed,
		OutOfService:              employee.OutOfService,
		IsArchived:                employee.IsArchived,
		ProfilePicture:            employee.ProfilePicture,
		RoleID:                    employee.RoleID,
	}, "Employee profile retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateEmployeeProfileRequest represents the request for UpdateEmployeeProfileApi

type UpdateEmployeeProfileRequest struct {
	FirstName                 *string `json:"first_name"`
	LastName                  *string `json:"last_name"`
	Position                  *string `json:"position"`
	Department                *string `json:"department"`
	EmployeeNumber            *string `json:"employee_number"`
	EmploymentNumber          *string `json:"employment_number"`
	PrivateEmailAddress       *string `json:"private_email_address"`
	Email                     *string `json:"email"`
	AuthenticationPhoneNumber *string `json:"authentication_phone_number"`
	PrivatePhoneNumber        *string `json:"private_phone_number"`
	WorkPhoneNumber           *string `json:"work_phone_number"`
	DateOfBirth               *string `json:"date_of_birth"`
	HomeTelephoneNumber       *string `json:"home_telephone_number"`
	IsSubcontractor           *bool   `json:"is_subcontractor"`
	Gender                    *string `json:"gender"`
	LocationID                *int64  `json:"location_id"`
	HasBorrowed               *bool   `json:"has_borrowed"`
	OutOfService              *bool   `json:"out_of_service"`
	IsArchived                *bool   `json:"is_archived"`
}

// UpdateEmployeeProfileResponse represents the response for UpdateEmployeeProfileApi

type UpdateEmployeeProfileResponse struct {
	ID                        int64              `json:"id"`
	UserID                    int64              `json:"user_id"`
	FirstName                 string             `json:"first_name"`
	LastName                  string             `json:"last_name"`
	Position                  *string            `json:"position"`
	Department                *string            `json:"department"`
	EmployeeNumber            *string            `json:"employee_number"`
	EmploymentNumber          *string            `json:"employment_number"`
	PrivateEmailAddress       *string            `json:"private_email_address"`
	Email                     string             `json:"email"`
	AuthenticationPhoneNumber *string            `json:"authentication_phone_number"`
	PrivatePhoneNumber        *string            `json:"private_phone_number"`
	WorkPhoneNumber           *string            `json:"work_phone_number"`
	DateOfBirth               pgtype.Date        `json:"date_of_birth"`
	HomeTelephoneNumber       *string            `json:"home_telephone_number"`
	CreatedAt                 pgtype.Timestamptz `json:"created_at"`
	IsSubcontractor           *bool              `json:"is_subcontractor"`
	Gender                    *string            `json:"gender"`
	LocationID                *int64             `json:"location_id"`
	HasBorrowed               bool               `json:"has_borrowed"`
	OutOfService              *bool              `json:"out_of_service"`
	IsArchived                bool               `json:"is_archived"`
}

// @Summary Update employee profile by ID
// @Description Update employee profile by ID
// @Tags employees
// @Produce json
// @Success 200 {object} Response[UpdateEmployeeProfileResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id} [put]
func (server *Server) UpdateEmployeeProfileApi(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var req UpdateEmployeeProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var parsedDate time.Time
	if req.DateOfBirth != nil {
		parsedDate, err = time.Parse("2006-01-02", *req.DateOfBirth)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	employee, err := server.store.UpdateEmployeeProfile(ctx, db.UpdateEmployeeProfileParams{
		ID:                        employeeID,
		FirstName:                 req.FirstName,
		LastName:                  req.LastName,
		Position:                  req.Position,
		Department:                req.Department,
		EmployeeNumber:            req.EmployeeNumber,
		EmploymentNumber:          req.EmploymentNumber,
		PrivateEmailAddress:       req.PrivateEmailAddress,
		Email:                     req.Email,
		AuthenticationPhoneNumber: req.AuthenticationPhoneNumber,
		PrivatePhoneNumber:        req.PrivatePhoneNumber,
		WorkPhoneNumber:           req.WorkPhoneNumber,
		DateOfBirth:               pgtype.Date{Time: parsedDate, Valid: true},
		HomeTelephoneNumber:       req.HomeTelephoneNumber,
		IsSubcontractor:           req.IsSubcontractor,
		Gender:                    req.Gender,
		LocationID:                req.LocationID,
		HasBorrowed:               req.HasBorrowed,
		OutOfService:              req.OutOfService,
		IsArchived:                req.IsArchived,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(UpdateEmployeeProfileResponse{
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
		DateOfBirth:               employee.DateOfBirth,
		HomeTelephoneNumber:       employee.HomeTelephoneNumber,
		CreatedAt:                 employee.CreatedAt,
		IsSubcontractor:           employee.IsSubcontractor,
		Gender:                    employee.Gender,
		LocationID:                employee.LocationID,
		HasBorrowed:               employee.HasBorrowed,
		OutOfService:              employee.OutOfService,
		IsArchived:                employee.IsArchived,
	}, "Employee profile updated successfully")
	ctx.JSON(http.StatusOK, res)

}

// AddEducationToEmployeeProfileRequest represents the request for AddEducationToEmployeeProfileApi
type AddEducationToEmployeeProfileRequest struct {
	InstitutionName string `json:"institution_name" binding:"required"`
	Degree          string `json:"degree" binding:"required"`
	FieldOfStudy    string `json:"field_of_study" binding:"required"`
	StartDate       string `json:"start_date" binding:"required" time_format:"2006-01-02"`
	EndDate         string `json:"end_date" binding:"required" time_format:"2006-01-02"`
}

// AddEducationToEmployeeProfileResponse represents the response for AddEducationToEmployeeProfileApi
type AddEducationToEmployeeProfileResponse struct {
	ID              int64  `json:"id"`
	EmployeeID      int64  `json:"employee_id"`
	InstitutionName string `json:"institution_name"`
	Degree          string `json:"degree"`
	FieldOfStudy    string `json:"field_of_study"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
}

// @Summary Add education to employee profile
// @Description Add education to employee profile
// @Tags employees
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Param request body AddEducationToEmployeeProfileRequest true "Education details"
// @Success 201 {object} Response[AddEducationToEmployeeProfileResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
func (server *Server) AddEducationToEmployeeProfileApi(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var req AddEducationToEmployeeProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	parsedStartDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	parsedEndDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.AddEducationToEmployeeProfileParams{
		EmployeeID:      employeeID,
		InstitutionName: req.InstitutionName,
		Degree:          req.Degree,
		FieldOfStudy:    req.FieldOfStudy,
		StartDate:       pgtype.Date{Time: parsedStartDate, Valid: true},
		EndDate:         pgtype.Date{Time: parsedEndDate, Valid: true},
	}
	education, err := server.store.AddEducationToEmployeeProfile(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(AddEducationToEmployeeProfileResponse{
		ID:              education.ID,
		EmployeeID:      education.EmployeeID,
		InstitutionName: education.InstitutionName,
		Degree:          education.Degree,
		FieldOfStudy:    education.FieldOfStudy,
		StartDate:       education.StartDate.Time.Format("2006-01-02"),
		EndDate:         education.EndDate.Time.Format("2006-01-02"),
	}, "Education added to employee profile successfully")
	ctx.JSON(http.StatusCreated, res)

}

// ListEmployeeEducationResponse represents the response for ListEmployeeEducationApi
type ListEmployeeEducationResponse struct {
	ID              int64  `json:"id"`
	EmployeeID      int64  `json:"employee_id"`
	InstitutionName string `json:"institution_name"`
	Degree          string `json:"degree"`
	FieldOfStudy    string `json:"field_of_study"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
}

// @Summary List education for employee profile
// @Description Get a list of education for employee profile
// @Tags employees
// @Produce json
// @Param id path int true "Employee ID"
// @Success 200 {object} Response[ListEmployeeEducationResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/education [get]
func (server *Server) ListEmployeeEducationApi(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	educations, err := server.store.ListEducations(ctx, employeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	responseEducations := make([]ListEmployeeEducationResponse, len(educations))
	for i, education := range educations {
		responseEducations[i] = ListEmployeeEducationResponse{
			ID:              education.ID,
			EmployeeID:      education.EmployeeID,
			InstitutionName: education.InstitutionName,
			Degree:          education.Degree,
			FieldOfStudy:    education.FieldOfStudy,
			StartDate:       education.StartDate.Time.Format("2006-01-02"),
			EndDate:         education.EndDate.Time.Format("2006-01-02"),
		}
	}
	res := SuccessResponse(responseEducations, "Employee education retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateEmployeeEducationRequest represents the request for UpdateEmployeeEducationApi
type UpdateEmployeeEducationRequest struct {
	InstitutionName *string `json:"institution_name"`
	Degree          *string `json:"degree"`
	FieldOfStudy    *string `json:"field_of_study"`
	StartDate       *string `json:"start_date" time_format:"2006-01-02"`
	EndDate         *string `json:"end_date" time_format:"2006-01-02"`
}

// UpdateEmployeeEducationResponse represents the response for UpdateEmployeeEducationApi
type UpdateEmployeeEducationResponse struct {
	ID              int64  `json:"id"`
	InstitutionName string `json:"institution_name"`
	Degree          string `json:"degree"`
	FieldOfStudy    string `json:"field_of_study"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
}

// @Summary Update education for employee profile
// @Description Update education for employee profile
// @Tags employees
// @Produce json
// @Param id path int true "Employee ID"
// @Param education_id path int true "Education ID"
// @Param request body UpdateEmployeeEducationRequest true "Education details"
// @Success 200 {object} Response[UpdateEmployeeEducationResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/education/{education_id} [put]
func (server *Server) UpdateEmployeeEducationApi(ctx *gin.Context) {
	id := ctx.Param(":education_id")
	educationID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var req UpdateEmployeeEducationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var parsedStartDate time.Time
	if req.StartDate != nil {
		parsedStartDate, err = time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	var parsedEndDate time.Time
	if req.EndDate != nil {
		parsedEndDate, err = time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	education, err := server.store.UpdateEmployeeEducation(ctx, db.UpdateEmployeeEducationParams{
		ID:              educationID,
		InstitutionName: req.InstitutionName,
		Degree:          req.Degree,
		FieldOfStudy:    req.FieldOfStudy,
		StartDate:       pgtype.Date{Time: parsedStartDate, Valid: true},
		EndDate:         pgtype.Date{Time: parsedEndDate, Valid: true},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(UpdateEmployeeEducationResponse{
		ID:              education.ID,
		InstitutionName: education.InstitutionName,
		Degree:          education.Degree,
		FieldOfStudy:    education.FieldOfStudy,
		StartDate:       education.StartDate.Time.Format("2006-01-02"),
		EndDate:         education.EndDate.Time.Format("2006-01-02"),
	}, "Education updated successfully")
	ctx.JSON(http.StatusOK, res)

}

// AddEmployeeExperienceRequest represents the request for AddEmployeeExperienceApi
type AddEmployeeExperienceRequest struct {
	JobTitle    string  `json:"job_title" binding:"required"`
	CompanyName string  `json:"company_name" binding:"required"`
	StartDate   string  `json:"start_date" binding:"required" time_format:"2006-01-02"`
	EndDate     string  `json:"end_date" binding:"required" time_format:"2006-01-02"`
	Description *string `json:"description"`
}

// AddEmployeeExperienceResponse represents the response for AddEmployeeExperienceApi
type AddEmployeeExperienceResponse struct {
	ID          int64   `json:"id"`
	EmployeeID  int64   `json:"employee_id"`
	JobTitle    string  `json:"job_title"`
	CompanyName string  `json:"company_name"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
	Description *string `json:"description"`
	CreatedAt   string  `json:"created_at"`
}

// @Summary Add experience to employee profile
// @Description Add experience to employee profile
// @Tags employees
// @Produce json
// @Param id path int true "Employee ID"
// @Param request body AddEmployeeExperienceRequest true "Experience details"
// @Success 201 {object} Response[AddEmployeeExperienceResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/experience [post]
func (server *Server) AddEmployeeExperienceApi(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var req AddEmployeeExperienceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	parsedStartDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	parsedEndDate, err := time.Parse("2006-01-02", req.EndDate)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.AddEmployeeExperienceParams{
		EmployeeID:  employeeID,
		JobTitle:    req.JobTitle,
		CompanyName: req.CompanyName,
		StartDate:   pgtype.Date{Time: parsedStartDate, Valid: true},
		EndDate:     pgtype.Date{Time: parsedEndDate, Valid: true},
		Description: req.Description,
	}
	experience, err := server.store.AddEmployeeExperience(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(AddEmployeeExperienceResponse{
		ID:          experience.ID,
		EmployeeID:  experience.EmployeeID,
		JobTitle:    experience.JobTitle,
		CompanyName: experience.CompanyName,
		StartDate:   experience.StartDate.Time.Format(time.RFC3339),
		EndDate:     experience.EndDate.Time.Format(time.RFC3339),
		Description: experience.Description,
		CreatedAt:   experience.CreatedAt.Time.Format(time.RFC3339),
	}, "Experience added to employee profile successfully")
	ctx.JSON(http.StatusCreated, res)

}

// ListEmployeeExperienceResponse represents the response for ListEmployeeExperienceApi
type ListEmployeeExperienceResponse struct {
	ID          int64   `json:"id"`
	EmployeeID  int64   `json:"employee_id"`
	JobTitle    string  `json:"job_title"`
	CompanyName string  `json:"company_name"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
	Description *string `json:"description"`
	CreatedAt   string  `json:"created_at"`
}

func (server *Server) ListEmployeeExperienceApi(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	experiences, err := server.store.ListEmployeeExperience(ctx, employeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	responseExperiences := make([]ListEmployeeExperienceResponse, len(experiences))
	for i, experience := range experiences {
		responseExperiences[i] = ListEmployeeExperienceResponse{
			ID:          experience.ID,
			EmployeeID:  experience.EmployeeID,
			JobTitle:    experience.JobTitle,
			CompanyName: experience.CompanyName,
			StartDate:   experience.StartDate.Time.Format(time.RFC3339),
			EndDate:     experience.EndDate.Time.Format(time.RFC3339),
			Description: experience.Description,
			CreatedAt:   experience.CreatedAt.Time.Format(time.RFC3339),
		}
	}
	res := SuccessResponse(responseExperiences, "Employee experience retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateEmployeeExperienceRequest represents the request for UpdateEmployeeExperienceApi
type UpdateEmployeeExperienceRequest struct {
	JobTitle    *string `json:"job_title"`
	CompanyName *string `json:"company_name"`
	StartDate   *string `json:"start_date" time_format:"2006-01-02"`
	EndDate     *string `json:"end_date" time_format:"2006-01-02"`
	Description *string `json:"description"`
}

// UpdateEmployeeExperienceResponse represents the response for UpdateEmployeeExperienceApi
type UpdateEmployeeExperienceResponse struct {
	ID          int64   `json:"id"`
	EmployeeID  int64   `json:"employee_id"`
	JobTitle    string  `json:"job_title"`
	CompanyName string  `json:"company_name"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
	Description *string `json:"description"`
	CreatedAt   string  `json:"created_at"`
}

// @Summary Update experience for employee profile
// @Description Update experience for employee profile
// @Tags employees
// @Produce json
// @Param id path int true "Employee ID"
// @Param experience_id path int true "Experience ID"
// @Param request body UpdateEmployeeExperienceRequest true "Experience details"
// @Success 200 {object} Response[UpdateEmployeeExperienceResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/experience/{experience_id} [put]
func (server *Server) UpdateEmployeeExperienceApi(ctx *gin.Context) {
	id := ctx.Param(":experience_id")
	experienceID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var req UpdateEmployeeExperienceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var parsedStartDate time.Time
	if req.StartDate != nil {
		parsedStartDate, err = time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	var parsedEndDate time.Time
	if req.EndDate != nil {
		parsedEndDate, err = time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	experience, err := server.store.UpdateEmployeeExperience(ctx, db.UpdateEmployeeExperienceParams{
		ID:          experienceID,
		JobTitle:    req.JobTitle,
		CompanyName: req.CompanyName,
		StartDate:   pgtype.Date{Time: parsedStartDate, Valid: true},
		EndDate:     pgtype.Date{Time: parsedEndDate, Valid: true},
		Description: req.Description,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(UpdateEmployeeExperienceResponse{
		ID:          experience.ID,
		EmployeeID:  experience.EmployeeID,
		JobTitle:    experience.JobTitle,
		CompanyName: experience.CompanyName,
		StartDate:   experience.StartDate.Time.Format("2006-01-02"),
		EndDate:     experience.EndDate.Time.Format("2006-01-02"),
		Description: experience.Description,
		CreatedAt:   experience.CreatedAt.Time.Format("2006-01-02"),
	}, "Experience updated successfully")
	ctx.JSON(http.StatusOK, res)
}
