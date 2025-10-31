package main

import (
    "context"
    "go-oracle/internal/config"
    "go-oracle/internal/providers"
    "go-oracle/internal/blockchain"
    "go-oracle/internal/utils"
    "go-oracle/internal/observability"
    "expvar"
    "net/http"
    "log/slog"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    cfg := config.Load()

    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    observability.InitLogger()
    // Expose /debug/vars
    go func() {
        mux := http.NewServeMux()
        mux.Handle("/debug/vars", expvar.Handler())
        _ = http.ListenAndServe(cfg.MetricsAddr, mux)
    }()

    oracle, db, multi, err := initializeComponents(cfg)
    if err != nil {
        slog.Error("init error", "err", err)
        return
    }

    runLoop(ctx, oracle, db, multi, 1*time.Second)
}

func initializeComponents(cfg config.Config) (*blockchain.OracleClient, *utils.JSONDB, *providers.MultiProvider, error) {
	oracle, err := blockchain.NewOracleClient(cfg.RPCURL, cfg.PrivateKeyHex, cfg.ContractAddress, cfg.DryRun)
	if err != nil {
		return nil, nil, nil, err
	}

	db, err := utils.NewJSONDB(cfg.DBPath)
	if err != nil {
		return nil, nil, nil, err
	}

	multi := &providers.MultiProvider{
		Providers: []providers.PriceProvider{
			providers.BinanceProvider{},
			providers.CoinbaseProvider{},
		},
	}

	return oracle, db, multi, nil
}

func runLoop(ctx context.Context, oracle *blockchain.OracleClient, db *utils.JSONDB, multi *providers.MultiProvider, interval time.Duration) {
    ticker := time.NewTicker(interval)
	defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            slog.Info("shutdown requested")
            return
        case <-ticker.C:
            if err := processPriceUpdate(ctx, oracle, db, multi); err != nil {
                slog.Error("price update error", "err", err)
            }
        }
    }
}

func processPriceUpdate(ctx context.Context, oracle *blockchain.OracleClient, db *utils.JSONDB, multi *providers.MultiProvider) error {
    price, err := multi.GetPrice(ctx)
	if err != nil {
		return err
	}

    if err := db.AddPrice(price, len(multi.Providers)); err != nil {
        slog.Error("db add price", "err", err)
    }

    txHash, err := oracle.UpdatePrice(ctx, price)
	if err != nil {
		return err
	}

	if txHash != "" {
        if err := db.AddTx(txHash, price); err != nil {
            slog.Error("db add tx", "err", err)
        }
	}
	return nil
}
