package providers

import (
    "encoding/json"
    "fmt"
    "context"
    "net/http"
    "time"
)

const defaultHTTPTimeout = 5 * time.Second

var defaultHTTPClient = &http.Client{Timeout: defaultHTTPTimeout}

func fetchApiClient(ctx context.Context, url string, out any) error {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return err
    }
    resp, err := defaultHTTPClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("code HTTP inattendu: %d", resp.StatusCode)
    }
    if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
        return err
    }
    return nil
}

func parseFloat(s string) (float64, error) {
    var v float64
    if _, err := fmt.Sscan(s, &v); err != nil {
        return 0, fmt.Errorf("conversion prix échouée: %w", err)
    }
    return v, nil
}


