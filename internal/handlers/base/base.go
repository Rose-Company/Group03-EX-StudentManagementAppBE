// internal/handlers/base.go
package base

import (
	"Group03-EX-StudentManagementAppBE/internal/services"
)

// Handler provides common handler functionality
type Handler struct {
	Service *services.Service
}

// NewHandler creates a new base handler
func NewHandler(service *services.Service) Handler {
	return Handler{
		Service: service,
	}
}
