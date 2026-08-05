// Hub
// -
package server

import (
	"log"
	"slices"

	"game_project/game"

	"github.com/google/uuid"
)

type Hub struct {
	clients map[uuid.UUID]*Client
	queue   []*Client

	register   chan *Client
	unregister chan *Client
	leave      chan *Client
}

// -- Hub Methods --
// Spawns a hub from main
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		leave:      make(chan *Client),
	}
}

// Registers clients with the Hub
func (hub *Hub) registerClient(client *Client) {
	hub.clients[client.id] = client
	log.Printf("Client Registered: ID =  %v", client.id)
}

// Adds client to the queue
func (hub *Hub) queueClient(client *Client) {
	hub.queue = append(hub.queue, client)
	log.Printf("Client Queued: ID = %v", client.id)
}

// Unregisters clients
func (hub *Hub) unregisterClient(client *Client) {
	if _, ok := hub.clients[client.id]; ok {
		delete(hub.clients, client.id)
		log.Printf("Client Unregistered: ID =  %v", client.id)
	}
}

// Removes clients from the queue
func (hub *Hub) removeFromQueue(client *Client) {
	for i, c := range hub.queue {
		if c.id == client.id {
			hub.queue = slices.Delete(hub.queue, i, i+1)
		}
	}
	log.Printf("Client Removed From Queue: ID = %v", client.id)
}

// Disconnects client from the queue.
func (hub *Hub) disconnectClient(client *Client) {
	hub.removeFromQueue(client)
	hub.unregisterClient(client)
}

// Removes two clients from the front of the queue to be placed into a match
func (hub *Hub) dequeue() (*Client, *Client) {
	client1 := hub.queue[0]
	client2 := hub.queue[1]
	hub.queue = hub.queue[2:]
	return client1, client2
}

// Main hub loop, started by main
func (hub *Hub) Run() {
	// queue for players waiting to be matched
	log.Println("Start Hub")
	for {
		// Place clients into a match
		if len(hub.queue) > 1 {
			c1, c2 := hub.dequeue()
			match:= game.NewMatch(c1.id, "p1", c2.id, "p2")

			// Associate Chanels
			match.Player1.Sender = c1
			match.Player2.Sender = c2
			c1.ActionSender = match.Player1.Receiver
			c2.ActionSender = match.Player2.Receiver

			log.Printf("Match Spawned: ID = %+v", match.MatchID)
			go match.Run()
		}
		select {
		case c := <-hub.leave: // remove the client from queue and registry
			hub.disconnectClient(c)
		case c := <-hub.register: // Add new client to registry and queue
			hub.registerClient(c)
			hub.queueClient(c)
		case c := <-hub.unregister: // Remove a client from registry
			hub.unregisterClient(c)
		}

	}
}
