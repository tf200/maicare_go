package api

import (
	"net/http"
	"strconv"

	"maicare_go/service/employees"

	"github.com/gin-gonic/gin"
)

// @Summary List working hours for an employee
// @Description List working hours for an employee in a given month and year
// @Tags Working Hours
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Param year query int true "Year"
// @Param week query int true "Week"
// @Success 200 {object} Response[ListWorkingHoursResponse] "
// @Failure 400 {object} Response[any] "Invalid request parameters"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /employees/{id}/working_hours [get]
func (server *Server) ListWorkingHours(ctx *gin.Context) {
	id := ctx.Param("id")
	employeeID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}
	var req employees.ListWorkingHoursRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	workingHoursResponse, err := server.businessService.EmployeeService.ListWorkingHours(ctx.Request.Context(), employeeID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(workingHoursResponse, "List working hours successfully")
	ctx.JSON(http.StatusOK, res)
}
