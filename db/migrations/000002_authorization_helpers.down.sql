DROP FUNCTION IF EXISTS public.can_access_client(UUID, TEXT);
DROP FUNCTION IF EXISTS public.is_assigned_to_client(UUID);
DROP FUNCTION IF EXISTS public.get_permission_scope(TEXT);
DROP FUNCTION IF EXISTS public.has_permission(TEXT);
DROP FUNCTION IF EXISTS public.get_current_user_id();

CREATE OR REPLACE FUNCTION public.get_current_employee_id()
RETURNS UUID
LANGUAGE plpgsql
STABLE
AS $$
BEGIN
    RETURN current_setting('myapp.current_employee_id', true)::UUID;
EXCEPTION
    WHEN OTHERS THEN RETURN NULL;
END;
$$;

GRANT EXECUTE ON FUNCTION public.get_current_employee_id() TO PUBLIC;
