package handbook

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"
	"maicare_go/service/deps"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type handbookService struct {
	*deps.ServiceDependencies
}

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
			StartedAt:   tsPtr(row.ProgressStartedAt),
			CompletedAt: tsPtr(row.ProgressCompletedAt),
			Response:    mustUnmarshalJSON(row.ProgressResponse),
		})
	}

	res := &GetMyActiveHandbookResponse{
		HandbookID:      hb.ID,
		Status:          string(hb.Status),
		AssignedAt:      hb.AssignedAt.Time,
		StartedAt:       tsPtr(hb.StartedAt),
		CompletedAt:     tsPtr(hb.CompletedAt),
		DueAt:           tsPtr(hb.DueAt),
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
	updated, err := s.Store.MarkEmployeeHandbookStarted(ctx, hb.ID)
	if err != nil {
		return nil, err
	}
	return &StartMyHandbookResponse{
		HandbookID: updated.ID,
		Status:     string(updated.Status),
		StartedAt:  tsPtr(updated.StartedAt),
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

	progress, err := s.Store.CompleteEmployeeHandbookStep(ctx, db.CompleteEmployeeHandbookStepParams{
		Response:           respBytes,
		EmployeeHandbookID: hb.ID,
		StepID:             stepID,
	})
	if err != nil {
		return nil, err
	}

	remaining, err := s.Store.CountRemainingRequiredHandbookSteps(ctx, hb.ID)
	if err != nil {
		return nil, err
	}

	handbookStatus := string(hb.Status)
	if remaining == 0 {
		completed, err := s.Store.MarkEmployeeHandbookCompleted(ctx, hb.ID)
		if err != nil {
			return nil, err
		}
		handbookStatus = string(completed.Status)
	}

	return &CompleteMyHandbookStepResponse{
		HandbookID:     progress.EmployeeHandbookID,
		StepID:         progress.StepID,
		StepStatus:     string(progress.Status),
		CompletedAt:    progress.CompletedAt.Time,
		HandbookStatus: handbookStatus,
	}, nil
}

func (s *handbookService) CreateTemplateForDepartment(ctx context.Context, actorEmployeeID uuid.UUID, req CreateTemplateForDepartmentRequest) (*CreateTemplateForDepartmentResponse, error) {
	t, err := s.Store.CreateHandbookTemplateForDepartment(ctx, db.CreateHandbookTemplateForDepartmentParams{
		DepartmentID:        req.DepartmentID,
		Title:               req.Title,
		Description:         req.Description,
		CreatedByEmployeeID: uuidPtrOrNil(actorEmployeeID),
	})
	if err != nil {
		return nil, err
	}
	return &CreateTemplateForDepartmentResponse{
		ID:           t.ID,
		DepartmentID: t.DepartmentID,
		Title:        t.Title,
		Description:  t.Description,
		Version:      t.Version,
		IsActive:     t.IsActive,
		CreatedAt:    t.CreatedAt.Time,
	}, nil
}

func (s *handbookService) ListTemplatesByDepartment(ctx context.Context, departmentID uuid.UUID) (*pagination.Response[ListTemplateResponse], error) {
	templates, err := s.Store.ListHandbookTemplatesByDepartment(ctx, departmentID)
	if err != nil {
		return nil, err
	}
	out := make([]ListTemplateResponse, 0, len(templates))
	for _, t := range templates {
		out = append(out, ListTemplateResponse{
			ID:           t.ID,
			DepartmentID: t.DepartmentID,
			Title:        t.Title,
			Description:  t.Description,
			Version:      t.Version,
			IsActive:     t.IsActive,
		})
	}

	// No real pagination yet (templates are small). Use a stable wrapper for the frontend.
	resp := &pagination.Response[ListTemplateResponse]{
		Next:     nil,
		Previous: nil,
		Count:    int64(len(out)),
		PageSize: int32(len(out)),
		Results:  out,
	}
	return resp, nil
}

func (s *handbookService) CreateStep(ctx context.Context, req CreateStepRequest) (*CreateStepResponse, error) {
	step, err := s.Store.CreateHandbookStep(ctx, db.CreateHandbookStepParams{
		TemplateID: req.TemplateID,
		SortOrder:  req.SortOrder,
		Kind:       req.Kind,
		Title:      req.Title,
		Body:       req.Body,
		Content:    req.Content,
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

func (s *handbookService) AssignTemplateToEmployee(ctx context.Context, actorEmployeeID uuid.UUID, req AssignTemplateToEmployeeRequest) (*AssignTemplateToEmployeeResponse, error) {
	// Quick sanity: ensure template exists.
	_, err := s.Store.GetHandbookTemplateByID(ctx, req.TemplateID)
	if err != nil {
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

func uuidPtrOrNil(id uuid.UUID) *uuid.UUID {
	if id == uuid.Nil {
		return nil
	}
	return &id
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

func tsPtr(ts pgtype.Timestamptz) *time.Time {
	if !ts.Valid {
		return nil
	}
	t := ts.Time
	return &t
}
