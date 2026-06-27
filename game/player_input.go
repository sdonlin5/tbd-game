package game

import "game_project/types"

// play_handler.go
type PlayerInput struct{}

func (ph PlayerInput) InputHandler(env *types.Envelope) {
	// Implements the types package GameHanlder interface
	switch {
	case env.Type == "shot":
		// call handle shot method
	case env.Type == "quit":
		// call confirm quit method

	}
}