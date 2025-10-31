package providers

import (
    "context"
	"testing"
)

type MockProvider struct {
	Price float64
	Err   error
}

func (m MockProvider) GetPrice(ctx context.Context) (float64, error) {
	return m.Price, m.Err
}

func TestMultiProviderWithMock(t *testing.T) {
	multi := MultiProvider{
		Providers: []PriceProvider{
			MockProvider{Price: 10000, Err: nil},
			MockProvider{Price: 20000, Err: nil},
		},
	}

    price, err := multi.GetPrice(context.Background())
	if err != nil {
		t.Fatalf("Erreur: %v", err)
	}

	if price != 15000 {
		t.Fatalf("Prix moyen attendu 15000, got %f", price)
	}
}
