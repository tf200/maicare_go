CREATE TABLE "group" (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);


CREATE TABLE location (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    address VARCHAR(100) NOT NULL,
    capacity INTEGER NULL
);

-- Table: Roles
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,       
    "name" VARCHAR(255) NOT NULL UNIQUE 
);

-- Table: Permissions
CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,       
    "name" VARCHAR(255) NOT NULL,  
    "resource" VARCHAR(255) NOT NULL,
    "method" VARCHAR(255) NOT NULL
);

-- Table: Role_Permissions
CREATE TABLE role_Permissions (
    role_id INT NOT NULL,        -- Role ID
    permission_id INT NOT NULL,  -- Permission ID
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES Roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES Permissions(id) ON DELETE CASCADE
);

CREATE TABLE custom_user (
    id BIGSERIAL PRIMARY KEY,
    "password" VARCHAR(128) NOT NULL,
    last_login TIMESTAMPTZ NULL,
    email VARCHAR(254) NOT NULL,
    role_id INT NOT NULL REFERENCES Roles(id) ON DELETE CASCADE DEFAULT 1,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    date_joined TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    profile_picture VARCHAR(100) NULL
);

CREATE INDEX custom_user_email_idx ON custom_user(email);


CREATE TABLE "sessions" (
    "id" uuid PRIMARY KEY,
    "refresh_token" varchar NOT NULL,
    "user_agent" varchar NOT NULL,
    "client_ip" varchar NOT NULL,
    "is_blocked" boolean NOT NULL DEFAULT false,
    "expires_at" timestamptz NOT NULL,
    "created_at" timestamptz NOT NULL,
    "user_id" BIGINT NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY ("user_id") REFERENCES custom_user("id")
);

-- Index on foreign key for faster joins
CREATE INDEX idx_sessions_user ON sessions("user_id");

-- Index on expires_at for efficient cleanup of expired sessions
CREATE INDEX idx_sessions_expires ON sessions("expires_at");

-- Composite index for session validation
CREATE INDEX idx_sessions_token_blocked ON sessions("refresh_token", "is_blocked");




CREATE TABLE assessment_domain (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

CREATE INDEX assessment_domain_name_idx ON assessment_domain(name);

CREATE TABLE assessment (
    id BIGSERIAL PRIMARY KEY,
    content TEXT NULL DEFAULT '',
    domain_id BIGINT NULL REFERENCES assessment_domain(id) ON DELETE CASCADE,
    "level" INTEGER NOT NULL CHECK (level BETWEEN 1 AND 5),
    updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE db_settings (
    id BIGSERIAL PRIMARY KEY,
    option_name VARCHAR(255) UNIQUE NOT NULL,
    option_value VARCHAR(255) NOT NULL DEFAULT '',
    option_type VARCHAR(5) NOT NULL CHECK (option_type IN ('str', 'int', 'float', 'bool')) DEFAULT 'str'
);

CREATE TABLE notification (
    id BIGSERIAL PRIMARY KEY,
    event VARCHAR(50) NOT NULL CHECK (event IN (
        'normal', 'login_send_credentials', 'appointment_created', 
        'appointment_updated', 'appointment_rescheduled', 
        'appointment_canceled', 'invoice_expired', 'invoice_created',
        'progress_report_available', 'progress_report_created',
        'medication_time', 'contract_reminder'
    )) DEFAULT 'normal',
    title VARCHAR(100) NULL,
    content TEXT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    metadata JSONB NULL DEFAULT '{}',
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE attachment_file (
    "uuid" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    "file" VARCHAR(255) NOT NULL,
    "size" INTEGER NOT NULL DEFAULT 0,
    is_used BOOLEAN NOT NULL DEFAULT FALSE,
    tag VARCHAR(100) NULL DEFAULT '',
    updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX attachment_file_is_used_idx ON attachment_file(is_used);
CREATE INDEX attachment_file_created_idx ON attachment_file(created);



CREATE TABLE expense (
    id BIGSERIAL PRIMARY KEY,
    amount DECIMAL(20,2) NOT NULL,
    tax FLOAT NOT NULL DEFAULT 0,
    "desc" TEXT NULL DEFAULT '',
    attachment_ids JSONB NOT NULL DEFAULT '[]',
    location_id BIGINT NULL REFERENCES location(id) ON DELETE SET NULL,
    updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX expense_location_id_idx ON expense(location_id);


CREATE TABLE protected_email (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(254) NOT NULL,
    subject VARCHAR(255) NULL,
    content TEXT NULL,
    email_type VARCHAR(20) NOT NULL CHECK (email_type IN ('incident_report', 'medical_report', 'progress_report')),
    expired_at TIMESTAMPTZ NOT NULL,
    metadata JSONB NULL DEFAULT '{}',
    passkey VARCHAR(30) NOT NULL DEFAULT '',
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX protected_email_email_idx ON protected_email(email);
CREATE INDEX protected_email_expired_at_idx ON protected_email(expired_at);


CREATE TABLE sender (
    id BIGSERIAL PRIMARY KEY,
    types VARCHAR(50) NOT NULL CHECK (types IN (
        'main_provider', 'local_authority', 
        'particular_party', 'healthcare_institution'
    )),
    name VARCHAR(20) NOT NULL,
    address VARCHAR(200) NULL,
    postal_code VARCHAR(20) NULL,
    place VARCHAR(20) NULL,
    land VARCHAR(20) NULL,
    kvknumber VARCHAR(20) NULL,
    btwnumber VARCHAR(20) NULL,
    phone_number VARCHAR(20) NULL,
    client_number VARCHAR(20) NULL,
    email_adress VARCHAR(20) NULL,
    contacts JSONB NOT NULL DEFAULT '[]',
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX sender_types_idx ON sender(types);

CREATE TABLE sender_audit (
    id BIGSERIAL PRIMARY KEY,
    sender_id BIGINT NOT NULL REFERENCES sender(id),
    changed_by VARCHAR(50) NOT NULL, -- 
    changed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    old_data JSONB NOT NULL,
    new_data JSONB NOT NULL

);

CREATE INDEX idx_sender_audit_sender_id ON sender_audit(sender_id);
CREATE INDEX idx_sender_audit_changed_at ON sender_audit(changed_at DESC);
CREATE INDEX idx_sender_audit_changed_by ON sender_audit(changed_by);


CREATE TABLE client_details (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES custom_user(id) ON DELETE CASCADE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE NULL,
    "identity" BOOLEAN NOT NULL DEFAULT FALSE,
    "status" VARCHAR(20) NULL CHECK (status IN ('In Care', 'On Waiting List', 'Out Of Care')) DEFAULT 'On Waiting List',
    bsn VARCHAR(50) NULL,        -- Reduced length assuming it's a social security number
    source VARCHAR(100) NULL,
    birthplace VARCHAR(100) NULL,
    email VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20) NULL,
    organisation VARCHAR(100) NULL,
    departement VARCHAR(100) NULL,
    gender VARCHAR(100) NOT NULL,
    filenumber INTEGER NOT NULL,
    profile_picture VARCHAR(100) NULL,
    infix VARCHAR(100) NULL,
    created TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP,
    sender_id BIGINT NULL REFERENCES sender(id) ON DELETE CASCADE,
    location_id BIGINT NULL REFERENCES location(id) ON DELETE SET NULL,
    identity_attachment_ids JSONB NOT NULL DEFAULT '[]',
    departure_reason VARCHAR(255) NULL,
    departure_report TEXT NULL,
    gps_position JSONB NOT NULL DEFAULT '[]',
    maturity_domains JSONB NOT NULL DEFAULT '[]',
    addresses JSONB NOT NULL DEFAULT '[]',
    legal_measure VARCHAR(255) NULL,
    has_untaken_medications BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX client_details_user_id_idx ON client_details(user_id);
CREATE INDEX client_details_sender_id_idx ON client_details(sender_id);
CREATE INDEX client_details_location_id_idx ON client_details(location_id);


CREATE TABLE client_current_level (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    domain_id BIGINT NOT NULL REFERENCES assessment_domain(id) ON DELETE CASCADE,
    "level" INTEGER NOT NULL DEFAULT 1 CHECK (level BETWEEN 1 AND 5),
    content TEXT NULL DEFAULT '',
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX client_current_level_client_id_idx ON client_current_level(client_id);
CREATE INDEX client_current_level_domain_id_idx ON client_current_level(domain_id);



CREATE TABLE client_status_history (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL CHECK (status IN ('In Care', 'On Waiting List', 'Out Of Care')),
    start_date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX client_status_history_client_id_idx ON client_status_history(client_id);
CREATE INDEX client_status_history_start_date_idx ON client_status_history(start_date);


CREATE TABLE client_state (
    id BIGSERIAL PRIMARY KEY,
    "value" INTEGER NOT NULL DEFAULT 0,
    "type" VARCHAR(20) NOT NULL CHECK (type IN ('emotional', 'physical')),
    content TEXT NULL DEFAULT '',
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX client_state_client_id_idx ON client_state(client_id);



CREATE TABLE client_diagnosis (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(50) NULL,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    diagnosis_code VARCHAR(10) NOT NULL,
    description TEXT NOT NULL,
    date_of_diagnosis TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    severity VARCHAR(50) NULL,
    status VARCHAR(100) NOT NULL,
    diagnosing_clinician VARCHAR(100) NOT NULL,
    notes TEXT NULL,
    created TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX client_diagnosis_client_id_idx ON client_diagnosis(client_id);
CREATE INDEX client_diagnosis_diagnosis_code_idx ON client_diagnosis(diagnosis_code);



CREATE TABLE contact_relationship (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    soft_delete BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX contact_relationship_soft_delete_idx ON contact_relationship(soft_delete);

CREATE TABLE client_emergency_contact (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    first_name VARCHAR(50) NULL,
    last_name VARCHAR(100) NULL,
    email VARCHAR(100) NULL,
    phone_number VARCHAR(20) NULL,
    address VARCHAR(100) NULL,
    relationship VARCHAR(100) NULL,
    relation_status VARCHAR(50) NULL CHECK (relation_status IN ('Primary Relationship', 'Secondary Relationship')),
    created TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    uuid UUID NULL,
    medical_reports BOOLEAN NOT NULL DEFAULT FALSE,
    incidents_reports BOOLEAN NOT NULL DEFAULT FALSE,
    goals_reports BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX client_emergency_contact_client_id_idx ON client_emergency_contact(client_id);
CREATE INDEX client_emergency_contact_uuid_idx ON client_emergency_contact(uuid);





CREATE TABLE treatments (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    treatment_name VARCHAR(500) NOT NULL,
    treatment_date VARCHAR(255) NOT NULL,
    created TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE client_allergy (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    allergy_type VARCHAR(100) NOT NULL,
    severity VARCHAR(100) NOT NULL,
    reaction TEXT NOT NULL,
    notes TEXT NULL,
    created TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE client_documents (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    documents VARCHAR(255) NOT NULL,
    uploaded_at TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP,
    original_filename VARCHAR(255) NULL,
    file_size BIGINT NULL,
    created TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP,
    label VARCHAR(100) NOT NULL CHECK (label IN (
        'registration_form', 'intake_form', 'consent_form',
        'risk_assessment', 'self_reliance_matrix', 'force_inventory',
        'care_plan', 'signaling_plan', 'cooperation_agreement', 'other'
    )) DEFAULT 'other'
);

CREATE INDEX treatments_user_id_idx ON treatments(user_id);
CREATE INDEX client_allergy_client_id_idx ON client_allergy(client_id);
CREATE INDEX client_documents_user_id_idx ON client_documents(user_id);
CREATE INDEX client_documents_label_idx ON client_documents(label);




CREATE TABLE contract_type (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE contract (
    id BIGSERIAL PRIMARY KEY,
    type_id BIGINT NULL REFERENCES contract_type(id) ON DELETE SET NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('approved', 'draft', 'terminated', 'stopped')) DEFAULT 'draft',
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    reminder_period INTEGER NOT NULL DEFAULT 10,
    tax INTEGER NULL DEFAULT -1,
    price DECIMAL(10,2) NOT NULL,
    price_frequency VARCHAR(20) NOT NULL CHECK (price_frequency IN ('minute', 'hourly', 'daily', 'weekly', 'monthly')) DEFAULT 'weekly',
    hours INTEGER NULL DEFAULT 0,
    hours_type VARCHAR(20) NOT NULL CHECK (hours_type IN ('weekly', 'all_period')) DEFAULT 'all_period',
    care_name VARCHAR(255) NOT NULL,
    care_type VARCHAR(20) NOT NULL CHECK (care_type IN ('ambulante', 'accommodation')),
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    sender_id BIGINT NULL REFERENCES sender(id) ON DELETE SET NULL,
    attachment_ids JSONB NOT NULL DEFAULT '[]',
    financing_act VARCHAR(10) NOT NULL CHECK (financing_act IN ('WMO', 'ZVW', 'WLZ', 'JW', 'WPG')) DEFAULT 'WMO',
    financing_option VARCHAR(10) NOT NULL CHECK (financing_option IN ('ZIN', 'PGB')) DEFAULT 'PGB',
    departure_reason VARCHAR(255) NULL,
    departure_report TEXT NULL,
    updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX contract_type_id_idx ON contract(type_id);
CREATE INDEX contract_client_id_idx ON contract(client_id);
CREATE INDEX contract_sender_id_idx ON contract(sender_id);
CREATE INDEX contract_status_idx ON contract(status);


CREATE TABLE contract_working_hours (
    id BIGSERIAL PRIMARY KEY,
    contract_id BIGINT NOT NULL REFERENCES contract(id) ON DELETE CASCADE,
    minutes INTEGER NOT NULL DEFAULT 0,
    "datetime" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    notes TEXT NULL DEFAULT ''
);

CREATE TABLE invoice (
    id BIGSERIAL PRIMARY KEY,
    invoice_number VARCHAR(10) NOT NULL UNIQUE,
    issue_date DATE NOT NULL DEFAULT CURRENT_DATE,
    due_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN (
        'outstanding', 'partially_paid', 'paid', 'expired',
        'overpaid', 'imported', 'concept'
    )) DEFAULT 'concept',
    invoice_details JSONB NULL DEFAULT '[]',
    total_amount DECIMAL(20,2) NOT NULL DEFAULT 0,
    pdf_attachment_id UUID NULL UNIQUE REFERENCES attachment_file("uuid") ON DELETE SET NULL,
    extra_content TEXT NULL DEFAULT '',
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX contract_working_hours_contract_id_idx ON contract_working_hours(contract_id);
CREATE INDEX contract_working_hours_datetime_idx ON contract_working_hours(datetime);
CREATE INDEX invoice_invoice_number_idx ON invoice(invoice_number);
CREATE INDEX invoice_client_id_idx ON invoice(client_id);
CREATE INDEX invoice_status_idx ON invoice(status);



CREATE TABLE invoice_history (
    id BIGSERIAL PRIMARY KEY,
    payment_method VARCHAR(20) NULL CHECK (payment_method IN (
        'bank_transfer', 'credit_card', 'check', 'cash', 'other'
    )),
    amount DECIMAL(20,2) NOT NULL DEFAULT 0,
    updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    invoice_id BIGINT NOT NULL REFERENCES invoice(id) ON DELETE CASCADE
);

CREATE INDEX invoice_history_invoice_id_idx ON invoice_history(invoice_id);



CREATE TABLE contract_attachment (
    id BIGSERIAL PRIMARY KEY,
    contract_id BIGINT NOT NULL REFERENCES contract(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    attachment VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE client_agreement (
    id BIGSERIAL PRIMARY KEY,
    contract_id BIGINT NOT NULL REFERENCES contract(id) ON DELETE CASCADE,
    agreement_details TEXT NOT NULL,
    created TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE provision (
    id BIGSERIAL PRIMARY KEY,
    contract_id BIGINT NOT NULL REFERENCES contract(id) ON DELETE CASCADE,
    provision_details TEXT NOT NULL,
    created TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE framework_agreement (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    agreement_details TEXT NOT NULL,
    created TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX contract_attachment_contract_id_idx ON contract_attachment(contract_id);
CREATE INDEX client_agreement_contract_id_idx ON client_agreement(contract_id);
CREATE INDEX provision_contract_id_idx ON provision(contract_id);
CREATE INDEX framework_agreement_client_id_idx ON framework_agreement(client_id);




CREATE TABLE contact (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    email VARCHAR(254) NOT NULL
);

CREATE TABLE sender_contact_relation (
    id BIGSERIAL PRIMARY KEY,
    client_type_id BIGINT NOT NULL REFERENCES sender(id) ON DELETE CASCADE,
    contact_id BIGINT NOT NULL REFERENCES contact(id) ON DELETE CASCADE
);

CREATE INDEX sender_contact_relation_client_type_id_idx ON sender_contact_relation(client_type_id);
CREATE INDEX sender_contact_relation_contact_id_idx ON sender_contact_relation(contact_id);


CREATE TABLE temporary_file (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file VARCHAR(255) NOT NULL,
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX temporary_file_uploaded_at_idx ON temporary_file(uploaded_at);


CREATE TABLE invoice_contract (
    id BIGSERIAL PRIMARY KEY,
    invoice_id BIGINT NULL REFERENCES invoice(id) ON DELETE SET NULL,
    contract_id BIGINT NULL REFERENCES contract(id) ON DELETE SET NULL,
    pre_vat_total DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    vat_rate DECIMAL(5,2) NOT NULL DEFAULT 20.00,
    vat_amount DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    total_amount DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX invoice_contract_invoice_id_idx ON invoice_contract(invoice_id);
CREATE INDEX invoice_contract_contract_id_idx ON invoice_contract(contract_id);
CREATE INDEX invoice_contract_updated_idx ON invoice_contract(updated);
CREATE INDEX invoice_contract_created_idx ON invoice_contract(created);



CREATE TABLE care_plan (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NULL REFERENCES client_details(id) ON DELETE SET NULL,
    description TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE care_plan_domains (
    care_plan_id BIGINT NOT NULL REFERENCES care_plan(id) ON DELETE CASCADE,
    assessment_domain_id BIGINT NOT NULL REFERENCES assessment_domain(id) ON DELETE CASCADE,
    PRIMARY KEY (care_plan_id, assessment_domain_id)
);

CREATE TABLE careplan_atachements (
    id BIGSERIAL PRIMARY KEY,
    careplan_id BIGINT NULL REFERENCES care_plan(id) ON DELETE SET NULL,
    attachement VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name VARCHAR(100) NULL
);

CREATE INDEX care_plan_client_id_idx ON care_plan(client_id);
CREATE INDEX careplan_atachements_careplan_id_idx ON careplan_atachements(careplan_id);



CREATE TABLE incident (
    id BIGSERIAL PRIMARY KEY,
    employee_fullname VARCHAR(100) NOT NULL,
    employee_position VARCHAR(100) NOT NULL,
    location_id BIGINT NULL REFERENCES location(id) ON DELETE SET NULL,
    reporter_involvement VARCHAR(100) NOT NULL CHECK (reporter_involvement IN (
        'directly_involved', 'witness', 'found_afterwards', 'alarmed'
    )),
    inform_who JSONB NOT NULL DEFAULT '[]',
    incident_date DATE NOT NULL,
    runtime_incident VARCHAR(100) NOT NULL,
    incident_type VARCHAR(100) NOT NULL,
    passing_away BOOLEAN NOT NULL DEFAULT FALSE,
    self_harm BOOLEAN NOT NULL DEFAULT FALSE,
    violence BOOLEAN NOT NULL DEFAULT FALSE,
    fire_water_damage BOOLEAN NOT NULL DEFAULT FALSE,
    accident BOOLEAN NOT NULL DEFAULT FALSE,
    client_absence BOOLEAN NOT NULL DEFAULT FALSE,
    medicines BOOLEAN NOT NULL DEFAULT FALSE,
    organization BOOLEAN NOT NULL DEFAULT FALSE,
    use_prohibited_substances BOOLEAN NOT NULL DEFAULT FALSE,
    other_notifications BOOLEAN NOT NULL DEFAULT FALSE,
    severity_of_incident VARCHAR(100) NOT NULL CHECK (severity_of_incident IN (
        'near_incident', 'less_serious', 'serious', 'fatal'
    )),
    incident_explanation TEXT NULL,
    recurrence_risk VARCHAR(100) NOT NULL CHECK (recurrence_risk IN (
        'very_low', 'means', 'high', 'very_high'
    )),
    incident_prevent_steps TEXT NULL,
    incident_taken_measures TEXT NULL,
    technical JSONB NOT NULL DEFAULT '[]',
    organizational JSONB NOT NULL DEFAULT '[]',
    mese_worker JSONB NOT NULL DEFAULT '[]',
    client_options JSONB NOT NULL DEFAULT '[]',
    other_cause VARCHAR(100) NULL,
    cause_explanation TEXT NULL DEFAULT '',
    physical_injury VARCHAR(100) NOT NULL CHECK (physical_injury IN (
        'no_injuries', 'not_noticeable_yet', 'bruising_swelling', 'skin_injury',
        'broken_bones', 'shortness_of_breath', 'death', 'other'
    )),
    physical_injury_desc TEXT NULL DEFAULT '',
    psychological_damage VARCHAR(100) NOT NULL CHECK (psychological_damage IN (
        'no', 'not_noticeable_yet', 'drowsiness', 'unrest', 'other'
    )),
    psychological_damage_desc TEXT NULL DEFAULT '',
    needed_consultation VARCHAR(100) NOT NULL CHECK (needed_consultation IN (
        'no', 'not_clear', 'hospitalization', 'consult_gp'
    )),
    succession JSONB NOT NULL DEFAULT '[]',
    succession_desc TEXT NULL DEFAULT '',
    other BOOLEAN NOT NULL DEFAULT FALSE,
    other_desc VARCHAR(100) NULL,
    additional_appointments TEXT NULL DEFAULT '',
    employee_absenteeism JSONB NOT NULL DEFAULT '[]',
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    soft_delete BOOLEAN NOT NULL DEFAULT FALSE,
    updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX incident_client_id_idx ON incident(client_id);
CREATE INDEX incident_location_id_idx ON incident(location_id);
CREATE INDEX incident_soft_delete_idx ON incident(soft_delete);




CREATE TABLE collaboration_agreement (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    client_full_name VARCHAR(100) NOT NULL,
    client_skn VARCHAR(100) NOT NULL,
    client_number VARCHAR(100) NOT NULL,
    client_phone VARCHAR(100) NOT NULL,
    probation_full_name VARCHAR(100) NOT NULL,
    probation_organization VARCHAR(100) NOT NULL,
    probation_phone VARCHAR(100) NOT NULL,
    healthcare_institution_name VARCHAR(100) NOT NULL,
    healthcare_institution_organization VARCHAR(100) NOT NULL,
    healthcare_institution_phone VARCHAR(100) NOT NULL,
    healthcare_institution_function VARCHAR(100) NOT NULL,
    contact_agreements TEXT NOT NULL,
    pdf_attachment_id UUID NULL UNIQUE REFERENCES attachment_file("uuid") ON DELETE SET NULL,
    attention_risks JSONB NOT NULL DEFAULT '[]',
    updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX collaboration_agreement_client_id_idx ON collaboration_agreement(client_id);



CREATE TABLE risk_assessment (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    date_of_birth DATE NOT NULL,
    gender VARCHAR(100) NOT NULL,
    date_of_intake TIMESTAMPTZ NOT NULL,
    intaker_position_name VARCHAR(100) NOT NULL,
    family_situation TEXT NOT NULL,
    education_work TEXT NOT NULL,
    current_living_situation TEXT NOT NULL,
    social_network TEXT NOT NULL,
    previous_assistance TEXT NOT NULL,
    behaviour_at_school_work TEXT NOT NULL,
    people_skills TEXT NOT NULL,
    emotional_status TEXT NOT NULL,
    self_image_self_confidence TEXT NOT NULL,
    stress_factors TEXT NOT NULL,
    committed_offences_description TEXT NOT NULL,
    offences_frequency_seriousness TEXT NOT NULL,
    age_first_offense TEXT NOT NULL,
    circumstances_surrounding_crimes TEXT NOT NULL,
    offenses_recations TEXT NOT NULL,
    personal_risk_factors TEXT NOT NULL,
    environmental_risk_factors TEXT NOT NULL,
    behaviour_recurrence_risk TEXT NOT NULL,
    abuse_substance_risk TEXT NOT NULL,
    person_strengths TEXT NOT NULL,
    positive_influences TEXT NOT NULL,
    available_support_assistance TEXT NOT NULL,
    person_strategies TEXT NOT NULL,
    specific_needs TEXT NOT NULL,
    recommended_interventions TEXT NOT NULL,
    other_agencies_involvement TEXT NOT NULL,
    risk_management_plan_of_actions TEXT NOT NULL,
    findings_summary TEXT NOT NULL,
    institution_advice TEXT NOT NULL,
    inclusion TEXT NOT NULL,
    intaker_name VARCHAR(100) NOT NULL,
    report_date DATE NOT NULL,
    regular_evaluation_plan VARCHAR(255) NOT NULL,
    success_criteria VARCHAR(255) NOT NULL,
    time_table VARCHAR(255) NOT NULL,
    pdf_attachment_id UUID NULL UNIQUE REFERENCES attachment_file("uuid") ON DELETE SET NULL,
    updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX risk_assessment_client_id_idx ON risk_assessment(client_id);


CREATE TABLE consent_declaration (
    id BIGSERIAL PRIMARY KEY,
    youth_name VARCHAR(255) NOT NULL,
    date_of_birth DATE NOT NULL,
    parent_guardian_name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    youth_care_institution VARCHAR(255) NOT NULL,
    proposed_assistance_description TEXT NOT NULL,
    statement_by_representative TEXT NOT NULL,
    parent_guardian_signature_date DATE NOT NULL,
    juvenile_name VARCHAR(255) NULL,
    juvenile_signature_date DATE NULL,
    representative_name VARCHAR(255) NOT NULL,
    representative_signature_date DATE NOT NULL,
    contact_person_name VARCHAR(255) NOT NULL,
    contact_phone_number VARCHAR(20) NOT NULL,
    contact_email VARCHAR(254) NOT NULL,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    pdf_attachment_id UUID NULL UNIQUE REFERENCES attachment_file("uuid") ON DELETE SET NULL
);

CREATE INDEX consent_declaration_client_id_idx ON consent_declaration(client_id);



CREATE TABLE youth_care_intake (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    date_of_birth DATE NOT NULL,
    gender VARCHAR(30) NOT NULL,
    nationality VARCHAR(100) NOT NULL,
    bsn VARCHAR(20) NOT NULL,
    address TEXT NOT NULL,
    postcode VARCHAR(20) NOT NULL,
    residence VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    email VARCHAR(254) NOT NULL,
    referrer_name VARCHAR(255) NOT NULL,
    referrer_organization VARCHAR(255) NOT NULL,
    referrer_function VARCHAR(255) NOT NULL,
    referrer_phone_number VARCHAR(20) NOT NULL,
    referrer_email VARCHAR(254) NOT NULL,
    service_choice VARCHAR(20) NOT NULL CHECK (service_choice IN ('outpatient_care', 'sheltered_housing', 'assisted_living')),
    financing_acts VARCHAR(20) NOT NULL CHECK (financing_acts IN ('WMO', 'ZVW', 'WLZ', 'JW', 'WPG')),
    financing_options VARCHAR(20) NOT NULL CHECK (financing_options IN ('ZIN', 'PGB')),
    financing_other VARCHAR(255) NULL,
    registration_reason TEXT NOT NULL,
    current_situation_background TEXT NOT NULL,
    previous_aid_agencies_involved BOOLEAN NOT NULL DEFAULT FALSE,
    previous_aid_agencies_details TEXT NULL,
    medical_conditions BOOLEAN NOT NULL DEFAULT FALSE,
    medical_conditions_details TEXT NULL,
    medication_use BOOLEAN NOT NULL DEFAULT FALSE,
    medication_details TEXT NULL,
    allergies_or_dietary_needs BOOLEAN NOT NULL DEFAULT FALSE,
    allergies_or_dietary_details TEXT NULL,
    addictions BOOLEAN NOT NULL DEFAULT FALSE,
    addictions_details TEXT NULL,
    school_or_daytime_activities BOOLEAN NOT NULL DEFAULT FALSE,
    school_daytime_name VARCHAR(255) NULL,
    current_class_level VARCHAR(100) NULL,
    school_contact_person VARCHAR(255) NULL,
    school_contact_phone VARCHAR(20) NULL,
    school_contact_email VARCHAR(254) NULL,
    important_people TEXT NULL,
    external_supervisors_involved BOOLEAN NOT NULL DEFAULT FALSE,
    external_supervisors_details TEXT NULL,
    special_circumstances TEXT NULL,
    updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX youth_care_intake_client_id_idx ON youth_care_intake(client_id);



CREATE TABLE data_sharing_statement (
    id BIGSERIAL PRIMARY KEY,
    youth_name VARCHAR(255) NOT NULL,
    date_of_birth DATE NOT NULL,
    parent_guardian_name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    youth_care_institution VARCHAR(255) NOT NULL,
    data_description TEXT NOT NULL,
    data_purpose TEXT NOT NULL,
    third_party_names TEXT NOT NULL,
    "statement" TEXT NOT NULL,
    parent_guardian_signature_name VARCHAR(255) NOT NULL,
    parent_guardian_signature VARCHAR(255) NOT NULL,
    parent_guardian_signature_date DATE NOT NULL,
    juvenile_name VARCHAR(255) NULL,
    juvenile_signature_date DATE NULL,
    institution_representative_name VARCHAR(255) NOT NULL,
    institution_representative_signature VARCHAR(255) NOT NULL,
    institution_representative_signature_date DATE NOT NULL,
    contact_person_name VARCHAR(255) NOT NULL,
    contact_phone_number VARCHAR(20) NOT NULL,
    contact_email VARCHAR(254) NOT NULL,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE maturity_matrix (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_approved BOOLEAN NOT NULL DEFAULT FALSE,
    updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE selected_maturity_matrix_assessment (
    id BIGSERIAL PRIMARY KEY,  -- Adding the id column
    maturity_matrix_id BIGINT NOT NULL REFERENCES maturity_matrix(id) ON DELETE CASCADE,
    assessment_id BIGINT NOT NULL REFERENCES assessment(id) ON DELETE CASCADE
    
);

CREATE INDEX data_sharing_statement_client_id_idx ON data_sharing_statement(client_id);
CREATE INDEX maturity_matrix_client_id_idx ON maturity_matrix(client_id);




CREATE TABLE ai_generated_report (
    id BIGSERIAL PRIMARY KEY,
    report_type VARCHAR(100) NOT NULL CHECK (report_type IN (
        'client_reports_summary',
        'client_profile_summary',
        'client_goals_and_objectives_summary',
        'employee_performance'
    )),
    title VARCHAR(255) NOT NULL,
    content TEXT NULL,
    user_id BIGINT NULL REFERENCES custom_user(id) ON DELETE CASCADE,
    user_type VARCHAR(20) NOT NULL CHECK (user_type IN ('client', 'employee')),
    start_date TIMESTAMPTZ NULL,
    end_date TIMESTAMPTZ NULL,
    updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX ai_generated_report_user_id_idx ON ai_generated_report(user_id);
CREATE INDEX ai_generated_report_created_idx ON ai_generated_report(created DESC);
CREATE INDEX ai_generated_report_updated_idx ON ai_generated_report(updated);



CREATE TABLE employee_profile (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES custom_user(id) ON DELETE CASCADE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    position VARCHAR(100) NULL,
    department VARCHAR(100) NULL,
    employee_number VARCHAR(50) NULL,
    employment_number VARCHAR(50) NULL,
    private_email_address VARCHAR(254) NULL,
    email VARCHAR(254) NOT NULL,
    authentication_phone_number VARCHAR(100) NULL,
    private_phone_number VARCHAR(100) NULL,
    work_phone_number VARCHAR(100) NULL,
    date_of_birth DATE NULL,
    home_telephone_number VARCHAR(100) NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_subcontractor BOOLEAN NULL,
    gender VARCHAR(100) NULL,
    location_id BIGINT NULL REFERENCES location(id) ON DELETE SET NULL,
    has_borrowed BOOLEAN NOT NULL DEFAULT FALSE,
    out_of_service BOOLEAN NULL DEFAULT FALSE,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE
);


CREATE INDEX employee_profile_user_id_idx ON employee_profile(user_id);
CREATE INDEX employee_profile_location_id_idx ON employee_profile(location_id);
CREATE INDEX employee_profile_id_desc_idx ON employee_profile(id DESC);
CREATE INDEX idx_employee_profile_is_archived ON employee_profile(is_archived);
CREATE INDEX idx_employee_profile_out_of_service ON employee_profile(out_of_service);



CREATE TABLE certification (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    issued_by VARCHAR(255) NOT NULL,
    date_issued DATE NOT NULL,
    created TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX certification_employee_id_idx ON certification(employee_id);


CREATE TABLE experience (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    job_title VARCHAR(255) NOT NULL,
    company_name VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NULL,
    description TEXT NULL,
    created TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE education (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    institution_name VARCHAR(255) NOT NULL,
    degree VARCHAR(100) NOT NULL,
    field_of_study VARCHAR(100) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NULL
);

CREATE TABLE assignment (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    start_datetime TIMESTAMPTZ NOT NULL,
    end_datetime TIMESTAMPTZ NOT NULL,
    status VARCHAR(50) NOT NULL CHECK (status IN ('Confirmed', 'Pending', 'Cancelled')),
    created TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX experience_employee_id_idx ON experience(employee_id);
CREATE INDEX education_employee_id_idx ON education(employee_id);
CREATE INDEX assignment_employee_id_idx ON assignment(employee_id);
CREATE INDEX assignment_client_id_idx ON assignment(client_id);



CREATE TABLE client_employee_assignment (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    employee_id BIGINT NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    role VARCHAR(100) NOT NULL,
    created TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE progress_report (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    title VARCHAR(50) NULL,
    report_text TEXT NOT NULL,
    author_id BIGINT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL CHECK (type IN (
        'morning_report', 'evening_report', 'night_report', 'shift_report',
        'one_to_one_report', 'process_report', 'contact_journal', 'other'
    )) DEFAULT 'other',
    emotional_state VARCHAR(20) NOT NULL CHECK (emotional_state IN (
        'normal', 'excited', 'happy', 'sad', 'angry', 'anxious', 'depressed'
    )) DEFAULT 'normal',
    created TIMESTAMPTZ NULL
);

CREATE INDEX client_employee_assignment_client_id_idx ON client_employee_assignment(client_id);
CREATE INDEX client_employee_assignment_employee_id_idx ON client_employee_assignment(employee_id);
CREATE INDEX progress_report_client_id_idx ON progress_report(client_id);
CREATE INDEX progress_report_author_id_idx ON progress_report(author_id);
CREATE INDEX progress_report_created_idx ON progress_report(created DESC);


CREATE TABLE measurement (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    measurement_type VARCHAR(100) NOT NULL,
    value FLOAT NOT NULL,
    created TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE observations (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    category VARCHAR(100) NULL,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    observation_text TEXT NOT NULL,
    created TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE feedback (
    id BIGSERIAL PRIMARY KEY,
    author_id BIGINT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    feedback_text TEXT NOT NULL,
    created TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE emotional_state (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    date TIMESTAMPTZ NOT NULL,
    state_description TEXT NOT NULL,
    intensity INTEGER NOT NULL,
    created TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE physical_state (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    date TIMESTAMPTZ NOT NULL,
    symptoms TEXT NOT NULL,
    severity INTEGER NOT NULL,
    created TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE client_medication (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    dosage VARCHAR(100) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NULL,
    notes TEXT NULL,
    self_administered BOOLEAN NOT NULL DEFAULT TRUE,
    slots JSONB NULL DEFAULT '[]',
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    administered_by_id BIGINT NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    is_critical BOOLEAN NOT NULL DEFAULT FALSE,
    updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX measurement_client_id_idx ON measurement(client_id);
CREATE INDEX observations_client_id_idx ON observations(client_id);
CREATE INDEX feedback_client_id_idx ON feedback(client_id);
CREATE INDEX feedback_author_id_idx ON feedback(author_id);
CREATE INDEX emotional_state_client_id_idx ON emotional_state(client_id);
CREATE INDEX physical_state_client_id_idx ON physical_state(client_id);
CREATE INDEX client_medication_client_id_idx ON client_medication(client_id);
CREATE INDEX client_medication_administered_by_id_idx ON client_medication(administered_by_id);
CREATE INDEX client_medication_updated_idx ON client_medication(updated);
CREATE INDEX client_medication_created_idx ON client_medication(created);



CREATE TABLE client_medication_record (
    id BIGSERIAL PRIMARY KEY,
    client_medication_id BIGINT NOT NULL REFERENCES client_medication(id) ON DELETE CASCADE,
    status VARCHAR(20) NULL CHECK (status IN ('taken', 'not_taken', 'awaiting')) DEFAULT 'awaiting',
    reason TEXT NULL DEFAULT '',
    time TIMESTAMPTZ NOT NULL,
    updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX client_medication_record_client_medication_id_idx 
    ON client_medication_record(client_medication_id);
CREATE INDEX client_medication_record_time_idx ON client_medication_record(time);
CREATE INDEX client_medication_record_updated_idx ON client_medication_record(updated);
CREATE INDEX client_medication_record_created_idx ON client_medication_record(created DESC);


CREATE TABLE client_goals (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    goal_name VARCHAR(100) NOT NULL,
    goal_details VARCHAR(500) NOT NULL,
    created_at TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP,
    administered_by_id BIGINT NULL REFERENCES employee_profile(id) ON DELETE SET NULL
);

CREATE TABLE goals_report (
    id BIGSERIAL PRIMARY KEY,
    goal_id BIGINT NULL REFERENCES client_goals(id) ON DELETE SET NULL,
    title VARCHAR(100) NOT NULL,
    report_text TEXT NOT NULL,
    rating INTEGER NULL,
    created_at_sys TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NULL
);

CREATE TABLE ai_generated_weekly_reports (
    id BIGSERIAL PRIMARY KEY,
    report_text TEXT NOT NULL,
    goal_id BIGINT NULL REFERENCES client_goals(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE weekly_report_summary (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    summary_text TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX client_goals_client_id_idx ON client_goals(client_id);
CREATE INDEX client_goals_administered_by_id_idx ON client_goals(administered_by_id);
CREATE INDEX goals_report_goal_id_idx ON goals_report(goal_id);
CREATE INDEX ai_generated_weekly_reports_goal_id_idx ON ai_generated_weekly_reports(goal_id);
CREATE INDEX weekly_report_summary_client_id_idx ON weekly_report_summary(client_id);
CREATE INDEX weekly_report_summary_created_at_idx ON weekly_report_summary(created_at DESC);


CREATE TABLE incident_details (
    id BIGSERIAL PRIMARY KEY,
    reported_by_id BIGINT NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    date_reported TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    date_of_incident TIMESTAMPTZ NOT NULL,
    location VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    action_taken TEXT NULL,
    follow_up_required BOOLEAN NOT NULL DEFAULT FALSE,
    follow_up_date DATE NULL,
    notes TEXT NULL,
    status VARCHAR(100) NOT NULL CHECK (status IN (
        'Reported', 'Under Investigation', 'Resolved', 'Closed'
    )) DEFAULT 'Reported'
);

CREATE TABLE incident_involved_children (
    incident_id BIGINT NOT NULL REFERENCES incident_details(id) ON DELETE CASCADE,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    PRIMARY KEY (incident_id, client_id)
);

CREATE INDEX incident_details_reported_by_id_idx ON incident_details(reported_by_id);
CREATE INDEX incident_details_date_of_incident_idx ON incident_details(date_of_incident);
CREATE INDEX incident_details_status_idx ON incident_details(status);


CREATE TABLE domain_goal (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    "desc" TEXT NULL DEFAULT '',
    domain_id BIGINT NULL REFERENCES assessment_domain(id) ON DELETE SET NULL,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    created_by_id BIGINT NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    reviewed_by_id BIGINT NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    selected_maturity_matrix_assessment_id BIGINT NULL REFERENCES selected_maturity_matrix_assessment(id) ON DELETE CASCADE,
    is_approved BOOLEAN NOT NULL DEFAULT FALSE,
    updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX domain_goal_domain_id_idx ON domain_goal(domain_id);
CREATE INDEX domain_goal_client_id_idx ON domain_goal(client_id);
CREATE INDEX domain_goal_created_by_id_idx ON domain_goal(created_by_id);
CREATE INDEX domain_goal_created_idx ON domain_goal(created DESC);



CREATE TABLE domain_objective (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    "desc" TEXT NULL DEFAULT '',
    rating FLOAT NOT NULL DEFAULT 0,
    goal_id BIGINT NULL REFERENCES domain_goal(id) ON DELETE SET NULL,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX domain_objective_goal_id_idx ON domain_objective(goal_id);
CREATE INDEX domain_objective_client_id_idx ON domain_objective(client_id);



CREATE TABLE objective_history (
    id BIGSERIAL PRIMARY KEY,
    rating FLOAT NOT NULL DEFAULT 0,
    week INTEGER NOT NULL,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    objective_id BIGINT NOT NULL REFERENCES domain_objective(id) ON DELETE CASCADE,
    content TEXT NULL DEFAULT '',
    UNIQUE(week, objective_id)
);

CREATE INDEX objective_history_objective_id_idx ON objective_history(objective_id);
CREATE INDEX objective_history_week_idx ON objective_history(week);
CREATE INDEX objective_history_date_idx ON objective_history(date);



CREATE TABLE objective_progress_report (
    id BIGSERIAL PRIMARY KEY,
    objective_id BIGINT NOT NULL REFERENCES domain_objective(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    report_text TEXT NULL,
    rating FLOAT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE goal_history (
    id BIGSERIAL PRIMARY KEY,
    rating FLOAT NOT NULL DEFAULT 0,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    goal_id BIGINT NOT NULL REFERENCES domain_goal(id) ON DELETE CASCADE
);

CREATE TABLE group_access (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    group_id BIGINT NOT NULL REFERENCES "group"(id) ON DELETE CASCADE,
    start_date TIMESTAMPTZ NULL,
    end_date TIMESTAMPTZ NULL,
    updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX objective_progress_report_objective_id_idx ON objective_progress_report(objective_id);
CREATE INDEX objective_progress_report_created_at_idx ON objective_progress_report(created_at DESC);
CREATE INDEX goal_history_goal_id_idx ON goal_history(goal_id);
CREATE INDEX goal_history_date_idx ON goal_history(date);
CREATE INDEX group_access_employee_id_idx ON group_access(employee_id);
CREATE INDEX group_access_group_id_idx ON group_access(group_id);
CREATE INDEX group_access_created_idx ON group_access(created DESC);