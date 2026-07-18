package main

import (
	"context"
	"fmt"
	"math"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/internal/domain"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Seeder) SeedInvoicesAndPaymentsForInCareClients(ctx context.Context, invoicesPerClient int, maxPaymentsPerInvoice int) error {
	if invoicesPerClient <= 0 || len(s.data.InCareClientIDs) == 0 {
		return nil
	}
	if maxPaymentsPerInvoice < 0 {
		maxPaymentsPerInvoice = 0
	}

	for i, clientID := range s.data.InCareClientIDs {
		if (i+1)%10 == 0 || i == 0 || i+1 == len(s.data.InCareClientIDs) {
			fmt.Printf("[seed] invoices/payments for in-care clients: %d/%d\n", i+1, len(s.data.InCareClientIDs))
		}

		employeeID, err := s.pickPaymentEmployee(ctx, clientID)
		if err != nil {
			return fmt.Errorf("pick payment employee for client %s: %w", clientID, err)
		}

		generated := 0
		maxAttempts := invoicesPerClient * 6
		basePeriodStart := time.Now().UTC().Truncate(24 * time.Hour)
		var lastGenerationErr error

		for attempt := 0; generated < invoicesPerClient && attempt < maxAttempts; attempt++ {
			// Move windows forward from "now" so they overlap the contract approved-effective period.
			periodStart := basePeriodStart.AddDate(0, 0, 28*attempt)
			periodEnd := periodStart.AddDate(0, 0, 28)

			if err := s.seedBillableAppointmentsForClient(ctx, clientID, employeeID, periodStart, periodEnd); err != nil {
				return fmt.Errorf("seed billable appointments for client %s: %w", clientID, err)
			}

			generatedInvoice, _, err := s.invoiceService.GenerateInvoice(ctx, domain.GenerateInvoiceParams{
				ClientID:        clientID,
				StartDate:       periodStart,
				EndDate:         periodEnd,
				BillingTimezone: "Europe/Amsterdam",
				BillingCycle:    "iso_4_week",
			})
			if err != nil {
				lastGenerationErr = err
				continue
			}

			s.data.InvoiceIDs = append(s.data.InvoiceIDs, generatedInvoice.ID)
			generated++

			if maxPaymentsPerInvoice == 0 || generatedInvoice.GrossTotal <= 0 {
				continue
			}

			if err := s.seedPaymentsForInvoice(ctx, generatedInvoice.ID, generatedInvoice.GrossTotal, maxPaymentsPerInvoice, employeeID); err != nil {
				return fmt.Errorf("seed payments for invoice %s: %w", generatedInvoice.ID, err)
			}
		}

		if generated < invoicesPerClient {
			return fmt.Errorf("generated %d/%d invoices for client %s: %w", generated, invoicesPerClient, clientID, lastGenerationErr)
		}
	}

	return nil
}

func (s *Seeder) pickPaymentEmployee(ctx context.Context, clientID uuid.UUID) (uuid.UUID, error) {
	if employeeID := s.pickCoordinatorForClient(clientID); employeeID != uuid.Nil {
		return employeeID, nil
	}
	if len(s.data.CoordinatorIDs) > 0 {
		return oneOf(s.data.CoordinatorIDs), nil
	}
	if len(s.data.EmployeeIDs) > 0 {
		return oneOf(s.data.EmployeeIDs), nil
	}
	return uuid.Nil, fmt.Errorf("no employee profiles available")
}

func (s *Seeder) seedBillableAppointmentsForClient(ctx context.Context, clientID uuid.UUID, employeeID uuid.UUID, periodStart, periodEnd time.Time) error {
	if !periodEnd.After(periodStart) {
		return nil
	}
	if employeeID == uuid.Nil {
		return fmt.Errorf("employee id is required for appointment seeding")
	}

	appointmentsToCreate := gofakeit.Number(2, 5)
	return s.store.ExecTx(ctx, func(q *db.Queries) error {
		for j := 0; j < appointmentsToCreate; j++ {
			startMin := gofakeit.Number(0, int(periodEnd.Sub(periodStart).Minutes())-180)
			startAt := periodStart.Add(time.Duration(startMin) * time.Minute)
			durationMinutes := gofakeit.Number(45, 150)
			endAt := startAt.Add(time.Duration(durationMinutes) * time.Minute)
			if !endAt.Before(periodEnd) {
				endAt = periodEnd.Add(-5 * time.Minute)
			}
			if !endAt.After(startAt) {
				continue
			}

			title := fmt.Sprintf("Seed care appointment %d", j+1)
			location := gofakeit.Street()
			description := "Seeded appointment for invoice generation"
			color := oneOf([]string{"#4A7C59", "#2F6D80", "#8F6A4A"})

			event, err := q.CreateCalendarEvent(ctx, db.CreateCalendarEventParams{
				OrganizerEmployeeID: employeeID,
				CreatedByEmployeeID: employeeID,
				Kind:                db.CalendarEventKindEnumAppointment,
				Status:              db.CalendarEventStatusEnumConfirmed,
				Title:               title,
				Description:         &description,
				Location:            &location,
				Color:               &color,
				StartAt:             pgtype.Timestamptz{Time: startAt, Valid: true},
				EndAt:               pgtype.Timestamptz{Time: endAt, Valid: true},
				Timezone:            "UTC",
				Rrule:               nil,
				RecurringEventID:    nil,
				RecurrenceID:        pgtype.Timestamptz{Valid: false},
			})
			if err != nil {
				return fmt.Errorf("create calendar event: %w", err)
			}

			if err := q.AddEventClientAttendee(ctx, db.AddEventClientAttendeeParams{
				EventID:  event.ID,
				ClientID: &clientID,
			}); err != nil {
				return fmt.Errorf("add client attendee: %w", err)
			}
			if err := q.AddEventEmployeeAttendee(ctx, db.AddEventEmployeeAttendeeParams{
				EventID:    event.ID,
				EmployeeID: &employeeID,
			}); err != nil {
				return fmt.Errorf("add employee attendee: %w", err)
			}

			if err := q.UpdateCalendarEventWorkApproval(ctx, db.UpdateCalendarEventWorkApprovalParams{
				WorkApprovalStatus: db.CalendarEventWorkApprovalStatusEnumApproved,
				ActorUserID:        nil,
				RejectionReason:    nil,
				EventID:            event.ID,
			}); err != nil {
				return fmt.Errorf("approve calendar event work: %w", err)
			}
		}
		return nil
	})
}

func (s *Seeder) seedPaymentsForInvoice(ctx context.Context, invoiceID uuid.UUID, grossTotal float64, maxPayments int, employeeID uuid.UUID) error {
	if maxPayments <= 0 {
		return nil
	}
	paymentCount := gofakeit.Number(1, maxPayments)
	if paymentCount < 1 {
		paymentCount = 1
	}

	remaining := grossTotal
	for i := 0; i < paymentCount; i++ {
		var amount float64
		if i == paymentCount-1 {
			amount = remaining
		} else {
			share := gofakeit.Float64Range(0.2, 0.7)
			amount = remaining * share
		}
		if amount <= 0 {
			break
		}
		if amount > remaining {
			amount = remaining
		}
		amount = roundMoney(amount)
		if amount <= 0 {
			continue
		}

		method := oneOf([]string{
			string(db.PaymentMethodEnumBankTransfer),
			string(db.PaymentMethodEnumCreditCard),
			string(db.PaymentMethodEnumCash),
			string(db.PaymentMethodEnumCheck),
		})

		paymentDate := randomRecentDate(45)
		ref := fmt.Sprintf("SEED-PMT-%s-%02d", invoiceID.String()[:8], i+1)
		notes := "Seeded payment"

		payment, err := s.invoiceService.CreatePayment(ctx, invoiceID, employeeID, domain.CreatePaymentParams{
			PaymentMethod:    method,
			PaymentStatus:    string(db.PaymentStatusEnumCompleted),
			Amount:           amount,
			PaymentDate:      paymentDate,
			PaymentReference: &ref,
			Notes:            &notes,
		})
		if err != nil {
			return fmt.Errorf("create payment %d: %w", i+1, err)
		}

		s.data.PaymentIDs = append(s.data.PaymentIDs, payment.PaymentID)
		remaining = roundMoney(remaining - amount)
		if remaining <= 0 {
			break
		}
	}
	return nil
}

func roundMoney(v float64) float64 {
	return math.Round(v*100) / 100
}
