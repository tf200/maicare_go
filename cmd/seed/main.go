package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"maicare_go/bucket"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/service/deps"
	invoicesvc "maicare_go/service/invoice"
	"maicare_go/util"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	organisationCount := flag.Int("organisations", 6, "number of organisations to seed")
	locationsPerOrg := flag.Int("locations-per-org", 2, "locations per organisation")
	senderCount := flag.Int("senders", 12, "number of senders to seed")
	count := flag.Int("count", 25, "number of registration forms to seed")
	otherIntakeForms := flag.Int("other-intake-forms", 10, "number of non-suitable intake forms to seed without promoting to clients")
	waitingListClients := flag.Int("waiting-list-clients", 12, "number of waiting list clients to seed via intake promotion flow")
	inCareClients := flag.Int("in-care-clients", 6, "number of in-care clients to seed via waiting-list to in-care promotion flow")
	outOfCareClients := flag.Int("out-of-care-clients", 4, "number of out-of-care clients to seed via waiting-list to in-care to out-of-care promotion flow")
	evaluationsPerInCareClient := flag.Int("evaluations-per-in-care-client", 2, "number of goal evaluations to seed for each in-care client")
	diagnosesPerClient := flag.Int("diagnoses-per-client", 2, "max number of diagnoses to seed per client")
	medicationOrdersPerClient := flag.Int("medication-orders-per-client", 2, "max number of medication orders to seed per client")
	incidentsPerClient := flag.Int("incidents-per-client", 2, "max number of incidents to seed per client")
	invoicesPerClient := flag.Int("invoices-per-client", 1, "number of generated invoices to seed per in-care client")
	paymentsPerInvoice := flag.Int("payments-per-invoice", 2, "max number of payments to seed per generated invoice")
	seedValue := flag.Int64("seed", time.Now().UnixNano(), "random seed")
	seedTimeout := flag.Duration("timeout", 10*time.Minute, "overall seed timeout")
	dataSource := flag.String("db", "", "database connection string (defaults to DB_SOURCE or local default)")
	flag.Parse()

	gofakeit.Seed(*seedValue)

	ctx, cancel := context.WithTimeout(context.Background(), *seedTimeout)
	defer cancel()

	dsn := strings.TrimSpace(*dataSource)
	if dsn == "" {
		appEnvDSN, err := dbSourceFromAppEnv("app.env")
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				log.Printf("[seed] warning: cannot read DB_SOURCE from app.env: %v", err)
			}
		} else {
			dsn = strings.TrimSpace(appEnvDSN)
		}
	}
	if dsn == "" {
		dsn = strings.TrimSpace(os.Getenv("DB_SOURCE"))
	}
	if dsn == "" {
		dsn = "postgres://maicare:maicare@localhost:5432/maicare?sslmode=disable"
	}

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("cannot parse db config: %v", err)
	}

	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		return db.RegisterEnumTypes(ctx, conn)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("cannot ping db: %v", err)
	}

	store := db.NewStore(pool)

	appLogger, err := logger.SetupLogger("development")
	if err != nil {
		log.Fatalf("cannot setup logger: %v", err)
	}

	serviceDeps := deps.NewServiceDependencies(
		store,
		nil,
		appLogger,
		&util.Config{},
		bucket.NewNoopObjectStorageClient(),
		nil,
		nil,
		nil,
	)
	invoiceService := invoicesvc.NewInvoiceService(serviceDeps)

	seeder := newSeeder(store, invoiceService)

	startedAt := time.Now()
	fmt.Printf("[seed] start organisations=%d locations_per_org=%d senders=%d registration_forms=%d other_intake_forms=%d waiting_list_clients=%d in_care_clients=%d out_of_care_clients=%d evaluations_per_in_care_client=%d diagnoses_per_client=%d medication_orders_per_client=%d incidents_per_client=%d invoices_per_client=%d payments_per_invoice=%d timeout=%s\n",
		*organisationCount, *locationsPerOrg, *senderCount, *count, *otherIntakeForms, *waitingListClients, *inCareClients, *outOfCareClients, *evaluationsPerInCareClient, *diagnosesPerClient, *medicationOrdersPerClient, *incidentsPerClient, *invoicesPerClient, *paymentsPerInvoice, (*seedTimeout).String())
	if err := seeder.SeedOrganisations(ctx, *organisationCount); err != nil {
		log.Fatalf("seeding organisations failed: %v", err)
	}

	if err := seeder.SeedLocations(ctx, *locationsPerOrg); err != nil {
		log.Fatalf("seeding locations failed: %v", err)
	}

	if err := seeder.SeedSenders(ctx, *senderCount); err != nil {
		log.Fatalf("seeding senders failed: %v", err)
	}

	if err := seeder.SeedRegistrationForms(ctx, *count); err != nil {
		log.Fatalf("seeding registration forms failed: %v", err)
	}

	if err := seeder.SeedWaitingListClients(ctx, *waitingListClients); err != nil {
		log.Fatalf("seeding waiting list clients failed: %v", err)
	}

	if err := seeder.SeedInCareClients(ctx, *inCareClients); err != nil {
		log.Fatalf("seeding in-care clients failed: %v", err)
	}

	if err := seeder.SeedOutOfCareClients(ctx, *outOfCareClients, *evaluationsPerInCareClient); err != nil {
		log.Fatalf("seeding out-of-care clients failed: %v", err)
	}

	if err := seeder.SeedOtherIntakeForms(ctx, *otherIntakeForms); err != nil {
		log.Fatalf("seeding non-suitable intake forms failed: %v", err)
	}

	if err := seeder.SeedGoalEvaluationsForInCareClients(ctx, *evaluationsPerInCareClient); err != nil {
		log.Fatalf("seeding goal evaluations for in-care clients failed: %v", err)
	}

	if err := seeder.SeedMedicalForClients(ctx, *diagnosesPerClient, *medicationOrdersPerClient); err != nil {
		log.Fatalf("seeding client diagnoses and medication orders failed: %v", err)
	}

	if err := seeder.SeedIncidentsForClients(ctx, *incidentsPerClient); err != nil {
		log.Fatalf("seeding incidents failed: %v", err)
	}

	if err := seeder.SeedInvoicesAndPaymentsForInCareClients(ctx, *invoicesPerClient, *paymentsPerInvoice); err != nil {
		log.Fatalf("seeding invoices and payments failed: %v", err)
	}

	fmt.Printf("Seeded %d organisations, %d locations, %d senders, %d registration forms, %d intake forms, %d total clients, %d waiting list clients, %d in-care clients, %d out-of-care clients, %d coordinators, %d goal evaluations, %d diagnoses, %d medication orders, %d incidents, %d invoices, %d payments in %s\n",
		len(seeder.data.OrganisationIDs),
		len(seeder.data.LocationIDs),
		len(seeder.data.SenderIDs),
		len(seeder.data.RegistrationFormIDs),
		len(seeder.data.IntakeFormIDs),
		len(seeder.data.ClientIDs),
		len(seeder.data.ClientIDs)-len(seeder.data.InCareClientIDs)-len(seeder.data.OutOfCareClientIDs),
		len(seeder.data.InCareClientIDs),
		len(seeder.data.OutOfCareClientIDs),
		len(seeder.data.CoordinatorIDs),
		len(seeder.data.EvaluationIDs),
		len(seeder.data.DiagnosisIDs),
		len(seeder.data.MedicationOrderIDs),
		len(seeder.data.IncidentIDs),
		len(seeder.data.InvoiceIDs),
		len(seeder.data.PaymentIDs),
		time.Since(startedAt).Round(time.Millisecond),
	)
	if len(seeder.data.RegistrationFormIDs) > 0 {
		fmt.Printf("First ID: %s\n", seeder.data.RegistrationFormIDs[0])
		fmt.Printf("Last ID:  %s\n", seeder.data.RegistrationFormIDs[len(seeder.data.RegistrationFormIDs)-1])
	}
}

func dbSourceFromAppEnv(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	for _, rawLine := range strings.Split(string(content), "\n") {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok || strings.TrimSpace(key) != "DB_SOURCE" {
			continue
		}

		value = strings.TrimSpace(value)
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		return strings.TrimSpace(value), nil
	}

	return "", nil
}
