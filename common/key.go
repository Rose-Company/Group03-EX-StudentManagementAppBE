package common

import "time"

const (
	PREFIX_MAIN_POSTGRES = "MAIN_POSTGRES"
)

const ( //must NOT edit this
	ENV_GIN_DEBUG  = "GIN_DEBUG"
	ENV_RABBIT_URI = "RABBIT"
)

const (
	ENVJWTSecretKey = "JWT__SECRET_KEY"
)

var (
	DATETIME_WITH_TIMEZONE = time.RFC3339
)

const (
	USER_JWT_KEY = "USER_JWT_PROFILE"
	UserId       = "user_id"
)
