package server

import (
	_ "log"
	"net/http"
	"github.com/gorilla/websocket"
	"game_project/game"
)

var upgrader = websocket.Upgrader{
	// Upgrades HTTP connection to websocket
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleConnection() {
	
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}
		p := game.Player{WS: ws,}
	}
}
