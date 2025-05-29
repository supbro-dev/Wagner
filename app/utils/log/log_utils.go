package log

import (
	"log/slog"
	"os"
)

var (
	SystemLogger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
)
