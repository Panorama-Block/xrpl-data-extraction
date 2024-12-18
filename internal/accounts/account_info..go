package accounts

import (
	"encoding/json"
	"log"

	"github.com/Panorama-Block/xrpl-data-extraction/internal/xrpl"
)

// FetchAccountInfo fetches historical account information via HTTP
func FetchAccountInfo(client *xrpl.HTTPClient, account string, ledgerIndex string, queue bool) ([]byte, error) {
	// Build the parameters
	params := AccountInfoParam{
		Account: account,
		Queue:   queue,
	}
	if ledgerIndex != "" {
		params.LedgerIndex = ledgerIndex
	}

	// Build the JSON-RPC payload
	payload := AccountInfoRequest{
		Method: "account_info",
		Params: []AccountInfoParam{params},
	}

	payloadJSON, _ := json.Marshal(payload) // Convert payload to JSON
	log.Printf("Sending payload: %s\n", payloadJSON)

	// Send the request
	return client.Post("", payload)
}

// StreamAccountInfo subscribes to real-time account information using WebSocket
func StreamAccountInfo(client *xrpl.WebSocketClient, account string, ledgerIndex string, queue bool, callback func(*AccountInfoWSResponse)) error {
	// Build the parameters
	params := AccountInfoParam{
		Account: account,
		Queue:   queue,
	}
	if ledgerIndex != "" {
		params.LedgerIndex = ledgerIndex
	}

	// Build the JSON-RPC payload
	payload := AccountInfoWSRequest{
		ID:          2,
		Command:     "account_info",
		Account:     account,
		LedgerIndex: ledgerIndex,
		Queue:       queue,
	}

	// Subscribe to WebSocket
	err := client.Subscribe(payload)
	if err != nil {
		return err
	}

	// Read messages from WebSocket
	client.ReadMessages(func(message []byte) {
		var response AccountInfoWSResponse
		err := json.Unmarshal(message, &response)
		if err != nil {
			return
		}
		callback(&response)
	})
	return nil
}
