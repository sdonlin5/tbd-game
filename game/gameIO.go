// game package IO handling
//

package game

import "game_project/types"

type HandleInput struct{}

func (ph HandleInput) InputHandler(env *types.Envelope) {
	// Implements the InputHandler interface declared in gameHandler.go

	switch {
	case env.Type == "shot":
		// call handle shot method
	case env.Type == "quit":
		// call confirm quit method

	}
}
