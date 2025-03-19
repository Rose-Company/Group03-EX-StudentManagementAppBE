-- Add nationality to students table
ALTER TABLE "PUBLIC".students
ADD COLUMN nationality TEXT;

-- Create a new table for addresses
CREATE TABLE "PUBLIC".student_addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES "PUBLIC".students(id),
    address_type TEXT NOT NULL CHECK (address_type IN ('Permanent', 'Temporary', 'Mailing')),
    street TEXT,
    ward TEXT,
    district TEXT,
    city TEXT,
    country TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    -- Each student can have one address of each type
    UNIQUE(student_id, address_type)
);

-- Create an index for faster lookups
CREATE INDEX idx_student_addresses_student_id ON "PUBLIC".student_addresses(student_id);



CREATE TABLE "PUBLIC".student_identity_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES "PUBLIC".students(id),
    document_type TEXT NOT NULL CHECK (document_type IN ('CMND', 'CCCD', 'Passport')),
    document_number TEXT NOT NULL,
    issue_date DATE NOT NULL,
    issue_place TEXT NOT NULL,
    expiry_date DATE,
    country_of_issue TEXT,
    has_chip BOOLEAN,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(student_id, document_type, document_number)
);

CREATE INDEX idx_student_identity_documents_student_id ON "PUBLIC".student_identity_documents(student_id);

CREATE TABLE "PUBLIC".activity_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id TEXT NOT NULL,
    action TEXT NOT NULL,
    resource TEXT NOT NULL,
    details TEXT,
    user_name TEXT,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_activity_logs_user_id ON "PUBLIC".activity_logs(user_id);
CREATE INDEX idx_activity_logs_resource ON "PUBLIC".activity_logs(resource);
CREATE INDEX idx_activity_logs_action ON "PUBLIC".activity_logs(action);
CREATE INDEX idx_activity_logs_created_at ON "PUBLIC".activity_logs(created_at);

ALTER TABLE "PUBLIC"."users"
ADD COLUMN program_id INTEGER;

-- Create student_program table
CREATE TABLE "PUBLIC"."student_programs" (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Add foreign key constraint (optional, but recommended)
ALTER TABLE "PUBLIC"."users"
ADD CONSTRAINT fk_users_program_id
FOREIGN KEY (program_id) REFERENCES "PUBLIC"."student_programs"(id);