package accounts

import "time"

// TransactionSchema define os campos para salvar transações no MongoDB
type TransactionSchema struct {
	Account      string `bson:"account"`
	Fee          string `bson:"fee"`
	TakerGets    string `bson:"taker_gets"` // Agora string
	TakerPays    string `bson:"taker_pays"` // Agora string
	Date         time.Time `bson:"date"`
	OwnerFunds   string    `bson:"owner_funds"`
	Type         string    `bson:"type"`
	Validated    bool      `bson:"validated"`
	Status       string    `bson:"status"`
	LedgerHash   string    `bson:"ledger_hash"`
	LedgerIndex  int       `bson:"ledger_index"`
	CreatedAt    time.Time `bson:"created_at"`
}
