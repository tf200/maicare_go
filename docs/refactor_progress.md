# Refactor Progress

## 2026-03-29 - GET /organizations/{id}/locations

Old flow studied:
- `api/location_handler.go`
- `service/organization/service.go`
- `service/organization/location.go`
- `service/organization/shifts.go`
- `db/query/location.sql`
- `db/query/shifts.sql`

New files added:
- `internal/domain/organization.go`
- `internal/domain/organization_repository.go`
- `internal/domain/organization_service.go`
- `internal/repository/organization_repo.go`
- `internal/service/organization_service.go`
- `internal/handler/organization_handler.go`
- `internal/handler/response.go`
- `pkg/conv/pg.go`

Completed:
- Added a parallel domain, repository, service, and handler flow for organization location listing
- Kept SQLC and `pgtype` usage inside the repository and `pkg/conv`
- Kept Gin out of the new service layer
- Kept `main.go` unchanged and left the new flow unwired

Notes:
- Pagination response shaping still uses the existing `pagination` package in the handler layer
- The repository keeps the current N+1 shift loading behavior to preserve the old endpoint behavior for this first slice

## 2026-03-29 - GET /organizations

Old flow studied:
- `api/location_handler.go`
- `service/organization/service.go`
- `service/organization/organisation.go`
- `service/organization/organisation_dto.go`
- `db/sqlc/location.sql.go`

New files added:
- `internal/domain/organization.go`
- `internal/domain/organization_repository.go`
- `internal/domain/organization_service.go`
- `internal/repository/organization_repo.go`
- `internal/service/organization_service.go`
- `internal/handler/organization_handler.go`

Completed:
- Added the organizations list endpoint to the existing organizations domain, repository, service, and handler flow
- Kept SQLC access inside the repository and preserved the existing pagination envelope in the handler
- Kept Gin out of the new service layer
- Kept `main.go` unchanged and left the new flow unwired

Notes:
- The list endpoint lives beside the other organizations methods instead of a separate list-only stack

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
