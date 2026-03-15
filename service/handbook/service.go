package handbook

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/goccy/go-json"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"
	"maicare_go/service/deps"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type handbookService struct {
	*deps.ServiceDependencies
}

var ErrDraftTemplateAlreadyExists = errors.New("a draft template already exists for this department")
var ErrTemplateNotFound = errors.New("template not found")
var ErrTemplateNotDraft = errors.New("template is not in draft status")
var ErrTemplateNotPublished = errors.New("template is not in published status")
var ErrTemplateHasNoSteps = errors.New("template must contain at least one step before publishing")
var ErrStepNotFound = errors.New("step not found")
var ErrInvalidStepReorder = errors.New("ordered_step_ids must match the template steps exactly")
var ErrInvalidStepContent = errors.New("invalid step content")
var ErrEmployeeHandbookNotFound = errors.New("employee handbook not found")
var ErrEmployeeHandbookNotActive = errors.New("employee handbook is not active")
var ErrInvalidAssignmentStatusFilter = errors.New("invalid assignment status filter")

func NewHandbookService(deps *deps.ServiceDependencies) HandbookService {
	return &handbookService{ServiceDependencies: deps}
}

func (s *handbookService) GetMyActiveHandbook(ctx context.Context, employeeID uuid.UUID) (*GetMyActiveHandbookResponse, error) {
	hb, err := s.Store.GetActiveEmployeeHandbookByEmployeeID(ctx, employeeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no active handbook assigned")
		}
		return nil, err
	}

	rows, err := s.Store.ListEmployeeHandbookStepsByHandbookID(ctx, hb.ID)
	if err != nil {
		return nil, err
	}

	steps := make([]MyHandbookStep, 0, len(rows))
	for _, row := range rows {
		steps = append(steps, MyHandbookStep{
			StepID:      row.StepID,
			SortOrder:   row.SortOrder,
			Kind:        row.Kind,
			Title:       row.Title,
			Body:        row.Body,
			Content:     mustUnmarshalJSON(row.Content),
			IsRequired:  row.IsRequired,
			Status:      string(row.ProgressStatus),
			StartedAt:   util.TsPtr(row.ProgressStartedAt),
			CompletedAt: util.TsPtr(row.ProgressCompletedAt),
			Response:    mustUnmarshalJSON(row.ProgressResponse),
		})
	}

	res := &GetMyActiveHandbookResponse{
		HandbookID:      hb.ID,
		Status:          string(hb.Status),
		AssignedAt:      hb.AssignedAt.Time,
		StartedAt:       util.TsPtr(hb.StartedAt),
		CompletedAt:     util.TsPtr(hb.CompletedAt),
		DueAt:           util.TsPtr(hb.DueAt),
		TemplateID:      hb.TemplateID,
		TemplateTitle:   hb.TemplateTitle,
		TemplateDesc:    hb.TemplateDescription,
		TemplateVersion: hb.TemplateVersion,
		DepartmentID:    hb.DepartmentID,
		DepartmentName:  hb.DepartmentName,
		Steps:           steps,
	}

	return res, nil
}

func (s *handbookService) StartMyHandbook(ctx context.Context, employeeID uuid.UUID) (*StartMyHandbookResponse, error) {
	hb, err := s.Store.GetActiveEmployeeHandbookByEmployeeID(ctx, employeeID)
	if err != nil {
		return nil, err
	}

	var updated db.EmployeeHandbook
	err = s.Store.ExecTx(ctx, func(q *db.Queries) error {
		row, err := q.MarkEmployeeHandbookStarted(ctx, hb.ID)
		if err != nil {
			return err
		}
		updated = row

		if hb.Status == db.HandbookAssignmentStatusEnumNotStarted {
			metadata, err := marshalJSONPayload(map[string]any{
				"source": "employee_self_service",
			})
			if err != nil {
				return err
			}
			_, err = q.CreateEmployeeHandbookAssignmentHistory(ctx, db.CreateEmployeeHandbookAssignmentHistoryParams{
				EmployeeHandbookID: &hb.ID,
				EmployeeID:         hb.EmployeeID,
				TemplateID:         hb.TemplateID,
				TemplateVersion:    hb.TemplateVersion,
				Event:              db.HandbookAssignmentEventEnumStarted,
				ActorEmployeeID:    &hb.EmployeeID,
				Metadata:           metadata,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return &StartMyHandbookResponse{
		HandbookID: updated.ID,
		Status:     string(updated.Status),
		StartedAt:  util.TsPtr(updated.StartedAt),
	}, nil
}

func (s *handbookService) CompleteMyHandbookStep(ctx context.Context, employeeID, stepID uuid.UUID, response any) (*CompleteMyHandbookStepResponse, error) {
	hb, err := s.Store.GetActiveEmployeeHandbookByEmployeeID(ctx, employeeID)
	if err != nil {
		return nil, err
	}

	respBytes, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("invalid response payload: %w", err)
	}

	var progress db.EmployeeHandbookStepProgress
	handbookStatus := string(hb.Status)
	err = s.Store.ExecTx(ctx, func(q *db.Queries) error {
		row, err := q.CompleteEmployeeHandbookStep(ctx, db.CompleteEmployeeHandbookStepParams{
			Response:           respBytes,
			EmployeeHandbookID: hb.ID,
			StepID:             stepID,
		})
		if err != nil {
			return err
		}
		progress = row

		remaining, err := q.CountRemainingRequiredHandbookSteps(ctx, hb.ID)
		if err != nil {
			return err
		}

		if remaining == 0 {
			completed, err := q.MarkEmployeeHandbookCompleted(ctx, hb.ID)
			if err != nil {
				return err
			}
			handbookStatus = string(completed.Status)

			metadata, err := marshalJSONPayload(map[string]any{
				"source":            "employee_self_service",
				"completed_step_id": stepID.String(),
			})
			if err != nil {
				return err
			}
			_, err = q.CreateEmployeeHandbookAssignmentHistory(ctx, db.CreateEmployeeHandbookAssignmentHistoryParams{
				EmployeeHandbookID: &hb.ID,
				EmployeeID:         hb.EmployeeID,
				TemplateID:         hb.TemplateID,
				TemplateVersion:    hb.TemplateVersion,
				Event:              db.HandbookAssignmentEventEnumCompleted,
				ActorEmployeeID:    &hb.EmployeeID,
				Metadata:           metadata,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &CompleteMyHandbookStepResponse{
		HandbookID:     progress.EmployeeHandbookID,
		StepID:         progress.StepID,
		StepStatus:     string(progress.Status),
		CompletedAt:    progress.CompletedAt.Time,
		HandbookStatus: handbookStatus,
	}, nil
}

func (s *handbookService) CreateTemplateForDepartment(ctx context.Context, actorEmployeeID uuid.UUID, req CreateTemplateForDepartmentRequest) (*HandbookTemplateAPI, error) {
	t, err := s.Store.CreateHandbookTemplateForDepartment(ctx, db.CreateHandbookTemplateForDepartmentParams{
		DepartmentID:        req.DepartmentID,
		Title:               req.Title,
		Description:         req.Description,
		CreatedByEmployeeID: uuidPtrOrNil(actorEmployeeID),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "handbook_templates_one_draft_per_department" {
			return nil, ErrDraftTemplateAlreadyExists
		}
		return nil, err
	}
	return &HandbookTemplateAPI{
		ID:           t.ID,
		DepartmentID: t.DepartmentID,
		Title:        t.Title,
		Description:  t.Description,
		Version:      t.Version,
		Status:       string(t.Status),
		PublishedAt:  util.TsPtr(t.PublishedAt),
		ArchivedAt:   util.TsPtr(t.ArchivedAt),
		CreatedAt:    t.CreatedAt.Time,
		UpdatedAt:    t.UpdatedAt.Time,
	}, nil
}

func (s *handbookService) CloneTemplateToDraft(ctx context.Context, actorEmployeeID uuid.UUID, req CloneTemplateToDraftRequest) (*HandbookTemplateAPI, error) {
	t, err := s.Store.CloneHandbookTemplateToDraft(ctx, db.CloneHandbookTemplateToDraftParams{
		SourceTemplateID:    req.SourceTemplateID,
		CreatedByEmployeeID: uuidPtrOrNil(actorEmployeeID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("source template not found")
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "handbook_templates_one_draft_per_department" {
			return nil, ErrDraftTemplateAlreadyExists
		}
		return nil, err
	}
	return &HandbookTemplateAPI{
		ID:           t.ID,
		DepartmentID: t.DepartmentID,
		Title:        t.Title,
		Description:  t.Description,
		Version:      t.Version,
		Status:       string(t.Status),
		PublishedAt:  util.TsPtr(t.PublishedAt),
		ArchivedAt:   util.TsPtr(t.ArchivedAt),
		CreatedAt:    t.CreatedAt.Time,
		UpdatedAt:    t.UpdatedAt.Time,
	}, nil
}

func (s *handbookService) UpdateTemplate(ctx context.Context, req UpdateTemplateRequest) (*HandbookTemplateAPI, error) {
	tmpl, err := s.Store.GetHandbookTemplateByID(ctx, req.TemplateID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTemplateNotFound
		}
		return nil, err
	}
	if tmpl.Status != db.HandbookTemplateStatusEnumDraft {
		return nil, ErrTemplateNotDraft
	}

	updated, err := s.Store.UpdateHandbookTemplateMetadata(ctx, db.UpdateHandbookTemplateMetadataParams{
		TemplateID:     req.TemplateID,
		Title:          req.Title,
		SetTitle:       req.SetTitle,
		Description:    req.Description,
		SetDescription: req.SetDescription,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTemplateNotDraft
		}
		return nil, err
	}
	api := mapTemplateAPI(updated)
	return &api, nil
}

func (s *handbookService) PublishTemplate(ctx context.Context, actorEmployeeID uuid.UUID, req PublishTemplateRequest) (*HandbookTemplateAPI, error) {
	tmpl, err := s.Store.GetHandbookTemplateByID(ctx, req.TemplateID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTemplateNotFound
		}
		return nil, err
	}
	if tmpl.Status != db.HandbookTemplateStatusEnumDraft {
		return nil, ErrTemplateNotDraft
	}

	stepCount, err := s.Store.CountHandbookStepsByTemplateID(ctx, req.TemplateID)
	if err != nil {
		return nil, err
	}
	if stepCount == 0 {
		return nil, ErrTemplateHasNoSteps
	}

	t, err := s.Store.PublishHandbookTemplate(ctx, db.PublishHandbookTemplateParams{
		TemplateID:            req.TemplateID,
		PublishedByEmployeeID: uuidPtrOrNil(actorEmployeeID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTemplateNotDraft
		}
		return nil, err
	}
	api := mapTemplateAPI(t)
	return &api, nil
}

func (s *handbookService) ListTemplatesByDepartment(ctx context.Context, departmentID uuid.UUID) (*pagination.Response[HandbookTemplateAPI], error) {
	templates, err := s.Store.ListHandbookTemplatesByDepartment(ctx, departmentID)
	if err != nil {
		return nil, err
	}
	out := make([]HandbookTemplateAPI, 0, len(templates))
	for _, t := range templates {
		out = append(out, mapTemplateAPI(t))
	}

	// No real pagination yet (templates are small). Use a stable wrapper for the frontend.
	resp := &pagination.Response[HandbookTemplateAPI]{
		Next:     nil,
		Previous: nil,
		Count:    int64(len(out)),
		PageSize: int32(len(out)),
		Results:  out,
	}
	return resp, nil
}

func (s *handbookService) CreateStep(ctx context.Context, req CreateStepRequest) (*CreateStepResponse, error) {
	tmpl, err := s.Store.GetHandbookTemplateByID(ctx, req.TemplateID)
	if err != nil {
		return nil, err
	}
	if tmpl.Status != db.HandbookTemplateStatusEnumDraft {
		return nil, ErrTemplateNotDraft
	}

	dbKind := mapFrontendKindToDB(req.Kind)

	normalizedContent, err := normalizeAndValidateContent(dbKind, req.Content)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidStepContent, err)
	}

	step, err := s.Store.CreateHandbookStep(ctx, db.CreateHandbookStepParams{
		TemplateID: req.TemplateID,
		SortOrder:  req.SortOrder,
		Kind:       dbKind,
		Title:      req.Title,
		Body:       req.Body,
		Content:    normalizedContent,
		IsRequired: req.IsRequired,
	})
	if err != nil {
		return nil, err
	}
	return &CreateStepResponse{
		ID:         step.ID,
		TemplateID: step.TemplateID,
		SortOrder:  step.SortOrder,
		Kind:       step.Kind,
		Title:      step.Title,
		Body:       step.Body,
		Content:    mustUnmarshalJSON(step.Content),
		IsRequired: step.IsRequired,
	}, nil
}

func (s *handbookService) ListStepsByTemplate(ctx context.Context, templateID uuid.UUID) ([]ListStepResponse, error) {
	steps, err := s.Store.ListHandbookStepsByTemplate(ctx, templateID)
	if err != nil {
		return nil, err
	}
	out := make([]ListStepResponse, 0, len(steps))
	for _, step := range steps {
		out = append(out, ListStepResponse{
			ID:         step.ID,
			SortOrder:  step.SortOrder,
			Kind:       step.Kind,
			Title:      step.Title,
			Body:       step.Body,
			Content:    mustUnmarshalJSON(step.Content),
			IsRequired: step.IsRequired,
		})
	}
	return out, nil
}

func (s *handbookService) UpdateStep(ctx context.Context, req UpdateStepRequest) (*UpdateStepResponse, error) {
	step, err := s.Store.GetHandbookStepByID(ctx, req.StepID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrStepNotFound
		}
		return nil, err
	}

	tmpl, err := s.Store.GetHandbookTemplateByID(ctx, step.TemplateID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTemplateNotFound
		}
		return nil, err
	}
	if tmpl.Status != db.HandbookTemplateStatusEnumDraft {
		return nil, ErrTemplateNotDraft
	}

	var content []byte
	if req.ContentProvided {
		normalizedContent, err := normalizeAndValidateContent(step.Kind, req.Content)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidStepContent, err)
		}
		content = normalizedContent
	}

	updated, err := s.Store.UpdateHandbookStepByID(ctx, db.UpdateHandbookStepByIDParams{
		StepID:        req.StepID,
		Title:         req.Title,
		SetTitle:      req.SetTitle,
		Body:          req.Body,
		SetBody:       req.SetBody,
		Content:       content,
		SetContent:    req.ContentProvided,
		IsRequired:    req.IsRequired,
		SetIsRequired: req.SetIsRequired,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrStepNotFound
		}
		return nil, err
	}

	return &UpdateStepResponse{
		ID:         updated.ID,
		TemplateID: updated.TemplateID,
		SortOrder:  updated.SortOrder,
		Kind:       updated.Kind,
		Title:      updated.Title,
		Body:       updated.Body,
		Content:    mustUnmarshalJSON(updated.Content),
		IsRequired: updated.IsRequired,
		UpdatedAt:  updated.UpdatedAt.Time,
	}, nil
}

func (s *handbookService) DeleteStep(ctx context.Context, req DeleteStepRequest) (*DeleteStepResponse, error) {
	step, err := s.Store.GetHandbookStepByID(ctx, req.StepID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrStepNotFound
		}
		return nil, err
	}

	tmpl, err := s.Store.GetHandbookTemplateByID(ctx, step.TemplateID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTemplateNotFound
		}
		return nil, err
	}
	if tmpl.Status != db.HandbookTemplateStatusEnumDraft {
		return nil, ErrTemplateNotDraft
	}

	err = s.Store.ExecTx(ctx, func(q *db.Queries) error {
		if err := q.DeleteHandbookStepByID(ctx, req.StepID); err != nil {
			return err
		}
		remaining, err := q.ListHandbookStepsByTemplate(ctx, step.TemplateID)
		if err != nil {
			return err
		}
		for i, st := range remaining {
			if err := q.UpdateHandbookStepSortOrder(ctx, db.UpdateHandbookStepSortOrderParams{
				StepID:    st.ID,
				SortOrder: int32(i + 1),
			}); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &DeleteStepResponse{
		StepID:  req.StepID,
		Deleted: true,
	}, nil
}

func (s *handbookService) ReorderTemplateSteps(ctx context.Context, req ReorderStepsRequest) (*ReorderStepsResponse, error) {
	tmpl, err := s.Store.GetHandbookTemplateByID(ctx, req.TemplateID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTemplateNotFound
		}
		return nil, err
	}
	if tmpl.Status != db.HandbookTemplateStatusEnumDraft {
		return nil, ErrTemplateNotDraft
	}

	currentSteps, err := s.Store.ListHandbookStepsByTemplate(ctx, req.TemplateID)
	if err != nil {
		return nil, err
	}
	if len(currentSteps) != len(req.OrderedStepIDs) {
		return nil, ErrInvalidStepReorder
	}

	allowed := make(map[uuid.UUID]struct{}, len(currentSteps))
	for _, s := range currentSteps {
		allowed[s.ID] = struct{}{}
	}

	seen := make(map[uuid.UUID]struct{}, len(req.OrderedStepIDs))
	for _, id := range req.OrderedStepIDs {
		if _, ok := allowed[id]; !ok {
			return nil, ErrInvalidStepReorder
		}
		if _, dup := seen[id]; dup {
			return nil, ErrInvalidStepReorder
		}
		seen[id] = struct{}{}
	}

	err = s.Store.ExecTx(ctx, func(q *db.Queries) error {
		// First pass uses temporary negative positions to avoid unique collisions.
		for i, stepID := range req.OrderedStepIDs {
			if err := q.UpdateHandbookStepSortOrder(ctx, db.UpdateHandbookStepSortOrderParams{
				StepID:    stepID,
				SortOrder: -int32(i + 1),
			}); err != nil {
				return err
			}
		}

		for i, stepID := range req.OrderedStepIDs {
			if err := q.UpdateHandbookStepSortOrder(ctx, db.UpdateHandbookStepSortOrderParams{
				StepID:    stepID,
				SortOrder: int32(i + 1),
			}); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	steps, err := s.Store.ListHandbookStepsByTemplate(ctx, req.TemplateID)
	if err != nil {
		return nil, err
	}
	out := make([]ListStepResponse, 0, len(steps))
	for _, step := range steps {
		out = append(out, ListStepResponse{
			ID:         step.ID,
			SortOrder:  step.SortOrder,
			Kind:       step.Kind,
			Title:      step.Title,
			Body:       step.Body,
			Content:    mustUnmarshalJSON(step.Content),
			IsRequired: step.IsRequired,
		})
	}

	return &ReorderStepsResponse{
		TemplateID: req.TemplateID,
		Steps:      out,
	}, nil
}

func (s *handbookService) AssignTemplateToEmployee(ctx context.Context, actorEmployeeID uuid.UUID, req AssignTemplateToEmployeeRequest) (*AssignTemplateToEmployeeResponse, error) {
	// Quick sanity: ensure template exists.
	tmpl, err := s.Store.GetHandbookTemplateByID(ctx, req.TemplateID)
	if err != nil {
		return nil, err
	}
	if tmpl.Status != db.HandbookTemplateStatusEnumPublished {
		return nil, ErrTemplateNotPublished
	}

	var previousActive *db.GetActiveEmployeeHandbookByEmployeeIDRow
	current, err := s.Store.GetActiveEmployeeHandbookByEmployeeID(ctx, req.EmployeeID)
	if err == nil {
		previousActive = &current
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	var result db.CreateEmployeeHandbookFromTemplateRow
	err = s.Store.ExecTx(ctx, func(q *db.Queries) error {
		if err := q.WaiveActiveEmployeeHandbooksByEmployeeID(ctx, req.EmployeeID); err != nil {
			return err
		}

		row, err := q.CreateEmployeeHandbookFromTemplate(ctx, db.CreateEmployeeHandbookFromTemplateParams{
			EmployeeID:           req.EmployeeID,
			TemplateID:           req.TemplateID,
			AssignedByEmployeeID: uuidPtrOrNil(actorEmployeeID),
		})
		if err != nil {
			return err
		}
		result = row

		if previousActive != nil {
			metadata, err := marshalJSONPayload(map[string]any{
				"source":                  "manual_assignment",
				"replaced_by_handbook_id": result.ID.String(),
				"new_template_id":         result.TemplateID.String(),
			})
			if err != nil {
				return err
			}
			_, err = q.CreateEmployeeHandbookAssignmentHistory(ctx, db.CreateEmployeeHandbookAssignmentHistoryParams{
				EmployeeHandbookID: &previousActive.ID,
				EmployeeID:         previousActive.EmployeeID,
				TemplateID:         previousActive.TemplateID,
				TemplateVersion:    previousActive.TemplateVersion,
				Event:              db.HandbookAssignmentEventEnumReassigned,
				ActorEmployeeID:    uuidPtrOrNil(actorEmployeeID),
				Metadata:           metadata,
			})
			if err != nil {
				return err
			}
		}

		metadata, err := marshalJSONPayload(map[string]any{
			"source":               "manual_assignment",
			"previous_handbook_id": uuidStringOrEmpty(previousActive),
			"previous_template_id": templateStringOrEmpty(previousActive),
		})
		if err != nil {
			return err
		}
		_, err = q.CreateEmployeeHandbookAssignmentHistory(ctx, db.CreateEmployeeHandbookAssignmentHistoryParams{
			EmployeeHandbookID: &result.ID,
			EmployeeID:         result.EmployeeID,
			TemplateID:         result.TemplateID,
			TemplateVersion:    result.TemplateVersion,
			Event:              db.HandbookAssignmentEventEnumAssigned,
			ActorEmployeeID:    uuidPtrOrNil(actorEmployeeID),
			Metadata:           metadata,
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "AssignTemplateToEmployee", "Failed to assign handbook", zap.Error(err))
		return nil, err
	}

	return &AssignTemplateToEmployeeResponse{
		EmployeeHandbookID: result.ID,
		EmployeeID:         result.EmployeeID,
		TemplateID:         result.TemplateID,
		AssignedAt:         result.AssignedAt.Time,
		Status:             string(result.Status),
	}, nil
}

func (s *handbookService) WaiveEmployeeHandbook(ctx context.Context, actorEmployeeID uuid.UUID, req WaiveEmployeeHandbookRequest) (*WaiveEmployeeHandbookResponse, error) {
	hb, err := s.Store.GetEmployeeHandbookByID(ctx, req.EmployeeHandbookID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEmployeeHandbookNotFound
		}
		return nil, err
	}
	if hb.Status != db.HandbookAssignmentStatusEnumNotStarted && hb.Status != db.HandbookAssignmentStatusEnumInProgress {
		return nil, ErrEmployeeHandbookNotActive
	}

	var waived db.EmployeeHandbook
	err = s.Store.ExecTx(ctx, func(q *db.Queries) error {
		row, err := q.WaiveEmployeeHandbookByID(ctx, req.EmployeeHandbookID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrEmployeeHandbookNotActive
			}
			return err
		}
		waived = row

		metadataMap := map[string]any{
			"source": "manual_waive",
		}
		if req.Reason != nil && strings.TrimSpace(*req.Reason) != "" {
			metadataMap["reason"] = strings.TrimSpace(*req.Reason)
		}
		metadata, err := marshalJSONPayload(metadataMap)
		if err != nil {
			return err
		}
		_, err = q.CreateEmployeeHandbookAssignmentHistory(ctx, db.CreateEmployeeHandbookAssignmentHistoryParams{
			EmployeeHandbookID: &waived.ID,
			EmployeeID:         waived.EmployeeID,
			TemplateID:         waived.TemplateID,
			TemplateVersion:    waived.TemplateVersion,
			Event:              db.HandbookAssignmentEventEnumWaived,
			ActorEmployeeID:    uuidPtrOrNil(actorEmployeeID),
			Metadata:           metadata,
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	return &WaiveEmployeeHandbookResponse{
		EmployeeHandbookID: waived.ID,
		EmployeeID:         waived.EmployeeID,
		Status:             string(waived.Status),
		CompletedAt:        util.TsPtr(waived.CompletedAt),
	}, nil
}

func (s *handbookService) ListEmployeeHandbookHistory(ctx context.Context, employeeID uuid.UUID) ([]HandbookAssignmentHistoryEntry, error) {
	rows, err := s.Store.ListEmployeeHandbookAssignmentHistoryByEmployeeID(ctx, db.ListEmployeeHandbookAssignmentHistoryByEmployeeIDParams{
		EmployeeID: employeeID,
		Limit:      50,
		Offset:     0,
	})
	if err != nil {
		return nil, err
	}

	history := make([]HandbookAssignmentHistoryEntry, 0, len(rows))
	for _, row := range rows {
		history = append(history, HandbookAssignmentHistoryEntry{
			ID:                 row.ID,
			EmployeeHandbookID: row.EmployeeHandbookID,
			EmployeeID:         row.EmployeeID,
			TemplateID:         row.TemplateID,
			TemplateVersion:    row.TemplateVersion,
			Event:              string(row.Event),
			ActorEmployeeID:    row.ActorEmployeeID,
			Metadata:           mustUnmarshalJSON(row.Metadata),
			CreatedAt:          row.CreatedAt.Time,
		})
	}

	return history, nil
}

func (s *handbookService) ListEmployeeHandbookAssignments(ctx *gin.Context, req ListEmployeeHandbookAssignmentsRequest) (*pagination.Response[EmployeeHandbookAssignmentSummary], error) {
	params := req.GetParams()
	statusFilter, err := normalizeAssignmentStatusFilter(req.Status)
	if err != nil {
		return nil, err
	}

	rows, err := s.Store.ListEmployeeHandbookAssignments(ctx, db.ListEmployeeHandbookAssignmentsParams{
		Limit:        params.Limit,
		Offset:       params.Offset,
		DepartmentID: req.DepartmentID,
		StatusFilter: statusFilter,
		Search:       req.Search,
	})
	if err != nil {
		return nil, err
	}

	totalCount, err := s.Store.CountEmployeeHandbookAssignments(ctx, db.CountEmployeeHandbookAssignmentsParams{
		DepartmentID: req.DepartmentID,
		StatusFilter: statusFilter,
		Search:       req.Search,
	})
	if err != nil {
		return nil, err
	}

	results := make([]EmployeeHandbookAssignmentSummary, 0, len(rows))
	for _, row := range rows {
		results = append(results, EmployeeHandbookAssignmentSummary{
			EmployeeID:             row.EmployeeID,
			FirstName:              row.FirstName,
			LastName:               row.LastName,
			DepartmentID:           row.EmployeeDepartmentID,
			DepartmentName:         row.DepartmentName,
			EmployeeHandbookID:     row.EmployeeHandbookID,
			TemplateID:             row.HandbookTemplateID,
			TemplateTitle:          row.TemplateTitle,
			TemplateVersion:        row.TemplateVersion,
			HandbookStatus:         interfaceString(row.EmployeeHandbookStatus),
			AssignedAt:             util.TsPtr(row.AssignedAt),
			StartedAt:              util.TsPtr(row.StartedAt),
			CompletedAt:            util.TsPtr(row.CompletedAt),
			DueAt:                  util.TsPtr(row.DueAt),
			RequiredStepsTotal:     row.RequiredStepsTotal,
			RequiredStepsCompleted: row.RequiredStepsCompleted,
		})
	}

	response := pagination.NewResponse(ctx, req.Request, results, totalCount)
	return &response, nil
}

func (s *handbookService) GetEmployeeHandbookDetails(ctx context.Context, handbookID uuid.UUID) (*GetEmployeeHandbookDetailsResponse, error) {
	hb, err := s.Store.GetEmployeeHandbookDetailsByID(ctx, handbookID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEmployeeHandbookNotFound
		}
		return nil, err
	}

	rows, err := s.Store.ListEmployeeHandbookStepsByHandbookID(ctx, handbookID)
	if err != nil {
		return nil, err
	}

	steps := make([]MyHandbookStep, 0, len(rows))
	for _, row := range rows {
		steps = append(steps, MyHandbookStep{
			StepID:      row.StepID,
			SortOrder:   row.SortOrder,
			Kind:        row.Kind,
			Title:       row.Title,
			Body:        row.Body,
			Content:     mustUnmarshalJSON(row.Content),
			IsRequired:  row.IsRequired,
			Status:      string(row.ProgressStatus),
			StartedAt:   util.TsPtr(row.ProgressStartedAt),
			CompletedAt: util.TsPtr(row.ProgressCompletedAt),
			Response:    mustUnmarshalJSON(row.ProgressResponse),
		})
	}

	return &GetEmployeeHandbookDetailsResponse{
		EmployeeHandbookID: hb.ID,
		EmployeeID:         hb.EmployeeID,
		FirstName:          hb.FirstName,
		LastName:           hb.LastName,
		Status:             string(hb.Status),
		AssignedAt:         hb.AssignedAt.Time,
		StartedAt:          util.TsPtr(hb.StartedAt),
		CompletedAt:        util.TsPtr(hb.CompletedAt),
		DueAt:              util.TsPtr(hb.DueAt),
		TemplateID:         hb.TemplateID,
		TemplateTitle:      hb.TemplateTitle,
		TemplateDesc:       hb.TemplateDescription,
		TemplateVersion:    hb.TemplateVersion,
		DepartmentID:       hb.DepartmentID,
		DepartmentName:     hb.DepartmentName,
		Steps:              steps,
	}, nil
}

func (s *handbookService) ListEligibleEmployees(ctx *gin.Context, actorEmployeeID uuid.UUID, req ListEligibleEmployeesRequest) (*pagination.Response[ListEligibleEmployeesResponse], error) {
	params := req.GetParams()
	effectiveDepartmentID := req.DepartmentID
	search := normalizeOptionalSearch(req.Search)

	userID, err := s.Store.GetUserIDByEmployeeID(ctx, actorEmployeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user id: %w", err)
	}

	hasViewAll, err := s.Store.CheckUserPermission(ctx, db.CheckUserPermissionParams{
		UserID: userID,
		Name:   "HANDBOOK.ELIGIBLE_EMPLOYEES.VIEW_ALL",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check handbook eligible employee permission: %w", err)
	}

	if !hasViewAll {
		actorProfile, err := s.Store.GetEmployeeProfileByID(ctx, actorEmployeeID)
		if err != nil {
			return nil, fmt.Errorf("failed to get actor employee profile: %w", err)
		}

		if actorProfile.DepartmentID == nil {
			response := pagination.NewResponse(ctx, req.Request, []ListEligibleEmployeesResponse{}, 0)
			return &response, nil
		}

		effectiveDepartmentID = actorProfile.DepartmentID
	}

	rows, err := s.Store.ListEligibleEmployeesForHandbookAssignment(ctx, db.ListEligibleEmployeesForHandbookAssignmentParams{
		Limit:        params.Limit,
		Offset:       params.Offset,
		DepartmentID: effectiveDepartmentID,
		Search:       search,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list eligible employees: %w", err)
	}

	totalCount, err := s.Store.CountEligibleEmployeesForHandbookAssignment(ctx, db.CountEligibleEmployeesForHandbookAssignmentParams{
		DepartmentID: effectiveDepartmentID,
		Search:       search,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to count eligible employees: %w", err)
	}

	results := make([]ListEligibleEmployeesResponse, 0, len(rows))
	for _, row := range rows {
		results = append(results, ListEligibleEmployeesResponse{
			EmployeeID:     row.EmployeeID,
			FirstName:      row.FirstName,
			LastName:       row.LastName,
			DepartmentID:   row.DepartmentID,
			DepartmentName: row.DepartmentName,
		})
	}

	response := pagination.NewResponse(ctx, req.Request, results, totalCount)
	return &response, nil
}

func uuidPtrOrNil(id uuid.UUID) *uuid.UUID {
	if id == uuid.Nil {
		return nil
	}
	return &id
}

func marshalJSONPayload(v any) ([]byte, error) {
	if v == nil {
		return []byte(`{}`), nil
	}
	return json.Marshal(v)
}

func uuidStringOrEmpty(hb *db.GetActiveEmployeeHandbookByEmployeeIDRow) string {
	if hb == nil {
		return ""
	}
	return hb.ID.String()
}

func templateStringOrEmpty(hb *db.GetActiveEmployeeHandbookByEmployeeIDRow) string {
	if hb == nil {
		return ""
	}
	return hb.TemplateID.String()
}

func interfaceString(v any) string {
	switch typed := v.(type) {
	case string:
		return typed
	case []byte:
		return string(typed)
	default:
		return fmt.Sprint(v)
	}
}

func normalizeAssignmentStatusFilter(status *string) (*string, error) {
	if status == nil {
		return nil, nil
	}
	normalized := strings.TrimSpace(strings.ToLower(*status))
	switch normalized {
	case "":
		return nil, nil
	case "unassigned", "not_started", "in_progress", "completed", "waived":
		return &normalized, nil
	default:
		return nil, ErrInvalidAssignmentStatusFilter
	}
}

func normalizeOptionalSearch(search *string) *string {
	if search == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*search)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func mustUnmarshalJSON(b []byte) any {
	if len(b) == 0 {
		return nil
	}
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		// Keep responses resilient; malformed JSON shouldn't take down an endpoint.
		return nil
	}
	return v
}

func mapFrontendKindToDB(kind string) db.HandbookStepKindEnum {
	switch kind {
	case "rich_text":
		return db.HandbookStepKindEnumContent
	case "link":
		return db.HandbookStepKindEnumLink
	case "quiz":
		return db.HandbookStepKindEnumQuiz
	case "content":
		return db.HandbookStepKindEnumContent
	case "ack":
		return db.HandbookStepKindEnumAck
	default:
		return db.HandbookStepKindEnum(kind)
	}
}

func normalizeAndValidateContent(kind db.HandbookStepKindEnum, raw json.RawMessage) ([]byte, error) {
	switch kind {
	case db.HandbookStepKindEnumLink:
		if len(raw) == 0 || string(raw) == "null" {
			return nil, fmt.Errorf("link content is required")
		}

		var content LinkStepContent
		if err := json.Unmarshal(raw, &content); err != nil {
			return nil, fmt.Errorf("invalid link content: %w", err)
		}

		urlStr := strings.TrimSpace(content.URL)
		if urlStr == "" {
			return nil, fmt.Errorf("link URL is required")
		}

		parsed, err := url.Parse(urlStr)
		if err != nil || parsed.Host == "" {
			return nil, fmt.Errorf("invalid link URL")
		}
		if parsed.Scheme != "http" && parsed.Scheme != "https" {
			return nil, fmt.Errorf("URL must use http or https")
		}

		return json.Marshal(LinkStepContent{URL: urlStr})

	case db.HandbookStepKindEnumQuiz:
		if len(raw) == 0 || string(raw) == "null" {
			return nil, fmt.Errorf("quiz content is required")
		}

		var content QuizStepContent
		if err := json.Unmarshal(raw, &content); err != nil {
			return nil, fmt.Errorf("invalid quiz content: %w", err)
		}

		if strings.TrimSpace(content.Question) == "" {
			return nil, fmt.Errorf("quiz question is required")
		}
		if len(content.Options) < 2 {
			return nil, fmt.Errorf("quiz must have at least 2 options")
		}
		if content.CorrectOptionIndex < 0 || content.CorrectOptionIndex >= len(content.Options) {
			return nil, fmt.Errorf("correct_option_index out of range")
		}

		options := make([]string, 0, len(content.Options))
		for _, opt := range content.Options {
			trimmed := strings.TrimSpace(opt)
			if trimmed == "" {
				return nil, fmt.Errorf("quiz options must be non-empty")
			}
			options = append(options, trimmed)
		}

		return json.Marshal(QuizStepContent{
			Question:           strings.TrimSpace(content.Question),
			Options:            options,
			CorrectOptionIndex: content.CorrectOptionIndex,
		})

	default:
		if len(raw) == 0 || string(raw) == "null" {
			return []byte("{}"), nil
		}
		return []byte("{}"), nil
	}
}

func mapTemplateAPI(t db.HandbookTemplate) HandbookTemplateAPI {
	return HandbookTemplateAPI{
		ID:           t.ID,
		DepartmentID: t.DepartmentID,
		Title:        t.Title,
		Description:  t.Description,
		Version:      t.Version,
		Status:       string(t.Status),
		PublishedAt:  util.TsPtr(t.PublishedAt),
		ArchivedAt:   util.TsPtr(t.ArchivedAt),
		CreatedAt:    t.CreatedAt.Time,
		UpdatedAt:    t.UpdatedAt.Time,
	}
}
