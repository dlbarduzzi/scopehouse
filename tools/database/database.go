package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"

	_ "github.com/lib/pq" // postgres database driver
)

const (
	defaultMaxOpenConns    int           = 10
	defaultMaxIdleConns    int           = 5
	defaultConnMaxIdleTime time.Duration = time.Minute * 3
)

type Config struct {
	Url             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxIdleTime time.Duration
}

func New(config Config) (*sql.DB, error) {
	c, err := config.parse()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", c.Url)
	if err != nil {
		return nil, fmt.Errorf("db open connection failed; %v", err)
	}

	db.SetMaxIdleConns(c.MaxIdleConns)
	db.SetMaxOpenConns(c.MaxOpenConns)
	db.SetConnMaxIdleTime(c.ConnMaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("db ping connection failed; %v", err)
	}

	return db, nil
}

func (c *Config) parse() (*Config, error) {
	_, err := url.ParseRequestURI(c.Url)
	if err != nil {
		// Intentionally suppressing the error to avoid exposing credentials.
		return nil, errors.New("invalid database url")
	}

	if c.MaxIdleConns < 1 {
		c.MaxIdleConns = defaultMaxIdleConns
	}

	if c.MaxOpenConns < 1 {
		c.MaxOpenConns = defaultMaxOpenConns
	}

	if c.ConnMaxIdleTime < time.Minute {
		c.ConnMaxIdleTime = defaultConnMaxIdleTime
	}

	return c, nil
}
