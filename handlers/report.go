package handlers

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"monitron-server/messaging"
	"monitron-server/models"
)

// CreateReport
// @Summary Create a new report
// @Description Create a new report entry and queue it for generation
// @Tags Reports
// @Accept json
// @Produce json
// @Param report body models.Report true "Report object to be created"
// @Success 201 {object} models.Report
// @Failure 400 {object} map[string]string "error": "Cannot parse JSON"
// @Failure 500 {object} map[string]string "error": "Could not create report"
// @Security ApiKeyAuth
// @Router /reports [post]
func CreateReport(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		report := new(models.Report)
		if err := c.BodyParser(report); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		report.ID = uuid.New()
		report.Status = "pending"
		report.CreatedAt = time.Now()
		report.UpdatedAt = time.Now()

		err := db.Create(report).Error
		if err != nil {
			log.Printf("Error inserting report: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create report"})
		}

		// Publish message to RabbitMQ for report generation
		reportJSON, _ := json.Marshal(report)
		err = messaging.PublishMessage("report_generation_queue", reportJSON)
		if err != nil {
			log.Printf("Error publishing report generation message: %v", err)
			// Optionally update report status to failed if message couldn't be published
		}

		return c.Status(fiber.StatusCreated).JSON(report)
	}
}

// GetReports
// @Summary Get all reports
// @Description Retrieve a list of all reports
// @Tags Reports
// @Produce json
// @Success 200 {array} models.Report
// @Failure 500 {object} map[string]string "error": "Could not retrieve reports"
// @Security ApiKeyAuth
// @Router /reports [get]
func GetReports(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		reports := []models.Report{}
		err := db.Select(&reports).Error
		if err != nil {
			log.Printf("Error fetching reports: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve reports"})
		}
		return c.JSON(reports)
	}
}

// GetReport
// @Summary Get report by ID
// @Description Retrieve a single report by its ID
// @Tags Reports
// @Produce json
// @Param id path string true "Report ID"
// @Success 200 {object} models.Report
// @Failure 400 {object} map[string]string "error": "Invalid report ID"
// @Failure 404 {object} map[string]string "error": "Report not found"
// @Failure 500 {object} map[string]string "error": "Could not retrieve report"
// @Security ApiKeyAuth
// @Router /reports/{id} [get]
func GetReport(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		_, err := uuid.Parse(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid report ID"})
		}

		report := models.Report{}
		err = db.Where("id = ?", id).First(&report).Error
		if err != nil {
			log.Printf("Error fetching report: %v", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Report not found"})
		}
		return c.JSON(report)
	}
}
