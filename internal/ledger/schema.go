package ledger

import "time"

// LedgerSchema define os campos que serão salvos no MongoDB
type LedgerSchema struct {
	LedgerIndex     int    `bson:"ledger_index"`
	LedgerHash      string    `bson:"ledger_hash"`
	CloseTimeHuman  string    `bson:"close_time_human"`
	TxnCount				int       `bson:"txn_count"`
	FeeBase				 int       `bson:"fee_base"`
	TotalCoins			string    `bson:"total_coins"`
	CreatedAt       time.Time `bson:"created_at"`
}


