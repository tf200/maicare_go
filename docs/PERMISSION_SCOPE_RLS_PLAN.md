# Permission Scope and Row-Level Security Plan

Last updated: 2026-08-20

Status: Phase 9 in progress; Groups A, B, and C completed

## Purpose

This document is the canonical implementation plan for permission scope and PostgreSQL row-level security (RLS). It is intended to preserve context across development sessions and to track progress as the work is completed incrementally.

The application handles healthcare information and must support access controls appropriate for NEN 7510. The implementation must follow least privilege and need-to-know principles, fail closed when authorization context is missing, and provide evidence through tests and audit records.

Do not implement this plan as one large change. Complete and verify one phase before starting the next phase.

## Progress Legend

- `[ ]` Not started
- `[-]` In progress
- `[x]` Completed and verified
- `[!]` Blocked or requires a decision

When completing a phase:

1. Mark its checklist items complete.
2. Record the verification commands and results.
3. Add a dated entry to the Progress Log.
4. Record any changed decision in the Decision Log.
5. Do not mark a phase complete until its acceptance criteria pass.

## Agreed Terminology

Use the generic name `scope`, not `client_scope`, in database columns, Go models, and JSON fields.

A permission answers:

> What action may the user perform?

Examples:

- `CLIENT.VIEW`
- `CLIENT.UPDATE`
- `CLIENT.MEDICATION.VIEW`
- `ROLES.VIEW`

A scope answers:

> Over which client records may the user perform that action?

Initial scope values:

- `assigned`: only clients actively assigned to the employee
- `all`: all clients in the system
- `NULL`: scope does not apply to this permission

`NULL` must not mean `all`. A scoped permission with a missing scope must fail closed.

## Target Authorization Rules

Permissions are assigned to roles. Scope belongs to the role-permission grant, not globally to the role.

Examples:

| Role | Permission | Scope |
|---|---|---|
| Administrator | `CLIENT.VIEW` | `all` |
| Administrator | `CLIENT.UPDATE` | `all` |
| Coordinator | `CLIENT.VIEW` | `assigned` |
| Coordinator | `CLIENT.UPDATE` | `assigned` |
| Manager | `CLIENT.VIEW` | `all` |
| Manager | `CLIENT.MEDICATION.VIEW` | `assigned` |
| Administrator | `ROLES.VIEW` | `NULL` |

This supports different scopes for different categories of client data. That distinction is important for least privilege: access to basic client details must not automatically grant access to medical, incident, financial, or document data.

Authorization has two layers:

1. The Go HTTP middleware checks whether the authenticated user has the required permission at all.
2. PostgreSQL RLS checks whether that permission and its effective scope allow access to each particular client row.

RLS is the final data-protection boundary. HTTP middleware improves API behavior but must not be the only authorization control.

## Intended Request Lifecycle

1. Receive the HTTP request.
2. Validate the JWT access token.
3. Extract the authenticated user ID, employee ID, and session ID.
4. Check the required route permission using current database authorization data.
5. Return `403 Forbidden` if the permission is absent.
6. Begin a database transaction for protected data access.
7. Set transaction-local `myapp.current_user_id` and `myapp.current_employee_id`.
8. Execute normal SQL without manually duplicating assignment filters.
9. RLS resolves the effective permission and scope for the requested operation.
10. For `assigned`, RLS verifies an active assignment for that client.
11. For `all`, RLS allows the row globally.
12. Commit or roll back the transaction.
13. Return only rows PostgreSQL allowed.

Roles, permissions, and scopes should not be trusted from JWT claims. The JWT identifies the actor; current authorization is loaded from the database so changes can take effect without waiting for token expiration.

## Current State

### Database RBAC

Defined in `db/migrations/000001_init.up.sql`:

- `roles` at approximately lines 475-479
- `permissions` at approximately lines 482-489
- `role_permissions` at approximately lines 493-500
- `user_roles` follows `custom_user`

Current limitations:

- `role_permissions` stores scope per grant.
- Permissions declare whether scope applies.
- Roles are the sole source of user permissions; direct per-user overrides are not supported.
- `user_roles.user_id` is the primary key, so one user can currently have only one role.
- Effective-permission queries return only permission ID and name.

### Current HTTP Authorization

Relevant files:

- `pkg/jwt/jwt.go`
- `internal/adapters/jwt_token_maker.go`
- `internal/middleware/auth.go`
- `internal/middleware/rbac.go`
- `internal/middleware/context.go`
- `internal/ctxkeys/ctxkeys.go`
- `internal/app/app.go`

Current behavior:

- JWT claims identify the user, employee, and session.
- `internal/middleware/rbac.go` checks a permission as a boolean.
- Effective permission resolution is wired in `internal/app/app.go`.
- Scope is not represented in middleware or context.
- JWT claims currently do not contain roles or permissions; keep this property.

### Current Role Management

Relevant files:

- `internal/domain/permission.go`
- `internal/domain/role.go`
- `internal/repository/role_repo.go`
- `internal/service/role_service.go`
- `internal/handler/role_dto.go`
- `internal/handler/role_handler.go`
- `internal/handler/role_routes.go`
- `db/query/roles.sql`
- `cmd/roles/main.go`
- `sqlc.yaml`
- Generated files under `db/sqlc/`

Current limitations:

- Role-permission requests and responses include grant scope.
- Replacing role permissions is transactional.
- Direct user permission mutation has been removed.
- The role synchronization command is additive and does not remove obsolete grants.

### Current RLS

Defined near the end of `db/migrations/000001_init.up.sql`:

- `get_current_employee_id()`
- `is_admin()`
- `is_coordinator()`
- `is_assigned_coordinator(client_id)`
- `apply_client_rls(table_name, client_id_col)`

Current policy behavior:

| Operation | Current behavior |
|---|---|
| `SELECT` | Any administrator or coordinator can see rows |
| `INSERT` | Any administrator or coordinator can insert rows |
| `UPDATE` | Administrator or assigned coordinator |
| `DELETE` | Administrator or assigned coordinator |

Current problems:

- Coordinator reads are not restricted to assigned clients.
- Coordinator inserts are not restricted by assignment.
- Policies use hard-coded role names instead of permissions.
- Custom future roles cannot participate correctly.
- Role grants are the sole permission source for future RLS policies.
- Different data categories cannot require different permissions.
- RLS is enabled but not forced.
- Runtime PostgreSQL ownership and `BYPASSRLS` behavior are not defined in migrations.

### Current Database Identity Handling

The standard transaction wrapper is `db/sqlc/store.go` and sets `myapp.current_employee_id` in `Store.ExecTx`.

Other relevant transaction paths include:

- `internal/repository/contract_repo.go`
- `internal/service/invoice_service.go`
- `internal/service/event_service.go`
- `internal/repository/notification_repo.go`
- Repositories that call generated sqlc methods directly through the pool

Current problems:

- Direct pool queries do not receive transaction-local identity.
- Some manual transactions set the employee ID and some do not.
- Workers do not naturally have an authenticated employee context.
- `Store.ExecTx` prints the employee ID to standard output.
- Session-level settings must not be used on pooled connections because identity could leak between requests.

### Current Client Assignment

The `assigned_employee` table is defined in `db/migrations/000001_init.up.sql` near lines 2110-2123.

Relevant application files include:

- `db/query/client_network.sql`
- `internal/domain/client.go`
- `internal/repository/client_repo.go`
- `internal/service/client_service.go`
- `internal/handler/client_dto.go`
- `internal/handler/client_handler.go`

Current RLS recognizes only assignments where `assigned_employee.role = 'coordinator'`. The future meaning of `assigned` must be made explicit and must account for assignment start and end dates.

## Phase 0: Confirm Design Decisions

Status: `[-]` Partially resolved; remaining decisions are required by later phases

- [x] Confirm whether `000001_init.up.sql` has been used in any deployed or shared database.
- [x] Confirm whether implementation requires a new migration or may modify the initial migration.
- [x] Confirm whether `assigned` means any active `assigned_employee` record or only selected assignment types.
- [x] Confirm how an assignment ends; the current table has `start_date` but no `end_date`.
- [x] Confirm whether `all` means all clients globally or all clients in the actor's organization.
- [ ] Confirm whether users will remain limited to one role or may receive multiple roles later.
- [x] Confirm that direct per-user permission overrides are not supported.
- [x] Confirm behavior for background workers and trusted system operations.
- [x] Confirm whether `CLIENT.CREATE` is unscoped initially, because a client does not have an assignment before creation.

Recommended decisions:

- Create a new migration if any non-disposable database has run migration `000001`.
- Treat `all` as all clients globally.
- Treat any active assignment as `assigned`; permissions determine which client data categories are accessible.
- Add an optional assignment end date.
- Keep roles as the sole source of permissions and scope.
- Treat `CLIENT.CREATE` as unscoped until an organization-aware creation rule is designed.
- Design for multiple roles before enabling them; do not accidentally combine a permission from one role with an unrelated broader scope from another role.

Acceptance criteria:

- [ ] Every item above has a recorded decision in the Decision Log.
- [ ] Scope semantics can be explained without referring to a role name.

## Phase 1: Add Scope Representation To The Database

Status: `[x]` Completed and verified on 2026-08-16

Goal: represent scope without changing current authorization behavior or RLS policies.

Likely affected files:

- `db/migrations/000001_init.up.sql`, only if safe to modify
- `db/migrations/000001_init.down.sql`, only if safe to modify
- Otherwise new files such as `db/migrations/000002_permission_scope.up.sql` and `db/migrations/000002_permission_scope.down.sql`
- `db/sqlc/models.go`, regenerated rather than manually edited

Planned schema:

```sql
CREATE TYPE permission_scope_enum AS ENUM ('assigned', 'all');

CREATE TABLE permissions (
    -- Existing fields omitted
    is_scoped BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE role_permissions (
    -- Existing fields omitted
    scope permission_scope_enum NULL
);
```

Required rules:

- Scoped permissions require a non-null scope when granted.
- Unscoped permissions require `scope IS NULL`.
- Missing scope never defaults to `all`.
- Existing production grants require an explicit and reviewed migration strategy.

Implementation decision: deferred constraint triggers validate both sides of the relationship. `role_permissions_validate_scope` validates grant writes, and `permissions_validate_grant_scopes` validates permission metadata writes. Deferring validation until transaction commit allows a future transaction to change `permissions.is_scoped` and its existing grants atomically while rejecting any invalid final state.

Checklist:

- [x] Add the scope enum.
- [x] Add `permissions.is_scoped`.
- [x] Add `role_permissions.scope`.
- [x] Add database-level validation of scoped versus unscoped grants.
- [x] Confirm no existing-grant data migration is needed because the database is disposable during early development.
- [x] Add matching cleanup to the down migration.
- [x] Regenerate sqlc models.
- [x] Verify migration up and down on a disposable PostgreSQL database.
- [x] Verify no RLS behavior changed in this phase.

Acceptance criteria:

- [x] The database stores `assigned`, `all`, or `NULL` correctly.
- [x] Invalid permission/scope combinations are rejected.
- [x] Existing authorization behavior is unchanged.

## Phase 2: Classify Permissions As Scoped Or Unscoped

Status: `[x]` Completed and verified on 2026-08-16

Goal: make permission metadata the source of truth for whether scope applies.

Affected files:

- `internal/domain/permission.go`
- `cmd/roles/main.go`
- `db/query/roles.sql`
- Generated files under `db/sqlc/`

Implemented classification principles:

- Every registered `CLIENT.*` permission operating on an existing client is scoped.
- `CLIENT.CREATE` is the only unscoped `CLIENT.*` permission.
- Permissions outside the `CLIENT.*` namespace remain unscoped until their resource area is reviewed.
- Medical, incident, document, evaluation, and basic client information remain separately permissioned.

Examples likely to be scoped:

- `CLIENT.VIEW`
- `CLIENT.UPDATE`
- `CLIENT.DELETE`
- `CLIENT.CARE_PLAN.*`
- `CLIENT.DOCUMENTS.*`
- `CLIENT.INCIDENT.*`
- `CLIENT.DIAGNOSIS.*`
- `CLIENT.MEDICATION.*`
- `CLIENT.EMERGENCY_CONTACT.*`
- `CLIENT.PROGRESS_REPORT.*`
- `CLIENT.EVALUATION.CREATE`
- `CLIENT.EVALUATION.VIEW`

Examples likely to remain unscoped:

- `ROLES.*`
- `PERMISSIONS.*`
- General settings permissions
- Authentication and self-service operations
- `CLIENT.CREATE`, initially

Checklist:

- [x] Extend permission metadata with `IsScoped`.
- [x] Review every permission in `internal/domain/permission.go`.
- [x] Classify all current `CLIENT.*` permissions using the documented namespace rule.
- [x] Rename active evaluation permissions to `CLIENT.EVALUATION.CREATE` and `CLIENT.EVALUATION.VIEW`.
- [x] Remove the unused `EVALUATION.DELETE` permission instead of creating an unused client equivalent.
- [x] Update evaluation routes to require the dedicated evaluation permissions.
- [x] Update `cmd/roles/main.go` to persist `is_scoped`.
- [x] Assign `all` to scoped administrator grants and `assigned` to scoped coordinator grants.
- [x] Ensure permission synchronization updates changed metadata and grant scopes.
- [x] Add automated classification and evaluation-registration tests.

Acceptance criteria:

- [x] Client permission classification follows one explicit registry rule with a tested `CLIENT.CREATE` exception.
- [x] Sensitive client-data categories are separately classified.
- [x] The database and Go registry agree.

## Phase 3: Update Role-Permission Management

Status: `[x]` Completed and verified on 2026-08-17

Goal: allow administrators to assign a permission and optional scope to a role.

Affected files:

- `db/query/roles.sql`
- `internal/domain/role.go`
- `internal/repository/role_repo.go`
- `internal/service/role_service.go`
- `internal/handler/role_dto.go`
- `internal/handler/role_handler.go`
- `internal/handler/role_routes.go`
- `internal/app/app.go`, if interface wiring changes
- Generated files under `db/sqlc/`

Target request shape:

```json
{
  "permissions": [
    {
      "permission_id": "00000000-0000-0000-0000-000000000001",
      "scope": "assigned"
    },
    {
      "permission_id": "00000000-0000-0000-0000-000000000002",
      "scope": null
    }
  ]
}
```

Target response fields:

```json
{
  "permission_id": "00000000-0000-0000-0000-000000000001",
  "permission_name": "CLIENT.VIEW",
  "is_scoped": true,
  "scope": "assigned"
}
```

Checklist:

- [x] Add a domain scope type with strict values.
- [x] Replace permission-ID-only mutation models with permission grant models.
- [x] Return `is_scoped` and `scope` from role detail APIs.
- [x] Return `is_scoped` from the permission catalog API.
- [x] Update sqlc role-permission queries to read and write scope.
- [x] Validate unknown scope values as client errors.
- [x] Validate missing scope on scoped permissions.
- [x] Validate non-null scope on unscoped permissions.
- [x] Reject unknown and duplicate permission IDs before replacement.
- [x] Make role-permission replacement one transaction.
- [x] Preserve the old grants if replacement fails.
- [x] Regenerate sqlc code.
- [x] Update the API contract in this plan.

Acceptance criteria:

- [x] A role can contain mixed `assigned`, `all`, and unscoped grants.
- [x] Invalid combinations return a clear 4xx response.
- [x] Replacement is atomic.
- [x] API reads return exactly what was saved.

## Phase 4: Remove Legacy User Permission Overrides

Status: `[x]` Completed and verified on 2026-08-17

Goal: make role grants the sole source of user permissions and eliminate the legacy bypass path.

Affected files:

- Database migration files
- `db/query/roles.sql`
- `internal/domain/role.go`
- `internal/repository/role_repo.go`
- `internal/service/role_service.go`
- `internal/handler/role_dto.go`
- Generated files under `db/sqlc/`

Checklist:

- [x] Remove `user_permission_overrides` and its effect enum from the schema.
- [x] Remove override SQL, generated models, domain models, and repository/service methods.
- [x] Remove the direct user-permission mutation endpoint and DTOs.
- [x] Remove the obsolete `PERMISSIONS.GRANT` permission.
- [x] Resolve effective permissions exclusively from the assigned role.
- [x] Remove override fields from role-management responses.
- [x] Update employee-profile permission aggregation to use role grants only.

Acceptance criteria:

- [x] No database or API path can store a direct user permission.
- [x] Middleware and employee-profile permissions come only from role grants.
- [x] Permission changes are managed through role grants and role assignment.

## Phase 5: Resolve Effective Permissions And Scope

Status: `[x]` Completed and verified on 2026-08-17

Goal: calculate the user's current effective permission grant and scope from database data.

Affected files:

- `db/query/roles.sql`
- `internal/domain/role.go`
- `internal/repository/role_repo.go`
- `internal/repository/handbook_repo.go`
- `internal/service/role_service.go`
- `internal/middleware/rbac.go`
- `internal/app/app.go`
- Generated files under `db/sqlc/`

Expected result shape:

```json
{
  "permission_name": "CLIENT.VIEW",
  "is_scoped": true,
  "scope": "assigned"
}
```

Rules:

- Unscoped effective permissions return `scope = NULL`.
- Missing scope on a scoped grant results in no usable grant.
- If multiple roles are introduced, scope combination must be explicit and tested. Never combine a permission from one role with `all` scope from an unrelated role that does not grant that permission.

Checklist:

- [x] Update role-derived permission queries to include scope.
- [x] Update effective permission queries to include scope.
- [x] Update direct permission checks.
- [x] Avoid loading all permissions merely to check one key where possible.
- [x] Keep role and scope data out of JWT claims.
- [x] Add effective-permission tests for role grants.

Acceptance criteria:

- [x] HTTP middleware still returns a clear `403` when permission is absent.
- [x] Effective scope matches the specific grant that provides the permission.
- [x] Database changes to grants affect new requests without issuing a new JWT.

## Phase 6: Centralize Protected Database Execution

Status: `[x]` Completed and verified on 2026-08-19

Goal: guarantee that every protected query executes with trusted transaction-local actor identity.

Primary affected files:

- `db/sqlc/store.go`
- `internal/ctxkeys/ctxkeys.go`
- `internal/middleware/auth.go`
- `internal/middleware/context.go`
- `internal/repository/client_repo.go`
- `internal/repository/incident_repo.go`
- `internal/repository/contract_repo.go`
- `internal/service/invoice_service.go`
- `internal/service/event_service.go`
- `internal/repository/notification_repo.go`
- Other repositories that query RLS-protected tables directly
- Worker entry points under `internal/worker/`

Required transaction-local settings:

```text
myapp.current_user_id
myapp.current_employee_id
```

Rules:

- Set identity only after JWT authentication.
- Set identity with transaction-local `set_config(..., true)`.
- Never set request identity at connection/session level on a pool.
- Missing identity must fail closed.
- Background work must use an explicit, documented system actor or service authorization path.
- Remove the employee-ID standard-output print from `Store.ExecTx`.

Checklist:

- [x] Add user ID to the standard internal request context if not already accessible there.
- [x] Define one authenticated transaction entry point.
- [x] Set both user and employee IDs in that transaction.
- [x] Inventory all direct pool access to protected tables.
- [x] Convert protected direct queries in the inventoried client, incident, intake, registration, contract, invoice, event-attendee, and worker paths to the standard transaction path.
- [x] Convert identified contract, invoice, event, and notification manual transactions to the shared initialization method.
- [x] Define worker behavior explicitly.
- [x] Test pooled connection reuse for identity leakage.
- [x] Test missing identity behavior.

Acceptance criteria:

- [x] No inventoried protected repository path can accidentally omit actor identity.
- [x] Identity does not leak to a later request on the same pooled connection.
- [x] Workers do not depend on table-owner RLS bypass.

## Phase 7: Add General Database Authorization Functions

Status: `[x]` Completed and verified on 2026-08-19

Goal: replace role-name-specific logic with permission-and-scope checks while leaving existing policies in place until the pilot phase.

Likely affected files:

- A new migration file, or initial migration only if confirmed safe
- Matching down migration

Conceptual functions:

```sql
get_current_user_id()
get_current_employee_id()
has_permission(permission_name)
get_permission_scope(permission_name)
is_assigned_to_client(client_id)
can_access_client(client_id, permission_name)
```

`can_access_client` must check:

1. Authenticated database request identity exists.
2. User is active.
3. Employee is active where applicable.
4. The user's assigned role has the requested permission.
5. The permission is scoped.
6. `all` grants global client access.
7. `assigned` has an active assignment for the client.
8. Unknown or invalid state returns false.

Security requirements for `SECURITY DEFINER` functions:

- Use a dedicated non-login owner where possible.
- Set a fixed safe `search_path`.
- Schema-qualify referenced objects.
- Revoke execution from `PUBLIC`.
- Grant only required execution rights to the runtime role.
- Avoid untrusted dynamic SQL.
- Avoid RLS recursion.

Checklist:

- [x] Implement identity getter functions.
- [x] Implement effective-permission lookup.
- [x] Implement scope lookup.
- [x] Implement active-assignment lookup.
- [x] Implement `can_access_client`.
- [x] Fix function search paths, schema-qualify objects, and revoke execution from `PUBLIC`; functions remain owned by the migration role until Phase 10 introduces separate database roles.
- [x] Add direct SQL tests for every allow and deny path.
- [x] Confirm missing context returns false rather than raising an information-leaking error.

Acceptance criteria:

- [x] No helper checks a business role name such as `admin` or `coordinator`.
- [x] Permission and scope changes affect authorization immediately on the next transaction.
- [x] Invalid or missing context denies access.

## Phase 8: Convert `client_details` As The RLS Pilot

Status: `[x]` Completed and verified on 2026-08-19

Goal: prove the new model on the main client table before changing all related tables.

Affected files:

- Database migration files
- `db/query/client.sql`
- Client repository/service paths if transaction handling must change
- RLS integration test files to be added

Operation mapping:

| SQL operation | Permission |
|---|---|
| `SELECT` | `CLIENT.VIEW` |
| `INSERT` | `CLIENT.CREATE` |
| `UPDATE` | `CLIENT.UPDATE` |
| `DELETE` | `CLIENT.DELETE` |

Notes:

- `SELECT`, `UPDATE`, and `DELETE` can evaluate an existing client ID.
- `INSERT` requires separate design because a new client cannot already be assigned. Initially treat `CLIENT.CREATE` as unscoped and enforce mandatory organization placement separately.
- A denied single-client read should normally appear as `404 Not Found` to avoid confirming that an inaccessible client exists.
- `UPDATE USING` must check access to the existing row.
- `UPDATE WITH CHECK` must check the resulting row and prevent moving data across unauthorized ownership boundaries.
- PostgreSQL also applies select visibility to update targets, so `CLIENT.UPDATE` requires `CLIENT.VIEW` for the same row.

Checklist:

- [x] Replace role-name policies on `client_details`.
- [x] Add separate policies for select, insert, update, and delete.
- [x] Verify list and aggregate queries filter rows automatically.
- [x] Verify single-record queries hide unauthorized clients.
- [x] Verify update and delete denial.
- [x] Verify `all` behavior.
- [x] Verify assigned behavior.
- [x] Verify missing permission behavior.
- [x] Verify missing identity behavior.
- [x] Verify inactive and mismatched actor denial.
- [x] Verify global `all` behavior.

Acceptance criteria:

- [x] An assigned actor sees only assigned clients.
- [x] An `all` grant sees all clients in the system.
- [x] A user without `CLIENT.VIEW` sees no client rows.
- [x] A non-owner, non-`BYPASSRLS` runtime role cannot bypass the forced policy.
- [x] Existing client reads and updates return policy-filtered results; client creation was adapted for safe `INSERT ... RETURNING` behavior.

## Phase 9: Convert Related Tables In Controlled Groups

Status: `[~]` In progress; Groups A, B, and C completed

Goal: apply permission-specific RLS to all client-owned information without one high-risk migration.

Each group requires:

1. A table-to-client ownership map.
2. A SQL-operation-to-permission map.
3. Correct transaction context in all repository paths.
4. RLS policies.
5. Integration tests.
6. Audit classification where access is sensitive.

### Group A: Client Network And Assignments

Likely tables and files:

- `assigned_employee`
- `client_emergency_contact`
- `db/query/client_network.sql`
- Client domain, repository, service, handler, and DTO files

Permissions include `CLIENT.INVOLVED_EMPLOYEE.*` and `CLIENT.EMERGENCY_CONTACT.*`.

- [x] Map operations.
- [x] Implement policies.
- [x] Test assignment-management edge cases.
- [x] Prevent users from granting themselves access through assignment changes.

Assignment reads use normal `assigned`/`all` scope. Assignment create, update, and delete require the corresponding `CLIENT.INVOLVED_EMPLOYEE.*` permission with `all` scope because assignment rows are themselves authorization inputs. Emergency contacts contain sensitive personal contact and disclosure-preference data; assignment rows are security-sensitive access-control metadata. Mutation audit coverage remains part of the later auditing phase.

### Group B: Progress Reports And AI Reports

Likely tables and files:

- `progress_report`
- `ai_generated_reports`
- Related client queries and service methods

Permissions include `CLIENT.PROGRESS_REPORT.*` and `CLIENT.AI_PROGRESS_REPORT.*`.

- [x] Map operations.
- [x] Implement policies.
- [x] Test read, create, update, delete, generate, and confirm flows.

AI generation reads progress-report source text under `CLIENT.AI_PROGRESS_REPORT.GENERATE` but does not persist an AI report. Confirmation is the only AI-report insert and uses `CLIENT.AI_PROGRESS_REPORT.CONFIRM`. AI reports have no update or delete operation or permission, so no such policies exist. Report text is sensitive care data; mutation audit coverage remains part of the later auditing phase.

Progress-report UPDATE and DELETE also require VIEW visibility for the target row because PostgreSQL applies SELECT policies while locating mutation targets. CREATE and AI confirmation remain independently usable through transaction-bound creation contexts. Report authorship remains caller-supplied existing business behavior and requires a separate field-level authorization decision rather than an implicit RLS change.

### Group C: Medical Data

Likely tables:

- `client_diagnosis`
- `client_medication_order`
- Related supporting tables

Permissions include `CLIENT.DIAGNOSIS.*` and `CLIENT.MEDICATION.*`.

- [x] Map operations.
- [x] Implement policies.
- [x] Add stricter access and audit tests for medical data.

Diagnosis and medication data require their own permissions; general `CLIENT.VIEW` does not expose either category. The combined overview requires both VIEW permissions. CREATE remains independent through actor-bound creation context, while UPDATE and DELETE also require category VIEW visibility under PostgreSQL mutation semantics. Creator/updater attribution comes from the transaction actor rather than request parameters.

Medication-to-diagnosis ownership is enforced by a composite client foreign key. `source_attachment_uuid` remains a bare attachment reference because `attachment_file` has no client ownership; validating that association is explicitly deferred to Group E documents rather than inventing an unverifiable policy here.

Deleting a diagnosis is restricted while medication orders reference it. Unlinking or deleting those orders requires the medication category's permissions and audit path, preventing diagnosis deletion from silently mutating medication data.

### Group D: Incidents

Likely tables and files:

- `incident`
- `db/query/incident.sql` or the applicable incident query files
- `internal/domain/incident.go`
- `internal/repository/incident_repo.go`
- `internal/service/incident_service.go`
- `internal/handler/incident_*`

Permissions include `CLIENT.INCIDENT.*` and any top-level incident permissions whose semantics must be reconciled.

- [x] Resolve duplicate or overlapping incident permission meanings.
- [x] Map operations.
- [x] Implement policies.
- [x] Test confirmation and file access.

### Group E: Documents, Care Plans, Goals, And Evaluations

Likely tables:

- `client_documents`
- `client_goals`
- `client_goal_evaluations`
- `client_goal_evaluation_items`
- `intake_topic_assessments`
- Care-plan-related tables

- [ ] Map operations and nested client ownership.
- [ ] Prefer direct indexed `client_id` where justified.
- [ ] Implement policies.
- [ ] Test nested records and helper functions.

### Group F: Contracts And Invoices

Likely tables and files:

- `contract`
- `client_agreement`
- `provision`
- `invoice`
- Invoice lines, payments, run items, calendar links, and billed events
- `internal/repository/contract_repo.go`
- `internal/service/invoice_service.go`

- [ ] Map each nested table to a client.
- [ ] Decide whether direct `client_id` should be added to deeply nested tables.
- [ ] Implement permission-specific policies.
- [ ] Test financial separation from general client access.

### Group G: Intake And Registration

Likely tables:

- `registration_form`
- `intake_forms`
- `intake_topic_assessments`
- `youth_care_intake`
- Related declarations and agreements

- [ ] Define when a pre-client record becomes client-owned.
- [ ] Define access before client creation.
- [ ] Replace nested helper policies carefully.
- [ ] Test public registration/intake routes remain functional without exposing protected records.

Acceptance criteria for Phase 9:

- [ ] Every RLS-protected table has a documented ownership and permission mapping.
- [ ] No policy uses hard-coded application role names.
- [ ] Sensitive categories require their own permissions.
- [ ] Nested-table policies are tested for correctness and performance.

## Phase 10: PostgreSQL Runtime Hardening

Status: `[ ]` Not started

Goal: ensure the application cannot bypass RLS accidentally or through excessive database privileges.

Likely affected files:

- Database migrations or deployment SQL
- `docker-compose.yml`
- `docker-compose.dev.yml`
- Environment documentation such as `.env.example`
- Deployment configuration outside this repository, if applicable

Required role separation:

- Migration/owner role owns schema objects.
- Runtime application role is `NOSUPERUSER` and `NOBYPASSRLS`.
- Runtime application role does not own protected tables.
- Background/service operations use an explicit controlled path.

Checklist:

- [ ] Define migration database role.
- [ ] Define runtime application database role.
- [ ] Apply least-privilege grants.
- [ ] Confirm runtime role is not a table owner.
- [ ] Confirm runtime role has no `BYPASSRLS`.
- [ ] Decide where `FORCE ROW LEVEL SECURITY` is appropriate.
- [ ] Test policies using the real runtime role, not a superuser.
- [ ] Document local, CI, staging, and production setup.

Acceptance criteria:

- [ ] Runtime queries cannot bypass RLS.
- [ ] Migrations can still be applied through a separate controlled role.
- [ ] Automated tests prove role configuration rather than assuming it.

## Phase 11: NEN 7510 Access-Control And Audit Controls

Status: `[ ]` Not started

Goal: support least privilege, accountability, access review, and auditable authorization changes.

Relevant existing audit implementation:

- Audit table and triggers in `db/migrations/000001_init.up.sql`
- `internal/domain/audit.go`
- `internal/audit/audit.go`
- Any audit query/repository files under `db/query/` and `internal/`

Required controls:

- Record role assignment and removal.
- Record role-permission and scope changes.
- Record assignment start and end.
- Record sensitive client-data reads where required by the audit policy.
- Record sensitive writes and actor identity.
- Preserve append-only audit behavior.
- Support periodic access review reports.
- Ensure inactive users and employees lose access immediately.
- Design a controlled and audited break-glass process if emergency access is required.
- Define retention and access rules for authorization audit records.

Checklist:

- [ ] Create an audit-event matrix.
- [ ] Add missing authorization-management events.
- [ ] Add required sensitive-read events.
- [ ] Include actor, subject, client, permission, scope, reason, request, session, and timestamp where applicable.
- [ ] Prevent ordinary users from changing audit records.
- [ ] Add access-review queries or reports.
- [ ] Document operational review responsibilities.

Acceptance criteria:

- [ ] It is possible to determine who changed access, what changed, and when.
- [ ] It is possible to explain why an employee had access to a client at a given time.
- [ ] Sensitive access produces the required audit evidence.

## Phase 12: Test Matrix And Rollout

Status: `[ ]` Not started

Goal: verify fail-closed behavior and roll out without exposing or unexpectedly hiding production data.

There are currently no discovered Go `*_test.go` files. This work requires PostgreSQL integration tests, not only mocked service tests.

Minimum authorization actors:

- Administrator with `all`
- Coordinator with `assigned`
- Assigned employee without the requested permission
- Unassigned employee with an `assigned` permission
- User whose role lacks the requested permission
- Inactive user
- Inactive employee
- Missing request identity
- Background/system actor
- Migration database role
- Runtime database role

Minimum operations:

- List rows
- Get one row
- Insert
- Update
- Delete
- Change a row's client ownership
- Read nested client data
- Access medical data
- Access incident data
- Access financial data
- Start and end an assignment
- Change role scope during an active session

Minimum infrastructure cases:

- Connection-pool reuse
- Transaction rollback
- Missing transaction context
- Concurrent permission change
- Runtime role ownership and bypass checks
- RLS helper recursion and query performance

Checklist:

- [ ] Add integration-test database setup.
- [ ] Seed deterministic authorization fixtures.
- [ ] Test every actor/operation combination required by each phase.
- [ ] Capture query plans for important list queries.
- [ ] Add indexes for permission and assignment lookups where demonstrated necessary.
- [ ] Run tests with the restricted runtime role.
- [ ] Define rollback steps for each production migration.
- [ ] Roll out one table group at a time.
- [ ] Monitor denied-query and application-error rates during rollout.

Acceptance criteria:

- [ ] All authorization tests pass against real PostgreSQL.
- [ ] Missing or malformed authorization state fails closed.
- [ ] No known path relies on application-side filtering alone.
- [ ] Rollback procedures are documented and tested.

## Table-To-Permission Mapping Tracker

Complete this tracker before converting each table. Add rows as tables are discovered.

| Table | Client reference | Select permission | Insert permission | Update permission | Delete permission | Status |
|---|---|---|---|---|---|---|
| `client_details` | `id` | `CLIENT.VIEW` | `CLIENT.CREATE` | `CLIENT.UPDATE` | `CLIENT.DELETE` | Phase 8 completed |
| `progress_report` | `client_id` | `CLIENT.PROGRESS_REPORT.VIEW` or source read through `CLIENT.AI_PROGRESS_REPORT.GENERATE` | `CLIENT.PROGRESS_REPORT.CREATE` | `CLIENT.PROGRESS_REPORT.UPDATE` | `CLIENT.PROGRESS_REPORT.DELETE` | Phase 9 Group B completed |
| `client_diagnosis` | `client_id` | `CLIENT.DIAGNOSIS.VIEW` | `CLIENT.DIAGNOSIS.CREATE` | `CLIENT.DIAGNOSIS.UPDATE` | `CLIENT.DIAGNOSIS.DELETE` | Phase 9 Group C completed |
| `client_medication_order` | `client_id` | `CLIENT.MEDICATION.VIEW` | `CLIENT.MEDICATION.CREATE` | `CLIENT.MEDICATION.UPDATE` | `CLIENT.MEDICATION.DELETE` | Phase 9 Group C completed |
| `client_emergency_contact` | `client_id` | `CLIENT.EMERGENCY_CONTACT.VIEW` | `CLIENT.EMERGENCY_CONTACT.CREATE` | `CLIENT.EMERGENCY_CONTACT.UPDATE` | `CLIENT.EMERGENCY_CONTACT.DELETE` | Phase 9 Group A completed |
| `assigned_employee` | `client_id` | `CLIENT.INVOLVED_EMPLOYEE.VIEW` | `CLIENT.INVOLVED_EMPLOYEE.CREATE` (`all` only) | `CLIENT.INVOLVED_EMPLOYEE.UPDATE` (`all` only) | `CLIENT.INVOLVED_EMPLOYEE.DELETE` (`all` only) | Phase 9 Group A completed |
| `incident` | `client_id` | `CLIENT.INCIDENT.VIEW` | `CLIENT.INCIDENT.CREATE` | `CLIENT.INCIDENT.UPDATE` | `CLIENT.INCIDENT.DELETE` | Completed; confirmation uses `CLIENT.INCIDENT.CONFIRM` through a narrow DB function |
| `ai_generated_reports` | `client_id` | `CLIENT.AI_PROGRESS_REPORT.VIEW` | `CLIENT.AI_PROGRESS_REPORT.CONFIRM` | No operation or permission | No operation or permission | Phase 9 Group B completed |
| `client_documents` | `client_id` | `CLIENT.DOCUMENTS.VIEW` | `CLIENT.DOCUMENTS.UPLOAD` | To decide | `CLIENT.DOCUMENTS.DELETE` | Pending review |

## Known Risks

1. Existing RLS may currently be bypassed if the application connects as table owner or a role with `BYPASSRLS`.
2. Enforcing RLS before centralizing transaction identity could hide valid data or break endpoints.
3. Defaulting existing scoped grants to `all` could overgrant access and violate least privilege.
4. Defaulting existing scoped grants to `assigned` could unexpectedly remove legitimate access. Existing assignments and role intent must be reviewed during migration.
5. Assignment-management permissions can become a privilege-escalation path if users can assign themselves or others without strict checks.
6. Deep client ownership lookups can cause RLS recursion or poor query performance.
7. Multiple roles can cause accidental scope widening if grants and scopes are combined incorrectly.
8. Background jobs can fail or bypass controls if no explicit service-actor model exists.
9. Public registration and intake flows need separate treatment from authenticated client records.
10. RLS does not replace auditing, route permissions, field-level response controls, encryption, session controls, or operational access reviews.

## Decision Log

Record finalized decisions here. Do not silently change an earlier decision; add a new dated entry that supersedes it.

| Date | Decision | Reason | Status |
|---|---|---|---|
| 2026-08-17 | Roles are the sole source of user permissions; direct per-user overrides are removed. | The override model is legacy and would create a second authorization path that can bypass role scope policy. | Confirmed |
| 2026-08-17 | Replace a role's permission grants through one repository transaction after validating the complete request. | Invalid IDs or scopes must not remove valid existing grants. | Confirmed |
| 2026-08-16 | Use the name `scope`, not `client_scope`, throughout the design. | Scope may be generalized and the chosen API terminology is simpler. | Confirmed |
| 2026-08-16 | Scope belongs to each role-permission grant. | Different client-data categories may require different reach under least privilege. | Confirmed |
| 2026-08-16 | Initial scope values are `assigned` and `all`; unscoped permissions use `NULL`. | Supports current administrator/coordinator needs without premature scope types. | Confirmed |
| 2026-08-16 | Implement the work in small independently verified phases. | The authorization surface is large and high risk. | Confirmed |
| 2026-08-16 | Modify `000001_init.up.sql` and its down migration directly. | The project is in early development, has no production database, and databases can be dropped and recreated. | Confirmed |
| 2026-08-16 | Enforce permission/grant scope consistency with deferred constraint triggers. | Cross-table rules cannot use a normal check constraint; deferred validation permits atomic metadata and grant changes while rejecting invalid committed state. | Confirmed |
| 2026-08-16 | Scope every registered `CLIENT.*` permission except `CLIENT.CREATE`. | Existing-client operations require row-level reach; client creation has no existing assignment to evaluate. | Confirmed |
| 2026-08-16 | Rename active evaluation permissions to `CLIENT.EVALUATION.CREATE` and `CLIENT.EVALUATION.VIEW`, and remove the unused delete permission. | Evaluations belong to clients, and no evaluation delete route currently exists. | Confirmed |
| 2026-08-16 | Seed scoped administrator grants as `all` and scoped coordinator grants as `assigned`. | These values match the current role responsibilities while avoiding role-name checks in the future RLS design. | Confirmed |
| 2026-08-19 | `assigned` means any started `assigned_employee` row; deleting the row ends the assignment. | The assignment table has no end date, and unassignment already removes the relationship. | Confirmed |
| 2026-08-19 | `all` grants access to all clients in the system without an organizational boundary. | This is the required operational meaning of the broad scope. | Confirmed |
| 2026-08-19 | User-triggered workers carry delegated actor identity; scheduled and public trusted operations use a validated persisted system actor. | Workers need explicit authorization identity without relying on table-owner bypass. | Confirmed |
| 2026-08-20 | Restrict assignment create, update, and delete to `all` scope while allowing scoped assignment reads. | `assigned_employee` determines `assigned` access, so allowing assigned-scope mutation would permit self-granted client access. | Confirmed |
| 2026-08-20 | Own assignment lookup and narrow recipient-email functions with a no-login `BYPASSRLS` policy-owner role. | This permits forced RLS on authorization-source and contact tables without recursive policies or broadening direct row visibility. | Confirmed |
| 2026-08-20 | Let `CLIENT.AI_PROGRESS_REPORT.GENERATE` read scoped progress-report source rows without granting AI-report persistence or visibility. | Generation processes progress text but does not write or read `ai_generated_reports`; confirmation and AI history use separate permissions. | Confirmed |
| 2026-08-20 | Give `ai_generated_reports` only VIEW and CONFIRM/insert policies. | There are no AI report update/delete operations or registered permissions, so inventing mutation authority would overgrant access. | Confirmed |
| 2026-08-20 | Require both diagnosis and medication VIEW permissions for the combined medical overview. | General client visibility must not imply access to either sensitive medical category or return misleading partial overviews. | Confirmed |
| 2026-08-20 | Derive medical creator/updater attribution from the database transaction actor. | Request parameters and direct repository callers must not be able to forge medical provenance. | Confirmed |
| 2026-08-20 | Enforce medication diagnosis ownership with a composite `(diagnosis_id, client_id)` foreign key. | A medication order must not reference a diagnosis belonging to another client. | Confirmed |
| 2026-08-20 | Restrict deletion of diagnoses referenced by medication orders. | Referential cascades must not mutate medication data without medication permission and audit coverage. | Confirmed |
| 2026-08-22 | Remove the unused top-level `INCIDENT.VIEW` permission and standardize on scoped `CLIENT.INCIDENT.*`. | All incident routes and role grants already use the client-scoped permission family; retaining an unused duplicate leaves its meaning ambiguous. | Confirmed |
| 2026-08-22 | Require matching `CLIENT.VIEW`, `CLIENT.INCIDENT.VIEW`, and `CLIENT.INCIDENT.CONFIRM` scope for confirmation. | Confirmation delegates a protected report read to the email worker, so the initiating actor must be able to load the same incident and client identity. | Confirmed |
| 2026-08-22 | Treat incident UPDATE as requiring matching incident VIEW visibility. | PostgreSQL applies SELECT visibility to mutation queries that read existing columns and return rows; the route makes this dependency explicit. | Confirmed |
| 2026-08-22 | Treat incident DELETE as requiring matching incident VIEW visibility. | Returning the owning client for accurate audit attribution applies incident SELECT visibility to the deletion; the route makes this dependency explicit. | Confirmed |

## Progress Log

Add the newest entry first.

### 2026-08-22 - Phase 9 Group D incident RLS completed

Status: Completed and verified

Changes:

- Replaced legacy role-name incident policies with forced operation-specific `CLIENT.INCIDENT.*` permission-and-scope RLS.
- Removed the unused top-level `INCIDENT.VIEW` permission and retained the scoped client incident permission family as canonical.
- Preserved CREATE without broad VIEW through an actor-bound incident creation context and database-enforced reporter attribution.
- Isolated confirmation and confirmation-email marking behind narrow security-definer functions instead of granting broad UPDATE authority.
- Added an expiring database claim so concurrent confirmation tasks cannot deliver duplicate incident reports.
- Required matching client VIEW, incident VIEW, and incident CONFIRM scope so delegated confirmation workers can reload only authorized reports.
- Corrected nullable confirmation filters and preserved pagination totals for empty pages.
- Made zero-row incident deletion return not-found instead of false success.
- Restricted actorless incident seed writes to an explicitly checked superuser/`BYPASSRLS` bootstrap role.
- Added best-effort incident access, export, confirmation, and mutation audits without narrative or clinical payloads.

Verification:

- Applied and rolled back the initial migration on a clean PostgreSQL 17 database.
- Tested through a temporary `NOSUPERUSER NOBYPASSRLS` non-owner runtime role.
- Verified assigned/all visibility, forced policy state, scoped aggregates, production CREATE, unforgeable creation context, reporter attribution, update/delete separation, confirmation isolation, and PDF source visibility.
- Verified incident audit classification and sensitive-payload exclusion.

Next action:

- Convert Phase 9 Group E documents, care plans, goals, and evaluations.

### 2026-08-20 - Phase 9 Group C medical RLS completed

Status: Completed and verified

Changes:

- Replaced role-name policies on `client_diagnosis` and `client_medication_order` with forced permission-and-scope RLS.
- Preserved independent medical CREATE permissions with actor-bound, unforgeable creation contexts.
- Required both medical VIEW permissions for the combined overview and kept general client VIEW out of medical policies.
- Derived medical creator/updater employee IDs from the transaction actor.
- Prevented cross-client medication-to-diagnosis references with a composite foreign key.
- Restricted linked-diagnosis deletion so medication changes remain explicit and separately authorized.
- Made medication deletion report zero-row denial/not-found instead of false success.
- Added best-effort NEN 7513 medical access/change events containing identifiers and counts but no clinical payloads.
- Restricted actorless medical seed writes to an explicitly checked superuser/`BYPASSRLS` bootstrap role; normal application writes remain actor-bound.

Verification:

- Applied and rolled back the initial migration on a clean PostgreSQL 16 database.
- Tested through a temporary `NOSUPERUSER NOBYPASSRLS` non-owner runtime role.
- Verified category isolation, assigned/all visibility, general-VIEW denial, CREATE without VIEW, actor attribution, mutation visibility requirements, cross-client diagnosis-reference denial, forced policy state, unforgeable creation context, and failed `row_security=off` bypasses.
- Verified medical audit classification and clinical-payload exclusion.

Next action:

- Convert Phase 9 Group D incidents.

### 2026-08-20 - Phase 9 Group B report RLS completed

Status: Completed and verified

Changes:

- Replaced role-name policies on `progress_report` and `ai_generated_reports` with forced permission-and-scope RLS.
- Mapped progress-report CRUD to its dedicated permissions and allowed generate-only actors to read scoped source reports.
- Mapped AI-report reads to VIEW and confirmation inserts to CONFIRM, with no update or delete policies.
- Added client/date and client/creation-time indexes for report list and generation queries.
- Classified progress and AI report text as sensitive care data.
- Preserved independent CREATE and CONFIRM permissions by restricting `RETURNING` visibility to owner-generated report UUIDs in the same actor transaction.

Verification:

- Applied the initial migration to a clean PostgreSQL 16 database.
- Tested through a temporary `NOSUPERUSER NOBYPASSRLS` non-owner runtime role.
- Verified isolated VIEW/UPDATE/DELETE permissions, assigned/all progress CRUD, cross-client movement denial, generate-only source reads, create/confirm without VIEW, scoped AI confirmation, unforgeable creation context, absent AI mutation authority, missing identity/permission denial, forced-policy catalog state, and failed `row_security=off` bypasses.

Next action:

- Convert Phase 9 Group C medical data.

### 2026-08-20 - Phase 9 Group A client network RLS completed

Status: Completed and verified

Changes:

- Replaced role-name policies on `assigned_employee` and `client_emergency_contact` with operation-specific permission policies.
- Restricted assignment mutation to `all` scope to prevent self-assignment and assignment-rewrite privilege escalation.
- Applied normal `assigned`/`all` client scope to assignment reads and emergency-contact CRUD.
- Preserved intake promotion by allowing emergency-contact creation and readback only for the owner-generated client UUID in the current transaction.
- Added narrow permission-checking functions for general related-email and incident-recipient flows so those operations do not grant direct emergency-contact visibility.
- Forced RLS on both Group A tables and isolated assignment lookup behind a no-login policy-owner role to avoid recursive assignment policies.
- Classified emergency contacts as sensitive personal/disclosure data and assignments as security-sensitive access-control metadata.

Verification:

- Applied and rolled back the initial migration on a clean PostgreSQL 16 database.
- Tested through a temporary `NOSUPERUSER NOBYPASSRLS` non-owner runtime role.
- Verified assigned and global visibility, missing identity and permission denial, scoped emergency-contact CRUD, cross-client update denial, safe new-client contact creation/readback, assignment mutation denial for assigned scope, assignment mutation for all scope, and permission-specific recipient-email access.

Next action:

- Convert Phase 9 Group B progress reports and AI-generated reports.

### 2026-08-19 - Phase 8 client RLS pilot completed

Status: Completed and verified

Changes:

- Folded the Phase 7 authorization helpers into the disposable initial migration so a clean schema is self-contained.
- Tightened permission and scope lookup to require matching active user and employee identities.
- Replaced role-name policies on `client_details` with operation-specific permission policies.
- Mapped select, insert, update, and delete to `CLIENT.VIEW`, `CLIENT.CREATE`, `CLIENT.UPDATE`, and `CLIENT.DELETE` respectively.
- Made client creation request a fresh UUID through an owner-only transaction context so `INSERT ... RETURNING` can return only the newly created row without broadening `CLIENT.VIEW`.
- Forced row-level security on `client_details`, including for a non-superuser table owner.
- Left related tables on their legacy policies until their controlled Phase 9 conversion.

Verification:

- Applied the consolidated initial migration to a clean PostgreSQL 16 database.
- Tested through a temporary `NOSUPERUSER NOBYPASSRLS` non-owner runtime role.
- Verified assigned list filtering, hidden unassigned single reads, global `all`, missing identity, missing permission, authorized creation, scoped update/delete, and failed `row_security=off` bypass.
- Existing Phase 6 and Phase 7 PostgreSQL integration tests continue to pass.

Next action:

- Begin Phase 9 by converting related client tables in controlled permission groups.

### 2026-08-19 - Phase 7 database authorization helpers added

Status: Completed and verified

Decisions:

- `all` grants access to all clients in the system; no organization boundary applies.
- `assigned` accepts any `assigned_employee` role after its `start_date`.
- An assignment remains active while its row exists; unassigning an employee deletes that row.

Changes:

- Added fail-closed transaction identity getters for current user and employee UUIDs.
- Added live role-grant lookup through `has_permission` and `get_permission_scope`.
- Added assignment and client-access helpers without business role-name checks.
- Required the configured user and employee identities to match active database records before client access.
- Added fixed function search paths, schema-qualified references, restricted parallel execution, and revoked default `PUBLIC` execution.
- Fixed system-actor startup validation to use the existing `is_archived` and `out_of_service` employee fields.
- Updated normal migration targets to apply all pending migrations.

Verification:

- Applied the authorization helpers as part of the consolidated initial migration to a clean PostgreSQL 16 database.
- Direct SQL integration tests cover missing and malformed context, absent grants, scoped and unscoped grants, future/current/deleted assignments, live scope changes, global `all`, inactive users, inactive employees, and mismatched identities.
- `go test ./...` completed successfully.

Next action:

- Begin Phase 8 by replacing the role-name policies on `client_details` with operation-specific permission and scope policies.

### 2026-08-19 - Phase 6 protected execution centralized

Status: Completed and verified

Changes:

- Added validated user-and-employee actor identity to the shared request context and reject incomplete access-token identities.
- Added strict `ExecActorTx` and `BeginActorTx` entry points that install both PostgreSQL settings with transaction-local `set_config`.
- Removed the employee-ID standard-output print and caller-controlled database identity setup.
- Converted existing client, incident, intake, contract, invoice, event, and notification transaction blocks to authenticated actor transactions.
- Preserved `ExecTx` for explicit bootstrap and maintenance work; it installs actor settings when valid identity is present.
- Added seed invoice actor propagation for service calls that now require authenticated transactions.
- Added delegated actor payloads for incident-confirmation jobs and a required configured system actor for scheduled maintenance jobs.
- Converted all standalone protected client, incident, contract, invoice, and relevant intake queries to authenticated transactions.
- Routed public registration and intake-option operations through the explicit configured service actor because those tables already have RLS enabled.
- Added startup validation that the configured service actor references the same active persisted user and employee.
- Converted protected dashboard, organization aggregate, and calendar attendee queries found during the final cross-table audit.

Verification:

- `go test ./...` completed successfully.
- Real PostgreSQL tests use a one-connection pool to verify both identity settings, missing-identity rejection, and no identity leakage after commit or rollback.

Next action:

- Provision `SYSTEM_ACTOR_USER_ID` and `SYSTEM_ACTOR_EMPLOYEE_ID` in each environment before deployment.
- Begin Phase 7 by adding general database authorization functions based on permission grants and scope.

### 2026-08-17 - Phase 5 effective permission scopes resolved

Status: Completed and verified

Changes:

- Effective permission lists now return `is_scoped` and the scope from the specific role grant.
- Direct permission checks now fetch one usable grant instead of loading every permission into Go.
- Scoped grants without a valid scope and unscoped grants with a scope fail closed.
- Middleware stores the resolved effective permission in request context for later client-access enforcement.
- Role context remains available temporarily for existing role-based RLS policies.
- Handbook permission checks now use the same scoped grant resolver.

Verification:

- `sqlc generate` completed successfully.
- `go test ./...` completed successfully.
- Middleware tests cover missing grants, exact scope propagation, malformed scoped grants, and grant changes between requests using the same JWT identity.
- A fresh PostgreSQL 16 migration returned `NULL` for an unscoped grant, returned `assigned` and then `all` after a grant update, denied an absent permission, and rolled back successfully.

Next action:

- Begin Phase 6 by centralizing protected database execution with transaction-local actor identity.

### 2026-08-17 - Phase 4 legacy user overrides removed

Status: Completed and verified

Changes:

- Removed the `user_permission_overrides` table and `permission_override_effect` enum from the initial migration.
- Removed all override queries, generated sqlc models, domain types, repository/service methods, DTOs, and the direct user-permission endpoint.
- Removed the obsolete `PERMISSIONS.GRANT` registry entry.
- Simplified effective permission checks and employee-profile permission aggregation to use role grants only.
- Removed override fields from the employee role-and-permission response.

Verification:

- `sqlc generate` completed successfully.
- `go test ./...` completed successfully.
- A fresh migration verified that the override table and enum are absent.
- Role-granted permissions remained available and permissions absent from the role remained denied.

Next action:

- Begin Phase 5 by returning each role-derived permission's scope and using a direct scoped permission check in middleware.

### 2026-08-17 - Phase 3 role-permission management completed

Status: Completed and verified

Changes:

- Replaced the permission-ID-only role mutation contract with permission grants containing `permission_id` and nullable `scope`.
- Added strict `assigned` and `all` scope validation in the role service.
- Reject unknown IDs, duplicate IDs, missing scopes on scoped permissions, and scopes on unscoped permissions before modifying grants.
- Return `is_scoped` in the permission catalog and return both `is_scoped` and `scope` for role grants.
- Replaced separate delete and bulk-insert repository calls with one transaction-owned replacement operation.
- Updated SQL and regenerated sqlc code to read and write role grant scopes.
- Added service validation/delegation tests and handler 400/500 mapping tests.

Verification:

- `sqlc generate` completed successfully.
- Focused domain, repository, service, and handler tests passed.
- A fresh migration completed on disposable PostgreSQL 16.
- Mixed `assigned` and unscoped grants were inserted and read back with matching metadata.
- A replacement containing a missing scope failed at the deferred database constraint and rolled back to the original two grants.

Next action:

- Phase 4 superseded the earlier override design by removing direct user permissions entirely.

### 2026-08-16 - Phase 2 client permission classification completed

Status: Completed and verified

Changes:

- Added `IsScoped` to permission definitions in `internal/domain/permission.go`.
- Classified every `CLIENT.*` permission as scoped except `CLIENT.CREATE`.
- Replaced `EVALUATION.CREATE` and `EVALUATION.VIEW` with `CLIENT.EVALUATION.CREATE` and `CLIENT.EVALUATION.VIEW`.
- Removed the unused `EVALUATION.DELETE` permission because no delete route exists.
- Updated evaluation routes in `internal/handler/client_handler.go` to require dedicated evaluation permissions instead of broad `CLIENT.VIEW` and `CLIENT.UPDATE` permissions.
- Added default role scope metadata: administrator uses `all`, coordinator uses `assigned`.
- Updated `cmd/roles/main.go` to persist `permissions.is_scoped` and upsert role grant scopes.
- Added `internal/domain/permission_test.go` for client classification and evaluation registry coverage.

Verification:

- `sqlc generate` completed successfully without unexpected generated changes.
- `go test ./...` completed successfully.
- A fresh migration and RBAC sync completed on disposable PostgreSQL 17.
- The registry produced 42 `CLIENT.*` permissions: 41 scoped and one unscoped.
- `CLIENT.CREATE` had `scope = NULL`.
- Administrator client grants had `scope = all`.
- Coordinator client grants had `scope = assigned`.
- The database contained zero invalid role-permission scope combinations.
- The database contained zero legacy `EVALUATION.*` permissions.
- Running the RBAC sync a second time completed successfully and preserved valid scopes.

Known transition state:

- The role-management API still accepts permission IDs without scopes. Editing scoped role grants through that API is intentionally deferred to Phase 3 and may be rejected by the database constraints until Phase 3 is complete.

Next action:

- Begin Phase 3 by updating role-permission domain models, SQL, API requests/responses, and transactional replacement to support `scope`.

### 2026-08-16 - Phase 1 database representation completed

Status: Completed and verified

Changes:

- Added `permission_scope_enum` with `assigned` and `all` values to `db/migrations/000001_init.up.sql`.
- Added `permissions.is_scoped BOOLEAN NOT NULL DEFAULT FALSE`.
- Added nullable `role_permissions.scope`.
- Added deferred constraint triggers that reject missing scope for scoped permissions and non-null scope for unscoped permissions.
- Updated `db/migrations/000001_init.down.sql` to remove the validation functions and enum.
- Regenerated `db/sqlc/models.go` and `db/sqlc/roles.sql.go` with sqlc 1.31.1.
- Did not change RLS policies or runtime authorization behavior.

Verification:

- `sqlc generate` completed successfully.
- `go test ./...` completed successfully; the repository currently reports no test files.
- `git diff --check` completed successfully before progress registration.
- Migration `up 1` completed successfully on a disposable PostgreSQL 17 container.
- Valid unscoped `NULL`, scoped `assigned`, and scoped `all` grants committed successfully.
- A scoped grant with `NULL` was rejected.
- An unscoped grant with `all` was rejected.
- Atomically changing a permission to scoped and its grant to `assigned` succeeded with deferred validation.
- Migration `down 1` completed successfully.

Next action:

- Begin Phase 2 by classifying every permission in `internal/domain/permission.go` as scoped or unscoped.

### 2026-08-16 - Planning document created

Status: Completed

Changes:

- Created `docs/PERMISSION_SCOPE_RLS_PLAN.md`.
- Recorded current authorization architecture and known gaps.
- Defined twelve implementation phases plus a design-decision phase.
- Added file paths, acceptance criteria, a table-to-permission tracker, risks, and decision tracking.

Verification:

- Documentation-only change; no application behavior changed.

Next action:

- Resolve the remaining Phase 0 decisions when their dependent phases begin.

## Session Restart Checklist

At the beginning of a future session:

1. Read this document completely.
2. Check `git status` and preserve unrelated changes.
3. Read the latest Progress Log and Decision Log entries.
4. Identify the first incomplete phase and its acceptance criteria.
5. Re-read the listed affected files because the codebase may have changed.
6. Mark only the current phase item as in progress.
7. Make the smallest change needed for that phase.
8. Run the phase verification before marking it complete.
9. Update this document in the same change.
