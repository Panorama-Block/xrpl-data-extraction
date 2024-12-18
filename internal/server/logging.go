package server

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LoggingMiddleware(c *fiber.Ctx) error {
	start := time.Now()

	// Process the request
	err := c.Next()

	// Capture duration and log request details
	duration := time.Since(start)
	log.Printf("[%s] %s | %d | %s | %s",
		c.Method(),     // HTTP method (e.g., GET, POST)
		c.Path(),       // Request path
		c.Response().StatusCode(), // Response status code
		c.IP(),         // Client IP
		duration,       // Request duration
	)

	return err
}
