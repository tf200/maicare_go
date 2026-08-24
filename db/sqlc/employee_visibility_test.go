package db

import (
	"strings"
	"testing"
)

func TestHumanFacingEmployeeQueriesExcludeSystemEmployee(t *testing.T) {
	queries := map[string]string{
		"employee list":           listEmployeeProfile,
		"employee count":          countEmployeeProfile,
		"employee search":         searchEmployeesByNameOrEmail,
		"handbook eligible":       listEligibleEmployeesForHandbookAssignment,
		"handbook eligible count": countEligibleEmployeesForHandbookAssignment,
		"dashboard totals":        getAdminDashboardStatCards,
		"employee totals":         getEmployeeCounts,
	}

	for name, query := range queries {
		if !strings.Contains(query, "id <> $") && !strings.Contains(query, "ep.id <> $") && !strings.Contains(query, "employee_profile.id <> $") {
			t.Errorf("%s query does not exclude the configured system employee", name)
		}
	}
}

func TestCountEmployeeProfileKeepsSearchFilter(t *testing.T) {
	if !strings.Contains(countEmployeeProfile, "first_name ILIKE") || !strings.Contains(countEmployeeProfile, "last_name ILIKE") {
		t.Error("employee count query does not apply the employee search filter")
	}
}
