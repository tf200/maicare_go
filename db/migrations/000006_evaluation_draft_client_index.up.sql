CREATE INDEX CONCURRENTLY client_goal_evaluations_draft_client_updated_idx
    ON public.client_goal_evaluations (client_id, updated_at DESC, id)
    WHERE status = 'draft';
