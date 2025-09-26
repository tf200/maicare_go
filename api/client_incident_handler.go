package api

import (
	"errors"
	"fmt"
	"maicare_go/async/aclient"
	db "maicare_go/db/sqlc"
	"maicare_go/notification"
	"maicare_go/pagination"
	"maicare_go/pdf"
	"net/http"
	"strconv"
	"time"

	"github.com/goccy/go-json"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateIncidentRequest represents a request to create an incident
type CreateIncidentRequest struct {
	EmployeeID              int64     `json:"employee_id"`
	LocationID              int64     `json:"location_id"`
	ReporterInvolvement     string    `json:"reporter_involvement" binding:"required" enums:"directly_involved,witness,found_afterwards,alarmed"`
	InformWho               []string  `json:"inform_who"`
	IncidentDate            time.Time `json:"incident_date"`
	RuntimeIncident         string    `json:"runtime_incident" binding:"required"`
	IncidentType            string    `json:"incident_type" binding:"required"`
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
	SeverityOfIncident      string    `json:"severity_of_incident" binding:"required" enums:"fatal,serious,less_serious,near_incident"`
	IncidentExplanation     *string   `json:"incident_explanation"`
	RecurrenceRisk          string    `json:"recurrence_risk" binding:"required" enums:"high,very_high,means,very_low"`
	IncidentPreventSteps    *string   `json:"incident_prevent_steps"`
	IncidentTakenMeasures   *string   `json:"incident_taken_measures"`
	Technical               []string  `json:"technical"`
	Organizational          []string  `json:"organizational"`
	MeseWorker              []string  `json:"mese_worker"`
	ClientOptions           []string  `json:"client_options"`
	OtherCause              *string   `json:"other_cause"`
	CauseExplanation        *string   `json:"cause_explanation"`
	PhysicalInjury          string    `json:"physical_injury" binding:"required" enums:"no_injuries,not_noticeable_yet,bruising_swelling,broken_bones,shortness_of_breath,death,other"`
	PhysicalInjuryDesc      *string   `json:"physical_injury_desc"`
	PsychologicalDamage     string    `json:"psychological_damage"`
	PsychologicalDamageDesc *string   `json:"psychological_damage_desc"`
	NeededConsultation      string    `json:"needed_consultation" binding:"required" enums:"no,not_clear,hospitalization,consult_gp"`
	Succession              []string  `json:"succession"`
	SuccessionDesc          *string   `json:"succession_desc"`
	Other                   bool      `json:"other"`
	OtherDesc               *string   `json:"other_desc"`
	AdditionalAppointments  *string   `json:"additional_appointments"`
	EmployeeAbsenteeism     string    `json:"employee_absenteeism"`
	ClientID                int64     `json:"client_id"`
	Emails                  []string  `json:"emails"`
}

// CreateIncidentResponse represents a response for CreateIncidentApi
type CreateIncidentResponse struct {
	ID                      int64     `json:"id"`
	EmployeeID              int64     `json:"employee_id"`
	LocationID              int64     `json:"location_id"`
	ReporterInvolvement     string    `json:"reporter_involvement"`
	InformWho               []string  `json:"inform_who"`
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
	Technical               []string  `json:"technical"`
	Organizational          []string  `json:"organizational"`
	MeseWorker              []string  `json:"mese_worker"`
	ClientOptions           []string  `json:"client_options"`
	OtherCause              *string   `json:"other_cause"`
	CauseExplanation        *string   `json:"cause_explanation"`
	PhysicalInjury          string    `json:"physical_injury"`
	PhysicalInjuryDesc      *string   `json:"physical_injury_desc"`
	PsychologicalDamage     string    `json:"psychological_damage"`
	PsychologicalDamageDesc *string   `json:"psychological_damage_desc"`
	NeededConsultation      string    `json:"needed_consultation"`
	Succession              []string  `json:"succession"`
	SuccessionDesc          *string   `json:"succession_desc"`
	Other                   bool      `json:"other"`
	OtherDesc               *string   `json:"other_desc"`
	AdditionalAppointments  *string   `json:"additional_appointments"`
	EmployeeAbsenteeism     string    `json:"employee_absenteeism"`
	ClientID                int64     `json:"client_id"`
	Emails                  []string  `json:"emails"`
	SoftDelete              bool      `json:"soft_delete"`
	UpdatedAt               time.Time `json:"updated"`
	CreatedAt               time.Time `json:"created"`
}

// CreateIncidentApi creates an incident
// @Summary Create an incident
// @Tags incidents
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body CreateIncidentRequest true "Incident data"
// @Success 201 {object} Response[CreateIncidentResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/incidents [post]
func (server *Server) CreateIncidentApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateIncidentApi", "Invalid client ID", zap.String("client_id", id), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid client ID")))
		return
	}

	var req CreateIncidentRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "CreateIncidentApi", "Invalid request body", zap.String("client_id", id), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	informWhoBytes, err := json.Marshal(req.InformWho)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateIncidentApi", "Failed to marshal informWho", zap.String("client_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process informWho")))
		return
	}

	technicalBytes, err := json.Marshal(req.Technical)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateIncidentApi", "Failed to marshal technical", zap.String("client_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process technical")))
		return
	}

	organizationalBytes, err := json.Marshal(req.Organizational)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateIncidentApi", "Failed to marshal organizational", zap.String("client_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process organizational")))
		return
	}

	meseWorkerBytes, err := json.Marshal(req.MeseWorker)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateIncidentApi", "Failed to marshal meseWorker", zap.String("client_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process meseWorker")))
		return
	}

	clientOptionsBytes, err := json.Marshal(req.ClientOptions)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateIncidentApi", "Failed to marshal clientOptions", zap.String("client_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process clientOptions")))
		return
	}

	successionBytes, err := json.Marshal(req.Succession)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateIncidentApi", "Failed to marshal succession", zap.String("client_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process succession")))
		return
	}

	arg := db.CreateIncidentParams{
		EmployeeID:              req.EmployeeID,
		LocationID:              req.LocationID,
		ReporterInvolvement:     req.ReporterInvolvement,
		InformWho:               informWhoBytes,
		IncidentDate:            pgtype.Date{Time: req.IncidentDate, Valid: true},
		RuntimeIncident:         req.RuntimeIncident,
		IncidentType:            req.IncidentType,
		PassingAway:             req.PassingAway,
		SelfHarm:                req.SelfHarm,
		Violence:                req.Violence,
		FireWaterDamage:         req.FireWaterDamage,
		Accident:                req.Accident,
		ClientAbsence:           req.ClientAbsence,
		Medicines:               req.Medicines,
		Organization:            req.Organization,
		UseProhibitedSubstances: req.UseProhibitedSubstances,
		OtherNotifications:      req.OtherNotifications,
		SeverityOfIncident:      req.SeverityOfIncident,
		IncidentExplanation:     req.IncidentExplanation,
		RecurrenceRisk:          req.RecurrenceRisk,
		IncidentPreventSteps:    req.IncidentPreventSteps,
		IncidentTakenMeasures:   req.IncidentTakenMeasures,
		Technical:               technicalBytes,
		Organizational:          organizationalBytes,
		MeseWorker:              meseWorkerBytes,
		ClientOptions:           clientOptionsBytes,
		OtherCause:              req.OtherCause,
		CauseExplanation:        req.CauseExplanation,
		PhysicalInjury:          req.PhysicalInjury,
		PhysicalInjuryDesc:      req.PhysicalInjuryDesc,
		PsychologicalDamage:     req.PsychologicalDamage,
		PsychologicalDamageDesc: req.PsychologicalDamageDesc,
		NeededConsultation:      req.NeededConsultation,
		Succession:              successionBytes,
		SuccessionDesc:          req.SuccessionDesc,
		Other:                   req.Other,
		OtherDesc:               req.OtherDesc,
		AdditionalAppointments:  req.AdditionalAppointments,
		EmployeeAbsenteeism:     req.EmployeeAbsenteeism,
		ClientID:                clientID,
		Emails:                  req.Emails,
	}

	incident, err := server.store.CreateIncident(ctx, arg)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateIncidentApi", "Failed to create incident", zap.String("client_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create incident")))
		return
	}

	var informWho []string
	err = json.Unmarshal(incident.InformWho, &informWho)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateIncidentApi", "Failed to unmarshal informWho", zap.String("client_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process informWho")))
		return
	}

	var technical []string
	err = json.Unmarshal(incident.Technical, &technical)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateIncidentApi", "Failed to unmarshal technical", zap.String("client_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process technical")))
		return
	}

	var organizational []string
	err = json.Unmarshal(incident.Organizational, &organizational)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateIncidentApi", "Failed to unmarshal organizational", zap.String("client_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process organizational")))
		return
	}

	var meseWorker []string
	err = json.Unmarshal(incident.MeseWorker, &meseWorker)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateIncidentApi", "Failed to unmarshal meseWorker", zap.String("client_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process meseWorker")))
		return
	}

	var clientOptions []string
	err = json.Unmarshal(incident.ClientOptions, &clientOptions)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateIncidentApi", "Failed to unmarshal clientOptions", zap.String("client_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process clientOptions")))
		return
	}

	var succession []string
	err = json.Unmarshal(incident.Succession, &succession)

	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateIncidentApi", "Failed to unmarshal succession", zap.String("client_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process succession")))
		return
	}
	err = server.asynqClient.EnqueueIncident(aclient.IncidentPayload{
		ID:                      incident.ID,
		EmployeeID:              incident.EmployeeID,
		EmployeeFirstName:       "",
		EmployeeLastName:        "",
		LocationID:              incident.LocationID,
		ReporterInvolvement:     incident.ReporterInvolvement,
		InformWho:               informWho,
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
		Technical:               technical,
		Organizational:          organizational,
		MeseWorker:              meseWorker,
		ClientOptions:           clientOptions,
		OtherCause:              incident.OtherCause,
		CauseExplanation:        incident.CauseExplanation,
		PhysicalInjury:          incident.PhysicalInjury,
		PhysicalInjuryDesc:      incident.PhysicalInjuryDesc,
		PsychologicalDamage:     incident.PsychologicalDamage,
		PsychologicalDamageDesc: incident.PsychologicalDamageDesc,
		NeededConsultation:      incident.NeededConsultation,
		Succession:              succession,
		SuccessionDesc:          incident.SuccessionDesc,
		Other:                   incident.Other,
		OtherDesc:               incident.OtherDesc,
		AdditionalAppointments:  incident.AdditionalAppointments,
		EmployeeAbsenteeism:     incident.EmployeeAbsenteeism,
		ClientID:                incident.ClientID,
		LocationName:            "",
		To:                      incident.Emails,
	}, ctx)

	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateIncidentApi", "Failed to start Task", zap.Error(err))
	}

	res := SuccessResponse(CreateIncidentResponse{
		ID:                      incident.ID,
		EmployeeID:              incident.EmployeeID,
		LocationID:              incident.LocationID,
		ReporterInvolvement:     incident.ReporterInvolvement,
		InformWho:               informWho,
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
		Technical:               technical,
		Organizational:          organizational,
		MeseWorker:              meseWorker,
		ClientOptions:           clientOptions,
		OtherCause:              incident.OtherCause,
		CauseExplanation:        incident.CauseExplanation,
		PhysicalInjury:          incident.PhysicalInjury,
		PhysicalInjuryDesc:      incident.PhysicalInjuryDesc,
		PsychologicalDamage:     incident.PsychologicalDamage,
		PsychologicalDamageDesc: incident.PsychologicalDamageDesc,
		NeededConsultation:      incident.NeededConsultation,
		Succession:              succession,
		SuccessionDesc:          incident.SuccessionDesc,
		Other:                   incident.Other,
		OtherDesc:               incident.OtherDesc,
		AdditionalAppointments:  incident.AdditionalAppointments,
		EmployeeAbsenteeism:     incident.EmployeeAbsenteeism,
		ClientID:                incident.ClientID,
		Emails:                  incident.Emails,
		SoftDelete:              incident.SoftDelete,
		UpdatedAt:               incident.UpdatedAt.Time,
		CreatedAt:               incident.CreatedAt.Time,
	}, "Incident created successfully")

	ctx.JSON(http.StatusCreated, res)

}

// ListIncidentsRequest defines the request for listing incidents
type ListIncidentsRequest struct {
	pagination.Request
}

// ListIncidentsResponse defines the response for listing incidents
type ListIncidentsResponse struct {
	ID                      int64     `json:"id"`
	EmployeeID              int64     `json:"employee_id"`
	EmployeeFirstName       string    `json:"employee_first_name"`
	EmployeeLastName        string    `json:"employee_last_name"`
	LocationID              int64     `json:"location_id"`
	ReporterInvolvement     string    `json:"reporter_involvement"`
	InformWho               []string  `json:"inform_who"`
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
	Technical               []string  `json:"technical"`
	Organizational          []string  `json:"organizational"`
	MeseWorker              []string  `json:"mese_worker"`
	ClientOptions           []string  `json:"client_options"`
	OtherCause              *string   `json:"other_cause"`
	CauseExplanation        *string   `json:"cause_explanation"`
	PhysicalInjury          string    `json:"physical_injury"`
	PhysicalInjuryDesc      *string   `json:"physical_injury_desc"`
	PsychologicalDamage     string    `json:"psychological_damage"`
	PsychologicalDamageDesc *string   `json:"psychological_damage_desc"`
	NeededConsultation      string    `json:"needed_consultation"`
	Succession              []string  `json:"succession"`
	SuccessionDesc          *string   `json:"succession_desc"`
	Other                   bool      `json:"other"`
	OtherDesc               *string   `json:"other_desc"`
	AdditionalAppointments  *string   `json:"additional_appointments"`
	EmployeeAbsenteeism     string    `json:"employee_absenteeism"`
	ClientID                int64     `json:"client_id"`
	Emails                  []string  `json:"emails"`
	SoftDelete              bool      `json:"soft_delete"`
	UpdatedAt               time.Time `json:"updated_at"`
	CreatedAt               time.Time `json:"created_at"`
	IsConfirmed             bool      `json:"is_confirmed"`
	EmployeeProfilePicture  *string   `json:"employee_profile_picture"`
	LocationName            string    `json:"location_name"`
}

// ListIncidentsApi lists all incidents
// @Summary List all incidents
// @Tags incidents
// @Produce json
// @Param id path int true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} Response[pagination.Response[ListIncidentsResponse]]
// @Router /clients/{id}/incidents [get]
func (server *Server) ListIncidentsApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "ListIncidentsApi", "Invalid client ID", zap.String("client_id", id), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid client ID")))
		return
	}

	var req ListIncidentsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "ListIncidentsApi", "Failed to bind query parameters", zap.String("client_id", id), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("failed to bind query parameters")))
		return
	}

	params := req.GetParams()

	incidents, err := server.store.ListIncidents(ctx, db.ListIncidentsParams{
		Limit:    params.Limit,
		Offset:   params.Offset,
		ClientID: clientID,
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "ListIncidentsApi", "Failed to list incidents", zap.String("client_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed to list incidents")))
		return
	}

	if len(incidents) == 0 {
		pag := pagination.NewResponse(ctx, req.Request, []ListIncidentsResponse{}, 0)
		res := SuccessResponse(pag, "No incidents found")
		ctx.JSON(http.StatusOK, res)
		return
	}

	totalCount := incidents[0].TotalCount

	incidentList := make([]ListIncidentsResponse, len(incidents))
	for i, incident := range incidents {
		var informWho []string
		err = json.Unmarshal(incident.InformWho, &informWho)
		if err != nil {
			server.logBusinessEvent(LogLevelError, "ListIncidentsApi", "Failed to unmarshal informWho", zap.String("client_id", id), zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process informWho")))
			return
		}

		var technical []string
		err = json.Unmarshal(incident.Technical, &technical)
		if err != nil {
			server.logBusinessEvent(LogLevelError, "ListIncidentsApi", "Failed to unmarshal technical", zap.String("client_id", id), zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process technical")))
			return
		}

		var organizational []string
		err = json.Unmarshal(incident.Organizational, &organizational)
		if err != nil {
			server.logBusinessEvent(LogLevelError, "ListIncidentsApi", "Failed to unmarshal organizational", zap.String("client_id", id), zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process organizational")))
			return
		}

		var meseWorker []string
		err = json.Unmarshal(incident.MeseWorker, &meseWorker)
		if err != nil {
			server.logBusinessEvent(LogLevelError, "ListIncidentsApi", "Failed to unmarshal meseWorker", zap.String("client_id", id), zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process meseWorker")))
			return
		}

		var clientOptions []string
		err = json.Unmarshal(incident.ClientOptions, &clientOptions)
		if err != nil {
			server.logBusinessEvent(LogLevelError, "ListIncidentsApi", "Failed to unmarshal clientOptions", zap.String("client_id", id), zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process clientOptions")))
			return
		}

		var succession []string
		err = json.Unmarshal(incident.Succession, &succession)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		incidentList[i] = ListIncidentsResponse{
			ID:                      incident.ID,
			EmployeeID:              incident.EmployeeID,
			EmployeeFirstName:       incident.EmployeeFirstName,
			EmployeeLastName:        incident.EmployeeLastName,
			LocationID:              incident.LocationID,
			ReporterInvolvement:     incident.ReporterInvolvement,
			InformWho:               informWho,
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
			Technical:               technical,
			Organizational:          organizational,
			MeseWorker:              meseWorker,
			ClientOptions:           clientOptions,
			OtherCause:              incident.OtherCause,
			CauseExplanation:        incident.CauseExplanation,
			PhysicalInjury:          incident.PhysicalInjury,
			PhysicalInjuryDesc:      incident.PhysicalInjuryDesc,
			PsychologicalDamage:     incident.PsychologicalDamage,
			PsychologicalDamageDesc: incident.PsychologicalDamageDesc,
			NeededConsultation:      incident.NeededConsultation,
			Succession:              succession,
			SuccessionDesc:          incident.SuccessionDesc,
			Other:                   incident.Other,
			OtherDesc:               incident.OtherDesc,
			AdditionalAppointments:  incident.AdditionalAppointments,
			EmployeeAbsenteeism:     incident.EmployeeAbsenteeism,
			ClientID:                incident.ClientID,
			Emails:                  incident.Emails,
			SoftDelete:              incident.SoftDelete,
			UpdatedAt:               incident.UpdatedAt.Time,
			CreatedAt:               incident.CreatedAt.Time,
			IsConfirmed:             incident.IsConfirmed,
			EmployeeProfilePicture:  incident.EmployeeProfilePicture,
			LocationName:            incident.LocationName,
		}
	}

	pag := pagination.NewResponse(ctx, req.Request, incidentList, totalCount)

	res := SuccessResponse(pag, "Incidents retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// GetIncidentResponse represents a response for GetIncidentApi
type GetIncidentResponse struct {
	ID                      int64     `json:"id"`
	EmployeeID              int64     `json:"employee_id"`
	EmployeeFirstName       string    `json:"employee_first_name"`
	EmployeeLastName        string    `json:"employee_last_name"`
	LocationID              int64     `json:"location_id"`
	ReporterInvolvement     string    `json:"reporter_involvement"`
	InformWho               []string  `json:"inform_who"`
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
	Technical               []string  `json:"technical"`
	Organizational          []string  `json:"organizational"`
	MeseWorker              []string  `json:"mese_worker"`
	ClientOptions           []string  `json:"client_options"`
	OtherCause              *string   `json:"other_cause"`
	CauseExplanation        *string   `json:"cause_explanation"`
	PhysicalInjury          string    `json:"physical_injury"`
	PhysicalInjuryDesc      *string   `json:"physical_injury_desc"`
	PsychologicalDamage     string    `json:"psychological_damage"`
	PsychologicalDamageDesc *string   `json:"psychological_damage_desc"`
	NeededConsultation      string    `json:"needed_consultation"`
	Succession              []string  `json:"succession"`
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
	LocationName            string    `json:"location_name"`
	Emails                  []string  `json:"emails"`
}

// GetIncidentApi retrieves an incident
// @Summary Retrieve an incident
// @Tags incidents
// @Produce json
// @Param id path int true "Client ID"
// @Param incident_id path int true "Incident ID"
// @Success 200 {object} Response[GetIncidentResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/incidents/{incident_id} [get]
func (server *Server) GetIncidentApi(ctx *gin.Context) {
	id := ctx.Param("incident_id")
	incidentID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetIncidentApi", "Invalid incident ID", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid incident ID")))
		return
	}

	incident, err := server.store.GetIncident(ctx, incidentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			server.logBusinessEvent(LogLevelError, "GetIncidentApi", "Incident not found", zap.String("incident_id", id), zap.Error(err))
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("incident not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed to retrieve incident")))
		return
	}

	var informWho []string
	err = json.Unmarshal(incident.InformWho, &informWho)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetIncidentApi", "Failed to unmarshal informWho", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process informWho")))
		return
	}

	var technical []string
	err = json.Unmarshal(incident.Technical, &technical)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetIncidentApi", "Failed to unmarshal technical", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process technical")))
		return
	}

	var organizational []string
	err = json.Unmarshal(incident.Organizational, &organizational)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetIncidentApi", "Failed to unmarshal organizational", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process organizational")))
		return
	}

	var meseWorker []string
	err = json.Unmarshal(incident.MeseWorker, &meseWorker)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetIncidentApi", "Failed to unmarshal meseWorker", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process meseWorker")))
		return
	}

	var clientOptions []string
	err = json.Unmarshal(incident.ClientOptions, &clientOptions)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetIncidentApi", "Failed to unmarshal clientOptions", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process clientOptions")))
		return
	}

	var succession []string
	err = json.Unmarshal(incident.Succession, &succession)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetIncidentApi", "Failed to unmarshal succession", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process succession")))
		return
	}

	res := SuccessResponse(GetIncidentResponse{
		ID:                      incident.ID,
		EmployeeID:              incident.EmployeeID,
		EmployeeFirstName:       incident.EmployeeFirstName,
		EmployeeLastName:        incident.EmployeeLastName,
		LocationID:              incident.LocationID,
		ReporterInvolvement:     incident.ReporterInvolvement,
		InformWho:               informWho,
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
		Technical:               technical,
		Organizational:          organizational,
		MeseWorker:              meseWorker,
		ClientOptions:           clientOptions,
		OtherCause:              incident.OtherCause,
		CauseExplanation:        incident.CauseExplanation,
		PhysicalInjury:          incident.PhysicalInjury,
		PhysicalInjuryDesc:      incident.PhysicalInjuryDesc,
		PsychologicalDamage:     incident.PsychologicalDamage,
		PsychologicalDamageDesc: incident.PsychologicalDamageDesc,
		NeededConsultation:      incident.NeededConsultation,
		Succession:              succession,
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
		LocationName:            incident.LocationName,
		Emails:                  incident.Emails,
	}, "Incident retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// UpdateIncidentRequest represents a request to update an incident
type UpdateIncidentRequest struct {
	ID                      int64     `json:"id"`
	EmployeeID              *int64    `json:"employee_id"`
	LocationID              *int64    `json:"location_id"`
	ReporterInvolvement     *string   `json:"reporter_involvement"`
	InformWho               []string  `json:"inform_who"`
	IncidentDate            time.Time `json:"incident_date"`
	RuntimeIncident         *string   `json:"runtime_incident"`
	IncidentType            *string   `json:"incident_type"`
	PassingAway             *bool     `json:"passing_away"`
	SelfHarm                *bool     `json:"self_harm"`
	Violence                *bool     `json:"violence"`
	FireWaterDamage         *bool     `json:"fire_water_damage"`
	Accident                *bool     `json:"accident"`
	ClientAbsence           *bool     `json:"client_absence"`
	Medicines               *bool     `json:"medicines"`
	Organization            *bool     `json:"organization"`
	UseProhibitedSubstances *bool     `json:"use_prohibited_substances"`
	OtherNotifications      *bool     `json:"other_notifications"`
	SeverityOfIncident      *string   `json:"severity_of_incident"`
	IncidentExplanation     *string   `json:"incident_explanation"`
	RecurrenceRisk          *string   `json:"recurrence_risk"`
	IncidentPreventSteps    *string   `json:"incident_prevent_steps"`
	IncidentTakenMeasures   *string   `json:"incident_taken_measures"`
	Technical               []string  `json:"technical"`
	Organizational          []string  `json:"organizational"`
	MeseWorker              []string  `json:"mese_worker"`
	ClientOptions           []string  `json:"client_options"`
	OtherCause              *string   `json:"other_cause"`
	CauseExplanation        *string   `json:"cause_explanation"`
	PhysicalInjury          *string   `json:"physical_injury"`
	PhysicalInjuryDesc      *string   `json:"physical_injury_desc"`
	PsychologicalDamage     *string   `json:"psychological_damage"`
	PsychologicalDamageDesc *string   `json:"psychological_damage_desc"`
	NeededConsultation      *string   `json:"needed_consultation"`
	Succession              []string  `json:"succession"`
	SuccessionDesc          *string   `json:"succession_desc"`
	Other                   *bool     `json:"other"`
	OtherDesc               *string   `json:"other_desc"`
	AdditionalAppointments  *string   `json:"additional_appointments"`
	EmployeeAbsenteeism     *string   `json:"employee_absenteeism"`
	Emails                  []string  `json:"emails"`
}

// UpdateIncidentResponse represents a response for UpdateIncidentApi
type UpdateIncidentResponse struct {
	ID                      int64     `json:"id"`
	EmployeeID              int64     `json:"employee_id"`
	LocationID              int64     `json:"location_id"`
	ReporterInvolvement     string    `json:"reporter_involvement"`
	InformWho               []string  `json:"inform_who"`
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
	Technical               []string  `json:"technical"`
	Organizational          []string  `json:"organizational"`
	MeseWorker              []string  `json:"mese_worker"`
	ClientOptions           []string  `json:"client_options"`
	OtherCause              *string   `json:"other_cause"`
	CauseExplanation        *string   `json:"cause_explanation"`
	PhysicalInjury          string    `json:"physical_injury"`
	PhysicalInjuryDesc      *string   `json:"physical_injury_desc"`
	PsychologicalDamage     string    `json:"psychological_damage"`
	PsychologicalDamageDesc *string   `json:"psychological_damage_desc"`
	NeededConsultation      string    `json:"needed_consultation"`
	Succession              []string  `json:"succession"`
	SuccessionDesc          *string   `json:"succession_desc"`
	Other                   bool      `json:"other"`
	OtherDesc               *string   `json:"other_desc"`
	AdditionalAppointments  *string   `json:"additional_appointments"`
	EmployeeAbsenteeism     string    `json:"employee_absenteeism"`
	ClientID                int64     `json:"client_id"`
	SoftDelete              bool      `json:"soft_delete"`
	UpdatedAt               time.Time `json:"updated"`
	CreatedAt               time.Time `json:"created"`
	IsConfirmed             bool      `json:"is_confirmed"`
	Emails                  []string  `json:"emails"`
}

// UpdateIncidentApi updates an incident
// @Summary Update an incident
// @Tags incidents
// @Produce json
// @Param id path int true "Client ID"
// @Param incident_id path int true "Incident ID"
// @Param incident body UpdateIncidentRequest true "Incident"
// @Success 200 {object} Response[UpdateIncidentResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/incidents/{incident_id} [put]
func (server *Server) UpdateIncidentApi(ctx *gin.Context) {
	id := ctx.Param("incident_id")
	incidentID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateIncidentApi", "Invalid incident ID", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid incident ID")))
		return
	}

	var req UpdateIncidentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateIncidentApi", "Failed to bind JSON", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("failed to process incident data")))
		return
	}

	arg := db.UpdateIncidentParams{
		ID:                      incidentID,
		EmployeeID:              req.EmployeeID,
		LocationID:              req.LocationID,
		ReporterInvolvement:     req.ReporterInvolvement,
		IncidentDate:            pgtype.Date{Time: req.IncidentDate, Valid: true},
		RuntimeIncident:         req.RuntimeIncident,
		IncidentType:            req.IncidentType,
		PassingAway:             req.PassingAway,
		SelfHarm:                req.SelfHarm,
		Violence:                req.Violence,
		FireWaterDamage:         req.FireWaterDamage,
		Accident:                req.Accident,
		ClientAbsence:           req.ClientAbsence,
		Medicines:               req.Medicines,
		Organization:            req.Organization,
		UseProhibitedSubstances: req.UseProhibitedSubstances,
		OtherNotifications:      req.OtherNotifications,
		SeverityOfIncident:      req.SeverityOfIncident,
		IncidentExplanation:     req.IncidentExplanation,
		RecurrenceRisk:          req.RecurrenceRisk,
		IncidentPreventSteps:    req.IncidentPreventSteps,
		IncidentTakenMeasures:   req.IncidentTakenMeasures,
		OtherCause:              req.OtherCause,
		CauseExplanation:        req.CauseExplanation,
		PhysicalInjury:          req.PhysicalInjury,
		PhysicalInjuryDesc:      req.PhysicalInjuryDesc,
		PsychologicalDamage:     req.PsychologicalDamage,
		PsychologicalDamageDesc: req.PsychologicalDamageDesc,
		NeededConsultation:      req.NeededConsultation,
		SuccessionDesc:          req.SuccessionDesc,
		Other:                   req.Other,
		OtherDesc:               req.OtherDesc,
		AdditionalAppointments:  req.AdditionalAppointments,
		EmployeeAbsenteeism:     req.EmployeeAbsenteeism,
		Emails:                  req.Emails,
	}

	if req.InformWho != nil {
		informWhoBytes, err := json.Marshal(req.InformWho)
		if err != nil {
			server.logBusinessEvent(LogLevelError, "UpdateIncidentApi", "Failed to marshal InformWho", zap.String("incident_id", id), zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process InformWho")))
			return
		}
		arg.InformWho = informWhoBytes
	}

	if req.Technical != nil {
		technicalBytes, err := json.Marshal(req.Technical)
		if err != nil {
			server.logBusinessEvent(LogLevelError, "UpdateIncidentApi", "Failed to marshal Technical", zap.String("incident_id", id), zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process Technical")))
			return
		}
		arg.Technical = technicalBytes
	}

	if req.Organizational != nil {
		organizationalBytes, err := json.Marshal(req.Organizational)
		if err != nil {
			server.logBusinessEvent(LogLevelError, "UpdateIncidentApi", "Failed to marshal Organizational", zap.String("incident_id", id), zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process Organizational")))
			return
		}
		arg.Organizational = organizationalBytes
	}

	if req.MeseWorker != nil {
		meseWorkerBytes, err := json.Marshal(req.MeseWorker)
		if err != nil {
			server.logBusinessEvent(LogLevelError, "UpdateIncidentApi", "Failed to marshal MeseWorker", zap.String("incident_id", id), zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process MeseWorker")))
			return
		}
		arg.MeseWorker = meseWorkerBytes
	}

	if req.ClientOptions != nil {
		clientOptionsBytes, err := json.Marshal(req.ClientOptions)
		if err != nil {
			server.logBusinessEvent(LogLevelError, "UpdateIncidentApi", "Failed to marshal ClientOptions", zap.String("incident_id", id), zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process ClientOptions")))
			return
		}
		arg.ClientOptions = clientOptionsBytes
	}

	if req.Succession != nil {
		successionBytes, err := json.Marshal(req.Succession)
		if err != nil {
			server.logBusinessEvent(LogLevelError, "UpdateIncidentApi", "Failed to marshal Succession", zap.String("incident_id", id), zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process Succession")))
			return
		}
		arg.Succession = successionBytes
	}

	incident, err := server.store.UpdateIncident(ctx, arg)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateIncidentApi", "Failed to update incident", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process incident data")))
		return
	}

	var informWho []string
	err = json.Unmarshal(incident.InformWho, &informWho)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateIncidentApi", "Failed to unmarshal InformWho", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process InformWho")))
		return
	}

	var technical []string
	err = json.Unmarshal(incident.Technical, &technical)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateIncidentApi", "Failed to unmarshal Technical", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process Technical")))
		return
	}

	var organizational []string
	err = json.Unmarshal(incident.Organizational, &organizational)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateIncidentApi", "Failed to unmarshal Organizational", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process Organizational")))
		return
	}

	var meseWorker []string
	err = json.Unmarshal(incident.MeseWorker, &meseWorker)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateIncidentApi", "Failed to unmarshal MeseWorker", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process MeseWorker")))
		return
	}

	var clientOptions []string
	err = json.Unmarshal(incident.ClientOptions, &clientOptions)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateIncidentApi", "Failed to unmarshal ClientOptions", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process ClientOptions")))
		return
	}

	var succession []string
	err = json.Unmarshal(incident.Succession, &succession)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateIncidentApi", "Failed to unmarshal Succession", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to process Succession")))
		return
	}

	err = server.asynqClient.EnqueueIncident(aclient.IncidentPayload{
		ID:                      incident.ID,
		EmployeeID:              incident.EmployeeID,
		EmployeeFirstName:       "",
		EmployeeLastName:        "",
		LocationID:              incident.LocationID,
		ReporterInvolvement:     incident.ReporterInvolvement,
		InformWho:               informWho,
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
		Technical:               technical,
		Organizational:          organizational,
		MeseWorker:              meseWorker,
		ClientOptions:           clientOptions,
		OtherCause:              incident.OtherCause,
		CauseExplanation:        incident.CauseExplanation,
		PhysicalInjury:          incident.PhysicalInjury,
		PhysicalInjuryDesc:      incident.PhysicalInjuryDesc,
		PsychologicalDamage:     incident.PsychologicalDamage,
		PsychologicalDamageDesc: incident.PsychologicalDamageDesc,
		NeededConsultation:      incident.NeededConsultation,
		Succession:              succession,
		SuccessionDesc:          incident.SuccessionDesc,
		Other:                   incident.Other,
		OtherDesc:               incident.OtherDesc,
		AdditionalAppointments:  incident.AdditionalAppointments,
		EmployeeAbsenteeism:     incident.EmployeeAbsenteeism,
		ClientID:                incident.ClientID,
		LocationName:            "",
		To:                      incident.Emails,
	}, ctx)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateIncidentApi", "Failed to enqueue incident", zap.String("incident_id", id), zap.Error(err))
	}

	res := SuccessResponse(UpdateIncidentResponse{
		ID:                      incident.ID,
		EmployeeID:              incident.EmployeeID,
		LocationID:              incident.LocationID,
		ReporterInvolvement:     incident.ReporterInvolvement,
		InformWho:               informWho,
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
		Technical:               technical,
		Organizational:          organizational,
		MeseWorker:              meseWorker,
		ClientOptions:           clientOptions,
		OtherCause:              incident.OtherCause,
		CauseExplanation:        incident.CauseExplanation,
		PhysicalInjury:          incident.PhysicalInjury,
		PhysicalInjuryDesc:      incident.PhysicalInjuryDesc,
		PsychologicalDamage:     incident.PsychologicalDamage,
		PsychologicalDamageDesc: incident.PsychologicalDamageDesc,
		NeededConsultation:      incident.NeededConsultation,
		Succession:              succession,
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
		Emails:                  incident.Emails,
	}, "Incident updated successfully")

	ctx.JSON(http.StatusOK, res)

}

// DeleteIncidentApi deletes an incident
// @Summary Delete an incident
// @Tags incidents
// @Produce json
// @Param id path int true "Client ID"
// @Param incident_id path int true "Incident ID"
// @Success 200 {object} Response[any]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/incidents/{incident_id} [delete]
func (server *Server) DeleteIncidentApi(ctx *gin.Context) {
	id := ctx.Param("incident_id")
	incidentID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DeleteIncidentApi", "Invalid incident ID", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid incident ID")))
		return
	}

	err = server.store.DeleteIncident(ctx, incidentID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DeleteIncidentApi", "Failed to delete incident", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed to delete incident")))
		return
	}

	res := SuccessResponse([]string{}, "Incident deleted successfully")

	ctx.JSON(http.StatusOK, res)
}

// GenerateIncidentFileResponse represents a response for GenerateIncidentFileApi
type GenerateIncidentFileResponse struct {
	FileUrl *string `json:"file_url"`
	ID      int64   `json:"incident_id"`
}

// GenerateIncidentFileApi generates an incident file
// @Summary Generate an incident file
// @Tags incidents
// @Produce json
// @Param incident_id path int true "Incident ID"
// @Param id path int true "Client ID"
// @Success 200 {object} Response[GenerateIncidentFileResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /clients/{id}/incidents/{incident_id}/file [get]
func (server *Server) GenerateIncidentFileApi(ctx *gin.Context) {
	id := ctx.Param("incident_id")
	incidentID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GenerateIncidentFileApi", "Invalid incident ID", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid incident ID")))
		return
	}

	incident, err := server.store.GetIncident(ctx, incidentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			server.logBusinessEvent(LogLevelError, "GenerateIncidentFileApi", "Incident not found", zap.String("incident_id", id), zap.Error(err))
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("incident not found")))
			return
		}
		server.logBusinessEvent(LogLevelError, "GenerateIncidentFileApi", "Failed to retrieve incident", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed to retrieve incident")))
		return
	}

	var informWho []string
	err = json.Unmarshal(incident.InformWho, &informWho)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GenerateIncidentFileApi", "Failed to unmarshal InformWho", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed to process InformWho")))
		return
	}

	var technical []string
	err = json.Unmarshal(incident.Technical, &technical)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GenerateIncidentFileApi", "Failed to unmarshal Technical", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed to process Technical")))
		return
	}

	var organizational []string
	err = json.Unmarshal(incident.Organizational, &organizational)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GenerateIncidentFileApi", "Failed to unmarshal Organizational", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed to process Organizational")))
		return
	}

	var meseWorker []string
	err = json.Unmarshal(incident.MeseWorker, &meseWorker)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GenerateIncidentFileApi", "Failed to unmarshal MeseWorker", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed to process MeseWorker")))
		return
	}

	var clientOptions []string
	err = json.Unmarshal(incident.ClientOptions, &clientOptions)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GenerateIncidentFileApi", "Failed to unmarshal ClientOptions", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed to process ClientOptions")))
		return
	}

	var succession []string
	err = json.Unmarshal(incident.Succession, &succession)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GenerateIncidentFileApi", "Failed to unmarshal Succession", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed to process Succession")))
		return
	}

	notificationData := notification.NewIncidentReportData{
		ID:                 incident.ID,
		EmployeeID:         incident.EmployeeID,
		EmployeeFirstName:  incident.EmployeeFirstName,
		EmployeeLastName:   incident.EmployeeLastName,
		LocationID:         incident.LocationID,
		LocationName:       incident.LocationName,
		ClientID:           incident.ClientID,
		ClientFirstName:    incident.ClientFirstName,
		ClientLastName:     incident.ClientLastName,
		SeverityOfIncident: incident.SeverityOfIncident,
	}

	receipients, err := server.store.GetAllAdminUsers(ctx)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GenerateIncidentFileApi", "Failed to retrieve admin users", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed to retrieve admin users")))
		return
	}

	var recipientUserIDs []int64
	for _, user := range receipients {
		recipientUserIDs = append(recipientUserIDs, user.ID)
	}

	server.asynqClient.EnqueueNotificationTask(ctx, notification.NotificationPayload{
		RecipientUserIDs: recipientUserIDs,
		Type:             notification.TypeNewClientAssignment,
		Data: notification.NotificationData{
			NewIncidentReport: &notificationData,
		},
		CreatedAt: time.Now(),
	})

	incidentData := pdf.IncidentReportData{
		ID:                      incident.ID,
		EmployeeID:              incident.EmployeeID,
		EmployeeFirstName:       incident.EmployeeFirstName,
		EmployeeLastName:        incident.EmployeeLastName,
		LocationID:              incident.LocationID,
		ReporterInvolvement:     incident.ReporterInvolvement,
		InformWho:               informWho,
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
		Technical:               technical,
		Organizational:          organizational,
		MeseWorker:              meseWorker,
		ClientOptions:           clientOptions,
		OtherCause:              incident.OtherCause,
		CauseExplanation:        incident.CauseExplanation,
		PhysicalInjury:          incident.PhysicalInjury,
		PhysicalInjuryDesc:      incident.PhysicalInjuryDesc,
		PsychologicalDamage:     incident.PsychologicalDamage,
		PsychologicalDamageDesc: incident.PsychologicalDamageDesc,
		NeededConsultation:      incident.NeededConsultation,
		Succession:              succession,
		SuccessionDesc:          incident.SuccessionDesc,
		Other:                   incident.Other,
		OtherDesc:               incident.OtherDesc,
		AdditionalAppointments:  incident.AdditionalAppointments,
		EmployeeAbsenteeism:     incident.EmployeeAbsenteeism,
		ClientID:                incident.ClientID,
		ClientFirstName:         incident.ClientFirstName,
		ClientLastName:          incident.ClientLastName,
		LocationName:            incident.LocationName,
	}

	fileUrl, err := pdf.GenerateAndUploadIncidentPDF(ctx, incidentData, server.b2Client)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GenerateIncidentFileApi", "Failed to generate incident PDF", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed to generate incident file")))
		return
	}

	incidentWithUpdatedFileUrl, err := server.store.UpdateIncidentFileUrl(ctx, db.UpdateIncidentFileUrlParams{
		ID:      incident.ID,
		FileUrl: &fileUrl,
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GenerateIncidentFileApi", "Failed to update incident file URL", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed to update incident file URL")))
		return
	}

	res := SuccessResponse(GenerateIncidentFileResponse{
		FileUrl: incidentWithUpdatedFileUrl,
		ID:      incident.ID,
	}, "Incident file generated successfully")

	ctx.JSON(http.StatusOK, res)
}

// ConfirmIncidentResponse represents a response for ConfirmIncidentApi
type ConfirmIncidentResponse struct {
	ID      int64   `json:"id"`
	FileUrl *string `json:"file_url"`
}

// ConfirmIncidentApi confirms an incident
// @Summary Confirm an incident
// @Tags incidents
// @Produce json
// @Param incident_id path int true "Incident ID"
// @Success 200 {object} Response[ConfirmIncidentResponse]
// @Failure 400,404,500 {object} Response[any]
// @Router /incidents/{incident_id}/confirm [put]
func (server *Server) ConfirmIncidentApi(ctx *gin.Context) {
	id := ctx.Param("incident_id")
	incidentID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "ConfirmIncidentApi", "Invalid incident ID", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid incident ID")))
		return
	}

	incident, err := server.store.ConfirmIncident(ctx, incidentID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "ConfirmIncidentApi", "Failed to confirm incident", zap.String("incident_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed to confirm incident")))
		return
	}

	res := SuccessResponse(ConfirmIncidentResponse{
		ID:      incident.ID,
		FileUrl: incident.FileUrl,
	}, "Incident confirmed successfully")

	// TODO: Send notification to the party responsivle for the client

	ctx.JSON(http.StatusOK, res)
}
