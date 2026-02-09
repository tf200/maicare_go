# Project Context: Care Coordination App (Backend)

## Goal
This is the backend for a care coordination application. We are currently modifying the codebase to adapt it to specific requirements, proceeding table by table and service by service.

## Current Status

### Verified Modules
The following components have been reviewed and verified/modified to fit our current needs:
- **Authentication & Authorization**: User roles, permissions, and overall auth logic.
- **Infrastructure**: Organizations and Locations.
- **Entities**: Senders and Employee profiles.
- **Intake Flow**:
    - Registration
    - Planning
    - Intake (including maturity matrix/goals)
    - Decision

### Work in Progress: Client Management
We are now focusing on client management:
- **Client Profiles**: Core details, BSN verification, and contact info.
- **Waitlist to Care Lifecycle**:
    - Status: completed.
    - We completed the transition from waiting list into care with explicit lifecycle states and strict DB guardrails.
    - Current lifecycle statuses are machine-friendly and defined as:
        - `on_waiting_list`
        - `scheduled_in_care`
        - `in_care`
        - `scheduled_out_of_care`
        - `out_of_care`
    - A client moving toward care must capture:
        - `placed_in_care_at`
        - `care_start_date`
        - `next_evaluation_date` (derived from start date + evaluation interval)
    - Main coordinator ownership is represented through `assigned_employee` with role `coordinator` (single coordinator per client).
    - Clients cannot be marked as `scheduled_in_care` or `in_care` without at least one active goal in `client_goals`.
- **In-Care Flow**:
    - Status: completed.
    - In-care client listing now supports:
        - statuses `in_care` and `scheduled_in_care`
        - search by client name
        - status filtering
        - sorting by days in care (asc/desc)
        - active contract flag (approved and currently active)

### Next Step
- **Client Dossier**:
    - Build and verify the complete client dossier retrieval flow.
    - Focus on the client detail GET endpoint and aggregation of all related client data.

## Technical Notes

### Database
- **Migration Policy**: We are in early development and currently allow direct edits to `db/migrations/000001_init.up.sql` for schema iteration.
- **Source of Truth**: The primary schema definition is maintained in `db/migrations/000001_init.up.sql`.
- **Status Scheduling Note**: `scheduled_status_changes` is treated as legacy for the waitlist->care start flow. The lifecycle state is now represented directly on `client_details.status`.

### Testing
- **Current State**: All old test files are considered obsolete and are currently broken due to significant structural changes.
- **Policy**: We are not fixing these tests yet, as the codebase is still highly volatile and undergoing rapid modifications. Tests will be addressed once the core logic stabilizes.

## How to Work on This Codebase
- Refer to `db/migrations/000001_init.up.sql` for the latest schema.
- Follow the service-based architecture in `service/` and handler patterns in `api/`.
- Prioritize client management implementation, especially in-care listing and related operational workflows.
