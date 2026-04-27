package repository

import (
	"context"

	"maicare_go/db/sqlc"
	"maicare_go/internal/domain"
	"maicare_go/pkg/conv"
	"maicare_go/util"

	"github.com/google/uuid"
)

type incidentRepository struct {
	queries db.Querier
	store   *db.Store
}

func NewIncidentRepository(queries db.Querier, store *db.Store) domain.IncidentRepository {
	return &incidentRepository{queries: queries, store: store}
}

func (r *incidentRepository) CreateIncident(ctx context.Context, params domain.CreateIncidentParams) (*domain.Incident, error) {
	var incident db.CreateIncidentRow

	err := r.store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		incident, err = q.CreateIncident(ctx, db.CreateIncidentParams{
			EmployeeID:              params.EmployeeID,
			LocationID:              params.LocationID,
			ReporterInvolvement:     db.IncidentReporterInvolvementEnum(params.ReporterInvolvement),
			InformedParties:         toInformedParties(params.InformedParties),
			OccurredAt:              conv.PgTimestamptzFromTime(params.OccurredAt),
			IncidentType:            db.IncidentTypeEnum(params.IncidentType),
			SeverityOfIncident:      db.SeverityOfIncidentEnum(params.SeverityOfIncident),
			IncidentExplanation:     params.IncidentExplanation,
			RecurrenceRisk:          db.RecurrenceRiskEnum(params.RecurrenceRisk),
			IncidentPreventSteps:    params.IncidentPreventSteps,
			IncidentTakenMeasures:   params.IncidentTakenMeasures,
			CauseCategories:         toCauseCategories(params.CauseCategories),
			CauseExplanation:        params.CauseExplanation,
			PhysicalInjury:          db.PhysicalInjuryEnum(params.PhysicalInjury),
			PhysicalInjuryDesc:      params.PhysicalInjuryDesc,
			PsychologicalDamage:     db.PsychologicalDamageEnum(params.PsychologicalDamage),
			PsychologicalDamageDesc: params.PsychologicalDamageDesc,
			NeededConsultation:      db.NeededConsultationEnum(params.NeededConsultation),
			FollowUpActions:         toFollowUpActions(params.FollowUpActions),
			FollowUpNotes:           params.FollowUpNotes,
			IsEmployeeAbsent:        params.IsEmployeeAbsent,
			AdditionalDetails:       params.AdditionalDetails,
			ClientID:                params.ClientID,
			Emails:                  params.Emails,
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	return toDomainIncident(&incident), nil
}

func (r *incidentRepository) ListIncidents(ctx context.Context, params domain.ListIncidentsParams) (*domain.ListIncidentsResult, error) {
	rows, err := r.queries.ListIncidents(ctx, db.ListIncidentsParams{
		ClientID: params.ClientID,
		Limit:    params.Limit,
		Offset:   params.Offset,
	})
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return &domain.ListIncidentsResult{
			Items:      []domain.IncidentListItem{},
			TotalCount: 0,
		}, nil
	}

	items := make([]domain.IncidentListItem, len(rows))
	for i, row := range rows {
		items[i] = domain.IncidentListItem{
			ID:                     row.ID,
			OccurredAt:             conv.TimeFromPgTimestamptz(row.OccurredAt),
			IncidentType:           string(row.IncidentType),
			SeverityOfIncident:     string(row.SeverityOfIncident),
			IsConfirmed:            row.IsConfirmed,
			EmployeeFirstName:      row.EmployeeFirstName,
			EmployeeLastName:       row.EmployeeLastName,
			EmployeeProfilePicture: row.EmployeeProfilePicture,
			LocationName:           row.LocationName,
		}
	}

	return &domain.ListIncidentsResult{
		Items:      items,
		TotalCount: rows[0].TotalCount,
	}, nil
}

func (r *incidentRepository) GetIncident(ctx context.Context, id uuid.UUID) (*domain.Incident, error) {
	row, err := r.queries.GetIncident(ctx, id)
	if err != nil {
		return nil, err
	}

	return toDomainIncidentFromGet(row), nil
}

func (r *incidentRepository) UpdateIncident(ctx context.Context, params domain.UpdateIncidentParams) (*domain.Incident, error) {
	var incident db.Incident

	err := r.store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		incident, err = q.UpdateIncident(ctx, db.UpdateIncidentParams{
			ID:                      params.ID,
			EmployeeID:              params.EmployeeID,
			LocationID:              params.LocationID,
			ReporterInvolvement:     nullIncidentReporterInvolvementFromPtr(params.ReporterInvolvement),
			InformedParties:         toInformedParties(params.InformedParties),
			OccurredAt:              conv.PgTimestamptzFromTime(params.OccurredAt),
			IncidentType:            nullIncidentTypeFromPtr(params.IncidentType),
			SeverityOfIncident:      db.NullSeverityOfIncidentFromPtr(params.SeverityOfIncident),
			IncidentExplanation:     params.IncidentExplanation,
			RecurrenceRisk:          db.NullRecurrenceRiskFromPtr(params.RecurrenceRisk),
			IncidentPreventSteps:    params.IncidentPreventSteps,
			IncidentTakenMeasures:   params.IncidentTakenMeasures,
			CauseCategories:         toCauseCategories(params.CauseCategories),
			CauseExplanation:        params.CauseExplanation,
			PhysicalInjury:          db.NullPhysicalInjuryFromPtr(params.PhysicalInjury),
			PhysicalInjuryDesc:      params.PhysicalInjuryDesc,
			PsychologicalDamage:     db.NullPsychologicalDamageFromPtr(params.PsychologicalDamage),
			PsychologicalDamageDesc: params.PsychologicalDamageDesc,
			NeededConsultation:      db.NullNeededConsultationFromPtr(params.NeededConsultation),
			FollowUpActions:         toFollowUpActions(params.FollowUpActions),
			FollowUpNotes:           params.FollowUpNotes,
			IsEmployeeAbsent:        params.IsEmployeeAbsent,
			AdditionalDetails:       params.AdditionalDetails,
			Emails:                  params.Emails,
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	return toDomainIncidentFromUpdate(incident), nil
}

func (r *incidentRepository) DeleteIncident(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteIncident(ctx, id)
}

func (r *incidentRepository) ConfirmIncident(ctx context.Context, id uuid.UUID, confirmedBy *uuid.UUID) (int64, error) {
	return r.queries.ConfirmIncident(ctx, db.ConfirmIncidentParams{
		ID:          id,
		ConfirmedBy: confirmedBy,
	})
}

func (r *incidentRepository) ListAllIncidents(ctx context.Context, params domain.ListAllIncidentsParams) (*domain.ListAllIncidentsResult, error) {
	var items []db.ListAllIncidentsRow
	var totalCount int64

	err := r.store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		items, err = q.ListAllIncidents(ctx, db.ListAllIncidentsParams{
			Limit:       params.Limit,
			Offset:      params.Offset,
			IsConfirmed: params.IsConfirmed,
			Search:      params.Search,
		})
		if err != nil {
			return err
		}

		totalCount, err = q.CountAllIncidents(ctx, db.CountAllIncidentsParams{
			IsConfirmed: params.IsConfirmed,
			Search:      params.Search,
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return &domain.ListAllIncidentsResult{
			Items:      []domain.IncidentSummary{},
			TotalCount: 0,
		}, nil
	}

	result := make([]domain.IncidentSummary, len(items))
	for i, row := range items {
		result[i] = domain.IncidentSummary{
			ID:                 row.ID,
			OccurredAt:         conv.TimeFromPgTimestamptz(row.OccurredAt),
			IncidentType:       string(row.IncidentType),
			SeverityOfIncident: string(row.SeverityOfIncident),
			IsConfirmed:        row.IsConfirmed,
			ClientFirstName:    row.ClientFirstName,
			ClientLastName:     row.ClientLastName,
			ClientBSN:          row.ClientBsn,
			EmployeeFirstName:  row.EmployeeFirstName,
			EmployeeLastName:   row.EmployeeLastName,
			LocationName:       row.LocationName,
		}
	}

	return &domain.ListAllIncidentsResult{
		Items:      result,
		TotalCount: totalCount,
	}, nil
}

func (r *incidentRepository) GetIncidentCounts(ctx context.Context) (*domain.IncidentCounts, error) {
	row, err := r.queries.GetIncidentCounts(ctx)
	if err != nil {
		return nil, err
	}

	return &domain.IncidentCounts{
		SeriousFatalCount:        row.SeriousFatalCount,
		PendingConfirmationCount: row.PendingConfirmationCount,
		Past24hCount:             row.Past24hCount,
	}, nil
}

func (r *incidentRepository) GetAllAdminUsers(ctx context.Context) ([]uuid.UUID, error) {
	rows, err := r.queries.GetAllAdminUsers(ctx)
	if err != nil {
		return nil, err
	}

	ids := make([]uuid.UUID, len(rows))
	for i, row := range rows {
		ids[i] = row.ID
	}
	return ids, nil
}

// Helper functions

func toInformedParties(values []string) []db.InformedPartyEnum {
	result := make([]db.InformedPartyEnum, 0, len(values))
	for _, value := range values {
		result = append(result, db.InformedPartyEnum(value))
	}
	return result
}

func informedPartiesToStrings(values []db.InformedPartyEnum) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, string(value))
	}
	return result
}

func toCauseCategories(values []string) []db.IncidentCauseCategoryEnum {
	result := make([]db.IncidentCauseCategoryEnum, 0, len(values))
	for _, value := range values {
		result = append(result, db.IncidentCauseCategoryEnum(value))
	}
	return result
}

func causeCategoriesToStrings(values []db.IncidentCauseCategoryEnum) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, string(value))
	}
	return result
}

func toFollowUpActions(values []string) []db.IncidentFollowUpActionEnum {
	result := make([]db.IncidentFollowUpActionEnum, 0, len(values))
	for _, value := range values {
		result = append(result, db.IncidentFollowUpActionEnum(value))
	}
	return result
}

func followUpActionsToStrings(values []db.IncidentFollowUpActionEnum) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, string(value))
	}
	return result
}

func nullIncidentTypeFromPtr(ptr *string) db.NullIncidentTypeEnum {
	if ptr == nil {
		return db.NullIncidentTypeEnum{Valid: false}
	}
	return db.NullIncidentTypeEnum{IncidentTypeEnum: db.IncidentTypeEnum(*ptr), Valid: true}
}

func nullIncidentReporterInvolvementFromPtr(ptr *string) db.NullIncidentReporterInvolvementEnum {
	if ptr == nil {
		return db.NullIncidentReporterInvolvementEnum{Valid: false}
	}
	return db.NullIncidentReporterInvolvementEnum{IncidentReporterInvolvementEnum: db.IncidentReporterInvolvementEnum(*ptr), Valid: true}
}

func toDomainIncident(row *db.CreateIncidentRow) *domain.Incident {
	if row == nil {
		return nil
	}

	return &domain.Incident{
		ID:                      row.ID,
		EmployeeID:              row.EmployeeID,
		EmployeeFirstName:       util.DerefString(row.EmployeeFirstName),
		EmployeeLastName:        util.DerefString(row.EmployeeLastName),
		LocationID:              row.LocationID,
		LocationName:            util.DerefString(row.LocationName),
		ClientID:                row.ClientID,
		ClientFirstName:         util.DerefString(row.ClientFirstName),
		ClientLastName:          util.DerefString(row.ClientLastName),
		ReporterInvolvement:     string(row.ReporterInvolvement),
		InformedParties:         informedPartiesToStrings(row.InformedParties),
		OccurredAt:              conv.TimeFromPgTimestamptz(row.OccurredAt),
		IncidentType:            string(row.IncidentType),
		SeverityOfIncident:      string(row.SeverityOfIncident),
		IncidentExplanation:     row.IncidentExplanation,
		RecurrenceRisk:          string(row.RecurrenceRisk),
		IncidentPreventSteps:    row.IncidentPreventSteps,
		IncidentTakenMeasures:   row.IncidentTakenMeasures,
		CauseCategories:         causeCategoriesToStrings(row.CauseCategories),
		CauseExplanation:        row.CauseExplanation,
		PhysicalInjury:          string(row.PhysicalInjury),
		PhysicalInjuryDesc:      row.PhysicalInjuryDesc,
		PsychologicalDamage:     string(row.PsychologicalDamage),
		PsychologicalDamageDesc: row.PsychologicalDamageDesc,
		NeededConsultation:      string(row.NeededConsultation),
		FollowUpActions:         followUpActionsToStrings(row.FollowUpActions),
		FollowUpNotes:           row.FollowUpNotes,
		IsEmployeeAbsent:        row.IsEmployeeAbsent,
		AdditionalDetails:       row.AdditionalDetails,
		Emails:                  row.Emails,
		UpdatedAt:               conv.TimeFromPgTimestamptz(row.UpdatedAt),
		CreatedAt:               conv.TimeFromPgTimestamptz(row.CreatedAt),
		IsConfirmed:             row.IsConfirmed,
		FileUrl:                 row.FileUrl,
		ConfirmedAt:             conv.TimePtrFromPgTimestamptz(row.ConfirmedAt),
		ConfirmedBy:             row.ConfirmedBy,
		ConfirmationEmailSentAt: conv.TimePtrFromPgTimestamptz(row.ConfirmationEmailSentAt),
	}
}

func toDomainIncidentFromGet(row db.GetIncidentRow) *domain.Incident {
	return &domain.Incident{
		ID:                      row.ID,
		EmployeeID:              row.EmployeeID,
		EmployeeFirstName:       row.EmployeeFirstName,
		EmployeeLastName:        row.EmployeeLastName,
		LocationID:              row.LocationID,
		LocationName:            row.LocationName,
		ClientID:                row.ClientID,
		ClientFirstName:         row.ClientFirstName,
		ClientLastName:          row.ClientLastName,
		ReporterInvolvement:     string(row.ReporterInvolvement),
		InformedParties:         informedPartiesToStrings(row.InformedParties),
		OccurredAt:              conv.TimeFromPgTimestamptz(row.OccurredAt),
		IncidentType:            string(row.IncidentType),
		SeverityOfIncident:      string(row.SeverityOfIncident),
		IncidentExplanation:     row.IncidentExplanation,
		RecurrenceRisk:          string(row.RecurrenceRisk),
		IncidentPreventSteps:    row.IncidentPreventSteps,
		IncidentTakenMeasures:   row.IncidentTakenMeasures,
		CauseCategories:         causeCategoriesToStrings(row.CauseCategories),
		CauseExplanation:        row.CauseExplanation,
		PhysicalInjury:          string(row.PhysicalInjury),
		PhysicalInjuryDesc:      row.PhysicalInjuryDesc,
		PsychologicalDamage:     string(row.PsychologicalDamage),
		PsychologicalDamageDesc: row.PsychologicalDamageDesc,
		NeededConsultation:      string(row.NeededConsultation),
		FollowUpActions:         followUpActionsToStrings(row.FollowUpActions),
		FollowUpNotes:           row.FollowUpNotes,
		IsEmployeeAbsent:        row.IsEmployeeAbsent,
		AdditionalDetails:       row.AdditionalDetails,
		Emails:                  row.Emails,
		UpdatedAt:               conv.TimeFromPgTimestamptz(row.UpdatedAt),
		CreatedAt:               conv.TimeFromPgTimestamptz(row.CreatedAt),
		IsConfirmed:             row.IsConfirmed,
		FileUrl:                 row.FileUrl,
		ConfirmedAt:             conv.TimePtrFromPgTimestamptz(row.ConfirmedAt),
		ConfirmedBy:             row.ConfirmedBy,
		ConfirmationEmailSentAt: conv.TimePtrFromPgTimestamptz(row.ConfirmationEmailSentAt),
	}
}

func toDomainIncidentFromUpdate(row db.Incident) *domain.Incident {
	return &domain.Incident{
		ID:                      row.ID,
		EmployeeID:              row.EmployeeID,
		LocationID:              row.LocationID,
		ClientID:                row.ClientID,
		ReporterInvolvement:     string(row.ReporterInvolvement),
		InformedParties:         informedPartiesToStrings(row.InformedParties),
		OccurredAt:              conv.TimeFromPgTimestamptz(row.OccurredAt),
		IncidentType:            string(row.IncidentType),
		SeverityOfIncident:      string(row.SeverityOfIncident),
		IncidentExplanation:     row.IncidentExplanation,
		RecurrenceRisk:          string(row.RecurrenceRisk),
		IncidentPreventSteps:    row.IncidentPreventSteps,
		IncidentTakenMeasures:   row.IncidentTakenMeasures,
		CauseCategories:         causeCategoriesToStrings(row.CauseCategories),
		CauseExplanation:        row.CauseExplanation,
		PhysicalInjury:          string(row.PhysicalInjury),
		PhysicalInjuryDesc:      row.PhysicalInjuryDesc,
		PsychologicalDamage:     string(row.PsychologicalDamage),
		PsychologicalDamageDesc: row.PsychologicalDamageDesc,
		NeededConsultation:      string(row.NeededConsultation),
		FollowUpActions:         followUpActionsToStrings(row.FollowUpActions),
		FollowUpNotes:           row.FollowUpNotes,
		IsEmployeeAbsent:        row.IsEmployeeAbsent,
		AdditionalDetails:       row.AdditionalDetails,
		Emails:                  row.Emails,
		UpdatedAt:               conv.TimeFromPgTimestamptz(row.UpdatedAt),
		CreatedAt:               conv.TimeFromPgTimestamptz(row.CreatedAt),
		IsConfirmed:             row.IsConfirmed,
		FileUrl:                 row.FileUrl,
		ConfirmedAt:             conv.TimePtrFromPgTimestamptz(row.ConfirmedAt),
		ConfirmedBy:             row.ConfirmedBy,
		ConfirmationEmailSentAt: conv.TimePtrFromPgTimestamptz(row.ConfirmationEmailSentAt),
	}
}
