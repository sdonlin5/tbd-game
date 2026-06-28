package game

import (
	"log"
)

func playGame() {
	// implements the game loop
	log.Printf("playGame Called")
	
}

//nolint:unused
func validateShot(s Shot) bool {
	// check for valid tiles
	log.Printf("validateShot called: %d, %d", s.X, s.Y)
	return true
}

//nolint:unused
func checkHit(s Shot) bool {
	log.Printf("checkHit called: %d, %d", s.X, s.Y)
	return true
}
