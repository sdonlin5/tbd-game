// Implements the player structure
package game

import (
	_ "encoding/json"

	"github.com/google/uuid"

	_ "github.com/gorilla/websocket"
)

// Type to hold player data
type Player struct {
	PlayerID 	uuid.UUID `json:"player_id"`
	Name 		string    `json:"name"`
	Color		string    `json:"color"`
	Board		Board
}

type Coordinates struct {
	X uint8
	Y uint8
}

type Ship struct {
	Health	uint8	// number of tiles
}

type Board struct {
	tiles [10][10]bool	// validation testing
	// tiles [][]*Ship
}

