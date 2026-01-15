package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/viper"

	"github.com/dlbarduzzi/scopehouse/internal/apis"
	"github.com/dlbarduzzi/scopehouse/internal/core"
	"github.com/dlbarduzzi/scopehouse/internal/tools/database"
	"github.com/dlbarduzzi/scopehouse/internal/tools/logging"
)

func main() {
	if err := start(); err != nil {
		fmt.Fprintf(os.Stderr, "[error] %s\n", err)
		os.Exit(1)
	}
}

func start() error {
	config := getConfig()

	logger := logging.NewLoggerWithConfig(config.logger)
	logger = logger.With(slog.String("app", "scopehouse"))

	db, err := database.New(config.db)
	if err != nil {
		return err
	}

	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("db close failed", slog.Any("error", err))
		}
	}()

	app := core.NewBaseApp(db, logger)

	if err := app.Bootstrap(); err != nil {
		return err
	}

	return apis.Serve(app, config.server)
}

type config struct {
	db     database.Config
	logger logging.Config
	server apis.ServerConfig
}

var (
	// default database configs
	defaultDatabaseMaxIdleConns    = 10
	defaultDatabaseMaxOpenConns    = 10
	defaultDatabaseConnMaxIdleTime = time.Minute * 3

	// default logger configs
	defaultLoggerLevel  = "info"
	defaultLoggerFormat = "json"

	// default server configs
	defaultServerPort         = 8090
	defaultServerIdleTimeout  = time.Second * 10
	defaultServerReadTimeout  = time.Second * 5
	defaultServerWriteTimeout = time.Second * 5
)

func getConfig() config {
	v := viper.New()
	v.AutomaticEnv()

	c := config{
		db: database.Config{
			Url:             v.GetString("SH_DATABASE_URL"),
			MaxIdleConns:    v.GetInt("SH_DATABASE_MAX_IDLE_CONNS"),
			MaxOpenConns:    v.GetInt("SH_DATABASE_MAX_OPEN_CONNS"),
			ConnMaxIdleTime: v.GetDuration("SH_DATABASE_CONN_MAX_IDLE_TIME"),
		},
		logger: logging.Config{
			Level:     v.GetString("SH_LOG_LEVEL"),
			Format:    v.GetString("SH_LOG_FORMAT"),
			UseNano:   v.GetBool("SH_LOG_USE_NANO"),
			UseSource: v.GetBool("SH_LOG_USE_SOURCE"),
		},
		server: apis.ServerConfig{
			Port:         v.GetInt("SH_SERVER_PORT"),
			IdleTimeout:  v.GetDuration("SH_SERVER_IDLE_TIMEOUT"),
			ReadTimeout:  v.GetDuration("SH_SERVER_READ_TIMEOUT"),
			WriteTimeout: v.GetDuration("SH_SERVER_WRITE_TIMEOUT"),
		},
	}

	// Database configs.
	if c.db.MaxIdleConns < 1 {
		c.db.MaxIdleConns = defaultDatabaseMaxIdleConns
	}

	if c.db.MaxOpenConns < 1 {
		c.db.MaxOpenConns = defaultDatabaseMaxOpenConns
	}

	if c.db.ConnMaxIdleTime < time.Minute {
		c.db.ConnMaxIdleTime = defaultDatabaseConnMaxIdleTime
	}

	// Logger configs.
	if c.logger.Level == "" {
		c.logger.Level = defaultLoggerLevel
	}

	if c.logger.Format == "" {
		c.logger.Format = defaultLoggerFormat
	}

	// Server configs.
	if c.server.Port < 1 {
		c.server.Port = defaultServerPort
	}

	if c.server.IdleTimeout < time.Second*10 {
		c.server.IdleTimeout = defaultServerIdleTimeout
	}

	if c.server.ReadTimeout < time.Second {
		c.server.ReadTimeout = defaultServerReadTimeout
	}

	if c.server.WriteTimeout < time.Second {
		c.server.WriteTimeout = defaultServerWriteTimeout
	}

	return c
}
