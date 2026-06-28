package server

import (
	"log"
	"net/http"

	"game_project/types"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RunServer(gh types.GameHandler) {
	http.HandleFunc(
		"/ws",
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Handler Called")

			// Upgrade to websocket connection
			ws, upgradeErr := upgrader.Upgrade(w, r, nil)
			switch {
			case upgradeErr != nil:
				log.Printf("Error: %v", upgradeErr)
				return
			default:
				log.Printf("Http upgraded to websocket")
			}

			defer ws.Close()

			// Server connection loop
			for {
				// Input
				var envIn types.Envelope
				readErr := ws.ReadJSON(&envIn)
				switch {
					case readErr != nil:
						log.Printf("Error: %v", readErr)
					default:
						log.Printf("JSON Read: %v", envIn)
				}
				gh.InputHandler(&envIn)

				// Output
				var envOut types.Envelope
				writeErr := ws.WriteJSON(&envOut)
					switch {
						case writeErr != nil:
							log.Printf("Error: %v", writeErr)
						default:
							log.Printf("JSON Write: %v", envOut)
					}


			}
		},
	)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
