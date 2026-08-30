# Evaluation Remediation Plan

This document tracks the evaluation workflow remediation across:

- Backend: `maicare_go`
- Frontend: `new_maicare_svelte`

It is based on the evaluation audit completed against the current backend and frontend implementations.

## Status Legend

- `[ ]` Not started
- `[~]` In progress or implemented but not fully verified
- `[x]` Completed and verified
- `[!]` Blocked by a decision or dependency

## Current Contract

The backend creates an evaluation as a scheduled client workflow identified by:

```text
client_id + client.next_evaluation_date
```

Current mutation behavior:

```text
POST  /clients/{clientId}/evaluations
PATCH /evaluations/{evaluationId}/draft
POST  /evaluations/{evaluationId}/submit
```

- Client-scoped POST creates or returns the current-cycle draft.
- PATCH updates exactly the draft identified by `evaluationId`.
- Submit transitions exactly the saved draft identified by `evaluationId` without accepting a request body.
- Historical drafts remain readable but PATCH and submit return HTTP 409 with `EVALUATION_NOT_CURRENT_CYCLE`.
- The legacy `submit` create field remains temporarily supported until its planned deprecation.

Existing drafts use an ID-based workflow:

```text
Bootstrap -> evaluation ID -> update the same ID -> submit the same ID
```

## P0: Data Integrity And Contract Safety

### P0.1 Bootstrap Must Return Only The Current-Cycle Draft

**Status:** `[x]` Completed and verified with PostgreSQL Testcontainers.

Backend tasks:

- [x] Replace bootstrap's latest-draft lookup with a lookup by `client_id`, `next_evaluation_date`, and `draft` status.
- [x] Return no `existing_draft` when the client has no `next_evaluation_date`.
- [x] Reuse `GetDraftGoalEvaluationByClientAndDate`; no migration or SQLC regeneration is required.
- [x] Run `go test ./...`.
- [x] Run `go vet ./...`.
- [x] Run database integration tests with the shared PostgreSQL Testcontainers harness.
- [x] Add focused integration coverage for a client with both historical and current-cycle drafts.
- [x] Commit the backend change.

Changed backend file:

```text
internal/repository/client_repo.go
```

Acceptance criteria:

- A historical draft is never returned as `existing_draft` by bootstrap.
- A current-cycle draft is returned regardless of whether a newer historical draft exists.
- A client without a scheduled evaluation date receives `existing_draft: null`.
- The returned draft's `evaluation_date` equals the client's `next_evaluation_date`.

### P0.2 Adopt ID-Based Update And Submit Operations

**Status:** `[x]` Backend and frontend completed and verified, including repository integration coverage and regenerated Swagger artifacts.

Target API:

```text
POST  /clients/{clientId}/evaluations
PATCH /evaluations/{evaluationId}/draft
POST  /evaluations/{evaluationId}/submit
```

Recommended responsibilities:

- Client-scoped POST creates or returns the current-cycle draft.
- PATCH updates exactly the draft identified by `evaluationId`.
- Submit transitions exactly the draft identified by `evaluationId`.
- Existing drafts must never be updated indirectly through a client/date lookup.

Backend tasks:

- [x] Add `UpdateGoalEvaluationDraft` to the domain repository and service interfaces.
- [x] Add `SubmitGoalEvaluationDraft` to the domain repository and service interfaces.
- [x] Add PATCH and submit handlers and route registration.
- [x] Validate UUIDs and request payloads at the handler boundary.
- [x] Enforce client access, evaluation permissions, creator ownership, and draft status.
- [x] Submit an already-saved revision without accepting a request body.
- [x] Keep client-scoped POST idempotent when a current draft already exists.
- [x] Preserve blocked-submit compatibility as HTTP 200 with `submit_error` until P1.2 finalizes response semantics.
- [x] Add focused service and handler tests for exact-ID delegation, validation, and error mapping.
- [x] Run repository/database integration coverage with the shared PostgreSQL Testcontainers harness.
- [ ] Deprecate and eventually remove `submit` from the create request after frontend migration.
- [x] Add Swagger annotations for the new routes.
- [x] Regenerate committed Swagger artifacts.

Frontend tasks:

- [x] Use client-scoped POST only for a new current-cycle workflow.
- [x] Use PATCH when saving an existing draft.
- [x] Save through PATCH before using the dedicated submit operation.
- [x] Keep the active `evaluationId` after the first successful draft creation.
- [x] Remove mutation logic in which `clientId` takes priority over `evaluationId`.
- [x] Preserve local form values across successful mutations and API failures.
- [x] Open recent drafts with their evaluation ID instead of re-entering the client create flow.
- [x] Send the submit request without a request body.
- [x] Run `bun run check` with zero errors; four unrelated pre-existing warnings remain.
- [x] Run `bun run build` successfully.
- [ ] Complete repository-wide lint cleanup; `bun run lint` currently fails on 111 pre-existing formatting issues outside this change.

Acceptance criteria:

- Editing evaluation A can never modify evaluation B.
- Save and submit operations use the ID displayed in the active workflow.
- Completed and archived evaluations cannot be modified.
- A mutation cannot silently resolve a different evaluation from the client's schedule.

### P0.3 Prevent Historical Drafts From Writing Into The Current Cycle

**Status:** `[x]` Backend and frontend completed and verified, including repository integration coverage.

Agreed policy:

```text
Historical drafts are read-only.
```

A historical draft is one whose `evaluation_date` does not equal the client's current `next_evaluation_date`.

Backend tasks:

- [x] Reject PATCH for historical drafts.
- [x] Reject submit for historical drafts.
- [x] Return HTTP 409 with `EVALUATION_NOT_CURRENT_CYCLE`.
- [x] Keep GET detail access available according to view permissions.
- [x] Ensure historical drafts are never copied into a current-cycle draft implicitly.
- [x] Validate the cycle inside the same transaction as each mutation.
- [x] Add database integration coverage with the shared PostgreSQL Testcontainers harness.

Frontend tasks:

- [x] Detect `EVALUATION_NOT_CURRENT_CYCLE` and preserve the user's current input.
- [x] Open known historical drafts in read-only mode.
- [x] Explain that the evaluation belongs to an earlier cycle.
- [x] Provide navigation to the current-cycle draft when one exists.

Acceptance criteria:

- Historical drafts can be reviewed but not changed or submitted.
- A historical draft mutation cannot create a current-cycle evaluation.
- The API returns a stable machine-readable conflict code.

## P1: Workflow Reliability And API Quality

### P1.1 Stable Evaluation Error Codes

**Status:** `[x]` Lifecycle and validation codes are complete and localized.

- [x] Return stable lifecycle codes for missing, non-owner, completed, and historical evaluations.
- [x] Localize lifecycle messages in the frontend based on codes.
- [x] Keep raw backend messages as diagnostics, not primary user copy.
- [x] Add stable codes for create validation and submission-window failures.

Proposed codes:

```text
EVALUATION_CLIENT_NOT_IN_CARE
EVALUATION_NO_ACTIVE_GOALS
EVALUATION_NO_DUE_DATE
EVALUATION_NOT_OWNER
EVALUATION_NOT_FOUND
EVALUATION_DUPLICATE_GOAL
EVALUATION_GOAL_NOT_ACTIVE
EVALUATION_INVALID_PROGRESS
EVALUATION_INCOMPLETE
EVALUATION_TOO_EARLY
EVALUATION_ALREADY_COMPLETED
EVALUATION_NOT_CURRENT_CYCLE
EVALUATION_CONFLICT
```

### P1.2 Submission Response Semantics

**Status:** `[x]` Blocked submission uses HTTP 422 with an explicit saved-draft payload.

- [x] Replace HTTP 200 `submit_error` compatibility responses.
- [x] Return HTTP 422 with a stable code for eligibility failures; reserve HTTP 409 for lifecycle and concurrency conflicts.
- [x] Return `data.draft_saved` and the persisted evaluation when submission is blocked after saving.
- [x] Update frontend success and error handling to match the final contract.

Blocked submission contract:

```json
{
  "success": false,
  "code": "EVALUATION_INCOMPLETE",
  "message": "all active goals must be evaluated before submission",
  "data": {
    "draft_saved": true,
    "evaluation": {
      "status": "draft"
    }
  }
}
```

### P1.3 Optimistic Concurrency

**Status:** `[x]` Draft updates and submissions use an `If-Match` `updated_at` compare-and-swap precondition.

- [x] Add an `updated_at` precondition to draft mutations.
- [x] Reject missing preconditions with HTTP 428 and stale writes with HTTP 409 `EVALUATION_CONFLICT`.
- [x] Return the complete current server evaluation, including its revision, in conflict responses.
- [x] Add frontend conflict UI that preserves local entries until an explicit confirmed reload.

### P1.4 Transaction Boundaries

- [x] Keep draft save and submission as separate transactions so a blocked submission preserves the saved draft.
- [x] Prevent another request from changing the draft between final save and submission by submitting the revision returned by the final save.
- [x] Ensure schedule advancement can happen only once for sequential submissions.
- [x] Verify concurrent submission behavior and single schedule advancement with PostgreSQL integration coverage.

### P1.5 Evaluation Permissions

**Status:** `[x]` Backend policy documented and frontend routes/actions aligned with the evaluation permissions.

Backend permissions currently used:

```text
CLIENT.VIEW
CLIENT.EVALUATION.VIEW
CLIENT.EVALUATION.CREATE
```

Final policy:

- Admin receives evaluation view and create permissions with `all` scope.
- Coordinator receives evaluation view and create permissions with `assigned` scope.
- `assigned` permits access only when the employee is assigned to the client; `all` permits access to every client allowed by the organization context.
- Evaluation routes also require `CLIENT.VIEW` where they expose general client details.
- `CLIENT.EVALUATION.VIEW` controls bootstrap, goals, listings, history, and detail access.
- `CLIENT.EVALUATION.CREATE` remains the mutation permission for create, update, and submit.
- Draft updates and submissions remain restricted to the employee who created the draft.
- Separate update and submit permissions are deferred until a distinct reviewer or approver workflow is required.

- [x] Document which roles receive evaluation view and create permissions.
- [x] Document `assigned` versus `all` permission scopes.
- [x] Confirm whether view and mutation need separate update/submit permissions.
- [x] Add matching constants to the frontend if missing.
- [x] Protect evaluation routes before API calls.
- [x] Protect create, update, submit, and history controls with `PermissionGuard`.
- [x] Do not use `CARE_COORDINATION.VIEW` as a substitute unless backend policy explicitly maps it.

### P1.6 Evaluation Statistics

- [ ] Define the intended scope of every KPI: assigned clients, current employee, organization, or active filters.
- [ ] Add a dedicated `GET /evaluations/stats` endpoint if global or mixed aggregates are required.
- [ ] Return stable typed counts and `as_of` metadata.
- [ ] Stop deriving totals from the current page's `results.length`.
- [ ] Use paginated response `count` only when its endpoint scope matches the displayed metric.

### P1.7 Date And Timezone Policy

- [ ] Choose the authoritative product timezone for evaluation calendar dates.
- [ ] Align PostgreSQL `CURRENT_DATE` and Go date calculations.
- [ ] Document date-only fields in the API contract.
- [ ] Prevent frontend local-time conversion from shifting evaluation dates.
- [ ] Verify the 14-day submission boundary around timezone and daylight-saving transitions.

### P1.8 Query And Listing Behavior

- [ ] Confirm that Upcoming, Drafts, and Submitted remain independently paginated listings.
- [ ] Add independent URL state if all three tables require navigation.
- [ ] Review lateral aggregate query plans at production-like volume.
- [ ] Add or adjust indexes based on measured query plans.

### P1.9 Direct Lifecycle Test Coverage

- [x] Add handler contract tests.
- [x] Add service validation tests.
- [x] Add repository/database lifecycle integration tests.
- [ ] Add frontend workflow tests for create, update, submit, conflict, and read-only states.

## P2: Validation, History, And Long-Term Hardening

### P2.1 Note Length Limits

- [ ] Decide maximum lengths for overall notes and goal notes.
- [ ] Enforce limits in backend request validation.
- [ ] Mirror limits in Valibot schemas and input attributes.
- [ ] Return field-specific errors.

### P2.2 Goal-Set Semantics

- [ ] Decide whether the required goal set is frozen when a draft is created.
- [ ] If frozen, persist an immutable evaluation goal snapshot.
- [ ] If dynamic, document that active-goal changes affect submission requirements.
- [ ] Define behavior for goals added, archived, or replaced during an active draft.

### P2.3 Archive Semantics

- [ ] Define who can archive an evaluation.
- [ ] Define valid transitions into `archived`.
- [ ] Define whether archived evaluations remain visible in history and statistics.
- [ ] Add API operations only after product rules are established.

### P2.4 Bootstrap And API Type Alignment

- [ ] Update frontend bootstrap types:

```ts
days_left: number | null;
priority: 'critical' | 'normal' | null;
```

- [ ] Review all optional and nullable evaluation fields against backend JSON.
- [ ] Keep date-only values distinct from timestamps.
- [ ] Remove frontend endpoint definitions that are not implemented, unless P0.2 implements them first.

### P2.5 Form And Request Validation

- [ ] Confirm that empty `items` is valid for draft creation.
- [ ] Keep duplicate goal validation in both service and repository layers.
- [ ] Reject inactive, foreign-client, and unknown goals.
- [ ] Ensure submission requires non-`no_progress` values for every required goal.
- [ ] Surface collection-level validation errors in the frontend.

### P2.6 Operational Observability

- [ ] Add structured audit events for create, update, blocked submit, successful submit, conflict, and archive.
- [ ] Include evaluation ID, client ID, cycle date, actor ID, and stable result code.
- [ ] Avoid logging note content or other sensitive payload data.
- [ ] Add metrics for blocked submissions and ownership conflicts.

## Cross-Repository Implementation Order

1. Complete and commit P0.1 bootstrap correction.
2. Finalize the P0.2 request/response and submission contract.
3. Implement backend ID-based update and submit operations.
4. Implement P0.3 historical-draft conflict enforcement.
5. Migrate the frontend form to the ID-based lifecycle.
6. Add stable error codes and localized frontend handling.
7. Add concurrency protection.
8. Align permissions and route/action gating.
9. Define and implement statistics.
10. Complete P2 validation and long-term hardening.

## Required End-To-End Scenarios

- [x] Bootstrap ignores a historical draft.
- [x] Bootstrap returns the current-cycle draft.
- [x] Bootstrap returns no draft without a scheduled cycle.
- [ ] Concurrent draft creation produces one current-cycle draft.
- [x] Saving draft A cannot modify draft B.
- [x] Historical drafts are read-only.
- [x] Another employee cannot update or submit the draft.
- [x] Completed evaluations cannot be updated or submitted again.
- [ ] Incomplete submission returns a stable error.
- [ ] Early submission returns a stable error.
- [x] Successful submission advances the schedule exactly once.
- [ ] Stale writes preserve the user's local input and show a conflict.
- [ ] Evaluation dates render consistently across supported timezones.
- [ ] KPI labels and values use the same documented scope.

## Verification Commands

Backend:

```bash
gofmt -w <changed-go-files>
go test ./...
go vet ./...
golangci-lint run
make test-integration
```

The integration suite starts a migrated PostgreSQL 17 container automatically. Set `TEST_DATABASE_URL` only to override the container with an already-migrated database for local debugging.

Frontend:

```bash
bun run check
bun run lint
```

Run the Svelte autofixer on every changed Svelte component before considering frontend work complete.

## Open Decisions

- [x] Submit transitions an already-saved revision; clients PATCH before POST when final edits exist.
- [ ] Does create return `201`, while returning an existing current draft uses `200`?
- [ ] Should current-cycle drafts owned by another employee be visible as read-only in bootstrap?
- [ ] Should blocked submission use HTTP 409 or 422?
- [ ] What exact scope should each dashboard KPI represent?
- [ ] What is the authoritative application timezone?
- [ ] Is the evaluation goal set frozen or dynamic?
- [ ] Are dedicated evaluation update and submit permissions required?
