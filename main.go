package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/cmd"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/controller"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/database"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/middleware"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/repository"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/routes"
	"github.com/Shabrinashsf/go-gin-gorm-boilerplate/service"
	"github.com/common-nighthawk/go-figure"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Println("cannot load .env:", err)
	} else {
		log.Println(".env loaded successfully")
	}

	// Initialized database
	log.Println("Setting up database connection...")
	db := database.SetUpDatabaseConnection()
	defer database.CloseDatabaseConnection(db)
	log.Println("Database connection established.")

	// Handle CLI command
	if len(os.Args) > 1 {
		log.Println("Running commands...")
		cmd.Command(db)
		return
	}

	// Dependency injection
	var (
		// Initilization package
		//jwtService service.JWTService = service.NewJWTService()

		// Repository
		transactionRepository repository.TransactionRepository = repository.NewTransactionRepository(db)

		// Service
		transactionService service.TransactionService = service.NewTransactionService(transactionRepository, db)

		// Controller
		transactionController controller.TransactionController = controller.NewTransactionController(transactionService)
	)

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Route Not Found",
		})
	})

	server.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong pong",
		})
	})

	// Initialize routes
	routes.Transaction(server, transactionController)

	server.Static("/assets", "./assets")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	var serve string
	if os.Getenv("APP_ENV") == "localhost" {
		serve = "127.0.0.1:" + port
	} else {
		serve = ":" + port
	}

	myFigure := figure.NewColorFigure("Backend Boilerplate", "", "blue", true)
	myFigure.Print()

	fmt.Printf("Starting server on %s\n", serve)
	if err := server.Run(serve); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
