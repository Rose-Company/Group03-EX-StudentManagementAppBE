package main

import (
	"Group03-EX-StudentManagementAppBE/config"
	"Group03-EX-StudentManagementAppBE/internal/app"
	"Group03-EX-StudentManagementAppBE/internal/repositories/faculty"
	"Group03-EX-StudentManagementAppBE/internal/repositories/student"
	"Group03-EX-StudentManagementAppBE/internal/repositories/user"
	"Group03-EX-StudentManagementAppBE/internal/services"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Load configuration from environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Create the PostgreSQL DSN
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
		cfg.Postgres.Host,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Name,
		cfg.Postgres.Port,
	)

	// Configure GORM logger for SQL query logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level (Info shows all queries)
			IgnoreRecordNotFoundError: false,       // Show record not found errors
			Colorful:                  true,        // Enable color
		},
	)

	// Initialize the database connection with logger
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic("Failed to connect to the database")
	}
	// Initialize repositories
	userRepo := user.NewRepository(db)
	studentRepo := student.NewRepository(db)
	facultyRepo := faculty.NewRepository(db)

	// Initialize services
	service := services.NewService(userRepo, studentRepo, facultyRepo)

	// Initialize the Gin router
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Setup the application (connect handlers, services, etc.)
	app.Setup(router, service)

	// Start the server
	router.Run("0.0.0.0:8080")
}
