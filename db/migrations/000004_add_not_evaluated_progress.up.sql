ALTER TYPE public.client_goal_progress_enum ADD VALUE IF NOT EXISTS 'not_evaluated' BEFORE 'no_progress';

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
              AND (item.id IS NULL OR item.progress = 'not_evaluated')
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
