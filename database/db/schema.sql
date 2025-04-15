-- ENUM TYPES
CREATE TYPE AP_STATUS AS ENUM ('None', 'Booked', 'Resolved', 'Cancelled');
CREATE TYPE MEDICAL_STATUS AS ENUM ('Ongoing', 'Recovered', 'None');
CREATE TYPE U_ROLE AS ENUM ('patient', 'doctor');

-- BASE TABLES
CREATE TABLE IF NOT EXISTS contact_details (
    id SERIAL PRIMARY KEY,
    primary_number TEXT,
    secondary_number TEXT,
    email VARCHAR(254) -- Email max RFC-compliant length
);

CREATE TABLE IF NOT EXISTS address (
    id SERIAL PRIMARY KEY,
    street VARCHAR(250),
    city VARCHAR(200),
    state VARCHAR(200),
    country VARCHAR(200),
    zip VARCHAR(20) -- To support ZIPs with dashes or leading 0s
);

-- MAIN USER TABLE
CREATE TABLE IF NOT EXISTS "user" (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    location VARCHAR(200),
    user_role U_ROLE NOT NULL,
    contact_id INTEGER,
    emergency_contact INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_contact_id FOREIGN KEY (contact_id) REFERENCES contact_details(id),
    CONSTRAINT fk_emergency_contact FOREIGN KEY (emergency_contact) REFERENCES contact_details(id)
);

-- AUTH TABLE
CREATE TABLE IF NOT EXISTS auth (
    user_id INTEGER PRIMARY KEY,
    user_name VARCHAR(100) NOT NULL,
    password VARCHAR(250) NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE
);

-- DOCTOR & PATIENT EXTENSION OF USER
CREATE TABLE IF NOT EXISTS doctor (
    id INTEGER PRIMARY KEY,
    specialty VARCHAR(200) NOT NULL,
    experience INTEGER NOT NULL, -- in years
    note TEXT,
    appointment_status AP_STATUS,
    is_available BOOLEAN DEFAULT false,
    appointment_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_doctor_user FOREIGN KEY (id) REFERENCES "user"(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS patient (
    id INTEGER PRIMARY KEY,
    appointment_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    medical_history_id INTEGER,
    CONSTRAINT fk_patient_user FOREIGN KEY (id) REFERENCES "user"(id) ON DELETE CASCADE,
    CONSTRAINT fk_appointment_id FOREIGN KEY (appointment_id) REFERENCES appointment(id),
    CONSTRAINT fk_medical_history_id FOREIGN KEY (medical_history_id) REFERENCES medical_history(id)
);

-- APPOINTMENT TABLE
CREATE TABLE IF NOT EXISTS appointment (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    reason TEXT NOT NULL,
    booked TIMESTAMP NOT NULL,
    appointment_status AP_STATUS,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_appointment FOREIGN KEY (user_id) REFERENCES "user"(id)
);

-- PAYMENTS
CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    patient_id INTEGER NOT NULL,
    doctor_id INTEGER NOT NULL,
    amount REAL NOT NULL CHECK (amount > 0),
    CONSTRAINT fk_payment_patient FOREIGN KEY (patient_id) REFERENCES patient(id),
    CONSTRAINT fk_payment_doctor FOREIGN KEY (doctor_id) REFERENCES doctor(id)
);

-- MEDICAL RELATED TABLES
CREATE TABLE IF NOT EXISTS medication_history (
    id SERIAL PRIMARY KEY,
    description TEXT,
    dosage TEXT, -- changed from TIMESTAMP to TEXT
    frequency INTEGER DEFAULT 1,
    start_date TIMESTAMP,
    end_date TIMESTAMP
);

CREATE TABLE IF NOT EXISTS allergy (
    id SERIAL PRIMARY KEY,
    description TEXT,
    reaction TEXT,
    severity TEXT -- Changed column name from 'Severity' for consistency
);

CREATE TABLE IF NOT EXISTS surgery (
    id SERIAL PRIMARY KEY,
    description TEXT,
    date_done TIMESTAMP,
    hospital TEXT
);

CREATE TABLE IF NOT EXISTS medical_history (
    id SERIAL PRIMARY KEY,
    patient_id INTEGER,
    description TEXT,
    diagnosed_date TIMESTAMP,
    status MEDICAL_STATUS,
    medication INTEGER,
    allergies INTEGER,
    surgeries INTEGER,
    CONSTRAINT fk_mh_patient FOREIGN KEY (patient_id) REFERENCES patient(id),
    CONSTRAINT fk_mh_medication FOREIGN KEY (medication) REFERENCES medication_history(id),
    CONSTRAINT fk_mh_allergies FOREIGN KEY (allergies) REFERENCES allergy(id),
    CONSTRAINT fk_mh_surgeries FOREIGN KEY (surgeries) REFERENCES surgery(id)
);
