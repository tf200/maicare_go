package api

import (
	"fmt"
	_ "maicare_go/pagination" // import for pagination.Response used in swagger
	"maicare_go/service/employees"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Get employee profile by user ID
// @Description Get employee profile by user ID
// @Tags employees
// @Produce json
// @Success 200 {object} Response[GetEmployeeProfileResponse]
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
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(profile, "Employee profile retrieved successfully")
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
// @Param department query string false "Filter by department"
// @Param position query string false "Filter by position"
// @Param location_id query integer false "Filter by location ID"
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

// GetEmployeeCountsResponse represents the response for GetEmployeeCountsApi
type GetEmployeeCountsResponse struct {
	TotalEmployees      int64 `json:"total_employees"`
	TotalSubcontractors int64 `json:"total_subcontractors"`
	TotalArchived       int64 `json:"total_archived"`
	TotalOutOfService   int64 `json:"total_out_of_service"`
}

// @Summary Get employee counts
// @Description Get total counts of employees, subcontractors, archived, and out of service employees
// @Tags employees
// @Produce json
// @Success 200 {object} Response[GetEmployeeCountsResponse]
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
// @Success 200 {object} Response[GetEmployeeProfileByIDApiResponse]
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
	employeeID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	profile, err := server.businessService.EmployeeService.GetEmployeeProfileByID(employeeID, currentUserID, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Generate presigned URL for profile picture
	profile.ProfilePicture = server.generateResponsePresignedURL(profile.ProfilePicture)

	res := SuccessResponse(profile, "Employee profile retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Update employee profile by ID
// @Description Update employee profile by ID
// @Tags employees
// @Produce json
// @Param id path int true "Employee ID"
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
// @Param id path int true "Employee ID"
// @Param request body SetEmployeeProfilePictureRequest true "Profile picture details"
// @Success 200 {object} Response[SetEmployeeProfilePictureResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/profile_picture [put]
func (server *Server) SetEmployeeProfilePictureApi(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := strconv.ParseInt(id, 10, 64)
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
// @Param id path int true "Employee ID"
// @Param request body employees.UpdateEmployeeIsSubcontractorRequest true "Subcontractor status details"
// @Success 200 {object} Response[employees.UpdateEmployeeIsSubcontractorResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/is_subcontractor [put]
func (server *Server) UpdateEmployeeIsSubcontractorApi(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := strconv.ParseInt(id, 10, 64)
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
// @Param id path int true "Employee ID"
// @Param request body AddEmployeeContractDetailsRequest true "Contract details"
// @Success 201 {object} Response[AddEmployeeContractDetailsResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/contract_details [put]
func (server *Server) AddEmployeeContractDetailsApi(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := strconv.ParseInt(id, 10, 64)
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
// @Param id path int true "Employee ID"
// @Success 200 {object} Response[GetEmployeeContractDetailsResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/contract_details [get]
func (server *Server) GetEmployeeContractDetailsApi(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := strconv.ParseInt(id, 10, 64)
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
// @Param id path int true "Employee ID"
// @Param request body AddEducationToEmployeeProfileRequest true "Education details"
// @Success 201 {object} Response[AddEducationToEmployeeProfileResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/education [post]
func (server *Server) AddEducationToEmployeeProfileApi(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := strconv.ParseInt(id, 10, 64)
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
// @Param id path int true "Employee ID"
// @Success 200 {object} Response[[]ListEmployeeEducationResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/education [get]
func (server *Server) ListEmployeeEducationApi(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := strconv.ParseInt(id, 10, 64)
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
// @Param id path int true "Employee ID"
// @Param education_id path int true "Education ID"
// @Param request body UpdateEmployeeEducationRequest true "Education details"
// @Success 200 {object} Response[UpdateEmployeeEducationResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/education/{education_id} [put]
func (server *Server) UpdateEmployeeEducationApi(ctx *gin.Context) {
	id := ctx.Param("education_id")
	educationID, err := strconv.ParseInt(id, 10, 64)
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
// @Param id path int true "Employee ID"
// @Param education_id path int true "Education ID"
// @Success 200 {object} Response[DeleteEmployeeEducationResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/education/{education_id} [delete]
func (server *Server) DeleteEmployeeEducationApi(ctx *gin.Context) {
	id := ctx.Param("education_id")
	educationID, err := strconv.ParseInt(id, 10, 64)
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
// @Param id path int true "Employee ID"
// @Success 200 {object} Response[[]ListEmployeeExperienceResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/experience [get]
func (server *Server) ListEmployeeExperienceApi(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := strconv.ParseInt(id, 10, 64)
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
// @Param id path int true "Employee ID"
// @Param experience_id path int true "Experience ID"
// @Param request body UpdateEmployeeExperienceRequest true "Experience details"
// @Success 200 {object} Response[UpdateEmployeeExperienceResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/experience/{experience_id} [put]
func (server *Server) UpdateEmployeeExperienceApi(ctx *gin.Context) {
	id := ctx.Param("experience_id")
	experienceID, err := strconv.ParseInt(id, 10, 64)
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
// @Param id path int true "Employee ID"
// @Param experience_id path int true "Experience ID"
// @Success 200 {object} Response[DeleteEmployeeExperienceResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/experience/{experience_id} [delete]
func (server *Server) DeleteEmployeeExperienceApi(ctx *gin.Context) {
	id := ctx.Param("experience_id")
	experienceID, err := strconv.ParseInt(id, 10, 64)
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
// @Param id path int true "Employee ID"
// @Param request body AddEmployeeCertificationRequest true "Certification details"
// @Success 201 {object} Response[AddEmployeeCertificationResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/certification [post]
func (server *Server) AddEmployeeCertificationApi(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := strconv.ParseInt(id, 10, 64)
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
// @Param id path int true "Employee ID"
// @Success 200 {object} Response[[]ListEmployeeCertificationResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/certification [get]
func (server *Server) ListEmployeeCertificationApi(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := strconv.ParseInt(id, 10, 64)
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
// @Param id path int true "Employee ID"
// @Param certification_id path int true "Certification ID"
// @Param request body UpdateEmployeeCertificationRequest true "Certification details"
// @Success 200 {object} Response[UpdateEmployeeCertificationResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/certification/{certification_id} [put]
func (server *Server) UpdateEmployeeCertificationApi(ctx *gin.Context) {
	id := ctx.Param("certification_id")
	certificationID, err := strconv.ParseInt(id, 10, 64)
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
// @Param id path int true "Employee ID"
// @Param certification_id path int true "Certification ID"
// @Success 200 {object} Response[DeleteEmployeeCertificationResponse]
// @Failure 400,401,404,409,500 {object} Response[any]
// @Router /employees/{id}/certification/{certification_id} [delete]
func (server *Server) DeleteEmployeeCertificationApi(ctx *gin.Context) {
	id := ctx.Param("certification_id")
	certificationID, err := strconv.ParseInt(id, 10, 64)
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
// @Success 200 {object} Response[[]SearchEmployeesByNameOrEmailResponse]
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
