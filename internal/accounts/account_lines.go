package accounts

import (
	"encoding/json"
	"log"

	"github.com/Panorama-Block/xrpl-data-extraction/internal/xrpl"
)

// FetchAccountLines fetches trust lines for an account using HTTP
func FetchAccountLines(client *xrpl.HTTPClient, account string, ledgerIndex string, limit int, marker string) ([]byte, error) {
	// Build the parameters
	params := AccountLinesParam{
		Account: account,
	}
	if ledgerIndex != "" {
		params.LedgerIndex = ledgerIndex
	}

	// Build the JSON-RPC payload
	payload := AccountLinesRequest{
		Method: "account_lines",
		Params: []AccountLinesParam{params},
	}

	payloadJSON, _ := json.Marshal(payload) // Convert payload to JSON
	log.Printf("Sending payload: %s\n", payloadJSON)

	// Send the request
	return client.Post("", payload)
}

// StreamAccountLines subscribes to real-time trust lines using WebSocket
func StreamAccountLines(client *xrpl.WebSocketClient, account string, ledgerIndex string, limit int, callback func(*AccountLinesWSResponse)) error {
	// Build the parameters
	params := AccountLinesParam{
		Account: account,
	}
	if ledgerIndex != "" {
		params.LedgerIndex = ledgerIndex
	}

	// Build the JSON-RPC payload
	payload := AccountLinesWSRequest{
		ID:          1,
		Command:     "account_lines",
		Account:     account,
		LedgerIndex: ledgerIndex,
	}

	// Subscribe to WebSocket
	err := client.Subscribe(payload)
	if err != nil {
		return err
	}

	// Read messages from WebSocket
	client.ReadMessages(func(message []byte) {
		var response AccountLinesWSResponse
		err := json.Unmarshal(message, &response)
		if err != nil {
			return
		}
		callback(&response)
	})
	return nil
}
