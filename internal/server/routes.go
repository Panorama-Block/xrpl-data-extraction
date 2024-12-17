package server

import (
	"xrpl-data-extraction/internal/account_channels"

	"github.com/gofiber/fiber/v2"
)

// Routes definition
func SetupRoutes(app *fiber.App) {
	app.Get("/xrp/accounts-channels", account_channels.AccountChannelsHandler)
}