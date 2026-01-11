package core

import (
	"database/sql"
	"errors"
	"log/slog"
	"time"
)

// Ensures that the BaseApp implements the App interface.
var _ App = (*BaseApp)(nil)

// BaseApp implements core.App and defines the base ScopeHouse app structure.
type BaseApp struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewBaseApp(db *sql.DB, logger *slog.Logger) *BaseApp {
	app := &BaseApp{
		db:     db,
		logger: logger,
	}

	return app
}

// Logger returns the default app logger.
func (app *BaseApp) Logger() *slog.Logger {
	return app.logger
}

// DB returns the default app data.db builder instance.
func (app *BaseApp) DB() *sql.DB {
	return app.db
}

// Bootstrap initializes the application.
func (app *BaseApp) Bootstrap() error {
	if app.logger == nil {
		return errors.New("logger not initialized")
	}

	if app.db == nil {
		return errors.New("database not initialized")
	}

	return nil
}

// OnShutdown run jobs before the application shuts down.
func (app *BaseApp) OnShutdown() {
	func() {
		app.Logger().Info("...running sample A shutdown func()...")
		time.Sleep(time.Millisecond * 10)
	}()

	func() {
		app.Logger().Info("...running sample B shutdown func()...")
		time.Sleep(time.Millisecond * 10)
	}()
}
