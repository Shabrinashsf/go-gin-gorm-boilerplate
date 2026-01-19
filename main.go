package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/cmd"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/controller"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/database"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/middleware"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/repository"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/routes"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/service"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/utils/logger"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/utils/mailer"
	"github.com/common-nighthawk/go-figure"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Server struct {
	// Configuration
	port string
	env  string

	// HTTP Server
	ginEngine  *gin.Engine
	httpServer *http.Server

	// Context Management
	rootCTX    context.Context
	cancelFunc context.CancelFunc

	// Database
	db *gorm.DB

	// Dependency injection
	jwtService service.JWTService
	mailer     mailer.Mailer

	// Repository
	transactionRepo repository.TransactionRepository
	userRepo        repository.UserRepository

	// Service
	transactionService service.TransactionService
	userService        service.UserService

	// Controller
	transactionController controller.TransactionController
	userController        controller.UserController
}

func NewServer(db *gorm.DB) *Server {
	jwtService := service.NewJWTService()
	mailer := mailer.NewMailer()

	// Repository
	transactionRepo := repository.NewTransactionRepository(db)
	userRepo := repository.NewUserController(db)

	// Service
	transactionService := service.NewTransactionService(transactionRepo, db)
	userService := service.NewUserService(userRepo, jwtService, mailer, db)

	// Controller
	transactionController := controller.NewTransactionController(transactionService)
	userController := controller.NewUserController(userService)

	// Get current mode
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8888"
	}

	mode := os.Getenv("APP_ENV")
	if mode == "" {
		mode = "localhost"
	}

	return &Server{
		port:                  port,
		env:                   mode,
		db:                    db,
		transactionRepo:       transactionRepo,
		transactionService:    transactionService,
		transactionController: transactionController,
		userRepo:              userRepo,
		userService:           userService,
		userController:        userController,
		jwtService:            jwtService,
	}
}

func main() {
	// Initialize environment variables
	if err := godotenv.Load(".env"); err != nil {
		logger.Errorf("cannot load .env: %v", err)
	} else {
		logger.Infof(".env loaded successfully")
	}

	// Initialized database
	logger.Infof("Setting up database connection...")
	db := database.SetUpDatabaseConnection()
	defer database.CloseDatabaseConnection(db)
	logger.Infof("Database connection established.")

	// Handle CLI command
	if len(os.Args) > 1 {
		logger.Infof("Running commands...")
		cmd.Command(db)
		return
	}

	// Create server instance
	server := NewServer(db)

	// Start server
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	sig := <-quit
	logger.Infof("Received signal: %v", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		logger.Errorf("Shutdown error: %v", err)
	}

	logger.Infof("Application exited")
}

func (s *Server) Start() error {
	// Create root context
	s.rootCTX, s.cancelFunc = context.WithCancel(context.Background())
	logger.Infof("Services initialized")
	logger.Infof("Setting up server...")

	// Setup Gin
	s.ginEngine = gin.Default()
	s.ginEngine.Use(middleware.CORSMiddleware())

	// No route handler
	s.ginEngine.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Route Not Found",
		})
	})

	// Health check
	s.ginEngine.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "aku sehat, kangen?",
		})
	})

	// Register routes
	routes.Transaction(s.ginEngine, s.transactionController)
	routes.User(s.ginEngine, s.userController, s.jwtService)

	s.ginEngine.Static("/assets", "./assets")

	// Create HTTP server
	var addr string
	if s.env == "localhost" {
		addr = "127.0.0.1:" + s.port
	} else {
		addr = ":" + s.port
	}

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.ginEngine,
	}

	// Start HTTP server in goroutine
	go func() {
		myFigure := figure.NewColorFigure("Backend Boilerplate", "", "blue", true)
		myFigure.Print()
		fmt.Printf("Starting server on %s\n", addr)

		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Server error: %v", err)
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	logger.Infof("Starting Graceful Shutdown")

	// Step 1: Cancel root context
	s.cancelFunc()

	// Step 2: Shutdown HTTP server
	logger.Infof("Shutting down HTTP server...")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		logger.Errorf("HTTP server shutdown error: %v", err)
		return err
	}
	logger.Infof("HTTP Server stopped")

	// Step 3: Close database connections
	logger.Infof("Closing database connections...")
	sqlDB, err := s.db.DB()
	if err == nil {
		sqlDB.Close()
	}
	logger.Infof("Database connections closed")
	logger.Infof("Graceful Shutdown completed")
	return nil
}
