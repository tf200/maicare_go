package clientp

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"
	"maicare_go/service/attachment"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *clientService) CreateClientDetails(req CreateClientDetailsRequest, ctx context.Context) (*CreateClientDetailsResponse, error) {
	parsedDateOfBirth, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientDetails",
			"Failed to parse date of birth", zap.Error(err))
		return nil, fmt.Errorf("failed to parse date of birth")
	}

	client, err := s.Store.CreateClientDetails(ctx, db.CreateClientDetailsParams{
		FirstName:                  req.FirstName,
		LastName:                   req.LastName,
		DateOfBirth:                pgtype.Date{Time: parsedDateOfBirth, Valid: true},
		Identity:                   true,
		Bsn:                        req.Bsn,
		BsnVerifiedBy:              req.BsnVerifiedBy,
		Email:                      req.Email,
		PhoneNumber:                req.PhoneNumber,
		CareType:                   db.NullIntakeCareTypeFromPtr(req.CareType),
		SenderID:                   req.SenderID,
		LocationID:                 req.LocationID,
		EducationCurrentlyEnrolled: req.EducationCurrentlyEnrolled,
		EducationInstitution:       req.EducationInstitution,
		EducationMentorName:        req.EducationMentorName,
		EducationMentorPhone:       req.EducationMentorPhone,
		EducationMentorEmail:       req.EducationMentorEmail,
		EducationAdditionalNotes:   req.EducationAdditionalNotes,
		WorkCurrentlyEmployed:      req.WorkCurrentlyEmployed,
		WorkCurrentEmployer:        req.WorkCurrentEmployer,
		WorkCurrentEmployerPhone:   req.WorkCurrentEmployerPhone,
		WorkCurrentEmployerEmail:   req.WorkCurrentEmployerEmail,
		WorkCurrentPosition:        req.WorkCurrentPosition,
		WorkStartDate:              pgtype.Date{Time: req.WorkStartDate, Valid: true},
		WorkAdditionalNotes:        req.WorkAdditionalNotes,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientDetails",
			"Failed to create client details", zap.Error(err))
		return nil, fmt.Errorf("failed to create client details")
	}

	result := &CreateClientDetailsResponse{
		ID:                         client.ID,
		FirstName:                  client.FirstName,
		LastName:                   client.LastName,
		DateOfBirth:                client.DateOfBirth.Time,
		Identity:                   client.Identity,
		Status:                     string(client.Status),
		Bsn:                        client.Bsn,
		Email:                      client.Email,
		PhoneNumber:                client.PhoneNumber,
		Gender:                     string(client.Gender),
		Filenumber:                 client.Filenumber,
		Created:                    client.CreatedAt.Time,
		SenderID:                   client.SenderID,
		LocationID:                 client.LocationID,
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
		WorkStartDate:              client.WorkStartDate.Time,
		WorkAdditionalNotes:        client.WorkAdditionalNotes,
	}

	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "CreateClientDetails",
		"Successfully created client details", zap.String("ClientID", client.ID.String()))
	return result, nil
}

func (s *clientService) ListClientDetails(ctx *gin.Context, req ListClientsApiParams) (*pagination.Response[ListClientsApiResponse], error) {
	params := req.GetParams()

	var clients []db.ListClientDetailsRow
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		clients, err = q.ListClientDetails(ctx, db.ListClientDetailsParams{
			Limit:      params.Limit,
			Offset:     params.Offset,
			Status:     db.NullClientStatusFromPtr(req.Status),
			LocationID: req.LocationID,
			Search:     req.Search,
		})
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListClientDetails",
			"Failed to list client details", zap.Error(err))
		return nil, fmt.Errorf("failed to list client details")
	}
	if len(clients) == 0 {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "ListClientDetails",
			"No clients found")
		pagObj := pagination.NewResponse(ctx, req.Request, []ListClientsApiResponse{}, 0)
		return &pagObj, nil
	}

	totalCount := clients[0].TotalCount

	clientList := make([]ListClientsApiResponse, len(clients))
	for i, client := range clients {
		clientList[i] = ListClientsApiResponse{
			ID:           client.ID,
			FirstName:    client.FirstName,
			LastName:     client.LastName,
			Bsn:          client.Bsn,
			Filenumber:   client.Filenumber,
			LocationName: client.LocationName,
			CareType:     db.IntakeCareTypePtrFromEnum(client.CareType),
			Status:       string(client.Status),
			GoalsCount:   client.GoalsCount,
			RiskCount:    client.RiskCount,
			CreatedAt:    client.CreatedAt.Time,
		}
	}
	pagObj := pagination.NewResponse(ctx, req.Request, clientList, totalCount)
	return &pagObj, nil
}

func (s *clientService) ListWaitingListClients(ctx *gin.Context, req ListWaitingListClientsParams) (*pagination.Response[ListWaitingListClientsResponse], error) {
	params := req.GetParams()
	sortDays := "desc"
	if req.SortDays != nil {
		sortDays = *req.SortDays
	}

	var clients []db.ListWaitingListClientsRow
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		clients, err = q.ListWaitingListClients(ctx, db.ListWaitingListClientsParams{
			Search:    req.Search,
			Placement: db.NullIntakeCareTypeFromPtr(req.Placement),
			SortDays:  sortDays,
			Offset:    params.Offset,
			Limit:     params.Limit,
		})
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListWaitingListClients",
			"Failed to list waiting list clients", zap.Error(err))
		return nil, fmt.Errorf("failed to list waiting list clients")
	}

	if len(clients) == 0 {
		pagObj := pagination.NewResponse(ctx, req.Request, []ListWaitingListClientsResponse{}, 0)
		return &pagObj, nil
	}

	totalCount := clients[0].TotalCount
	response := make([]ListWaitingListClientsResponse, len(clients))
	for i, client := range clients {
		var admissionType *string
		if client.AdmissionType.Valid {
			value := string(client.AdmissionType.AdmissionTypeEnum)
			admissionType = &value
		}

		response[i] = ListWaitingListClientsResponse{
			ID:             client.ID,
			FirstName:      client.FirstName,
			LastName:       client.LastName,
			CareType:       db.IntakeCareTypePtrFromEnum(client.CareType),
			SenderName:     client.SenderName,
			DaysInWaitlist: client.DaysInWaitlist,
			AdmissionType:  admissionType,
			Bsn:            client.Bsn,
		}
	}

	pagObj := pagination.NewResponse(ctx, req.Request, response, totalCount)
	return &pagObj, nil
}

func (s *clientService) ListInCareClients(ctx *gin.Context, req ListInCareClientsParams) (*pagination.Response[ListInCareClientsResponse], error) {
	params := req.GetParams()
	sortDaysInCare := "desc"
	if req.SortDaysInCare != nil {
		sortDaysInCare = *req.SortDaysInCare
	}

	var statusFilters []db.ClientStatusEnum
	if len(req.Status) > 0 {
		statusFilters = make([]db.ClientStatusEnum, len(req.Status))
		for i, status := range req.Status {
			statusFilters[i] = db.ClientStatusEnum(status)
		}
	}

	var clients []db.ListInCareClientsRow
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		clients, err = q.ListInCareClients(ctx, db.ListInCareClientsParams{
			Search:         req.Search,
			Status:         statusFilters,
			SortDaysInCare: sortDaysInCare,
			Limit:          params.Limit,
			Offset:         params.Offset,
		})
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListInCareClients",
			"Failed to list in-care clients", zap.Error(err))
		return nil, fmt.Errorf("failed to list in-care clients")
	}

	if len(clients) == 0 {
		pagObj := pagination.NewResponse(ctx, req.Request, []ListInCareClientsResponse{}, 0)
		return &pagObj, nil
	}

	totalCount := clients[0].TotalCount
	response := make([]ListInCareClientsResponse, len(clients))
	for i, client := range clients {
		var careStartDate *time.Time
		if client.CareStartDate.Valid {
			t := client.CareStartDate.Time
			careStartDate = &t
		}

		var coordinatorName *string
		if client.CoordinatorName != "" {
			name := client.CoordinatorName
			coordinatorName = &name
		}

		response[i] = ListInCareClientsResponse{
			ID:                client.ID,
			Bsn:               client.Bsn,
			FirstName:         client.FirstName,
			LastName:          client.LastName,
			CoordinatorName:   coordinatorName,
			LocationName:      client.LocationName,
			Status:            string(client.Status),
			CareStartDate:     careStartDate,
			DaysInCare:        client.DaysInCare,
			HasActiveContract: client.HasActiveContract,
		}
	}

	pagObj := pagination.NewResponse(ctx, req.Request, response, totalCount)
	return &pagObj, nil
}

func (s *clientService) GetClientsCount(ctx context.Context) (*GetClientsCountResponse, error) {
	var count db.GetClientCountsRow
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		count, err = q.GetClientCounts(ctx)
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetClientsCount",
			"Failed to get clients count", zap.Error(err))
		return nil, fmt.Errorf("failed to get clients count")
	}
	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "GetClientsCount",
		"Successfully retrieved clients count")
	return &GetClientsCountResponse{
		TotalClients:         count.TotalClients,
		ClientsInCare:        count.ClientsInCare,
		ClientsOnWaitingList: count.ClientsOnWaitingList,
		ClientsOutOfCare:     count.ClientsOutOfCare,
	}, nil
}

func (s *clientService) GetClientStatusCounts(ctx context.Context) (*GetClientStatusCountsResponse, error) {
	var count db.GetClientStatusCountsRow
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		count, err = q.GetClientStatusCounts(ctx)
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetClientStatusCounts",
			"Failed to get client status counts", zap.Error(err))
		return nil, fmt.Errorf("failed to get client status counts")
	}

	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "GetClientStatusCounts",
		"Successfully retrieved client status counts")
	return &GetClientStatusCountsResponse{
		ClientsInOrScheduledInCare:     count.ClientsInOrScheduledInCare,
		ClientsOnWaitingList:           count.ClientsOnWaitingList,
		ClientsOutOrScheduledOutOfCare: count.ClientsOutOrScheduledOutOfCare,
	}, nil
}

func (s *clientService) GetClientDetails(ctx context.Context, clientID uuid.UUID) (*GetClientApiResponse, error) {
	var client db.GetClientDetailsRow
	var goals []db.ListActiveGoalSummariesByClientIDRow
	var emergencyContacts []db.ListTopEmergencyContactsByClientIDRow
	var existingDocumentLabels []string
	var missingDocumentLabels []string
	var counts db.GetClientPageCountsRow
	var coordinatorRows []db.GetClientCoordinatorRow
	var latestStatusHistory db.GetClientLatestStatusHistoryRow
	var activeContracts []db.ListClientActiveApprovedContractsRow
	var latestDraft db.GetLatestDraftEvaluationByClientRow
	var hasLatestDraft bool
	var latestCompleted db.GetLatestCompletedEvaluationByClientRow
	var hasLatestCompleted bool
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		client, err = q.GetClientDetails(ctx, clientID)
		if err != nil {
			return err
		}

		goals, err = q.ListActiveGoalSummariesByClientID(ctx, clientID)
		if err != nil {
			return err
		}

		emergencyContacts, err = q.ListTopEmergencyContactsByClientID(ctx, clientID)
		if err != nil {
			return err
		}

		existingDocumentLabels, err = q.ListExistingClientDocumentLabels(ctx, clientID)
		if err != nil {
			return err
		}

		missingDocumentLabels, err = q.GetMissingClientDocuments(ctx, clientID)
		if err != nil {
			return err
		}

		counts, err = q.GetClientPageCounts(ctx, clientID)
		if err != nil {
			return err
		}

		coordinatorRows, err = q.GetClientCoordinator(ctx, clientID)
		if err != nil {
			return err
		}

		latestStatusHistory, err = q.GetClientLatestStatusHistory(ctx, clientID)
		if err != nil {
			return err
		}

		activeContracts, err = q.ListClientActiveApprovedContracts(ctx, clientID)
		if err != nil {
			return err
		}

		latestDraft, err = q.GetLatestDraftEvaluationByClient(ctx, clientID)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return err
			}
		} else {
			hasLatestDraft = true
		}

		latestCompleted, err = q.GetLatestCompletedEvaluationByClient(ctx, clientID)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return err
			}
		} else {
			hasLatestCompleted = true
		}

		return nil
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetClientDetails",
			"Failed to get client details", zap.Error(err), zap.String("ClientID", clientID.String()))
		return nil, fmt.Errorf("failed to get client details")
	}

	var sender *ClientSenderMinimalResponse
	if client.SenderName != nil || client.SenderEmailAddress != nil || client.SenderPhoneNumber != nil {
		senderName := ""
		if client.SenderName != nil {
			senderName = *client.SenderName
		}

		sender = &ClientSenderMinimalResponse{
			Name:         senderName,
			EmailAddress: client.SenderEmailAddress,
			PhoneNumber:  client.SenderPhoneNumber,
		}
	}

	goalsResponse := make([]ClientGoalSummaryResponse, len(goals))
	for i, goal := range goals {
		var topicName *string
		if goal.TopicName != "" {
			topicName = &goal.TopicName
		}

		goalsResponse[i] = ClientGoalSummaryResponse{
			Title:     goal.Title,
			Priority:  string(goal.Priority),
			TopicName: topicName,
		}
	}

	emergencyContactsResponse := make([]ClientEmergencySummary, len(emergencyContacts))
	for i, contact := range emergencyContacts {
		emergencyContactsResponse[i] = ClientEmergencySummary{
			ID:           contact.ID,
			FirstName:    contact.FirstName,
			LastName:     contact.LastName,
			Relationship: contact.Relationship,
			PhoneNumber:  contact.PhoneNumber,
			Email:        contact.Email,
		}
	}

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

	age := calculateAge(client.DateOfBirth)

	alerts := make([]ClientPageAlert, 0, 4)
	if len(missingDocumentLabels) > 0 {
		alerts = append(alerts, ClientPageAlert{
			Code:     "missing_documents",
			Severity: "warning",
			Message:  fmt.Sprintf("%d required documents are missing", len(missingDocumentLabels)),
		})
	}
	if len(goalsResponse) == 0 {
		alerts = append(alerts, ClientPageAlert{
			Code:     "missing_goals",
			Severity: "info",
			Message:  "No active goals defined",
		})
	}
	if counts.IncidentsCount > 0 {
		alerts = append(alerts, ClientPageAlert{
			Code:     "has_incidents",
			Severity: "warning",
			Message:  fmt.Sprintf("%d incident(s) linked to client", counts.IncidentsCount),
		})
	}

	var coordinator *ClientCoordinatorResponse
	if len(coordinatorRows) > 0 {
		coordinator = &ClientCoordinatorResponse{
			EmployeeID: &coordinatorRows[0].EmployeeID,
			FirstName:  &coordinatorRows[0].FirstName,
			LastName:   &coordinatorRows[0].LastName,
			StartDate:  util.DatePtr(coordinatorRows[0].StartDate),
		}
	}

	var careSchedule *ClientCareScheduleResponse
	if client.Status == db.ClientStatusEnumScheduledInCare {
		var careStartDate *time.Time
		var placedInCareAt *time.Time
		var daysUntilStart *int32

		if client.CareStartDate.Valid {
			careStartDate = util.DatePtr(client.CareStartDate)

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
			placedInCareAt = &placedAt
		}

		shouldBeActiveNow := false
		if client.CareStartDate.Valid {
			shouldBeActiveNow = !client.CareStartDate.Time.After(time.Now().UTC())
		}

		if shouldBeActiveNow {
			alerts = append(alerts, ClientPageAlert{
				Code:     "start_date_reached_not_activated",
				Severity: "warning",
				Message:  "Care start date has been reached but client is still scheduled_in_care",
			})
		}

		if coordinator == nil {
			alerts = append(alerts, ClientPageAlert{
				Code:     "missing_coordinator",
				Severity: "warning",
				Message:  "No coordinator assigned",
			})
		}

		careSchedule = &ClientCareScheduleResponse{
			CareStartDate:      careStartDate,
			PlacedInCareAt:     placedInCareAt,
			DaysUntilStart:     daysUntilStart,
			ShouldBeActiveNow:  shouldBeActiveNow,
			NextEvaluationDate: util.DatePtr(client.NextEvaluationDate),
		}
	}

	var care *ClientInCareResponse
	var dischargeSchedule *ClientDischargeScheduleResponse
	var dischargeSummary *ClientDischargeSummaryResponse
	var contractSummary *ClientContractSummaryResponse
	var evaluationSummary *ClientEvaluationSummaryResponse
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

		care = &ClientInCareResponse{
			CareStartDate:            util.DatePtr(client.CareStartDate),
			PlacedInCareAt:           timestamptzPtr(client.PlacedInCareAt),
			DaysInCare:               daysInCare,
			EvaluationIntervalsWeeks: client.EvaluationIntervalsWeeks,
			LastEvaluationAnchorDate: util.DatePtr(client.LastEvaluationAnchorDate),
			NextEvaluationDate:       util.DatePtr(client.NextEvaluationDate),
		}

		var activeContract *ClientActiveContractResponse
		var daysUntilContractEnd *int32
		if len(activeContracts) > 0 {
			active := activeContracts[0]
			activeContract = &ClientActiveContractResponse{
				ID:              active.ID,
				Status:          stringPtrOrNil(active.Status),
				StartDate:       timestamptzPtr(active.StartDate),
				EndDate:         timestamptzPtr(active.EndDate),
				FinancingAct:    stringPtrOrNil(active.FinancingAct),
				FinancingOption: stringPtrOrNil(active.FinancingOption),
				CareType:        stringPtrOrNil(active.CareType),
			}
			daysUntilContractEnd = &active.DaysUntilContractEnd
		}

		hasActiveApprovedContract := len(activeContracts) > 0
		contractSummary = &ClientContractSummaryResponse{
			HasActiveApprovedContract: hasActiveApprovedContract,
			ActiveContract:            activeContract,
			DaysUntilContractEnd:      daysUntilContractEnd,
		}

		var draft *ClientEvaluationDraftSummaryResponse
		if hasLatestDraft {
			draft = &ClientEvaluationDraftSummaryResponse{
				ID:        latestDraft.ID,
				UpdatedAt: timestamptzPtr(latestDraft.UpdatedAt),
			}
		}

		var lastCompleted *ClientEvaluationLastCompletedResponse
		if hasLatestCompleted {
			creatorName := strings.TrimSpace(strings.Join([]string{derefString(latestCompleted.CreatorFirstName), derefString(latestCompleted.CreatorLastName)}, " "))
			lastCompleted = &ClientEvaluationLastCompletedResponse{
				ID:                  latestCompleted.ID,
				SubmittedAt:         timestamptzPtr(latestCompleted.SubmittedAt),
				CreatedByEmployeeID: latestCompleted.CreatedByEmployeeID,
				CreatorName:         stringPtrOrNil(creatorName),
			}
		}

		var daysLeft *int32
		var priority *string
		nextEvaluationDate := util.DatePtr(client.NextEvaluationDate)
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

		evaluationSummary = &ClientEvaluationSummaryResponse{
			NextEvaluationDate: nextEvaluationDate,
			DaysLeft:           daysLeft,
			Priority:           priority,
			Draft:              draft,
			LastCompleted:      lastCompleted,
		}

		if coordinator == nil {
			alerts = append(alerts, ClientPageAlert{
				Code:     "missing_coordinator",
				Severity: "warning",
				Message:  "No coordinator assigned",
			})
		}

		if !hasActiveApprovedContract {
			alerts = append(alerts, ClientPageAlert{
				Code:     "missing_active_contract",
				Severity: "warning",
				Message:  "No active approved contract",
			})
		} else if daysUntilContractEnd != nil && *daysUntilContractEnd <= 30 {
			alerts = append(alerts, ClientPageAlert{
				Code:     "contract_ends_soon",
				Severity: "warning",
				Message:  fmt.Sprintf("Active contract ends in %d day(s)", *daysUntilContractEnd),
			})
		}

		if daysLeft != nil {
			if *daysLeft < 0 {
				alerts = append(alerts, ClientPageAlert{
					Code:     "evaluation_overdue",
					Severity: "warning",
					Message:  fmt.Sprintf("Evaluation is overdue by %d day(s)", -*daysLeft),
				})
			} else if *daysLeft <= 3 {
				alerts = append(alerts, ClientPageAlert{
					Code:     "evaluation_due_soon",
					Severity: "warning",
					Message:  fmt.Sprintf("Evaluation due in %d day(s)", *daysLeft),
				})
			}
		}
	}

	if client.Status == db.ClientStatusEnumScheduledOutOfCare {
		dischargeDate := util.DatePtr(client.DischargeDate)
		var dischargeReason *string
		if client.DischargeReason.Valid {
			reason := string(client.DischargeReason.DischargeReasonEnum)
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
		dischargeSchedule = &ClientDischargeScheduleResponse{
			DischargeDate:          dischargeDate,
			DischargeReason:        dischargeReason,
			FinalEvaluation:        client.FinalEvaluation,
			DaysUntilDischarge:     daysUntilDischarge,
			IsDue:                  isDue,
			MissingFinalEvaluation: missingFinalEvaluation,
		}

		if isDue && missingFinalEvaluation {
			alerts = append(alerts, ClientPageAlert{
				Code:     "discharge_due_missing_final_evaluation",
				Severity: "warning",
				Message:  "Discharge is due but final evaluation is missing",
			})
		}
	}

	if client.Status == db.ClientStatusEnumOutOfCare {
		var dischargeReason *string
		if client.DischargeReason.Valid {
			reason := string(client.DischargeReason.DischargeReasonEnum)
			dischargeReason = &reason
		}

		dischargeSummary = &ClientDischargeSummaryResponse{
			DischargeDate:   util.DatePtr(client.DischargeDate),
			DischargeReason: dischargeReason,
			FinalEvaluation: client.FinalEvaluation,
		}
	}

	var statusTimeline *ClientStatusTimelineResponse
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

		statusTimeline = &ClientStatusTimelineResponse{
			LastChangeReason: latestStatusHistory.LastChangeReason,
			LastChangedAt:    lastChangedAt,
			LastStatus:       lastStatus,
		}
	}

	clientAddress := ClientAddressResponse{
		Street:              client.Street,
		HouseNumber:         client.HouseNumber,
		HouseNumberAddition: client.HouseNumberAddition,
		PostalCode:          client.PostalCode,
		City:                client.City,
	}

	var location *ClientLocationResponse
	if client.LocationID != nil && client.LocationName != nil {
		location = &ClientLocationResponse{
			ID:   *client.LocationID,
			Name: *client.LocationName,
		}
	}

	dateOfBirth := util.DatePtr(client.DateOfBirth)

	waitlistSince := time.Now().UTC()
	if client.CreatedAt.Valid {
		waitlistSince = client.CreatedAt.Time
	}

	return &GetClientApiResponse{
		SchemaVersion: 1,
		Status:        string(client.Status),
		Client: ClientPageClientResponse{
			ID:          client.ID,
			FirstName:   client.FirstName,
			LastName:    client.LastName,
			Bsn:         client.Bsn,
			FileNumber:  client.Filenumber,
			Gender:      string(client.Gender),
			DateOfBirth: dateOfBirth,
			Age:         age,
			CareType:    db.IntakeCareTypePtrFromEnum(client.CareType),
			Address:     clientAddress,
			Location:    location,
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
		Documents: ClientDocumentsSummary{
			Existing: existingDocumentLabels,
			Missing:  missingDocumentLabels,
		},
		Goals: goalsResponse,
		Intake: ClientIntakeResponse{
			SelfSufficiencyScore: client.IntakeSelfSufficiency,
			Conclusion:           db.IntakeConclusionPtrFromEnum(client.IntakeConclusion),
			ConclusionNotes:      client.IntakeConclusionNotes,
		},
		Risks: ClientRiskSummary{
			Flags: riskFlags,
			Notes: client.RiskAdditionalNotes,
		},
		Counts: ClientPageCounts{
			Contracts:    counts.ContractsCount,
			Incidents:    counts.IncidentsCount,
			Reports:      counts.ReportsCount,
			Evaluations:  counts.EvaluationsCount,
			Documents:    counts.DocumentsCount,
			Appointments: counts.AppointmentsCount,
		},
		Alerts: alerts,
		Meta: ClientPageMetaResponse{
			WaitlistSince: waitlistSince,
			LastUpdatedAt: time.Now().UTC(),
		},
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

func timestamptzPtr(ts pgtype.Timestamptz) *time.Time {
	if !ts.Valid {
		return nil
	}
	t := ts.Time
	return &t
}

func stringPtrOrNil(v string) *string {
	trimmed := strings.TrimSpace(v)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func derefString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func (s *clientService) GetClientAddresses(ctx context.Context, clientID uuid.UUID) (*GetClientAddressesApiResponse, error) {

	return &GetClientAddressesApiResponse{}, nil
}

func (s *clientService) UpdateClientDetails(ctx context.Context, req UpdateClientDetailsRequest, clientID uuid.UUID) (*UpdateClientDetailsResponse, error) {

	result := &UpdateClientDetailsResponse{}
	return result, nil
}

func (s *clientService) UpdateClientStatus(ctx context.Context, req UpdateClientStatusRequest, clientID uuid.UUID) (*UpdateClientStatusResponse, error) {
	switch req.IsSchedueled {
	case true:
		return s.handleSchedueledStatusUpdates(ctx, req, clientID)
	case false:
		return s.handleNormalStatusUpdates(ctx, req, clientID)
	default:
		return nil, fmt.Errorf("invalid is_schedueled value")
	}
}

func (s *clientService) handleSchedueledStatusUpdates(ctx context.Context, req UpdateClientStatusRequest, clientID uuid.UUID) (*UpdateClientStatusResponse, error) {
	// if req.SchedueledFor.Before(time.Now()) {
	// 	return nil, fmt.Errorf("scheduled time must be in the future")
	// }

	// var schedueledChange db.ScheduledStatusChange
	// err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
	// 	var err error
	// 	schedueledChange, err = q.CreateSchedueledClientStatusChange(ctx, db.CreateSchedueledClientStatusChangeParams{
	// 		ClientID:      clientID,
	// 		NewStatus:     db.NullClientStatusFromPtr(&req.Status),
	// 		Reason:        &req.Reason,
	// 		ScheduledDate: pgtype.Date{Time: req.SchedueledFor, Valid: true},
	// 	})
	// 	return err
	// })
	// if err != nil {
	// 	s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateClientStatus",
	// 		"Failed to create scheduled status change", zap.Error(err), zap.String("ClientID", clientID.String()))
	// 	return nil, fmt.Errorf("failed to create scheduled status change")
	// }

	// s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "UpdateClientStatus",
	// 	"Successfully created scheduled status change", zap.String("ClientID", clientID.String()),
	// 	zap.String("NewStatus", req.Status), zap.Time("ScheduledFor", req.SchedueledFor))

	// if !schedueledChange.NewStatus.Valid {
	// 	return nil, fmt.Errorf("scheduled status change not created properly")
	// }

	return &UpdateClientStatusResponse{
		ID: clientID,
		// Status: string(schedueledChange.NewStatus.ClientStatusEnum),
	}, nil
}

func (s *clientService) ListStatusHistory(ctx context.Context, clientID uuid.UUID) ([]ListStatusHistoryApiResponse, error) {
	arg := db.ListClientStatusHistoryParams{
		ClientID: clientID,
		Limit:    10,
		Offset:   0,
	}
	var histories []db.ClientStatusHistory
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		histories, err = q.ListClientStatusHistory(ctx, arg)
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListStatusHistory",
			"Failed to list status history", zap.Error(err), zap.String("ClientID", clientID.String()))
		return nil, fmt.Errorf("failed to list status history")
	}
	if len(histories) == 0 {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "ListStatusHistory",
			"No status history found", zap.String("ClientID", clientID.String()))
		return []ListStatusHistoryApiResponse{}, nil
	}
	var historyList []ListStatusHistoryApiResponse
	for _, history := range histories {
		historyList = append(historyList, ListStatusHistoryApiResponse{
			ID:        history.ID,
			ClientID:  history.ClientID,
			OldStatus: history.OldStatus,
			NewStatus: history.NewStatus,
			Reason:    history.Reason,
			ChangedAt: history.ChangedAt.Time,
			ChangedBy: history.ChangedBy,
		})
	}
	return historyList, nil
}

func (s *clientService) handleNormalStatusUpdates(ctx context.Context, req UpdateClientStatusRequest, clientID uuid.UUID) (*UpdateClientStatusResponse, error) {
	var result *UpdateClientStatusResponse
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		// 1. Fetch old details (automatically subject to RLS visibility)
		oldClient, err := q.GetClientDetails(ctx, clientID)
		if err != nil {
			return fmt.Errorf("failed to get client details: %w", err)
		}
		// 2. Perform the update (protected by RLS UPDATE policy)
		client, err := q.UpdateClientStatus(ctx, db.UpdateClientStatusParams{
			ID:     clientID,
			Status: db.ClientStatusEnum(req.Status),
		})
		if err != nil {
			return fmt.Errorf("failed to update client status: %w", err)
		}
		// 3. Create history record (protected by RLS INSERT policy)
		_, err = q.CreateClientStatusHistory(ctx, db.CreateClientStatusHistoryParams{
			ClientID:  clientID,
			OldStatus: util.StringPtr(string(oldClient.Status)),
			NewStatus: req.Status,
			Reason:    &req.Reason,
		})
		if err != nil {
			return fmt.Errorf("failed to create client status history: %w", err)
		}
		// Prepare response object
		result = &UpdateClientStatusResponse{
			ID:     client.ID,
			Status: string(client.Status),
		}
		return nil
	})
	if err != nil {
		// Log the error returned by the transaction
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateClientStatus",
			"Transaction failed", zap.Error(err), zap.String("ClientID", clientID.String()))
		return nil, fmt.Errorf("failed to update client status: %v", err)
	}
	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "UpdateClientStatus",
		"Successfully updated client status", zap.String("ClientID", clientID.String()),
		zap.String("NewStatus", req.Status))
	return result, nil
}

func (s *clientService) AddClientDocument(ctx context.Context, req AddClientDocumentApiRequest, clientID uuid.UUID) (*AddClientDocumentApiResponse, error) {
	if len(req.Documents) == 0 {
		return nil, fmt.Errorf("at least one document is required")
	}

	txDocs := make([]db.AddClientDocumentTxParams, 0, len(req.Documents))
	seen := make(map[uuid.UUID]struct{}, len(req.Documents))
	attachmentIDs := make([]uuid.UUID, 0, len(req.Documents))
	for _, document := range req.Documents {
		if _, ok := seen[document.AttachmentID]; ok {
			return nil, fmt.Errorf("duplicate attachment_id %s in request", document.AttachmentID)
		}
		seen[document.AttachmentID] = struct{}{}
		attachmentIDs = append(attachmentIDs, document.AttachmentID)
	}

	attachments, err := s.Store.GetAttachmentsByUUIDs(ctx, attachmentIDs)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "AddClientDocument",
			"Failed to fetch attachments for client documents", zap.Error(err), zap.String("ClientID", clientID.String()))
		return nil, fmt.Errorf("failed to add client document")
	}

	if len(attachments) != len(req.Documents) {
		return nil, fmt.Errorf("one or more attachments were not found")
	}

	attachmentsByID := make(map[uuid.UUID]db.AttachmentFile, len(attachments))
	objectKeys := make([]string, 0, len(attachments))
	for _, attachmentRecord := range attachments {
		attachmentsByID[attachmentRecord.Uuid] = attachmentRecord
		objectKeys = append(objectKeys, attachmentRecord.File)
	}

	fileSizes, err := s.B2Client.GetFileInfos(ctx, objectKeys)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "AddClientDocument",
			"Uploaded file verification failed for client documents", zap.Error(err), zap.String("ClientID", clientID.String()))
		return nil, fmt.Errorf("failed to verify uploaded file")
	}

	for _, document := range req.Documents {
		attachmentRecord, ok := attachmentsByID[document.AttachmentID]
		if !ok {
			return nil, fmt.Errorf("attachment %s was not found", document.AttachmentID)
		}

		size, ok := fileSizes[attachmentRecord.File]
		if !ok {
			return nil, fmt.Errorf("uploaded file is missing")
		}
		if size <= 0 {
			return nil, fmt.Errorf("uploaded file is empty")
		}
		if size > attachment.MaxFileSize {
			return nil, fmt.Errorf("file size exceeds maximum limit of %dMB", attachment.MaxFileSize>>20)
		}

		txDocs = append(txDocs, db.AddClientDocumentTxParams{
			ClientID:     clientID,
			AttachmentID: document.AttachmentID,
			Label:        document.Label,
		})
	}

	clientDocs, err := s.Store.AddClientDocumentsTx(ctx, db.AddClientDocumentsTxParams{
		ClientID:  clientID,
		Documents: txDocs,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "AddClientDocument",
			"Failed to add client document", zap.Error(err), zap.String("ClientID", clientID.String()))
		return nil, fmt.Errorf("failed to add client document")
	}

	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "AddClientDocument",
		"Successfully added client documents", zap.String("ClientID", clientID.String()),
		zap.Int("Count", len(clientDocs.Documents)))
	responseDocs := make([]AddClientDocumentResult, 0, len(clientDocs.Documents))
	for _, clientDoc := range clientDocs.Documents {
		responseDocs = append(responseDocs, AddClientDocumentResult{
			ID:           clientDoc.ClientDocument.ID,
			AttachmentID: clientDoc.ClientDocument.AttachmentUuid,
			ClientID:     clientDoc.ClientDocument.ClientID,
			Label:        string(clientDoc.ClientDocument.Label),
			Name:         clientDoc.Attachment.Name,
			File:         clientDoc.Attachment.File,
			Size:         clientDoc.Attachment.Size,
			IsUsed:       clientDoc.Attachment.IsUsed,
			Tag:          clientDoc.Attachment.Tag,
			UpdatedAt:    clientDoc.Attachment.Updated.Time,
			CreatedAt:    clientDoc.Attachment.Created.Time,
		})
	}
	return &AddClientDocumentApiResponse{
		Documents: responseDocs,
	}, nil
}

func (s *clientService) ListClientDocuments(ctx *gin.Context, req ListClientDocumentsApiRequest, clientID uuid.UUID) (*pagination.Response[ListClientDocumentsApiResponse], error) {
	params := req.GetParams()
	var clientDocs []db.ListClientDocumentsRow
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		clientDocs, err = q.ListClientDocuments(ctx, db.ListClientDocumentsParams{
			ClientID: clientID,
			Offset:   params.Offset,
			Limit:    params.Limit,
		})
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListClientDocuments",
			"Failed to list client documents", zap.Error(err), zap.String("ClientID", clientID.String()))
		return nil, fmt.Errorf("failed to list client documents")
	}
	if len(clientDocs) == 0 {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "ListClientDocuments",
			"No client documents found", zap.String("ClientID", clientID.String()))
		pag := pagination.NewResponse(ctx, req.Request, []ListClientDocumentsApiResponse{}, 0)
		return &pag, nil
	}

	totalCount := clientDocs[0].TotalCount

	var docList []ListClientDocumentsApiResponse
	for _, doc := range clientDocs {
		docList = append(docList, ListClientDocumentsApiResponse{
			ID:             doc.ID,
			AttachmentUuid: doc.AttachmentUuid,
			ClientID:       doc.ClientID,
			Label:          string(doc.Label),
			Name:           doc.Name,
			File:           s.GenerateResponsePresignedURL(&doc.File, ctx),
			Size:           doc.Size,
			IsUsed:         doc.IsUsed,
			Tag:            doc.Tag,
			UpdatedAt:      doc.Updated.Time,
			CreatedAt:      doc.Created.Time,
		})
	}
	pag := pagination.NewResponse(ctx, req.Request, docList, totalCount)
	return &pag, nil
}

func (s *clientService) DeleteClientDocument(ctx context.Context, clientID uuid.UUID, documentID uuid.UUID) (*DeleteClientDocumentApiResponse, error) {
	clientDoc, err := s.Store.DeleteClientDocumentTx(ctx, db.DeleteClientDocumentTxParams{
		ClientID:   clientID,
		DocumentID: documentID,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteClientDocument",
			"Failed to delete client document", zap.Error(err), zap.String("ClientID", clientID.String()))
		return nil, fmt.Errorf("failed to delete client document")
	}
	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "DeleteClientDocument",
		"Successfully deleted client document", zap.String("ClientID", clientID.String()),
		zap.String("DocumentID", documentID.String()))
	return &DeleteClientDocumentApiResponse{
		ID:           clientDoc.ClientDocument.ID,
		AttachmentID: clientDoc.ClientDocument.AttachmentUuid,
	}, nil
}

func (s *clientService) GetMissingClientDocuments(ctx context.Context, clientID uuid.UUID) (*GetMissingClientDocumentsApiResponse, error) {
	var missingDocs []string
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		missingDocs, err = q.GetMissingClientDocuments(ctx, clientID)
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetMissingClientDocuments",
			"Failed to get missing client documents", zap.Error(err), zap.String("ClientID", clientID.String()))
		return nil, fmt.Errorf("failed to get missing client documents")
	}
	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "GetMissingClientDocuments",
		"Successfully retrieved missing client documents", zap.String("ClientID", clientID.String()))
	return &GetMissingClientDocumentsApiResponse{
		MissingDocs: missingDocs,
	}, nil
}
