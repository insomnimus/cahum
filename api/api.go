package api

import (
	"github.com/insomnimus/cahum/api/event"
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
	// A message to this channel will start the game if at least 2 players are ready
	// and the source of the message is the lobby owner (clients[0]).
	start chan bool
}

func NewGame(c *Client) *Game {
	c.player.ID = 1
	return &Game{
		events:     make(chan Event, 2),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		start:      make(chan bool),
		clients:    []*Client{c},
	}
}

// `run` starts the lobby, registering players until the start
// message is received.
func (g *Game) Run() {
	idCounter := uint32(1) // There's one player already.
	// register clients until the start message is received
INIT:
	for {
		select {
		case c := <-g.register:
			idCounter++
			c.player.ID = idCounter
			g.clients = append(g.clients, c)
		case client := <-g.unregister:
			g.unregisterClient(client)
			if len(g.clients) == 0 {
				return
			}
		case <-g.start:
			if len(g.clients) > 1 {
				close(g.register)
				close(g.start)
				break INIT
			}
		// TODO: Send an error message to g.clients[0].send.
		case _ = <-g.events:
			// TODO: Send an error message to _.source.send; the game is not ready.
		}
	}

	// Start the game.
	g.begin()
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

func (g *Game) getClientByPlayerID(id uint32) *Client {
	for _, c := range g.clients {
		if c.player.ID == id {
			return c
		}
	}
	return nil
}

// Deals cards to every player, returning true if the deck is out of cards.
func (g *Game) dealCards() (outOfCards bool) {
	panic("not yet implemented")
}

// Starts the game.
func (g *Game) begin() {
	// TODO: Deal cards and broadcast a ready message.

	var turn int
	// Map of players and the card they played in a turn.
	// The values will be reset to nil each turn.
	played := make(map[*Client]*cah.White)
	// Map of players and who they voted for.
	// The values will be reset to nil each turn.
	voted := make(map[*Client]*Client)

	for _, c := range g.clients {
		played[c] = nil
		voted[c] = nil
	}

	// Game loop.
	for {
		for c := range played {
			played[c] = nil
			voted[c] = nil
		}

		// Deal cards every 10 turns.
		if turn > 0 && turn%10 == 0 {
			if outOfCards := g.dealCards(); outOfCards {
				// TODO: Send game over to everyone.
				return
			}
		}

		turn++
		i := 0
		// Take input from every player.
		for i < len(played) {
			select {
			case c := <-g.unregister:
				g.unregisterClient(c)
				// Return if not enough players are left.
				if len(g.clients) < 2 {
					// TODO: Inform the remaining player.
					return
				}
				// Decrement `i` if the player had played a card so the loop still works.
				if card := played[c]; card != nil {
					i--
				}
				delete(played, c)
				delete(voted, c)

			case e := <-g.events:
				if e.Type != event.PlayCard {
					// TODO: Send an error message to e.source.send.
					continue
				}
				// Do not accept if the player already played this turn.
				if card := played[e.source]; card != nil {
					// TODO: Send an error message to e.source.send.
					continue
				}
				i++
				played[e.source] = e.Card
				e.source.player.RemoveCard(e.Card.ID)
			}
		}

		// Every player played a card.
		// TODO: Broadcast event to start voting.

		i = 0
		for i < len(voted) {
			select {
			case c := <-g.unregister:
				g.unregisterClient(c)
				if len(g.clients) < 2 {
					// TODO: Inform the remaining client.
					return
				}
				if val := played[c]; val != nil {
					i--
				}
				delete(played, c)
				delete(voted, c)

			case e := <-g.events:
				if e.Type != event.Vote {
					// TODO: Send an error message to e.source.send.
					continue
				}
				if val := voted[e.source]; val != nil {
					// TODO: Send an error message to e.source.send.
					continue
				}
				c := g.getClientByPlayerID(e.VoteFor)
				if c == nil {
					// TODO: Send an error message to e.source.send.
					continue
				}
				i++
				voted[e.source] = c
			}

			// Everyone voted, calculate scores.
			// TODO: Do the above.
			// TODO: Broadcast a new black card.
		}
	}
}
