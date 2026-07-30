package main

import (
	"game_project/server"
	"log"
	"net/http"
)

func main() {
	log.Println("MAIN")
	hub := server.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Println("http.HandleFunc called")
		server.ServeWS(hub, w, r)

	})

	http.ListenAndServe(":8080", nil)
}
