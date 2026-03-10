package main

import (
	db "maicare_go/db/sqlc"
	invoicesvc "maicare_go/service/invoice"

	"github.com/google/uuid"
)

type SeedData struct {
	OrganisationIDs     []uuid.UUID
	DepartmentIDs       []uuid.UUID
	HandbookTemplateIDs []uuid.UUID
	RegistrationFormIDs []uuid.UUID
	NextRegistrationIdx int
	IntakeFormIDs       []uuid.UUID
	ClientIDs           []uuid.UUID
	IncidentIDs         []uuid.UUID
	DiagnosisIDs        []uuid.UUID
	MedicationOrderIDs  []uuid.UUID
	InCareClientIDs     []uuid.UUID
	OutOfCareClientIDs  []uuid.UUID
	EvaluationIDs       []uuid.UUID
	EmployeeIDs         []uuid.UUID
	CoordinatorIDs      []uuid.UUID
	NextCoordinatorIdx  int
	ClientCoordinators  map[uuid.UUID]uuid.UUID
	DepartmentHandbooks map[uuid.UUID]uuid.UUID
	AttachmentIDs       []uuid.UUID
	LocationIDs         []uuid.UUID
	SenderIDs           []uuid.UUID
	InvoiceIDs          []uuid.UUID
	PaymentIDs          []uuid.UUID
}

type Seeder struct {
	store          *db.Store
	invoiceService invoicesvc.InvoiceService
	data           *SeedData
}

func newSeeder(store *db.Store, invoiceService invoicesvc.InvoiceService) *Seeder {
	return &Seeder{
		store:          store,
		invoiceService: invoiceService,
		data: &SeedData{
			ClientCoordinators:  make(map[uuid.UUID]uuid.UUID),
			DepartmentHandbooks: make(map[uuid.UUID]uuid.UUID),
		},
	}
}
