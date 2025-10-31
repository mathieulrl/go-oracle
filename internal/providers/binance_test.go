package providers

import (
	"testing"
)

// TestBinanceProvider vérifie que BinanceProvider retourne un prix positif
func TestBinanceProvider(t *testing.T) {
	provider := BinanceProvider{}
	price, err := provider.GetPrice()
	if err != nil {
		t.Fatalf("Erreur GetPrice: %v", err)
	}
	if price <= 0 {
		t.Fatalf("Prix invalide: %f", price)
	}
}
