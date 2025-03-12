// internal/handlers/base.go
package handlers

import (
	"Group03-EX-StudentManagementAppBE/internal/services"
)

// BaseHandler provides common handler functionality
type BaseHandler struct {
	Service *services.Service
}

// NewBaseHandler creates a new base handler
func NewBaseHandler(service *services.Service) BaseHandler {
	return BaseHandler{
		Service: service,
	}
}
