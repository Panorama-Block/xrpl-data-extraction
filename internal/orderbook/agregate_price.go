package orderbook

import (
	"encoding/json"
	"github.com/Panorama-Block/xrpl-data-extraction/internal/xrpl"
)

// FetchAggregatePrice fetches aggregate price data.
func FetchAggregatePrice(client *xrpl.HTTPClient, params GetAggregatePriceParams) ([]byte, error) {
	request := struct {
		Method string               `json:"method"`
		Params []GetAggregatePriceParams `json:"params"`
	}{
		Method: "get_aggregate_price",
		Params: []GetAggregatePriceParams{params},
	}

	return client.Post("", request)
}

// StreamAggregatePrice streams aggregate price data using WebSocket.
func StreamAggregatePrice(wsClient *xrpl.WebSocketClient, params GetAggregatePriceParams, callback func(*GetAggregatePriceResponse)) error {
	request := GetAggregatePriceParams{
		LedgerIndex: params.LedgerIndex,
		BaseAsset:   params.BaseAsset,
		QuoteAsset:  params.QuoteAsset,
		Trim:        params.Trim,
		Oracles:     params.Oracles,
	}

	if err := wsClient.Subscribe(request); err != nil {
		return err
	}

	wsClient.ReadMessages(func(msg []byte) {
		var response GetAggregatePriceResponse
		if err := json.Unmarshal(msg, &response); err == nil {
			callback(&response)
		}
	})
	return nil
}
