CREATE FUNCTION public.protect_intake_form_provenance()
RETURNS TRIGGER
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = pg_catalog, public
AS $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM public.client_details AS client
        WHERE client.intake_form_id = OLD.id
    ) THEN
        RAISE EXCEPTION 'promoted intake source records are immutable';
    END IF;

    IF TG_OP = 'UPDATE'
       AND NEW.registration_form_id IS DISTINCT FROM OLD.registration_form_id THEN
        RAISE EXCEPTION 'intake registration ownership is immutable';
    END IF;

    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    END IF;
    RETURN NEW;
END;
$$;

CREATE FUNCTION public.protect_intake_topic_assessment_provenance()
RETURNS TRIGGER
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = pg_catalog, public
AS $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM public.client_details AS client
        WHERE client.intake_form_id = OLD.intake_form_id
    ) THEN
        RAISE EXCEPTION 'promoted intake source records are immutable';
    END IF;

    IF TG_OP = 'UPDATE'
       AND NEW.intake_form_id IS DISTINCT FROM OLD.intake_form_id THEN
        RAISE EXCEPTION 'intake assessment ownership is immutable';
    END IF;

    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    END IF;
    RETURN NEW;
END;
$$;

REVOKE ALL ON FUNCTION public.protect_intake_form_provenance() FROM PUBLIC;
REVOKE ALL ON FUNCTION public.protect_intake_topic_assessment_provenance() FROM PUBLIC;

DO $$
DECLARE
    policy_owner_name TEXT := 'maicare_rls_policy_owner_' || (
        SELECT oid::TEXT FROM pg_catalog.pg_database WHERE datname = current_database()
    );
BEGIN
    EXECUTE format('GRANT %I TO %I', policy_owner_name, current_user);
    EXECUTE format('ALTER FUNCTION public.protect_intake_form_provenance() OWNER TO %I', policy_owner_name);
    EXECUTE format('ALTER FUNCTION public.protect_intake_topic_assessment_provenance() OWNER TO %I', policy_owner_name);
    EXECUTE format('REVOKE %I FROM %I', policy_owner_name, current_user);
END;
$$;

DROP TRIGGER IF EXISTS intake_forms_protect_provenance ON public.intake_forms;
CREATE TRIGGER intake_forms_protect_provenance
BEFORE UPDATE OR DELETE ON public.intake_forms
FOR EACH ROW EXECUTE FUNCTION public.protect_intake_form_provenance();

DROP TRIGGER IF EXISTS intake_topic_assessments_protect_provenance ON public.intake_topic_assessments;
CREATE TRIGGER intake_topic_assessments_protect_provenance
BEFORE UPDATE OR DELETE ON public.intake_topic_assessments
FOR EACH ROW EXECUTE FUNCTION public.protect_intake_topic_assessment_provenance();
