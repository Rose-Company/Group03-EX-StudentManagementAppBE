package models

type LoginRequest struct {
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
}

type LoginResponse struct {
	ID    string `json:"id"`
	Code  int    `json:"code"`
	Token string `json:"token"`
}
