package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"
)

// contextKey is the logger string type used to avoid context collisions.
type contextKey string

// loggerKey identifies the logger value stored in the context.
const loggerKey = contextKey("logger")

var (
	defaultLogger     *slog.Logger
	defaultLoggerOnce sync.Once
)

type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

type LogFormat string

const (
	FormatText LogFormat = "text"
	FormatJson LogFormat = "json"
)

const (
	DefaultLevel  LogLevel  = LevelInfo
	DefaultFormat LogFormat = FormatText
)

type Config struct {
	Level     LogLevel
	Format    LogFormat
	UseNano   bool
	AddSource bool
	Disabled  bool
}

func NewLogger() *slog.Logger {
	return NewLoggerWithConfig(Config{
		Level:     DefaultLevel,
		Format:    DefaultFormat,
		UseNano:   false,
		AddSource: true,
	})
}

func NewLoggerWithConfig(config Config) *slog.Logger {
	if config.Level == "" {
		config.Level = DefaultLevel
	}

	if config.Format == "" {
		config.Format = DefaultFormat
	}

	options := &slog.HandlerOptions{
		Level:       SlogLevel(config.Level),
		AddSource:   config.AddSource,
		ReplaceAttr: ReplaceAttr(config.UseNano),
	}

	if config.Disabled {
		return slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, options))

	if config.Format == FormatJson {
		logger = slog.New(slog.NewJSONHandler(os.Stderr, options))
	}

	return logger
}

func DefaultLogger() *slog.Logger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = NewLogger()
	})
	return defaultLogger
}

func LoggerWithContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return DefaultLogger()
}

type slogAttr func(groups []string, attr slog.Attr) slog.Attr

func ReplaceAttr(useNano bool) slogAttr {
	return func(_ []string, attr slog.Attr) slog.Attr {
		if attr.Key == slog.TimeKey {
			attr.Key = "time"
			attr.Value = slog.StringValue(
				TimeFormat(attr.Value.Time().UTC(), useNano),
			)
		}
		if attr.Key == slog.LevelKey {
			if level, ok := attr.Value.Any().(slog.Level); ok {
				attr.Value = slog.StringValue(strings.ToLower(level.String()))
			}
		}
		if attr.Key == slog.MessageKey {
			attr.Key = "message"
		}
		if attr.Key == slog.SourceKey {
			source := attr.Value.Any().(*slog.Source)
			attr.Key = "caller"
			attr.Value = slog.StringValue(fmt.Sprintf("%s:%d", source.File, source.Line))
		}
		return attr
	}
}

func SlogLevel(level LogLevel) slog.Level {
	switch level {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func TimeFormat(t time.Time, useNano bool) string {
	if useNano {
		return t.Format("2006-01-02T15:04:05.000000000Z")
	}
	return t.Format("2006-01-02T15:04:05.000Z")
}
