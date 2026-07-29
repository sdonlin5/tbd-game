// Implements client
// -
package server

import (
	"encoding/json"
	"game_project/game"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Constat values for ping / pong between peer and server
const turnTime = 60 * time.Second
const pongWait = time.Second * 60
const pingPeriod = (pongWait * 9) / 10

// Struct to receive input JSON from the websocket connection.
type Envelope struct {
	Type    string          // type of input the user entered
	Payload json.RawMessage // data of the input, e.g. shot coords
}

// Intermediary between websocket connection and hub.
type Client struct {
	id   uuid.UUID // unique identifier for each client
	hub  *Hub
	conn *websocket.Conn

	// Write only permission channel for Match to read
	// Writes the Move received from websocket connection to be
	// processed by Match
	ActionSender chan<- *game.Action

	// Read/Write Bidirectional channel.
	// Writes response from Match to be read by outputPump
	// Outbound to WS: write permission - Response -> conn
	ResponseReceiver chan *game.Response
}

// Satisfies EventNotifier interface
func (c *Client) Notify(resp *game.Response) {
	c.ResponseReceiver <- resp
}

// inputPump pumps input received from the websocket connection to the hub.
//
// Called on a pointer to a Client (c) the server runs inputPump on a per-connection goroutine. Server
// ensures that there is only one input on a connection by executing all reads from this goroutine.
func (c *Client) inputPump() {

	// Defers functions until inputPump finishes executing
	defer func() {
		c.conn.Close()
	}()

	// Sets the read deadline before connection times out
	c.conn.SetReadDeadline(time.Now().Add(pongWait))

	// Pong handler configuration for the connection
	c.conn.SetPongHandler(func(string) error {
		// After each pong is received, reset the the deadline
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		var input Envelope
		// Blocks until data arrives
		readErr := c.conn.ReadJSON(&input)

		// Checks for network error
		if readErr != nil {
			log.Printf("ERROR: %v", readErr)
			return
		}

		// Gaurds against the client not having a match
		// If true, skips switch until next websocket frame arrives
		if c.ActionSender == nil {
			log.Printf("Error: %v sent payload before joining match!", c.id)
			continue
		}

		// Intermediate storage
		var move game.Move
		var jsonError error

		// Parses the payload to concrete type
		switch input.Type {
		case "Shot":
			var s game.Shot
			jsonError = json.Unmarshal(input.Payload, &s)
			move = s

		case "Quit":
			var q game.Quit
			jsonError = json.Unmarshal(input.Payload, &q)
			move = q

		// Add other cases as needed
		default:
			log.Printf("Unknown input type: %v", input.Type)
			continue // prevent sending nil to match
		}

		// Handle malformed input
		if jsonError != nil {
			log.Printf("Error %v", jsonError)
			continue

			// Send the input to the match channel
		} else {
			c.ActionSender <- &game.Action{SenderID: c.id, Move: move}
			c.ActionSender <- &game.Action{Move: move}
		}
	}
}

func (c *Client) outputPump() {
	// Pumps messages from the hub to the websocket connection.

	// Sets the heartbeat interval
	// Sends control frame on each tick
	ticker := time.NewTicker(pingPeriod)

	// Defer until after outputPump finishes executing
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	// Output loop
	for {
		select {
		// Response { }
		case resp, ok := <-c.ResponseReceiver:
			c.conn.SetWriteDeadline(time.Now().Add(pingPeriod))
			if !ok {
				// Close the websocket connection gracefully
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Printf("Error: %v", ok)
				return
			}
			// transmit to peer
			if err := c.conn.WriteJSON(resp); err != nil {
				log.Printf("Error: %v", err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(pingPeriod))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Upgrades http to websocket connection
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Handles client connections
func serveWS(hub *Hub, w http.ResponseWriter, r *http.Request) {

	// Upgrade the connection to websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return
	}

	// Initialize a new client instance
	client := &Client{
		id:               uuid.New(),
		hub:              hub,
		conn:             ws,
		ResponseReceiver: make(chan *game.Response, 256),
	}

	// Register the client with the hub
	client.hub.register <- client

	// Start the client input and output loops in go routines
	//go client.outputPump()
	go client.inputPump()
}
