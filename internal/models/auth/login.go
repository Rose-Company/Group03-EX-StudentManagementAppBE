package models

type LoginRequest struct {
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
}

type LoginResponse struct {
	Code  int    `json:"code"`
	Token string `json:"token"`
}
