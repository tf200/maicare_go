CREATE OR REPLACE FUNCTION public.can_access_client(client_id UUID, permission_name TEXT)
RETURNS BOOLEAN
LANGUAGE sql
STABLE
SECURITY DEFINER
PARALLEL RESTRICTED
SET search_path = pg_catalog, public
AS $$
    SELECT COALESCE(
        $1 IS NOT NULL
        AND $2 IS NOT NULL
        AND public.has_permission($2)
        AND CASE public.get_permission_scope($2)
            WHEN 'all'::public.permission_scope_enum THEN TRUE
            WHEN 'assigned'::public.permission_scope_enum THEN public.is_assigned_to_client($1)
            ELSE FALSE
        END,
        FALSE
    );
$$;

CREATE OR REPLACE FUNCTION public.can_manage_invoice_run(invoice_run_id UUID)
RETURNS BOOLEAN
LANGUAGE sql
STABLE
SECURITY DEFINER
SET search_path = pg_catalog, public
AS $$
    SELECT COALESCE(EXISTS (
        SELECT 1
        FROM public.invoice_run AS managed_run
        WHERE managed_run.id = $1
          AND managed_run.created_by = public.get_current_employee_id()
          AND public.get_permission_scope('INVOICE.CREATE') = 'all'
    ), FALSE);
$$;

DROP POLICY IF EXISTS invoice_run_insert ON public.invoice_run;
CREATE POLICY invoice_run_insert ON public.invoice_run FOR INSERT
    WITH CHECK (
        created_by = public.get_current_employee_id()
        AND public.get_permission_scope('INVOICE.CREATE') = 'all'
    );

DROP POLICY IF EXISTS invoice_run_update ON public.invoice_run;
CREATE POLICY invoice_run_update ON public.invoice_run FOR UPDATE
    USING (
        created_by = public.get_current_employee_id()
        AND public.get_permission_scope('INVOICE.CREATE') = 'all'
    )
    WITH CHECK (created_by = public.get_current_employee_id());
