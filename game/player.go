// Implements the player structure
package game

import (
	_ "encoding/json"

	"github.com/google/uuid"
	_ "github.com/gorilla/websocket"
)

// Type to hold player data
type Player struct {
	ClientID uuid.UUID `json:"client_id"`
	Name     string    `json:"name"`
	// Color		string    `json:"color"`
	PlayerBoard Board
	Shots       Board
	Turn        bool
}

func NewPlayer(id uuid.UUID, name string) *Player {
	return &Player{ClientID: id, Name: name}
}

type Coordinates struct {
	X uint8
	Y uint8
}

type Ship struct {
	Health uint8 `json:"health"` // number of tiles ship uses
}

type Board struct {
	Tiles [10][10]*Ship
	// tiles [20][20]*Ship production board size
}

func (b *Board) MarkHit() {
	//
}
