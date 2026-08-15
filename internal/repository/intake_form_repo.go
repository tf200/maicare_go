package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	db "maicare_go/db/sqlc"
	"maicare_go/internal/domain"
)

type IntakeFormRepository struct {
	store *db.Store
}

func NewIntakeFormRepository(store *db.Store) domain.IntakeFormRepository {
	return &IntakeFormRepository{store: store}
}

func (r *IntakeFormRepository) ExecTx(ctx context.Context, fn func(*db.Queries) error) error {
	return r.store.ExecTx(ctx, fn)
}

func (r *IntakeFormRepository) CreateIntakeForm(ctx context.Context, params domain.CreateIntakeFormParams) (*domain.IntakeForm, error) {
	// Validate registration form exists and status is processed
	regForm, err := r.store.GetRegistrationForm(ctx, params.RegistrationFormID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("registration form not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get registration form: %w", err)
	}

	if regForm.FormStatus != db.FormStatusEnumProcessed {
		return nil, fmt.Errorf("registration form must be processed to create intake form")
	}

	// Check if intake form already exists for this registration form
	_, err = r.store.GetIntakeFormByRegistrationFormID(ctx, params.RegistrationFormID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to check existing intake form: %w", err)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("intake form already exists for this registration form")
	}

	dbParams := db.CreateIntakeFormParams{
		RegistrationFormID:       params.RegistrationFormID,
		DateOfIntake:             pgtype.Timestamptz{Time: params.DateOfIntake, Valid: true},
		CareType:                 params.CareType,
		IntakeParticipants:       params.IntakeParticipants,
		FamilySituation:          params.FamilySituation,
		PsychologicalState:       params.PsychologicalState,
		SelfSufficiency:          params.SelfSufficiency,
		SenderID:                 params.SenderID,
		AssignedLocationID:       params.AssignedLocationID,
		RiskAssessment:           params.RiskAssessment,
		IntakeConclusion:         params.IntakeConclusion,
		IntakeConclusionNotes:    params.IntakeConclusionNotes,
		EvaluationIntervalsWeeks: params.EvaluationIntervalsWeeks,
		Signature:                params.Signature,
	}

	intakeForm, err := r.store.CreateIntakeForm(ctx, dbParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create intake form: %w", err)
	}

	return toDomainIntakeForm(intakeForm), nil
}

func (r *IntakeFormRepository) ListIntakeForms(ctx context.Context, params domain.ListIntakeFormsParams) ([]domain.IntakeFormListItem, int64, error) {
	var status *db.IntakeConclusionEnum
	if params.Status != nil {
		status = params.Status
	}

	sortBy := ""
	sortOrder := ""
	if params.SortOrder != nil {
		parts := strings.SplitN(*params.SortOrder, " ", 2)
		if len(parts) == 2 {
			sortByStr := strings.ToLower(parts[0])
			sortOrderStr := strings.ToLower(parts[1])
			if sortByStr == "created_at" && (sortOrderStr == "asc" || sortOrderStr == "desc") {
				sortBy = sortByStr
				sortOrder = sortOrderStr
			}
		}
	}

	search := ""
	if params.Search != nil {
		search = *params.Search
	}

	dbParams := db.ListIntakeFormsParams{
		Limit:     params.Limit,
		Offset:    params.Offset,
		Search:    search,
		Status:    status,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	rows, err := r.store.ListIntakeForms(ctx, dbParams)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list intake forms: %w", err)
	}

	if len(rows) == 0 {
		return []domain.IntakeFormListItem{}, 0, nil
	}

	items := make([]domain.IntakeFormListItem, len(rows))
	for i, row := range rows {
		items[i] = toDomainIntakeFormListItem(row)
	}

	totalCount := rows[0].TotalCount
	return items, totalCount, nil
}

func (r *IntakeFormRepository) GetIntakeFormTotals(ctx context.Context) (*domain.IntakeFormTotals, error) {
	row, err := r.store.GetIntakeFormTotals(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get intake form totals: %w", err)
	}

	return &domain.IntakeFormTotals{
		FurtherInvestigationTotal: row.FurtherInvestigationTotal,
		WithoutGoalsTotal:         row.WithoutGoalsTotal,
	}, nil
}

func (r *IntakeFormRepository) GetIntakeFormDetail(ctx context.Context, id uuid.UUID) (*domain.IntakeFormDetail, error) {
	row, err := r.store.GetIntakeFormDetails(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrIntakeFormNotFound
		}
		return nil, fmt.Errorf("failed to get intake form detail: %w", err)
	}

	// Get topic assessments
	assessmentsRows, err := r.store.GetIntakeTopicsAssessments(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get intake topic assessments: %w", err)
	}

	intakeGoalsAssigned := make([]domain.IntakeGoalTopic, 0, len(assessmentsRows))
	for _, assessment := range assessmentsRows {
		var goals []domain.IntakeAssessmentGoal
		if len(assessment.ProposedGoals) > 0 {
			if err := json.Unmarshal(assessment.ProposedGoals, &goals); err != nil {
				goals = []domain.IntakeAssessmentGoal{}
			}
		}
		if goals == nil {
			goals = []domain.IntakeAssessmentGoal{}
		}

		intakeGoalsAssigned = append(intakeGoalsAssigned, domain.IntakeGoalTopic{
			AssessmentID:  assessment.ID,
			TopicID:       assessment.TopicID,
			TopicName:     assessment.TopicName,
			CurrentLevel:  assessment.CurrentLevel,
			ProposedGoals: goals,
			Notes:         assessment.Notes,
		})
	}

	// Build location details
	var location *domain.IntakeFormLocationDetails
	if row.LocationName != nil || row.LocationStreet != nil || row.LocationCity != nil {
		location = &domain.IntakeFormLocationDetails{
			Name:                row.LocationName,
			Street:              row.LocationStreet,
			HouseNumber:         row.LocationHouseNumber,
			HouseNumberAddition: row.LocationHouseNumberAddition,
			PostalCode:          row.LocationPostalCode,
			City:                row.LocationCity,
		}
	}

	return &domain.IntakeFormDetail{
		ID:                       row.ID,
		RegistrationFormID:       row.RegistrationFormID,
		DateOfIntake:             row.DateOfIntake.Time,
		CareType:                 row.CareType,
		IntakeParticipants:       row.IntakeParticipants,
		FamilySituation:          row.FamilySituation,
		PsychologicalState:       row.PsychologicalState,
		SelfSufficiency:          row.SelfSufficiency,
		SenderID:                 row.SenderID,
		AssignedLocationID:       row.AssignedLocationID,
		RiskAssessment:           row.RiskAssessment,
		IntakeConclusion:         row.IntakeConclusion,
		IntakeConclusionNotes:    row.IntakeConclusionNotes,
		EvaluationIntervalsWeeks: row.EvaluationIntervalsWeeks,
		Signature:                row.Signature,
		CreatedAt:                row.CreatedAt.Time,
		UpdatedAt:                row.UpdatedAt.Time,
		ClientFirstName:          row.ClientFirstName,
		ClientLastName:           row.ClientLastName,
		ClientBsnNumber:          row.ClientBsnNumber,
		DesiredGoals:             row.ClientGoals,
		SenderName:               row.SenderName,
		Location:                 location,
		IntakeGoalsAssigned:      intakeGoalsAssigned,
		HasClient:                row.HasClient,
	}, nil
}

func (r *IntakeFormRepository) UpdateIntakeForm(ctx context.Context, params domain.UpdateIntakeFormParams) (*domain.IntakeForm, error) {
	dbParams := db.UpdateIntakeFormParams{
		ID:                       params.ID,
		DateOfIntake:             pgtype.Timestamptz{Time: *params.DateOfIntake, Valid: params.DateOfIntake != nil},
		CareType:                 db.NullIntakeCareTypeFromPtr((*string)(params.CareType)),
		IntakeParticipants:       *params.IntakeParticipants,
		SelfSufficiency:          params.SelfSufficiency,
		SenderID:                 params.SenderID,
		AssignedLocationID:       params.AssignedLocationID,
		RiskAssessment:           params.RiskAssessment,
		EvaluationIntervalsWeeks: params.EvaluationIntervalsWeeks,
		Signature:                params.Signature,
		FamilySituation:          params.FamilySituation,
		PsychologicalState:       params.PsychologicalState,
	}

	// Handle clear fields
	for _, field := range params.ClearFields {
		switch strings.ToLower(field) {
		case "family_situation":
			dbParams.ClearFamilySituation = true
		case "psychological_state":
			dbParams.ClearPsychologicalState = true
		case "sender_id":
			dbParams.ClearSenderID = true
		case "assigned_location_id":
			dbParams.ClearAssignedLocationID = true
		case "risk_assessment":
			dbParams.ClearRiskAssessment = true
		case "signature":
			dbParams.ClearSignature = true
		default:
			return nil, domain.ErrInvalidIntakeFormClearField
		}
	}

	intakeForm, err := r.store.UpdateIntakeForm(ctx, dbParams)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Check if there's an active client blocking the update
			hasActiveClient, checkErr := r.store.HasActiveClientByIntakeFormID(ctx, &params.ID)
			if checkErr == nil && hasActiveClient {
				return nil, domain.ErrIntakeFormUpdateBlockedByActiveClient
			}
			// Check if the intake form exists at all
			_, existsErr := r.store.GetIntakeForm(ctx, params.ID)
			if existsErr != nil {
				if errors.Is(existsErr, pgx.ErrNoRows) {
					return nil, err // Return original error for not found
				}
				return nil, existsErr
			}
			return nil, domain.ErrIntakeFormUpdateConflict
		}
		return nil, fmt.Errorf("failed to update intake form: %w", err)
	}

	return toDomainIntakeForm(intakeForm), nil
}

func (r *IntakeFormRepository) UpdateIntakeConclusion(ctx context.Context, id uuid.UUID, params domain.UpdateIntakeConclusionParams) (*domain.IntakeFormConclusion, error) {
	// Validate decision
	var conclusion db.IntakeConclusionEnum
	switch strings.ToLower(params.Decision) {
	case "accept":
		conclusion = db.IntakeConclusionEnumSuitable
	case "refuse":
		conclusion = db.IntakeConclusionEnumUnsuitable
	default:
		return nil, domain.ErrInvalidIntakeFormUpdate
	}

	dbParams := db.UpdateIntakeConclusionParams{
		ID:                    id,
		IntakeConclusion:      conclusion,
		IntakeConclusionNotes: params.IntakeConclusionNotes,
	}

	intakeForm, err := r.store.UpdateIntakeConclusion(ctx, dbParams)
	if err != nil {
		return nil, fmt.Errorf("failed to update intake conclusion: %w", err)
	}

	return &domain.IntakeFormConclusion{
		ID:                    intakeForm.ID,
		IntakeConclusion:      intakeForm.IntakeConclusion,
		IntakeConclusionNotes: intakeForm.IntakeConclusionNotes,
		UpdatedAt:             intakeForm.UpdatedAt.Time,
	}, nil
}

func (r *IntakeFormRepository) ReplaceIntakeFormGoals(ctx context.Context, intakeFormID uuid.UUID, params domain.ReplaceIntakeFormGoalsParams) (*domain.IntakeFormGoalsResult, error) {
	var result *domain.IntakeFormGoalsResult

	err := r.store.ExecTx(ctx, func(q *db.Queries) error {
		// Step 1: Lock intake form
		_, err := q.LockIntakeFormByID(ctx, intakeFormID)
		if err != nil {
			return fmt.Errorf("failed to lock intake form: %w", err)
		}

		// Step 2: Check for active client
		hasActiveClient, err := q.HasActiveClientByIntakeFormID(ctx, &intakeFormID)
		if err != nil {
			return fmt.Errorf("failed to check active client: %w", err)
		}
		if hasActiveClient {
			return domain.ErrIntakeGoalsUpdateBlockedByActiveClient
		}

		// Step 3: Delete existing assessments
		if err := q.DeleteIntakeTopicAssessmentsByIntakeForm(ctx, intakeFormID); err != nil {
			return fmt.Errorf("failed to delete existing assessments: %w", err)
		}

		// Step 4: If assessments not empty, create new ones
		if len(params.Assessments) > 0 {
			// Marshal assessments to JSON
			assessmentsJSON, err := json.Marshal(params.Assessments)
			if err != nil {
				return fmt.Errorf("failed to marshal assessments: %w", err)
			}

			batchParams := db.CreateIntakeTopicAssessmentsBatchParams{
				IntakeFormID: intakeFormID,
				Items:        assessmentsJSON,
			}

			rows, err := q.CreateIntakeTopicAssessmentsBatch(ctx, batchParams)
			if err != nil {
				return fmt.Errorf("failed to create topic assessments batch: %w", err)
			}

			// Parse returned rows
			assessments := make([]domain.IntakeMaturityAssessment, 0, len(rows))
			for _, row := range rows {
				var goals []domain.IntakeAssessmentGoal
				if len(row.ProposedGoals) > 0 {
					if err := json.Unmarshal(row.ProposedGoals, &goals); err != nil {
						goals = []domain.IntakeAssessmentGoal{}
					}
				}
				if goals == nil {
					goals = []domain.IntakeAssessmentGoal{}
				}

				assessments = append(assessments, domain.IntakeMaturityAssessment{
					ID:            row.ID,
					IntakeFormID:  row.IntakeFormID,
					TopicID:       row.TopicID,
					TopicName:     row.TopicName,
					CurrentLevel:  row.CurrentLevel,
					ProposedGoals: goals,
					Notes:         row.Notes,
					CreatedAt:     row.CreatedAt.Time,
				})
			}

			result = &domain.IntakeFormGoalsResult{
				Assessments: assessments,
			}
		} else {
			result = &domain.IntakeFormGoalsResult{
				Assessments: []domain.IntakeMaturityAssessment{},
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *IntakeFormRepository) PromoteIntakeToClient(ctx context.Context, params domain.PromoteIntakeToClientParams) (*domain.PromoteIntakeToClientResult, error) {
	var result *domain.PromoteIntakeToClientResult

	err := r.store.ExecTx(ctx, func(q *db.Queries) error {
		// Step 1: Lock intake form
		_, err := q.LockIntakeFormByID(ctx, params.IntakeFormID)
		if err != nil {
			return fmt.Errorf("failed to lock intake form: %w", err)
		}

		// Step 2: Get intake form and validate conclusion
		intakeForm, err := q.GetIntakeForm(ctx, params.IntakeFormID)
		if err != nil {
			return fmt.Errorf("failed to get intake form: %w", err)
		}
		if intakeForm.IntakeConclusion != db.IntakeConclusionEnumSuitable {
			return domain.ErrIntakeNotSuitable
		}

		// Step 3: Check if client already exists
		clientDetail, err := q.GetClientByIntakeFormID(ctx, &params.IntakeFormID)
		if err == nil && clientDetail.ID != uuid.Nil {
			// Client already exists, return existing result
			result = &domain.PromoteIntakeToClientResult{
				ClientID:           clientDetail.ID,
				IntakeFormID:       params.IntakeFormID,
				RegistrationFormID: intakeForm.RegistrationFormID,
				Message:            "Client already exists",
			}
			return nil
		}

		// Step 4: Get registration form
		regForm, err := q.GetRegistrationForm(ctx, intakeForm.RegistrationFormID)
		if err != nil {
			return fmt.Errorf("failed to get registration form: %w", err)
		}

		// Step 5: Create client details
		createdClient, err := q.CreateClientDetails(ctx, db.CreateClientDetailsParams{
			IntakeFormID:               &params.IntakeFormID,
			RegistrationFormID:         &intakeForm.RegistrationFormID,
			FirstName:                  regForm.ClientFirstName,
			LastName:                   regForm.ClientLastName,
			DateOfBirth:                regForm.ClientDateOfBirth,
			Identity:                   true,
			Bsn:                        &regForm.ClientBsnNumber,
			BsnVerifiedBy:              &params.EmployeeID,
			Gender:                     regForm.ClientGender,
			Email:                      regForm.ClientEmail,
			PhoneNumber:                &regForm.ClientPhoneNumber,
			CareType:                   &intakeForm.CareType,
			SenderID:                   intakeForm.SenderID,
			LocationID:                 intakeForm.AssignedLocationID,
			Street:                     regForm.ClientStreet,
			HouseNumber:                regForm.ClientHouseNumber,
			HouseNumberAddition:        regForm.ClientHouseNumberAddition,
			PostalCode:                 regForm.ClientPostalCode,
			City:                       regForm.ClientCity,
			EducationCurrentlyEnrolled: regForm.EducationCurrentlyEnrolled,
			EducationInstitution:       regForm.EducationInstitution,
			EducationMentorName:        regForm.EducationMentorName,
			EducationMentorPhone:       regForm.EducationMentorPhone,
			EducationMentorEmail:       regForm.EducationMentorEmail,
			EducationAdditionalNotes:   regForm.EducationAdditionalNotes,
			EducationLevel:             regForm.EducationLevel,
			WorkCurrentlyEmployed:      regForm.WorkCurrentlyEmployed,
			WorkCurrentEmployer:        regForm.WorkCurrentEmployer,
			WorkCurrentEmployerPhone:   regForm.WorkEmployerPhone,
			WorkCurrentEmployerEmail:   regForm.WorkEmployerEmail,
			WorkCurrentPosition:        regForm.WorkCurrentPosition,
			WorkStartDate:              regForm.WorkStartDate,
			WorkAdditionalNotes:        regForm.WorkAdditionalNotes,
			Nationality:                &regForm.ClientNationality,
			EvaluationIntervalsWeeks:   intakeForm.EvaluationIntervalsWeeks,
		})

		if err != nil {
			// Check for unique violation - on 23505, retry fetching client
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				clientDetail, retryErr := q.GetClientByIntakeFormID(ctx, &params.IntakeFormID)
				if retryErr != nil {
					return fmt.Errorf("unique constraint violation and failed to retry: %w", retryErr)
				}
				createdClient = clientDetail
			} else {
				return fmt.Errorf("failed to create client details: %w", err)
			}
		}

		// Step 6: Create client goals from intake assessments
		goalRows, err := q.CreateClientGoalsFromIntakeAssessments(ctx, db.CreateClientGoalsFromIntakeAssessmentsParams{
			IntakeFormID: params.IntakeFormID,
			ClientID:     createdClient.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to create client goals: %w", err)
		}

		maturityAssessmentsCreated := len(goalRows)

		// Step 7: Create emergency contacts from guardians
		emergencyContactsCreated := 0

		// Guardian 1
		if regForm.Guardian1FirstName != nil && *regForm.Guardian1FirstName != "" {
			_, err := q.CreateEmemrgencyContact(ctx, db.CreateEmemrgencyContactParams{
				ClientID:         createdClient.ID,
				FirstName:        regForm.Guardian1FirstName,
				LastName:         regForm.Guardian1LastName,
				Relationship:     regForm.Guardian1Relationship,
				PhoneNumber:      regForm.Guardian1PhoneNumber,
				Email:            regForm.Guardian1Email,
				MedicalReports:   true,
				IncidentsReports: true,
				GoalsReports:     true,
			})
			if err == nil {
				emergencyContactsCreated++
			}
		}

		// Guardian 2
		if regForm.Guardian2FirstName != nil && *regForm.Guardian2FirstName != "" {
			_, err := q.CreateEmemrgencyContact(ctx, db.CreateEmemrgencyContactParams{
				ClientID:         createdClient.ID,
				FirstName:        regForm.Guardian2FirstName,
				LastName:         regForm.Guardian2LastName,
				Relationship:     regForm.Guardian2Relationship,
				PhoneNumber:      regForm.Guardian2PhoneNumber,
				Email:            regForm.Guardian2Email,
				MedicalReports:   true,
				IncidentsReports: true,
				GoalsReports:     true,
			})
			if err == nil {
				emergencyContactsCreated++
			}
		}

		result = &domain.PromoteIntakeToClientResult{
			ClientID:                   createdClient.ID,
			IntakeFormID:               params.IntakeFormID,
			RegistrationFormID:         intakeForm.RegistrationFormID,
			Message:                    "Intake promoted to client successfully",
			MaturityAssessmentsCreated: maturityAssessmentsCreated,
			EmergencyContactsCreated:   emergencyContactsCreated,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *IntakeFormRepository) GetIntakeForm(ctx context.Context, id uuid.UUID) (db.IntakeForm, error) {
	return r.store.GetIntakeForm(ctx, id)
}

func (r *IntakeFormRepository) GetIntakeFormByRegistrationFormID(ctx context.Context, registrationFormID uuid.UUID) (db.IntakeForm, error) {
	return r.store.GetIntakeFormByRegistrationFormID(ctx, registrationFormID)
}

func (r *IntakeFormRepository) GetRegistrationForm(ctx context.Context, id uuid.UUID) (db.GetRegistrationFormRow, error) {
	return r.store.GetRegistrationForm(ctx, id)
}

func (r *IntakeFormRepository) GetIntakeMaturityAssessment(ctx context.Context, id uuid.UUID) (db.GetIntakeMaturityAssessmentRow, error) {
	return r.store.GetIntakeMaturityAssessment(ctx, id)
}

func (r *IntakeFormRepository) GetTopicLevel(ctx context.Context, params db.GetTopicLevelParams) (db.GetTopicLevelRow, error) {
	return r.store.GetTopicLevel(ctx, params)
}

func (r *IntakeFormRepository) HasActiveClientByIntakeFormID(ctx context.Context, intakeFormID *uuid.UUID) (bool, error) {
	return r.store.HasActiveClientByIntakeFormID(ctx, intakeFormID)
}

func (r *IntakeFormRepository) LockIntakeFormByID(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	return r.store.LockIntakeFormByID(ctx, id)
}

func (r *IntakeFormRepository) GetClientByIntakeFormID(ctx context.Context, intakeFormID *uuid.UUID) (db.ClientDetail, error) {
	return r.store.GetClientByIntakeFormID(ctx, intakeFormID)
}

func (r *IntakeFormRepository) CreateClientDetails(ctx context.Context, params db.CreateClientDetailsParams) (db.ClientDetail, error) {
	return r.store.CreateClientDetails(ctx, params)
}

func (r *IntakeFormRepository) CreateClientGoalsFromIntakeAssessments(ctx context.Context, params db.CreateClientGoalsFromIntakeAssessmentsParams) ([]db.CreateClientGoalsFromIntakeAssessmentsRow, error) {
	return r.store.CreateClientGoalsFromIntakeAssessments(ctx, params)
}

func (r *IntakeFormRepository) CreateEmergencyContact(ctx context.Context, params db.CreateEmemrgencyContactParams) (db.ClientEmergencyContact, error) {
	return r.store.CreateEmemrgencyContact(ctx, params)
}

func (r *IntakeFormRepository) CreateIntakeTopicAssessmentsBatch(ctx context.Context, params db.CreateIntakeTopicAssessmentsBatchParams) ([]db.CreateIntakeTopicAssessmentsBatchRow, error) {
	return r.store.CreateIntakeTopicAssessmentsBatch(ctx, params)
}

func (r *IntakeFormRepository) DeleteIntakeTopicAssessmentsByIntakeForm(ctx context.Context, intakeFormID uuid.UUID) error {
	return r.store.DeleteIntakeTopicAssessmentsByIntakeForm(ctx, intakeFormID)
}

func (r *IntakeFormRepository) DeleteIntakeForm(ctx context.Context, id uuid.UUID) error {
	return r.store.ExecTx(ctx, func(q *db.Queries) error {
		_, err := q.LockIntakeFormByID(ctx, id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return domain.ErrIntakeFormNotFound
			}
			return fmt.Errorf("failed to lock intake form: %w", err)
		}

		hasClient, err := q.HasActiveClientByIntakeFormID(ctx, &id)
		if err != nil {
			return fmt.Errorf("failed to check client for intake form: %w", err)
		}
		if hasClient {
			return domain.ErrIntakeFormDeleteBlockedByActiveClient
		}

		if err := q.DeleteIntakeForm(ctx, id); err != nil {
			return fmt.Errorf("failed to delete intake form: %w", err)
		}

		return nil
	})
}

// Helper functions

func toDomainIntakeForm(row db.IntakeForm) *domain.IntakeForm {
	return &domain.IntakeForm{
		ID:                       row.ID,
		RegistrationFormID:       row.RegistrationFormID,
		DateOfIntake:             row.DateOfIntake.Time,
		CareType:                 row.CareType,
		IntakeParticipants:       row.IntakeParticipants,
		FamilySituation:          row.FamilySituation,
		PsychologicalState:       row.PsychologicalState,
		SelfSufficiency:          row.SelfSufficiency,
		SenderID:                 row.SenderID,
		AssignedLocationID:       row.AssignedLocationID,
		RiskAssessment:           row.RiskAssessment,
		IntakeConclusion:         row.IntakeConclusion,
		IntakeConclusionNotes:    row.IntakeConclusionNotes,
		EvaluationIntervalsWeeks: row.EvaluationIntervalsWeeks,
		Signature:                row.Signature,
		UpdatedAt:                row.UpdatedAt.Time,
	}
}

func toDomainIntakeFormListItem(row db.ListIntakeFormsRow) domain.IntakeFormListItem {
	var address *domain.AssignedLocationAddress
	if row.AssignedLocationStreet != nil || row.AssignedLocationCity != nil {
		address = &domain.AssignedLocationAddress{
			Street:              row.AssignedLocationStreet,
			HouseNumber:         row.AssignedLocationHouseNumber,
			HouseNumberAddition: row.AssignedLocationHouseNumberAddition,
			PostalCode:          row.AssignedLocationPostalCode,
			City:                row.AssignedLocationCity,
		}
	}

	return domain.IntakeFormListItem{
		ID:                      row.ID,
		RegistrationFormID:      row.RegistrationFormID,
		DateOfIntake:            row.DateOfIntake.Time,
		ClientFirstName:         row.ClientFirstName,
		ClientLastName:          row.ClientLastName,
		ClientBsnNumber:         row.ClientBsnNumber,
		IntakeStatus:            row.IntakeConclusion,
		GoalAssessmentDone:      row.GoalAssessmentDone,
		HasClient:               row.HasClient,
		CareType:                row.CareType,
		AssignedLocationID:      row.AssignedLocationID,
		AssignedLocationAddress: address,
	}
}
