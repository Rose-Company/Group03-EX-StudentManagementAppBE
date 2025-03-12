-- Complete database schema script for the entire student management system

-- Drop tables if they exist (in reverse order of dependencies)
DROP TABLE IF EXISTS "PUBLIC".students CASCADE;
DROP TABLE IF EXISTS "PUBLIC".users CASCADE;
DROP TABLE IF EXISTS "PUBLIC".student_statuses CASCADE;
DROP TABLE IF EXISTS "PUBLIC".faculties CASCADE;
DROP TABLE IF EXISTS "PUBLIC".roles CASCADE;

-- Create sequences for auto-incrementing IDs
CREATE SEQUENCE IF NOT EXISTS faculty_id_seq START 1;
CREATE SEQUENCE IF NOT EXISTS student_status_id_seq START 1;
CREATE SEQUENCE IF NOT EXISTS role_id_seq START 1;

-- Create the roles table
CREATE TABLE "PUBLIC".roles (
    id INTEGER PRIMARY KEY DEFAULT nextval('role_id_seq'),
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create the faculties table
CREATE TABLE "PUBLIC".faculties (
    id INTEGER PRIMARY KEY DEFAULT nextval('faculty_id_seq'),
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create the student_statuses table
CREATE TABLE "PUBLIC".student_statuses (
    id INTEGER PRIMARY KEY DEFAULT nextval('student_status_id_seq'),
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create the users table
CREATE TABLE "PUBLIC".users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role_id INTEGER REFERENCES "PUBLIC".roles(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create the students table
CREATE TABLE "PUBLIC".students (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_code INTEGER UNIQUE NOT NULL,
    fullname TEXT NOT NULL,
    date_of_birth DATE NOT NULL,
    gender TEXT CHECK (gender IN ('Male', 'Female', 'Other')),
    faculty_id INTEGER REFERENCES "PUBLIC".faculties(id),
    batch TEXT NOT NULL,
    program TEXT NOT NULL,
    address TEXT,
    email TEXT UNIQUE,
    phone TEXT,
    status_id INTEGER REFERENCES "PUBLIC".student_statuses(id),
    user_id UUID REFERENCES "PUBLIC".users(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX idx_students_student_code ON "PUBLIC".students(student_code);
CREATE INDEX idx_students_faculty_id ON "PUBLIC".students(faculty_id);
CREATE INDEX idx_students_status_id ON "PUBLIC".students(status_id);
CREATE INDEX idx_students_user_id ON "PUBLIC".students(user_id);
CREATE INDEX idx_users_role_id ON "PUBLIC".users(role_id);
CREATE INDEX idx_users_email ON "PUBLIC".users(email);

-- Initialize roles
INSERT INTO "PUBLIC".roles (name) VALUES 
('Admin'),
('Teacher'),
('Student'),
('Staff');

-- Initialize faculties
INSERT INTO "PUBLIC".faculties (name) VALUES 
('School of Law'),
('School of Business English'),
('School of Japanese Studies'),
('School of French Studies');

-- Initialize student statuses
INSERT INTO "PUBLIC".student_statuses (name) VALUES 
('Enrolled'),
('Graduated'),
('Withdrawn'),
('On Leave');

-- Initialize some users
INSERT INTO "PUBLIC".users (email, password, role_id) VALUES
('admin@university.edu', '$2a$10$hKDVYxLefVHV/vtuPhWD3OigtRyOykRLDdUAp80Z1crSoS1lFqaFS', 1), -- Admin user
('teacher@university.edu', '$2a$10$hKDVYxLefVHV/vtuPhWD3OigtRyOykRLDdUAp80Z1crSoS1lFqaFS', 2), -- Teacher user
('student@university.edu', '$2a$10$hKDVYxLefVHV/vtuPhWD3OigtRyOykRLDdUAp80Z1crSoS1lFqaFS', 3); -- Student user

-- Initialize students (50 records)
INSERT INTO "PUBLIC".students 
(id, student_code, fullname, date_of_birth, gender, faculty_id, batch, program, address, email, phone, status_id, user_id)
VALUES
-- School of Law (faculty_id = 1, faculty code = 10)
(gen_random_uuid(), 22100001, 'John Smith', '2000-05-15', 'Male', 1, '2022', 'Bachelor of Law', '123 Main St, Apt 4B', 'john.smith@email.com', '555-123-4567', 1, NULL),
(gen_random_uuid(), 22100002, 'Emily Johnson', '2001-07-22', 'Female', 1, '2022', 'Bachelor of Law', '456 Oak Ave', 'emily.j@email.com', '555-234-5678', 1, NULL),
(gen_random_uuid(), 22100003, 'Michael Brown', '1999-11-30', 'Male', 1, '2022', 'Bachelor of Law', '789 Pine St', 'mbrown@email.com', '555-345-6789', 1, NULL),
(gen_random_uuid(), 21100001, 'Sarah Lee', '2000-02-14', 'Female', 1, '2021', 'Bachelor of Law', '101 Maple Dr', 'slee@email.com', '555-456-7890', 1, NULL),
(gen_random_uuid(), 21100002, 'David Wilson', '1999-08-19', 'Male', 1, '2021', 'Bachelor of Law', '202 Elm St', 'dwilson@email.com', '555-567-8901', 1, NULL),
(gen_random_uuid(), 21100003, 'Jessica Taylor', '2000-04-05', 'Female', 1, '2021', 'Bachelor of Law', '303 Cedar Ln', 'jtaylor@email.com', '555-678-9012', 1, NULL),
(gen_random_uuid(), 20100001, 'Daniel Martinez', '1998-12-11', 'Male', 1, '2020', 'Bachelor of Law', '404 Birch Rd', 'dmartinez@email.com', '555-789-0123', 1, NULL),
(gen_random_uuid(), 20100002, 'Amanda Garcia', '1999-01-25', 'Female', 1, '2020', 'Bachelor of Law', '505 Walnut Ave', 'agarcia@email.com', '555-890-1234', 1, NULL),
(gen_random_uuid(), 20100003, 'James Rodriguez', '1998-06-30', 'Male', 1, '2020', 'Bachelor of Law', '606 Cherry St', 'jrodriguez@email.com', '555-901-2345', 1, NULL),
(gen_random_uuid(), 19100001, 'Jennifer Lopez', '1997-09-17', 'Female', 1, '2019', 'Bachelor of Law', '707 Pine Ave', 'jlopez@email.com', '555-012-3456', 2, NULL),
(gen_random_uuid(), 19100002, 'Robert Smith', '1998-03-22', 'Male', 1, '2019', 'Bachelor of Law', '808 Oak Dr', 'rsmith@email.com', '555-123-4567', 2, NULL),
(gen_random_uuid(), 19100003, 'Lisa Johnson', '1997-07-15', 'Female', 1, '2019', 'Bachelor of Law', '909 Maple St', 'ljohnson@email.com', '555-234-5678', 3, NULL),

-- School of Business English (faculty_id = 2, faculty code = 20)
(gen_random_uuid(), 22200001, 'Kevin Williams', '2001-10-12', 'Male', 2, '2022', 'Bachelor of Business English', '111 Spruce Ave', 'kwilliams@email.com', '555-345-6789', 1, NULL),
(gen_random_uuid(), 22200002, 'Michelle Brown', '2002-01-30', 'Female', 2, '2022', 'Bachelor of Business English', '222 Fir St', 'mbrown2@email.com', '555-456-7890', 1, NULL),
(gen_random_uuid(), 22200003, 'Thomas Jones', '2000-04-05', 'Male', 2, '2022', 'Bachelor of Business English', '333 Willow Ln', 'tjones@email.com', '555-567-8901', 1, NULL),
(gen_random_uuid(), 21200001, 'Rachel Miller', '2001-06-18', 'Female', 2, '2021', 'Bachelor of Business English', '444 Aspen Rd', 'rmiller@email.com', '555-678-9012', 1, NULL),
(gen_random_uuid(), 21200002, 'Christopher Davis', '1999-11-22', 'Male', 2, '2021', 'Bachelor of Business English', '555 Birch Ave', 'cdavis@email.com', '555-789-0123', 1, NULL),
(gen_random_uuid(), 21200003, 'Laura Garcia', '2000-03-14', 'Female', 2, '2021', 'Bachelor of Business English', '666 Cedar St', 'lgarcia@email.com', '555-890-1234', 1, NULL),
(gen_random_uuid(), 20200001, 'Andrew Wilson', '1999-05-27', 'Male', 2, '2020', 'Bachelor of Business English', '777 Elm Dr', 'awilson@email.com', '555-901-2345', 1, NULL),
(gen_random_uuid(), 20200002, 'Nicole Martinez', '1998-08-09', 'Female', 2, '2020', 'Bachelor of Business English', '888 Maple Ln', 'nmartinez@email.com', '555-012-3456', 1, NULL),
(gen_random_uuid(), 20200003, 'Steven Taylor', '1999-02-11', 'Male', 2, '2020', 'Bachelor of Business English', '999 Oak Ave', 'staylor@email.com', '555-123-4567', 4, NULL),
(gen_random_uuid(), 19200001, 'Elizabeth Smith', '1997-04-30', 'Female', 2, '2019', 'Bachelor of Business English', '123 Pine St', 'esmith@email.com', '555-234-5678', 2, NULL),
(gen_random_uuid(), 19200002, 'Brandon Johnson', '1998-09-15', 'Male', 2, '2019', 'Bachelor of Business English', '234 Cedar Ave', 'bjohnson@email.com', '555-345-6789', 2, NULL),
(gen_random_uuid(), 19200003, 'Samantha Brown', '1997-12-03', 'Female', 2, '2019', 'Bachelor of Business English', '345 Elm St', 'sbrown@email.com', '555-456-7890', 3, NULL),

-- School of Japanese Studies (faculty_id = 3, faculty code = 30)
(gen_random_uuid(), 22300001, 'Ryan Chen', '2002-02-28', 'Male', 3, '2022', 'Bachelor of Japanese Studies', '456 Cherry Dr', 'rchen@email.com', '555-567-8901', 1, NULL),
(gen_random_uuid(), 22300002, 'Sophia Lee', '2001-05-17', 'Female', 3, '2022', 'Bachelor of Japanese Studies', '567 Spruce Ln', 'slee2@email.com', '555-678-9012', 1, NULL),
(gen_random_uuid(), 22300003, 'William Kim', '2000-08-21', 'Male', 3, '2022', 'Bachelor of Japanese Studies', '678 Willow Ave', 'wkim@email.com', '555-789-0123', 1, NULL),
(gen_random_uuid(), 21300001, 'Julia Wong', '2001-10-06', 'Female', 3, '2021', 'Bachelor of Japanese Studies', '789 Maple Rd', 'jwong@email.com', '555-890-1234', 1, NULL),
(gen_random_uuid(), 21300002, 'Daniel Park', '1999-12-19', 'Male', 3, '2021', 'Bachelor of Japanese Studies', '890 Pine St', 'dpark@email.com', '555-901-2345', 1, NULL),
(gen_random_uuid(), 21300003, 'Emma Nguyen', '2000-03-30', 'Female', 3, '2021', 'Bachelor of Japanese Studies', '901 Oak Ln', 'enguyen@email.com', '555-012-3456', 1, NULL),
(gen_random_uuid(), 20300001, 'Alex Tanaka', '1999-07-13', 'Male', 3, '2020', 'Bachelor of Japanese Studies', '012 Cedar Dr', 'atanaka@email.com', '555-123-4567', 1, NULL),
(gen_random_uuid(), 20300002, 'Olivia Chang', '1998-11-27', 'Female', 3, '2020', 'Bachelor of Japanese Studies', '123 Elm Ave', 'ochang@email.com', '555-234-5678', 4, NULL),
(gen_random_uuid(), 20300003, 'Matthew Suzuki', '1999-01-14', 'Male', 3, '2020', 'Bachelor of Japanese Studies', '234 Birch St', 'msuzuki@email.com', '555-345-6789', 1, NULL),
(gen_random_uuid(), 19300001, 'Hannah Kim', '1997-04-09', 'Female', 3, '2019', 'Bachelor of Japanese Studies', '345 Aspen Dr', 'hkim@email.com', '555-456-7890', 2, NULL),
(gen_random_uuid(), 19300002, 'Justin Lee', '1998-06-22', 'Male', 3, '2019', 'Bachelor of Japanese Studies', '456 Spruce Ave', 'jlee@email.com', '555-567-8901', 2, NULL),
(gen_random_uuid(), 19300003, 'Megan Chen', '1997-09-05', 'Female', 3, '2019', 'Bachelor of Japanese Studies', '567 Pine Ln', 'mchen@email.com', '555-678-9012', 3, NULL),

-- School of French Studies (faculty_id = 4, faculty code = 40)
(gen_random_uuid(), 22400001, 'Nathan Martin', '2001-12-11', 'Male', 4, '2022', 'Bachelor of French Studies', '678 Oak St', 'nmartin@email.com', '555-789-0123', 1, NULL),
(gen_random_uuid(), 22400002, 'Victoria Bernard', '2002-03-24', 'Female', 4, '2022', 'Bachelor of French Studies', '789 Elm Rd', 'vbernard@email.com', '555-890-1234', 1, NULL),
(gen_random_uuid(), 22400003, 'Adam Dubois', '2000-06-15', 'Male', 4, '2022', 'Bachelor of French Studies', '890 Cedar Ln', 'adubois@email.com', '555-901-2345', 1, NULL),
(gen_random_uuid(), 21400001, 'Natalie Moreau', '2001-09-28', 'Female', 4, '2021', 'Bachelor of French Studies', '901 Maple Ave', 'nmoreau@email.com', '555-012-3456', 1, NULL),
(gen_random_uuid(), 21400002, 'Eric Lambert', '1999-11-07', 'Male', 4, '2021', 'Bachelor of French Studies', '012 Pine Dr', 'elambert@email.com', '555-123-4567', 1, NULL),
(gen_random_uuid(), 21400003, 'Sophie Rousseau', '2000-02-19', 'Female', 4, '2021', 'Bachelor of French Studies', '123 Birch St', 'srousseau@email.com', '555-234-5678', 1, NULL),
(gen_random_uuid(), 20400001, 'Paul Lefebvre', '1999-05-03', 'Male', 4, '2020', 'Bachelor of French Studies', '234 Cherry Ave', 'plefebvre@email.com', '555-345-6789', 1, NULL),
(gen_random_uuid(), 20400002, 'Claire Dupont', '1998-07-20', 'Female', 4, '2020', 'Bachelor of French Studies', '345 Willow Dr', 'cdupont@email.com', '555-456-7890', 1, NULL),
(gen_random_uuid(), 20400003, 'Lucas Mercier', '1999-10-16', 'Male', 4, '2020', 'Bachelor of French Studies', '456 Aspen Ln', 'lmercier@email.com', '555-567-8901', 4, NULL),
(gen_random_uuid(), 19400001, 'Isabelle Fournier', '1997-03-01', 'Female', 4, '2019', 'Bachelor of French Studies', '567 Spruce St', 'ifournier@email.com', '555-678-9012', 2, NULL),
(gen_random_uuid(), 19400002, 'Antoine Girard', '1998-06-14', 'Male', 4, '2019', 'Bachelor of French Studies', '678 Fir Ave', 'agirard@email.com', '555-789-0123', 2, NULL),
(gen_random_uuid(), 19400003, 'Camille Laurent', '1997-08-27', 'Female', 4, '2019', 'Bachelor of French Studies', '789 Cedar Rd', 'claurent@email.com', '555-890-1234', 3, NULL);

-- Link a student to a user account (example)
UPDATE "PUBLIC".students 
SET user_id = (SELECT id FROM "PUBLIC".users WHERE email = 'student@university.edu')
WHERE student_code = 22100001;

-- Verify the data was successfully inserted
SELECT 'Roles count: ' || COUNT(*) as count FROM "PUBLIC".roles
UNION ALL
SELECT 'Faculties count: ' || COUNT(*) as count FROM "PUBLIC".faculties
UNION ALL
SELECT 'Student statuses count: ' || COUNT(*) as count FROM "PUBLIC".student_statuses
UNION ALL
SELECT 'Users count: ' || COUNT(*) as count FROM "PUBLIC".users
UNION ALL
SELECT 'Students count: ' || COUNT(*) as count FROM "PUBLIC".students;