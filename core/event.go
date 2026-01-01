package core

import "github.com/dlbarduzzi/scopehouse/tools/event"

type EventRequest struct {
	App App
	event.Event
}
