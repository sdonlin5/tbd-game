package main

import (
	"game_project/game"
	"game_project/server"
)

func main() {
	playHandler := game.PlayerInput{}
	server.RunServer(playHandler)
}

