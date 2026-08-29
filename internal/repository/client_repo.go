package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/internal/domain"
	"maicare_go/pkg/conv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type ClientRepository struct {
	store *db.Store
}

func NewClientRepository(store *db.Store) domain.ClientRepository {
	return &ClientRepository{store: store}
}

func (r *ClientRepository) CreateClient(ctx context.Context, params domain.CreateClientParams) (*domain.Client, error) {
	client, err := actorQuery(ctx, r.store, func(q *db.Queries) (db.ClientDetail, error) {
		return q.CreateClientDetails(ctx, db.CreateClientDetailsParams{
			FirstName:                  params.FirstName,
			LastName:                   params.LastName,
			DateOfBirth:                conv.PgDateFromTime(params.DateOfBirth),
			Identity:                   true,
			Bsn:                        params.Bsn,
			BsnVerifiedBy:              params.BsnVerifiedBy,
			Email:                      params.Email,
			PhoneNumber:                params.PhoneNumber,
			CareType:                   db.NullIntakeCareTypeFromPtr(params.CareType),
			SenderID:                   params.SenderID,
			LocationID:                 params.LocationID,
			EducationCurrentlyEnrolled: params.EducationCurrentlyEnrolled,
			EducationInstitution:       params.EducationInstitution,
			EducationMentorName:        params.EducationMentorName,
			EducationMentorPhone:       params.EducationMentorPhone,
			EducationMentorEmail:       params.EducationMentorEmail,
			EducationAdditionalNotes:   params.EducationAdditionalNotes,
			WorkCurrentlyEmployed:      params.WorkCurrentlyEmployed,
			WorkCurrentEmployer:        params.WorkCurrentEmployer,
			WorkCurrentEmployerPhone:   params.WorkCurrentEmployerPhone,
			WorkCurrentEmployerEmail:   params.WorkCurrentEmployerEmail,
			WorkCurrentPosition:        params.WorkCurrentPosition,
			WorkStartDate:              conv.PgDateFromTime(params.WorkStartDate),
			WorkAdditionalNotes:        params.WorkAdditionalNotes,
		})
	})
	if err != nil {
		return nil, err
	}

	return toDomainClient(client), nil
}

func (r *ClientRepository) ListClients(ctx context.Context, params domain.ListClientsParams) (*domain.ClientPage, error) {
	rows, err := actorQuery(ctx, r.store, func(q *db.Queries) ([]db.ListClientDetailsRow, error) {
		return q.ListClientDetails(ctx, db.ListClientDetailsParams{
			Status:     db.NullClientStatusFromPtr(params.Status),
			LocationID: params.LocationID,
			Search:     params.Search,
			Offset:     params.Offset,
			Limit:      params.Limit,
		})
	})
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return &domain.ClientPage{Items: []domain.ClientListItem{}, TotalCount: 0}, nil
	}

	items := make([]domain.ClientListItem, len(rows))
	for i, row := range rows {
		items[i] = toDomainClientListItem(row)
	}

	return &domain.ClientPage{
		Items:      items,
		TotalCount: rows[0].TotalCount,
	}, nil
}

func (r *ClientRepository) ListWaitingListClients(ctx context.Context, params domain.ListWaitingListClientsParams) (*domain.WaitingListClientPage, error) {
	rows, err := actorQuery(ctx, r.store, func(q *db.Queries) ([]db.ListWaitingListClientsRow, error) {
		return q.ListWaitingListClients(ctx, db.ListWaitingListClientsParams{
			Search:    params.Search,
			Placement: db.NullIntakeCareTypeFromPtr(params.Placement),
			SortDays:  params.SortDays,
			Offset:    params.Offset,
			Limit:     params.Limit,
		})
	})
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return &domain.WaitingListClientPage{Items: []domain.WaitingListClient{}, TotalCount: 0}, nil
	}

	items := make([]domain.WaitingListClient, len(rows))
	for i, row := range rows {
		items[i] = toDomainWaitingListClient(row)
	}

	return &domain.WaitingListClientPage{
		Items:      items,
		TotalCount: rows[0].TotalCount,
	}, nil
}

func (r *ClientRepository) ListInCareClients(ctx context.Context, params domain.ListInCareClientsParams) (*domain.InCareClientPage, error) {
	var statusFilters []db.ClientStatusEnum
	if len(params.Status) > 0 {
		statusFilters = make([]db.ClientStatusEnum, len(params.Status))
		for i, status := range params.Status {
			statusFilters[i] = db.ClientStatusEnum(status)
		}
	}

	rows, err := actorQuery(ctx, r.store, func(q *db.Queries) ([]db.ListInCareClientsRow, error) {
		return q.ListInCareClients(ctx, db.ListInCareClientsParams{
			Search:         params.Search,
			Status:         statusFilters,
			SortDaysInCare: params.SortDaysInCare,
			Offset:         params.Offset,
			Limit:          params.Limit,
		})
	})
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return &domain.InCareClientPage{Items: []domain.InCareClient{}, TotalCount: 0}, nil
	}

	items := make([]domain.InCareClient, len(rows))
	for i, row := range rows {
		items[i] = toDomainInCareClient(row)
	}

	return &domain.InCareClientPage{
		Items:      items,
		TotalCount: rows[0].TotalCount,
	}, nil
}

func (r *ClientRepository) GetClientCounts(ctx context.Context) (*domain.ClientCounts, error) {
	counts, err := actorQuery(ctx, r.store, func(q *db.Queries) (db.GetClientCountsRow, error) {
		return q.GetClientCounts(ctx)
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &domain.ClientCounts{}, nil
		}
		return nil, err
	}

	return &domain.ClientCounts{
		TotalClients:         counts.TotalClients,
		ClientsInCare:        counts.ClientsInCare,
		ClientsOnWaitingList: counts.ClientsOnWaitingList,
		ClientsOutOfCare:     counts.ClientsOutOfCare,
	}, nil
}

func (r *ClientRepository) GetClientStatusCounts(ctx context.Context) (*domain.ClientStatusCounts, error) {
	counts, err := actorQuery(ctx, r.store, func(q *db.Queries) (db.GetClientStatusCountsRow, error) {
		return q.GetClientStatusCounts(ctx)
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &domain.ClientStatusCounts{}, nil
		}
		return nil, err
	}

	return &domain.ClientStatusCounts{
		ClientsInOrScheduledInCare:     counts.ClientsInOrScheduledInCare,
		ClientsOnWaitingList:           counts.ClientsOnWaitingList,
		ClientsOutOrScheduledOutOfCare: counts.ClientsOutOrScheduledOutOfCare,
	}, nil
}

func (r *ClientRepository) GetInCareStats(ctx context.Context) (*domain.InCareStats, error) {
	query := `
		SELECT
			COUNT(*) FILTER (WHERE c.status = 'in_care') AS clients_in_care,
			COUNT(*) FILTER (WHERE c.status = 'scheduled_in_care') AS clients_scheduled_in_care,
			COUNT(*) FILTER (
				WHERE c.status IN ('in_care', 'scheduled_in_care')
				  AND EXISTS (
					SELECT 1 FROM contract ct
					WHERE ct.client_id = c.id
					  AND ct.status = 'approved'
					  AND ct.start_date <= CURRENT_TIMESTAMP
					  AND ct.end_date >= CURRENT_TIMESTAMP
					  AND ct.end_date <= CURRENT_TIMESTAMP + INTERVAL '30 days'
				  )
			) AS contracts_ending_soon,
			COUNT(*) FILTER (WHERE c.status IN ('in_care', 'scheduled_in_care')) AS total
		FROM client_details c
	`

	var stats domain.InCareStats
	tx, err := r.store.BeginActorTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	err = tx.QueryRow(ctx, query).Scan(
		&stats.ClientsInCare,
		&stats.ClientsScheduledInCare,
		&stats.ContractsEndingSoon,
		&stats.Total,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &domain.InCareStats{}, nil
		}
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *ClientRepository) GetWaitingListStats(ctx context.Context) (*domain.WaitingListStats, error) {
	query := `
		SELECT
			COUNT(*) AS total_clients,
			COUNT(*) FILTER (WHERE rf.addmission_type = 'crisis_admission') AS total_crisis,
			COUNT(*) FILTER (WHERE rf.addmission_type = 'regular_placement') AS total_regular,
			COALESCE(ROUND(AVG(CURRENT_DATE - c.created_at::date), 1), 0)::float8 AS avg_days_in_waitlist
		FROM client_details c
		LEFT JOIN registration_form rf ON c.registration_form_id = rf.id
		WHERE c.status = 'on_waiting_list'
	`

	var stats domain.WaitingListStats
	tx, err := r.store.BeginActorTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	err = tx.QueryRow(ctx, query).Scan(
		&stats.TotalClients,
		&stats.TotalCrisis,
		&stats.TotalRegular,
		&stats.AvgDaysInWaitlist,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &domain.WaitingListStats{}, nil
		}
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *ClientRepository) GetClientByID(ctx context.Context, id uuid.UUID) (*domain.ClientPageDetail, error) {
	tx, err := r.store.BeginActorTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	q := r.store.WithTx(tx)

	client, err := q.GetClientDetails(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrClientNotFound
		}
		return nil, err
	}

	goals, err := q.ListActiveGoalSummariesByClientID(ctx, id)
	if err != nil {
		return nil, err
	}

	emergencyContacts, err := q.ListTopEmergencyContactsByClientID(ctx, id)
	if err != nil {
		return nil, err
	}

	existingDocumentLabels, err := q.ListExistingClientDocumentLabels(ctx, id)
	if err != nil {
		return nil, err
	}

	missingDocumentLabels, err := q.GetMissingClientDocuments(ctx, id)
	if err != nil {
		return nil, err
	}

	counts, err := q.GetClientPageCounts(ctx, id)
	if err != nil {
		return nil, err
	}

	coordinatorRows, err := q.GetClientCoordinator(ctx, id)
	if err != nil {
		return nil, err
	}

	latestStatusHistory, err := q.GetClientLatestStatusHistory(ctx, id)
	if err != nil {
		return nil, err
	}

	activeContracts, err := q.ListClientActiveApprovedContracts(ctx, id)
	if err != nil {
		return nil, err
	}

	latestDraft, err := q.GetLatestDraftEvaluationByClient(ctx, id)
	hasLatestDraft := false
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	} else {
		hasLatestDraft = true
	}

	latestCompleted, err := q.GetLatestCompletedEvaluationByClient(ctx, id)
	hasLatestCompleted := false
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	} else {
		hasLatestCompleted = true
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	goalsResponse := make([]domain.ClientGoalSummary, len(goals))
	for i, goal := range goals {
		var topicName *string
		if goal.TopicName != "" {
			topicName = &goal.TopicName
		}
		goalsResponse[i] = domain.ClientGoalSummary{
			Title:     goal.Title,
			Priority:  string(goal.Priority),
			TopicName: topicName,
		}
	}

	emergencyContactsResponse := make([]domain.ClientEmergencyContact, len(emergencyContacts))
	for i, contact := range emergencyContacts {
		emergencyContactsResponse[i] = domain.ClientEmergencyContact{
			ID:           contact.ID,
			FirstName:    contact.FirstName,
			LastName:     contact.LastName,
			Relationship: contact.Relationship,
			PhoneNumber:  contact.PhoneNumber,
			Email:        contact.Email,
		}
	}

	var coordinator *domain.ClientCoordinator
	if len(coordinatorRows) > 0 {
		coordinator = &domain.ClientCoordinator{
			EmployeeID: &coordinatorRows[0].EmployeeID,
			FirstName:  &coordinatorRows[0].FirstName,
			LastName:   &coordinatorRows[0].LastName,
			StartDate:  conv.TimePtrFromPgDate(coordinatorRows[0].StartDate),
		}
	}

	var statusTimeline *domain.ClientStatusTimeline
	if latestStatusHistory.LastChangeReason != nil || latestStatusHistory.LastChangedAt.Valid || latestStatusHistory.LastStatus != "" {
		var lastChangedAt *time.Time
		if latestStatusHistory.LastChangedAt.Valid {
			t := latestStatusHistory.LastChangedAt.Time
			lastChangedAt = &t
		}
		var lastStatus *string
		if latestStatusHistory.LastStatus != "" {
			lastStatus = &latestStatusHistory.LastStatus
		}
		statusTimeline = &domain.ClientStatusTimeline{
			LastChangeReason: latestStatusHistory.LastChangeReason,
			LastChangedAt:    lastChangedAt,
			LastStatus:       lastStatus,
		}
	}

	var sender *domain.ClientSenderMinimal
	if client.SenderName != nil || client.SenderEmailAddress != nil || client.SenderPhoneNumber != nil {
		senderName := ""
		if client.SenderName != nil {
			senderName = *client.SenderName
		}
		sender = &domain.ClientSenderMinimal{
			Name:         senderName,
			EmailAddress: client.SenderEmailAddress,
			PhoneNumber:  client.SenderPhoneNumber,
		}
	}

	var location *domain.ClientLocation
	if client.LocationID != nil && client.LocationName != nil {
		location = &domain.ClientLocation{
			ID:   *client.LocationID,
			Name: *client.LocationName,
		}
	}

	dateOfBirth := conv.TimePtrFromPgDate(client.DateOfBirth)
	age := calculateAge(client.DateOfBirth)

	riskFlags := make([]string, 0, 10)
	if client.RiskAggressiveBehavior != nil && *client.RiskAggressiveBehavior {
		riskFlags = append(riskFlags, "risk_aggressive_behavior")
	}
	if client.RiskSuicidalSelfharm != nil && *client.RiskSuicidalSelfharm {
		riskFlags = append(riskFlags, "risk_suicidal_selfharm")
	}
	if client.RiskSubstanceAbuse != nil && *client.RiskSubstanceAbuse {
		riskFlags = append(riskFlags, "risk_substance_abuse")
	}
	if client.RiskPsychiatricIssues != nil && *client.RiskPsychiatricIssues {
		riskFlags = append(riskFlags, "risk_psychiatric_issues")
	}
	if client.RiskCriminalHistory != nil && *client.RiskCriminalHistory {
		riskFlags = append(riskFlags, "risk_criminal_history")
	}
	if client.RiskFlightBehavior != nil && *client.RiskFlightBehavior {
		riskFlags = append(riskFlags, "risk_flight_behavior")
	}
	if client.RiskWeaponPossession != nil && *client.RiskWeaponPossession {
		riskFlags = append(riskFlags, "risk_weapon_possession")
	}
	if client.RiskSexualBehavior != nil && *client.RiskSexualBehavior {
		riskFlags = append(riskFlags, "risk_sexual_behavior")
	}
	if client.RiskDayNightRhythm != nil && *client.RiskDayNightRhythm {
		riskFlags = append(riskFlags, "risk_day_night_rhythm")
	}
	if client.RiskOther != nil && *client.RiskOther {
		riskFlags = append(riskFlags, "risk_other")
	}

	alerts := make([]domain.ClientPageAlert, 0, 4)
	if len(missingDocumentLabels) > 0 {
		alerts = append(alerts, domain.ClientPageAlert{
			Code:     "missing_documents",
			Severity: "warning",
			Message:  fmt.Sprintf("%d required documents are missing", len(missingDocumentLabels)),
		})
	}
	if len(goalsResponse) == 0 {
		alerts = append(alerts, domain.ClientPageAlert{
			Code:     "missing_goals",
			Severity: "info",
			Message:  "No active goals defined",
		})
	}
	if counts.IncidentsCount > 0 {
		alerts = append(alerts, domain.ClientPageAlert{
			Code:     "has_incidents",
			Severity: "warning",
			Message:  fmt.Sprintf("%d incident(s) linked to client", counts.IncidentsCount),
		})
	}

	var careSchedule *domain.ClientCareSchedule
	if client.Status == db.ClientStatusEnumScheduledInCare {
		var careStartDatePtr *time.Time
		var placedInCareAtPtr *time.Time
		var daysUntilStart *int32

		if client.CareStartDate.Valid {
			careStartDatePtr = conv.TimePtrFromPgDate(client.CareStartDate)
			today := time.Now().UTC()
			todayDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
			startDate := time.Date(client.CareStartDate.Time.Year(), client.CareStartDate.Time.Month(), client.CareStartDate.Time.Day(), 0, 0, 0, 0, time.UTC)
			days := int32(startDate.Sub(todayDate).Hours() / 24)
			if days < 0 {
				days = 0
			}
			daysUntilStart = &days
		}

		if client.PlacedInCareAt.Valid {
			placedAt := client.PlacedInCareAt.Time
			placedInCareAtPtr = &placedAt
		}

		shouldBeActiveNow := false
		if client.CareStartDate.Valid {
			shouldBeActiveNow = !client.CareStartDate.Time.After(time.Now().UTC())
		}

		if shouldBeActiveNow {
			alerts = append(alerts, domain.ClientPageAlert{
				Code:     "start_date_reached_not_activated",
				Severity: "warning",
				Message:  "Care start date has been reached but client is still scheduled_in_care",
			})
		}

		if coordinator == nil {
			alerts = append(alerts, domain.ClientPageAlert{
				Code:     "missing_coordinator",
				Severity: "warning",
				Message:  "No coordinator assigned",
			})
		}

		careSchedule = &domain.ClientCareSchedule{
			CareStartDate:      careStartDatePtr,
			PlacedInCareAt:     placedInCareAtPtr,
			DaysUntilStart:     daysUntilStart,
			ShouldBeActiveNow:  shouldBeActiveNow,
			NextEvaluationDate: conv.TimePtrFromPgDate(client.NextEvaluationDate),
		}
	}

	var care *domain.ClientInCare
	var dischargeSchedule *domain.ClientDischargeSchedule
	var dischargeSummary *domain.ClientDischargeSummary
	var contractSummary *domain.ClientContractSummary
	var evaluationSummary *domain.ClientEvaluationSummary
	if client.Status == db.ClientStatusEnumInCare {
		var daysInCare *int32
		if client.CareStartDate.Valid {
			today := time.Now().UTC()
			todayDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
			startDate := time.Date(client.CareStartDate.Time.Year(), client.CareStartDate.Time.Month(), client.CareStartDate.Time.Day(), 0, 0, 0, 0, time.UTC)
			days := int32(todayDate.Sub(startDate).Hours() / 24)
			if days < 0 {
				days = 0
			}
			daysInCare = &days
		}

		care = &domain.ClientInCare{
			CareStartDate:            conv.TimePtrFromPgDate(client.CareStartDate),
			PlacedInCareAt:           conv.TimePtrFromPgTimestamptz(client.PlacedInCareAt),
			DaysInCare:               daysInCare,
			EvaluationIntervalsWeeks: client.EvaluationIntervalsWeeks,
			LastEvaluationAnchorDate: conv.TimePtrFromPgDate(client.LastEvaluationAnchorDate),
			NextEvaluationDate:       conv.TimePtrFromPgDate(client.NextEvaluationDate),
		}

		var activeContractPtr *domain.ClientActiveContract
		var daysUntilContractEndPtr *int32
		if len(activeContracts) > 0 {
			active := activeContracts[0]
			activeContractPtr = &domain.ClientActiveContract{
				ID:              active.ID,
				Status:          stringPtrOrNil(active.Status),
				StartDate:       conv.TimePtrFromPgTimestamptz(active.StartDate),
				EndDate:         conv.TimePtrFromPgTimestamptz(active.EndDate),
				FinancingAct:    stringPtrOrNil(active.FinancingAct),
				FinancingOption: stringPtrOrNil(active.FinancingOption),
				CareType:        stringPtrOrNil(active.CareType),
			}
			daysUntilContractEndPtr = &active.DaysUntilContractEnd
		}

		hasActiveApprovedContract := len(activeContracts) > 0
		contractSummary = &domain.ClientContractSummary{
			HasActiveApprovedContract: hasActiveApprovedContract,
			ActiveContract:            activeContractPtr,
			DaysUntilContractEnd:      daysUntilContractEndPtr,
		}

		var draftPtr *domain.ClientEvaluationDraft
		if hasLatestDraft {
			draftPtr = &domain.ClientEvaluationDraft{
				ID:        latestDraft.ID,
				UpdatedAt: conv.TimePtrFromPgTimestamptz(latestDraft.UpdatedAt),
			}
		}

		var lastCompletedPtr *domain.ClientEvaluationLastCompleted
		if hasLatestCompleted {
			lastCompletedPtr = &domain.ClientEvaluationLastCompleted{
				ID:                  latestCompleted.ID,
				SubmittedAt:         conv.TimePtrFromPgTimestamptz(latestCompleted.SubmittedAt),
				CreatedByEmployeeID: latestCompleted.CreatedByEmployeeID,
			}
		}

		var daysLeft *int32
		var priority *string
		nextEvaluationDate := conv.TimePtrFromPgDate(client.NextEvaluationDate)
		if client.NextEvaluationDate.Valid {
			today := time.Now().UTC()
			todayDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
			nextDate := time.Date(client.NextEvaluationDate.Time.Year(), client.NextEvaluationDate.Time.Month(), client.NextEvaluationDate.Time.Day(), 0, 0, 0, 0, time.UTC)
			delta := int32(nextDate.Sub(todayDate).Hours() / 24)
			daysLeft = &delta

			p := "normal"
			if delta <= 3 {
				p = "critical"
			}
			priority = &p
		}

		evaluationSummary = &domain.ClientEvaluationSummary{
			NextEvaluationDate: nextEvaluationDate,
			DaysLeft:           daysLeft,
			Priority:           priority,
			Draft:              draftPtr,
			LastCompleted:      lastCompletedPtr,
		}

		if coordinator == nil {
			alerts = append(alerts, domain.ClientPageAlert{
				Code:     "missing_coordinator",
				Severity: "warning",
				Message:  "No coordinator assigned",
			})
		}

		if !hasActiveApprovedContract {
			alerts = append(alerts, domain.ClientPageAlert{
				Code:     "missing_active_contract",
				Severity: "warning",
				Message:  "No active approved contract",
			})
		} else if daysUntilContractEndPtr != nil && *daysUntilContractEndPtr <= 30 {
			alerts = append(alerts, domain.ClientPageAlert{
				Code:     "contract_ends_soon",
				Severity: "warning",
				Message:  fmt.Sprintf("Active contract ends in %d day(s)", *daysUntilContractEndPtr),
			})
		}

		if daysLeft != nil {
			if *daysLeft < 0 {
				alerts = append(alerts, domain.ClientPageAlert{
					Code:     "evaluation_overdue",
					Severity: "warning",
					Message:  fmt.Sprintf("Evaluation is overdue by %d day(s)", -*daysLeft),
				})
			} else if *daysLeft <= 3 {
				alerts = append(alerts, domain.ClientPageAlert{
					Code:     "evaluation_due_soon",
					Severity: "warning",
					Message:  fmt.Sprintf("Evaluation due in %d day(s)", *daysLeft),
				})
			}
		}
	}

	if client.Status == db.ClientStatusEnumScheduledOutOfCare {
		dischargeDate := conv.TimePtrFromPgDate(client.DischargeDate)
		var dischargeReason *string
		if client.DischargeReason != nil {
			reason := string(*client.DischargeReason)
			dischargeReason = &reason
		}

		var daysUntilDischarge *int32
		isDue := false
		if client.DischargeDate.Valid {
			today := time.Now().UTC()
			todayDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
			dischargeDay := time.Date(client.DischargeDate.Time.Year(), client.DischargeDate.Time.Month(), client.DischargeDate.Time.Day(), 0, 0, 0, 0, time.UTC)
			days := int32(dischargeDay.Sub(todayDate).Hours() / 24)
			daysUntilDischarge = &days
			isDue = days <= 0
		}

		missingFinalEvaluation := strings.TrimSpace(derefString(client.FinalEvaluation)) == ""
		dischargeSchedule = &domain.ClientDischargeSchedule{
			DischargeDate:          dischargeDate,
			DischargeReason:        dischargeReason,
			FinalEvaluation:        client.FinalEvaluation,
			DaysUntilDischarge:     daysUntilDischarge,
			IsDue:                  isDue,
			MissingFinalEvaluation: missingFinalEvaluation,
		}

		if isDue && missingFinalEvaluation {
			alerts = append(alerts, domain.ClientPageAlert{
				Code:     "discharge_due_missing_final_evaluation",
				Severity: "warning",
				Message:  "Discharge is due but final evaluation is missing",
			})
		}
	}

	if client.Status == db.ClientStatusEnumOutOfCare {
		var dischargeReason *string
		if client.DischargeReason != nil {
			reason := string(*client.DischargeReason)
			dischargeReason = &reason
		}

		dischargeSummary = &domain.ClientDischargeSummary{
			DischargeDate:   conv.TimePtrFromPgDate(client.DischargeDate),
			DischargeReason: dischargeReason,
			FinalEvaluation: client.FinalEvaluation,
		}
	}

	waitlistSince := time.Now().UTC()
	if client.CreatedAt.Valid {
		waitlistSince = client.CreatedAt.Time
	}

	// Build bsn_verified_by full name
	var bsnVerifiedByName *string
	if client.BsnVerifiedByFirstName != nil || client.BsnVerifiedByLastName != nil {
		name := ""
		if client.BsnVerifiedByFirstName != nil {
			name = *client.BsnVerifiedByFirstName
		}
		if client.BsnVerifiedByLastName != nil {
			if name != "" {
				name += " "
			}
			name += *client.BsnVerifiedByLastName
		}
		bsnVerifiedByName = &name
	}

	return &domain.ClientPageDetail{
		SchemaVersion: 1,
		Status:        string(client.Status),
		Client: domain.ClientPageClient{
			ID:                client.ID,
			FirstName:         client.FirstName,
			LastName:          client.LastName,
			Bsn:               client.Bsn,
			BsnVerifiedBy:     client.BsnVerifiedBy,
			BsnVerifiedByName: bsnVerifiedByName,
			FileNumber:        client.Filenumber,
			Gender:            string(client.Gender),
			DateOfBirth:       dateOfBirth,
			Age:               age,
			CareType:          db.IntakeCareTypePtrFromEnum(client.CareType),
			Address: domain.ClientAddress{
				Street:              client.Street,
				HouseNumber:         client.HouseNumber,
				HouseNumberAddition: client.HouseNumberAddition,
				PostalCode:          client.PostalCode,
				City:                client.City,
			},
			Location:                   location,
			EducationCurrentlyEnrolled: client.EducationCurrentlyEnrolled,
			EducationInstitution:       client.EducationInstitution,
			EducationMentorName:        client.EducationMentorName,
			EducationMentorPhone:       client.EducationMentorPhone,
			EducationMentorEmail:       client.EducationMentorEmail,
			EducationAdditionalNotes:   client.EducationAdditionalNotes,
			EducationLevel:             string(client.EducationLevel),
			WorkCurrentlyEmployed:      client.WorkCurrentlyEmployed,
			WorkCurrentEmployer:        client.WorkCurrentEmployer,
			WorkCurrentEmployerPhone:   client.WorkCurrentEmployerPhone,
			WorkCurrentEmployerEmail:   client.WorkCurrentEmployerEmail,
			WorkCurrentPosition:        client.WorkCurrentPosition,
			WorkStartDate:              conv.TimeFromPgDate(client.WorkStartDate),
			WorkAdditionalNotes:        client.WorkAdditionalNotes,
		},
		Care:              care,
		CareSchedule:      careSchedule,
		DischargeSchedule: dischargeSchedule,
		DischargeSummary:  dischargeSummary,
		Sender:            sender,
		Coordinator:       coordinator,
		ContractSummary:   contractSummary,
		EvaluationSummary: evaluationSummary,
		EmergencyContacts: emergencyContactsResponse,
		Documents: domain.ClientDocuments{
			Existing: existingDocumentLabels,
			Missing:  missingDocumentLabels,
		},
		Goals: goalsResponse,
		Intake: domain.ClientIntake{
			SelfSufficiencyScore: client.IntakeSelfSufficiency,
			Conclusion:           db.IntakeConclusionPtrFromEnum(client.IntakeConclusion),
			ConclusionNotes:      client.IntakeConclusionNotes,
		},
		Risks: domain.ClientRiskSummary{
			Flags: riskFlags,
			Notes: client.RiskAdditionalNotes,
		},
		Counts: domain.ClientPageCounts{
			Contracts:    counts.ContractsCount,
			Incidents:    counts.IncidentsCount,
			Reports:      counts.ReportsCount,
			Evaluations:  counts.EvaluationsCount,
			Documents:    counts.DocumentsCount,
			Appointments: counts.AppointmentsCount,
		},
		Alerts:         alerts,
		Meta:           domain.ClientPageMeta{WaitlistSince: waitlistSince, LastUpdatedAt: time.Now().UTC()},
		StatusTimeline: statusTimeline,
	}, nil
}

func calculateAge(dateOfBirth pgtype.Date) *int32 {
	if !dateOfBirth.Valid {
		return nil
	}

	today := time.Now().UTC()
	age := int32(today.Year() - dateOfBirth.Time.Year())
	birthdayThisYear := time.Date(today.Year(), dateOfBirth.Time.Month(), dateOfBirth.Time.Day(), 0, 0, 0, 0, time.UTC)
	if today.Before(birthdayThisYear) {
		age--
	}

	if age < 0 {
		return nil
	}

	return &age
}

func (r *ClientRepository) UpdateClient(ctx context.Context, id uuid.UUID, params domain.UpdateClientParams) (*domain.Client, error) {
	var updated db.ClientDetail

	var gender *db.GenderEnum
	if params.Gender != nil {
		value := db.GenderEnum(*params.Gender)
		gender = &value
	}
	var educationLevel *db.EducationLevelEnum
	if params.EducationLevel != nil {
		value := db.EducationLevelEnum(*params.EducationLevel)
		educationLevel = &value
	}

	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		var err error
		updated, err = q.UpdateClientDetailsV2(ctx, db.UpdateClientDetailsV2Params{
			ID:                         id,
			FirstName:                  params.FirstName,
			LastName:                   params.LastName,
			DateOfBirth:                conv.PgDateFromTime(params.DateOfBirth),
			Identity:                   params.Identity,
			Bsn:                        params.Bsn,
			BsnVerifiedBy:              params.BsnVerifiedBy,
			Email:                      params.Email,
			PhoneNumber:                params.PhoneNumber,
			Gender:                     gender,
			Filenumber:                 params.Filenumber,
			SenderID:                   params.SenderID,
			LocationID:                 params.LocationID,
			EducationCurrentlyEnrolled: params.EducationCurrentlyEnrolled,
			EducationInstitution:       params.EducationInstitution,
			EducationMentorName:        params.EducationMentorName,
			EducationMentorPhone:       params.EducationMentorPhone,
			EducationMentorEmail:       params.EducationMentorEmail,
			EducationAdditionalNotes:   params.EducationAdditionalNotes,
			EducationLevel:             educationLevel,
			WorkCurrentlyEmployed:      params.WorkCurrentlyEmployed,
			WorkCurrentEmployer:        params.WorkCurrentEmployer,
			WorkCurrentEmployerPhone:   params.WorkCurrentEmployerPhone,
			WorkCurrentEmployerEmail:   params.WorkCurrentEmployerEmail,
			WorkCurrentPosition:        params.WorkCurrentPosition,
			WorkStartDate:              conv.PgDateFromTime(params.WorkStartDate),
			WorkAdditionalNotes:        params.WorkAdditionalNotes,
			Nationality:                params.Nationality,
		})
		if err != nil {
			return err
		}

		if params.CoordinatorEmployeeID != nil {
			if err := q.UpsertMainCoordinatorByEmployee(ctx, db.UpsertMainCoordinatorByEmployeeParams{
				ClientID:   id,
				EmployeeID: *params.CoordinatorEmployeeID,
			}); err != nil {
				return fmt.Errorf("failed to update main coordinator: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return toDomainClient(updated), nil
}

func (r *ClientRepository) GetClientAddresses(ctx context.Context, id uuid.UUID) ([]domain.ClientAddress, error) {
	return []domain.ClientAddress{}, nil
}

// CreateProgressReport creates a new progress report.
func (r *ClientRepository) CreateProgressReport(ctx context.Context, params domain.CreateProgressReportParams) (*domain.ProgressReport, error) {
	var report db.ProgressReport
	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		var err error
		report, err = q.CreateProgressReport(ctx, db.CreateProgressReportParams{
			ClientID:       params.ClientID,
			EmployeeID:     params.EmployeeID,
			Title:          params.Title,
			Date:           conv.PgTimestamptzFromTime(params.Date),
			ReportText:     params.ReportText,
			Type:           db.ProgressReportTypeEnum(params.Type),
			EmotionalState: db.EmotionalStateEnum(params.EmotionalState),
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	return toDomainProgressReport(report), nil
}

// ListProgressReports lists progress reports for a client.
func (r *ClientRepository) ListProgressReports(ctx context.Context, params domain.ListProgressReportsParams) (*domain.ListProgressReportsResult, error) {
	rows, err := actorQuery(ctx, r.store, func(q *db.Queries) ([]db.ListProgressReportsRow, error) {
		return q.ListProgressReports(ctx, db.ListProgressReportsParams{
			ClientID: params.ClientID,
			Type:     db.NullProgressReportTypeFromPtr(params.Type),
			Offset:   params.Offset,
			Limit:    params.Limit,
		})
	})
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return &domain.ListProgressReportsResult{
			Items:      []domain.ProgressReport{},
			TotalCount: 0,
		}, nil
	}

	items := make([]domain.ProgressReport, len(rows))
	for i, row := range rows {
		items[i] = toDomainProgressReportFromRow(row)
	}

	return &domain.ListProgressReportsResult{
		Items:      items,
		TotalCount: rows[0].TotalCount,
	}, nil
}

// GetProgressReport retrieves a single progress report by ID.
func (r *ClientRepository) GetProgressReport(ctx context.Context, reportID uuid.UUID) (*domain.ProgressReport, error) {
	row, err := actorQuery(ctx, r.store, func(q *db.Queries) (db.GetProgressReportRow, error) {
		return q.GetProgressReport(ctx, reportID)
	})
	if err != nil {
		return nil, err
	}
	return toDomainProgressReportFromGetRow(row), nil
}

// UpdateProgressReport updates an existing progress report.
func (r *ClientRepository) UpdateProgressReport(ctx context.Context, params domain.UpdateProgressReportParams) (*domain.ProgressReport, error) {
	var report db.ProgressReport
	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		var err error
		report, err = q.UpdateProgressReport(ctx, db.UpdateProgressReportParams{
			ID:             params.ID,
			EmployeeID:     params.EmployeeID,
			Title:          params.Title,
			Date:           conv.PgTimestamptzFromTime(params.Date),
			ReportText:     params.ReportText,
			Type:           db.NullProgressReportTypeFromPtr(params.Type),
			EmotionalState: db.NullEmotionalStateFromPtr(params.EmotionalState),
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	return toDomainProgressReport(report), nil
}

// DeleteProgressReport deletes a progress report by ID.
func (r *ClientRepository) DeleteProgressReport(ctx context.Context, reportID uuid.UUID) error {
	return r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		return q.DeleteProgressReport(ctx, reportID)
	})
}

// GetProgressReportsByDateRange retrieves progress reports within a date range.
func (r *ClientRepository) GetProgressReportsByDateRange(ctx context.Context, params domain.GetProgressReportsByDateRangeParams) ([]domain.ProgressReport, error) {
	rows, err := actorQuery(ctx, r.store, func(q *db.Queries) ([]db.ProgressReport, error) {
		return q.GetProgressReportsByDateRange(ctx, db.GetProgressReportsByDateRangeParams{
			ClientID:  params.ClientID,
			StartDate: conv.PgTimestamptzFromTime(params.StartDate),
			EndDate:   conv.PgTimestamptzFromTime(params.EndDate),
		})
	})
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return []domain.ProgressReport{}, nil
	}

	reports := make([]domain.ProgressReport, len(rows))
	for i, row := range rows {
		reports[i] = *toDomainProgressReport(row)
	}

	return reports, nil
}

// CreateAiGeneratedReport creates a new AI-generated report.
func (r *ClientRepository) CreateAiGeneratedReport(ctx context.Context, params domain.CreateAiGeneratedReportParams) (*domain.AiGeneratedReport, error) {
	var report db.AiGeneratedReport
	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		var err error
		report, err = q.CreateAiGeneratedReport(ctx, db.CreateAiGeneratedReportParams{
			ClientID:   params.ClientID,
			ReportText: params.ReportText,
			StartDate:  conv.PgDateFromTime(params.StartDate),
			EndDate:    conv.PgDateFromTime(params.EndDate),
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	return toDomainAiGeneratedReport(report), nil
}

// ListAiGeneratedReports lists AI-generated reports for a client.
func (r *ClientRepository) ListAiGeneratedReports(ctx context.Context, params domain.ListAiGeneratedReportsParams) (*domain.ListAiGeneratedReportsResult, error) {
	rows, err := actorQuery(ctx, r.store, func(q *db.Queries) ([]db.ListAiGeneratedReportsRow, error) {
		return q.ListAiGeneratedReports(ctx, db.ListAiGeneratedReportsParams{
			ClientID: params.ClientID,
			Limit:    params.Limit,
			Offset:   params.Offset,
		})
	})
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return &domain.ListAiGeneratedReportsResult{
			Items:      []domain.AiGeneratedReport{},
			TotalCount: 0,
		}, nil
	}

	items := make([]domain.AiGeneratedReport, len(rows))
	for i, row := range rows {
		items[i] = toDomainAiGeneratedReportFromRow(row)
	}

	return &domain.ListAiGeneratedReportsResult{
		Items:      items,
		TotalCount: rows[0].TotalCount,
	}, nil
}

func stringPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func toDomainClient(row db.ClientDetail) *domain.Client {
	return &domain.Client{
		ID:                         row.ID,
		FirstName:                  row.FirstName,
		LastName:                   row.LastName,
		DateOfBirth:                conv.TimeFromPgDate(row.DateOfBirth),
		Identity:                   row.Identity,
		Status:                     string(row.Status),
		Bsn:                        row.Bsn,
		BsnVerifiedBy:              row.BsnVerifiedBy,
		Email:                      row.Email,
		PhoneNumber:                row.PhoneNumber,
		Gender:                     string(row.Gender),
		Filenumber:                 row.Filenumber,
		CreatedAt:                  conv.TimeFromPgTimestamptz(row.CreatedAt),
		SenderID:                   row.SenderID,
		LocationID:                 row.LocationID,
		EducationCurrentlyEnrolled: row.EducationCurrentlyEnrolled,
		EducationInstitution:       row.EducationInstitution,
		EducationMentorName:        row.EducationMentorName,
		EducationMentorPhone:       row.EducationMentorPhone,
		EducationMentorEmail:       row.EducationMentorEmail,
		EducationAdditionalNotes:   row.EducationAdditionalNotes,
		EducationLevel:             string(row.EducationLevel),
		WorkCurrentlyEmployed:      row.WorkCurrentlyEmployed,
		WorkCurrentEmployer:        row.WorkCurrentEmployer,
		WorkCurrentEmployerPhone:   row.WorkCurrentEmployerPhone,
		WorkCurrentEmployerEmail:   row.WorkCurrentEmployerEmail,
		WorkCurrentPosition:        row.WorkCurrentPosition,
		WorkStartDate:              conv.TimeFromPgDate(row.WorkStartDate),
		WorkAdditionalNotes:        row.WorkAdditionalNotes,
	}
}

func toDomainClientListItem(row db.ListClientDetailsRow) domain.ClientListItem {
	return domain.ClientListItem{
		ID:           row.ID,
		FirstName:    row.FirstName,
		LastName:     row.LastName,
		Bsn:          row.Bsn,
		Filenumber:   row.Filenumber,
		LocationName: row.LocationName,
		CareType:     db.IntakeCareTypePtrFromEnum(row.CareType),
		Status:       string(row.Status),
		GoalsCount:   row.GoalsCount,
		RiskCount:    row.RiskCount,
		CreatedAt:    conv.TimeFromPgTimestamptz(row.CreatedAt),
	}
}

func toDomainWaitingListClient(row db.ListWaitingListClientsRow) domain.WaitingListClient {
	var admissionType *string
	if row.AdmissionType != nil {
		s := string(*row.AdmissionType)
		admissionType = &s
	}

	return domain.WaitingListClient{
		ID:             row.ID,
		FirstName:      row.FirstName,
		LastName:       row.LastName,
		CareType:       db.IntakeCareTypePtrFromEnum(row.CareType),
		Bsn:            row.Bsn,
		SenderName:     row.SenderName,
		DaysInWaitlist: row.DaysInWaitlist,
		AdmissionType:  admissionType,
	}
}

func toDomainInCareClient(row db.ListInCareClientsRow) domain.InCareClient {
	var coordinatorName *string
	if row.CoordinatorName != "" {
		c := row.CoordinatorName
		coordinatorName = &c
	}

	return domain.InCareClient{
		ID:                row.ID,
		Bsn:               row.Bsn,
		FirstName:         row.FirstName,
		LastName:          row.LastName,
		CoordinatorName:   coordinatorName,
		LocationName:      row.LocationName,
		Status:            string(row.Status),
		CareStartDate:     conv.TimePtrFromPgDate(row.CareStartDate),
		DaysInCare:        row.DaysInCare,
		HasActiveContract: row.HasActiveContract,
	}
}

func (r *ClientRepository) AddClientDocuments(ctx context.Context, clientID uuid.UUID, params domain.AddClientDocumentParams) ([]domain.AddClientDocumentResult, error) {
	txDocs := make([]db.AddClientDocumentTxParams, 0, len(params.Documents))
	for _, doc := range params.Documents {
		txDocs = append(txDocs, db.AddClientDocumentTxParams{
			ClientID:     clientID,
			AttachmentID: doc.AttachmentID,
			Label:        doc.Label,
		})
	}

	result, err := r.store.AddClientDocumentsTx(ctx, db.AddClientDocumentsTxParams{
		ClientID:  clientID,
		Documents: txDocs,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrClientNotFound
		}
		return nil, fmt.Errorf("failed to add client documents: %w", err)
	}

	docs := make([]domain.AddClientDocumentResult, 0, len(result.Documents))
	for _, item := range result.Documents {
		attachmentID := item.ClientDocument.AttachmentUuid
		docs = append(docs, domain.AddClientDocumentResult{
			ID:           item.ClientDocument.ID,
			AttachmentID: &attachmentID,
			ClientID:     item.ClientDocument.ClientID,
			Label:        string(item.ClientDocument.Label),
			Name:         item.Attachment.Name,
			File:         item.Attachment.File,
			Size:         item.Attachment.Size,
			IsUsed:       item.Attachment.IsUsed,
			Tag:          item.Attachment.Tag,
			UpdatedAt:    conv.TimeFromPgTimestamptz(item.Attachment.Updated),
			CreatedAt:    conv.TimeFromPgTimestamptz(item.Attachment.Created),
		})
	}
	return docs, nil
}

func (r *ClientRepository) ListClientDocuments(ctx context.Context, params domain.ListClientDocumentsParams) (*domain.ListClientDocumentsResult, error) {
	rows, err := actorQuery(ctx, r.store, func(q *db.Queries) ([]db.ListClientDocumentsRow, error) {
		return q.ListClientDocuments(ctx, db.ListClientDocumentsParams{
			ClientID: params.ClientID,
			Limit:    params.Limit,
			Offset:   params.Offset,
		})
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list client documents: %w", err)
	}

	if len(rows) == 0 {
		return &domain.ListClientDocumentsResult{
			Documents:  []domain.ClientDocument{},
			TotalCount: 0,
		}, nil
	}

	totalCount := rows[0].TotalCount
	docs := make([]domain.ClientDocument, 0, len(rows))
	for _, row := range rows {
		attachmentID := row.AttachmentUuid
		docs = append(docs, domain.ClientDocument{
			ID:           row.ID,
			AttachmentID: &attachmentID,
			Uuid:         row.Uuid,
			ClientID:     row.ClientID,
			Label:        string(row.Label),
			Name:         row.Name,
			File:         row.File,
			Size:         row.Size,
			IsUsed:       row.IsUsed,
			Tag:          row.Tag,
			UpdatedAt:    conv.TimeFromPgTimestamptz(row.Updated),
			CreatedAt:    conv.TimeFromPgTimestamptz(row.Created),
		})
	}

	return &domain.ListClientDocumentsResult{
		Documents:  docs,
		TotalCount: totalCount,
	}, nil
}

func (r *ClientRepository) DeleteClientDocument(ctx context.Context, clientID uuid.UUID, documentID uuid.UUID) (*domain.DeleteClientDocumentResult, error) {
	result, err := r.store.DeleteClientDocumentTx(ctx, db.DeleteClientDocumentTxParams{
		ClientID:   clientID,
		DocumentID: documentID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrClientDocumentNotFound
		}
		return nil, fmt.Errorf("failed to delete client document: %w", err)
	}

	attachmentID := result.ClientDocument.AttachmentUuid
	return &domain.DeleteClientDocumentResult{
		ID:           result.ClientDocument.ID,
		AttachmentID: &attachmentID,
	}, nil
}

func (r *ClientRepository) GetMissingClientDocuments(ctx context.Context, clientID uuid.UUID) ([]string, error) {
	missingDocs, err := actorQuery(ctx, r.store, func(q *db.Queries) ([]string, error) {
		return q.GetMissingClientDocuments(ctx, clientID)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get missing client documents: %w", err)
	}
	return missingDocs, nil
}

func (r *ClientRepository) GetAttachmentsByUUIDs(ctx context.Context, ids []uuid.UUID) ([]domain.AttachmentFile, error) {
	attachments, err := actorQuery(ctx, r.store, func(q *db.Queries) ([]db.AttachmentFile, error) {
		return q.GetActorAttachmentsByUUIDs(ctx, ids)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get attachments: %w", err)
	}

	result := make([]domain.AttachmentFile, 0, len(attachments))
	for _, a := range attachments {
		result = append(result, domain.AttachmentFile{
			UUID:      a.Uuid,
			Name:      a.Name,
			File:      a.File,
			Size:      a.Size,
			IsUsed:    a.IsUsed,
			Tag:       a.Tag,
			UpdatedAt: conv.TimeFromPgTimestamptz(a.Updated),
			CreatedAt: conv.TimeFromPgTimestamptz(a.Created),
		})
	}
	return result, nil
}

// CreateClientGoal creates a new client goal.
func (r *ClientRepository) CreateClientGoal(ctx context.Context, clientID uuid.UUID, params domain.CreateClientGoalParams) (*domain.ClientGoal, error) {
	var createdGoal db.ClientGoal

	title := strings.TrimSpace(params.Title)
	if title == "" {
		return nil, domain.ErrClientGoalTitleRequired
	}

	priority := db.ClientGoalPriorityEnumMedium
	if params.Priority != nil && strings.TrimSpace(*params.Priority) != "" {
		priority = db.ClientGoalPriorityEnum(strings.TrimSpace(*params.Priority))
	}

	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		if _, err := q.GetClientDetails(ctx, clientID); err != nil {
			return fmt.Errorf("%w: %w", domain.ErrClientGoalClientNotFound, err)
		}

		topic, err := q.GetTopicByID(ctx, params.TopicID)
		if err != nil {
			return fmt.Errorf("%w: %w", domain.ErrClientGoalTopicNotFound, err)
		}

		sortOrder := int32(0)
		if params.SortOrder != nil {
			sortOrder = *params.SortOrder
		} else {
			nextSortOrder, err := q.GetNextActiveClientGoalSortOrder(ctx, clientID)
			if err != nil {
				return fmt.Errorf("failed to compute next goal sort order: %w", err)
			}
			sortOrder = nextSortOrder
		}

		createdGoal, err = q.CreateManualClientGoal(ctx, db.CreateManualClientGoalParams{
			ClientID:          clientID,
			Title:             title,
			Description:       params.Description,
			Priority:          priority,
			TopicID:           &topic.ID,
			TopicNameSnapshot: &topic.TopicName,
			SortOrder:         sortOrder,
		})
		if err != nil {
			return fmt.Errorf("failed to create client goal: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return toDomainClientGoal(createdGoal), nil
}

// UpdateClientGoal updates an existing client goal.
func (r *ClientRepository) UpdateClientGoal(ctx context.Context, clientID uuid.UUID, goalID uuid.UUID, params domain.UpdateClientGoalParams) (*domain.UpdateClientGoalResult, error) {
	if params.Title == nil && params.Description == nil && params.Priority == nil && params.TopicID == nil && params.SortOrder == nil {
		return nil, domain.ErrClientGoalEmptyPatch
	}

	var (
		activeGoal        db.ClientGoal
		replacementGoalID *uuid.UUID
		mutationType      = "updated"
	)

	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		if _, err := q.GetClientDetails(ctx, clientID); err != nil {
			return fmt.Errorf("%w: %w", domain.ErrClientGoalClientNotFound, err)
		}

		existingGoal, err := q.GetClientGoalByIDAndClientID(ctx, db.GetClientGoalByIDAndClientIDParams{
			ID:       goalID,
			ClientID: clientID,
		})
		if err != nil {
			return fmt.Errorf("%w: %w", domain.ErrClientGoalGoalNotFound, err)
		}
		if existingGoal.Status != db.ClientGoalStatusEnumActive {
			return domain.ErrClientGoalGoalNotFound
		}

		hasDraft, err := q.ClientHasDraftEvaluationForGoalUpdate(ctx, clientID)
		if err != nil {
			return fmt.Errorf("failed to check for draft evaluations: %w", err)
		}
		if hasDraft {
			return domain.ErrClientGoalDraftEvaluationExists
		}

		nextTitle := existingGoal.Title
		if params.Title != nil {
			trimmedTitle := strings.TrimSpace(*params.Title)
			if trimmedTitle == "" {
				return domain.ErrClientGoalTitleRequired
			}
			nextTitle = trimmedTitle
		}

		nextDescription := existingGoal.Description
		if params.Description != nil {
			nextDescription = normalizeOptionalTrimmedString(params.Description)
		}

		nextPriority := existingGoal.Priority
		if params.Priority != nil && strings.TrimSpace(*params.Priority) != "" {
			nextPriority = db.ClientGoalPriorityEnum(strings.TrimSpace(*params.Priority))
		}

		nextTopicID := existingGoal.TopicID
		nextTopicNameSnapshot := existingGoal.TopicNameSnapshot
		if params.TopicID != nil {
			topic, err := q.GetTopicByID(ctx, *params.TopicID)
			if err != nil {
				return fmt.Errorf("%w: %w", domain.ErrClientGoalTopicNotFound, err)
			}
			nextTopicID = &topic.ID
			nextTopicNameSnapshot = &topic.TopicName
		}

		nextSortOrder := existingGoal.SortOrder
		if params.SortOrder != nil {
			nextSortOrder = *params.SortOrder
		}

		contentChanged := nextTitle != existingGoal.Title ||
			!equalOptionalString(nextDescription, existingGoal.Description) ||
			!equalOptionalUUID(nextTopicID, existingGoal.TopicID)

		priorityChanged := nextPriority != existingGoal.Priority
		sortOrderChanged := nextSortOrder != existingGoal.SortOrder

		if !contentChanged && !priorityChanged && !sortOrderChanged {
			activeGoal = existingGoal
			return nil
		}

		hasHistory, err := q.GoalHasEvaluationItems(ctx, db.GoalHasEvaluationItemsParams{
			GoalID:   goalID,
			ClientID: clientID,
		})
		if err != nil {
			return fmt.Errorf("failed to check goal evaluation history: %w", err)
		}

		if hasHistory && contentChanged {
			if _, err := q.CancelClientGoalByID(ctx, db.CancelClientGoalByIDParams{
				ID:       goalID,
				ClientID: clientID,
			}); err != nil {
				return fmt.Errorf("failed to retire client goal: %w", err)
			}

			activeGoal, err = q.CreateReviewUpdatedClientGoal(ctx, db.CreateReviewUpdatedClientGoalParams{
				ClientID:          clientID,
				Title:             nextTitle,
				Description:       nextDescription,
				Priority:          nextPriority,
				TopicID:           nextTopicID,
				TopicNameSnapshot: nextTopicNameSnapshot,
				SortOrder:         nextSortOrder,
			})
			if err != nil {
				return fmt.Errorf("failed to create replacement client goal: %w", err)
			}

			mutationType = "versioned"
			replacementGoalID = &activeGoal.ID
			return nil
		}

		activeGoal, err = q.UpdateClientGoalByID(ctx, db.UpdateClientGoalByIDParams{
			ID:                goalID,
			ClientID:          clientID,
			Title:             nextTitle,
			Description:       nextDescription,
			Priority:          nextPriority,
			TopicID:           nextTopicID,
			TopicNameSnapshot: nextTopicNameSnapshot,
			SortOrder:         nextSortOrder,
		})
		if err != nil {
			return fmt.Errorf("failed to update client goal: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &domain.UpdateClientGoalResult{
		MutationType:      mutationType,
		GoalID:            goalID,
		ReplacementGoalID: replacementGoalID,
		Goal:              *toDomainClientGoal(activeGoal),
	}, nil
}

// GetClientGoalsForEvaluationPage returns goals for the evaluation page.
func (r *ClientRepository) GetClientGoalsForEvaluationPage(ctx context.Context, clientID uuid.UUID, employeeID uuid.UUID) (*domain.ClientGoalsForEvaluationPage, error) {
	tx, err := r.store.BeginActorTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	q := r.store.WithTx(tx)

	client, err := q.GetClientDetails(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get client details: %w", err)
	}

	response := &domain.ClientGoalsForEvaluationPage{
		CanUpdateGoals: true,
		Goals:          []domain.ClientGoalForEvaluationPage{},
	}

	if client.NextEvaluationDate.Valid {
		nextDate := client.NextEvaluationDate.Time
		response.NextEvaluationDate = &nextDate
	}

	coordinatorRows, err := q.GetClientCoordinator(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get client coordinator: %w", err)
	}
	response.IsResponsibleEmployee = len(coordinatorRows) > 0 && coordinatorRows[0].EmployeeID == employeeID

	activeGoals, err := q.ListActiveGoalsByClientID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to list active goals: %w", err)
	}

	latestProgressRows, err := q.ListLatestCompletedGoalProgressByClient(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to list latest completed goal progress: %w", err)
	}

	latestProgressByGoal := make(map[uuid.UUID]db.ListLatestCompletedGoalProgressByClientRow, len(latestProgressRows))
	for _, row := range latestProgressRows {
		latestProgressByGoal[row.GoalID] = row
	}

	response.Goals = make([]domain.ClientGoalForEvaluationPage, 0, len(activeGoals))
	for _, goal := range activeGoals {
		goalResponse := domain.ClientGoalForEvaluationPage{
			ID:        goal.ID,
			TopicName: goal.TopicNameSnapshot,
			Title:     goal.Title,
			Priority:  string(goal.Priority),
		}

		if latestGoalProgress, ok := latestProgressByGoal[goal.ID]; ok {
			p := string(latestGoalProgress.Progress)
			goalResponse.LastEvaluationProgress = &p
		}

		response.Goals = append(response.Goals, goalResponse)
	}

	_, err = q.GetLatestDraftEvaluationByClient(ctx, clientID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("failed to check for existing draft evaluation: %w", err)
		}
	} else {
		response.CanUpdateGoals = false
		blockReason := "draft_evaluation_exists"
		response.GoalUpdateBlockReason = &blockReason
	}

	draftEval, err := q.GetCurrentCycleDraftEvaluationByClientAndEmployee(ctx, db.GetCurrentCycleDraftEvaluationByClientAndEmployeeParams{
		ClientID:            clientID,
		CreatedByEmployeeID: &employeeID,
	})
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("failed to get current cycle draft evaluation: %w", err)
		}
	} else {
		response.MyDraftEvaluationID = &draftEval.ID
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return response, nil
}

// CreateGoalEvaluation creates or updates a goal evaluation draft and optionally submits it.
func (r *ClientRepository) CreateGoalEvaluation(ctx context.Context, clientID uuid.UUID, employeeID uuid.UUID, params domain.CreateGoalEvaluationParams) (*domain.GoalEvaluation, error) {
	var evaluation db.ClientGoalEvaluation
	var items []db.GetGoalEvaluationItemsRow

	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		client, err := q.GetClientDetails(ctx, clientID)
		if err != nil {
			return fmt.Errorf("failed to get client details: %w", err)
		}

		if client.Status != db.ClientStatusEnumInCare {
			return fmt.Errorf("evaluations can only be created for clients in care")
		}

		activeGoals, err := q.ListActiveGoalsByClientID(ctx, clientID)
		if err != nil {
			return fmt.Errorf("failed to list active goals: %w", err)
		}
		if len(activeGoals) == 0 {
			return fmt.Errorf("client must have at least one active goal to start an evaluation")
		}

		if !client.NextEvaluationDate.Valid {
			return fmt.Errorf("client has no next evaluation date configured")
		}

		reqItemsByGoal, err := mapDraftItemsByGoalID(params.Items)
		if err != nil {
			return err
		}

		evaluationDate := client.NextEvaluationDate
		periodStart := client.LastEvaluationAnchorDate
		periodEnd := client.NextEvaluationDate
		intervalWeeks := client.EvaluationIntervalsWeeks
		if intervalWeeks <= 0 {
			intervalWeeks = 12
		}

		evaluation, err = q.GetGoalEvaluationByClientAndDate(ctx, db.GetGoalEvaluationByClientAndDateParams{
			ClientID:       clientID,
			EvaluationDate: evaluationDate,
		})
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("failed to get current evaluation: %w", err)
			}

			evaluation, err = q.CreateGoalEvaluation(ctx, db.CreateGoalEvaluationParams{
				ClientID:                clientID,
				EvaluationDate:          evaluationDate,
				PeriodStart:             periodStart,
				PeriodEnd:               periodEnd,
				EvaluationIntervalWeeks: intervalWeeks,
				Status:                  db.EvaluationStatusEnumDraft,
				OverallNotes:            params.OverallNotes,
				CreatedByEmployeeID:     &employeeID,
			})
			if err != nil {
				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) && pgErr.Code == "23505" {
					evaluation, err = q.GetGoalEvaluationByClientAndDate(ctx, db.GetGoalEvaluationByClientAndDateParams{
						ClientID:       clientID,
						EvaluationDate: evaluationDate,
					})
					if err != nil {
						if errors.Is(err, pgx.ErrNoRows) {
							return domain.ErrGoalEvaluationNotDraft
						}
						return fmt.Errorf("failed to load concurrent draft evaluation: %w", err)
					}
				} else {
					return fmt.Errorf("failed to create evaluation header: %w", err)
				}
			}
		}
		if evaluation.Status != db.EvaluationStatusEnumDraft {
			return domain.ErrGoalEvaluationNotDraft
		}
		if evaluation.CreatedByEmployeeID == nil || *evaluation.CreatedByEmployeeID != employeeID {
			return domain.ErrGoalEvaluationOwnedByOther
		}

		evaluation, err = q.UpdateGoalEvaluation(ctx, db.UpdateGoalEvaluationParams{
			ID:           evaluation.ID,
			OverallNotes: params.OverallNotes,
			Status:       nil,
		})
		if err != nil {
			return fmt.Errorf("failed to update evaluation header: %w", err)
		}

		activeGoalIDs := make(map[uuid.UUID]struct{}, len(activeGoals))
		for _, goal := range activeGoals {
			activeGoalIDs[goal.ID] = struct{}{}
		}

		existingItems, err := q.GetGoalEvaluationItems(ctx, evaluation.ID)
		if err != nil {
			return fmt.Errorf("failed to load evaluation items: %w", err)
		}

		existingByGoalID := make(map[uuid.UUID]struct{}, len(existingItems))
		for _, item := range existingItems {
			existingByGoalID[item.GoalID] = struct{}{}
		}

		for _, goal := range activeGoals {
			if _, exists := existingByGoalID[goal.ID]; exists {
				continue
			}

			progress := db.ClientGoalProgressEnumNoProgress
			var notes *string
			if reqItem, ok := reqItemsByGoal[goal.ID]; ok {
				parsedProgress, parseErr := parseProgress(reqItem.Progress)
				if parseErr != nil {
					return parseErr
				}
				progress = parsedProgress
				notes = reqItem.Notes
			}

			_, err = q.UpsertGoalEvaluationItem(ctx, db.UpsertGoalEvaluationItemParams{
				ClientID:     clientID,
				EvaluationID: evaluation.ID,
				GoalID:       goal.ID,
				Progress:     progress,
				Notes:        notes,
			})
			if err != nil {
				return fmt.Errorf("failed to ensure evaluation item for goal %s: %w", goal.ID, err)
			}
		}

		for _, reqItem := range params.Items {
			if _, isActive := activeGoalIDs[reqItem.GoalID]; !isActive {
				return fmt.Errorf("goal %s is not an active goal for this client", reqItem.GoalID)
			}

			parsedProgress, parseErr := parseProgress(reqItem.Progress)
			if parseErr != nil {
				return parseErr
			}

			_, err = q.UpsertGoalEvaluationItem(ctx, db.UpsertGoalEvaluationItemParams{
				ClientID:     clientID,
				EvaluationID: evaluation.ID,
				GoalID:       reqItem.GoalID,
				Progress:     parsedProgress,
				Notes:        reqItem.Notes,
			})
			if err != nil {
				return fmt.Errorf("failed to save item for goal %s: %w", reqItem.GoalID, err)
			}
		}

		items, err = q.GetGoalEvaluationItems(ctx, evaluation.ID)
		if err != nil {
			return fmt.Errorf("failed to reload evaluation items: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	var submitErrMessage *string
	if params.Submit {
		submittedEvaluation, submitErr := r.trySubmitGoalEvaluationDraft(ctx, evaluation.ID)
		if submitErr != nil {
			var blockedErr *goalEvaluationSubmitBlockedError
			if errors.As(submitErr, &blockedErr) {
				msg := blockedErr.Error()
				submitErrMessage = &msg
			} else {
				return nil, submitErr
			}
		} else {
			evaluation = submittedEvaluation
		}
	}

	res := toDomainGoalEvaluation(evaluation, items, nil)
	res.SubmitError = submitErrMessage
	return res, nil
}

// UpdateGoalEvaluationDraft updates the exact draft identified by evaluationID.
func (r *ClientRepository) UpdateGoalEvaluationDraft(ctx context.Context, evaluationID uuid.UUID, employeeID uuid.UUID, params domain.UpdateGoalEvaluationDraftParams) (*domain.GoalEvaluation, error) {
	var evaluation db.ClientGoalEvaluation
	var items []db.GetGoalEvaluationItemsRow

	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		row, err := q.GetGoalEvaluationByID(ctx, evaluationID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return domain.ErrGoalEvaluationNotFound
			}
			return fmt.Errorf("failed to load goal evaluation: %w", err)
		}
		if row.Status != db.EvaluationStatusEnumDraft {
			return domain.ErrGoalEvaluationNotDraft
		}
		if row.CreatedByEmployeeID == nil || *row.CreatedByEmployeeID != employeeID {
			return domain.ErrGoalEvaluationOwnedByOther
		}
		if err := ensureCurrentCycleEvaluation(ctx, q, row.ClientID, row.EvaluationDate); err != nil {
			return err
		}

		activeGoals, err := q.ListActiveGoalsByClientID(ctx, row.ClientID)
		if err != nil {
			return fmt.Errorf("failed to list active goals: %w", err)
		}
		if len(activeGoals) == 0 {
			return fmt.Errorf("client must have at least one active goal to update an evaluation")
		}

		evaluation, err = q.UpdateGoalEvaluation(ctx, db.UpdateGoalEvaluationParams{
			ID:           evaluationID,
			OverallNotes: params.OverallNotes,
		})
		if err != nil {
			return fmt.Errorf("failed to update evaluation header: %w", err)
		}

		items, err = saveGoalEvaluationDraftItems(ctx, q, evaluationID, row.ClientID, activeGoals, params.Items)
		return err
	})
	if err != nil {
		return nil, err
	}

	return toDomainGoalEvaluation(evaluation, items, nil), nil
}

// SubmitGoalEvaluationDraft submits the exact saved draft identified by evaluationID.
func (r *ClientRepository) SubmitGoalEvaluationDraft(ctx context.Context, evaluationID uuid.UUID, employeeID uuid.UUID) (*domain.GoalEvaluation, error) {
	var evaluation db.ClientGoalEvaluation
	var items []db.GetGoalEvaluationItemsRow

	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		row, err := q.GetGoalEvaluationByID(ctx, evaluationID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return domain.ErrGoalEvaluationNotFound
			}
			return fmt.Errorf("failed to load goal evaluation: %w", err)
		}
		if row.Status != db.EvaluationStatusEnumDraft {
			return domain.ErrGoalEvaluationNotDraft
		}
		if row.CreatedByEmployeeID == nil || *row.CreatedByEmployeeID != employeeID {
			return domain.ErrGoalEvaluationOwnedByOther
		}
		if err := ensureCurrentCycleEvaluation(ctx, q, row.ClientID, row.EvaluationDate); err != nil {
			return err
		}
		if err := ensureAllGoalsEvaluated(ctx, q, evaluationID); err != nil {
			return &goalEvaluationSubmitBlockedError{message: err.Error()}
		}

		completedStatus := db.EvaluationStatusEnumCompleted
		evaluation, err = q.UpdateGoalEvaluation(ctx, db.UpdateGoalEvaluationParams{
			ID:     evaluationID,
			Status: &completedStatus,
		})
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "P0001" {
				return &goalEvaluationSubmitBlockedError{message: fmt.Sprintf("cannot submit evaluation yet: %s", pgErr.Message)}
			}
			return fmt.Errorf("failed to submit evaluation: %w", err)
		}

		items, err = q.GetGoalEvaluationItems(ctx, evaluationID)
		if err != nil {
			return fmt.Errorf("failed to reload evaluation items: %w", err)
		}
		return nil
	})
	if err != nil {
		var blockedErr *goalEvaluationSubmitBlockedError
		if !errors.As(err, &blockedErr) {
			return nil, err
		}

		result, loadErr := r.GetGoalEvaluation(ctx, evaluationID)
		if loadErr != nil {
			return nil, loadErr
		}
		message := blockedErr.Error()
		result.SubmitError = &message
		return result, nil
	}

	return toDomainGoalEvaluation(evaluation, items, nil), nil
}

func ensureCurrentCycleEvaluation(ctx context.Context, q *db.Queries, clientID uuid.UUID, evaluationDate pgtype.Date) error {
	client, err := q.GetClientDetails(ctx, clientID)
	if err != nil {
		return fmt.Errorf("failed to get client details: %w", err)
	}
	if !evaluationDate.Valid || !client.NextEvaluationDate.Valid || !evaluationDate.Time.Equal(client.NextEvaluationDate.Time) {
		return domain.ErrGoalEvaluationNotCurrentCycle
	}
	return nil
}

func saveGoalEvaluationDraftItems(ctx context.Context, q *db.Queries, evaluationID, clientID uuid.UUID, activeGoals []db.ClientGoal, requestItems []domain.GoalEvaluationItemParams) ([]db.GetGoalEvaluationItemsRow, error) {
	requestByGoalID, err := mapDraftItemsByGoalID(requestItems)
	if err != nil {
		return nil, err
	}

	activeGoalIDs := make(map[uuid.UUID]struct{}, len(activeGoals))
	for _, goal := range activeGoals {
		activeGoalIDs[goal.ID] = struct{}{}
	}

	existingItems, err := q.GetGoalEvaluationItems(ctx, evaluationID)
	if err != nil {
		return nil, fmt.Errorf("failed to load evaluation items: %w", err)
	}
	existingByGoalID := make(map[uuid.UUID]struct{}, len(existingItems))
	for _, item := range existingItems {
		existingByGoalID[item.GoalID] = struct{}{}
	}

	for _, goal := range activeGoals {
		if _, exists := existingByGoalID[goal.ID]; exists {
			continue
		}
		progress := db.ClientGoalProgressEnumNoProgress
		var notes *string
		if requestItem, ok := requestByGoalID[goal.ID]; ok {
			progress, err = parseProgress(requestItem.Progress)
			if err != nil {
				return nil, err
			}
			notes = requestItem.Notes
		}
		if _, err := q.UpsertGoalEvaluationItem(ctx, db.UpsertGoalEvaluationItemParams{
			ClientID: clientID, EvaluationID: evaluationID, GoalID: goal.ID, Progress: progress, Notes: notes,
		}); err != nil {
			return nil, fmt.Errorf("failed to ensure evaluation item for goal %s: %w", goal.ID, err)
		}
	}

	for _, requestItem := range requestItems {
		if _, active := activeGoalIDs[requestItem.GoalID]; !active {
			return nil, fmt.Errorf("goal %s is not an active goal for this client", requestItem.GoalID)
		}
		progress, err := parseProgress(requestItem.Progress)
		if err != nil {
			return nil, err
		}
		if _, err := q.UpsertGoalEvaluationItem(ctx, db.UpsertGoalEvaluationItemParams{
			ClientID: clientID, EvaluationID: evaluationID, GoalID: requestItem.GoalID, Progress: progress, Notes: requestItem.Notes,
		}); err != nil {
			return nil, fmt.Errorf("failed to save item for goal %s: %w", requestItem.GoalID, err)
		}
	}

	items, err := q.GetGoalEvaluationItems(ctx, evaluationID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload evaluation items: %w", err)
	}
	return items, nil
}

type goalEvaluationSubmitBlockedError struct {
	message string
}

func (e *goalEvaluationSubmitBlockedError) Error() string {
	return e.message
}

func (r *ClientRepository) trySubmitGoalEvaluationDraft(ctx context.Context, evaluationID uuid.UUID) (db.ClientGoalEvaluation, error) {
	var evaluation db.ClientGoalEvaluation

	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		if err := ensureAllGoalsEvaluated(ctx, q, evaluationID); err != nil {
			return &goalEvaluationSubmitBlockedError{message: err.Error()}
		}

		completedStatus := db.EvaluationStatusEnumCompleted
		updatedEvaluation, err := q.UpdateGoalEvaluation(ctx, db.UpdateGoalEvaluationParams{
			ID:     evaluationID,
			Status: &completedStatus,
		})
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "P0001" {
				return &goalEvaluationSubmitBlockedError{message: fmt.Sprintf("cannot submit evaluation yet: %s", pgErr.Message)}
			}
			return fmt.Errorf("failed to submit evaluation: %w", err)
		}

		evaluation = updatedEvaluation
		return nil
	})
	if err != nil {
		return db.ClientGoalEvaluation{}, err
	}

	return evaluation, nil
}

func ensureAllGoalsEvaluated(ctx context.Context, q *db.Queries, evaluationID uuid.UUID) error {
	items, err := q.GetGoalEvaluationItems(ctx, evaluationID)
	if err != nil {
		return fmt.Errorf("failed to load evaluation items: %w", err)
	}
	if len(items) == 0 {
		return fmt.Errorf("evaluation has no goal items")
	}

	for _, item := range items {
		if item.Progress == db.ClientGoalProgressEnumNoProgress {
			return fmt.Errorf("all goals must be evaluated before submit (missing goal: %s)", item.GoalTitle)
		}
	}

	return nil
}

func mapDraftItemsByGoalID(items []domain.GoalEvaluationItemParams) (map[uuid.UUID]domain.GoalEvaluationItemParams, error) {
	itemsByGoal := make(map[uuid.UUID]domain.GoalEvaluationItemParams, len(items))
	for _, item := range items {
		if _, exists := itemsByGoal[item.GoalID]; exists {
			return nil, fmt.Errorf("goal %s appears multiple times in request", item.GoalID)
		}
		itemsByGoal[item.GoalID] = item
	}
	return itemsByGoal, nil
}

func parseProgress(value string) (db.ClientGoalProgressEnum, error) {
	progress := strings.TrimSpace(value)
	if progress == "" {
		return db.ClientGoalProgressEnumNoProgress, nil
	}

	allowed := map[string]db.ClientGoalProgressEnum{
		string(db.ClientGoalProgressEnumNoProgress):      db.ClientGoalProgressEnumNoProgress,
		string(db.ClientGoalProgressEnumRegression):      db.ClientGoalProgressEnumRegression,
		string(db.ClientGoalProgressEnumLimitedProgress): db.ClientGoalProgressEnumLimitedProgress,
		string(db.ClientGoalProgressEnumGoodProgress):    db.ClientGoalProgressEnumGoodProgress,
		string(db.ClientGoalProgressEnumAchieved):        db.ClientGoalProgressEnumAchieved,
		string(db.ClientGoalProgressEnumBlocked):         db.ClientGoalProgressEnumBlocked,
	}

	parsed, ok := allowed[progress]
	if !ok {
		return "", fmt.Errorf("invalid progress value: %s", value)
	}
	return parsed, nil
}

// GetGoalEvaluationBootstrap returns bootstrap data for the goal evaluation page.
func (r *ClientRepository) GetGoalEvaluationBootstrap(ctx context.Context, clientID uuid.UUID) (*domain.GoalEvaluationBootstrap, error) {
	tx, err := r.store.BeginActorTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	q := r.store.WithTx(tx)

	client, err := q.GetClientDetails(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get client details: %w", err)
	}

	response := &domain.GoalEvaluationBootstrap{
		ClientID:        client.ID,
		ClientFirstName: client.FirstName,
		ClientLastName:  client.LastName,
		ActiveGoals:     []domain.GoalEvaluationBootstrapActiveGoal{},
	}

	if client.NextEvaluationDate.Valid {
		nextDate := client.NextEvaluationDate.Time
		response.NextEvaluationDate = &nextDate
		today := time.Now().UTC().Truncate(24 * time.Hour)
		due := nextDate.UTC().Truncate(24 * time.Hour)
		daysLeft := int32(due.Sub(today).Hours() / 24)
		response.DaysLeft = &daysLeft
		priority := "normal"
		if daysLeft <= 3 {
			priority = "critical"
		}
		response.Priority = &priority

		currentDraft, err := q.GetDraftGoalEvaluationByClientAndDate(ctx, db.GetDraftGoalEvaluationByClientAndDateParams{
			ClientID:       clientID,
			EvaluationDate: client.NextEvaluationDate,
		})
		if err == nil {
			response.ExistingDraft = &domain.GoalEvaluationBootstrapDraft{
				ID:             currentDraft.ID,
				EvaluationDate: conv.TimeFromPgDate(currentDraft.EvaluationDate),
				UpdatedAt:      conv.TimeFromPgTimestamptz(currentDraft.UpdatedAt),
			}
		} else if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("failed to get current cycle draft evaluation: %w", err)
		}
	}

	latestCompleted, err := q.GetLatestCompletedEvaluationByClient(ctx, clientID)
	if err == nil {
		response.LastCompletedEvaluation = &domain.GoalEvaluationBootstrapCompleted{
			ID:                  latestCompleted.ID,
			EvaluationDate:      conv.TimeFromPgDate(latestCompleted.EvaluationDate),
			SubmittedAt:         conv.TimeFromPgTimestamptz(latestCompleted.SubmittedAt),
			OverallNotes:        latestCompleted.OverallNotes,
			CreatedByEmployeeID: latestCompleted.CreatedByEmployeeID,
			CreatorName:         composeCreatorName(latestCompleted.CreatorFirstName, latestCompleted.CreatorLastName),
		}
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to get latest completed evaluation: %w", err)
	}

	activeGoals, err := q.ListActiveGoalsByClientID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to list active goals: %w", err)
	}

	latestProgressRows, err := q.ListLatestCompletedGoalProgressByClient(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to list latest completed goal progress: %w", err)
	}

	latestProgressByGoal := make(map[uuid.UUID]db.ListLatestCompletedGoalProgressByClientRow, len(latestProgressRows))
	for _, row := range latestProgressRows {
		latestProgressByGoal[row.GoalID] = row
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	response.ActiveGoals = make([]domain.GoalEvaluationBootstrapActiveGoal, 0, len(activeGoals))
	for _, goal := range activeGoals {
		goalResponse := domain.GoalEvaluationBootstrapActiveGoal{
			GoalID:            goal.ID,
			Title:             goal.Title,
			TopicNameSnapshot: goal.TopicNameSnapshot,
			Priority:          string(goal.Priority),
			SortOrder:         goal.SortOrder,
		}

		if latestGoalProgress, ok := latestProgressByGoal[goal.ID]; ok {
			p := string(latestGoalProgress.Progress)
			goalResponse.LastProgress = &p
			goalResponse.LastNotes = latestGoalProgress.Notes
		}

		response.ActiveGoals = append(response.ActiveGoals, goalResponse)
	}

	return response, nil
}

// ListClientSubmittedEvaluations returns list of submitted evaluations for a client.
func (r *ClientRepository) ListClientSubmittedEvaluations(ctx context.Context, params domain.ListClientSubmittedEvaluationsParams) (*domain.ListClientSubmittedEvaluationsResult, error) {
	rows, err := actorQuery(ctx, r.store, func(q *db.Queries) ([]db.ListSubmittedEvaluationsByClientRow, error) {
		return q.ListSubmittedEvaluationsByClient(ctx, db.ListSubmittedEvaluationsByClientParams{
			ClientID: params.ClientID,
			Limit:    params.Limit,
			Offset:   params.Offset,
		})
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list submitted evaluations: %w", err)
	}

	if len(rows) == 0 {
		return &domain.ListClientSubmittedEvaluationsResult{
			Items:      []domain.ListClientSubmittedEvaluationsItem{},
			TotalCount: 0,
		}, nil
	}

	items := make([]domain.ListClientSubmittedEvaluationsItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, domain.ListClientSubmittedEvaluationsItem{
			EvaluationID:        row.ID,
			EvaluationDate:      conv.TimeFromPgDate(row.EvaluationDate),
			SubmittedAt:         conv.TimeFromPgTimestamptz(row.SubmittedAt),
			FilledGoalsCount:    row.FilledGoalsCount,
			TotalGoalsCount:     row.TotalGoalsCount,
			CreatedByEmployeeID: row.CreatedByEmployeeID,
			CreatorName:         composeCreatorName(row.CreatorFirstName, row.CreatorLastName),
		})
	}

	return &domain.ListClientSubmittedEvaluationsResult{
		Items:      items,
		TotalCount: rows[0].TotalCount,
	}, nil
}

// ListGoalEvaluationHistory returns history of goal evaluations.
func (r *ClientRepository) ListGoalEvaluationHistory(ctx context.Context, params domain.ListGoalEvaluationHistoryParams) (*domain.ListGoalEvaluationHistoryResult, error) {
	rows, err := actorQuery(ctx, r.store, func(q *db.Queries) ([]db.ListGoalEvaluationHistoryByClientAndGoalRow, error) {
		return q.ListGoalEvaluationHistoryByClientAndGoal(ctx, db.ListGoalEvaluationHistoryByClientAndGoalParams{
			ClientID: params.ClientID,
			GoalID:   params.GoalID,
			Limit:    params.Limit,
			Offset:   params.Offset,
		})
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list goal evaluation history: %w", err)
	}

	if len(rows) == 0 {
		return &domain.ListGoalEvaluationHistoryResult{
			Items:      []domain.ListGoalEvaluationHistoryItem{},
			TotalCount: 0,
		}, nil
	}

	items := make([]domain.ListGoalEvaluationHistoryItem, 0, len(rows))
	for _, row := range rows {
		item := domain.ListGoalEvaluationHistoryItem{
			EvaluationID:        row.EvaluationID,
			EvaluationDate:      conv.TimeFromPgDate(row.EvaluationDate),
			SubmittedAt:         conv.TimeFromPgTimestamptz(row.SubmittedAt),
			Progress:            string(row.Progress),
			Notes:               row.Notes,
			CreatedByEmployeeID: row.CreatedByEmployeeID,
			CreatorName:         composeCreatorName(row.CreatorFirstName, row.CreatorLastName),
		}
		if row.PeriodStart.Valid {
			t := conv.TimeFromPgDate(row.PeriodStart)
			item.PeriodStart = &t
		}
		if row.PeriodEnd.Valid {
			t := conv.TimeFromPgDate(row.PeriodEnd)
			item.PeriodEnd = &t
		}
		items = append(items, item)
	}

	return &domain.ListGoalEvaluationHistoryResult{
		Items:      items,
		TotalCount: rows[0].TotalCount,
	}, nil
}

// Helper functions

func toDomainClientGoal(goal db.ClientGoal) *domain.ClientGoal {
	return &domain.ClientGoal{
		ID:                goal.ID,
		ClientID:          goal.ClientID,
		Title:             goal.Title,
		Description:       goal.Description,
		Priority:          string(goal.Priority),
		Status:            string(goal.Status),
		TopicID:           goal.TopicID,
		TopicNameSnapshot: goal.TopicNameSnapshot,
		Source:            string(goal.Source),
		SortOrder:         goal.SortOrder,
		CreatedAt:         conv.TimeFromPgTimestamptz(goal.CreatedAt),
		UpdatedAt:         conv.TimeFromPgTimestamptz(goal.UpdatedAt),
	}
}

func toDomainGoalEvaluation(eval db.ClientGoalEvaluation, items []db.GetGoalEvaluationItemsRow, creatorName *string) *domain.GoalEvaluation {
	res := &domain.GoalEvaluation{
		ID:                      eval.ID,
		ClientID:                eval.ClientID,
		EvaluationDate:          conv.TimeFromPgDate(eval.EvaluationDate),
		EvaluationIntervalWeeks: eval.EvaluationIntervalWeeks,
		Status:                  string(eval.Status),
		OverallNotes:            eval.OverallNotes,
		CreatedByEmployeeID:     eval.CreatedByEmployeeID,
		CreatorName:             creatorName,
		CreatedAt:               conv.TimeFromPgTimestamptz(eval.CreatedAt),
		UpdatedAt:               conv.TimeFromPgTimestamptz(eval.UpdatedAt),
	}

	if eval.PeriodStart.Valid {
		res.PeriodStart = conv.TimePtrFromPgDate(eval.PeriodStart)
	}
	if eval.PeriodEnd.Valid {
		res.PeriodEnd = conv.TimePtrFromPgDate(eval.PeriodEnd)
	}

	res.Items = make([]domain.GoalEvaluationItem, 0, len(items))
	for _, item := range items {
		res.Items = append(res.Items, domain.GoalEvaluationItem{
			ID:                item.ID,
			EvaluationID:      item.EvaluationID,
			GoalID:            item.GoalID,
			GoalTitle:         item.GoalTitle,
			GoalDescription:   item.GoalDescription,
			TopicNameSnapshot: item.TopicNameSnapshot,
			Progress:          string(item.Progress),
			Notes:             item.Notes,
			CreatedAt:         conv.TimeFromPgTimestamptz(item.CreatedAt),
			UpdatedAt:         conv.TimeFromPgTimestamptz(item.UpdatedAt),
		})
	}

	return res
}

func composeCreatorName(firstName, lastName *string) *string {
	if firstName == nil && lastName == nil {
		return nil
	}
	parts := make([]string, 0, 2)
	if firstName != nil && strings.TrimSpace(*firstName) != "" {
		parts = append(parts, strings.TrimSpace(*firstName))
	}
	if lastName != nil && strings.TrimSpace(*lastName) != "" {
		parts = append(parts, strings.TrimSpace(*lastName))
	}
	if len(parts) == 0 {
		return nil
	}
	name := strings.Join(parts, " ")
	return &name
}

func normalizeOptionalTrimmedString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

func equalOptionalString(a *string, b *string) bool {
	switch {
	case a == nil && b == nil:
		return true
	case a == nil || b == nil:
		return false
	default:
		return *a == *b
	}
}

func equalOptionalUUID(a *uuid.UUID, b *uuid.UUID) bool {
	switch {
	case a == nil && b == nil:
		return true
	case a == nil || b == nil:
		return false
	default:
		return *a == *b
	}
}

func (r *ClientRepository) CreateLocationTransfer(ctx context.Context, clientID uuid.UUID, params domain.CreateLocationTransferParams) error {
	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		return q.CreateClientLocationTransfer(ctx, db.CreateClientLocationTransferParams{
			ClientID:       clientID,
			FromLocationID: &params.FromLocationID,
			ToLocationID:   &params.ToLocationID,
			RequestDate:    conv.PgTimestamptzFromTime(time.Now()),
			NewMentorID:    params.NewMentorID,
			Reason:         &params.Reason,
		})
	})
	if err != nil {
		return fmt.Errorf("failed to create location transfer request: %w", err)
	}
	return nil
}

func (r *ClientRepository) ApproveLocationTransfer(ctx context.Context, employeeID uuid.UUID, params domain.ApproveLocationTransferParams) error {
	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		return q.ApproveOrRejectClientLocationTransfer(ctx, db.ApproveOrRejectClientLocationTransferParams{
			ID:                 params.TransferID,
			Status:             db.ClientLocationTransferStatusEnum(params.Status),
			ApprovedRejectedBy: &employeeID,
		})
	})
	if err != nil {
		return fmt.Errorf("failed to approve location transfer request: %w", err)
	}
	return nil
}

func (r *ClientRepository) ListLocationTransferRequests(ctx context.Context, params domain.ListLocationTransferParams) (*domain.ListLocationTransferResult, error) {
	rows, err := actorQuery(ctx, r.store, func(q *db.Queries) ([]db.ListClientLocationTransferRow, error) {
		return q.ListClientLocationTransfer(ctx, db.ListClientLocationTransferParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list location transfer requests: %w", err)
	}

	if len(rows) == 0 {
		return &domain.ListLocationTransferResult{
			Items:      []domain.LocationTransfer{},
			TotalCount: 0,
		}, nil
	}

	items := make([]domain.LocationTransfer, 0, len(rows))
	for _, row := range rows {
		item := domain.LocationTransfer{
			ID:                 row.ID,
			ClientID:           row.ClientID,
			FromLocationID:     row.FromLocationID,
			ToLocationID:       row.ToLocationID,
			NewMentorID:        row.NewMentorID,
			RequestDate:        conv.TimeFromPgTimestamptz(row.RequestDate),
			Status:             string(row.Status),
			ApprovedRejectedBy: row.ApprovedRejectedBy,
			Reason:             row.Reason,
			MentorFirstName:    row.MentorFirstName,
			MentorLastName:     row.MentorLastName,
		}
		if row.ApprovedRejectedAt.Valid {
			t := conv.TimeFromPgTimestamptz(row.ApprovedRejectedAt)
			item.ApprovedRejectedAt = &t
		}
		items = append(items, item)
	}

	return &domain.ListLocationTransferResult{
		Items:      items,
		TotalCount: rows[0].TotalCount,
	}, nil
}

func (r *ClientRepository) GetGoalEvaluation(ctx context.Context, evaluationID uuid.UUID) (*domain.GoalEvaluation, error) {
	var row db.GetGoalEvaluationByIDRow
	var items []db.GetGoalEvaluationItemsRow
	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		var err error
		row, err = q.GetGoalEvaluationByID(ctx, evaluationID)
		if err != nil {
			return err
		}
		items, err = q.GetGoalEvaluationItems(ctx, evaluationID)
		return err
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrGoalEvaluationNotFound
		}
		return nil, fmt.Errorf("failed to get goal evaluation: %w", err)
	}

	eval := db.ClientGoalEvaluation{
		ID:                      row.ID,
		ClientID:                row.ClientID,
		EvaluationDate:          row.EvaluationDate,
		PeriodStart:             row.PeriodStart,
		PeriodEnd:               row.PeriodEnd,
		EvaluationIntervalWeeks: row.EvaluationIntervalWeeks,
		Status:                  row.Status,
		OverallNotes:            row.OverallNotes,
		CreatedByEmployeeID:     row.CreatedByEmployeeID,
		CreatedAt:               row.CreatedAt,
		UpdatedAt:               row.UpdatedAt,
	}

	res := toDomainGoalEvaluation(eval, items, composeCreatorName(row.CreatorFirstName, row.CreatorLastName))
	return res, nil
}

func (r *ClientRepository) ListUpcomingEvaluations(ctx context.Context, params domain.ListUpcomingEvaluationsParams) (*domain.ListUpcomingEvaluationsResult, error) {
	rows, err := actorQuery(ctx, r.store, func(q *db.Queries) ([]db.ListUpcomingEvaluationsForCoordinatorRow, error) {
		return q.ListUpcomingEvaluationsForCoordinator(ctx, db.ListUpcomingEvaluationsForCoordinatorParams{
			EmployeeID: params.EmployeeID,
			Limit:      params.Limit,
			Offset:     params.Offset,
		})
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list upcoming evaluations: %w", err)
	}

	if len(rows) == 0 {
		return &domain.ListUpcomingEvaluationsResult{
			Items:      []domain.UpcomingEvaluation{},
			TotalCount: 0,
		}, nil
	}

	items := make([]domain.UpcomingEvaluation, 0, len(rows))
	for _, row := range rows {
		items = append(items, domain.UpcomingEvaluation{
			ClientID:         row.ClientID,
			ClientFirstName:  row.ClientFirstName,
			ClientLastName:   row.ClientLastName,
			DueDate:          conv.TimeFromPgDate(row.NextEvaluationDate),
			DaysLeft:         row.DaysLeft,
			Priority:         row.Priority,
			HasDraft:         row.HasDraft,
			FilledGoalsCount: row.FilledGoalsCount,
			TotalGoalsCount:  row.TotalGoalsCount,
		})
	}

	return &domain.ListUpcomingEvaluationsResult{
		Items:      items,
		TotalCount: rows[0].TotalCount,
	}, nil
}

func (r *ClientRepository) ListRecentSubmittedEvaluations(ctx context.Context, params domain.ListRecentSubmittedEvaluationsParams) (*domain.ListRecentSubmittedEvaluationsResult, error) {
	rows, err := actorQuery(ctx, r.store, func(q *db.Queries) ([]db.ListRecentSubmittedEvaluationsByEmployeeRow, error) {
		return q.ListRecentSubmittedEvaluationsByEmployee(ctx, db.ListRecentSubmittedEvaluationsByEmployeeParams{
			CreatedByEmployeeID: &params.EmployeeID,
			Limit:               params.Limit,
			Offset:              params.Offset,
		})
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list recent submitted evaluations: %w", err)
	}

	if len(rows) == 0 {
		return &domain.ListRecentSubmittedEvaluationsResult{
			Items:      []domain.RecentSubmittedEvaluation{},
			TotalCount: 0,
		}, nil
	}

	items := make([]domain.RecentSubmittedEvaluation, 0, len(rows))
	for _, row := range rows {
		item := domain.RecentSubmittedEvaluation{
			EvaluationID:     row.ID,
			ClientID:         row.ClientID,
			ClientFirstName:  row.ClientFirstName,
			ClientLastName:   row.ClientLastName,
			EvaluationDate:   conv.TimeFromPgDate(row.EvaluationDate),
			SubmittedAt:      conv.TimeFromPgTimestamptz(row.SubmittedAt),
			FilledGoalsCount: row.FilledGoalsCount,
			TotalGoalsCount:  row.TotalGoalsCount,
		}
		if row.NextEvaluationDate.Valid {
			t := conv.TimeFromPgDate(row.NextEvaluationDate)
			item.NextEvaluationDate = &t
		}
		items = append(items, item)
	}

	return &domain.ListRecentSubmittedEvaluationsResult{
		Items:      items,
		TotalCount: rows[0].TotalCount,
	}, nil
}

func (r *ClientRepository) ListRecentDraftEvaluations(ctx context.Context, params domain.ListRecentDraftEvaluationsParams) (*domain.ListRecentDraftEvaluationsResult, error) {
	rows, err := actorQuery(ctx, r.store, func(q *db.Queries) ([]db.ListRecentDraftEvaluationsByEmployeeRow, error) {
		return q.ListRecentDraftEvaluationsByEmployee(ctx, db.ListRecentDraftEvaluationsByEmployeeParams{
			CreatedByEmployeeID: &params.EmployeeID,
			Limit:               params.Limit,
			Offset:              params.Offset,
		})
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list recent draft evaluations: %w", err)
	}

	if len(rows) == 0 {
		return &domain.ListRecentDraftEvaluationsResult{
			Items:      []domain.RecentDraftEvaluation{},
			TotalCount: 0,
		}, nil
	}

	items := make([]domain.RecentDraftEvaluation, 0, len(rows))
	for _, row := range rows {
		items = append(items, domain.RecentDraftEvaluation{
			EvaluationID:     row.ID,
			ClientID:         row.ClientID,
			ClientFirstName:  row.ClientFirstName,
			ClientLastName:   row.ClientLastName,
			DueDate:          conv.TimeFromPgDate(row.EvaluationDate),
			UpdatedAt:        conv.TimeFromPgTimestamptz(row.UpdatedAt),
			DaysLeft:         row.DaysLeft,
			Priority:         row.Priority,
			FilledGoalsCount: row.FilledGoalsCount,
			TotalGoalsCount:  row.TotalGoalsCount,
		})
	}

	return &domain.ListRecentDraftEvaluationsResult{
		Items:      items,
		TotalCount: rows[0].TotalCount,
	}, nil
}

// =====================
// Medical - Diagnoses
// =====================

func (r *ClientRepository) CreateClientDiagnosis(ctx context.Context, params domain.CreateClientDiagnosisParams) (*domain.ClientDiagnosis, error) {
	var diagnosis db.ClientDiagnosis

	status := db.DiagnosisStatusEnum("confirmed")
	if params.Status != nil {
		status = db.DiagnosisStatusEnum(*params.Status)
	}

	severity := db.DiagnosisSeverityEnum("unknown")
	if params.Severity != nil {
		severity = db.DiagnosisSeverityEnum(*params.Severity)
	}

	var diagnosedOn pgtype.Date
	if params.DiagnosedOn != nil {
		diagnosedOn = pgtype.Date{Time: *params.DiagnosedOn, Valid: true}
	} else {
		diagnosedOn = pgtype.Date{Valid: false}
	}

	var resolvedOn pgtype.Date
	if params.ResolvedOn != nil {
		resolvedOn = pgtype.Date{Time: *params.ResolvedOn, Valid: true}
	} else {
		resolvedOn = pgtype.Date{Valid: false}
	}

	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		var err error
		diagnosis, err = q.CreateClientDiagnosis(ctx, db.CreateClientDiagnosisParams{
			ClientID:            params.ClientID,
			CodeSystem:          params.CodeSystem,
			Code:                params.Code,
			Title:               params.Title,
			Description:         params.Description,
			Status:              status,
			Severity:            severity,
			DiagnosedOn:         diagnosedOn,
			ResolvedOn:          resolvedOn,
			DiagnosingClinician: params.DiagnosingClinician,
			Notes:               params.Notes,
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	return toDomainClientDiagnosis(diagnosis), nil
}

func (r *ClientRepository) ListClientDiagnoses(ctx context.Context, params domain.ListClientDiagnosesParams) (*domain.ListClientDiagnosesResult, error) {
	rows, err := actorQuery(ctx, r.store, func(q *db.Queries) ([]db.ListClientDiagnosesRow, error) {
		return q.ListClientDiagnoses(ctx, db.ListClientDiagnosesParams{
			ClientID: params.ClientID,
			Limit:    params.Limit,
			Offset:   params.Offset,
		})
	})
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return &domain.ListClientDiagnosesResult{
			Items:      []domain.ClientDiagnosis{},
			TotalCount: 0,
		}, nil
	}

	items := make([]domain.ClientDiagnosis, 0, len(rows))
	for _, row := range rows {
		items = append(items, toDomainClientDiagnosisFromRow(row))
	}

	return &domain.ListClientDiagnosesResult{
		Items:      items,
		TotalCount: rows[0].TotalCount,
	}, nil
}

func (r *ClientRepository) GetClientDiagnosis(ctx context.Context, clientID, diagnosisID uuid.UUID) (*domain.ClientDiagnosis, error) {
	diagnosis, err := actorQuery(ctx, r.store, func(q *db.Queries) (db.ClientDiagnosis, error) {
		return q.GetClientDiagnosis(ctx, db.GetClientDiagnosisParams{
			ClientID: clientID,
			ID:       diagnosisID,
		})
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("diagnosis not found")
		}
		return nil, err
	}

	return toDomainClientDiagnosis(diagnosis), nil
}

func (r *ClientRepository) UpdateClientDiagnosis(ctx context.Context, params domain.UpdateClientDiagnosisParams) (*domain.ClientDiagnosis, error) {
	var diagnosis db.ClientDiagnosis

	var diagnosedOn pgtype.Date
	if params.DiagnosedOn != nil {
		diagnosedOn = pgtype.Date{Time: *params.DiagnosedOn, Valid: true}
	} else {
		diagnosedOn = pgtype.Date{Valid: false}
	}

	var resolvedOn pgtype.Date
	if params.ResolvedOn != nil {
		resolvedOn = pgtype.Date{Time: *params.ResolvedOn, Valid: true}
	} else {
		resolvedOn = pgtype.Date{Valid: false}
	}

	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		var err error
		diagnosis, err = q.UpdateClientDiagnosis(ctx, db.UpdateClientDiagnosisParams{
			CodeSystem:          params.CodeSystem,
			Code:                params.Code,
			Title:               params.Title,
			Description:         params.Description,
			Status:              diagnosisStatusFromPtr(params.Status),
			Severity:            diagnosisSeverityFromPtr(params.Severity),
			DiagnosedOn:         diagnosedOn,
			ResolvedOn:          resolvedOn,
			DiagnosingClinician: params.DiagnosingClinician,
			Notes:               params.Notes,
			ClientID:            params.ClientID,
			ID:                  params.ID,
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	return toDomainClientDiagnosis(diagnosis), nil
}

func (r *ClientRepository) DeleteClientDiagnosis(ctx context.Context, clientID, diagnosisID uuid.UUID) (*domain.DeleteClientDiagnosisResult, error) {
	var deleted db.ClientDiagnosis

	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		var err error
		deleted, err = q.DeleteClientDiagnosis(ctx, db.DeleteClientDiagnosisParams{
			ClientID: clientID,
			ID:       diagnosisID,
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	return &domain.DeleteClientDiagnosisResult{ID: deleted.ID}, nil
}

// =====================
// Medical - Medication Orders
// =====================

func (r *ClientRepository) CreateClientMedicationOrder(ctx context.Context, params domain.CreateClientMedicationOrderParams) (*domain.ClientMedicationOrder, error) {
	var order db.ClientMedicationOrder

	status := db.MedicationOrderStatusEnum("active")
	if params.Status != nil {
		status = db.MedicationOrderStatusEnum(*params.Status)
	}

	adminMode := db.MedicationAdminModeEnum("self")
	if params.AdminMode != nil {
		adminMode = db.MedicationAdminModeEnum(*params.AdminMode)
	}

	startDate := conv.PgDateFromTime(params.StartDate)

	var endDate pgtype.Date
	if params.EndDate != nil {
		endDate = pgtype.Date{Time: *params.EndDate, Valid: true}
	} else {
		endDate = pgtype.Date{Valid: false}
	}

	schedule := params.Schedule
	if len(schedule) == 0 {
		schedule = []byte("[]")
	}

	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		var err error
		order, err = q.CreateClientMedicationOrder(ctx, db.CreateClientMedicationOrderParams{
			ClientID:              params.ClientID,
			DiagnosisID:           params.DiagnosisID,
			MedicationName:        params.MedicationName,
			DosageText:            params.DosageText,
			DoseAmount:            params.DoseAmount,
			DoseUnit:              params.DoseUnit,
			Route:                 params.Route,
			FrequencyText:         params.FrequencyText,
			Schedule:              schedule,
			IsPrn:                 params.IsPrn,
			PrnIndication:         params.PrnIndication,
			MaxDosesPer24h:        params.MaxDosesPer24h,
			StartDate:             startDate,
			EndDate:               endDate,
			Status:                status,
			AdminMode:             adminMode,
			ResponsibleEmployeeID: params.ResponsibleEmployeeID,
			IsCritical:            params.IsCritical,
			Notes:                 params.Notes,
			SourceAttachmentUuid:  params.SourceAttachmentUUID,
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	return toDomainClientMedicationOrder(order), nil
}

func (r *ClientRepository) ListClientMedicationOrders(ctx context.Context, params domain.ListClientMedicationOrdersParams) (*domain.ListClientMedicationOrdersResult, error) {
	rows, err := actorQuery(ctx, r.store, func(q *db.Queries) ([]db.ListClientMedicationOrdersRow, error) {
		return q.ListClientMedicationOrders(ctx, db.ListClientMedicationOrdersParams{
			ClientID:    params.ClientID,
			Status:      medicationOrderStatusFromPtr(params.Status),
			AdminMode:   medicationAdminModeFromPtr(params.AdminMode),
			DiagnosisID: params.DiagnosisID,
			Search:      params.Search,
			Offset:      params.Offset,
			Limit:       params.Limit,
		})
	})
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return &domain.ListClientMedicationOrdersResult{
			Items:      []domain.ClientMedicationOrder{},
			TotalCount: 0,
		}, nil
	}

	items := make([]domain.ClientMedicationOrder, 0, len(rows))
	for _, row := range rows {
		items = append(items, toDomainClientMedicationOrderFromRow(row))
	}

	return &domain.ListClientMedicationOrdersResult{
		Items:      items,
		TotalCount: rows[0].TotalCount,
	}, nil
}

func (r *ClientRepository) GetClientMedicationOrder(ctx context.Context, clientID, orderID uuid.UUID) (*domain.ClientMedicationOrder, error) {
	row, err := actorQuery(ctx, r.store, func(q *db.Queries) (db.GetClientMedicationOrderRow, error) {
		return q.GetClientMedicationOrder(ctx, db.GetClientMedicationOrderParams{
			ClientID: clientID,
			ID:       orderID,
		})
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("medication order not found")
		}
		return nil, err
	}

	return toDomainClientMedicationOrderFromGetRow(row), nil
}

func (r *ClientRepository) UpdateClientMedicationOrder(ctx context.Context, params domain.UpdateClientMedicationOrderParams) (*domain.ClientMedicationOrder, error) {
	var order db.ClientMedicationOrder

	var startDate pgtype.Date
	if params.StartDate != nil {
		startDate = pgtype.Date{Time: *params.StartDate, Valid: true}
	} else {
		startDate = pgtype.Date{Valid: false}
	}

	var endDate pgtype.Date
	if params.EndDate != nil {
		endDate = pgtype.Date{Time: *params.EndDate, Valid: true}
	} else {
		endDate = pgtype.Date{Valid: false}
	}

	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		var err error
		order, err = q.UpdateClientMedicationOrder(ctx, db.UpdateClientMedicationOrderParams{
			DiagnosisID:           params.DiagnosisID,
			MedicationName:        params.MedicationName,
			DosageText:            params.DosageText,
			DoseAmount:            params.DoseAmount,
			DoseUnit:              params.DoseUnit,
			Route:                 params.Route,
			FrequencyText:         params.FrequencyText,
			Schedule:              params.Schedule,
			IsPrn:                 params.IsPrn,
			PrnIndication:         params.PrnIndication,
			MaxDosesPer24h:        params.MaxDosesPer24h,
			StartDate:             startDate,
			EndDate:               endDate,
			Status:                medicationOrderStatusFromPtr(params.Status),
			AdminMode:             medicationAdminModeFromPtr(params.AdminMode),
			ResponsibleEmployeeID: params.ResponsibleEmployeeID,
			IsCritical:            params.IsCritical,
			Notes:                 params.Notes,
			SourceAttachmentUuid:  params.SourceAttachmentUUID,
			ClientID:              params.ClientID,
			ID:                    params.ID,
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	return toDomainClientMedicationOrder(order), nil
}

func (r *ClientRepository) DeleteClientMedicationOrder(ctx context.Context, clientID, orderID uuid.UUID) (*domain.DeleteClientMedicationOrderResult, error) {
	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		rowsAffected, err := q.DeleteClientMedicationOrder(ctx, db.DeleteClientMedicationOrderParams{
			ClientID: clientID,
			ID:       orderID,
		})
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return pgx.ErrNoRows
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrClientMedicationOrderNotFound
		}
		return nil, err
	}

	return &domain.DeleteClientMedicationOrderResult{ID: orderID}, nil
}

// =====================
// Medical - Overview
// =====================

func (r *ClientRepository) GetClientMedicalOverview(ctx context.Context, clientID uuid.UUID) (*domain.ClientMedicalOverview, error) {
	var diagnoses []db.ListClientDiagnosesRow
	var medicationOrders []db.ListClientMedicationOrdersRow

	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		var err error
		diagnoses, err = q.ListClientDiagnoses(ctx, db.ListClientDiagnosesParams{
			ClientID: clientID,
			Limit:    500,
			Offset:   0,
		})
		if err != nil {
			return err
		}

		activeStatus := db.MedicationOrderStatusEnumActive
		medicationOrders, err = q.ListClientMedicationOrders(ctx, db.ListClientMedicationOrdersParams{
			ClientID:    clientID,
			Status:      &activeStatus,
			AdminMode:   nil,
			DiagnosisID: nil,
			Search:      nil,
			Offset:      0,
			Limit:       500,
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	domainDiagnoses := make([]domain.ClientDiagnosis, 0, len(diagnoses))
	for _, d := range diagnoses {
		domainDiagnoses = append(domainDiagnoses, toDomainClientDiagnosisFromRow(d))
	}

	domainMedicationOrders := make([]domain.ClientMedicationOrder, 0, len(medicationOrders))
	for _, mo := range medicationOrders {
		domainMedicationOrders = append(domainMedicationOrders, toDomainClientMedicationOrderFromRow(mo))
	}

	return &domain.ClientMedicalOverview{
		Diagnoses:        domainDiagnoses,
		MedicationOrders: domainMedicationOrders,
	}, nil
}

// =====================
// Network - Sender
// =====================

func (r *ClientRepository) GetClientSender(ctx context.Context, clientID uuid.UUID) (*domain.Sender, error) {
	sender, err := actorQuery(ctx, r.store, func(q *db.Queries) (db.Sender, error) {
		return q.GetClientSender(ctx, clientID)
	})
	if err != nil {
		return nil, err
	}
	return toDomainSender(sender), nil
}

// =====================
// Network - Emergency Contacts
// =====================

func (r *ClientRepository) CreateClientEmergencyContact(ctx context.Context, params domain.CreateClientEmergencyContactParams) (*domain.ClientEmergencyContact, error) {
	var contact db.ClientEmergencyContact

	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		var err error
		contact, err = q.CreateEmemrgencyContact(ctx, db.CreateEmemrgencyContactParams{
			ClientID:         params.ClientID,
			FirstName:        params.FirstName,
			LastName:         params.LastName,
			Email:            params.Email,
			PhoneNumber:      params.PhoneNumber,
			Address:          params.Address,
			Relationship:     params.Relationship,
			RelationStatus:   db.NullRelationStatusFromPtr(params.RelationStatus),
			MedicalReports:   params.MedicalReports,
			IncidentsReports: params.IncidentsReports,
			GoalsReports:     params.GoalsReports,
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	return toDomainClientEmergencyContact(contact), nil
}

func (r *ClientRepository) ListClientEmergencyContacts(ctx context.Context, params domain.ListClientEmergencyContactsParams) (*domain.ListClientEmergencyContactsResult, error) {
	search := params.Search
	if search == "" {
		search = ""
	}

	rows, err := actorQuery(ctx, r.store, func(q *db.Queries) ([]db.ListEmergencyContactsRow, error) {
		return q.ListEmergencyContacts(ctx, db.ListEmergencyContactsParams{
			ClientID: params.ClientID,
			Search:   search,
			Offset:   params.Offset,
			Limit:    params.Limit,
		})
	})
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return &domain.ListClientEmergencyContactsResult{
			Items:      []domain.ClientEmergencyContact{},
			TotalCount: 0,
		}, nil
	}

	items := make([]domain.ClientEmergencyContact, 0, len(rows))
	for _, row := range rows {
		items = append(items, toDomainClientEmergencyContactFromRow(row))
	}

	return &domain.ListClientEmergencyContactsResult{
		Items:      items,
		TotalCount: rows[0].TotalCount,
	}, nil
}

func (r *ClientRepository) GetClientEmergencyContact(ctx context.Context, contactID uuid.UUID) (*domain.ClientEmergencyContact, error) {
	contact, err := actorQuery(ctx, r.store, func(q *db.Queries) (db.ClientEmergencyContact, error) {
		return q.GetEmergencyContact(ctx, contactID)
	})
	if err != nil {
		return nil, err
	}
	return toDomainClientEmergencyContact(contact), nil
}

func (r *ClientRepository) UpdateClientEmergencyContact(ctx context.Context, params domain.UpdateClientEmergencyContactParams) (*domain.ClientEmergencyContact, error) {
	var contact db.ClientEmergencyContact

	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		var err error
		contact, err = q.UpdateEmergencyContact(ctx, db.UpdateEmergencyContactParams{
			ID:               params.ID,
			FirstName:        params.FirstName,
			LastName:         params.LastName,
			Email:            params.Email,
			PhoneNumber:      params.PhoneNumber,
			Address:          params.Address,
			Relationship:     params.Relationship,
			RelationStatus:   db.NullRelationStatusFromPtr(params.RelationStatus),
			MedicalReports:   params.MedicalReports,
			IncidentsReports: params.IncidentsReports,
			GoalsReports:     params.GoalsReports,
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	return toDomainClientEmergencyContact(contact), nil
}

func (r *ClientRepository) DeleteClientEmergencyContact(ctx context.Context, contactID uuid.UUID) (*domain.DeleteClientEmergencyContactResult, error) {
	var deleted db.ClientEmergencyContact

	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		var err error
		deleted, err = q.DeleteEmergencyContact(ctx, contactID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return &domain.DeleteClientEmergencyContactResult{ID: deleted.ID}, nil
}

// =====================
// Network - Assigned Employees
// =====================

func (r *ClientRepository) CreateAssignedEmployee(ctx context.Context, params domain.CreateAssignedEmployeeParams) (*domain.AssignedEmployee, error) {
	var row db.AssignEmployeeRow

	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		var err error
		row, err = q.AssignEmployee(ctx, db.AssignEmployeeParams{
			ClientID:   params.ClientID,
			EmployeeID: params.EmployeeID,
			StartDate:  conv.PgDateFromTime(params.StartDate),
			Role:       params.Role,
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	return toDomainAssignedEmployeeFromAssignRow(row), nil
}

func (r *ClientRepository) ListAssignedEmployees(ctx context.Context, params domain.ListAssignedEmployeesParams) (*domain.ListAssignedEmployeesResult, error) {
	rows, err := actorQuery(ctx, r.store, func(q *db.Queries) ([]db.ListAssignedEmployeesRow, error) {
		return q.ListAssignedEmployees(ctx, db.ListAssignedEmployeesParams{
			ClientID: params.ClientID,
			Offset:   params.Offset,
			Limit:    params.Limit,
		})
	})
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return &domain.ListAssignedEmployeesResult{
			Items:      []domain.AssignedEmployee{},
			TotalCount: 0,
		}, nil
	}

	items := make([]domain.AssignedEmployee, 0, len(rows))
	for _, row := range rows {
		items = append(items, toDomainAssignedEmployeeFromListRow(row))
	}

	return &domain.ListAssignedEmployeesResult{
		Items:      items,
		TotalCount: rows[0].TotalCount,
	}, nil
}

func (r *ClientRepository) GetAssignedEmployee(ctx context.Context, assignmentID uuid.UUID) (*domain.AssignedEmployee, error) {
	row, err := actorQuery(ctx, r.store, func(q *db.Queries) (db.GetAssignedEmployeeRow, error) {
		return q.GetAssignedEmployee(ctx, assignmentID)
	})
	if err != nil {
		return nil, err
	}
	return toDomainAssignedEmployeeFromGetRow(row), nil
}

func (r *ClientRepository) UpdateAssignedEmployee(ctx context.Context, params domain.UpdateAssignedEmployeeParams) (*domain.AssignedEmployee, error) {
	var row db.AssignedEmployee

	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		var err error
		row, err = q.UpdateAssignedEmployee(ctx, db.UpdateAssignedEmployeeParams{
			ID:         params.ID,
			EmployeeID: params.EmployeeID,
			StartDate:  conv.PgDateFromTime(params.StartDate),
			Role:       params.Role,
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	return toDomainAssignedEmployee(row), nil
}

func (r *ClientRepository) DeleteAssignedEmployee(ctx context.Context, assignmentID uuid.UUID) (*domain.DeleteAssignedEmployeeResult, error) {
	var deleted db.AssignedEmployee

	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		var err error
		deleted, err = q.DeleteAssignedEmployee(ctx, assignmentID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return &domain.DeleteAssignedEmployeeResult{ID: deleted.ID}, nil
}

// =====================
// Network - Related Emails
// =====================

func (r *ClientRepository) GetClientRelatedEmails(ctx context.Context, clientID uuid.UUID) (*domain.ClientRelatedEmails, error) {
	emails, err := actorQuery(ctx, r.store, func(q *db.Queries) ([]string, error) {
		return q.GetClientRelatedEmails(ctx, clientID)
	})
	if err != nil {
		return nil, err
	}
	emailPointers := make([]*string, len(emails))
	for i := range emails {
		emailPointers[i] = &emails[i]
	}
	return &domain.ClientRelatedEmails{Emails: emailPointers}, nil
}

// ListStatusHistory lists client status history records.
func (r *ClientRepository) ListStatusHistory(ctx context.Context, params domain.ListStatusHistoryParams) ([]domain.ClientStatusHistory, error) {
	rows, err := actorQuery(ctx, r.store, func(q *db.Queries) ([]db.ClientStatusHistory, error) {
		return q.ListClientStatusHistory(ctx, db.ListClientStatusHistoryParams{
			ClientID: params.ClientID,
			Limit:    params.Limit,
			Offset:   params.Offset,
		})
	})
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return []domain.ClientStatusHistory{}, nil
	}

	history := make([]domain.ClientStatusHistory, len(rows))
	for i, h := range rows {
		history[i] = domain.ClientStatusHistory{
			ID:        h.ID,
			ClientID:  h.ClientID,
			OldStatus: h.OldStatus,
			NewStatus: h.NewStatus,
			ChangedAt: h.ChangedAt.Time,
			ChangedBy: h.ChangedBy,
			Reason:    h.Reason,
		}
	}

	return history, nil
}

func (r *ClientRepository) UpdateClientStatus(ctx context.Context, clientID uuid.UUID, params domain.UpdateClientStatusParams) (*domain.UpdateClientStatusResult, error) {
	updated, err := actorQuery(ctx, r.store, func(q *db.Queries) (db.ClientDetail, error) {
		return q.UpdateClientStatus(ctx, db.UpdateClientStatusParams{
			ID:     clientID,
			Status: db.ClientStatusEnum(params.Status),
		})
	})
	if err != nil {
		return nil, err
	}
	return &domain.UpdateClientStatusResult{
		ID:     updated.ID,
		Status: string(updated.Status),
	}, nil
}

func (r *ClientRepository) PutClientInCare(ctx context.Context, clientID uuid.UUID, targetStatus db.ClientStatusEnum, careStartDay time.Time, params domain.PutClientInCareParams) (*domain.PutClientInCareResult, error) {
	var result domain.PutClientInCareResult
	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		existingClient, err := q.GetClientDetailsForUpdate(ctx, clientID)
		if err != nil {
			return fmt.Errorf("failed to get client details: %w", err)
		}
		if existingClient.Status != db.ClientStatusEnumOnWaitingList {
			return fmt.Errorf("client must be in on_waiting_list to be put in care")
		}

		activeGoalsCount, err := q.CountActiveGoalsByClientID(ctx, clientID)
		if err != nil {
			return fmt.Errorf("failed to validate active goals: %w", err)
		}
		if activeGoalsCount == 0 {
			return fmt.Errorf("client cannot be put in care without at least one active goal")
		}

		if _, err := q.GetActiveEmployeeForCare(ctx, params.CoordinatorEmployeeID); err != nil {
			return fmt.Errorf("selected coordinator is not available: %w", err)
		}

		coordinatorAssignment, err := q.GetMainCoordinator(ctx, clientID)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("failed to get main coordinator: %w", err)
		}
		if errors.Is(err, pgx.ErrNoRows) {
			coordinatorAssignment, err = q.CreateMainCoordinator(ctx, db.CreateMainCoordinatorParams{
				ClientID:   clientID,
				EmployeeID: params.CoordinatorEmployeeID,
				StartDate:  pgtype.Date{Time: careStartDay, Valid: true},
			})
			if err != nil {
				return fmt.Errorf("failed to assign main coordinator: %w", err)
			}
		} else if coordinatorAssignment.EmployeeID != params.CoordinatorEmployeeID {
			coordinatorAssignment, err = q.UpdateAssignedEmployee(ctx, db.UpdateAssignedEmployeeParams{
				ID:         coordinatorAssignment.ID,
				EmployeeID: &params.CoordinatorEmployeeID,
				StartDate:  pgtype.Date{Time: careStartDay, Valid: true},
			})
			if err != nil {
				return fmt.Errorf("failed to update main coordinator: %w", err)
			}
		}

		placedInCareAt := pgtype.Timestamptz{Valid: false}
		if params.PlacedInCareAt != nil {
			placedInCareAt = pgtype.Timestamptz{Time: *params.PlacedInCareAt, Valid: true}
		}

		updatedClient, err := q.PutClientInCare(ctx, db.PutClientInCareParams{
			ID:             clientID,
			Status:         targetStatus,
			CareStartDate:  pgtype.Date{Time: careStartDay, Valid: true},
			PlacedInCareAt: placedInCareAt,
		})
		if err != nil {
			return fmt.Errorf("failed to put client in care: %w", err)
		}

		historyReason := "client_put_in_care"
		if params.Reason != nil {
			trimmed := strings.TrimSpace(*params.Reason)
			if trimmed != "" {
				historyReason = trimmed
			}
		}

		oldStatus := string(existingClient.Status)
		_, err = q.CreateClientStatusHistory(ctx, db.CreateClientStatusHistoryParams{
			ClientID:  clientID,
			OldStatus: &oldStatus,
			NewStatus: string(updatedClient.Status),
			Reason:    &historyReason,
		})
		if err != nil {
			return fmt.Errorf("failed to create client status history: %w", err)
		}

		result = domain.PutClientInCareResult{
			ID:                      updatedClient.ID,
			Status:                  string(updatedClient.Status),
			CareStartDate:           updatedClient.CareStartDate.Time,
			PlacedInCareAt:          updatedClient.PlacedInCareAt.Time,
			NextEvaluationDate:      nil,
			CoordinatorAssignmentID: coordinatorAssignment.ID,
		}
		if updatedClient.NextEvaluationDate.Valid {
			t := updatedClient.NextEvaluationDate.Time
			result.NextEvaluationDate = &t
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *ClientRepository) PutClientOutOfCare(ctx context.Context, clientID uuid.UUID, targetStatus db.ClientStatusEnum, dischargeDay time.Time, dischargeReason db.DischargeReasonEnum, finalEvaluation *string, params domain.PutClientOutOfCareParams) (*domain.PutClientOutOfCareResult, error) {
	var result domain.PutClientOutOfCareResult
	err := r.store.ExecActorTx(ctx, func(q *db.Queries) error {
		existingClient, err := q.GetClientDetails(ctx, clientID)
		if err != nil {
			return fmt.Errorf("failed to get client details: %w", err)
		}
		if existingClient.Status != db.ClientStatusEnumInCare && existingClient.Status != db.ClientStatusEnumScheduledOutOfCare {
			return fmt.Errorf("client must be in in_care or scheduled_out_of_care to be put out of care")
		}

		updatedClient, err := q.PutClientOutOfCare(ctx, db.PutClientOutOfCareParams{
			ID:              clientID,
			Status:          targetStatus,
			DischargeDate:   pgtype.Date{Time: dischargeDay, Valid: true},
			DischargeReason: &dischargeReason,
			FinalEvaluation: finalEvaluation,
		})
		if err != nil {
			return fmt.Errorf("failed to put client out of care: %w", err)
		}

		historyReason := fmt.Sprintf("client_put_out_of_care:%s", dischargeReason)
		if targetStatus == db.ClientStatusEnumScheduledOutOfCare {
			historyReason = fmt.Sprintf("client_scheduled_out_of_care:%s", dischargeReason)
		}
		if params.Reason != nil {
			trimmed := strings.TrimSpace(*params.Reason)
			if trimmed != "" {
				historyReason = trimmed
			}
		}

		oldStatus := string(existingClient.Status)
		_, err = q.CreateClientStatusHistory(ctx, db.CreateClientStatusHistoryParams{
			ClientID:  clientID,
			OldStatus: &oldStatus,
			NewStatus: string(updatedClient.Status),
			Reason:    &historyReason,
		})
		if err != nil {
			return fmt.Errorf("failed to create client status history: %w", err)
		}

		result = domain.PutClientOutOfCareResult{
			ID:              updatedClient.ID,
			Status:          string(updatedClient.Status),
			DischargeDate:   updatedClient.DischargeDate.Time,
			DischargeReason: string(*updatedClient.DischargeReason),
			FinalEvaluation: updatedClient.FinalEvaluation,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// =====================
// Helper Functions
// =====================

func toDomainSender(d db.Sender) *domain.Sender {
	return &domain.Sender{
		ID:                  d.ID,
		Types:               string(d.Types),
		Name:                d.Name,
		Street:              d.Street,
		HouseNumber:         d.HouseNumber,
		HouseNumberAddition: d.HouseNumberAddition,
		PostalCode:          d.PostalCode,
		City:                d.City,
		Land:                d.Land,
		Kvknumber:           d.Kvknumber,
		Btwnumber:           d.Btwnumber,
		PhoneNumber:         d.PhoneNumber,
		ClientNumber:        d.ClientNumber,
		EmailAddress:        d.EmailAddress,
		Contacts:            d.Contacts,
		IsArchived:          d.IsArchived,
		CreatedAt:           conv.TimeFromPgTimestamptz(d.CreatedAt),
		UpdatedAt:           conv.TimeFromPgTimestamptz(d.UpdatedAt),
	}
}

func toDomainClientEmergencyContact(d db.ClientEmergencyContact) *domain.ClientEmergencyContact {
	var relationStatus *string
	if d.RelationStatus != nil {
		s := string(*d.RelationStatus)
		relationStatus = &s
	}

	return &domain.ClientEmergencyContact{
		ID:               d.ID,
		ClientID:         d.ClientID,
		FirstName:        d.FirstName,
		LastName:         d.LastName,
		Email:            d.Email,
		PhoneNumber:      d.PhoneNumber,
		Address:          d.Address,
		Relationship:     d.Relationship,
		RelationStatus:   relationStatus,
		CreatedAt:        conv.TimeFromPgTimestamptz(d.CreatedAt),
		IsVerified:       d.IsVerified,
		MedicalReports:   d.MedicalReports,
		IncidentsReports: d.IncidentsReports,
		GoalsReports:     d.GoalsReports,
	}
}

func toDomainClientEmergencyContactFromRow(row db.ListEmergencyContactsRow) domain.ClientEmergencyContact {
	var relationStatus *string
	if row.RelationStatus != nil {
		s := string(*row.RelationStatus)
		relationStatus = &s
	}

	return domain.ClientEmergencyContact{
		ID:               row.ID,
		ClientID:         row.ClientID,
		FirstName:        row.FirstName,
		LastName:         row.LastName,
		Email:            row.Email,
		PhoneNumber:      row.PhoneNumber,
		Address:          row.Address,
		Relationship:     row.Relationship,
		RelationStatus:   relationStatus,
		CreatedAt:        conv.TimeFromPgTimestamptz(row.CreatedAt),
		IsVerified:       row.IsVerified,
		MedicalReports:   row.MedicalReports,
		IncidentsReports: row.IncidentsReports,
		GoalsReports:     row.GoalsReports,
	}
}

func toDomainAssignedEmployeeFromAssignRow(row db.AssignEmployeeRow) *domain.AssignedEmployee {
	return &domain.AssignedEmployee{
		ID:                 row.ID,
		ClientID:           row.ClientID,
		EmployeeID:         row.EmployeeID,
		StartDate:          conv.TimeFromPgDate(row.StartDate),
		Role:               row.Role,
		CreatedAt:          conv.TimeFromPgTimestamptz(row.CreatedAt),
		EmployeeFirstName:  "",
		EmployeeLastName:   "",
		UserID:             row.UserID,
		ClientFirstName:    row.ClientFirstName,
		ClientLastName:     row.ClientLastName,
		ClientLocationName: row.ClientLocationName,
	}
}

func toDomainAssignedEmployeeFromListRow(row db.ListAssignedEmployeesRow) domain.AssignedEmployee {
	return domain.AssignedEmployee{
		ID:                 row.ID,
		ClientID:           row.ClientID,
		EmployeeID:         row.EmployeeID,
		StartDate:          conv.TimeFromPgDate(row.StartDate),
		Role:               row.Role,
		CreatedAt:          conv.TimeFromPgTimestamptz(row.CreatedAt),
		EmployeeFirstName:  row.EmployeeFirstName,
		EmployeeLastName:   row.EmployeeLastName,
		UserID:             uuid.Nil,
		ClientFirstName:    "",
		ClientLastName:     "",
		ClientLocationName: nil,
	}
}

func toDomainAssignedEmployeeFromGetRow(row db.GetAssignedEmployeeRow) *domain.AssignedEmployee {
	return &domain.AssignedEmployee{
		ID:                 row.ID,
		ClientID:           row.ClientID,
		EmployeeID:         row.EmployeeID,
		StartDate:          conv.TimeFromPgDate(row.StartDate),
		Role:               row.Role,
		CreatedAt:          conv.TimeFromPgTimestamptz(row.CreatedAt),
		EmployeeFirstName:  row.EmployeeFirstName,
		EmployeeLastName:   row.EmployeeLastName,
		UserID:             uuid.Nil,
		ClientFirstName:    "",
		ClientLastName:     "",
		ClientLocationName: nil,
	}
}

func toDomainAssignedEmployee(d db.AssignedEmployee) *domain.AssignedEmployee {
	return &domain.AssignedEmployee{
		ID:                 d.ID,
		ClientID:           d.ClientID,
		EmployeeID:         d.EmployeeID,
		StartDate:          conv.TimeFromPgDate(d.StartDate),
		Role:               d.Role,
		CreatedAt:          conv.TimeFromPgTimestamptz(d.CreatedAt),
		EmployeeFirstName:  "",
		EmployeeLastName:   "",
		UserID:             uuid.Nil,
		ClientFirstName:    "",
		ClientLastName:     "",
		ClientLocationName: nil,
	}
}

func toDomainClientDiagnosis(d db.ClientDiagnosis) *domain.ClientDiagnosis {
	var diagnosedOn *time.Time
	if d.DiagnosedOn.Valid {
		t := d.DiagnosedOn.Time
		diagnosedOn = &t
	}

	var resolvedOn *time.Time
	if d.ResolvedOn.Valid {
		t := d.ResolvedOn.Time
		resolvedOn = &t
	}

	return &domain.ClientDiagnosis{
		ID:                  d.ID,
		ClientID:            d.ClientID,
		CodeSystem:          d.CodeSystem,
		Code:                d.Code,
		Title:               d.Title,
		Description:         d.Description,
		Status:              string(d.Status),
		Severity:            string(d.Severity),
		DiagnosedOn:         diagnosedOn,
		ResolvedOn:          resolvedOn,
		DiagnosingClinician: d.DiagnosingClinician,
		Notes:               d.Notes,
		CreatedAt:           conv.TimeFromPgTimestamptz(d.CreatedAt),
		UpdatedAt:           conv.TimeFromPgTimestamptz(d.UpdatedAt),
	}
}

func toDomainClientDiagnosisFromRow(row db.ListClientDiagnosesRow) domain.ClientDiagnosis {
	var diagnosedOn *time.Time
	if row.DiagnosedOn.Valid {
		t := row.DiagnosedOn.Time
		diagnosedOn = &t
	}

	var resolvedOn *time.Time
	if row.ResolvedOn.Valid {
		t := row.ResolvedOn.Time
		resolvedOn = &t
	}

	return domain.ClientDiagnosis{
		ID:                  row.ID,
		ClientID:            row.ClientID,
		CodeSystem:          row.CodeSystem,
		Code:                row.Code,
		Title:               row.Title,
		Description:         row.Description,
		Status:              string(row.Status),
		Severity:            string(row.Severity),
		DiagnosedOn:         diagnosedOn,
		ResolvedOn:          resolvedOn,
		DiagnosingClinician: row.DiagnosingClinician,
		Notes:               row.Notes,
		CreatedAt:           conv.TimeFromPgTimestamptz(row.CreatedAt),
		UpdatedAt:           conv.TimeFromPgTimestamptz(row.UpdatedAt),
	}
}

func toDomainClientMedicationOrder(d db.ClientMedicationOrder) *domain.ClientMedicationOrder {
	var endDate *time.Time
	if d.EndDate.Valid {
		t := d.EndDate.Time
		endDate = &t
	}

	return &domain.ClientMedicationOrder{
		ID:                           d.ID,
		ClientID:                     d.ClientID,
		DiagnosisID:                  d.DiagnosisID,
		MedicationName:               d.MedicationName,
		DosageText:                   d.DosageText,
		DoseAmount:                   d.DoseAmount,
		DoseUnit:                     d.DoseUnit,
		Route:                        d.Route,
		FrequencyText:                d.FrequencyText,
		Schedule:                     d.Schedule,
		IsPrn:                        d.IsPrn,
		PrnIndication:                d.PrnIndication,
		MaxDosesPer24h:               d.MaxDosesPer24h,
		StartDate:                    conv.TimeFromPgDate(d.StartDate),
		EndDate:                      endDate,
		Status:                       string(d.Status),
		AdminMode:                    string(d.AdminMode),
		ResponsibleEmployeeID:        d.ResponsibleEmployeeID,
		ResponsibleEmployeeFirstName: nil,
		ResponsibleEmployeeLastName:  nil,
		IsCritical:                   d.IsCritical,
		Notes:                        d.Notes,
		SourceAttachmentUUID:         d.SourceAttachmentUuid,
		DiagnosisTitle:               nil,
		DiagnosisCodeSystem:          nil,
		DiagnosisCode:                nil,
		CreatedAt:                    conv.TimeFromPgTimestamptz(d.CreatedAt),
		UpdatedAt:                    conv.TimeFromPgTimestamptz(d.UpdatedAt),
	}
}

func toDomainClientMedicationOrderFromRow(row db.ListClientMedicationOrdersRow) domain.ClientMedicationOrder {
	var endDate *time.Time
	if row.EndDate.Valid {
		t := row.EndDate.Time
		endDate = &t
	}

	return domain.ClientMedicationOrder{
		ID:                           row.ID,
		ClientID:                     row.ClientID,
		DiagnosisID:                  row.DiagnosisID,
		MedicationName:               row.MedicationName,
		DosageText:                   row.DosageText,
		DoseAmount:                   row.DoseAmount,
		DoseUnit:                     row.DoseUnit,
		Route:                        row.Route,
		FrequencyText:                row.FrequencyText,
		Schedule:                     row.Schedule,
		IsPrn:                        row.IsPrn,
		PrnIndication:                row.PrnIndication,
		MaxDosesPer24h:               row.MaxDosesPer24h,
		StartDate:                    conv.TimeFromPgDate(row.StartDate),
		EndDate:                      endDate,
		Status:                       string(row.Status),
		AdminMode:                    string(row.AdminMode),
		ResponsibleEmployeeID:        row.ResponsibleEmployeeID,
		ResponsibleEmployeeFirstName: row.ResponsibleEmployeeFirstName,
		ResponsibleEmployeeLastName:  row.ResponsibleEmployeeLastName,
		IsCritical:                   row.IsCritical,
		Notes:                        row.Notes,
		SourceAttachmentUUID:         row.SourceAttachmentUuid,
		DiagnosisTitle:               row.DiagnosisTitle,
		DiagnosisCodeSystem:          row.DiagnosisCodeSystem,
		DiagnosisCode:                row.DiagnosisCode,
		CreatedAt:                    conv.TimeFromPgTimestamptz(row.CreatedAt),
		UpdatedAt:                    conv.TimeFromPgTimestamptz(row.UpdatedAt),
	}
}

func toDomainClientMedicationOrderFromGetRow(row db.GetClientMedicationOrderRow) *domain.ClientMedicationOrder {
	var endDate *time.Time
	if row.EndDate.Valid {
		t := row.EndDate.Time
		endDate = &t
	}

	return &domain.ClientMedicationOrder{
		ID:                           row.ID,
		ClientID:                     row.ClientID,
		DiagnosisID:                  row.DiagnosisID,
		MedicationName:               row.MedicationName,
		DosageText:                   row.DosageText,
		DoseAmount:                   row.DoseAmount,
		DoseUnit:                     row.DoseUnit,
		Route:                        row.Route,
		FrequencyText:                row.FrequencyText,
		Schedule:                     row.Schedule,
		IsPrn:                        row.IsPrn,
		PrnIndication:                row.PrnIndication,
		MaxDosesPer24h:               row.MaxDosesPer24h,
		StartDate:                    conv.TimeFromPgDate(row.StartDate),
		EndDate:                      endDate,
		Status:                       string(row.Status),
		AdminMode:                    string(row.AdminMode),
		ResponsibleEmployeeID:        row.ResponsibleEmployeeID,
		ResponsibleEmployeeFirstName: row.ResponsibleEmployeeFirstName,
		ResponsibleEmployeeLastName:  row.ResponsibleEmployeeLastName,
		IsCritical:                   row.IsCritical,
		Notes:                        row.Notes,
		SourceAttachmentUUID:         row.SourceAttachmentUuid,
		DiagnosisTitle:               row.DiagnosisTitle,
		DiagnosisCodeSystem:          row.DiagnosisCodeSystem,
		DiagnosisCode:                row.DiagnosisCode,
		CreatedAt:                    conv.TimeFromPgTimestamptz(row.CreatedAt),
		UpdatedAt:                    conv.TimeFromPgTimestamptz(row.UpdatedAt),
	}
}

func toDomainProgressReport(row db.ProgressReport) *domain.ProgressReport {
	return &domain.ProgressReport{
		ID:                     row.ID,
		ClientID:               row.ClientID,
		Date:                   conv.TimeFromPgTimestamptz(row.Date),
		Title:                  row.Title,
		ReportText:             row.ReportText,
		EmployeeID:             row.EmployeeID,
		Type:                   string(row.Type),
		EmotionalState:         string(row.EmotionalState),
		CreatedAt:              conv.TimeFromPgTimestamptz(row.CreatedAt),
		EmployeeFirstName:      "",
		EmployeeLastName:       "",
		EmployeeProfilePicture: nil,
	}
}

func toDomainProgressReportFromRow(row db.ListProgressReportsRow) domain.ProgressReport {
	return domain.ProgressReport{
		ID:                     row.ID,
		ClientID:               row.ClientID,
		Date:                   conv.TimeFromPgTimestamptz(row.Date),
		Title:                  row.Title,
		ReportText:             row.ReportText,
		EmployeeID:             row.EmployeeID,
		Type:                   string(row.Type),
		EmotionalState:         string(row.EmotionalState),
		CreatedAt:              conv.TimeFromPgTimestamptz(row.CreatedAt),
		EmployeeFirstName:      row.EmployeeFirstName,
		EmployeeLastName:       row.EmployeeLastName,
		EmployeeProfilePicture: row.EmployeeProfilePicture,
	}
}

func toDomainProgressReportFromGetRow(row db.GetProgressReportRow) *domain.ProgressReport {
	return &domain.ProgressReport{
		ID:                     row.ID,
		ClientID:               row.ClientID,
		Date:                   conv.TimeFromPgTimestamptz(row.Date),
		Title:                  row.Title,
		ReportText:             row.ReportText,
		EmployeeID:             row.EmployeeID,
		Type:                   string(row.Type),
		EmotionalState:         string(row.EmotionalState),
		CreatedAt:              conv.TimeFromPgTimestamptz(row.CreatedAt),
		EmployeeFirstName:      row.EmployeeFirstName,
		EmployeeLastName:       row.EmployeeLastName,
		EmployeeProfilePicture: row.EmployeeProfilePicture,
	}
}

// =====================
// Appointment Card
// =====================

func (r *ClientRepository) GetAppointmentCard(ctx context.Context, clientID uuid.UUID) (*domain.AppointmentCard, error) {
	row, err := actorQuery(ctx, r.store, func(q *db.Queries) (db.GetAppointmentCardRow, error) {
		return q.GetAppointmentCard(ctx, clientID)
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get appointment card: %w", err)
	}
	return &domain.AppointmentCard{
		ID:                     row.ID,
		ClientID:               row.ClientID,
		GeneralInformation:     row.GeneralInformation,
		ImportantContacts:      row.ImportantContacts,
		HouseholdInfo:          row.HouseholdInfo,
		OrganizationAgreements: row.OrganizationAgreements,
		YouthOfficerAgreements: row.YouthOfficerAgreements,
		TreatmentAgreements:    row.TreatmentAgreements,
		SmokingRules:           row.SmokingRules,
		Work:                   row.Work,
		SchoolInternship:       row.SchoolInternship,
		Travel:                 row.Travel,
		Leave:                  row.Leave,
		CreatedAt:              conv.TimeFromPgTimestamptz(row.CreatedAt),
		UpdatedAt:              conv.TimeFromPgTimestamptz(row.UpdatedAt),
		ClientFirstName:        row.FirstName,
		ClientLastName:         row.LastName,
	}, nil
}

func (r *ClientRepository) CreateAppointmentCard(ctx context.Context, clientID uuid.UUID, params domain.UpdateAppointmentCardParams) (*domain.AppointmentCard, error) {
	card, err := actorQuery(ctx, r.store, func(q *db.Queries) (db.AppointmentCard, error) {
		return q.CreateAppointmentCard(ctx, db.CreateAppointmentCardParams{
			ClientID:               clientID,
			GeneralInformation:     params.GeneralInformation,
			ImportantContacts:      params.ImportantContacts,
			HouseholdInfo:          params.HouseholdInfo,
			OrganizationAgreements: params.OrganizationAgreements,
			YouthOfficerAgreements: params.YouthOfficerAgreements,
			TreatmentAgreements:    params.TreatmentAgreements,
			SmokingRules:           params.SmokingRules,
			Work:                   params.Work,
			SchoolInternship:       params.SchoolInternship,
			Travel:                 params.Travel,
			Leave:                  params.Leave,
		})
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create appointment card: %w", err)
	}
	return toDomainAppointmentCard(card), nil
}

func (r *ClientRepository) UpdateAppointmentCard(ctx context.Context, clientID uuid.UUID, params domain.UpdateAppointmentCardParams) (*domain.AppointmentCard, error) {
	card, err := actorQuery(ctx, r.store, func(q *db.Queries) (db.AppointmentCard, error) {
		return q.UpdateAppointmentCard(ctx, db.UpdateAppointmentCardParams{
			ClientID:               clientID,
			GeneralInformation:     params.GeneralInformation,
			ImportantContacts:      params.ImportantContacts,
			HouseholdInfo:          params.HouseholdInfo,
			OrganizationAgreements: params.OrganizationAgreements,
			YouthOfficerAgreements: params.YouthOfficerAgreements,
			TreatmentAgreements:    params.TreatmentAgreements,
			SmokingRules:           params.SmokingRules,
			Work:                   params.Work,
			SchoolInternship:       params.SchoolInternship,
			Travel:                 params.Travel,
			Leave:                  params.Leave,
		})
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update appointment card: %w", err)
	}
	return toDomainAppointmentCard(card), nil
}

func toDomainAppointmentCard(card db.AppointmentCard) *domain.AppointmentCard {
	return &domain.AppointmentCard{
		ID:                     card.ID,
		ClientID:               card.ClientID,
		GeneralInformation:     card.GeneralInformation,
		ImportantContacts:      card.ImportantContacts,
		HouseholdInfo:          card.HouseholdInfo,
		OrganizationAgreements: card.OrganizationAgreements,
		YouthOfficerAgreements: card.YouthOfficerAgreements,
		TreatmentAgreements:    card.TreatmentAgreements,
		SmokingRules:           card.SmokingRules,
		Work:                   card.Work,
		SchoolInternship:       card.SchoolInternship,
		Travel:                 card.Travel,
		Leave:                  card.Leave,
		CreatedAt:              conv.TimeFromPgTimestamptz(card.CreatedAt),
		UpdatedAt:              conv.TimeFromPgTimestamptz(card.UpdatedAt),
	}
}

func toDomainAiGeneratedReport(row db.AiGeneratedReport) *domain.AiGeneratedReport {
	return &domain.AiGeneratedReport{
		ID:         row.ID,
		ClientID:   row.ClientID,
		ReportText: row.ReportText,
		StartDate:  conv.TimeFromPgDate(row.StartDate),
		EndDate:    conv.TimeFromPgDate(row.EndDate),
		CreatedAt:  conv.TimeFromPgTimestamptz(row.CreatedAt),
	}
}

func toDomainAiGeneratedReportFromRow(row db.ListAiGeneratedReportsRow) domain.AiGeneratedReport {
	return domain.AiGeneratedReport{
		ID:         row.ID,
		ClientID:   row.ClientID,
		ReportText: row.ReportText,
		StartDate:  conv.TimeFromPgDate(row.StartDate),
		EndDate:    conv.TimeFromPgDate(row.EndDate),
		CreatedAt:  conv.TimeFromPgTimestamptz(row.CreatedAt),
	}
}

func diagnosisStatusFromPtr(value *string) *db.DiagnosisStatusEnum {
	if value == nil {
		return nil
	}
	status := db.DiagnosisStatusEnum(*value)
	return &status
}

func diagnosisSeverityFromPtr(value *string) *db.DiagnosisSeverityEnum {
	if value == nil {
		return nil
	}
	severity := db.DiagnosisSeverityEnum(*value)
	return &severity
}

func medicationOrderStatusFromPtr(value *string) *db.MedicationOrderStatusEnum {
	if value == nil {
		return nil
	}
	status := db.MedicationOrderStatusEnum(*value)
	return &status
}

func medicationAdminModeFromPtr(value *string) *db.MedicationAdminModeEnum {
	if value == nil {
		return nil
	}
	mode := db.MedicationAdminModeEnum(*value)
	return &mode
}
