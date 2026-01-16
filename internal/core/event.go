package core

import "github.com/dlbarduzzi/scopehouse/internal/tools/event"

type EventRequest struct {
	App App
	event.Event
}
