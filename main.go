package main

import (
	"game_project/server"
	"game_project/game"
)

func main() {
	gameHandler := game.IOHandler{}
	server.RunServer(gameHandler)
}
