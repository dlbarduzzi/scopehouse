package apis

import (
	"log/slog"
	"time"

	"github.com/dlbarduzzi/scopehouse/core"
)

const (
	DefaultPort         = 8090
	DefaultIdleTimeout  = time.Second * 30
	DefaultReadTimeout  = time.Second * 5
	DefaultWriteTimeout = time.Second * 5
)

type ServerConfig struct {
	Port         int
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Serve(app core.App, config ServerConfig) error {
	if config.Port < 1 {
		config.Port = DefaultPort
	}

	if config.IdleTimeout < time.Second*10 {
		config.IdleTimeout = DefaultIdleTimeout
	}

	if config.ReadTimeout < time.Second*1 {
		config.ReadTimeout = DefaultReadTimeout
	}

	if config.WriteTimeout < time.Second*1 {
		config.WriteTimeout = DefaultWriteTimeout
	}

	app.Logger().Info("server starting", slog.Int("port", config.Port))

	return nil
}
