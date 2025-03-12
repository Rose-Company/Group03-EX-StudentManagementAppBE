package main

import (
	"Group03-EX-StudentManagementAppBE/config"
	"Group03-EX-StudentManagementAppBE/internal/handlers"
	"Group03-EX-StudentManagementAppBE/internal/repositories"
	"Group03-EX-StudentManagementAppBE/internal/services"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration from environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Create the PostgreSQL DSN (Data Source Name) using the loaded config
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
		cfg.Postgres.Host,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Name,
		cfg.Postgres.Port,
	)

	// Initialize the database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database")
	}

	// Initialize repositories
	studentRepo := repositories.NewStudentRepository(db)

	// Initialize services
	studentService := services.NewStudentService(studentRepo)

	// Initialize the Gin router and register routes
	router := gin.Default()

	// Create API group
	api := router.Group("")

	// Initialize handlers and register routes
	studentHandler := handlers.NewStudentHandler(studentService)
	studentHandler.RegisterRoutes(api)

	// Start the server
	router.Run("0.0.0.0:" + fmt.Sprintf("%d", 8080))
}
