// Defines the player move: Quit
// --
package game

type EndMatch struct{}

func (q EndMatch) Play() string { return "End" }

func (q *EndMatch) PlayEndMatch() *EndMatchResult {
	return &EndMatchResult{
		Type: "End Match Result",
		Quit: true,
	}
}

type EndMatchResult struct {
	Type string
	Quit bool
}

// Satisfies result interface
func (res *EndMatchResult) Played() string { return "EndMatchResult" }
