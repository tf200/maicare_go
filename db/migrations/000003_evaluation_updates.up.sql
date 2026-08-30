DO $$
DECLARE
    policy_owner_name TEXT := 'maicare_rls_policy_owner_' || (
        SELECT oid::TEXT FROM pg_catalog.pg_database WHERE datname = current_database()
    );
BEGIN
    EXECUTE format('GRANT %I TO %I', policy_owner_name, current_user);
    ALTER FUNCTION public.enforce_evaluation_submission_window() OWNER TO CURRENT_USER;
    EXECUTE format('REVOKE %I FROM %I', policy_owner_name, current_user);
END;
$$;

CREATE OR REPLACE FUNCTION evaluation_business_date(
    p_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
)
RETURNS DATE
LANGUAGE sql
STABLE
SECURITY DEFINER
SET search_path = pg_catalog, public
AS $$
    SELECT (p_at AT TIME ZONE COALESCE(
        (SELECT default_timezone FROM public.app_organization_profile WHERE singleton = TRUE),
        'Europe/Amsterdam'
    ))::date;
$$;

CREATE OR REPLACE FUNCTION enforce_evaluation_submission_window()
RETURNS TRIGGER AS $$
DECLARE
    v_next_eval_date DATE;
BEGIN
    IF NEW.status = 'completed' AND (OLD.status IS NULL OR OLD.status <> 'completed') THEN
        IF NOT EXISTS (
            SELECT 1
            FROM public.client_goals AS goal
            WHERE goal.client_id = NEW.client_id
              AND goal.status = 'active'
        ) OR EXISTS (
            SELECT 1
            FROM public.client_goals AS goal
            LEFT JOIN public.client_goal_evaluation_items AS item
              ON item.evaluation_id = NEW.id
             AND item.goal_id = goal.id
             AND item.client_id = NEW.client_id
            WHERE goal.client_id = NEW.client_id
              AND goal.status = 'active'
              AND (item.id IS NULL OR item.progress = 'no_progress')
        ) THEN
            RAISE EXCEPTION 'All active goals must be evaluated before submission';
        END IF;

        SELECT next_evaluation_date INTO v_next_eval_date
        FROM public.client_details
        WHERE id = NEW.client_id;

        IF public.evaluation_business_date() < (v_next_eval_date - 14) THEN
            RAISE EXCEPTION 'Evaluation cannot be completed more than 14 days before the due date (%)', v_next_eval_date;
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER SET search_path = pg_catalog, public;

REVOKE ALL ON FUNCTION public.enforce_evaluation_submission_window() FROM PUBLIC;

DO $$
DECLARE
    policy_owner_name TEXT := 'maicare_rls_policy_owner_' || (
        SELECT oid::TEXT FROM pg_catalog.pg_database WHERE datname = current_database()
    );
BEGIN
    EXECUTE format('GRANT %I TO %I', policy_owner_name, current_user);
    EXECUTE format('ALTER FUNCTION public.enforce_evaluation_submission_window() OWNER TO %I', policy_owner_name);
    EXECUTE format('REVOKE %I FROM %I', policy_owner_name, current_user);
END;
$$;

CREATE INDEX client_goal_evaluations_employee_completed_updated_idx
    ON public.client_goal_evaluations (created_by_employee_id, updated_at DESC, id)
    WHERE status = 'completed';

CREATE INDEX client_goal_evaluations_employee_draft_updated_idx
    ON public.client_goal_evaluations (created_by_employee_id, updated_at DESC, id)
    WHERE status = 'draft';

CREATE INDEX client_goal_evaluations_draft_client_updated_idx
    ON public.client_goal_evaluations (client_id, updated_at DESC, id)
    WHERE status = 'draft';
