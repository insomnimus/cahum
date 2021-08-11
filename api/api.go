package api

import (
	"github.com/insomnimus/cahum/cah"
)

type Game struct {
	clients []*Client
	// the player that created the game
	owner *Client
	deck  cah.Deck
	turn  uint32

	// broadcasts game events to every client
	events     chan []byte
	register   chan *Client
	unregister chan *Client
	// a message to this channel will start the game if at least 2 players are ready
	start chan bool
}

func NewGame(c *Client) *Game {
	return &Game{
		owner:      c,
		events:     make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		start:      make(chan bool),
		clients:    []*Client{},
	}
}

func (g *Game) Run() {
	// register clients until the start message is received
INIT:
	for {
		select {
		case client := <-g.register:
			g.clients = append(g.clients, client)
		case client := <-g.unregister:
			if client == g.owner {
				close(client.send)
				// terminate the game if there is no one left
				if len(g.clients) == 0 {
					close(g.register)
					close(g.unregister)
					close(g.start)
					close(g.events)
					return
				}

				// move the ownership to the oldest player
				g.owner = g.clients[0]
				g.clients = g.clients[1:]
			} else {
				g.unregisterClient(client)
			}

		case <-g.start:
			if len(g.clients) > 0 {
				close(g.register)
				close(g.start)
				break INIT
			}
		}
	}

	// Start the game.
	g.begin()
}

// Starts the game.
func (g *Game) begin() {
	// TODO: Deal cards and start the game.
}

// Unregisters a non-owner client.
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
