package api

import (
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ListAllIncidentsRequest represents the request body for listing all incidents
type ListAllIncidentsRequest struct {
	pagination.Request
	IsConfirmed bool `form:"is_confirmed" json:"is_confirmed"`
}

// ListAllIncidentsResponse represents the response body for listing all incidents
type ListAllIncidentsResponse struct {
	ID                      int64     `json:"id"`
	EmployeeID              int64     `json:"employee_id"`
	LocationID              int64     `json:"location_id"`
	ReporterInvolvement     string    `json:"reporter_involvement"`
	IncidentDate            time.Time `json:"incident_date"`
	RuntimeIncident         string    `json:"runtime_incident"`
	IncidentType            string    `json:"incident_type"`
	PassingAway             bool      `json:"passing_away"`
	SelfHarm                bool      `json:"self_harm"`
	Violence                bool      `json:"violence"`
	FireWaterDamage         bool      `json:"fire_water_damage"`
	Accident                bool      `json:"accident"`
	ClientAbsence           bool      `json:"client_absence"`
	Medicines               bool      `json:"medicines"`
	Organization            bool      `json:"organization"`
	UseProhibitedSubstances bool      `json:"use_prohibited_substances"`
	OtherNotifications      bool      `json:"other_notifications"`
	SeverityOfIncident      string    `json:"severity_of_incident"`
	IncidentExplanation     *string   `json:"incident_explanation"`
	RecurrenceRisk          string    `json:"recurrence_risk"`
	IncidentPreventSteps    *string   `json:"incident_prevent_steps"`
	IncidentTakenMeasures   *string   `json:"incident_taken_measures"`
	OtherCause              *string   `json:"other_cause"`
	CauseExplanation        *string   `json:"cause_explanation"`
	PhysicalInjury          string    `json:"physical_injury"`
	PhysicalInjuryDesc      *string   `json:"physical_injury_desc"`
	PsychologicalDamage     string    `json:"psychological_damage"`
	PsychologicalDamageDesc *string   `json:"psychological_damage_desc"`
	NeededConsultation      string    `json:"needed_consultation"`
	SuccessionDesc          *string   `json:"succession_desc"`
	Other                   bool      `json:"other"`
	OtherDesc               *string   `json:"other_desc"`
	AdditionalAppointments  *string   `json:"additional_appointments"`
	EmployeeAbsenteeism     string    `json:"employee_absenteeism"`
	ClientID                int64     `json:"client_id"`
	SoftDelete              bool      `json:"soft_delete"`
	UpdatedAt               time.Time `json:"updated_at"`
	CreatedAt               time.Time `json:"created_at"`
	IsConfirmed             bool      `json:"is_confirmed"`
	FileUrl                 *string   `json:"file_url"`
	Emails                  []string  `json:"emails"`
	ClientFirstName         string    `json:"client_first_name"`
	ClientLastName          string    `json:"client_last_name"`
	EmployeeFirstName       string    `json:"employee_first_name"`
	EmployeeLastName        string    `json:"employee_last_name"`
}

// ListAllIncidentsApi handles the API request to list all incidents
// @Summary List all incidents
// @Description List all incidents with pagination and filtering options
// @Tags incidents
// @Produce json
// @Param is_confirmed query bool false "Filter by confirmation status"
// @Param page query int false "Page number"
// @Param page_size query int false "Number of items per page"
// @Success 200 {object} Response[pagination.Response[ListAllIncidentsResponse]]
// @Failure 400 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /incidents [get]
func (serevr *Server) ListAllIncidentsApi(ctx *gin.Context) {
	var req ListAllIncidentsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := req.GetParams()

	arg := db.ListAllIncidentsParams{
		Limit:       params.Limit,
		Offset:      params.Offset,
		IsConfirmed: req.IsConfirmed,
	}

	incidents, err := serevr.store.ListAllIncidents(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	count, err := serevr.store.CountAllIncidents(ctx, req.IsConfirmed)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	incidentList := make([]ListAllIncidentsResponse, len(incidents))
	for i, incident := range incidents {
		incidentList[i] = ListAllIncidentsResponse{
			ID:                      incident.ID,
			EmployeeID:              incident.EmployeeID,
			LocationID:              incident.LocationID,
			ReporterInvolvement:     incident.ReporterInvolvement,
			IncidentDate:            incident.IncidentDate.Time,
			RuntimeIncident:         incident.RuntimeIncident,
			IncidentType:            incident.IncidentType,
			PassingAway:             incident.PassingAway,
			SelfHarm:                incident.SelfHarm,
			Violence:                incident.Violence,
			FireWaterDamage:         incident.FireWaterDamage,
			Accident:                incident.Accident,
			ClientAbsence:           incident.ClientAbsence,
			Medicines:               incident.Medicines,
			Organization:            incident.Organization,
			UseProhibitedSubstances: incident.UseProhibitedSubstances,
			OtherNotifications:      incident.OtherNotifications,
			SeverityOfIncident:      incident.SeverityOfIncident,
			IncidentExplanation:     incident.IncidentExplanation,
			RecurrenceRisk:          incident.RecurrenceRisk,
			IncidentPreventSteps:    incident.IncidentPreventSteps,
			IncidentTakenMeasures:   incident.IncidentTakenMeasures,
			OtherCause:              incident.OtherCause,
			CauseExplanation:        incident.CauseExplanation,
			PhysicalInjury:          incident.PhysicalInjury,
			PhysicalInjuryDesc:      incident.PhysicalInjuryDesc,
			PsychologicalDamage:     incident.PsychologicalDamage,
			PsychologicalDamageDesc: incident.PsychologicalDamageDesc,
			NeededConsultation:      incident.NeededConsultation,
			SuccessionDesc:          incident.SuccessionDesc,
			Other:                   incident.Other,
			OtherDesc:               incident.OtherDesc,
			AdditionalAppointments:  incident.AdditionalAppointments,
			EmployeeAbsenteeism:     incident.EmployeeAbsenteeism,
			ClientID:                incident.ClientID,
			SoftDelete:              incident.SoftDelete,
			UpdatedAt:               incident.UpdatedAt.Time,
			CreatedAt:               incident.CreatedAt.Time,
			IsConfirmed:             incident.IsConfirmed,
			FileUrl:                 incident.FileUrl,
			Emails:                  incident.Emails,
			ClientFirstName:         incident.ClientFirstName,
			ClientLastName:          incident.ClientLastName,
			EmployeeFirstName:       incident.EmployeeFirstName,
			EmployeeLastName:        incident.EmployeeLastName,
		}
	}

	pag := pagination.NewResponse(ctx, req.Request, incidentList, count)
	res := SuccessResponse(pag, "Incidents fetched successfully")
	ctx.JSON(http.StatusOK, res)
}
