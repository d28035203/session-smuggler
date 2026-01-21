// Package handlers contains HTTP handlers for authentication and health endpoints.
package handlers

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/d28035203/session-smuggler/database"
	"github.com/d28035203/session-smuggler/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// secretKey reads the JWT signing secret at call time so tests/env reloads work.
func secretKey() string {
	return os.Getenv("TOKEN_SECRET")
}

// respond sends a successful or informational JSON envelope.
func respond(c *fiber.Ctx, status int, message string, data interface{}) error {
	body := models.BuildResponse(http.StatusText(status), message, data, "")
	return c.Status(status).JSON(body)
}

// respondErr sends a failure envelope and includes err.Error() when present.
func respondErr(c *fiber.Ctx, status int, message string, err error) error {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	body := models.BuildResponse(http.StatusText(status), message, nil, errMsg)
	return c.Status(status).JSON(body)
}

// HandleRegister creates a new user with a bcrypt-hashed password.
// Expected JSON body: { "username": "...", "password": "..." }
func HandleRegister(c *fiber.Ctx, db *database.Database) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return respond(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	username := strings.TrimSpace(data["username"])
	password := data["password"]
	if username == "" || password == "" {
		return respond(c, fiber.StatusBadRequest, "Username and password are required", nil)
	}

	// Reject duplicate usernames before hashing to keep the happy path simple.
	var existing models.User
	if err := db.DB.Where("username = ?", username).First(&existing).Error; err == nil {
		return respond(c, fiber.StatusConflict, "Username already exists", nil)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return respondErr(c, fiber.StatusInternalServerError, "Failed to hash password", err)
	}

	user := models.User{Username: username, Password: hashed}
	if err := db.DB.Create(&user).Error; err != nil {
		return respondErr(c, fiber.StatusInternalServerError, "Failed to create user", err)
	}

	// Never return the password hash to the client.
	return respond(c, fiber.StatusCreated, "User created successfully", fiber.Map{
		"username": user.Username,
	})
}

// HandleLogin validates credentials and issues a JWT stored in user sessions.
// Expected JSON body: { "username": "...", "password": "..." }
func HandleLogin(c *fiber.Ctx, db *database.Database) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return respond(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	var user models.User
	if err := db.DB.Where("username = ?", data["username"]).First(&user).Error; err != nil {
		// Generic message avoids username enumeration.
		return respond(c, fiber.StatusUnauthorized, "Invalid credentials", nil)
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		return respond(c, fiber.StatusUnauthorized, "Invalid credentials", nil)
	}

	// Replace any existing session so only the latest login is valid.
	var session models.UserSessions
	if err := db.DB.Where("username = ?", user.Username).First(&session).Error; err == nil {
		_ = db.DB.Where("username = ?", user.Username).Delete(&models.UserSessions{})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    user.Username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	token, err := claims.SignedString([]byte(secretKey()))
	if err != nil {
		return respondErr(c, fiber.StatusInternalServerError, "Could not create token", err)
	}

	session = models.UserSessions{Username: user.Username, Token: token}
	if err := db.DB.Create(&session).Error; err != nil {
		return respondErr(c, fiber.StatusInternalServerError, "Error creating session", err)
	}

	// Also expose the token via Authorization header for convenience.
	c.Set("Authorization", "Bearer "+token)
	return respond(c, fiber.StatusOK, "Logged in", fiber.Map{
		"username": user.Username,
		"token":    token,
	})
}

// HandleLogout invalidates the caller's stored session (requires Bearer token).
func HandleLogout(c *fiber.Ctx, db *database.Database) error {
	authorized, claims, err := isAuthenticated(c, db)
	if err != nil {
		return respondErr(c, fiber.StatusInternalServerError, "Internal server error", err)
	}
	if !authorized {
		return respond(c, fiber.StatusUnauthorized, "Not authorized", nil)
	}

	if err := db.DB.Where("username = ?", claims.Issuer).Delete(&models.UserSessions{}).Error; err != nil {
		return respondErr(c, fiber.StatusInternalServerError, "Error logging out", err)
	}

	return respond(c, fiber.StatusOK, "Logged out", nil)
}

// HandleIsAuthenticated reports whether the request carries a valid session token.
func HandleIsAuthenticated(c *fiber.Ctx, db *database.Database) error {
	authorized, _, err := isAuthenticated(c, db)
	if err != nil {
		return respondErr(c, fiber.StatusInternalServerError, "Internal server error", err)
	}
	if !authorized {
		return respond(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	return respond(c, fiber.StatusOK, "Authorized", nil)
}

// isAuthenticated validates the Bearer JWT against the sessions table and expiry.
// Returns (false, nil, nil) for missing/invalid tokens rather than hard errors.
func isAuthenticated(c *fiber.Ctx, db *database.Database) (bool, *jwt.RegisteredClaims, error) {
	raw := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if raw == "" {
		return false, nil, nil
	}

	jwtToken, claims, err := parseToken(raw)
	if err != nil {
		return false, nil, nil
	}

	var userSession models.UserSessions
	if err := db.DB.Where("username = ?", claims.Issuer).First(&userSession).Error; err != nil {
		return false, nil, nil
	}

	// Token string must match what we stored at login (supports logout revocation).
	if userSession.Token != jwtToken.Raw {
		return false, nil, nil
	}

	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return false, nil, nil
	}

	return true, claims, nil
}

// parseToken parses and validates a JWT using HS256 and TOKEN_SECRET.
func parseToken(token string) (*jwt.Token, *jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}
	jwtToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey()), nil
	})
	if err != nil {
		return nil, nil, err
	}
	return jwtToken, claims, nil
}
