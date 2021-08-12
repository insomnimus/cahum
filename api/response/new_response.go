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
	"encoding/json"
	"fmt"

	"github.com/insomnimus/cahum/cah"
)

// Error creates a new error object and returns the json encoded bytes.
func Error(format string, args ...interface{}) []byte {
	data, err := json.Marshal(errorMessage{
		Type:  mtError,
		Error: fmt.Sprintf(format, args...),
	})
	if err != nil {
		panic(err)
	}

	return data
}

// DrawCard creates a new `DrawCard` message and returns the json encoded bytes.
func DrawCard(cards []cah.White) []byte {
	data, err := json.Marshal(drawCard{
		Type:  mtDrawCard,
		Cards: cards,
	})
	if err != nil {
		panic(err)
	}
	return data
}

// NewRound creates a new `NewRound` message and returns the json encoded bytes.
//
// The parameter `p []cah.Player` is filtered,
// meaning that the players' hands will not be included in the JSON.
func NewRound(card cah.Black, players []cah.Player) []byte {
	// Erase the cards from every player
	// so the hands are not revealed.
	for i := range players {
		players[i].Cards = nil
	}

	data, err := json.Marshal(newRound{
		Type:    mtNewRound,
		Card:    card,
		Players: players,
	})
	if err != nil {
		panic(err)
	}
	return data
}

// PlayerJoined creates a new `PlayerJoined` message and returns the json encoded bytes.
func PlayerJoined(p cah.Player) []byte {
	data, err := json.Marshal(player{
		Type: mtPlayerJoined,
		ID:   p.ID,
		Name: p.Name,
	})
	if err != nil {
		panic(err)
	}
	return data
}

// PlayerLeft creates a new `PlayerLeft` message and returns the json encoded bytes.
func PlayerLeft(p cah.Player) []byte {
	data, err := json.Marshal(player{
		Type: mtPlayerLeft,
		ID:   p.ID,
		Name: p.Name,
	})
	if err != nil {
		panic(err)
	}
	return data
}

// PlayerPlayedCard creates a new `PlayerPlayedCard` message and returns the json encoded bytes.
func PlayerPlayedCard(p cah.Player) []byte {
	data, err := json.Marshal(player{
		Type: mtPlayerPlayedCard,
		ID:   p.ID,
		Name: p.Name,
	})
	if err != nil {
		panic(err)
	}
	return data
}

// PlayerVoted creates a new `PlayerVoted` message and returns the json encoded bytes.
func PlayerVoted(who, votedFor cah.Player) []byte {
	data, err := json.Marshal(playerVoted{
		Type: mtPlayerVoted,
		Player: player{
			ID:   who.ID,
			Name: who.Name,
		},
		VotedFor: player{
			ID:   votedFor.ID,
			Name: votedFor.Name,
		},
	})
	if err != nil {
		panic(err)
	}

	return data
}

// GameOver creates a `GameOver` message and returns the json encoded bytes.

func GameOver() []byte {
	return []byte(
		fmt.Sprintf(`{"type":%d}`, mtError),
	)
}
