CREATE INDEX CONCURRENTLY client_goal_evaluations_employee_draft_updated_idx
    ON public.client_goal_evaluations (created_by_employee_id, updated_at DESC, id)
    WHERE status = 'draft';
