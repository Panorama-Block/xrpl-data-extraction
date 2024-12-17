package main

import (
	"log"
	"xrpl-api-backend/internal/server"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := Fiber.New()

	server.SetupRoutes(app)

	log.Println("Server is running on port 3000")
	log.Fatal(app.Listen(":3000"))
}