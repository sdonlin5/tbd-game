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

	// Outbound channel to the match for the client connection
	match chan <- game.Action

	// Outbound channel to the websocket for the client connection
	response chan <- game.Response // Receives



}

// Sets the initial deadline for pong
func (c *Client) setWaitTime() error {
	return c.conn.SetReadDeadline(time.Now().Add(pongWait))
}

// Sets the pong handler for messages received from the peer
func (c *Client) initKeepAlive() {
	c.conn.SetPongHandler(func(string) error {
		// After each pong is received, reset the the deadline
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})
}

func (c *Client) setTurnTime() {
	c.conn.SetWriteDeadline(time.Now().Add(turnTime))
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

	// Wrapper to set the initial heartbeat
	c.setWaitTime()
	// Wrapper to configure the pong handler
	c.initKeepAlive()

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
		if c.match == nil {
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

		// Send the input to the channel
		} else {
			c.match <- game.Action { SenderID: c.id, Move: move }
		}
	}
}

// Pumps messages from the hub to the websocket connection.
func (c *Client) ouputPump() {

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
		//select {
		//case update, ok := <- c.send:
		//	c.setTurnTime()
		//	if !ok {
		//		// Channel closed by Hub
		//		c.conn.Close()
		//	}
		//}




	}

}

// Upgrades http to websocket connection
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Handles websocket requests from peer application
func serveWS(hub *Hub, w http.ResponseWriter, r *http.Request) {

	// Upgrade the connection to websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return
	}

	client := &Client{id: uuid.New(), hub: hub, conn: ws} // creates client reference
	client.hub.register <- client                 // registers the new client

	// Start the client input and output loops in go routines
	//go client.outputPump()
	go client.inputPump()
}
