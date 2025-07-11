package handlers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"monitron-server/models"
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
func GetUsers(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		users := []models.User{}
		err := db.Select(&users, `SELECT id, username, email, role, status, last_login, created_at, updated_at FROM users`)
		if err != nil {
			log.Printf("Error fetching users: %v", err)
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
func GetUser(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuidID, err := uuid.Parse(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
		}

		user := models.User{}
		err = db.Get(&user, `SELECT id, username, email, role, status, last_login, created_at, updated_at FROM users WHERE id = $1`, uuidID)
		if err != nil {
			log.Printf("Error fetching user: %v", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
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
func UpdateUser(db *sqlx.DB) fiber.Handler {
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

		user.ID = uuidID
		user.UpdatedAt = time.Now()

		query := `UPDATE users SET username = :username, email = :email, role = :role, status = :status, updated_at = :updated_at WHERE id = :id`

		result, err := db.NamedExec(query, user)
		if err != nil {
			log.Printf("Error updating user: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update user"})
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}

		return c.JSON(user)
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
func DeleteUser(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuidID, err := uuid.Parse(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
		}

		result, err := db.Exec(`DELETE FROM users WHERE id = $1`, uuidID)
		if err != nil {
			log.Printf("Error deleting user: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete user"})
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
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
func ChangePassword(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uuid.UUID)

		passwordChange := struct {
			CurrentPassword string `json:"current_password"`
			NewPassword     string `json:"new_password"`
		}{}

		if err := c.BodyParser(&passwordChange); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		user := models.User{}
		err := db.Get(&user, `SELECT password FROM users WHERE id = $1`, userID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
		}

		// Verify current password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordChange.CurrentPassword))
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
		_, err = db.Exec(`UPDATE users SET password = $1, updated_at = $2 WHERE id = $3`, string(hashedPassword), time.Now(), userID)
		if err != nil {
			log.Printf("Error updating password: %v", err)
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
func ForgotPassword(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		resetRequest := struct {
			Email string `json:"email"`
		}{}

		if err := c.BodyParser(&resetRequest); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		user := models.User{}
		err := db.Get(&user, `SELECT id FROM users WHERE email = $1`, resetRequest.Email)
		if err != nil {
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

		query := `INSERT INTO password_reset_tokens (id, user_id, token, expires_at, created_at)
				  VALUES (:id, :user_id, :token, :expires_at, :created_at)`

		_, err = db.NamedExec(query, passwordResetToken)
		if err != nil {
			log.Printf("Error saving password reset token: %v", err)
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
func ResetPassword(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		resetRequest := struct {
			Token       string `json:"token"`
			NewPassword string `json:"new_password"`
		}{}

		if err := c.BodyParser(&resetRequest); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		passwordResetToken := models.PasswordResetToken{}
		err := db.Get(&passwordResetToken, `SELECT * FROM password_reset_tokens WHERE token = $1 AND expires_at > $2`, resetRequest.Token, time.Now())
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		// Hash new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(resetRequest.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing new password: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not reset password"})
		}

		// Update user\'s password
		_, err = db.Exec(`UPDATE users SET password = $1, updated_at = $2 WHERE id = $3`, string(hashedPassword), time.Now(), passwordResetToken.UserID)
		if err != nil {
			log.Printf("Error updating user password: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not reset password"})
		}

		// Invalidate the token
		_, err = db.Exec(`DELETE FROM password_reset_tokens WHERE id = $1`, passwordResetToken.ID)
		if err != nil {
			log.Printf("Error deleting password reset token: %v", err)
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Password reset successfully"})
	}
}

