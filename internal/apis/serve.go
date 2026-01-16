package apis

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dlbarduzzi/scopehouse/internal/core"
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

	router := newRouter(app)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		Handler:      router.handler(),
		IdleTimeout:  time.Second * config.IdleTimeout,
		ReadTimeout:  time.Second * config.ReadTimeout,
		WriteTimeout: time.Second * config.WriteTimeout,
	}

	shutdownErr := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		inSignal := <-quit
		app.Logger().Info("server received shutdown signal",
			slog.String("signal", inSignal.String()),
		)

		// Finish running jobs before the server shuts down.
		app.OnShutdown()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			shutdownErr <- err
		}

		shutdownErr <- nil
	}()

	app.Logger().Info("server starting", slog.Int("port", config.Port))

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErr
	if err != nil {
		return err
	}

	app.Logger().Info("server stopped")

	return nil
}
