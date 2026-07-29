package main

import (
	"game_project/server"
)

func main() {

	hub := server.NewHub()
	go hub.Run()
}
