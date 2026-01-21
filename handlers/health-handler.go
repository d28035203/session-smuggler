package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// HandleHealthCheck returns a small JSON payload used by load balancers and orchestrators
// to confirm the process is up (does not check database connectivity).
func HandleHealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": "session-smuggler",
	})
}
