package orderbook

import (
	"encoding/json"
	"github.com/Panorama-Block/xrpl-data-extraction/internal/xrpl"
)

// FetchBookChanges fetches book changes via HTTP
func FetchBookChanges(client *xrpl.HTTPClient, ledgerIndex int) ([]byte, error) {
	params := BookChangesParam{
		LedgerIndex: ledgerIndex,
	}

	request := BookChangesRequest{
		Method: "book_changes",
		Params: []BookChangesParam{params},
	}

	return client.Post("", request)
}

// StreamBookChanges streams book changes via WebSocket
func StreamBookChanges(wsClient *xrpl.WebSocketClient, ledgerIndex int, callback func(*BookChangesWSResponse)) error {
	request := BookChangesWSRequest{
		ID:          2,
		Command:     "book_changes",
		LedgerIndex: ledgerIndex,
	}

	err := wsClient.Subscribe(request)
	if err != nil {
		return err
	}

	wsClient.ReadMessages(func(msg []byte) {
		var response BookChangesWSResponse
		if err := json.Unmarshal(msg, &response); err == nil {
			callback(&response)
		}
	})
	return nil
}
