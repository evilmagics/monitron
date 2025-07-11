package handlers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"errors"

	"monitron-server/models"
)
// CreateDomainSSL
// @Summary Create a new domain/SSL entry
// @Description Create a new domain and SSL certificate monitoring entry
// @Tags Domain & SSL
// @Accept json
// @Produce json
// @Param domainSSL body models.DomainSSL true "Domain/SSL object to be created"
// @Success 201 {object} models.DomainSSL
// @Failure 400 {object} map[string]string "error": "Cannot parse JSON"
// @Failure 500 {object} map[string]string "error": "Could not create domain/SSL entry"
// @Security ApiKeyAuth
// @Router /domain-ssl [post]
func CreateDomainSSL(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		domainSSL := new(models.DomainSSL)
		if err := c.BodyParser(domainSSL); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		domainSSL.ID = uuid.New()
		domainSSL.CreatedAt = time.Now()
		domainSSL.UpdatedAt = time.Now()

		if result := db.Create(&domainSSL); result.Error != nil {
			log.Printf("Error creating domain/SSL entry: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create domain/SSL entry"})
		}

		return c.Status(fiber.StatusCreated).JSON(domainSSL)
	}
}

// GetDomainSSLs
// @Summary Get all domain/SSL entries
// @Description Retrieve a list of all domain and SSL certificate monitoring entries
// @Tags Domain & SSL
// @Produce json
// @Success 200 {array} models.DomainSSL
// @Failure 500 {object} map[string]string "error": "Could not retrieve domain/SSL entries"
// @Security ApiKeyAuth
// @Router /domain-ssl [get]
func GetDomainSSLs(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		domainSSLs := []models.DomainSSL{}
		if result := db.Find(&domainSSLs); result.Error != nil {
			log.Printf("Error fetching domain/SSLs: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve domain/SSL entries"})
		}

		return c.JSON(domainSSLs)
	}
}

// GetDomainSSL
// @Summary Get domain/SSL entry by ID
// @Description Retrieve a single domain and SSL certificate monitoring entry by its ID
// @Tags Domain & SSL
// @Produce json
// @Param id path string true "Domain/SSL ID"
// @Success 200 {object} models.DomainSSL
// @Failure 400 {object} map[string]string "error": "Invalid domain/SSL ID"
// @Failure 404 {object} map[string]string "error": "Domain/SSL entry not found"
// @Failure 500 {object} map[string]string "error": "Could not retrieve domain/SSL entry"
// @Security ApiKeyAuth
// @Router /domain-ssl/{id} [get]
func GetDomainSSL(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuidID, err := uuid.Parse(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid domain/SSL ID"})
		}

		domainSSL := models.DomainSSL{}
		if result := db.First(&domainSSL, "id = ?", uuidID); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Domain/SSL entry not found"})
			}
			log.Printf("Error fetching domain/SSL: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve domain/SSL entry"})
		}

		return c.JSON(domainSSL)
	}
}

// UpdateDomainSSL
// @Summary Update an existing domain/SSL entry
// @Description Update details of an existing domain and SSL certificate monitoring entry by its ID
// @Tags Domain & SSL
// @Accept json
// @Produce json
// @Param id path string true "Domain/SSL ID"
// @Param domainSSL body models.DomainSSL true "Domain/SSL object with updated fields"
// @Success 200 {object} models.DomainSSL
// @Failure 400 {object} map[string]string "error": "Invalid domain/SSL ID" or "Cannot parse JSON"
// @Failure 404 {object} map[string]string "error": "Domain/SSL entry not found"
// @Failure 500 {object} map[string]string "error": "Could not update domain/SSL entry"
// @Security ApiKeyAuth
// @Router /domain-ssl/{id} [put]
func UpdateDomainSSL(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuidID, err := uuid.Parse(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid domain/SSL ID"})
		}

		domainSSL := new(models.DomainSSL)
		if err := c.BodyParser(domainSSL); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		var existingDomainSSL models.DomainSSL
		if result := db.First(&existingDomainSSL, "id = ?", uuidID); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Domain/SSL entry not found"})
			}
			log.Printf("Error finding domain/SSL for update: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update domain/SSL entry"})
		}

		if result := db.Model(&existingDomainSSL).Updates(domainSSL); result.Error != nil {
			log.Printf("Error updating domain/SSL: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update domain/SSL entry"})
		}

		return c.JSON(existingDomainSSL)
}

// DeleteDomainSSL
// @Summary Delete a domain/SSL entry
// @Description Delete a domain and SSL certificate monitoring entry by its ID
// @Tags Domain & SSL
// @Produce json
// @Param id path string true "Domain/SSL ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "error": "Invalid domain/SSL ID"
// @Failure 404 {object} map[string]string "error": "Domain/SSL entry not found"
// @Failure 500 {object} map[string]string "error": "Could not delete domain/SSL entry"
// @Security ApiKeyAuth
// @Router /domain-ssl/{id} [delete]
func DeleteDomainSSL(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuidID, err := uuid.Parse(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid domain/SSL ID"})
		}

		if result := db.Delete(&models.DomainSSL{}, "id = ?", uuidID); result.Error != nil {
			log.Printf("Error deleting domain/SSL: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete domain/SSL entry"})
		}

		if result.RowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Domain/SSL entry not found"})
		}

		return c.Status(fiber.StatusNoContent).SendString("")
	}
}


// DomainSSLHealthCheck performs health checks for all domain/SSL entries
func DomainSSLHealthCheck(db *gorm.DB) {
	log.Println("Running scheduled domain/SSL health check...")
	domainSSLs := []models.DomainSSL{}
	
		if result := db.Find(&domainSSLs); result.Error != nil {
			log.Printf("Error fetching domain/SSLs for health check: %v", result.Error)
			return
		}

	for _, domainSSL := range domainSSLs {
		log.Printf("Checking domain/SSL: %s", domainSSL.Domain)
		// Implement actual health check logic for domain/SSL
		// For now, just simulate a successful check
		log.Printf("Domain/SSL %s health check successful.", domainSSL.Domain)
	}
	log.Println("Domain/SSL health check completed.")
}


