package providers

import (
    "fmt"
    "context"
    "time"
    "go-oracle/internal/observability"
)

type CoinbaseProvider struct{}

func (c CoinbaseProvider) GetPrice(ctx context.Context) (float64, error) {
    start := time.Now()
    var retErr error
    defer func() { observability.ObserveProviderRequest(time.Since(start), retErr) }()
	url := "https://api.coinbase.com/v2/prices/BTC-USD/spot"
    var data CoinbaseResponse
    if err := fetchApiClient(ctx, url, &data); err != nil {
        retErr = fmt.Errorf("erreur appel/JSON Coinbase: %w", err)
        return 0, retErr
    }
    price, err := parseFloat(data.Data.Amount)
    if err != nil {
        retErr = err
        return 0, retErr
    }
    return price, nil
}
