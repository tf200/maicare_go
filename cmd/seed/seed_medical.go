package main

import (
	"context"
	"github.com/goccy/go-json"
	"fmt"
	"time"

	db "maicare_go/db/sqlc"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type diagnosisSeedPreset struct {
	codeSystem  string
	code        string
	title       string
	description string
}

type medicationSeedPreset struct {
	name       string
	unit       string
	route      string
	minDose    float64
	maxDose    float64
	critical   bool
	indication string
}

var diagnosisSeedPresets = []diagnosisSeedPreset{
	{codeSystem: "DSM-5", code: "F41.1", title: "Generalized anxiety disorder", description: "Persistent anxiety affecting daily functioning."},
	{codeSystem: "DSM-5", code: "F33.1", title: "Major depressive disorder, recurrent", description: "Episodes of depressed mood with impaired functioning."},
	{codeSystem: "ICD-10", code: "F90.0", title: "ADHD", description: "Attention and concentration difficulties with impulsivity."},
	{codeSystem: "ICD-10", code: "F43.1", title: "Post-traumatic stress disorder", description: "Trauma-related stress symptoms and hypervigilance."},
	{codeSystem: "ICD-10", code: "F84.0", title: "Autism spectrum condition", description: "Neurodevelopmental condition with social communication differences."},
	{codeSystem: "ICD-10", code: "F91.3", title: "Oppositional defiant disorder", description: "Pattern of irritable mood and defiant behavior."},
	{codeSystem: "ICD-10", code: "G47.0", title: "Insomnia", description: "Difficulty initiating or maintaining sleep."},
}

var medicationSeedPresets = []medicationSeedPreset{
	{name: "Sertraline", unit: "mg", route: "oral", minDose: 25, maxDose: 150, critical: false, indication: "mood and anxiety stabilization"},
	{name: "Methylphenidate", unit: "mg", route: "oral", minDose: 5, maxDose: 40, critical: false, indication: "attention and impulse control"},
	{name: "Risperidone", unit: "mg", route: "oral", minDose: 0.5, maxDose: 4, critical: true, indication: "behavioral dysregulation"},
	{name: "Melatonin", unit: "mg", route: "oral", minDose: 1, maxDose: 5, critical: false, indication: "sleep onset support"},
	{name: "Diazepam", unit: "mg", route: "oral", minDose: 2, maxDose: 10, critical: true, indication: "acute severe anxiety or agitation"},
	{name: "Ibuprofen", unit: "mg", route: "oral", minDose: 200, maxDose: 600, critical: false, indication: "pain or inflammation"},
}

func (s *Seeder) SeedMedicalForClients(ctx context.Context, maxDiagnosesPerClient int, maxMedicationOrdersPerClient int) error {
	if maxDiagnosesPerClient <= 0 && maxMedicationOrdersPerClient <= 0 {
		return nil
	}
	if len(s.data.ClientIDs) == 0 {
		return fmt.Errorf("no clients available; seed waiting-list or in-care clients first")
	}

	for i, clientID := range s.data.ClientIDs {
		if (i+1)%10 == 0 || i == 0 || i+1 == len(s.data.ClientIDs) {
			fmt.Printf("[seed] client medical records: %d/%d\n", i+1, len(s.data.ClientIDs))
		}

		if err := s.seedMedicalForClient(ctx, clientID, maxDiagnosesPerClient, maxMedicationOrdersPerClient); err != nil {
			return fmt.Errorf("seed medical records for client %s: %w", clientID, err)
		}
	}

	return nil
}

func (s *Seeder) seedMedicalForClient(ctx context.Context, clientID uuid.UUID, maxDiagnosesPerClient int, maxMedicationOrdersPerClient int) error {
	return s.store.ExecTx(ctx, func(q *db.Queries) error {
		actorID := s.pickCoordinatorForClient(clientID)
		var actor *uuid.UUID
		if actorID != uuid.Nil {
			actor = &actorID
		}

		diagnosisCount := 0
		if maxDiagnosesPerClient > 0 {
			diagnosisCount = gofakeit.Number(1, maxDiagnosesPerClient)
		}

		diagnosisIDs := make([]uuid.UUID, 0, diagnosisCount)
		for range diagnosisCount {
			createdDiagnosis, err := q.CreateClientDiagnosis(ctx, randomDiagnosisParams(clientID, actor))
			if err != nil {
				return fmt.Errorf("create client diagnosis: %w", err)
			}
			diagnosisIDs = append(diagnosisIDs, createdDiagnosis.ID)
			s.data.DiagnosisIDs = append(s.data.DiagnosisIDs, createdDiagnosis.ID)
		}

		medicationCount := 0
		if maxMedicationOrdersPerClient > 0 {
			medicationCount = gofakeit.Number(1, maxMedicationOrdersPerClient)
		}

		for range medicationCount {
			createdOrder, err := q.CreateClientMedicationOrder(ctx, randomMedicationOrderParams(clientID, diagnosisIDs, actor))
			if err != nil {
				return fmt.Errorf("create client medication order: %w", err)
			}
			s.data.MedicationOrderIDs = append(s.data.MedicationOrderIDs, createdOrder.ID)
		}

		return nil
	})
}

func randomDiagnosisParams(clientID uuid.UUID, actor *uuid.UUID) db.CreateClientDiagnosisParams {
	preset := oneOf(diagnosisSeedPresets)
	status := oneOf([]db.DiagnosisStatusEnum{
		db.DiagnosisStatusEnumConfirmed,
		db.DiagnosisStatusEnumConfirmed,
		db.DiagnosisStatusEnumSuspected,
		db.DiagnosisStatusEnumResolved,
		db.DiagnosisStatusEnumRuledOut,
	})
	severity := oneOf([]db.DiagnosisSeverityEnum{
		db.DiagnosisSeverityEnumMild,
		db.DiagnosisSeverityEnumModerate,
		db.DiagnosisSeverityEnumModerate,
		db.DiagnosisSeverityEnumSevere,
		db.DiagnosisSeverityEnumUnknown,
	})

	diagnosedAt := randomRecentDate(900).UTC().Truncate(24 * time.Hour)
	diagnosedOn := pgDate(diagnosedAt)

	resolvedOn := pgtype.Date{}
	if status == db.DiagnosisStatusEnumResolved {
		resolvedAt := diagnosedAt.AddDate(0, 0, gofakeit.Number(14, 360))
		now := time.Now().UTC().Truncate(24 * time.Hour)
		if resolvedAt.After(now) {
			resolvedAt = now
		}
		if resolvedAt.Before(diagnosedAt) {
			resolvedAt = diagnosedAt
		}
		resolvedOn = pgDate(resolvedAt)
	}

	notes := nullableString(gofakeit.Sentence(12), 0.35)
	description := preset.description
	title := preset.title
	clinician := nullableString(fmt.Sprintf("Dr. %s", gofakeit.LastName()), 0.25)

	return db.CreateClientDiagnosisParams{
		ClientID:            clientID,
		CodeSystem:          preset.codeSystem,
		Code:                preset.code,
		Title:               &title,
		Description:         &description,
		Status:              status,
		Severity:            severity,
		DiagnosedOn:         diagnosedOn,
		ResolvedOn:          resolvedOn,
		DiagnosingClinician: clinician,
		Notes:               notes,
		CreatedByEmployeeID: actor,
		UpdatedByEmployeeID: actor,
	}
}

func randomMedicationOrderParams(clientID uuid.UUID, diagnosisIDs []uuid.UUID, actor *uuid.UUID) db.CreateClientMedicationOrderParams {
	preset := oneOf(medicationSeedPresets)

	status := oneOf([]db.MedicationOrderStatusEnum{
		db.MedicationOrderStatusEnumActive,
		db.MedicationOrderStatusEnumActive,
		db.MedicationOrderStatusEnumPaused,
		db.MedicationOrderStatusEnumCompleted,
		db.MedicationOrderStatusEnumStopped,
	})

	startAt := randomRecentDate(540).UTC().Truncate(24 * time.Hour)
	startDate := pgDate(startAt)
	endDate := pgtype.Date{}
	if status == db.MedicationOrderStatusEnumCompleted || status == db.MedicationOrderStatusEnumStopped {
		endedAt := startAt.AddDate(0, 0, gofakeit.Number(14, 270))
		now := time.Now().UTC().Truncate(24 * time.Hour)
		if endedAt.After(now) {
			endedAt = now
		}
		if endedAt.Before(startAt) {
			endedAt = startAt
		}
		endDate = pgDate(endedAt)
	}

	adminMode := db.MedicationAdminModeEnumSelf
	if actor != nil {
		adminMode = oneOf([]db.MedicationAdminModeEnum{
			db.MedicationAdminModeEnumSelf,
			db.MedicationAdminModeEnumStaff,
			db.MedicationAdminModeEnumShared,
		})
	}
	var responsibleEmployeeID *uuid.UUID
	if adminMode == db.MedicationAdminModeEnumStaff || adminMode == db.MedicationAdminModeEnumShared {
		responsibleEmployeeID = actor
	}

	isPrn := chance(0.22)
	var prnIndication *string
	var maxDosesPer24h *int32
	if isPrn {
		prnIndication = stringPtr(preset.indication)
		maxDoses := int32(gofakeit.Number(1, 4))
		maxDosesPer24h = &maxDoses
	}

	doseAmount := gofakeit.Float64Range(preset.minDose, preset.maxDose)
	doseAmount = float64(int(doseAmount*10)) / 10
	doseUnit := preset.unit
	route := preset.route

	dosageText := fmt.Sprintf("%.1f %s", doseAmount, doseUnit)
	frequencyText := "once daily"
	schedule, _ := json.Marshal([]map[string]string{{
		"time": oneOf([]string{"08:00", "12:00", "18:00", "21:00"}),
		"dose": dosageText,
	}})
	if isPrn {
		frequencyText = "as needed"
		schedule = []byte("[]")
		dosageText = fmt.Sprintf("up to %.1f %s as needed", doseAmount, doseUnit)
	}

	var diagnosisID *uuid.UUID
	if len(diagnosisIDs) > 0 && chance(0.75) {
		picked := oneOf(diagnosisIDs)
		diagnosisID = &picked
	}

	notes := nullableString(gofakeit.Sentence(10), 0.45)

	return db.CreateClientMedicationOrderParams{
		ClientID:              clientID,
		DiagnosisID:           diagnosisID,
		MedicationName:        preset.name,
		DosageText:            dosageText,
		DoseAmount:            &doseAmount,
		DoseUnit:              &doseUnit,
		Route:                 &route,
		FrequencyText:         &frequencyText,
		Schedule:              schedule,
		IsPrn:                 isPrn,
		PrnIndication:         prnIndication,
		MaxDosesPer24h:        maxDosesPer24h,
		StartDate:             startDate,
		EndDate:               endDate,
		Status:                status,
		AdminMode:             adminMode,
		ResponsibleEmployeeID: responsibleEmployeeID,
		IsCritical:            preset.critical || chance(0.15),
		Notes:                 notes,
		SourceAttachmentUuid:  nil,
		CreatedByEmployeeID:   actor,
		UpdatedByEmployeeID:   actor,
	}
}
