# Permission Scope and Row-Level Security Plan

Last updated: 2026-08-16

Status: Phase 1 completed; Phase 2 not started

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
- `all`: all clients within the applicable organizational boundary
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
11. For `all`, RLS allows the row within the mandatory organizational boundary.
12. Commit or roll back the transaction.
13. Return only rows PostgreSQL allowed.

Roles, permissions, and scopes should not be trusted from JWT claims. The JWT identifies the actor; current authorization is loaded from the database so changes can take effect without waiting for token expiration.

## Current State

### Database RBAC

Defined in `db/migrations/000001_init.up.sql`:

- `roles` at approximately lines 475-479
- `permissions` at approximately lines 482-489
- `role_permissions` at approximately lines 493-500
- `user_permission_overrides` at approximately lines 520-528
- `user_roles` at approximately lines 531-536

Current limitations:

- `role_permissions` contains only `role_id` and `permission_id`.
- Permissions do not indicate whether scope applies.
- User permission overrides have no scope.
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

- Role-permission requests accept only arrays of permission UUIDs.
- Responses do not expose scope.
- Replacing role permissions uses delete followed by insert without one transaction.
- Replacing user overrides also uses multiple operations without one transaction.
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
- User permission overrides do not affect RLS.
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
- [ ] Confirm whether `assigned` means any active `assigned_employee` record or only selected assignment types.
- [ ] Confirm how an assignment ends; the current table has `start_date` but no `end_date`.
- [ ] Confirm whether `all` means all clients globally or all clients in the actor's organization.
- [ ] Confirm whether users will remain limited to one role or may receive multiple roles later.
- [ ] Confirm scoped behavior for direct per-user `allow` overrides.
- [ ] Confirm behavior for background workers and trusted system operations.
- [ ] Confirm whether `CLIENT.CREATE` is unscoped initially, because a client does not have an assignment before creation.

Recommended decisions:

- Create a new migration if any non-disposable database has run migration `000001`.
- Treat `all` as all clients inside the user's organization, not all organizations.
- Treat any active assignment as `assigned`; permissions determine which client data categories are accessible.
- Add an optional assignment end date.
- Require `assigned` or `all` for scoped per-user allow overrides.
- Keep deny overrides unscoped because deny removes the permission entirely.
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

Status: `[ ]` Not started

Goal: make permission metadata the source of truth for whether scope applies.

Affected files:

- `internal/domain/permission.go`
- `cmd/roles/main.go`
- `db/query/roles.sql`
- Generated files under `db/sqlc/`

Initial classification principles:

- Permissions operating on an existing client's records are normally scoped.
- General system permissions are unscoped.
- Client creation is initially unscoped unless a clear target-client scope exists.
- Medical, incident, document, financial, and basic client information are classified separately.

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
- Relevant invoice and contract permissions after ownership is mapped

Examples likely to remain unscoped:

- `ROLES.*`
- `PERMISSIONS.*`
- General settings permissions
- Authentication and self-service operations
- `CLIENT.CREATE`, initially

Checklist:

- [ ] Extend permission metadata with `IsScoped`.
- [ ] Review every permission in `internal/domain/permission.go`.
- [ ] Produce a permission classification table in this document or a linked document.
- [ ] Update `cmd/roles/main.go` to persist `is_scoped`.
- [ ] Ensure permission synchronization updates changed metadata.
- [ ] Verify every registered permission has an explicit classification.

Acceptance criteria:

- [ ] No permission relies on an implicit scoped/unscoped default in Go metadata.
- [ ] Sensitive client-data categories are separately classified.
- [ ] The database and Go registry agree.

## Phase 3: Update Role-Permission Management

Status: `[ ]` Not started

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

- [ ] Add a domain scope type with strict values.
- [ ] Replace permission-ID-only mutation models with permission grant models.
- [ ] Return `is_scoped` and `scope` from role detail APIs.
- [ ] Update sqlc role-permission queries to read and write scope.
- [ ] Validate unknown scope values as client errors.
- [ ] Validate missing scope on scoped permissions.
- [ ] Validate non-null scope on unscoped permissions.
- [ ] Make role-permission replacement one transaction.
- [ ] Validate all requested permission IDs before replacing existing grants.
- [ ] Preserve the old grants if replacement fails.
- [ ] Regenerate sqlc code.
- [ ] Update API documentation or frontend contract if maintained elsewhere.

Acceptance criteria:

- [ ] A role can contain mixed `assigned`, `all`, and unscoped grants.
- [ ] Invalid combinations return a clear 4xx response.
- [ ] Replacement is atomic.
- [ ] API reads return exactly what was saved.

## Phase 4: Add Scope To User Permission Overrides

Status: `[ ]` Not started

Goal: prevent direct user allowances from becoming unintended unrestricted access.

Affected files:

- Database migration files
- `db/query/roles.sql`
- `internal/domain/role.go`
- `internal/repository/role_repo.go`
- `internal/service/role_service.go`
- `internal/handler/role_dto.go`
- Generated files under `db/sqlc/`

Target rules:

- `deny` removes the permission and uses `scope = NULL`.
- Allowing an unscoped permission requires `scope = NULL`.
- Allowing a scoped permission requires `assigned` or `all`.
- Explicit deny wins over role inheritance and explicit allow.

Checklist:

- [ ] Add nullable `scope` to `user_permission_overrides`.
- [ ] Add database-level validation for effect, permission type, and scope.
- [ ] Update override request and response DTOs.
- [ ] Update domain and repository models.
- [ ] Make override replacement one transaction.
- [ ] Reject overlapping allow and deny entries.
- [ ] Add precedence tests.

Acceptance criteria:

- [ ] A scoped user allowance cannot be stored without a scope.
- [ ] A deny always removes effective access.
- [ ] Failed replacement preserves the previous overrides.

## Phase 5: Resolve Effective Permissions And Scope

Status: `[ ]` Not started

Goal: calculate the user's current effective permission grant and scope from database data.

Affected files:

- `db/query/roles.sql`
- `internal/domain/role.go`
- `internal/repository/role_repo.go`
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

- Deny overrides win.
- Unscoped effective permissions return `scope = NULL`.
- Missing scope on a scoped grant results in no usable grant.
- If multiple roles are introduced, scope combination must be explicit and tested. Never combine a permission from one role with `all` scope from an unrelated role that does not grant that permission.

Checklist:

- [ ] Update inherited permission queries to include scope.
- [ ] Update effective permission queries to include scope.
- [ ] Update direct permission checks.
- [ ] Avoid loading all permissions merely to check one key where possible.
- [ ] Keep role and scope data out of JWT claims.
- [ ] Add effective-permission tests for role grants and user overrides.

Acceptance criteria:

- [ ] HTTP middleware still returns a clear `403` when permission is absent.
- [ ] Effective scope matches the specific grant that provides the permission.
- [ ] Database changes to grants affect new requests without issuing a new JWT.

## Phase 6: Centralize Protected Database Execution

Status: `[ ]` Not started

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

- [ ] Add user ID to the standard internal request context if not already accessible there.
- [ ] Define one authenticated transaction entry point.
- [ ] Set both user and employee IDs in that transaction.
- [ ] Inventory all direct pool access to protected tables.
- [ ] Convert protected direct queries to the standard transaction path.
- [ ] Convert manual transactions to the shared initialization method.
- [ ] Define worker behavior explicitly.
- [ ] Test pooled connection reuse for identity leakage.
- [ ] Test missing identity behavior.

Acceptance criteria:

- [ ] No protected repository path can accidentally omit actor identity.
- [ ] Identity does not leak to a later request on the same pooled connection.
- [ ] Workers do not depend on table-owner RLS bypass.

## Phase 7: Add General Database Authorization Functions

Status: `[ ]` Not started

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
4. User has the requested permission after overrides.
5. The permission is scoped.
6. `all` satisfies the organizational boundary.
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

- [ ] Implement identity getter functions.
- [ ] Implement effective-permission lookup.
- [ ] Implement scope lookup.
- [ ] Implement active-assignment lookup.
- [ ] Implement `can_access_client`.
- [ ] Harden function ownership and execution grants.
- [ ] Add direct SQL tests for every allow and deny path.
- [ ] Confirm missing context returns false rather than raising an information-leaking error.

Acceptance criteria:

- [ ] No helper checks a business role name such as `admin` or `coordinator`.
- [ ] Permission and scope changes affect authorization immediately on the next transaction.
- [ ] Invalid or missing context denies access.

## Phase 8: Convert `client_details` As The RLS Pilot

Status: `[ ]` Not started

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

Checklist:

- [ ] Replace role-name policies on `client_details`.
- [ ] Add separate policies for select, insert, update, and delete.
- [ ] Verify list queries filter rows automatically.
- [ ] Verify single-record queries hide unauthorized clients.
- [ ] Verify update and delete denial.
- [ ] Verify administrator/all behavior.
- [ ] Verify assigned behavior.
- [ ] Verify missing permission behavior.
- [ ] Verify missing identity behavior.
- [ ] Verify explicit user deny behavior.
- [ ] Verify organization boundary behavior once defined.

Acceptance criteria:

- [ ] An assigned coordinator sees only assigned clients.
- [ ] An `all` grant sees all clients inside its permitted organization.
- [ ] A user without `CLIENT.VIEW` sees no client rows.
- [ ] The application runtime role cannot bypass the policy.
- [ ] Existing client API behavior remains correct for authorized users.

## Phase 9: Convert Related Tables In Controlled Groups

Status: `[ ]` Not started

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

- [ ] Map operations.
- [ ] Implement policies.
- [ ] Test assignment-management edge cases.
- [ ] Prevent users from granting themselves access through assignment changes.

### Group B: Progress Reports And AI Reports

Likely tables and files:

- `progress_report`
- `ai_generated_reports`
- Related client queries and service methods

Permissions include `CLIENT.PROGRESS_REPORT.*` and `CLIENT.AI_PROGRESS_REPORT.*`.

- [ ] Map operations.
- [ ] Implement policies.
- [ ] Test read, create, update, delete, generate, and confirm flows.

### Group C: Medical Data

Likely tables:

- `client_diagnosis`
- `client_medication_order`
- Related supporting tables

Permissions include `CLIENT.DIAGNOSIS.*` and `CLIENT.MEDICATION.*`.

- [ ] Map operations.
- [ ] Implement policies.
- [ ] Add stricter access and audit tests for medical data.

### Group D: Incidents

Likely tables and files:

- `incident`
- `db/query/incident.sql` or the applicable incident query files
- `internal/domain/incident.go`
- `internal/repository/incident_repo.go`
- `internal/service/incident_service.go`
- `internal/handler/incident_*`

Permissions include `CLIENT.INCIDENT.*` and any top-level incident permissions whose semantics must be reconciled.

- [ ] Resolve duplicate or overlapping incident permission meanings.
- [ ] Map operations.
- [ ] Implement policies.
- [ ] Test confirmation and file access.

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
- Record user override changes.
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
- User with explicit allow
- User with explicit deny
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
| `client_details` | `id` | `CLIENT.VIEW` | `CLIENT.CREATE` | `CLIENT.UPDATE` | `CLIENT.DELETE` | Planned pilot |
| `progress_report` | `client_id` | `CLIENT.PROGRESS_REPORT.VIEW` | `CLIENT.PROGRESS_REPORT.CREATE` | `CLIENT.PROGRESS_REPORT.UPDATE` | `CLIENT.PROGRESS_REPORT.DELETE` | Pending review |
| `client_diagnosis` | `client_id` | `CLIENT.DIAGNOSIS.VIEW` | `CLIENT.DIAGNOSIS.CREATE` | `CLIENT.DIAGNOSIS.UPDATE` | `CLIENT.DIAGNOSIS.DELETE` | Pending review |
| `client_medication_order` | `client_id` | `CLIENT.MEDICATION.VIEW` | `CLIENT.MEDICATION.CREATE` | `CLIENT.MEDICATION.UPDATE` | `CLIENT.MEDICATION.DELETE` | Pending review |
| `client_emergency_contact` | `client_id` | `CLIENT.EMERGENCY_CONTACT.VIEW` | `CLIENT.EMERGENCY_CONTACT.CREATE` | `CLIENT.EMERGENCY_CONTACT.UPDATE` | `CLIENT.EMERGENCY_CONTACT.DELETE` | Pending review |
| `assigned_employee` | `client_id` | `CLIENT.INVOLVED_EMPLOYEE.VIEW` | `CLIENT.INVOLVED_EMPLOYEE.CREATE` | `CLIENT.INVOLVED_EMPLOYEE.UPDATE` | `CLIENT.INVOLVED_EMPLOYEE.DELETE` | Pending review |
| `incident` | `client_id` | `CLIENT.INCIDENT.VIEW` | `CLIENT.INCIDENT.CREATE` | `CLIENT.INCIDENT.UPDATE` | `CLIENT.INCIDENT.DELETE` | Pending review |
| `ai_generated_reports` | `client_id` | `CLIENT.AI_PROGRESS_REPORT.VIEW` | `CLIENT.AI_PROGRESS_REPORT.CONFIRM` or dedicated permission | To decide | To decide | Decision required |
| `client_documents` | `client_id` | `CLIENT.DOCUMENTS.VIEW` | `CLIENT.DOCUMENTS.UPLOAD` | To decide | `CLIENT.DOCUMENTS.DELETE` | Pending review |

## Known Risks

1. Existing RLS may currently be bypassed if the application connects as table owner or a role with `BYPASSRLS`.
2. Enforcing RLS before centralizing transaction identity could hide valid data or break endpoints.
3. Defaulting existing scoped grants to `all` could overgrant access and violate least privilege.
4. Defaulting existing scoped grants to `assigned` could unexpectedly remove legitimate access. Existing assignments and role intent must be reviewed during migration.
5. Direct per-user allow overrides can overgrant unless they receive explicit scope.
6. Assignment-management permissions can become a privilege-escalation path if users can assign themselves or others without strict checks.
7. Deep client ownership lookups can cause RLS recursion or poor query performance.
8. Multiple roles can cause accidental scope widening if grants and scopes are combined incorrectly.
9. Background jobs can fail or bypass controls if no explicit service-actor model exists.
10. Public registration and intake flows need separate treatment from authenticated client records.
11. RLS does not replace auditing, route permissions, field-level response controls, encryption, session controls, or operational access reviews.

## Decision Log

Record finalized decisions here. Do not silently change an earlier decision; add a new dated entry that supersedes it.

| Date | Decision | Reason | Status |
|---|---|---|---|
| 2026-08-16 | Use the name `scope`, not `client_scope`, throughout the design. | Scope may be generalized and the chosen API terminology is simpler. | Confirmed |
| 2026-08-16 | Scope belongs to each role-permission grant. | Different client-data categories may require different reach under least privilege. | Confirmed |
| 2026-08-16 | Initial scope values are `assigned` and `all`; unscoped permissions use `NULL`. | Supports current administrator/coordinator needs without premature scope types. | Confirmed |
| 2026-08-16 | Implement the work in small independently verified phases. | The authorization surface is large and high risk. | Confirmed |
| 2026-08-16 | Modify `000001_init.up.sql` and its down migration directly. | The project is in early development, has no production database, and databases can be dropped and recreated. | Confirmed |
| 2026-08-16 | Enforce permission/grant scope consistency with deferred constraint triggers. | Cross-table rules cannot use a normal check constraint; deferred validation permits atomic metadata and grant changes while rejecting invalid committed state. | Confirmed |
| TBD | Exact meaning and lifetime of an active assignment. | Needed for `assigned` scope. | Open |
| TBD | Organizational boundary for `all`. | Needed to prevent cross-organization access. | Open |
| TBD | Scoped user-allow override behavior. | Needed before Phase 4. | Open |
| TBD | Background worker authorization model. | Needed before strict RLS enforcement. | Open |

## Progress Log

Add the newest entry first.

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
