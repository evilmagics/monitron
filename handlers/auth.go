package handlers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
	"github.com/go-playground/validator/v10"

	"monitron-server/models"
	"monitron-server/utils"
)

var validate = validator.New()

// RegisterUser
// @Summary Register a new user
// @Description Register a new user with username, email, and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body models.User true "User object for registration"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string "error": "Cannot parse JSON"
// @Failure 500 {object} map[string]string "error": "Could not register user"
// @Router /auth/register [post]
func RegisterUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := new(models.User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		if err := validate.Struct(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not hash password"})
		}
		user.Password = string(hashedPassword)

		user.ID = uuid.New()
		user.Role = "user" // Default role
		user.Status = "active" // Default status
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		if result := db.Create(&user); result.Error != nil {
			log.Printf("Error inserting user: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not register user"})
		}

		// Do not return password hash
		user.Password = ""
		return c.Status(fiber.StatusCreated).JSON(user)
	}
}

// LoginUser
// @Summary Log in a user
// @Description Authenticate user and return a JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body object{email:string,password:string} true "User credentials"
// @Success 200 {object} map[string]interface{} "message": "Login successful", "token": "<JWT_TOKEN>", "user": "<USER_OBJECT>"
// @Failure 400 {object} map[string]string "error": "Cannot parse JSON"
// @Failure 401 {object} map[string]string "error": "Invalid credentials"
// @Failure 500 {object} map[string]string "error": "Could not generate token"
// @Router /auth/login [post]
func LoginUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		loginRequest := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{}

		if err := c.BodyParser(&loginRequest); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		if err := validate.Struct(loginRequest); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		user := models.User{}
		if result := db.First(&user, "email = ?", loginRequest.Email); result.Error != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		// Compare password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		// Update last login time
		now := time.Now()
		user.LastLogin = &now
		if result := db.Save(&user); result.Error != nil {
			log.Printf("Error updating last login: %v", result.Error)
		}

		token, err := utils.GenerateJWT(user.ID, user.Role)
		if err != nil {
			log.Printf("Error generating JWT: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
		}

		user.Password = ""
		return c.JSON(fiber.Map{"message": "Login successful", "token": token, "user": user})
	}
}

