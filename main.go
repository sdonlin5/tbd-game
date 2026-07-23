package main

import (
	_ "game_project/game"
	"game_project/server"
)

func main() {

	hub := server.NewHub()
	go hub.Run()
}
