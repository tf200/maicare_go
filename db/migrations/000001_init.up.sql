-- ==========================================
-- INFRASTRUCTURE & ORGANIZATIONS
-- ==========================================

-- Maturity matrix topics and levels
CREATE TABLE maturity_matrix (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    topic_name VARCHAR(255) NOT NULL,
    level_description JSONB NOT NULL DEFAULT '[]'
);

-- Insert maturity matrix data
INSERT INTO maturity_matrix (topic_name, level_description)
VALUES
    ('Finances', '[
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
    ('Work & Education', '[
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
    ('Use of Time', '[
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
    ('Housing', '[
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
    ('Domestic Relationships', '[
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
    ('Mental Health', '[
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
    ('Physical Health', '[
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
    ('Substance Use', '[
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
    ('Basic ADL', '[
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
    ('Instrumental ADL', '[
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
    ('Social Network', '[
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
    ('Social Participation', '[
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
    ('Justice', '[
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

-- Organizations and their locations
CREATE TABLE organisations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
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


-- Create ENUM type for location_type
CREATE TYPE location_type_enum AS ENUM ('care_home', 'office', 'other');
-- Location represents a physical place (care home, apartment building, etc.) for the youth intake
CREATE TABLE location (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organisation_id UUID NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    address VARCHAR(100) NOT NULL,
    capacity INTEGER NULL,
    location_type location_type_enum NOT NULL DEFAULT 'other',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Rooms within a location (studios, apartments, meeting rooms, etc.)
CREATE TABLE room (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    location_id UUID NOT NULL,
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
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    location_id UUID NOT NULL REFERENCES location(id) ON DELETE CASCADE,
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
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE
);

-- System permissions
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    resource VARCHAR(255) NOT NULL,
    method VARCHAR(255) NOT NULL
);

-- Role-to-Permission mapping (template)
CREATE TABLE role_permissions (
    role_id UUID NOT NULL,
    permission_id UUID NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

-- User authentication data
CREATE TABLE custom_user (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
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
    user_id UUID NOT NULL,
    permission_id UUID NOT NULL,
    PRIMARY KEY (user_id, permission_id),
    FOREIGN KEY (user_id) REFERENCES custom_user(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

-- Track which role templates were given to a user
CREATE TABLE user_roles (
    user_id UUID NOT NULL PRIMARY KEY,
    role_id UUID NOT NULL,
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
    "user_id" UUID NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY ("user_id") REFERENCES custom_user("id") ON DELETE CASCADE
);

CREATE INDEX idx_sessions_user ON sessions("user_id");
CREATE INDEX idx_sessions_expires ON sessions("expires_at");
CREATE INDEX idx_sessions_token_blocked ON sessions("refresh_token", "is_blocked");


-- Notification types ENUM
CREATE TYPE notification_type_enum AS ENUM (
    'new_appointment', 'appointment_update', 'new_client_assigned',
    'client_goal_update', 'incident_report', 'client_contract_reminder',
    'new_schedule_notification'
);

-- Notifications for users
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES custom_user(id) ON DELETE CASCADE,
    type notification_type_enum NOT NULL,
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

-- Employee Gender ENUM
CREATE TYPE employee_gender_enum AS ENUM ('male', 'female', 'not_specified');
-- Employee Contract Type ENUM
CREATE TYPE employee_contract_type_enum AS ENUM ('loondienst', 'ZZP', 'none');
-- Employee profile (linked to custom_user)
CREATE TABLE employee_profile (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES custom_user(id) ON DELETE CASCADE,
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
    gender employee_gender_enum NOT NULL,
    location_id UUID NULL REFERENCES location(id) ON DELETE SET NULL,
    has_borrowed BOOLEAN NOT NULL DEFAULT FALSE,
    out_of_service BOOLEAN NULL DEFAULT FALSE,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    contract_hours FLOAT NULL DEFAULT 0.0,
    contract_end_date DATE NULL,
    contract_start_date DATE NULL,
    contract_type employee_contract_type_enum NOT NULL DEFAULT 'none',
    contract_rate DECIMAL(10,2) NULL DEFAULT 0.00
);

CREATE INDEX employee_profile_user_id_idx ON employee_profile(user_id);
CREATE INDEX employee_profile_location_id_idx ON employee_profile(location_id);
CREATE INDEX employee_profile_id_desc_idx ON employee_profile(id DESC);
CREATE INDEX idx_employee_profile_is_archived ON employee_profile(is_archived);
CREATE INDEX idx_employee_profile_out_of_service ON employee_profile(out_of_service);

-- Employee education records
CREATE TABLE employee_education (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
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
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    issued_by VARCHAR(255) NOT NULL,
    date_issued DATE NOT NULL,
    created_at TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX certification_employee_id_idx ON certification(employee_id);

-- Employee work experience
CREATE TABLE employee_experience (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
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

-- Sender types ENUM
CREATE TYPE sender_types_enum AS ENUM (
    'main_provider', 'local_authority',
    'particular_party', 'healthcare_institution'
);
-- Sender organizations (municipalities, authorities, etc.)
CREATE TABLE sender (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    types sender_types_enum NOT NULL,
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
    invoice_template UUID[] NULL,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX sender_types_idx ON sender(types);


-- Registration forms
-- Form Status ENUM
CREATE TYPE form_status_enum AS ENUM ('pending', 'approved', 'rejected');
-- Client Gender ENUM
CREATE TYPE client_gender_enum AS ENUM ('male', 'female', 'other');
-- Client Education Level ENUM
CREATE TYPE client_education_level_enum AS ENUM ('primary', 'secondary', 'higher', 'none');
CREATE TABLE registration_form (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_first_name VARCHAR(255) NOT NULL,
    client_last_name VARCHAR(255) NOT NULL,
    client_bsn_number VARCHAR(50) NOT NULL,
    client_gender client_gender_enum NOT NULL,
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
    education_level client_education_level_enum NULL,
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
    form_status form_status_enum NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    submitted_at TIMESTAMPTZ NULL,
    processed_at TIMESTAMPTZ NULL,
    processed_by_employee_id UUID NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    status TEXT NOT NULL CHECK (status IN ('new', 'in_review', 'approved', 'rejected')) DEFAULT 'new',
    intake_appointment_datetime TIMESTAMPTZ NULL,
    intake_appointment_location VARCHAR(255) NULL,
    addmission_type VARCHAR(50) NULL CHECK (addmission_type IN ('crisis_admission', 'regular_placement'))
);


-- Care type ENUM
CREATE TYPE intake_care_type_enum AS ENUM ('protected_living', 'training_center', 'supported_independent_living', 'ambulatory_support', 'other');
-- Intake participants ENUM
CREATE TYPE intake_participants_enum AS ENUM ('client', 'referrer', 'parents/guardians', 'care_coordinator', 'other');
-- Intake conclusion ENUM
CREATE TYPE intake_conclusion_enum AS ENUM ('suitable', 'unsuitable', 'further_investigation', 'possible_palcement_date', 'other');
-- Initial intake forms
CREATE TABLE intake_forms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    registration_form_id UUID NOT NULL UNIQUE REFERENCES registration_form(id) ON DELETE CASCADE,
    date_of_intake TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    care_type intake_care_type_enum NOT NULL,
    intake_participants intake_participants_enum[] NOT NULL DEFAULT '{}',
    family_situation TEXT NULL,
    psychological_state TEXT NULL,
    self_sufficiency INT NOT NULL CHECK (self_sufficiency BETWEEN 0 AND 5),
    maturity_matrix_id UUID NULL REFERENCES maturity_matrix(id) ON DELETE SET NULL,
    goals TEXT NULL,
    risk_assessment TEXT NULL,
    intake_conclusion intake_conclusion_enum NOT NULL,
    intake_conclusion_notes TEXT NULL,
    signature TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- Client STATUS ENUM
CREATE TYPE client_status_enum AS ENUM ('In Care', 'On Waiting List', 'Out Of Care');


-- Clients living Situation ENUM
CREATE TYPE client_living_situation_enum AS ENUM ('home', 'foster_care', 'youth_care_institution', 'other');

-- Main client details table
CREATE TABLE client_details (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    intake_form_id UUID NULL REFERENCES intake_forms(id) ON DELETE SET NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE NULL,
    "identity" BOOLEAN NOT NULL DEFAULT FALSE,
    "status" client_status_enum NOT NULL DEFAULT 'On Waiting List',
    bsn VARCHAR(50) NULL,
    bsn_verified_by UUID NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    source VARCHAR(100) NULL,
    birthplace VARCHAR(100) NULL,
    email VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20) NULL,
    organization_id UUID NULL REFERENCES organisations(id) ON DELETE SET NULL,
    departement VARCHAR(100) NULL,
    gender client_gender_enum NOT NULL,
    filenumber VARCHAR(100) NOT NULL,
    profile_picture VARCHAR(600) NULL,
    infix VARCHAR(100) NULL,
    created_at TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP,
    sender_id UUID NULL REFERENCES sender(id) ON DELETE SET NULL DEFAULT NULL,
    location_id UUID NULL REFERENCES location(id) ON DELETE SET NULL DEFAULT NULL,
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
    education_level client_education_level_enum NOT NULL DEFAULT 'none',
    -- Work
    work_currently_employed BOOLEAN NOT NULL DEFAULT FALSE,
    work_current_employer VARCHAR(255) NULL,
    work_current_employer_phone VARCHAR(50) NULL,
    work_current_employer_email VARCHAR(255) NULL,
    work_current_position VARCHAR(255) NULL,
    work_start_date DATE NULL,
    work_additional_notes TEXT NULL,
    -- Living situation
    living_situation client_living_situation_enum NULL DEFAULT NULL,
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
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    old_status VARCHAR(50),
    new_status VARCHAR(50) NOT NULL,
    changed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    changed_by UUID NULL,
    reason VARCHAR(255)
);

CREATE INDEX idx_client_status_history_client_id ON client_status_history(client_id);
CREATE INDEX idx_client_status_history_changed_at ON client_status_history(changed_at DESC);

-- Scheduled status changes
CREATE TABLE scheduled_status_changes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    new_status VARCHAR(50) NULL,
    reason TEXT,
    scheduled_date DATE NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Client diagnoses
CREATE TABLE client_diagnosis (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(50) NULL,
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
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
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    soft_delete BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX contact_relationship_soft_delete_idx ON contact_relationship(soft_delete);

-- Client Emergency Contacts Relationship Status ENUM
CREATE TYPE relation_status_enum AS ENUM ('Primary Relationship', 'Secondary Relationship');
-- Client emergency contacts
CREATE TABLE client_emergency_contact (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    first_name VARCHAR(50) NULL,
    last_name VARCHAR(100) NULL,
    email VARCHAR(100) NULL,
    phone_number VARCHAR(20) NULL,
    address VARCHAR(100) NULL,
    relationship VARCHAR(100) NULL,
    relation_status relation_status_enum NULL,
    created_at TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    medical_reports BOOLEAN NOT NULL DEFAULT FALSE,
    incidents_reports BOOLEAN NOT NULL DEFAULT FALSE,
    goals_reports BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX client_emergency_contact_client_id_idx ON client_emergency_contact(client_id);

-- Client Documents Labels ENUM
CREATE TYPE client_document_label_enum AS ENUM (
    'registration_form', 'intake_form', 'consent_form',
    'risk_assessment', 'self_reliance_matrix', 'force_inventory',
    'care_plan', 'signaling_plan', 'cooperation_agreement', 'other'
);
-- Client documents
CREATE TABLE client_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attachment_uuid UUID NULL REFERENCES attachment_file("uuid") ON DELETE SET NULL,
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    label client_document_label_enum NOT NULL DEFAULT 'other'
);

CREATE INDEX client_documents_user_id_idx ON client_documents(client_id);
CREATE INDEX client_documents_label_idx ON client_documents(label);

-- Client medications
CREATE TABLE client_medication (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    diagnosis_id UUID NULL REFERENCES client_diagnosis(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    dosage VARCHAR(100) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NULL,
    notes TEXT NULL,
    self_administered BOOLEAN NOT NULL DEFAULT TRUE,
    slots JSONB NULL DEFAULT '[]',
    administered_by_id UUID NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    is_critical BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

-- Client location transfer status ENUM
CREATE TYPE  client_location_transfer_status_enum AS ENUM ('pending', 'approved', 'rejected');

CREATE TABLE client_location_transfer (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    from_location_id UUID NULL REFERENCES location(id) ON DELETE SET NULL,
    to_location_id UUID NULL REFERENCES location(id) ON DELETE SET NULL,
    new_mentor_id UUID NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    request_date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status client_location_transfer_status_enum NOT NULL DEFAULT 'pending',
    approved_rejected_by UUID NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    approved_rejected_at TIMESTAMPTZ NULL,
    reason TEXT NULL
);

-- ==========================================
-- CONTRACTS & FINANCIAL MANAGEMENT
-- ==========================================

-- Contract types
CREATE TABLE contract_type (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL
);
-- Contract status ENUM
CREATE TYPE contract_status_enum AS ENUM ('approved', 'draft', 'terminated', 'stopped', 'expired');
-- Price time unit ENUM
CREATE TYPE price_time_unit_enum AS ENUM ('minute', 'hourly', 'daily', 'weekly', 'monthly');
-- Hours type ENUM
CREATE TYPE hours_type_enum AS ENUM ('weekly', 'all_period');
-- Care type ENUM
CREATE TYPE care_type_enum AS ENUM ('ambulante', 'accommodation');
-- Financing act ENUM
CREATE TYPE financing_act_enum AS ENUM ('WMO', 'ZVW', 'WLZ', 'JW', 'WPG');
-- Financing option ENUM
CREATE TYPE financing_option_enum AS ENUM ('ZIN', 'PGB');
-- Main contracts table
CREATE TABLE contract (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type_id UUID NULL REFERENCES contract_type(id) ON DELETE SET NULL,
    status contract_status_enum NOT NULL DEFAULT 'draft',
    approved_at TIMESTAMPTZ NULL,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    reminder_period INTEGER NOT NULL DEFAULT 90,
    VAT INTEGER NULL DEFAULT -1,
    price DECIMAL(10,2) NOT NULL,
    price_time_unit price_time_unit_enum NOT NULL DEFAULT 'weekly',
    hours DECIMAL(10,2) NULL DEFAULT 0,
    hours_type hours_type_enum NOT NULL DEFAULT NULL,
    care_name VARCHAR(255) NOT NULL,
    care_type care_type_enum NOT NULL,
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    sender_id UUID NULL REFERENCES sender(id) ON DELETE SET NULL,
    attachment_ids UUID[] NOT NULL DEFAULT '{}',
    financing_act financing_act_enum NOT NULL DEFAULT 'WMO',
    financing_option financing_option_enum NOT NULL DEFAULT 'PGB',
    departure_reason VARCHAR(255) NULL,
    departure_report TEXT NULL,
    updated_at TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX contract_type_id_idx ON contract(type_id);
CREATE INDEX contract_client_id_idx ON contract(client_id);
CREATE INDEX contract_sender_id_idx ON contract(sender_id);
CREATE INDEX contract_status_idx ON contract(status);

-- Contract audit operations ENUM
CREATE TYPE contract_audit_operation_enum AS ENUM ('INSERT', 'UPDATE', 'DELETE');

-- Contract audit table
CREATE TABLE contract_audit (
    audit_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contract_id UUID NOT NULL,
    operation contract_audit_operation_enum NOT NULL,
    changed_by UUID NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
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
    current_user_id UUID;
BEGIN
    BEGIN
        current_user_id := current_setting('myapp.current_employee_id')::UUID;
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

-- Contract reminder reminder_type ENUM
CREATE TYPE contract_reminder_type_enum AS ENUM ('initial', 'follow_up', 'none');
-- Contract-related tables
CREATE TABLE contract_reminder (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contract_id UUID NOT NULL REFERENCES contract(id) ON DELETE CASCADE,
    reminder_sent_at TIMESTAMPTZ NULL,
    reminder_type contract_reminder_type_enum NOT NULL DEFAULT 'none'
);

CREATE TABLE contract_working_hours (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contract_id UUID NOT NULL REFERENCES contract(id) ON DELETE CASCADE,
    minutes INTEGER NOT NULL DEFAULT 0,
    "datetime" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    notes TEXT NULL DEFAULT ''
);

CREATE INDEX contract_working_hours_contract_id_idx ON contract_working_hours(contract_id);
CREATE INDEX contract_working_hours_datetime_idx ON contract_working_hours(datetime);

CREATE TABLE contract_attachment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contract_id UUID NOT NULL REFERENCES contract(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    attachment VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX contract_attachment_contract_id_idx ON contract_attachment(contract_id);

CREATE TABLE client_agreement (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contract_id UUID NOT NULL REFERENCES contract(id) ON DELETE CASCADE,
    agreement_details TEXT NOT NULL,
    created TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX client_agreement_contract_id_idx ON client_agreement(contract_id);

CREATE TABLE provision (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contract_id UUID NOT NULL REFERENCES contract(id) ON DELETE CASCADE,
    provision_details TEXT NOT NULL,
    created TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX provision_contract_id_idx ON provision(contract_id);

CREATE TABLE framework_agreement (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    agreement_details TEXT NOT NULL,
    created TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX framework_agreement_client_id_idx ON framework_agreement(client_id);


-- Invoice management

-- Invoice status ENUM
CREATE TYPE invoice_status_enum AS ENUM (
    'outstanding', 'partially_paid', 'paid', 'expired',
    'overpaid', 'imported', 'concept', 'canceled'
);
-- Invoice type ENUM
CREATE TYPE invoice_type_enum AS ENUM ('standard', 'credit_note');
-- Main invoice table
CREATE TABLE invoice (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_number VARCHAR(50) NOT NULL UNIQUE,
    invoice_sequence BIGINT NOT NULL DEFAULT 1,
    issue_date DATE NOT NULL DEFAULT CURRENT_DATE,
    due_date DATE NOT NULL,
    status invoice_status_enum NOT NULL DEFAULT 'concept',
    invoice_type invoice_type_enum NOT NULL DEFAULT 'standard',
    original_invoice_id UUID NULL REFERENCES invoice(id) ON DELETE SET NULL,
    invoice_details JSONB NULL DEFAULT '[]',
    total_amount DECIMAL(20,2) NOT NULL DEFAULT 0,
    pdf_attachment_id UUID NULL UNIQUE REFERENCES attachment_file("uuid") ON DELETE SET NULL,
    extra_content JSONB NULL DEFAULT '{}',
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    sender_id UUID NULL REFERENCES sender(id) ON DELETE SET NULL,
    warning_count INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX invoice_invoice_number_idx ON invoice(invoice_number);
CREATE INDEX invoice_client_id_idx ON invoice(client_id);
CREATE INDEX invoice_status_idx ON invoice(status);

-- Invoice audit table
-- Invoice audit operations ENUM
CREATE TYPE invoice_audit_operation_enum AS ENUM ('INSERT', 'UPDATE', 'DELETE');
CREATE TABLE invoice_audit (
    audit_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL,
    operation invoice_audit_operation_enum NOT NULL,
    changed_by UUID REFERENCES employee_profile(id) ON DELETE SET NULL,
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
    current_employee_id UUID;
BEGIN
    BEGIN
        current_employee_id := current_setting('myapp.current_employee_id')::UUID;
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

-- Payment method ENUM
CREATE TYPE payment_method_enum AS ENUM (
    'bank_transfer', 'credit_card', 'check', 'cash', 'other'
);
-- Payment status ENUM
CREATE TYPE payment_status_enum AS ENUM (
    'completed', 'pending', 'failed', 'reversed', 'refunded'
);
CREATE TABLE invoice_payment_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL REFERENCES invoice(id) ON DELETE CASCADE,
    payment_method payment_method_enum NOT NULL DEFAULT 'bank_transfer',
    payment_status payment_status_enum NOT NULL DEFAULT 'completed',
    amount DECIMAL(20,2) NOT NULL DEFAULT 0,
    payment_date DATE NOT NULL DEFAULT CURRENT_DATE,
    payment_reference VARCHAR(100) NULL,
    notes TEXT NULL,
    recorded_by UUID NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_invoice_payment_history_invoice_id ON invoice_payment_history(invoice_id);
CREATE INDEX idx_invoice_payment_history_payment_date ON invoice_payment_history(payment_date);
CREATE INDEX idx_invoice_payment_history_payment_status ON invoice_payment_history(payment_status);

-- Invoice-Contract relationship
CREATE TABLE invoice_contract (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NULL REFERENCES invoice(id) ON DELETE SET NULL,
    contract_id UUID NULL REFERENCES contract(id) ON DELETE SET NULL,
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


-- Client maturity matrix assessments
-- Care plan status ENUM
CREATE TYPE care_plan_status_enum AS ENUM ('pending', 'generated', 'approved', 'active', 'completed', 'discontinued');
CREATE TABLE client_maturity_matrix_assessment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    maturity_matrix_id UUID NOT NULL REFERENCES maturity_matrix(id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    initial_level INT NOT NULL CHECK (initial_level BETWEEN 1 AND 5),
    target_level INT NOT NULL CHECK (target_level BETWEEN 1 AND 5),
    current_level INT NOT NULL CHECK (current_level BETWEEN 1 AND 5),
    care_plan_generated_at TIMESTAMPTZ NULL DEFAULT NULL,
    care_plan_status care_plan_status_enum NOT NULL DEFAULT 'pending',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE(client_id, maturity_matrix_id)
);

-- Level change history
CREATE TABLE level_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_maturity_matrix_assessment_id UUID NOT NULL REFERENCES client_maturity_matrix_assessment(id) ON DELETE CASCADE,
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
        VALUES (OLD.id, OLD.current_level, NEW.current_level, 'Automatic logging of level change.');
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
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_id UUID NOT NULL REFERENCES client_maturity_matrix_assessment(id) ON DELETE CASCADE,
    generated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    generated_by_employee_id UUID REFERENCES employee_profile(id),
    approved_by_employee_id UUID REFERENCES employee_profile(id),
    approved_at TIMESTAMP,
    status care_plan_status_enum NOT NULL DEFAULT 'generated',
    assessment_summary TEXT NOT NULL,
    raw_llm_response JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version INT NOT NULL DEFAULT 1,
    UNIQUE(assessment_id)
);

-- Care plan objectives
-- Timeframe ENUM
CREATE TYPE care_plan_timeframe_enum AS ENUM ('short_term', 'medium_term', 'long_term');
-- Care plan objective status ENUM
CREATE TYPE care_plan_objective_status_enum AS ENUM ('not_started', 'in_progress', 'completed', 'discontinued', 'draft');
CREATE TABLE care_plan_objectives (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    care_plan_id UUID NOT NULL REFERENCES care_plans(id) ON DELETE CASCADE,
    timeframe care_plan_timeframe_enum NOT NULL,
    goal_title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    target_date DATE,
    status care_plan_objective_status_enum NOT NULL DEFAULT 'not_started',
    completion_date DATE,
    completion_notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Care plan actions
CREATE TABLE care_plan_actions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    objective_id UUID NOT NULL REFERENCES care_plan_objectives(id) ON DELETE CASCADE,
    action_description TEXT NOT NULL,
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    completed_at TIMESTAMP,
    completed_by_employee_id UUID REFERENCES employee_profile(id),
    notes TEXT,
    sort_order INT NOT NULL DEFAULT 0
);

-- Care plan interventions
-- Frequency ENUM
CREATE TYPE care_plan_intervention_frequency_enum AS ENUM ('daily', 'weekly', 'monthly');
CREATE TABLE care_plan_interventions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    care_plan_id UUID NOT NULL REFERENCES care_plans(id) ON DELETE CASCADE,
    frequency care_plan_intervention_frequency_enum NOT NULL,
    intervention_description TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_completed_date DATE,
    total_completions INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Care plan metrics
CREATE TABLE care_plan_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    care_plan_id UUID NOT NULL REFERENCES care_plans(id) ON DELETE CASCADE,
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
-- Risk level ENUM
CREATE TYPE care_plan_risk_level_enum AS ENUM ('low', 'medium', 'high');
CREATE TABLE care_plan_risks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    care_plan_id UUID NOT NULL REFERENCES care_plans(id) ON DELETE CASCADE,
    risk_description TEXT NOT NULL,
    mitigation_strategy TEXT NOT NULL,
    risk_level care_plan_risk_level_enum NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Care plan support network
CREATE TABLE care_plan_support_network (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    care_plan_id UUID NOT NULL REFERENCES care_plans(id) ON DELETE CASCADE,
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
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    care_plan_id UUID NOT NULL REFERENCES care_plans(id) ON DELETE CASCADE,
    resource_description TEXT NOT NULL,
    is_obtained BOOLEAN NOT NULL DEFAULT FALSE,
    obtained_date DATE,
    cost_estimate DECIMAL(10,2),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Care plan reports
-- Report type ENUM
CREATE TYPE care_plan_report_type_enum AS ENUM ('progress', 'concern', 'achievement', 'modification');
CREATE TABLE care_plan_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    care_plan_id UUID NOT NULL REFERENCES care_plans(id) ON DELETE CASCADE,
    report_type care_plan_report_type_enum NOT NULL,
    report_content TEXT NOT NULL,
    created_by_employee_id UUID NOT NULL REFERENCES employee_profile(id),
    is_critical BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ==========================================
-- INCIDENTS & REPORTING
-- ==========================================

-- Incident reports
-- Incident reporter involvement ENUM
CREATE TYPE incident_reporter_involvement_enum AS ENUM ('directly_involved', 'witness', 'found_afterwards', 'alarmed');
-- Severity of incident ENUM
CREATE TYPE severity_of_incident_enum AS ENUM ('near_incident', 'less_serious', 'serious', 'fatal');
-- Recurrence risk ENUM
CREATE TYPE recurrence_risk_enum AS ENUM ('very_low', 'means', 'high', 'very_high');
-- Physical injury ENUM
CREATE TYPE physical_injury_enum AS ENUM ('no_injuries', 'not_noticeable_yet', 'bruising_swelling', 'skin_injury', 'broken_bones', 'shortness_of_breath', 'death', 'other');
-- Psychological damage ENUM
CREATE TYPE psychological_damage_enum AS ENUM ('no', 'not_noticeable_yet', 'drowsiness', 'unrest', 'other');
-- Needed consultation ENUM
CREATE TYPE needed_consultation_enum AS ENUM ('no', 'not_clear', 'hospitalization', 'consult_gp');
CREATE TABLE incident (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    location_id UUID NOT NULL REFERENCES location(id) ON DELETE CASCADE,
    reporter_involvement incident_reporter_involvement_enum NOT NULL,
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
    severity_of_incident severity_of_incident_enum NOT NULL,
    incident_explanation TEXT NULL,
    recurrence_risk recurrence_risk_enum NOT NULL,
    incident_prevent_steps TEXT NULL,
    incident_taken_measures TEXT NULL,
    technical VARCHAR(255)[] NOT NULL DEFAULT '{}',
    organizational VARCHAR(255)[] NOT NULL DEFAULT '{}',
    mese_worker VARCHAR(255)[] NOT NULL DEFAULT '{}',
    client_options VARCHAR(255)[] NOT NULL DEFAULT '{}',
    other_cause VARCHAR(100) NULL,
    cause_explanation TEXT NULL DEFAULT '',
    physical_injury physical_injury_enum NOT NULL,
    physical_injury_desc TEXT NULL DEFAULT '',
    psychological_damage psychological_damage_enum NOT NULL,
    psychological_damage_desc TEXT NULL DEFAULT '',
    needed_consultation needed_consultation_enum NOT NULL,
    succession VARCHAR(255)[] NOT NULL DEFAULT '{}',
    succession_desc TEXT NULL DEFAULT '',
    other BOOLEAN NOT NULL DEFAULT FALSE,
    other_desc VARCHAR(100) NULL,
    additional_appointments TEXT NULL DEFAULT '',
    employee_absenteeism VARCHAR(100) NOT NULL DEFAULT '',
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
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
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    start_datetime TIMESTAMPTZ NOT NULL,
    end_datetime TIMESTAMPTZ NOT NULL,
    "status" VARCHAR(50) NOT NULL CHECK (status IN ('Confirmed', 'Pending', 'Cancelled')),
    created_at TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX assignment_employee_id_idx ON assignment(employee_id);
CREATE INDEX assignment_client_id_idx ON assignment(client_id);

-- Assigned employees to clients
CREATE TABLE assigned_employee (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    role VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX assigned_employee_client_id_idx ON assigned_employee(client_id);
CREATE INDEX assigned_employee_employee_id_idx ON assigned_employee(employee_id);

-- Progress reports
-- Progress report types ENUM
CREATE TYPE progress_report_type_enum AS ENUM (
    'morning_report', 'evening_report', 'night_report', 'shift_report',
    'one_to_one_report', 'process_report', 'contact_journal', 'other'
);
-- Emotional state ENUM
CREATE TYPE emotional_state_enum AS ENUM (
    'normal', 'excited', 'happy', 'sad', 'angry', 'anxious', 'depressed'
);
CREATE TABLE progress_report (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    date TIMESTAMPTZ NOT NULL,
    title VARCHAR(50) NULL,
    report_text TEXT NOT NULL,
    employee_id UUID NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    type progress_report_type_enum NOT NULL,
    emotional_state emotional_state_enum NOT NULL DEFAULT 'normal',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX progress_report_client_id_idx ON progress_report(client_id);
CREATE INDEX progress_report_author_id_idx ON progress_report(employee_id);
CREATE INDEX progress_report_created_idx ON progress_report(created_at DESC);

-- AI generated reports
CREATE TABLE ai_generated_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_text TEXT NOT NULL,
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
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
    employee_id UUID NOT NULL REFERENCES employee_profile(id),
    color VARCHAR(20) DEFAULT '#0000FF',
    location_id UUID NOT NULL REFERENCES location(id),
    location_shift_id UUID NULL REFERENCES location_shift(id),
    is_custom BOOLEAN NOT NULL DEFAULT FALSE,
    start_datetime TIMESTAMP NOT NULL,
    end_datetime TIMESTAMP NOT NULL,
    created_by_employee_id UUID NOT NULL REFERENCES employee_profile(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_timeframe CHECK (end_datetime > start_datetime)
);

-- Appointment templates
-- Recurrence type ENUM
CREATE TYPE recurrence_type_enum AS ENUM ('DAILY', 'WEEKLY', 'MONTHLY');
CREATE TABLE appointment_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    creator_employee_id UUID NOT NULL REFERENCES employee_profile(id),
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    location VARCHAR(255),
    description TEXT,
    color VARCHAR(20) DEFAULT '#0000FF',
    recurrence_type recurrence_type_enum NOT NULL DEFAULT 'DAILY',
    recurrence_interval INT NULL,
    recurrence_end_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Scheduled appointments
-- Appointment status ENUM
CREATE TYPE appointment_status_enum AS ENUM ('PENDING', 'CONFIRMED', 'CANCELLED');
CREATE TABLE scheduled_appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_templates_id UUID NULL REFERENCES appointment_templates(id) ON DELETE CASCADE,
    creator_employee_id UUID NULL REFERENCES employee_profile(id),
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    location VARCHAR(255),
    description TEXT,
    status appointment_status_enum NOT NULL DEFAULT 'PENDING',
    color VARCHAR(20) DEFAULT '#0000FF',
    is_confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    confirmed_by_employee_id UUID REFERENCES employee_profile(id),
    confirmed_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_scheduled_appointments_time_range ON scheduled_appointments (start_time, end_time);
CREATE INDEX idx_scheduled_appointments_template_id ON scheduled_appointments (appointment_templates_id);

-- Appointment participants
CREATE TABLE appointment_participants (
    appointment_participant_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID NOT NULL REFERENCES scheduled_appointments(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employee_profile(id),
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (appointment_id, employee_id)
);

-- Appointment clients
CREATE TABLE appointment_clients (
    appointment_client_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID NOT NULL REFERENCES scheduled_appointments(id) ON DELETE CASCADE,
    client_id UUID NOT NULL REFERENCES client_details(id),
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (appointment_id, client_id)
);

-- ==========================================
-- FORMS & DOCUMENTATION
-- ==========================================

-- Appointment cards
CREATE TABLE appointment_card (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE UNIQUE,
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



-- Collaboration agreements
CREATE TABLE collaboration_agreement (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
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
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
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
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
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
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    pdf_attachment_id UUID NULL UNIQUE REFERENCES attachment_file("uuid") ON DELETE SET NULL
);

CREATE INDEX consent_declaration_client_id_idx ON consent_declaration(client_id);

-- Youth care intake forms
CREATE TABLE youth_care_intake (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
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
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
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
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ==========================================
-- UTILITY TABLES
-- ==========================================

-- Template items for document generation for custom data fields to include in invoice generation
CREATE TABLE template_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
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

-- ===============================================
-- AUDIT LOGGING
-- ===============================================
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

-- ===============================================
-- ROW LEVEL SECURITY (RLS)
-- ===============================================

-- Helper function to get current employee ID
CREATE OR REPLACE FUNCTION get_current_employee_id() RETURNS UUID AS $$
BEGIN
    RETURN current_setting('myapp.current_employee_id', true)::UUID;
EXCEPTION
    WHEN OTHERS THEN RETURN NULL;
END;
$$ LANGUAGE plpgsql STABLE;

-- Check if current employee is Admin
CREATE OR REPLACE FUNCTION is_admin() RETURNS BOOLEAN AS $$
BEGIN
    RETURN EXISTS (
        SELECT 1 FROM user_roles ur
        JOIN roles r ON ur.role_id = r.id
        JOIN employee_profile ep ON ep.user_id = ur.user_id
        WHERE ep.id = get_current_employee_id()
        AND r.name = 'admin'
    );
END;
$$ LANGUAGE plpgsql STABLE SECURITY DEFINER;

-- Check if current employee is Coordinator
CREATE OR REPLACE FUNCTION is_coordinator() RETURNS BOOLEAN AS $$
BEGIN
    RETURN EXISTS (
        SELECT 1 FROM user_roles ur
        JOIN roles r ON ur.role_id = r.id
        JOIN employee_profile ep ON ep.user_id = ur.user_id
        WHERE ep.id = get_current_employee_id()
        AND r.name = 'coordinator'
    );
END;
$$ LANGUAGE plpgsql STABLE SECURITY DEFINER;

-- Check if current employee is assigned as coordinator to a client
CREATE OR REPLACE FUNCTION is_assigned_coordinator(cid UUID) RETURNS BOOLEAN AS $$
BEGIN
    IF cid IS NULL THEN RETURN FALSE; END IF;
    RETURN EXISTS (
        SELECT 1 FROM assigned_employee
        WHERE client_id = cid
        AND employee_id = get_current_employee_id()
        AND role = 'coordinator'
    );
END;
$$ LANGUAGE plpgsql STABLE SECURITY DEFINER;

-- Helper to get client_id from various sub-tables
CREATE OR REPLACE FUNCTION get_client_id_from_care_plan(cp_id UUID) RETURNS UUID AS $$
    SELECT a.client_id FROM client_maturity_matrix_assessment a JOIN care_plans cp ON cp.assessment_id = a.id WHERE cp.id = cp_id;
$$ LANGUAGE sql STABLE;

CREATE OR REPLACE FUNCTION get_client_id_from_objective(obj_id UUID) RETURNS UUID AS $$
    SELECT get_client_id_from_care_plan(care_plan_id) FROM care_plan_objectives WHERE id = obj_id;
$$ LANGUAGE sql STABLE;

CREATE OR REPLACE FUNCTION get_client_id_from_diagnosis(diag_id UUID) RETURNS UUID AS $$
    SELECT client_id FROM client_diagnosis WHERE id = diag_id;
$$ LANGUAGE sql STABLE;

CREATE OR REPLACE FUNCTION get_client_id_from_contract(cont_id UUID) RETURNS UUID AS $$
    SELECT client_id FROM contract WHERE id = cont_id;
$$ LANGUAGE sql STABLE;

CREATE OR REPLACE FUNCTION get_client_id_from_registration_form(reg_id UUID) RETURNS UUID AS $$
    SELECT cd.id FROM client_details cd JOIN intake_forms iform ON iform.id = cd.intake_form_id WHERE iform.registration_form_id = reg_id;
$$ LANGUAGE sql STABLE;

-- Function to apply policies to a table with 'client_id' (or custom column)
CREATE OR REPLACE FUNCTION apply_client_rls(table_name TEXT, client_id_col TEXT DEFAULT 'client_id') 
RETURNS VOID AS $$
BEGIN
    EXECUTE format('ALTER TABLE %I ENABLE ROW LEVEL SECURITY', table_name);
    EXECUTE format('DROP POLICY IF EXISTS coordinator_select ON %I', table_name);
    EXECUTE format('DROP POLICY IF EXISTS coordinator_insert ON %I', table_name);
    EXECUTE format('DROP POLICY IF EXISTS coordinator_update ON %I', table_name);
    EXECUTE format('DROP POLICY IF EXISTS coordinator_delete ON %I', table_name);
    
    EXECUTE format('CREATE POLICY coordinator_select ON %I FOR SELECT USING (is_admin() OR is_coordinator())', table_name);
    EXECUTE format('CREATE POLICY coordinator_insert ON %I FOR INSERT WITH CHECK (is_admin() OR is_coordinator())', table_name);
    EXECUTE format('CREATE POLICY coordinator_update ON %I FOR UPDATE USING (is_admin() OR is_coordinator()) WITH CHECK (is_admin() OR is_coordinator())', table_name);
    EXECUTE format('CREATE POLICY coordinator_delete ON %I FOR DELETE USING (is_admin() OR is_assigned_coordinator(%I))', table_name, client_id_col);
END;
$$ LANGUAGE plpgsql;

-- Apply RLS to tables with direct client_id or id
SELECT apply_client_rls('client_details', 'id');
SELECT apply_client_rls('progress_report');
SELECT apply_client_rls('incident');
SELECT apply_client_rls('client_maturity_matrix_assessment');
SELECT apply_client_rls('client_documents');
SELECT apply_client_rls('client_status_history');
SELECT apply_client_rls('scheduled_status_changes');
SELECT apply_client_rls('client_diagnosis');
SELECT apply_client_rls('client_emergency_contact');
SELECT apply_client_rls('client_location_transfer');
SELECT apply_client_rls('contract');
SELECT apply_client_rls('invoice');
SELECT apply_client_rls('assignment');
SELECT apply_client_rls('assigned_employee');
SELECT apply_client_rls('ai_generated_reports');
SELECT apply_client_rls('appointment_clients');
SELECT apply_client_rls('appointment_card');
SELECT apply_client_rls('collaboration_agreement');
SELECT apply_client_rls('risk_assessment');
SELECT apply_client_rls('consent_declaration');
SELECT apply_client_rls('youth_care_intake');
SELECT apply_client_rls('data_sharing_statement');
SELECT apply_client_rls('framework_agreement');

-- Special cases for nested tables
-- Registration and Intake
SELECT apply_client_rls('registration_form', 'get_client_id_from_registration_form(id)');
SELECT apply_client_rls('intake_forms', 'get_client_id_from_registration_form(registration_form_id)');

-- Medication
SELECT apply_client_rls('client_medication', 'get_client_id_from_diagnosis(diagnosis_id)');

-- Contract sub-tables
SELECT apply_client_rls('client_agreement', 'get_client_id_from_contract(contract_id)');
SELECT apply_client_rls('provision', 'get_client_id_from_contract(contract_id)');

-- Care Plans and sub-tables
SELECT apply_client_rls('care_plans', '(SELECT client_id FROM client_maturity_matrix_assessment WHERE id = assessment_id)');
SELECT apply_client_rls('care_plan_objectives', 'get_client_id_from_care_plan(care_plan_id)');
SELECT apply_client_rls('care_plan_interventions', 'get_client_id_from_care_plan(care_plan_id)');
SELECT apply_client_rls('care_plan_metrics', 'get_client_id_from_care_plan(care_plan_id)');
SELECT apply_client_rls('care_plan_risks', 'get_client_id_from_care_plan(care_plan_id)');
SELECT apply_client_rls('care_plan_support_network', 'get_client_id_from_care_plan(care_plan_id)');
SELECT apply_client_rls('care_plan_resources', 'get_client_id_from_care_plan(care_plan_id)');
SELECT apply_client_rls('care_plan_reports', 'get_client_id_from_care_plan(care_plan_id)');
SELECT apply_client_rls('care_plan_actions', 'get_client_id_from_objective(objective_id)');
SELECT apply_client_rls('level_history', '(SELECT client_id FROM client_maturity_matrix_assessment WHERE id = client_maturity_matrix_assessment_id)');

-- Clean up helper functions if desired, or keep them for future use.
-- DROP FUNCTION apply_client_rls(TEXT, TEXT);




