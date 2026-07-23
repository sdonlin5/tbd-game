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
	broadcast  chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:  make(map[uuid.UUID]*Client),
		register: make(chan *Client),
	}
}

// Runs a hub
func (hub *Hub) Run() {
	// Queue for players waiting to be matched
	queue := make([]*Client, 0, 2)

	for {
		// Pop clients from queue
		if len(queue) > 1 {
			// spawn new match and adds UUIDs for players
			match := &game.Match{
				PlayerOne: queue[0].id,
				PlayerTwo: queue[1].id,
			}
			// assign acttion channel to clients
			queue[0].match = match.Actions
			queue[1].match = match.Actions

			// clear the queue
			clear(queue)
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
				close(client.response)
			}
		}
	}
}
