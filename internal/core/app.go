package core

import (
	"log/slog"

	"github.com/dlbarduzzi/scopehouse/internal/data"
)

type App interface {
	// Logger returns the default app logger.
	Logger() *slog.Logger

	// Models returns the default app data.models instance.
	Models() *data.Models

	// Bootstrap initializes the application.
	Bootstrap() error

	// OnShutdown run jobs before the application shuts down.
	OnShutdown()
}
