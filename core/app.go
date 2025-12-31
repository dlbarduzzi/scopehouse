package core

import "log/slog"

type App interface {
	// Logger returns the default app logger.
	Logger() *slog.Logger

	// Bootstrap initializes the application.
	Bootstrap() error

	// OnShutdown run jobs before the application shuts down.
	OnShutdown()
}
