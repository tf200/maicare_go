package tasks

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/goccy/go-json"

	"github.com/hibiken/asynq"
)

const (
	TypeEmailDelivery = "email:deliver"
	TypeIncident      = "incident:process"
)

type EmailDeliveryPayload struct {
	To           string `json:"to"`
	UserEmail    string `json:"user_email"`
	UserPassword string `json:"user_password"`
}

func (c *AsynqClient) EnqueueEmailDelivery(
	payload EmailDeliveryPayload,
	ctx context.Context,
	opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("json.Marshal failed: %v", err)
	}
	task := asynq.NewTask(TypeEmailDelivery, jsonPayload)
	info, err := c.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		return fmt.Errorf("client.EnqueueContext failed: %v", err)
	}
	log.Printf("task enqueued: id=%s queue=%s", info.ID, info.Queue)
	return nil

}

type IncidentPayload struct {
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
	LocationName            string    `json:"location_name"`
	To                      []string  `json:"to"`
}

func (c *AsynqClient) EnqueueIncident(
	payload IncidentPayload,
	ctx context.Context,
	opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("json.Marshal failed: %v", err)
	}
	task := asynq.NewTask(TypeIncident, jsonPayload)
	info, err := c.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		return fmt.Errorf("client.EnqueueContext failed: %v", err)
	}
	log.Printf("task enqueued: id=%s queue=%s", info.ID, info.Queue)
	return nil
}
