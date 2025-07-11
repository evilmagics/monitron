package handlers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"monitron-server/models"
)
// CreateService
// @Summary Create a new service
// @Description Create a new monitoring service
// @Tags Services
// @Accept json
// @Produce json
// @Param service body models.Service true "Service object to be created"
// @Success 201 {object} models.Service
// @Failure 400 {object} map[string]string "error": "Cannot parse JSON"
// @Failure 500 {object} map[string]string "error": "Could not create service"
// @Security ApiKeyAuth
// @Router /services [post]
func CreateService(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		service := new(models.Service)
		if err := c.BodyParser(service); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		service.ID = uuid.New()
		service.CreatedAt = time.Now()
		service.UpdatedAt = time.Now()

		if result := db.Create(&service); result.Error != nil {
			log.Printf("Error creating service: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create service"})
		}

		return c.Status(fiber.StatusCreated).JSON(service)
	}
}

// GetServices
// @Summary Get all services
// @Description Retrieve a list of all monitoring services
// @Tags Services
// @Produce json
// @Success 200 {array} models.Service
// @Failure 500 {object} map[string]string "error": "Could not retrieve services"
// @Security ApiKeyAuth
// @Router /services [get]
func GetServices(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		services := []models.Service{}
		if result := db.Find(&services); result.Error != nil {
			log.Printf("Error fetching services: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve services"})
		}

		return c.JSON(services)
	}
}

// GetService
// @Summary Get service by ID
// @Description Retrieve a single monitoring service by its ID
// @Tags Services
// @Produce json
// @Param id path string true "Service ID"
// @Success 200 {object} models.Service
// @Failure 400 {object} map[string]string "error": "Invalid service ID"
// @Failure 404 {object} map[string]string "error": "Service not found"
// @Failure 500 {object} map[string]string "error": "Could not retrieve service"
// @Security ApiKeyAuth
// @Router /services/{id} [get]
func GetService(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuidID, err := uuid.Parse(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid service ID"})
		}

		service := models.Service{}
		if result := db.First(&service, "id = ?", uuidID); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Service not found"})
			}
			log.Printf("Error fetching service: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve service"})
		}

		return c.JSON(service)
	}
}

// UpdateService
// @Summary Update an existing service
// @Description Update details of an existing monitoring service by its ID
// @Tags Services
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Param service body models.Service true "Service object with updated fields"
// @Success 200 {object} models.Service
// @Failure 400 {object} map[string]string "error": "Invalid service ID" or "Cannot parse JSON"
// @Failure 404 {object} map[string]string "error": "Service not found"
// @Failure 500 {object} map[string]string "error": "Could not update service"
// @Security ApiKeyAuth
// @Router /services/{id} [put]
func UpdateService(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuidID, err := uuid.Parse(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid service ID"})
		}

		service := new(models.Service)
		if err := c.BodyParser(service); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		var existingService models.Service
		if result := db.First(&existingService, "id = ?", uuidID); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Service not found"})
			}
			log.Printf("Error finding service for update: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update service"})
		}

		if result := db.Model(&existingService).Updates(service); result.Error != nil {
			log.Printf("Error updating service: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update service"})
		}

		return c.JSON(existingService)
}

// DeleteService
// @Summary Delete a service
// @Description Delete a monitoring service by its ID
// @Tags Services
// @Produce json
// @Param id path string true "Service ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "error": "Invalid service ID"
// @Failure 404 {object} map[string]string "error": "Service not found"
// @Failure 500 {object} map[string]string "error": "Could not delete service"
// @Security ApiKeyAuth
// @Router /services/{id} [delete]
func DeleteService(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuidID, err := uuid.Parse(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid service ID"})
		}

		if result := db.Delete(&models.Service{}, "id = ?", uuidID); result.Error != nil {
			log.Printf("Error deleting service: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete service"})
		}

		if result.RowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Service not found"})
		}

		return c.Status(fiber.StatusNoContent).SendString("")
	}
}



// ServiceHealthCheck performs health checks for all services
func ServiceHealthCheck(db *gorm.DB) {
	log.Println("Running scheduled service health check...")
	services := []models.Service{}
	
		if result := db.Find(&services); result.Error != nil {
			log.Printf("Error fetching services for health check: %v", result.Error)
			return
		}

	for _, service := range services {
		log.Printf("Checking service: %s (Type: %s)", service.Name, service.APIType)
		// Implement actual health check logic based on service.APIType
		// For now, just simulate a successful check
		log.Printf("Service %s health check successful.", service.Name)
	}
	log.Println("Service health check completed.")
}


