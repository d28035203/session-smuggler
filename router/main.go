// Package router registers HTTP endpoints and binds them to handler functions.
package router

import (
	"github.com/d28035203/session-smuggler/database"
	"github.com/d28035203/session-smuggler/handlers"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes mounts versioned API routes under /api/v1.
// The database instance is closed over so handlers can query without globals.
func SetupRoutes(app *fiber.App, db *database.Database) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Liveness probe for Docker/Kubernetes health checks.
	v1.Get("/health", handlers.HandleHealthCheck)

	// Auth endpoints — each handler receives the shared DB connection.
	v1.Post("/register", func(c *fiber.Ctx) error {
		return handlers.HandleRegister(c, db)
	})
	v1.Post("/login", func(c *fiber.Ctx) error {
		return handlers.HandleLogin(c, db)
	})
	v1.Post("/logout", func(c *fiber.Ctx) error {
		return handlers.HandleLogout(c, db)
	})
	v1.Get("/authentication", func(c *fiber.Ctx) error {
		return handlers.HandleIsAuthenticated(c, db)
	})
}
