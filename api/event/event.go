package event

import (
	"github.com/insomnimus/cahum/cah"
)

type EventType string

const (
	PlayCard  EventType = "play"
	StartGame EventType = "start"
	Vote      EventType = "vote"
)

type Event struct {
	Type    EventType  `json:"type"`
	Card    *cah.White `json:"card,omitempty"`
	VoteFor uint32     `json:"vote_for,omitempty"`
}
