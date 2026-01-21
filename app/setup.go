// Package app wires configuration, middleware, database, and HTTP routing together.
package app

import (
	"fmt"
	"os"

	"github.com/d28035203/session-smuggler/database"
	"github.com/d28035203/session-smuggler/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

// SetupAndRunApp loads environment variables, validates required settings,
// opens a PostgreSQL connection, mounts middleware and routes, then starts Fiber.
func SetupAndRunApp() error {
	// Load .env if present; production can rely on real environment variables instead.
	_ = godotenv.Load()

	// Fail fast when critical config is missing rather than panicking later at runtime.
	required := []string{
		"APP_HOST", "APP_PORT", "TOKEN_SECRET",
		"POSTGRES_HOST", "POSTGRES_PORT", "POSTGRES_USER",
		"POSTGRES_PASSWORD", "POSTGRES_DBNAME",
	}
	for _, key := range required {
		if os.Getenv(key) == "" {
			return fmt.Errorf("environment variable %s is required", key)
		}
	}

	db, err := database.Connect()
	if err != nil {
		return err
	}

	app := fiber.New(fiber.Config{
		AppName: "session-smuggler",
	})

	// CORS for browser clients; recover so panics become 500s instead of crashing the process.
	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path} ${latency}\n",
	}))

	// Register /api/v1 routes and inject the shared DB handle into handlers.
	router.SetupRoutes(app, db)

	addr := fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))
	return app.Listen(addr)
}
