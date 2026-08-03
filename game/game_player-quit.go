// Defines player quit move
// --
package game

// Implements Move interface
func (q PlayerQuit) Play() string {return "PlayerQuit"}

// Implements Result interface
func (res PlayerQuitResult) Played() string { return "PlayerQuitResult" }

type PlayerQuit struct{
	QuitMatch bool
}

type PlayerQuitResult struct {
	Type	string
	Quit  	bool
}
