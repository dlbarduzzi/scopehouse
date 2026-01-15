package tests

import (
	"database/sql"
	"log/slog"

	"github.com/dlbarduzzi/scopehouse/internal/core"
)

type TestApp struct {
	*core.BaseApp
}

func NewTestApp() (*TestApp, error) {
	db := &sql.DB{}
	logger := slog.New(slog.DiscardHandler)

	app := core.NewBaseApp(db, logger)

	if err := app.Bootstrap(); err != nil {
		return nil, err
	}

	t := &TestApp{
		BaseApp: app,
	}

	return t, nil
}
