package transactions

// ---------- HTTP Types ----------

// TransactionEntryRequest defines the structure for HTTP /transaction_entry
type TransactionEntryRequest struct {
	Method string                  `json:"method"`
	Params []TransactionEntryParam `json:"params"`
}

type TransactionEntryParam struct {
	TxHash      string `json:"tx_hash"`
	LedgerIndex string `json:"ledger_index,omitempty"`
}

// TransactionRequest defines the structure for HTTP /tx
type TransactionRequest struct {
	Method string              `json:"method"`
	Params []TransactionParam `json:"params"`
}

type TransactionParam struct {
	Transaction string `json:"transaction"`
	Binary      bool   `json:"binary,omitempty"`
	APIVersion  int    `json:"api_version,omitempty"`
}

// TransactionEntryResponse defines the HTTP response structure for /transaction_entry
type TransactionEntryResponse struct {
	Result struct {
		TxJSON map[string]interface{} `json:"tx_json"`
		Meta   map[string]interface{} `json:"meta"`
	} `json:"result"`
}

// TransactionResponse defines the HTTP response structure for /tx
type TransactionResponse struct {
	Result struct {
		TxJSON map[string]interface{} `json:"tx_json"`
		Meta   map[string]interface{} `json:"meta"`
	} `json:"result"`
}

// ---------- WebSocket Types ----------

// TransactionEntryWSRequest defines the WebSocket request for transaction_entry
type TransactionEntryWSRequest struct {
	ID          int    `json:"id"`
	Command     string `json:"command"`
	TxHash      string `json:"tx_hash"`
	LedgerIndex string `json:"ledger_index,omitempty"`
}

// TransactionWSRequest defines the WebSocket request for tx
type TransactionWSRequest struct {
	ID          int    `json:"id"`
	Command     string `json:"command"`
	Transaction string `json:"transaction"`
	Binary      bool   `json:"binary,omitempty"`
	APIVersion  int    `json:"api_version,omitempty"`
}

// TransactionEntryWSResponse defines the WebSocket response for transaction_entry
type TransactionEntryWSResponse struct {
	ID     int                    `json:"id"`
	Status string                 `json:"status"`
	Result map[string]interface{} `json:"result"`
}

// TransactionWSResponse defines the WebSocket response for tx
type TransactionWSResponse struct {
	ID     int                    `json:"id"`
	Status string                 `json:"status"`
	Result map[string]interface{} `json:"result"`
}
