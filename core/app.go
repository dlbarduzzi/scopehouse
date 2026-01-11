package core

import (
	"database/sql"
	"log/slog"
)

type App interface {
	// Logger returns the default app logger.
	Logger() *slog.Logger

	// DB returns the default app data.db builder instance.
	DB() *sql.DB

	// Bootstrap initializes the application.
	Bootstrap() error

	// OnShutdown run jobs before the application shuts down.
	OnShutdown()
}
