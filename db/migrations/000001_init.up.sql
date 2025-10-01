-- ==========================================
-- INFRASTRUCTURE & ORGANIZATIONS
-- ==========================================

-- Organizations and their locations
CREATE TABLE organisations (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    address VARCHAR(200) NOT NULL,
    postal_code VARCHAR(20) NOT NULL,
    city VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20) NULL,
    email VARCHAR(100) NULL,
    kvk_number VARCHAR(20) NULL,
    btw_number VARCHAR(20) NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Location represents a physical place (care home, apartment building, etc.) for the youth intake
CREATE TABLE location (
    id BIGSERIAL PRIMARY KEY,
    organisation_id BIGINT NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    address VARCHAR(100) NOT NULL,
    capacity INTEGER NULL,
    location_type VARCHAR(50) NOT NULL CHECK (location_type IN ('care_home', 'office', 'other')) DEFAULT 'other',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Rooms within a location (studios, apartments, meeting rooms, etc.)
CREATE TABLE room (
    id BIGSERIAL PRIMARY KEY,
    location_id BIGINT NOT NULL,
    room_number VARCHAR(20) NOT NULL,
    room_name VARCHAR(100),
    room_type VARCHAR(50),
    capacity INTEGER,
    is_occupied BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (location_id) REFERENCES location(id) ON DELETE CASCADE,
    UNIQUE(location_id, room_number)
);

-- Standard shifts for locations
CREATE TABLE location_shift (
    id BIGSERIAL PRIMARY KEY,
    location_id BIGINT NOT NULL REFERENCES location(id) ON DELETE CASCADE,
    shift_name VARCHAR(50) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(location_id, shift_name)
);

-- Function to insert default shifts for new locations
CREATE OR REPLACE FUNCTION insert_default_shifts()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO location_shift (location_id, shift_name, start_time, end_time)
    VALUES 
        (NEW.id, 'Ochtenddienst', TIME '07:30:00', TIME '15:30:00'),
        (NEW.id, 'Avonddienst', TIME '15:00:00', TIME '23:00:00'),
        (NEW.id, 'Slaapdienst of Waakdienst', TIME '23:00:00', TIME '07:30:00');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_insert_default_shifts
AFTER INSERT ON location
FOR EACH ROW
EXECUTE FUNCTION insert_default_shifts();

-- ==========================================
-- USER AUTHENTICATION & PERMISSIONS
-- ==========================================

-- Role templates
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

-- System permissions
CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    resource VARCHAR(255) NOT NULL,
    method VARCHAR(255) NOT NULL
);

-- Role-to-Permission mapping (template)
CREATE TABLE role_permissions (
    role_id INT NOT NULL,
    permission_id INT NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

-- User authentication data
CREATE TABLE custom_user (
    id BIGSERIAL PRIMARY KEY,
    password VARCHAR(128) NOT NULL,
    last_login TIMESTAMPTZ,
    email VARCHAR(254) NOT NULL UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    date_joined TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    profile_picture VARCHAR(100),
    two_factor_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    two_factor_secret VARCHAR(100),
    two_factor_secret_temp VARCHAR(100),
    recovery_codes TEXT[] NOT NULL DEFAULT '{}'
);

CREATE INDEX custom_user_email_idx ON custom_user(email);
CREATE INDEX custom_user_id_idx ON custom_user(id);

-- Direct user-to-permission assignments
CREATE TABLE user_permissions (
    user_id BIGINT NOT NULL,
    permission_id INT NOT NULL,
    PRIMARY KEY (user_id, permission_id),
    FOREIGN KEY (user_id) REFERENCES custom_user(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

-- Track which role templates were given to a user
CREATE TABLE user_roles (
    user_id BIGINT NOT NULL PRIMARY KEY,
    role_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES custom_user(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

-- Session management for refresh tokens
CREATE TABLE "sessions" (
    "id" uuid PRIMARY KEY,
    "refresh_token" varchar NOT NULL,
    "user_agent" varchar NOT NULL,
    "client_ip" varchar NOT NULL,
    "is_blocked" boolean NOT NULL DEFAULT false,
    "expires_at" timestamptz NOT NULL,
    "created_at" timestamptz NOT NULL,
    "user_id" BIGINT NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY ("user_id") REFERENCES custom_user("id") ON DELETE CASCADE
);

CREATE INDEX idx_sessions_user ON sessions("user_id");
CREATE INDEX idx_sessions_expires ON sessions("expires_at");
CREATE INDEX idx_sessions_token_blocked ON sessions("refresh_token", "is_blocked");

-- Notifications for users
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id BIGINT NOT NULL REFERENCES custom_user(id) ON DELETE CASCADE,
    type VARCHAR(100) NOT NULL CHECK (type IN (
        'new_appointment', 'appointment_update', 'new_client_assigned',
        'client_goal_update', 'incident_report', 'client_contract_reminder',
        'new_schedule_notification'
    )),
    message TEXT NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    data JSONB NULL,
    read_at TIMESTAMPTZ NULL DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_user_id_created_at ON notifications (user_id, created_at DESC);
CREATE INDEX idx_notifications_user_id_read_at ON notifications (user_id, read_at);

-- ==========================================
-- FILE MANAGEMENT
-- ==========================================

-- Attachment files
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

-- Temporary file storage
CREATE TABLE temporary_file (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file VARCHAR(255) NOT NULL,
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX temporary_file_uploaded_at_idx ON temporary_file(uploaded_at);

-- ==========================================
-- EMPLOYEE MANAGEMENT
-- ==========================================

-- Employee profile (linked to custom_user)
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
    gender VARCHAR(20) NULL CHECK (gender IN ('male', 'female', 'not_specified')),
    location_id BIGINT NULL REFERENCES location(id) ON DELETE SET NULL,
    has_borrowed BOOLEAN NOT NULL DEFAULT FALSE,
    out_of_service BOOLEAN NULL DEFAULT FALSE,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    contract_hours FLOAT NULL DEFAULT 0.0,
    contract_end_date DATE NULL,
    contract_start_date DATE NULL,
    contract_type VARCHAR(50) NULL CHECK (contract_type IN ('loondienst', 'ZZP', 'none')) DEFAULT 'none',
    contract_rate DECIMAL(10,2) NULL DEFAULT 0.00
);

CREATE INDEX employee_profile_user_id_idx ON employee_profile(user_id);
CREATE INDEX employee_profile_location_id_idx ON employee_profile(location_id);
CREATE INDEX employee_profile_id_desc_idx ON employee_profile(id DESC);
CREATE INDEX idx_employee_profile_is_archived ON employee_profile(is_archived);
CREATE INDEX idx_employee_profile_out_of_service ON employee_profile(out_of_service);

-- Employee education records
CREATE TABLE employee_education (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    institution_name VARCHAR(255) NOT NULL,
    degree VARCHAR(100) NOT NULL,
    field_of_study VARCHAR(100) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX education_employee_id_idx ON employee_education(employee_id);

-- Employee certifications
CREATE TABLE certification (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    issued_by VARCHAR(255) NOT NULL,
    date_issued DATE NOT NULL,
    created_at TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX certification_employee_id_idx ON certification(employee_id);

-- Employee work experience
CREATE TABLE employee_experience (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    job_title VARCHAR(255) NOT NULL,
    company_name VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NULL,
    description TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX experience_employee_id_idx ON employee_experience(employee_id);

-- ==========================================
-- CLIENT MANAGEMENT & INTAKE
-- ==========================================

-- Sender organizations (municipalities, authorities, etc.)
CREATE TABLE sender (
    id BIGSERIAL PRIMARY KEY,
    types VARCHAR(50) NOT NULL CHECK (types IN (
        'main_provider', 'local_authority', 
        'particular_party', 'healthcare_institution'
    )),
    name VARCHAR(60) NOT NULL,
    address VARCHAR(200) NULL,
    postal_code VARCHAR(20) NULL,
    place VARCHAR(20) NULL,
    land VARCHAR(20) NULL,
    kvknumber VARCHAR(20) NULL,
    btwnumber VARCHAR(20) NULL, 
    phone_number VARCHAR(20) NULL,
    client_number VARCHAR(20) NULL,
    email_address VARCHAR(20) NULL,
    contacts JSONB NOT NULL DEFAULT '[]',
    invoice_template BIGINT[] NULL,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX sender_types_idx ON sender(types);

-- Contact information
CREATE TABLE contact (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    email VARCHAR(254) NOT NULL
);

-- Sender-Contact relationships
CREATE TABLE sender_contact_relation (
    id BIGSERIAL PRIMARY KEY,
    client_type_id BIGINT NOT NULL REFERENCES sender(id) ON DELETE CASCADE,
    contact_id BIGINT NOT NULL REFERENCES contact(id) ON DELETE CASCADE
);

CREATE INDEX sender_contact_relation_client_type_id_idx ON sender_contact_relation(client_type_id);
CREATE INDEX sender_contact_relation_contact_id_idx ON sender_contact_relation(contact_id);

-- Initial intake forms
CREATE TABLE intake_forms (
    id BIGSERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE NOT NULL,
    nationality VARCHAR(100) NOT NULL,
    bsn VARCHAR(20) NOT NULL,
    address TEXT NOT NULL,
    city VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    gender VARCHAR(20) NOT NULL,
    email VARCHAR(254) NOT NULL,
    id_type VARCHAR(100) NOT NULL CHECK (id_type IN ('passport', 'id_card', 'residence_permit')),
    id_number VARCHAR(100) NOT NULL,
    -- Referrer information
    referrer_name VARCHAR(255),
    referrer_organization VARCHAR(255),
    referrer_function VARCHAR(150),
    referrer_phone VARCHAR(50),
    referrer_email VARCHAR(255),
    signed_by VARCHAR(50) CHECK (signed_by IN ('Referrer', 'Parent/Guardian', 'Client')),
    -- Care indication details
    has_valid_indication BOOLEAN NOT NULL DEFAULT FALSE,
    law_type VARCHAR(50) CHECK (law_type IN ('Youth Act', 'WLZ', 'WMO', 'Other')),
    other_law_specification VARCHAR(255),
    main_provider_name VARCHAR(255),
    main_provider_contact VARCHAR(255),
    indication_start_date DATE,
    indication_end_date DATE,
    registration_reason TEXT,
    guidance_goals TEXT,
    registration_type VARCHAR(50) CHECK (registration_type IN ('Protected Living', 'Supervised Independent Living', 'Outpatient Guidance')),
    -- Client current situation
    living_situation VARCHAR(50) CHECK (living_situation IN ('Home', 'Foster care', 'Youth care institution', 'Other')),
    other_living_situation VARCHAR(255),
    parental_authority BOOLEAN NOT NULL DEFAULT FALSE,
    current_school VARCHAR(255),
    mentor_name VARCHAR(255),
    mentor_phone VARCHAR(50),
    mentor_email VARCHAR(255),
    previous_care TEXT,
    guardian_details JSONB NOT NULL DEFAULT '{}',
    -- Medical and psychosocial information
    diagnoses TEXT,
    uses_medication BOOLEAN NOT NULL DEFAULT FALSE,
    medication_details TEXT,
    addiction_issues BOOLEAN NOT NULL DEFAULT FALSE,
    judicial_involvement BOOLEAN NOT NULL DEFAULT FALSE,
    -- Risk factors
    risk_aggression BOOLEAN NOT NULL DEFAULT FALSE,
    risk_suicidality BOOLEAN NOT NULL DEFAULT FALSE,
    risk_running_away BOOLEAN NOT NULL DEFAULT FALSE,
    risk_self_harm BOOLEAN NOT NULL DEFAULT FALSE,
    risk_weapon_possession BOOLEAN NOT NULL DEFAULT FALSE,
    risk_drug_dealing BOOLEAN NOT NULL DEFAULT FALSE,
    other_risks TEXT,
    -- Consent and signatures
    sharing_permission BOOLEAN NOT NULL DEFAULT FALSE,
    truth_declaration BOOLEAN NOT NULL DEFAULT FALSE,
    client_signature BOOLEAN NOT NULL DEFAULT FALSE,
    guardian_signature BOOLEAN,
    referrer_signature BOOLEAN,
    signature_date DATE,
    -- Status and processing
    status VARCHAR(20) NOT NULL CHECK (status IN ('submitted', 'under review', 'approved', 'rejected')) DEFAULT 'submitted',
    urgency_score varchar(20) NOT NULL DEFAULT 'low' CHECK (urgency_score IN ('low', 'medium', 'high')),
    description TEXT,
    attachement_ids UUID[] NULL DEFAULT '{}',
    is_in_waiting_list BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Main client details table
CREATE TABLE client_details (
    id BIGSERIAL PRIMARY KEY,
    intake_form_id BIGINT NULL REFERENCES intake_forms(id) ON DELETE SET NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE NULL,
    "identity" BOOLEAN NOT NULL DEFAULT FALSE,
    "status" VARCHAR(20) NULL CHECK (status IN ('In Care', 'On Waiting List', 'Out Of Care')) DEFAULT 'On Waiting List',
    bsn VARCHAR(50) NULL,
    bsn_verified_by BIGINT NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    source VARCHAR(100) NULL,
    birthplace VARCHAR(100) NULL,
    email VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20) NULL,
    organisation VARCHAR(100) NULL,
    departement VARCHAR(100) NULL,
    gender VARCHAR(100) NOT NULL CHECK (gender IN ('male', 'female', 'other')),
    filenumber VARCHAR(100) NOT NULL,
    profile_picture VARCHAR(600) NULL,
    infix VARCHAR(100) NULL,
    created_at TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP,
    sender_id BIGINT NULL REFERENCES sender(id) ON DELETE SET NULL DEFAULT NULL,
    location_id BIGINT NULL REFERENCES location(id) ON DELETE SET NULL DEFAULT NULL,
    departure_reason VARCHAR(255) NULL,
    departure_report TEXT NULL,
    gps_position JSONB NOT NULL DEFAULT '[]',
    maturity_domains JSONB NOT NULL DEFAULT '[]',
    addresses JSONB NOT NULL DEFAULT '[]',
    legal_measure VARCHAR(255) NULL,
    has_untaken_medications BOOLEAN NOT NULL DEFAULT FALSE,
    -- Education
    education_currently_enrolled BOOLEAN NOT NULL DEFAULT FALSE,
    education_institution VARCHAR(255) NULL,
    education_mentor_name VARCHAR(255) NULL,
    education_mentor_phone VARCHAR(50) NULL,
    education_mentor_email VARCHAR(255) NULL,
    education_additional_notes TEXT NULL,
    education_level VARCHAR(50) NULL CHECK (education_level IN ('primary', 'secondary', 'higher', 'none')) DEFAULT 'none',
    -- Work
    work_currently_employed BOOLEAN NOT NULL DEFAULT FALSE,
    work_current_employer VARCHAR(255) NULL,
    work_current_employer_phone VARCHAR(50) NULL,
    work_current_employer_email VARCHAR(255) NULL,
    work_current_position VARCHAR(255) NULL,
    work_start_date DATE NULL,
    work_additional_notes TEXT NULL,
    -- Living situation
    living_situation VARCHAR(50) NULL CHECK (living_situation IN ('home', 'foster_care', 'youth_care_institution', 'other')),
    living_situation_notes TEXT NULL,

    -- Risks 
    risk_aggressive_behavior BOOLEAN DEFAULT FALSE,
    risk_suicidal_selfharm BOOLEAN DEFAULT FALSE,
    risk_substance_abuse BOOLEAN DEFAULT FALSE,
    risk_psychiatric_issues BOOLEAN DEFAULT FALSE,
    risk_criminal_history BOOLEAN DEFAULT FALSE,
    risk_flight_behavior BOOLEAN DEFAULT FALSE,
    risk_weapon_possession BOOLEAN DEFAULT FALSE,
    risk_sexual_behavior BOOLEAN DEFAULT FALSE,
    risk_day_night_rhythm BOOLEAN DEFAULT FALSE,
    risk_other BOOLEAN DEFAULT FALSE,
    risk_other_description TEXT,
    risk_additional_notes TEXT
);

CREATE INDEX client_details_sender_id_idx ON client_details(sender_id);
CREATE INDEX client_details_location_id_idx ON client_details(location_id);

-- Client status history tracking
CREATE TABLE client_status_history (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    old_status VARCHAR(50),
    new_status VARCHAR(50) NOT NULL,
    changed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    changed_by BIGINT NULL,
    reason VARCHAR(255)
);

CREATE INDEX idx_client_status_history_client_id ON client_status_history(client_id);
CREATE INDEX idx_client_status_history_changed_at ON client_status_history(changed_at DESC);

-- Scheduled status changes
CREATE TABLE scheduled_status_changes (
    id SERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    new_status VARCHAR(50) NULL,
    reason TEXT,
    scheduled_date DATE NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Client diagnoses
CREATE TABLE client_diagnosis (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(50) NULL,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    diagnosis_code VARCHAR(10) NOT NULL,
    description TEXT NOT NULL,
    severity VARCHAR(50) NULL,
    status VARCHAR(100) NOT NULL,
    diagnosing_clinician VARCHAR(100) NULL,
    notes TEXT NULL,
    created_at TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX client_diagnosis_client_id_idx ON client_diagnosis(client_id);
CREATE INDEX client_diagnosis_diagnosis_code_idx ON client_diagnosis(diagnosis_code);

-- Contact relationships
CREATE TABLE contact_relationship (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    soft_delete BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX contact_relationship_soft_delete_idx ON contact_relationship(soft_delete);

-- Client emergency contacts
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
    created_at TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    medical_reports BOOLEAN NOT NULL DEFAULT FALSE,
    incidents_reports BOOLEAN NOT NULL DEFAULT FALSE,
    goals_reports BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX client_emergency_contact_client_id_idx ON client_emergency_contact(client_id);

-- Client documents
CREATE TABLE client_documents (
    id BIGSERIAL PRIMARY KEY,
    attachment_uuid UUID NULL REFERENCES attachment_file("uuid") ON DELETE SET NULL,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    label VARCHAR(100) NOT NULL CHECK (label IN (
        'registration_form', 'intake_form', 'consent_form',
        'risk_assessment', 'self_reliance_matrix', 'force_inventory',
        'care_plan', 'signaling_plan', 'cooperation_agreement', 'other'
    )) DEFAULT 'other'
);

CREATE INDEX client_documents_user_id_idx ON client_documents(client_id);
CREATE INDEX client_documents_label_idx ON client_documents(label);

-- Client medications
CREATE TABLE client_medication (
    id BIGSERIAL PRIMARY KEY,
    diagnosis_id BIGINT NULL REFERENCES client_diagnosis(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    dosage VARCHAR(100) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NULL,
    notes TEXT NULL,
    self_administered BOOLEAN NOT NULL DEFAULT TRUE,
    slots JSONB NULL DEFAULT '[]',
    administered_by_id BIGINT NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    is_critical BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

-- ==========================================
-- CONTRACTS & FINANCIAL MANAGEMENT
-- ==========================================

-- Contract types
CREATE TABLE contract_type (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

-- Main contracts table
CREATE TABLE contract (
    id BIGSERIAL PRIMARY KEY,
    type_id BIGINT NULL REFERENCES contract_type(id) ON DELETE SET NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('approved', 'draft', 'terminated', 'stopped', 'expired')) DEFAULT 'draft',
    approved_at TIMESTAMPTZ NULL,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    reminder_period INTEGER NOT NULL DEFAULT 90,
    VAT INTEGER NULL DEFAULT -1,
    price DECIMAL(10,2) NOT NULL,
    price_time_unit VARCHAR(20) NOT NULL CHECK (price_time_unit IN ('minute', 'hourly', 'daily', 'weekly', 'monthly')) DEFAULT 'weekly',
    hours DECIMAL(10,2) NULL DEFAULT 0,
    hours_type VARCHAR(20) NULL DEFAULT NULL CHECK (hours_type IN ('weekly', 'all_period') OR hours_type IS NULL),
    care_name VARCHAR(255) NOT NULL,
    care_type VARCHAR(20) NOT NULL CHECK (care_type IN ('ambulante', 'accommodation')),
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE, 
    sender_id BIGINT NULL REFERENCES sender(id) ON DELETE SET NULL,
    attachment_ids UUID[] NOT NULL DEFAULT '{}',  
    financing_act VARCHAR(50) NOT NULL CHECK (financing_act IN ('WMO', 'ZVW', 'WLZ', 'JW', 'WPG')) DEFAULT 'WMO',
    financing_option VARCHAR(50) NOT NULL CHECK (financing_option IN ('ZIN', 'PGB')) DEFAULT 'PGB',
    departure_reason VARCHAR(255) NULL,
    departure_report TEXT NULL,
    updated_at TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX contract_type_id_idx ON contract(type_id);
CREATE INDEX contract_client_id_idx ON contract(client_id);
CREATE INDEX contract_sender_id_idx ON contract(sender_id);
CREATE INDEX contract_status_idx ON contract(status);

-- Contract audit table
CREATE TABLE contract_audit (
    audit_id BIGSERIAL PRIMARY KEY,
    contract_id BIGINT NOT NULL,
    operation VARCHAR(10) NOT NULL CHECK (operation IN ('INSERT', 'UPDATE', 'DELETE')),
    changed_by BIGINT NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    changed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    old_values JSONB NULL,
    new_values JSONB NULL,
    changed_fields TEXT[] NULL
);

CREATE INDEX idx_contract_audit_contract_id ON contract_audit(contract_id);
CREATE INDEX idx_contract_audit_changed_at ON contract_audit(changed_at);
CREATE INDEX idx_contract_audit_operation ON contract_audit(operation);

-- Contract audit trigger function
CREATE OR REPLACE FUNCTION contract_audit_trigger_func() 
RETURNS TRIGGER AS $$ 
DECLARE 
    old_row JSONB;
    new_row JSONB;
    changed_fields TEXT[] := '{}';
    field_name TEXT;
    current_user_id BIGINT;
BEGIN
    BEGIN
        current_user_id := current_setting('myapp.current_employee_id')::BIGINT;
    EXCEPTION
        WHEN OTHERS THEN
            current_user_id := NULL;
    END;

    IF TG_OP = 'DELETE' THEN
        old_row := to_jsonb(OLD);
        INSERT INTO contract_audit (contract_id, operation, old_values, changed_by, changed_at)
        VALUES (OLD.id, 'DELETE', old_row, current_user_id, CURRENT_TIMESTAMP);
        RETURN OLD;
        
    ELSIF TG_OP = 'INSERT' THEN
        new_row := to_jsonb(NEW);
        INSERT INTO contract_audit (contract_id, operation, new_values, changed_by, changed_at)
        VALUES (NEW.id, 'INSERT', new_row, current_user_id, CURRENT_TIMESTAMP);
        RETURN NEW;
        
    ELSIF TG_OP = 'UPDATE' THEN
        old_row := to_jsonb(OLD);
        new_row := to_jsonb(NEW);
        
        FOR field_name IN SELECT jsonb_object_keys(new_row) LOOP
            IF old_row->>field_name IS DISTINCT FROM new_row->>field_name THEN
                changed_fields := array_append(changed_fields, field_name);
            END IF;
        END LOOP;
        
        IF array_length(changed_fields, 1) > 0 AND 
           NOT (array_length(changed_fields, 1) = 1 AND 'updated_at' = ANY(changed_fields)) THEN
            INSERT INTO contract_audit (contract_id, operation, old_values, new_values, changed_fields, changed_by, changed_at)
            VALUES (NEW.id, 'UPDATE', old_row, new_row, changed_fields, current_user_id, CURRENT_TIMESTAMP);
        END IF;
        RETURN NEW;
    END IF;
    
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER contract_audit_trigger
    AFTER INSERT OR UPDATE OR DELETE ON contract
    FOR EACH ROW EXECUTE FUNCTION contract_audit_trigger_func();

-- Contract-related tables
CREATE TABLE contract_reminder (
    id BIGSERIAL PRIMARY KEY,
    contract_id BIGINT NOT NULL REFERENCES contract(id) ON DELETE CASCADE,
    reminder_sent_at TIMESTAMPTZ NULL,
    reminder_type VARCHAR(20) NOT NULL CHECK (reminder_type IN ('initial', 'follow_up', 'none')) DEFAULT 'none'
);

CREATE TABLE contract_working_hours (
    id BIGSERIAL PRIMARY KEY,
    contract_id BIGINT NOT NULL REFERENCES contract(id) ON DELETE CASCADE,
    minutes INTEGER NOT NULL DEFAULT 0,
    "datetime" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    notes TEXT NULL DEFAULT ''
);

CREATE INDEX contract_working_hours_contract_id_idx ON contract_working_hours(contract_id);
CREATE INDEX contract_working_hours_datetime_idx ON contract_working_hours(datetime);

CREATE TABLE contract_attachment (
    id BIGSERIAL PRIMARY KEY,
    contract_id BIGINT NOT NULL REFERENCES contract(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    attachment VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX contract_attachment_contract_id_idx ON contract_attachment(contract_id);

CREATE TABLE client_agreement (
    id BIGSERIAL PRIMARY KEY,
    contract_id BIGINT NOT NULL REFERENCES contract(id) ON DELETE CASCADE,
    agreement_details TEXT NOT NULL,
    created TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX client_agreement_contract_id_idx ON client_agreement(contract_id);

CREATE TABLE provision (
    id BIGSERIAL PRIMARY KEY,
    contract_id BIGINT NOT NULL REFERENCES contract(id) ON DELETE CASCADE,
    provision_details TEXT NOT NULL,
    created TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX provision_contract_id_idx ON provision(contract_id);

CREATE TABLE framework_agreement (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    agreement_details TEXT NOT NULL,
    created TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX framework_agreement_client_id_idx ON framework_agreement(client_id);

-- Invoice management
CREATE TABLE invoice (
    id BIGSERIAL PRIMARY KEY,
    invoice_number VARCHAR(50) NOT NULL UNIQUE,
    invoice_sequence BIGINT NOT NULL DEFAULT 1,
    issue_date DATE NOT NULL DEFAULT CURRENT_DATE,
    due_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN (
        'outstanding', 'partially_paid', 'paid', 'expired',
        'overpaid', 'imported', 'concept', 'canceled'
    )) DEFAULT 'concept',
    invoice_type VARCHAR(20) NOT NULL CHECK (invoice_type IN ('standard', 'credit_note')) DEFAULT 'standard',
    original_invoice_id BIGINT NULL REFERENCES invoice(id) ON DELETE SET NULL,
    invoice_details JSONB NULL DEFAULT '[]',
    total_amount DECIMAL(20,2) NOT NULL DEFAULT 0,
    pdf_attachment_id UUID NULL UNIQUE REFERENCES attachment_file("uuid") ON DELETE SET NULL,
    extra_content JSONB NULL DEFAULT '{}',
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    sender_id BIGINT NULL REFERENCES sender(id) ON DELETE SET NULL,
    warning_count INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX invoice_invoice_number_idx ON invoice(invoice_number);
CREATE INDEX invoice_client_id_idx ON invoice(client_id);
CREATE INDEX invoice_status_idx ON invoice(status);

-- Invoice audit table
CREATE TABLE invoice_audit (
    audit_id BIGSERIAL PRIMARY KEY,
    invoice_id BIGINT NOT NULL,
    operation VARCHAR(10) NOT NULL CHECK (operation IN ('INSERT', 'UPDATE', 'DELETE')),
    changed_by BIGINT REFERENCES employee_profile(id) ON DELETE SET NULL,
    changed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    old_values JSONB NULL,
    new_values JSONB NULL,
    changed_fields TEXT[] NULL
);

CREATE INDEX idx_invoice_audit_invoice_id ON invoice_audit(invoice_id);
CREATE INDEX idx_invoice_audit_changed_at ON invoice_audit(changed_at);
CREATE INDEX idx_invoice_audit_operation ON invoice_audit(operation);

-- Invoice audit trigger function
CREATE OR REPLACE FUNCTION invoice_audit_trigger_func()
RETURNS TRIGGER AS $$
DECLARE
    old_row JSONB;
    new_row JSONB;
    changed_fields_arr TEXT[] := '{}';
    field_name TEXT;
    current_employee_id BIGINT;
BEGIN
    BEGIN
        current_employee_id := current_setting('myapp.current_employee_id')::BIGINT;
    EXCEPTION
        WHEN OTHERS THEN
            current_employee_id := NULL;
    END;

    IF TG_OP = 'DELETE' THEN
        old_row := to_jsonb(OLD);
        INSERT INTO invoice_audit (invoice_id, operation, old_values, changed_by, changed_at)
        VALUES (OLD.id, 'DELETE', old_row, current_employee_id, CURRENT_TIMESTAMP);
        RETURN OLD;

    ELSIF TG_OP = 'INSERT' THEN
        new_row := to_jsonb(NEW);
        INSERT INTO invoice_audit (invoice_id, operation, new_values, changed_by, changed_at)
        VALUES (NEW.id, 'INSERT', new_row, current_employee_id, CURRENT_TIMESTAMP);
        RETURN NEW;

    ELSIF TG_OP = 'UPDATE' THEN
        old_row := to_jsonb(OLD);
        new_row := to_jsonb(NEW);
        
        FOR field_name IN SELECT jsonb_object_keys(new_row) LOOP
            IF old_row->>field_name IS DISTINCT FROM new_row->>field_name THEN
                changed_fields_arr := array_append(changed_fields_arr, field_name);
            END IF;
        END LOOP;
        
        IF array_length(changed_fields_arr, 1) > 0 AND NOT (array_length(changed_fields_arr, 1) = 1 AND 'updated_at' = ANY(changed_fields_arr)) THEN
            INSERT INTO invoice_audit (invoice_id, operation, old_values, new_values, changed_fields, changed_by, changed_at)
            VALUES (NEW.id, 'UPDATE', old_row, new_row, changed_fields_arr, current_employee_id, CURRENT_TIMESTAMP);
        END IF;
        
        RETURN NEW;
    END IF;
    
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER invoice_audit_trigger
AFTER INSERT OR UPDATE OR DELETE ON invoice
FOR EACH ROW EXECUTE FUNCTION invoice_audit_trigger_func();

-- Invoice payment history
CREATE TABLE invoice_payment_history (
    id BIGSERIAL PRIMARY KEY,
    invoice_id BIGINT NOT NULL REFERENCES invoice(id) ON DELETE CASCADE,
    payment_method VARCHAR(20) NULL CHECK (payment_method IN (
        'bank_transfer', 'credit_card', 'check', 'cash', 'other'
    )),
    payment_status VARCHAR(20) NOT NULL CHECK (payment_status IN (
        'completed', 'pending', 'failed', 'reversed', 'refunded'
    )) DEFAULT 'completed',
    amount DECIMAL(20,2) NOT NULL DEFAULT 0,
    payment_date DATE NOT NULL DEFAULT CURRENT_DATE,
    payment_reference VARCHAR(100) NULL,
    notes TEXT NULL,
    recorded_by BIGINT NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_invoice_payment_history_invoice_id ON invoice_payment_history(invoice_id);
CREATE INDEX idx_invoice_payment_history_payment_date ON invoice_payment_history(payment_date);
CREATE INDEX idx_invoice_payment_history_payment_status ON invoice_payment_history(payment_status);

-- Invoice-Contract relationship
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

-- ==========================================
-- CARE PLANS & ASSESSMENTS
-- ==========================================

-- Maturity matrix topics and levels
CREATE TABLE maturity_matrix (
    id BIGINT PRIMARY KEY,
    topic_name VARCHAR(255) NOT NULL,
    level_description JSONB NOT NULL DEFAULT '[]'
);

-- Insert maturity matrix data
INSERT INTO maturity_matrix (id, topic_name, level_description)
VALUES
    (1, 'Finances', '[
  {
    "level": 1,
    "name": "Acute problematiek",
    "description": "Groeiende complexe schulden"
  },
  {
    "level": 2,
    "name": "Niet zelfredzaam",
    "description": "Beschikt niet over vrij besteedbaar inkomen of groeiende schulden door spontaan of ongepast uitgeven"
  },
  {
    "level": 3,
    "name": "Beperkt zelfredzaam",
    "description": "Beschikt over vrij besteedbaar inkomen van ouders zonder verantwoordelijkheid voor noodzakelijke behoeften (zakgeld). Eventuele schulden zijn stabiel en onder beheer"
  },
  {
    "level": 4,
    "name": "Voldoende zelfredzaam",
    "description": "Beschikt over vrij besteedbaar inkomen van ouders met enige verantwoordelijkheid voor noodzakelijke behoeften (zakgeld, en kleed-/lunchgeld). Gepast uitgeven. Eventuele schulden verminderen"
  },
  {
    "level": 5,
    "name": "Volledig zelfredzaam",
    "description": "Beschikt over vrij besteedbaar inkomen (uit klusjes of (bij)baan) met enige verantwoordelijkheid voor noodzakelijke behoeften. Aan het eind van de maand is geld over. Geen schulden"
  }
]'),
    (2, 'Work & Education', '[
  {
    "level": 1,
    "name": "Geen opleiding of werk",
    "description": "Geen (traject naar) opleiding/werk of werk zonder adequate toerusting/verzekering. Geen zoekactiviteiten naar opleiding/werk."
  },
  {
    "level": 2,
    "name": "Zoekende maar instabiel",
    "description": "Geen (traject naar) opleiding/werk, maar wel zoekactiviteiten gericht op opleiding/werk of \"papieren\" opleiding (ingeschreven maar niet volgend) of veel schoolverzuim/dreigend ontslag of dreigende drop-out."
  },
  {
    "level": 3,
    "name": "Instabiele opleiding of werk",
    "description": "Volgt opleiding maar loopt achter of heeft geregeld verzuim van opleiding/werk of volgt traject naar opleiding (trajectbegeleiding, coaching voor schoolverlaters)."
  },
  {
    "level": 4,
    "name": "Op schema",
    "description": "Op schema met opleiding of heeft startkwalificatie met tijdelijke baan/traject naar opleiding/traject naar werk. Zelden ongeoorloofd verzuim."
  },
  {
    "level": 5,
    "name": "Succesvol in opleiding of werk",
    "description": "Presteert zeer goed op opleiding of heeft startkwalificatie met vaste baan. Geen ongeoorloofd verzuim."
  }
]'),
    (3, 'Use of Time', '[
  {
    "level": 1,
    "name": "Geen structuur of activiteiten",
    "description": "Afwezigheid van activiteiten die plezierig/nuttig zijn. Geen structuur in de dag. Onregelmatig dag-nacht ritme."
  },
  {
    "level": 2,
    "name": "Zeer beperkte activiteiten en structuur",
    "description": "Nauwelijks activiteiten die plezierig/nuttig zijn. Nauwelijks structuur in de dag. Afwijkend dag-nacht ritme."
  },
  {
    "level": 3,
    "name": "Onvoldoende maar acceptabel",
    "description": "Onvoldoende activiteiten die plezierig/nuttig zijn, maar voldoende structuur in de dag. Enige afwijkingen in het dag-nacht ritme."
  },
  {
    "level": 4,
    "name": "Voldoende activiteiten en structuur",
    "description": "Voldoende activiteiten die plezierig/nuttig zijn. Dag-nacht ritme heeft geen negatieve invloed op het dagelijks functioneren."
  },
  {
    "level": 5,
    "name": "Gezonde balans en structuur",
    "description": "Tijd is overwegend gevuld met plezierige/nuttige activiteiten. Gezond dag-nacht ritme."
  }
]'),
    (4, 'Housing', '[
    {
        "level": 1,
        "name": "Acute problematiek",
        "description": "dakloos of in crisisopvang"
    },
    {
        "level": 2,
        "name": "Niet zelfredzaam",
        "description": "voor wonen ongeschikte huisvesting of dreigende huisuitzetting"
    },
    {
        "level": 3,
        "name": "Beperkt zelfredzaam",
        "description": "veilige, stabiele huisvesting maar slechts marginaal toereikend of verblijft in niet-autonome huisvesting (instelling)"
    },
    {
        "level": 4,
        "name": "Voldoende zelfredzaam",
        "description": "veilige, stabiele en toereikende huisvesting, gedeeltelijk autonome huisvesting (begeleid wonen)"
    },
    {
        "level": 5,
        "name": "Volledig zelfredzaam",
        "description": "veilige, stabiele en toereikende huisvesting, autonome huisvesting (zelfstandig wonen), woont bij ouders/verzorgers"
    }
]'),
    (5, 'Domestic Relationships', '[
    {
        "level": 1,
        "name": "Acute problematiek",
        "description": "geweld in huiselijke kring/ kindermishandeling/ misbruik/ verwaarlozing"
    },
    {
        "level": 2,
        "name": "Niet zelfredzaam",
        "description": "relationele problemen met leden van het huishouden of dreigend geweld in huiselijke kring/ kindermishandeling/ misbruik/ verwaarlozing"
    },
    {
        "level": 3,
        "name": "Beperkt zelfredzaam",
        "description": "spanningen in relatie(s) met leden van het huishouden, probeert eigen negatief relationeel gedrag te veranderen"
    },
    {
        "level": 4,
        "name": "Voldoende zelfredzaam",
        "description": "relationele problemen met leden van het huishouden of spanningen tussen leden van het huishouden zijn niet (meer) aanwezig"
    },
    {
        "level": 5,
        "name": "Volledig zelfredzaam",
        "description": "wordt gesteund en steunt binnen het huishouden, communicatie met leden van het huishouden is consistent open"
    }
]'),
    (6, 'Mental Health', '[
    {
        "level": 1,
        "name": "Acute problematiek",
        "description": "geestelijke noodsituatie, een gevaar voor zichzelf/anderen"
    },
    {
        "level": 2,
        "name": "Niet zelfredzaam",
        "description": "(chronische) geestelijke aandoening maar geen gevaar voor zichzelf/anderen, functioneren is ernstig beperkt door geestelijk gezondheidsprobleem (incl. gedrags-ontwikkelingsproblematiek), geen behandeling"
    },
    {
        "level": 3,
        "name": "Beperkt zelfredzaam",
        "description": "geestelijke aandoening, functioneren is beperkt door geestelijk gezondheidsprobleem (incl. gedrags- en ontwikkelingsproblematiek), behandeltrouw is minimaal of beperking bestaat ondanks goede behandeltrouw"
    },
    {
        "level": 4,
        "name": "Voldoende zelfredzaam",
        "description": "minimale tekenen van geestelijke onrust die voorspelbare reactie zijn op stressoren in het leven (ook puberteit), functioneren is marginaal beperkt door geestelijke onrust, goede behandeltrouw of geen behandeling nodig"
    },
    {
        "level": 5,
        "name": "Volledig zelfredzaam",
        "description": "geestelijk gezond, niet meer dan de dagelijkse beslommeringen/zorgen"
    }
]'),
    (7, 'Physical Health', '[
    {
        "level": 1,
        "name": "Acute problematiek",
        "description": "een noodgeval/ kritieke situatie, direct medische aandacht nodig"
    },
    {
        "level": 2,
        "name": "Niet zelfredzaam",
        "description": "(chronische) lichamelijke aandoening die medische behandeling vereist, functioneren is ernstig beperkt door lichamelijk gezondheidsprobleem, geen behandeling"
    },
    {
        "level": 3,
        "name": "Beperkt zelfredzaam",
        "description": "lichamelijke aandoening, functioneren is beperkt door lichamelijk gezondheidsprobleem, behandeltrouw is minimaal of beperking bestaat ondanks goede behandeltrouw"
    },
    {
        "level": 4,
        "name": "Voldoende zelfredzaam",
        "description": "minimaal lichamelijk ongemak dat samenhangt met dagelijkse activiteiten, functioneren is marginaal beperkt door lichamelijk ongemak, goede behandeltrouw of geen behandeling nodig"
    },
    {
        "level": 5,
        "name": "Volledig zelfredzaam",
        "description": "lichamelijk gezond, gezonde leefstijl (gezonde voeding en voldoende bewegen)"
    }
]'),
    (8, 'Substance Use', '[
    {
        "level": 1,
        "name": "Acute problematiek",
        "description": "(gedrags-) stoornis/afhankelijkheid van het gebruik van middelen of van games/gokken/seks/internet, gebruik veroorzaakt/verergert lichamelijke/geestelijke problemen die behandeling vereisen"
    },
    {
        "level": 2,
        "name": "Niet zelfredzaam",
        "description": "gebruik van middelen of problematisch \"gebruik\" van games/gokken/seks/internet, aan gebruik gerelateerde lichamelijke/geestelijke problemen of problemen thuis/op school/op het werk, geen behandeling"
    },
    {
        "level": 3,
        "name": "Beperkt zelfredzaam",
        "description": "gebruik van middelen, geen aan middelengebruik gerelateerde problemen, behandeltrouw is minimaal of beperking bestaat ondanks goede behandeltrouw"
    },
    {
        "level": 4,
        "name": "Voldoende zelfredzaam",
        "description": "geen middelengebruik ondanks sterke drang of behandeling met potentieel verslavende middelen zonder bijgebruik, goede behandeltrouw of geen behandeling nodig"
    },
    {
        "level": 5,
        "name": "Volledig zelfredzaam",
        "description": "geen middelengebruik, geen sterke drang naar gebruik van middelen"
    }
]'),
    (9, 'Basic ADL', '[
    {
        "level": 1,
        "name": "Acute problematiek",
        "description": "een gebied van de basale ADL wordt niet uitgevoerd, verhongering of uitdroging of bevulling/vervulling"
    },
    {
        "level": 2,
        "name": "Niet zelfredzaam",
        "description": "meerdere gebieden van de basale ADL worden beperkt uitgevoerd"
    },
    {
        "level": 3,
        "name": "Beperkt zelfredzaam",
        "description": "alle gebieden van de basale ADL worden uitgevoerd maar een enkel gebied van de basale ADL wordt beperkt uitgevoerd"
    },
    {
        "level": 4,
        "name": "Voldoende zelfredzaam",
        "description": "geen beperkingen in de uitvoering van de basale ADL, krijgt hulp of gebruikt hulpmiddel"
    },
    {
        "level": 5,
        "name": "Volledig zelfredzaam",
        "description": "geen beperkingen in de uitvoering van de basale ADL, zoals eten, wassen en aankleden, geen gebruik van hulpmiddelen"
    }
]'),
    (10, 'Instrumental ADL', '[
    {
        "level": 1,
        "name": "Acute problematiek",
        "description": "meerdere gebieden van de instrumentele ADL worden niet uitgevoerd, woningvervulling of onder-/over-medicatie of geen administratie of voedselvergiftiging"
    },
    {
        "level": 2,
        "name": "Niet zelfredzaam",
        "description": "een enkel gebied van de instrumentele ADL wordt niet uitgevoerd of uitvoering op meerdere gebieden is beperkt, weet gezien de leeftijd te weinig van welke instanties er zijn, wat je er mee moet doen en hoe ze te benaderen"
    },
    {
        "level": 3,
        "name": "Beperkt zelfredzaam",
        "description": "alle gebieden van de instrumentele ADL worden uitgevoerd, uitvoering van een enkel gebied van de instrumentele ADL is beperkt, weet beperkt van instanties af en krijgt gezien de leeftijd veel hulp bij het contact met instanties"
    },
    {
        "level": 4,
        "name": "Voldoende zelfredzaam",
        "description": "geen beperkingen in de uitvoering van de instrumentele ADL, krijgt hulp van buiten het huishouden of gebruikt hulpmiddel, weet van instanties af, maar krijgt gezien de leeftijd enige hulp bij het contact leggen met en het gebruik maken van instanties"
    },
    {
        "level": 5,
        "name": "Volledig zelfredzaam",
        "description": "geen beperkingen in de uitvoering van de instrumentele ADL, krijgt geen hulp van buiten huishouden en maakt geen gebruik van hulpmiddelen, maakt leeftijdsadequaat gebruik van instanties"
    }
]'),
    (11, 'Social Network', '[
    {
        "level": 1,
        "name": "Acute problematiek",
        "description": "ernstig sociaal isolement, geen steunend contact met familie of met volwassen steunfiguur buiten gezin, geen steunend contact met leeftijdgenoten"
    },
    {
        "level": 2,
        "name": "Niet zelfredzaam",
        "description": "geen steunend contact met familie of met volwassen steunfiguur buiten gezin, weinig steunend contact met leeftijdgenoten, veel belemmerend contact"
    },
    {
        "level": 3,
        "name": "Beperkt zelfredzaam",
        "description": "enig steunend contact met familie of met één volwassen steunfiguur buiten het huishouden, enig steunend contact met leeftijdgenoten, weinig belemmerend contact"
    },
    {
        "level": 4,
        "name": "Voldoende zelfredzaam",
        "description": "voldoende steunend contact met familie of met volwassen steunfiguren buiten het huishouden, voldoende steunend contact met leeftijdgenoten, nauwelijks belemmerend contact"
    },
    {
        "level": 5,
        "name": "Volledig zelfredzaam",
        "description": "gezond sociaal netwerk, veel steunend contact met familie of met volwassen steunfiguur buiten het huishouden, veel steunend contact met leeftijdgenoten, geen belemmerend contact"
    }
]'),
    (12, 'Social Participation', '[
    {
        "level": 1,
        "name": "Acute problematiek",
        "description": "niet van toepassing door crisissituatie of in overlevingsmodus of veroorzaakt ernstige overlast"
    },
    {
        "level": 2,
        "name": "Niet zelfredzaam",
        "description": "geen maatschappelijke participatie of veroorzaakt overlast"
    },
    {
        "level": 3,
        "name": "Beperkt zelfredzaam",
        "description": "nauwelijks participerend in maatschappij, logistieke, financiële of sociaal-maatschappelijke hindernissen om meer te participeren"
    },
    {
        "level": 4,
        "name": "Voldoende zelfredzaam",
        "description": "enige maatschappelijke participatie (meedoen), persoonlijke hindernis (motivatie) om meer te participeren"
    },
    {
        "level": 5,
        "name": "Volledig zelfredzaam",
        "description": "actief participerend in de maatschappij (bijdragen)"
    }
]'),
    (13, 'Justice', '[
    {
        "level": 1,
        "name": "Acute problematiek",
        "description": "zeer regelmatig (maandelijks) contact met politie of openstaande zaken bij justitie"
    },
    {
        "level": 2,
        "name": "Niet zelfredzaam",
        "description": "regelmatig (meerdere keren per jaar) contact met politie of lopende zaken bij justitie"
    },
    {
        "level": 3,
        "name": "Beperkt zelfredzaam",
        "description": "incidenteel (eens per jaar) contact met politie of voorwaardelijke straf/voorwaardelijke invrijheidsstelling"
    },
    {
        "level": 4,
        "name": "Voldoende zelfredzaam",
        "description": "zelden (minder dan eens per jaar) contact met politie of strafblad"
    },
    {
        "level": 5,
        "name": "Volledig zelfredzaam",
        "description": "geen contact met politie, geen strafblad"
    }
]');

-- Client maturity matrix assessments
CREATE TABLE client_maturity_matrix_assessment (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    maturity_matrix_id BIGINT NOT NULL REFERENCES maturity_matrix(id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    initial_level INT NOT NULL CHECK (initial_level BETWEEN 1 AND 5),
    target_level INT NOT NULL CHECK (target_level BETWEEN 1 AND 5),
    current_level INT NOT NULL CHECK (current_level BETWEEN 1 AND 5),
    care_plan_generated_at TIMESTAMPTZ NULL DEFAULT NULL,
    care_plan_status VARCHAR(20) NOT NULL CHECK (care_plan_status IN ('pending', 'generated', 'approved', 'active', 'completed', 'discontinued')) DEFAULT 'pending',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE(client_id, maturity_matrix_id)
);

-- Level change history
CREATE TABLE level_history (
    id BIGSERIAL PRIMARY KEY,
    client_maturity_matrix_assessment_id BIGINT NOT NULL REFERENCES client_maturity_matrix_assessment(id) ON DELETE CASCADE,
    change_date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    old_level INT NOT NULL CHECK (old_level BETWEEN 1 AND 5),
    new_level INT NOT NULL CHECK (new_level BETWEEN 1 AND 5),
    comment TEXT NOT NULL
);

-- Trigger function for level history
CREATE OR REPLACE FUNCTION trg_after_update_client_maturity_matrix_assessment_func()
RETURNS trigger AS $$
BEGIN
    IF NEW.current_level <> OLD.current_level THEN
        INSERT INTO level_history (client_maturity_matrix_assessment_id, old_level, new_level, comment)
        VALUES (OLD.client_maturity_matrix_assessment_id, OLD.current_level, NEW.current_level, 'Automatic logging of level change.');
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_after_update_client_maturity_matrix_assessment
AFTER UPDATE ON client_maturity_matrix_assessment
FOR EACH ROW
EXECUTE FUNCTION trg_after_update_client_maturity_matrix_assessment_func();

-- Care plans
CREATE TABLE care_plans (
    id BIGSERIAL PRIMARY KEY,
    assessment_id BIGINT NOT NULL REFERENCES client_maturity_matrix_assessment(id) ON DELETE CASCADE,
    generated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    generated_by_employee_id BIGINT REFERENCES employee_profile(id),
    approved_by_employee_id BIGINT REFERENCES employee_profile(id),
    approved_at TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'pending_approval', 'approved', 'active', 'completed', 'discontinued')),
    assessment_summary TEXT NOT NULL,
    raw_llm_response JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version INT NOT NULL DEFAULT 1,
    UNIQUE(assessment_id)
);

-- Care plan objectives
CREATE TABLE care_plan_objectives (
    id BIGSERIAL PRIMARY KEY,
    care_plan_id BIGINT NOT NULL REFERENCES care_plans(id) ON DELETE CASCADE,
    timeframe VARCHAR(20) NOT NULL CHECK (timeframe IN ('short_term', 'medium_term', 'long_term')),
    goal_title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    target_date DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'not_started' CHECK (status IN ('not_started', 'in_progress', 'completed', 'discontinued', 'draft')),
    completion_date DATE,
    completion_notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Care plan actions
CREATE TABLE care_plan_actions (
    id BIGSERIAL PRIMARY KEY,
    objective_id BIGINT NOT NULL REFERENCES care_plan_objectives(id) ON DELETE CASCADE,
    action_description TEXT NOT NULL,
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    completed_at TIMESTAMP,
    completed_by_employee_id BIGINT REFERENCES employee_profile(id),
    notes TEXT,
    sort_order INT NOT NULL DEFAULT 0
);

-- Care plan interventions
CREATE TABLE care_plan_interventions (
    id BIGSERIAL PRIMARY KEY,
    care_plan_id BIGINT NOT NULL REFERENCES care_plans(id) ON DELETE CASCADE,
    frequency VARCHAR(20) NOT NULL CHECK (frequency IN ('daily', 'weekly', 'monthly')),
    intervention_description TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_completed_date DATE,
    total_completions INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP 
);

-- Care plan metrics
CREATE TABLE care_plan_metrics (
    id BIGSERIAL PRIMARY KEY,
    care_plan_id BIGINT NOT NULL REFERENCES care_plans(id) ON DELETE CASCADE,
    metric_name VARCHAR(255) NOT NULL,
    target_value VARCHAR(255) NOT NULL,
    measurement_method TEXT NOT NULL,
    current_value VARCHAR(255),
    last_measured_date DATE,
    is_achieved BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Care plan risks
CREATE TABLE care_plan_risks (
    id BIGSERIAL PRIMARY KEY,
    care_plan_id BIGINT NOT NULL REFERENCES care_plans(id) ON DELETE CASCADE,
    risk_description TEXT NOT NULL,
    mitigation_strategy TEXT NOT NULL,
    risk_level VARCHAR(20) CHECK (risk_level IN ('low', 'medium', 'high')),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Care plan support network
CREATE TABLE care_plan_support_network (
    id BIGSERIAL PRIMARY KEY,
    care_plan_id BIGINT NOT NULL REFERENCES care_plans(id) ON DELETE CASCADE,
    role_title VARCHAR(255) NOT NULL,
    responsibility_description TEXT NOT NULL,
    contact_person VARCHAR(255),
    contact_details TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Care plan resources
CREATE TABLE care_plan_resources (
    id BIGSERIAL PRIMARY KEY,
    care_plan_id BIGINT NOT NULL REFERENCES care_plans(id) ON DELETE CASCADE,
    resource_description TEXT NOT NULL,
    is_obtained BOOLEAN NOT NULL DEFAULT FALSE,
    obtained_date DATE,
    cost_estimate DECIMAL(10,2),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Care plan reports
CREATE TABLE care_plan_reports (
    id BIGSERIAL PRIMARY KEY,
    care_plan_id BIGINT NOT NULL REFERENCES care_plans(id) ON DELETE CASCADE,
    report_type VARCHAR(50) NOT NULL CHECK (report_type IN ('progress', 'concern', 'achievement', 'modification')),
    report_content TEXT NOT NULL,
    created_by_employee_id BIGINT NOT NULL REFERENCES employee_profile(id),
    is_critical BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ==========================================
-- INCIDENTS & REPORTING
-- ==========================================

-- Incident reports
CREATE TABLE incident (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    location_id BIGINT NOT NULL REFERENCES location(id) ON DELETE CASCADE,
    reporter_involvement VARCHAR(100) NOT NULL CHECK (reporter_involvement IN (
        'directly_involved', 'witness', 'found_afterwards', 'alarmed'
    )),
    inform_who VARCHAR(255)[] NOT NULL DEFAULT '{}',
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
    technical VARCHAR(255)[] NOT NULL DEFAULT '{}',
    organizational VARCHAR(255)[] NOT NULL DEFAULT '{}',
    mese_worker VARCHAR(255)[] NOT NULL DEFAULT '{}',
    client_options VARCHAR(255)[] NOT NULL DEFAULT '{}',
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
    succession VARCHAR(255)[] NOT NULL DEFAULT '{}',
    succession_desc TEXT NULL DEFAULT '',
    other BOOLEAN NOT NULL DEFAULT FALSE,
    other_desc VARCHAR(100) NULL,
    additional_appointments TEXT NULL DEFAULT '',
    employee_absenteeism VARCHAR(100) NOT NULL DEFAULT '',
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    soft_delete BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    file_url VARCHAR(255) NULL,
    emails TEXT[] NULL DEFAULT '{}'
);

CREATE INDEX incident_client_id_idx ON incident(client_id);
CREATE INDEX incident_location_id_idx ON incident(location_id);
CREATE INDEX incident_soft_delete_idx ON incident(soft_delete);

-- Employee assignments
CREATE TABLE assignment (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    start_datetime TIMESTAMPTZ NOT NULL,
    end_datetime TIMESTAMPTZ NOT NULL,
    "status" VARCHAR(50) NOT NULL CHECK (status IN ('Confirmed', 'Pending', 'Cancelled')),
    created_at TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX assignment_employee_id_idx ON assignment(employee_id);
CREATE INDEX assignment_client_id_idx ON assignment(client_id);

-- Assigned employees to clients
CREATE TABLE assigned_employee (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    employee_id BIGINT NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    role VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX assigned_employee_client_id_idx ON assigned_employee(client_id);
CREATE INDEX assigned_employee_employee_id_idx ON assigned_employee(employee_id);

-- Progress reports
CREATE TABLE progress_report (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    date TIMESTAMPTZ NOT NULL,
    title VARCHAR(50) NULL,
    report_text TEXT NOT NULL,
    employee_id BIGINT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL CHECK (type IN (
        'morning_report', 'evening_report', 'night_report', 'shift_report',
        'one_to_one_report', 'process_report', 'contact_journal', 'other'
    )) DEFAULT 'other',
    emotional_state VARCHAR(20) NOT NULL CHECK (emotional_state IN (
        'normal', 'excited', 'happy', 'sad', 'angry', 'anxious', 'depressed'
    )) DEFAULT 'normal',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX progress_report_client_id_idx ON progress_report(client_id);
CREATE INDEX progress_report_author_id_idx ON progress_report(employee_id);
CREATE INDEX progress_report_created_idx ON progress_report(created_at DESC);

-- AI generated reports
CREATE TABLE ai_generated_reports (
    id BIGSERIAL PRIMARY KEY,
    report_text TEXT NOT NULL,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ==========================================
-- SCHEDULING & APPOINTMENTS
-- ==========================================

-- Employee schedules
CREATE TABLE schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id BIGINT NOT NULL REFERENCES employee_profile(id),
    color VARCHAR(20) DEFAULT '#0000FF',
    location_id BIGINT NOT NULL REFERENCES location(id),
    location_shift_id BIGINT NULL REFERENCES location_shift(id),
    is_custom BOOLEAN NOT NULL DEFAULT FALSE,
    start_datetime TIMESTAMP NOT NULL,
    end_datetime TIMESTAMP NOT NULL,
    created_by_employee_id BIGINT NOT NULL REFERENCES employee_profile(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_timeframe CHECK (end_datetime > start_datetime)
);

-- Appointment templates
CREATE TABLE appointment_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    creator_employee_id BIGINT NOT NULL REFERENCES employee_profile(id),
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    location VARCHAR(255),
    description TEXT,
    color VARCHAR(20) DEFAULT '#0000FF', 
    recurrence_type VARCHAR(50) DEFAULT 'DAILY' CHECK (recurrence_type IN ('DAILY', 'WEEKLY', 'MONTHLY')),
    recurrence_interval INT NULL,
    recurrence_end_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Scheduled appointments
CREATE TABLE scheduled_appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_templates_id UUID NULL REFERENCES appointment_templates(id) ON DELETE CASCADE,
    creator_employee_id BIGINT NULL REFERENCES employee_profile(id),
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    location VARCHAR(255),
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'CONFIRMED', 'CANCELLED')),
    color VARCHAR(20) DEFAULT '#0000FF',
    is_confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    confirmed_by_employee_id BIGINT REFERENCES employee_profile(id),
    confirmed_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_scheduled_appointments_time_range ON scheduled_appointments (start_time, end_time);
CREATE INDEX idx_scheduled_appointments_template_id ON scheduled_appointments (appointment_templates_id);

-- Appointment participants
CREATE TABLE appointment_participants (
    appointment_participant_id BIGSERIAL PRIMARY KEY,
    appointment_id UUID NOT NULL REFERENCES scheduled_appointments(id) ON DELETE CASCADE,
    employee_id BIGINT NOT NULL REFERENCES employee_profile(id),
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (appointment_id, employee_id)
);

-- Appointment clients
CREATE TABLE appointment_clients (
    appointment_client_id BIGSERIAL PRIMARY KEY,
    appointment_id UUID NOT NULL REFERENCES scheduled_appointments(id) ON DELETE CASCADE,
    client_id BIGINT NOT NULL REFERENCES client_details(id),
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (appointment_id, client_id)
);

-- ==========================================
-- FORMS & DOCUMENTATION
-- ==========================================

-- Appointment cards
CREATE TABLE appointment_card (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES client_details(id) ON DELETE CASCADE UNIQUE,
    general_information TEXT[] NOT NULL DEFAULT '{}',
    important_contacts TEXT[] NOT NULL DEFAULT '{}',
    household_info TEXT[] NOT NULL DEFAULT '{}',
    organization_agreements TEXT[] NOT NULL DEFAULT '{}',
    youth_officer_agreements TEXT[] NOT NULL DEFAULT '{}',
    treatment_agreements TEXT[] NOT NULL DEFAULT '{}',
    smoking_rules TEXT[] NOT NULL DEFAULT '{}',
    work TEXT[] NOT NULL DEFAULT '{}',
    school_internship TEXT[] NOT NULL DEFAULT '{}',
    travel TEXT[] NOT NULL DEFAULT '{}',
    leave TEXT[] NOT NULL DEFAULT '{}',
    file_url VARCHAR(255) NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Registration forms
CREATE TABLE registration_form (
    id BIGSERIAL PRIMARY KEY,
    client_first_name VARCHAR(255) NOT NULL,
    client_last_name VARCHAR(255) NOT NULL,
    client_bsn_number VARCHAR(50) NOT NULL,
    client_gender VARCHAR(10) NOT NULL CHECK (client_gender IN ('male', 'female', 'other')),
    client_nationality VARCHAR(100) NOT NULL,
    client_phone_number VARCHAR(20) NOT NULL,
    client_email VARCHAR(255) NOT NULL,
    -- Client address
    client_street VARCHAR(255) NOT NULL,
    client_house_number VARCHAR(20) NOT NULL,
    client_postal_code VARCHAR(20) NOT NULL,
    client_city VARCHAR(100) NOT NULL,
    -- Referrer details
    referrer_first_name VARCHAR(255) NOT NULL,
    referrer_last_name VARCHAR(255) NOT NULL,
    referrer_organization VARCHAR(255) NOT NULL,
    referrer_job_title VARCHAR(255) NOT NULL,
    referrer_phone_number VARCHAR(20) NOT NULL,
    referrer_email VARCHAR(255) NOT NULL,
    -- Guardian details
    guardian1_first_name VARCHAR(255) NOT NULL,
    guardian1_last_name VARCHAR(255) NOT NULL,
    guardian1_relationship VARCHAR(100) NOT NULL,
    guardian1_phone_number VARCHAR(20) NOT NULL,
    guardian1_email VARCHAR(255) NOT NULL,
    guardian2_first_name VARCHAR(255) NOT NULL,
    guardian2_last_name VARCHAR(255) NOT NULL,
    guardian2_relationship VARCHAR(100) NOT NULL,
    guardian2_phone_number VARCHAR(20) NOT NULL,
    guardian2_email VARCHAR(255) NOT NULL,
    -- Education
    education_institution VARCHAR(255) NULL,
    education_mentor_name VARCHAR(255) NULL,
    education_mentor_phone VARCHAR(20) NULL,
    education_mentor_email VARCHAR(255) NULL,
    education_currently_enrolled BOOLEAN NOT NULL DEFAULT FALSE,
    education_additional_notes TEXT NULL,
    education_level VARCHAR(50) NULL CHECK (education_level IN ('primary', 'secondary', 'higher', 'none')) DEFAULT 'none',
    -- Work
    work_current_employer VARCHAR(255) NULL,
    work_employer_phone VARCHAR(20) NULL,
    work_employer_email VARCHAR(255) NULL,
    work_current_position VARCHAR(255) NULL,
    work_currently_employed BOOLEAN NOT NULL DEFAULT FALSE,
    work_start_date DATE NULL,
    work_additional_notes TEXT NULL,
    -- Care type
    care_protected_living BOOLEAN DEFAULT FALSE,
    care_assisted_independent_living BOOLEAN DEFAULT FALSE,
    care_room_training_center BOOLEAN DEFAULT FALSE,
    care_ambulatory_guidance BOOLEAN DEFAULT FALSE,
    -- Additional information
    application_reason TEXT,
    client_goals TEXT,
    -- Risks
    risk_aggressive_behavior BOOLEAN DEFAULT FALSE,
    risk_suicidal_selfharm BOOLEAN DEFAULT FALSE,
    risk_substance_abuse BOOLEAN DEFAULT FALSE,
    risk_psychiatric_issues BOOLEAN DEFAULT FALSE,
    risk_criminal_history BOOLEAN DEFAULT FALSE,
    risk_flight_behavior BOOLEAN DEFAULT FALSE,
    risk_weapon_possession BOOLEAN DEFAULT FALSE,
    risk_sexual_behavior BOOLEAN DEFAULT FALSE,
    risk_day_night_rhythm BOOLEAN DEFAULT FALSE,
    risk_other BOOLEAN DEFAULT FALSE,
    risk_other_description TEXT,
    risk_additional_notes TEXT,
    -- Document attachments
    document_referral UUID NULL REFERENCES attachment_file(uuid) ON DELETE SET NULL,
    document_education_report UUID NULL REFERENCES attachment_file(uuid) ON DELETE SET NULL,
    document_action_plan UUID NULL REFERENCES attachment_file(uuid) ON DELETE SET NULL,
    document_psychiatric_report UUID NULL REFERENCES attachment_file(uuid) ON DELETE SET NULL,
    document_diagnosis UUID NULL REFERENCES attachment_file(uuid) ON DELETE SET NULL,
    document_safety_plan UUID NULL REFERENCES attachment_file(uuid) ON DELETE SET NULL,
    document_id_copy UUID NULL REFERENCES attachment_file(uuid) ON DELETE SET NULL,
    -- Signatures and processing
    application_date DATE,
    referrer_signature BOOLEAN DEFAULT FALSE,
    form_status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (form_status IN ('pending', 'approved', 'rejected')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    submitted_at TIMESTAMPTZ NULL,
    processed_at TIMESTAMPTZ NULL,
    processed_by_employee_id BIGINT NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    intake_appointment_datetime TIMESTAMPTZ NULL,
    intake_appointment_location VARCHAR(255) NULL,
    addmission_type VARCHAR(50) NULL CHECK (addmission_type IN ('crisis_admission', 'regular_placement'))
);

-- Collaboration agreements
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

-- Risk assessments
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

-- Consent declarations
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

-- Youth care intake forms
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

-- Data sharing statements
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

-- ==========================================
-- UTILITY TABLES
-- ==========================================

-- Template items for document generation for custom data fields to include in invoice generation 
CREATE TABLE template_items (
    id BIGSERIAL PRIMARY KEY,
    item_tag VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    source_table VARCHAR(64) NOT NULL,
    source_column VARCHAR(64) NOT NULL
);

INSERT INTO template_items (item_tag, description, source_table, source_column) VALUES
('client.date_of_birth', 'Date of birth', 'client_details', 'date_of_birth'),
('client.filenumber', 'File number', 'client_details', 'filenumber'),
('contract.financing_act', 'Financing act', 'contract', 'financing_act'),
('contract.financing_option', 'Financing option', 'contract', 'financing_option');