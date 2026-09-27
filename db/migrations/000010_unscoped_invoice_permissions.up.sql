-- Invoice permissions are unscoped. A grant applies to invoices for every client.
-- Keep the client ID check so callers cannot authorize records without an owner.
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
        AND (
            $2 LIKE 'INVOICE.%'
            OR CASE public.get_permission_scope($2)
                WHEN 'all'::public.permission_scope_enum THEN TRUE
                WHEN 'assigned'::public.permission_scope_enum THEN public.is_assigned_to_client($1)
                ELSE FALSE
            END
        ),
        FALSE
    );
$$;

-- Multi-client runs remain creator-owned, but no longer require a scope that
-- an unscoped INVOICE.CREATE grant cannot have.
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
          AND public.has_permission('INVOICE.CREATE')
    ), FALSE);
$$;

DROP POLICY IF EXISTS invoice_run_insert ON public.invoice_run;
CREATE POLICY invoice_run_insert ON public.invoice_run FOR INSERT
    WITH CHECK (
        created_by = public.get_current_employee_id()
        AND public.has_permission('INVOICE.CREATE')
    );

DROP POLICY IF EXISTS invoice_run_update ON public.invoice_run;
CREATE POLICY invoice_run_update ON public.invoice_run FOR UPDATE
    USING (
        created_by = public.get_current_employee_id()
        AND public.has_permission('INVOICE.CREATE')
    )
    WITH CHECK (created_by = public.get_current_employee_id());
