package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	db "maicare_go/db/sqlc"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Seeder) SeedOutOfCareClients(ctx context.Context, count int, evaluationsPerClient int) error {
	if count <= 0 {
		return nil
	}

	if evaluationsPerClient <= 0 {
		evaluationsPerClient = 1
	}

	startClientCount := len(s.data.ClientIDs)
	if err := s.SeedInCareClients(ctx, count); err != nil {
		return fmt.Errorf("seed base in-care clients for out-of-care flow: %w", err)
	}

	if len(s.data.ClientIDs) < startClientCount+count {
		return fmt.Errorf("expected %d newly seeded clients for out-of-care flow, got %d", count, len(s.data.ClientIDs)-startClientCount)
	}

	outOfCareClientIDs := append([]uuid.UUID(nil), s.data.ClientIDs[startClientCount:]...)
	for i, clientID := range outOfCareClientIDs {
		if (i+1)%10 == 0 || i == 0 || i+1 == len(outOfCareClientIDs) {
			fmt.Printf("[seed] out-of-care clients: %d/%d\n", i+1, len(outOfCareClientIDs))
		}

		if err := s.seedEvaluationsForClient(ctx, clientID, evaluationsPerClient); err != nil {
			return fmt.Errorf("seed goal evaluations before out-of-care transition for client %s: %w", clientID, err)
		}

		err := s.store.ExecTx(ctx, func(q *db.Queries) error {
			client, err := q.GetClientDetails(ctx, clientID)
			if err != nil {
				return fmt.Errorf("get client details: %w", err)
			}

			careStartDate := time.Now().UTC().AddDate(0, 0, -90)
			if client.CareStartDate.Valid {
				careStartDate = client.CareStartDate.Time
			}

			latestPossibleDischarge := time.Now().UTC().AddDate(0, 0, -1)
			earliestPossibleDischarge := careStartDate.AddDate(0, 0, gofakeit.Number(30, 150))
			if earliestPossibleDischarge.After(latestPossibleDischarge) {
				earliestPossibleDischarge = latestPossibleDischarge
			}
			dischargeDate := gofakeit.DateRange(earliestPossibleDischarge, latestPossibleDischarge)

			finalEvaluation := strings.TrimSpace(gofakeit.Paragraph(1, 3, 12, " "))
			dischargeReason := oneOf([]db.DischargeReasonEnum{
				db.DischargeReasonEnumTreatmentCompleted,
				db.DischargeReasonEnumTerminatedByMutualAgreement,
				db.DischargeReasonEnumTerminatedByClient,
				db.DischargeReasonEnumTerminatedByProvider,
				db.DischargeReasonEnumTerminatedDueToExternalFactors,
				db.DischargeReasonEnumOther,
			})

			updatedClient, err := q.PutClientOutOfCare(ctx, db.PutClientOutOfCareParams{
				ID:            client.ID,
				Status:        db.ClientStatusEnumOutOfCare,
				DischargeDate: pgDate(dischargeDate),
				DischargeReason: db.NullDischargeReasonEnum{
					DischargeReasonEnum: dischargeReason,
					Valid:               true,
				},
				FinalEvaluation: stringPtr(finalEvaluation),
			})
			if err != nil {
				return fmt.Errorf("put client out of care: %w", err)
			}

			historyReason := fmt.Sprintf("seed_client_put_out_of_care:%s", dischargeReason)
			if _, err := q.CreateClientStatusHistory(ctx, db.CreateClientStatusHistoryParams{
				ClientID:  updatedClient.ID,
				OldStatus: stringPtr(string(client.Status)),
				NewStatus: string(updatedClient.Status),
				Reason:    &historyReason,
			}); err != nil {
				return fmt.Errorf("create client status history: %w", err)
			}

			return nil
		})
		if err != nil {
			return fmt.Errorf("seed out-of-care client %d: %w", i+1, err)
		}
	}

	s.data.OutOfCareClientIDs = append(s.data.OutOfCareClientIDs, outOfCareClientIDs...)
	s.data.InCareClientIDs = excludeClientIDs(s.data.InCareClientIDs, outOfCareClientIDs)

	return nil
}

func (s *Seeder) SeedInCareClients(ctx context.Context, count int) error {
	if count <= 0 {
		return nil
	}

	startClientCount := len(s.data.ClientIDs)
	if err := s.SeedWaitingListClients(ctx, count); err != nil {
		return fmt.Errorf("seed base waiting-list clients for in-care flow: %w", err)
	}

	if len(s.data.ClientIDs) < startClientCount+count {
		return fmt.Errorf("expected %d newly seeded clients, got %d", count, len(s.data.ClientIDs)-startClientCount)
	}

	promotedClientIDs := append([]uuid.UUID(nil), s.data.ClientIDs[startClientCount:]...)

	for i, clientID := range promotedClientIDs {
		if (i+1)%10 == 0 || i == 0 || i+1 == len(promotedClientIDs) {
			fmt.Printf("[seed] in-care clients: %d/%d\n", i+1, len(promotedClientIDs))
		}

		err := s.store.ExecTx(ctx, func(q *db.Queries) error {
			client, err := q.GetClientDetails(ctx, clientID)
			if err != nil {
				return fmt.Errorf("get client details: %w", err)
			}

			careStartDate := randomRecentDate(90)
			placedInCareAt := careStartDate.AddDate(0, 0, -gofakeit.Number(1, 14))

			if _, err := q.PutClientInCare(ctx, db.PutClientInCareParams{
				ID:             client.ID,
				Status:         db.ClientStatusEnumInCare,
				CareStartDate:  pgDate(careStartDate),
				PlacedInCareAt: pgTimestamptz(placedInCareAt),
			}); err != nil {
				return fmt.Errorf("put client in care: %w", err)
			}

			coordinatorID, err := s.createSeedCoordinatorProfile(ctx, q, client.LocationID)
			if err != nil {
				return fmt.Errorf("create coordinator profile: %w", err)
			}

			if _, err := q.UpsertMainCoordinator(ctx, db.UpsertMainCoordinatorParams{ClientID: client.ID, EmployeeID: coordinatorID, StartDate: pgDate(careStartDate)}); err != nil {
				return fmt.Errorf("upsert main coordinator: %w", err)
			}

			contractCareType, priceUnit, hours, hoursType := contractSettingsFromIntakeCareType(client.CareType)
			price := float64(gofakeit.Number(450, 1800))
			careName := "Residential care placement"
			if contractCareType == db.CareTypeEnumAmbulante {
				price = float64(gofakeit.Number(45, 125))
				careName = "Ambulatory guidance"
			}

			contractSenderID := oneOf(s.data.SenderIDs)
			if client.SenderID != nil {
				contractSenderID = *client.SenderID
			}

			if _, err := q.CreateContract(ctx, db.CreateContractParams{
				TypeID:          nil,
				Status:          db.ContractStatusEnumApproved,
				StartDate:       pgTimestamptz(careStartDate),
				EndDate:         pgTimestamptz(careStartDate.AddDate(0, gofakeit.Number(4, 12), 0)),
				ReminderPeriod:  90,
				Vat:             nil,
				Price:           price,
				PriceTimeUnit:   priceUnit,
				Hours:           hours,
				HoursType:       hoursType,
				CareName:        careName,
				CareType:        contractCareType,
				ClientID:        client.ID,
				SenderID:        contractSenderID,
				AttachmentIds:   []uuid.UUID{},
				FinancingAct:    db.FinancingActEnumWMO,
				FinancingOption: db.FinancingOptionEnumPGB,
			}); err != nil {
				return fmt.Errorf("create approved contract: %w", err)
			}

			return nil
		})
		if err != nil {
			return fmt.Errorf("seed in-care client %d: %w", i+1, err)
		}

		s.data.InCareClientIDs = append(s.data.InCareClientIDs, clientID)
		s.data.ClientCoordinators[clientID] = s.data.CoordinatorIDs[len(s.data.CoordinatorIDs)-1]
	}

	return nil
}

func (s *Seeder) SeedGoalEvaluationsForInCareClients(ctx context.Context, evaluationsPerClient int) error {
	if evaluationsPerClient <= 0 || len(s.data.InCareClientIDs) == 0 {
		return nil
	}

	for i, clientID := range s.data.InCareClientIDs {
		if (i+1)%10 == 0 || i == 0 || i+1 == len(s.data.InCareClientIDs) {
			fmt.Printf("[seed] goal evaluations for in-care clients: %d/%d\n", i+1, len(s.data.InCareClientIDs))
		}

		if err := s.seedEvaluationsForClient(ctx, clientID, evaluationsPerClient); err != nil {
			return fmt.Errorf("seed goal evaluations for in-care client %s: %w", clientID, err)
		}
	}

	return nil
}

func (s *Seeder) seedEvaluationsForClient(ctx context.Context, clientID uuid.UUID, evaluationsPerClient int) error {
	return s.store.ExecTx(ctx, func(q *db.Queries) error {
		client, err := q.GetClientDetails(ctx, clientID)
		if err != nil {
			return fmt.Errorf("get client details: %w", err)
		}

		goals, err := q.ListActiveGoalsByClientID(ctx, clientID)
		if err != nil {
			return fmt.Errorf("list active goals: %w", err)
		}
		if len(goals) == 0 {
			return nil
		}

		employeeID := s.pickCoordinatorForClient(clientID)
		var createdByEmployeeID *uuid.UUID
		if employeeID != uuid.Nil {
			createdByEmployeeID = &employeeID
		}
		nowDate := time.Now().UTC().Truncate(24 * time.Hour)
		intervalWeeks := client.EvaluationIntervalsWeeks
		if intervalWeeks <= 0 {
			intervalWeeks = 12
		}

		createdCount := 0
		if evaluationsPerClient >= 2 {
			if evalID, ok, err := s.createCompletedEvaluationIfAllowed(ctx, q, client, createdByEmployeeID, goals, intervalWeeks, nowDate); err != nil {
				return err
			} else if ok {
				s.data.EvaluationIDs = append(s.data.EvaluationIDs, evalID)
				createdCount++
				client, err = q.GetClientDetails(ctx, clientID)
				if err != nil {
					return fmt.Errorf("refresh client details after completed evaluation: %w", err)
				}
			}
		}

		if createdCount < evaluationsPerClient {
			evalID, err := s.createDraftEvaluation(ctx, q, client, createdByEmployeeID, goals, intervalWeeks)
			if err != nil {
				return err
			}
			s.data.EvaluationIDs = append(s.data.EvaluationIDs, evalID)
		}

		return nil
	})
}

func (s *Seeder) createCompletedEvaluationIfAllowed(ctx context.Context, q *db.Queries, client db.GetClientDetailsRow, createdByEmployeeID *uuid.UUID, goals []db.ClientGoal, intervalWeeks int32, nowDate time.Time) (uuid.UUID, bool, error) {
	if !client.NextEvaluationDate.Valid {
		return uuid.Nil, false, nil
	}

	dueDate := client.NextEvaluationDate.Time
	windowStart := dueDate.AddDate(0, 0, -14)
	if nowDate.Before(windowStart) {
		return uuid.Nil, false, nil
	}

	periodEnd := dueDate
	periodStart := periodEnd.AddDate(0, 0, -int(intervalWeeks*7))
	overallNotes := "Seeded completed evaluation"

	evalDate := pgDate(dueDate)
	eval, err := q.GetDraftGoalEvaluationByClientAndDate(ctx, db.GetDraftGoalEvaluationByClientAndDateParams{ClientID: client.ID, EvaluationDate: evalDate})
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, false, fmt.Errorf("get existing draft evaluation: %w", err)
		}

		eval, err = q.CreateGoalEvaluation(ctx, db.CreateGoalEvaluationParams{
			ClientID:                client.ID,
			EvaluationDate:          evalDate,
			PeriodStart:             pgDate(periodStart),
			PeriodEnd:               pgDate(periodEnd),
			EvaluationIntervalWeeks: intervalWeeks,
			Status:                  db.EvaluationStatusEnumDraft,
			OverallNotes:            &overallNotes,
			CreatedByEmployeeID:     createdByEmployeeID,
		})
		if err != nil {
			return uuid.Nil, false, fmt.Errorf("create completed-candidate evaluation: %w", err)
		}
	} else {
		eval, err = q.UpdateGoalEvaluation(ctx, db.UpdateGoalEvaluationParams{
			ID:                      eval.ID,
			EvaluationDate:          evalDate,
			PeriodStart:             pgDate(periodStart),
			PeriodEnd:               pgDate(periodEnd),
			EvaluationIntervalWeeks: &intervalWeeks,
			Status:                  db.NullEvaluationStatusEnum{Valid: false},
			OverallNotes:            &overallNotes,
		})
		if err != nil {
			return uuid.Nil, false, fmt.Errorf("update completed-candidate evaluation: %w", err)
		}
	}

	for _, goal := range goals {
		progress := oneOf([]db.ClientGoalProgressEnum{db.ClientGoalProgressEnumRegression, db.ClientGoalProgressEnumLimitedProgress, db.ClientGoalProgressEnumGoodProgress, db.ClientGoalProgressEnumAchieved, db.ClientGoalProgressEnumBlocked})
		notes := nullableString(gofakeit.Sentence(10), 0.15)
		if _, err := q.UpsertGoalEvaluationItem(ctx, db.UpsertGoalEvaluationItemParams{EvaluationID: eval.ID, GoalID: goal.ID, Progress: progress, Notes: notes}); err != nil {
			return uuid.Nil, false, fmt.Errorf("upsert completed evaluation item: %w", err)
		}
	}

	if _, err := q.UpdateGoalEvaluation(ctx, db.UpdateGoalEvaluationParams{ID: eval.ID, Status: db.NullEvaluationStatusEnum{EvaluationStatusEnum: db.EvaluationStatusEnumCompleted, Valid: true}, EvaluationDate: pgtype.Date{}, PeriodStart: pgtype.Date{}, PeriodEnd: pgtype.Date{}}); err != nil {
		return uuid.Nil, false, fmt.Errorf("mark evaluation as completed: %w", err)
	}

	return eval.ID, true, nil
}

func (s *Seeder) createDraftEvaluation(ctx context.Context, q *db.Queries, client db.GetClientDetailsRow, createdByEmployeeID *uuid.UUID, goals []db.ClientGoal, intervalWeeks int32) (uuid.UUID, error) {
	evaluationDate := time.Now().UTC().Truncate(24 * time.Hour)
	if client.NextEvaluationDate.Valid {
		evaluationDate = client.NextEvaluationDate.Time
	}
	evalDate := pgDate(evaluationDate)
	periodEnd := evaluationDate
	periodStart := periodEnd.AddDate(0, 0, -int(intervalWeeks*7))
	overallNotes := "Seeded draft evaluation"

	eval, err := q.GetDraftGoalEvaluationByClientAndDate(ctx, db.GetDraftGoalEvaluationByClientAndDateParams{ClientID: client.ID, EvaluationDate: evalDate})
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("get existing draft evaluation: %w", err)
		}

		eval, err = q.CreateGoalEvaluation(ctx, db.CreateGoalEvaluationParams{
			ClientID:                client.ID,
			EvaluationDate:          evalDate,
			PeriodStart:             pgDate(periodStart),
			PeriodEnd:               pgDate(periodEnd),
			EvaluationIntervalWeeks: intervalWeeks,
			Status:                  db.EvaluationStatusEnumDraft,
			OverallNotes:            &overallNotes,
			CreatedByEmployeeID:     createdByEmployeeID,
		})
		if err != nil {
			return uuid.Nil, fmt.Errorf("create draft evaluation: %w", err)
		}
	} else {
		eval, err = q.UpdateGoalEvaluation(ctx, db.UpdateGoalEvaluationParams{
			ID:                      eval.ID,
			EvaluationDate:          evalDate,
			PeriodStart:             pgDate(periodStart),
			PeriodEnd:               pgDate(periodEnd),
			EvaluationIntervalWeeks: &intervalWeeks,
			Status:                  db.NullEvaluationStatusEnum{Valid: false},
			OverallNotes:            &overallNotes,
		})
		if err != nil {
			return uuid.Nil, fmt.Errorf("update draft evaluation: %w", err)
		}
	}

	for _, goal := range goals {
		progressPool := []db.ClientGoalProgressEnum{db.ClientGoalProgressEnumNoProgress, db.ClientGoalProgressEnumLimitedProgress, db.ClientGoalProgressEnumGoodProgress, db.ClientGoalProgressEnumBlocked}
		if chance(0.35) {
			progressPool = []db.ClientGoalProgressEnum{db.ClientGoalProgressEnumLimitedProgress, db.ClientGoalProgressEnumGoodProgress, db.ClientGoalProgressEnumAchieved}
		}

		if _, err := q.UpsertGoalEvaluationItem(ctx, db.UpsertGoalEvaluationItemParams{EvaluationID: eval.ID, GoalID: goal.ID, Progress: oneOf(progressPool), Notes: nullableString(gofakeit.Sentence(9), 0.35)}); err != nil {
			return uuid.Nil, fmt.Errorf("upsert draft evaluation item: %w", err)
		}
	}

	return eval.ID, nil
}

func (s *Seeder) pickCoordinatorForClient(clientID uuid.UUID) uuid.UUID {
	if coordinatorID, ok := s.data.ClientCoordinators[clientID]; ok {
		return coordinatorID
	}
	if len(s.data.CoordinatorIDs) > 0 {
		return oneOf(s.data.CoordinatorIDs)
	}
	if len(s.data.EmployeeIDs) > 0 {
		return oneOf(s.data.EmployeeIDs)
	}
	return uuid.Nil
}

func excludeClientIDs(existing []uuid.UUID, remove []uuid.UUID) []uuid.UUID {
	if len(existing) == 0 || len(remove) == 0 {
		return existing
	}

	removeSet := make(map[uuid.UUID]struct{}, len(remove))
	for _, id := range remove {
		removeSet[id] = struct{}{}
	}

	kept := make([]uuid.UUID, 0, len(existing))
	for _, id := range existing {
		if _, shouldRemove := removeSet[id]; shouldRemove {
			continue
		}
		kept = append(kept, id)
	}
	return kept
}

func (s *Seeder) createSeedCoordinatorProfile(ctx context.Context, q *db.Queries, locationID *uuid.UUID) (uuid.UUID, error) {
	resolvedLocationID := locationID
	if resolvedLocationID == nil && len(s.data.LocationIDs) > 0 {
		picked := oneOf(s.data.LocationIDs)
		resolvedLocationID = &picked
	}

	email := fmt.Sprintf("seed.coordinator.%s@maicare.local", strings.ToLower(gofakeit.LetterN(8)))
	user, err := q.CreateUser(ctx, db.CreateUserParams{Password: "seed-password", Email: email, IsActive: true, ProfilePicture: nil})
	if err != nil {
		return uuid.Nil, fmt.Errorf("create user: %w", err)
	}

	contractHours := 36.0
	contractRate := 58.0
	employeeNumber := fmt.Sprintf("EMP-%06d", gofakeit.Number(1, 999999))
	employmentNumber := fmt.Sprintf("CONT-%06d", gofakeit.Number(1, 999999))
	workEmail := email
	privateEmail := strings.ToLower(gofakeit.Email())
	workPhone := fakePhone()
	privatePhone := fakePhone()
	homePhone := fakePhone()

	employee, err := q.CreateEmployeeProfile(ctx, db.CreateEmployeeProfileParams{
		UserID:              user.ID,
		FirstName:           gofakeit.FirstName(),
		LastName:            gofakeit.LastName(),
		Bsn:                 fmt.Sprintf("%09d", gofakeit.Number(100000000, 999999999)),
		Street:              gofakeit.StreetName(),
		HouseNumber:         fmt.Sprintf("%d", gofakeit.Number(1, 350)),
		HouseNumberAddition: nullableString(strings.ToUpper(gofakeit.LetterN(1)), 0.8),
		PostalCode:          fakePostalCodeNL(),
		City:                gofakeit.City(),
		Position:            stringPtr("Care Coordinator"),
		Department:          stringPtr("Youth Care"),
		EmployeeNumber:      &employeeNumber,
		EmploymentNumber:    &employmentNumber,
		PrivateEmailAddress: &privateEmail,
		WorkEmailAddress:    &workEmail,
		WorkPhoneNumber:     &workPhone,
		PrivatePhoneNumber:  &privatePhone,
		DateOfBirth:         pgDate(randomDate(1975, 1998)),
		HomeTelephoneNumber: &homePhone,
		Gender:              oneOf([]db.GenderEnum{db.GenderEnumMale, db.GenderEnumFemale, db.GenderEnumOther}),
		LocationID:          resolvedLocationID,
		ContractHours:       &contractHours,
		ContractEndDate:     pgtype.Date{},
		ContractStartDate:   pgDate(time.Now().AddDate(-1, 0, 0)),
		ContractType:        db.EmployeeContractTypeEnumLoondienst,
		ContractRate:        &contractRate,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("create employee profile: %w", err)
	}

	s.data.EmployeeIDs = append(s.data.EmployeeIDs, employee.ID)
	s.data.CoordinatorIDs = append(s.data.CoordinatorIDs, employee.ID)

	return employee.ID, nil
}
