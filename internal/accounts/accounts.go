package accounts

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Panorama-Block/xrpl-data-extraction/internal/database"
	"github.com/Panorama-Block/xrpl-data-extraction/internal/xrpl"
)

func SubscribeAccounts(wsClient *xrpl.WebSocketClient, accounts []string, stopChan chan struct{}) error {
	request := SubscribeAccountsRequest{
		ID:       "subscribe_accounts",
		Command:  "subscribe",
		Accounts: accounts,
	}

	log.Printf("🔗 Enviando comando de inscrição para contas: %+v", accounts)

	if err := wsClient.Subscribe(request); err != nil {
		log.Printf("❌ Erro ao enviar o comando subscribe: %v", err)
		return err
	}

	go func() {
		for {
			select {
			case <-stopChan:
				log.Println("⛔ Encerrando a inscrição para contas.")
				wsClient.Connection.Close()
				return
			default:
				wsClient.ReadMessages(func(msg []byte) {
					log.Printf("📩 Mensagem recebida do WebSocket: %s", string(msg))

					var accountMessage AccountTransactionMessage
					if err := json.Unmarshal(msg, &accountMessage); err != nil {
						log.Printf("⚠️ Erro ao interpretar mensagem de conta: %v", err)
						return
					}

					// Salvar transação no banco de dados
					err := SaveTransactionToDB(&accountMessage)
					if err != nil {
						log.Printf("❌ Erro ao salvar transação no MongoDB: %v", err)
					} else {
						log.Printf("✅ Transação salva no banco de dados: %+v", accountMessage)
					}
				})
			}
		}
	}()
	return nil
}


// SaveTransactionToDB salva transações específicas no banco de dados
func SaveTransactionToDB(data *AccountTransactionMessage) error {
	collection := database.GetTransactionCollection()

	// Tratamento do campo TakerGets (interface{})
	takerGets, err := json.Marshal(data.Tx.TakerGets)
	if err != nil {
		log.Printf("❌ Erro ao processar TakerGets: %v", err)
		return err
	}

	// Tratamento do campo TakerPays (interface{})
	takerPays, err := json.Marshal(data.Tx.TakerPays)
	if err != nil {
		log.Printf("❌ Erro ao processar TakerPays: %v", err)
		return err
	}

	// Converter timestamp para time.Time
	transactionDate := time.Unix(data.Tx.Date, 0)

	// Criar objeto para salvar no banco
	transactionData := TransactionSchema{
		Account:      data.Tx.Account,
		Fee:          data.Tx.Fee,
		TakerGets:    string(takerGets),
		TakerPays:    string(takerPays),
		Date:         transactionDate,
		OwnerFunds:   data.Tx.OwnerFunds,
		Type:         data.Type,
		Validated:    data.Validated,
		Status:       data.Status,
		LedgerHash:   data.LedgerHash,
		LedgerIndex:  data.LedgerIndex,
		CreatedAt:    time.Now(),
	}

	// Verificar se os dados são válidos antes de salvar
	if transactionData.Account == "" || transactionData.Type == "" {
		log.Printf("⚠️ Dados incompletos ou inválidos: %+v", transactionData)
		return nil
	}

	log.Printf("💾 Salvando transação no MongoDB: %+v", transactionData)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, transactionData)
	if err != nil {
		log.Printf("❌ Erro ao salvar transação no MongoDB: %v", err)
		return err
	}

	log.Printf("✅ Transação salva no MongoDB com sucesso: %+v", transactionData)
	return nil
}
