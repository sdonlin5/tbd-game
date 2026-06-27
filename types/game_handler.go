package types

// game_handler.go
type GameHandler interface {
	// Handler to pass input from server package to game package
	InputHandler(env *Envelope)
}
