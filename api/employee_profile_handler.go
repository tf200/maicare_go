package api

import (
	"net/http"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateEmployeeProfileRequest struct {
	EmployeeNumber            *string `json:"employee_number"`
	EmploymentNumber          *string `json:"employment_number" binding:"required"`
	Location                  *int64  `json:"location" binding:"required"`
	IsSubcontractor           *bool   `json:"is_subcontractor"`
	FirstName                 string  `json:"first_name" binding:"required"`
	LastName                  string  `json:"last_name" binding:"required"`
	DateOfBirth               *string `json:"date_of_birth" binding:"required"`
	Gender                    *string `json:"gender" binding:"required"`
	EmailAddress              *string `json:"email_address" binding:"required,email"`
	PrivateEmailAddress       *string `json:"private_email_address" binding:"required,email"`
	AuthenticationPhoneNumber *string `json:"authentication_phone_number" binding:"required"`
	WorkPhoneNumber           *string `json:"work_phone_number" binding:"required"`
	PrivatePhoneNumber        *string `json:"private_phone_number" binding:"required"`
	HomeTelephoneNumber       *string `json:"home_telephone_number" binding:"required"`
	OutOfService              *bool   `json:"out_of_service"`
}

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
	EmailAddress              *string   `json:"email_address"`
	AuthenticationPhoneNumber *string   `json:"authentication_phone_number"`
	PrivatePhoneNumber        *string   `json:"private_phone_number"`
	WorkPhoneNumber           *string   `json:"work_phone_number"`
	DateOfBirth               time.Time `json:"date_of_birth"`
	HomeTelephoneNumber       *string   `json:"home_telephone_number"`
	Created                   time.Time `json:"created"`
	IsSubcontractor           *bool     `json:"is_subcontractor"`
	Gender                    *string   `json:"gender"`
	LocationID                *int64    `json:"location_id"`
	HasBorrowed               bool      `json:"has_borrowed"`
	OutOfService              *bool     `json:"out_of_service"`
	IsArchived                bool      `json:"is_archived"`
}

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
				Username:    util.StringPtr(util.GenerateUsername(req.FirstName, req.LastName)),
				Password:    hashedPassword,
				FirstName:   req.FirstName,
				LastName:    req.LastName,
				Email:       *req.EmailAddress,
				IsActive:    true,
				PhoneNumber: util.IntPtr(562), // TO DO
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
				EmailAddress:              req.EmailAddress,
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
	res := CreateEmployeeProfileResponse{
		ID:                        employee.Employee.ID,
		EmployeeNumber:            employee.Employee.EmployeeNumber,
		EmploymentNumber:          employee.Employee.EmploymentNumber,
		FirstName:                 employee.Employee.FirstName,
		LastName:                  employee.Employee.LastName,
		IsSubcontractor:           employee.Employee.IsSubcontractor,
		DateOfBirth:               employee.Employee.DateOfBirth.Time,
		Gender:                    employee.Employee.Gender,
		EmailAddress:              employee.Employee.EmailAddress,
		PrivateEmailAddress:       employee.Employee.PrivateEmailAddress,
		AuthenticationPhoneNumber: employee.Employee.AuthenticationPhoneNumber,
		WorkPhoneNumber:           employee.Employee.WorkPhoneNumber,
		PrivatePhoneNumber:        employee.Employee.PrivatePhoneNumber,
		HomeTelephoneNumber:       employee.Employee.HomeTelephoneNumber,
		OutOfService:              employee.Employee.OutOfService,
		HasBorrowed:               employee.Employee.HasBorrowed,
		UserID:                    employee.User.ID,
		Created:                   employee.Employee.Created.Time,
		IsArchived:                employee.Employee.IsArchived,
		LocationID:                employee.Employee.LocationID,
	}

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
