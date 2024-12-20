package orderbook

import (
	"github.com/Panorama-Block/xrpl-data-extraction/internal/xrpl"
)

// FetchNFTBuyOffers fetches buy offers for a specific NFT.
func FetchNFTBuyOffers(client *xrpl.HTTPClient, params NFTBuyOffersParams) ([]byte, error) {
	request := struct {
		Method string           `json:"method"`
		Params []NFTBuyOffersParams `json:"params"`
	}{
		Method: "nft_buy_offers",
		Params: []NFTBuyOffersParams{params},
	}

	return client.Post("", request)
}

// FetchNFTSellOffers fetches sell offers for a specific NFT.
func FetchNFTSellOffers(client *xrpl.HTTPClient, params NFTSellOffersParams) ([]byte, error) {
	request := struct {
		Method string            `json:"method"`
		Params []NFTSellOffersParams `json:"params"`
	}{
		Method: "nft_sell_offers",
		Params: []NFTSellOffersParams{params},
	}

	return client.Post("", request)
}
