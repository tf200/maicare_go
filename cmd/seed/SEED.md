# Seed Script Guide

This folder contains the project seed command used to create mock data for local development.

- Entry point: `cmd/seed/main.go`
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
go run cmd/seed/main.go -organisations 6 -locations-per-org 2 -senders 12 -count 25
```

Example with waiting list clients:

```bash
go run cmd/seed/main.go -organisations 6 -locations-per-org 2 -senders 12 -count 30 -waiting-list-clients 20
```

## CLI Flags

- `-organisations`: number of organisations to create (default: `6`)
- `-locations-per-org`: locations per organisation (default: `2`)
- `-senders`: number of senders to create (default: `12`)
- `-count`: number of registration forms to create (default: `25`)
- `-waiting-list-clients`: number of clients to create through intake->client promotion flow (default: `12`)
- `-seed`: random seed value (default: current timestamp)
- `-timeout`: overall seed timeout duration (default: `10m`)
- `-db`: explicit DB connection string

DB connection fallback order:

1. `-db` flag
2. `DB_SOURCE` environment variable
3. local default: `postgres://maicare:maicare@localhost:5432/maicare?sslmode=disable`

## Extending for Future Seeding

When adding a new table:

1. Add ID collection in `SeedData` if needed.
2. Add a dedicated `Seed<Table>()` method in `main.go`.
3. Use SQLC `Create...` methods from `db.Store`.
4. Place the new method call in `main()` in dependency-safe order.
5. Reuse IDs from `SeedData` for foreign key fields.

## Conventions

- Keep seeding logic deterministic when needed by passing `-seed`.
- Keep fake data realistic enough for API and UI testing.
- Use pointer values (`*string`, `*bool`) where DB fields are nullable.
- Avoid raw SQL in this script; prefer generated DB methods.
