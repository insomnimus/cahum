package api

import (
	"github.com/insomnimus/cahum/api/event"
)

type Event struct {
	source *Client
	event.Event
}
