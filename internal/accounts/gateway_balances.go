package accounts

import (
	"encoding/json"
	"log"

	"github.com/Panorama-Block/xrpl-data-extraction/internal/xrpl"
)

// Fetch historical gateway balances using HTTP
func FetchGatewayBalances(client *xrpl.HTTPClient, account string, hotWallet []string, ledgerIndex string, strict bool) ([]byte, error) {
	params := GatewayBalancesParam{
		Account:     account,
		// HotWallet:   hotWallet,
		// LedgerIndex: ledgerIndex,
		// Strict:      strict,
	}

	if len(hotWallet) > 0 {
		params.HotWallet = hotWallet
	}
	if ledgerIndex != "" {
		params.LedgerIndex = ledgerIndex
	}
	params.Strict = strict 


	payload := GatewayBalancesRequest{
		Method: "gateway_balances",
		Params: []GatewayBalancesParam{params},
	}

	payloadJSON, _ := json.Marshal(payload) // Convert payload to JSON
	log.Printf("Sending payload: %s\n",
		payloadJSON)

	return client.Post("", payloadJSON)
}

// Stream real-time gateway balances using WebSocket
func StreamGatewayBalances(client *xrpl.WebSocketClient, account string, hotWallet []string, ledgerIndex string, strict bool, callback func(*GatewayBalancesResponse)) error {
	params := GatewayBalancesWSRequest{
		Command:    "gateway_balances",
		Account:    account,
		HotWallet:  hotWallet,
		LedgerIndex: ledgerIndex,
		Strict:      strict,
	}

	err := client.Subscribe(params)
	if err != nil {
		return err
	}

	client.ReadMessages(func(message []byte) {
		var response GatewayBalancesResponse
		err := json.Unmarshal(message, &response)
		if err != nil {
			log.Println("Failed to unmarshal WebSocket response:", err)
			return
		}
		callback(&response)
	})

	return nil
}
