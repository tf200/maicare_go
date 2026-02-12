package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	db "maicare_go/db/sqlc"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SeedData struct {
	OrganisationIDs     []uuid.UUID
	RegistrationFormIDs []uuid.UUID
	NextRegistrationIdx int
	IntakeFormIDs       []uuid.UUID
	ClientIDs           []uuid.UUID
	InCareClientIDs     []uuid.UUID
	EvaluationIDs       []uuid.UUID
	EmployeeIDs         []uuid.UUID
	CoordinatorIDs      []uuid.UUID
	ClientCoordinators  map[uuid.UUID]uuid.UUID
	AttachmentIDs       []uuid.UUID
	LocationIDs         []uuid.UUID
	SenderIDs           []uuid.UUID
}

type Seeder struct {
	store *db.Store
	data  *SeedData
}

func newSeeder(store *db.Store) *Seeder {
	return &Seeder{
		store: store,
		data: &SeedData{
			ClientCoordinators: make(map[uuid.UUID]uuid.UUID),
		},
	}
}

func (s *Seeder) SeedRegistrationForms(ctx context.Context, count int) error {
	for i := range count {
		if (i+1)%25 == 0 || i == 0 || i+1 == count {
			fmt.Printf("[seed] registration forms: %d/%d\n", i+1, count)
		}
		params := randomRegistrationFormParams(i + 1)

		created, err := s.store.CreateRegistrationForm(ctx, params)
		if err != nil {
			return fmt.Errorf("create registration form %d: %w", i+1, err)
		}

		s.data.RegistrationFormIDs = append(s.data.RegistrationFormIDs, created.ID)
	}

	return nil
}

func (s *Seeder) SeedWaitingListClients(ctx context.Context, count int) error {
	if count <= 0 {
		return nil
	}
	if len(s.data.SenderIDs) == 0 {
		return fmt.Errorf("no senders available; seed senders first")
	}
	if len(s.data.LocationIDs) == 0 {
		return fmt.Errorf("no locations available; seed locations first")
	}

	requiredRegistrationForms := s.data.NextRegistrationIdx + count
	if len(s.data.RegistrationFormIDs) < requiredRegistrationForms {
		missing := requiredRegistrationForms - len(s.data.RegistrationFormIDs)
		if err := s.SeedRegistrationForms(ctx, missing); err != nil {
			return fmt.Errorf("create additional registration forms: %w", err)
		}
	}

	topics, err := s.store.ListCarePlanTopics(ctx)
	if err != nil {
		return fmt.Errorf("list care plan topics: %w", err)
	}
	if len(topics) == 0 {
		return fmt.Errorf("no topics found; cannot create intake topic assessments")
	}

	for i := 0; i < count; i++ {
		if (i+1)%10 == 0 || i == 0 || i+1 == count {
			fmt.Printf("[seed] waiting-list clients: %d/%d\n", i+1, count)
		}
		registrationFormID := s.data.RegistrationFormIDs[s.data.NextRegistrationIdx+i]

		var createdIntakeID uuid.UUID
		var createdClientID uuid.UUID
		err := s.store.ExecTx(ctx, func(q *db.Queries) error {
			registrationForm, err := q.GetRegistrationForm(ctx, registrationFormID)
			if err != nil {
				return fmt.Errorf("get registration form: %w", err)
			}

			admissionType := deriveAdmissionType(registrationForm)
			_, err = q.UpdateRegistrationFormStatus(ctx, db.UpdateRegistrationFormStatusParams{
				ID:                        registrationFormID,
				FormStatus:                db.FormStatusEnumProcessed,
				ProcessedByEmployeeID:     nil,
				IntakeAppointmentLocation: stringPtr("Office"),
				AddmissionType: db.NullAdmissionTypeEnum{
					AdmissionTypeEnum: admissionType,
					Valid:             true,
				},
				IntakeOptions:   []byte("[]"),
				IntakeToken:     stringPtr(uuid.NewString()),
				RejectionReason: nil,
			})
			if err != nil {
				return fmt.Errorf("update registration form status: %w", err)
			}

			careType := deriveCareType(registrationForm)
			senderID := oneOf(s.data.SenderIDs)
			assignedLocationID := oneOf(s.data.LocationIDs)
			selfSufficiency := deriveSelfSufficiency(registrationForm)
			evaluationWeeks := int32(8)
			if admissionType == db.AdmissionTypeEnumCrisisAdmission {
				evaluationWeeks = 4
			}

			intakeForm, err := q.CreateSeedIntakeForm(ctx, db.CreateSeedIntakeFormParams{
				RegistrationFormID:       registrationFormID,
				DateOfIntake:             pgTimestamptz(randomRecentDate(40)),
				CareType:                 careType,
				FamilySituation:          nullableString(gofakeit.Sentence(10), 0.2),
				PsychologicalState:       nullableString(gofakeit.Sentence(10), 0.15),
				SelfSufficiency:          selfSufficiency,
				SenderID:                 &senderID,
				AssignedLocationID:       &assignedLocationID,
				RiskAssessment:           nullableString(gofakeit.Sentence(12), 0.1),
				IntakeConclusion:         db.IntakeConclusionEnumSuitable,
				IntakeConclusionNotes:    stringPtr("Suitable for waiting list placement"),
				EvaluationIntervalsWeeks: evaluationWeeks,
				Signature:                stringPtr("seed-system"),
			})
			if err != nil {
				return fmt.Errorf("create intake form: %w", err)
			}
			createdIntakeID = intakeForm.ID

			assessmentTopics := pickUniqueTopics(topics, gofakeit.Number(2, 4))
			for _, topic := range assessmentTopics {
				proposedGoalsPayload, err := json.Marshal(buildProposedGoals(topic.TopicName))
				if err != nil {
					return fmt.Errorf("marshal proposed goals: %w", err)
				}

				currentLevel := boundedLevel(selfSufficiency + int32(gofakeit.Number(-1, 1)))
				_, err = q.CreateIntakeTopicAssessment(ctx, db.CreateIntakeTopicAssessmentParams{
					IntakeFormID:  intakeForm.ID,
					TopicID:       topic.ID,
					CurrentLevel:  currentLevel,
					ProposedGoals: proposedGoalsPayload,
					Notes:         nullableString(gofakeit.Sentence(10), 0.35),
				})
				if err != nil {
					return fmt.Errorf("create intake topic assessment: %w", err)
				}
			}

			client, err := q.CreateClientDetails(ctx, db.CreateClientDetailsParams{
				IntakeFormID:       &intakeForm.ID,
				RegistrationFormID: &registrationForm.ID,
				FirstName:          registrationForm.ClientFirstName,
				LastName:           registrationForm.ClientLastName,
				DateOfBirth:        registrationForm.ClientDateOfBirth,
				Identity:           false,
				Bsn:                &registrationForm.ClientBsnNumber,
				BsnVerifiedBy:      nil,
				Email:              registrationForm.ClientEmail,
				PhoneNumber:        &registrationForm.ClientPhoneNumber,
				Gender:             registrationForm.ClientGender,
				CareType: db.NullIntakeCareTypeEnum{
					IntakeCareTypeEnum: intakeForm.CareType,
					Valid:              true,
				},
				SenderID:                   intakeForm.SenderID,
				LocationID:                 intakeForm.AssignedLocationID,
				Street:                     registrationForm.ClientStreet,
				HouseNumber:                registrationForm.ClientHouseNumber,
				HouseNumberAddition:        registrationForm.ClientHouseNumberAddition,
				PostalCode:                 registrationForm.ClientPostalCode,
				City:                       registrationForm.ClientCity,
				EducationCurrentlyEnrolled: registrationForm.EducationCurrentlyEnrolled,
				EducationInstitution:       registrationForm.EducationInstitution,
				EducationMentorName:        registrationForm.EducationMentorName,
				EducationMentorPhone:       registrationForm.EducationMentorPhone,
				EducationMentorEmail:       registrationForm.EducationMentorEmail,
				EducationAdditionalNotes:   registrationForm.EducationAdditionalNotes,
				EducationLevel:             registrationForm.EducationLevel,
				WorkCurrentlyEmployed:      registrationForm.WorkCurrentlyEmployed,
				WorkCurrentEmployer:        registrationForm.WorkCurrentEmployer,
				WorkCurrentEmployerPhone:   registrationForm.WorkEmployerPhone,
				WorkCurrentEmployerEmail:   registrationForm.WorkEmployerEmail,
				WorkCurrentPosition:        registrationForm.WorkCurrentPosition,
				WorkStartDate:              registrationForm.WorkStartDate,
				WorkAdditionalNotes:        registrationForm.WorkAdditionalNotes,
				Nationality:                &registrationForm.ClientNationality,
				RiskAggressiveBehavior:     registrationForm.RiskAggressiveBehavior,
				RiskSuicidalSelfharm:       registrationForm.RiskSuicidalSelfharm,
				RiskSubstanceAbuse:         registrationForm.RiskSubstanceAbuse,
				RiskPsychiatricIssues:      registrationForm.RiskPsychiatricIssues,
				RiskCriminalHistory:        registrationForm.RiskCriminalHistory,
				RiskFlightBehavior:         registrationForm.RiskFlightBehavior,
				RiskWeaponPossession:       registrationForm.RiskWeaponPossession,
				RiskSexualBehavior:         registrationForm.RiskSexualBehavior,
				RiskDayNightRhythm:         registrationForm.RiskDayNightRhythm,
				RiskOther:                  registrationForm.RiskOther,
				RiskOtherDescription:       registrationForm.RiskOtherDescription,
				RiskAdditionalNotes:        registrationForm.RiskAdditionalNotes,
				EvaluationIntervalsWeeks:   intakeForm.EvaluationIntervalsWeeks,
			})
			if err != nil {
				return fmt.Errorf("create client details: %w", err)
			}
			createdClientID = client.ID

			if _, err := q.CreateClientGoalsFromIntakeAssessments(ctx, db.CreateClientGoalsFromIntakeAssessmentsParams{
				ClientID:     client.ID,
				IntakeFormID: intakeForm.ID,
			}); err != nil {
				return fmt.Errorf("create client goals from intake: %w", err)
			}

			if registrationForm.Guardian1FirstName != "" {
				relationStatus := db.NullRelationStatusEnum{RelationStatusEnum: db.RelationStatusEnumPrimaryRelationship, Valid: true}
				if _, err := q.CreateEmemrgencyContact(ctx, db.CreateEmemrgencyContactParams{
					ClientID:         client.ID,
					FirstName:        &registrationForm.Guardian1FirstName,
					LastName:         &registrationForm.Guardian1LastName,
					Email:            &registrationForm.Guardian1Email,
					PhoneNumber:      &registrationForm.Guardian1PhoneNumber,
					Relationship:     &registrationForm.Guardian1Relationship,
					RelationStatus:   relationStatus,
					MedicalReports:   true,
					IncidentsReports: true,
					GoalsReports:     true,
				}); err != nil {
					return fmt.Errorf("create guardian 1 emergency contact: %w", err)
				}
			}

			if registrationForm.Guardian2FirstName != "" {
				relationStatus := db.NullRelationStatusEnum{RelationStatusEnum: db.RelationStatusEnumSecondaryRelationship, Valid: true}
				if _, err := q.CreateEmemrgencyContact(ctx, db.CreateEmemrgencyContactParams{
					ClientID:         client.ID,
					FirstName:        &registrationForm.Guardian2FirstName,
					LastName:         &registrationForm.Guardian2LastName,
					Email:            &registrationForm.Guardian2Email,
					PhoneNumber:      &registrationForm.Guardian2PhoneNumber,
					Relationship:     &registrationForm.Guardian2Relationship,
					RelationStatus:   relationStatus,
					MedicalReports:   false,
					IncidentsReports: true,
					GoalsReports:     true,
				}); err != nil {
					return fmt.Errorf("create guardian 2 emergency contact: %w", err)
				}
			}

			return nil
		})
		if err != nil {
			return fmt.Errorf("seed waiting list client %d: %w", i+1, err)
		}

		s.data.IntakeFormIDs = append(s.data.IntakeFormIDs, createdIntakeID)
		s.data.ClientIDs = append(s.data.ClientIDs, createdClientID)
	}

	s.data.NextRegistrationIdx += count

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

			if _, err := q.UpsertMainCoordinator(ctx, db.UpsertMainCoordinatorParams{
				ClientID:   client.ID,
				EmployeeID: coordinatorID,
				StartDate:  pgDate(careStartDate),
			}); err != nil {
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

func (s *Seeder) createCompletedEvaluationIfAllowed(
	ctx context.Context,
	q *db.Queries,
	client db.GetClientDetailsRow,
	createdByEmployeeID *uuid.UUID,
	goals []db.ClientGoal,
	intervalWeeks int32,
	nowDate time.Time,
) (uuid.UUID, bool, error) {
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
	eval, err := q.GetDraftGoalEvaluationByClientAndDate(ctx, db.GetDraftGoalEvaluationByClientAndDateParams{
		ClientID:       client.ID,
		EvaluationDate: evalDate,
	})
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
		progress := oneOf([]db.ClientGoalProgressEnum{
			db.ClientGoalProgressEnumRegression,
			db.ClientGoalProgressEnumLimitedProgress,
			db.ClientGoalProgressEnumGoodProgress,
			db.ClientGoalProgressEnumAchieved,
			db.ClientGoalProgressEnumBlocked,
		})
		notes := nullableString(gofakeit.Sentence(10), 0.15)
		if _, err := q.UpsertGoalEvaluationItem(ctx, db.UpsertGoalEvaluationItemParams{
			EvaluationID: eval.ID,
			GoalID:       goal.ID,
			Progress:     progress,
			Notes:        notes,
		}); err != nil {
			return uuid.Nil, false, fmt.Errorf("upsert completed evaluation item: %w", err)
		}
	}

	if _, err := q.UpdateGoalEvaluation(ctx, db.UpdateGoalEvaluationParams{
		ID:             eval.ID,
		Status:         db.NullEvaluationStatusEnum{EvaluationStatusEnum: db.EvaluationStatusEnumCompleted, Valid: true},
		EvaluationDate: pgtype.Date{},
		PeriodStart:    pgtype.Date{},
		PeriodEnd:      pgtype.Date{},
	}); err != nil {
		return uuid.Nil, false, fmt.Errorf("mark evaluation as completed: %w", err)
	}

	return eval.ID, true, nil
}

func (s *Seeder) createDraftEvaluation(
	ctx context.Context,
	q *db.Queries,
	client db.GetClientDetailsRow,
	createdByEmployeeID *uuid.UUID,
	goals []db.ClientGoal,
	intervalWeeks int32,
) (uuid.UUID, error) {
	evaluationDate := time.Now().UTC().Truncate(24 * time.Hour)
	if client.NextEvaluationDate.Valid {
		evaluationDate = client.NextEvaluationDate.Time
	}
	evalDate := pgDate(evaluationDate)
	periodEnd := evaluationDate
	periodStart := periodEnd.AddDate(0, 0, -int(intervalWeeks*7))
	overallNotes := "Seeded draft evaluation"

	eval, err := q.GetDraftGoalEvaluationByClientAndDate(ctx, db.GetDraftGoalEvaluationByClientAndDateParams{
		ClientID:       client.ID,
		EvaluationDate: evalDate,
	})
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
		progressPool := []db.ClientGoalProgressEnum{
			db.ClientGoalProgressEnumNoProgress,
			db.ClientGoalProgressEnumLimitedProgress,
			db.ClientGoalProgressEnumGoodProgress,
			db.ClientGoalProgressEnumBlocked,
		}
		if chance(0.35) {
			progressPool = []db.ClientGoalProgressEnum{
				db.ClientGoalProgressEnumLimitedProgress,
				db.ClientGoalProgressEnumGoodProgress,
				db.ClientGoalProgressEnumAchieved,
			}
		}

		if _, err := q.UpsertGoalEvaluationItem(ctx, db.UpsertGoalEvaluationItemParams{
			EvaluationID: eval.ID,
			GoalID:       goal.ID,
			Progress:     oneOf(progressPool),
			Notes:        nullableString(gofakeit.Sentence(9), 0.35),
		}); err != nil {
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

func (s *Seeder) createSeedCoordinatorProfile(ctx context.Context, q *db.Queries, locationID *uuid.UUID) (uuid.UUID, error) {
	resolvedLocationID := locationID
	if resolvedLocationID == nil && len(s.data.LocationIDs) > 0 {
		picked := oneOf(s.data.LocationIDs)
		resolvedLocationID = &picked
	}

	email := fmt.Sprintf("seed.coordinator.%s@maicare.local", strings.ToLower(gofakeit.LetterN(8)))
	user, err := q.CreateUser(ctx, db.CreateUserParams{
		Password:       "seed-password",
		Email:          email,
		IsActive:       true,
		ProfilePicture: nil,
	})
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

func (s *Seeder) SeedOrganisations(ctx context.Context, count int) error {
	for i := range count {
		if i == 0 || i+1 == count {
			fmt.Printf("[seed] organisations: %d/%d\n", i+1, count)
		}
		name := fmt.Sprintf("%s Care %d", gofakeit.Company(), i+1)
		phone := fakePhone()
		email := strings.ToLower(fmt.Sprintf("info.%d@%s.nl", i+1, sanitizeDomain(name)))
		kvk := fmt.Sprintf("%08d", gofakeit.Number(10000000, 99999999))
		btw := fmt.Sprintf("NL%09dB01", gofakeit.Number(100000000, 999999999))

		created, err := s.store.CreateOrganisation(ctx, db.CreateOrganisationParams{
			Name:                name,
			Street:              gofakeit.StreetName(),
			HouseNumber:         fmt.Sprintf("%d", gofakeit.Number(1, 350)),
			HouseNumberAddition: nullableString(strings.ToUpper(gofakeit.LetterN(1)), 0.8),
			PostalCode:          fakePostalCodeNL(),
			City:                gofakeit.City(),
			PhoneNumber:         &phone,
			Email:               &email,
			KvkNumber:           &kvk,
			BtwNumber:           &btw,
		})
		if err != nil {
			return fmt.Errorf("create organisation %d: %w", i+1, err)
		}

		s.data.OrganisationIDs = append(s.data.OrganisationIDs, created.ID)
	}

	return nil
}

func (s *Seeder) SeedLocations(ctx context.Context, perOrganisation int) error {
	if len(s.data.OrganisationIDs) == 0 {
		return fmt.Errorf("no organisations available; seed organisations first")
	}

	for _, organisationID := range s.data.OrganisationIDs {
		for i := range perOrganisation {
			if i == 0 {
				fmt.Printf("[seed] locations for organisation %s\n", organisationID)
			}
			capacity := int32(gofakeit.Number(6, 60))
			created, err := s.store.CreateLocation(ctx, db.CreateLocationParams{
				OrganisationID:      organisationID,
				Name:                fmt.Sprintf("%s %d", oneOf([]string{"Main", "North", "South", "West", "East"}), i+1),
				Street:              gofakeit.StreetName(),
				HouseNumber:         fmt.Sprintf("%d", gofakeit.Number(1, 350)),
				HouseNumberAddition: nullableString(strings.ToUpper(gofakeit.LetterN(1)), 0.85),
				PostalCode:          fakePostalCodeNL(),
				City:                gofakeit.City(),
				Capacity:            &capacity,
			})
			if err != nil {
				return fmt.Errorf("create location for organisation %s: %w", organisationID, err)
			}

			s.data.LocationIDs = append(s.data.LocationIDs, created.ID)
		}
	}

	return nil
}

func (s *Seeder) SeedSenders(ctx context.Context, count int) error {
	types := []db.SenderTypesEnum{
		db.SenderTypesEnumMainProvider,
		db.SenderTypesEnumLocalAuthority,
		db.SenderTypesEnumParticularParty,
		db.SenderTypesEnumHealthcareInstitution,
	}

	for i := range count {
		if (i+1)%20 == 0 || i == 0 || i+1 == count {
			fmt.Printf("[seed] senders: %d/%d\n", i+1, count)
		}
		name := fmt.Sprintf("%s Sender %d", gofakeit.Company(), i+1)
		email := strings.ToLower(fmt.Sprintf("contact.%d@%s.nl", i+1, sanitizeDomain(name)))
		street := gofakeit.StreetName()
		houseNumber := fmt.Sprintf("%d", gofakeit.Number(1, 350))
		houseNumberAddition := nullableString(strings.ToUpper(gofakeit.LetterN(1)), 0.85)
		postalCode := fakePostalCodeNL()
		city := gofakeit.City()
		land := "Netherlands"
		kvk := fmt.Sprintf("%08d", gofakeit.Number(10000000, 99999999))
		btw := fmt.Sprintf("NL%09dB01", gofakeit.Number(100000000, 999999999))
		phone := fakePhone()
		clientNumber := fmt.Sprintf("CL-%06d", gofakeit.Number(1, 999999))

		created, err := s.store.CreateSender(ctx, db.CreateSenderParams{
			Types:               oneOf(types),
			Name:                name,
			Street:              &street,
			HouseNumber:         &houseNumber,
			HouseNumberAddition: houseNumberAddition,
			PostalCode:          &postalCode,
			City:                &city,
			Land:                &land,
			Kvknumber:           &kvk,
			Btwnumber:           &btw,
			PhoneNumber:         &phone,
			ClientNumber:        &clientNumber,
			EmailAddress:        &email,
			Contacts:            []byte("[]"),
		})
		if err != nil {
			return fmt.Errorf("create sender %d: %w", i+1, err)
		}

		s.data.SenderIDs = append(s.data.SenderIDs, created.ID)
	}

	return nil
}

func randomRegistrationFormParams(index int) db.CreateRegistrationFormParams {
	genders := []db.GenderEnum{
		db.GenderEnumMale,
		db.GenderEnumFemale,
		db.GenderEnumOther,
		db.GenderEnumUnknown,
	}
	educationLevels := []db.EducationLevelEnum{
		db.EducationLevelEnumPrimary,
		db.EducationLevelEnumSecondary,
		db.EducationLevelEnumHigher,
		db.EducationLevelEnumNone,
	}

	firstName := gofakeit.FirstName()
	lastName := gofakeit.LastName()
	refFirst := gofakeit.FirstName()
	refLast := gofakeit.LastName()
	guardian1First := gofakeit.FirstName()
	guardian1Last := gofakeit.LastName()
	guardian2First := gofakeit.FirstName()
	guardian2Last := gofakeit.LastName()

	clientDOB := randomDate(2000, 2012)
	applicationDate := randomDate(2024, 2026)

	clientEmail := gofakeit.Email()
	referrerEmail := gofakeit.Email()
	guardian1Email := gofakeit.Email()
	guardian2Email := gofakeit.Email()

	clientGoals := randomGoals()
	educationEnrolled := gofakeit.Bool()
	workEmployed := gofakeit.Bool()

	var workStartDate pgtype.Date
	var workEmployer, workEmployerPhone, workEmployerEmail, workCurrentPosition, workAdditionalNotes *string
	if workEmployed {
		workStartDate = pgDate(randomDate(2023, 2026))
		workEmployer = stringPtr(gofakeit.Company())
		workEmployerPhone = stringPtr(fakePhone())
		workEmployerEmail = stringPtr(strings.ToLower(fmt.Sprintf("hr.%d@%s.com", index, sanitizeDomain(gofakeit.Company()))))
		workCurrentPosition = stringPtr(gofakeit.JobTitle())
		workAdditionalNotes = stringPtr(gofakeit.Sentence(8))
	} else {
		workStartDate = pgtype.Date{}
	}

	educationInstitution := nullableString(gofakeit.School(), 0.2)
	educationMentorName := nullableString(gofakeit.Name(), 0.3)
	educationMentorPhone := nullableString(fakePhone(), 0.3)
	educationMentorEmail := nullableString(strings.ToLower(fmt.Sprintf("mentor.%d@school.example.com", index)), 0.3)
	educationNotes := nullableString(gofakeit.Sentence(10), 0.4)

	careProtected := boolPtr(gofakeit.Bool())
	careAssisted := boolPtr(gofakeit.Bool())
	careTraining := boolPtr(gofakeit.Bool())
	careAmbulatory := boolPtr(gofakeit.Bool())

	riskAggressive := boolPtr(chance(0.2))
	riskSuicidal := boolPtr(chance(0.1))
	riskSubstance := boolPtr(chance(0.15))
	riskPsychiatric := boolPtr(chance(0.2))
	riskCriminal := boolPtr(chance(0.1))
	riskFlight := boolPtr(chance(0.15))
	riskWeapon := boolPtr(chance(0.05))
	riskSexual := boolPtr(chance(0.08))
	riskDayNight := boolPtr(chance(0.2))
	riskOther := boolPtr(chance(0.05))

	clientHouseNumberAddition := nullableString(gofakeit.LetterN(1), 0.8)
	referrerSignature := boolPtr(chance(0.9))

	return db.CreateRegistrationFormParams{
		ClientFirstName:               firstName,
		ClientLastName:                lastName,
		ClientDateOfBirth:             pgDate(clientDOB),
		ClientBsnNumber:               fakeBSN(index),
		ClientGender:                  oneOf(genders),
		ClientNationality:             gofakeit.Country(),
		ClientPhoneNumber:             fakePhone(),
		ClientEmail:                   clientEmail,
		ClientStreet:                  gofakeit.StreetName(),
		ClientHouseNumber:             fmt.Sprintf("%d", gofakeit.Number(1, 350)),
		ClientHouseNumberAddition:     clientHouseNumberAddition,
		ClientPostalCode:              fakePostalCodeNL(),
		ClientCity:                    gofakeit.City(),
		ReferrerFirstName:             refFirst,
		ReferrerLastName:              refLast,
		ReferrerOrganization:          gofakeit.Company(),
		ReferrerJobTitle:              gofakeit.JobTitle(),
		ReferrerPhoneNumber:           fakePhone(),
		ReferrerEmail:                 referrerEmail,
		Guardian1FirstName:            guardian1First,
		Guardian1LastName:             guardian1Last,
		Guardian1Relationship:         oneOf([]string{"Mother", "Father", "Guardian", "Sibling", "Aunt", "Uncle"}),
		Guardian1PhoneNumber:          fakePhone(),
		Guardian1Email:                guardian1Email,
		Guardian2FirstName:            guardian2First,
		Guardian2LastName:             guardian2Last,
		Guardian2Relationship:         oneOf([]string{"Mother", "Father", "Guardian", "Sibling", "Aunt", "Uncle"}),
		Guardian2PhoneNumber:          fakePhone(),
		Guardian2Email:                guardian2Email,
		EducationInstitution:          educationInstitution,
		EducationMentorName:           educationMentorName,
		EducationMentorPhone:          educationMentorPhone,
		EducationMentorEmail:          educationMentorEmail,
		EducationCurrentlyEnrolled:    educationEnrolled,
		EducationAdditionalNotes:      educationNotes,
		EducationLevel:                oneOf(educationLevels),
		WorkCurrentEmployer:           workEmployer,
		WorkEmployerPhone:             workEmployerPhone,
		WorkEmployerEmail:             workEmployerEmail,
		WorkCurrentPosition:           workCurrentPosition,
		WorkCurrentlyEmployed:         workEmployed,
		WorkStartDate:                 workStartDate,
		WorkAdditionalNotes:           workAdditionalNotes,
		CareProtectedLiving:           careProtected,
		CareAssistedIndependentLiving: careAssisted,
		CareRoomTrainingCenter:        careTraining,
		CareAmbulatoryGuidance:        careAmbulatory,
		ApplicationReason:             stringPtr(gofakeit.Sentence(16)),
		ClientGoals:                   clientGoals,
		RiskAggressiveBehavior:        riskAggressive,
		RiskSuicidalSelfharm:          riskSuicidal,
		RiskSubstanceAbuse:            riskSubstance,
		RiskPsychiatricIssues:         riskPsychiatric,
		RiskCriminalHistory:           riskCriminal,
		RiskFlightBehavior:            riskFlight,
		RiskWeaponPossession:          riskWeapon,
		RiskSexualBehavior:            riskSexual,
		RiskDayNightRhythm:            riskDayNight,
		RiskOther:                     riskOther,
		RiskOtherDescription:          nullableString(gofakeit.Sentence(8), 0.7),
		RiskAdditionalNotes:           nullableString(gofakeit.Sentence(12), 0.3),
		DocumentReferral:              nil,
		DocumentEducationReport:       nil,
		DocumentPsychiatricReport:     nil,
		DocumentDiagnosis:             nil,
		DocumentSafetyPlan:            nil,
		DocumentIDCopy:                nil,
		ApplicationDate:               pgDate(applicationDate),
		ReferrerSignature:             referrerSignature,
	}
}

func main() {
	organisationCount := flag.Int("organisations", 6, "number of organisations to seed")
	locationsPerOrg := flag.Int("locations-per-org", 2, "locations per organisation")
	senderCount := flag.Int("senders", 12, "number of senders to seed")
	count := flag.Int("count", 25, "number of registration forms to seed")
	waitingListClients := flag.Int("waiting-list-clients", 12, "number of waiting list clients to seed via intake promotion flow")
	inCareClients := flag.Int("in-care-clients", 6, "number of in-care clients to seed via waiting-list to in-care promotion flow")
	evaluationsPerInCareClient := flag.Int("evaluations-per-in-care-client", 2, "number of goal evaluations to seed for each in-care client")
	seedValue := flag.Int64("seed", time.Now().UnixNano(), "random seed")
	seedTimeout := flag.Duration("timeout", 10*time.Minute, "overall seed timeout")
	dataSource := flag.String("db", "", "database connection string (defaults to DB_SOURCE or local default)")
	flag.Parse()

	gofakeit.Seed(*seedValue)

	ctx, cancel := context.WithTimeout(context.Background(), *seedTimeout)
	defer cancel()

	dsn := strings.TrimSpace(*dataSource)
	if dsn == "" {
		dsn = strings.TrimSpace(os.Getenv("DB_SOURCE"))
	}
	if dsn == "" {
		dsn = "postgres://maicare:maicare@localhost:5432/maicare?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("cannot ping db: %v", err)
	}

	store := db.NewStore(pool)
	seeder := newSeeder(store)

	startedAt := time.Now()
	fmt.Printf("[seed] start organisations=%d locations_per_org=%d senders=%d registration_forms=%d waiting_list_clients=%d in_care_clients=%d evaluations_per_in_care_client=%d timeout=%s\n",
		*organisationCount, *locationsPerOrg, *senderCount, *count, *waitingListClients, *inCareClients, *evaluationsPerInCareClient, (*seedTimeout).String())
	if err := seeder.SeedOrganisations(ctx, *organisationCount); err != nil {
		log.Fatalf("seeding organisations failed: %v", err)
	}

	if err := seeder.SeedLocations(ctx, *locationsPerOrg); err != nil {
		log.Fatalf("seeding locations failed: %v", err)
	}

	if err := seeder.SeedSenders(ctx, *senderCount); err != nil {
		log.Fatalf("seeding senders failed: %v", err)
	}

	if err := seeder.SeedRegistrationForms(ctx, *count); err != nil {
		log.Fatalf("seeding registration forms failed: %v", err)
	}

	if err := seeder.SeedWaitingListClients(ctx, *waitingListClients); err != nil {
		log.Fatalf("seeding waiting list clients failed: %v", err)
	}

	if err := seeder.SeedInCareClients(ctx, *inCareClients); err != nil {
		log.Fatalf("seeding in-care clients failed: %v", err)
	}

	if err := seeder.SeedGoalEvaluationsForInCareClients(ctx, *evaluationsPerInCareClient); err != nil {
		log.Fatalf("seeding goal evaluations for in-care clients failed: %v", err)
	}

	fmt.Printf("Seeded %d organisations, %d locations, %d senders, %d registration forms, %d intake forms, %d total clients, %d waiting list clients, %d in-care clients, %d coordinators, %d goal evaluations in %s\n",
		len(seeder.data.OrganisationIDs),
		len(seeder.data.LocationIDs),
		len(seeder.data.SenderIDs),
		len(seeder.data.RegistrationFormIDs),
		len(seeder.data.IntakeFormIDs),
		len(seeder.data.ClientIDs),
		len(seeder.data.ClientIDs)-len(seeder.data.InCareClientIDs),
		len(seeder.data.InCareClientIDs),
		len(seeder.data.CoordinatorIDs),
		len(seeder.data.EvaluationIDs),
		time.Since(startedAt).Round(time.Millisecond),
	)
	if len(seeder.data.RegistrationFormIDs) > 0 {
		fmt.Printf("First ID: %s\n", seeder.data.RegistrationFormIDs[0])
		fmt.Printf("Last ID:  %s\n", seeder.data.RegistrationFormIDs[len(seeder.data.RegistrationFormIDs)-1])
	}
}

func pgDate(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Valid: true}
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func fakePhone() string {
	return fmt.Sprintf("06%08d", gofakeit.Number(0, 99999999))
}

func fakeBSN(index int) string {
	return fmt.Sprintf("%09d", (index%900000000)+100000000)
}

func fakePostalCodeNL() string {
	return fmt.Sprintf("%04d%s", gofakeit.Number(1000, 9999), strings.ToUpper(gofakeit.LetterN(2)))
}

func randomDate(startYear, endYear int) time.Time {
	start := time.Date(startYear, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(endYear, 12, 31, 0, 0, 0, 0, time.UTC)
	return gofakeit.DateRange(start, end)
}

func randomRecentDate(maxDaysAgo int) time.Time {
	now := time.Now().UTC()
	start := now.AddDate(0, 0, -maxDaysAgo)
	return gofakeit.DateRange(start, now)
}

func pgTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func sanitizeDomain(s string) string {
	clean := strings.ToLower(strings.TrimSpace(s))
	clean = strings.ReplaceAll(clean, "&", "and")
	clean = strings.ReplaceAll(clean, " ", "")
	if clean == "" {
		return "company"
	}
	return clean
}

func chance(probability float64) bool {
	return gofakeit.Float64() < probability
}

func nullableString(value string, nilProbability float64) *string {
	if chance(nilProbability) {
		return nil
	}
	return &value
}

func randomGoals() []string {
	pool := []string{
		"Improve school attendance",
		"Build stable daily routine",
		"Improve communication with guardians",
		"Reduce stress and anxiety",
		"Increase self-reliance in daily tasks",
		"Find suitable education pathway",
		"Develop healthy sleep schedule",
		"Improve social network support",
	}

	goalCount := gofakeit.Number(1, 3)
	goals := make([]string, 0, goalCount)
	used := map[int]struct{}{}

	for len(goals) < goalCount {
		idx := gofakeit.Number(0, len(pool)-1)
		if _, exists := used[idx]; exists {
			continue
		}
		used[idx] = struct{}{}
		goals = append(goals, pool[idx])
	}

	return goals
}

func oneOf[T any](values []T) T {
	return values[gofakeit.Number(0, len(values)-1)]
}

func boolValue(v *bool) bool {
	return v != nil && *v
}

func deriveCareType(form db.GetRegistrationFormRow) db.IntakeCareTypeEnum {
	switch {
	case boolValue(form.CareProtectedLiving):
		return db.IntakeCareTypeEnumProtectedLiving
	case boolValue(form.CareRoomTrainingCenter):
		return db.IntakeCareTypeEnumTrainingCenter
	case boolValue(form.CareAssistedIndependentLiving):
		return db.IntakeCareTypeEnumSupportedIndependentLiving
	case boolValue(form.CareAmbulatoryGuidance):
		return db.IntakeCareTypeEnumAmbulatorySupport
	default:
		if form.WorkCurrentlyEmployed || form.EducationCurrentlyEnrolled {
			return db.IntakeCareTypeEnumAmbulatorySupport
		}
		return db.IntakeCareTypeEnumProtectedLiving
	}
}

func deriveAdmissionType(form db.GetRegistrationFormRow) db.AdmissionTypeEnum {
	highUrgency := boolValue(form.RiskSuicidalSelfharm) || boolValue(form.RiskAggressiveBehavior) || boolValue(form.RiskWeaponPossession)
	if highUrgency || chance(0.18) {
		return db.AdmissionTypeEnumCrisisAdmission
	}
	return db.AdmissionTypeEnumRegularPlacement
}

func deriveSelfSufficiency(form db.GetRegistrationFormRow) int32 {
	riskCount := 0
	for _, risk := range []*bool{
		form.RiskAggressiveBehavior,
		form.RiskSuicidalSelfharm,
		form.RiskSubstanceAbuse,
		form.RiskPsychiatricIssues,
		form.RiskCriminalHistory,
		form.RiskFlightBehavior,
		form.RiskWeaponPossession,
		form.RiskSexualBehavior,
		form.RiskDayNightRhythm,
		form.RiskOther,
	} {
		if boolValue(risk) {
			riskCount++
		}
	}

	base := int32(4)
	switch {
	case riskCount >= 6:
		base = 1
	case riskCount >= 4:
		base = 2
	case riskCount >= 2:
		base = 3
	}

	if form.WorkCurrentlyEmployed || form.EducationCurrentlyEnrolled {
		base++
	}
	return boundedLevel(base)
}

func boundedLevel(level int32) int32 {
	if level < 1 {
		return 1
	}
	if level > 5 {
		return 5
	}
	return level
}

func pickUniqueTopics(topics []db.Topic, n int) []db.Topic {
	if n >= len(topics) {
		return topics
	}

	chosen := make([]db.Topic, 0, n)
	used := make(map[uuid.UUID]struct{}, n)
	for len(chosen) < n {
		topic := oneOf(topics)
		if _, exists := used[topic.ID]; exists {
			continue
		}
		used[topic.ID] = struct{}{}
		chosen = append(chosen, topic)
	}
	return chosen
}

func buildProposedGoals(topicName string) []map[string]string {
	return []map[string]string{
		{
			"title":       fmt.Sprintf("Improve %s stability", strings.ToLower(topicName)),
			"description": gofakeit.Sentence(10),
			"priority":    oneOf([]string{"medium", "high"}),
		},
		{
			"title":       fmt.Sprintf("Reach next level in %s", strings.ToLower(topicName)),
			"description": gofakeit.Sentence(8),
			"priority":    "medium",
		},
	}
}

func contractSettingsFromIntakeCareType(careType db.NullIntakeCareTypeEnum) (db.CareTypeEnum, db.PriceTimeUnitEnum, *float64, db.NullHoursTypeEnum) {
	if careType.Valid && careType.IntakeCareTypeEnum == db.IntakeCareTypeEnumAmbulatorySupport {
		hours := float64(gofakeit.Number(4, 24))
		return db.CareTypeEnumAmbulante, db.PriceTimeUnitEnumHourly, &hours, db.NullHoursTypeEnum{
			HoursTypeEnum: db.HoursTypeEnumWeekly,
			Valid:         true,
		}
	}

	return db.CareTypeEnumAccommodation, db.PriceTimeUnitEnumWeekly, nil, db.NullHoursTypeEnum{Valid: false}
}
