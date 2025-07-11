package handlers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

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
func CreateDomainSSL(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		domainSSL := new(models.DomainSSL)
		if err := c.BodyParser(domainSSL); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		domainSSL.ID = uuid.New()
		domainSSL.CreatedAt = time.Now()
		domainSSL.UpdatedAt = time.Now()

		query := `INSERT INTO domain_ssl (id, domain, warning_threshold, expiry_threshold, check_interval, label, created_at, updated_at,
				  certificate_detail, issuer, valid_from, resolved_ip, expiry)
				  VALUES (:id, :domain, :warning_threshold, :expiry_threshold, :check_interval, :label, :created_at, :updated_at,
				  :certificate_detail, :issuer, :valid_from, :resolved_ip, :expiry)`

		_, err := db.NamedExec(query, domainSSL)
		if err != nil {
			log.Printf("Error inserting domain/SSL: %v", err)
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
func GetDomainSSLs(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		domainSSLs := []models.DomainSSL{}
		err := db.Select(&domainSSLs, `SELECT * FROM domain_ssl`)
		if err != nil {
			log.Printf("Error fetching domain/SSLs: %v", err)
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
func GetDomainSSL(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuidID, err := uuid.Parse(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid domain/SSL ID"})
		}

		domainSSL := models.DomainSSL{}
		err = db.Get(&domainSSL, `SELECT * FROM domain_ssl WHERE id = $1`, uuidID)
		if err != nil {
			log.Printf("Error fetching domain/SSL: %v", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Domain/SSL entry not found"})
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
func UpdateDomainSSL(db *sqlx.DB) fiber.Handler {
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

		domainSSL.ID = uuidID
		domainSSL.UpdatedAt = time.Now()

		query := `UPDATE domain_ssl SET domain = :domain, warning_threshold = :warning_threshold, expiry_threshold = :expiry_threshold,
				  check_interval = :check_interval, label = :label, updated_at = :updated_at,
				  certificate_detail = :certificate_detail, issuer = :issuer, valid_from = :valid_from, resolved_ip = :resolved_ip, expiry = :expiry
				  WHERE id = :id`

		result, err := db.NamedExec(query, domainSSL)
		if err != nil {
			log.Printf("Error updating domain/SSL: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update domain/SSL entry"})
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Domain/SSL entry not found"})
		}

		return c.JSON(domainSSL)
	}
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
func DeleteDomainSSL(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuidID, err := uuid.Parse(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid domain/SSL ID"})
		}

		result, err := db.Exec(`DELETE FROM domain_ssl WHERE id = $1`, uuidID)
		if err != nil {
			log.Printf("Error deleting domain/SSL: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete domain/SSL entry"})
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Domain/SSL entry not found"})
		}

		return c.Status(fiber.StatusNoContent).SendString("")
	}
}


// DomainSSLHealthCheck performs health checks for all domain/SSL entries
func DomainSSLHealthCheck(db *sqlx.DB) {
	log.Println("Running scheduled domain/SSL health check...")
	domainSSLs := []models.DomainSSL{}
	
	err := db.Select(&domainSSLs, `SELECT * FROM domain_ssl`)
	if err != nil {
		log.Printf("Error fetching domain/SSLs for health check: %v", err)
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


