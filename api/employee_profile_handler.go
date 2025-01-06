package api

import (
	"log"
	"net/http"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// GetEmployeeProfileResponse represents the response for GetEmployeeProfileApi
type GetEmployeeProfileResponse struct {
	UserID     int64  `json:"user_id"`
	EmployeeID int64  `json:"employee_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	RoleID     int32  `json:"role_id"`
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
	log.Print("Profile: ", profile)
	res := SuccessResponse(GetEmployeeProfileResponse{
		UserID:     profile.UserID,
		EmployeeID: profile.EmployeeID,
		FirstName:  profile.FirstName,
		LastName:   profile.LastName,
		Email:      profile.Email,
		RoleID:     profile.RoleID,
	}, "Employee profile retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// CreateEmployeeProfileRequest represents the request for CreateEmployeeProfileApi
type CreateEmployeeProfileRequest struct {
	EmployeeNumber            *string `json:"employee_number"`
	EmploymentNumber          *string `json:"employment_number"`
	Location                  *int64  `json:"location" example:"1"`
	IsSubcontractor           *bool   `json:"is_subcontractor"`
	FirstName                 string  `json:"first_name" binding:"required"`
	LastName                  string  `json:"last_name" binding:"required"`
	DateOfBirth               *string `json:"date_of_birth"`
	Gender                    *string `json:"gender"`
	Email                     string  `json:"email_address" binding:"required,email"`
	PrivateEmailAddress       *string `json:"private_email_address" binding:"email"`
	AuthenticationPhoneNumber *string `json:"authentication_phone_number"`
	WorkPhoneNumber           *string `json:"work_phone_number"`
	PrivatePhoneNumber        *string `json:"private_phone_number"`
	HomeTelephoneNumber       *string `json:"home_telephone_number"`
	OutOfService              *bool   `json:"out_of_service"`
	RoleID                    int32   `json:"role_id" binding:"required" example:"1"`
}

// CreateEmployeeProfileResponse represents the response for CreateEmployeeProfileApi
type CreateEmployeeProfileResponse struct {
	ID                        int64     `json:"id"`
	UserID                    int64     `json:"user_id"`
	FirstName                 string    `json:"first_name"`
	LastName                  string    `json:"last_name"`
	Position                  *string   `json:"position"`
	Department                *string   `json:"department"`
	EmployeeNumber            *string   `json:"employee_number"`
	EmploymentNumber          *string   `json:"employment_number"`
	PrivateEmailAddress       *string   `json:"private_email_address"`
	Email                     string    `json:"email_address"`
	AuthenticationPhoneNumber *string   `json:"authentication_phone_number"`
	PrivatePhoneNumber        *string   `json:"private_phone_number"`
	WorkPhoneNumber           *string   `json:"work_phone_number"`
	DateOfBirth               time.Time `json:"date_of_birth"`
	HomeTelephoneNumber       *string   `json:"home_telephone_number"`
	CreatedAt                 time.Time `json:"created"`
	IsSubcontractor           *bool     `json:"is_subcontractor"`
	Gender                    *string   `json:"gender"`
	LocationID                *int64    `json:"location_id"`
	HasBorrowed               bool      `json:"has_borrowed"`
	OutOfService              *bool     `json:"out_of_service"`
	IsArchived                bool      `json:"is_archived"`
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

	parsedDate, err := time.Parse("2006-01-02", *req.DateOfBirth)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
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
	}, "Employee profile created successfully")

	ctx.JSON(http.StatusCreated, res)
}

type ListEmployeeRequest struct {
	pagination.Request
	IncludeArchived     *bool   `form:"is_archived"`
	IncludeOutOfService *bool   `form:"out_of_service"`
	Department          *string `form:"department"`
	Position            *string `form:"position"`
	LocationID          *int64  `form:"location_id"`
	Search              *string `form:"search"`
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
// @Success 200 {object} pagination.Response[db.EmployeeProfile]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Security BearerAuth
// @Router /employees/employees_list [get]
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

	response := pagination.NewResponse(ctx, req.Request, employees, totalCount)
	ctx.JSON(http.StatusOK, response)
}
