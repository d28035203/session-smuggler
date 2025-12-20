package app

import (
	"fmt"
	"os"

	"github.com/d28035203/legendary-succotash/database"
	"github.com/d28035203/legendary-succotash/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

// SetupAndRunApp loads config, connects to the database, and starts Fiber.
func SetupAndRunApp() error {
	_ = godotenv.Load()

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
		AppName: "legendary-succotash",
	})
	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path} ${latency}\n",
	}))

	router.SetupRoutes(app, db)

	addr := fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))
	return app.Listen(addr)
}
