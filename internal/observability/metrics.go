package observability

import (
    "expvar"
    "time"
)

var (
    providersMap  = expvar.NewMap("providers")
    blockchainMap = expvar.NewMap("blockchain")
)

func ObserveProviderRequest(d time.Duration, err error) {
    providersMap.Add("requests_total", 1)
    if err != nil {
        providersMap.Add("errors_total", 1)
    }
    // store last latency in ms
    providersMap.Set("last_latency_ms", expvar.Func(func() any { return float64(d.Milliseconds()) }))
}

func ObserveBlockchainTx(d time.Duration, err error) {
    blockchainMap.Add("tx_total", 1)
    if err != nil {
        blockchainMap.Add("tx_errors_total", 1)
    }
    blockchainMap.Set("tx_last_duration_ms", expvar.Func(func() any { return float64(d.Milliseconds()) }))
}


