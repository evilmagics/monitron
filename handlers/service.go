package handlers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

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
func CreateService(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		service := new(models.Service)
		if err := c.BodyParser(service); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		service.ID = uuid.New()
		service.CreatedAt = time.Now()
		service.UpdatedAt = time.Now()

		query := `INSERT INTO services (id, name, api_type, check_interval, timeout, description, label, "group", created_at, updated_at,
				  http_method, http_health_url, http_expected_status,
				  grpc_host, grpc_port, grpc_auth, grpc_proto,
				  mqtt_host, mqtt_port, mqtt_qos, mqtt_topic, mqtt_auth,
				  tcp_host, tcp_port,
				  dns_domain_name,
				  ping_host)
				  VALUES (:id, :name, :api_type, :check_interval, :timeout, :description, :label, :group, :created_at, :updated_at,
				  :http_method, :http_health_url, :http_expected_status,
				  :grpc_host, :grpc_port, :grpc_auth, :grpc_proto,
				  :mqtt_host, :mqtt_port, :mqtt_qos, :mqtt_topic, :mqtt_auth,
				  :tcp_host, :tcp_port,
				  :dns_domain_name,
				  :ping_host)`

		_, err := db.NamedExec(query, service)
		if err != nil {
			log.Printf("Error inserting service: %v", err)
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
func GetServices(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		services := []models.Service{}
		err := db.Select(&services, `SELECT * FROM services`)
		if err != nil {
			log.Printf("Error fetching services: %v", err)
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
func GetService(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuidID, err := uuid.Parse(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid service ID"})
		}

		service := models.Service{}
		err = db.Get(&service, `SELECT * FROM services WHERE id = $1`, uuidID)
		if err != nil {
			log.Printf("Error fetching service: %v", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Service not found"})
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
func UpdateService(db *sqlx.DB) fiber.Handler {
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

		service.ID = uuidID
		service.UpdatedAt = time.Now()

		query := `UPDATE services SET name = :name, api_type = :api_type, check_interval = :check_interval, timeout = :timeout,
				  description = :description, label = :label, "group" = :group, updated_at = :updated_at,
				  http_method = :http_method, http_health_url = :http_health_url, http_expected_status = :http_expected_status,
				  grpc_host = :grpc_host, grpc_port = :grpc_port, grpc_auth = :grpc_auth, grpc_proto = :grpc_proto,
				  mqtt_host = :mqtt_host, mqtt_port = :mqtt_port, mqtt_qos = :mqtt_qos, mqtt_topic = :mqtt_topic, mqtt_auth = :mqtt_auth,
				  tcp_host = :tcp_host, tcp_port = :tcp_port,
				  dns_domain_name = :dns_domain_name,
				  ping_host = :ping_host
				  WHERE id = :id`

		result, err := db.NamedExec(query, service)
		if err != nil {
			log.Printf("Error updating service: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update service"})
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Service not found"})
		}

		return c.JSON(service)
	}
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
func DeleteService(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuidID, err := uuid.Parse(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid service ID"})
		}

		result, err := db.Exec(`DELETE FROM services WHERE id = $1`, uuidID)
		if err != nil {
			log.Printf("Error deleting service: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete service"})
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Service not found"})
		}

		return c.Status(fiber.StatusNoContent).SendString("")
	}
}



// ServiceHealthCheck performs health checks for all services
func ServiceHealthCheck(db *sqlx.DB) {
	log.Println("Running scheduled service health check...")
	services := []models.Service{}
	
	err := db.Select(&services, `SELECT * FROM services`)
	if err != nil {
		log.Printf("Error fetching services for health check: %v", err)
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


