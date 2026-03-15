package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/goccy/go-json"

	"maicare_go/pagination"
	"maicare_go/service/handbook"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	errCodeHandbookDraftAlreadyExists   = "HANDBOOK_DRAFT_ALREADY_EXISTS"
	errCodeHandbookTemplateNotFound     = "HANDBOOK_TEMPLATE_NOT_FOUND"
	errCodeHandbookTemplateNotDraft     = "HANDBOOK_TEMPLATE_NOT_DRAFT"
	errCodeHandbookTemplateNotPublished = "HANDBOOK_TEMPLATE_NOT_PUBLISHED"
	errCodeHandbookTemplateNoSteps      = "HANDBOOK_TEMPLATE_NO_STEPS"
	errCodeHandbookStepNotFound         = "HANDBOOK_STEP_NOT_FOUND"
	errCodeEmployeeHandbookNotFound     = "EMPLOYEE_HANDBOOK_NOT_FOUND"
	errCodeEmployeeHandbookNotActive    = "EMPLOYEE_HANDBOOK_NOT_ACTIVE"
	errCodeHandbookLinkURLInvalid       = "HANDBOOK_LINK_URL_INVALID"
	errCodeHandbookQuizContentInvalid   = "HANDBOOK_QUIZ_CONTENT_INVALID"
	errCodeHandbookStepReorderMismatch  = "HANDBOOK_STEP_REORDER_SET_MISMATCH"
	errCodeInvalidRequest               = "INVALID_REQUEST"
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
		ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("invalid step_id"), errCodeInvalidRequest))
		return
	}

	var req completeMyHandbookStepRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("invalid request body"), errCodeInvalidRequest))
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
// @Success 201 {object} Response[handbook.HandbookTemplateAPI]
// @Failure 400,401,403,409,500 {object} Response[any]
// @Router /handbook/templates [post]
func (server *Server) CreateHandbookTemplateApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var req handbook.CreateTemplateForDepartmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("invalid request body"), errCodeInvalidRequest))
		return
	}

	res, err := server.businessService.HandbookService.CreateTemplateForDepartment(ctx, payload.EmployeeID, req)
	if err != nil {
		if errors.Is(err, handbook.ErrDraftTemplateAlreadyExists) {
			ctx.JSON(http.StatusConflict, errorResponseWithCode(err, errCodeHandbookDraftAlreadyExists))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, SuccessResponse(res, "Template created"))
}

// @Summary Clone handbook template to draft (copies steps)
// @Tags handbook
// @Accept json
// @Produce json
// @Param request body handbook.CloneTemplateToDraftRequest true "Clone payload"
// @Success 201 {object} Response[handbook.HandbookTemplateAPI]
// @Failure 400,401,403,409,500 {object} Response[any]
// @Router /handbook/templates/clone [post]
func (server *Server) CloneHandbookTemplateApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var req handbook.CloneTemplateToDraftRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("invalid request body"), errCodeInvalidRequest))
		return
	}

	res, err := server.businessService.HandbookService.CloneTemplateToDraft(ctx, payload.EmployeeID, req)
	if err != nil {
		if errors.Is(err, handbook.ErrDraftTemplateAlreadyExists) {
			ctx.JSON(http.StatusConflict, errorResponseWithCode(err, errCodeHandbookDraftAlreadyExists))
			return
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, SuccessResponse(res, "Template cloned"))
}

type updateHandbookTemplateRequest struct {
	Title       *json.RawMessage `json:"title"`
	Description *json.RawMessage `json:"description"`
}

// @Summary Update handbook template metadata (draft only)
// @Tags handbook
// @Accept json
// @Produce json
// @Param template_id path uuid true "Template ID"
// @Param request body updateHandbookTemplateRequest true "Template metadata payload"
// @Success 200 {object} Response[handbook.HandbookTemplateAPI]
// @Failure 400,401,403,404,500 {object} Response[any]
// @Router /handbook/templates/{template_id} [patch]
func (server *Server) UpdateHandbookTemplateApi(ctx *gin.Context) {
	templateID, err := uuid.Parse(ctx.Param("template_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("invalid template_id"), errCodeInvalidRequest))
		return
	}

	var req updateHandbookTemplateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("invalid request body"), errCodeInvalidRequest))
		return
	}
	if req.Title == nil && req.Description == nil {
		ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("at least one field is required: title or description"), errCodeInvalidRequest))
		return
	}

	var (
		titleVal *string
		setTitle bool
		descVal  *string
		setDesc  bool
	)
	if req.Title != nil {
		setTitle = true
		if string(*req.Title) != "null" {
			var title string
			if err := json.Unmarshal(*req.Title, &title); err != nil {
				ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("title must be a string"), errCodeInvalidRequest))
				return
			}
			if strings.TrimSpace(title) == "" {
				ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("title cannot be empty"), errCodeInvalidRequest))
				return
			}
			titleVal = &title
		}
	}
	if setTitle && titleVal == nil {
		ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("title cannot be null"), errCodeInvalidRequest))
		return
	}
	if req.Description != nil {
		setDesc = true
		if string(*req.Description) != "null" {
			var description string
			if err := json.Unmarshal(*req.Description, &description); err != nil {
				ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("description must be a string or null"), errCodeInvalidRequest))
				return
			}
			descVal = &description
		}
	}

	res, err := server.businessService.HandbookService.UpdateTemplate(ctx, handbook.UpdateTemplateRequest{
		TemplateID:     templateID,
		Title:          titleVal,
		SetTitle:       setTitle,
		Description:    descVal,
		SetDescription: setDesc,
	})
	if err != nil {
		switch {
		case errors.Is(err, handbook.ErrTemplateNotFound):
			ctx.JSON(http.StatusNotFound, errorResponseWithCode(err, errCodeHandbookTemplateNotFound))
		case errors.Is(err, handbook.ErrTemplateNotDraft):
			ctx.JSON(http.StatusBadRequest, errorResponseWithCode(err, errCodeHandbookTemplateNotDraft))
		default:
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(res, "Template updated"))
}

// @Summary Publish handbook template
// @Tags handbook
// @Accept json
// @Produce json
// @Param request body handbook.PublishTemplateRequest true "Publish payload"
// @Success 200 {object} Response[handbook.HandbookTemplateAPI]
// @Failure 400,401,403,404,500 {object} Response[any]
// @Router /handbook/templates/publish [post]
func (server *Server) PublishHandbookTemplateApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var req handbook.PublishTemplateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("invalid request body"), errCodeInvalidRequest))
		return
	}

	res, err := server.businessService.HandbookService.PublishTemplate(ctx, payload.EmployeeID, req)
	if err != nil {
		switch {
		case errors.Is(err, handbook.ErrTemplateNotFound):
			ctx.JSON(http.StatusNotFound, errorResponseWithCode(err, errCodeHandbookTemplateNotFound))
		case errors.Is(err, handbook.ErrTemplateNotDraft):
			ctx.JSON(http.StatusBadRequest, errorResponseWithCode(err, errCodeHandbookTemplateNotDraft))
		case errors.Is(err, handbook.ErrTemplateHasNoSteps):
			ctx.JSON(http.StatusBadRequest, errorResponseWithCode(err, errCodeHandbookTemplateNoSteps))
		default:
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
		}
		return
	}
	ctx.JSON(http.StatusOK, SuccessResponse(res, "Template published"))
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
		ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("invalid department_id"), errCodeInvalidRequest))
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
	Kind       string          `json:"kind" binding:"required,oneof=content ack link quiz rich_text"`
	Title      string          `json:"title" binding:"required"`
	Body       *string         `json:"body"`
	Content    json.RawMessage `json:"content"`
	IsRequired *bool           `json:"is_required"`
}

type updateHandbookStepRequest struct {
	Title      *json.RawMessage `json:"title"`
	Body       *json.RawMessage `json:"body"`
	Content    *json.RawMessage `json:"content"`
	IsRequired *json.RawMessage `json:"is_required"`
}

type reorderHandbookStepsRequest struct {
	OrderedStepIDs []uuid.UUID `json:"ordered_step_ids" binding:"required,min=1"`
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
		ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("invalid request body"), errCodeInvalidRequest))
		return
	}

	var content any
	if len(req.Content) > 0 {
		if err := json.Unmarshal(req.Content, &content); err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("invalid content JSON"), errCodeInvalidRequest))
			return
		}
	}

	res, err := server.businessService.HandbookService.CreateStep(ctx, handbook.CreateStepRequest{
		TemplateID: req.TemplateID,
		SortOrder:  req.SortOrder,
		Kind:       req.Kind,
		Title:      req.Title,
		Body:       req.Body,
		Content:    req.Content,
		IsRequired: req.IsRequired,
	})
	if err != nil {
		switch {
		case errors.Is(err, handbook.ErrTemplateNotDraft):
			ctx.JSON(http.StatusBadRequest, errorResponseWithCode(err, errCodeHandbookTemplateNotDraft))
		case errors.Is(err, handbook.ErrInvalidStepContent):
			ctx.JSON(http.StatusBadRequest, errorResponseWithCode(err, stepContentErrorCode(err)))
		default:
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}
	ctx.JSON(http.StatusCreated, SuccessResponse(res, "Step created"))
}

// @Summary Update handbook step (draft templates only)
// @Tags handbook
// @Accept json
// @Produce json
// @Param step_id path uuid true "Step ID"
// @Param request body updateHandbookStepRequest true "Step update payload"
// @Success 200 {object} Response[handbook.UpdateStepResponse]
// @Failure 400,401,403,404,500 {object} Response[any]
// @Router /handbook/steps/{step_id} [patch]
func (server *Server) UpdateHandbookStepApi(ctx *gin.Context) {
	stepID, err := uuid.Parse(ctx.Param("step_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("invalid step_id"), errCodeInvalidRequest))
		return
	}

	var req updateHandbookStepRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("invalid request body"), errCodeInvalidRequest))
		return
	}
	if req.Title == nil && req.Body == nil && req.Content == nil && req.IsRequired == nil {
		ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("at least one field is required: title, body, content or is_required"), errCodeInvalidRequest))
		return
	}

	var (
		titleVal      *string
		setTitle      bool
		bodyVal       *string
		setBody       bool
		content       []byte
		contentSet    bool
		isRequiredVal *bool
		setIsRequired bool
	)
	if req.Title != nil {
		setTitle = true
		if string(*req.Title) == "null" {
			ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("title cannot be null"), errCodeInvalidRequest))
			return
		}
		var title string
		if err := json.Unmarshal(*req.Title, &title); err != nil || strings.TrimSpace(title) == "" {
			ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("title must be a non-empty string"), errCodeInvalidRequest))
			return
		}
		titleVal = &title
	}
	if req.Body != nil {
		setBody = true
		if string(*req.Body) != "null" {
			var body string
			if err := json.Unmarshal(*req.Body, &body); err != nil {
				ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("body must be a string or null"), errCodeInvalidRequest))
				return
			}
			bodyVal = &body
		}
	}
	if req.Content != nil {
		contentSet = true
		var payload any
		if err := json.Unmarshal(*req.Content, &payload); err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("invalid content JSON"), errCodeInvalidRequest))
			return
		}
		content = append([]byte(nil), (*req.Content)...)
	}
	if req.IsRequired != nil {
		setIsRequired = true
		if string(*req.IsRequired) == "null" {
			ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("is_required cannot be null"), errCodeInvalidRequest))
			return
		}
		var isRequired bool
		if err := json.Unmarshal(*req.IsRequired, &isRequired); err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("is_required must be a boolean"), errCodeInvalidRequest))
			return
		}
		isRequiredVal = &isRequired
	}

	res, err := server.businessService.HandbookService.UpdateStep(ctx, handbook.UpdateStepRequest{
		StepID:          stepID,
		Title:           titleVal,
		SetTitle:        setTitle,
		Body:            bodyVal,
		SetBody:         setBody,
		Content:         content,
		ContentProvided: contentSet,
		IsRequired:      isRequiredVal,
		SetIsRequired:   setIsRequired,
	})
	if err != nil {
		switch {
		case errors.Is(err, handbook.ErrStepNotFound):
			ctx.JSON(http.StatusNotFound, errorResponseWithCode(err, errCodeHandbookStepNotFound))
		case errors.Is(err, handbook.ErrTemplateNotDraft):
			ctx.JSON(http.StatusBadRequest, errorResponseWithCode(err, errCodeHandbookTemplateNotDraft))
		case errors.Is(err, handbook.ErrInvalidStepContent):
			ctx.JSON(http.StatusBadRequest, errorResponseWithCode(err, stepContentErrorCode(err)))
		default:
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}
	ctx.JSON(http.StatusOK, SuccessResponse(res, "Step updated"))
}

// @Summary Delete handbook step (draft templates only)
// @Tags handbook
// @Produce json
// @Param step_id path uuid true "Step ID"
// @Success 200 {object} Response[handbook.DeleteStepResponse]
// @Failure 400,401,403,404,500 {object} Response[any]
// @Router /handbook/steps/{step_id} [delete]
func (server *Server) DeleteHandbookStepApi(ctx *gin.Context) {
	stepID, err := uuid.Parse(ctx.Param("step_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("invalid step_id"), errCodeInvalidRequest))
		return
	}

	res, err := server.businessService.HandbookService.DeleteStep(ctx, handbook.DeleteStepRequest{StepID: stepID})
	if err != nil {
		switch {
		case errors.Is(err, handbook.ErrStepNotFound):
			ctx.JSON(http.StatusNotFound, errorResponseWithCode(err, errCodeHandbookStepNotFound))
		case errors.Is(err, handbook.ErrTemplateNotDraft):
			ctx.JSON(http.StatusBadRequest, errorResponseWithCode(err, errCodeHandbookTemplateNotDraft))
		default:
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}
	ctx.JSON(http.StatusOK, SuccessResponse(res, "Step deleted"))
}

// @Summary Reorder handbook steps (draft templates only)
// @Tags handbook
// @Accept json
// @Produce json
// @Param template_id path uuid true "Template ID"
// @Param request body reorderHandbookStepsRequest true "Step reorder payload"
// @Success 200 {object} Response[handbook.ReorderStepsResponse]
// @Failure 400,401,403,404,500 {object} Response[any]
// @Router /handbook/templates/{template_id}/steps/reorder [post]
func (server *Server) ReorderHandbookStepsApi(ctx *gin.Context) {
	templateID, err := uuid.Parse(ctx.Param("template_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("invalid template_id"), errCodeInvalidRequest))
		return
	}

	var req reorderHandbookStepsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("invalid request body"), errCodeInvalidRequest))
		return
	}

	res, err := server.businessService.HandbookService.ReorderTemplateSteps(ctx, handbook.ReorderStepsRequest{
		TemplateID:     templateID,
		OrderedStepIDs: req.OrderedStepIDs,
	})
	if err != nil {
		switch {
		case errors.Is(err, handbook.ErrTemplateNotFound):
			ctx.JSON(http.StatusNotFound, errorResponseWithCode(err, errCodeHandbookTemplateNotFound))
		case errors.Is(err, handbook.ErrTemplateNotDraft):
			ctx.JSON(http.StatusBadRequest, errorResponseWithCode(err, errCodeHandbookTemplateNotDraft))
		case errors.Is(err, handbook.ErrInvalidStepReorder):
			ctx.JSON(http.StatusBadRequest, errorResponseWithCode(err, errCodeHandbookStepReorderMismatch))
		default:
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}
	ctx.JSON(http.StatusOK, SuccessResponse(res, "Steps reordered"))
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
		ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("invalid template_id"), errCodeInvalidRequest))
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
		ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("invalid request body"), errCodeInvalidRequest))
		return
	}
	res, err := server.businessService.HandbookService.AssignTemplateToEmployee(ctx, payload.EmployeeID, req)
	if err != nil {
		if errors.Is(err, handbook.ErrTemplateNotPublished) {
			ctx.JSON(http.StatusBadRequest, errorResponseWithCode(err, errCodeHandbookTemplateNotPublished))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, SuccessResponse(res, "Handbook assigned"))
}

// @Summary List handbook assignment history for an employee
// @Tags handbook
// @Produce json
// @Param employee_id path uuid true "Employee ID"
// @Success 200 {object} Response[[]handbook.HandbookAssignmentHistoryEntry]
// @Failure 400,401,403,500 {object} Response[any]
// @Router /handbook/employees/{employee_id}/history [get]
func (server *Server) ListEmployeeHandbookHistoryApi(ctx *gin.Context) {
	employeeID, err := uuid.Parse(ctx.Param("employee_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("invalid employee_id"), errCodeInvalidRequest))
		return
	}

	res, err := server.businessService.HandbookService.ListEmployeeHandbookHistory(ctx, employeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(res, "Handbook history retrieved"))
}

type waiveEmployeeHandbookRequest struct {
	Reason *string `json:"reason"`
}

type listHandbookAssignmentsQuery struct {
	pagination.Request
	DepartmentID *string `form:"department_id"`
	Search       *string `form:"search"`
	Status       *string `form:"status"`
}

type listEligibleEmployeesQuery struct {
	pagination.Request
	DepartmentID *string `form:"department_id"`
	Search       *string `form:"search"`
}

// @Summary List employee handbook assignments
// @Tags handbook
// @Produce json
// @Param page query int true "Page number"
// @Param page_size query int true "Page size"
// @Param department_id query uuid false "Department ID"
// @Param search query string false "Employee name search"
// @Param status query string false "Assignment status: unassigned, not_started, in_progress, completed, waived"
// @Success 200 {object} Response[pagination.Response[handbook.EmployeeHandbookAssignmentSummary]]
// @Failure 400,401,403,500 {object} Response[any]
// @Router /handbook/assignments [get]
func (server *Server) ListEmployeeHandbookAssignmentsApi(ctx *gin.Context) {
	var query listHandbookAssignmentsQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("invalid query parameters"), errCodeInvalidRequest))
		return
	}

	departmentID, err := parseOptionalUUIDQueryParam(query.DepartmentID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("invalid department_id"), errCodeInvalidRequest))
		return
	}

	req := handbook.ListEmployeeHandbookAssignmentsRequest{
		Request:      query.Request,
		DepartmentID: departmentID,
		Search:       query.Search,
		Status:       query.Status,
	}

	res, err := server.businessService.HandbookService.ListEmployeeHandbookAssignments(ctx, req)
	if err != nil {
		if errors.Is(err, handbook.ErrInvalidAssignmentStatusFilter) {
			ctx.JSON(http.StatusBadRequest, errorResponseWithCode(err, errCodeInvalidRequest))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(res, "Employee handbook assignments retrieved"))
}

// @Summary List employees eligible for handbook assignment
// @Tags handbook
// @Produce json
// @Param page query int true "Page number"
// @Param page_size query int true "Page size"
// @Param department_id query uuid false "Department ID"
// @Param search query string false "Employee name search"
// @Success 200 {object} Response[pagination.Response[handbook.ListEligibleEmployeesResponse]]
// @Failure 400,401,403,500 {object} Response[any]
// @Router /handbook/assignments/eligible-employees [get]
func (server *Server) ListEligibleEmployeesApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var query listEligibleEmployeesQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("invalid query parameters"), errCodeInvalidRequest))
		return
	}

	departmentID, err := parseOptionalUUIDQueryParam(query.DepartmentID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("invalid department_id"), errCodeInvalidRequest))
		return
	}

	req := handbook.ListEligibleEmployeesRequest{
		Request:      query.Request,
		DepartmentID: departmentID,
		Search:       query.Search,
	}

	res, err := server.businessService.HandbookService.ListEligibleEmployees(ctx, payload.EmployeeID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(res, "Eligible employees retrieved"))
}

// @Summary Get employee handbook details
// @Tags handbook
// @Produce json
// @Param handbook_id path uuid true "Employee handbook ID"
// @Success 200 {object} Response[handbook.GetEmployeeHandbookDetailsResponse]
// @Failure 400,401,403,404,500 {object} Response[any]
// @Router /handbook/assignments/{handbook_id} [get]
func (server *Server) GetEmployeeHandbookDetailsApi(ctx *gin.Context) {
	handbookID, err := uuid.Parse(ctx.Param("handbook_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("invalid handbook_id"), errCodeInvalidRequest))
		return
	}

	res, err := server.businessService.HandbookService.GetEmployeeHandbookDetails(ctx, handbookID)
	if err != nil {
		if errors.Is(err, handbook.ErrEmployeeHandbookNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponseWithCode(err, errCodeEmployeeHandbookNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(res, "Employee handbook details retrieved"))
}

// @Summary Waive an active employee handbook
// @Tags handbook
// @Accept json
// @Produce json
// @Param handbook_id path uuid true "Employee handbook ID"
// @Param request body waiveEmployeeHandbookRequest false "Waive payload"
// @Success 200 {object} Response[handbook.WaiveEmployeeHandbookResponse]
// @Failure 400,401,403,404,500 {object} Response[any]
// @Router /handbook/assignments/{handbook_id}/waive [post]
func (server *Server) WaiveEmployeeHandbookApi(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	handbookID, err := uuid.Parse(ctx.Param("handbook_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("invalid handbook_id"), errCodeInvalidRequest))
		return
	}

	var req waiveEmployeeHandbookRequest
	if ctx.Request.ContentLength > 0 {
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponseWithCode(fmt.Errorf("invalid request body"), errCodeInvalidRequest))
			return
		}
	}

	res, err := server.businessService.HandbookService.WaiveEmployeeHandbook(ctx, payload.EmployeeID, handbook.WaiveEmployeeHandbookRequest{
		EmployeeHandbookID: handbookID,
		Reason:             req.Reason,
	})
	if err != nil {
		switch {
		case errors.Is(err, handbook.ErrEmployeeHandbookNotFound):
			ctx.JSON(http.StatusNotFound, errorResponseWithCode(err, errCodeEmployeeHandbookNotFound))
		case errors.Is(err, handbook.ErrEmployeeHandbookNotActive):
			ctx.JSON(http.StatusBadRequest, errorResponseWithCode(err, errCodeEmployeeHandbookNotActive))
		default:
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(res, "Handbook waived"))
}

func stepContentErrorCode(err error) string {
	if strings.Contains(strings.ToLower(err.Error()), "quiz") {
		return errCodeHandbookQuizContentInvalid
	}
	return errCodeHandbookLinkURLInvalid
}

func parseOptionalUUIDQueryParam(raw *string) (*uuid.UUID, error) {
	if raw == nil {
		return nil, nil
	}

	trimmed := strings.TrimSpace(*raw)
	if trimmed == "" {
		return nil, nil
	}

	parsed, err := uuid.Parse(trimmed)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}
