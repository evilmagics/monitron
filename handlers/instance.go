package handlers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"monitron-server/config"
	"monitron-server/models"
	"monitron-server/utils"
)

// CreateInstance
// @Summary Create a new instance
// @Description Create a new monitoring instance
// @Tags Instances
// @Accept json
// @Produce json
// @Param instance body models.Instance true "Instance object to be created"
// @Success 201 {object} models.Instance
// @Failure 400 {object} map[string]string "error": "Cannot parse JSON"
// @Failure 500 {object} map[string]string "error": "Could not create instance"
// @Security ApiKeyAuth
// @Router /instances [post]
func CreateInstance(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		instance := new(models.Instance)
		if err := c.BodyParser(instance); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		cfg := config.LoadConfig()
		if instance.AgentAuth != "" {
			encryptedAuth, err := utils.Encrypt([]byte(instance.AgentAuth), cfg)
			if err != nil {
				log.Printf("Error encrypting agent auth: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not encrypt agent authentication"})
			}
			instance.AgentAuth = encryptedAuth
		}

		instance.ID = uuid.New()
		instance.CreatedAt = time.Now()
		instance.UpdatedAt = time.Now()

		query := `INSERT INTO instances (id, name, host, check_interval, check_timeout, agent_port, agent_auth, description, label, "group", created_at, updated_at)
				  VALUES (:id, :name, :host, :check_interval, :check_timeout, :agent_port, :agent_auth, :description, :label, :group, :created_at, :updated_at)`

		_, err := db.NamedExec(query, instance)
		if err != nil {
			log.Printf("Error inserting instance: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create instance"})
		}

		// Decrypt agent_auth before returning
		if instance.AgentAuth != "" {
			decryptedAuth, err := utils.Decrypt(instance.AgentAuth, cfg)
			if err != nil {
				log.Printf("Error decrypting agent auth: %v", err)
				// Continue without agent_auth if decryption fails to avoid exposing encrypted string
				instance.AgentAuth = ""
			} else {
				instance.AgentAuth = string(decryptedAuth)
			}
		}

		return c.Status(fiber.StatusCreated).JSON(instance)
	}
}

// GetInstances
// @Summary Get all instances
// @Description Retrieve a list of all monitoring instances
// @Tags Instances
// @Produce json
// @Success 200 {array} models.Instance
// @Failure 500 {object} map[string]string "error": "Could not retrieve instances"
// @Security ApiKeyAuth
// @Router /instances [get]
func GetInstances(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		instances := []models.Instance{}
		err := db.Select(&instances, `SELECT * FROM instances`)
		if err != nil {
			log.Printf("Error fetching instances: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve instances"})
		}

		cfg := config.LoadConfig()
		for i := range instances {
			if instances[i].AgentAuth != "" {
				decryptedAuth, err := utils.Decrypt(instances[i].AgentAuth, cfg)
				if err != nil {
					log.Printf("Error decrypting agent auth for instance %s: %v", instances[i].ID, err)
					instances[i].AgentAuth = ""
				} else {
					instances[i].AgentAuth = string(decryptedAuth)
				}
			}
		}

		return c.JSON(instances)
	}
}

// GetInstance
// @Summary Get instance by ID
// @Description Retrieve a single monitoring instance by its ID
// @Tags Instances
// @Produce json
// @Param id path string true "Instance ID"
// @Success 200 {object} models.Instance
// @Failure 400 {object} map[string]string "error": "Invalid instance ID"
// @Failure 404 {object} map[string]string "error": "Instance not found"
// @Failure 500 {object} map[string]string "error": "Could not retrieve instance"
// @Security ApiKeyAuth
// @Router /instances/{id} [get]
func GetInstance(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuidID, err := uuid.Parse(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid instance ID"})
		}

		instance := models.Instance{}
		err = db.Get(&instance, `SELECT * FROM instances WHERE id = $1`, uuidID)
		if err != nil {
			log.Printf("Error fetching instance: %v", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Instance not found"})
		}

		cfg := config.LoadConfig()
		if instance.AgentAuth != "" {
			decryptedAuth, err := utils.Decrypt(instance.AgentAuth, cfg)
			if err != nil {
				log.Printf("Error decrypting agent auth for instance %s: %v", instance.ID, err)
				instance.AgentAuth = ""
			} else {
				instance.AgentAuth = string(decryptedAuth)
			}
		}

		return c.JSON(instance)
	}
}

// UpdateInstance
// @Summary Update an existing instance
// @Description Update details of an existing monitoring instance by its ID
// @Tags Instances
// @Accept json
// @Produce json
// @Param id path string true "Instance ID"
// @Param instance body models.Instance true "Instance object with updated fields"
// @Success 200 {object} models.Instance
// @Failure 400 {object} map[string]string "error": "Invalid instance ID" or "Cannot parse JSON"
// @Failure 404 {object} map[string]string "error": "Instance not found"
// @Failure 500 {object} map[string]string "error": "Could not update instance"
// @Security ApiKeyAuth
// @Router /instances/{id} [put]
func UpdateInstance(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuidID, err := uuid.Parse(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid instance ID"})
		}

		instance := new(models.Instance)
		if err := c.BodyParser(instance); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		cfg := config.LoadConfig()
		if instance.AgentAuth != "" {
			encryptedAuth, err := utils.Encrypt([]byte(instance.AgentAuth), cfg)
			if err != nil {
				log.Printf("Error encrypting agent auth: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not encrypt agent authentication"})
			}
			instance.AgentAuth = encryptedAuth
		}

		instance.ID = uuidID
		instance.UpdatedAt = time.Now()

		query := `UPDATE instances SET name = :name, host = :host, check_interval = :check_interval, check_timeout = :check_timeout, 
				  agent_port = :agent_port, agent_auth = :agent_auth, description = :description, label = :label, "group" = :group, updated_at = :updated_at
				  WHERE id = :id`

		result, err := db.NamedExec(query, instance)
		if err != nil {
			log.Printf("Error updating instance: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update instance"})
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Instance not found"})
		}

		// Decrypt agent_auth before returning
		if instance.AgentAuth != "" {
			decryptedAuth, err := utils.Decrypt(instance.AgentAuth, cfg)
			if err != nil {
				log.Printf("Error decrypting agent auth: %v", err)
				instance.AgentAuth = ""
			} else {
				instance.AgentAuth = string(decryptedAuth)
			}
		}

		return c.JSON(instance)
	}
}

// DeleteInstance
// @Summary Delete an instance
// @Description Delete a monitoring instance by its ID
// @Tags Instances
// @Produce json
// @Param id path string true "Instance ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "error": "Invalid instance ID"
// @Failure 404 {object} map[string]string "error": "Instance not found"
// @Failure 500 {object} map[string]string "error": "Could not delete instance"
// @Security ApiKeyAuth
// @Router /instances/{id} [delete]
func DeleteInstance(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuidID, err := uuid.Parse(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid instance ID"})
		}

		result, err := db.Exec(`DELETE FROM instances WHERE id = $1`, uuidID)
		if err != nil {
			log.Printf("Error deleting instance: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete instance"})
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Instance not found"})
		}

		return c.Status(fiber.StatusNoContent).SendString("")
	}
}



// InstanceHealthCheck performs health checks for all instances
func InstanceHealthCheck(db *sqlx.DB) {
	log.Println("Running scheduled instance health check...")
	instances := []models.Instance{}
	
	err := db.Select(&instances, `SELECT * FROM instances`)
	if err != nil {
		log.Printf("Error fetching instances for health check: %v", err)
		return
	}

	for _, instance := range instances {
		log.Printf("Checking instance: %s (Host: %s)", instance.Name, instance.Host)
		// Implement actual health check logic for instances
		// For now, just simulate a successful check
		log.Printf("Instance %s health check successful.", instance.Name)
	}
	log.Println("Instance health check completed.")
}


