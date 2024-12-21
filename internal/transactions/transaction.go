package transactions

import (
	"encoding/json"
	"github.com/Panorama-Block/xrpl-data-extraction/internal/xrpl"
)

// FetchTransactionEntry fetches transaction information for a specific ledger version
func FetchTransactionEntry(client *xrpl.HTTPClient, txHash string, ledgerIndex string) ([]byte, error) {
	params := TransactionEntryParam{
		TxHash:      txHash,
		LedgerIndex: ledgerIndex,
	}

	request := TransactionEntryRequest{
		Method: "transaction_entry",
		Params: []TransactionEntryParam{params},
	}

	return client.Post("", request)
}

// FetchTransaction retrieves transaction details by transaction hash
func FetchTransaction(client *xrpl.HTTPClient, txHash string, binary bool) ([]byte, error) {
	params := TransactionParam{
		Transaction: txHash,
		Binary:      binary,
		APIVersion:  2,
	}

	request := TransactionRequest{
		Method: "tx",
		Params: []TransactionParam{params},
	}

	return client.Post("", request)
}

// StreamTransactionEntry fetches transaction entry data in real-time via WebSocket
func StreamTransactionEntry(wsClient *xrpl.WebSocketClient, txHash string, ledgerIndex string, callback func(*TransactionEntryWSResponse)) error {
	request := TransactionEntryWSRequest{
		ID:          5,
		Command:     "transaction_entry",
		TxHash:      txHash,
		LedgerIndex: ledgerIndex,
	}

	err := wsClient.Subscribe(request)
	if err != nil {
		return err
	}

	wsClient.ReadMessages(func(msg []byte) {
		var response TransactionEntryWSResponse
		if err := json.Unmarshal(msg, &response); err == nil {
			callback(&response)
		}
	})
	return nil
}

// StreamTransaction fetches transaction data in real-time via WebSocket
func StreamTransaction(wsClient *xrpl.WebSocketClient, txHash string, binary bool, callback func(*TransactionWSResponse)) error {
	request := TransactionWSRequest{
		ID:          6,
		Command:     "tx",
		Transaction: txHash,
		Binary:      binary,
		APIVersion:  2,
	}

	err := wsClient.Subscribe(request)
	if err != nil {
		return err
	}

	wsClient.ReadMessages(func(msg []byte) {
		var response TransactionWSResponse
		if err := json.Unmarshal(msg, &response); err == nil {
			callback(&response)
		}
	})
	return nil
}
