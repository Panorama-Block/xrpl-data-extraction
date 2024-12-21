package orderbook

// ---------- HTTP Request/Response Types ----------

// AssetParam defines an asset in an AMM
type AssetParam struct {
	Currency string `json:"currency"`
	Issuer   string `json:"issuer,omitempty"`
}

// AMMInfoRequest defines the structure for the amm_info HTTP request
type AMMInfoRequest struct {
	Method string         `json:"method"`
	Params []AMMInfoParam `json:"params"`
}

type AMMInfoParam struct {
	AMMAccount string     `json:"amm_account,omitempty"`
	Asset      AssetParam `json:"asset,omitempty"`
	Asset2     AssetParam `json:"asset2,omitempty"`
}

// AMMInfoResponse defines the response for amm_info
type AMMInfoResponse struct {
	Result struct {
		AMM                AMMDetails `json:"amm"`
		LedgerCurrentIndex int        `json:"ledger_current_index"`
		Validated          bool       `json:"validated"`
	} `json:"result"`
}

type AMMDetails struct {
	Account   string      `json:"account"`
	Amount    string      `json:"amount"`
	Amount2   AssetParam  `json:"amount2"`
	TradingFee int        `json:"trading_fee"`
	AuctionSlot AuctionSlot `json:"auction_slot,omitempty"`
}

type AuctionSlot struct {
	Account         string `json:"account"`
	DiscountedFee   int    `json:"discounted_fee"`
	Expiration      string `json:"expiration"`
	TimeInterval    int    `json:"time_interval"`
}

// BookChangesRequest defines the request for book_changes
type BookChangesRequest struct {
	Method string             `json:"method"`
	Params []BookChangesParam `json:"params"`
}

type BookChangesParam struct {
	LedgerIndex int `json:"ledger_index"`
}

// BookChangesResponse defines the response for book_changes
type BookChangesResponse struct {
	Result struct {
		Changes     []BookChange `json:"changes"`
		LedgerIndex int          `json:"ledger_index"`
		LedgerHash  string       `json:"ledger_hash"`
		LedgerTime  int          `json:"ledger_time"`
		Validated   bool         `json:"validated"`
	} `json:"result"`
}

type BookChange struct {
	CurrencyA string `json:"currency_a"`
	CurrencyB string `json:"currency_b"`
	VolumeA   string `json:"volume_a"`
	VolumeB   string `json:"volume_b"`
	High      string `json:"high"`
	Low       string `json:"low"`
	Open      string `json:"open"`
	Close     string `json:"close"`
}

// ---------- WebSocket Request/Response Types ----------

type AMMInfoWSRequest struct {
	ID         int        `json:"id"`
	Command    string     `json:"command"`
	AMMAccount string     `json:"amm_account,omitempty"`
	Asset      AssetParam `json:"asset,omitempty"`
	Asset2     AssetParam `json:"asset2,omitempty"`
}

type AMMInfoWSResponse struct {
	ID     int          `json:"id"`
	Result AMMInfoResponse `json:"result"`
}

type BookChangesWSRequest struct {
	ID          int    `json:"id"`
	Command     string `json:"command"`
	LedgerIndex int    `json:"ledger_index"`
}

type BookChangesWSResponse struct {
	ID     int                  `json:"id"`
	Result BookChangesResponse `json:"result"`
}

// ---------- HTTP Types ----------

// BookOffersParams defines the request parameters for the HTTP `book_offers` method.
type BookOffersParams struct {
	Taker       string     `json:"taker"`
	TakerGets   AssetParam `json:"taker_gets"`
	TakerPays   AssetParam `json:"taker_pays"`
	Limit       int        `json:"limit,omitempty"`
}

// BookOffersResponse represents the HTTP response for `book_offers`.
type BookOffersResponse struct {
	Result struct {
		LedgerCurrentIndex int `json:"ledger_current_index"`
		Offers             []struct {
			Account       string      `json:"Account"`
			BookDirectory string      `json:"BookDirectory"`
			BookNode      string      `json:"BookNode"`
			Flags         int         `json:"Flags"`
			LedgerEntryType string    `json:"LedgerEntryType"`
			OwnerNode     string      `json:"OwnerNode"`
			PreviousTxnID string      `json:"PreviousTxnID"`
			Sequence      int         `json:"Sequence"`
			TakerGets     interface{} `json:"TakerGets"`
			TakerPays     interface{} `json:"TakerPays"`
			Quality       string      `json:"quality"`
		} `json:"offers"`
		Validated bool `json:"validated"`
	} `json:"result"`
}

// GetAggregatePriceParams defines the request parameters for the HTTP `get_aggregate_price` method.
type GetAggregatePriceParams struct {
	LedgerIndex string      `json:"ledger_index"`
	BaseAsset   string      `json:"base_asset"`
	QuoteAsset  string      `json:"quote_asset"`
	Trim        int         `json:"trim,omitempty"`
	Oracles     []Oracle    `json:"oracles"`
}

// Oracle represents the structure for an oracle in `get_aggregate_price`.
type Oracle struct {
	Account           string `json:"account"`
	OracleDocumentID  int    `json:"oracle_document_id"`
}

// GetAggregatePriceResponse represents the HTTP response for `get_aggregate_price`.
type GetAggregatePriceResponse struct {
	Result struct {
		EntireSet struct {
			Mean               string `json:"mean"`
			Size               int    `json:"size"`
			StandardDeviation  string `json:"standard_deviation"`
		} `json:"entire_set"`
		LedgerCurrentIndex int    `json:"ledger_current_index"`
		Median             string `json:"median"`
		Time               int64  `json:"time"`
		TrimmedSet         struct {
			Mean               string `json:"mean"`
			Size               int    `json:"size"`
			StandardDeviation  string `json:"standard_deviation"`
		} `json:"trimmed_set"`
		Validated bool `json:"validated"`
	} `json:"result"`
}

// NFTBuyOffersParams defines the request parameters for the HTTP `nft_buy_offers` method.
type NFTBuyOffersParams struct {
	NFTID       string `json:"nft_id"`
	LedgerIndex string `json:"ledger_index"`
}

// NFTBuyOffersResponse represents the HTTP response for `nft_buy_offers`.
type NFTBuyOffersResponse struct {
	Result struct {
		NFTID  string `json:"nft_id"`
		Offers []struct {
			Amount          string `json:"amount"`
			Flags           int    `json:"flags"`
			NFTOfferIndex   string `json:"nft_offer_index"`
			Owner           string `json:"owner"`
		} `json:"offers"`
	} `json:"result"`
}

// NFTSellOffersParams defines the request parameters for the HTTP `nft_sell_offers` method.
type NFTSellOffersParams struct {
	NFTID       string `json:"nft_id"`
	LedgerIndex string `json:"ledger_index"`
}

// NFTSellOffersResponse represents the HTTP response for `nft_sell_offers`.
type NFTSellOffersResponse struct {
	Result struct {
		NFTID  string `json:"nft_id"`
		Offers []struct {
			Amount          string `json:"amount"`
			Flags           int    `json:"flags"`
			NFTOfferIndex   string `json:"nft_offer_index"`
			Owner           string `json:"owner"`
		} `json:"offers"`
	} `json:"result"`
}

// ---------- WebSocket Types ----------

// BookOffersWSRequest represents the WebSocket request for `book_offers`.
type BookOffersWSRequest struct {
	ID        int        `json:"id"`
	Command   string     `json:"command"`
	Taker     string     `json:"taker"`
	TakerGets AssetParam `json:"taker_gets"`
	TakerPays AssetParam `json:"taker_pays"`
	Limit     int        `json:"limit,omitempty"`
}

// BookOffersWSResponse represents the WebSocket response for `book_offers`.
type BookOffersWSResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Result struct {
		LedgerCurrentIndex int `json:"ledger_current_index"`
		Offers             []struct {
			Account       string      `json:"Account"`
			BookDirectory string      `json:"BookDirectory"`
			BookNode      string      `json:"BookNode"`
			Flags         int         `json:"Flags"`
			LedgerEntryType string    `json:"LedgerEntryType"`
			OwnerNode     string      `json:"OwnerNode"`
			PreviousTxnID string      `json:"PreviousTxnID"`
			Sequence      int         `json:"Sequence"`
			TakerGets     interface{} `json:"TakerGets"`
			TakerPays     interface{} `json:"TakerPays"`
			Quality       string      `json:"quality"`
		} `json:"offers"`
	} `json:"result"`
}

// NFTBuyOffersWSRequest represents the WebSocket request for `nft_buy_offers`.
type NFTBuyOffersWSRequest struct {
	ID        int    `json:"id"`
	Command   string `json:"command"`
	NFTID     string `json:"nft_id"`
	LedgerIndex string `json:"ledger_index,omitempty"`
}

// NFTBuyOffersWSResponse represents the WebSocket response for `nft_buy_offers`.
type NFTBuyOffersWSResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Result struct {
		NFTID  string `json:"nft_id"`
		Offers []struct {
			Amount          string `json:"amount"`
			Flags           int    `json:"flags"`
			NFTOfferIndex   string `json:"nft_offer_index"`
			Owner           string `json:"owner"`
		} `json:"offers"`
	} `json:"result"`
}

// NFTSellOffersWSRequest represents the WebSocket request for `nft_sell_offers`.
type NFTSellOffersWSRequest struct {
	ID        int    `json:"id"`
	Command   string `json:"command"`
	NFTID     string `json:"nft_id"`
	LedgerIndex string `json:"ledger_index,omitempty"`
}

// NFTSellOffersWSResponse represents the WebSocket response for `nft_sell_offers`.
type NFTSellOffersWSResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Result struct {
		NFTID  string `json:"nft_id"`
		Offers []struct {
			Amount          string `json:"amount"`
			Flags           int    `json:"flags"`
			NFTOfferIndex   string `json:"nft_offer_index"`
			Owner           string `json:"owner"`
		} `json:"offers"`
	} `json:"result"`
}


