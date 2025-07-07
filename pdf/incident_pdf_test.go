package pdf

import (
	"context"
	"log"
	"maicare_go/bucket"
	"maicare_go/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGenerateIncidentPDF(t *testing.T) {

	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatalf("Could not load conf %v", err)
	}

	testb2Client, err := bucket.NewB2Client(config)
	if err != nil {
		log.Fatal("cannot create b2 client:", err)
	}

	mockIncident := IncidentReportData{
		ID:                1,
		EmployeeID:        100,
		EmployeeFirstName: "John",
		EmployeeLastName:  "Doe",
		LocationID:        200,
		LocationName:      "Main Building - Floor 3",

		// Basic incident information
		ReporterInvolvement: "directly_involved",
		InformWho:           []string{"supervisor", "medical_staff"},
		IncidentDate:        time.Date(2024, 2, 9, 14, 30, 0, 0, time.UTC),
		RuntimeIncident:     "14:30",
		IncidentType:        "workplace_accident",

		// Incident categories
		PassingAway:             false,
		SelfHarm:                false,
		Violence:                false,
		FireWaterDamage:         true,
		Accident:                true,
		ClientAbsence:           false,
		Medicines:               false,
		Organization:            false,
		UseProhibitedSubstances: false,
		OtherNotifications:      false,

		// Severity and risk assessment
		SeverityOfIncident:  "serious",
		IncidentExplanation: util.StringPtr("Water leak from ceiling caused slippery floor condition. Employee slipped while attending to routine duties."),
		RecurrenceRisk:      "means",

		// Prevention and measures
		IncidentPreventSteps:  util.StringPtr("1. Regular ceiling maintenance checks\n2. Installation of water detection systems\n3. Implementation of emergency response protocols"),
		IncidentTakenMeasures: util.StringPtr("1. Area immediately cordoned off\n2. Emergency maintenance called\n3. First aid administered\n4. Incident report filed"),

		// Contributing factors
		Technical:      []string{"faulty_plumbing", "inadequate_drainage"},
		Organizational: []string{"maintenance_schedule", "emergency_response"},
		MeseWorker:     []string{"following_protocol", "proper_reporting"},
		ClientOptions:  []string{"none_applicable"},

		// Cause analysis
		OtherCause:       util.StringPtr("Unexpected pipe burst"),
		CauseExplanation: util.StringPtr("Recent temperature fluctuations may have contributed to pipe stress"),

		// Injury assessment
		PhysicalInjury:          "bruising_swelling",
		PhysicalInjuryDesc:      util.StringPtr("Minor bruising on right hip and elbow"),
		PsychologicalDamage:     "unrest",
		PsychologicalDamageDesc: util.StringPtr("Employee expressed anxiety about returning to work area"),
		NeededConsultation:      "consult_gp",

		// Follow-up
		Succession:     []string{"medical_checkup", "facility_inspection"},
		SuccessionDesc: util.StringPtr("Follow-up medical appointment scheduled for next week"),

		// Additional information
		Other:                  false,
		OtherDesc:              nil,
		AdditionalAppointments: util.StringPtr("GP appointment scheduled for 2024-02-16"),
		EmployeeAbsenteeism:    "1_week",

		// References
		ClientID:        300,
		ClientFirstName: "Jane",
		ClientLastName:  "Smith",
	}

	pdfBytes, err := GenerateAndUploadIncidentPDF(context.Background(), mockIncident, testb2Client)
	require.NoError(t, err)
	require.NotEmpty(t, pdfBytes)

}
