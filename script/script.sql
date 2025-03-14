-- Complete database schema script for the entire student management system

-- Drop tables if they exist (in reverse order of dependencies)
DROP TABLE IF EXISTS "PUBLIC".students CASCADE;
DROP TABLE IF EXISTS "PUBLIC".users CASCADE;
DROP TABLE IF EXISTS "PUBLIC".student_statuses CASCADE;
DROP TABLE IF EXISTS "PUBLIC".faculties CASCADE;
DROP TABLE IF EXISTS "PUBLIC".roles CASCADE;

-- Drop sequences if they exist
DROP SEQUENCE IF EXISTS faculty_id_seq CASCADE;
DROP SEQUENCE IF EXISTS student_status_id_seq CASCADE;

-- Create sequences for auto-incrementing IDs
CREATE SEQUENCE IF NOT EXISTS faculty_id_seq START 1;
CREATE SEQUENCE IF NOT EXISTS student_status_id_seq START 1;

-- Create the roles table
CREATE TABLE "PUBLIC".roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
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
    role_id UUID REFERENCES "PUBLIC".roles(id),
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
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX idx_students_student_code ON "PUBLIC".students(student_code);
CREATE INDEX idx_students_faculty_id ON "PUBLIC".students(faculty_id);
CREATE INDEX idx_students_status_id ON "PUBLIC".students(status_id);
CREATE INDEX idx_users_role_id ON "PUBLIC".users(role_id);
CREATE INDEX idx_users_email ON "PUBLIC".users(email);

-- Initialize roles
INSERT INTO "PUBLIC".roles (id, name) VALUES 
('a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6', 'Admin'),
('b2c3d4e5-f6a7-48b9-c0d1-e2f3a4b5c6d7', 'Teacher'),
('c3d4e5f6-a7b8-49c0-d1e2-f3a4b5c6d7e8', 'Student'),
('d4e5f6a7-b8c9-40d1-e2f3-a4b5c6d7e8f9', 'Staff');

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

-- Initialize users
INSERT INTO "PUBLIC"."users" ("id", "email", "password", "role", "created_at", "updated_at") VALUES
('8a0f7a89-cac7-48b3-8f6e-cdb1786fa953', 'test@gmail.com', '$2y$10$REycf2QaxAVgJC6ZyIQqMeh51ADnPBZKRXiKyDR31GPdH0nIW.1r2', 'a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6', '2025-03-12 04:12:05.918871+00', '2025-03-12 04:12:05.918871+00');

-- Initialize sample students
INSERT INTO "PUBLIC"."students" ("id", "student_code", "fullname", "date_of_birth", "gender", "faculty_id", "batch", "program", "address", "email", "phone", "status_id", "created_at", "updated_at") VALUES
('103c01c6-6a3e-43a8-a7a8-f8b1903c8c41', 22100002, 'Emily Johnson', '2001-07-22', 'Female', 1, '2022', 'Bachelor of Law', '456 Oak Ave', 'emily.j@email.com', '555-234-5678', 1, '2025-03-12 04:29:10.292482+00', '2025-03-12 04:29:10.292482+00'),
('3baf1ca3-33e4-414c-a00b-af18ae890137', 22100003, 'Michael Brown', '1999-11-30', 'Male', 1, '2022', 'Bachelor of Law', '789 Pine St', 'mbrown@email.com', '555-345-6789', 1, '2025-03-12 04:29:10.292482+00', '2025-03-12 04:29:10.292482+00'),
('ef58c28d-06d5-469e-a022-4b9cf95c94e1', 21100001, 'Sarah Lee', '2000-02-14', 'Female', 1, '2021', 'Bachelor of Law', '101 Maple Dr', 'slee@email.com', '555-456-7890', 1, '2025-03-12 04:29:10.292482+00', '2025-03-12 04:29:10.292482+00'),
('3622d20d-69bf-46d2-8af1-a68ee4a56c3d', 21200003, 'Laura Garcia', '2000-03-14', 'Female', 2, '2021', 'Bachelor of Business English', '666 Cedar St', 'lgarcia@email.com', '555-890-1234', 1, '2025-03-12 04:29:10.292482+00', '2025-03-14 08:05:42.157608+00'),
('3aae81a0-f11c-47f4-ab0a-499eccbf7b98', 20400003, 'Lucas Mercier', '1999-10-16', 'Male', 4, '2020', 'Bachelor of French Studies', '456 Aspen Ln', 'lmercier@email.com', '555-567-8901', 4, '2025-03-12 04:29:10.292482+00', '2025-03-14 08:27:10.73392+00'),
('9e80fa1b-848b-480a-8b86-775b34b279a6', 1234567999, 'Tran Thai Nhat', '1998-11-27', 'Female', 3, '2025', 'Bachelor of Japanese Studies', '123 Elm Ave', 'ochang@email.com', '555-234-5678-9', 4, '2025-03-12 04:29:10.292482+00', '2025-03-14 06:17:15.991985+00'),
('9513ac4a-fa9b-4e1c-a9cc-e1b2d004d9d9', 204000012, 'Paul Lefebvre', '1999-05-03', 'Male', 4, '2020', 'Bachelor of French Studies', '234 Cherry Ave', 'plefebvre@email.com', '555-345-6789', 1, '2025-03-12 04:29:10.292482+00', '2025-03-14 06:17:42.814721+00'),
('9ad9fb5d-6d14-4d8d-ac6f-2fc4410fd1f2', 22120222, 'Trần Thái Nhật', '1111-11-11', 'Male', 2, '2020', 'Bachelor of Business English', 'Thôn Tân Lập, xã Cam Phước Tây, huyện Cam Lâm, tỉnh Khánh Hòa', 'Sadmin@example.com', '0584606179', 1, '2025-03-14 07:11:57.612517+00', '2025-03-14 07:11:57.612517+00'),
('02d958fd-f7f1-444f-97a8-c1c8ed891f6f', 22200001, 'Kevin Williams', '2001-10-12', 'Male', 2, '2022', 'Bachelor of Business English', '111 Spruce Ave', 'kwilliams@email.com', '555-345-6789', 1, '2025-03-12 04:29:10.292482+00', '2025-03-12 04:29:10.292482+00'),
('e74b5e48-f4f8-4f43-bf7b-c8ad18675751', 22200002, 'Michelle Brown', '2002-01-30', 'Female', 2, '2022', 'Bachelor of Business English', '222 Fir St', 'mbrown2@email.com', '555-456-7890', 1, '2025-03-12 04:29:10.292482+00', '2025-03-12 04:29:10.292482+00'),
('dfd222fd-1b24-4fc6-941a-cd3dbdeb6f93', 22200003, 'Thomas Jones', '2000-04-05', 'Male', 2, '2022', 'Bachelor of Business English', '333 Willow Ln', 'tjones@email.com', '555-567-8901', 1, '2025-03-12 04:29:10.292482+00', '2025-03-12 04:29:10.292482+00'),
('f4ea68a2-5b79-414b-9563-195b450f3f38', 21100003, 'Jessica Taylor', '2000-04-05', 'Female', 1, '2021', 'Bachelor of Law', '303 Cedar Lnnnnnnnnn', 'jtaylor@email.com', '555-678-9012', 1, '2025-03-12 04:29:10.292482+00', '2025-03-14 08:04:46.846317+00'),
('821330aa-a925-4e47-88e6-7f1f7be35ffd', 20400002, 'Claire Dupont', '1998-07-20', 'Female', 4, '2021', 'Bachelor of French Studies', '345 Willow Dr', 'cdupont@email.com', '555-456-7890', 1, '2025-03-12 04:29:10.292482+00', '2025-03-14 08:30:11.409165+00'),
('cb9a3acf-08b8-4f27-8d20-06f1920d0e7a', 21200001, 'Rachel Miller', '2001-06-18', 'Female', 2, '2021', 'Bachelor of Business English', '444 Aspen Rd', 'rmiller@email.com', '555-678-9012', 1, '2025-03-12 04:29:10.292482+00', '2025-03-14 08:08:43.278838+00'),
('c18beb73-7760-4c82-bc2f-7dccc96bbe1b', 21100002, 'David Wilson', '1999-08-19', 'Male', 1, '2021', 'Bachelor of Law', '202 Elm St', 'dwilson@email.com', '555-567-8901', 1, '2025-03-12 04:29:10.292482+00', '2025-03-14 08:10:25.295804+00'),
('62ac94ee-fc24-4cd2-bc51-99781f2bbec3', 22300001, 'Ryan Chen', '2002-02-28', 'Male', 3, '2022', 'Bachelor of Japanese Studies', '456 Cherry Dr', 'rchen@email.com', '555-567-8901', 1, '2025-03-12 04:29:10.292482+00', '2025-03-12 04:29:10.292482+00'),
('69ffb754-a12a-4d4a-9dd8-5a5e71f2b21c', 22300002, 'Sophia Lee', '2001-05-17', 'Female', 3, '2022', 'Bachelor of Japanese Studies', '567 Spruce Ln', 'slee2@email.com', '555-678-9012', 1, '2025-03-12 04:29:10.292482+00', '2025-03-12 04:29:10.292482+00'),
('8c7706a1-94d0-4092-b207-e050ae829146', 22300003, 'William Kim', '2000-08-21', 'Male', 3, '2022', 'Bachelor of Japanese Studies', '678 Willow Ave', 'wkim@email.com', '555-789-0123', 1, '2025-03-12 04:29:10.292482+00', '2025-03-12 04:29:10.292482+00'),
('596f7470-97c4-4d05-a0b7-31f56097a74f', 21300001, 'Julia Wong', '2001-10-06', 'Female', 3, '2021', 'Bachelor of Japanese Studies', '789 Maple Rd', 'jwong@email.com', '555-890-1234', 1, '2025-03-12 04:29:10.292482+00', '2025-03-12 04:29:10.292482+00'),
('70fba6fa-a3d9-480e-9443-457d8f7fa607', 21300002, 'Daniel Park', '1999-12-19', 'Male', 3, '2021', 'Bachelor of Japanese Studies', '890 Pine St', 'dpark@email.com', '555-901-2345', 1, '2025-03-12 04:29:10.292482+00', '2025-03-12 04:29:10.292482+00'),
('65606311-69b0-41cc-8689-f3601a5d5590', 21200002, 'Christopher Davis', '1999-11-22', 'Male', 2, '2021', 'Bachelor of Business English', '555 Birch Ave', 'cdavis@email.com', '555-789-0123', 1, '2025-03-12 04:29:10.292482+00', '2025-03-14 08:10:29.727017+00'),
('d018983f-08e1-4a6f-ae11-6184b2fa24d9', 22400001, 'Nathan Martin', '2001-12-11', 'Male', 4, '2022', 'Bachelor of French Studies', '678 Oak St', 'nmartin@email.com', '555-789-0123', 1, '2025-03-12 04:29:10.292482+00', '2025-03-12 04:29:10.292482+00'),
('f7e1bcc2-3e48-4bef-b85c-5beac26ed95b', 22400002, 'Victoria Bernard', '2002-03-24', 'Female', 4, '2022', 'Bachelor of French Studies', '789 Elm Rd', 'vbernard@email.com', '555-890-1234', 1, '2025-03-12 04:29:10.292482+00', '2025-03-12 04:29:10.292482+00'),
('dfe07721-8a4c-48fe-a27e-f16358d34546', 21400002, 'Eric Lambert', '1999-11-07', 'Male', 4, '2021', 'Bachelor of French Studies', '012 Pine Dr', 'elambert@email.com', '555-123-4567', 1, '2025-03-12 04:29:10.292482+00', '2025-03-12 04:29:10.292482+00'),
('6abce18a-8aa5-4a6c-ba96-174bed8e8a15', 21400003, 'Sophie Rousseau', '2000-02-19', 'Female', 4, '2021', 'Bachelor of French Studies', '123 Birch St', 'srousseau@email.com', '555-234-5678', 1, '2025-03-12 04:29:10.292482+00', '2025-03-12 04:29:10.292482+00'),
('defd7813-b30e-4c86-aaaa-a7c4250a95a7', 22100001, 'John Smith nr', '2000-05-15', 'Male', 1, '2022', 'Bachelor of Law', '123 Main St, Apt 4B', 'john.smith@email.com', '555-123-4567', 1, '2025-03-12 04:29:10.292482+00', '2025-03-12 04:29:10.292482+00'),
('8b9e87b1-78c7-42c8-91f0-d142922a14a4', 190999990, 'Lisa Johnson', '1997-07-15', 'Other', 1, '2019', 'Bachelor of Laww', '909 Maple St', 'ljohnson@email.com', '555-234-5678', 3, '2025-03-12 04:29:10.292482+00', '2025-03-13 12:08:19.681587+00'),
('baa7a7f9-532a-4d26-babc-eb4bc87d5d2a', 22120263, 'Elizabeth Smith', '1997-04-30', 'Female', 2, '2019', 'Bachelor of Business English', '123 Pine St', 'Test@email.com', '555-234-5678', 2, '2025-03-12 04:29:10.292482+00', '2025-03-13 12:37:45.433902+00'),
('0bd888bc-f0ae-47e3-b850-7d0ddd21cd48', 191000111, 'Robert Smith', '1998-03-22', 'Male', 2, '2019', 'Bachelor of Law', '8/18 Dinh Bo Linh', 'nguyenthanhphat@gmail.com', '555-123-4567', 2, '2025-03-12 04:29:10.292482+00', '2025-03-13 12:25:24.897662+00'),
('11557270-417d-4150-aae9-139b5621b022', 22120195, 'Brandon Johnson', '1998-09-15', 'Male', 2, '2019', 'Bachelor of Business English', '234 Cedar Ave', 'bjohnson@email.com', '555-345-6789', 2, '2025-03-12 04:29:10.292482+00', '2025-03-13 12:52:52.727614+00'),
('94eb78e6-f597-4c98-acaa-9d2e91a8f9be', 22120256, 'Samantha Brown', '1997-12-03', 'Female', 2, '2019', 'Bachelor of Business English', '345 Elm St', 'sbrown@email.com', '555-456-7890', 3, '2025-03-12 04:29:10.292482+00', '2025-03-13 12:55:25.106002+00'),
('b0791121-bb5d-4d3b-9d77-a236b0503193', 21220263, 'Hannah Kim', '1997-04-09', 'Female', 3, '2019', 'Bachelor of Japanese Studies', '345 Aspen Dr', 'hkim@email.com', '555-456-7890', 2, '2025-03-12 04:29:10.292482+00', '2025-03-13 13:00:13.560558+00'),
('b801d38a-2ec4-4be0-97b5-060af1157f7e', 221204444, 'An An An', '2004-11-04', 'Female', 2, '2020', 'Bachelor of Business English', 'tramtnasdn', 'thainhat04@1111.com', '123456', 2, '2025-03-14 04:33:41.655667+00', '2025-03-14 06:20:56.587072+00'),
('9b2ee70b-ec7c-4a6c-82f8-503430628500', 22120216, 'Antoine Girard', '1998-06-14', 'Male', 4, '2019', 'Bachelor of French Studies', '678 Fir Ave', 'agirard@email.com', '555-789-0123', 2, '2025-03-12 04:29:10.292482+00', '2025-03-13 13:58:20.376193+00'),
('8d459ce9-c448-4a68-a655-f730426b9147', 22120255, 'Trần Thái Nhật', '1111-11-11', 'Female', 2, '2020', 'Bachelor of Business English', 'Thôn Tân Lập, xã Cam Phước Tây, huyện Cam Lâm, tỉnh Khánh Hòa', 'sadmin@example.com', '0584606179', 2, '2025-03-14 04:19:48.15617+00', '2025-03-14 04:19:48.15617+00'),
('4708f151-f042-48c4-b660-3296b4598ab1', 22120244, 'Trần Thái Nhật', '1111-11-11', 'Female', 2, '2020', 'Bachelor of Business English', 'Thôn Tân Lập, xã Cam Phước Tây, huyện Cam Lâm, tỉnh Khánh Hòa', 'tranthainhat2k4@gmail.com', '0584606179', 2, '2025-03-14 04:23:45.671318+00', '2025-03-14 04:23:45.671318+00'),
('51afc9c3-781b-4f88-8051-50a7b0c8b705', 21212122, 'Trần Thái Nhật', '1111-05-04', 'Female', 3, '2020', 'Bachelor of Business English', 'Thôn Tân Lập, xã Cam Phước Tây, huyện Cam Lâm, tỉnh Khánh Hòa', 'test@gmail.com', '0584606179', 2, '2025-03-14 04:42:11.651462+00', '2025-03-14 04:42:11.651462+00');

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