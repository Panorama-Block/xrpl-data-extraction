package ledger

import "time"

// LedgerSchema define os campos que serão salvos no MongoDB
type LedgerSchema struct {
	LedgerIndex     int    `bson:"ledger_index"`
	LedgerHash      string    `bson:"ledger_hash"`
	CloseTimeHuman  string    `bson:"close_time_human"`
	TxnCount				int       `bson:"txn_count"`
	FeeBase				 int       `bson:"fee_base"`
	CreatedAt       time.Time `bson:"created_at"`
}

// // Solicitação WebSocket para filtrar os campos
// type LedgerWSRequest2 struct {
// 	ID          int    `json:"id"`
// 	Command     string `json:"command"`
// 	LedgerIndex string `json:"ledger_index,omitempty"`
// 	Transactions bool  `json:"transactions,omitempty"`
// 	Expand       bool  `json:"expand,omitempty"`
// 	OwnerFunds   bool  `json:"owner_funds,omitempty"`
// }

// // Resposta WebSocket já filtrada
// type LedgerWSResponse2 struct {
// 	ID     int    `json:"id"`
// 	Status string `json:"status"`
// 	Type   string `json:"type"`
// 	Result struct {
// 		Ledger struct {
// 			AccountHash     string `json:"account_hash"`
// 			CloseTime       int64  `json:"close_time"`
// 			CloseTimeHuman  string `json:"close_time_human"`
// 			LedgerHash      string `json:"ledger_hash"`
// 			LedgerIndex     string `json:"ledger_index"`
// 			ParentHash      string `json:"parent_hash"`
// 			TotalCoins      string `json:"total_coins"`
// 			TransactionHash string `json:"transaction_hash"`
// 		} `json:"ledger"`
// 		Validated bool `json:"validated"`
// 	} `json:"result"`
// }

