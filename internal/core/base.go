package core

import (
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/dlbarduzzi/scopehouse/internal/data"
)

// Ensures that the ScopeHouse implements the App interface.
var _ App = (*BaseApp)(nil)

type BaseApp struct {
	logger *slog.Logger
	models *data.Models
}

func NewBaseApp(db *sql.DB, logger *slog.Logger) *BaseApp {
	app := &BaseApp{
		logger: logger,
		models: data.NewModels(db),
	}

	return app
}

// Logger returns the default app logger.
func (app *BaseApp) Logger() *slog.Logger {
	return app.logger
}

// Models returns the default app data.models instance.
func (app *BaseApp) Models() *data.Models {
	return app.models
}

// Bootstrap initializes the application.
func (app *BaseApp) Bootstrap() error {
	if app.logger == nil {
		return errors.New("logger not initialized")
	}

	if app.models == nil {
		return errors.New("models not initialized")
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
