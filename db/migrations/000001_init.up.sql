-- ==========================================
-- INFRASTRUCTURE & ORGANIZATIONS
-- ==========================================

-- topics topics and levels
CREATE TABLE topics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    topic_name VARCHAR(255) NOT NULL,
    level_description JSONB NOT NULL DEFAULT '[]'
);

-- Insert topics data
INSERT INTO topics (topic_name, level_description)
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
    street VARCHAR(200) NOT NULL,
    house_number VARCHAR(20) NOT NULL,
    house_number_addition VARCHAR(20) NULL,
    postal_code VARCHAR(20) NOT NULL,
    city VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20) NULL,
    email VARCHAR(100) NULL,
    kvk_number VARCHAR(20) NULL,
    btw_number VARCHAR(20) NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE app_organization_profile (
    singleton BOOLEAN PRIMARY KEY DEFAULT TRUE CHECK (singleton),
    name TEXT NOT NULL DEFAULT '',
    default_timezone TEXT NOT NULL DEFAULT 'Europe/Amsterdam',
    email TEXT NULL,
    phone_number TEXT NULL,
    website TEXT NULL,
    hq_street TEXT NULL,
    hq_house_number TEXT NULL,
    hq_house_number_addition TEXT NULL,
    hq_postal_code TEXT NULL,
    hq_city TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO app_organization_profile (singleton)
VALUES (TRUE);


-- Create ENUM type for location_type
CREATE TYPE location_type_enum AS ENUM ('care_home', 'office', 'other');
-- Location represents a physical place (care home, apartment building, etc.) for the youth intake
CREATE TABLE location (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organisation_id UUID NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    street VARCHAR(200) NOT NULL,
    house_number VARCHAR(20) NOT NULL,
    house_number_addition VARCHAR(20) NULL,
    postal_code VARCHAR(20) NOT NULL,
    city VARCHAR(100) NOT NULL,
    timezone TEXT NOT NULL DEFAULT 'Europe/Amsterdam',
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
    slot SMALLINT NOT NULL CHECK (slot BETWEEN 1 AND 4),
    shift_name VARCHAR(50) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(location_id, slot),
    UNIQUE(location_id, shift_name)
);

-- Function to insert default shifts for new locations
CREATE OR REPLACE FUNCTION insert_default_shifts()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO location_shift (location_id, slot, shift_name, start_time, end_time)
    VALUES
        (NEW.id, 1, 'Ochtenddienst', TIME '07:30:00', TIME '15:30:00'),
        (NEW.id, 2, 'Avonddienst', TIME '15:00:00', TIME '23:00:00'),
        (NEW.id, 3, 'Slaapdienst of Waakdienst', TIME '23:00:00', TIME '07:30:00');
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
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT NULL
);

-- System permissions
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    resource VARCHAR(255) NOT NULL,
    method VARCHAR(255) NOT NULL,
    group_key VARCHAR(100) NOT NULL DEFAULT 'general',
    section_key VARCHAR(100) NOT NULL DEFAULT 'general',
    display_name VARCHAR(255) NOT NULL DEFAULT '',
    description TEXT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0
);

CREATE TYPE permission_override_effect AS ENUM ('allow', 'deny');

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

-- Explicit per-user permission exceptions layered on top of role inheritance
CREATE TABLE user_permission_overrides (
    user_id UUID NOT NULL,
    permission_id UUID NOT NULL,
    effect permission_override_effect NOT NULL,
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
    'new_schedule_notification', 'system_reminder'
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

-- Shared Gender ENUM
CREATE TYPE gender_enum AS ENUM ('male', 'female', 'other', 'unknown');
-- Employee Contract Type ENUM
CREATE TYPE employee_contract_type_enum AS ENUM ('loondienst', 'ZZP', 'none');

-- Departments (used for employee assignment and handbook templates)
CREATE TABLE departments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT NULL,
    department_head_employee_id UUID NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX departments_name_idx ON departments(name);
CREATE INDEX departments_department_head_employee_id_idx ON departments(department_head_employee_id);

-- Employee profile (linked to custom_user)
CREATE TABLE employee_profile (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES custom_user(id) ON DELETE CASCADE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    bsn TEXT NOT NULL,
    street TEXT NOT NULL,
    house_number TEXT NOT NULL,
    house_number_addition TEXT NULL,
    postal_code TEXT NOT NULL,
    city TEXT NOT NULL,
    position VARCHAR(100) NULL,
    employee_number VARCHAR(50) NULL,
    employment_number VARCHAR(50) NULL,
    private_email_address VARCHAR(254) NULL,
    work_email_address VARCHAR(254) NULL,
    private_phone_number VARCHAR(100) NULL,
    work_phone_number VARCHAR(100) NULL,
    date_of_birth DATE NULL,
    home_telephone_number VARCHAR(100) NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    gender gender_enum NOT NULL,
    location_id UUID NULL REFERENCES location(id) ON DELETE SET NULL,
    department_id UUID NULL REFERENCES departments(id) ON DELETE SET NULL,
    manager_employee_id UUID NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    has_borrowed BOOLEAN NOT NULL DEFAULT FALSE,
    out_of_service BOOLEAN NULL DEFAULT FALSE,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    contract_hours FLOAT NULL DEFAULT 0.0,
    contract_end_date DATE NULL,
    contract_start_date DATE NULL,
    contract_type employee_contract_type_enum NOT NULL DEFAULT 'none',
    contract_rate DECIMAL(10,2) NULL DEFAULT 0.00,
    CONSTRAINT employee_profile_manager_not_self
        CHECK (manager_employee_id IS NULL OR manager_employee_id <> id)
);

CREATE INDEX employee_profile_user_id_idx ON employee_profile(user_id);
CREATE INDEX employee_profile_location_id_idx ON employee_profile(location_id);
CREATE INDEX idx_employee_profile_department_id ON employee_profile(department_id);
CREATE INDEX idx_employee_profile_manager_employee_id ON employee_profile(manager_employee_id);
CREATE INDEX employee_profile_id_desc_idx ON employee_profile(id DESC);
CREATE INDEX idx_employee_profile_is_archived ON employee_profile(is_archived);
CREATE INDEX idx_employee_profile_out_of_service ON employee_profile(out_of_service);

ALTER TABLE departments
    ADD CONSTRAINT departments_department_head_employee_id_fkey
    FOREIGN KEY (department_head_employee_id) REFERENCES employee_profile(id) ON DELETE SET NULL;

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
-- EMPLOYEE HANDBOOKS (ONBOARDING)
-- ==========================================

CREATE TYPE handbook_step_kind_enum AS ENUM ('content', 'ack', 'link', 'quiz');
CREATE TYPE handbook_assignment_status_enum AS ENUM ('not_started', 'in_progress', 'completed', 'waived');
CREATE TYPE handbook_step_status_enum AS ENUM ('pending', 'completed', 'skipped');
CREATE TYPE handbook_template_status_enum AS ENUM ('draft', 'published', 'archived');
CREATE TYPE handbook_assignment_event_enum AS ENUM ('assigned', 'reassigned', 'waived', 'started', 'completed');

-- A template is department-specific and can be versioned. Assignments point to a specific version.
CREATE TABLE handbook_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    department_id UUID NOT NULL REFERENCES departments(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT NULL,
    version INT NOT NULL,
    status handbook_template_status_enum NOT NULL DEFAULT 'draft',
    created_by_employee_id UUID NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    published_by_employee_id UUID NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    published_at TIMESTAMPTZ NULL,
    archived_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (department_id, version)
);

-- Enforce at most one published template per department.
CREATE UNIQUE INDEX handbook_templates_one_published_per_department
    ON handbook_templates(department_id)
    WHERE status = 'published';

-- Enforce at most one draft template per department.
CREATE UNIQUE INDEX handbook_templates_one_draft_per_department
    ON handbook_templates(department_id)
    WHERE status = 'draft';

CREATE INDEX idx_handbook_templates_department_id ON handbook_templates(department_id);

CREATE TABLE handbook_steps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID NOT NULL REFERENCES handbook_templates(id) ON DELETE CASCADE,
    sort_order INT NOT NULL CHECK (sort_order > 0),
    kind handbook_step_kind_enum NOT NULL DEFAULT 'content',
    title TEXT NOT NULL,
    body TEXT NULL,
    content JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_required BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (template_id, sort_order)
);

CREATE INDEX idx_handbook_steps_template_sort ON handbook_steps(template_id, sort_order);

CREATE TABLE employee_handbooks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    template_id UUID NOT NULL REFERENCES handbook_templates(id) ON DELETE RESTRICT,
    template_version INT NOT NULL,
    assigned_by_employee_id UUID NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    status handbook_assignment_status_enum NOT NULL DEFAULT 'not_started',
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMPTZ NULL,
    completed_at TIMESTAMPTZ NULL,
    due_at TIMESTAMPTZ NULL
);

-- At most one active handbook (not_started/in_progress) per employee.
CREATE UNIQUE INDEX employee_handbooks_one_active_per_employee
    ON employee_handbooks(employee_id)
    WHERE status IN ('not_started', 'in_progress');

CREATE INDEX idx_employee_handbooks_employee_assigned_at ON employee_handbooks(employee_id, assigned_at DESC);

CREATE TABLE employee_handbook_assignment_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_handbook_id UUID NULL REFERENCES employee_handbooks(id) ON DELETE SET NULL,
    employee_id UUID NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    template_id UUID NOT NULL REFERENCES handbook_templates(id) ON DELETE RESTRICT,
    template_version INT NOT NULL,
    event handbook_assignment_event_enum NOT NULL,
    actor_employee_id UUID NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_employee_handbook_assignment_history_employee_created_at
    ON employee_handbook_assignment_history(employee_id, created_at DESC);

CREATE TABLE employee_handbook_step_progress (
    employee_handbook_id UUID NOT NULL REFERENCES employee_handbooks(id) ON DELETE CASCADE,
    step_id UUID NOT NULL REFERENCES handbook_steps(id) ON DELETE RESTRICT,
    status handbook_step_status_enum NOT NULL DEFAULT 'pending',
    started_at TIMESTAMPTZ NULL,
    completed_at TIMESTAMPTZ NULL,
    response JSONB NOT NULL DEFAULT '{}'::jsonb,
    PRIMARY KEY (employee_handbook_id, step_id)
);

CREATE INDEX idx_employee_handbook_step_progress_handbook_id
    ON employee_handbook_step_progress(employee_handbook_id);

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
    name TEXT NOT NULL,
    street TEXT NULL,
    house_number TEXT NULL,
    house_number_addition TEXT NULL,
    postal_code TEXT NULL,
    city TEXT NULL,
    land TEXT NULL,
    kvknumber TEXT NULL,
    btwnumber TEXT NULL,
    phone_number TEXT NULL,
    client_number TEXT NULL,
    email_address  TEXT NULL,
    contacts JSONB NOT NULL DEFAULT '[]',
    invoice_template UUID[] NULL,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX sender_types_idx ON sender(types);


-- Registration forms
-- Form Status ENUM
CREATE TYPE form_status_enum AS ENUM ('pending', 'processed', 'rejected');
-- Shared Education Level ENUM
CREATE TYPE education_level_enum AS ENUM ('primary', 'secondary', 'higher', 'none');
-- Admission Type ENUM
CREATE TYPE admission_type_enum AS ENUM ('crisis_admission', 'regular_placement');
CREATE TABLE registration_form (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_first_name VARCHAR(255) NOT NULL,
    client_last_name VARCHAR(255) NOT NULL,
    client_date_of_birth DATE NULL,
    client_bsn_number VARCHAR(50) NOT NULL,
    client_gender gender_enum NOT NULL,
    client_nationality VARCHAR(100) NOT NULL,
    client_phone_number VARCHAR(20) NOT NULL,
    client_email VARCHAR(255) NOT NULL,
    -- Client address
    client_street VARCHAR(255) NOT NULL,
    client_house_number VARCHAR(20) NOT NULL,
    client_house_number_addition VARCHAR(255) NULL,
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
    education_level education_level_enum NOT NULL DEFAULT 'none',
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
    client_goals TEXT[] DEFAULT '{}',
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
    intake_options JSONB DEFAULT '[]',
    intake_token VARCHAR(255) UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    submitted_at TIMESTAMPTZ NULL,
    processed_at TIMESTAMPTZ NULL,
    processed_by_employee_id UUID NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    intake_appointment_datetime TIMESTAMPTZ NULL,
    intake_appointment_location VARCHAR(255) NULL,
    addmission_type admission_type_enum NOT NULL DEFAULT 'regular_placement',
    rejection_reason TEXT NULL
);

CREATE INDEX registration_form_intake_token_idx ON registration_form(intake_token);


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
    sender_id UUID NULL REFERENCES sender(id) ON DELETE SET NULL,
    assigned_location_id UUID NULL REFERENCES location(id) ON DELETE SET NULL,
    risk_assessment TEXT NULL,
    intake_conclusion intake_conclusion_enum NOT NULL,
    intake_conclusion_notes TEXT NULL,
    evaluation_intervals_weeks INT NOT NULL DEFAULT 0,
    signature TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Intake topics assessments (goals per topic)
CREATE TABLE intake_topic_assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    intake_form_id UUID NOT NULL REFERENCES intake_forms(id) ON DELETE CASCADE,
    topic_id UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    current_level INT NOT NULL CHECK (current_level BETWEEN 1 AND 5),
    proposed_goals JSONB NOT NULL DEFAULT '[]',
    notes TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(intake_form_id, topic_id)
);


-- Client STATUS ENUM
CREATE TYPE client_status_enum AS ENUM (
    'in_care',
    'on_waiting_list',
    'scheduled_in_care',
    'scheduled_out_of_care',
    'out_of_care'
);

CREATE TYPE discharge_reason_enum AS ENUM (
    'treatment_completed',
    'terminated_by_mutual_agreement',
    'terminated_by_client',
    'terminated_by_provider',
    'terminated_due_to_external_factors',
    'other'
);


-- Clients living Situation ENUM
CREATE TYPE client_living_situation_enum AS ENUM ('home', 'foster_care', 'youth_care_institution', 'other');

-- Client file number sequence and formatter (YYYY-00001)
CREATE SEQUENCE client_filenumber_seq START 1;

CREATE OR REPLACE FUNCTION generate_client_filenumber()
RETURNS TEXT AS $$
BEGIN
    RETURN to_char(CURRENT_DATE, 'YYYY') || '-' || lpad(nextval('client_filenumber_seq')::TEXT, 5, '0');
END;
$$ LANGUAGE plpgsql;

-- Main client details table
CREATE TABLE client_details (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    intake_form_id UUID NULL REFERENCES intake_forms(id) ON DELETE SET NULL,
    registration_form_id UUID NULL REFERENCES registration_form(id) ON DELETE SET NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE NULL,
    "identity" BOOLEAN NOT NULL DEFAULT FALSE,
    "status" client_status_enum NOT NULL DEFAULT 'on_waiting_list',
    bsn VARCHAR(50) NULL,
    bsn_verified_by UUID NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    evaluation_intervals_weeks INT NOT NULL DEFAULT 0,
    care_type intake_care_type_enum NULL,
    -- source VARCHAR(100) NULL, -- Not needed now
    -- birthplace VARCHAR(100) NULL, -- Not needed now
    email VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20) NULL,
    -- organization_id UUID NULL REFERENCES organisations(id) ON DELETE SET NULL, -- Not needed now
    -- departement VARCHAR(100) NULL, -- Not needed now
    gender gender_enum NOT NULL,
    filenumber VARCHAR(100) NOT NULL DEFAULT generate_client_filenumber() UNIQUE,
    -- profile_picture VARCHAR(600) NULL, -- Not needed now
    -- infix VARCHAR(100) NULL, -- Not needed now
    created_at TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP,
    placed_in_care_at TIMESTAMPTZ NULL,
    care_start_date DATE NULL,
    last_evaluation_anchor_date DATE NULL,
    next_evaluation_date DATE NULL,
    discharge_date DATE NULL,
    discharge_reason discharge_reason_enum NULL,
    final_evaluation TEXT NULL,
    sender_id UUID NULL REFERENCES sender(id) ON DELETE SET NULL DEFAULT NULL,
    location_id UUID NULL REFERENCES location(id) ON DELETE SET NULL DEFAULT NULL,
    -- departure_reason VARCHAR(255) NULL, -- Not needed now
    -- departure_report TEXT NULL, -- Not needed now
    -- gps_position JSONB NOT NULL DEFAULT '[]', -- Not needed now
    -- maturity_domains JSONB NOT NULL DEFAULT '[]', -- Not needed now
    street VARCHAR(255) NOT NULL,
    house_number VARCHAR(20) NOT NULL,
    house_number_addition VARCHAR(255) NULL,
    postal_code VARCHAR(20) NOT NULL,
    city VARCHAR(100) NOT NULL,
    -- legal_measure VARCHAR(255) NULL, -- Not needed now
    -- has_untaken_medications BOOLEAN NOT NULL DEFAULT FALSE, -- Not needed now
    -- Education
    education_currently_enrolled BOOLEAN NOT NULL DEFAULT FALSE,
    education_institution VARCHAR(255) NULL,
    education_mentor_name VARCHAR(255) NULL,
    education_mentor_phone VARCHAR(50) NULL,
    education_mentor_email VARCHAR(255) NULL,
    education_additional_notes TEXT NULL,
    education_level education_level_enum NOT NULL DEFAULT 'none',
    -- Work
    work_currently_employed BOOLEAN NOT NULL DEFAULT FALSE,
    work_current_employer VARCHAR(255) NULL,
    work_current_employer_phone VARCHAR(50) NULL,
    work_current_employer_email VARCHAR(255) NULL,
    work_current_position VARCHAR(255) NULL,
    work_start_date DATE NULL,
    work_additional_notes TEXT NULL,
    -- Nationality
    nationality VARCHAR(100) NULL,
    -- Living situation
    -- living_situation client_living_situation_enum NULL DEFAULT NULL, -- Not needed now
    -- living_situation_notes TEXT NULL, -- Not needed now

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
    CONSTRAINT client_details_care_dates_required_for_care_status
        CHECK (
            status NOT IN ('scheduled_in_care', 'in_care')
            OR (care_start_date IS NOT NULL AND placed_in_care_at IS NOT NULL)
        ),
    CONSTRAINT client_details_waiting_list_has_no_care_dates
        CHECK (
            status <> 'on_waiting_list'
            OR (care_start_date IS NULL AND placed_in_care_at IS NULL)
        ),
    CONSTRAINT client_details_discharge_fields_match_status
        CHECK (
            (
                status = 'out_of_care'
                AND discharge_date IS NOT NULL
                AND discharge_reason IS NOT NULL
                AND final_evaluation IS NOT NULL
            )
            OR (
                status = 'scheduled_out_of_care'
                AND discharge_date IS NOT NULL
                AND discharge_reason IS NOT NULL
            )
            OR (
                status NOT IN ('scheduled_out_of_care', 'out_of_care')
                AND discharge_date IS NULL
                AND discharge_reason IS NULL
                AND final_evaluation IS NULL
            )
        ),
    UNIQUE (intake_form_id)
);

CREATE INDEX client_details_sender_id_idx ON client_details(sender_id);
CREATE INDEX client_details_location_id_idx ON client_details(location_id);
CREATE INDEX client_details_registration_form_id_idx ON client_details(registration_form_id);


-- Client goals and grouped evaluations (new model)
CREATE TYPE client_goal_priority_enum AS ENUM ('low', 'medium', 'high');
CREATE TYPE client_goal_status_enum AS ENUM ('active', 'achieved', 'cancelled');
CREATE TYPE client_goal_source_enum AS ENUM ('intake', 'manual', 'review_update');
CREATE TYPE client_goal_progress_enum AS ENUM (
    'no_progress',
    'regression',
    'limited_progress',
    'good_progress',
    'achieved',
    'blocked'
);

CREATE TABLE client_goals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT NULL,
    priority client_goal_priority_enum NOT NULL DEFAULT 'medium',
    status client_goal_status_enum NOT NULL DEFAULT 'active',
    topic_id UUID NULL REFERENCES topics(id) ON DELETE SET NULL,
    topic_name_snapshot VARCHAR(255) NULL,
    source client_goal_source_enum NOT NULL DEFAULT 'intake',
    origin_intake_assessment_id UUID NULL REFERENCES intake_topic_assessments(id) ON DELETE SET NULL,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    archived_at TIMESTAMPTZ NULL
);

CREATE INDEX client_goals_client_status_idx ON client_goals(client_id, status);
CREATE INDEX client_goals_client_sort_order_idx ON client_goals(client_id, sort_order);

CREATE TYPE evaluation_status_enum AS ENUM ('draft', 'completed', 'archived');

CREATE TABLE client_goal_evaluations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    evaluation_date DATE NOT NULL DEFAULT CURRENT_DATE,
    period_start DATE NULL,
    period_end DATE NULL,
    evaluation_interval_weeks INT NOT NULL,
    status evaluation_status_enum NOT NULL DEFAULT 'draft',
    overall_notes TEXT NULL,
    created_by_employee_id UUID NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX client_goal_evaluations_client_date_idx ON client_goal_evaluations(client_id, evaluation_date DESC);
CREATE INDEX client_goal_evaluations_client_created_idx ON client_goal_evaluations(client_id, created_at DESC);
CREATE UNIQUE INDEX client_goal_evaluations_unique_draft_client_date_idx
    ON client_goal_evaluations(client_id, evaluation_date)
    WHERE status = 'draft';

CREATE TABLE client_goal_evaluation_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    evaluation_id UUID NOT NULL REFERENCES client_goal_evaluations(id) ON DELETE CASCADE,
    goal_id UUID NOT NULL REFERENCES client_goals(id) ON DELETE CASCADE,
    progress client_goal_progress_enum NOT NULL,
    notes TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(evaluation_id, goal_id)
);

CREATE INDEX client_goal_evaluation_items_goal_created_idx ON client_goal_evaluation_items(goal_id, created_at DESC);
CREATE INDEX client_goal_evaluation_items_evaluation_idx ON client_goal_evaluation_items(evaluation_id);

-- ==========================================
-- EVALUATION SCHEDULING & CADENCE LOGIC
-- ==========================================

-- Function to initialize evaluation dates when a client enters care
CREATE OR REPLACE FUNCTION initialize_client_evaluation_dates()
RETURNS TRIGGER AS $$
BEGIN
    -- If care_start_date is being set for the first time
    IF NEW.care_start_date IS NOT NULL AND (OLD.care_start_date IS NULL OR OLD.care_start_date <> NEW.care_start_date) THEN
        -- Initialize anchor as the start date
        NEW.last_evaluation_anchor_date := NEW.care_start_date;

        -- Set next evaluation date based on intervals
        IF NEW.evaluation_intervals_weeks > 0 THEN
            NEW.next_evaluation_date := NEW.care_start_date + (NEW.evaluation_intervals_weeks * INTERVAL '1 week');
        ELSE
            -- Default to 12 weeks if not specified
            NEW.next_evaluation_date := NEW.care_start_date + INTERVAL '12 weeks';
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_initialize_client_evaluation_dates
BEFORE UPDATE OF care_start_date ON client_details
FOR EACH ROW
EXECUTE FUNCTION initialize_client_evaluation_dates();

-- Function to prevent premature evaluation completion
CREATE OR REPLACE FUNCTION enforce_evaluation_submission_window()
RETURNS TRIGGER AS $$
DECLARE
    v_next_eval_date DATE;
BEGIN
    -- Only check when transitioning to 'completed'
    IF NEW.status = 'completed' AND (OLD.status IS NULL OR OLD.status <> 'completed') THEN
        SELECT next_evaluation_date INTO v_next_eval_date
        FROM client_details
        WHERE id = NEW.client_id;

        -- Refuse if more than 14 days before the due date
        IF CURRENT_DATE < (v_next_eval_date - INTERVAL '14 days') THEN
            RAISE EXCEPTION 'Evaluation cannot be completed more than 14 days before the due date (%)', v_next_eval_date;
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_enforce_evaluation_submission_window
BEFORE UPDATE OF status ON client_goal_evaluations
FOR EACH ROW
EXECUTE FUNCTION enforce_evaluation_submission_window();

-- Function to update next evaluation date upon completion
CREATE OR REPLACE FUNCTION update_client_evaluation_cadence()
RETURNS TRIGGER AS $$
DECLARE
    v_current_next DATE;
    v_interval_weeks INT;
BEGIN
    -- When an evaluation is completed
    IF NEW.status = 'completed' AND (OLD.status IS NULL OR OLD.status <> 'completed') THEN
        SELECT next_evaluation_date
        INTO v_current_next
        FROM client_details
        WHERE id = NEW.client_id;

        v_interval_weeks := COALESCE(NULLIF(NEW.evaluation_interval_weeks, 0), 12);

        -- Update the client record
        UPDATE client_details
        SET
            last_evaluation_anchor_date = v_current_next,
            next_evaluation_date = v_current_next + (v_interval_weeks * INTERVAL '1 week')
        WHERE id = NEW.client_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_client_evaluation_cadence
AFTER UPDATE OF status ON client_goal_evaluations
FOR EACH ROW
EXECUTE FUNCTION update_client_evaluation_cadence();

-- Guardrail: cannot move a client into/scheduled for care without active goals
CREATE OR REPLACE FUNCTION ensure_client_has_active_goals_before_care_status()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status IN ('scheduled_in_care', 'in_care') THEN
        IF NOT EXISTS (
            SELECT 1
            FROM client_goals cg
            WHERE cg.client_id = NEW.id
              AND cg.status = 'active'
        ) THEN
            RAISE EXCEPTION 'client must have at least one active goal before status %', NEW.status;
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER SET search_path = public;

CREATE TRIGGER trigger_ensure_client_has_active_goals_before_care_status
BEFORE INSERT OR UPDATE OF status ON client_details
FOR EACH ROW
EXECUTE FUNCTION ensure_client_has_active_goals_before_care_status();

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

-- -- Scheduled status changes
-- CREATE TABLE scheduled_status_changes (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
--     new_status client_status_enum NULL,
--     reason TEXT,
--     scheduled_date DATE NULL,
--     created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
-- );

-- Client diagnoses (v2)
CREATE TYPE diagnosis_status_enum AS ENUM (
    'suspected',
    'confirmed',
    'resolved',
    'ruled_out'
);

CREATE TYPE diagnosis_severity_enum AS ENUM ('mild', 'moderate', 'severe', 'unknown');

CREATE TABLE client_diagnosis (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,

    code_system TEXT NOT NULL,
    code TEXT NOT NULL,
    title TEXT NULL,
    description TEXT NULL,

    status diagnosis_status_enum NOT NULL DEFAULT 'confirmed',
    severity diagnosis_severity_enum NOT NULL DEFAULT 'unknown',

    diagnosed_on DATE NULL,
    resolved_on DATE NULL,
    diagnosing_clinician TEXT NULL,
    notes TEXT NULL,

    created_by_employee_id UUID NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    updated_by_employee_id UUID NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    archived_at TIMESTAMPTZ NULL,

    CONSTRAINT client_diagnosis_dates CHECK (resolved_on IS NULL OR diagnosed_on IS NULL OR resolved_on >= diagnosed_on)
);

CREATE INDEX client_diagnosis_client_status_idx ON client_diagnosis(client_id, status);
CREATE INDEX client_diagnosis_client_code_idx ON client_diagnosis(client_id, code_system, code);



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
CREATE UNIQUE INDEX client_documents_attachment_uuid_unique_idx
ON client_documents(attachment_uuid)
WHERE attachment_uuid IS NOT NULL;

-- Client medication orders (v2)
CREATE TYPE medication_order_status_enum AS ENUM ('active', 'paused', 'stopped', 'completed');
CREATE TYPE medication_admin_mode_enum AS ENUM ('self', 'staff', 'shared');

CREATE TABLE client_medication_order (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    diagnosis_id UUID NULL REFERENCES client_diagnosis(id) ON DELETE SET NULL,

    medication_name TEXT NOT NULL,
    dosage_text TEXT NOT NULL,

    dose_amount NUMERIC NULL,
    dose_unit TEXT NULL,
    route TEXT NULL,
    frequency_text TEXT NULL,
    schedule JSONB NOT NULL DEFAULT '[]',

    is_prn BOOLEAN NOT NULL DEFAULT FALSE,
    prn_indication TEXT NULL,
    max_doses_per_24h INT NULL,

    start_date DATE NOT NULL,
    end_date DATE NULL,
    status medication_order_status_enum NOT NULL DEFAULT 'active',

    admin_mode medication_admin_mode_enum NOT NULL DEFAULT 'self',
    responsible_employee_id UUID NULL REFERENCES employee_profile(id) ON DELETE SET NULL,

    is_critical BOOLEAN NOT NULL DEFAULT FALSE,
    notes TEXT NULL,
    source_attachment_uuid UUID NULL REFERENCES attachment_file(uuid) ON DELETE SET NULL,

    created_by_employee_id UUID NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    updated_by_employee_id UUID NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    archived_at TIMESTAMPTZ NULL,

    CONSTRAINT client_med_order_dates CHECK (end_date IS NULL OR end_date >= start_date),
    CONSTRAINT client_med_order_admin_responsible CHECK (
        (admin_mode IN ('staff', 'shared') AND responsible_employee_id IS NOT NULL)
        OR (admin_mode = 'self' AND responsible_employee_id IS NULL)
    ),
    CONSTRAINT client_med_order_prn_rules CHECK (
        (is_prn = FALSE)
        OR (is_prn = TRUE AND prn_indication IS NOT NULL)
    )
);

CREATE INDEX client_med_order_client_status_idx ON client_medication_order(client_id, status);
CREATE INDEX client_med_order_client_dates_idx ON client_medication_order(client_id, start_date DESC);
CREATE INDEX client_med_order_diag_idx ON client_medication_order(diagnosis_id);

-- Auto-update updated_at on row updates
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_set_updated_at_client_diagnosis
BEFORE UPDATE ON client_diagnosis
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trigger_set_updated_at_client_medication_order
BEFORE UPDATE ON client_medication_order
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

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
CREATE TYPE price_time_unit_enum AS ENUM ('minute', 'hourly', 'daily', 'weekly');
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
    VAT INTEGER NULL DEFAULT 20,
    price DECIMAL(10,2) NOT NULL,
    price_time_unit price_time_unit_enum NOT NULL DEFAULT 'weekly',
    hours DECIMAL(10,2) NULL,
    hours_type hours_type_enum NULL,
    care_name VARCHAR(255) NOT NULL,
    care_type care_type_enum NOT NULL,
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES sender(id) ON DELETE RESTRICT,
    attachment_ids UUID[] NOT NULL DEFAULT '{}',
    financing_act financing_act_enum NOT NULL DEFAULT 'WMO',
    financing_option financing_option_enum NOT NULL DEFAULT 'PGB',
    departure_reason VARCHAR(255) NULL,
    departure_report TEXT NULL,
    updated_at TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT contract_care_type_pricing_and_hours_consistency CHECK (
        (
            care_type = 'ambulante'
            AND price_time_unit IN ('minute', 'hourly')
            AND hours IS NOT NULL
            AND hours > 0
            AND hours_type IS NOT NULL
        )
        OR
        (
            care_type = 'accommodation'
            AND price_time_unit IN ('daily', 'weekly')
            AND hours IS NULL
            AND hours_type IS NULL
        )
    )
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

-- Invoice source ENUM (how the invoice was created)
CREATE TYPE invoice_source_enum AS ENUM ('auto', 'manual', 'imported');

-- Invoice line type ENUM (what kind of line it is)
CREATE TYPE invoice_line_type_enum AS ENUM ('contract', 'manual', 'adjustment');

-- Invoice run status ENUM (batch generation)
CREATE TYPE invoice_run_status_enum AS ENUM ('running', 'completed', 'completed_with_errors', 'failed');
CREATE TYPE invoice_run_item_status_enum AS ENUM ('created', 'skipped', 'failed');

-- Batch run (4-ISO-week automation)
CREATE TABLE invoice_run (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    billing_cycle TEXT NOT NULL, -- e.g. 'iso_4_week'
    period_start TIMESTAMPTZ NOT NULL,
    period_end TIMESTAMPTZ NOT NULL,
    timezone TEXT NOT NULL DEFAULT 'UTC',
    dry_run BOOLEAN NOT NULL DEFAULT FALSE,
    status invoice_run_status_enum NOT NULL DEFAULT 'running',
    params JSONB NOT NULL DEFAULT '{}'::jsonb,
    started_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    finished_at TIMESTAMPTZ NULL,
    created_by UUID NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    CHECK (period_start < period_end)
);

CREATE INDEX invoice_run_period_idx ON invoice_run(period_start, period_end);
-- Main invoice table
CREATE TABLE invoice (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_number VARCHAR(50) NOT NULL UNIQUE,
    invoice_sequence BIGINT NOT NULL DEFAULT 1,
    issue_date DATE NOT NULL DEFAULT CURRENT_DATE,
    due_date DATE NOT NULL,
    status invoice_status_enum NOT NULL DEFAULT 'concept',
    invoice_type invoice_type_enum NOT NULL DEFAULT 'standard',
    source invoice_source_enum NOT NULL DEFAULT 'manual',
    original_invoice_id UUID NULL REFERENCES invoice(id) ON DELETE SET NULL,
    replaces_invoice_id UUID NULL REFERENCES invoice(id) ON DELETE SET NULL,

    -- Billing window (used for auto-generation idempotency)
    period_start TIMESTAMPTZ NULL,
    period_end TIMESTAMPTZ NULL,
    billing_cycle TEXT NULL, -- e.g. 'iso_4_week'
    billing_timezone TEXT NOT NULL DEFAULT 'UTC',

    -- Optional snapshots at issuance time (avoid depending on mutable sender/client details)
    bill_to_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
    client_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,

    currency CHAR(3) NOT NULL DEFAULT 'EUR',
    details_snapshot JSONB NOT NULL DEFAULT '[]',
    net_total_amount NUMERIC(20,2) NOT NULL DEFAULT 0,
    vat_total_amount NUMERIC(20,2) NOT NULL DEFAULT 0,
    gross_total_amount NUMERIC(20,2) NOT NULL DEFAULT 0,
    pdf_attachment_id UUID NULL UNIQUE REFERENCES attachment_file("uuid") ON DELETE SET NULL,
    extra_content JSONB NULL DEFAULT '{}',
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES sender(id) ON DELETE RESTRICT,
    warning_count INTEGER NOT NULL DEFAULT 0,
    run_id UUID NULL REFERENCES invoice_run(id) ON DELETE SET NULL,
    locked_at TIMESTAMPTZ NULL,
    calc_version INTEGER NOT NULL DEFAULT 1,
    calc_metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK ((period_start IS NULL AND period_end IS NULL) OR (period_start < period_end)),
    CHECK (invoice_type <> 'credit_note' OR original_invoice_id IS NOT NULL)
);

CREATE INDEX invoice_invoice_number_idx ON invoice(invoice_number);
CREATE INDEX invoice_client_id_idx ON invoice(client_id);
CREATE INDEX invoice_status_idx ON invoice(status);
CREATE INDEX invoice_sender_id_idx ON invoice(sender_id);
CREATE INDEX invoice_period_idx ON invoice(period_start, period_end);

-- Only one active auto standard invoice per (sender, client, period)
CREATE UNIQUE INDEX invoice_auto_standard_period_uq
ON invoice(sender_id, client_id, period_start, period_end)
WHERE source = 'auto' AND invoice_type = 'standard' AND status <> 'canceled';

CREATE TABLE invoice_run_item (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    run_id UUID NOT NULL REFERENCES invoice_run(id) ON DELETE CASCADE,
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES sender(id) ON DELETE RESTRICT,
    status invoice_run_item_status_enum NOT NULL DEFAULT 'created',
    invoice_id UUID NULL REFERENCES invoice(id) ON DELETE SET NULL,
    error TEXT NULL,
    warnings JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (run_id, client_id, sender_id)
);

CREATE INDEX invoice_run_item_run_idx ON invoice_run_item(run_id);
CREATE INDEX invoice_run_item_client_idx ON invoice_run_item(client_id);
CREATE INDEX invoice_run_item_sender_idx ON invoice_run_item(sender_id);

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

-- Canonical invoice lines (queryable & enforceable), replaces invoice_contract and JSON-only models
CREATE TABLE invoice_line (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL REFERENCES invoice(id) ON DELETE CASCADE,
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES sender(id) ON DELETE RESTRICT,

    line_no INTEGER NOT NULL,
    line_type invoice_line_type_enum NOT NULL DEFAULT 'contract',

    contract_id UUID NULL REFERENCES contract(id) ON DELETE SET NULL,

    service_type TEXT NOT NULL, -- e.g. 'accommodation' | 'ambulante'
    description TEXT NOT NULL,

    period_start TIMESTAMPTZ NULL,
    period_end TIMESTAMPTZ NULL,

    quantity NUMERIC(20,4) NOT NULL DEFAULT 0,
    unit TEXT NOT NULL, -- e.g. 'day' | 'hour' | 'minute' | 'item'
    unit_price NUMERIC(20,4) NOT NULL DEFAULT 0,

    net_amount NUMERIC(20,2) NOT NULL DEFAULT 0,
    vat_rate NUMERIC(5,2) NOT NULL DEFAULT 0,
    vat_amount NUMERIC(20,2) NOT NULL DEFAULT 0,
    gross_amount NUMERIC(20,2) NOT NULL DEFAULT 0,

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (invoice_id, line_no),
    CHECK ((period_start IS NULL AND period_end IS NULL) OR (period_start < period_end))
);

CREATE INDEX invoice_line_invoice_id_idx ON invoice_line(invoice_id);
CREATE INDEX invoice_line_client_id_idx ON invoice_line(client_id);
CREATE INDEX invoice_line_sender_id_idx ON invoice_line(sender_id);
CREATE INDEX invoice_line_contract_id_idx ON invoice_line(contract_id);

-- ==========================================
-- CARE PLANS & ASSESSMENTS
-- ==========================================

/*
Deprecated (intentionally commented out):
- client_topic_assessment
- level_history
- care_plans and care_plan_* tables
- care_plan_* enums and level-history trigger
This flow is replaced by client_goals + client_goal_evaluations.
*/

-- ==========================================
-- INCIDENTS & REPORTING
-- ==========================================

-- Incident reports
-- Incident reporter involvement ENUM
CREATE TYPE incident_reporter_involvement_enum AS ENUM ('directly_involved', 'witness', 'found_afterwards', 'alarmed');
-- Incident category ENUM (primary incident classification)
CREATE TYPE incident_type_enum AS ENUM (
    'passing_away',
    'self_harm',
    'violence',
    'fire_water_damage',
    'accident',
    'client_absence',
    'medicines',
    'organization',
    'use_prohibited_substances',
    'other'
);
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
-- Parties that can be informed about an incident
CREATE TYPE informed_party_enum AS ENUM (
    'parents_guardians',
    'care_coordinator',
    'referrer',
    'healthcare_provider',
    'inspectorate',
    'police',
    'other'
);
-- Cause categories for incident analysis
CREATE TYPE incident_cause_category_enum AS ENUM (
    'technical',
    'organizational',
    'employee_related',
    'client_related',
    'external',
    'other'
);
-- Follow-up actions after incident handling
CREATE TYPE incident_follow_up_action_enum AS ENUM (
    'notify_parents_guardians',
    'notify_referrer',
    'notify_inspectorate',
    'medical_consultation',
    'care_plan_adjustment',
    'team_evaluation',
    'other'
);
CREATE TABLE incident (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    location_id UUID NOT NULL REFERENCES location(id) ON DELETE CASCADE,
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    reporter_involvement incident_reporter_involvement_enum NOT NULL,
    informed_parties informed_party_enum[] NOT NULL DEFAULT '{}',
    occurred_at TIMESTAMPTZ NOT NULL,
    incident_type incident_type_enum NOT NULL,
    severity_of_incident severity_of_incident_enum NOT NULL,
    incident_explanation TEXT NULL,
    recurrence_risk recurrence_risk_enum NOT NULL,
    incident_prevent_steps TEXT NULL,
    incident_taken_measures TEXT NULL,
    cause_categories incident_cause_category_enum[] NOT NULL DEFAULT '{}',
    cause_explanation TEXT NULL,
    physical_injury physical_injury_enum NOT NULL,
    physical_injury_desc TEXT NULL,
    psychological_damage psychological_damage_enum NOT NULL,
    psychological_damage_desc TEXT NULL,
    needed_consultation needed_consultation_enum NOT NULL,
    follow_up_actions incident_follow_up_action_enum[] NOT NULL DEFAULT '{}',
    follow_up_notes TEXT NULL,
    is_employee_absent BOOLEAN NOT NULL DEFAULT FALSE,
    additional_details TEXT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    file_url VARCHAR(255) NULL,
    emails TEXT[] NULL DEFAULT '{}',
    confirmed_at TIMESTAMPTZ NULL,
    confirmed_by UUID NULL REFERENCES custom_user(id) ON DELETE SET NULL,
    confirmation_email_sent_at TIMESTAMPTZ NULL
);

CREATE INDEX incident_client_id_idx ON incident(client_id);
CREATE INDEX incident_location_id_idx ON incident(location_id);

CREATE TRIGGER trigger_set_updated_at_incident
BEFORE UPDATE ON incident
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

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
CREATE UNIQUE INDEX assigned_employee_one_coordinator_per_client_idx
    ON assigned_employee(client_id)
    WHERE role = 'coordinator';

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
    location_id UUID NOT NULL REFERENCES location(id),
    location_shift_id UUID NULL REFERENCES location_shift(id),
    shift_name_snapshot VARCHAR(50) NULL,
    shift_start_time_snapshot TIME NULL,
    shift_end_time_snapshot TIME NULL,
    is_custom BOOLEAN NOT NULL DEFAULT FALSE,
    start_datetime TIMESTAMPTZ NOT NULL,
    end_datetime TIMESTAMPTZ NOT NULL,
    created_by_employee_id UUID NOT NULL REFERENCES employee_profile(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_timeframe CHECK (end_datetime > start_datetime),
    CONSTRAINT schedules_custom_shift_link_consistency CHECK (
        (is_custom = TRUE AND location_shift_id IS NULL)
        OR (is_custom = FALSE AND location_shift_id IS NOT NULL)
    ),
    CONSTRAINT schedules_shift_snapshot_consistency CHECK (
        (location_shift_id IS NULL
         AND shift_name_snapshot IS NULL
         AND shift_start_time_snapshot IS NULL
         AND shift_end_time_snapshot IS NULL)
        OR
        (location_shift_id IS NOT NULL
         AND shift_name_snapshot IS NOT NULL
         AND shift_start_time_snapshot IS NOT NULL
         AND shift_end_time_snapshot IS NOT NULL)
    )
);

-- ==========================================
-- LATE ARRIVALS
-- ==========================================

CREATE TABLE late_arrivals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    schedule_id UUID NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    created_by_employee_id UUID NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    arrival_date DATE NOT NULL,
    arrival_time TIME NOT NULL,
    reason TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT late_arrivals_unique_schedule UNIQUE (schedule_id)
);

CREATE INDEX idx_late_arrivals_employee_id ON late_arrivals(employee_id);
CREATE INDEX idx_late_arrivals_arrival_date_desc ON late_arrivals(arrival_date DESC);
CREATE INDEX idx_late_arrivals_employee_date ON late_arrivals(employee_id, arrival_date DESC);

CREATE TYPE shift_swap_status_enum AS ENUM (
    'pending_recipient',
    'recipient_rejected',
    'pending_admin',
    'admin_rejected',
    'confirmed',
    'cancelled',
    'expired'
);

CREATE TABLE shift_swap_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    requester_employee_id UUID NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    recipient_employee_id UUID NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    requester_schedule_id UUID NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
    recipient_schedule_id UUID NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
    status shift_swap_status_enum NOT NULL DEFAULT 'pending_recipient',
    requested_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    recipient_responded_at TIMESTAMPTZ NULL,
    admin_decided_at TIMESTAMPTZ NULL,
    recipient_response_note TEXT NULL,
    admin_decision_note TEXT NULL,
    admin_employee_id UUID NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    expires_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT shift_swap_requester_not_recipient CHECK (requester_employee_id <> recipient_employee_id),
    CONSTRAINT shift_swap_schedule_pair_not_same CHECK (requester_schedule_id <> recipient_schedule_id)
);

CREATE INDEX idx_shift_swap_requests_requester_employee_id ON shift_swap_requests(requester_employee_id);
CREATE INDEX idx_shift_swap_requests_recipient_employee_id ON shift_swap_requests(recipient_employee_id);
CREATE INDEX idx_shift_swap_requests_status ON shift_swap_requests(status);
CREATE INDEX idx_shift_swap_requests_requested_at_desc ON shift_swap_requests(requested_at DESC);
CREATE INDEX idx_shift_swap_requests_expires_at ON shift_swap_requests(expires_at);

CREATE UNIQUE INDEX uq_shift_swap_active_requester_schedule
    ON shift_swap_requests(requester_schedule_id)
    WHERE status IN ('pending_recipient', 'pending_admin');

CREATE UNIQUE INDEX uq_shift_swap_active_recipient_schedule
    ON shift_swap_requests(recipient_schedule_id)
    WHERE status IN ('pending_recipient', 'pending_admin');

CREATE OR REPLACE FUNCTION enforce_shift_swap_active_schedule_uniqueness()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status IN ('pending_recipient', 'pending_admin') THEN
        IF EXISTS (
            SELECT 1
            FROM shift_swap_requests ssr
            WHERE ssr.id <> NEW.id
              AND ssr.status IN ('pending_recipient', 'pending_admin')
              AND (
                ssr.requester_schedule_id IN (NEW.requester_schedule_id, NEW.recipient_schedule_id)
                OR ssr.recipient_schedule_id IN (NEW.requester_schedule_id, NEW.recipient_schedule_id)
              )
        ) THEN
            RAISE EXCEPTION 'one of the schedules is already in an active swap request'
                USING ERRCODE = '23505', CONSTRAINT = 'uq_shift_swap_active_schedule_any';
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_shift_swap_active_schedule_uniqueness
BEFORE INSERT OR UPDATE OF requester_schedule_id, recipient_schedule_id, status
ON shift_swap_requests
FOR EACH ROW
EXECUTE FUNCTION enforce_shift_swap_active_schedule_uniqueness();

-- ==========================================
-- LEAVE REQUESTS
-- ==========================================

CREATE TYPE leave_request_type_enum AS ENUM (
    'vacation',
    'personal',
    'sick',
    'pregnancy',
    'unpaid',
    'other'
);

CREATE TYPE leave_request_status_enum AS ENUM (
    'pending',
    'approved',
    'rejected',
    'cancelled',
    'expired'
);

CREATE TABLE leave_policies (
    leave_type leave_request_type_enum PRIMARY KEY,
    requires_approval BOOLEAN NOT NULL DEFAULT TRUE,
    deducts_balance BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO leave_policies (leave_type, requires_approval, deducts_balance, is_active) VALUES
    ('vacation', TRUE, TRUE, TRUE),
    ('personal', TRUE, TRUE, TRUE),
    ('sick', FALSE, FALSE, TRUE),
    ('pregnancy', FALSE, FALSE, TRUE),
    ('unpaid', TRUE, FALSE, TRUE),
    ('other', TRUE, FALSE, TRUE);

CREATE TABLE leave_balances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    year INT NOT NULL,
    legal_total_days INT NOT NULL DEFAULT 0,
    extra_total_days INT NOT NULL DEFAULT 0,
    legal_used_days INT NOT NULL DEFAULT 0,
    extra_used_days INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT leave_balances_unique_employee_year UNIQUE (employee_id, year),
    CONSTRAINT leave_balances_non_negative CHECK (
        legal_total_days >= 0
        AND extra_total_days >= 0
        AND legal_used_days >= 0
        AND extra_used_days >= 0
    ),
    CONSTRAINT leave_balances_usage_not_exceed_total CHECK (
        legal_used_days <= legal_total_days
        AND extra_used_days <= extra_total_days
    )
);

CREATE INDEX idx_leave_balances_employee_year ON leave_balances(employee_id, year);

CREATE TABLE leave_balance_adjustments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    leave_balance_id UUID NOT NULL REFERENCES leave_balances(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    year INT NOT NULL,
    legal_days_delta INT NOT NULL DEFAULT 0,
    extra_days_delta INT NOT NULL DEFAULT 0,
    reason TEXT NOT NULL,
    adjusted_by_employee_id UUID NOT NULL REFERENCES employee_profile(id) ON DELETE RESTRICT,
    legal_total_days_before INT NOT NULL,
    extra_total_days_before INT NOT NULL,
    legal_total_days_after INT NOT NULL,
    extra_total_days_after INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT leave_balance_adjustments_non_zero_delta CHECK (
        legal_days_delta <> 0 OR extra_days_delta <> 0
    )
);

CREATE INDEX idx_leave_balance_adjustments_employee_year_created_at
ON leave_balance_adjustments(employee_id, year, created_at DESC);

CREATE OR REPLACE FUNCTION initialize_leave_balance_on_employee_insert()
RETURNS TRIGGER AS $$
DECLARE
    computed_legal_days INT;
    current_year INT;
BEGIN
    current_year := EXTRACT(YEAR FROM CURRENT_DATE)::INT;
    computed_legal_days := GREATEST(0, ROUND(COALESCE(NEW.contract_hours, 0)::numeric / 2.0)::INT);

    INSERT INTO leave_balances (
        employee_id,
        year,
        legal_total_days,
        extra_total_days,
        legal_used_days,
        extra_used_days
    ) VALUES (
        NEW.id,
        current_year,
        computed_legal_days,
        0,
        0,
        0
    )
    ON CONFLICT (employee_id, year) DO NOTHING;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_initialize_leave_balance_on_employee_insert
AFTER INSERT ON employee_profile
FOR EACH ROW
EXECUTE FUNCTION initialize_leave_balance_on_employee_insert();

CREATE TABLE leave_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    created_by_employee_id UUID NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    leave_type leave_request_type_enum NOT NULL,
    status leave_request_status_enum NOT NULL DEFAULT 'pending',
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    reason TEXT NULL,
    decision_note TEXT NULL,
    decided_by_employee_id UUID NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    requested_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    decided_at TIMESTAMPTZ NULL,
    cancelled_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT leave_requests_date_order CHECK (end_date >= start_date)
);

CREATE INDEX idx_leave_requests_employee_id ON leave_requests(employee_id);
CREATE INDEX idx_leave_requests_status ON leave_requests(status);
CREATE INDEX idx_leave_requests_leave_type ON leave_requests(leave_type);
CREATE INDEX idx_leave_requests_requested_at_desc ON leave_requests(requested_at DESC);
CREATE INDEX idx_leave_requests_employee_status ON leave_requests(employee_id, status);

CREATE TYPE calendar_event_kind_enum AS ENUM ('appointment', 'reminder');
CREATE TYPE calendar_event_status_enum AS ENUM ('confirmed', 'cancelled');
-- Work approval status for appointments (hours are counted/billed only after admin approval)
CREATE TYPE calendar_event_work_approval_status_enum AS ENUM ('pending', 'approved', 'rejected');
CREATE TYPE attendee_response_enum AS ENUM ('needs_action', 'accepted', 'declined', 'tentative');
CREATE TYPE reminder_channel_enum AS ENUM ('in_app');

CREATE TABLE calendar_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organizer_employee_id UUID NOT NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    created_by_employee_id UUID NOT NULL REFERENCES employee_profile(id) ON DELETE SET NULL,
    kind calendar_event_kind_enum NOT NULL,
    status calendar_event_status_enum NOT NULL DEFAULT 'confirmed',
    work_approval_status calendar_event_work_approval_status_enum NOT NULL DEFAULT 'pending',
    work_approved_by UUID NULL REFERENCES custom_user(id) ON DELETE SET NULL,
    work_approved_at TIMESTAMPTZ NULL,
    work_rejected_by UUID NULL REFERENCES custom_user(id) ON DELETE SET NULL,
    work_rejected_at TIMESTAMPTZ NULL,
    work_rejection_reason TEXT NULL,
    title TEXT NOT NULL DEFAULT '',
    description TEXT NULL,
    location TEXT NULL,
    color VARCHAR(20) NULL,
    start_at TIMESTAMPTZ NOT NULL,
    end_at TIMESTAMPTZ NOT NULL,
    timezone TEXT NOT NULL DEFAULT 'UTC',
    rrule TEXT NULL,
    recurring_event_id UUID NULL REFERENCES calendar_events(id) ON DELETE CASCADE,
    recurrence_id TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT calendar_events_time_order CHECK (end_at > start_at),
    CONSTRAINT calendar_events_exception_shape CHECK (
        (recurring_event_id IS NULL AND recurrence_id IS NULL)
        OR
        (recurring_event_id IS NOT NULL AND recurrence_id IS NOT NULL)
    ),
    CONSTRAINT calendar_events_exception_no_rrule CHECK (
        recurring_event_id IS NULL OR rrule IS NULL
    )
);

CREATE INDEX idx_calendar_events_organizer_start ON calendar_events(organizer_employee_id, start_at);
CREATE INDEX idx_calendar_events_start ON calendar_events(start_at);
CREATE UNIQUE INDEX uq_calendar_events_exception ON calendar_events(recurring_event_id, recurrence_id)
    WHERE recurring_event_id IS NOT NULL;

-- If appointment times change, previously approved hours should be re-approved.
CREATE OR REPLACE FUNCTION calendar_event_reset_work_approval_on_time_change()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.kind = 'appointment' THEN
        IF (NEW.start_at IS DISTINCT FROM OLD.start_at) OR (NEW.end_at IS DISTINCT FROM OLD.end_at) THEN
            NEW.work_approval_status := 'pending';
            NEW.work_approved_by := NULL;
            NEW.work_approved_at := NULL;
            NEW.work_rejected_by := NULL;
            NEW.work_rejected_at := NULL;
            NEW.work_rejection_reason := NULL;
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_calendar_event_reset_work_approval
BEFORE UPDATE ON calendar_events
FOR EACH ROW
EXECUTE FUNCTION calendar_event_reset_work_approval_on_time_change();

CREATE TABLE calendar_event_attendees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES calendar_events(id) ON DELETE CASCADE,
    employee_id UUID NULL REFERENCES employee_profile(id) ON DELETE CASCADE,
    client_id UUID NULL REFERENCES client_details(id) ON DELETE CASCADE,
    email TEXT NULL,
    response attendee_response_enum NOT NULL DEFAULT 'needs_action',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT calendar_event_attendees_one_target CHECK (
        ((employee_id IS NOT NULL)::INT + (client_id IS NOT NULL)::INT + (email IS NOT NULL)::INT) = 1
    )
);

CREATE UNIQUE INDEX uq_event_attendee_employee ON calendar_event_attendees(event_id, employee_id)
    WHERE employee_id IS NOT NULL;
CREATE UNIQUE INDEX uq_event_attendee_client ON calendar_event_attendees(event_id, client_id)
    WHERE client_id IS NOT NULL;
CREATE UNIQUE INDEX uq_event_attendee_email ON calendar_event_attendees(event_id, email)
    WHERE email IS NOT NULL;
CREATE INDEX idx_event_attendees_employee ON calendar_event_attendees(employee_id);
CREATE INDEX idx_event_attendees_client ON calendar_event_attendees(client_id);

CREATE TABLE calendar_event_reminders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES calendar_events(id) ON DELETE CASCADE,
    channel reminder_channel_enum NOT NULL DEFAULT 'in_app',
    minutes_before INT NULL,
    remind_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT calendar_event_reminders_one_mode CHECK (
        (minutes_before IS NOT NULL) <> (remind_at IS NOT NULL)
    ),
    CONSTRAINT calendar_event_reminders_minutes_positive CHECK (
        minutes_before IS NULL OR minutes_before >= 0
    )
);

CREATE INDEX idx_event_reminders_event ON calendar_event_reminders(event_id);

-- ==========================================
-- INVOICE SOURCES (APPOINTMENTS)
-- ==========================================

-- Links invoice lines to the underlying appointments billed for a client.
-- Billing rule: "start-within" period selection should be enforced by the application query logic.
CREATE TABLE invoice_line_calendar_event (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_line_id UUID NOT NULL REFERENCES invoice_line(id) ON DELETE CASCADE,
    calendar_event_id UUID NOT NULL REFERENCES calendar_events(id) ON DELETE RESTRICT,
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,

    -- Snapshot values used for billing (do not depend on mutable calendar_events rows later)
    start_at TIMESTAMPTZ NOT NULL,
    end_at TIMESTAMPTZ NOT NULL,
    minutes_billed NUMERIC(20,4) NOT NULL,

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (invoice_line_id, calendar_event_id, client_id),
    CHECK (end_at > start_at),
    CHECK (minutes_billed >= 0)
);

CREATE INDEX invoice_line_calendar_event_line_idx ON invoice_line_calendar_event(invoice_line_id);
CREATE INDEX invoice_line_calendar_event_event_idx ON invoice_line_calendar_event(calendar_event_id);
CREATE INDEX invoice_line_calendar_event_client_idx ON invoice_line_calendar_event(client_id);

-- Prevents double-billing the same appointment for the same client (unless voided via credit/cancel).
CREATE TABLE billed_calendar_event (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    calendar_event_id UUID NOT NULL REFERENCES calendar_events(id) ON DELETE RESTRICT,
    client_id UUID NOT NULL REFERENCES client_details(id) ON DELETE CASCADE,
    invoice_id UUID NOT NULL REFERENCES invoice(id) ON DELETE RESTRICT,
    invoice_line_id UUID NOT NULL REFERENCES invoice_line(id) ON DELETE RESTRICT,
    voided_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX billed_calendar_event_active_uq
    ON billed_calendar_event(calendar_event_id, client_id)
    WHERE voided_at IS NULL;

CREATE INDEX billed_calendar_event_invoice_idx ON billed_calendar_event(invoice_id);
CREATE INDEX billed_calendar_event_line_idx ON billed_calendar_event(invoice_line_id);
CREATE INDEX billed_calendar_event_client_idx ON billed_calendar_event(client_id);

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
CREATE OR REPLACE FUNCTION get_client_id_from_goal(goal_id UUID) RETURNS UUID AS $$
    SELECT client_id FROM client_goals WHERE id = goal_id;
$$ LANGUAGE sql STABLE;

CREATE OR REPLACE FUNCTION get_client_id_from_goal_evaluation(eval_id UUID) RETURNS UUID AS $$
    SELECT client_id FROM client_goal_evaluations WHERE id = eval_id;
$$ LANGUAGE sql STABLE;

CREATE OR REPLACE FUNCTION get_client_id_from_contract(cont_id UUID) RETURNS UUID AS $$
    SELECT client_id FROM contract WHERE id = cont_id;
$$ LANGUAGE sql STABLE;

CREATE OR REPLACE FUNCTION get_client_id_from_invoice(inv_id UUID) RETURNS UUID AS $$
    SELECT client_id FROM invoice WHERE id = inv_id;
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
    EXECUTE format('CREATE POLICY coordinator_update ON %I FOR UPDATE USING (is_admin() OR is_assigned_coordinator(%s)) WITH CHECK (is_admin() OR is_assigned_coordinator(%s))', table_name, client_id_col, client_id_col);
    EXECUTE format('CREATE POLICY coordinator_delete ON %I FOR DELETE USING (is_admin() OR is_assigned_coordinator(%s))', table_name, client_id_col);
END;
$$ LANGUAGE plpgsql;

-- Apply RLS to tables with direct client_id or id
SELECT apply_client_rls('client_details', 'id');
SELECT apply_client_rls('progress_report');
SELECT apply_client_rls('incident');
SELECT apply_client_rls('client_documents');
SELECT apply_client_rls('client_status_history');
SELECT apply_client_rls('client_diagnosis');
SELECT apply_client_rls('client_medication_order');
SELECT apply_client_rls('client_emergency_contact');
SELECT apply_client_rls('client_location_transfer');
	SELECT apply_client_rls('contract');
	SELECT apply_client_rls('invoice');
	SELECT apply_client_rls('invoice_run_item');
	SELECT apply_client_rls('invoice_line');
	SELECT apply_client_rls('invoice_line_calendar_event');
	SELECT apply_client_rls('billed_calendar_event');
	SELECT apply_client_rls('invoice_payment_history', 'get_client_id_from_invoice(invoice_id)');
	SELECT apply_client_rls('assignment');
	SELECT apply_client_rls('assigned_employee');
	SELECT apply_client_rls('ai_generated_reports');
	SELECT apply_client_rls('calendar_event_attendees');
SELECT apply_client_rls('appointment_card');
SELECT apply_client_rls('collaboration_agreement');
SELECT apply_client_rls('risk_assessment');
SELECT apply_client_rls('consent_declaration');
SELECT apply_client_rls('youth_care_intake');
SELECT apply_client_rls('data_sharing_statement');
SELECT apply_client_rls('framework_agreement');
SELECT apply_client_rls('client_goals');
SELECT apply_client_rls('client_goal_evaluations');

-- Special cases for nested tables
-- Registration and Intake
SELECT apply_client_rls('registration_form', 'get_client_id_from_registration_form(id)');
SELECT apply_client_rls('intake_forms', 'get_client_id_from_registration_form(registration_form_id)');

-- Medication orders have direct client_id

-- Contract sub-tables
SELECT apply_client_rls('client_agreement', 'get_client_id_from_contract(contract_id)');
SELECT apply_client_rls('provision', 'get_client_id_from_contract(contract_id)');

-- Goal evaluations (new model)
SELECT apply_client_rls('client_goal_evaluation_items', 'get_client_id_from_goal_evaluation(evaluation_id)');

-- Helper to get client_id from intake_form
CREATE OR REPLACE FUNCTION get_client_id_from_intake_form(intake_id UUID) RETURNS UUID AS $$
    SELECT id FROM client_details WHERE intake_form_id = intake_id;
$$ LANGUAGE sql STABLE;

SELECT apply_client_rls('intake_topic_assessments', 'get_client_id_from_intake_form(intake_form_id)');

-- Clean up helper functions if desired, or keep them for future use.
-- DROP FUNCTION apply_client_rls(TEXT, TEXT);
