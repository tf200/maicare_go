package main

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"strings"

	db "maicare_go/db/sqlc"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

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
				AddmissionType:            &admissionType,
				IntakeOptions:             []byte("[]"),
				IntakeToken:               stringPtr(uuid.NewString()),
				RejectionReason:           nil,
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
				IntakeFormID:               &intakeForm.ID,
				RegistrationFormID:         &registrationForm.ID,
				FirstName:                  registrationForm.ClientFirstName,
				LastName:                   registrationForm.ClientLastName,
				DateOfBirth:                registrationForm.ClientDateOfBirth,
				Identity:                   false,
				Bsn:                        &registrationForm.ClientBsnNumber,
				BsnVerifiedBy:              nil,
				Email:                      registrationForm.ClientEmail,
				PhoneNumber:                &registrationForm.ClientPhoneNumber,
				Gender:                     registrationForm.ClientGender,
				CareType:                   &intakeForm.CareType,
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

			if _, err := q.CreateClientGoalsFromIntakeAssessments(ctx, db.CreateClientGoalsFromIntakeAssessmentsParams{ClientID: client.ID, IntakeFormID: intakeForm.ID}); err != nil {
				return fmt.Errorf("create client goals from intake: %w", err)
			}

			if registrationForm.Guardian1FirstName != "" {
				relationStatus := db.RelationStatusEnumPrimaryRelationship
				if _, err := q.CreateEmemrgencyContact(ctx, db.CreateEmemrgencyContactParams{
					ClientID:         client.ID,
					FirstName:        &registrationForm.Guardian1FirstName,
					LastName:         &registrationForm.Guardian1LastName,
					Email:            &registrationForm.Guardian1Email,
					PhoneNumber:      &registrationForm.Guardian1PhoneNumber,
					Relationship:     &registrationForm.Guardian1Relationship,
					RelationStatus:   &relationStatus,
					MedicalReports:   true,
					IncidentsReports: true,
					GoalsReports:     true,
				}); err != nil {
					return fmt.Errorf("create guardian 1 emergency contact: %w", err)
				}
			}

			if registrationForm.Guardian2FirstName != "" {
				relationStatus := db.RelationStatusEnumSecondaryRelationship
				if _, err := q.CreateEmemrgencyContact(ctx, db.CreateEmemrgencyContactParams{
					ClientID:         client.ID,
					FirstName:        &registrationForm.Guardian2FirstName,
					LastName:         &registrationForm.Guardian2LastName,
					Email:            &registrationForm.Guardian2Email,
					PhoneNumber:      &registrationForm.Guardian2PhoneNumber,
					Relationship:     &registrationForm.Guardian2Relationship,
					RelationStatus:   &relationStatus,
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

func (s *Seeder) SeedOtherIntakeForms(ctx context.Context, count int) error {
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

	conclusions := []db.IntakeConclusionEnum{
		db.IntakeConclusionEnumUnsuitable,
		db.IntakeConclusionEnumFurtherInvestigation,
		db.IntakeConclusionEnumPossiblePalcementDate,
		db.IntakeConclusionEnumOther,
	}

	for i := 0; i < count; i++ {
		if (i+1)%10 == 0 || i == 0 || i+1 == count {
			fmt.Printf("[seed] other intake forms: %d/%d\n", i+1, count)
		}

		registrationFormID := s.data.RegistrationFormIDs[s.data.NextRegistrationIdx+i]

		var createdIntakeID uuid.UUID
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
				AddmissionType:            &admissionType,
				IntakeOptions:             []byte("[]"),
				IntakeToken:               stringPtr(uuid.NewString()),
				RejectionReason:           nil,
			})
			if err != nil {
				return fmt.Errorf("update registration form status: %w", err)
			}

			conclusion := oneOf(conclusions)
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
				IntakeConclusion:         conclusion,
				IntakeConclusionNotes:    stringPtr(otherIntakeConclusionNote(conclusion)),
				EvaluationIntervalsWeeks: evaluationWeeks,
				Signature:                stringPtr("seed-system"),
			})
			if err != nil {
				return fmt.Errorf("create intake form: %w", err)
			}
			createdIntakeID = intakeForm.ID

			return nil
		})
		if err != nil {
			return fmt.Errorf("seed other intake form %d: %w", i+1, err)
		}

		s.data.IntakeFormIDs = append(s.data.IntakeFormIDs, createdIntakeID)
	}

	s.data.NextRegistrationIdx += count

	return nil
}

func otherIntakeConclusionNote(conclusion db.IntakeConclusionEnum) string {
	switch conclusion {
	case db.IntakeConclusionEnumUnsuitable:
		return "Not suitable at this time"
	case db.IntakeConclusionEnumFurtherInvestigation:
		return "Additional information required before placement decision"
	case db.IntakeConclusionEnumPossiblePalcementDate:
		return "Potential placement date discussed"
	default:
		return "Additional notes captured during intake"
	}
}

func randomRegistrationFormParams(index int) db.CreateRegistrationFormParams {
	genders := []db.GenderEnum{db.GenderEnumMale, db.GenderEnumFemale, db.GenderEnumOther, db.GenderEnumUnknown}
	educationLevels := []db.EducationLevelEnum{db.EducationLevelEnumPrimary, db.EducationLevelEnumSecondary, db.EducationLevelEnumHigher, db.EducationLevelEnumNone}

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
