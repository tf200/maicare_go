# Database Test Structure Documentation

## Overview
This document provides a template and guidelines for writing database tests in the `maicare_go` project. All database tests should follow this consistent structure to ensure maintainability and reliability.

## File Structure
- Test files should be named `<entity>_test.go` (e.g., `custom_user_test.go`)
- Place test files in the `db/sqlc/` directory alongside the generated code
- Import required packages at the top of the file

## Required Imports
```go
package db

import (
    "context"
    "testing"

    "maicare_go/util"

    "github.com/google/uuid"
    "github.com/stretchr/testify/require"
)
```

## Test Function Structure

### 1. Table-Driven Tests Pattern
All tests MUST use the table-driven test pattern with the following structure:

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name   string
        setup  func(ctx context.Context, qtx *Queries) <ReturnType>
        checks func(t *testing.T, <result>, err error)
    }{
        {
            name: "descriptive test case name",
            setup: func(ctx context.Context, qtx *Queries) <ReturnType> {
                // Setup test data and return parameters
            },
            checks: func(t *testing.T, <result>, err error) {
                // Verify results using require assertions
            },
        },
        // Additional test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx := context.Background()

            tx, err := testDB.Begin(ctx)
            require.NoError(t, err, "failed to begin transaction")
            defer tx.Rollback(ctx)

            qtx := testQueries.WithTx(tx)

            // Execute test-specific logic
            tt.checks(t, result, err)
        })
    }
}
```

### 2. Test Cases to Include
For each database function, include AT LEAST these test scenarios:

#### For CREATE operations:
- ✅ Successful creation with minimal required fields
- ✅ Successful creation with all optional fields populated
- ✅ Creation with edge cases (empty strings, nil values, defaults)
- ❌ (Optional) Invalid data that should fail

#### For UPDATE operations:
- ✅ Successful update with valid data
- ✅ Update with empty/nil values where applicable
- ⚠️ Update for non-existent entity (should not error but affect 0 rows)

#### For GET/READ operations:
- ✅ Successfully retrieve existing entity
- ❌ Error when entity doesn't exist
- ✅ Retrieve with all fields populated correctly
- ✅ Retrieve when optional fields are nil

#### For DELETE operations:
- ✅ Successfully delete existing entity
- ⚠️ Delete non-existent entity (should not error)

## Test Structure Components

### 1. Test Case Definition
```go
{
    name: "descriptive name that explains what is being tested",
    setup: func(ctx context.Context, qtx *Queries) <ParamType> {
        // Create any prerequisite data
        // Return the parameters needed for the test
    },
    checks: func(t *testing.T, <result>, err error) {
        // Verify the outcome
    },
}
```

### 2. Setup Function
The `setup` function should:
- Create any prerequisite data using helper functions
- Use `createRandomEntity()` helpers for dependent entities
- Return the exact parameters needed for the function being tested
- Use `require.NoError()` for setup operations that must succeed

Example:
```go
setup: func(ctx context.Context, qtx *Queries) CreateUserParams {
    return CreateUserParams{
        Password:       "hashedpassword123",
        Email:          util.RandomEmail(),
        IsActive:       true,
        ProfilePicture: nil,
    }
}
```

For functions that need existing entities:
```go
setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
    user := createRandomUser(ctx, qtx)
    return user.ID
}
```

### 3. Checks Function
The `checks` function should:
- Use `require` package for all assertions
- Always include descriptive messages in assertions
- Check for expected errors with `require.Error()` or `require.NoError()`
- Verify all important fields of returned data
- Check nil/non-nil status of optional fields

Examples:
```go
// Success case
checks: func(t *testing.T, user CustomUser, err error) {
    require.NoError(t, err, "CreateUser() should not error")
    require.Equal(t, "test@example.com", user.Email)
    require.True(t, user.IsActive)
    require.Nil(t, user.ProfilePicture)
}

// Error case
checks: func(t *testing.T, user CustomUser, err error) {
    require.Error(t, err, "GetUserByID() should error for non-existent ID")
}

// Non-nil field check
checks: func(t *testing.T, user CustomUser, err error) {
    require.NoError(t, err)
    require.NotNil(t, user.ProfilePicture)
    require.Equal(t, "https://example.com/pic.jpg", *user.ProfilePicture)
}
```

### 4. Transaction Management
**CRITICAL**: Every test MUST use transactions that are rolled back:

```go
ctx := context.Background()

tx, err := testDB.Begin(ctx)
require.NoError(t, err, "failed to begin transaction")
defer tx.Rollback(ctx)

qtx := testQueries.WithTx(tx)
```

This ensures:
- Tests don't interfere with each other
- Database remains clean after tests
- Tests can run in parallel safely

## Helper Functions

### Purpose
Helper functions create random valid entities for use in tests. Place them at the bottom of the test file.

### Structure
```go
// Random data generators for tests

func createRandomEntity(ctx context.Context, qtx *Queries) Entity {
    entity, err := qtx.CreateEntity(ctx, CreateEntityParams{
        Field1: util.RandomString(10),
        Field2: util.RandomEmail(),
        Field3: true,
        // Use sensible defaults
    })
    if err != nil {
        panic("failed to create random entity: " + err.Error())
    }
    return entity
}
```

### Guidelines for Helpers:
- Name: `createRandom<EntityName>`
- Accept: `ctx context.Context, qtx *Queries`
- Return: The created entity
- Use `panic()` on error (setup failures should stop the test)
- Use `util.Random*()` functions for generating random data
- Keep it simple with valid default values

## Naming Conventions

### Test Functions
- Format: `Test<FunctionName>` (e.g., `TestCreateUser`, `TestGetUserByEmail`)
- Match the exact name of the function being tested

### Test Cases
Use descriptive names that explain the scenario:
- ✅ "successful creation"
- ✅ "get existing user by email"
- ✅ "update password for non-existent user"
- ✅ "with profile picture"
- ✅ "user inactive by default"
- ❌ "test 1", "case 2" (too vague)

## Assertion Guidelines

### Use `require` not `assert`
- `require` stops test execution on failure
- `assert` continues (can lead to panics/confusing errors)

### Always Include Messages
```go
// ✅ Good
require.NoError(t, err, "CreateUser() should not error")
require.Equal(t, expected, actual, "email should match")

// ❌ Bad
require.NoError(t, err)
require.Equal(t, expected, actual)
```

### Common Assertions
```go
// Error checking
require.NoError(t, err, "function should not error")
require.Error(t, err, "function should error for invalid input")

// Equality
require.Equal(t, expected, actual, "description")
require.NotEqual(t, unexpected, actual, "description")

// Nil checks
require.Nil(t, value, "should be nil")
require.NotNil(t, value, "should not be nil")

// Boolean checks
require.True(t, condition, "should be true")
require.False(t, condition, "should be false")

// Empty checks
require.Empty(t, slice, "should be empty")
require.NotEmpty(t, value, "should not be empty")
```

## Complete Example Template

```go
package db

import (
    "context"
    "testing"

    "maicare_go/util"

    "github.com/google/uuid"
    "github.com/stretchr/testify/require"
)

func TestCreateEntity(t *testing.T) {
    tests := []struct {
        name   string
        params CreateEntityParams
        checks func(t *testing.T, entity Entity)
    }{
        {
            name: "successful creation",
            params: CreateEntityParams{
                Field1: "value1",
                Field2: "value2",
                Field3: true,
            },
            checks: func(t *testing.T, entity Entity) {
                require.Equal(t, "value1", entity.Field1)
                require.Equal(t, "value2", entity.Field2)
                require.True(t, entity.Field3)
            },
        },
        {
            name: "with optional field",
            params: CreateEntityParams{
                Field1:         "value1",
                Field2:         "value2",
                Field3:         true,
                OptionalField:  util.StringPtr("optional"),
            },
            checks: func(t *testing.T, entity Entity) {
                require.NotNil(t, entity.OptionalField)
                require.Equal(t, "optional", *entity.OptionalField)
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx := context.Background()

            tx, err := testDB.Begin(ctx)
            require.NoError(t, err, "failed to begin transaction")
            defer tx.Rollback(ctx)

            qtx := testQueries.WithTx(tx)

            entity, err := qtx.CreateEntity(ctx, tt.params)
            require.NoError(t, err, "CreateEntity() should not error")

            tt.checks(t, entity)
        })
    }
}

func TestGetEntityByID(t *testing.T) {
    tests := []struct {
        name   string
        setup  func(ctx context.Context, qtx *Queries) uuid.UUID
        checks func(t *testing.T, entity Entity, err error)
    }{
        {
            name: "get existing entity by ID",
            setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
                entity := createRandomEntity(ctx, qtx)
                return entity.ID
            },
            checks: func(t *testing.T, entity Entity, err error) {
                require.NoError(t, err, "GetEntityByID() should not error")
                require.NotEmpty(t, entity.ID)
            },
        },
        {
            name: "get non-existent entity by ID",
            setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
                return uuid.New()
            },
            checks: func(t *testing.T, entity Entity, err error) {
                require.Error(t, err, "GetEntityByID() should error for non-existent ID")
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx := context.Background()

            tx, err := testDB.Begin(ctx)
            require.NoError(t, err, "failed to begin transaction")
            defer tx.Rollback(ctx)

            qtx := testQueries.WithTx(tx)

            id := tt.setup(ctx, qtx)
            entity, err := qtx.GetEntityByID(ctx, id)
            tt.checks(t, entity, err)
        })
    }
}

func TestUpdateEntity(t *testing.T) {
    tests := []struct {
        name   string
        setup  func(ctx context.Context, qtx *Queries) UpdateEntityParams
        checks func(t *testing.T, err error)
    }{
        {
            name: "successful update",
            setup: func(ctx context.Context, qtx *Queries) UpdateEntityParams {
                entity := createRandomEntity(ctx, qtx)
                return UpdateEntityParams{
                    ID:     entity.ID,
                    Field1: "updated_value",
                }
            },
            checks: func(t *testing.T, err error) {
                require.NoError(t, err, "UpdateEntity() should not error")
            },
        },
        {
            name: "update non-existent entity",
            setup: func(ctx context.Context, qtx *Queries) UpdateEntityParams {
                return UpdateEntityParams{
                    ID:     uuid.New(),
                    Field1: "value",
                }
            },
            checks: func(t *testing.T, err error) {
                require.NoError(t, err, "UpdateEntity() should not error even for non-existent entity")
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx := context.Background()

            tx, err := testDB.Begin(ctx)
            require.NoError(t, err, "failed to begin transaction")
            defer tx.Rollback(ctx)

            qtx := testQueries.WithTx(tx)

            params := tt.setup(ctx, qtx)
            err = qtx.UpdateEntity(ctx, params)
            tt.checks(t, err)
        })
    }
}

// Random data generators for tests

func createRandomEntity(ctx context.Context, qtx *Queries) Entity {
    entity, err := qtx.CreateEntity(ctx, CreateEntityParams{
        Field1: util.RandomString(10),
        Field2: util.RandomEmail(),
        Field3: true,
    })
    if err != nil {
        panic("failed to create random entity: " + err.Error())
    }
    return entity
}
```

## Checklist for New Tests

Before submitting database tests, verify:

- [ ] File named correctly: `<entity>_test.go`
- [ ] All imports included
- [ ] Table-driven test pattern used for all test functions
- [ ] Each function has at least 3 test cases (success, error, edge case)
- [ ] All test cases have descriptive names
- [ ] Setup functions create prerequisite data correctly
- [ ] Checks functions use `require` with descriptive messages
- [ ] Every test uses transaction with rollback
- [ ] Helper functions created for entity creation
- [ ] Helper functions use `panic()` on error
- [ ] No hardcoded IDs or emails (use `util.Random*()`)
- [ ] Optional fields tested with both nil and non-nil values
- [ ] Error cases verified with `require.Error()`
- [ ] Success cases verified with `require.NoError()`

## Common Patterns Reference

### Testing with Params Struct
```go
params CreateEntityParams
checks func(t *testing.T, entity Entity)
```
Use when you want to directly specify all parameters inline.

### Testing with Setup Function  
```go
setup  func(ctx context.Context, qtx *Queries) Params
checks func(t *testing.T, result, err error)
```
Use when you need to create dependent data or generate dynamic values.

### Testing :exec Functions (no return value)
```go
setup  func(ctx context.Context, qtx *Queries) UpdateParams
checks func(t *testing.T, err error)
```

### Testing :one Functions (single return)
```go
setup  func(ctx context.Context, qtx *Queries) ID
checks func(t *testing.T, entity Entity, err error)
```

### Testing :many Functions (slice return)
```go
setup  func(ctx context.Context, qtx *Queries) Params
checks func(t *testing.T, entities []Entity, err error)
```

## Tips and Best Practices

1. **Keep tests isolated**: Each test should be independent and not rely on other tests
2. **Use meaningful test data**: Avoid "test", "foo", "bar" - use realistic values
3. **Test edge cases**: Empty strings, nil values, UUID.Nil, empty slices
4. **Don't over-test**: Focus on function behavior, not database internals
5. **Use helper functions liberally**: They make tests more readable
6. **Rollback everything**: Never commit test transactions
7. **Descriptive error messages**: Future you will thank you
8. **Follow the pattern**: Consistency across test files helps maintainability

## Reference: SQL Function Annotations

When writing tests, understand what the SQL function does:

- `:one` - Returns a single row (test with exists/not exists)
- `:many` - Returns multiple rows (test with 0, 1, many results)
- `:exec` - No return value (test that it doesn't error)
- `:execrows` - Returns affected row count (test the count)
