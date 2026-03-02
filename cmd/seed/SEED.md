# Seed Script Guide

This folder contains the project seed command used to create mock data for local development.

- Entry point: package directory `cmd/seed` (run with `go run ./cmd/seed ...`)
- Faker library: `github.com/brianvoe/gofakeit/v7`
- DB access: app-generated SQLC methods through `db.Store`

## Current Seeding Scope

The script currently seeds the following, in this order:

1. `organisations`
2. `location` (depends on organisation IDs)
3. `sender`
4. `registration_form`
5. `intake_forms` + `intake_topic_assessments` (suitable only)
6. `client_details` (promoted to **On Waiting List**)
7. `client_details` promoted to **In Care** with:
   - care dates (`placed_in_care_at`, `care_start_date`)
   - coordinator assignment (`assigned_employee` with role `coordinator`)
   - approved active contract (`contract`)
8. `client_details` promoted to **Out Of Care** with:
   - full intake + waiting-list + in-care history
   - seeded evaluations before discharge
   - discharge data (`discharge_date`, `discharge_reason`, `final_evaluation`)
   - status history transition to `out_of_care`
9. `client_goal_evaluations` + `client_goal_evaluation_items` for active in-care clients
10. `client_diagnosis` + `client_medication_order` for seeded clients
11. `invoice` + `invoice_line` + `invoice_payment_history` for in-care clients via invoice service logic:
   - seeds billable `calendar_events` appointments per in-care client in 4-week windows
   - calls `GenerateInvoice` (auto logic) to create invoices from approved contracts
   - calls `CreatePayment` to create completed payments and trigger invoice status transitions

This order is intentional and should be kept for FK safety when future tables are added.

## How It Works

- A `Seeder` struct wraps the `db.Store` and a `SeedData` cache.
- `SeedData` stores created IDs so later seed steps can reuse earlier records.
- Each table has its own method (for example `SeedOrganisations`, `SeedLocations`, `SeedSenders`, `SeedRegistrationForms`).
- Random values are generated with helper functions (`fakePhone`, `fakePostalCodeNL`, `randomDate`, etc.).

## Run the Seeder

### Make target

```bash
make seed
```

### Direct run

```bash
go run ./cmd/seed -organisations 6 -locations-per-org 2 -senders 12 -count 25
```

Example with waiting list clients:

```bash
go run ./cmd/seed -organisations 6 -locations-per-org 2 -senders 12 -count 30 -waiting-list-clients 20 -in-care-clients 8
```

Example with out-of-care lifecycle clients:

```bash
go run ./cmd/seed -waiting-list-clients 20 -in-care-clients 8 -out-of-care-clients 6
```

Example with in-care evaluations:

```bash
go run ./cmd/seed -in-care-clients 8 -evaluations-per-in-care-client 2
```

Example with diagnosis and medication data:

```bash
go run ./cmd/seed -waiting-list-clients 20 -in-care-clients 8 -diagnoses-per-client 2 -medication-orders-per-client 3
```

Example with invoice/payment data:

```bash
go run ./cmd/seed -in-care-clients 8 -invoices-per-client 2 -payments-per-invoice 2
```

## CLI Flags

- `-organisations`: number of organisations to create (default: `6`)
- `-locations-per-org`: locations per organisation (default: `2`)
- `-senders`: number of senders to create (default: `12`)
- `-count`: number of registration forms to create (default: `25`)
- `-waiting-list-clients`: number of clients to create through intake->client promotion flow (default: `12`)
- `-in-care-clients`: number of clients to create and promote to in-care pattern (default: `6`)
- `-out-of-care-clients`: number of clients to create and promote through in-care to out-of-care (default: `4`)
- `-evaluations-per-in-care-client`: goal evaluations to create per in-care client (default: `2`)
- `-diagnoses-per-client`: max diagnoses to generate per client (default: `2`)
- `-medication-orders-per-client`: max medication orders to generate per client (default: `2`)
- `-invoices-per-client`: number of generated invoices to seed per in-care client (default: `1`)
- `-payments-per-invoice`: max number of completed payments to seed per generated invoice (default: `2`)
- `-seed`: random seed value (default: current timestamp)
- `-timeout`: overall seed timeout duration (default: `10m`)
- `-db`: explicit DB connection string

## In-Care Data Pattern

For each requested in-care client, the seeder performs a realistic flow:

1. Seeds a regular waiting-list client through registration + intake + goals.
2. Promotes that client to status `in_care` using `PutClientInCare`.
3. Creates a dedicated coordinator user/profile and assigns them as main coordinator.
4. Creates an `approved` contract with pricing/hours fields aligned to care type constraints.
5. Seeds compact goal-evaluation history per in-care client:
   - tries to create one `completed` evaluation when the 14-day completion window allows it
   - fills the remainder as `draft` evaluations
   - adds `client_goal_evaluation_items` for each active goal

This keeps seeded in-care clients compatible with in-care listing endpoints and dashboard logic.

DB connection fallback order:

1. `-db` flag
2. `DB_SOURCE` environment variable
3. local default: `postgres://maicare:maicare@localhost:5432/maicare?sslmode=disable`

## Extending for Future Seeding

## File Layout

- `cmd/seed/main.go`: CLI flags, DB bootstrap, and execution order.
- `cmd/seed/types.go`: `SeedData`, `Seeder`, and `newSeeder`.
- `cmd/seed/seed_org_sender.go`: organisation/location/sender seeding methods.
- `cmd/seed/seed_registration_waiting.go`: registration form + waiting-list seeding.
- `cmd/seed/seed_incare_eval.go`: in-care promotion and goal-evaluation seeding, including coordinator/evaluation helpers.
- `cmd/seed/seed_medical.go`: diagnosis and medication-order seeding for clients.
- `cmd/seed/helpers_fake.go`: generic fake data/date/pointer helper functions.
- `cmd/seed/helpers_domain.go`: domain derivation and selector/helper logic.

Add new seed logic in the closest matching `seed_*.go` file, keep helper utilities in the `helpers_*.go` files, and keep `main.go` focused on orchestration only.

When adding a new table:

1. Add ID collection in `SeedData` if needed.
2. Add a dedicated `Seed<Table>()` method in the relevant `seed_*.go` file.
3. Use SQLC `Create...` methods from `db.Store`.
4. Place the new method call in `main()` in dependency-safe order.
5. Reuse IDs from `SeedData` for foreign key fields.

## Conventions

- Keep seeding logic deterministic when needed by passing `-seed`.
- Keep fake data realistic enough for API and UI testing.
- Use pointer values (`*string`, `*bool`) where DB fields are nullable.
- Avoid raw SQL in this script; prefer generated DB methods.
