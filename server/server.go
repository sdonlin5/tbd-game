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
			ws, err := upgrader.Upgrade(w, r, nil)
			switch {
			case err != nil:
				log.Printf("Error: %v", err)
				return
			default:
				log.Printf("Http upgraded to websocket")
			}
			defer ws.Close()
			for {
				var env types.Envelope
				err := ws.ReadJSON(&env)
				switch {
					case err != nil:
						log.Printf("Error: %v", err)
					default:
						log.Printf("JSON Read: %v", env)
				}
				gh.InputHandler(&env)
			}
		},
	)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
