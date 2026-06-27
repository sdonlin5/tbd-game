package game

import (
	"fmt"
)

// player_moves.go
type PlayerMove interface {
	// Interface to execute on player inputs
	Execute()
}

type Shot struct {
	// Location of player shot
	X uint8 `json:"x"`
	Y uint8 `json:"y"`
}

// receiver for execute
func (s Shot) Execute() {
	// Shot s receives Execute() command to play move
	fmt.Printf("Shot Fired: %d, %d", s.X, s.Y)
	// 1. validate coordinates
	// 2. check if coordinates have enemy
	// 3. update enemy state on enemy board
}
