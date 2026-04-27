package handler

import (
	"encoding/json"
	"fmt"
	"time"

	"maicare_go/internal/domain"
	"maicare_go/internal/httpapi"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ==================== Request DTOs ====================

type createContractTypeRequest struct {
	Name string `json:"name" binding:"required"`
}

type createContractRequest struct {
	ClientID        uuid.UUID   `json:"client_id" binding:"required"`
	TypeID          *uuid.UUID  `json:"type_id"`
	StartDate       time.Time   `json:"start_date" binding:"required"`
	EndDate         time.Time   `json:"end_date" binding:"required"`
	ReminderPeriod  *int32      `json:"reminder_period" binding:"omitempty,gte=0"`
	Vat             *int32      `json:"VAT"`
	Price           float64     `json:"price"`
	PriceTimeUnit   string      `json:"price_time_unit" binding:"required,oneof=minute hourly daily weekly"`
	Hours           *float64    `json:"hours"`
	HoursType       *string     `json:"hours_type" binding:"omitempty,oneof=weekly all_period"`
	CareName        string      `json:"care_name"`
	CareType        string      `json:"care_type" binding:"required,oneof=ambulante accommodation"`
	SenderID        uuid.UUID   `json:"sender_id" binding:"required"`
	AttachmentIds   []uuid.UUID `json:"attachment_ids"`
	FinancingAct    string      `json:"financing_act" binding:"required,oneof=WMO ZVW WLZ JW WPG"`
	FinancingOption string      `json:"financing_option" binding:"required,oneof=ZIN PGB"`
}

type updateContractRequest struct {
	TypeID          *uuid.UUID  `json:"type_id"`
	StartDate       *time.Time  `json:"start_date"`
	EndDate         *time.Time  `json:"end_date"`
	ReminderPeriod  *int32      `json:"reminder_period" binding:"omitempty,gte=0"`
	Vat             *int32      `json:"VAT"`
	Price           *float64    `json:"price"`
	PriceTimeUnit   *string     `json:"price_time_unit" binding:"omitempty,oneof=minute hourly daily weekly"`
	Hours           *float64    `json:"hours"`
	HoursType       *string     `json:"hours_type" binding:"omitempty,oneof=weekly all_period"`
	CareName        *string     `json:"care_name"`
	SenderID        *uuid.UUID  `json:"sender_id"`
	AttachmentIds   []uuid.UUID `json:"attachment_ids"`
	FinancingAct    *string     `json:"financing_act" binding:"omitempty,oneof=WMO ZVW WLZ JW WPG"`
	FinancingOption *string     `json:"financing_option" binding:"omitempty,oneof=ZIN PGB"`
}

type updateContractStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=approved draft terminated stopped expired"`
}

type listContractsRequest struct {
	httpapi.PageRequest
	Search          *string    `form:"search" binding:"omitempty"`
	Status          []string   `form:"status" binding:"omitempty,dive,oneof=approved draft terminated stopped expired"`
	CareType        []string   `form:"care_type" binding:"omitempty,dive,oneof=ambulante accommodation"`
	FinancingAct    []string   `form:"financing_act" binding:"omitempty,dive,oneof=WMO ZVW WLZ JW WPG"`
	FinancingOption *string    `form:"financing_option" binding:"omitempty,oneof=ZIN PGB"`
	EndDateFrom     *time.Time `form:"end_date_from" time_format:"2006-01-02"`
	EndDateTo       *time.Time `form:"end_date_to" time_format:"2006-01-02"`
}

// ==================== Response DTOs ====================

type contractTypeResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type attachmentDetailResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	DownloadURL string    `json:"download_url"`
}

type contractResponse struct {
	ID              uuid.UUID                  `json:"id"`
	TypeID          *uuid.UUID                 `json:"type_id"`
	Status          string                     `json:"status"`
	StartDate       time.Time                  `json:"start_date"`
	EndDate         time.Time                  `json:"end_date"`
	ReminderPeriod  int32                      `json:"reminder_period"`
	Vat             *int32                     `json:"VAT"`
	Price           float64                    `json:"price"`
	PriceTimeUnit   string                     `json:"price_time_unit"`
	Hours           *float64                   `json:"hours"`
	HoursType       *string                    `json:"hours_type"`
	CareName        string                     `json:"care_name"`
	CareType        string                     `json:"care_type"`
	ClientID        uuid.UUID                  `json:"client_id"`
	SenderID        uuid.UUID                  `json:"sender_id"`
	AttachmentIds   []uuid.UUID                `json:"attachment_ids"`
	Attachments     []attachmentDetailResponse `json:"attachments"`
	FinancingAct    string                     `json:"financing_act"`
	FinancingOption string                     `json:"financing_option"`
	DepartureReason *string                    `json:"departure_reason"`
	DepartureReport *string                    `json:"departure_report"`
	UpdatedAt       time.Time                  `json:"updated_at"`
	CreatedAt       time.Time                  `json:"created_at"`
}

type updateContractResponse struct {
	ID              uuid.UUID   `json:"id"`
	TypeID          *uuid.UUID  `json:"type_id"`
	Status          string      `json:"status"`
	StartDate       time.Time   `json:"start_date"`
	EndDate         time.Time   `json:"end_date"`
	ReminderPeriod  int32       `json:"reminder_period"`
	Vat             *int32      `json:"VAT"`
	Price           float64     `json:"price"`
	PriceFrequency  string      `json:"price_frequency"`
	Hours           *float64    `json:"hours"`
	HoursType       *string     `json:"hours_type"`
	CareName        string      `json:"care_name"`
	CareType        string      `json:"care_type"`
	ClientID        uuid.UUID   `json:"client_id"`
	SenderID        uuid.UUID   `json:"sender_id"`
	AttachmentIds   []uuid.UUID `json:"attachment_ids"`
	FinancingAct    string      `json:"financing_act"`
	FinancingOption string      `json:"financing_option"`
	DepartureReason *string     `json:"departure_reason"`
	DepartureReport *string     `json:"departure_report"`
	UpdatedAt       time.Time   `json:"updated_at"`
	CreatedAt       time.Time   `json:"created_at"`
}

type updateContractStatusResponse struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

type contractClientListItemResponse struct {
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	DaysLeft        int32     `json:"days_left"`
	CareName        string    `json:"care_name"`
	CareType        string    `json:"care_type"`
	FinancingAct    string    `json:"financing_act"`
	FinancingOption string    `json:"financing_option"`
}

type contractListItemResponse struct {
	ID               uuid.UUID  `json:"id"`
	ClientID         uuid.UUID  `json:"client_id"`
	ClientFirstName  string     `json:"client_first_name"`
	ClientLastName   string     `json:"client_last_name"`
	ClientFilenumber string     `json:"client_filenumber"`
	SenderID         uuid.UUID  `json:"sender_id"`
	SenderName       string     `json:"sender_name"`
	CareName         string     `json:"care_name"`
	CareType         string     `json:"care_type"`
	Price            float64    `json:"price"`
	PriceTimeUnit    string     `json:"price_time_unit"`
	Hours            *float64   `json:"hours"`
	HoursType        *string    `json:"hours_type"`
	FinancingAct     string     `json:"financing_act"`
	FinancingOption  string     `json:"financing_option"`
	StartDate        time.Time  `json:"start_date"`
	EndDate          time.Time  `json:"end_date"`
	DaysLeft         int32      `json:"days_left"`
	Status           string     `json:"status"`
	ApprovedAt       *time.Time `json:"approved_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type contractDetailResponse struct {
	// Contract fields
	ID              uuid.UUID   `json:"id"`
	TypeID          *uuid.UUID  `json:"type_id"`
	TypeName        *string     `json:"type_name"`
	Status          string      `json:"status"`
	ApprovedAt      *time.Time  `json:"approved_at"`
	StartDate       time.Time   `json:"start_date"`
	EndDate         time.Time   `json:"end_date"`
	ReminderPeriod  int32       `json:"reminder_period"`
	Vat             *int32      `json:"VAT"`
	Price           float64     `json:"price"`
	PriceTimeUnit   string      `json:"price_time_unit"`
	Hours           *float64    `json:"hours"`
	HoursType       *string     `json:"hours_type"`
	CareName        string      `json:"care_name"`
	CareType        string      `json:"care_type"`
	AttachmentIds   []uuid.UUID `json:"attachment_ids"`
	FinancingAct    string      `json:"financing_act"`
	FinancingOption string      `json:"financing_option"`
	DepartureReason *string     `json:"departure_reason"`
	DepartureReport *string     `json:"departure_report"`
	UpdatedAt       time.Time   `json:"updated_at"`
	CreatedAt       time.Time   `json:"created_at"`

	// Client fields
	ClientID         uuid.UUID `json:"client_id"`
	ClientFirstName  string    `json:"client_first_name"`
	ClientLastName   string    `json:"client_last_name"`
	ClientFilenumber string    `json:"client_filenumber"`
	ClientBsn        *string   `json:"client_bsn"`

	// Sender fields
	SenderID                  uuid.UUID `json:"sender_id"`
	SenderName                string    `json:"sender_name"`
	SenderType                string    `json:"sender_type"`
	SenderStreet              *string   `json:"sender_street"`
	SenderHouseNumber         *string   `json:"sender_house_number"`
	SenderHouseNumberAddition *string   `json:"sender_house_number_addition"`
	SenderPostalCode          *string   `json:"sender_postal_code"`
	SenderCity                *string   `json:"sender_city"`
	SenderLand                *string   `json:"sender_land"`
	SenderKvknumber           *string   `json:"sender_kvknumber"`
	SenderBtwnumber           *string   `json:"sender_btwnumber"`
	SenderPhoneNumber         *string   `json:"sender_phone_number"`
	SenderClientNumber        *string   `json:"sender_client_number"`
	SenderEmailAddress        *string   `json:"sender_email_address"`
}

type contractAuditLogResponse struct {
	AuditID            uuid.UUID       `json:"audit_id"`
	ContractID         uuid.UUID       `json:"contract_id"`
	Operation          string          `json:"operation"`
	ChangedBy          *uuid.UUID      `json:"changed_by"`
	ChangedAt          time.Time       `json:"changed_at"`
	OldValues          map[string]interface{} `json:"old_values"`
	NewValues          map[string]interface{} `json:"new_values"`
	ChangedFields      []string        `json:"changed_fields"`
	ChangedByFirstName *string         `json:"changed_by_first_name"`
	ChangedByLastName  *string         `json:"changed_by_last_name"`
}

// ==================== Mappers ====================

func toCreateContractTypeResult(res domain.CreateContractTypeResult) contractTypeResponse {
	return contractTypeResponse{ID: res.ID, Name: res.Name}
}

func toContractResponse(res domain.CreateContractResult) contractResponse {
	return contractResponse{
		ID:              res.ID,
		TypeID:          res.TypeID,
		Status:          res.Status,
		StartDate:       res.StartDate,
		EndDate:         res.EndDate,
		ReminderPeriod:  res.ReminderPeriod,
		Vat:             res.Vat,
		Price:           res.Price,
		PriceTimeUnit:   res.PriceTimeUnit,
		Hours:           res.Hours,
		HoursType:       res.HoursType,
		CareName:        res.CareName,
		CareType:        res.CareType,
		ClientID:        res.ClientID,
		SenderID:        res.SenderID,
		AttachmentIds:   res.AttachmentIds,
		Attachments:     toAttachmentDetailResponses(res.Attachments),
		FinancingAct:    res.FinancingAct,
		FinancingOption: res.FinancingOption,
		DepartureReason: res.DepartureReason,
		DepartureReport: res.DepartureReport,
		UpdatedAt:       res.UpdatedAt,
		CreatedAt:       res.CreatedAt,
	}
}

func toUpdateContractResponse(res domain.UpdateContractResult) updateContractResponse {
	return updateContractResponse{
		ID:              res.ID,
		TypeID:          res.TypeID,
		Status:          res.Status,
		StartDate:       res.StartDate,
		EndDate:         res.EndDate,
		ReminderPeriod:  res.ReminderPeriod,
		Vat:             res.Vat,
		Price:           res.Price,
		PriceFrequency:  res.PriceTimeUnit,
		Hours:           res.Hours,
		HoursType:       res.HoursType,
		CareName:        res.CareName,
		CareType:        res.CareType,
		ClientID:        res.ClientID,
		SenderID:        res.SenderID,
		AttachmentIds:   res.AttachmentIds,
		FinancingAct:    res.FinancingAct,
		FinancingOption: res.FinancingOption,
		DepartureReason: res.DepartureReason,
		DepartureReport: res.DepartureReport,
		UpdatedAt:       res.UpdatedAt,
		CreatedAt:       res.CreatedAt,
	}
}

func toUpdateContractStatusResponse(res domain.UpdateContractStatusResult) updateContractStatusResponse {
	return updateContractStatusResponse{ID: res.ID, Status: res.Status}
}

func toContractDetailResponse(d *domain.ContractDetail) contractDetailResponse {
	return contractDetailResponse{
		ID:                        d.ID,
		TypeID:                    d.TypeID,
		TypeName:                  d.TypeName,
		Status:                    d.Status,
		ApprovedAt:                d.ApprovedAt,
		StartDate:                 d.StartDate,
		EndDate:                   d.EndDate,
		ReminderPeriod:            d.ReminderPeriod,
		Vat:                       d.Vat,
		Price:                     d.Price,
		PriceTimeUnit:             d.PriceTimeUnit,
		Hours:                     d.Hours,
		HoursType:                 d.HoursType,
		CareName:                  d.CareName,
		CareType:                  d.CareType,
		AttachmentIds:             d.AttachmentIds,
		FinancingAct:              d.FinancingAct,
		FinancingOption:           d.FinancingOption,
		DepartureReason:           d.DepartureReason,
		DepartureReport:           d.DepartureReport,
		UpdatedAt:                 d.UpdatedAt,
		CreatedAt:                 d.CreatedAt,
		ClientID:                  d.ClientID,
		ClientFirstName:           d.ClientFirstName,
		ClientLastName:            d.ClientLastName,
		ClientFilenumber:          d.ClientFilenumber,
		ClientBsn:                 d.ClientBsn,
		SenderID:                  d.SenderID,
		SenderName:                d.SenderName,
		SenderType:                d.SenderType,
		SenderStreet:              d.SenderStreet,
		SenderHouseNumber:         d.SenderHouseNumber,
		SenderHouseNumberAddition: d.SenderHouseNumberAddition,
		SenderPostalCode:          d.SenderPostalCode,
		SenderCity:                d.SenderCity,
		SenderLand:                d.SenderLand,
		SenderKvknumber:           d.SenderKvknumber,
		SenderBtwnumber:           d.SenderBtwnumber,
		SenderPhoneNumber:         d.SenderPhoneNumber,
		SenderClientNumber:        d.SenderClientNumber,
		SenderEmailAddress:        d.SenderEmailAddress,
	}
}

func toContractClientListItemResponses(items []domain.ClientContractListItem) []contractClientListItemResponse {
	result := make([]contractClientListItemResponse, len(items))
	for i, item := range items {
		result[i] = contractClientListItemResponse{
			StartDate:       item.StartDate,
			EndDate:         item.EndDate,
			DaysLeft:        item.DaysLeft,
			CareName:        item.CareName,
			CareType:        item.CareType,
			FinancingAct:    item.FinancingAct,
			FinancingOption: item.FinancingOption,
		}
	}
	return result
}

func toContractListItemResponses(items []domain.ContractListItem) []contractListItemResponse {
	result := make([]contractListItemResponse, len(items))
	for i, item := range items {
		result[i] = contractListItemResponse{
			ID:               item.ID,
			ClientID:         item.ClientID,
			ClientFirstName:  item.ClientFirstName,
			ClientLastName:   item.ClientLastName,
			ClientFilenumber: item.ClientFilenumber,
			SenderID:         item.SenderID,
			SenderName:       item.SenderName,
			CareName:         item.CareName,
			CareType:         item.CareType,
			Price:            item.Price,
			PriceTimeUnit:    item.PriceTimeUnit,
			Hours:            item.Hours,
			HoursType:        item.HoursType,
			FinancingAct:     item.FinancingAct,
			FinancingOption:  item.FinancingOption,
			StartDate:        item.StartDate,
			EndDate:          item.EndDate,
			DaysLeft:         item.DaysLeft,
			Status:           item.Status,
			ApprovedAt:       item.ApprovedAt,
			UpdatedAt:        item.UpdatedAt,
		}
	}
	return result
}

func toAttachmentDetailResponses(items []domain.ContractAttachment) []attachmentDetailResponse {
	result := make([]attachmentDetailResponse, len(items))
	for i, item := range items {
		result[i] = attachmentDetailResponse{
			ID:          item.ID,
			Name:        item.Name,
			Size:        item.Size,
			DownloadURL: item.DownloadURL,
		}
	}
	return result
}

func toContractAuditLogResponses(items []domain.ContractAuditLog) []contractAuditLogResponse {
	result := make([]contractAuditLogResponse, len(items))
	for i, item := range items {
		result[i] = contractAuditLogResponse{
			AuditID:            item.AuditID,
			ContractID:         item.ContractID,
			Operation:          item.Operation,
			ChangedBy:          item.ChangedBy,
			ChangedAt:          item.ChangedAt,
			OldValues:          parseJSONToObject(item.OldValues),
			NewValues:          parseJSONToObject(item.NewValues),
			ChangedFields:      item.ChangedFields,
			ChangedByFirstName: item.ChangedByFirstName,
			ChangedByLastName:  item.ChangedByLastName,
		}
	}
	return result
}

// ==================== Request Mappers ====================

func toCreateContractParams(req createContractRequest) domain.CreateContractParams {
	return domain.CreateContractParams{
		ClientID:        req.ClientID,
		TypeID:          req.TypeID,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		ReminderPeriod:  req.ReminderPeriod,
		Vat:             req.Vat,
		Price:           req.Price,
		PriceTimeUnit:   req.PriceTimeUnit,
		Hours:           req.Hours,
		HoursType:       req.HoursType,
		CareName:        req.CareName,
		CareType:        req.CareType,
		SenderID:        req.SenderID,
		AttachmentIds:   req.AttachmentIds,
		FinancingAct:    req.FinancingAct,
		FinancingOption: req.FinancingOption,
	}
}

func toUpdateContractParams(contractID uuid.UUID, req updateContractRequest) domain.UpdateContractParams {
	return domain.UpdateContractParams{
		ContractID:      contractID,
		TypeID:          req.TypeID,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		ReminderPeriod:  req.ReminderPeriod,
		Vat:             req.Vat,
		Price:           req.Price,
		PriceTimeUnit:   req.PriceTimeUnit,
		Hours:           req.Hours,
		HoursType:       req.HoursType,
		CareName:        req.CareName,
		SenderID:        req.SenderID,
		AttachmentIds:   req.AttachmentIds,
		FinancingAct:    req.FinancingAct,
		FinancingOption: req.FinancingOption,
	}
}

func toListContractsParams(req listContractsRequest) domain.ListContractsParams {
	return domain.ListContractsParams{
		Limit:           req.PageSize,
		Offset:          (req.Page - 1) * req.PageSize,
		Search:          req.Search,
		Status:          req.Status,
		CareType:        req.CareType,
		FinancingAct:    req.FinancingAct,
		FinancingOption: req.FinancingOption,
		EndDateFrom:     req.EndDateFrom,
		EndDateTo:       req.EndDateTo,
	}
}

func getEmployeeIDFromContextForContracts(ctx *gin.Context) (uuid.UUID, error) {
	val, exists := ctx.Get("employee_id")
	if !exists {
		return uuid.Nil, fmt.Errorf("employee_id not found in context")
	}
	switch v := val.(type) {
	case uuid.UUID:
		return v, nil
	case string:
		return uuid.Parse(v)
	default:
		return uuid.Nil, fmt.Errorf("invalid employee_id type")
	}
}

func parseJSONToObject(data []byte) map[string]interface{} {
	if len(data) == 0 {
		return nil
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return map[string]interface{}{}
	}
	return result
}
