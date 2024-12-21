package orderbook

import (
	"encoding/json"
	"github.com/Panorama-Block/xrpl-data-extraction/internal/xrpl"
)

// FetchBookOffers sends a JSON-RPC request to fetch book offers.
func FetchBookOffers(client *xrpl.HTTPClient, params BookOffersParams) ([]byte, error) {
	request := struct {
		Method string           `json:"method"`
		Params []BookOffersParams `json:"params"`
	}{
		Method: "book_offers",
		Params: []BookOffersParams{params},
	}

	return client.Post("", request)
}

// StreamBookOffers streams book offers using WebSocket.
func StreamBookOffers(wsClient *xrpl.WebSocketClient, params BookOffersParams, callback func(*BookOffersWSResponse)) error {
	request := BookOffersWSRequest{
		ID:        4,
		Command:   "book_offers",
		Taker:     params.Taker,
		TakerGets: params.TakerGets,
		TakerPays: params.TakerPays,
		Limit:     params.Limit,
	}

	if err := wsClient.Subscribe(request); err != nil {
		return err
	}

	wsClient.ReadMessages(func(msg []byte) {
		var response BookOffersWSResponse
		if err := json.Unmarshal(msg, &response); err == nil {
			callback(&response)
		}
	})
	return nil
}
