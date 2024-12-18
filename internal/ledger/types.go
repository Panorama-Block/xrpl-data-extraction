package ledger

// ---------- HTTP Types ----------

// LedgerRequest defines the structure for HTTP /ledger requests
type LedgerRequest struct {
	Method string        `json:"method"`
	Params []LedgerParam `json:"params"`
}

type LedgerParam struct {
	LedgerIndex  string `json:"ledger_index,omitempty"`
	Transactions bool   `json:"transactions,omitempty"`
	Expand       bool   `json:"expand,omitempty"`
	OwnerFunds   bool   `json:"owner_funds,omitempty"`
}

type LedgerResponse struct {
	Result struct {
		Ledger struct {
			AccountHash     string `json:"account_hash"`
			CloseTime       int64  `json:"close_time"`
			CloseTimeHuman  string `json:"close_time_human"`
			LedgerHash      string `json:"ledger_hash"`
			LedgerIndex     string `json:"ledger_index"`
			ParentHash      string `json:"parent_hash"`
			TotalCoins      string `json:"total_coins"`
			TransactionHash string `json:"transaction_hash"`
		} `json:"ledger"`
		Validated bool `json:"validated"`
	} `json:"result"`
}

type LedgerClosedRequest struct {
	Method string   `json:"method"`
	Params []struct{} `json:"params"`
}

type LedgerClosedResponse struct {
	Result struct {
		LedgerHash  string `json:"ledger_hash"`
		LedgerIndex int    `json:"ledger_index"`
	} `json:"result"`
}

type LedgerCurrentRequest struct {
	Method string   `json:"method"`
	Params []struct{} `json:"params"`
}

type LedgerCurrentResponse struct {
	Result struct {
		LedgerCurrentIndex int `json:"ledger_current_index"`
	} `json:"result"`
}

type LedgerDataRequest struct {
	Method string            `json:"method"`
	Params []LedgerDataParam `json:"params"`
}

type LedgerDataParam struct {
	Binary     bool   `json:"binary,omitempty"`
	LedgerHash string `json:"ledger_hash,omitempty"`
	LedgerIndex string `json:"ledger_index,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	Marker     string `json:"marker,omitempty"`
}

type LedgerDataResponse struct {
	Result struct {
		LedgerHash string `json:"ledger_hash"`
		Marker     string `json:"marker"`
		State      []struct {
			Data  string `json:"data"`
			Index string `json:"index"`
		} `json:"state"`
	} `json:"result"`
}

// ---------- WebSocket Types ----------

// WS Ledger Request
type LedgerWSRequest struct {
	ID          int    `json:"id"`
	Command     string `json:"command"`
	LedgerIndex string `json:"ledger_index,omitempty"`
	Transactions bool  `json:"transactions,omitempty"`
	Expand       bool  `json:"expand,omitempty"`
	OwnerFunds   bool  `json:"owner_funds,omitempty"`
}

// WS Ledger Response
type LedgerWSResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Result struct {
		Ledger struct {
			AccountHash     string `json:"account_hash"`
			CloseTime       int64  `json:"close_time"`
			CloseTimeHuman  string `json:"close_time_human"`
			LedgerHash      string `json:"ledger_hash"`
			LedgerIndex     string `json:"ledger_index"`
			ParentHash      string `json:"parent_hash"`
			TotalCoins      string `json:"total_coins"`
			TransactionHash string `json:"transaction_hash"`
		} `json:"ledger"`
		Validated bool `json:"validated"`
	} `json:"result"`
}

// WS Ledger Closed Request
type LedgerClosedWSRequest struct {
	ID      int    `json:"id"`
	Command string `json:"command"`
}

// WS Ledger Closed Response
type LedgerClosedWSResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Result struct {
		LedgerHash  string `json:"ledger_hash"`
		LedgerIndex int    `json:"ledger_index"`
	} `json:"result"`
}

// WS Ledger Current Request
type LedgerCurrentWSRequest struct {
	ID      int    `json:"id"`
	Command string `json:"command"`
}

// WS Ledger Current Response
type LedgerCurrentWSResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Result struct {
		LedgerCurrentIndex int `json:"ledger_current_index"`
	} `json:"result"`
}

// WS Ledger Data Request
type LedgerDataWSRequest struct {
	ID          int    `json:"id"`
	Command     string `json:"command"`
	LedgerHash  string `json:"ledger_hash,omitempty"`
	LedgerIndex string `json:"ledger_index,omitempty"`
	Binary      bool   `json:"binary,omitempty"`
	Limit       int    `json:"limit,omitempty"`
	Marker      string `json:"marker,omitempty"`
}

// WS Ledger Data Response
type LedgerDataWSResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Result struct {
		LedgerHash string `json:"ledger_hash"`
		Marker     string `json:"marker"`
		State      []struct {
			Data  string `json:"data"`
			Index string `json:"index"`
		} `json:"state"`
	} `json:"result"`
}
