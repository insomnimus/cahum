package api

import (
	"github.com/insomnimus/cahum/api/response"
	"github.com/insomnimus/cahum/cah"
)

type Game struct {
	clients []*Client
	deck    cah.Deck

	// Incoming events such as a player playing a card.
	events chan Event

	// Clients are registered by sending them through this channel.
	/// The channel is closed once the game begins.
	register   chan *Client
	unregister chan *Client
}

func NewGame(c *Client, deck cah.Deck) *Game {
	c.player.ID = 1
	return &Game{
		deck:       deck,
		events:     make(chan Event, 2),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    []*Client{c},
	}
}

// Unregisters a client.
func (g *Game) unregisterClient(client *Client) {
	n := -1
	for i, c := range g.clients {
		if c == client {
			n = i
			break
		}
	}

	if n < 0 {
		return
	}

	close(client.send)
	g.clients = append(g.clients[:n], g.clients[n+1:]...)
}

// getPlayerByID returns the player from the list of clients with the given ID.
//
// Returns nil if no player has the given ID.
func (g *Game) getPlayerByID(id uint32) *cah.Player {
	for _, c := range g.clients {
		if c.player.ID == id {
			return &c.player
		}
	}
	return nil
}

// broadcast sends a message to every client.
func (g *Game) broadcast(data []byte) {
	for _, c := range g.clients {
		c.send <- data
	}
}

// dealCards deals cards to every player, returning true if the deck is out of cards.
func (g *Game) dealCards(n int) (outOfCards bool) {
	// Do we have enough cards?
	if n*len(g.clients) > len(g.deck.White) {
		return true
	}

	for _, c := range g.clients {
		cards := g.deck.DrawWhites(n)
		c.player.Cards = append(c.player.Cards, cards...)
		c.send <- response.DrawCard(cards)
	}

	return
}

// players is a convenience method that returns every `g.client.player` as a slice.
func (g *Game) players() []cah.Player {
	players := make([]cah.Player, 0, len(g.clients))
	for _, c := range g.clients {
		players = append(players, c.player)
	}

	return players
}

// newRound is a convenience method that sends
// a `NewRound` message to every player.
func (g *Game) newRound() {
	card := g.deck.DrawBlack()
	players := g.players()
	data := response.NewRound(card, players)
	g.broadcast(data)
}

// broadcastExcept broadcasts a message to every client except the given one.
func (g *Game) broadcastExcept(client *Client, msg []byte) {
	for _, c := range g.clients {
		if c != client {
			c.send <- msg
		}
	}
}
