package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	db "maicare_go/db/sqlc"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SeedData struct {
	OrganisationIDs     []uuid.UUID
	RegistrationFormIDs []uuid.UUID
	IntakeFormIDs       []uuid.UUID
	ClientIDs           []uuid.UUID
	EmployeeIDs         []uuid.UUID
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
		data:  &SeedData{},
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

	if len(s.data.RegistrationFormIDs) < count {
		missing := count - len(s.data.RegistrationFormIDs)
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
		registrationFormID := s.data.RegistrationFormIDs[i]

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
				EvaluationIntarvalsWeeks:   intakeForm.EvaluationIntervalsWeeks,
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

	return nil
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
	fmt.Printf("[seed] start organisations=%d locations_per_org=%d senders=%d registration_forms=%d waiting_list_clients=%d timeout=%s\n",
		*organisationCount, *locationsPerOrg, *senderCount, *count, *waitingListClients, (*seedTimeout).String())
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

	fmt.Printf("Seeded %d organisations, %d locations, %d senders, %d registration forms, %d intake forms, %d waiting list clients in %s\n",
		len(seeder.data.OrganisationIDs),
		len(seeder.data.LocationIDs),
		len(seeder.data.SenderIDs),
		len(seeder.data.RegistrationFormIDs),
		len(seeder.data.IntakeFormIDs),
		len(seeder.data.ClientIDs),
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
