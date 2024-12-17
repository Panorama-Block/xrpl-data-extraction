package account_channels

type XRPAccountChannels struct {
	ID 				int 			`json:"id"`
	Status		string 		`json:"status"`
	Type      string    `json:"type"`
	
	Result struct {
		Account     	string       `json:"account"`
		Channels    	[]Channel	   `json:"channels"`
		LedgerHash  	string			 `json:"ledger_hash"`
		LedgerIndex 	int64			 	 `json:"ledger_index"`
		Validated   	bool				 `json:"validated"`
		Limit 		 		*int				 `json:"limit, omitempty"`
		Marker 				*string			 `json:"marker, omitempty"`		
	} `json:"result"`
}

type Channel struct {
	Account 						string 		`json:"account"`
	Amount 							string 		`json:"amount"`
	Balance 						string 		`json:"balance"`
	ChannelID   				string  	`json:"channel_id"`
	DestinationAccount 	string 		`json:"destination_account"`
	PublicKey 					string 		`json:"public_key, omitempty"`
	PublicKeyHex 				string 		`json:"public_key_hex, omitempty"`
	SettleDelay 				int 			`json:"settle_delay"`	
}