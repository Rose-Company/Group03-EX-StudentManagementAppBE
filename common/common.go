package common

// Define New Table Name here
const (
	POSTGRES_TABLE_NAME_USERS                      = "PUBLIC.users"
	POSTGRES_TABLE_NAME_STUDENTS                   = "PUBLIC.students"
	POSTGRES_TABLE_NAME_FACULTY                    = "PUBLIC.faculties"
	POSTGRES_TABLE_NAME_STUDENTS_STATUSES          = "PUBLIC.student_statuses"
	POSTGRES_TABLE_NAME_STUDENT_ADDRESSES          = "PUBLIC.student_addresses"
	POSTGRES_TABLE_NAME_STUDENT_IDENTITY_DOCUMENTS = "PUBLIC.student_identity_documents"
	POSTGRES_TABLE_NAME_STUDENT_PROGRAMS           = "PUBLIC.student_programs"
	POSTGRES_TABLE_NAME_FILES                      = "PUBLIC.files"
	POSTGRES_TABLE_NAME_STATUS_TRANSITION_RULES      = "PUBLIC.status_transition_rules"
	POSTGRES_TABLE_NAME_VALIDTION_SETTINGS         = "PUBLIC.validation_settings"
	POSTGRES_TABLE_NAME_VALIDATION_RULES           = "PUBLIC.validation_rules"
)

const (
	ROLE_END_USER        = "END_USER"
	ROLE_END_USER_UUID   = ""
	ROLE_ADMIN           = "ADMIN"
	ROLE_ADMIN_USER_UUID = ""

	USER_PROVIDER_GOOGLE = ""
)
