package accounts

import (
	"encoding/json"
	"log"

	"github.com/Panorama-Block/xrpl-data-extraction/internal/xrpl"
)

// Fetch historical account channels data from the XRPL using HTTP client
func FetchHistoricalAccountChannels(client *xrpl.HTTPClient, account string, destinationAccount string, ledgerIndex string) ([]byte, error) {
	// Build the params
	params := AccountChannelsParam{
		Account: account,
	}
	if destinationAccount != "" {
		params.DestinationAccount = destinationAccount
	}
	if ledgerIndex != "" {
		params.LedgerIndex = ledgerIndex
	}
	// Build the JSON-RPC payload
	payload := AccountChannelsRequest{
		Method: "account_channels",
		Params: []AccountChannelsParam{params},
	}
	payloadJSON, _ := json.Marshal(payload) // Convert payload to JSON
	log.Printf("Sending payload: %s\n", payloadJSON)
	// Send the request
	return client.Post("", payload)
}

// StreamRecentAccountChannels subscribes to real-time data using WebSocket
func StreamRecentAccountChannels(client *xrpl.WebSocketClient, account string, destinationAccount string, callback func(*AccountChannelsWSResponse)) error {
	// Build the params
	params := AccountChannelsParam{
		Account: account,
	}
	if destinationAccount != "" {
		params.DestinationAccount = destinationAccount
	}

	// Build the JSON-RPC payload
	payload := AccountChannelsWSRequest{
		ID:                1,
		Command:           "account_channels",
		Account:           account,
		DestinationAccount: destinationAccount,
		LedgerIndex:       "validated",
	}

	// Subscribe to WebSocket
	err := client.Subscribe(payload)
	if err != nil {
		return err
	}

	// Read messages from WebSocket
	client.ReadMessages(func(message []byte) {
		var response AccountChannelsWSResponse
		err := json.Unmarshal(message, &response)
		if err != nil {
			return
		}
		callback(&response)
	})
	return nil
}
