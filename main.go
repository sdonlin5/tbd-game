package main

import (
	_ "fmt"
	"log"

	"game_project/game"
	_ "game_project/interfaces"
	"game_project/server"
)

func main() {
	log.Println("Main Called")
	var gh game.GameHandler
	server.RunServer(gh)
}
