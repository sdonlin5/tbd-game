// Implements client
// -
package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"game_project/game"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Constat values for ping / pong between peer and server
const (
	turnTime   = 60 * time.Second
	pongWait   = time.Second * 60
	pingPeriod = (pongWait * 9) / 10
)

// Struct to receive input JSON from the websocket connection.
type Envelope struct {
	Type    string          // type of input the user entered
	Payload json.RawMessage // data of the input, e.g. shot coords
}

// Intermediary between websocket connection and hub.
type Client struct {
	id               uuid.UUID // unique identifier for each client
	hub              *Hub
	conn             *websocket.Conn
	ActionSender     chan<- *game.PlayerTurn
	ResponseReceiver chan *game.Response
}

// Satisfies EventNotifier interface
func (c *Client) Notify(resp *game.Response) {
	select {
	case c.ResponseReceiver <- resp:
	default:
	}
}

// inputPump pumps input received from the websocket connection to the hub.
func (c *Client) inputPump() {
	defer func() {
		c.hub.leave <- c
		c.conn.Close()
	}()

	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("%v", err)
		return
	}
	c.conn.SetPongHandler(func(string) error {
		// After each pong is received, reset the the deadline
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		var input Envelope
		readErr := c.conn.ReadJSON(&input)
		log.Printf("ReadJSON: %v", input)

		// Gaurds against sending to a match who's goroutine ended
		if readErr != nil {
			select {
				case c.ActionSender <- &game.PlayerTurn{
				SenderID: 	c.id,
				Disconnected: true,
			}:
				default:
			}
			return
		}

		// Gaurds against the client not having a match
		if c.ActionSender == nil {
			log.Printf("Error: %v sent payload before joining match!", c.id)
			continue
		}
		// Intermediate storage
		var move game.Move
		var inputError error
		log.Printf("ClientID: %v  - Input: %v", c.id, input)

		// parse the payload
		switch input.Type {
		case "Shot":
			var s game.Shot
			inputError = json.Unmarshal(input.Payload, &s)
			move = &s

		case "PlayerQuit":
			var q game.PlayerQuit
			inputError = json.Unmarshal(input.Payload, &q)
			log.Printf("Received: %v", input)
			move = &q

		default:
			log.Printf("Unknown input type: %v", input.Type)
			continue
		}
		// Error if json malformed
		if inputError != nil {
			log.Printf("Error %v", inputError)
			continue

		} else {
			select {
			case c.ActionSender <- &game.PlayerTurn{SenderID: c.id, Move: move}:
				log.Printf("Input sent to match")
			default:
				return
			}
		}
	}
}

// Pumps messages from the hub to the websocket connection
func (c *Client) outputPump() {
	// sets heartbeat interval
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	// Output loop
	for {
		select {
		case resp, ok := <-c.ResponseReceiver:
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{});
				err != nil {
					log.Printf("Error: %v", ok)
				}
				return
			}

			if e := c.conn.SetWriteDeadline(time.Now().Add(pingPeriod)); e != nil {
				log.Printf("%v", e)
				return
			}

			// transmit to peer
			if err := c.conn.WriteJSON(resp); err != nil {
				log.Printf("Error: %v", err)
				return
			} // Handle WriteJSON error

		case <-ticker.C:
			if e := c.conn.SetWriteDeadline(time.Now().Add(pingPeriod)); e != nil {
				log.Printf("%v", e)
				return
			} // handle SetWriteDeadline Error

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Error: %v", err)
				return
			} // handle WriteMessage Error
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
func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
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
	log.Printf("client spawned: %+v", client.id)

	// Register the client with the hub
	client.hub.register <- client

	// Start the client input and output loops in go routines
	go client.outputPump()
	go client.inputPump()
}
