package response

import (
	"github.com/insomnimus/cahum/cah"
)

type messageType uint8

//go:generate stringer -type=messageType -trimprefix=mt

const (
	_ messageType = iota
	mtError
	mtDrawCard
	mtPlayerJoined
	mtPlayerLeft
	mtPlayerPlayedCard
	mtPlayerVoted
	mtNewRound
	mtGameOver
)

type errorMessage struct {
	Type  messageType `json:"type"`
	Error string      `json:"error"`
}

type drawCard struct {
	Type  messageType `json:"type"`
	Cards []cah.White `json:"cards"`
}

type newRound struct {
	Type    messageType  `json:"type"`
	Card    cah.Black    `json:"card"`
	Players []cah.Player `json:"players"`
}

type player struct {
	// This field is for convenience.
	Type messageType `json:"type,omitempty"`
	Name string      `json:"name"`
	ID   uint32      `json:"id"`
}

type playerVoted struct {
	Type     messageType `json:"type"`
	Player   player      `json:"player"`
	VotedFor player      `json:"voted_for"`
}
