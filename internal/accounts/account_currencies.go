package accounts

import (
	"encoding/json"

	"github.com/Panorama-Block/xrpl-data-extraction/internal/xrpl"
)

func FetchAccountCurrencies(client *xrpl.HTTPClient, account string, ledgerIndex string) ([]byte, error) {
	params := AccountCurrenciesParam{
		Account: account,
	}
	if ledgerIndex != "" {
		params.LedgerIndex = ledgerIndex
	}

	payload := AccountCurrenciesRequest{
		Method: "account_currencies",
		Params: []AccountCurrenciesParam{params},
	}

	return client.Post("", payload)
}

func StreamAccountCurrencies(client *xrpl.WebSocketClient, account string, ledgerIndex string, callback func(*AccountCurrenciesWSResponse)) error {
	params := AccountCurrenciesParam{
		Account: account,
	}
	if ledgerIndex != "" {
		params.LedgerIndex = ledgerIndex
	}

	payload := AccountCurrenciesWSRequest{
		ID:          1,
		Command:     "account_currencies",
		Account:     account,
		LedgerIndex: ledgerIndex,
	}

	err := client.Subscribe(payload)
	if err != nil {
		return err
	}

	client.ReadMessages(func(message []byte) {
		var response AccountCurrenciesWSResponse
		err := json.Unmarshal(message, &response)
		if err != nil {
			return
		}
		callback(&response)
	})

	return nil
}
