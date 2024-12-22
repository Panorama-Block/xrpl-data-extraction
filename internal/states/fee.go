package serverinfo

import (
	"encoding/json"
	"github.com/Panorama-Block/xrpl-data-extraction/internal/xrpl"
)

// FetchFee fetches the current state of fees via HTTP
func FetchFee(client *xrpl.HTTPClient) ([]byte, error) {
	request := struct {
		Method string   `json:"method"`
		Params []struct{} `json:"params"`
	}{
		Method: "fee",
		Params: []struct{}{{}}, // Empty object in "params"
	}

	return client.Post("", request)
}

// StreamFee streams the current state of fees via WebSocket
func StreamFee(wsClient *xrpl.WebSocketClient, callback func(*FeeWSResponse)) error {
	request := struct {
		ID      string `json:"id"`
		Command string `json:"command"`
	}{
		ID:      "fee_websocket_example",
		Command: "fee",
	}

	err := wsClient.Subscribe(request)
	if err != nil {
		return err
	}

	wsClient.ReadMessages(func(msg []byte) {
		var response FeeWSResponse
		if err := json.Unmarshal(msg, &response); err == nil {
			callback(&response)
		}
	})
	return nil
}

// FeeWSResponse represents the WebSocket response for the fee command
type FeeWSResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Result FeeResponse `json:"result"`
}

// FeeResponse defines the structure for the fee command response
type FeeResponse struct {
	CurrentLedgerSize string `json:"current_ledger_size"`
	CurrentQueueSize  string `json:"current_queue_size"`
	Drops             struct {
		BaseFee       string `json:"base_fee"`
		MedianFee     string `json:"median_fee"`
		MinimumFee    string `json:"minimum_fee"`
		OpenLedgerFee string `json:"open_ledger_fee"`
	} `json:"drops"`
	ExpectedLedgerSize string `json:"expected_ledger_size"`
	LedgerCurrentIndex int    `json:"ledger_current_index"`
	Levels             struct {
		MedianLevel     string `json:"median_level"`
		MinimumLevel    string `json:"minimum_level"`
		OpenLedgerLevel string `json:"open_ledger_level"`
		ReferenceLevel  string `json:"reference_level"`
	} `json:"levels"`
	MaxQueueSize string `json:"max_queue_size"`
	Status       string `json:"status"`
}
