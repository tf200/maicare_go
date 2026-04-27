# Refactor Progress

# 2026-04-02 - Late arrival routes scaffold
Old: `api/late_arrival_router.go`, `api/late_arrival_handler.go`, `service/late_arrival/*`, `db/sqlc/late_arrival.sql.go`
New: `internal/domain/late_arrival.go`, `internal/repository/late_arrival_repo.go`, `internal/service/late_arrival_service.go`, `internal/handler/late_arrival_handler.go`, `internal/handler/late_arrival_dto.go`
Status: scaffolded, unwired
Notes: migrated late-arrival endpoint surface into the internal stack with domain-owned errors/contracts, repository SQLC mapping via `pkg/conv`, service-level shift/timezone validation, and handler-local DTO mapping

# 2026-03-31 - Handbook routes scaffold
Old: `api/handbook_router.go`, `api/handbook_handler.go`, `service/handbook/*`
New: `internal/domain/handbook.go`, `internal/repository/handbook_repo.go`, `internal/service/handbook_service.go`, `internal/handler/handbook_handler.go`, `internal/handler/handbook_dto.go`
Status: scaffolded, unwired
Notes: handbook flow now follows the internal refactor pattern: domain models/interfaces own feature contracts and errors, repository maps SQLC rows to domain models, service depends on the domain repository instead of `db.Store`, and handler uses handler-local DTOs instead of binding directly into domain structs

# 2026-03-31 - Shift swap routes scaffold
Old: `api/shift_swap_router.go`, `api/shift_swap_handler.go`, `service/schedule/shift_swap*`, `db/sqlc/shift_swap.sql.go`
New: `internal/domain/schedule.go`, `internal/repository/schedule_repo.go`, `internal/service/schedule_shift_swap.go`, `internal/handler/shift_swap_handler.go`, `internal/handler/shift_swap_dto.go`
Status: scaffolded, unwired
Notes: migrated shift swap endpoint surface and transactional swap decision flow into the internal schedule stack

# 2026-03-31 - Leave routes scaffold
Old: `api/leave_router.go`, `api/leave_handler.go`, `service/leave/*`, `db/sqlc/leave_request.sql.go`, `db/sqlc/leave_balance.sql.go`, `db/sqlc/leave_policy.sql.go`
New: `internal/domain/leave.go`, `internal/repository/leave_repo.go`, `internal/service/leave_service.go`, `internal/handler/leave_handler.go`, `internal/handler/leave_dto.go`
Status: scaffolded, unwired
Notes: migrated leave request and leave balance endpoint surface into the internal stack with equivalent permissions/routes and transactional leave decision/adjustment flows

# 2026-03-31 - Schedule routes scaffold
Old: `api/schedule_router.go`, `api/schedule_handler.go`, `service/schedule/*`
New: `internal/domain/schedule.go`, `internal/repository/schedule_repo.go`, `internal/service/schedule_service.go`, `internal/handler/schedule_handler.go`
Status: scaffolded, unwired
Notes: duplicated the schedule endpoint surface into the internal stack; route wiring is intentionally deferred

# 2026-03-30 - POST /auth/refresh

Old: `api/auth_router.go`, `api/auth_handler.go`, `service/auth/*`, `db/sqlc/session.sql.go`
New: `internal/domain/auth.go`, `internal/repository/auth_repo.go`, `internal/service/auth_service.go`, `internal/handler/auth_*`
Status: migrated, unwired
Notes: refresh token validation and access-token reissue added

# 2026-03-30 - POST /auth/logout

Old: `api/auth_router.go`, `api/auth_handler.go`, `service/auth/*`, `db/sqlc/session.sql.go`
New: `internal/domain/auth.go`, `internal/repository/auth_repo.go`, `internal/service/auth_service.go`, `internal/handler/auth_*`
Status: migrated, unwired
Notes: session deletion added with auth payload lookup in the handler

# 2026-03-30 - POST /auth/token

Old: `api/auth_router.go`, `api/auth_handler.go`, `service/auth/*`, `db/sqlc/custom_user.sql.go`, `db/sqlc/session.sql.go`
New: `internal/domain/auth.go`, `internal/repository/auth_repo.go`, `internal/service/auth_service.go`, `internal/handler/auth_*`
Status: migrated, unwired
Notes: login only, with shared token/password helpers

# 2026-03-30 - POST /auth/verify_2fa

Old: `api/auth_handler.go`, `service/auth/twofa.go`, `db/sqlc/custom_user.sql.go`
New: `internal/domain/auth.go`, `internal/repository/auth_repo.go`, `internal/service/auth_service.go`, `internal/handler/auth_*`
Status: migrated, unwired
Notes: 2FA code verification with temp token, returns final access+refresh tokens

# 2026-03-30 - POST /auth/setup_2fa

Old: `service/auth/twofa.go`, `db/sqlc/custom_user.sql.go`
New: `internal/domain/auth.go`, `internal/repository/auth_repo.go`, `internal/service/auth_service.go`, `internal/handler/auth_*`
Status: migrated, unwired
Notes: generates TOTP secret and QR code, stores temp secret

# 2026-03-30 - POST /auth/enable_2fa

Old: `service/auth/twofa.go`, `db/sqlc/custom_user.sql.go`
New: `internal/domain/auth.go`, `internal/repository/auth_repo.go`, `internal/service/auth_service.go`, `internal/handler/auth_*`
Status: migrated, unwired
Notes: validates code, enables 2FA and returns recovery codes

# 2026-03-30 - POST /auth/change_password

Old: `service/auth/auth.go`, `db/sqlc/custom_user.sql.go`
New: `internal/domain/auth.go`, `internal/repository/auth_repo.go`, `internal/service/auth_service.go`, `internal/handler/auth_*`
Status: migrated, unwired
Notes: password change with old password verification

# 2026-03-30 - GET/PUT/DELETE /locations/:id/shifts

Old: `api/shift_handler.go`, `service/organization/*`, `db/sqlc/shifts.sql.go`
New: `internal/domain/organization.go`, `internal/repository/organization_repo.go`, `internal/service/organization_service.go`, `internal/handler/organization_*`
Status: migrated, unwired
Notes: shift list, update, and delete added together

# 2026-03-30 - POST /locations/:id/shifts

Old: `api/shift_handler.go`, `service/organization/*`, `db/sqlc/shifts.sql.go`
New: `internal/domain/organization.go`, `internal/repository/organization_repo.go`, `internal/service/organization_service.go`, `internal/handler/organization_*`
Status: migrated, unwired
Notes: shift create added with max-4 slot rule

# 2026-03-30 - DELETE /locations/:id

Old: `api/location_handler.go`, `service/organization/*`, `db/sqlc/location.sql.go`
New: `internal/domain/organization.go`, `internal/repository/organization_repo.go`, `internal/service/organization_service.go`, `internal/handler/organization_*`
Status: migrated, unwired
Notes: location delete added with shared delete DTO

# 2026-03-30 - PUT /locations/:id

Old: `api/location_handler.go`, `service/organization/*`, `db/sqlc/location.sql.go`
New: `internal/domain/organization.go`, `internal/repository/organization_repo.go`, `internal/service/organization_service.go`, `internal/handler/organization_*`
Status: migrated, unwired
Notes: location update added with shared location DTO

# 2026-03-30 - GET /locations/:id

Old: `api/location_handler.go`, `service/organization/*`, `db/sqlc/location.sql.go`
New: `internal/domain/organization.go`, `internal/repository/organization_repo.go`, `internal/service/organization_service.go`, `internal/handler/organization_*`
Status: migrated, unwired
Notes: location detail added with shared location DTO

# 2026-03-30 - GET /locations

Old: `api/location_handler.go`, `service/organization/*`, `db/sqlc/location.sql.go`
New: `internal/domain/organization.go`, `internal/repository/organization_repo.go`, `internal/service/organization_service.go`, `internal/handler/organization_*`
Status: migrated, unwired
Notes: global locations list added with shared location DTO

# 2026-03-30 - POST /organisations/:id/locations

Old: `api/location_handler.go`, `service/organization/*`, `db/sqlc/location.sql.go`
New: `internal/domain/organization.go`, `internal/repository/organization_repo.go`, `internal/service/organization_service.go`, `internal/handler/organization_*`
Status: migrated, unwired
Notes: create location added under organization scope

# 2026-03-30 - DELETE /organisations/:id

Old: `api/location_handler.go`, `service/organization/*`, `db/sqlc/location.sql.go`
New: `internal/domain/organization.go`, `internal/repository/organization_repo.go`, `internal/service/organization_service.go`, `internal/handler/organization_*`
Status: migrated, unwired
Notes: delete organization added with explicit RBAC

# 2026-03-30 - PUT /organisations/:id

Old: `api/location_handler.go`, `service/organization/*`, `db/sqlc/location.sql.go`
New: `internal/domain/organization.go`, `internal/repository/organization_repo.go`, `internal/service/organization_service.go`, `internal/handler/organization_*`
Status: migrated, unwired
Notes: update organization added with shared response envelope

## 2026-03-30 - GET /organisations/count

Old: `api/location_handler.go`, `service/organization/*`, `db/sqlc/location.sql.go`
New: `internal/domain/organization.go`, `internal/repository/organization_repo.go`, `internal/service/organization_service.go`, `internal/handler/organization_*`
Status: migrated, unwired
Notes: global organization counts added

## 2026-03-30 - GET /organizations/{id}/counts

Old: `api/location_handler.go`, `service/organization/*`, `db/sqlc/location.sql.go`
New: `internal/domain/organization.go`, `internal/repository/organization_repo.go`, `internal/service/organization_service.go`, `internal/handler/organization_*`
Status: migrated, unwired
Notes: organization aggregate counts added

## 2026-03-30 - GET /organizations/{id}

Old: `api/location_handler.go`, `service/organization/*`, `db/sqlc/location.sql.go`
New: `internal/domain/organization.go`, `internal/repository/organization_repo.go`, `internal/service/organization_service.go`, `internal/handler/organization_*`
Status: migrated, unwired
Notes: added get-by-id, created/updated timestamps, and service logging

## 2026-03-30 - POST /organizations

Old: `api/location_handler.go`, `api/location_router.go`, `service/organization/*`, `db/sqlc/location.sql.go`
New: `internal/domain/organization.go`, `internal/repository/organization_repo.go`, `internal/service/organization_service.go`, `internal/handler/organization_*`
Status: migrated, unwired
Notes: create endpoint added with logger, auth, RBAC

## 2026-03-29 - GET /organizations/{id}/locations

Old: `api/location_handler.go`, `service/organization/*`, `db/query/location.sql`, `db/query/shifts.sql`
New: `internal/domain/organization.go`, `internal/repository/organization_repo.go`, `internal/service/organization_service.go`, `internal/handler/organization_handler.go`
Status: migrated, unwired
Notes: list locations with shift loading preserved

## 2026-03-29 - GET /organizations

Old: `api/location_handler.go`, `service/organization/*`, `db/sqlc/location.sql.go`
New: `internal/domain/organization.go`, `internal/repository/organization_repo.go`, `internal/service/organization_service.go`, `internal/handler/organization_handler.go`
Status: migrated, unwired
Notes: list endpoint with pagination preserved

## 2026-03-29 - Organization refactor cleanup

Completed:
- Added organization name search to the new `GET /organizations` flow
- Split organization handler DTOs and mapper functions into a dedicated handler DTO file
- Renamed the new organization service and handler types to organization-scoped names

Notes:
- This cleanup only affects the new unwired `internal/` implementation

## 2026-03-29 - Middleware foundation

New files added:
- `internal/middleware/context.go`
- `internal/middleware/response.go`
- `internal/middleware/request_context.go`
- `internal/middleware/auth.go`
- `internal/middleware/rbac.go`
- `internal/middleware/request_logging.go`
- `internal/domain/logger.go`
- `pkg/logger/logger.go`

Completed:
- Added a new middleware foundation for the refactored stack without touching the old `api/` middleware
- Added request context propagation that works through plain `context.Context`
- Added new auth, RBAC, and request logging middleware isolated from the old `Server` type
- Added the new logger interface and `pkg/logger` implementation for refactored flows
- Added a new domain token contract and `pkg/jwt` implementation so refactored middleware does not depend on the legacy `token/` package interface

Not done:
- Audit middleware is intentionally deferred to a second pass because the old audit flow needs a safer redesign
