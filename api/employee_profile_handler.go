package api

import (
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"maicare_go/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateEmployeeProfileRequest struct {
	EmployeeNumber            string `json:"employee_number" binding:"required"`
	EmploymentNumber          string `json:"employment_number" binding:"required"`
	Location                  int64  `json:"location" binding:"required"`
	IsSubcontractor           bool   `json:"is_subcontractor"`
	FirstName                 string `json:"first_name" binding:"required"`
	LastName                  string `json:"last_name" binding:"required"`
	DateOfBirth               string `json:"date_of_birth" binding:"required"` // Could be time.Time if you parse it
	Gender                    string `json:"gender" binding:"required"`
	EmailAddress              string `json:"email_address" binding:"required,email"`
	PrivateEmailAddress       string `json:"private_email_address" binding:"required,email"`
	AuthenticationPhoneNumber string `json:"authentication_phone_number" binding:"required"`
	WorkPhoneNumber           string `json:"work_phone_number" binding:"required"`
	PrivatePhoneNumber        string `json:"private_phone_number" binding:"required"`
	HomeTelephoneNumber       string `json:"home_telephone_number" binding:"required"`
	OutOfService              bool   `json:"out_of_service"`
}

type CreateEmployeeProfileResponse struct {
	ID                        int    `json:"id"`
	EmployeeNumber            string `json:"employee_number"`
	EmploymentNumber          string `json:"employment_number"`
	Location                  int    `json:"location"`
	IsSubcontractor           bool   `json:"is_subcontractor"`
	FirstName                 string `json:"first_name"`
	LastName                  string `json:"last_name"`
	DateOfBirth               string `json:"date_of_birth"`
	Gender                    string `json:"gender"`
	EmailAddress              string `json:"email_address"`
	PrivateEmailAddress       string `json:"private_email_address"`
	AuthenticationPhoneNumber string `json:"authentication_phone_number"`
	WorkPhoneNumber           string `json:"work_phone_number"`
	PrivatePhoneNumber        string `json:"private_phone_number"`
	HomeTelephoneNumber       string `json:"home_telephone_number"`
	OutOfService              bool   `json:"out_of_service"`
	HasBorrowed               bool   `json:"has_borrowed"`
	IsArchived                bool   `json:"is_archived"`
	User                      int    `json:"user"`
	Created                   string `json:"created"`
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

	parsedDate, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	employee, err := server.store.CreateEmployeeWithAccountTx(
		ctx,
		db.CreateEmployeeWithAccountTxParams{
			CreateUserParams: db.CreateUserParams{
				Username:    util.GenerateUsername(req.FirstName, req.LastName),
				Password:    hashedPassword,
				FirstName:   req.FirstName,
				LastName:    req.LastName,
				Email:       req.EmailAddress,
				IsActive:    true,
				PhoneNumber: pgtype.Int8{Int64: 8}, // TO DO
			},

			CreateEmployeeParams: db.CreateEmployeeProfileParams{
				FirstName: req.FirstName,
				LastName:  req.LastName,
				EmployeeNumber: pgtype.Text{
					String: req.EmployeeNumber,
					Valid:  true,
				},
				EmploymentNumber: pgtype.Text{
					String: req.EmploymentNumber,
					Valid:  true,
				},
				LocationID: pgtype.Int8{Int64: req.Location, Valid: true}, // To do
				IsSubcontractor: pgtype.Bool{
					Bool:  req.IsSubcontractor,
					Valid: true,
				},
				DateOfBirth: pgtype.Date{
					Time:  parsedDate,
					Valid: true,
				},
				Gender: pgtype.Text{
					String: req.Gender,
					Valid:  true,
				},
				EmailAddress: pgtype.Text{
					String: req.EmailAddress,
					Valid:  true,
				}, // TO DO PARSING
				PrivateEmailAddress: pgtype.Text{
					String: req.PrivateEmailAddress,
					Valid:  true,
				},
				AuthenticationPhoneNumber: pgtype.Text{
					String: req.AuthenticationPhoneNumber,
					Valid:  true,
				},
				WorkPhoneNumber: pgtype.Text{
					String: req.WorkPhoneNumber,
					Valid:  true,
				},
				PrivatePhoneNumber: pgtype.Text{
					String: req.PrivatePhoneNumber,
					Valid:  true,
				},
				HomeTelephoneNumber: pgtype.Text{
					String: req.HomeTelephoneNumber,
					Valid:  true,
				},
				OutOfService: pgtype.Bool{
					Bool:  req.OutOfService,
					Valid: true,
				},
			},
		},
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := CreateEmployeeProfileResponse{
		ID:                        int(employee.Employee.ID),
		EmployeeNumber:            employee.Employee.EmployeeNumber.String,
		EmploymentNumber:          employee.Employee.EmploymentNumber.String,
		FirstName:                 employee.Employee.FirstName,
		LastName:                  employee.Employee.LastName,
		IsSubcontractor:           employee.Employee.IsSubcontractor.Bool,
		DateOfBirth:               employee.Employee.DateOfBirth.Time.String(),
		Gender:                    employee.Employee.Gender.String,
		EmailAddress:              employee.Employee.EmailAddress.String,
		PrivateEmailAddress:       employee.Employee.PrivateEmailAddress.String,
		AuthenticationPhoneNumber: employee.Employee.AuthenticationPhoneNumber.String,
		WorkPhoneNumber:           employee.Employee.WorkPhoneNumber.String,
		PrivatePhoneNumber:        employee.Employee.PrivatePhoneNumber.String,
		HomeTelephoneNumber:       employee.Employee.HomeTelephoneNumber.String,
		OutOfService:              employee.Employee.OutOfService.Bool,
		HasBorrowed:               employee.Employee.HasBorrowed,
		User:                      int(employee.User.ID),
		Created:                   employee.Employee.Created.Time.String(),
		IsArchived:                employee.Employee.IsArchived.Bool,
		Location:                  int(employee.Employee.LocationID.Int64),
	}

	ctx.JSON(http.StatusCreated, res)

}

type listEmployeeRequest struct {
	pagination.Request
	IncludeArchived     *bool   `form:"is_archived"`
	IncludeOutOfService *bool   `form:"out_of_service"`
	Department          *string `form:"department"`
	Position            *string `form:"position"`
	LocationID          *int64  `form:"location_id"`
}

func (server *Server) ListEmployeeProfileApi(ctx *gin.Context) {
	var req listEmployeeRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get pagination params
	params := req.GetParams()

	// Prepare database query params
	arg := db.ListEmployeeProfileParams{
		Limit:  params.Limit,
		Offset: params.Offset,
		IncludeArchived: pgtype.Bool{
			Valid: req.IncludeArchived != nil,
			Bool:  req.IncludeArchived != nil && *req.IncludeArchived,
		},
		IncludeOutOfService: pgtype.Bool{
			Valid: req.IncludeOutOfService != nil,
			Bool:  req.IncludeOutOfService != nil && *req.IncludeOutOfService,
		},
	}

	// Add optional filters...

	// Get employees
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

	// Create paginated response
	response := pagination.NewResponse(ctx, req.Request, employees, totalCount)
	ctx.JSON(http.StatusOK, response)
}
