package clientp

import (
	"context"
	"errors"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

var ErrIntakeNotSuitable = errors.New("intake conclusion is not suitable")

// PromoteIntakeToClient promotes an intake form and its related data to a full client record
func (s *clientService) PromoteIntakeToClient(ctx context.Context, req *PromoteIntakeToClientRequest) (*PromoteIntakeToClientResponse, error) {
	var result PromoteIntakeToClientResponse

	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		if _, err := q.LockIntakeFormByID(ctx, req.IntakeFormID); err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "PromoteIntakeToClient", "Failed to lock intake form", zap.Error(err))
			return err
		}

		// 1. Get the intake form
		intakeForm, err := q.GetIntakeForm(ctx, req.IntakeFormID)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "PromoteIntakeToClient", "Failed to get intake form", zap.Error(err))
			return err
		}

		if intakeForm.IntakeConclusion != db.IntakeConclusionEnumSuitable {
			err := fmt.Errorf("%w: only intakes with conclusion 'suitable' can be promoted", ErrIntakeNotSuitable)
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelWarn, "PromoteIntakeToClient", "Intake promotion blocked by conclusion", zap.String("IntakeFormID", req.IntakeFormID.String()), zap.String("IntakeConclusion", string(intakeForm.IntakeConclusion)))
			return err
		}

		// Idempotency: if this intake is already promoted, return existing client.
		existingClient, err := q.GetClientByIntakeFormID(ctx, &req.IntakeFormID)
		if err == nil {
			result.ClientID = existingClient.ID
			result.IntakeFormID = req.IntakeFormID
			result.RegistrationFormID = intakeForm.RegistrationFormID
			result.Message = "Intake already promoted; returning existing client"
			return nil
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "PromoteIntakeToClient", "Failed to check existing client for intake", zap.Error(err))
			return err
		}

		// 2. Get the registration form
		registrationForm, err := q.GetRegistrationForm(ctx, intakeForm.RegistrationFormID)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "PromoteIntakeToClient", "Failed to get registration form", zap.Error(err))
			return err
		}

		// 3. Intake goals are migrated after client creation.

		// 5. Create the client record
		createClientParams := db.CreateClientDetailsParams{
			IntakeFormID:       &req.IntakeFormID,
			RegistrationFormID: &intakeForm.RegistrationFormID,
			FirstName:          registrationForm.ClientFirstName,
			LastName:           registrationForm.ClientLastName,
			DateOfBirth:        registrationForm.ClientDateOfBirth,
			Identity:           false,
			Bsn:                &registrationForm.ClientBsnNumber,
			BsnVerifiedBy:      nil,
			Nationality:        &registrationForm.ClientNationality,
			Email:              registrationForm.ClientEmail,
			PhoneNumber:        &registrationForm.ClientPhoneNumber,
			Gender:             registrationForm.ClientGender,
			CareType: db.NullIntakeCareTypeEnum{
				IntakeCareTypeEnum: intakeForm.CareType,
				Valid:              true,
			},
			SenderID:                   intakeForm.SenderID,
			LocationID:                 intakeForm.AssignedLocationID,
			EvaluationIntervalsWeeks:   intakeForm.EvaluationIntervalsWeeks,
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
		}

		client, err := q.CreateClientDetails(ctx, createClientParams)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				existingClient, getErr := q.GetClientByIntakeFormID(ctx, &req.IntakeFormID)
				if getErr == nil {
					result.ClientID = existingClient.ID
					result.IntakeFormID = req.IntakeFormID
					result.RegistrationFormID = intakeForm.RegistrationFormID
					result.Message = "Intake already promoted; returning existing client"
					return nil
				}
			}
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "PromoteIntakeToClient", "Failed to create client", zap.Error(err))
			return err
		}

		result.ClientID = client.ID
		result.IntakeFormID = req.IntakeFormID
		result.RegistrationFormID = intakeForm.RegistrationFormID

		createdGoals, err := q.CreateClientGoalsFromIntakeAssessments(ctx, db.CreateClientGoalsFromIntakeAssessmentsParams{
			ClientID:     client.ID,
			IntakeFormID: req.IntakeFormID,
		})
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "PromoteIntakeToClient", "Failed to migrate intake goals into client goals", zap.Error(err))
			return err
		}
		result.MaturityAssessmentsCreated = len(createdGoals)

		// 7. Create emergency contacts from guardians
		// Guardian 1
		if registrationForm.Guardian1FirstName != "" {
			relationStatus1 := db.NullRelationStatusEnum{
				RelationStatusEnum: db.RelationStatusEnumPrimaryRelationship,
				Valid:              true,
			}
			_, err := q.CreateEmemrgencyContact(ctx, db.CreateEmemrgencyContactParams{
				ClientID:       client.ID,
				FirstName:      &registrationForm.Guardian1FirstName,
				LastName:       &registrationForm.Guardian1LastName,
				Email:          &registrationForm.Guardian1Email,
				PhoneNumber:    &registrationForm.Guardian1PhoneNumber,
				Relationship:   &registrationForm.Guardian1Relationship,
				RelationStatus: relationStatus1,
			})
			if err != nil {
				s.Logger.LogBusinessEvent(ctx, logger.LogLevelWarn, "PromoteIntakeToClient", "Failed to create guardian 1 emergency contact", zap.Error(err))
			} else {
				result.EmergencyContactsCreated++
			}
		}

		// Guardian 2
		if registrationForm.Guardian2FirstName != "" {
			relationStatus2 := db.NullRelationStatusEnum{
				RelationStatusEnum: db.RelationStatusEnumSecondaryRelationship,
				Valid:              true,
			}
			_, err := q.CreateEmemrgencyContact(ctx, db.CreateEmemrgencyContactParams{
				ClientID:       client.ID,
				FirstName:      &registrationForm.Guardian2FirstName,
				LastName:       &registrationForm.Guardian2LastName,
				Email:          &registrationForm.Guardian2Email,
				PhoneNumber:    &registrationForm.Guardian2PhoneNumber,
				Relationship:   &registrationForm.Guardian2Relationship,
				RelationStatus: relationStatus2,
			})
			if err != nil {
				s.Logger.LogBusinessEvent(ctx, logger.LogLevelWarn, "PromoteIntakeToClient", "Failed to create guardian 2 emergency contact", zap.Error(err))
			} else {
				result.EmergencyContactsCreated++
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if result.Message == "" {
		result.Message = "Intake successfully promoted to client"
	}
	return &result, nil
}
