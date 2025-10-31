package observability

import (
    "log/slog"
    "os"
)

var defaultLogger *slog.Logger

func InitLogger() {
    if defaultLogger != nil {
        return
    }
    // Text handler for readability; switch to JSON if needed
    defaultLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

func Logger() *slog.Logger {
    if defaultLogger == nil {
        InitLogger()
    }
    return defaultLogger
}


