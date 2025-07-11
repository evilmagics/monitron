package handlers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"monitron-server/models"
)

// CreateLogEntry handles the creation of a new log entry
func CreateLogEntry(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logEntry := new(models.LogEntry)
		if err := c.BodyParser(logEntry); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		logEntry.ID = uuid.New()
		logEntry.Timestamp = time.Now()

		// Insert log entry into the database
		err := db.Create(logEntry).Error
		if err != nil {
			log.Printf("Error inserting log entry: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create log entry"})
		}

		return c.Status(fiber.StatusCreated).JSON(logEntry)
	}
}

// GetLogEntries handles fetching all log entries
func GetLogEntries(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logEntries := []models.LogEntry{}
		err := db.Select(&logEntries, `SELECT * FROM log_entries ORDER BY timestamp DESC`)
		if err != nil {
			log.Printf("Error fetching log entries: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve log entries"})
		}

		return c.JSON(logEntries)
	}
}

// GetLogEntry handles fetching a single log entry by ID
func GetLogEntry(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuidID, err := uuid.Parse(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid log entry ID"})
		}

		logEntry := models.LogEntry{}
		err = db.Where("id = ?", uuidID).First(&logEntry).Error
		if err != nil {
			log.Printf("Error fetching log entry: %v", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Log entry not found"})
		}

		return c.JSON(logEntry)
	}
}
