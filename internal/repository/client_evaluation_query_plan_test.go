package repository

import (
	"context"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
)

func TestEvaluationListingQueryPlans(t *testing.T) {
	fixture := seedEvaluationLifecycleFixture(t)
	ctx := context.Background()
	tx, err := evaluationIntegrationPool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin plan transaction: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, `
		CREATE TEMP TABLE plan_clients (id uuid PRIMARY KEY) ON COMMIT DROP;
		INSERT INTO plan_clients
		SELECT gen_random_uuid() FROM generate_series(1, 5000);

		INSERT INTO client_details (
			id, first_name, last_name, email, gender, filenumber, street, house_number,
			postal_code, city
		)
		SELECT id, 'Plan', row_number() OVER ()::text, id::text || '@test.invalid', 'unknown',
			id::text, 'Test', '1', '1000AA', 'Test'
		FROM plan_clients;

		INSERT INTO client_goals (id, client_id, title, source, status)
		SELECT gen_random_uuid(), id, 'Plan goal', 'manual', 'active' FROM plan_clients;

		UPDATE client_details AS client SET
			status = 'in_care', care_start_date = CURRENT_DATE - 28,
			placed_in_care_at = CURRENT_TIMESTAMP - INTERVAL '28 days',
			evaluation_intervals_weeks = 4, last_evaluation_anchor_date = CURRENT_DATE - 28,
			next_evaluation_date = CURRENT_DATE + (ordinal.n % 60)
		FROM (
			SELECT id, row_number() OVER ()::int AS n FROM plan_clients
		) AS ordinal
		WHERE client.id = ordinal.id;

		INSERT INTO assigned_employee (client_id, employee_id, start_date, role)
		SELECT id, CASE WHEN row_number() OVER ()::int % 20 = 0 THEN $1::uuid ELSE $2::uuid END,
			CURRENT_DATE - 60, 'coordinator'
		FROM plan_clients;

		INSERT INTO client_goal_evaluations (
			id, client_id, evaluation_date, evaluation_interval_weeks, status,
			created_by_employee_id, updated_at
		)
		SELECT gen_random_uuid(), id, CURRENT_DATE + ((row_number() OVER ())::int % 60), 4,
			CASE WHEN (row_number() OVER ())::int % 2 = 0 THEN 'draft'::evaluation_status_enum
			ELSE 'completed'::evaluation_status_enum END,
			CASE WHEN (row_number() OVER ())::int % 20 < 2 THEN $1::uuid ELSE $2::uuid END,
			CURRENT_TIMESTAMP - ((row_number() OVER ())::int * INTERVAL '1 minute')
		FROM plan_clients;

		INSERT INTO client_goal_evaluations (
			id, client_id, evaluation_date, evaluation_interval_weeks, status,
			created_by_employee_id, updated_at
		)
		SELECT gen_random_uuid(), clients.id, CURRENT_DATE - series.n, 4,
			CASE WHEN series.n % 2 = 0 THEN 'draft'::evaluation_status_enum
			ELSE 'completed'::evaluation_status_enum END,
			CASE WHEN row_number() OVER ()::int % 20 < 2 THEN $1::uuid ELSE $2::uuid END,
			CURRENT_TIMESTAMP - (row_number() OVER ()::int * INTERVAL '1 minute')
		FROM plan_clients AS clients
		CROSS JOIN generate_series(1, 4) AS series(n);

		ANALYZE assigned_employee;
		ANALYZE client_details;
		ANALYZE client_goal_evaluations;
	`, pgx.QueryExecModeSimpleProtocol, fixture.owner.employeeID, fixture.other.employeeID); err != nil {
		t.Fatalf("seed plan fixture: %v", err)
	}
	if _, err := tx.Exec(ctx, `
		CREATE ROLE evaluation_plan_runtime NOLOGIN NOSUPERUSER NOBYPASSRLS;
		GRANT USAGE ON SCHEMA public TO evaluation_plan_runtime;
		GRANT SELECT ON assigned_employee, client_details, client_goal_evaluations,
			client_goal_evaluation_items TO evaluation_plan_runtime;
		GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO evaluation_plan_runtime;
		SET LOCAL ROLE evaluation_plan_runtime;
		SELECT set_config('myapp.current_user_id', $1, true),
			set_config('myapp.current_employee_id', $2, true);
	`, pgx.QueryExecModeSimpleProtocol, fixture.owner.userID, fixture.owner.employeeID); err != nil {
		t.Fatalf("initialize restricted plan actor: %v", err)
	}

	plans := map[string]struct {
		query     string
		wantIndex string
	}{
		"upcoming": {query: `WITH paged_clients AS MATERIALIZED (
			SELECT c.id AS client_id, c.first_name, c.last_name, c.next_evaluation_date
			FROM assigned_employee ae
			JOIN client_details c ON c.id = ae.client_id
			WHERE ae.employee_id = $1 AND ae.role = 'coordinator'
			AND c.status = 'in_care' AND c.next_evaluation_date IS NOT NULL
			ORDER BY c.next_evaluation_date, c.first_name, c.last_name, c.id LIMIT 10 OFFSET 0
		)
			SELECT c.client_id, c.first_name, c.last_name, c.next_evaluation_date,
			COALESCE(d.id IS NOT NULL, false), COALESCE(di.total_goals_count, 0)
			FROM paged_clients c
			LEFT JOIN LATERAL (
				SELECT e.id FROM client_goal_evaluations e
				WHERE e.client_id = c.client_id AND e.status = 'draft'
				ORDER BY e.updated_at DESC, e.id LIMIT 1
			) d ON true
			LEFT JOIN LATERAL (
				SELECT COUNT(*) AS total_goals_count
				FROM client_goal_evaluation_items i WHERE i.evaluation_id = d.id
			) di ON true
			ORDER BY c.next_evaluation_date, c.first_name, c.last_name, c.client_id LIMIT 10 OFFSET 0`,
			wantIndex: "client_goal_evaluations_draft_client_updated_idx"},
		"submitted": {query: `WITH paged_evaluations AS MATERIALIZED (
			SELECT e.id, e.client_id, c.first_name, c.last_name, e.evaluation_date, e.updated_at
			FROM client_goal_evaluations e JOIN client_details c ON c.id = e.client_id
			WHERE e.created_by_employee_id = $1 AND e.status = 'completed'
			ORDER BY e.updated_at DESC, e.id LIMIT 10 OFFSET 0
		)
			SELECT e.id, e.client_id, e.first_name, e.last_name,
			e.evaluation_date, e.updated_at, COALESCE(di.total_goals_count, 0)
			FROM paged_evaluations e
			LEFT JOIN LATERAL (
				SELECT COUNT(*) AS total_goals_count
				FROM client_goal_evaluation_items i WHERE i.evaluation_id = e.id
			) di ON true
			ORDER BY e.updated_at DESC, e.id LIMIT 10 OFFSET 0`,
			wantIndex: "client_goal_evaluations_employee_completed_updated_idx"},
		"drafts": {query: `WITH paged_evaluations AS MATERIALIZED (
			SELECT e.id, e.client_id, c.first_name, c.last_name, e.evaluation_date, e.updated_at
			FROM client_goal_evaluations e JOIN client_details c ON c.id = e.client_id
			WHERE e.created_by_employee_id = $1 AND e.status = 'draft'
			ORDER BY e.updated_at DESC, e.id LIMIT 10 OFFSET 0
		)
			SELECT e.id, e.client_id, e.first_name, e.last_name,
			e.evaluation_date, e.updated_at, COALESCE(di.total_goals_count, 0)
			FROM paged_evaluations e
			LEFT JOIN LATERAL (
				SELECT COUNT(*) AS total_goals_count
				FROM client_goal_evaluation_items i WHERE i.evaluation_id = e.id
			) di ON true
			ORDER BY e.updated_at DESC, e.id LIMIT 10 OFFSET 0`,
			wantIndex: "client_goal_evaluations_employee_draft_updated_idx"},
	}

	for name, plan := range plans {
		t.Run(name, func(t *testing.T) {
			rows, err := tx.Query(ctx, "EXPLAIN (ANALYZE, BUFFERS) "+plan.query, fixture.owner.employeeID)
			if err != nil {
				t.Fatalf("explain %s: %v", name, err)
			}
			defer rows.Close()
			var output strings.Builder
			for rows.Next() {
				var line string
				if err := rows.Scan(&line); err != nil {
					t.Fatalf("scan %s plan: %v", name, err)
				}
				output.WriteString(line)
				output.WriteByte('\n')
				t.Log(line)
			}
			if err := rows.Err(); err != nil {
				t.Fatalf("read %s plan: %v", name, err)
			}
			if !strings.Contains(output.String(), plan.wantIndex) {
				t.Fatalf("%s plan did not use %s", name, plan.wantIndex)
			}
		})
	}
}
