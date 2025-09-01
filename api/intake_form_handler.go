package api

import (
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// // IntakeFormUploadHandlerResponse represents a response from the intake form upload handler
// type IntakeFormUploadHandlerResponse struct {
// 	FileURL   string    `json:"file_url"`
// 	FileID    uuid.UUID `json:"file_id"`
// 	CreatedAt time.Time `json:"created_at"`
// 	Size      int64     `json:"size"`
// }

// // @Summary Upload a file for intake form
// // @Description Upload a file for intake form
// // @Tags intake_form
// // @Accept mpfd
// // @Produce json
// // @Param file formData file true "File to upload"
// // @Success 201 {object} Response[IntakeFormUploadHandlerResponse]
// // @Failure 400 {object} Response[any] "Bad request"
// // @Failure 401 {object} Response[any] "Unauthorized"
// // @Failure 413 {object} Response[any] "Request entity too large"
// // @Failure 500 {object} Response[any] "Internal server error"
// // @Router /intake_form/upload [post]
// // @Security -
// func (server *Server) IntakeFormUploadHandlerApi(ctx *gin.Context) {

// 	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxFileSize)
// 	file, header, err := ctx.Request.FormFile("file")
// 	if err != nil {
// 		if strings.Contains(err.Error(), "request body too large") {
// 			ctx.JSON(http.StatusRequestEntityTooLarge, errorResponse(fmt.Errorf("file size exceeds maximum limit of 10MB")))
// 			return
// 		}
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	// Basic validations
// 	if err := bucket.ValidateFile(header, maxFileSize); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	filename := bucket.GenerateUniqueFilename(header.Filename)

// 	buff := make([]byte, 512)
// 	_, err = file.Read(buff)
// 	if err != nil && err != io.EOF {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("error reading file: %v", err)))
// 		return
// 	}

// 	// Reset file pointer after reading
// 	if _, err := file.Seek(0, 0); err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("error resetting file: %v", err)))
// 		return
// 	}

// 	// Verify content type
// 	contentType := http.DetectContentType(buff)
// 	if !allowedMimeTypes[contentType] {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("unsupported file type: %s", contentType)))
// 		return
// 	}

// 	err = server.b2Client.UploadToB2(ctx.Request.Context(), file, filename)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	fileURL := fmt.Sprintf("%s/file/%s/%s",
// 		server.b2Client.Bucket.BaseURL(),
// 		server.b2Client.Bucket.Name(),
// 		filename)

// 	arg := db.CreateAttachmentParams{
// 		Name: filename,
// 		File: fileURL,
// 		Size: int32(header.Size),
// 		Tag:  util.StringPtr(""),
// 	}
// 	attachment, err := server.store.CreateAttachment(ctx, arg)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	res := SuccessResponse(UploadHandlerResponse{
// 		FileURL:   fileURL,
// 		FileID:    attachment.Uuid,
// 		CreatedAt: attachment.Created.Time,
// 		Size:      int64(attachment.Size),
// 	}, "File uploaded successfully")

// 	ctx.JSON(http.StatusCreated, res)
// }

type GuardionInfo struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Address     string `json:"address"`
}

// CreateIntakeFormRequest represents a request to create an intake form
type CreateIntakeFormRequest struct {
	FirstName             string         `json:"first_name"`
	LastName              string         `json:"last_name"`
	DateOfBirth           time.Time      `json:"date_of_birth"`
	Nationality           string         `json:"nationality"`
	Bsn                   string         `json:"bsn"`
	Address               string         `json:"address"`
	City                  string         `json:"city"`
	PostalCode            string         `json:"postal_code"`
	PhoneNumber           string         `json:"phone_number"`
	Gender                string         `json:"gender"`
	Email                 string         `json:"email"`
	IDType                string         `json:"id_type"`
	IDNumber              string         `json:"id_number"`
	ReferrerName          *string        `json:"referrer_name"`
	ReferrerOrganization  *string        `json:"referrer_organization"`
	ReferrerFunction      *string        `json:"referrer_function"`
	ReferrerPhone         *string        `json:"referrer_phone"`
	ReferrerEmail         *string        `json:"referrer_email"`
	SignedBy              *string        `json:"signed_by"`
	HasValidIndication    bool           `json:"has_valid_indication"`
	LawType               *string        `json:"law_type"`
	OtherLawSpecification *string        `json:"other_law_specification"`
	MainProviderName      *string        `json:"main_provider_name"`
	MainProviderContact   *string        `json:"main_provider_contact"`
	IndicationStartDate   time.Time      `json:"indication_start_date"`
	IndicationEndDate     time.Time      `json:"indication_end_date"`
	RegistrationReason    *string        `json:"registration_reason"`
	GuidanceGoals         *string        `json:"guidance_goals"`
	RegistrationType      *string        `json:"registration_type"`
	LivingSituation       *string        `json:"living_situation"`
	OtherLivingSituation  *string        `json:"other_living_situation"`
	ParentalAuthority     bool           `json:"parental_authority"`
	CurrentSchool         *string        `json:"current_school"`
	MentorName            *string        `json:"mentor_name"`
	MentorPhone           *string        `json:"mentor_phone"`
	MentorEmail           *string        `json:"mentor_email"`
	PreviousCare          *string        `json:"previous_care"`
	GuardianDetails       []GuardionInfo `json:"guardian_details"`
	Diagnoses             *string        `json:"diagnoses"`
	UsesMedication        bool           `json:"uses_medication"`
	MedicationDetails     *string        `json:"medication_details"`
	AddictionIssues       bool           `json:"addiction_issues"`
	JudicialInvolvement   bool           `json:"judicial_involvement"`
	RiskAggression        bool           `json:"risk_aggression"`
	RiskSuicidality       bool           `json:"risk_suicidality"`
	RiskRunningAway       bool           `json:"risk_running_away"`
	RiskSelfHarm          bool           `json:"risk_self_harm"`
	RiskWeaponPossession  bool           `json:"risk_weapon_possession"`
	RiskDrugDealing       bool           `json:"risk_drug_dealing"`
	OtherRisks            *string        `json:"other_risks"`
	SharingPermission     bool           `json:"sharing_permission"`
	TruthDeclaration      bool           `json:"truth_declaration"`
	ClientSignature       bool           `json:"client_signature"`
	GuardianSignature     *bool          `json:"guardian_signature"`
	ReferrerSignature     *bool          `json:"referrer_signature"`
	SignatureDate         time.Time      `json:"signature_date"`
	AttachementIds        []uuid.UUID    `json:"attachement_ids"`
	UrgencyScore          string         `json:"urgency_score" binding:"oneof=low medium high"`
}

// CreateIntakeFormResponse represents a response from the create intake form handler
type CreateIntakeFormResponse struct {
	ID                    int64          `json:"id"`
	FirstName             string         `json:"first_name"`
	LastName              string         `json:"last_name"`
	DateOfBirth           time.Time      `json:"date_of_birth"`
	Nationality           string         `json:"nationality"`
	Bsn                   string         `json:"bsn"`
	Address               string         `json:"address"`
	City                  string         `json:"city"`
	PostalCode            string         `json:"postal_code"`
	PhoneNumber           string         `json:"phone_number"`
	Gender                string         `json:"gender"`
	Email                 string         `json:"email"`
	IDType                string         `json:"id_type"`
	IDNumber              string         `json:"id_number"`
	ReferrerName          *string        `json:"referrer_name"`
	ReferrerOrganization  *string        `json:"referrer_organization"`
	ReferrerFunction      *string        `json:"referrer_function"`
	ReferrerPhone         *string        `json:"referrer_phone"`
	ReferrerEmail         *string        `json:"referrer_email"`
	SignedBy              *string        `json:"signed_by"`
	HasValidIndication    bool           `json:"has_valid_indication"`
	LawType               *string        `json:"law_type"`
	OtherLawSpecification *string        `json:"other_law_specification"`
	MainProviderName      *string        `json:"main_provider_name"`
	MainProviderContact   *string        `json:"main_provider_contact"`
	IndicationStartDate   time.Time      `json:"indication_start_date"`
	IndicationEndDate     time.Time      `json:"indication_end_date"`
	RegistrationReason    *string        `json:"registration_reason"`
	GuidanceGoals         *string        `json:"guidance_goals"`
	RegistrationType      *string        `json:"registration_type"`
	LivingSituation       *string        `json:"living_situation"`
	OtherLivingSituation  *string        `json:"other_living_situation"`
	ParentalAuthority     bool           `json:"parental_authority"`
	CurrentSchool         *string        `json:"current_school"`
	MentorName            *string        `json:"mentor_name"`
	MentorPhone           *string        `json:"mentor_phone"`
	MentorEmail           *string        `json:"mentor_email"`
	PreviousCare          *string        `json:"previous_care"`
	GuardianDetails       []GuardionInfo `json:"guardian_details"`
	Diagnoses             *string        `json:"diagnoses"`
	UsesMedication        bool           `json:"uses_medication"`
	MedicationDetails     *string        `json:"medication_details"`
	AddictionIssues       bool           `json:"addiction_issues"`
	JudicialInvolvement   bool           `json:"judicial_involvement"`
	RiskAggression        bool           `json:"risk_aggression"`
	RiskSuicidality       bool           `json:"risk_suicidality"`
	RiskRunningAway       bool           `json:"risk_running_away"`
	RiskSelfHarm          bool           `json:"risk_self_harm"`
	RiskWeaponPossession  bool           `json:"risk_weapon_possession"`
	RiskDrugDealing       bool           `json:"risk_drug_dealing"`
	OtherRisks            *string        `json:"other_risks"`
	SharingPermission     bool           `json:"sharing_permission"`
	TruthDeclaration      bool           `json:"truth_declaration"`
	ClientSignature       bool           `json:"client_signature"`
	GuardianSignature     *bool          `json:"guardian_signature"`
	ReferrerSignature     *bool          `json:"referrer_signature"`
	SignatureDate         time.Time      `json:"signature_date"`
	AttachementIds        []uuid.UUID    `json:"attachement_ids"`
	UrgencyScore          string         `json:"urgency_score"`
}

// @Summary Create an intake form
// @Description Create an intake form
// @Tags intake_form
// @Accept json
// @Produce json
// @Param request body CreateIntakeFormRequest true "Intake form request"
// @Success 201 {object} Response[CreateIntakeFormResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /intake_form [post]
// @Security -
func (server *Server) CreateIntakeFormApi(ctx *gin.Context) {
	var req CreateIntakeFormRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	guardianDetailsBytes, err := json.Marshal(req.GuardianDetails)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	tx, err := server.store.ConnPool.Begin(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer tx.Rollback(ctx)

	qtx := server.store.WithTx(tx)

	arg := db.CreateIntakeFormParams{
		FirstName:             req.FirstName,
		LastName:              req.LastName,
		DateOfBirth:           pgtype.Date{Time: req.DateOfBirth, Valid: true},
		Nationality:           req.Nationality,
		Bsn:                   req.Bsn,
		Address:               req.Address,
		City:                  req.City,
		PostalCode:            req.PostalCode,
		PhoneNumber:           req.PhoneNumber,
		Gender:                req.Gender,
		Email:                 req.Email,
		IDType:                req.IDType,
		IDNumber:              req.IDNumber,
		ReferrerName:          req.ReferrerName,
		ReferrerOrganization:  req.ReferrerOrganization,
		ReferrerFunction:      req.ReferrerFunction,
		ReferrerPhone:         req.ReferrerPhone,
		ReferrerEmail:         req.ReferrerEmail,
		SignedBy:              req.SignedBy,
		HasValidIndication:    req.HasValidIndication,
		LawType:               req.LawType,
		OtherLawSpecification: req.OtherLawSpecification,
		MainProviderName:      req.MainProviderName,
		MainProviderContact:   req.MainProviderContact,
		IndicationStartDate:   pgtype.Date{Time: req.IndicationStartDate, Valid: true},
		IndicationEndDate:     pgtype.Date{Time: req.IndicationEndDate, Valid: true},
		RegistrationReason:    req.RegistrationReason,
		GuidanceGoals:         req.GuidanceGoals,
		RegistrationType:      req.RegistrationType,
		LivingSituation:       req.LivingSituation,
		OtherLivingSituation:  req.OtherLivingSituation,
		ParentalAuthority:     req.ParentalAuthority,
		CurrentSchool:         req.CurrentSchool,
		MentorName:            req.MentorName,
		MentorPhone:           req.MentorPhone,
		MentorEmail:           req.MentorEmail,
		PreviousCare:          req.PreviousCare,
		GuardianDetails:       guardianDetailsBytes,
		Diagnoses:             req.Diagnoses,
		UsesMedication:        req.UsesMedication,
		MedicationDetails:     req.MedicationDetails,
		AddictionIssues:       req.AddictionIssues,
		JudicialInvolvement:   req.JudicialInvolvement,
		RiskAggression:        req.RiskAggression,
		RiskSuicidality:       req.RiskSuicidality,
		RiskRunningAway:       req.RiskRunningAway,
		RiskSelfHarm:          req.RiskSelfHarm,
		RiskWeaponPossession:  req.RiskWeaponPossession,
		RiskDrugDealing:       req.RiskDrugDealing,
		OtherRisks:            req.OtherRisks,
		SharingPermission:     req.SharingPermission,
		TruthDeclaration:      req.TruthDeclaration,
		ClientSignature:       req.ClientSignature,
		GuardianSignature:     req.GuardianSignature,
		ReferrerSignature:     req.ReferrerSignature,
		SignatureDate:         pgtype.Date{Time: req.SignatureDate, Valid: true},
		UrgencyScore:          req.UrgencyScore,
		AttachementIds:        req.AttachementIds,
	}

	form, err := qtx.CreateIntakeForm(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	for _, attachmentID := range req.AttachementIds {
		_, err = qtx.SetAttachmentAsUsedorUnused(ctx, db.SetAttachmentAsUsedorUnusedParams{
			Uuid:   attachmentID,
			IsUsed: true,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var GuardianDetails []GuardionInfo
	err = json.Unmarshal(form.GuardianDetails, &GuardianDetails)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(CreateIntakeFormResponse{
		ID:                    form.ID,
		FirstName:             form.FirstName,
		LastName:              form.LastName,
		DateOfBirth:           form.DateOfBirth.Time,
		Nationality:           form.Nationality,
		Bsn:                   form.Bsn,
		Address:               form.Address,
		City:                  form.City,
		PostalCode:            form.PostalCode,
		Gender:                form.Gender,
		Email:                 form.Email,
		IDType:                form.IDType,
		IDNumber:              form.IDNumber,
		ReferrerName:          form.ReferrerName,
		ReferrerOrganization:  form.ReferrerOrganization,
		ReferrerFunction:      form.ReferrerFunction,
		ReferrerPhone:         form.ReferrerPhone,
		ReferrerEmail:         form.ReferrerEmail,
		SignedBy:              form.SignedBy,
		HasValidIndication:    form.HasValidIndication,
		LawType:               form.LawType,
		OtherLawSpecification: form.OtherLawSpecification,
		MainProviderName:      form.MainProviderName,
		MainProviderContact:   form.MainProviderContact,
		IndicationStartDate:   form.IndicationStartDate.Time,
		IndicationEndDate:     form.IndicationEndDate.Time,
		RegistrationReason:    form.RegistrationReason,
		GuidanceGoals:         form.GuidanceGoals,
		RegistrationType:      form.RegistrationType,
		LivingSituation:       form.LivingSituation,
		OtherLivingSituation:  form.OtherLivingSituation,
		ParentalAuthority:     form.ParentalAuthority,
		CurrentSchool:         form.CurrentSchool,
		MentorName:            form.MentorName,
		MentorPhone:           form.MentorPhone,
		MentorEmail:           form.MentorEmail,
		PreviousCare:          form.PreviousCare,
		GuardianDetails:       GuardianDetails,
		Diagnoses:             form.Diagnoses,
		UsesMedication:        form.UsesMedication,
		MedicationDetails:     form.MedicationDetails,
		AddictionIssues:       form.AddictionIssues,
		JudicialInvolvement:   form.JudicialInvolvement,
		RiskAggression:        form.RiskAggression,
		RiskSuicidality:       form.RiskSuicidality,
		RiskRunningAway:       form.RiskRunningAway,
		RiskSelfHarm:          form.RiskSelfHarm,
		RiskWeaponPossession:  form.RiskWeaponPossession,
		RiskDrugDealing:       form.RiskDrugDealing,
		OtherRisks:            form.OtherRisks,
		SharingPermission:     form.SharingPermission,
		TruthDeclaration:      form.TruthDeclaration,
		ClientSignature:       form.ClientSignature,
		GuardianSignature:     form.GuardianSignature,
		ReferrerSignature:     form.ReferrerSignature,
		SignatureDate:         form.SignatureDate.Time,
		AttachementIds:        form.AttachementIds,
		UrgencyScore:          form.UrgencyScore,
	}, "Intake form created successfully")

	ctx.JSON(http.StatusCreated, res)
}

// ListIntakeFormsRequest represents a request to list intake forms
type ListIntakeFormsRequest struct {
	pagination.Request
	Search    string `form:"search"`
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=created_at urgency_score"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// ListIntakeFormsResponse represents a response from the list intake forms handler
type ListIntakeFormsResponse struct {
	ID                    int64          `json:"id"`
	FirstName             string         `json:"first_name"`
	LastName              string         `json:"last_name"`
	DateOfBirth           time.Time      `json:"date_of_birth"`
	Nationality           string         `json:"nationality"`
	Bsn                   string         `json:"bsn"`
	Address               string         `json:"address"`
	City                  string         `json:"city"`
	PostalCode            string         `json:"postal_code"`
	PhoneNumber           string         `json:"phone_number"`
	Gender                string         `json:"gender"`
	Email                 string         `json:"email"`
	IDType                string         `json:"id_type"`
	IDNumber              string         `json:"id_number"`
	ReferrerName          *string        `json:"referrer_name"`
	ReferrerOrganization  *string        `json:"referrer_organization"`
	ReferrerFunction      *string        `json:"referrer_function"`
	ReferrerPhone         *string        `json:"referrer_phone"`
	ReferrerEmail         *string        `json:"referrer_email"`
	SignedBy              *string        `json:"signed_by"`
	HasValidIndication    bool           `json:"has_valid_indication"`
	LawType               *string        `json:"law_type"`
	OtherLawSpecification *string        `json:"other_law_specification"`
	MainProviderName      *string        `json:"main_provider_name"`
	MainProviderContact   *string        `json:"main_provider_contact"`
	IndicationStartDate   time.Time      `json:"indication_start_date"`
	IndicationEndDate     time.Time      `json:"indication_end_date"`
	RegistrationReason    *string        `json:"registration_reason"`
	GuidanceGoals         *string        `json:"guidance_goals"`
	RegistrationType      *string        `json:"registration_type"`
	LivingSituation       *string        `json:"living_situation"`
	OtherLivingSituation  *string        `json:"other_living_situation"`
	ParentalAuthority     bool           `json:"parental_authority"`
	CurrentSchool         *string        `json:"current_school"`
	MentorName            *string        `json:"mentor_name"`
	MentorPhone           *string        `json:"mentor_phone"`
	MentorEmail           *string        `json:"mentor_email"`
	PreviousCare          *string        `json:"previous_care"`
	GuardianDetails       []GuardionInfo `json:"guardian_details"`
	Diagnoses             *string        `json:"diagnoses"`
	UsesMedication        bool           `json:"uses_medication"`
	MedicationDetails     *string        `json:"medication_details"`
	AddictionIssues       bool           `json:"addiction_issues"`
	JudicialInvolvement   bool           `json:"judicial_involvement"`
	RiskAggression        bool           `json:"risk_aggression"`
	RiskSuicidality       bool           `json:"risk_suicidality"`
	RiskRunningAway       bool           `json:"risk_running_away"`
	RiskSelfHarm          bool           `json:"risk_self_harm"`
	RiskWeaponPossession  bool           `json:"risk_weapon_possession"`
	RiskDrugDealing       bool           `json:"risk_drug_dealing"`
	OtherRisks            *string        `json:"other_risks"`
	SharingPermission     bool           `json:"sharing_permission"`
	TruthDeclaration      bool           `json:"truth_declaration"`
	ClientSignature       bool           `json:"client_signature"`
	GuardianSignature     *bool          `json:"guardian_signature"`
	ReferrerSignature     *bool          `json:"referrer_signature"`
	SignatureDate         time.Time      `json:"signature_date"`
	AttachementIds        []uuid.UUID    `json:"attachement_ids"`
	TimeSinceSubmission   string         `json:"time_since_submission"`
	UrgencyScore          string         `json:"urgency_score"`
}

// @Summary List intake forms
// @Description List intake forms
// @Tags intake_form
// @Accept json
// @Produce json
// @Param search query string false "Search query"
// @Param page query integer false "Page"
// @Param page_size query integer false "Page size"
// @Param sort_by query string false "Sort by (options: created_at, urgency_score)"
// @Param sort_order query string false "Sort order (options: asc, desc)"
// @Success 200 {object} Response[pagination.Response[ListIntakeFormsResponse]]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /intake_form [get]
func (server *Server) ListIntakeFormsApi(ctx *gin.Context) {
	var req ListIntakeFormsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	param := req.GetParams()

	forms, err := server.store.ListIntakeForms(ctx, db.ListIntakeFormsParams{
		Search:    req.Search,
		Limit:     param.Limit,
		Offset:    param.Offset,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	formList := make([]ListIntakeFormsResponse, len(forms))
	for i, form := range forms {
		var GuardianDetails []GuardionInfo
		err = json.Unmarshal(form.GuardianDetails, &GuardianDetails)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		timeSinceSubmission := time.Since(form.CreatedAt.Time)
		daysElapsed := timeSinceSubmission.Hours() / 24

		formList[i] = ListIntakeFormsResponse{
			ID:                    form.ID,
			FirstName:             form.FirstName,
			LastName:              form.LastName,
			DateOfBirth:           form.DateOfBirth.Time,
			Nationality:           form.Nationality,
			Bsn:                   form.Bsn,
			Address:               form.Address,
			City:                  form.City,
			PostalCode:            form.PostalCode,
			PhoneNumber:           form.PhoneNumber,
			Gender:                form.Gender,
			Email:                 form.Email,
			IDType:                form.IDType,
			IDNumber:              form.IDNumber,
			ReferrerName:          form.ReferrerName,
			ReferrerOrganization:  form.ReferrerOrganization,
			ReferrerFunction:      form.ReferrerFunction,
			ReferrerPhone:         form.ReferrerPhone,
			ReferrerEmail:         form.ReferrerEmail,
			SignedBy:              form.SignedBy,
			HasValidIndication:    form.HasValidIndication,
			LawType:               form.LawType,
			OtherLawSpecification: form.OtherLawSpecification,
			MainProviderName:      form.MainProviderName,
			MainProviderContact:   form.MainProviderContact,
			IndicationStartDate:   form.IndicationStartDate.Time,
			IndicationEndDate:     form.IndicationEndDate.Time,
			RegistrationReason:    form.RegistrationReason,
			GuidanceGoals:         form.GuidanceGoals,
			RegistrationType:      form.RegistrationType,
			LivingSituation:       form.LivingSituation,
			OtherLivingSituation:  form.OtherLivingSituation,
			ParentalAuthority:     form.ParentalAuthority,
			CurrentSchool:         form.CurrentSchool,
			MentorName:            form.MentorName,
			MentorPhone:           form.MentorPhone,
			MentorEmail:           form.MentorEmail,
			PreviousCare:          form.PreviousCare,
			GuardianDetails:       GuardianDetails,
			Diagnoses:             form.Diagnoses,
			UsesMedication:        form.UsesMedication,
			MedicationDetails:     form.MedicationDetails,
			AddictionIssues:       form.AddictionIssues,
			JudicialInvolvement:   form.JudicialInvolvement,
			RiskAggression:        form.RiskAggression,
			RiskSuicidality:       form.RiskSuicidality,
			RiskRunningAway:       form.RiskRunningAway,
			RiskSelfHarm:          form.RiskSelfHarm,
			RiskWeaponPossession:  form.RiskWeaponPossession,
			RiskDrugDealing:       form.RiskDrugDealing,
			OtherRisks:            form.OtherRisks,
			SharingPermission:     form.SharingPermission,
			TruthDeclaration:      form.TruthDeclaration,
			ClientSignature:       form.ClientSignature,
			GuardianSignature:     form.GuardianSignature,
			ReferrerSignature:     form.ReferrerSignature,
			SignatureDate:         form.SignatureDate.Time,
			AttachementIds:        form.AttachementIds,
			TimeSinceSubmission:   fmt.Sprintf("%d days", int(daysElapsed)),
			UrgencyScore:          form.UrgencyScore,
		}
	}
	if len(forms) == 0 {
		ctx.JSON(http.StatusOK, SuccessResponse([]ListIntakeFormsResponse{}, "No intake forms found"))
		return
	}

	totalCount := forms[0].TotalCount
	pag := pagination.NewResponse(ctx, req.Request, formList, totalCount)
	res := SuccessResponse(pag, "Intake forms retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// GetIntakeFormResponse represents a response from the get intake form handler
type GetIntakeFormResponse struct {
	ID                    int64          `json:"id"`
	FirstName             string         `json:"first_name"`
	LastName              string         `json:"last_name"`
	DateOfBirth           time.Time      `json:"date_of_birth"`
	Nationality           string         `json:"nationality"`
	Bsn                   string         `json:"bsn"`
	Address               string         `json:"address"`
	City                  string         `json:"city"`
	PostalCode            string         `json:"postal_code"`
	PhoneNumber           string         `json:"phone_number"`
	Gender                string         `json:"gender"`
	Email                 string         `json:"email"`
	IDType                string         `json:"id_type"`
	IDNumber              string         `json:"id_number"`
	ReferrerName          *string        `json:"referrer_name"`
	ReferrerOrganization  *string        `json:"referrer_organization"`
	ReferrerFunction      *string        `json:"referrer_function"`
	ReferrerPhone         *string        `json:"referrer_phone"`
	ReferrerEmail         *string        `json:"referrer_email"`
	SignedBy              *string        `json:"signed_by"`
	HasValidIndication    bool           `json:"has_valid_indication"`
	LawType               *string        `json:"law_type"`
	OtherLawSpecification *string        `json:"other_law_specification"`
	MainProviderName      *string        `json:"main_provider_name"`
	MainProviderContact   *string        `json:"main_provider_contact"`
	IndicationStartDate   time.Time      `json:"indication_start_date"`
	IndicationEndDate     time.Time      `json:"indication_end_date"`
	RegistrationReason    *string        `json:"registration_reason"`
	GuidanceGoals         *string        `json:"guidance_goals"`
	RegistrationType      *string        `json:"registration_type"`
	LivingSituation       *string        `json:"living_situation"`
	OtherLivingSituation  *string        `json:"other_living_situation"`
	ParentalAuthority     bool           `json:"parental_authority"`
	CurrentSchool         *string        `json:"current_school"`
	MentorName            *string        `json:"mentor_name"`
	MentorPhone           *string        `json:"mentor_phone"`
	MentorEmail           *string        `json:"mentor_email"`
	PreviousCare          *string        `json:"previous_care"`
	GuardianDetails       []GuardionInfo `json:"guardian_details"`
	Diagnoses             *string        `json:"diagnoses"`
	UsesMedication        bool           `json:"uses_medication"`
	MedicationDetails     *string        `json:"medication_details"`
	AddictionIssues       bool           `json:"addiction_issues"`
	JudicialInvolvement   bool           `json:"judicial_involvement"`
	RiskAggression        bool           `json:"risk_aggression"`
	RiskSuicidality       bool           `json:"risk_suicidality"`
	RiskRunningAway       bool           `json:"risk_running_away"`
	RiskSelfHarm          bool           `json:"risk_self_harm"`
	RiskWeaponPossession  bool           `json:"risk_weapon_possession"`
	RiskDrugDealing       bool           `json:"risk_drug_dealing"`
	OtherRisks            *string        `json:"other_risks"`
	SharingPermission     bool           `json:"sharing_permission"`
	TruthDeclaration      bool           `json:"truth_declaration"`
	ClientSignature       bool           `json:"client_signature"`
	GuardianSignature     *bool          `json:"guardian_signature"`
	ReferrerSignature     *bool          `json:"referrer_signature"`
	SignatureDate         time.Time      `json:"signature_date"`
	AttachementIds        []uuid.UUID    `json:"attachement_ids"`
	TimeSinceSubmission   string         `json:"time_since_submission"`
	UrgencyScore          string         `json:"urgency_score"`
}

// @Summary Get an intake form
// @Description Get an intake form
// @Tags intake_form
// @Accept json
// @Produce json
// @Param id path string true "Intake form ID"
// @Success 200 {object} Response[GetIntakeFormResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /intake_form/{id} [get]
func (server *Server) GetIntakeFormApi(ctx *gin.Context) {
	id := ctx.Param("id")
	formID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	form, err := server.store.GetIntakeForm(ctx, formID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var GuardianDetails []GuardionInfo
	err = json.Unmarshal(form.GuardianDetails, &GuardianDetails)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	timeSinceSubmission := time.Since(form.CreatedAt.Time)
	daysElapsed := timeSinceSubmission.Hours() / 24

	res := SuccessResponse(GetIntakeFormResponse{
		ID:                    form.ID,
		FirstName:             form.FirstName,
		LastName:              form.LastName,
		DateOfBirth:           form.DateOfBirth.Time,
		Nationality:           form.Nationality,
		Bsn:                   form.Bsn,
		Address:               form.Address,
		City:                  form.City,
		PostalCode:            form.PostalCode,
		PhoneNumber:           form.PhoneNumber,
		Gender:                form.Gender,
		Email:                 form.Email,
		IDType:                form.IDType,
		IDNumber:              form.IDNumber,
		ReferrerName:          form.ReferrerName,
		ReferrerOrganization:  form.ReferrerOrganization,
		ReferrerFunction:      form.ReferrerFunction,
		ReferrerPhone:         form.ReferrerPhone,
		ReferrerEmail:         form.ReferrerEmail,
		SignedBy:              form.SignedBy,
		HasValidIndication:    form.HasValidIndication,
		LawType:               form.LawType,
		OtherLawSpecification: form.OtherLawSpecification,
		MainProviderName:      form.MainProviderName,
		MainProviderContact:   form.MainProviderContact,
		IndicationStartDate:   form.IndicationStartDate.Time,
		IndicationEndDate:     form.IndicationEndDate.Time,
		RegistrationReason:    form.RegistrationReason,
		GuidanceGoals:         form.GuidanceGoals,
		RegistrationType:      form.RegistrationType,
		LivingSituation:       form.LivingSituation,
		OtherLivingSituation:  form.OtherLivingSituation,
		ParentalAuthority:     form.ParentalAuthority,
		CurrentSchool:         form.CurrentSchool,
		MentorName:            form.MentorName,
		MentorPhone:           form.MentorPhone,
		MentorEmail:           form.MentorEmail,
		PreviousCare:          form.PreviousCare,
		GuardianDetails:       GuardianDetails,
		Diagnoses:             form.Diagnoses,
		UsesMedication:        form.UsesMedication,
		MedicationDetails:     form.MedicationDetails,
		AddictionIssues:       form.AddictionIssues,
		JudicialInvolvement:   form.JudicialInvolvement,
		RiskAggression:        form.RiskAggression,
		RiskSuicidality:       form.RiskSuicidality,
		RiskRunningAway:       form.RiskRunningAway,
		RiskSelfHarm:          form.RiskSelfHarm,
		RiskWeaponPossession:  form.RiskWeaponPossession,
		RiskDrugDealing:       form.RiskDrugDealing,
		OtherRisks:            form.OtherRisks,
		SharingPermission:     form.SharingPermission,
		TruthDeclaration:      form.TruthDeclaration,
		ClientSignature:       form.ClientSignature,
		GuardianSignature:     form.GuardianSignature,
		ReferrerSignature:     form.ReferrerSignature,
		SignatureDate:         form.SignatureDate.Time,
		AttachementIds:        form.AttachementIds,
		TimeSinceSubmission:   fmt.Sprintf("%d days", int(daysElapsed)),
		UrgencyScore:          form.UrgencyScore,
	}, "Intake form retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// AddUrgencyScoreRequest represents a request to add urgency score to an intake form
type AddUrgencyScoreRequest struct {
	UrgencyScore string `json:"urgency_score"`
}

// AddUrgencyScoreResponse represents a response from the add urgency score handler
type AddUrgencyScoreResponse struct {
	ID           int64  `json:"id"`
	UrgencyScore string `json:"urgency_score"`
}

// @Summary Add urgency score to an intake form
// @Description Add urgency score to an intake form
// @Tags intake_form
// @Accept json
// @Produce json
// @Param id path string true "Intake form ID"
// @Param request body AddUrgencyScoreRequest true "Urgency score request"
// @Success 200 {object} Response[AddUrgencyScoreResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /intake_form/{id}/urgency_score [post]
func (server *Server) AddUrgencyScoreApi(ctx *gin.Context) {
	id := ctx.Param("id")
	formID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req AddUrgencyScoreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	form, err := server.store.AddUrgencyScore(ctx, db.AddUrgencyScoreParams{
		ID:           formID,
		UrgencyScore: req.UrgencyScore,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(AddUrgencyScoreResponse{
		ID:           form.ID,
		UrgencyScore: form.UrgencyScore,
	}, "Urgency score added successfully")
	ctx.JSON(http.StatusOK, res)
}

// MoveToWaitingListResponse represents a response from the move to waiting list handler
type MoveToWaitingListResponse struct {
	ClientID int64 `json:"client_id"`
}

// @Summary Move an intake form to waiting list
// @Description Move an intake form to waiting list
// @Tags intake_form
// @Accept json
// @Produce json
// @Param id path string true "Intake form ID"
// @Success 200 {object} Response[MoveToWaitingListResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /intake_form/{id}/move_to_waiting_list [post]
// func (server *Server) MoveToWaitingList(ctx *gin.Context) {
// 	id := ctx.Param("id")
// 	formID, err := strconv.ParseInt(id, 10, 64)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	form, err := server.store.MoveToWaitingList(ctx, formID)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	address := []Address{
// 		{
// 			Address: &form.Address,
// 			City:    &form.City,
// 			ZipCode: &form.PostalCode,
// 		},
// 	}

// 	AddressesJSON, err := json.Marshal(address)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	arg := db.CreateClientDetailsParams{
// 		IntakeFormID:          &form.ID,
// 		FirstName:             form.FirstName,
// 		LastName:              form.LastName,
// 		DateOfBirth:           form.DateOfBirth,
// 		Bsn:                   &form.Bsn,
// 		Birthplace:            &form.Nationality,
// 		Email:                 form.Email,
// 		PhoneNumber:           &form.PhoneNumber,
// 		Gender:                form.Gender,
// 		Addresses:             AddressesJSON,
// 	}
// 	client, err := server.store.CreateClientDetails(ctx, arg)

// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	res := SuccessResponse(MoveToWaitingListResponse{
// 		ClientID: client.ID,
// 	}, "Client moved to waiting list successfully")
// 	ctx.JSON(http.StatusOK, res)
// }
