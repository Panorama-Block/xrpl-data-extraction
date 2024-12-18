package accounts

import (
	"encoding/json"
	"log"

	"github.com/Panorama-Block/xrpl-data-extraction/internal/xrpl"
)

// Fetch historical account NFTs using HTTP
func FetchHistoricalAccountNFTs(client *xrpl.HTTPClient, account string, ledgerIndex string, limit int) ([]byte, error) {
	params := AccountNFTsParam{
		Account: account,
		// LedgerIndex: ledgerIndex,
		// Limit: limit,
	}
	if ledgerIndex != "" {
		params.LedgerIndex = ledgerIndex
	}
	if limit != 0 {
		params.Limit = limit
	}

	payload := AccountNFTsRequest{
		Method: "account_nfts",
		Params: []AccountNFTsParam{params},
	}
	payloadJSON, _ := json.Marshal(payload) // Convert payload to JSON
	log.Printf("Sending payload: %s\n", payloadJSON)

	return client.Post("", payloadJSON)
}

// Stream real-time account NFTs using WebSocket
func StreamAccountNFTs(client *xrpl.WebSocketClient, account string, ledgerIndex string, limit int, callback func(*AccountNFTsResponse)) error {
	params := AccountNFTsWSRequest{
		Command: "account_nfts",
		Account: account,
		LedgerIndex: ledgerIndex,
		Limit: limit,
	}

	err := client.Subscribe(params)
	if err != nil {
		return err
	}

	client.ReadMessages(func(message []byte) {
		var response AccountNFTsResponse
		err := json.Unmarshal(message, &response)
		if err != nil {
			log.Println("Failed to unmarshal WebSocket response:", err)
			return
		}
		callback(&response)
	})

	return nil
}
