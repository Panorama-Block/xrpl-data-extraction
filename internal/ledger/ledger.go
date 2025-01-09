package ledger

import (
	"log"

	"encoding/json"
	"github.com/Panorama-Block/xrpl-data-extraction/internal/xrpl"
)

// FetchLedger fetches ledger information (HTTP)
func FetchLedger(client *xrpl.HTTPClient, ledgerIndex string, transactions, expand, ownerFunds bool) ([]byte, error) {
	params := LedgerParam{
		LedgerIndex:  ledgerIndex,
		Transactions: transactions,
		Expand:       expand,
		OwnerFunds:   ownerFunds,
	}

	request := LedgerRequest{
		Method: "ledger",
		Params: []LedgerParam{params},
	}

	return client.Post("", request)
}

// FetchLedgerClosed fetches the most recently closed ledger
func FetchLedgerClosed(client *xrpl.HTTPClient) ([]byte, error) {
	request := LedgerClosedRequest{
		Method: "ledger_closed",
		Params: []struct{}{{}},
	}
	return client.Post("", request)
}

// FetchLedgerCurrent fetches the current in-progress ledger
func FetchLedgerCurrent(client *xrpl.HTTPClient) ([]byte, error) {
	request := LedgerCurrentRequest{
		Method: "ledger_current",
		Params: []struct{}{{}},
	}
	return client.Post("", request)
}

// FetchLedgerData fetches state data from the specified ledger
func FetchLedgerData(client *xrpl.HTTPClient, ledgerHash string, binary bool, limit int, marker string) ([]byte, error) {
	params := LedgerDataParam{
		LedgerHash: ledgerHash,
		Binary:     binary,
		Limit:      limit,
		Marker:     marker,
	}

	request := LedgerDataRequest{
		Method: "ledger_data",
		Params: []LedgerDataParam{params},
	}

	return client.Post("", request)
}

// StreamLedger envia o comando de subscrição para o WebSocket e processa mensagens continuamente
func StreamLedger(wsClient *xrpl.WebSocketClient, ledgerIndex string, callback func(*LedgerWSResponse), stopChan chan struct{}) error {
	// Comando de subscrição para o WebSocket
	request := map[string]interface{}{
		"id":      1,
		"command": "subscribe",
		"streams": []string{"ledger"}, // Subscrição no stream "ledger"
	}

	// Enviar comando de subscrição
	if err := wsClient.Subscribe(request); err != nil {
		log.Printf("Erro ao enviar comando de subscrição: %v", err)
		return err
	}

	// Processar mensagens continuamente
	go func() {
		wsClient.ReadMessages(func(msg []byte) {
			select {
			case <-stopChan: // Parar quando o canal de parada for fechado
				log.Println("Encerrando streaming de ledger")
				return
			default:
				// Desserializar mensagem para a estrutura LedgerWSResponse
				var response LedgerWSResponse
				if err := json.Unmarshal(msg, &response); err != nil {
					log.Printf("Erro ao desserializar mensagem: %v", err)
					return
				}

				log.Printf("Mensagem processada: %+v", response)
				callback(&response)
			}
		})
	}()

	return nil
}

// StreamLedgerClosed fetches the most recent closed ledger via WebSocket
func StreamLedgerClosed(wsClient *xrpl.WebSocketClient, callback func(*LedgerClosedWSResponse)) error {
	request := LedgerClosedWSRequest{
		ID:      2,
		Command: "ledger_closed",
	}

	err := wsClient.Subscribe(request)
	if err != nil {
		return err
	}

	wsClient.ReadMessages(func(msg []byte) {
		var response LedgerClosedWSResponse
		if err := json.Unmarshal(msg, &response); err == nil {
			callback(&response)
		}
	})
	return nil
}

// StreamLedgerCurrent fetches the current ledger index via WebSocket
func StreamLedgerCurrent(wsClient *xrpl.WebSocketClient, callback func(*LedgerCurrentWSResponse)) error {
	request := LedgerCurrentWSRequest{
		ID:      3,
		Command: "ledger_current",
	}

	err := wsClient.Subscribe(request)
	if err != nil {
		return err
	}

	wsClient.ReadMessages(func(msg []byte) {
		var response LedgerCurrentWSResponse
		if err := json.Unmarshal(msg, &response); err == nil {
			callback(&response)
		}
	})
	return nil
}

// StreamLedgerData fetches ledger data via WebSocket
func StreamLedgerData(wsClient *xrpl.WebSocketClient, ledgerHash string, binary bool, limit int, marker string, callback func(*LedgerDataWSResponse)) error {
	request := LedgerDataWSRequest{
		ID:          4, // Unique ID for the request
		Command:     "ledger_data",
		LedgerHash:  ledgerHash,
		Binary:      binary,
		Limit:       limit,
		Marker:      marker,
	}

	err := wsClient.Subscribe(request)
	if err != nil {
		return err
	}

	wsClient.ReadMessages(func(msg []byte) {
		var response LedgerDataWSResponse
		if err := json.Unmarshal(msg, &response); err == nil {
			callback(&response)
		}
	})
	return nil
}

