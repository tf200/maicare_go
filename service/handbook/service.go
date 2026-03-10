package handbook

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"
	"maicare_go/service/deps"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
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
		content := normalizeStepContentForOutput(row.Kind, mustUnmarshalJSON(row.Content))
		response := mustUnmarshalJSON(row.ProgressResponse)
		steps = append(steps, MyHandbookStep{
			StepID:      row.StepID,
			SortOrder:   row.SortOrder,
			Kind:        row.Kind,
			Title:       row.Title,
			Body:        row.Body,
			Content:     content,
			IsRequired:  row.IsRequired,
			Status:      string(row.ProgressStatus),
			StartedAt:   tsPtr(row.ProgressStartedAt),
			CompletedAt: tsPtr(row.ProgressCompletedAt),
			Response:    response,
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
		PublishedAt:  tsPtr(t.PublishedAt),
		ArchivedAt:   tsPtr(t.ArchivedAt),
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
		PublishedAt:  tsPtr(t.PublishedAt),
		ArchivedAt:   tsPtr(t.ArchivedAt),
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
	if err := validateStepContentByKind(req.Kind, req.Content); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidStepContent, err)
	}
	normalizedContent, err := normalizeStepContentByKind(req.Kind, req.Content)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidStepContent, err)
	}

	step, err := s.Store.CreateHandbookStep(ctx, db.CreateHandbookStepParams{
		TemplateID: req.TemplateID,
		SortOrder:  req.SortOrder,
		Kind:       req.Kind,
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
		Content:    normalizeStepContentForOutput(step.Kind, mustUnmarshalJSON(step.Content)),
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
			Content:    normalizeStepContentForOutput(step.Kind, mustUnmarshalJSON(step.Content)),
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
		if err := validateStepContentByKind(step.Kind, req.Content); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidStepContent, err)
		}
		normalizedContent, normErr := normalizeStepContentByKind(step.Kind, req.Content)
		if normErr != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidStepContent, normErr)
		}
		contentBytes, marshalErr := json.Marshal(normalizedContent)
		if marshalErr != nil {
			return nil, fmt.Errorf("invalid content payload: %w", marshalErr)
		}
		content = contentBytes
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
		Content:    normalizeStepContentForOutput(updated.Kind, mustUnmarshalJSON(updated.Content)),
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
			Content:    normalizeStepContentForOutput(step.Kind, mustUnmarshalJSON(step.Content)),
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

		_, err = q.CreateEmployeeHandbookAssignmentHistory(ctx, db.CreateEmployeeHandbookAssignmentHistoryParams{
			EmployeeHandbookID: &result.ID,
			EmployeeID:         result.EmployeeID,
			TemplateID:         result.TemplateID,
			TemplateVersion:    result.TemplateVersion,
			Event:              db.HandbookAssignmentEventEnumAssigned,
			ActorEmployeeID:    uuidPtrOrNil(actorEmployeeID),
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

func validateLinkStepContent(content any) error {
	if content == nil {
		return fmt.Errorf("link step content is required")
	}

	obj, ok := content.(map[string]any)
	if !ok {
		return fmt.Errorf("link step content must be an object containing a URL")
	}

	var rawURL any
	for _, key := range []string{"url", "href", "link"} {
		if v, exists := obj[key]; exists {
			rawURL = v
			break
		}
	}
	if rawURL == nil {
		return fmt.Errorf("link step content must include one of: url, href, link")
	}

	urlStr, ok := rawURL.(string)
	if !ok || strings.TrimSpace(urlStr) == "" {
		return fmt.Errorf("link URL must be a non-empty string")
	}
	urlStr = strings.TrimSpace(urlStr)

	parsed, err := url.Parse(urlStr)
	if err != nil || parsed.Host == "" {
		return fmt.Errorf("link URL must be a valid absolute URL")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("link URL must use http or https")
	}

	return nil
}

func validateQuizStepContent(content any) error {
	if content == nil {
		return fmt.Errorf("quiz step content is required")
	}

	obj, ok := content.(map[string]any)
	if !ok {
		return fmt.Errorf("quiz step content must be an object")
	}

	questionRaw, ok := obj["question"]
	if !ok {
		return fmt.Errorf("quiz content must include question")
	}
	question, ok := questionRaw.(string)
	if !ok || strings.TrimSpace(question) == "" {
		return fmt.Errorf("quiz question must be a non-empty string")
	}

	optionsRaw, ok := obj["options"]
	if !ok {
		return fmt.Errorf("quiz content must include options")
	}
	optionsAny, ok := optionsRaw.([]any)
	if !ok || len(optionsAny) < 2 {
		return fmt.Errorf("quiz options must be an array with at least 2 items")
	}

	options := make([]string, 0, len(optionsAny))
	for _, opt := range optionsAny {
		str, ok := opt.(string)
		if !ok || strings.TrimSpace(str) == "" {
			return fmt.Errorf("quiz options must be non-empty strings")
		}
		options = append(options, strings.TrimSpace(str))
	}

	idxRaw, exists := obj["correct_option_index"]
	if !exists {
		return fmt.Errorf("quiz content must include correct_option_index")
	}
	switch v := idxRaw.(type) {
	case float64:
		if v != float64(int(v)) {
			return fmt.Errorf("correct_option_index must be an integer")
		}
		i := int(v)
		if i < 0 || i >= len(options) {
			return fmt.Errorf("correct_option_index out of range")
		}
	default:
		return fmt.Errorf("correct_option_index must be an integer")
	}

	return nil
}

func validateStepContentByKind(kind db.HandbookStepKindEnum, content any) error {
	switch kind {
	case db.HandbookStepKindEnumLink:
		return validateLinkStepContent(content)
	case db.HandbookStepKindEnumQuiz:
		return validateQuizStepContent(content)
	default:
		return nil
	}
}

func normalizeStepContentByKind(kind db.HandbookStepKindEnum, content any) (any, error) {
	switch kind {
	case db.HandbookStepKindEnumLink:
		return normalizeLinkStepContent(content)
	case db.HandbookStepKindEnumQuiz:
		return normalizeQuizStepContent(content)
	default:
		return content, nil
	}
}

func normalizeStepContentForOutput(kind db.HandbookStepKindEnum, content any) any {
	normalized, err := normalizeStepContentByKind(kind, content)
	if err != nil {
		return content
	}
	return normalized
}

func normalizeLinkStepContent(content any) (any, error) {
	if err := validateLinkStepContent(content); err != nil {
		return nil, err
	}
	obj := content.(map[string]any)
	var rawURL any
	for _, key := range []string{"url", "href", "link"} {
		if v, exists := obj[key]; exists {
			rawURL = v
			break
		}
	}
	urlStr := strings.TrimSpace(rawURL.(string))
	return map[string]any{"url": urlStr}, nil
}

func normalizeQuizStepContent(content any) (any, error) {
	if err := validateQuizStepContent(content); err != nil {
		return nil, err
	}
	obj := content.(map[string]any)
	question := strings.TrimSpace(obj["question"].(string))
	optionsAny := obj["options"].([]any)
	options := make([]string, 0, len(optionsAny))
	for _, opt := range optionsAny {
		options = append(options, strings.TrimSpace(opt.(string)))
	}
	idx := int(obj["correct_option_index"].(float64))
	return map[string]any{
		"question":             question,
		"options":              options,
		"correct_option_index": idx,
	}, nil
}

func mapTemplateAPI(t db.HandbookTemplate) HandbookTemplateAPI {
	return HandbookTemplateAPI{
		ID:           t.ID,
		DepartmentID: t.DepartmentID,
		Title:        t.Title,
		Description:  t.Description,
		Version:      t.Version,
		Status:       string(t.Status),
		PublishedAt:  tsPtr(t.PublishedAt),
		ArchivedAt:   tsPtr(t.ArchivedAt),
		CreatedAt:    t.CreatedAt.Time,
		UpdatedAt:    t.UpdatedAt.Time,
	}
}

func tsPtr(ts pgtype.Timestamptz) *time.Time {
	if !ts.Valid {
		return nil
	}
	t := ts.Time
	return &t
}
