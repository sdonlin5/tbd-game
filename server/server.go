package server

import (
	"log"
	"net/http"

	"game_project/interfaces"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// Upgrades HTTP connection to websocket
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RunServer(gh interfaces.Handler) {
	http.HandleFunc("/ws", func(w http.ResponseWriter,
		r *http.Request,
	) {
		// upgrade connection
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Upgrade Failed: %v", err)
			return
		} else {
			log.Printf("Upgraded to websocket")
		}

		defer ws.Close() // ws only closes afer hanlder function ends

		for {
			// Input
			var input interfaces.Envelope
			readError := ws.ReadJSON(&input)
			if readError != nil {
				log.Printf("JSON errRead: %v", readError)
			}
			gh.InputHandler(&input)

			// output
			var output interfaces.Envelope
			gh.OutputHandler(&output) // updates envelope

		}
	})
}
