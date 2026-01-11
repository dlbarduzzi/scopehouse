package logging

import (
	"context"
	"fmt"
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

const (
	levelDebug string = "debug"
	levelInfo  string = "info"
	levelWarn  string = "warn"
	levelError string = "error"
)

const (
	formatText string = "text"
	formatJson string = "json"
)

const (
	defaultLevel  string = levelInfo
	defaultFormat string = formatText
)

type Config struct {
	Level     string
	Format    string
	UseNano   bool
	UseSource bool
}

func NewLogger() *slog.Logger {
	return NewLoggerWithConfig(Config{
		Level:     defaultLevel,
		Format:    defaultFormat,
		UseNano:   false,
		UseSource: true,
	})
}

func NewLoggerWithConfig(config Config) *slog.Logger {
	if config.Level == "" {
		config.Level = defaultLevel
	}

	if config.Format == "" {
		config.Format = defaultFormat
	}

	options := &slog.HandlerOptions{
		Level:       toSlogLevel(config.Level),
		AddSource:   config.UseSource,
		ReplaceAttr: ReplaceAttr(config.UseNano),
	}

	if config.Format == formatJson {
		return slog.New(slog.NewJSONHandler(os.Stderr, options))
	}

	return slog.New(slog.NewTextHandler(os.Stderr, options))
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
				formatTime(attr.Value.Time().UTC(), useNano),
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

func toSlogLevel(level string) slog.Level {
	switch strings.TrimSpace(strings.ToLower(level)) {
	case levelDebug:
		return slog.LevelDebug
	case levelInfo:
		return slog.LevelInfo
	case levelWarn:
		return slog.LevelWarn
	case levelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func formatTime(t time.Time, useNano bool) string {
	if useNano {
		return t.Format("2006-01-02T15:04:05.000000000Z")
	}
	return t.Format("2006-01-02T15:04:05.000Z")
}
