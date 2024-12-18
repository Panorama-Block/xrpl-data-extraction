package server

import (
	"log" 
	"encoding/json" //for json operations

	"github.com/Panorama-Block/xrpl-data-extraction/internal/accounts"
	"github.com/Panorama-Block/xrpl-data-extraction/internal/xrpl"
	"github.com/Panorama-Block/xrpl-data-extraction/internal/ledger"
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

		go accounts.StreamRecentAccountChannels(wsClient, account, destinationAccount, func(data *accounts.AccountChannelsWSResponse) {
			// TODO ---> Handle real-time data
		}) // function to stream real-time account channels

		return c.JSON(fiber.Map{"message": "Subscribed to account_channels"}) // return message
	})

	
	// Account Currencies - Historical
	app.Post("/accounts/currencies/historical", func(c *fiber.Ctx) error {
		var payload accounts.AccountCurrenciesRequest
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}
		
		response, err := httpClient.Post("", payload)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		
		var decodedResponse map[string]interface{}
		err = json.Unmarshal(response, &decodedResponse)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode JSON response"})
		}
		
		return c.JSON(decodedResponse)
	})
	
	app.Get("/accounts/:account/currencies/realtime", func(c *fiber.Ctx) error {
		account := c.Params("account") 
		ledgerIndex := c.Query("ledger_index", "validated")

		go accounts.StreamAccountCurrencies(wsClient, account, ledgerIndex, func(data *accounts.AccountCurrenciesWSResponse) {
			// TODO ---> Handle real-time data
		}) // function to stream real-time account currencies

		return c.JSON(fiber.Map{"message": "Subscribed to account_currencies"}) // return message
	})

	// Account Info - Historical
	app.Post("/accounts/info/historical", func(c *fiber.Ctx) error {
		var payload accounts.AccountInfoRequest
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}

		response, err := httpClient.Post("", payload)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		var decodedResponse map[string]interface{}
		err = json.Unmarshal(response, &decodedResponse)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode JSON response"})
		}

		return c.JSON(decodedResponse)
	})

	// Account Info - Real-time WebSocket
	app.Get("/accounts/:account/info/realtime", func(c *fiber.Ctx) error {
		account := c.Params("account") // Extract the account parameter
		ledgerIndex := c.Query("ledger_index", "validated") // Extract ledger_index query parameter
		queue := c.Query("queue", "false") == "true"       // Extract queue query parameter

		go accounts.StreamAccountInfo(wsClient, account, ledgerIndex, queue, func(data *accounts.AccountInfoWSResponse) {
			// Log real-time data for debugging purposes
			log.Printf("Real-time account info: %+v", data)
		})

		return c.JSON(fiber.Map{"message": "Subscribed to account_info"}) // Return a success message
	})

	// Account Lines - Historical
	app.Post("/accounts/lines/historical", func(c *fiber.Ctx) error {
		var payload accounts.AccountLinesRequest
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}

		response, err := httpClient.Post("", payload)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		var decodedResponse map[string]interface{}
		err = json.Unmarshal(response, &decodedResponse)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode JSON response"})
		}

		return c.JSON(decodedResponse)
	})

	// Account Lines - Real-time
	app.Post("/accounts/lines/realtime", func(c *fiber.Ctx) error {
		var payload accounts.AccountLinesWSRequest
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}

		err := wsClient.Subscribe(payload)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		wsClient.ReadMessages(func(message []byte) {
			var response accounts.AccountLinesWSResponse
			err := json.Unmarshal(message, &response)
			if err != nil {
				return
			}
			log.Printf("Real-time data: %+v", response)
		})

		return c.JSON(fiber.Map{"message": "Subscribed to account_lines"})
	})

	app.Post("/accounts/nfts/historical", func(c *fiber.Ctx) error {
		var payload accounts.AccountNFTsRequest
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}

		response, err := httpClient.Post("", payload)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		var decodedResponse map[string]interface{}
		err = json.Unmarshal(response, &decodedResponse)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode JSON response"})
		}

		return c.JSON(decodedResponse)
	})

	app.Get("/accounts/:account/nfts/realtime", func(c *fiber.Ctx) error {
		account := c.Params("account")
		ledgerIndex := c.Query("ledger_index", "validated")
		limit := c.QueryInt("limit", 100)

		go accounts.StreamAccountNFTs(wsClient, account, ledgerIndex, limit, func(data *accounts.AccountNFTsResponse) {
			log.Printf("Real-time data: %+v", data)
		})

		return c.JSON(fiber.Map{"message": "Subscribed to account_nfts"})
	})

	// Similarly add routes for gateway balances
	app.Post("/accounts/balances/historical", func(c *fiber.Ctx) error {
		var payload accounts.GatewayBalancesRequest
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}

		response, err := httpClient.Post("", payload)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		var decodedResponse map[string]interface{}
		err = json.Unmarshal(response, &decodedResponse)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode JSON response"})
		}

		return c.JSON(decodedResponse)
	})

	// Gateway Balances - Real-time
	app.Get("/accounts/:account/balances/realtime", func(c *fiber.Ctx) error {
		account := c.Params("account")
		ledgerIndex := c.Query("ledger_index", "validated")
		hotWallet := c.Query("hotwallet", "") // Example: "wallet1,wallet2"
		strict := c.QueryBool("strict", false)

		hotWallets := []string{}
		if hotWallet != "" {
			hotWallets = append(hotWallets, hotWallet)
		}

		go accounts.StreamGatewayBalances(wsClient, account, hotWallets, ledgerIndex, strict, func(data *accounts.GatewayBalancesResponse) {
			log.Printf("Real-time data: %+v", data)
		})

		return c.JSON(fiber.Map{"message": "Subscribed to gateway_balances"})
	})

// Ledger Information - HTTP POST
	app.Post("/ledger", func(c *fiber.Ctx) error {
		var payload ledger.LedgerParam
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}

		response, err := ledger.FetchLedger(httpClient, payload.LedgerIndex, payload.Transactions, payload.Expand, payload.OwnerFunds)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		var result map[string]interface{}
		if err := json.Unmarshal(response, &result); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse response"})
		}

		return c.JSON(result)
	})

	// Most Recently Closed Ledger - HTTP POST
	app.Post("/ledger/closed", func(c *fiber.Ctx) error {
		response, err := ledger.FetchLedgerClosed(httpClient)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		var result map[string]interface{}
		if err := json.Unmarshal(response, &result); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse response"})
		}

		return c.JSON(result)
	})

	// Current Ledger Index - HTTP POST
	app.Post("/ledger/current", func(c *fiber.Ctx) error {
		response, err := ledger.FetchLedgerCurrent(httpClient)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		var result map[string]interface{}
		if err := json.Unmarshal(response, &result); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse response"})
		}

		return c.JSON(result)
	})

	// Ledger Data - HTTP POST
	app.Post("/ledger/data", func(c *fiber.Ctx) error {
		var payload ledger.LedgerDataParam
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}

		response, err := ledger.FetchLedgerData(httpClient, payload.LedgerHash, payload.Binary, payload.Limit, payload.Marker)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		var result map[string]interface{}
		if err := json.Unmarshal(response, &result); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse response"})
		}

		return c.JSON(result)
	})

	// ---------- WebSocket Routes ----------
	// Ledger Streaming Routes
	app.Get("/ledger/realtime", func(c *fiber.Ctx) error {
		ledgerIndex := c.Query("ledger_index", "validated")

		go ledger.StreamLedger(wsClient, ledgerIndex, func(data *ledger.LedgerWSResponse) {
			log.Printf("Ledger Real-Time Data: %+v", data)
		})


		return c.JSON(fiber.Map{"message": "Subscribed to ledger stream"})
	})

	app.Get("/ledger/closed/realtime", func(c *fiber.Ctx) error {
		go ledger.StreamLedgerClosed(wsClient, func(data *ledger.LedgerClosedWSResponse) {
			log.Printf("Ledger Closed Data: %+v", data)
		})

		return c.JSON(fiber.Map{"message": "Subscribed to ledger_closed stream"})
	})

	app.Get("/ledger/current/realtime", func(c *fiber.Ctx) error {
		go ledger.StreamLedgerCurrent(wsClient, func(data *ledger.LedgerCurrentWSResponse) {
			log.Printf("Current Ledger Data: %+v", data)
		})

		return c.JSON(fiber.Map{"message": "Subscribed to ledger_current stream"})
	})

	// Ledger Data Real-Time - WS
	app.Get("/ledger/data/realtime", func(c *fiber.Ctx) error {
		ledgerHash := c.Query("ledger_hash", "")
		binary := c.QueryBool("binary", true)
		limit := c.QueryInt("limit", 5)
		marker := c.Query("marker", "")

		go ledger.StreamLedgerData(wsClient, ledgerHash, binary, limit, marker, func(data *ledger.LedgerDataWSResponse) {
			log.Printf("Ledger Data Real-Time: %+v", data)
		})


		return c.JSON(fiber.Map{"message": "Subscribed to ledger_data stream"})
	})
}
