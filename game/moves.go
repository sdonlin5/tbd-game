package game

import (
	"log"
)

// player_moves.go
type PlayerMove interface {
	// Interface to execute on player inputs
	Execute()
}

type Shot struct {
	// Location of player shot
	// Receives execute
	X uint8 `json:"x"`
	Y uint8 `json:"y"`
}

type Quit struct{}

// receiver for execute
func (s Shot) Execute() {
	// Shot s receives Execute() command to play move
	log.Printf("Shot Fired: %d, %d", s.X, s.Y)
	v := validateShot(s)
	if v {
	}
}

func (q Quit) Execute() {
	log.Printf("Player Quit")
}
