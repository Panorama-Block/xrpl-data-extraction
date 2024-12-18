package server

import (
	"encoding/json" //for json operations
	"github.com/Panorama-Block/xrpl-data-extraction/internal/accounts"
	"github.com/Panorama-Block/xrpl-data-extraction/internal/xrpl"
	"github.com/gofiber/fiber/v2"
)

// setup routes
func SetupRoutes(app *fiber.App, httpClient *xrpl.HTTPClient, wsClient *xrpl.WebSocketClient) {
	// Historical data account channels Endpoint
	app.Get("/accounts/:account/channels/historical", func(c *fiber.Ctx) error {
		account := c.Params("account") // extract account part from url
		destinationAccount := c.Query("destination_account", "") // extract destination account from query
		ledgerIndex := c.Query("ledger_index", "validated") // extract ledger index from query

		// call fetch historical account channels function in accounts package passing the http client, account, destination account and ledger index
		response, err := accounts.FetchHistoricalAccountChannels(httpClient, account, destinationAccount, ledgerIndex)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// decode the response
		var decodedResponse map[string]interface{}
		err = json.Unmarshal(response, &decodedResponse) // converts json to go object
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode response"})
		}

		return c.JSON(decodedResponse) // return the decoded response
	})

	// WS account channels endpoint
	app.Get("/accounts/:account/channels/realtime", func(c *fiber.Ctx) error {
		account := c.Params("account") 
		destinationAccount := c.Query("destination_account", "")

		go accounts.StreamRecentAccountChannels(wsClient, account, destinationAccount, func(data *accounts.AccountChannelsResponse) {
			// TODO ---> Handle real-time data
		}) // function to stream real-time account channels

		return c.JSON(fiber.Map{"message": "Subscribed to account_channels"}) // return message
	})
}
