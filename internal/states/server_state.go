package serverinfo

import (
	"encoding/json"
	"github.com/Panorama-Block/xrpl-data-extraction/internal/xrpl"
)

// FetchServerState fetches the server state via HTTP
func FetchServerState(client *xrpl.HTTPClient) ([]byte, error) {
	request := struct {
		Method string `json:"method"`
		Params []struct {
			LedgerIndex string `json:"ledger_index"`
		} `json:"params"`
	}{
		Method: "server_state",
		Params: []struct {
			LedgerIndex string `json:"ledger_index"`
		}{{LedgerIndex: "current"}},
	}

	return client.Post("", request)
}

// StreamServerState streams the server state via WebSocket
func StreamServerState(wsClient *xrpl.WebSocketClient, callback func(*ServerStateWSResponse)) error {
	request := struct {
		ID          int    `json:"id"`
		Command     string `json:"command"`
		LedgerIndex string `json:"ledger_index"`
	}{
		ID:          2,
		Command:     "server_state",
		LedgerIndex: "current",
	}

	err := wsClient.Subscribe(request)
	if err != nil {
		return err
	}

	wsClient.ReadMessages(func(msg []byte) {
		var response ServerStateWSResponse
		if err := json.Unmarshal(msg, &response); err == nil {
			callback(&response)
		}
	})
	return nil
}

// ServerStateWSResponse represents the WebSocket response for the server_state command
type ServerStateWSResponse struct {
	ID     int              `json:"id"`
	Status string           `json:"status"`
	Type   string           `json:"type"`
	Result ServerStateResult `json:"result"`
}

// ServerStateResult defines the structure for the server_state response
type ServerStateResult struct {
	State struct {
		BuildVersion      string `json:"build_version"`
		CompleteLedgers   string `json:"complete_ledgers"`
		IO_LatencyMS      int    `json:"io_latency_ms"`
		JQTransOverflow   string `json:"jq_trans_overflow"`
		LastClose         struct {
			ConvergeTime int `json:"converge_time"`
			Proposers    int `json:"proposers"`
		} `json:"last_close"`
		LoadBase          int    `json:"load_base"`
		LoadFactor        int    `json:"load_factor"`
		LoadFactorFeeEsc  int    `json:"load_factor_fee_escalation"`
		LoadFactorFeeQueue int   `json:"load_factor_fee_queue"`
		LoadFactorFeeRef  int    `json:"load_factor_fee_reference"`
		LoadFactorServer  int    `json:"load_factor_server"`
		Peers             int    `json:"peers"`
		ServerState       string `json:"server_state"`
		ServerStateDurUS  string `json:"server_state_duration_us"`
		StateAccounting   map[string]struct {
			DurationUS   string `json:"duration_us"`
			Transitions  string `json:"transitions"`
		} `json:"state_accounting"`
		Time             string `json:"time"`
		Uptime           int    `json:"uptime"`
		ValidatedLedger  struct {
			BaseFee      int    `json:"base_fee"`
			CloseTime    int    `json:"close_time"`
			Hash         string `json:"hash"`
			ReserveBase  int    `json:"reserve_base"`
			ReserveInc   int    `json:"reserve_inc"`
			Seq          int    `json:"seq"`
		} `json:"validated_ledger"`
		ValidationQuorum int `json:"validation_quorum"`
	} `json:"state"`
}
