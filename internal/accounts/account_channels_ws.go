package accounts

import (
	"encoding/json"

	"github.com/Panorama-Block/xrpl-data-extraction/internal/xrpl"
)

// StreamRecentAccountChannels subscribes to real-time data using WebSocket
func StreamRecentAccountChannels(client *xrpl.WebSocketClient, account string, destinationAccount string, callback func(*AccountChannelsResponse)) error {
	// Build the params
	params := AccountChannelsParam{
		Account: account,
	}
	if destinationAccount != "" {
		params.DestinationAccount = destinationAccount
	}

	// Build the JSON-RPC payload
	payload := AccountChannelsRequest{
		Method: "account_channels",
		Params: []AccountChannelsParam{params},
	}

	// Subscribe to WebSocket
	err := client.Subscribe(payload)
	if err != nil {
		return err
	}

	// Read messages from WebSocket
	client.ReadMessages(func(message []byte) {
		var response AccountChannelsResponse
		err := json.Unmarshal(message, &response)
		if err != nil {
			return
		}
		callback(&response)
	})
	return nil
}
