package apis

import "github.com/dlbarduzzi/scopehouse/internal/core"

type service struct {
	app core.App
}

func newService(app core.App) *service {
	return &service{app}
}
