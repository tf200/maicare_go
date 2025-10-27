package care

import (
	"time"

	"maicare_go/pagination"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateClientMaturityMatrixAssessmentRequest represents a request to create a client maturity matrix assessment
type CreateClientCarePlanRequest struct {
	MaturityMatrixID int64 `json:"maturity_matrix_id"`
	InitialLevel     int32 `json:"initial_level"`
	TargetLevel      int32 `json:"target_level"`
}

// CreateClientMaturityMatrixAssessmentResponse represents a response for CreateClientMaturityMatrixAssessmentApi
type CreateClientCarePlanResponse struct {
	ClientID   uuid.UUID `json:"client_id"`
	CarePlanID int64     `json:"care_plan_id"`
}

type Level struct {
	Level       int32  `json:"level"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Details struct {
	TopicName        string  `json:"topic_name"`
	LivingSituation  *string `json:"living_situation"`
	EducationLevel   *string `json:"education_level"`
	Age              int32   `json:"age"`
	LevelDescription []Level `json:"level_description"`
}

// ListClientMaturityMatrixAssessmentsRequest represents a request to list client maturity matrix assessments
type ListClientCarePlansRequest struct {
	pagination.Request
}

// ListClientMaturityMatrixAssessmentsResponse represents a response for ListClientMaturityMatrixAssessmentsApi
type ListClientCarePlansResponse struct {
	CarePlanID   int64       `json:"care_plan_id"`
	ClientID     uuid.UUID   `json:"client_id"`
	StartDate    pgtype.Date `json:"start_date"`
	EndDate      pgtype.Date `json:"end_date"`
	InitialLevel int32       `json:"initial_level"`
	CurrentLevel int32       `json:"current_level"`
	IsActive     bool        `json:"is_active"`
	TopicName    string      `json:"topic_name"`
}

// care_plan represents the response for the GetCarePlanOverview API
type GetCarePlanOverviewResponse struct {
	ID                int64     `json:"id"`
	Domain            string    `json:"domain"`
	CurrentLevel      int32     `json:"current_level"`
	TargetLevel       int32     `json:"target_level"`
	Status            string    `json:"status"`
	GeneratedAt       time.Time `json:"generated_at"`
	AssessmentSummary string    `json:"assessment_summary"`
	RawLlmResponse    string    `json:"raw_llm_response"`
}

// UpdateCarePlanOverviewRequest represents the request body for updating the care plan overview
type UpdateCarePlanOverviewRequest struct {
	AssessmentSummary *string `json:"assessment_summary"`
}

// UpdateCarePlanOverviewResponse represents the response for the UpdateCarePlanOverview API
type UpdateCarePlanOverviewResponse struct {
	CarePlanID        int64  `json:"care_plan_id"`
	AssessmentID      int64  `json:"assessment_id"`
	AssessmentSummary string `json:"assessment_summary"`
}

// =============================== Care Plan Objectives and Actions ===============================

// CreateCarePlanObjectiveRequest represents the request body for creating a care plan objective
type CreateCarePlanObjectiveRequest struct {
	TimeFrame   string `json:"timeframe" binding:"required,oneof=short_term medium_term long_term"`
	GoalTitle   string `json:"goal_title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// CreateCarePlanObjectiveResponse represents the response for the CreateCarePlanObjective API
type CreateCarePlanObjectiveResponse struct {
	ID              int64     `json:"id"`
	CarePlanID      int64     `json:"care_plan_id"`
	Timeframe       string    `json:"timeframe"`
	GoalTitle       string    `json:"goal_title"`
	Description     string    `json:"description"`
	TargetDate      time.Time `json:"target_date"`
	Status          string    `json:"status"`
	CompletionDate  time.Time `json:"completion_date"`
	CompletionNotes *string   `json:"completion_notes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// CarePlanActions represents the actions in a care plan objective
type CarePlanActions struct {
	ActionID          int64  `json:"action_id"`
	SortOrder         int32  `json:"sort_order"`
	ActionDescription string `json:"action_description"`
	IsCompleted       bool   `json:"is_completed"`
	Notes             string `json:"notes"`
}

// CarePlanObjectives represents the objectives in a care plan
type CarePlanObjectives struct {
	ObjectiveID int64             `json:"objective_id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	TimeFrame   string            `json:"timeframe"`
	Status      string            `json:"status"`
	Actions     []CarePlanActions `json:"actions"`
}

// GetCarePlanObjectivesResponse represents the response for the GetCarePlanObjectives API
type GetCarePlanObjectivesResponse struct {
	ShortTermGoals  []CarePlanObjectives `json:"short_term_goals"`
	MediumTermGoals []CarePlanObjectives `json:"medium_term_goals"`
	LongTermGoals   []CarePlanObjectives `json:"long_term_goals"`
}

// UpdateCarePlanObjectiveRequest represents the request body for updating a care plan objective
type UpdateCarePlanObjectiveRequest struct {
	TimeFrame   *string `json:"timeframe" binding:"oneof=short_term medium_term long_term"`
	GoalTitle   *string `json:"goal_title"`
	Description *string `json:"description"`
	Status      *string `json:"status" binding:"oneof=not_started in_progress completed discontinued"`
}

// UpdateCarePlanObjectiveResponse represents the response for the UpdateCarePlanObjective API
type UpdateCarePlanObjectiveResponse struct {
	ObjectiveID int64 `json:"goal_id"`
	CarePlanId  int64 `json:"care_plan_id"`
}

// CreateCarePlanActionsRequest represents the request body for creating a care plan action
type CreateCarePlanActionsRequest struct {
	ActionDescription string `json:"action_description" binding:"required"`
}

// CreateCarePlanActionsResponse represents the response for the CreateCarePlanActions API
type CreateCarePlanActionsResponse struct {
	ActionID          int64  `json:"action_id"`
	ObjectiveID       int64  `json:"objective_id"`
	ActionDescription string `json:"action_description"`
}

// UpdateCarePlanActionsRequest represents the request body for updating a care plan action
type UpdateCarePlanActionsRequest struct {
	ActionDescription *string `json:"action_description"`
}

// UpdateCarePlanActionsResponse represents the response for the UpdateCarePlanActions API
type UpdateCarePlanActionsResponse struct {
	ActionID          int64  `json:"action_id"`
	ObjectiveID       int64  `json:"objective_id"`
	ActionDescription string `json:"action_description"`
}

// CreateCarePlanInterventionRequest represents the request body for creating a care plan intervention
type CreateCarePlanInterventionRequest struct {
	Frequency               string `json:"frequency" binding:"required,oneof=daily weekly monthly"`
	InterventionDescription string `json:"intervention_description" binding:"required"`
}

// CreateCarePlanInterventionResponse represents the response for the CreateCarePlanIntervention API
type CreateCarePlanInterventionResponse struct {
	InterventionID          int64  `json:"intervention_id"`
	CarePlanID              int64  `json:"care_plan_id"`
	Frequency               string `json:"frequency"`
	InterventionDescription string `json:"intervention_description"`
}

type Intervention struct {
	InterventionID          int64  `json:"intervention_id"`
	InterventionDescription string `json:"intervention_description"`
}

// GetCarePlanInterventionsResponse represents the response for the GetCarePlanInterventions API
type GetCarePlanInterventionsResponse struct {
	DailyActivities   []Intervention `json:"daily_activities"`
	WeeklyActivities  []Intervention `json:"weekly_activities"`
	MonthlyActivities []Intervention `json:"monthly_activities"`
}

// UpdateCarePlanInterventionRequest represents the request body for updating a care plan intervention
type UpdateCarePlanInterventionRequest struct {
	Frequency               *string `json:"frequency" binding:"oneof=daily weekly monthly"`
	InterventionDescription *string `json:"intervention_description"`
}

// UpdateCarePlanInterventionApi updates a care plan intervention by its ID
type UpdateCarePlanInterventionResponse struct {
	InterventionID          int64  `json:"intervention_id"`
	CarePlanID              int64  `json:"care_plan_id"`
	Frequency               string `json:"frequency"`
	InterventionDescription string `json:"intervention_description"`
}

// CreateCarePlanSuccessMetricsRequest represents the request body for creating a care plan success metric
type CreateCarePlanSuccessMetricsRequest struct {
	MetricName        string  `json:"metric_name" binding:"required"`
	TargetValue       string  `json:"target_value" binding:"required"`
	MeasurementMethod string  `json:"measurement_method" binding:"required"`
	CurrentValue      *string `json:"current_value"` // Optional, can be nil if not set
}

// CreateCarePlanSuccessMetricsResponse represents the response for the CreateCarePlanSuccessMetrics API
type CreateCarePlanSuccessMetricsResponse struct {
	MetricID          int64   `json:"metric_id"`
	MetricName        string  `json:"metric_name"`
	CurrentValue      *string `json:"current_value"`
	TargetValue       string  `json:"target_value"`
	MeasurementMethod string  `json:"measurement_method"`
}

// GetCarePlanSuccessMetricsResponse represents the response for the GetCarePlanSuccessMetrics AP
type GetCarePlanSuccessMetricsResponse struct {
	MetricID          int64   `json:"metric_id"`
	MetricName        string  `json:"metric_name"`
	CurrentValue      *string `json:"current_value"`
	TargetValue       string  `json:"target_value"`
	MeasurementMethod string  `json:"measurement_method"`
}

// UpdateCarePlanSuccessMetricsRequest represents the request body for updating a care plan success metric
type UpdateCarePlanSuccessMetricsRequest struct {
	MetricName        *string `json:"metric_name"`
	TargetValue       *string `json:"target_value"`
	MeasurementMethod *string `json:"measurement_method"`
	CurrentValue      *string `json:"current_value"` // Optional, can be nil if not set
}

// UpdateCarePlanSuccessMetricsResponse represents the response for the UpdateCarePlanSuccessMetrics API
type UpdateCarePlanSuccessMetricsResponse struct {
	MetricID          int64   `json:"metric_id"`
	MetricName        string  `json:"metric_name"`
	CurrentValue      *string `json:"current_value"`
	TargetValue       string  `json:"target_value"`
	MeasurementMethod string  `json:"measurement_method"`
}

// CreateCarePlanRisksRequest represents the request body for creating a care plan risk
type CreateCarePlanRisksRequest struct {
	RiskDescription    string  `json:"risk_description" binding:"required"`
	MitigationStrategy string  `json:"mitigation_strategy" binding:"required"`
	RiskLevel          *string `json:"risk_level"`
}

// CreateCarePlanRisksResponse represents the response for the CreateCarePlanRisks API
type CreateCarePlanRisksResponse struct {
	RiskID             int64   `json:"risk_id"`
	RiskDescription    string  `json:"risk_description"`
	MitigationStrategy string  `json:"mitigation_strategy"`
	RiskLevel          *string `json:"risk_level"`
}

// GetCarePlanRisksResponse represents the response for the GetCarePlanRisks API
type GetCarePlanRisksResponse struct {
	RiskID             int64   `json:"risk_id"`
	RiskDescription    string  `json:"risk_description"`
	MitigationStrategy string  `json:"mitigation_strategy"`
	RiskLevel          *string `json:"risk_level"`
}

// UpdateCarePlanRisksRequest represents the request body for updating a care plan risk
type UpdateCarePlanRisksRequest struct {
	RiskDescription    *string `json:"risk_description"`
	MitigationStrategy *string `json:"mitigation_strategy"`
	RiskLevel          *string `json:"risk_level"`
}

// UpdateCarePlanRisksResponse represents the response for the UpdateCarePlanRisks API
type UpdateCarePlanRisksResponse struct {
	RiskID             int64   `json:"risk_id"`
	RiskDescription    string  `json:"risk_description"`
	MitigationStrategy string  `json:"mitigation_strategy"`
	RiskLevel          *string `json:"risk_level"`
}

// CreateCarePlanSupportNetworkRequest represents the request body for creating a care plan support network
type CreateCarePlanSupportNetworkRequest struct {
	RoleTitle                 string `json:"role_title" binding:"required"`
	ResponsibilityDescription string `json:"responsibility_description"`
}

// CreateCarePlanSupportNetworkResponse represents the response for the CreateCarePlanSupportNetwork API
type CreateCarePlanSupportNetworkResponse struct {
	SupportNetworkID          int64  `json:"support_network_id"`
	RoleTitle                 string `json:"role_title"`
	ResponsibilityDescription string `json:"responsibility_description"`
}

// GetCarePlanSupportNetworkResponse represents the response for the GetCarePlanSupportNetwork API
type GetCarePlanSupportNetworkResponse struct {
	SupportNetworkID          int64   `json:"support_network_id"`
	RoleTitle                 string  `json:"role_title"`
	ResponsibilityDescription *string `json:"responsibility_description"`
}

// UpdateCarePlanSupportNetworkRequest represents the request body for updating a care plan support network
type UpdateCarePlanSupportNetworkRequest struct {
	RoleTitle                 *string `json:"role_title"`
	ResponsibilityDescription *string `json:"responsibility_description"`
}

// UpdateCarePlanSupportNetworkResponse represents the response for the UpdateCarePlanSupportNetwork API
type UpdateCarePlanSupportNetworkResponse struct {
	SupportNetworkID          int64  `json:"support_network_id"`
	RoleTitle                 string `json:"role_title"`
	ResponsibilityDescription string `json:"responsibility_description"`
}

// CreateCarePlanResourcesRequest represents the request body for creating a care plan resource
type CreateCarePlanResourcesRequest struct {
	ResourceDescription string     `json:"resource_description" binding:"required"`
	IsObtained          *bool      `json:"is_obtained"`
	ObtainedDate        *time.Time `json:"obtained_date"`
}

// CreateCarePlanResourcesResponse represents the response for the CreateCarePlanResources API
type CreateCarePlanResourcesResponse struct {
	ID                  int64      `json:"id"`
	ResourceDescription string     `json:"resource_description"`
	IsObtained          bool       `json:"is_obtained"`
	ObtainedDate        *time.Time `json:"obtained_date"`
}

// GetCarePlanResourcesResponse represents the response for the GetCarePlanResources API
type GetCarePlanResourcesResponse struct {
	ID                  int64      `json:"id"`
	ResourceDescription string     `json:"resource_description"`
	IsObtained          bool       `json:"is_obtained"`
	ObtainedDate        *time.Time `json:"obtained_date"`
}

// UpdateCarePlanResourcesRequest represents the request body for updating a care plan resource
type UpdateCarePlanResourcesRequest struct {
	ResourceDescription *string   `json:"resource_description"`
	IsObtained          *bool     `json:"is_obtained"`
	ObtainedDate        time.Time `json:"obtained_date"`
}

// UpdateCarePlanResourcesResponse represents the response for the UpdateCarePlanResources API
type UpdateCarePlanResourcesResponse struct {
	ID                  int64     `json:"id"`
	ResourceDescription string    `json:"resource_description"`
	IsObtained          bool      `json:"is_obtained"`
	ObtainedDate        time.Time `json:"obtained_date"`
}

// CreateCarePlanReportRequest represents the request body for creating a care plan report
type CreateCarePlanReportRequest struct {
	ReportType    string `json:"report_type" binding:"required" oneof:"progress concern achievement modification"`
	ReportContent string `json:"report_content" binding:"required"`
	IsCritical    bool   `json:"is_critical"`
}

// CreateCarePlanReportResponse represents the response for the CreateCarePlanReport API
type CreateCarePlanReportResponse struct {
	ID            int64     `json:"id"`
	CarePlanID    int64     `json:"care_plan_id"`
	ReportType    string    `json:"report_type"`
	ReportContent string    `json:"report_content"`
	IsCritical    bool      `json:"is_critical"`
	CreatedAt     time.Time `json:"created_at"`
}

// ListCarePlanReportsRequest represents the request body for listing care plan reports
type ListCarePlanReportsRequest struct {
	pagination.Request
}

// CarePlanReportsResponse represents the response for the ListCarePlanReports API
type ListCarePlanReportsResponse struct {
	ID                 int64     `json:"id"`
	CarePlanID         int64     `json:"care_plan_id"`
	ReportType         string    `json:"report_type"`
	ReportContent      string    `json:"report_content"`
	CreatedByFirstName string    `json:"created_by_first_name"`
	CreatedByLastName  string    `json:"created_by_last_name"`
	IsCritical         bool      `json:"is_critical"`
	CreatedAt          time.Time `json:"created_at"`
}

// UpdateCarePlanReportRequest represents the request body for updating a care plan report
type UpdateCarePlanReportRequest struct {
	ReportType    *string `json:"report_type"`
	ReportContent *string `json:"report_content"`
	IsCritical    *bool   `json:"is_critical"`
}

// UpdateCarePlanReportResponse represents the response for the UpdateCarePlanReport API
type UpdateCarePlanReportResponse struct {
	ID            int64     `json:"id"`
	CarePlanID    int64     `json:"care_plan_id"`
	ReportType    string    `json:"report_type"`
	ReportContent string    `json:"report_content"`
	IsCritical    bool      `json:"is_critical"`
	CreatedAt     time.Time `json:"created_at"`
}
