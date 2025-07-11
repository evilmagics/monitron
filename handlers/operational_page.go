package handlers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"monitron-server/models"
)

// CreateOperationalPage handles the creation of a new operational page
func CreateOperationalPage(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page := new(models.OperationalPage)
		if err := c.BodyParser(page); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		page.ID = uuid.New()
		page.CreatedAt = time.Now()
		page.UpdatedAt = time.Now()

		query := `INSERT INTO operational_pages (id, slug, name, description, is_public, created_at, updated_at)
				  VALUES (:id, :slug, :name, :description, :is_public, :created_at, :updated_at)`

		_, err := db.NamedExec(query, page)
		if err != nil {
			log.Printf("Error inserting operational page: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create operational page"})
		}

		return c.Status(fiber.StatusCreated).JSON(page)
	}
}

// GetOperationalPages handles fetching all operational pages
func GetOperationalPages(db *sqlx.DB) fiber.Handler {
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
func GetOperationalPage(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idOrSlug := c.Params("idOrSlug")

		page := models.OperationalPage{}
		err := db.Get(&page, `SELECT * FROM operational_pages WHERE id = $1 OR slug = $1`, idOrSlug)
		if err != nil {
			log.Printf("Error fetching operational page: %v", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Operational page not found"})
		}

		return c.JSON(page)
	}
}

// UpdateOperationalPage handles updating an existing operational page
func UpdateOperationalPage(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idOrSlug := c.Params("idOrSlug")

		page := new(models.OperationalPage)
		if err := c.BodyParser(page); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		page.UpdatedAt = time.Now()

		query := `UPDATE operational_pages SET slug = :slug, name = :name, description = :description, is_public = :is_public, updated_at = :updated_at
				  WHERE id = :idOrSlug OR slug = :idOrSlug`

		result, err := db.NamedExec(query, map[string]interface{}{
			"slug": page.Slug,
			"name": page.Name,
			"description": page.Description,
			"is_public": page.IsPublic,
			"updated_at": page.UpdatedAt,
			"idOrSlug": idOrSlug,
		})
		if err != nil {
			log.Printf("Error updating operational page: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update operational page"})
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Operational page not found"})
		}

		return c.JSON(page)
	}
}

// DeleteOperationalPage handles deleting an operational page by ID or slug
func DeleteOperationalPage(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idOrSlug := c.Params("idOrSlug")

		result, err := db.Exec(`DELETE FROM operational_pages WHERE id = $1 OR slug = $1`, idOrSlug)
		if err != nil {
			log.Printf("Error deleting operational page: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete operational page"})
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Operational page not found"})
		}

		return c.Status(fiber.StatusNoContent).SendString("")
	}
}

// AddComponentToOperationalPage handles adding a component to an operational page
func AddComponentToOperationalPage(db *sqlx.DB) fiber.Handler {
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

		query := `INSERT INTO operational_page_components (id, page_id, component_type, component_id, component_name, display_order, description, created_at, updated_at)
				  VALUES (:id, :page_id, :component_type, :component_id, :component_name, :display_order, :description, :created_at, :updated_at)`

		_, err = db.NamedExec(query, component)
		if err != nil {
			log.Printf("Error adding component to operational page: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not add component to operational page"})
		}

		return c.Status(fiber.StatusCreated).JSON(component)
	}
}

// GetComponentsForOperationalPage handles fetching components for a specific operational page
func GetComponentsForOperationalPage(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pageID := c.Params("pageID")
		uuidPageID, err := uuid.Parse(pageID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid page ID"})
		}

		components := []models.OperationalPageComponent{}
		err = db.Select(&components, `SELECT * FROM operational_page_components WHERE page_id = $1 ORDER BY display_order ASC`, uuidPageID)
		if err != nil {
			log.Printf("Error fetching components for operational page: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve components"})
		}

		return c.JSON(components)
	}
}

// RemoveComponentFromOperationalPage handles removing a component from an operational page
func RemoveComponentFromOperationalPage(db *sqlx.DB) fiber.Handler {
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

		result, err := db.Exec(`DELETE FROM operational_page_components WHERE page_id = $1 AND id = $2`, uuidPageID, uuidComponentID)
		if err != nil {
			log.Printf("Error removing component from operational page: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not remove component"})
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Component not found on this page"})
		}

		return c.Status(fiber.StatusNoContent).SendString("")
	}
}

