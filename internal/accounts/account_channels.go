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
