package core

import (
	"log/slog"
	"time"
)

const (
	DefaultLogLevel  = "info"
	DefaultLogFormat = "json"
)

// BaseAppConfig defines a BaseApp configuration option.
type BaseAppConfig struct {
	LogLevel    string
	LogFormat   string
	LogDisabled bool
}

// Ensures that the BaseApp implements the App interface.
var _ App = (*BaseApp)(nil)

// BaseApp implements core.App and defines the base ScopeHouse app structure.
type BaseApp struct {
	logger *slog.Logger
	config *BaseAppConfig
}

func NewBaseApp(config BaseAppConfig) *BaseApp {
	app := &BaseApp{
		config: &config,
	}

	if app.config.LogLevel == "" {
		app.config.LogLevel = DefaultLogLevel
	}

	if app.config.LogFormat == "" {
		app.config.LogFormat = DefaultLogFormat
	}

	return app
}

// Logger returns the default app logger.
func (app *BaseApp) Logger() *slog.Logger {
	if app.logger == nil {
		return slog.Default()
	}
	return app.logger
}

// Bootstrap initializes the application.
func (app *BaseApp) Bootstrap() error {
	if err := app.initLogger(); err != nil {
		return err
	}

	return nil
}

// OnShutdown run jobs before the application shuts down.
func (app *BaseApp) OnShutdown() {
	func() {
		app.Logger().Info("...running first shutdown func()...")
		time.Sleep(time.Millisecond * 10)
	}()
	func() {
		app.Logger().Info("...running second shutdown func()...")
		time.Sleep(time.Millisecond * 10)
	}()
}

func (app *BaseApp) initLogger() error {
	return nil
}
