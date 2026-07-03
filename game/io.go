package game

import (
	"log"

	"game_project/interfaces"
)

type GameHandler struct{}

func (gh GameHandler) InputHandler(env *interfaces.Envelope) {
	switch {
	case env.Type == "shot":
		// call play shot
		log.Printf("Input Type: %v", env.Type)
	case env.Type == "quit":
		log.Printf("Input Type: %v", env.Type)
	default:
		log.Printf("Input.Handler() no case found: %v", env.Type)
	}
}

func (gh GameHandler) OutputHandler(env *interfaces.Envelope) {
}

func (gh GameHandler) GetOutput() *string {
	// UserOutput Interface method
	// Creates envelope to be sent back to server via Output Handler
	return nil
}
