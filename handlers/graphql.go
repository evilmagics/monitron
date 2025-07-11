package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/graphql-go/handler"
	"gorm.io/gorm"

	monitrongraphql "monitron-server/graphql"
	"monitron-server/utils"
)

// GraphQLHandler handles GraphQL requests
// @Summary GraphQL Endpoint
// @Description Access the GraphQL API for querying data
// @Tags GraphQL
// @Accept json
// @Produce json
// @Param query body string true "GraphQL query"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "error": "Invalid request"
// @Router /graphql [post]
func GraphQLHandler(db *gorm.DB) fiber.Handler {
	// Create a GraphQL schema
	schema, err := monitrongraphql.CreateSchema(db)
	if err != nil {
		log.Fatalf("failed to create graphql schema, error: %v", err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	return func(c *fiber.Ctx) error {
		w, r := utils.AdaptFiberToHTTP(c)
		h.ServeHTTP(w, r)
		return nil
	}
}
