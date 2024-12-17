package account_channels

import (
	"github.com/gofiber/fiber/v2"
)

// WS url definiton
const xrplWebSocketURL = "wss://ripple.com"

func AccountChannelsHandler(c *fiber.Ctx) error {
	// Prepare ws client
	client, err := NewWebSocketClient(xrplWebSocketURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to connect to WebSocket Server",
		})
	}
	defer client.Connection.Close()

	// Request example
	response, err := map[string]interface{}{
		"id":     1,
		"command": "account_channels",
		"account": "r9cZA1mLK5R5Am25ArfXFmqgNwjZgnfk59",
	}

	// Send request
	response, err := client.SendResquest(response)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send request from XRPLedger",
		})
	}

	return c.JSON(response)
}

