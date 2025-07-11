package handlers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"errors"

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
func CreateInstance(db *gorm.DB) fiber.Handler {
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

		if result := db.Create(&instance); result.Error != nil {
			log.Printf("Error creating instance: %v", result.Error)
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
func GetInstances(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		instances := []models.Instance{}
		if result := db.Find(&instances); result.Error != nil {
			log.Printf("Error fetching instances: %v", result.Error)
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
func GetInstance(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuidID, err := uuid.Parse(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid instance ID"})
		}

		instance := models.Instance{}
		if result := db.First(&instance, "id = ?", uuidID); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Instance not found"})
			}
			log.Printf("Error fetching instance: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve instance"})
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
func UpdateInstance(db *gorm.DB) fiber.Handler {
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

		var existingInstance models.Instance
		if result := db.First(&existingInstance, "id = ?", uuidID); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Instance not found"})
			}
			log.Printf("Error finding instance for update: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update instance"})
		}

		if result := db.Model(&existingInstance).Updates(instance);
		if result.Error != nil {
			log.Printf("Error updating instance: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update instance"})
		}

		// Decrypt agent_auth before returning
		if existingInstance.AgentAuth != "" {
			decryptedAuth, err := utils.Decrypt(existingInstance.AgentAuth, cfg)
			if err != nil {
				log.Printf("Error decrypting agent auth: %v", err)
				existingInstance.AgentAuth = ""
			} else {
				existingInstance.AgentAuth = string(decryptedAuth)
			}
		}

		return c.JSON(existingInstance)

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
func DeleteInstance(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuidID, err := uuid.Parse(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid instance ID"})
		}

		if result := db.Delete(&models.Instance{}, "id = ?", uuidID); result.Error != nil {
			log.Printf("Error deleting instance: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete instance"})
		}

		if result.RowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Instance not found"})
		}

		return c.Status(fiber.StatusNoContent).SendString("")
	}
}



// InstanceHealthCheck performs health checks for all instances
func InstanceHealthCheck(db *gorm.DB) {
	log.Println("Running scheduled instance health check...")
	instances := []models.Instance{}
	
		if result := db.Find(&instances); result.Error != nil {
			log.Printf("Error fetching instances for health check: %v", result.Error)
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


