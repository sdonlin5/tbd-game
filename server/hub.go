// Hub
// -
package server

import (
	"game_project/game"

	"github.com/google/uuid"
)

type Hub struct {
	clients    map[uuid.UUID]*Client
	register   chan *Client // channel to register clients
	unregister chan *Client
	//broadcast  chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Runs a hub
func (hub *Hub) Run() {
	// Queue for players waiting to be matched
	queue := make([]*Client, 0, 2)

	for {
		// If queue has 2 clients, spawn new match,assign channel, clear queue, reset length
		if len(queue) > 1 {
			match := &game.Match{
				PlayerOne:      game.NewPlayer(queue[0].id),
				PlayerTwo:      game.NewPlayer(queue[1].id),
				ActionReceiver: make(chan game.Action),
			}
			queue[0].ActionSender = match.ActionReceiver
			queue[1].ActionSender = match.ActionReceiver
			clear(queue)
			queue = queue[:0]
		}

		select {
		// Registers a client & adds to queue
		case client := <-hub.register:
			hub.clients[client.id] = client
			queue = append(queue, client)

			// Unregisters client and closes connection
		case client := <-hub.unregister:
			if _, ok := hub.clients[client.id]; ok {
				delete(hub.clients, client.id)
				close(client.ResponseWriter)
			}
		}
	}
}
