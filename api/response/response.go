/*
The response package defines functions and types
for generating messages that will be sent from
the server to a client.

Since the communication is one way, the only exported
items of this package are the functions.

Every function in this package returns
a JSON encoded representation of a message as []byte.

Since the input to this package is trusted (comes from the server itself)
any errors that may arise from marshaling JSON are
considered a bug and will cause a panic.
*/

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

type playedInfo struct {
	player
	Card cah.White `json:"card"`
}

type startVoting struct {
	Played []playedInfo `json:"played"`
}
