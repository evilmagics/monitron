package handlers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"monitron-server/models"
)

// CreateOperationalPage handles the creation of a new operational page
func CreateOperationalPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page := new(models.OperationalPage)
		if err := c.BodyParser(page); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		page.ID = uuid.New()
		page.CreatedAt = time.Now()
		page.UpdatedAt = time.Now()

		err := db.Create(page).Error
		if err != nil {
			log.Printf("Error inserting operational page: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create operational page"})
		}

		return c.Status(fiber.StatusCreated).JSON(page)
	}
}

// GetOperationalPages handles fetching all operational pages
func GetOperationalPages(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pages := []models.OperationalPage{}
		err := db.Select(&pages, `SELECT * FROM operational_pages`)
		if err != nil {
			log.Printf("Error fetching operational pages: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve operational pages"})
		}

		return c.JSON(pages)
	}
}

// GetOperationalPage handles fetching a single operational page by ID or slug
func GetOperationalPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idOrSlug := c.Params("idOrSlug")

		page := models.OperationalPage{}
		err := db.Where("id = ? OR slug = ?", idOrSlug, idOrSlug).Error
		if err != nil {
			log.Printf("Error fetching operational page: %v", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Operational page not found"})
		}

		return c.JSON(page)
	}
}

// UpdateOperationalPage handles updating an existing operational page
func UpdateOperationalPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idOrSlug := c.Params("idOrSlug")

		page := new(models.OperationalPage)
		if err := c.BodyParser(page); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		page.UpdatedAt = time.Now()

		result := db.Where("id = ? OR slug = ?", idOrSlug, idOrSlug).Updates(page)
		if result.Error != nil {
			log.Printf("Error updating operational page: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update operational page"})
		}

		if result.RowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Operational page not found"})
		}

		return c.JSON(page)
	}
}

// DeleteOperationalPage handles deleting an operational page by ID or slug
func DeleteOperationalPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idOrSlug := c.Params("idOrSlug")

		result := db.Where("id = ? OR slug = ?", idOrSlug, idOrSlug).Delete(&models.OperationalPage{})
		if result.Error != nil {
			log.Printf("Error deleting operational page: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete operational page"})
		}

		if result.RowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Operational page not found"})
		}

		return c.Status(fiber.StatusNoContent).SendString("")
	}
}

// AddComponentToOperationalPage handles adding a component to an operational page
func AddComponentToOperationalPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pageID := c.Params("pageID")
		uuidPageID, err := uuid.Parse(pageID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid page ID"})
		}

		component := new(models.OperationalPageComponent)
		if err := c.BodyParser(component); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		component.ID = uuid.New()
		component.PageID = uuidPageID
		component.CreatedAt = time.Now()
		component.UpdatedAt = time.Now()

		err = db.Create(&component).Error
		if err != nil {
			log.Printf("Error adding component to operational page: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not add component to operational page"})
		}

		return c.Status(fiber.StatusCreated).JSON(component)
	}
}

// GetComponentsForOperationalPage handles fetching components for a specific operational page
func GetComponentsForOperationalPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pageID := c.Params("pageID")
		uuidPageID, err := uuid.Parse(pageID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid page ID"})
		}

		components := []models.OperationalPageComponent{}
		err = db.Where("page_id = ?", uuidPageID).Order("display_order ASC").Find(&components).Error
		if err != nil {
			log.Printf("Error fetching components for operational page: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve components"})
		}

		return c.JSON(components)
	}
}

// RemoveComponentFromOperationalPage handles removing a component from an operational page
func RemoveComponentFromOperationalPage(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pageID := c.Params("pageID")
		componentID := c.Params("componentID")

		uuidPageID, err := uuid.Parse(pageID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid page ID"})
		}
		uuidComponentID, err := uuid.Parse(componentID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid component ID"})
		}

		// result, err := db.Exec(`DELETE FROM operational_page_components WHERE page_id = $1 AND id = $2`, uuidPageID, uuidComponentID)
		result := db.Where("page_id = ? AND id = ?", uuidPageID, uuidComponentID).Delete(&models.OperationalPageComponent{})
		if result.Error != nil {
			log.Printf("Error removing component from operational page: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not remove component"})
		}

		if result.RowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Component not found on this page"})
		}

		return c.Status(fiber.StatusNoContent).SendString("")
	}
}
