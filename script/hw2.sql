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
