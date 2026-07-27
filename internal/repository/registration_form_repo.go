package repository

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	db "maicare_go/db/sqlc"
	"maicare_go/internal/domain"
	"maicare_go/pkg/conv"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type RegistrationFormRepository struct {
	store *db.Store
}

func NewRegistrationFormRepository(store *db.Store) domain.RegistrationFormRepository {
	return &RegistrationFormRepository{
		store: store,
	}
}

func (r *RegistrationFormRepository) CreateRegistrationForm(ctx context.Context, params domain.CreateRegistrationFormParams) (*domain.RegistrationForm, error) {
	arg := db.CreateRegistrationFormParams{
		ClientFirstName:               params.ClientFirstName,
		ClientLastName:                params.ClientLastName,
		ClientDateOfBirth:             pgDateFromPtr(params.ClientDateOfBirth),
		ClientBsnNumber:               params.ClientBsnNumber,
		ClientGender:                  db.GenderEnum(params.ClientGender),
		ClientNationality:             params.ClientNationality,
		ClientPhoneNumber:             params.ClientPhoneNumber,
		ClientEmail:                   params.ClientEmail,
		ClientStreet:                  params.ClientStreet,
		ClientHouseNumber:             params.ClientHouseNumber,
		ClientHouseNumberAddition:     params.ClientHouseNumberAddition,
		ClientPostalCode:              params.ClientPostalCode,
		ClientCity:                    params.ClientCity,
		ReferrerFirstName:             params.ReferrerFirstName,
		ReferrerLastName:              params.ReferrerLastName,
		ReferrerOrganization:          params.ReferrerOrganization,
		ReferrerJobTitle:              params.ReferrerJobTitle,
		ReferrerPhoneNumber:           params.ReferrerPhoneNumber,
		ReferrerEmail:                 params.ReferrerEmail,
		Guardian1FirstName:            params.Guardian1FirstName,
		Guardian1LastName:             params.Guardian1LastName,
		Guardian1Relationship:         params.Guardian1Relationship,
		Guardian1PhoneNumber:          params.Guardian1PhoneNumber,
		Guardian1Email:                params.Guardian1Email,
		Guardian2FirstName:            params.Guardian2FirstName,
		Guardian2LastName:             params.Guardian2LastName,
		Guardian2Relationship:         params.Guardian2Relationship,
		Guardian2PhoneNumber:          params.Guardian2PhoneNumber,
		Guardian2Email:                params.Guardian2Email,
		EducationInstitution:          params.EducationInstitution,
		EducationMentorName:           params.EducationMentorName,
		EducationMentorPhone:          params.EducationMentorPhone,
		EducationMentorEmail:          params.EducationMentorEmail,
		EducationCurrentlyEnrolled:    params.EducationCurrentlyEnrolled,
		EducationAdditionalNotes:      params.EducationAdditionalNotes,
		EducationLevel:                db.EducationLevelEnum(defaultEducationLevel(params.EducationLevel)),
		WorkCurrentEmployer:           params.WorkCurrentEmployer,
		WorkEmployerPhone:             params.WorkEmployerPhone,
		WorkEmployerEmail:             params.WorkEmployerEmail,
		WorkCurrentPosition:           params.WorkCurrentPosition,
		WorkCurrentlyEmployed:         params.WorkCurrentlyEmployed,
		WorkStartDate:                 pgDateFromPtr(params.WorkStartDate),
		WorkAdditionalNotes:           params.WorkAdditionalNotes,
		CareProtectedLiving:           params.CareProtectedLiving,
		CareAssistedIndependentLiving: params.CareAssistedIndependentLiving,
		CareRoomTrainingCenter:        params.CareRoomTrainingCenter,
		CareAmbulatoryGuidance:        params.CareAmbulatoryGuidance,
		ApplicationReason:             params.ApplicationReason,
		ClientGoals:                   params.ClientGoals,
		RiskAggressiveBehavior:        params.RiskAggressiveBehavior,
		RiskSuicidalSelfharm:          params.RiskSuicidalSelfharm,
		RiskSubstanceAbuse:            params.RiskSubstanceAbuse,
		RiskPsychiatricIssues:         params.RiskPsychiatricIssues,
		RiskCriminalHistory:           params.RiskCriminalHistory,
		RiskFlightBehavior:            params.RiskFlightBehavior,
		RiskWeaponPossession:          params.RiskWeaponPossession,
		RiskSexualBehavior:            params.RiskSexualBehavior,
		RiskDayNightRhythm:            params.RiskDayNightRhythm,
		RiskOther:                     params.RiskOther,
		RiskOtherDescription:          params.RiskOtherDescription,
		RiskAdditionalNotes:           params.RiskAdditionalNotes,
		DocumentReferral:              params.DocumentReferral,
		DocumentEducationReport:       params.DocumentEducationReport,
		DocumentPsychiatricReport:     params.DocumentPsychiatricReport,
		DocumentDiagnosis:             params.DocumentDiagnosis,
		DocumentSafetyPlan:            params.DocumentSafetyPlan,
		DocumentIDCopy:                params.DocumentIDCopy,
		ApplicationDate:               conv.PgDateFromTime(params.ApplicationDate),
		ReferrerSignature:             params.ReferrerSignature,
	}

	form, err := r.store.CreateRegistrationForm(ctx, arg)
	if err != nil {
		return nil, err
	}

	return toDomainRegistrationForm(form), nil
}

func (r *RegistrationFormRepository) ListRegistrationForms(ctx context.Context, params domain.ListRegistrationFormsParams) (*domain.ListResult[domain.RegistrationFormListItem], error) {
	rows, err := r.store.ListRegistrationForms(ctx, db.ListRegistrationFormsParams{
		Limit:                  params.Limit,
		Offset:                 params.Offset,
		Status:                 db.NullFormStatusFromPtr(params.Status),
		RiskAggressiveBehavior: params.RiskAggressiveBehavior,
		RiskSuicidalSelfharm:   params.RiskSuicidalSelfharm,
		RiskSubstanceAbuse:     params.RiskSubstanceAbuse,
		RiskPsychiatricIssues:  params.RiskPsychiatricIssues,
		RiskCriminalHistory:    params.RiskCriminalHistory,
		RiskFlightBehavior:     params.RiskFlightBehavior,
		RiskWeaponPossession:   params.RiskWeaponPossession,
		RiskSexualBehavior:     params.RiskSexualBehavior,
		RiskDayNightRhythm:     params.RiskDayNightRhythm,
		RiskOther:              nil,
	})
	if err != nil {
		return nil, err
	}

	totalCount, err := r.store.CountRegistrationForms(ctx, db.CountRegistrationFormsParams{
		Status:                 db.NullFormStatusFromPtr(params.Status),
		RiskAggressiveBehavior: params.RiskAggressiveBehavior,
		RiskSuicidalSelfharm:   params.RiskSuicidalSelfharm,
		RiskSubstanceAbuse:     params.RiskSubstanceAbuse,
		RiskPsychiatricIssues:  params.RiskPsychiatricIssues,
		RiskCriminalHistory:    params.RiskCriminalHistory,
		RiskFlightBehavior:     params.RiskFlightBehavior,
		RiskWeaponPossession:   params.RiskWeaponPossession,
		RiskSexualBehavior:     params.RiskSexualBehavior,
		RiskDayNightRhythm:     params.RiskDayNightRhythm,
	})
	if err != nil {
		return nil, err
	}

	items := make([]domain.RegistrationFormListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, toDomainRegistrationFormListItem(row))
	}

	return &domain.ListResult[domain.RegistrationFormListItem]{
		Items:      items,
		TotalCount: totalCount,
	}, nil
}

func (r *RegistrationFormRepository) GetRegistrationFormCounts(ctx context.Context) (domain.RegistrationFormCounts, error) {
	counts, err := r.store.GetRegistrationFormCounts(ctx)
	if err != nil {
		return domain.RegistrationFormCounts{}, err
	}
	return domain.RegistrationFormCounts{
		Total:         counts.Total,
		PendingReview: counts.PendingReview,
		Processed:     counts.Processed,
		HighRisk:      counts.HighRisk,
	}, nil
}

func (r *RegistrationFormRepository) GetRegistrationForm(ctx context.Context, id uuid.UUID) (*domain.RegistrationForm, error) {
	row, err := r.store.GetRegistrationForm(ctx, id)
	if err != nil {
		if isDBNotFound(err) {
			return nil, domain.ErrRegistrationFormNotFound
		}
		return nil, err
	}

	// Collect all document UUIDs
	var documentIDs []uuid.UUID
	if row.DocumentReferral != nil {
		documentIDs = append(documentIDs, *row.DocumentReferral)
	}
	if row.DocumentEducationReport != nil {
		documentIDs = append(documentIDs, *row.DocumentEducationReport)
	}
	if row.DocumentActionPlan != nil {
		documentIDs = append(documentIDs, *row.DocumentActionPlan)
	}
	if row.DocumentPsychiatricReport != nil {
		documentIDs = append(documentIDs, *row.DocumentPsychiatricReport)
	}
	if row.DocumentDiagnosis != nil {
		documentIDs = append(documentIDs, *row.DocumentDiagnosis)
	}
	if row.DocumentSafetyPlan != nil {
		documentIDs = append(documentIDs, *row.DocumentSafetyPlan)
	}
	if row.DocumentIDCopy != nil {
		documentIDs = append(documentIDs, *row.DocumentIDCopy)
	}

	// Fetch documents if any
	documentMap := make(map[uuid.UUID]domain.Document)
	if len(documentIDs) > 0 {
		attachments, err := r.store.GetAttachmentsByUUIDs(ctx, documentIDs)
		if err != nil {
			return nil, err
		}
		for _, att := range attachments {
			documentMap[att.Uuid] = domain.Document{
				ID:   att.Uuid,
				Name: att.Name,
				File: att.File,
				Size: att.Size,
			}
		}
	}

	return toDomainRegistrationFormFromGetRow(row, documentMap), nil
}

func (r *RegistrationFormRepository) UpdateRegistrationForm(ctx context.Context, params domain.UpdateRegistrationFormParams) (*domain.RegistrationForm, error) {
	arg := db.UpdateRegistrationFormParams{
		ID:                            params.ID,
		ClientFirstName:               params.ClientFirstName,
		ClientLastName:                params.ClientLastName,
		ClientDateOfBirth:             pgDateFromPtr(params.ClientDateOfBirth),
		ClientBsnNumber:               params.ClientBsnNumber,
		ClientGender:                  db.NullGenderFromPtr(params.ClientGender),
		ClientNationality:             params.ClientNationality,
		ClientPhoneNumber:             params.ClientPhoneNumber,
		ClientEmail:                   params.ClientEmail,
		ClientStreet:                  params.ClientStreet,
		ClientHouseNumber:             params.ClientHouseNumber,
		ClientHouseNumberAddition:     params.ClientHouseNumberAddition,
		ClientPostalCode:              params.ClientPostalCode,
		ClientCity:                    params.ClientCity,
		ReferrerFirstName:             params.ReferrerFirstName,
		ReferrerLastName:              params.ReferrerLastName,
		ReferrerOrganization:          params.ReferrerOrganization,
		ReferrerJobTitle:              params.ReferrerJobTitle,
		ReferrerPhoneNumber:           params.ReferrerPhoneNumber,
		ReferrerEmail:                 params.ReferrerEmail,
		Guardian1FirstName:            params.Guardian1FirstName,
		Guardian1LastName:             params.Guardian1LastName,
		Guardian1Relationship:         params.Guardian1Relationship,
		Guardian1PhoneNumber:          params.Guardian1PhoneNumber,
		Guardian1Email:                params.Guardian1Email,
		Guardian2FirstName:            params.Guardian2FirstName,
		Guardian2LastName:             params.Guardian2LastName,
		Guardian2Relationship:         params.Guardian2Relationship,
		Guardian2PhoneNumber:          params.Guardian2PhoneNumber,
		Guardian2Email:                params.Guardian2Email,
		EducationInstitution:          params.EducationInstitution,
		EducationMentorName:           params.EducationMentorName,
		EducationMentorPhone:          params.EducationMentorPhone,
		EducationMentorEmail:          params.EducationMentorEmail,
		EducationCurrentlyEnrolled:    params.EducationCurrentlyEnrolled,
		EducationAdditionalNotes:      params.EducationAdditionalNotes,
		EducationLevel:                db.NullClientEducationLevelFromPtr(params.EducationLevel),
		WorkCurrentEmployer:           params.WorkCurrentEmployer,
		WorkEmployerPhone:             params.WorkEmployerPhone,
		WorkEmployerEmail:             params.WorkEmployerEmail,
		WorkCurrentPosition:           params.WorkCurrentPosition,
		WorkCurrentlyEmployed:         params.WorkCurrentlyEmployed,
		WorkStartDate:                 pgDateFromPtr(params.WorkStartDate),
		WorkAdditionalNotes:           params.WorkAdditionalNotes,
		CareProtectedLiving:           params.CareProtectedLiving,
		CareAssistedIndependentLiving: params.CareAssistedIndependentLiving,
		CareRoomTrainingCenter:        params.CareRoomTrainingCenter,
		CareAmbulatoryGuidance:        params.CareAmbulatoryGuidance,
		ApplicationReason:             params.ApplicationReason,
		ClientGoals:                   params.ClientGoals,
		RiskAggressiveBehavior:        params.RiskAggressiveBehavior,
		RiskSuicidalSelfharm:          params.RiskSuicidalSelfharm,
		RiskSubstanceAbuse:            params.RiskSubstanceAbuse,
		RiskPsychiatricIssues:         params.RiskPsychiatricIssues,
		RiskCriminalHistory:           params.RiskCriminalHistory,
		RiskFlightBehavior:            params.RiskFlightBehavior,
		RiskWeaponPossession:          params.RiskWeaponPossession,
		RiskSexualBehavior:            params.RiskSexualBehavior,
		RiskDayNightRhythm:            params.RiskDayNightRhythm,
		RiskOther:                     params.RiskOther,
		RiskOtherDescription:          params.RiskOtherDescription,
		RiskAdditionalNotes:           params.RiskAdditionalNotes,
		ApplicationDate:               pgDateFromPtr(params.ApplicationDate),
		ReferrerSignature:             params.ReferrerSignature,
	}

	form, err := r.store.UpdateRegistrationForm(ctx, arg)
	if err != nil {
		if isDBNotFound(err) {
			return nil, domain.ErrRegistrationFormNotFound
		}
		return nil, err
	}

	return toDomainRegistrationForm(form), nil
}

func (r *RegistrationFormRepository) ReplaceRegistrationFormDocument(ctx context.Context, params domain.ReplaceRegistrationFormDocumentParams) (*domain.RegistrationForm, error) {
	form, err := r.store.ReplaceRegistrationFormDocument(ctx, db.ReplaceRegistrationFormDocumentParams{
		ID:           params.ID,
		DocumentType: params.DocumentType,
		FileID:       &params.FileID,
	})
	if err != nil {
		if isDBNotFound(err) {
			return nil, domain.ErrRegistrationFormNotFound
		}
		return nil, err
	}
	return toDomainRegistrationForm(form), nil
}

func (r *RegistrationFormRepository) DeleteRegistrationForm(ctx context.Context, id uuid.UUID) error {
	err := r.store.DeleteRegistrationForm(ctx, id)
	if err != nil {
		if isDBNotFound(err) {
			return domain.ErrRegistrationFormNotFound
		}
		return err
	}
	return nil
}

func (r *RegistrationFormRepository) UpdateRegistrationFormStatus(ctx context.Context, params domain.UpdateRegistrationFormStatusParams) error {
	arg := db.UpdateRegistrationFormStatusParams{
		ID:                        params.ID,
		FormStatus:                db.FormStatusEnum(params.Status),
		ProcessedByEmployeeID:     params.ProcessedByEmployeeID,
		IntakeAppointmentLocation: params.IntakeAppointmentLocation,
		AddmissionType:            nullAdmissionTypeEnumFromPtr(params.AddmissionType),
		RejectionReason:           params.RejectionReason,
	}

	_, err := r.store.UpdateRegistrationFormStatus(ctx, arg)
	if err != nil {
		if isDBNotFound(err) {
			return domain.ErrRegistrationFormNotFound
		}
		return err
	}

	return nil
}

func (r *RegistrationFormRepository) ProcessRegistrationForm(ctx context.Context, params domain.ProcessRegistrationFormParams) (*domain.RegistrationForm, string, error) {
	optionsJSON, err := json.Marshal(params.ProposedDates)
	if err != nil {
		return nil, "", err
	}

	token := randomString(32)

	arg := db.UpdateRegistrationFormStatusParams{
		ID:                        params.ID,
		FormStatus:                db.FormStatusEnum("processed"),
		ProcessedByEmployeeID:     &params.EmployeeID,
		IntakeAppointmentLocation: &params.IntakeAppointmentLocation,
		AddmissionType:            nullAdmissionTypeEnumFromPtr(&params.AddmissionType),
		IntakeOptions:             optionsJSON,
		IntakeToken:               &token,
	}

	form, err := r.store.UpdateRegistrationFormStatus(ctx, arg)
	if err != nil {
		if isDBNotFound(err) {
			return nil, "", domain.ErrRegistrationFormNotFound
		}
		return nil, "", err
	}

	return toDomainRegistrationForm(form), token, nil
}

func (r *RegistrationFormRepository) GetPublicIntakeOptions(ctx context.Context, token string) (*domain.PublicIntakeOptions, error) {
	form, err := r.store.GetRegistrationFormByToken(ctx, &token)
	if err != nil {
		if isDBNotFound(err) {
			return nil, domain.ErrRegistrationFormNotFound
		}
		return nil, err
	}

	var dates []string
	if len(form.IntakeOptions) > 0 {
		if err := json.Unmarshal(form.IntakeOptions, &dates); err != nil {
			return nil, err
		}
	}

	return &domain.PublicIntakeOptions{
		ClientFirstName: form.ClientFirstName,
		IntakeLocation:  derefString(form.IntakeAppointmentLocation),
		ProposedDates:   dates,
	}, nil
}

func (r *RegistrationFormRepository) SelectIntakeDate(ctx context.Context, params domain.SelectIntakeDateParams) error {
	form, err := r.store.GetRegistrationFormByToken(ctx, &params.Token)
	if err != nil {
		if isDBNotFound(err) {
			return domain.ErrRegistrationFormNotFound
		}
		return err
	}

	// Validate selected date is in options (permissive - doesn't strictly fail)
	if len(form.IntakeOptions) > 0 {
		var dates []string
		if err := json.Unmarshal(form.IntakeOptions, &dates); err != nil {
			// Permissive - continue even if unmarshal fails
		} else {
			// Validation is permissive per old implementation
			_ = dates
		}
	}

	arg := db.UpdateRegistrationFormIntakeDateParams{
		ID:                        form.ID,
		IntakeAppointmentDatetime: pgtype.Timestamptz{Time: params.SelectedDate, Valid: true},
	}

	_, err = r.store.UpdateRegistrationFormIntakeDate(ctx, arg)
	if err != nil {
		return err
	}

	return nil
}

// Helper functions

func toDomainRegistrationForm(row db.RegistrationForm) *domain.RegistrationForm {
	return &domain.RegistrationForm{
		ID:                            row.ID,
		ClientFirstName:               row.ClientFirstName,
		ClientLastName:                row.ClientLastName,
		ClientDateOfBirth:             conv.TimePtrFromPgDate(row.ClientDateOfBirth),
		ClientBsnNumber:               row.ClientBsnNumber,
		ClientGender:                  string(row.ClientGender),
		ClientNationality:             row.ClientNationality,
		ClientPhoneNumber:             row.ClientPhoneNumber,
		ClientEmail:                   row.ClientEmail,
		ClientStreet:                  row.ClientStreet,
		ClientHouseNumber:             row.ClientHouseNumber,
		ClientHouseNumberAddition:     row.ClientHouseNumberAddition,
		ClientPostalCode:              row.ClientPostalCode,
		ClientCity:                    row.ClientCity,
		ReferrerFirstName:             row.ReferrerFirstName,
		ReferrerLastName:              row.ReferrerLastName,
		ReferrerOrganization:          row.ReferrerOrganization,
		ReferrerJobTitle:              row.ReferrerJobTitle,
		ReferrerPhoneNumber:           row.ReferrerPhoneNumber,
		ReferrerEmail:                 row.ReferrerEmail,
		Guardian1FirstName:            row.Guardian1FirstName,
		Guardian1LastName:             row.Guardian1LastName,
		Guardian1Relationship:         row.Guardian1Relationship,
		Guardian1PhoneNumber:          row.Guardian1PhoneNumber,
		Guardian1Email:                row.Guardian1Email,
		Guardian2FirstName:            row.Guardian2FirstName,
		Guardian2LastName:             row.Guardian2LastName,
		Guardian2Relationship:         row.Guardian2Relationship,
		Guardian2PhoneNumber:          row.Guardian2PhoneNumber,
		Guardian2Email:                row.Guardian2Email,
		EducationInstitution:          row.EducationInstitution,
		EducationMentorName:           row.EducationMentorName,
		EducationMentorPhone:          row.EducationMentorPhone,
		EducationMentorEmail:          row.EducationMentorEmail,
		EducationCurrentlyEnrolled:    row.EducationCurrentlyEnrolled,
		EducationAdditionalNotes:      row.EducationAdditionalNotes,
		EducationLevel:                string(row.EducationLevel),
		WorkCurrentEmployer:           row.WorkCurrentEmployer,
		WorkEmployerPhone:             row.WorkEmployerPhone,
		WorkEmployerEmail:             row.WorkEmployerEmail,
		WorkCurrentPosition:           row.WorkCurrentPosition,
		WorkCurrentlyEmployed:         row.WorkCurrentlyEmployed,
		WorkStartDate:                 conv.TimePtrFromPgDate(row.WorkStartDate),
		WorkAdditionalNotes:           row.WorkAdditionalNotes,
		CareProtectedLiving:           row.CareProtectedLiving,
		CareAssistedIndependentLiving: row.CareAssistedIndependentLiving,
		CareRoomTrainingCenter:        row.CareRoomTrainingCenter,
		CareAmbulatoryGuidance:        row.CareAmbulatoryGuidance,
		ApplicationReason:             row.ApplicationReason,
		ClientGoals:                   row.ClientGoals,
		RiskAggressiveBehavior:        row.RiskAggressiveBehavior,
		RiskSuicidalSelfharm:          row.RiskSuicidalSelfharm,
		RiskSubstanceAbuse:            row.RiskSubstanceAbuse,
		RiskPsychiatricIssues:         row.RiskPsychiatricIssues,
		RiskCriminalHistory:           row.RiskCriminalHistory,
		RiskFlightBehavior:            row.RiskFlightBehavior,
		RiskWeaponPossession:          row.RiskWeaponPossession,
		RiskSexualBehavior:            row.RiskSexualBehavior,
		RiskDayNightRhythm:            row.RiskDayNightRhythm,
		RiskOther:                     row.RiskOther,
		RiskOtherDescription:          row.RiskOtherDescription,
		RiskAdditionalNotes:           row.RiskAdditionalNotes,
		DocumentReferral:              nil,
		DocumentEducationReport:       nil,
		DocumentActionPlan:            nil,
		DocumentPsychiatricReport:     nil,
		DocumentDiagnosis:             nil,
		DocumentSafetyPlan:            nil,
		DocumentIDCopy:                nil,
		ApplicationDate:               conv.TimeFromPgDate(row.ApplicationDate),
		ReferrerSignature:             row.ReferrerSignature,
		FormStatus:                    string(row.FormStatus),
		CreatedAt:                     conv.TimeFromPgTimestamptz(row.CreatedAt),
		UpdatedAt:                     conv.TimeFromPgTimestamptz(row.UpdatedAt),
		SubmittedAt:                   conv.TimeFromPgTimestamptz(row.SubmittedAt),
		ProcessedAt:                   conv.TimeFromPgTimestamptz(row.ProcessedAt),
		ProcessedByEmployeeID:         row.ProcessedByEmployeeID,
		ProcessedByEmployeeName:       nil,
		IntakeAppointmentDate:         conv.TimeFromPgTimestamptz(row.IntakeAppointmentDatetime),
		IntakeAppointmentLocation:     row.IntakeAppointmentLocation,
		AddmissionType:                string(row.AddmissionType),
		IntakeOptions:                 nil,
		IntakeFormID:                  nil,
		RejectionReason:               row.RejectionReason,
	}
}

func toDomainRegistrationFormFromGetRow(row db.GetRegistrationFormRow, documentMap map[uuid.UUID]domain.Document) *domain.RegistrationForm {
	form := &domain.RegistrationForm{
		ID:                            row.ID,
		ClientFirstName:               row.ClientFirstName,
		ClientLastName:                row.ClientLastName,
		ClientDateOfBirth:             conv.TimePtrFromPgDate(row.ClientDateOfBirth),
		ClientBsnNumber:               row.ClientBsnNumber,
		ClientGender:                  string(row.ClientGender),
		ClientNationality:             row.ClientNationality,
		ClientPhoneNumber:             row.ClientPhoneNumber,
		ClientEmail:                   row.ClientEmail,
		ClientStreet:                  row.ClientStreet,
		ClientHouseNumber:             row.ClientHouseNumber,
		ClientHouseNumberAddition:     row.ClientHouseNumberAddition,
		ClientPostalCode:              row.ClientPostalCode,
		ClientCity:                    row.ClientCity,
		ReferrerFirstName:             row.ReferrerFirstName,
		ReferrerLastName:              row.ReferrerLastName,
		ReferrerOrganization:          row.ReferrerOrganization,
		ReferrerJobTitle:              row.ReferrerJobTitle,
		ReferrerPhoneNumber:           row.ReferrerPhoneNumber,
		ReferrerEmail:                 row.ReferrerEmail,
		Guardian1FirstName:            row.Guardian1FirstName,
		Guardian1LastName:             row.Guardian1LastName,
		Guardian1Relationship:         row.Guardian1Relationship,
		Guardian1PhoneNumber:          row.Guardian1PhoneNumber,
		Guardian1Email:                row.Guardian1Email,
		Guardian2FirstName:            row.Guardian2FirstName,
		Guardian2LastName:             row.Guardian2LastName,
		Guardian2Relationship:         row.Guardian2Relationship,
		Guardian2PhoneNumber:          row.Guardian2PhoneNumber,
		Guardian2Email:                row.Guardian2Email,
		EducationInstitution:          row.EducationInstitution,
		EducationMentorName:           row.EducationMentorName,
		EducationMentorPhone:          row.EducationMentorPhone,
		EducationMentorEmail:          row.EducationMentorEmail,
		EducationCurrentlyEnrolled:    row.EducationCurrentlyEnrolled,
		EducationAdditionalNotes:      row.EducationAdditionalNotes,
		EducationLevel:                string(row.EducationLevel),
		WorkCurrentEmployer:           row.WorkCurrentEmployer,
		WorkEmployerPhone:             row.WorkEmployerPhone,
		WorkEmployerEmail:             row.WorkEmployerEmail,
		WorkCurrentPosition:           row.WorkCurrentPosition,
		WorkCurrentlyEmployed:         row.WorkCurrentlyEmployed,
		WorkStartDate:                 conv.TimePtrFromPgDate(row.WorkStartDate),
		WorkAdditionalNotes:           row.WorkAdditionalNotes,
		CareProtectedLiving:           row.CareProtectedLiving,
		CareAssistedIndependentLiving: row.CareAssistedIndependentLiving,
		CareRoomTrainingCenter:        row.CareRoomTrainingCenter,
		CareAmbulatoryGuidance:        row.CareAmbulatoryGuidance,
		ApplicationReason:             row.ApplicationReason,
		ClientGoals:                   row.ClientGoals,
		RiskAggressiveBehavior:        row.RiskAggressiveBehavior,
		RiskSuicidalSelfharm:          row.RiskSuicidalSelfharm,
		RiskSubstanceAbuse:            row.RiskSubstanceAbuse,
		RiskPsychiatricIssues:         row.RiskPsychiatricIssues,
		RiskCriminalHistory:           row.RiskCriminalHistory,
		RiskFlightBehavior:            row.RiskFlightBehavior,
		RiskWeaponPossession:          row.RiskWeaponPossession,
		RiskSexualBehavior:            row.RiskSexualBehavior,
		RiskDayNightRhythm:            row.RiskDayNightRhythm,
		RiskOther:                     row.RiskOther,
		RiskOtherDescription:          row.RiskOtherDescription,
		RiskAdditionalNotes:           row.RiskAdditionalNotes,
		ApplicationDate:               conv.TimeFromPgDate(row.ApplicationDate),
		ReferrerSignature:             row.ReferrerSignature,
		FormStatus:                    string(row.FormStatus),
		CreatedAt:                     conv.TimeFromPgTimestamptz(row.CreatedAt),
		UpdatedAt:                     conv.TimeFromPgTimestamptz(row.UpdatedAt),
		SubmittedAt:                   conv.TimeFromPgTimestamptz(row.SubmittedAt),
		ProcessedAt:                   conv.TimeFromPgTimestamptz(row.ProcessedAt),
		ProcessedByEmployeeID:         row.ProcessedByEmployeeID,
		IntakeAppointmentDate:         conv.TimeFromPgTimestamptz(row.IntakeAppointmentDatetime),
		IntakeAppointmentLocation:     row.IntakeAppointmentLocation,
		AddmissionType:                string(row.AddmissionType),
		IntakeFormID:                  row.IntakeFormID,
		RejectionReason:               row.RejectionReason,
	}

	// Build processed by employee name
	if row.ProcessedByFirstName != nil && row.ProcessedByLastName != nil {
		name := *row.ProcessedByFirstName + " " + *row.ProcessedByLastName
		form.ProcessedByEmployeeName = &name
	}

	// Unmarshal intake options
	if len(row.IntakeOptions) > 0 {
		var options []string
		if err := json.Unmarshal(row.IntakeOptions, &options); err == nil {
			form.IntakeOptions = options
		}
	}

	// Map documents
	getDocumentPtr := func(id *uuid.UUID) *domain.Document {
		if id == nil {
			return nil
		}
		if doc, ok := documentMap[*id]; ok {
			return &doc
		}
		return nil
	}

	form.DocumentReferral = getDocumentPtr(row.DocumentReferral)
	form.DocumentEducationReport = getDocumentPtr(row.DocumentEducationReport)
	form.DocumentActionPlan = getDocumentPtr(row.DocumentActionPlan)
	form.DocumentPsychiatricReport = getDocumentPtr(row.DocumentPsychiatricReport)
	form.DocumentDiagnosis = getDocumentPtr(row.DocumentDiagnosis)
	form.DocumentSafetyPlan = getDocumentPtr(row.DocumentSafetyPlan)
	form.DocumentIDCopy = getDocumentPtr(row.DocumentIDCopy)

	return form
}

func toDomainRegistrationFormListItem(row db.ListRegistrationFormsRow) domain.RegistrationFormListItem {
	return domain.RegistrationFormListItem{
		ID:                            row.ID,
		ClientFirstName:               row.ClientFirstName,
		ClientLastName:                row.ClientLastName,
		ClientBsnNumber:               row.ClientBsnNumber,
		ReferrerFirstName:             row.ReferrerFirstName,
		ReferrerLastName:              row.ReferrerLastName,
		CareProtectedLiving:           row.CareProtectedLiving,
		CareAssistedIndependentLiving: row.CareAssistedIndependentLiving,
		CareRoomTrainingCenter:        row.CareRoomTrainingCenter,
		CareAmbulatoryGuidance:        row.CareAmbulatoryGuidance,
		RiskCount:                     calculateRiskCount(row.RiskAggressiveBehavior, row.RiskSuicidalSelfharm, row.RiskSubstanceAbuse, row.RiskPsychiatricIssues, row.RiskCriminalHistory, row.RiskFlightBehavior, row.RiskWeaponPossession, row.RiskSexualBehavior, row.RiskDayNightRhythm, row.RiskOther),
		FormStatus:                    string(row.FormStatus),
		SubmittedAt:                   conv.TimeFromPgTimestamptz(row.SubmittedAt),
		IntakeFormID:                  row.IntakeFormID,
	}
}

func calculateRiskCount(
	riskAggressiveBehavior *bool,
	riskSuicidalSelfharm *bool,
	riskSubstanceAbuse *bool,
	riskPsychiatricIssues *bool,
	riskCriminalHistory *bool,
	riskFlightBehavior *bool,
	riskWeaponPossession *bool,
	riskSexualBehavior *bool,
	riskDayNightRhythm *bool,
	riskOther *bool,
) int {
	count := 0
	isTruePtr := func(b *bool) bool {
		return b != nil && *b
	}

	if isTruePtr(riskAggressiveBehavior) {
		count++
	}
	if isTruePtr(riskSuicidalSelfharm) {
		count++
	}
	if isTruePtr(riskSubstanceAbuse) {
		count++
	}
	if isTruePtr(riskPsychiatricIssues) {
		count++
	}
	if isTruePtr(riskCriminalHistory) {
		count++
	}
	if isTruePtr(riskFlightBehavior) {
		count++
	}
	if isTruePtr(riskWeaponPossession) {
		count++
	}
	if isTruePtr(riskSexualBehavior) {
		count++
	}
	if isTruePtr(riskDayNightRhythm) {
		count++
	}
	if isTruePtr(riskOther) {
		count++
	}

	return count
}

func nullAdmissionTypeEnumFromPtr(value *string) *db.AdmissionTypeEnum {
	if value == nil {
		return nil
	}
	admissionType := db.AdmissionTypeEnum(*value)
	return &admissionType
}

func randomString(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)[:n]
}

func defaultEducationLevel(level string) string {
	if level == "" {
		return "none"
	}
	return level
}

var _ domain.RegistrationFormRepository = (*RegistrationFormRepository)(nil)
