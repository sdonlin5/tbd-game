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

		// defer close until RunServer ends
		defer ws.Close()

		// IO loop
		for {
			var input interfaces.Envelope
			var output interfaces.Envelope

			// get input sent over websocket
			readError := ws.ReadJSON(&input)
			// check for input error
			if readError != nil {
				log.Printf("Error: %v", readError)
				break
			}
			// handle input
			gh.InputHandler(&input)

			// update the output
			gh.OutputHandler(&output)

			// write output over websocket
			writeError := ws.WriteJSON(&output)
			if writeError != nil {
				log.Printf("Error: %v", writeError)
				break
			}
		}
	})
}
