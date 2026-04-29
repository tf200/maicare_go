# Maicare Go - Refactoring Guidelines & Workflow

## Objective
Refactor the current codebase, focusing primarily on the `service` layer methods, to improve performance, code quality, and project structure.

## Core Focus Areas

1.  **Performance**
    *   Identify and eliminate redundant code.
    *   Remove useless or duplicated calls.
    *   Detect and fix N+1 query problems.
    *   Identify poorly optimized SQL queries (referencing the `db/query` and `db/sqlc` packages).
    *   Address any other bottlenecks that hinder execution speed or increase resource usage.
    *   Remove unecessary or redundant queries from the `service` layer like exrea fetches that can just be removed
2.  **Code Quality & Consistency**
    *   Ensure the code is well-written, idiomatic Go, and highly consistent across the project.
    *   Extract general helper methods and move them to the `util` package.
    *   Ensure that code in the `service` layer is strictly specific to business logic; decouple unrelated concerns.
    *   Verify that Data Transfer Objects (DTOs) and request structs have proper validation in place.
    *   **Logging:** Replace verbose `Logger.LogBusinessEvent` calls with the simplified `Logger.LogError`, `Logger.LogWarn`, and `Logger.LogInfo` methods to declutter business logic.

3.  **Folder Structure**
    *   Identify and fix any bad or non-standard folder structures within the project to maintain a clean domain or modular layout.

## Workflow Protocol

1.  **Selection:** The user provides a specific file (typically in the `service/` directory) to be verified.
2.  **Analysis:** The assistant reviews the file method-by-method against the core focus areas (Performance, Code Quality, Folder Structure).
3.  **Reporting:** The assistant reports the findings, clearly highlighting issues, bottlenecks, and areas for improvement. **No code changes are made at this stage.**
4.  **Direction:** The user reviews the findings and instructs the assistant on which specific issues to fix.
5.  **Execution:** The assistant implements the requested fixes, adhering to the project's standards.

## Tracked Files

- [x] `service/auth/auth.go`
- [x] `service/auth/roles.go`
- [x] Invoice service refactored from `service/invoice/` → `internal/service/invoice_service.go`, `internal/domain/invoice.go`, `internal/handler/invoice_handler.go`, `internal/handler/invoice_dto.go`, `internal/handler/invoice_routes.go`, wired in `internal/app/app.go`
