CREATE OR REPLACE FUNCTION public.get_current_user_id()
RETURNS UUID
LANGUAGE plpgsql
STABLE
PARALLEL SAFE
AS $$
BEGIN
    RETURN NULLIF(current_setting('myapp.current_user_id', true), '')::UUID;
EXCEPTION
    WHEN invalid_text_representation THEN RETURN NULL;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_current_employee_id()
RETURNS UUID
LANGUAGE plpgsql
STABLE
PARALLEL SAFE
AS $$
BEGIN
    RETURN NULLIF(current_setting('myapp.current_employee_id', true), '')::UUID;
EXCEPTION
    WHEN invalid_text_representation THEN RETURN NULL;
END;
$$;

CREATE OR REPLACE FUNCTION public.has_permission(permission_name TEXT)
RETURNS BOOLEAN
LANGUAGE sql
STABLE
SECURITY DEFINER
PARALLEL RESTRICTED
SET search_path = pg_catalog, public
AS $$
    SELECT COALESCE(EXISTS (
        SELECT 1
        FROM public.custom_user AS cu
        JOIN public.user_roles AS ur ON ur.user_id = cu.id
        JOIN public.role_permissions AS rp ON rp.role_id = ur.role_id
        JOIN public.permissions AS p ON p.id = rp.permission_id
        WHERE cu.id = public.get_current_user_id()
          AND cu.is_active
          AND p.name = $1
          AND (
              (p.is_scoped AND rp.scope IS NOT NULL)
              OR (NOT p.is_scoped AND rp.scope IS NULL)
          )
    ), FALSE);
$$;

CREATE OR REPLACE FUNCTION public.get_permission_scope(permission_name TEXT)
RETURNS public.permission_scope_enum
LANGUAGE sql
STABLE
SECURITY DEFINER
PARALLEL RESTRICTED
SET search_path = pg_catalog, public
AS $$
    SELECT rp.scope
    FROM public.custom_user AS cu
    JOIN public.user_roles AS ur ON ur.user_id = cu.id
    JOIN public.role_permissions AS rp ON rp.role_id = ur.role_id
    JOIN public.permissions AS p ON p.id = rp.permission_id
    WHERE cu.id = public.get_current_user_id()
      AND cu.is_active
      AND p.name = $1
      AND p.is_scoped
      AND rp.scope IS NOT NULL;
$$;

CREATE OR REPLACE FUNCTION public.is_assigned_to_client(client_id UUID)
RETURNS BOOLEAN
LANGUAGE sql
STABLE
SECURITY DEFINER
PARALLEL RESTRICTED
SET search_path = pg_catalog, public
AS $$
    SELECT COALESCE($1 IS NOT NULL AND EXISTS (
        SELECT 1
        FROM public.assigned_employee AS ae
        WHERE ae.client_id = $1
          AND ae.employee_id = public.get_current_employee_id()
          AND ae.start_date <= CURRENT_DATE
    ), FALSE);
$$;

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
        AND EXISTS (
            SELECT 1
            FROM public.custom_user AS cu
            JOIN public.employee_profile AS ep ON ep.user_id = cu.id
            WHERE cu.id = public.get_current_user_id()
              AND ep.id = public.get_current_employee_id()
              AND cu.is_active
              AND NOT ep.is_archived
              AND NOT COALESCE(ep.out_of_service, FALSE)
        )
        AND public.has_permission($2)
        AND CASE public.get_permission_scope($2)
            WHEN 'all'::public.permission_scope_enum THEN TRUE
            WHEN 'assigned'::public.permission_scope_enum THEN public.is_assigned_to_client($1)
            ELSE FALSE
        END,
        FALSE
    );
$$;

REVOKE ALL ON FUNCTION public.get_current_user_id() FROM PUBLIC;
REVOKE ALL ON FUNCTION public.get_current_employee_id() FROM PUBLIC;
REVOKE ALL ON FUNCTION public.has_permission(TEXT) FROM PUBLIC;
REVOKE ALL ON FUNCTION public.get_permission_scope(TEXT) FROM PUBLIC;
REVOKE ALL ON FUNCTION public.is_assigned_to_client(UUID) FROM PUBLIC;
REVOKE ALL ON FUNCTION public.can_access_client(UUID, TEXT) FROM PUBLIC;

GRANT EXECUTE ON FUNCTION public.get_current_user_id() TO CURRENT_USER;
GRANT EXECUTE ON FUNCTION public.get_current_employee_id() TO CURRENT_USER;
GRANT EXECUTE ON FUNCTION public.has_permission(TEXT) TO CURRENT_USER;
GRANT EXECUTE ON FUNCTION public.get_permission_scope(TEXT) TO CURRENT_USER;
GRANT EXECUTE ON FUNCTION public.is_assigned_to_client(UUID) TO CURRENT_USER;
GRANT EXECUTE ON FUNCTION public.can_access_client(UUID, TEXT) TO CURRENT_USER;
