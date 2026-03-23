package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	_ "maicare_go/pagination" // import for pagination.Response used in swagger
	"maicare_go/service/employees"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Get employee profile by user ID
// @Description Get employee profile by user ID
// @Tags employees
// @Produce json
// @Success 200 {object} Response[employees.GetEmployeeProfileResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/profile [get]
func (server *Server) GetEmployeeProfileApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	profile, err := server.businessService.EmployeeService.GetEmployeeProfile(payload.UserId, ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		case errors.Is(err, context.DeadlineExceeded):
			ctx.JSON(http.StatusRequestTimeout, errorResponse(err))
		default:
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	res := SuccessResponse(profile, "Employee profile retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Get detailed employee profile by user ID
// @Description Get detailed employee profile by user ID
// @Tags employees
// @Produce json
// @Success 200 {object} Response[employees.GetEmployeeProfileDetailsResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/profile/details [get]
func (server *Server) GetEmployeeProfileDetailsApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	profile, err := server.businessService.EmployeeService.GetEmployeeProfileDetails(payload.UserId, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(profile, "Employee profile details retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Get my schedules timeline
// @Description Get logged-in employee schedules and calendar events (excluding reminders) in a date range
// @Tags employees
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} Response[[]employees.GetMyScheduleTimelineDayResponse]
// @Failure 400,401,500 {object} Response[any]
// @Router /employees/profile/schedules [get]
func (server *Server) GetMyScheduleTimelineApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var req employees.GetMyScheduleTimelineRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	response, err := server.businessService.EmployeeService.GetMyScheduleTimeline(ctx.Request.Context(), payload.EmployeeID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(response, "Schedules timeline retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Create employee profile
// @Description Create a new employee profile with associated user account
// @Tags employees
// @Accept json
// @Produce json
// @Param request body employees.CreateEmployeeProfileRequest true "Employee profile details"
// @Success 201 {object} Response[employees.CreateEmployeeProfileResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees [post]
func (server *Server) CreateEmployeeProfileApi(ctx *gin.Context) {
	var req employees.CreateEmployeeProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	employee, err := server.businessService.EmployeeService.CreateEmployee(req, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create employee profile")))
		return
	}

	res := SuccessResponse(employee, "Employee profile created successfully")

	ctx.JSON(http.StatusCreated, res)
}

// @Summary List employee profiles
// @Description Get a paginated list of employee profiles with optional filters
// @Tags employees
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param is_archived query bool false "Include archived employees"
// @Param out_of_service query bool false "Include out of service employees"
// @Param location_id query integer false "Filter by location ID"
// @Param contract_type query string false "Filter by contract type (loondienst, ZZP, none)"
// @Param search query string false "Search term for employee name or number"
// @Success 200 {object} Response[pagination.Response[employees.ListEmployeeResponse]]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees [get]
func (server *Server) ListEmployeeProfileApi(ctx *gin.Context) {
	var req employees.ListEmployeeRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid query parameters")))
		return
	}
	response, err := server.businessService.EmployeeService.ListEmployees(req, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to list employee profiles")))
		return
	}

	res := SuccessResponse(response, "Employee profiles retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Get employee counts
// @Description Get total counts of employees, subcontractors, archived, and out of service employees
// @Tags employees
// @Produce json
// @Success 200 {object} Response[employees.GetEmployeeCountsResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/counts [get]
func (server *Server) GetEmployeeCountsApi(ctx *gin.Context) {
	counts, err := server.businessService.EmployeeService.GetEmployeeCounts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(counts, "Employee counts retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Get employee profile by  ID
// @Description Get employee profile by ID
// @Tags employees
// @Produce json
// @Success 200 {object} Response[employees.GetEmployeeProfileByIDResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id} [get]
func (server *Server) GetEmployeeProfileByIDApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	currentUserID := payload.UserId

	id := ctx.Param("id")
	employeeID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	profile, err := server.businessService.EmployeeService.GetEmployeeProfileByID(employeeID, currentUserID, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(profile, "Employee profile retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Update employee profile by ID
// @Description Update employee profile by ID
// @Tags employees
// @Produce json
// @Param id path uuid true "Employee ID"
// @Success 200 {object} Response[employees.UpdateEmployeeProfileResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id} [put]
func (server *Server) UpdateEmployeeProfileApi(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req employees.UpdateEmployeeProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	profile, err := server.businessService.EmployeeService.UpdateEmployeeProfile(req, employeeID, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(profile, "Employee profile updated successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Set employee profile picture by ID
// @Description Set employee profile picture by ID
// @Tags employees
// @Accept json
// @Produce json
// @Param id path uuid true "Employee ID"
// @Param request body employees.SetEmployeeProfilePictureRequest true "Profile picture details"
// @Success 200 {object} Response[employees.SetEmployeeProfilePictureResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/profile_picture [put]
func (server *Server) SetEmployeeProfilePictureApi(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req employees.SetEmployeeProfilePictureRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.EmployeeService.SetEmployeeProfilePicture(req, employeeID, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Employee profile picture updated successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Update employee's subcontractor status
// @Description Update an employee's subcontractor status and adjust contract details accordingly
// @Tags employees
// @Accept json
// @Produce json
// @Param id path uuid true "Employee ID"
// @Param request body employees.UpdateEmployeeIsSubcontractorRequest true "Subcontractor status details"
// @Success 200 {object} Response[employees.UpdateEmployeeIsSubcontractorResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/is_subcontractor [put]
func (server *Server) UpdateEmployeeIsSubcontractorApi(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid employee ID: %w", err)))
		return
	}
	var req employees.UpdateEmployeeIsSubcontractorRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body: %w", err)))
		return
	}

	result, err := server.businessService.EmployeeService.UpdateEmployeeIsSubcontractor(
		req, employeeID, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to update subcontractor status: %w", err)))
		return
	}

	res := SuccessResponse(result, "Employee subcontractor status updated successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Add contract details to employee profile
// @Description Add contract details to employee profile
// @Tags employees
// @Accept json
// @Produce json
// @Param id path uuid true "Employee ID"
// @Param request body employees.AddEmployeeContractDetailsRequest true "Contract details"
// @Success 201 {object} Response[employees.AddEmployeeContractDetailsResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/contract_details [put]
func (server *Server) AddEmployeeContractDetailsApi(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req employees.AddEmployeeContractDetailsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.EmployeeService.AddEmployeeContractDetails(req, employeeID, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Contract details added to employee profile successfully")
	ctx.JSON(http.StatusCreated, res)
}

// @Summary Get employee contract details by ID
// @Description Get employee contract details by ID
// @Tags employees
// @Produce json
// @Param id path uuid true "Employee ID"
// @Success 200 {object} Response[employees.GetEmployeeContractDetailsResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/contract_details [get]
func (server *Server) GetEmployeeContractDetailsApi(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	contractDetails, err := server.businessService.EmployeeService.GetEmployeeContractDetails(employeeID, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(contractDetails, "Employee contract details retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Add education to employee profile
// @Description Add education to employee profile
// @Tags employees
// @Accept json
// @Produce json
// @Param id path uuid true "Employee ID"
// @Param request body employees.AddEducationToEmployeeProfileRequest true "Education details"
// @Success 201 {object} Response[employees.AddEducationToEmployeeProfileResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/education [post]
func (server *Server) AddEducationToEmployeeProfileApi(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req employees.AddEducationToEmployeeProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.EmployeeService.AddEducationToEmployeeProfile(req, employeeID, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Education added to employee profile successfully")
	ctx.JSON(http.StatusCreated, res)
}

// @Summary List education for employee profile
// @Description Get a list of education for employee profile
// @Tags employees
// @Produce json
// @Param id path uuid true "Employee ID"
// @Success 200 {object} Response[[]employees.ListEmployeeEducationResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/education [get]
func (server *Server) ListEmployeeEducationApi(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	educations, err := server.businessService.EmployeeService.ListEmployeeEducation(employeeID, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(educations, "Employee education retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Update education for employee profile
// @Description Update education for employee profile
// @Tags employees
// @Produce json
// @Param id path uuid true "Employee ID"
// @Param education_id path uuid true "Education ID"
// @Param request body employees.UpdateEmployeeEducationRequest true "Education details"
// @Success 200 {object} Response[employees.UpdateEmployeeEducationResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/education/{education_id} [put]
func (server *Server) UpdateEmployeeEducationApi(ctx *gin.Context) {
	id := ctx.Param("education_id")
	educationID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req employees.UpdateEmployeeEducationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	education, err := server.businessService.EmployeeService.UpdateEmployeeEducation(req, educationID, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(education, "Education updated successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Delete education for employee profile
// @Description Delete education for employee profile
// @Tags employees
// @Produce json
// @Param id path uuid true "Employee ID"
// @Param education_id path uuid true "Education ID"
// @Success 200 {object} Response[employees.DeleteEmployeeEducationResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/education/{education_id} [delete]
func (server *Server) DeleteEmployeeEducationApi(ctx *gin.Context) {
	id := ctx.Param("education_id")
	educationID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	education, err := server.businessService.EmployeeService.DeleteEmployeeEducation(educationID, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(education, "Education deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Add experience to employee profile
// @Description Add experience to employee profile
// @Tags employees
// @Produce json
// @Param id path uuid true "Employee ID"
// @Param request body employees.AddEmployeeExperienceRequest true "Experience details"
// @Success 201 {object} Response[employees.AddEmployeeExperienceResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/experience [post]
func (server *Server) AddEmployeeExperienceApi(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var req employees.AddEmployeeExperienceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	result, err := server.businessService.EmployeeService.AddEmployeeExperience(req, employeeID, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Experience added to employee profile successfully")
	ctx.JSON(http.StatusCreated, res)
}

// @Summary List experience for employee profile
// @Description Get a list of experience for employee profile
// @Tags employees
// @Produce json
// @Param id path uuid true "Employee ID"
// @Success 200 {object} Response[[]employees.ListEmployeeExperienceResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/experience [get]
func (server *Server) ListEmployeeExperienceApi(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	responseExperiences, err := server.businessService.EmployeeService.ListEmployeeExperience(employeeID, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(responseExperiences, "Employee experience retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Update experience for employee profile
// @Description Update experience for employee profile
// @Tags employees
// @Produce json
// @Param id path uuid true "Employee ID"
// @Param experience_id path uuid true "Experience ID"
// @Param request body employees.UpdateEmployeeExperienceRequest true "Experience details"
// @Success 200 {object} Response[employees.UpdateEmployeeExperienceResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/experience/{experience_id} [put]
func (server *Server) UpdateEmployeeExperienceApi(ctx *gin.Context) {
	id := ctx.Param("experience_id")
	experienceID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var req employees.UpdateEmployeeExperienceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	experience, err := server.businessService.EmployeeService.UpdateEmployeeExperience(req, experienceID, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(experience, "Experience updated successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Delete experience for employee profile
// @Description Delete experience for employee profile
// @Tags employees
// @Produce json
// @Param id path uuid true "Employee ID"
// @Param experience_id path uuid true "Experience ID"
// @Success 200 {object} Response[employees.DeleteEmployeeExperienceResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/experience/{experience_id} [delete]
func (server *Server) DeleteEmployeeExperienceApi(ctx *gin.Context) {
	id := ctx.Param("experience_id")
	experienceID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	result, err := server.businessService.EmployeeService.DeleteEmployeeExperience(experienceID, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(result, "Experience deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Add certification to employee profile
// @Description Add certification to employee profile
// @Tags employees
// @Produce json
// @Param id path uuid true "Employee ID"
// @Param request body employees.AddEmployeeCertificationRequest true "Certification details"
// @Success 201 {object} Response[employees.AddEmployeeCertificationResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/certification [post]
func (server *Server) AddEmployeeCertificationApi(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var req employees.AddEmployeeCertificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	result, err := server.businessService.EmployeeService.AddEmployeeCertification(req, employeeID, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(result, "Certification added to employee profile successfully")
	ctx.JSON(http.StatusCreated, res)
}

// @Summary List certifications for employee profile
// @Description Get a list of certifications for employee profile
// @Tags employees
// @Produce json
// @Param id path uuid true "Employee ID"
// @Success 200 {object} Response[[]employees.ListEmployeeCertificationResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/certification [get]
func (server *Server) ListEmployeeCertificationApi(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.EmployeeService.ListEmployeeCertification(employeeID, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Employee certifications retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Update certification for employee profile
// @Description Update certification for employee profile
// @Tags employees
// @Produce json
// @Param id path uuid true "Employee ID"
// @Param certification_id path uuid true "Certification ID"
// @Param request body employees.UpdateEmployeeCertificationRequest true "Certification details"
// @Success 200 {object} Response[employees.UpdateEmployeeCertificationResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/certification/{certification_id} [put]
func (server *Server) UpdateEmployeeCertificationApi(ctx *gin.Context) {
	id := ctx.Param("certification_id")
	certificationID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var req employees.UpdateEmployeeCertificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	certification, err := server.businessService.EmployeeService.UpdateEmployeeCertification(req, certificationID, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(certification, "Certification updated successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Delete certification for employee profile
// @Description Delete certification for employee profile
// @Tags employees
// @Produce json
// @Param id path uuid true "Employee ID"
// @Param certification_id path uuid true "Certification ID"
// @Success 200 {object} Response[employees.DeleteEmployeeCertificationResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/certification/{certification_id} [delete]
func (server *Server) DeleteEmployeeCertificationApi(ctx *gin.Context) {
	id := ctx.Param("certification_id")
	certificationID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	result, err := server.businessService.EmployeeService.DeleteEmployeeCertification(certificationID, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(result, "Certification deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Search employees by name or email
// @Description Search employees by name or email
// @Tags employees
// @Produce json
// @Param search query string true "Search query"
// @Success 200 {object} Response[[]employees.SearchEmployeesByNameOrEmailResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/emails [get]
func (server *Server) SearchEmployeesByNameOrEmailApi(ctx *gin.Context) {
	var req employees.SearchEmployeesByNameOrEmailRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	employees, err := server.businessService.EmployeeService.SearchEmployeesByNameOrEmail(req, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(employees, "Employees retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}
