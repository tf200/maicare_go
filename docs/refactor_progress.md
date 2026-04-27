# Refactor Progress

## 2026-04-27 - WebSocket infrastructure (Fragment 21)

Old: `api/websocket_router.go`, `api/websocket_handler.go`, `api/ws_ticket_manager.go`, `hub/hub.go`, `hub/client.go`
New: `internal/ws/ticket_manager.go`, `internal/handler/ws_handler.go`, `internal/handler/ws_routes.go`
Updated: `internal/ws/hub.go`, `internal/ws/client.go` (already existed, now canonical), `service/deps/deps.go`, `service/service.go`, `main.go`, `api/server.go`, `api/websocket_handler.go`
Removed: `hub/` package (merged into `internal/ws/`), `api/ws_ticket_manager.go`
Status: scaffolded, unwired
Notes: Fully migrated the WebSocket infrastructure into the internal stack. `internal/ws/` is now the canonical Hub implementation (was a dead duplicate, `hub/` was the live one — swapped). Ticket manager moved from `api/` to `internal/ws/ticket_manager.go`. New `WebSocketHandler` in `internal/handler/ws_handler.go` handles both `GET /ws` (ticket-authenticated upgrade) and `POST /auth/ws-ticket` (ticket creation). All 5 files importing `hub` updated to `ws`. Old `api/ws_ticket_manager.go` and `hub/` package deleted.

## 2026-04-27 - Settings routes (Fragment 20)

Old: `api/settings_handler.go`, `api/settings_router.go`, `service/settings/service.go`, `service/settings/dto.go`, `service/settings/interface.go`
New: `internal/domain/setting.go`, `internal/repository/setting_repo.go`, `internal/service/setting_service.go`, `internal/handler/setting_dto.go`, `internal/handler/setting_handler.go`, `internal/handler/setting_routes.go`
Status: scaffolded, unwired
Notes: migrated all 5 settings endpoints into the internal stack. Endpoints covered: GET/POST /settings/departments (SETTINGS.DEPARTMENT.VIEW/CREATE), PUT /settings/departments/:id (SETTINGS.DEPARTMENT.UPDATE), GET /settings/organization-profile (SETTINGS.ORGANIZATION_PROFILE.VIEW), PUT /settings/organization-profile (SETTINGS.ORGANIZATION_PROFILE.UPDATE). Split into two sub-domains: `DepartmentService`/`DepartmentRepository` and `OrganizationProfileService`/`OrganizationProfileRepository`. `ListDepartments` returns non-paginated `PageResponse` matching old response shape. `UpdateOrganizationProfile` validates name, timezone, email, and website in the repository layer (mirroring old service logic). `UpdateDepartment` validates non-empty name via `util.OtpString`.

## 2026-04-27 - Sender routes (Fragment 19)

Old: `api/sender_handler.go`, `api/sender_router.go`, `service/sender/sender.go`, `service/sender/sender_dto.go`
New: `internal/domain/sender.go`, `internal/repository/sender_repo.go`, `internal/service/sender_service.go`, `internal/handler/sender_handler.go`, `internal/handler/sender_dto.go`, `internal/handler/sender_routes.go`
Status: scaffolded, unwired
Notes: migrated all 6 sender endpoints into the internal stack. Endpoints covered: POST /senders (SENDER.CREATE), GET /senders (SENDER.VIEW), GET /senders/:id (SENDER.VIEW), PUT /senders/:id (SENDER.UPDATE), DELETE /senders/:id (SENDER.DELETE), POST /senders/:id/invoice_template (SENDER.CREATE). `GetSenderByID` fetches invoice template items via `GetTemplateItemsBySourceTable`. `UpdateSender` handles optional contacts marshaling. `CreateSenderInvoiceTemplate` validates template IDs before saving. Used `SenderEntity`/`SenderContactInfo` domain types to avoid collision with existing `domain.Sender`/`domain.SenderContact` in `client.go` (client network/referral sender).

## 2026-04-27 - Roles and permissions routes (Fragment 18)

Old: `api/roles_handler.go`, `api/roles_router.go`, `service/settings/roles.go`, `service/settings/roles_dto.go`
New: `internal/domain/role.go`, `internal/repository/role_repo.go`, `internal/service/role_service.go`, `internal/handler/role_handler.go`, `internal/handler/role_dto.go`, `internal/handler/role_routes.go`
Status: scaffolded, unwired
Notes: migrated all 8 role/permission endpoints into the internal stack. Endpoints covered: GET /roles (ROLES.VIEW), POST /roles (ROLES.CREATE), GET /roles/:role_id/permissions (PERMISSIONS.VIEW), POST /roles/:role_id/permissions (PERMISSIONS.CREATE), GET /permissions (PERMISSIONS.VIEW), POST /employees/:id/roles (ROLES.ASSIGN), GET /employees/:id/roles_permissions (PERMISSIONS.VIEW), POST /employees/:id/permissions (PERMISSIONS.GRANT). `ListAllPermissions` groups permissions by `GroupKey` and `SectionKey` with humanized labels. `ListUserRolesAndPermissions` returns role + inherited + override allows/denies + effective permissions. `ReplaceUserPermissionOverrides` validates no overlap between allow and deny lists, then deletes existing overrides and inserts new ones. `AddPermissionsToRole` replaces all permissions for a role (delete then insert). `CreateRole` creates a role with optional description. Used `SystemPermission` domain type to avoid collision with existing `domain.Permission` in `employee.go`.

## 2026-04-27 - Registration form routes (Fragment 17)

Old: `api/registration_form_handler.go`, `api/registration_form_router.go`, `service/client/registration_form.go`, `service/client/registration_form_dto.go`
New: `internal/domain/registration_form.go`, `internal/repository/registration_form_repo.go`, `internal/service/registration_form_service.go`, `internal/handler/registration_form_handler.go`, `internal/handler/registration_form_dto.go`, `internal/handler/registration_form_routes.go`
Status: scaffolded, unwired
Notes: migrated all 9 registration form endpoints into the internal stack. Endpoints covered: POST /registration_forms (public, no auth), GET/GET by ID/PUT/DELETE /registration_forms, POST /registration_forms/:id/status, POST /registration_forms/:id/process, GET/POST /public/intake-options/:token. `CreateRegistrationForm` accepts optional empty body. `ListRegistrationForms` supports pagination + status filter + 9 risk boolean filters. `GetRegistrationForm` fetches document attachments via `GetAttachmentsByUUIDs` and maps to `domain.Document`. `UpdateRegistrationForm` supports partial updates with all optional fields. `ProcessRegistrationForm` generates random token, marshals proposed dates to JSON, updates status to "processed", and enqueues process registration form email task. `GetPublicIntakeOptions` and `SelectIntakeDate` are public token-based endpoints. Added `domain.ErrRegistrationFormNotFound` to `internal/domain/registration_form.go`.

## 2026-04-27 - Notification routes (Fragment 16)

Old: `api/notification_handler.go`, `api/notification_router.go`, `service/notification/notification.go`, `service/notification/notification_dto.go`
New: `internal/domain/notification.go`, `internal/repository/notification_repo.go`, `internal/service/notification_service.go`, `internal/handler/notification_handler.go`, `internal/handler/notification_dto.go`, `internal/handler/notification_routes.go`
Status: scaffolded, unwired
Notes: migrated `GET /notifications` and `POST /notifications/:id/read` into the internal stack. Domain types: `Notification`, `NotificationData`, and nested data structs (appointment, client assignment, contract reminder, incident report, schedule notification). `ListNotifications` supports pagination via `httpapi.PageRequest`. `MarkNotificationAsRead` uses a transaction with ownership check (returns 403 if notification does not belong to user). The old `service/notification` package with `CreateAndDeliver` is left untouched since the Asynq worker (`internal/worker`) depends on it.

## 2026-04-27 - Maturity matrix endpoint (Fragment 15)

Old: `api/maturity_matrix_handler.go` → `ListMaturityMatrixApi`, `service/care/topics.go`
New: `internal/domain/maturity_matrix.go`, `internal/repository/maturity_matrix_repo.go`, `internal/service/maturity_matrix_service.go`, `internal/handler/maturity_matrix_handler.go`, `internal/handler/maturity_matrix_routes.go`
Status: scaffolded, unwired
Notes: migrated `GET /maturity_matrix` into the internal stack. Domain types: `MaturityMatrix`, `MaturityMatrixLevel`. Repository unmarshals `level_description` JSONB into structured levels. Service is a thin passthrough. Handler returns `httpapi.OK` wrapper. Permission: `CARE_PLAN.VIEW`.

## 2026-04-27 - Intake form routes scaffold (Fragment 14)

Old: `api/intake_form_handler.go`, `api/intake_form_router.go`, `service/client/intake_form.go`, `service/client/intake_form_dto.go`, `service/client/intake_maturity.go`, `service/client/intake_maturity_dto.go`, `service/client/promote_intake.go`, `service/client/promote_intake_dto.go`
New: `internal/domain/intake_form.go`, `internal/domain/ai.go`, `internal/domain/common.go`, `internal/repository/intake_form_repo.go`, `internal/service/intake_form_service.go`, `internal/handler/intake_form_handler.go`, `internal/handler/intake_form_dto.go`, `internal/handler/intake_form_routes.go`
Status: scaffolded, unwired
Notes: migrated all 9 intake form endpoints into a new standalone intake form domain stack. Endpoints covered: POST/GET/GET totals/GET by ID/PATCH /intake_forms, POST /intake_forms/:id/generate_goals, PUT /intake_forms/:id/goals, PATCH /intake_forms/:id/conclusion, POST /intake_forms/:id/promote. `CreateIntakeForm` validates registration form exists and is processed. `UpdateIntakeForm` supports clear_fields for nullable fields, validates self_sufficiency (0-5) and evaluation_intervals_weeks (>=0). `ReplaceIntakeFormGoals` runs in a transaction with intake form locking, checks for active client blocking. `GenerateIntakeGoals` supports both assessment-based and intake-form-based generation, builds risk context from registration form flags + intake risk assessment, calls AI service. `PromoteIntakeToClient` runs in a transaction: locks intake form, validates suitable conclusion, checks idempotency (existing client), creates client details from registration+intake data, migrates goals to client goals, creates emergency contacts from guardians. Added `domain.AIService` and `domain.CarePlanResponse` to `internal/domain/ai.go`. Added generic `domain.ListResult[T]` to `internal/domain/common.go`.

## 2026-04-27 - Employee profile details endpoint (Fragment 13)

Old: `api/employee_profile_handler.go` → `GetEmployeeProfileDetailsApi`
New: `internal/handler/employee_handler.go` → `GetEmployeeProfileDetails`, `internal/handler/employee_dto.go`, `internal/handler/employee_routes.go`
Status: scaffolded, unwired
Notes: migrated `GET /employees/profile/details` into the internal stack. `GetEmployeeProfileDetails` aggregates 7 SQLC queries (`GetEmployeeProfileByUserID`, `GetEmployeeProfileByID`, `GetUserRoles`, `ListActiveSessionsByUserID`, `ListEducations`, `ListEmployeeExperience`, `GetLocation`/`GetOrganisation`) in the repository layer to build a rich `EmployeeProfileDetails` domain struct. Added domain types: `EmployeeProfileDetails`, `EmployeeRole`, `BriefEducationDetail`, `BriefExperienceDetail`, `ActiveSessionDetail`. Fixed pre-existing bug in `EmployeeRepository.GetEmployeeByUserID` where `toDomainEmployeeProfile` (which returns `(*EmployeeProfile, error)`) was incorrectly wrapped with `, nil`.

## 2026-04-27 - Contract routes scaffold (Fragment 12)

Old: `api/contract_handler.go`, `service/contract/contract.go`, `service/contract/contract_dto.go`
New: `internal/domain/contract.go`, `internal/repository/contract_repo.go`, `internal/service/contract_service.go`, `internal/handler/contract_handler.go`, `internal/handler/contract_dto.go`
Status: scaffolded, unwired
Notes: migrated contract endpoints into a new standalone contract domain stack. Endpoints covered: POST/GET/DELETE /contract_types, POST/GET/PUT/PUT status/GET audit /contracts, GET /clients/:id/contracts. `CreateContract` validates care_name, price, date range, care_type pricing rules, and attachment IDs. `UpdateContract` re-fetches existing contract for COALESCE-style patching, validates care_type cannot change, and re-validates attachments. `UpdateContractStatus` prevents approving expired contracts. `ListContracts` supports search, status, care_type, financing_act, financing_option, and end_date range filters. `UpdateContract` and `UpdateContractStatus` use raw `ConnPool.Begin` tx with `SET LOCAL myapp.current_employee_id` for audit logging (matches old behavior). Attachment validation (`GetAttachmentFiles`) lives in repository; presigned URL generation lives in service via `domain.Storage`. Renamed `ClientContractSummary` to `ClientContractListItem` and `clientContractSummaryResponse` to `contractClientListItemResponse` to avoid collisions with existing client domain types.

## 2026-04-27 - Client progress report routes scaffold (Fragment 11)

Old: `api/client_progress_report_router.go`, `api/client_progress_report_handler.go`, `service/client/reports.go`, `service/client/reports_dto.go`
New: `internal/domain/client.go`, `internal/repository/client_repo.go`, `internal/service/client_service.go`, `internal/handler/client_handler.go`, `internal/handler/client_dto.go`
Status: scaffolded, unwired
Notes: migrated client progress report endpoints into the existing client domain stack. Endpoints covered: POST/GET/GET/PUT/DELETE /clients/:id/progress_reports, POST /clients/:id/ai_progress_reports, POST /clients/:id/ai_progress_reports/confirm, GET /clients/:id/ai_progress_reports. `CreateProgressReport` validates report_text, type, and emotional_state via binding tags. `GenerateAutoReports` fetches reports by date range, builds text summary, and delegates to `domain.AutoReportGenerator.GenerateAutoReports`. `ConfirmAiProgressReport` creates an `AiGeneratedReport` record. `ListProgressReports` supports type filter via `NullProgressReportTypeFromPtr`. `UpdateProgressReport` uses `NullProgressReportTypeFromPtr` and `NullEmotionalStateFromPtr` for optional enum fields. Added `AutoReportGenerator` interface to domain layer to decouple from gRPC client. Added `reportGen` field to `ClientService` struct and updated `NewClientService` constructor. Fixed pre-existing repository gaps: added `UpdateClientStatus`, `PutClientInCare`, `PutClientOutOfCare` methods that were required by `domain.ClientRepository` interface but never implemented. Fixed `*db.ClientStatusEnum` to `*string` type mismatch in status history creation.

## 2026-04-27 - Client network routes scaffold (Fragment 10)

Old: `api/client_network_router.go`, `api/client_network_handler.go`, `service/client/network.go`, `service/client/network_dto.go`
New: `internal/domain/client.go`, `internal/repository/client_repo.go`, `internal/service/client_service.go`, `internal/handler/client_handler.go`, `internal/handler/client_dto.go`
Status: scaffolded, unwired
Notes: migrated client network endpoints into the existing client domain stack. Endpoints covered: GET /clients/:id/sender, POST/GET/GET/PUT/DELETE /clients/:id/emergency_contacts, POST/GET/GET/PUT/DELETE /clients/:id/involved_employees, GET /clients/:id/related_emails. `GetClientSender` unmarshals JSON contacts from `db.Sender.Contacts`. `CreateAssignedEmployee` enqueues a `new_client_assignment` notification to the assigned employee's user ID (non-fatal if enqueue fails). `ListClientEmergencyContacts` supports search parameter. `ListAssignedEmployees` constructs employee name from first + last name. All enum conversions for `RelationStatusEnum` use existing `db.NullRelationStatusFromPtr` / `db.RelationStatusPtrFromEnum` helpers. Removed duplicate minimal `ClientEmergencyContact` struct from domain that was leftover from Fragment 8. Fixed broken `toDomainClientGoal` duplicate and missing `toDomainClientDiagnosis` helper in repository.

## 2026-04-27 - Client medical routes scaffold (Fragment 9)

Old: `api/client_medical_router.go`, `api/client_medical_handler.go`, `service/client/medical.go`, `service/client/medical_dto.go`
New: `internal/domain/client.go`, `internal/repository/client_repo.go`, `internal/service/client_service.go`, `internal/handler/client_handler.go`, `internal/handler/client_dto.go`
Status: scaffolded, unwired
Notes: migrated client medical endpoints into the existing client domain stack. Endpoints covered: GET /clients/:id/medical/overview, POST /clients/:id/medical/diagnoses, GET /clients/:id/medical/diagnoses, GET /clients/:id/medical/diagnoses/:diagnosis_id, PUT /clients/:id/medical/diagnoses/:diagnosis_id, DELETE /clients/:id/medical/diagnoses/:diagnosis_id, POST /clients/:id/medical/medication-orders, GET /clients/:id/medical/medication-orders, GET /clients/:id/medical/medication-orders/:order_id, PUT /clients/:id/medical/medication-orders/:order_id, DELETE /clients/:id/medical/medication-orders/:order_id. `GetClientMedicalOverview` fetches diagnoses + active medication orders in a transaction with hardcoded limit 500. `CreateClientDiagnosis` defaults status to "confirmed" and severity to "unknown". `UpdateClientMedicationOrder` re-fetches enriched view after update for consistent response. `ListClientMedicationOrders` supports filtering by status, admin_mode, diagnosis_id, and search. All medical routes registered under `/clients/:id/medical/*` in `RegisterClientRoutes`. Enum conversions for `DiagnosisStatusEnum`, `DiagnosisSeverityEnum`, `MedicationOrderStatusEnum`, `MedicationAdminModeEnum` handled inline in repository (no helpers exist in `enum_helpers.go`).

## 2026-04-27 - Incident routes scaffold (Fragment 8)

Old: `api/client_incident_router.go`, `api/incidents_router.go`, `api/client_incident_handler.go`, `api/incidents_handler.go`, `service/client/incidents.go`, `service/client/incidents_dto.go`
New: `internal/domain/incident.go`, `internal/repository/incident_repo.go`, `internal/service/incident_service.go`, `internal/handler/incident_handler.go`, `internal/handler/incident_dto.go`
Status: scaffolded, unwired
Notes: migrated all incident endpoints into a new `incidents` domain. Endpoints covered: POST /incidents, GET /incidents, GET /incidents/counts, GET /incidents/:incident_id, PUT /incidents/:incident_id, DELETE /incidents/:incident_id, GET /incidents/:incident_id/file, PUT /incidents/:incident_id/confirm, GET /clients/:id/incidents. `RegisterIncidentRoutes` registers both `/clients/:id/incidents` and `/incidents/*` routes. `CreateIncident` sends notifications to all admin users via task queue (non-fatal if enqueue fails). `GenerateIncidentFile` maps domain incident to `IncidentPDFData` and delegates to `domain.IncidentPDFGenerator`. `ConfirmIncident` enqueues confirmation email when rowsAffected > 0. `ListAllIncidents` uses transaction for list + count queries. Repository uses `store.ExecTx` for write operations and `queries` directly for reads. All enum conversions handled inline in repository. Handler uses `getUserIDFromContext` helper for confirm endpoint.

## 2026-04-27 - Top-level evaluations routes scaffold (Fragment 7)

Old: `api/client_details_router.go`, `api/client_details_handler.go`, `service/client/evaluation.go`, `service/client/evaluation_dto.go`
New: `internal/domain/client.go`, `internal/repository/client_repo.go`, `internal/service/client_service.go`, `internal/handler/client_handler.go`, `internal/handler/client_dto.go`
Status: scaffolded, unwired
Notes: migrated top-level evaluation endpoints (GET /evaluations/:evaluation_id, GET /evaluations/upcoming, GET /evaluations/recent-submitted, GET /evaluations/recent-drafts) into the internal stack. `GetGoalEvaluation` fetches evaluation header plus items by ID, returns 404 if not found. `ListUpcomingEvaluations` uses coordinator employee ID from auth context. `ListRecentSubmittedEvaluations` and `ListRecentDraftEvaluations` filter by logged-in employee ID. Added separate `RegisterEvaluationRoutes` function for `/evaluations` group since these routes are top-level, not under `/clients`. Reused existing `domain.GoalEvaluation` type and `toDomainGoalEvaluation` repository helper. All list endpoints use `httpapi.PageRequest.Params()` for pagination.

## 2026-04-27 - Client location transfer routes scaffold (Fragment 6)

Old: `api/client_details_router.go`, `api/client_details_handler.go`, `service/client/location_transfer.go`, `service/client/location_transfer_dto.go`
New: `internal/domain/client.go`, `internal/repository/client_repo.go`, `internal/service/client_service.go`, `internal/handler/client_handler.go`, `internal/handler/client_dto.go`
Status: scaffolded, unwired
Notes: migrated location transfer endpoints (POST /clients/:id/location_transfer, POST /clients/location_transfer/approve_reject, GET /clients/location_transfer) into the internal stack. `CreateLocationTransfer` uses `time.Now()` for request_date and stores the transfer record. `ApproveLocationTransfer` validates status is "approved" or "rejected" in service, then updates status and sets approved_rejected_at/approved_rejected_by in repository. `ListLocationTransferRequests` returns paginated results with mentor name joins. Handler uses `getEmployeeIDFromContext` for the approve/reject endpoint. All three endpoints follow the existing pagination pattern with `httpapi.PageRequest.Params()`.

## 2026-04-27 - Client goals & evaluations routes scaffold (Fragment 5)

Old: `api/client_details_router.go`, `api/client_details_handler.go`, `service/client/client_goals.go`, `service/client/evaluation.go`
New: `internal/domain/client.go`, `internal/repository/client_repo.go`, `internal/service/client_service.go`, `internal/handler/client_handler.go`, `internal/handler/client_dto.go`
Status: scaffolded, unwired
Notes: migrated goal and evaluation endpoints (GET /clients/:id/evaluations/bootstrap, POST /clients/:id/goals, PATCH /clients/:id/goals/:goal_id, GET /clients/:id/goals, GET /clients/:id/goals/:goal_id/history, GET /clients/:id/evaluations/submitted, POST /clients/:id/evaluations) into the internal stack. `CreateClientGoal` validates title and delegates to repository transaction. `UpdateClientGoal` validates empty patch, handles domain error mapping (not found, conflict). `CreateGoalEvaluation` validates item duplicates and progress values in service, repository handles complex draft/submit transaction. `GetGoalEvaluationBootstrap` computes days-left and priority. `ListClientSubmittedEvaluations` and `ListGoalEvaluationHistory` are paginated. Handler uses `getEmployeeIDFromContext` helper to extract employee ID from gin context.

## 2026-04-27 - Client status transition routes scaffold (Fragment 3)

Old: `api/client_details_router.go`, `api/client_details_handler.go`, `service/client/client_details.go`, `service/client/put_in_care.go`, `service/client/put_out_of_care.go`
New: `internal/domain/client.go`, `internal/repository/client_repo.go`, `internal/service/client_service.go`, `internal/handler/client_handler.go`, `internal/handler/client_dto.go`
Status: scaffolded, unwired
Notes: migrated client status endpoints (PUT /clients/:id/status, PUT /clients/:id/put-in-care, PUT /clients/:id/put-out-of-care, GET /clients/:id/status_history) into the internal stack. All write operations use transactions via `db.Store.ExecTx`. Service layer preserves date parsing, target-status determination, discharge-reason validation, and goal-existence validation. Scheduled status updates are stubbed (matching old commented-out behavior). `ClientService` now depends on `domain.TaskQueue` for future notification support.

## 2026-04-27 - Client detail & addresses routes scaffold (Fragment 2)

Old: `api/client_details_router.go`, `api/client_details_handler.go`, `service/client/client_details.go`
New: `internal/domain/client.go`, `internal/repository/client_repo.go`, `internal/service/client_service.go`, `internal/handler/client_handler.go`, `internal/handler/client_dto.go`
Status: scaffolded, unwired
Notes: migrated client detail endpoints (GET /clients/:id, PUT /clients/:id, GET /clients/:id/addresses) into the internal stack. GET /clients/:id now aggregates 11 SQLC queries in the repository with full business logic porting (age calculation, risk flags, alerts, status-conditional care/discharge/evaluation/contract summaries). PUT and addresses are stubs matching the old service stubs. Repository now depends on `*db.Store` for future transaction support.

## 2026-04-27 - Client core routes scaffold (Fragment 1)

Old: `api/client_details_router.go`, `api/client_details_handler.go`, `service/client/client_details.go`, `service/client/client_details_dto.go`, `db/sqlc/client.sql.go`
New: `internal/domain/client.go`, `internal/repository/client_repo.go`, `internal/service/client_service.go`, `internal/handler/client_handler.go`, `internal/handler/client_dto.go`
Status: scaffolded, unwired
Notes: migrated client core endpoints (POST /clients, GET /clients, GET /clients/waiting-list, GET /clients/in-care, GET /clients/counts, GET /clients/status-counts) into the internal stack with domain-owned contracts, repository SQLC mapping via `pkg/conv`, service-level pagination defaults, and handler-local DTO mapping

## 2026-04-02 - Late arrival routes scaffold
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
