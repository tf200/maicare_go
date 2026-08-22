package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	db "maicare_go/db/sqlc"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

func (s *Seeder) SeedIncidentsForClients(ctx context.Context, maxIncidentsPerClient int) error {
	if maxIncidentsPerClient <= 0 {
		return nil
	}
	canBypassRLS, err := s.store.CurrentRoleBypassesRLS(ctx)
	if err != nil {
		return fmt.Errorf("check incident bootstrap database role: %w", err)
	}
	if !canBypassRLS {
		return fmt.Errorf("incident bootstrap requires a database role with superuser or BYPASSRLS capability")
	}
	if len(s.data.ClientIDs) == 0 {
		return fmt.Errorf("no clients available; seed waiting-list or in-care clients first")
	}

	for i, clientID := range s.data.ClientIDs {
		if (i+1)%10 == 0 || i == 0 || i+1 == len(s.data.ClientIDs) {
			fmt.Printf("[seed] client incidents: %d/%d\n", i+1, len(s.data.ClientIDs))
		}

		if err := s.seedIncidentsForClient(ctx, clientID, maxIncidentsPerClient); err != nil {
			return fmt.Errorf("seed incidents for client %s: %w", clientID, err)
		}
	}

	return nil
}

func (s *Seeder) seedIncidentsForClient(ctx context.Context, clientID uuid.UUID, maxIncidentsPerClient int) error {
	createdIncidentIDs := make([]uuid.UUID, 0)

	err := s.store.ExecTx(ctx, func(q *db.Queries) error {
		client, err := q.GetClientDetails(ctx, clientID)
		if err != nil {
			return fmt.Errorf("get client details: %w", err)
		}

		incidentCount := gofakeit.Number(1, maxIncidentsPerClient)
		for range incidentCount {
			employeeID, err := s.pickIncidentEmployee()
			if err != nil {
				return err
			}

			locationID, err := s.resolveIncidentLocation(client.LocationID)
			if err != nil {
				return err
			}

			createdID, err := q.SeedIncident(ctx, randomIncidentParams(client, employeeID, locationID))
			if err != nil {
				return fmt.Errorf("create incident: %w", err)
			}
			createdIncidentIDs = append(createdIncidentIDs, createdID)
		}

		return nil
	})
	if err != nil {
		return err
	}

	s.data.IncidentIDs = append(s.data.IncidentIDs, createdIncidentIDs...)
	return nil
}

func (s *Seeder) pickIncidentEmployee() (uuid.UUID, error) {
	if len(s.data.CoordinatorIDs) > 0 {
		return oneOf(s.data.CoordinatorIDs), nil
	}
	if len(s.data.EmployeeIDs) > 0 {
		return oneOf(s.data.EmployeeIDs), nil
	}

	return uuid.Nil, fmt.Errorf("no employees available for incidents; seed coordinators or employees first")
}

func (s *Seeder) resolveIncidentLocation(clientLocationID *uuid.UUID) (uuid.UUID, error) {
	if clientLocationID != nil {
		return *clientLocationID, nil
	}
	if len(s.data.LocationIDs) == 0 {
		return uuid.Nil, fmt.Errorf("no locations available for incidents; seed locations first")
	}

	return oneOf(s.data.LocationIDs), nil
}

func randomIncidentParams(client db.GetClientDetailsRow, employeeID uuid.UUID, locationID uuid.UUID) db.SeedIncidentParams {
	reporterInvolvement := oneOf([]db.IncidentReporterInvolvementEnum{
		db.IncidentReporterInvolvementEnumDirectlyInvolved,
		db.IncidentReporterInvolvementEnumWitness,
		db.IncidentReporterInvolvementEnumFoundAfterwards,
		db.IncidentReporterInvolvementEnumAlarmed,
	})

	incidentType := oneOf([]db.IncidentTypeEnum{
		db.IncidentTypeEnumAccident,
		db.IncidentTypeEnumViolence,
		db.IncidentTypeEnumSelfHarm,
		db.IncidentTypeEnumClientAbsence,
		db.IncidentTypeEnumMedicines,
		db.IncidentTypeEnumOrganization,
		db.IncidentTypeEnumUseProhibitedSubstances,
		db.IncidentTypeEnumFireWaterDamage,
		db.IncidentTypeEnumOther,
	})

	severity := oneOf([]db.SeverityOfIncidentEnum{
		db.SeverityOfIncidentEnumNearIncident,
		db.SeverityOfIncidentEnumNearIncident,
		db.SeverityOfIncidentEnumLessSerious,
		db.SeverityOfIncidentEnumLessSerious,
		db.SeverityOfIncidentEnumSerious,
		db.SeverityOfIncidentEnumFatal,
	})

	recurrenceRisk := oneOf([]db.RecurrenceRiskEnum{
		db.RecurrenceRiskEnumVeryLow,
		db.RecurrenceRiskEnumMeans,
		db.RecurrenceRiskEnumMeans,
		db.RecurrenceRiskEnumHigh,
		db.RecurrenceRiskEnumVeryHigh,
	})

	physicalInjury := oneOf([]db.PhysicalInjuryEnum{
		db.PhysicalInjuryEnumNoInjuries,
		db.PhysicalInjuryEnumNotNoticeableYet,
		db.PhysicalInjuryEnumBruisingSwelling,
		db.PhysicalInjuryEnumBrokenBones,
		db.PhysicalInjuryEnumShortnessOfBreath,
		db.PhysicalInjuryEnumDeath,
		db.PhysicalInjuryEnumOther,
	})

	psychologicalDamage := oneOf([]db.PsychologicalDamageEnum{
		db.PsychologicalDamageEnumNo,
		db.PsychologicalDamageEnumNo,
		db.PsychologicalDamageEnumNotNoticeableYet,
		db.PsychologicalDamageEnumDrowsiness,
		db.PsychologicalDamageEnumUnrest,
		db.PsychologicalDamageEnumOther,
	})

	neededConsultation := oneOf([]db.NeededConsultationEnum{
		db.NeededConsultationEnumNo,
		db.NeededConsultationEnumNotClear,
		db.NeededConsultationEnumConsultGp,
		db.NeededConsultationEnumHospitalization,
	})

	occurredAt := randomRecentDate(240).UTC().Truncate(time.Minute)
	informedParties := randomInformedParties()
	causeCategories := randomIncidentCauseCategories()
	followUpActions := randomIncidentFollowUpActions()

	incidentExplanation := nullableString(gofakeit.Sentence(16), 0.2)
	incidentPreventSteps := nullableString(gofakeit.Sentence(12), 0.35)
	incidentTakenMeasures := nullableString(gofakeit.Sentence(12), 0.35)
	causeExplanation := nullableString(gofakeit.Sentence(12), 0.4)
	followUpNotes := nullableString(gofakeit.Sentence(10), 0.45)
	additionalDetails := nullableString(gofakeit.Paragraph(1, 2, 12, " "), 0.55)

	var physicalInjuryDesc *string
	if physicalInjury != db.PhysicalInjuryEnumNoInjuries {
		physicalInjuryDesc = nullableString(gofakeit.Sentence(10), 0.25)
	}

	var psychologicalDamageDesc *string
	if psychologicalDamage != db.PsychologicalDamageEnumNo {
		psychologicalDamageDesc = nullableString(gofakeit.Sentence(10), 0.3)
	}

	isEmployeeAbsent := severity == db.SeverityOfIncidentEnumSerious || severity == db.SeverityOfIncidentEnumFatal || chance(0.2)

	emails := randomIncidentEmails(client)

	return db.SeedIncidentParams{
		EmployeeID:              employeeID,
		LocationID:              locationID,
		ReporterInvolvement:     reporterInvolvement,
		InformedParties:         informedParties,
		OccurredAt:              pgTimestamptz(occurredAt),
		IncidentType:            incidentType,
		SeverityOfIncident:      severity,
		IncidentExplanation:     incidentExplanation,
		RecurrenceRisk:          recurrenceRisk,
		IncidentPreventSteps:    incidentPreventSteps,
		IncidentTakenMeasures:   incidentTakenMeasures,
		CauseCategories:         causeCategories,
		CauseExplanation:        causeExplanation,
		PhysicalInjury:          physicalInjury,
		PhysicalInjuryDesc:      physicalInjuryDesc,
		PsychologicalDamage:     psychologicalDamage,
		PsychologicalDamageDesc: psychologicalDamageDesc,
		NeededConsultation:      neededConsultation,
		FollowUpActions:         followUpActions,
		FollowUpNotes:           followUpNotes,
		IsEmployeeAbsent:        isEmployeeAbsent,
		AdditionalDetails:       additionalDetails,
		ClientID:                client.ID,
		Emails:                  emails,
	}
}

func randomInformedParties() []db.InformedPartyEnum {
	pool := []db.InformedPartyEnum{
		db.InformedPartyEnumParentsGuardians,
		db.InformedPartyEnumCareCoordinator,
		db.InformedPartyEnumReferrer,
		db.InformedPartyEnumHealthcareProvider,
		db.InformedPartyEnumInspectorate,
		db.InformedPartyEnumPolice,
		db.InformedPartyEnumOther,
	}

	return pickUniqueEnums(pool, gofakeit.Number(0, 3))
}

func randomIncidentCauseCategories() []db.IncidentCauseCategoryEnum {
	pool := []db.IncidentCauseCategoryEnum{
		db.IncidentCauseCategoryEnumTechnical,
		db.IncidentCauseCategoryEnumOrganizational,
		db.IncidentCauseCategoryEnumEmployeeRelated,
		db.IncidentCauseCategoryEnumClientRelated,
		db.IncidentCauseCategoryEnumExternal,
		db.IncidentCauseCategoryEnumOther,
	}

	return pickUniqueEnums(pool, gofakeit.Number(1, 3))
}

func randomIncidentFollowUpActions() []db.IncidentFollowUpActionEnum {
	pool := []db.IncidentFollowUpActionEnum{
		db.IncidentFollowUpActionEnumNotifyParentsGuardians,
		db.IncidentFollowUpActionEnumNotifyReferrer,
		db.IncidentFollowUpActionEnumNotifyInspectorate,
		db.IncidentFollowUpActionEnumMedicalConsultation,
		db.IncidentFollowUpActionEnumCarePlanAdjustment,
		db.IncidentFollowUpActionEnumTeamEvaluation,
		db.IncidentFollowUpActionEnumOther,
	}

	return pickUniqueEnums(pool, gofakeit.Number(0, 3))
}

func randomIncidentEmails(client db.GetClientDetailsRow) []string {
	emails := make([]string, 0, 3)
	if strings.TrimSpace(client.Email) != "" && chance(0.85) {
		emails = append(emails, client.Email)
	}
	if client.SenderEmailAddress != nil && strings.TrimSpace(*client.SenderEmailAddress) != "" && chance(0.35) {
		emails = append(emails, *client.SenderEmailAddress)
	}
	if len(emails) == 0 && chance(0.4) {
		emails = append(emails, strings.ToLower(gofakeit.Email()))
	}

	return emails
}

func pickUniqueEnums[T comparable](pool []T, count int) []T {
	if count <= 0 || len(pool) == 0 {
		return []T{}
	}
	if count >= len(pool) {
		result := make([]T, len(pool))
		copy(result, pool)
		return result
	}

	result := make([]T, 0, count)
	used := make(map[T]struct{}, count)
	for len(result) < count {
		item := oneOf(pool)
		if _, exists := used[item]; exists {
			continue
		}
		used[item] = struct{}{}
		result = append(result, item)
	}

	return result
}
