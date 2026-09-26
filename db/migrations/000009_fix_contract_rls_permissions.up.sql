-- Contract permissions are global (admin-only) and intentionally unscoped, so
-- can_access_client() can never succeed for them: it requires an 'all' or
-- 'assigned' scope. Gate the contract domain on has_permission() instead.

-- contract
DROP POLICY IF EXISTS contract_select ON public.contract;
CREATE POLICY contract_select ON public.contract
    FOR SELECT
    USING (public.has_permission('CONTRACT.VIEW'));
DROP POLICY IF EXISTS contract_insert ON public.contract;
CREATE POLICY contract_insert ON public.contract
    FOR INSERT
    WITH CHECK (public.has_permission('CONTRACT.CREATE'));
DROP POLICY IF EXISTS contract_update ON public.contract;
CREATE POLICY contract_update ON public.contract
    FOR UPDATE
    USING (public.has_permission('CONTRACT.UPDATE'))
    WITH CHECK (public.has_permission('CONTRACT.UPDATE'));
DROP POLICY IF EXISTS contract_delete ON public.contract;
CREATE POLICY contract_delete ON public.contract
    FOR DELETE
    USING (public.has_permission('CONTRACT.DELETE'));

-- contract_audit
DROP POLICY IF EXISTS contract_audit_select ON public.contract_audit;
CREATE POLICY contract_audit_select ON public.contract_audit
    FOR SELECT
    USING (public.has_permission('CONTRACT.VIEW'));

-- contract_reminder
DROP POLICY IF EXISTS contract_reminder_select ON public.contract_reminder;
CREATE POLICY contract_reminder_select ON public.contract_reminder
    FOR SELECT
    USING (public.has_permission('CONTRACT.VIEW'));
DROP POLICY IF EXISTS contract_reminder_insert ON public.contract_reminder;
CREATE POLICY contract_reminder_insert ON public.contract_reminder
    FOR INSERT
    WITH CHECK (public.has_permission('CONTRACT.UPDATE'));

-- contract_working_hours
DROP POLICY IF EXISTS contract_working_hours_select ON public.contract_working_hours;
CREATE POLICY contract_working_hours_select ON public.contract_working_hours
    FOR SELECT
    USING (public.has_permission('CONTRACT.VIEW'));
DROP POLICY IF EXISTS contract_working_hours_insert ON public.contract_working_hours;
CREATE POLICY contract_working_hours_insert ON public.contract_working_hours
    FOR INSERT
    WITH CHECK (public.has_permission('CONTRACT.UPDATE'));
DROP POLICY IF EXISTS contract_working_hours_update ON public.contract_working_hours;
CREATE POLICY contract_working_hours_update ON public.contract_working_hours
    FOR UPDATE
    USING (public.has_permission('CONTRACT.UPDATE'))
    WITH CHECK (public.has_permission('CONTRACT.UPDATE'));
DROP POLICY IF EXISTS contract_working_hours_delete ON public.contract_working_hours;
CREATE POLICY contract_working_hours_delete ON public.contract_working_hours
    FOR DELETE
    USING (public.has_permission('CONTRACT.UPDATE'));

-- client_agreement
DROP POLICY IF EXISTS client_agreement_select ON public.client_agreement;
CREATE POLICY client_agreement_select ON public.client_agreement
    FOR SELECT
    USING (public.has_permission('CONTRACT.VIEW'));
DROP POLICY IF EXISTS client_agreement_insert ON public.client_agreement;
CREATE POLICY client_agreement_insert ON public.client_agreement
    FOR INSERT
    WITH CHECK (public.has_permission('CONTRACT.CREATE'));
DROP POLICY IF EXISTS client_agreement_update ON public.client_agreement;
CREATE POLICY client_agreement_update ON public.client_agreement
    FOR UPDATE
    USING (public.has_permission('CONTRACT.UPDATE'))
    WITH CHECK (public.has_permission('CONTRACT.UPDATE'));
DROP POLICY IF EXISTS client_agreement_delete ON public.client_agreement;
CREATE POLICY client_agreement_delete ON public.client_agreement
    FOR DELETE
    USING (public.has_permission('CONTRACT.DELETE'));

-- provision
DROP POLICY IF EXISTS provision_select ON public.provision;
CREATE POLICY provision_select ON public.provision
    FOR SELECT
    USING (public.has_permission('CONTRACT.VIEW'));
DROP POLICY IF EXISTS provision_insert ON public.provision;
CREATE POLICY provision_insert ON public.provision
    FOR INSERT
    WITH CHECK (public.has_permission('CONTRACT.CREATE'));
DROP POLICY IF EXISTS provision_update ON public.provision;
CREATE POLICY provision_update ON public.provision
    FOR UPDATE
    USING (public.has_permission('CONTRACT.UPDATE'))
    WITH CHECK (public.has_permission('CONTRACT.UPDATE'));
DROP POLICY IF EXISTS provision_delete ON public.provision;
CREATE POLICY provision_delete ON public.provision
    FOR DELETE
    USING (public.has_permission('CONTRACT.DELETE'));

-- framework_agreement
DROP POLICY IF EXISTS framework_agreement_select ON public.framework_agreement;
CREATE POLICY framework_agreement_select ON public.framework_agreement
    FOR SELECT
    USING (public.has_permission('CONTRACT.VIEW'));
DROP POLICY IF EXISTS framework_agreement_insert ON public.framework_agreement;
CREATE POLICY framework_agreement_insert ON public.framework_agreement
    FOR INSERT
    WITH CHECK (public.has_permission('CONTRACT.CREATE'));
DROP POLICY IF EXISTS framework_agreement_update ON public.framework_agreement;
CREATE POLICY framework_agreement_update ON public.framework_agreement
    FOR UPDATE
    USING (public.has_permission('CONTRACT.UPDATE'))
    WITH CHECK (public.has_permission('CONTRACT.UPDATE'));
DROP POLICY IF EXISTS framework_agreement_delete ON public.framework_agreement;
CREATE POLICY framework_agreement_delete ON public.framework_agreement
    FOR DELETE
    USING (public.has_permission('CONTRACT.DELETE'));

-- Calendar attendees: CONTRACT.VIEW is meant to cover invoice-generation
-- appointment reads for a client, but the can_access_client() branch could
-- never succeed for the unscoped CONTRACT.VIEW permission.
DROP POLICY IF EXISTS calendar_event_attendees_select ON public.calendar_event_attendees;
CREATE POLICY calendar_event_attendees_select ON public.calendar_event_attendees
    FOR SELECT
    USING (
        public.can_participate_in_event(event_id)
        OR (
            client_id IS NOT NULL
            AND (
                public.can_access_client(client_id, 'CLIENT.VIEW')
                OR public.has_permission('CONTRACT.VIEW')
            )
        )
    );
DROP POLICY IF EXISTS calendar_event_attendees_insert ON public.calendar_event_attendees;
CREATE POLICY calendar_event_attendees_insert ON public.calendar_event_attendees
    FOR INSERT
    WITH CHECK (
        public.can_organize_event(event_id)
        AND (
            client_id IS NULL
            OR public.can_access_client(client_id, 'CLIENT.VIEW')
            OR public.has_permission('CONTRACT.VIEW')
        )
    );
