// Defines interfaces for player moves and their results
// Handles types of moves a player can make
// --
package game

import (
	"github.com/google/uuid"
)

// Struct to receive data
type PlayerTurn struct {
	SenderID uuid.UUID `json:"senderID"`
	// Type     string    `json:"type"`
	Move         Move `json:"move"`
	Disconnected bool
}

// Move interface
type Move interface {
	Play()
}

// Interface for outcomes of player moves

type Result interface {
	Played()
}

// Interface between client and match
type EventNotifier interface {
	Notify(resp *Response)
}

// Structured response to send to server
type Response struct {
	Type   string `json:"type"`
	Result Result `json:"result"`
}


// Player attack
// Shot X,Y coordinates player attacked at
type Shot struct {
	X uint8 `json:"x"`
	Y uint8 `json:"y"`
}

// Implements Move interface
func (s *Shot) Play() {}


// Results of player attack
type ShotResult struct {
	Kind  string
	Shot  *Shot
	Valid bool
	Hit   bool
}

// Implements Result interface
func (res *ShotResult) Played() {}


func (s *Shot) shotHandler(shooter, defender *Player) Result {
	// Shot not valid
	if !s.validateShot(defender) {
		return &InvalidShotResult{}
		}
	return s.playShot(shooter, defender)
}

// Validates the values of player attack
func (s *Shot) validateShot(defender *Player) bool {
	if int(s.X) >= len(defender.playerBoard) ||
		int(s.Y) >= len(defender.playerBoard[0]) {
		return false
	}
	return true
}

// Hit detection it detection
func (s *Shot) playShot(shooter, defender *Player) *ShotResult {
	switch {
	// location occupired
	case defender.playerBoard[s.X][s.Y].occupied:
		// previously hit - should not happen
		if defender.playerBoard[s.X][s.Y].hit {
			return &ShotResult{
				Kind:  "ShotResult",
				Shot:  s,
				Valid: true,
				Hit:   false,
			}
		} else {
			return &ShotResult{
				Kind:  "ShotResult",
				Shot:  s,
				Valid: true,
				Hit:  true,
			}
		}
	// location no occupied, miss
	default:
		// not hit
		return &ShotResult{
			Kind:  "ShotResult",
			Shot:  s,
			Valid: true,
			Hit:   false,
		}
	}
}

// -- End Attack Moves

// Signal types
// Use where a default sholdn't be hit
type NullResult struct{}
func (nr *NullResult) Played() {}

// Signal player quitting the match
type PlayerQuit struct{}
func (q *PlayerQuit) Play() {}

type PlayerQuitResult struct{}
func (res *PlayerQuitResult) Played() {}

// Invalid shot received signal - shouldn't happen
type InvalidShotResult struct{}
func (i *InvalidShotResult) Played() {}

// Signal that a shot was played out of turn - shouldn't happen
type OutOfTurnResult struct{}
func (oot *OutOfTurnResult) Played() {}

// Used to signal
type Disconnect struct{}
func (d *Disconnect) Played() {}

//-- End Disconnection
