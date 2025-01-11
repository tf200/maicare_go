package db

import (
	"context"
	"sync"
	"testing"
	"time"

	"maicare_go/util"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomEmployee(t *testing.T) (EmployeeProfile, *CustomUser) {
	// Create prerequisite records
	location := CreateRandomLocation(t)
	user := CreateRandomUser(t)

	arg := CreateEmployeeProfileParams{
		UserID:                    user.ID,
		FirstName:                 util.RandomString(5),
		LastName:                  util.RandomString(5),
		Position:                  util.StringPtr(util.RandomString(5)),
		Department:                util.StringPtr(util.RandomString(5)),
		EmployeeNumber:            util.StringPtr(util.RandomString(5)),
		EmploymentNumber:          util.StringPtr(util.RandomString(5)),
		PrivateEmailAddress:       util.StringPtr(util.RandomString(5)),
		Email:                     util.RandomEmail(),
		AuthenticationPhoneNumber: util.StringPtr(util.RandomString(5)),
		PrivatePhoneNumber:        util.StringPtr(util.RandomString(5)),
		WorkPhoneNumber:           util.StringPtr(util.RandomString(5)),
		DateOfBirth:               pgtype.Date{Time: time.Now(), Valid: true},
		HomeTelephoneNumber:       util.StringPtr(util.RandomString(5)),
		IsSubcontractor:           util.BoolPtr(true),
		Gender:                    util.StringPtr(util.RandomString(5)),
		LocationID:                util.IntPtr(location.ID),
		HasBorrowed:               false,
		OutOfService:              util.BoolPtr(util.RandomBool()),
		IsArchived:                util.RandomBool(),
	}

	employee, err := testQueries.CreateEmployeeProfile(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, employee)

	// Verify all fields match
	require.Equal(t, arg.UserID, employee.UserID)
	require.Equal(t, arg.FirstName, employee.FirstName)
	require.Equal(t, arg.LastName, employee.LastName)
	require.Equal(t, arg.Position, employee.Position)
	require.Equal(t, arg.Department, employee.Department)
	require.Equal(t, arg.EmployeeNumber, employee.EmployeeNumber)
	require.Equal(t, arg.EmploymentNumber, employee.EmploymentNumber)
	require.Equal(t, arg.PrivateEmailAddress, employee.PrivateEmailAddress)
	require.Equal(t, arg.Email, employee.Email)
	require.Equal(t, arg.AuthenticationPhoneNumber, employee.AuthenticationPhoneNumber)
	require.Equal(t, arg.PrivatePhoneNumber, employee.PrivatePhoneNumber)
	require.Equal(t, arg.WorkPhoneNumber, employee.WorkPhoneNumber)
	require.Equal(t, arg.DateOfBirth.Time.Format("2006-01-02"), employee.DateOfBirth.Time.Format("2006-01-02"))
	require.Equal(t, arg.HomeTelephoneNumber, employee.HomeTelephoneNumber)
	require.Equal(t, arg.IsSubcontractor, employee.IsSubcontractor)
	require.Equal(t, arg.Gender, employee.Gender)
	require.Equal(t, arg.LocationID, employee.LocationID)
	require.Equal(t, arg.HasBorrowed, employee.HasBorrowed)
	require.Equal(t, arg.OutOfService, employee.OutOfService)
	require.Equal(t, arg.IsArchived, employee.IsArchived)

	require.NotZero(t, employee.ID)
	require.NotZero(t, employee.CreatedAt)

	require.Equal(t, util.IntPtr(location.ID), employee.LocationID)

	return employee, user
}

func TestCreateEmployeeProfile(t *testing.T) {
	createRandomEmployee(t)
}

func TestListEmployeeProfile(t *testing.T) {
	// Create multiple employees with random statuses
	numEmployees := 20
	var wg sync.WaitGroup
	errCh := make(chan error, numEmployees)
	employeeCh := make(chan EmployeeProfile, numEmployees)

	// Create employees concurrently
	for i := 0; i < numEmployees; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Randomly set archived and out of service
			emp, _ := createRandomEmployee(t)
			employeeCh <- emp
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errCh)
	close(employeeCh)

	// Check for any errors during creation
	for err := range errCh {
		require.NoError(t, err, "error creating employee")
	}

	testCases := []struct {
		name         string
		params       ListEmployeeProfileParams
		expectedLen  int
		checkResults func(t *testing.T, results []ListEmployeeProfileRow)
	}{
		{
			name: "List all employees with limit 5",
			params: ListEmployeeProfileParams{
				IncludeArchived:     util.BoolPtr(true),
				IncludeOutOfService: util.BoolPtr(true),
				Limit:               5,
				Offset:              0,
			},
			expectedLen: 5,
			checkResults: func(t *testing.T, results []ListEmployeeProfileRow) {
				require.NotEmpty(t, results)
			},
		},
		{
			name: "List with offset",
			params: ListEmployeeProfileParams{
				IncludeArchived:     util.BoolPtr(true),
				IncludeOutOfService: util.BoolPtr(true),
				Limit:               5,
				Offset:              5,
			},
			expectedLen: 5,
			checkResults: func(t *testing.T, results []ListEmployeeProfileRow) {
				require.NotEmpty(t, results)
			},
		},
		{
			name: "Exclude archived only",
			params: ListEmployeeProfileParams{
				IncludeArchived:     util.BoolPtr(false),
				IncludeOutOfService: util.BoolPtr(true),
				Limit:               10,
				Offset:              0,
			},
			expectedLen: 10,
			checkResults: func(t *testing.T, results []ListEmployeeProfileRow) {
				for _, emp := range results {
					require.False(t, emp.IsArchived, "should not include archived employees")
				}
			},
		},
		{
			name: "Exclude out of service only",
			params: ListEmployeeProfileParams{
				IncludeArchived:     util.BoolPtr(true),
				IncludeOutOfService: util.BoolPtr(false),
				Limit:               10,
				Offset:              0,
			},
			expectedLen: 10,
			checkResults: func(t *testing.T, results []ListEmployeeProfileRow) {
				for _, emp := range results {
					require.False(t, *emp.OutOfService, "should not include out of service employees")
				}
			},
		},
		{
			name: "Exclude both archived and out of service",
			params: ListEmployeeProfileParams{
				IncludeArchived:     util.BoolPtr(false),
				IncludeOutOfService: util.BoolPtr(false),
				Limit:               10,
				Offset:              0,
			},
			expectedLen: 10,
			checkResults: func(t *testing.T, results []ListEmployeeProfileRow) {
				for _, emp := range results {
					require.False(t, emp.IsArchived, "should not include archived employees")
					require.False(t, *emp.OutOfService, "should not include out of service employees")
				}
			},
		},
		{
			name: "Check ordering",
			params: ListEmployeeProfileParams{
				Limit:  10,
				Offset: 0,
			},
			expectedLen: 10,
			checkResults: func(t *testing.T, results []ListEmployeeProfileRow) {
				for i := 1; i < len(results); i++ {
					require.True(t, results[i-1].CreatedAt.Time.After(results[i].CreatedAt.Time) ||
						results[i-1].CreatedAt.Time.Equal(results[i].CreatedAt.Time),
						"results should be ordered by created DESC")
				}
			},
		},
		{
			name: "Check offset beyond total",
			params: ListEmployeeProfileParams{
				Limit:  10,
				Offset: 1000, // very large offset
			},
			expectedLen: 0,
			checkResults: func(t *testing.T, results []ListEmployeeProfileRow) {
				require.Empty(t, results)
			},
		},
	}

	// Run all test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			results, err := testQueries.ListEmployeeProfile(context.Background(), tc.params)
			require.NoError(t, err)

			// Check length matches expected
			require.Len(t, results, tc.expectedLen)

			// Run test-specific checks
			tc.checkResults(t, results)

			// Common validations
			for _, emp := range results {
				require.NotEmpty(t, emp.ID)
				require.NotEmpty(t, emp.UserID)
				require.NotEmpty(t, emp.FirstName)
				require.NotEmpty(t, emp.LastName)
				require.NotZero(t, emp.CreatedAt)
			}
		})
	}
}

func TestCountEmployeeProfile(t *testing.T) {
	// Get initial count before adding test data
	initialCount, err := testQueries.CountEmployeeProfile(context.Background(), CountEmployeeProfileParams{
		IncludeArchived:     util.BoolPtr(true),
		IncludeOutOfService: util.BoolPtr(true),
	})
	require.NoError(t, err)

	// Create test data
	numEmployees := 20
	var wg sync.WaitGroup
	for i := 0; i < numEmployees; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			createRandomEmployee(t)
		}()
	}
	wg.Wait()

	testCases := []struct {
		name       string
		params     CountEmployeeProfileParams
		checkCount func(t *testing.T, count int64)
	}{
		{
			name: "Count all employees",
			params: CountEmployeeProfileParams{
				IncludeArchived:     util.BoolPtr(true),
				IncludeOutOfService: util.BoolPtr(true),
			},
			checkCount: func(t *testing.T, count int64) {
				require.Equal(t, initialCount+int64(numEmployees), count,
					"should match initial count plus newly created employees")
			},
		},
		{
			name: "Count non-archived employees",
			params: CountEmployeeProfileParams{
				IncludeArchived:     util.BoolPtr(false),
				IncludeOutOfService: util.BoolPtr(true),
			},
			checkCount: func(t *testing.T, count int64) {
				require.Less(t, count, initialCount+int64(numEmployees))
			},
		},
		{
			name: "Count non-out-of-service employees",
			params: CountEmployeeProfileParams{
				IncludeArchived:     util.BoolPtr(true),
				IncludeOutOfService: util.BoolPtr(false),
			},
			checkCount: func(t *testing.T, count int64) {
				require.Less(t, count, initialCount+int64(numEmployees))
			},
		},
		{
			name: "Filter by department",
			params: CountEmployeeProfileParams{
				Department: util.StringPtr("IT"),
			},
			checkCount: func(t *testing.T, count int64) {
				require.GreaterOrEqual(t, count, int64(0))
				require.LessOrEqual(t, count, initialCount+int64(numEmployees))
			},
		},
		{
			name: "Filter by non-existent department",
			params: CountEmployeeProfileParams{
				Department: util.StringPtr("NonExistentDepartment"),
			},
			checkCount: func(t *testing.T, count int64) {
				require.Equal(t, int64(0), count)
			},
		},
		{
			name: "Count with all filters",
			params: CountEmployeeProfileParams{
				IncludeArchived:     util.BoolPtr(false),
				IncludeOutOfService: util.BoolPtr(false),
				Department:          util.StringPtr("IT"),
				Position:            util.StringPtr("Developer"),
			},
			checkCount: func(t *testing.T, count int64) {
				require.GreaterOrEqual(t, count, int64(0))
				require.LessOrEqual(t, count, initialCount+int64(numEmployees))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			count, err := testQueries.CountEmployeeProfile(context.Background(), tc.params)
			require.NoError(t, err)
			tc.checkCount(t, count)
		})
	}
}

func TestGetEmployeeProfileByUserID(t *testing.T) {
	employee, user := createRandomEmployee(t)
	employee2, err := testQueries.GetEmployeeProfileByUserID(context.Background(), employee.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, employee2)
	require.NotEmpty(t, employee2.Permissions)
	require.Equal(t, employee.ID, employee2.EmployeeID)
	require.Equal(t, employee.UserID, employee2.UserID)
	require.Equal(t, employee.FirstName, employee2.FirstName)
	require.Equal(t, employee.LastName, employee2.LastName)
	require.Equal(t, user.RoleID, employee2.RoleID)
}

func TestGetEmployeeProfileByID(t *testing.T) {
	employee, user := createRandomEmployee(t)
	employee2, err := testQueries.GetEmployeeProfileByID(context.Background(), employee.ID)
	require.NoError(t, err)
	require.NotEmpty(t, employee2)
	require.Equal(t, employee.ID, employee2.ID)
	require.Equal(t, employee.UserID, employee2.UserID)
	require.Equal(t, employee.FirstName, employee2.FirstName)
	require.Equal(t, employee.LastName, employee2.LastName)
	require.Equal(t, user.RoleID, employee2.RoleID)
	require.Equal(t, employee.Email, employee2.Email)
}

func TestUpdateEmployeeProfile(t *testing.T) {
	employee, _ := createRandomEmployee(t)

	arg := UpdateEmployeeProfileParams{
		ID:          employee.ID,
		FirstName:   util.StringPtr(util.RandomString(5)),
		DateOfBirth: pgtype.Date{Time: time.Date(1958, 1, 15, 0, 0, 0, 0, time.UTC), Valid: true},
	}
	updatedEmployee, err := testQueries.UpdateEmployeeProfile(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedEmployee)
	require.NotEqual(t, employee.FirstName, updatedEmployee.FirstName)
	require.NotEqual(t, employee.DateOfBirth, updatedEmployee.DateOfBirth)
	require.Equal(t, employee.ID, updatedEmployee.ID)
	require.Equal(t, employee.UserID, updatedEmployee.UserID)
	require.Equal(t, employee.LastName, updatedEmployee.LastName)
}

func addRandomEducation(t *testing.T, employeeID int64) EmployeeEducation {

	educationArg := AddEducationToEmployeeProfileParams{
		EmployeeID:      employeeID,
		InstitutionName: util.RandomString(5),
		Degree:          util.RandomString(5),
		FieldOfStudy:    util.RandomString(5),
		StartDate:       pgtype.Date{Time: time.Now(), Valid: true},
		EndDate:         pgtype.Date{Time: time.Now(), Valid: true},
	}

	education, err := testQueries.AddEducationToEmployeeProfile(context.Background(), educationArg)
	require.NoError(t, err)
	require.NotEmpty(t, education)
	require.Equal(t, educationArg.EmployeeID, education.EmployeeID)
	require.Equal(t, educationArg.InstitutionName, education.InstitutionName)
	require.Equal(t, educationArg.Degree, education.Degree)
	return education
}

func TestAddEducationToEmployeeProfile(t *testing.T) {
	employee, _ := createRandomEmployee(t)
	addRandomEducation(t, employee.ID)
}

func TestListEducations(t *testing.T) {
	employee, _ := createRandomEmployee(t)
	numEducations := 5
	var wg sync.WaitGroup
	for i := 0; i < numEducations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			addRandomEducation(t, employee.ID)
		}()
	}
	wg.Wait()

	educations, err := testQueries.ListEducations(context.Background(), employee.ID)
	require.NoError(t, err)
	require.Len(t, educations, numEducations)

}

func TestUpdateEducation(t *testing.T) {
	employee, _ := createRandomEmployee(t)
	education := addRandomEducation(t, employee.ID)

	arg := UpdateEmployeeEducationParams{
		ID:              education.ID,
		InstitutionName: util.StringPtr(util.RandomString(5)),
		Degree:          util.StringPtr(util.RandomString(5)),
		FieldOfStudy:    util.StringPtr(util.RandomString(5)),
	}

	updatedEducation, err := testQueries.UpdateEmployeeEducation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedEducation)
	require.NotEqual(t, education.InstitutionName, updatedEducation.InstitutionName)
	require.NotEqual(t, education.Degree, updatedEducation.Degree)
	require.NotEqual(t, education.FieldOfStudy, updatedEducation.FieldOfStudy)
}

func addRandomExperience(t *testing.T, employeeID int64) EmployeeExperience {

	experienceArg := AddEmployeeExperienceParams{
		EmployeeID:  employeeID,
		JobTitle:    util.RandomString(5),
		CompanyName: util.RandomString(5),
		StartDate:   pgtype.Date{Time: time.Now(), Valid: true},
		EndDate:     pgtype.Date{Time: time.Now(), Valid: true},
		Description: util.StringPtr(util.RandomString(5)),
	}

	experience, err := testQueries.AddEmployeeExperience(context.Background(), experienceArg)
	require.NoError(t, err)
	require.NotEmpty(t, experience)
	require.Equal(t, experienceArg.EmployeeID, experience.EmployeeID)
	require.Equal(t, experienceArg.JobTitle, experience.JobTitle)
	require.Equal(t, experienceArg.CompanyName, experience.CompanyName)
	return experience
}

func TestAddEmployeeExperience(t *testing.T) {
	employee, _ := createRandomEmployee(t)
	addRandomExperience(t, employee.ID)
}

func TestListEmployeeExperience(t *testing.T) {
	employee, _ := createRandomEmployee(t)
	numExperiences := 5
	var wg sync.WaitGroup
	for i := 0; i < numExperiences; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			addRandomExperience(t, employee.ID)
		}()
	}
	wg.Wait()

	experiences, err := testQueries.ListEmployeeExperience(context.Background(), employee.ID)
	require.NoError(t, err)
	require.Len(t, experiences, numExperiences)
}

func TestUpdateEmployeeExperience(t *testing.T) {
	employee, _ := createRandomEmployee(t)
	experience := addRandomExperience(t, employee.ID)

	arg := UpdateEmployeeExperienceParams{
		ID:          experience.ID,
		JobTitle:    util.StringPtr(util.RandomString(5)),
		CompanyName: util.StringPtr(util.RandomString(5)),
	}

	updatedExperience, err := testQueries.UpdateEmployeeExperience(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedExperience)
	require.NotEqual(t, experience.JobTitle, updatedExperience.JobTitle)
	require.NotEqual(t, experience.CompanyName, updatedExperience.CompanyName)

}
