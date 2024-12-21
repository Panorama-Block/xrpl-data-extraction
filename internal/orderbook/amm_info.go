package orderbook

import (
	"encoding/json"
	"github.com/Panorama-Block/xrpl-data-extraction/internal/xrpl"
)

// FetchAMMInfo fetches AMM information via HTTP
func FetchAMMInfo(client *xrpl.HTTPClient, ammAccount string, asset, asset2 AssetParam) ([]byte, error) {
	params := AMMInfoParam{
		AMMAccount: ammAccount,
		Asset:      asset,
		Asset2:     asset2,
	}

	request := AMMInfoRequest{
		Method: "amm_info",
		Params: []AMMInfoParam{params},
	}

	return client.Post("", request)
}

// StreamAMMInfo streams AMM information via WebSocket
func StreamAMMInfo(wsClient *xrpl.WebSocketClient, ammAccount string, asset, asset2 AssetParam, callback func(*AMMInfoWSResponse)) error {
	request := AMMInfoWSRequest{
		ID:         1,
		Command:    "amm_info",
		AMMAccount: ammAccount,
		Asset:      asset,
		Asset2:     asset2,
	}

	err := wsClient.Subscribe(request)
	if err != nil {
		return err
	}

	wsClient.ReadMessages(func(msg []byte) {
		var response AMMInfoWSResponse
		if err := json.Unmarshal(msg, &response); err == nil {
			callback(&response)
		}
	})
	return nil
}
