package tests

import "github.com/dlbarduzzi/scopehouse/core"

type TestApp struct {
	*core.BaseApp
}

func NewTestApp() (*TestApp, error) {
	return NewTestAppWithConfig(core.BaseAppConfig{
		LogDisabled: true,
	})
}

func NewTestAppWithConfig(config core.BaseAppConfig) (*TestApp, error) {
	app := core.NewBaseApp(config)

	if err := app.Bootstrap(); err != nil {
		return nil, err
	}

	t := &TestApp{
		BaseApp: app,
	}

	return t, nil
}
