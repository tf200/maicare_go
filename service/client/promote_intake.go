package clientp

import (
	"context"
	"encoding/json"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

// PromoteIntakeToClient promotes an intake form and its related data to a full client record
func (s *clientService) PromoteIntakeToClient(ctx context.Context, req *PromoteIntakeToClientRequest) (*PromoteIntakeToClientResponse, error) {
	var result PromoteIntakeToClientResponse

	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		// 1. Get the intake form
		intakeForm, err := q.GetIntakeForm(ctx, req.IntakeFormID)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "PromoteIntakeToClient", "Failed to get intake form", zap.Error(err))
			return err
		}

		// 2. Get the registration form
		registrationForm, err := q.GetRegistrationForm(ctx, intakeForm.RegistrationFormID)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "PromoteIntakeToClient", "Failed to get registration form", zap.Error(err))
			return err
		}

		// 3. Get all maturity assessments for this intake
		intakeAssessments, err := q.GetIntakeMaturityAssessments(ctx, req.IntakeFormID)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "PromoteIntakeToClient", "Failed to get maturity assessments", zap.Error(err))
			return err
		}

		// 4. Create addresses JSON from registration form
		addresses := []Address{
			{
				Street:      &registrationForm.ClientStreet,
				HouseNumber: &registrationForm.ClientHouseNumber,
				PostalCode:  &registrationForm.ClientPostalCode,
				City:        &registrationForm.ClientCity,
				PhoneNumber: &registrationForm.ClientPhoneNumber,
			},
		}
		addressesJSON, err := json.Marshal(addresses)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "PromoteIntakeToClient", "Failed to marshal addresses", zap.Error(err))
			return err
		}

		// Map gender from registration to client gender enum
		var clientGender db.ClientGenderEnum
		switch registrationForm.ClientGender {
		case "male":
			clientGender = db.ClientGenderEnumMale
		case "female":
			clientGender = db.ClientGenderEnumFemale
		default:
			clientGender = db.ClientGenderEnumOther
		}

		// Map education level
		var educationLevel db.ClientEducationLevelEnum
		if registrationForm.EducationLevel.Valid {
			educationLevel = registrationForm.EducationLevel.ClientEducationLevelEnum
		} else {
			educationLevel = db.ClientEducationLevelEnumNone
		}

		// Parse work start date
		var workStartDate pgtype.Date
		if registrationForm.WorkStartDate.Valid {
			workStartDate = registrationForm.WorkStartDate
		}

		// 5. Create the client record
		createClientParams := db.CreateClientDetailsParams{
			IntakeFormID:               &req.IntakeFormID,
			FirstName:                  registrationForm.ClientFirstName,
			LastName:                   registrationForm.ClientLastName,
			DateOfBirth:                registrationForm.ClientDateOfBirth,
			Bsn:                        &registrationForm.ClientBsnNumber,
			Source:                     &registrationForm.ReferrerOrganization,
			Nationality:                &registrationForm.ClientNationality,
			Email:                      registrationForm.ClientEmail,
			PhoneNumber:                &registrationForm.ClientPhoneNumber,
			Gender:                     clientGender,
			Filenumber:                 registrationForm.ClientBsnNumber,
			SenderID:                   intakeForm.SenderID,
			LocationID:                 intakeForm.AssignedLocationID,
			Addresses:                  addressesJSON,
			EducationCurrentlyEnrolled: registrationForm.EducationCurrentlyEnrolled,
			EducationInstitution:       registrationForm.EducationInstitution,
			EducationMentorName:        registrationForm.EducationMentorName,
			EducationMentorPhone:       registrationForm.EducationMentorPhone,
			EducationMentorEmail:       registrationForm.EducationMentorEmail,
			EducationAdditionalNotes:   registrationForm.EducationAdditionalNotes,
			EducationLevel:             educationLevel,
			WorkCurrentlyEmployed:      registrationForm.WorkCurrentlyEmployed,
			WorkCurrentEmployer:        registrationForm.WorkCurrentEmployer,
			WorkCurrentEmployerPhone:   registrationForm.WorkEmployerPhone,
			WorkCurrentEmployerEmail:   registrationForm.WorkEmployerEmail,
			WorkCurrentPosition:        registrationForm.WorkCurrentPosition,
			WorkStartDate:              workStartDate,
			WorkAdditionalNotes:        registrationForm.WorkAdditionalNotes,
		}

		client, err := q.CreateClientDetails(ctx, createClientParams)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "PromoteIntakeToClient", "Failed to create client", zap.Error(err))
			return err
		}

		result.ClientID = client.ID
		result.IntakeFormID = req.IntakeFormID
		result.RegistrationFormID = intakeForm.RegistrationFormID

		// 6. Create maturity matrix assessments and care plans from intake assessments
		for _, intakeAssessment := range intakeAssessments {
			// Create client maturity matrix assessment
			endDate := pgtype.Date{Time: intakeForm.DateOfIntake.Time.AddDate(1, 0, 0), Valid: true}

			createAssessmentParams := db.CreateClientMaturityMatrixAssessmentParams{
				ClientID:     client.ID,
				TopicID:      intakeAssessment.TopicID,
				StartDate:    pgtype.Date{Time: intakeForm.DateOfIntake.Time, Valid: true},
				EndDate:      endDate,
				InitialLevel: intakeAssessment.CurrentLevel,
				TargetLevel:  5,
				CurrentLevel: intakeAssessment.CurrentLevel,
			}

			clientAssessment, err := q.CreateClientMaturityMatrixAssessment(ctx, createAssessmentParams)
			if err != nil {
				s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "PromoteIntakeToClient", "Failed to create maturity assessment", zap.Error(err))
				return err
			}
			result.MaturityAssessmentsCreated++

			// Create care plan for this assessment
			createCarePlanParams := db.CreateCarePlanParams{
				AssessmentID:          clientAssessment.ID,
				GeneratedByEmployeeID: nil,
				AssessmentSummary:     fmt.Sprintf("Initial care plan for %s based on intake assessment", intakeAssessment.TopicName),
			}

			carePlan, err := q.CreateCarePlan(ctx, createCarePlanParams)
			if err != nil {
				s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "PromoteIntakeToClient", "Failed to create care plan", zap.Error(err))
				return err
			}

			// Create objectives from goals
			var goals []IntakeAssessmentGoal
			if len(intakeAssessment.ProposedGoals) > 0 {
				if err := json.Unmarshal(intakeAssessment.ProposedGoals, &goals); err != nil {
					s.Logger.LogBusinessEvent(ctx, logger.LogLevelWarn, "PromoteIntakeToClient", "Failed to unmarshal goals, continuing without objectives", zap.Error(err))
				} else {
					for _, goal := range goals {
						createObjectiveParams := db.CreateCarePlanObjectiveParams{
							CarePlanID:  carePlan.ID,
							Timeframe:   db.CarePlanTimeframeEnumShortTerm,
							GoalTitle:   goal.Title,
							Description: goal.Description,
							Priority:    goal.Priority,
						}

						_, err := q.CreateCarePlanObjective(ctx, createObjectiveParams)
						if err != nil {
							s.Logger.LogBusinessEvent(ctx, logger.LogLevelWarn, "PromoteIntakeToClient", "Failed to create care plan objective", zap.Error(err))
						}
					}
				}
			}
		}

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

		// 8. Link documents from registration to client
		documentFields := []struct {
			attachmentID *uuid.UUID
			label        db.ClientDocumentLabelEnum
		}{
			{registrationForm.DocumentReferral, db.ClientDocumentLabelEnumRegistrationForm},
			{registrationForm.DocumentEducationReport, db.ClientDocumentLabelEnumOther},
			{registrationForm.DocumentActionPlan, db.ClientDocumentLabelEnumOther},
			{registrationForm.DocumentPsychiatricReport, db.ClientDocumentLabelEnumOther},
			{registrationForm.DocumentDiagnosis, db.ClientDocumentLabelEnumOther},
			{registrationForm.DocumentSafetyPlan, db.ClientDocumentLabelEnumOther},
			{registrationForm.DocumentIDCopy, db.ClientDocumentLabelEnumOther},
		}

		for _, doc := range documentFields {
			if doc.attachmentID != nil {
				_, err := q.CreateClientDocument(ctx, db.CreateClientDocumentParams{
					ClientID:       client.ID,
					AttachmentUuid: doc.attachmentID,
					Label:          doc.label,
				})
				if err != nil {
					s.Logger.LogBusinessEvent(ctx, logger.LogLevelWarn, "PromoteIntakeToClient", "Failed to link document to client", zap.Error(err))
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	result.Message = "Intake successfully promoted to client"
	return &result, nil
}
