package logging

import (
	"context"
	"log/slog"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger()

	if logger == nil {
		t.Fatal("expected logger not to be nil")
	}

	handler := logger.Handler()

	if handler.Enabled(t.Context(), slog.LevelDebug) {
		t.Fatal("expected logger level debug to be disabled")
	}

	if !handler.Enabled(t.Context(), slog.LevelInfo) {
		t.Fatal("expected logger level info to be enabled")
	}
}

func TestNewLoggerWithConfig(t *testing.T) {
	logger := NewLoggerWithConfig(Config{
		Level:  LevelError,
		Format: FormatJson,
	})

	if logger == nil {
		t.Fatal("expected logger not to be nil")
	}

	handler := logger.Handler()

	if handler.Enabled(t.Context(), slog.LevelWarn) {
		t.Fatal("expected logger level warn to be disabled")
	}

	if !handler.Enabled(t.Context(), slog.LevelError) {
		t.Fatal("expected logger level error to be enabled")
	}
}

func TestNewLoggerWithConfigEmpty(t *testing.T) {
	logger := NewLoggerWithConfig(Config{})

	if logger == nil {
		t.Fatal("expected logger not to be nil")
	}

	handler := logger.Handler()

	if handler.Enabled(t.Context(), slog.LevelDebug) {
		t.Fatal("expected logger level debug to be disabled")
	}

	if !handler.Enabled(t.Context(), slog.LevelInfo) {
		t.Fatal("expected logger level info to be enabled")
	}
}

func TestNewLoggerWithConfigDisabled(t *testing.T) {
	logger := NewLoggerWithConfig(Config{Disabled: true})

	if logger == nil {
		t.Fatal("expected logger not to be nil")
	}

	handler := logger.Handler()

	if handler.Enabled(t.Context(), slog.LevelDebug) {
		t.Fatal("expected logger level debug to be disabled")
	}
}

func TestDefaultLogger(t *testing.T) {
	logger1 := DefaultLogger()
	if logger1 == nil {
		t.Fatal("expected logger not to be nil")
	}

	logger2 := DefaultLogger()
	if logger2 == nil {
		t.Fatal("expected logger not to be nil")
	}

	if logger1 != logger2 {
		t.Errorf("expected logger %#v to be equal %#v", logger1, logger2)
	}
}

func TestLoggerContext(t *testing.T) {
	ctx := context.Background()

	logger1 := LoggerFromContext(ctx)
	if logger1 == nil {
		t.Fatal("expected logger not to be nil")
	}

	ctx = LoggerWithContext(ctx, logger1)

	logger2 := LoggerFromContext(ctx)
	if logger1 != logger2 {
		t.Errorf("expected logger %#v to be equal %#v", logger1, logger2)
	}
}

func TestReplaceAttr(t *testing.T) {
	tm := time.Now().UTC()
	sr := &slog.Source{Function: "main.main", File: "/path/to/file", Line: 12}

	testCases := []struct {
		name          string
		useNano       bool
		loggerKey     string
		loggerValue   slog.Value
		expectedKey   string
		expectedValue string
	}{
		{
			name:          "use default timestamp",
			loggerKey:     slog.TimeKey,
			loggerValue:   slog.TimeValue(tm),
			expectedKey:   "time",
			expectedValue: tm.Format("2006-01-02T15:04:05.000Z"),
		},
		{
			name:          "use nano timestamp",
			useNano:       true,
			loggerKey:     slog.TimeKey,
			loggerValue:   slog.TimeValue(tm),
			expectedKey:   "time",
			expectedValue: tm.Format("2006-01-02T15:04:05.000000000Z"),
		},
		{
			name:          "apply level value",
			loggerKey:     slog.LevelKey,
			loggerValue:   slog.AnyValue(slog.LevelInfo),
			expectedKey:   "level",
			expectedValue: "info",
		},
		{
			name:          "apply message key",
			loggerKey:     slog.MessageKey,
			loggerValue:   slog.AnyValue("hello world"),
			expectedKey:   "message",
			expectedValue: "hello world",
		},
		{
			name:          "apply path value",
			loggerKey:     slog.SourceKey,
			loggerValue:   slog.AnyValue(sr),
			expectedKey:   "caller",
			expectedValue: "/path/to/file:12",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fn := ReplaceAttr(tc.useNano)

			attr := slog.Attr{
				Key:   tc.loggerKey,
				Value: tc.loggerValue,
			}

			resp := fn(nil, attr)

			if resp.Key != tc.expectedKey {
				t.Errorf(
					"expected key to be %s; got %s",
					tc.expectedKey, resp.Key,
				)
			}

			if resp.Value.String() != tc.expectedValue {
				t.Errorf(
					"expected value to be %s; got %v",
					tc.expectedValue, resp.Value,
				)
			}
		})
	}
}

func TestGetLogLevel(t *testing.T) {
	testCases := []struct {
		name          string
		level         LogLevel
		expectedLevel slog.Level
	}{
		{
			name:          "empty log level",
			level:         "",
			expectedLevel: slog.LevelInfo,
		},
		{
			name:          "invalid log level",
			level:         "invalid",
			expectedLevel: slog.LevelInfo,
		},
		{
			name:          "debug log level",
			level:         "debug",
			expectedLevel: slog.LevelDebug,
		},
		{
			name:          "info log level",
			level:         "info",
			expectedLevel: slog.LevelInfo,
		},
		{
			name:          "warn log level",
			level:         "warn",
			expectedLevel: slog.LevelWarn,
		},
		{
			name:          "error log level",
			level:         "error",
			expectedLevel: slog.LevelError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			level := SlogLevel(tc.level)
			if level != tc.expectedLevel {
				t.Errorf(
					"expected logger level to be %v; got %v",
					tc.expectedLevel, level,
				)
			}
		})
	}
}
