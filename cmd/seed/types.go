package main

import (
	"context"
	"fmt"

	db "maicare_go/db/sqlc"
	"maicare_go/internal/service"

	"github.com/google/uuid"
)

type SeedData struct {
	OrganisationIDs                []uuid.UUID
	DepartmentIDs                  []uuid.UUID
	HandbookTemplateIDs            []uuid.UUID
	RegistrationFormIDs            []uuid.UUID
	UnprocessedRegistrationFormIDs []uuid.UUID
	NextRegistrationIdx            int
	IntakeFormIDs                  []uuid.UUID
	ClientIDs                      []uuid.UUID
	IncidentIDs                    []uuid.UUID
	DiagnosisIDs                   []uuid.UUID
	MedicationOrderIDs             []uuid.UUID
	InCareClientIDs                []uuid.UUID
	OutOfCareClientIDs             []uuid.UUID
	EvaluationIDs                  []uuid.UUID
	EmployeeIDs                    []uuid.UUID
	CoordinatorIDs                 []uuid.UUID
	NextCoordinatorIdx             int
	ClientCoordinators             map[uuid.UUID]uuid.UUID
	DepartmentHandbooks            map[uuid.UUID]uuid.UUID
	AttachmentIDs                  []uuid.UUID
	LocationIDs                    []uuid.UUID
	SenderIDs                      []uuid.UUID
	InvoiceIDs                     []uuid.UUID
	PaymentIDs                     []uuid.UUID
	HandbookAssignmentCount        int
}

func (s *Seeder) requireRLSBootstrapRole(ctx context.Context, resource string) error {
	canBypassRLS, err := s.store.CurrentRoleBypassesRLS(ctx)
	if err != nil {
		return fmt.Errorf("check %s bootstrap database role: %w", resource, err)
	}
	if !canBypassRLS {
		return fmt.Errorf("%s bootstrap requires a database role with superuser or BYPASSRLS capability", resource)
	}
	return nil
}

type Seeder struct {
	store          *db.Store
	data           *SeedData
	invoiceService *service.InvoiceService
}

func newSeeder(store *db.Store) *Seeder {
	return &Seeder{
		store:          store,
		invoiceService: service.NewInvoiceService(store, nil, nil, nil),
		data: &SeedData{
			ClientCoordinators:  make(map[uuid.UUID]uuid.UUID),
			DepartmentHandbooks: make(map[uuid.UUID]uuid.UUID),
		},
	}
}
