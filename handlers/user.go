package handlers

import (
	"log"
	"time"

	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"monitron-server/models"
	"monitron-server/utils/validate"
)

// GetUsers
// @Summary Get all users (Admin Only)
// @Description Retrieve a list of all registered users
// @Tags User Management
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} map[string]string "error": "Could not retrieve users"
// @Security ApiKeyAuth
// @Router /users [get]
func GetUsers(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		users := []models.User{}
		if result := db.Select("id", "username", "email", "role", "status", "last_login", "created_at", "updated_at").Find(&users); result.Error != nil {
			log.Printf("Error fetching users: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve users"})
		}

		return c.JSON(users)
	}
}

// GetUser
// @Summary Get user by ID (Admin Only)
// @Description Retrieve a single user by their ID
// @Tags User Management
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string "error": "Invalid user ID"
// @Failure 404 {object} map[string]string "error": "User not found"
// @Failure 500 {object} map[string]string "error": "Could not retrieve user"
// @Security ApiKeyAuth
// @Router /users/{id} [get]
func GetUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuidID, err := uuid.Parse(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
		}

		user := models.User{}
		if result := db.Select("id", "username", "email", "role", "status", "last_login", "created_at", "updated_at").First(&user, "id = ?", uuidID); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
			}
			log.Printf("Error fetching user: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve user"})
		}

		return c.JSON(user)
	}
}

// UpdateUser
// @Summary Update an existing user (Admin Only)
// @Description Update details of an existing user by their ID
// @Tags User Management
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.User true "User object with updated fields"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string "error": "Invalid user ID" or "Cannot parse JSON"
// @Failure 404 {object} map[string]string "error": "User not found"
// @Failure 500 {object} map[string]string "error": "Could not update user"
// @Security ApiKeyAuth
// @Router /users/{id} [put]
func UpdateUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuidID, err := uuid.Parse(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
		}

		user := new(models.User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		var existingUser models.User
		if result := db.First(&existingUser, "id = ?", uuidID); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
			}
			log.Printf("Error finding user for update: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update user"})
		}

		if err := validate.V.Struct(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if result := db.Model(&existingUser).Updates(user); result.Error != nil {
			log.Printf("Error updating user: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update user"})
		}

		return c.JSON(existingUser)
	}
}

// DeleteUser
// @Summary Delete a user (Admin Only)
// @Description Delete a user by their ID
// @Tags User Management
// @Produce json
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "error": "Invalid user ID"
// @Failure 404 {object} map[string]string "error": "User not found"
// @Failure 500 {object} map[string]string "error": "Could not delete user"
// @Security ApiKeyAuth
// @Router /users/{id} [delete]
func DeleteUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuidID, err := uuid.Parse(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
		}

		if result := db.Delete(&models.User{}, "id = ?", uuidID); result.Error != nil {
			log.Printf("Error deleting user: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete user"})
		} else if result.RowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}

		return c.Status(fiber.StatusNoContent).SendString("")
	}
}

// ChangePassword
// @Summary Change password for authenticated user
// @Description Allows an authenticated user to change their password
// @Tags User Management
// @Accept json
// @Produce json
// @Param passwordChange body object{current_password:string,new_password:string} true "Current and new passwords"
// @Success 200 {object} map[string]string "message": "Password changed successfully"
// @Failure 400 {object} map[string]string "error": "Cannot parse JSON"
// @Failure 401 {object} map[string]string "error": "Invalid current password"
// @Failure 500 {object} map[string]string "error": "Could not change password"
// @Security ApiKeyAuth
// @Router /user/change-password [put]
func ChangePassword(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uuid.UUID)

		passwordChange := struct {
			CurrentPassword string `json:"current_password"`
			NewPassword     string `json:"new_password"`
		}{}

		if err := validate.V.Struct(passwordChange); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		user := models.User{}
		if result := db.First(&user, "id = ?", userID); result.Error != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
		}

		// Verify current password
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordChange.CurrentPassword))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid current password"})
		}

		// Hash new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordChange.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing new password: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not hash new password"})
		}

		// Update password in DB
		user.Password = string(hashedPassword)
		user.UpdatedAt = time.Now()
		if result := db.Save(&user); result.Error != nil {
			log.Printf("Error updating password: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not change password"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Password changed successfully"})
	}
}

// ForgotPassword
// @Summary Initiate password reset
// @Description Sends a password reset link to the user\'s email
// @Tags User Management
// @Accept json
// @Produce json
// @Param email body object{email:string} true "User email"
// @Success 200 {object} map[string]string "message": "If an account with that email exists, a password reset link has been sent."
// @Failure 400 {object} map[string]string "error": "Cannot parse JSON"
// @Failure 500 {object} map[string]string "error": "Could not initiate password reset"
// @Router /password/forgot [post]
func ForgotPassword(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		resetRequest := struct {
			Email string `json:"email"`
		}{}

		if err := validate.V.Struct(resetRequest); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		user := models.User{}
		if result := db.First(&user, "email = ?", resetRequest.Email); result.Error != nil {
			// For security, always return a generic success message even if user not found
			log.Printf("Forgot password request for non-existent email: %s", resetRequest.Email)
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "If an account with that email exists, a password reset link has been sent."})
		}

		// Generate token
		token := uuid.New().String()
		expiresAt := time.Now().Add(30 * time.Minute) // Token valid for 30 minutes

		passwordResetToken := models.PasswordResetToken{
			ID:        uuid.New(),
			UserID:    user.ID,
			Token:     token,
			ExpiresAt: expiresAt,
			CreatedAt: time.Now(),
		}

		if result := db.Create(&passwordResetToken); result.Error != nil {
			log.Printf("Error saving password reset token: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not initiate password reset"})
		}

		// TODO: Send email with reset link (e.g., http://your-frontend/reset-password?token=TOKEN)
		log.Printf("Password reset token for %s: %s", user.Email, token)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "If an account with that email exists, a password reset link has been sent."})
	}
}

// ResetPassword
// @Summary Reset password with token
// @Description Resets user password using a valid token
// @Tags User Management
// @Accept json
// @Produce json
// @Param resetRequest body object{token:string,new_password:string} true "Token and new password"
// @Success 200 {object} map[string]string "message": "Password reset successfully"
// @Failure 400 {object} map[string]string "error": "Invalid or expired token" or "Cannot parse JSON"
// @Failure 500 {object} map[string]string "error": "Could not reset password"
// @Router /password/reset [post]
func ResetPassword(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		resetRequest := struct {
			Token       string `json:"token"`
			NewPassword string `json:"new_password"`
		}{}

		if err := c.BodyParser(&resetRequest); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		if err := validate.V.Struct(resetRequest); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		passwordResetToken := models.PasswordResetToken{}
		if result := db.First(&passwordResetToken, "token = ? AND expires_at > ?", resetRequest.Token, time.Now()); result.Error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		// Hash new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(resetRequest.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing new password: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not reset password"})
		}

		// Update user\'s password
		var user models.User
		if result := db.First(&user, "id = ?", passwordResetToken.UserID); result.Error != nil {
			log.Printf("Error finding user to reset password: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not reset password"})
		}

		user.Password = string(hashedPassword)
		user.UpdatedAt = time.Now()
		if result := db.Save(&user); result.Error != nil {
			log.Printf("Error updating user password: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not reset password"})
		}

		// Invalidate.V the token
		if result := db.Delete(&passwordResetToken); result.Error != nil {
			log.Printf("Error deleting password reset token: %v", result.Error)
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Password reset successfully"})
	}
}
