DROP TRIGGER IF EXISTS intake_forms_protect_provenance ON public.intake_forms;
CREATE TRIGGER intake_forms_protect_provenance
BEFORE UPDATE OR DELETE ON public.intake_forms
FOR EACH ROW EXECUTE FUNCTION public.protect_intake_provenance();

DROP TRIGGER IF EXISTS intake_topic_assessments_protect_provenance ON public.intake_topic_assessments;
CREATE TRIGGER intake_topic_assessments_protect_provenance
BEFORE UPDATE OR DELETE ON public.intake_topic_assessments
FOR EACH ROW EXECUTE FUNCTION public.protect_intake_provenance();

DO $$
DECLARE
    policy_owner_name TEXT := 'maicare_rls_policy_owner_' || (
        SELECT oid::TEXT FROM pg_catalog.pg_database WHERE datname = current_database()
    );
BEGIN
    EXECUTE format('GRANT %I TO %I', policy_owner_name, current_user);
    ALTER FUNCTION public.protect_intake_form_provenance() OWNER TO CURRENT_USER;
    ALTER FUNCTION public.protect_intake_topic_assessment_provenance() OWNER TO CURRENT_USER;
    EXECUTE format('REVOKE %I FROM %I', policy_owner_name, current_user);
END;
$$;

DROP FUNCTION public.protect_intake_form_provenance();
DROP FUNCTION public.protect_intake_topic_assessment_provenance();
