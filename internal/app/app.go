// internal/app/app.go
package app

import (
	"Group03-EX-StudentManagementAppBE/internal/handlers/auth"
	"Group03-EX-StudentManagementAppBE/internal/services"

	"github.com/gin-gonic/gin"
)

// Setup initializes and connects all components
func Setup(router *gin.Engine, service *services.Service) {
	// Create API group
	authen := router.Group("")
	{
		authHandler := auth.NewHandler(service)
		authHandler.RegisterRoutes(authen)
	}

	student := router.Group("")
	{
		studentHandler := student.NewStudentHandler(service)
		studentHandler.RegisterRoutes(student)
	}

	// Register feature-specific routes

	//When you add more features, register their routes here:

}
