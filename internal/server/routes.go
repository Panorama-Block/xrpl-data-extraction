package server

import (
	"log" 
	"sync"
	"encoding/json" //for json operations

	"github.com/Panorama-Block/xrpl-data-extraction/internal/accounts"
	"github.com/Panorama-Block/xrpl-data-extraction/internal/xrpl"
	"github.com/Panorama-Block/xrpl-data-extraction/internal/ledger"
	"github.com/Panorama-Block/xrpl-data-extraction/internal/transactions"
	"github.com/Panorama-Block/xrpl-data-extraction/internal/orderbook"
	"github.com/Panorama-Block/xrpl-data-extraction/internal/states"
	"github.com/gofiber/fiber/v2"
)

// Controle de conexões WebSocket ativas
var (
	stopChans = make(map[string]chan struct{})
	mu        sync.Mutex
)


// setup routes
func SetupRoutes(app *fiber.App, httpClient *xrpl.HTTPClient, wsClient *xrpl.WebSocketClient) {

// Inscrição em contas e monitoramento de transações em tempo real
app.Post("/accounts/subscribe", func(c *fiber.Ctx) error {
	// Extrair lista de contas do corpo da requisição
	var payload struct {
		Accounts []string `json:"accounts"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if len(payload.Accounts) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Accounts list cannot be empty",
		})
	}

	// Criar canal para gerenciar o stream
	mu.Lock()
	stopChan := make(chan struct{})
	stopChans["accounts_subscribe"] = stopChan
	mu.Unlock()

	// Iniciar o streaming de contas
	go accounts.SubscribeAccounts(wsClient, payload.Accounts, stopChan)

	return c.JSON(fiber.Map{
		"message":  "Subscription to accounts started",
		"accounts": payload.Accounts,
	})
})

// Cancelar inscrição no stream de contas
app.Delete("/accounts/subscribe", func(c *fiber.Ctx) error {
	mu.Lock()
	stopChan, exists := stopChans["accounts_subscribe"]
	if exists {
		close(stopChan) // Enviar sinal de parada
		delete(stopChans, "accounts_subscribe")
		mu.Unlock()
		return c.JSON(fiber.Map{"message": "Subscription to accounts stopped"})
	}
	mu.Unlock()

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": "No active subscription to accounts",
	})
})



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

	
	
	// ==================================================================================================LEDGER DATA===============================================================================================================
app.Post("/ledger", func(c *fiber.Ctx) error {
	var payload ledger.LedgerParam
	if err := c.BodyParser(&payload); err != nil {
		log.Printf("❌ Erro ao interpretar o payload: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}
	if payload.LedgerIndex == "" {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "LedgerIndex é obrigatório"})
	}
	response, err := ledger.FetchLedgerInfo(httpClient, payload.LedgerIndex)
	if err != nil {
		log.Printf("❌ Erro ao buscar dados do ledger: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching ledger data"})
	}
	return c.JSON(response)
})
	
var stopChan chan struct{}
	
app.Get("/ledger/realtime", func(c *fiber.Ctx) error {
	stopChan = make(chan struct{})
	go ledger.StreamLedger(wsClient, httpClient, func(data *ledger.LedgerSubscribeClosedResponse) {
		
	}, stopChan)
	
	return c.JSON(fiber.Map{"message": "📡 Streaming de ledgers iniciado!"})
})

app.Get("/ledger/stop", func(c *fiber.Ctx) error {
	close(stopChan)
	return c.JSON(fiber.Map{"message": "⛔ Streaming de ledgers encerrado!"})
})


// ==================================================================================================TRANSACTIONS===============================================================================================================
app.Post("/transactions/entry", func(c *fiber.Ctx) error {
	var payload transactions.TransactionEntryParam
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}
	
	response, err := transactions.FetchTransactionEntry(httpClient, payload.TxHash, payload.LedgerIndex)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse response"})
	}
	
	return c.JSON(result)
})

// Transaction Details - HTTP POST
app.Post("/transactions", func(c *fiber.Ctx) error {
	var payload transactions.TransactionParam
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}
	
	response, err := transactions.FetchTransaction(httpClient, payload.Transaction, payload.Binary)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse response"})
	}
	return c.JSON(result)
})

// Transaction Entry - WebSocket
app.Get("/transactions/entry/realtime", func(c *fiber.Ctx) error {
	txHash := c.Query("tx_hash")
	ledgerIndex := c.Query("ledger_index", "")
	
	go transactions.StreamTransactionEntry(wsClient, txHash, ledgerIndex, func(data *transactions.TransactionEntryWSResponse) {
		log.Printf("Transaction Entry Real-Time: %+v", data)
	})
	return c.JSON(fiber.Map{"message": "Subscribed to transaction_entry"})
})

	// Transaction Details - WebSocket
	app.Get("/transactions/realtime", func(c *fiber.Ctx) error {
		txHash := c.Query("transaction")
		binary := c.QueryBool("binary", false)
		
		go transactions.StreamTransaction(wsClient, txHash, binary, func(data *transactions.TransactionWSResponse) {
			log.Printf("Transaction Real-Time: %+v", data)
		})
		
		return c.JSON(fiber.Map{"message": "Subscribed to tx"})
	})
	
	// AMM Info - HTTP and WebSocket
	app.Post("/orderbook/amm_info", func(c *fiber.Ctx) error {
		var request struct {
			AMMAccount string     `json:"amm_account,omitempty"`
			Asset      orderbook.AssetParam `json:"asset,omitempty"`
			Asset2     orderbook.AssetParam `json:"asset2,omitempty"`
		}
		
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
		}
		
		response, err := orderbook.FetchAMMInfo(httpClient, request.AMMAccount, request.Asset, request.Asset2)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		
		var result map[string]interface{}
		if err := json.Unmarshal(response, &result); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode response"})
		}
		
		return c.JSON(result)
	})
	
	app.Get("/orderbook/amm_info/realtime", func(c *fiber.Ctx) error {
		ammAccount := c.Query("amm_account", "")
		var asset, asset2 orderbook.AssetParam
		c.QueryParser(&asset)
		c.QueryParser(&asset2)
		
		go orderbook.StreamAMMInfo(wsClient, ammAccount, asset, asset2, func(data *orderbook.AMMInfoWSResponse) {
			log.Printf("AMM Info Real-Time: %+v", data)
		})
		
		return c.JSON(fiber.Map{"message": "Subscribed to AMM Info"})
	})
	
	// Book Changes - HTTP and WebSocket
	app.Post("/orderbook/book_changes", func(c *fiber.Ctx) error {
		var payload struct {
			LedgerIndex int `json:"ledger_index"`
		}
		
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
		}
		
		response, err := orderbook.FetchBookChanges(httpClient, payload.LedgerIndex)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		
		var result map[string]interface{}
		if err := json.Unmarshal(response, &result); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode response"})
		}
		
		return c.JSON(result)
	})
	
	app.Get("/orderbook/book_changes/realtime", func(c *fiber.Ctx) error {
		ledgerIndex := c.QueryInt("ledger_index", 0)
		
		go orderbook.StreamBookChanges(wsClient, ledgerIndex, func(data *orderbook.BookChangesWSResponse) {
			log.Printf("Book Changes Real-Time: %+v", data)
		})
		
		return c.JSON(fiber.Map{"message": "Subscribed to Book Changes"})
	})
	
	app.Post("/orderbook/book_offers", func(c *fiber.Ctx) error {
		var params orderbook.BookOffersParams
		if err := c.BodyParser(&params); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}
		
		response, err := orderbook.FetchBookOffers(httpClient, params)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		var decodedResponse map[string]interface{}
		if err := json.Unmarshal(response, &decodedResponse); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode response"})
		}
		
		return c.JSON(decodedResponse)
	})

	// Book Offers - WebSocket
	app.Get("/orderbook/book_offers/realtime", func(c *fiber.Ctx) error {
		params := orderbook.BookOffersParams{
			Taker: c.Query("taker"),
			TakerGets: orderbook.AssetParam{
				Currency: c.Query("taker_gets_currency"),
				Issuer:   c.Query("taker_gets_issuer"),
			},
			TakerPays: orderbook.AssetParam{
				Currency: c.Query("taker_pays_currency"),
				Issuer:   c.Query("taker_pays_issuer"),
			},
			Limit: c.QueryInt("limit", 10),
		}
		
		go orderbook.StreamBookOffers(wsClient, params, func(data *orderbook.BookOffersWSResponse) {
			log.Printf("Book Offers Data: %+v", data)
		})
		
		return c.JSON(fiber.Map{"message": "Subscribed to book_offers stream"})
	})
	
	// Get Aggregate Price - HTTP
	app.Post("/orderbook/aggregate_price", func(c *fiber.Ctx) error {
		var params orderbook.GetAggregatePriceParams
		if err := c.BodyParser(&params); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}

		response, err := orderbook.FetchAggregatePrice(httpClient, params)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		
		var decodedResponse map[string]interface{}
		if err := json.Unmarshal(response, &decodedResponse); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode response"})
		}
		
		return c.JSON(decodedResponse)
	})
	
	// Get Aggregate Price - WebSocket
	app.Get("/orderbook/aggregate_price/realtime", func(c *fiber.Ctx) error {
		params := orderbook.GetAggregatePriceParams{
			LedgerIndex: c.Query("ledger_index"),
			BaseAsset:   c.Query("base_asset"),
			QuoteAsset:  c.Query("quote_asset"),
			Trim:        c.QueryInt("trim", 20),
			Oracles:     []orderbook.Oracle{}, // Fill dynamically from request if needed
		}
		
		go orderbook.StreamAggregatePrice(wsClient, params, func(data *orderbook.GetAggregatePriceResponse) {
			log.Printf("Aggregate Price Data: %+v", data)
		})
		
		return c.JSON(fiber.Map{"message": "Subscribed to aggregate_price stream"})
	})

	// NFT Buy Offers - HTTP
	app.Post("/orderbook/nft_buy_offers", func(c *fiber.Ctx) error {
		var params orderbook.NFTBuyOffersParams
		if err := c.BodyParser(&params); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}

		response, err := orderbook.FetchNFTBuyOffers(httpClient, params)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		
		var decodedResponse map[string]interface{}
		if err := json.Unmarshal(response, &decodedResponse); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode response"})
		}
		
		return c.JSON(decodedResponse)
	})
	
	// NFT Buy Offers - WebSocket
	app.Get("/orderbook/nft_buy_offers/realtime", func(c *fiber.Ctx) error {
		params := orderbook.NFTBuyOffersWSRequest{
			ID:        1,
			Command:   "nft_buy_offers",
			NFTID:     c.Query("nft_id"),
			LedgerIndex: c.Query("ledger_index", "validated"),
		}
		
		go wsClient.Subscribe(params)
		
		wsClient.ReadMessages(func(message []byte) {
			var response orderbook.NFTBuyOffersWSResponse
			if err := json.Unmarshal(message, &response); err == nil {
				log.Printf("NFT Buy Offers Real-Time: %+v", response)
			}
		})

		return c.JSON(fiber.Map{"message": "Subscribed to nft_buy_offers"})
	})

	// NFT Sell Offers - HTTP
	app.Post("/orderbook/nft_sell_offers", func(c *fiber.Ctx) error {
		var params orderbook.NFTSellOffersParams
		if err := c.BodyParser(&params); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}

		response, err := orderbook.FetchNFTSellOffers(httpClient, params)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		
		var decodedResponse map[string]interface{}
		if err := json.Unmarshal(response, &decodedResponse); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode response"})
		}
		
		return c.JSON(decodedResponse)
	})
	
	// NFT Sell Offers - WebSocket
	app.Get("/orderbook/nft_sell_offers/realtime", func(c *fiber.Ctx) error {
		params := orderbook.NFTSellOffersWSRequest{
			ID:        1,
			Command:   "nft_sell_offers",
			NFTID:     c.Query("nft_id"),
			LedgerIndex: c.Query("ledger_index", "validated"),
		}

		go wsClient.Subscribe(params)
		
		wsClient.ReadMessages(func(message []byte) {
			var response orderbook.NFTSellOffersWSResponse
			if err := json.Unmarshal(message, &response); err == nil {
				log.Printf("NFT Sell Offers Real-Time: %+v", response)
			}
		})

		return c.JSON(fiber.Map{"message": "Subscribed to nft_sell_offers"})
	})
	
	// Fee - HTTP
	app.Post("/server/fee", func(c *fiber.Ctx) error {
		response, err := serverinfo.FetchFee(httpClient)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		var result map[string]interface{}
		if err := json.Unmarshal(response, &result); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode response"})
		}
		
		return c.JSON(result)
	})

	// Fee - WebSocket
	app.Get("/server/fee/realtime", func(c *fiber.Ctx) error {
		go serverinfo.StreamFee(wsClient, func(data *serverinfo.FeeWSResponse) {
			log.Printf("Fee Real-Time Data: %+v", data)
		})
		
		return c.JSON(fiber.Map{"message": "Subscribed to fee stream"})
	})
	
	// Server State - HTTP
	app.Post("/server/state", func(c *fiber.Ctx) error {
		response, err := serverinfo.FetchServerState(httpClient)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		
		var result map[string]interface{}
		if err := json.Unmarshal(response, &result); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode response"})
		}
		
		return c.JSON(result)
	})
	
	// Server State - WebSocket
	app.Get("/server/state/realtime", func(c *fiber.Ctx) error {
		go serverinfo.StreamServerState(wsClient, func(data *serverinfo.ServerStateWSResponse) {
			log.Printf("Server State Real-Time Data: %+v", data)
		})

		return c.JSON(fiber.Map{"message": "Subscribed to server_state stream"})
	})
}
