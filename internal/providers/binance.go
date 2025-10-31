package providers

import (
	"fmt"
    "context"
    "time"
    "go-oracle/internal/observability"
)

// BinanceProvider implémente PriceProvider pour Binance
type BinanceProvider struct{}

// GetPrice récupère le prix du BTC/USDT sur Binance
func (b BinanceProvider) GetPrice(ctx context.Context) (float64, error) {
    start := time.Now()
    var retErr error
    defer func() { observability.ObserveProviderRequest(time.Since(start), retErr) }()
	url := "https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT"
    var data BinanceResponse
    if err := fetchApiClient(ctx, url, &data); err != nil {
        retErr = fmt.Errorf("erreur appel/JSON Binance: %w", err)
        return 0, retErr
    }
    price, err := parseFloat(data.Price)
    if err != nil {
        retErr = err
        return 0, retErr
    }
    return price, nil
}
