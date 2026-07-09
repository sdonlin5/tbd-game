package server

import (
	_ "log"
	"net/http"
	"github.com/gorilla/websocket"
	"game_project/game"
)

var upgrader = websocket.Upgrader {}

type Client struct {
	// Connects websocket to game instance

	match *game.Match
	conn *websocket.Conn
}

func (c *Client) ReadInput() {

	for {

		
	}

}
