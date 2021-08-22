package api

import (
	"github.com/insomnimus/cahum/api/event"
	"github.com/insomnimus/cahum/api/response"
	"github.com/insomnimus/cahum/cah"
)

// run starts the lobby, registering players until the start
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
			// Inform other players.
			g.broadcast(response.PlayerJoined(c.player))
			g.clients = append(g.clients, c)
		case c := <-g.unregister:
			g.unregisterClient(c)
			if len(g.clients) == 0 {
				return
			}
			g.broadcast(response.PlayerLeft(c.player))
		case e := <-g.events:
			if e.Type != event.StartGame {
				e.source.send <- response.Error("the game has not yet started")
				continue
			}
			if e.source != g.clients[0] {
				e.source.send <- response.Error("only the lobby owner can start the game")
				continue
			}
			if len(g.clients) < 2 {
				e.source.send <- response.Error("not enough players")
				continue
			}
			break INIT
		}
	}

	close(g.register)

	// Start the game.
	g.begin()
}

// Starts the game.
func (g *Game) begin() {
	g.dealCards(10)
	g.newRound()

	var turn int
	// Map of players and the card they played in a turn.
	// The values will be reset to nil each turn.
	played := make(map[*cah.Player]*cah.White)
	// Map of players and who they voted for.
	// The values will be reset to nil each turn.
	voted := make(map[*cah.Player]*cah.Player)

	for _, c := range g.clients {
		played[&c.player] = nil
		voted[&c.player] = nil
	}

	// Game loop.
	for len(g.deck.Black) > 0 {
		for p := range played {
			played[p] = nil
			voted[p] = nil
		}

		// Deal cards every 5 turns.
		if turn > 0 && turn%5 == 0 {
			if outOfCards := g.dealCards(5); outOfCards {
				g.broadcast(response.GameOver())
				return
			}
		}

		turn++
		g.newRound()
		i := 0
		// Take input from every player.
		for i < len(played) {
			select {
			case c := <-g.unregister:
				g.unregisterClient(c)
				g.broadcast(response.PlayerLeft(c.player))

				// Return if not enough players are left.
				if len(g.clients) < 2 {
					g.broadcast(response.GameOver())
					return
				}
				// Decrement `i` if the player had played a card so the loop still works.
				if card := played[&c.player]; card != nil {
					i--
				}
				delete(played, &c.player)
				delete(voted, &c.player)

			case e := <-g.events:
				if e.Type != event.PlayCard {
					e.source.send <- response.Error("invalid request type, expected play-card, got %s", e.Type)
					continue
				}
				// Do not accept if the player already played this turn.
				if card := played[&e.source.player]; card != nil {
					e.source.send <- response.Error("already played a card this turn")
					continue
				}
				// Check if the player has the card, remove it.
				if !e.source.player.RemoveCard(e.Card.ID) {
					// TODO: Send the players hand to the client?
					e.source.send <- response.Error("you don't have that card")
					continue
				}
				// Inform other players that this player played a card.
				g.broadcastExcept(e.source, response.PlayerPlayedCard(e.source.player))
				i++
				played[&e.source.player] = e.Card
			}
		}

		// Every player played a card.
		// Reveal cards and inform clients to start voting.
		g.broadcast(response.StartVoting(played))

		i = 0
		for i < len(voted) {
			select {
			case c := <-g.unregister:
				g.unregisterClient(c)
				g.broadcast(response.PlayerLeft(c.player))

				if len(g.clients) < 2 {
					g.broadcast(response.GameOver())
					return
				}

				if val := voted[&c.player]; val != nil {
					i--
				}
				delete(played, &c.player)
				delete(voted, &c.player)

			case e := <-g.events:
				if e.Type != event.Vote {
					e.source.send <- response.Error("invalid request: expected vote; got %s", e.Type)
					continue
				}
				if val := voted[&e.source.player]; val != nil {
					e.source.send <- response.Error("you already played this turn")
					continue
				}

				p := g.getPlayerByID(e.VoteFor)
				if p == nil {
					// TODO: Send a list of players?
					e.source.send <- response.Error("player does not exist")
					continue
				}

				i++
				// Inform other players that a vote happened.
				g.broadcastExcept(e.source, response.PlayerVoted(e.source.player, *p))
				voted[&e.source.player] = p
			}
		}

		// Everyone voted; calculate scores.
		g.updateScores(voted)
	}

	// Game over.
	g.broadcast(response.GameOver())
}
