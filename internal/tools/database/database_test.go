package database

import (
	"testing"
	"time"
)

func TestParseConfig(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		Url: "postgres://user:pass@127.0.0.1:5432/db",
	}

	cfg, err := cfg.parse()
	if err != nil {
		t.Fatalf("expected error to be nil, got %v", err)
	}

	maxOpenConns := 10
	maxIdleConns := 5
	connMaxIdleTime := time.Minute * 3

	if cfg.MaxOpenConns != maxOpenConns {
		t.Errorf("expected max-open-conns to be %d, got %d", maxOpenConns, cfg.MaxOpenConns)
	}

	if cfg.MaxIdleConns != maxIdleConns {
		t.Errorf("expected max-idle-conns to be %d, got %d", maxIdleConns, cfg.MaxIdleConns)
	}

	if cfg.ConnMaxIdleTime != connMaxIdleTime {
		t.Errorf(
			"expected conn-max-idle-time to be %v, got %v",
			connMaxIdleTime, cfg.ConnMaxIdleTime,
		)
	}

	cfg.Url = ""

	_, err = cfg.parse()
	wantErr := "invalid database url"

	if err == nil || err.Error() != wantErr {
		t.Fatalf("expected error to be %v, got %v", wantErr, err)
	}
}
