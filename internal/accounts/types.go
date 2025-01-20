package accounts

// ---------- HTTP Types ----------

type Channel struct {
	Account           string `json:"account"`
	Amount            string `json:"amount"`
	Balance           string `json:"balance"`
	ChannelID         string `json:"channel_id"`
	DestinationAccount string `json:"destination_account"`
	SettleDelay       int    `json:"settle_delay"`
}

// HTTP: Account Channels
type AccountChannelsRequest struct {
	Method string                 `json:"method"` // The JSON-RPC command, e.g., "account_channels".
	Params []AccountChannelsParam `json:"params"`
}

type AccountChannelsParam struct {
	Account            string `json:"account"`                        // The account to query.
	DestinationAccount string `json:"destination_account,omitempty"` // Optional: Filter by destination account.
	LedgerIndex        string `json:"ledger_index,omitempty"`        // Optional: Specify a ledger index or shortcut.
}

type AccountChannelsResponse struct {
	ID     int    `json:"id"`     // Request identifier.
	Status string `json:"status"` // Status of the response, e.g., "success".
	Type   string `json:"type"`   // Type of the response, e.g., "response".
	Result struct {
		Account     string    `json:"account"`     // Account queried.
		Channels    []Channel `json:"channels"`    // List of payment channels.
		LedgerIndex int64     `json:"ledger_index"` // Ledger index used for the query.
		Validated   bool      `json:"validated"`   // Indicates if the data is from a validated ledger.
	} `json:"result"`
}

// HTTP: Account Currencies
type AccountCurrenciesRequest struct {
	Method string                  `json:"method"` // The JSON-RPC command, e.g., "account_currencies".
	Params []AccountCurrenciesParam `json:"params"`
}

type AccountCurrenciesParam struct {
	Account     string `json:"account"`                // The account to query for currencies.
	LedgerIndex string `json:"ledger_index,omitempty"` // Optional: Specify a ledger index or shortcut.
}

type AccountCurrenciesResponse struct {
	ID     int    `json:"id"`     // Request identifier.
	Status string `json:"status"` // Status of the response, e.g., "success".
	Type   string `json:"type"`   // Type of the response, e.g., "response".
	Result struct {
		LedgerIndex      int      `json:"ledger_index"`       // Ledger index used for the query.
		ReceiveCurrencies []string `json:"receive_currencies"` // Currencies the account can receive.
		SendCurrencies    []string `json:"send_currencies"`    // Currencies the account can send.
		Validated         bool     `json:"validated"`          // Indicates if the data is from a validated ledger.
	} `json:"result"`
}

// HTTP: Account Info
type AccountInfoRequest struct {
	Method string             `json:"method"` // The JSON-RPC command, e.g., "account_info".
	Params []AccountInfoParam `json:"params"`
}

type AccountInfoParam struct {
	Account     string `json:"account"`                // The account to retrieve information for.
	LedgerIndex string `json:"ledger_index,omitempty"` // Optional: Specify a ledger index or shortcut.
	Queue       bool   `json:"queue,omitempty"`        // Optional: Include information about queued transactions.
}

type AccountInfoResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Result struct {
		AccountData struct {
			Account           string `json:"Account"`           // Account address.
			Balance           string `json:"Balance"`           // Current XRP balance in drops.
			OwnerCount        int    `json:"OwnerCount"`        // Number of objects owned by the account.
			Sequence          int    `json:"Sequence"`          // Current transaction sequence number.
			PreviousTxnID     string `json:"PreviousTxnID"`     // Last transaction affecting the account.
			PreviousTxnLgrSeq int    `json:"PreviousTxnLgrSeq"` // Ledger sequence of the last transaction affecting the account.
		} `json:"account_data"`
		LedgerIndex int  `json:"ledger_index"` // Ledger index used for the query.
		Validated   bool `json:"validated"`   // Indicates if the data is from a validated ledger.
	} `json:"result"`
}

// HTTP: Account Lines
type AccountLinesRequest struct {
	Method string              `json:"method"` // The JSON-RPC command, e.g., "account_lines".
	Params []AccountLinesParam `json:"params"`
}

type AccountLinesParam struct {
	Account     string `json:"account"`                // The account to query for trust lines.
	LedgerIndex string `json:"ledger_index,omitempty"` // Optional: Specify a ledger index or shortcut.
}

type AccountLinesResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Result struct {
		Account string `json:"account"` // Account queried.
		Lines   []struct {
			Account    string `json:"account"`    // Counterparty account.
			Balance    string `json:"balance"`    // Current balance on the trust line.
			Currency   string `json:"currency"`   // Currency code.
			Limit      string `json:"limit"`      // Max amount the account is willing to owe.
			LimitPeer  string `json:"limit_peer"` // Max amount the peer is willing to owe.
			QualityIn  int    `json:"quality_in"` // Quality in value for this trust line.
			QualityOut int    `json:"quality_out"` // Quality out value for this trust line.
		} `json:"lines"`
		LedgerIndex int  `json:"ledger_index"` // Ledger index used for the query.
		Validated   bool `json:"validated"`   // Indicates if the data is from a validated ledger.
	} `json:"result"`
}

// ---------- WebSocket Types ----------

// WebSocket: Account Channels
type AccountChannelsWSRequest struct {
	ID                int    `json:"id"`                // Request identifier.
	Command           string `json:"command"`           // The JSON-RPC command, e.g., "account_channels".
	Account           string `json:"account"`           // The account to query.
	DestinationAccount string `json:"destination_account,omitempty"` // (Optional) Destination account filter.
	LedgerIndex       string `json:"ledger_index,omitempty"`         // (Optional) Specify a ledger index or shortcut.
}

type AccountChannelsWSResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Result struct {
		Account  string    `json:"account"`
		Channels []Channel `json:"channels"` // Reuse HTTP Channel type
	} `json:"result"`
}

// WebSocket: Account Currencies
type AccountCurrenciesWSRequest struct {
	ID           int    `json:"id"`            // Request identifier.
	Command      string `json:"command"`       // The JSON-RPC command, e.g., "account_currencies".
	Account      string `json:"account"`       // The account to query.
	LedgerIndex  string `json:"ledger_index"`  // Specify a ledger index or shortcut.
}

type AccountCurrenciesWSResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Result struct {
		LedgerIndex      int      `json:"ledger_index"`
		ReceiveCurrencies []string `json:"receive_currencies"`
		SendCurrencies    []string `json:"send_currencies"`
		Validated         bool     `json:"validated"`
	} `json:"result"`
}

// WebSocket: Account Info
type AccountInfoWSRequest struct {
	ID           int    `json:"id"`            // Request identifier.
	Command      string `json:"command"`       // The JSON-RPC command, e.g., "account_info".
	Account      string `json:"account"`       // The account to retrieve information for.
	LedgerIndex  string `json:"ledger_index"`  // Specify a ledger index or shortcut.
	Queue        bool   `json:"queue"`         // Include information about queued transactions.
}

type AccountInfoWSResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Result struct {
		AccountData struct {
			Account           string `json:"Account"`
			Balance           string `json:"Balance"`
			Flags             int64  `json:"Flags"`
			LedgerEntryType   string `json:"LedgerEntryType"`
			OwnerCount        int    `json:"OwnerCount"`
			PreviousTxnID     string `json:"PreviousTxnID"`
			PreviousTxnLgrSeq int    `json:"PreviousTxnLgrSeq"`
			Sequence          int    `json:"Sequence"`
			Index             string `json:"index"`
		} `json:"account_data"`
		LedgerCurrentIndex int  `json:"ledger_current_index"`
		QueueData          struct {
			AuthChangeQueued bool   `json:"auth_change_queued"`
			HighestSequence  int    `json:"highest_sequence"`
			LowestSequence   int    `json:"lowest_sequence"`
			MaxSpendDrops    string `json:"max_spend_drops_total"`
			Transactions     []struct {
				AuthChange       bool   `json:"auth_change"`
				Fee              string `json:"fee"`
				FeeLevel         string `json:"fee_level"`
				MaxSpendDrops    string `json:"max_spend_drops"`
				Seq              int    `json:"seq"`
				LastLedgerSeq    int    `json:"LastLedgerSequence,omitempty"`
			} `json:"transactions"`
			TxnCount int `json:"txn_count"`
		} `json:"queue_data"`
		Validated bool `json:"validated"`
	} `json:"result"`
}

// WebSocket: Account Lines
type AccountLinesWSRequest struct {
	ID           int    `json:"id"`            // Request identifier.
	Command      string `json:"command"`       // The JSON-RPC command, e.g., "account_lines".
	Account      string `json:"account"`       // The account to query for trust lines.
	LedgerIndex  string `json:"ledger_index"`  // Specify a ledger index or shortcut.
}

type AccountLinesWSResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Result struct {
		Account string `json:"account"`
		Lines   []struct {
			Account       string `json:"account"`
			Balance       string `json:"balance"`
			Currency      string `json:"currency"`
			Limit         string `json:"limit"`
			LimitPeer     string `json:"limit_peer"`
			NoRipple      bool   `json:"no_ripple,omitempty"`
			NoRipplePeer  bool   `json:"no_ripple_peer,omitempty"`
			QualityIn     int    `json:"quality_in"`
			QualityOut    int    `json:"quality_out"`
		} `json:"lines"`
	} `json:"result"`
}

// HTTP Request and Response for Account_NFTs
type AccountNFTsRequest struct {
	Method string             `json:"method"`
	Params []AccountNFTsParam `json:"params"`
}

type AccountNFTsParam struct {
	Account     string `json:"account"`
	LedgerIndex string `json:"ledger_index,omitempty"`
	Limit       int    `json:"limit,omitempty"`
	Marker      string `json:"marker,omitempty"`
}

type AccountNFTsResponse struct {
	ID      int    `json:"id"`
	Status  string `json:"status"`
	Type    string `json:"type"`
	Result  struct {
		Account     string   `json:"account"`
		AccountNFTs []NFT    `json:"account_nfts"`
		LedgerHash  string   `json:"ledger_hash,omitempty"`
		LedgerIndex int      `json:"ledger_index,omitempty"`
		Validated   bool     `json:"validated"`
		Marker      *string  `json:"marker,omitempty"`
	} `json:"result"`
}

type NFT struct {
	Flags       int    `json:"Flags"`
	Issuer      string `json:"Issuer"`
	NFTokenID   string `json:"NFTokenID"`
	NFTokenTaxon int    `json:"NFTokenTaxon"`
	URI         string `json:"URI"`
	NFTSerial   int    `json:"nft_serial"`
}

// WebSocket Request and Response for Account_NFTs
type AccountNFTsWSRequest struct {
	Command string `json:"command"`
	Account string `json:"account"`
	LedgerIndex string `json:"ledger_index,omitempty"`
	Limit int `json:"limit,omitempty"`
	Marker string `json:"marker,omitempty"`
}


// HTTP Request and Response for Gateway Balances
type GatewayBalancesRequest struct {
	Method string                  `json:"method"`
	Params []GatewayBalancesParam  `json:"params"`
}

type GatewayBalancesParam struct {
	Account     string   `json:"account"`
	HotWallet   []string `json:"hotwallet,omitempty"`
	LedgerIndex string   `json:"ledger_index,omitempty"`
	Strict      bool     `json:"strict,omitempty"`
}

type GatewayBalancesResponse struct {
	ID      int    `json:"id"`
	Status  string `json:"status"`
	Type    string `json:"type"`
	Result  struct {
		Account      string            `json:"account"`
		Obligations  map[string]string `json:"obligations,omitempty"`
		Balances     map[string][]struct {
			Currency string `json:"currency"`
			Value    string `json:"value"`
		} `json:"balances,omitempty"`
		Assets map[string][]struct {
			Currency string `json:"currency"`
			Value    string `json:"value"`
		} `json:"assets,omitempty"`
		LedgerHash  string `json:"ledger_hash,omitempty"`
		LedgerIndex int    `json:"ledger_index,omitempty"`
		Validated   bool   `json:"validated"`
	} `json:"result"`
}

// WebSocket Request and Response for Gateway Balances
type GatewayBalancesWSRequest struct {
	Command    string   `json:"command"`
	Account    string   `json:"account"`
	HotWallet  []string `json:"hotwallet,omitempty"`
	LedgerIndex string  `json:"ledger_index,omitempty"`
	Strict      bool    `json:"strict,omitempty"`
}

type GatewayBalancesWSResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Result struct {
		Account     string            `json:"account"`
		Obligations map[string]string `json:"obligations,omitempty"`
		Balances    map[string][]struct {
			Currency string `json:"currency"`
			Value    string `json:"value"`
		} `json:"balances,omitempty"`
		Assets map[string][]struct {
			Currency string `json:"currency"`
			Value    string `json:"value"`
		} `json:"assets,omitempty"`
		LedgerHash  string `json:"ledger_hash,omitempty"`
		LedgerIndex int    `json:"ledger_index,omitempty"`
		Validated   bool   `json:"validated"`
	} `json:"result"`
} 

// ACCOUNT SUBSCRIBE TYPES

type SubscribeAccountsRequest struct {
	ID       string   `json:"id"`
	Command  string   `json:"command"`
	Accounts []string `json:"accounts"`
}

// AccountTransactionMessage representa a mensagem de transação recebida do WebSocket
type AccountTransactionMessage struct {
	Type        string `json:"type"`
	Account     string `json:"Account"`
	LedgerHash  string `json:"ledger_hash"`
	LedgerIndex int    `json:"ledger_index"`
	Tx          struct {
		TransactionType string      `json:"TransactionType"`
		Account         string      `json:"Account"`
		Fee             string      `json:"Fee"`
		TakerGets       interface{} `json:"TakerGets"`
		TakerPays       interface{} `json:"TakerPays"`
		Date            int64       `json:"date"`
		OwnerFunds      string      `json:"owner_funds"`
	} `json:"transaction"`
	Validated bool   `json:"validated"`
	Status    string `json:"status"`
}