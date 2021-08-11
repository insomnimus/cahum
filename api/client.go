package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/insomnimus/cahum/cah"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	game   *Game
	con    *websocket.Conn
	player *cah.Player

	// for sending messages to the player
	send chan []byte
	// true if the client is the owner of the game
	isOwner bool
}

func (c *Client) readPump() {
	defer func() {
		c.game.unregister <- c
		c.con.Close()
	}()

	c.con.SetReadLimit(maxMessageSize)
	c.con.SetReadDeadline(time.Now().Add(pongWait))
	c.con.SetPongHandler(func(string) error {
		c.con.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := c.con.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Generate an event based on the message.
		_ = c.parseMessage(msg)
		// TODO: Do stuff based on the event.
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.con.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.con.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The game closed the channel.
				c.con.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.con.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(msg)

			// write queued messages
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.con.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.con.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWs(game *Game, w http.ResponseWriter, r *http.Request) {
	con, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		game: game,
		con:  con,
		send: make(chan []byte, 256),
	}
	client.game.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

func (c *Client) parseMessage(msg []byte) *Event {
	panic("unimplemented")
}
