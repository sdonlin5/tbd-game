// Hub
// -
package server

import (
	"game_project/game"
	"log"

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
	log.Println("Hub Started")
	queue := make([]*Client, 0, 2)
	for {

		// If queue has 2 clients, spawn new match,assign channel, clear queue, reset length
		if len(queue) > 1 {
			log.Println("Queue > 1")
			client1 := queue[0]
			client2 := queue[1]
			match := &game.Match{

				PlayerOne:      game.NewPlayer(client1.id, "p1"),
				PlayerTwo:      game.NewPlayer(client2.id, "p2"),
				ActionReceiver: make(chan *game.Action),
				P1Notifier:     client1,
				P2Notifier:     client2,
			}
			log.Println("Match Created")

			queue[0].ActionSender = match.ActionReceiver
			queue[1].ActionSender = match.ActionReceiver
			clear(queue)
			queue = queue[:0]
			match.Run()
		}
		select {
		// Registers a client & adds to queue
		case client := <-hub.register:
			log.Println("REGISTER")
			hub.clients[client.id] = client
			queue = append(queue, client)
			log.Printf("Client Registered")

		// Unregisters client and closes connection
		case client := <-hub.unregister:
			log.Println("UNREGISTER")
			if _, ok := hub.clients[client.id]; ok {
				delete(hub.clients, client.id)
				log.Println("Client Unregistered")
				close(client.ResponseReceiver)
			}
		}

	}

}
