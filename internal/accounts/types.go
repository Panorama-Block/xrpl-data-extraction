package accounts

type AccountChannelsRequest struct {
	Method string                 `json:"method"`
	Params []AccountChannelsParam `json:"params"`
}

type AccountChannelsParam struct {
	Account            string `json:"account"`
	DestinationAccount string `json:"destination_account,omitempty"`
	LedgerIndex        string `json:"ledger_index,omitempty"`
}

type AccountChannelsResponse struct {
	ID      int    `json:"id"`
	Status  string `json:"status"`
	Type    string `json:"type"`
	Result  struct {
		Account      string    `json:"account"`
		Channels     []Channel `json:"channels"`
		LedgerIndex  int64     `json:"ledger_index"`
		Validated    bool      `json:"validated"`
	} `json:"result"`
}

type Channel struct {
	Account           string `json:"account"`
	Amount            string `json:"amount"`
	Balance           string `json:"balance"`
	ChannelID         string `json:"channel_id"`
	DestinationAccount string `json:"destination_account"`
	SettleDelay       int    `json:"settle_delay"`
}
