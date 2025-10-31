package providers

import "testing"

func TestCoinbaseProvider(t *testing.T) {
	provider := CoinbaseProvider{}
	price, err := provider.GetPrice()
	if err != nil {
		t.Fatalf("Erreur GetPrice: %v", err)
	}
	if price <= 0 {
		t.Fatalf("Prix invalide: %f", price)
	}
}
