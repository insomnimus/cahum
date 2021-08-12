package event

import (
	"github.com/insomnimus/cahum/cah"
)

type EventType uint8

const (
	_ EventType = iota
	PlayCard
	StartGame
	Vote
)

type Event struct {
	Type    EventType  `json:"type"`
	Card    *cah.White `json:"card,omitempty"`
	VoteFor uint32     `json:"vote_for,omitempty"`
}
