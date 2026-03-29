# Refactoring Guide — Clean Architecture with Gin + SQLC

## Goal

Refactor the existing codebase into a clean, layered architecture using Gin and SQLC.
The goal is not to rewrite business logic — it is to reorganize existing code into the
correct layers with clear boundaries and consistent patterns.

## Stack

- **Gin** — HTTP router and handler layer
- **SQLC** — type-safe SQL query generation
- **pgx / pgtype** — postgres driver and type system
- **golang-jwt** — JWT signing and verification
- **aws-sdk-go** — S3 storage

## Definition of Done

A flow is considered complete when:

- Domain structs and interfaces are defined in `internal/domain/`
- Repository implements the domain interface, all SQLC types are converted via `pkg/conv`, no pgtype leaks out
- Service contains all business logic and depends only on domain interfaces
- Handler defines its own request/response DTOs, maps to/from domain models via private functions
- The endpoint is reachable and returns the correct response
- No layer imports something it should not (see rules section)
- Progress is recorded in the refactor progress document after the flow is finished

---

## Folder Structure

```
myapp/
├── cmd/
│   └── server/
│       └── main.go                  # entry point, wires everything (DI happens here)
├── internal/
│   ├── domain/
│   │   ├── user.go                  # User (lean), UserDetail (rich with joins)
│   │   ├── repository.go            # UserRepository interface (List, GetByID)
│   │   ├── service.go               # UserService interface
│   │   ├── storage.go               # Storage interface (wraps S3)
│   │   ├── auth.go                  # TokenClient interface (wraps JWT)
│   │   └── errors.go                # ErrUserNotFound, ErrUnauthorized etc.
│   ├── repository/
│   │   ├── db/                      # SQLC generated output (never touched manually)
│   │   │   ├── models.go
│   │   │   ├── queries.sql.go
│   │   │   └── db.go
│   │   └── user_repo.go             # implements domain.UserRepository, uses pkg/conv
│   ├── service/
│   │   └── user_service.go          # business logic, depends on domain interfaces only
│   ├── handler/
│   │   └── user_handler.go          # gin handler, request/response DTOs, mapper funcs
│   └── middleware/
│       └── auth.go                  # depends on domain.TokenClient
├── pkg/
│   ├── conv/
│   │   └── pg.go                    # all pgtype helpers (ToPgText, FromPgTimestamp etc.)
│   ├── s3/
│   │   └── client.go                # wraps aws sdk, implements domain.Storage
│   ├── jwt/
│   │   └── jwt.go                   # wraps golang-jwt, implements domain.TokenClient
│   └── sms/
│       └── client.go                # wraps twilio etc., implements domain.SMSSender
├── config/
│   └── config.go                    # loads env/yaml, holds typed config structs
├── migrations/
│   └── 001_create_users.sql
└── sqlc.yaml
```

**Notes:**
- `internal/repository/db/` is never manually edited — always SQLC generated
- `cmd/server/main.go` is where all the concrete types get wired together — it is the only place that knows about everything
- During migration, keep the current `main.go` unchanged and unwire the new stack until the selected migration milestone is complete
- Build the new structure in parallel with the old one; do not mass-move the old codebase first

---

## Data Flow

```
HTTP Request
    ↓
Handler     (parse request → call service with primitive types or request DTO)
    ↓
Service     (business logic → call repository with primitive types or domain model)
    ↓
Repository  (convert to SQLC params → query DB)
    ↓
DB (SQLC)   (returns generated row)
    ↑
Repository  (map SQLC row → domain model via private toDomainX funcs)
    ↑
Service     (apply business logic → return domain model)
    ↑
Handler     (map domain model → response DTO via private toXResponse funcs)
    ↓
HTTP Response
```

---

## Domain Entities

Domain structs represent the **full truth** of an entity. They have no JSON tags, no pgtype, no HTTP concerns.

When a `GET` endpoint requires heavy joins that a `LIST` endpoint does not need, define **two separate domain structs**:

```
User          → lean, used for list queries (no joins)
UserDetail    → rich, used for get by id (with joins)
```

This is a **data fetching decision**, not a response shaping one. The rule is:

- Different **response shape only** → same domain model, different handler response DTO
- Different **data fetched** (joins, extra tables) → different domain model, different repository method

---

## Mapping Rules

### db → domain (repository layer)

- Handled by private functions in the repository file (`toDomainUser`, `toDomainUserDetail`)
- Uses `pkg/conv` for all pgtype conversions
- SQLC types **never leave** the repository layer

```go
func toDomainUser(row db.User) *domain.User { ... }
func toDomainUserDetail(row db.GetUserByIDRow) *domain.UserDetail { ... }
```

### domain → response (handler layer)

- Handled by private functions in the handler file (`toGetUserResponse`, `toListUserResponse`)
- Each endpoint has its **own response DTO** defined in the handler package
- Response DTOs **never enter** the service or repository
- When a handler file starts to grow, move DTOs and mapper helpers into a dedicated `*_dto.go` file in the same handler package

```go
func toGetUserResponse(u *domain.UserDetail) getUserResponse { ... }
func toListUserResponse(u *domain.User) listUserResponse { ... }
```

### Summary

```
db.Row        →  toDomainX()      →  domain.Model    (repository)
domain.Model  →  toXResponse()    →  ResponseDTO     (handler)
```

Additional rules:
- Generic primitive or `pgtype` conversion helpers belong in `pkg/conv`
- Boundary-specific mapping stays private in the repository or handler file that owns it
- Service code should not keep ad hoc SQLC or `pgtype` conversion helpers
- Do not copy old helpers into the new stack when they can be centralized or removed
- Keep type names aligned with actual scope; do not leave narrow names on types that now serve a broader feature

---

## Utilities

### pkg/conv — pgtype conversions

All pgtype conversion helpers live in `pkg/conv/pg.go`. This is the **only place** in the codebase that deals with pgtype conversions.

```
pkg/conv/pg.go   → ToPgText, FromPgText, ToPgTimestamp, FromPgTimestamp etc.
```

Rules:
- `pkg/conv` is **only imported by the repository layer**
- Service, handler, and domain never see pgtype
- Covers all cases: required fields, nullable fields (`*string`, `*time.Time`)
- Existing scattered conversion helpers should be moved here or deleted if they are boundary-specific

### pkg/ — third party library wrappers

Every third party library gets a **thin wrapper** in `pkg/`. App code never imports the library directly.

```
pkg/s3/       → wraps aws sdk
pkg/jwt/      → wraps golang-jwt
pkg/sms/      → wraps twilio etc.
```

Each wrapper follows the same pattern:

1. **Domain defines the interface** — lives in `internal/domain/`
2. **`pkg/` implements it** — wraps the third party lib
3. **App code depends on the interface** — never the concrete wrapper

```
third party lib    →   pkg/ wrapper    →   domain interface   →   app code
aws-sdk-go         →   pkg/s3          →   domain.Storage     →   service
golang-jwt         →   pkg/jwt         →   domain.TokenClient →   middleware
twilio-go          →   pkg/sms         →   domain.SMSSender   →   service
```

The same pattern applies to shared infrastructure such as logging:

```
internal/domain/logger.go  → logger interface
pkg/logger/                → logger implementation
internal/service/...       → depends on domain.Logger
```

### Legacy `util/` package

The current `util/` package is treated as a legacy package during migration.
Do not carry it forward as-is and do not create new helpers inside it for refactored flows.

Rules:
- `util/config.go` should move to `config/`
- Reusable `pgtype` and primitive conversion helpers should move to `pkg/conv`
- Third-party or infrastructure helpers should move to a narrowly named `pkg/` package
- Feature-specific helpers should move close to the owning feature
- Test-only helpers should not stay in production utility packages
- Pointer helper functions should not be copied into the new architecture unless clearly justified

For new `internal/` code, prefer proper ownership over convenience:
- do not import `util` by default
- move, replace, or delete a legacy helper when a migrated flow depends on it

### Struct mapping helpers

Inline mapping is avoided in favour of **private mapper functions**. This applies to both repository and handler layers:

```
repository/user_repo.go  →  toDomainUser(), toDomainUserDetail()
handler/user_handler.go  →  toGetUserResponse(), toListUserResponse()
```

Reasons:
- Handlers and repositories stay readable
- Mappers are independently testable
- Consistent pattern across the codebase makes it easy for tooling and LLMs to follow

---

## Refactoring Strategy

### Approach — Endpoint by Endpoint, Not Layer by Layer

Refactor one complete endpoint at a time, from domain to handler.
Do not refactor all repositories first, then all services, or batch multiple
endpoints together — that creates broken intermediate states that cannot be
tested or verified.

Each endpoint must be fully covered before moving to the next endpoint.

### Migration Constraints

- Keep the old runtime wiring in place until the new stack is ready to be wired
- Do not make the old `api/` and `service/` packages depend on `internal/`
- Do not make the new `internal/` packages depend on the old `api/` or `service/` packages
- Avoid temporary duplication unless it is strictly needed to keep the migration moving
- Carry forward existing query filters, sorting inputs, and pagination behavior unless there is an explicit decision to change the endpoint contract

### Order Within Each Flow

For every endpoint, always work in this order:

1. **domain/** — define structs, interfaces, and errors first. Everything else depends on this.
2. **repository/** — implement the domain interface, write mapper functions, use pkg/conv
3. **service/** — implement business logic against the domain interface
4. **handler/** — define request/response DTOs, mapper functions, wire up gin route

### Grouping — By Domain Entity

Group endpoint work by entity for planning, but still finish one endpoint before
starting another. This keeps domain files stable without turning the migration
into a batch refactor.

```
Round 1 — users:    GET /users/:id,  GET /users,  POST /users
Round 2 — orders:   GET /orders/:id, POST /orders
Round 3 — auth:     POST /login,     POST /refresh
```

### Verification Checklist Per Flow

Before marking a flow as done and moving to the next, verify:

- [ ] Endpoint is reachable and returns correct response
- [ ] No pgtype in service, handler, or domain
- [ ] No JSON tags on domain structs
- [ ] No business logic in handler
- [ ] No Gin imports in service or repository
- [ ] SQLC types do not leave the repository
- [ ] Response DTOs do not enter the service or repository
- [ ] All db→domain mapping done via private `toDomainX` functions in repository
- [ ] All domain→response mapping done via private `toXResponse` functions in handler
- [ ] Third party libs are not imported directly outside of `pkg/`

---

## Rules

These rules are non-negotiable. The agent must check every file it produces against
this list before considering it done.

### Layer Import Rules

| Layer | Allowed imports | Never import |
|---|---|---|
| `domain/` | standard library only | gin, pgtype, sqlc, pkg/conv, any third party lib |
| `repository/` | domain, pkg/conv, sqlc generated db/, pgx | gin, service, handler |
| `service/` | domain | gin, pgtype, sqlc, pkg/conv, handler |
| `handler/` | domain, gin | pgtype, sqlc, pkg/conv, repository directly |
| `pkg/conv/` | pgtype, standard library | domain, internal/ |
| `pkg/s3`, `pkg/jwt` etc. | third party lib, standard library | internal/ |

### Domain Rules

- No JSON tags on domain structs
- No pgtype on domain structs
- No HTTP concerns (status codes, gin context) in domain
- Domain errors are plain Go errors defined in `internal/domain/errors.go`
- Lean vs rich structs (`User` vs `UserDetail`) are a data fetching decision, not a response shaping one

### Repository Rules

- SQLC output in `internal/repository/db/` is never manually edited
- All pgtype conversions go through `pkg/conv` — no inline pgtype construction
- Every db→domain mapping is done via a private `toDomainX` function, never inline
- DB errors are mapped to domain errors (`sql.ErrNoRows` → `domain.ErrNotFound`)
- SQLC types never leave this layer

### Service Rules

- No knowledge of HTTP — no status codes, no gin, no request/response structs
- Depends only on domain interfaces, never on concrete repository or pkg/ types
- All business logic and validation lives here, nowhere else
- For logging in new `internal/` flows, use `domain.Logger` with `LogError`, `LogWarn`, and `LogInfo`
- Do not introduce new `LogBusinessEvent` usage in refactored flows

### Handler Rules

- No business logic — if it feels like a rule or a decision, it belongs in the service
- Each endpoint defines its own request and response DTO in the handler layer
- Prefer splitting large handler files into `<feature>_handler.go` and `<feature>_dto.go`
- Route registration belongs to the handler layer and may live in the handler file by default
- If route registration becomes noisy, split it into a dedicated `<feature>_routes.go` file
- Route registration code should only map paths to handler methods; it must not construct dependencies
- Swagger/OpenAPI annotations belong in the handler layer only and should stay attached to the handler method for the endpoint
- Service, repository, and domain files must not contain Swagger/OpenAPI annotations
- During parallel migration, temporary duplicate Swagger comments are acceptable only until the old route is replaced
- After wiring a new route, remove or update the old Swagger comments so the endpoint has a single source of truth
- All domain→response mapping done via private `toXResponse` functions, never inline
- Domain errors are mapped to HTTP status codes here and only here
- Response DTOs never leave the handler layer

### Middleware Rules

- New middleware belongs in `internal/middleware/`; keep the old `api/` middleware unchanged until the new stack is wired
- Middleware must not depend on the old `Server` type; inject only the small dependencies it actually needs
- Build middleware as separate concerns: auth, RBAC, request context, request logging, and audit
- New middleware should prefer plain `context.Context` values over Gin-only access patterns so services can consume the same request metadata
- Auth middleware may reuse the current token maker abstraction during migration
- New middleware logging must use the new logger interface and `LogError`, `LogWarn`, `LogInfo`
- Do not introduce new `LogBusinessEvent` usage in middleware for refactored flows
- Audit middleware must not use the live `*gin.Context` in background goroutines

### pkg/ Rules

- Every third party library is wrapped in `pkg/` and never imported directly by app code
- Domain defines the interface, `pkg/` implements it
- All library config (keys, secrets, timeouts) stays inside the wrapper, never passed around

---

## Agent Quick Reference — Per Flow Checklist

Use this before starting and after completing every flow.

### Before You Start a Flow

- [ ] Identify the endpoint and its HTTP method and path (`GET /users/:id`)
- [ ] Identify the domain entity it belongs to (`user`)
- [ ] Check if a domain file for this entity already exists — if yes, extend it, do not recreate it
- [ ] Check if a repository file for this entity already exists — if yes, add the method, do not recreate it
- [ ] Identify whether the query requires joins — if yes, a separate domain struct is needed (`UserDetail`)
- [ ] Identify whether a third party lib is needed — if yes, check if a `pkg/` wrapper already exists

### As You Work Through Each Layer

**domain/**
- [ ] Define the entity struct — no JSON tags, no pgtype, no HTTP concerns
- [ ] Define the repository interface method for this flow
- [ ] Define the service interface method for this flow
- [ ] Add any new domain errors to `errors.go`

**repository/**
- [ ] Write the SQLC query in the `.sql` file and regenerate — do not manually edit `db/`
- [ ] Implement the repository interface method
- [ ] Write a private `toDomainX` mapper function — do not map inline
- [ ] Use `pkg/conv` for every pgtype conversion — no exceptions

**service/**
- [ ] Implement the service interface method
- [ ] All validation and business rules go here
- [ ] Depend only on domain interfaces — no concrete types from repository or pkg/

**handler/**
- [ ] Define a request DTO if the endpoint accepts a body or query params
- [ ] Define a response DTO specific to this endpoint
- [ ] Add or extend the feature route registration function in the handler layer
- [ ] Write a private `toXResponse` mapper function — do not map inline

**middleware/**
- [ ] Keep auth, RBAC, request context, request logging, and audit as separate middleware units
- [ ] Inject only required interfaces or infrastructure into middleware; do not depend on the old `Server`
- [ ] Store request-scoped values in a form that plain `context.Context` can carry forward
- [ ] Use the new logger interface in refactored middleware
- [ ] Avoid background work that retains the live `*gin.Context`
- [ ] Map domain errors to HTTP status codes here
- [ ] No business logic — if unsure, it belongs in the service

### After You Complete a Flow

- [ ] Endpoint is reachable and returns the correct response
- [ ] No pgtype outside of `repository/` and `pkg/conv/`
- [ ] No JSON tags on domain structs
- [ ] No business logic in handler
- [ ] No Gin imports in service or repository
- [ ] SQLC types do not leave the repository
- [ ] Response DTOs do not enter the service or repository
- [ ] All db→domain mapping via private `toDomainX` functions
- [ ] All domain→response mapping via private `toXResponse` functions
- [ ] Third party libs not imported directly outside of `pkg/`
- [ ] No layer imports something from the layer above it
