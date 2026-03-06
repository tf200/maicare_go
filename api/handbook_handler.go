package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "maicare_go/db/sqlc"
	"maicare_go/service/handbook"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Get my active handbook
// @Tags handbook
// @Produce json
// @Success 200 {object} Response[handbook.GetMyActiveHandbookResponse]
// @Failure 400,401,403,404,500 {object} Response[any]
// @Router /handbook/me [get]
func (server *Server) GetMyActiveHandbookApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	res, err := server.businessService.HandbookService.GetMyActiveHandbook(ctx, payload.EmployeeID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, SuccessResponse(res, "Active handbook retrieved"))
}

// @Summary Start my handbook
// @Tags handbook
// @Produce json
// @Success 200 {object} Response[handbook.StartMyHandbookResponse]
// @Failure 400,401,403,404,500 {object} Response[any]
// @Router /handbook/me/start [post]
func (server *Server) StartMyHandbookApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	res, err := server.businessService.HandbookService.StartMyHandbook(ctx, payload.EmployeeID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, SuccessResponse(res, "Handbook started"))
}

type completeMyHandbookStepRequest struct {
	Response json.RawMessage `json:"response"`
}

// @Summary Complete a handbook step
// @Tags handbook
// @Accept json
// @Produce json
// @Param step_id path uuid true "Step ID"
// @Param request body completeMyHandbookStepRequest true "Step completion payload"
// @Success 200 {object} Response[handbook.CompleteMyHandbookStepResponse]
// @Failure 400,401,403,404,500 {object} Response[any]
// @Router /handbook/me/steps/{step_id}/complete [post]
func (server *Server) CompleteMyHandbookStepApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	stepID, err := uuid.Parse(ctx.Param("step_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid step_id")))
		return
	}

	var req completeMyHandbookStepRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	res, err := server.businessService.HandbookService.CompleteMyHandbookStep(ctx, payload.EmployeeID, stepID, req.Response)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, SuccessResponse(res, "Step completed"))
}

// @Summary Create handbook template (new version for department)
// @Tags handbook
// @Accept json
// @Produce json
// @Param request body handbook.CreateTemplateForDepartmentRequest true "Template payload"
// @Success 201 {object} Response[handbook.CreateTemplateForDepartmentResponse]
// @Failure 400,401,403,500 {object} Response[any]
// @Router /handbook/templates [post]
func (server *Server) CreateHandbookTemplateApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var req handbook.CreateTemplateForDepartmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	res, err := server.businessService.HandbookService.CreateTemplateForDepartment(ctx, payload.EmployeeID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, SuccessResponse(res, "Template created"))
}

// @Summary List templates by department
// @Tags handbook
// @Produce json
// @Param department_id path uuid true "Department ID"
// @Success 200 {object} Response[any]
// @Failure 400,401,403,500 {object} Response[any]
// @Router /handbook/departments/{department_id}/templates [get]
func (server *Server) ListHandbookTemplatesByDepartmentApi(ctx *gin.Context) {
	deptID, err := uuid.Parse(ctx.Param("department_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid department_id")))
		return
	}
	res, err := server.businessService.HandbookService.ListTemplatesByDepartment(ctx, deptID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, SuccessResponse(res, "Templates retrieved"))
}

type createHandbookStepRequest struct {
	TemplateID uuid.UUID       `json:"template_id" binding:"required"`
	SortOrder  int32           `json:"sort_order" binding:"required,min=1"`
	Kind       string          `json:"kind" binding:"required,oneof=content ack link quiz"`
	Title      string          `json:"title" binding:"required"`
	Body       *string         `json:"body"`
	Content    json.RawMessage `json:"content"`
	IsRequired *bool           `json:"is_required"`
}

// @Summary Create handbook step
// @Tags handbook
// @Accept json
// @Produce json
// @Param request body createHandbookStepRequest true "Step payload"
// @Success 201 {object} Response[handbook.CreateStepResponse]
// @Failure 400,401,403,500 {object} Response[any]
// @Router /handbook/steps [post]
func (server *Server) CreateHandbookStepApi(ctx *gin.Context) {
	var req createHandbookStepRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	var content any
	if len(req.Content) > 0 {
		if err := json.Unmarshal(req.Content, &content); err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid content JSON")))
			return
		}
	}

	res, err := server.businessService.HandbookService.CreateStep(ctx, handbook.CreateStepRequest{
		TemplateID: req.TemplateID,
		SortOrder:  req.SortOrder,
		Kind:       db.HandbookStepKindEnum(req.Kind),
		Title:      req.Title,
		Body:       req.Body,
		Content:    content,
		IsRequired: req.IsRequired,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, SuccessResponse(res, "Step created"))
}

// @Summary List steps by template
// @Tags handbook
// @Produce json
// @Param template_id path uuid true "Template ID"
// @Success 200 {object} Response[any]
// @Failure 400,401,403,500 {object} Response[any]
// @Router /handbook/templates/{template_id}/steps [get]
func (server *Server) ListHandbookStepsByTemplateApi(ctx *gin.Context) {
	templateID, err := uuid.Parse(ctx.Param("template_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid template_id")))
		return
	}
	res, err := server.businessService.HandbookService.ListStepsByTemplate(ctx, templateID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, SuccessResponse(res, "Steps retrieved"))
}

// @Summary Assign handbook template to an employee
// @Tags handbook
// @Accept json
// @Produce json
// @Param request body handbook.AssignTemplateToEmployeeRequest true "Assignment payload"
// @Success 201 {object} Response[handbook.AssignTemplateToEmployeeResponse]
// @Failure 400,401,403,500 {object} Response[any]
// @Router /handbook/assignments [post]
func (server *Server) AssignHandbookTemplateToEmployeeApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var req handbook.AssignTemplateToEmployeeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}
	res, err := server.businessService.HandbookService.AssignTemplateToEmployee(ctx, payload.EmployeeID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, SuccessResponse(res, "Handbook assigned"))
}
