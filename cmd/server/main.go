package main

import (
	"log"
	"monitoring-backend/config"
	"monitoring-backend/handlers"
	"monitoring-backend/middleware"
	"monitoring-backend/models"
	"monitoring-backend/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := models.InitDB(cfg.GenerateDatabaseURL())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := models.Migrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize services
	instanceService := services.NewInstanceService(db)
	serviceService := services.NewServiceService(db)
	dnsService := services.NewDNSService(db)
	notificationService := services.NewNotificationService(db)
	monitoringService := services.NewMonitoringService(db, serviceService, dnsService, notificationService)

	// Start monitoring scheduler
	go monitoringService.StartScheduler()

	// Set database for auth handlers
	handlers.SetDB(db)

	// Initialize handlers
	instanceHandler := handlers.NewInstanceHandler(instanceService)
	serviceHandler := handlers.NewServiceHandler(serviceService)
	dnsHandler := handlers.NewDNSHandler(dnsService)
	metricsHandler := handlers.NewMetricsHandler(instanceService)
	dashboardHandler := handlers.NewDashboardHandler(instanceService, serviceService, dnsService)

	// Setup router
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Auth middleware for protected routes
	authMiddleware := middleware.AuthMiddleware()

	// Public routes
	api := r.Group("/api/v1")
	{
		// Metrics endpoint for instances to send data
		api.POST("/metrics", metricsHandler.ReceiveMetrics)

		// Auth routes
		api.POST("/auth/login", handlers.Login)
		api.POST("/auth/register", handlers.Register)
	}

	// Protected routes
	protected := api.Group("/")
	protected.Use(authMiddleware)
	{
		// Dashboard
		protected.GET("/dashboard", dashboardHandler.GetDashboard)

		// Instances
		instances := protected.Group("/instances")
		{
			instances.GET("", instanceHandler.GetInstances)
			instances.POST("", instanceHandler.CreateInstance)
			instances.GET("/:id", instanceHandler.GetInstance)
			instances.PUT("/:id", instanceHandler.UpdateInstance)
			instances.DELETE("/:id", instanceHandler.DeleteInstance)
			instances.GET("/:id/metrics", instanceHandler.GetInstanceMetrics)
		}

		// Services
		services := protected.Group("/services")
		{
			services.GET("", serviceHandler.GetServices)
			services.POST("", serviceHandler.CreateService)
			services.GET("/:id", serviceHandler.GetService)
			services.PUT("/:id", serviceHandler.UpdateService)
			services.DELETE("/:id", serviceHandler.DeleteService)
			services.GET("/:id/checks", serviceHandler.GetServiceChecks)
		}

		// DNS
		dns := protected.Group("/dns")
		{
			dns.GET("", dnsHandler.GetDNSRecords)
			dns.POST("", dnsHandler.CreateDNSRecord)
			dns.GET("/:id", dnsHandler.GetDNSRecord)
			dns.PUT("/:id", dnsHandler.UpdateDNSRecord)
			dns.DELETE("/:id", dnsHandler.DeleteDNSRecord)
			dns.GET("/:id/checks", dnsHandler.GetDNSChecks)
		}
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run("0.0.0.0:" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
