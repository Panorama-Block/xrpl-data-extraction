package main

import (
	"log"
	"os"

	"github.com/Panorama-Block/xrpl-data-extraction/config"
	"github.com/Panorama-Block/xrpl-data-extraction/internal/server"
	"github.com/Panorama-Block/xrpl-data-extraction/internal/xrpl"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize XRPL Manager
	manager, err := xrpl.NewXRPLManager(cfg.APIBaseURL, cfg.WebSocketURL)
	if err != nil {
		log.Fatalf("Failed to initialize XRPL manager: %v", err)
	}

	// Initialize Fiber app
	app := fiber.New()

	// Apply logging middleware globally
	app.Use(server.LoggingMiddleware)

	// Setup routes with manager's clients
	server.SetupRoutes(app, manager.GetHTTPClient(), manager.GetWSClient())

	// Start the server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server running at http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}
