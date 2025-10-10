# API Refactoring Rules

This document outlines the rules for refactoring API handlers to separate controllers from business logic.

## Overview
The goal is to separate HTTP handling concerns from business logic by:
- Moving business logic from `api/` handlers to `service/` layer
- Keeping only HTTP-specific logic in `api/` handlers
- Maintaining identical request/response types
- Preserving all existing functionality

## Migration Pattern

### 1. File Structure
```
api/
├── [feature]_handler.go     # Controllers only (HTTP handling)
service/[feature]/
├── [feature].go            # Business logic implementation
├── [feature]_dto.go        # Request/Response types
```

### 2. Controller Responsibilities (`api/[feature]_handler.go`)
- **HTTP Request handling**: Parse path parameters, query params, JSON body
- **Response formatting**: Create SuccessResponse/errorResponse and set HTTP status codes
- **Input validation**: Basic binding validation using `ctx.ShouldBindJSON()`
- **HTTP status codes**: Set appropriate status codes (200, 201, 400, 404, 500)
- **Error handling**: Convert business logic errors to HTTP responses using `errorResponse()`

**Controller functions should:**
- Have names ending with `Api` (e.g., `CreateProgressReportApi`)
- Accept `*gin.Context` as parameter
- Extract and parse request data
- Call corresponding business logic method
- Format and return HTTP response
- Be thin and delegate all business logic to service layer

### 3. Business Logic Responsibilities (`service/[feature]/[feature].go`)
- **Core business logic**: All business rules and data processing
- **Database operations**: All database calls via store
- **Data transformation**: Convert between database models and DTOs
- **Business validation**: Complex validation beyond basic input checking
- **Logging**: Business event logging with appropriate context
- **Notifications**: Trigger notifications or other side effects

**Business logic functions should:**
- Have descriptive names without `Api` suffix (e.g., `CreateProgressReport`)
- Accept `context.Context` and business parameters
- Return business response types and errors
- Handle all database operations
- Implement business rules and validation
- Log business events appropriately

### 4. DTO Types (`service/[feature]/[feature]_dto.go`)
- **Request types**: All input data structures with validation tags
- **Response types**: All output data structures with JSON tags
- **Pagination types**: Include pagination.Request in list request types
- **No business logic**: Pure data structures only

## Refactoring Rules

### Rule 1: Preserve Request/Response Types
- All request and response types must remain exactly the same
- Types should be moved to `service/[feature]/[feature]_dto.go`
- Import path changes: `api` package imports from `service/[feature]`

### Rule 2: Controller Signature Pattern
```go
// Before (mixed logic)
func (server *Server) CreateFeature(ctx *gin.Context) {
    // HTTP handling + business logic mixed
}

// After (separated)
func (server *Server) CreateFeatureApi(ctx *gin.Context) {
    // Only HTTP handling
}

func (s *featureService) CreateFeature(ctx context.Context, req CreateFeatureRequest, id int64) (*CreateFeatureResponse, error) {
    // Only business logic
}
```

### Rule 3: Error Handling Pattern
- Controllers: Use `errorResponse(err)` for all errors
- Business logic: Return errors with context, don't handle HTTP directly
- Logging: Business logic handles structured logging

### Rule 4: Parameter Extraction
- Path parameters: Extract in controller, pass to business logic
- Query parameters: Parse in controller using ShouldBindQuery()
- JSON body: Parse in controller using ShouldBindJSON()
- Validation: Basic validation in controller, business validation in service

### Rule 5: Response Pattern
```go
// Controller pattern
result, err := server.businessService.FeatureService.CreateFeature(ctx, req, id)
if err != nil {
    ctx.JSON(http.StatusInternalServerError, errorResponse(err))
    return
}
res := SuccessResponse(result, "Feature created successfully")
ctx.JSON(http.StatusCreated, res)
```

### Rule 6: Import Organization
- Controllers import business service packages
- Business logic imports database, logger, notification packages
- DTO files should have minimal dependencies

## Migration Steps

1. **Create service structure**: Make `service/[feature]/` directory
2. **Move types**: Extract request/response types to `_dto.go`
3. **Create business logic**: Extract business logic to `[feature].go`
4. **Update controller**: Refactor handler to only handle HTTP concerns
5. **Update imports**: Fix all import paths
6. **Verify functionality**: Ensure API behavior is identical

## Example Migration

### Before (`api/client_progress_report_handler.go`)
```go
func (server *Server) CreateProgressReportApi(ctx *gin.Context) {
    // Parameter extraction
    id := ctx.Param("id")
    clientID, err := strconv.ParseInt(id, 10, 64)
    // ... validation

    // Business logic mixed with HTTP handling
    var req clientp.CreateProgressReportRequest
    // ... database calls
    // ... business rules
    // ... data transformation

    // Response formatting
    res := SuccessResponse(progressReport, "Progress Report created successfully")
    ctx.JSON(http.StatusCreated, res)
}
```

### After
**Controller (`api/client_progress_report_handler.go`)**:
```go
func (server *Server) CreateProgressReportApi(ctx *gin.Context) {
    // Parameter extraction
    id := ctx.Param("id")
    clientID, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    // Request parsing
    var req clientp.CreateProgressReportRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    // Business logic delegation
    progressReport, err := server.businessService.ClientService.CreateProgressReport(ctx, &req, clientID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    // Response formatting
    res := SuccessResponse(progressReport, "Progress Report created successfully")
    ctx.JSON(http.StatusCreated, res)
}
```

**Business Logic (`service/client/progress_report.go`)**:
```go
func (s *clientService) CreateProgressReport(ctx context.Context, req *CreateProgressReportRequest, clientID int64) (*CreateProgressReportResponse, error) {
    // Database operations
    arg := db.CreateProgressReportParams{...}
    report, err := s.Store.CreateProgressReport(ctx, arg)
    if err != nil {
        s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateProgressReport", "Failed to create progress report", zap.Int64("client_id", clientID), zap.Error(err))
        return nil, err
    }

    // Data transformation
    response := &CreateProgressReportResponse{...}
    return response, nil
}
```

**DTOs (`service/client/progress_report_dto.go`)**:
```go
type CreateProgressReportRequest struct {
    // Request fields with validation tags
}

type CreateProgressReportResponse struct {
    // Response fields with JSON tags
}
```

## Key Principles

1. **No breaking changes**: API contracts must remain identical
2. **Single responsibility**: Each function has one clear purpose
3. **Separation of concerns**: HTTP handling vs business logic
4. **Testability**: Business logic can be tested without HTTP context
5. **Maintainability**: Clear organization and dependency flow