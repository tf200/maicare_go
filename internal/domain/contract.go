package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ==================== ContractType ====================

type ContractType struct {
	ID   uuid.UUID
	Name string
}

// ==================== Contract ====================

type Contract struct {
	ID              uuid.UUID
	TypeID          *uuid.UUID
	Status          string
	ApprovedAt      *time.Time
	StartDate       time.Time
	EndDate         time.Time
	ReminderPeriod  int32
	Vat             *int32
	Price           float64
	PriceTimeUnit   string
	Hours           *float64
	HoursType       *string
	CareName        string
	CareType        string
	ClientID        uuid.UUID
	SenderID        uuid.UUID
	AttachmentIds   []uuid.UUID
	FinancingAct    string
	FinancingOption string
	DepartureReason *string
	DepartureReport *string
	UpdatedAt       time.Time
	CreatedAt       time.Time
}

// ContractDetail includes joined client and sender info
type ContractDetail struct {
	Contract
	TypeName                  *string
	ClientFirstName           string
	ClientLastName            string
	ClientFilenumber          string
	ClientBsn                 *string
	SenderName                string
	SenderType                string
	SenderStreet              *string
	SenderHouseNumber         *string
	SenderHouseNumberAddition *string
	SenderPostalCode          *string
	SenderCity                *string
	SenderLand                *string
	SenderKvknumber           *string
	SenderBtwnumber           *string
	SenderPhoneNumber         *string
	SenderClientNumber        *string
	SenderEmailAddress        *string
}

// ==================== ContractAttachment ====================

type ContractAttachment struct {
	ID          uuid.UUID
	Name        string
	Size        int64
	DownloadURL string
}

// ==================== ContractAuditLog ====================

type ContractAuditLog struct {
	AuditID            uuid.UUID
	ContractID         uuid.UUID
	Operation          string
	ChangedBy          *uuid.UUID
	ChangedAt          time.Time
	OldValues          []byte
	NewValues          []byte
	ChangedFields      []string
	ChangedByFirstName *string
	ChangedByLastName  *string
}

// ==================== Parameters ====================

type CreateContractTypeParams struct {
	Name string
}

type CreateContractParams struct {
	ClientID        uuid.UUID
	TypeID          *uuid.UUID
	StartDate       time.Time
	EndDate         time.Time
	ReminderPeriod  *int32
	Vat             *int32
	Price           float64
	PriceTimeUnit   string
	Hours           *float64
	HoursType       *string
	CareName        string
	CareType        string
	SenderID        uuid.UUID
	AttachmentIds   []uuid.UUID
	FinancingAct    string
	FinancingOption string
}

type UpdateContractParams struct {
	ContractID      uuid.UUID
	TypeID          *uuid.UUID
	StartDate       *time.Time
	EndDate         *time.Time
	ReminderPeriod  *int32
	Vat             *int32
	Price           *float64
	PriceTimeUnit   *string
	Hours           *float64
	HoursType       *string
	CareName        *string
	SenderID        *uuid.UUID
	AttachmentIds   []uuid.UUID
	FinancingAct    *string
	FinancingOption *string
}

type UpdateContractStatusParams struct {
	ContractID uuid.UUID
	Status     string
}

type ListClientContractsParams struct {
	ClientID uuid.UUID
	Limit    int32
	Offset   int32
}

type ListContractsParams struct {
	Limit           int32
	Offset          int32
	Search          *string
	Status          []string
	CareType        []string
	FinancingAct    []string
	FinancingOption *string
	EndDateFrom     *time.Time
	EndDateTo       *time.Time
}

// ==================== Results ====================

type CreateContractTypeResult struct {
	ID   uuid.UUID
	Name string
}

type CreateContractResult struct {
	Contract
	Attachments []ContractAttachment
}

type UpdateContractResult struct {
	Contract
}

type UpdateContractStatusResult struct {
	ID     uuid.UUID
	Status string
}

type ListClientContractsResult struct {
	Contracts  []ClientContractListItem
	TotalCount int64
}

type ClientContractListItem struct {
	StartDate       time.Time
	EndDate         time.Time
	DaysLeft        int32
	CareName        string
	CareType        string
	FinancingAct    string
	FinancingOption string
}

type ListContractsResult struct {
	Contracts  []ContractListItem
	TotalCount int64
}

type ContractListItem struct {
	ID               uuid.UUID
	ClientID         uuid.UUID
	ClientFirstName  string
	ClientLastName   string
	ClientFilenumber string
	SenderID         uuid.UUID
	SenderName       string
	CareName         string
	CareType         string
	Price            float64
	PriceTimeUnit    string
	Hours            *float64
	HoursType        *string
	FinancingAct     string
	FinancingOption  string
	StartDate        time.Time
	EndDate          time.Time
	DaysLeft         int32
	Status           string
	ApprovedAt       *time.Time
	UpdatedAt        time.Time
}

// ==================== Interfaces ====================

type ContractRepository interface {
	CreateContractType(ctx context.Context, params CreateContractTypeParams) (*ContractType, error)
	ListContractTypes(ctx context.Context) ([]ContractType, error)
	DeleteContractType(ctx context.Context, contractTypeID uuid.UUID) error

	CreateContract(ctx context.Context, params CreateContractParams) (*Contract, error)
	GetContractByID(ctx context.Context, contractID uuid.UUID) (*ContractDetail, error)
	UpdateContract(ctx context.Context, params UpdateContractParams, employeeID uuid.UUID) (*Contract, error)
	UpdateContractStatus(ctx context.Context, params UpdateContractStatusParams, employeeID uuid.UUID) (*Contract, error)
	ListClientContracts(ctx context.Context, params ListClientContractsParams) ([]ClientContractListItem, int64, error)
	ListContracts(ctx context.Context, params ListContractsParams) ([]ContractListItem, int64, error)
	GetContractAuditLog(ctx context.Context, contractID uuid.UUID) ([]ContractAuditLog, error)

	GetAttachmentFiles(ctx context.Context, ids []uuid.UUID) ([]AttachmentFile, error)
	GetActorAttachmentFiles(ctx context.Context, ids []uuid.UUID) ([]AttachmentFile, error)
}

type ContractService interface {
	CreateContractType(ctx context.Context, params CreateContractTypeParams) (*CreateContractTypeResult, error)
	ListContractTypes(ctx context.Context) ([]CreateContractTypeResult, error)
	DeleteContractType(ctx context.Context, contractTypeID uuid.UUID) (*CreateContractTypeResult, error)

	CreateContract(ctx context.Context, params CreateContractParams) (*CreateContractResult, error)
	GetContractByID(ctx context.Context, contractID uuid.UUID) (*ContractDetail, error)
	UpdateContract(ctx context.Context, params UpdateContractParams, employeeID uuid.UUID) (*UpdateContractResult, error)
	UpdateContractStatus(ctx context.Context, params UpdateContractStatusParams, employeeID uuid.UUID) (*UpdateContractStatusResult, error)
	ListClientContracts(ctx context.Context, params ListClientContractsParams) (*ListClientContractsResult, error)
	ListContracts(ctx context.Context, params ListContractsParams) (*ListContractsResult, error)
	GetContractAuditLog(ctx context.Context, contractID uuid.UUID) ([]ContractAuditLog, error)
}
