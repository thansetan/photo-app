package helpers

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func NewLogger() *slog.Logger {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		AddSource:  true,
		TimeFormat: "2006-Jan-02 15:04 MST",
	}))

	return logger
}
