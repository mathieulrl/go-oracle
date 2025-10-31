package providers

import "context"

// PriceProvider lit un prix depuis une source externe.
type PriceProvider interface {
    GetPrice(ctx context.Context) (float64, error)
}

// Structure correspondant à la réponse de l’API Binance
type BinanceResponse struct {
    Symbol string `json:"symbol"`
    Price  string `json:"price"`
}

type CoinbaseResponse struct {
    Data struct {
        Amount string `json:"amount"`
    } `json:"data"`
}


