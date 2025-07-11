package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"monitron-server/handlers"
	"monitron-server/middleware"
)

func SetupRoutes(app *fiber.App, db *sqlx.DB) {
	api := app.Group("/api/v1")

	// Instance Management Routes
	instances := api.Group("/instances")
	instances.Post("/", handlers.CreateInstance(db))
	instances.Get("/", handlers.GetInstances(db))
	instances.Get("/:id", handlers.GetInstance(db))
	instances.Put("/:id", handlers.UpdateInstance(db))
	instances.Delete("/:id", handlers.DeleteInstance(db))

	// Service Management Routes
	services := api.Group("/services")
	services.Post("/", handlers.CreateService(db))
	services.Get("/", handlers.GetServices(db))
	services.Get("/:id", handlers.GetService(db))
	services.Put("/:id", handlers.UpdateService(db))
	services.Delete("/:id", handlers.DeleteService(db))

	// Domain & SSL Management Routes
	domainSSL := api.Group("/domain-ssl")
	domainSSL.Post("/", handlers.CreateDomainSSL(db))
	domainSSL.Get("/", handlers.GetDomainSSLs(db))
	domainSSL.Get("/:id", handlers.GetDomainSSL(db))
	domainSSL.Put("/:id", handlers.UpdateDomainSSL(db))
	domainSSL.Delete("/:id", handlers.DeleteDomainSSL(db))

	// Authentication Routes
	auth := api.Group("/auth")
	auth.Post("/register", handlers.RegisterUser(db))
	auth.Post("/login", handlers.LoginUser(db))

	// User Management Routes (Admin Only)
	users := api.Group("/users", middleware.JWTAuth(), middleware.AdminAuth())
	users.Get("/", handlers.GetUsers(db))
	users.Get("/:id", handlers.GetUser(db))
	users.Put("/:id", handlers.UpdateUser(db))
	users.Delete("/:id", handlers.DeleteUser(db))

	// Authenticated User Routes
	userAuth := api.Group("/user", middleware.JWTAuth())
	userAuth.Put("/change-password", handlers.ChangePassword(db))

	// Password Reset Routes (No Auth Required)
	passwordReset := api.Group("/password")
	passwordReset.Post("/forgot", handlers.ForgotPassword(db))
	passwordReset.Post("/reset", handlers.ResetPassword(db))

	// Report Routes (Admin Only for now, or specific user reports)
	reports := api.Group("/reports", middleware.JWTAuth())
	reports.Post("/", handlers.CreateReport(db))
	reports.Get("/", handlers.GetReports(db))
	reports.Get("/:id", handlers.GetReport(db))

	// Log Routes (Admin Only)
	logs := api.Group("/logs", middleware.JWTAuth(), middleware.AdminAuth())
	logs.Post("/", handlers.CreateLogEntry(db))
	logs.Get("/", handlers.GetLogEntries(db))
	logs.Get("/:id", handlers.GetLogEntry(db))

	// Operational Page Routes
	opPages := api.Group("/operational-pages")
	opPages.Post("/", middleware.JWTAuth(), handlers.CreateOperationalPage(db))
	opPages.Get("/", handlers.GetOperationalPages(db))
	opPages.Get("/:idOrSlug", handlers.GetOperationalPage(db))
	opPages.Put("/:idOrSlug", middleware.JWTAuth(), handlers.UpdateOperationalPage(db))
	opPages.Delete("/:idOrSlug", middleware.JWTAuth(), handlers.DeleteOperationalPage(db))

	// Operational Page Components Routes
	opPageComponents := api.Group("/operational-pages/:pageID/components")
	opPageComponents.Post("/", middleware.JWTAuth(), handlers.AddComponentToOperationalPage(db))
	opPageComponents.Get("/", handlers.GetComponentsForOperationalPage(db))
	opPageComponents.Delete("/:componentID", middleware.JWTAuth(), handlers.RemoveComponentFromOperationalPage(db))

	// GraphQL Route
	api.Post("/graphql", handlers.GraphQLHandler(db))
}

