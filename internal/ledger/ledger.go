package ledger

import (
	"log"
	"context"
	"time"
	"encoding/json"
	"fmt"
	
	"github.com/Panorama-Block/xrpl-data-extraction/internal/database"
	"github.com/Panorama-Block/xrpl-data-extraction/internal/xrpl"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

)

// FetchLedgerInfo sends a properly structured HTTP request to fetch ledger information
func FetchLedgerInfo(client *xrpl.HTTPClient, ledgerIndex string) (*LedgerResponse, error) {
	params := LedgerParam{
		LedgerIndex:  ledgerIndex,
		Transactions: true,
		Expand:       true,
		OwnerFunds:   false,
	}

	request := LedgerRequest{
		Method: "ledger",
		Params: []LedgerParam{params},
	}

	// Send the request using the HTTP client with the direct URL
	responseData, err := client.Post("https://s1.ripple.com:51234", request)
	if err != nil {
		log.Printf("❌ Failed to fetch ledger info: %v", err)
		return nil, err
	}

	// Parse the response into the defined struct
	var ledgerResponse LedgerResponse
	err = json.Unmarshal(responseData, &ledgerResponse)
	if err != nil {
		log.Printf("❌ Error unmarshalling ledger response: %v", err)
		return nil, err
	}

	log.Printf("✅ Successfully fetched ledger info: LedgerIndex: %s", ledgerResponse.Result.Ledger.LedgerIndex)
	return &ledgerResponse, nil
}

func StreamLedger(wsClient *xrpl.WebSocketClient, httpClient *xrpl.HTTPClient, callback func(*LedgerSubscribeClosedResponse), stopChan chan struct{}) error {
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
					var closedResponse LedgerSubscribeClosedResponse

					// Interpretar como mensagem de ledger fechado
					if err := json.Unmarshal(msg, &closedResponse); err == nil && closedResponse.Type == "ledgerClosed" {
						log.Printf("✅ Novo ledger recebido: %+v", closedResponse)

						// Chamar FetchLedgerInfo para obter totalCoins
						ledgerIndex := fmt.Sprintf("%d", closedResponse.LedgerIndex)
						ledgerInfo, err := FetchLedgerInfo(httpClient, ledgerIndex)
						if err != nil {
							log.Printf("❌ Erro ao buscar informações adicionais do ledger: %v", err)
						} else {
							// Extraindo o campo totalCoins
							totalCoins := ledgerInfo.Result.Ledger.TotalCoins
							log.Printf("✅ TotalCoins extraído: %s", totalCoins)

							// Adicionar totalCoins aos dados WebSocket
							closedResponse.TotalCoins = totalCoins
						}

						// Salvar no banco de dados
						if err := SaveLedgerToDB(&closedResponse); err != nil {
							log.Printf("❌ Erro ao salvar no banco de dados: %v", err)
						}
						
						// Invocar o callback com os dados atualizados
						callback(&closedResponse)
					} else {
						log.Printf("⚠️ Mensagem desconhecida recebida: %s", string(msg))
					}
				})
			}
		}
	}()
	return nil
}



func SaveLedgerToDB(data *LedgerSubscribeClosedResponse) error {
	collection := database.GetLedgerCollection()

	// Dados incompletos ou inválidos
	if data.LedgerIndex == 0 || data.LedgerHash == "" || data.TxnCount <= 0 || data.TotalCoins == "" {
		log.Println("⚠️ Dados incompletos ou inválidos. Ignorando salvamento.")
		return nil
	}

	ledgerData := LedgerSchema{
		LedgerIndex:    data.LedgerIndex,
		LedgerHash:     data.LedgerHash,
		CloseTimeHuman: time.Unix(data.LedgerTime, 0).Format(time.RFC3339),
		TxnCount:       data.TxnCount,
		FeeBase:        data.FeeBase,
		TotalCoins:     data.TotalCoins,
		CreatedAt:      time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Verificar se o ledger já existe no banco de dados
	filter := bson.M{"ledger_index": data.LedgerIndex}
	existing := collection.FindOne(ctx, filter)
	if existing.Err() == nil {
		log.Printf("⚠️ Ledger com ledger_index %d já existe no banco de dados. Ignorando.", data.LedgerIndex)
		return nil
	} else if existing.Err() != mongo.ErrNoDocuments {
		log.Printf("❌ Erro ao verificar existência do ledger no banco de dados: %v", existing.Err())
		return existing.Err()
	}

	// Inserir o ledger no banco de dados
	_, err := collection.InsertOne(ctx, ledgerData)
	if err != nil {
		log.Printf("❌ Erro ao salvar no banco de dados: %v", err)
		return err
	}

	log.Printf("✅ Ledger salvo no banco de dados: %+v", ledgerData)
	return nil
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

