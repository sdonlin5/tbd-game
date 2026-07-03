package game

import (
	"encoding/json"
	"log"

	"game_project/interfaces"
)

type GameHandler struct{}

func (gh GameHandler) InputHandler(env *interfaces.Envelope) {
	switch {
	case env.Type == "shot":
		// call PlayMove
		log.Printf("Input Type: %v", env.Type)
		var s Shot
		err := json.Unmarshal(env.Payload, &s)
		if err != nil {
			log.Printf("Error % v", err)
		}




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
