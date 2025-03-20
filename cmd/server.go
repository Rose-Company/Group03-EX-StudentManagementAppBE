package cmd

import (
	"Group03-EX-StudentManagementAppBE/config"
	"Group03-EX-StudentManagementAppBE/internal/app"
	"Group03-EX-StudentManagementAppBE/internal/repositories/faculty"
	"Group03-EX-StudentManagementAppBE/internal/repositories/program"
	"Group03-EX-StudentManagementAppBE/internal/repositories/student"
	student_addresses "Group03-EX-StudentManagementAppBE/internal/repositories/student_addresses"
	student_identity_documents "Group03-EX-StudentManagementAppBE/internal/repositories/student_documents"
	"Group03-EX-StudentManagementAppBE/internal/repositories/student_status"
	"Group03-EX-StudentManagementAppBE/internal/repositories/user"
	"Group03-EX-StudentManagementAppBE/internal/services"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// C:\Users\Admin\Desktop\project\Group03-EX-StudentManagementAppBE\internal\repositories\student_addresses.go
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the API server",
	Long:  "Start the Student Management API server with specified configuration",
	Run: func(cmd *cobra.Command, args []string) {
		start, _ := cmd.Flags().GetBool("start")
		port, _ := cmd.Flags().GetString("port")
		mode, _ := cmd.Flags().GetString("mode")

		if start {
			runServer(port, mode)
		} else {
			cmd.Help()
		}
	},
}

func runServer(port, mode string) {
	// Set Gin mode
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Initialize database connection
	db := initDatabase(cfg)

	// Initialize repositories
	repositories := initRepositories(db)

	// Initialize services
	service := initServices(repositories)

	// Initialize router
	router := initRouter()

	// Setup application routes and handlers
	app.Setup(router, service)

	// Start the server
	serverAddr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("Starting server at %s in %s mode\n", serverAddr, mode)
	router.Run(serverAddr)
}

func initDatabase(cfg *config.Config) *gorm.DB {
	// Create the PostgreSQL DSN
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
		cfg.Postgres.Host,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Name,
		cfg.Postgres.Port,
	)

	// Configure GORM logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		},
	)

	// Initialize the database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to the database: %v", err))
	}

	return db
}

type repositoriesContainer struct {
	userRepo            user.Repository
	studentRepo         student.Repository
	facultyRepo         faculty.Repository
	studentStatusRepo   student_status.Repository
	StudentAddressRepo  student_addresses.Repository
	StudentDocumentRepo student_identity_documents.Repository
	programRepo         program.Repository
}

func initRepositories(db *gorm.DB) *repositoriesContainer {
	return &repositoriesContainer{
		userRepo:            user.NewRepository(db),
		studentRepo:         student.NewRepository(db),
		facultyRepo:         faculty.NewRepository(db),
		studentStatusRepo:   student_status.NewRepository(db),
		StudentAddressRepo:  student_addresses.NewRepository(db),
		StudentDocumentRepo: student_identity_documents.NewRepository(db),
		programRepo:         program.NewRepository(db),
	}
}

func initServices(repos *repositoriesContainer) *services.Service {
	return services.NewService(
		repos.userRepo,
		repos.studentRepo,
		repos.facultyRepo,
		repos.studentStatusRepo,
		repos.StudentAddressRepo,
		repos.StudentDocumentRepo,
		repos.programRepo,
	)
}

func initRouter() *gin.Engine {
	// Initialize the Gin router
	router := gin.New()

	// Add recovery and logger middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Apply CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	return router
}
