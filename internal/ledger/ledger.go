package ledger

import (
	"log"
	"context"
	"time"
	"encoding/json"

	"github.com/Panorama-Block/xrpl-data-extraction/internal/database"
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


// StreamLedger usando o comando subscribe corretamente e printando os dados completos
func StreamLedger(wsClient *xrpl.WebSocketClient, callback func(*LedgerSubscribeClosedResponse), stopChan chan struct{}) error {
	request := map[string]interface{}{
		"id":      "1",
		"command": "subscribe",
		"streams": []string{"ledger"},
	}

	if err := wsClient.Subscribe(request); err != nil {
		log.Printf("❌ Erro ao enviar o comando de subscribe: %v", err)
		return err
	}

	go func() {
		for {
			select {
			case <-stopChan:
				log.Println("⛔ Encerrando o streaming de ledgers.")
				wsClient.Connection.Close()
				return
			default:
				wsClient.ReadMessages(func(msg []byte) {
					var initialResponse LedgerSubscribeResponse
					var closedResponse LedgerSubscribeClosedResponse

					// Tentar interpretar como resposta inicial
					if err := json.Unmarshal(msg, &initialResponse); err == nil && initialResponse.Type == "response" {
						log.Printf("✅ Conexão estabelecida. Ledger Atual: %+v", initialResponse.Result)
						return
					}

					// Tentar interpretar como mensagem de ledger fechado
					if err := json.Unmarshal(msg, &closedResponse); err == nil && closedResponse.Type == "ledgerClosed" {
						log.Printf("✅ Novo ledger fechado recebido: %+v", closedResponse)
						callback(&closedResponse)
						return
					}

					// Caso não caia em nenhum dos dois
					log.Printf("⚠️ Mensagem desconhecida recebida: %s", string(msg))
				})
			}
		}
	}()
	return nil
}


func SaveLedgerToDB(data *LedgerSubscribeClosedResponse) error {
	collection := database.GetLedgerCollection()

	if data.LedgerIndex == 0 || data.LedgerHash == "" {
		log.Println("⚠️ Dados incompletos. Ignorando salvamento.")
		return nil
	}

	ledgerData := LedgerSchema{
		LedgerIndex:     data.LedgerIndex,
		LedgerHash:      data.LedgerHash,
		// CloseTimeHuman:  time.Unix(data.LedgerTime, 0).Format(time.RFC3339),
		TxnCount:        data.TxnCount,
		FeeBase:         data.FeeBase,
		CreatedAt:       time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, ledgerData)
	if err != nil {
		log.Printf("❌ Erro ao salvar no banco de dados: %v", err)
		return err
	}

	log.Printf("✅ Ledger salvo no banco de dados: %+v", ledgerData)
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

