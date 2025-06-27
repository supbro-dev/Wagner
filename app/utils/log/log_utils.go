package log

import (
	"log/slog"
	"os"
	"wagner/app/global/business_error"
)

var (
	SystemLogger        = slog.New(slog.NewTextHandler(os.Stdout, nil))
	BusinessErrorLogger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	ComputeLogger       = slog.New(slog.NewTextHandler(os.Stdout, nil))
	InfoLogger          = slog.New(slog.NewTextHandler(os.Stdout, nil))
)

func LogBusinessError(businessError *business_error.BusinessError) {
	BusinessErrorLogger.Error(businessError.Type, businessError.Message)
}
