package types

// game_handler.go
type InputHandler interface {
	// Handler to pass input from server package to game package
	InputHandler(env *Envelope)
}

type OutputHandler interface {
	OutputHandler()
}