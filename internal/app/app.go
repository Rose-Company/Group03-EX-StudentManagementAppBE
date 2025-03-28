// internal/app/app.go
package app

import (
	"Group03-EX-StudentManagementAppBE/internal/handlers/admin"
	"Group03-EX-StudentManagementAppBE/internal/handlers/auth"
	"Group03-EX-StudentManagementAppBE/internal/handlers/faculty"
	"Group03-EX-StudentManagementAppBE/internal/handlers/program"
	"Group03-EX-StudentManagementAppBE/internal/handlers/status_transition_rule"
	"Group03-EX-StudentManagementAppBE/internal/handlers/student"
	"Group03-EX-StudentManagementAppBE/internal/handlers/validation"
	"Group03-EX-StudentManagementAppBE/internal/services"

	"github.com/gin-gonic/gin"
)

// Setup initializes and connects all components
func Setup(router *gin.Engine, service *services.Service) {
	api := router.Group("")

	// Initialize handlers and register their routes
	authHandler := auth.NewHandler(service)
	authHandler.RegisterRoutes(api)

	studentHandler := student.NewHandler(service)
	studentHandler.RegisterRoutes(api)

	facultyHandler := faculty.NewHandler(service)
	facultyHandler.RegisterRoutes(api)

	statusTransitionRuleHandler := status_transition_rule.NewHandler(service)
	statusTransitionRuleHandler.RegisterRoutes(api)

	programHandler := program.NewHandler(service)
	programHandler.RegisterRoutes(api)

	validationHandler := validation.NewHandler(service)
	validationHandler.RegisterRoutes(api)


	if service.Admin != nil {
		adminHandler := admin.NewHandler(service)
		adminHandler.RegisterRoutes(api)
	}

}
