-- Rename CLIENT.INCIDENT.* permissions to top-level INCIDENT.*.
-- Keeps is_scoped=TRUE so can_access_client() scope checks keep working.
-- Derived metadata matches internal/domain/permission.go:
--   group_key='incident', section_key='general', display names unchanged.

UPDATE public.permissions
SET name = 'INCIDENT.CREATE',
    group_key = 'incident',
    section_key = 'general',
    display_name = 'Create Incident',
    is_scoped = TRUE
WHERE name = 'CLIENT.INCIDENT.CREATE';

UPDATE public.permissions
SET name = 'INCIDENT.VIEW',
    group_key = 'incident',
    section_key = 'general',
    display_name = 'View Incident',
    is_scoped = TRUE
WHERE name = 'CLIENT.INCIDENT.VIEW';

UPDATE public.permissions
SET name = 'INCIDENT.UPDATE',
    group_key = 'incident',
    section_key = 'general',
    display_name = 'Update Incident',
    is_scoped = TRUE
WHERE name = 'CLIENT.INCIDENT.UPDATE';

UPDATE public.permissions
SET name = 'INCIDENT.DELETE',
    group_key = 'incident',
    section_key = 'general',
    display_name = 'Delete Incident',
    is_scoped = TRUE
WHERE name = 'CLIENT.INCIDENT.DELETE';

UPDATE public.permissions
SET name = 'INCIDENT.CONFIRM',
    group_key = 'incident',
    section_key = 'general',
    display_name = 'Confirm Incident',
    is_scoped = TRUE
WHERE name = 'CLIENT.INCIDENT.CONFIRM';

-- begin_incident_creation: check INCIDENT.CREATE
CREATE OR REPLACE FUNCTION public.begin_incident_creation(client_id UUID)
RETURNS UUID
LANGUAGE plpgsql
VOLATILE
SECURITY DEFINER
PARALLEL UNSAFE
SET search_path = pg_catalog, public
AS $$
DECLARE
    new_incident_id UUID;
BEGIN
    IF NOT public.can_access_client($1, 'INCIDENT.CREATE') THEN
        RETURN NULL;
    END IF;

    DELETE FROM public.rls_report_creation_context
    WHERE backend_pid = pg_backend_pid()
      AND report_kind = 'incident';

    new_incident_id := gen_random_uuid();
    INSERT INTO public.rls_report_creation_context (
        backend_pid, transaction_id, report_kind, report_id, user_id, employee_id
    ) VALUES (
        pg_backend_pid(), pg_current_xact_id(), 'incident', new_incident_id,
        public.get_current_user_id(), public.get_current_employee_id()
    );

    RETURN new_incident_id;
END;
$$;

-- get_authorized_incident_recipient_emails: check INCIDENT.CONFIRM
CREATE OR REPLACE FUNCTION public.get_authorized_incident_recipient_emails(client_id UUID)
RETURNS TABLE (email TEXT)
LANGUAGE sql
STABLE
SECURITY DEFINER
PARALLEL RESTRICTED
SET search_path = pg_catalog, public
AS $$
    SELECT cec.email::TEXT
    FROM public.client_emergency_contact AS cec
    WHERE cec.client_id = $1
      AND cec.incidents_reports
      AND cec.is_verified
      AND cec.email IS NOT NULL
      AND public.can_access_client($1, 'INCIDENT.CONFIRM')
    ORDER BY cec.created_at ASC;
$$;

-- confirm_incident: check CLIENT.VIEW + INCIDENT.VIEW + INCIDENT.CONFIRM
CREATE OR REPLACE FUNCTION public.confirm_incident(incident_id UUID)
RETURNS BIGINT
LANGUAGE plpgsql
VOLATILE
SECURITY DEFINER
PARALLEL UNSAFE
SET search_path = pg_catalog, public
AS $$
DECLARE
    owning_client_id UUID;
    affected BIGINT;
BEGIN
    SELECT client_id INTO owning_client_id
    FROM public.incident
    WHERE id = $1;

    IF owning_client_id IS NULL
       OR NOT public.can_access_client(owning_client_id, 'CLIENT.VIEW')
       OR NOT public.can_access_client(owning_client_id, 'INCIDENT.VIEW')
       OR NOT public.can_access_client(owning_client_id, 'INCIDENT.CONFIRM') THEN
        RETURN 0;
    END IF;

    UPDATE public.incident
    SET is_confirmed = TRUE,
        confirmed_at = NOW(),
        confirmed_by = public.get_current_user_id()
    WHERE id = $1
      AND is_confirmed = FALSE;
    GET DIAGNOSTICS affected = ROW_COUNT;
    RETURN affected;
END;
$$;

-- mark_incident_confirmation_email_sent
CREATE OR REPLACE FUNCTION public.mark_incident_confirmation_email_sent(incident_id UUID, claim_token UUID)
RETURNS BIGINT
LANGUAGE plpgsql
VOLATILE
SECURITY DEFINER
PARALLEL UNSAFE
SET search_path = pg_catalog, public
AS $$
DECLARE
    owning_client_id UUID;
    affected BIGINT;
BEGIN
    SELECT client_id INTO owning_client_id
    FROM public.incident
    WHERE id = $1;

    IF owning_client_id IS NULL
       OR NOT public.can_access_client(owning_client_id, 'CLIENT.VIEW')
       OR NOT public.can_access_client(owning_client_id, 'INCIDENT.VIEW')
       OR NOT public.can_access_client(owning_client_id, 'INCIDENT.CONFIRM') THEN
        RETURN 0;
    END IF;

    UPDATE public.incident
    SET confirmation_email_sent_at = NOW(),
        confirmation_email_claimed_at = NULL,
        confirmation_email_claim_token = NULL
    WHERE id = $1
      AND is_confirmed
      AND confirmation_email_claim_token = $2
      AND confirmation_email_claimed_at IS NOT NULL
      AND confirmation_email_sent_at IS NULL;
    GET DIAGNOSTICS affected = ROW_COUNT;
    RETURN affected;
END;
$$;

-- claim_incident_confirmation_email
CREATE OR REPLACE FUNCTION public.claim_incident_confirmation_email(incident_id UUID)
RETURNS UUID
LANGUAGE plpgsql
VOLATILE
SECURITY DEFINER
PARALLEL UNSAFE
SET search_path = pg_catalog, public
AS $$
DECLARE
    owning_client_id UUID;
    new_claim_token UUID := gen_random_uuid();
BEGIN
    SELECT client_id INTO owning_client_id FROM public.incident WHERE id = $1;
    IF owning_client_id IS NULL
       OR NOT public.can_access_client(owning_client_id, 'CLIENT.VIEW')
       OR NOT public.can_access_client(owning_client_id, 'INCIDENT.VIEW')
       OR NOT public.can_access_client(owning_client_id, 'INCIDENT.CONFIRM') THEN
        RETURN NULL;
    END IF;

    UPDATE public.incident
    SET confirmation_email_claimed_at = NOW(),
        confirmation_email_claim_token = new_claim_token
    WHERE id = $1
      AND is_confirmed
      AND confirmation_email_sent_at IS NULL
      AND (
          confirmation_email_claimed_at IS NULL
          OR confirmation_email_claimed_at < NOW() - INTERVAL '15 minutes'
      );
    IF FOUND THEN
        RETURN new_claim_token;
    END IF;
    RETURN NULL;
END;
$$;

-- release_incident_confirmation_email
CREATE OR REPLACE FUNCTION public.release_incident_confirmation_email(incident_id UUID, claim_token UUID)
RETURNS BIGINT
LANGUAGE plpgsql
VOLATILE
SECURITY DEFINER
PARALLEL UNSAFE
SET search_path = pg_catalog, public
AS $$
DECLARE
    owning_client_id UUID;
    affected BIGINT;
BEGIN
    SELECT client_id INTO owning_client_id FROM public.incident WHERE id = $1;
    IF owning_client_id IS NULL
       OR NOT public.can_access_client(owning_client_id, 'CLIENT.VIEW')
       OR NOT public.can_access_client(owning_client_id, 'INCIDENT.VIEW')
       OR NOT public.can_access_client(owning_client_id, 'INCIDENT.CONFIRM') THEN
        RETURN 0;
    END IF;

    UPDATE public.incident
    SET confirmation_email_claimed_at = NULL,
        confirmation_email_claim_token = NULL
    WHERE id = $1
      AND confirmation_email_claim_token = $2
      AND confirmation_email_sent_at IS NULL
      AND confirmation_email_claimed_at IS NOT NULL;
    GET DIAGNOSTICS affected = ROW_COUNT;
    RETURN affected;
END;
$$;

-- RLS policies on incident table
DROP POLICY IF EXISTS incident_select ON public.incident;
CREATE POLICY incident_select ON public.incident
    FOR SELECT
    USING (
        public.can_access_client(client_id, 'INCIDENT.VIEW')
        OR public.can_read_created_incident(id)
    );

DROP POLICY IF EXISTS incident_insert ON public.incident;
CREATE POLICY incident_insert ON public.incident
    FOR INSERT
    WITH CHECK (public.can_access_client(client_id, 'INCIDENT.CREATE'));

DROP POLICY IF EXISTS incident_update ON public.incident;
CREATE POLICY incident_update ON public.incident
    FOR UPDATE
    USING (public.can_access_client(client_id, 'INCIDENT.UPDATE'))
    WITH CHECK (public.can_access_client(client_id, 'INCIDENT.UPDATE'));

DROP POLICY IF EXISTS incident_delete ON public.incident;
CREATE POLICY incident_delete ON public.incident
    FOR DELETE
    USING (public.can_access_client(client_id, 'INCIDENT.DELETE'));
