 Refactor the next endpoint, and only that endpoint, into the new parallel architecture by following the existing reference implementation already added to this repo.

  Reference implementation to imitate:

  - internal/domain/organization.go
  - internal/domain/organization_repository.go
  - internal/domain/organization_service.go
  - internal/repository/organization_repo.go
  - internal/service/organization_service.go
  - internal/handler/organization_handler.go
  - pkg/conv/pg.go
  - docs/refactor_progress.md

  Task:
  Refactor this endpoint only:

  - GET /organizations

  Old code to study:

  - <old handler file>
  - <old service interface file>
  - <old service implementation file>
  - <sqlc query file(s)>

  Hard constraints:

  - Do not modify main.go
  - Do not wire the new implementation into the running app
  - Keep the old implementation untouched unless a tiny shared extraction is strictly necessary
  - Do not make old api/ or service/ depend on new internal/ packages
  - Do not make new internal/ packages depend on old api/ or service/ packages
  - Follow the same layering pattern as the existing organization refactor
  - Repository is the only layer allowed to import SQLC types
  - pgtype conversions belong in pkg/conv if reusable, otherwise keep them private in repository mappers
  - Service must not depend on Gin
  - Handler must own request/response DTOs and private mapping functions
  - Do not import util in the new flow unless absolutely necessary
  - Avoid duplication where possible
  - After finishing, append a short entry to docs/refactor_progress.md

  Deliverables:

  1. Create only the new files needed under:

  - internal/domain/
  - internal/repository/
  - internal/service/
  - internal/handler/
  - pkg/conv/ only if needed

  2. Keep the implementation limited to this one endpoint flow
  3. At the end, report:

  - files created
  - files studied
  - whether main.go changed
  - any remaining gaps or compromises

  Implementation rules:

  - Domain defines the model and repository/service interface
  - Repository maps db/sqlc rows to domain via private mapper functions
  - Service orchestrates and returns domain models only
  - Handler maps domain models to response DTOs and returns the standard response envelope
  - If pagination is needed, keep it in the handler layer like the organization reference
  - Match existing behavior first; do not optimize behavior unless required for correctness

  Validation:

  - Run GOCACHE=/tmp/go-build go test ./internal/... ./pkg/conv
  - If the new flow needs other packages to compile, run the smallest targeted go test command possible
  - Include the validation result in the final response

  Shorter version
  If you want a compact reusable prompt:

  Follow the existing internal/organization refactor as the template and refactor only <endpoint> into the new parallel architecture. Do not modify main.go, do not wire the new code, do not let new code depend on old api/service,
  keep SQLC and pgtype in repository or pkg/conv, keep Gin out of service, let handler own request/response DTOs, avoid util, update docs/refactor_progress.md, and validate with GOCACHE=/tmp/go-build go test ./internal/... ./pkg/
  conv.

  Best way to use it
  Fill in these fields each time for one endpoint at a time:

  - endpoint
  - old handler file
  - old service files
  - SQL files
