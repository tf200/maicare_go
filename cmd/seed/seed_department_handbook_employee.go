package main

import (
	"context"
	"github.com/goccy/go-json"
	"errors"
	"fmt"
	"strings"
	"time"

	db "maicare_go/db/sqlc"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Seeder) SeedDepartments(ctx context.Context, count int) error {
	if count <= 0 {
		return nil
	}

	baseNames := []string{
		"Youth Care",
		"Family Support",
		"Behavioral Health",
		"Ambulatory Care",
		"Crisis Response",
	}

	for i := 0; i < count; i++ {
		if i == 0 || i+1 == count {
			fmt.Printf("[seed] departments: %d/%d\n", i+1, count)
		}

		name := ""
		if i < len(baseNames) {
			name = baseNames[i]
		} else {
			name = fmt.Sprintf("Seed Department %d", i+1)
		}

		dept, err := s.store.GetDepartmentByName(ctx, name)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("get department by name %q: %w", name, err)
			}

			description := fmt.Sprintf("Seeded department %d", i+1)
			dept, err = s.store.CreateDepartment(ctx, db.CreateDepartmentParams{
				Name:                     name,
				Description:              &description,
				DepartmentHeadEmployeeID: nil,
			})
			if err != nil {
				return fmt.Errorf("create department %q: %w", name, err)
			}
		}

		if !containsUUID(s.data.DepartmentIDs, dept.ID) {
			s.data.DepartmentIDs = append(s.data.DepartmentIDs, dept.ID)
		}
	}

	return nil
}

func (s *Seeder) SeedHandbookTemplates(ctx context.Context, perDepartment int) error {
	if perDepartment <= 0 {
		return nil
	}
	if len(s.data.DepartmentIDs) == 0 {
		return fmt.Errorf("no departments available; seed departments first")
	}

	for _, departmentID := range s.data.DepartmentIDs {
		for i := 0; i < perDepartment; i++ {
			if i == 0 {
				fmt.Printf("[seed] handbook templates for department %s (%d)\n", departmentID, perDepartment)
			}

			var publishedTemplateID uuid.UUID
			err := s.store.ExecTx(ctx, func(q *db.Queries) error {
				templates, err := q.ListHandbookTemplatesByDepartment(ctx, departmentID)
				if err != nil {
					return fmt.Errorf("list handbook templates by department: %w", err)
				}

				var draftTemplate *db.HandbookTemplate
				var existingPublishedTemplate *db.HandbookTemplate
				for idx := range templates {
					if templates[idx].Status == db.HandbookTemplateStatusEnumDraft {
						draftTemplate = &templates[idx]
					}
					if templates[idx].Status == db.HandbookTemplateStatusEnumPublished && existingPublishedTemplate == nil {
						existingPublishedTemplate = &templates[idx]
					}
				}

				// Keep the seed idempotent. If a department already has an active published
				// template and there is no draft to publish, reuse the active template.
				if draftTemplate == nil && existingPublishedTemplate != nil {
					publishedTemplateID = existingPublishedTemplate.ID
					return nil
				}

				if draftTemplate == nil {
					title := fmt.Sprintf("Department %s Onboarding v%s", departmentID.String()[:8], uuid.NewString()[:6])
					description := "Default onboarding handbook for seeded data"
					created, err := q.CreateHandbookTemplateForDepartment(ctx, db.CreateHandbookTemplateForDepartmentParams{
						DepartmentID:        departmentID,
						Title:               title,
						Description:         &description,
						CreatedByEmployeeID: nil,
					})
					if err != nil {
						return fmt.Errorf("create handbook template: %w", err)
					}
					draftTemplate = &created
				}

				if err := ensureDefaultHandbookSteps(ctx, q, draftTemplate.ID); err != nil {
					return err
				}

				published, err := q.PublishHandbookTemplate(ctx, db.PublishHandbookTemplateParams{
					PublishedByEmployeeID: nil,
					TemplateID:            draftTemplate.ID,
				})
				if err != nil {
					return fmt.Errorf("publish handbook template: %w", err)
				}
				publishedTemplateID = published.ID
				return nil
			})
			if err != nil {
				return fmt.Errorf("seed handbook template %d for department %s: %w", i+1, departmentID, err)
			}

			s.data.HandbookTemplateIDs = append(s.data.HandbookTemplateIDs, publishedTemplateID)
			s.data.DepartmentHandbooks[departmentID] = publishedTemplateID
		}
	}

	return nil
}

func (s *Seeder) SeedCoordinators(ctx context.Context, count int) error {
	if count <= 0 {
		return nil
	}
	if len(s.data.DepartmentIDs) == 0 {
		return fmt.Errorf("no departments available; seed departments first")
	}

	for i := 0; i < count; i++ {
		if (i+1)%10 == 0 || i == 0 || i+1 == count {
			fmt.Printf("[seed] coordinators: %d/%d\n", i+1, count)
		}

		departmentID := s.data.DepartmentIDs[i%len(s.data.DepartmentIDs)]

		var resolvedLocationID *uuid.UUID
		if len(s.data.LocationIDs) > 0 {
			picked := s.data.LocationIDs[i%len(s.data.LocationIDs)]
			resolvedLocationID = &picked
		}

		var employeeID uuid.UUID
		err := s.store.ExecTx(ctx, func(q *db.Queries) error {
			templateID, err := s.resolveDepartmentHandbookTemplateID(ctx, q, departmentID)
			if err != nil {
				return err
			}

			createdEmployeeID, err := s.createSeedCoordinatorProfileForDepartment(ctx, q, departmentID, templateID, resolvedLocationID)
			if err != nil {
				return err
			}

			employeeID = createdEmployeeID
			return nil
		})
		if err != nil {
			return fmt.Errorf("seed coordinator %d: %w", i+1, err)
		}

		s.data.EmployeeIDs = append(s.data.EmployeeIDs, employeeID)
		s.data.CoordinatorIDs = append(s.data.CoordinatorIDs, employeeID)
		s.data.HandbookAssignmentCount++
	}

	return nil
}

func (s *Seeder) resolveDepartmentHandbookTemplateID(ctx context.Context, q *db.Queries, departmentID uuid.UUID) (uuid.UUID, error) {
	if templateID, ok := s.data.DepartmentHandbooks[departmentID]; ok {
		return templateID, nil
	}

	template, err := q.GetActiveHandbookTemplateByDepartment(ctx, departmentID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("get active handbook template for department %s: %w", departmentID, err)
	}

	s.data.DepartmentHandbooks[departmentID] = template.ID
	if !containsUUID(s.data.HandbookTemplateIDs, template.ID) {
		s.data.HandbookTemplateIDs = append(s.data.HandbookTemplateIDs, template.ID)
	}

	return template.ID, nil
}

func ensureDefaultHandbookSteps(ctx context.Context, q *db.Queries, templateID uuid.UUID) error {
	existingSteps, err := q.ListHandbookStepsByTemplate(ctx, templateID)
	if err != nil {
		return fmt.Errorf("list handbook steps: %w", err)
	}
	if len(existingSteps) > 0 {
		return nil
	}

	defaultSteps := []struct {
		order int32
		kind  db.HandbookStepKindEnum
		title string
		body  string
		req   bool
	}{
		{1, db.HandbookStepKindEnumContent, "Welcome", "Welcome to the team. This handbook will guide you through your first days.", true},
		{2, db.HandbookStepKindEnumContent, "Mission & Vision", "Our mission and vision, and how your work contributes.", true},
		{3, db.HandbookStepKindEnumContent, "Your Role", "Your responsibilities, expectations, and how success is measured.", true},
		{4, db.HandbookStepKindEnumContent, "Processes", "Key processes: scheduling, reporting, incidents, and communication.", true},
		{5, db.HandbookStepKindEnumAck, "Acknowledgement", "Please acknowledge you have read and understood the handbook.", true},
	}

	for _, step := range defaultSteps {
		body := step.body
		isRequired := step.req
		if _, err := q.CreateHandbookStep(ctx, db.CreateHandbookStepParams{
			TemplateID: templateID,
			SortOrder:  step.order,
			Kind:       step.kind,
			Title:      step.title,
			Body:       &body,
			Content:    nil,
			IsRequired: &isRequired,
		}); err != nil {
			return fmt.Errorf("create handbook step: %w", err)
		}
	}

	return nil
}

func (s *Seeder) createSeedCoordinatorProfileForDepartment(
	ctx context.Context,
	q *db.Queries,
	departmentID uuid.UUID,
	templateID uuid.UUID,
	locationID *uuid.UUID,
) (uuid.UUID, error) {
	email := fmt.Sprintf("seed.coordinator.%s@maicare.local", strings.ToLower(gofakeit.LetterN(8)))
	user, err := q.CreateUser(ctx, db.CreateUserParams{
		Password:       "seed-password",
		Email:          email,
		IsActive:       true,
		ProfilePicture: nil,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("create user: %w", err)
	}

	contractHours := 36.0
	contractRate := 58.0
	employeeNumber := fmt.Sprintf("EMP-%06d", gofakeit.Number(1, 999999))
	workEmail := email
	privateEmail := strings.ToLower(gofakeit.Email())
	workPhone := fakePhone()
	privatePhone := fakePhone()
	homePhone := fakePhone()

	var managerEmployeeID *uuid.UUID
	if len(s.data.EmployeeIDs) > 0 {
		managerID := oneOf(s.data.EmployeeIDs)
		managerEmployeeID = &managerID
	}

	employee, err := q.CreateEmployeeProfile(ctx, db.CreateEmployeeProfileParams{
		UserID:              user.ID,
		FirstName:           gofakeit.FirstName(),
		LastName:            gofakeit.LastName(),
		Bsn:                 fmt.Sprintf("%09d", gofakeit.Number(100000000, 999999999)),
		Street:              gofakeit.StreetName(),
		HouseNumber:         fmt.Sprintf("%d", gofakeit.Number(1, 350)),
		HouseNumberAddition: nullableString(strings.ToUpper(gofakeit.LetterN(1)), 0.8),
		PostalCode:          fakePostalCodeNL(),
		City:                gofakeit.City(),
		Position:            stringPtr("Care Coordinator"),
		DepartmentID:        &departmentID,
		ManagerEmployeeID:   managerEmployeeID,
		EmployeeNumber:      &employeeNumber,
		PrivateEmailAddress: &privateEmail,
		WorkEmailAddress:    &workEmail,
		WorkPhoneNumber:     &workPhone,
		PrivatePhoneNumber:  &privatePhone,
		DateOfBirth:         pgDate(randomDate(1975, 1998)),
		HomeTelephoneNumber: &homePhone,
		Gender:              oneOf([]db.GenderEnum{db.GenderEnumMale, db.GenderEnumFemale, db.GenderEnumOther}),
		LocationID:          locationID,
		ContractHours:       &contractHours,
		ContractEndDate:     pgtype.Date{},
		ContractStartDate:   pgDate(time.Now().AddDate(-1, 0, 0)),
		ContractType:        db.EmployeeContractTypeEnumLoondienst,
		ContractRate:        &contractRate,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("create employee profile: %w", err)
	}

	if err := createSeedHandbookAssignment(ctx, q, employee.ID, templateID, "seed_coordinator_creation"); err != nil {
		return uuid.Nil, err
	}

	return employee.ID, nil
}

func (s *Seeder) SeedEmployeeHandbookAssignments(ctx context.Context, perDepartment int) error {
	if perDepartment <= 0 {
		return nil
	}
	if len(s.data.DepartmentIDs) == 0 {
		return fmt.Errorf("no departments available; seed departments first")
	}

	for _, departmentID := range s.data.DepartmentIDs {
		var assignedCount int
		err := s.store.ExecTx(ctx, func(q *db.Queries) error {
			templateID, err := s.resolveDepartmentHandbookTemplateID(ctx, q, departmentID)
			if err != nil {
				return err
			}

			employeeIDs, err := q.ListEmployeesEligibleForDepartmentHandbookSeed(ctx, db.ListEmployeesEligibleForDepartmentHandbookSeedParams{
				DepartmentID: &departmentID,
				Limit:        int32(perDepartment),
			})
			if err != nil {
				return fmt.Errorf("list eligible employees for department handbook seed: %w", err)
			}

			for _, employeeID := range employeeIDs {
				if err := createSeedHandbookAssignment(ctx, q, employeeID, templateID, "seed_department_assignment"); err != nil {
					return err
				}
				assignedCount++
			}

			return nil
		})
		if err != nil {
			return fmt.Errorf("seed employee handbook assignments for department %s: %w", departmentID, err)
		}

		fmt.Printf("[seed] employee handbook assignments: department=%s assigned=%d requested=%d\n", departmentID, assignedCount, perDepartment)
		s.data.HandbookAssignmentCount += assignedCount
	}

	return nil
}

func createSeedHandbookAssignment(ctx context.Context, q *db.Queries, employeeID, templateID uuid.UUID, source string) error {
	assigned, err := q.CreateEmployeeHandbookFromTemplate(ctx, db.CreateEmployeeHandbookFromTemplateParams{
		EmployeeID:           employeeID,
		TemplateID:           templateID,
		AssignedByEmployeeID: nil,
	})
	if err != nil {
		return fmt.Errorf("assign employee handbook: %w", err)
	}

	metadata, err := json.Marshal(map[string]any{
		"source": source,
	})
	if err != nil {
		return fmt.Errorf("marshal handbook assignment history metadata: %w", err)
	}

	_, err = q.CreateEmployeeHandbookAssignmentHistory(ctx, db.CreateEmployeeHandbookAssignmentHistoryParams{
		EmployeeHandbookID: &assigned.ID,
		EmployeeID:         assigned.EmployeeID,
		TemplateID:         assigned.TemplateID,
		TemplateVersion:    assigned.TemplateVersion,
		Event:              db.HandbookAssignmentEventEnumAssigned,
		ActorEmployeeID:    nil,
		Metadata:           metadata,
	})
	if err != nil {
		return fmt.Errorf("create handbook assignment history: %w", err)
	}

	return nil
}

func (s *Seeder) nextCoordinatorID() (uuid.UUID, error) {
	if len(s.data.CoordinatorIDs) == 0 {
		return uuid.Nil, fmt.Errorf("no coordinators available; seed coordinators first")
	}

	idx := s.data.NextCoordinatorIdx % len(s.data.CoordinatorIDs)
	coordinatorID := s.data.CoordinatorIDs[idx]
	s.data.NextCoordinatorIdx++
	return coordinatorID, nil
}

func containsUUID(values []uuid.UUID, needle uuid.UUID) bool {
	for _, value := range values {
		if value == needle {
			return true
		}
	}
	return false
}
