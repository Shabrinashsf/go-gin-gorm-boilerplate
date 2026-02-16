package config

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/controller"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/repository"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/routes"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/service"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/utils/logger"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/utils/mailer"
	"github.com/common-nighthawk/go-figure"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RestConfig struct {
	server     *gin.Engine
	httpServer *http.Server
	port       string
	env        string
}

func NewRestConfig(db *gorm.DB) *RestConfig {
	app := gin.Default()
	server := NewRouter(app)

	// =========== (SERVICES) ===========
	jwtService := service.NewJWTService()
	mailerService := mailer.NewMailer()

	// =========== (REPOSITORY) ===========
	transactionRepo := repository.NewTransactionRepository(db)
	userRepo := repository.NewUserController(db)

	// =========== (SERVICE) ===========
	transactionService := service.NewTransactionService(transactionRepo, db)
	userService := service.NewUserService(userRepo, jwtService, mailerService, db)

	// =========== (CONTROLLER) ===========
	transactionController := controller.NewTransactionController(transactionService)
	userController := controller.NewUserController(userService)

	// =========== (ROUTES) ===========
	routes.Transaction(server, transactionController)
	routes.User(server, userController, jwtService)

	// Get configuration
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8888"
	}

	mode := os.Getenv("APP_ENV")
	if mode == "" {
		mode = "localhost"
	}

	return &RestConfig{
		server:     server,
		httpServer: nil,
		port:       port,
		env:        mode,
	}
}

func (rc *RestConfig) Start() error {
	var addr string
	if rc.env == "localhost" {
		addr = "127.0.0.1:" + rc.port
	} else {
		addr = ":" + rc.port
	}

	rc.httpServer = &http.Server{
		Addr:    addr,
		Handler: rc.server,
	}

	go func() {
		myFigure := figure.NewColorFigure("Backend Boilerplate", "", "blue", true)
		myFigure.Print()
		fmt.Printf("Starting server on %s\n", addr)

		if err := rc.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Server error: %v", err)
		}
	}()

	return nil
}

func (rc *RestConfig) GetEngine() *gin.Engine {
	return rc.server
}

func (rc *RestConfig) GetPort() string {
	return rc.port
}

func (rc *RestConfig) GetEnv() string {
	return rc.env
}

func (rc *RestConfig) Shutdown(ctx context.Context) error {
	logger.Infof("Starting graceful shutdown...")

	if rc.httpServer != nil {
		logger.Infof("Shutting down HTTP server...")
		if err := rc.httpServer.Shutdown(ctx); err != nil {
			logger.Errorf("HTTP server shutdown error: %v", err)
			return err
		}
		logger.Infof("HTTP server stopped")
	}

	logger.Infof("Graceful shutdown completed")
	return nil
}
