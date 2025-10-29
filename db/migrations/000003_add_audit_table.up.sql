CREATE TABLE audit (
    event_id UUID PRIMARY KEY,
    event_type TEXT NOT NULL,
    occured_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    actor_role TEXT[] NOT NULL,
    actor_id UUID NOT NULL,
    subject_type TEXT NOT NULL,
    subject_id UUID NOT NULL,
    access_reason TEXT NOT NULL,
    action TEXT NOT NULL,
    result TEXT NOT NULL,
    module TEXT NOT NULL,
    tenant_id TEXT NOT NULL,
    details JSONB,
    ip INET,
    user_agent TEXT,
    hash_prev TEXT NOT NULL,
    hash_self TEXT NOT NULL
);