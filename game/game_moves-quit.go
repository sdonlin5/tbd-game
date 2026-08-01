// Defines the player move: Quit
// --
package game

type Quit struct{}

func (q Quit) Play() string { return "Quit" }

type QuitResult struct {
	Quit bool
}

// Satisfies result interface
func (res *QuitResult) Played() string { return "QuitResult" }
