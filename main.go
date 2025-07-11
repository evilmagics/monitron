package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron/v3"

	"monitron-server/config"
	"monitron-server/database"
	"monitron-server/messaging"
	"monitron-server/router"
)

// @title Monitron API
// @version 1.0
// @description This is the API documentation for the Monitron application.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	cfg := config.LoadConfig()

	// Initialize database
	db := database.InitDB(cfg)
	defer database.CloseDB(db)

	// Initialize RabbitMQ
	messaging.InitRabbitMQ(cfg)
	defer messaging.CloseRabbitMQ()

	// Setup and start RabbitMQ consumers in a goroutine
	go messaging.SetupConsumers()

	app := fiber.New()

	// Setup API routes
	router.SetupRoutes(app, db)
	// Initialize and start cron scheduler
	c := cron.New()
	c.Start()
	defer c.Stop()

	// Start server in a goroutine
	go func() {
		log.Printf("Server is running on %s:%d", cfg.App.Host, cfg.App.Port)
		if err := app.Listen(cfg.App.Host + ":" + strconv.Itoa(cfg.App.Port)); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}
	log.Println("Server gracefully stopped.")
}
