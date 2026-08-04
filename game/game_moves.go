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
	//Type     string    `json:"type"`
	Move     Move      `json:"move"`
	Disconnected bool
}

// Move interface
type Move interface {
	Play() string
}

// Interface for outcomes of player moves
type Result interface {
	Played() string
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

// Implements Move interface
func (s Shot) Play() string { return "Shot" }

// Implements Result interface
func (res ShotResult) Played() string { return "ShotResult" }

// Player attack
// Shot X,Y coordinates player attacked at
type Shot struct {
	X uint8 `json:"x"`
	Y uint8 `json:"y"`
}

// Results of player attack
type ShotResult struct {
	Type  string
	Shot  *Shot
	Valid bool
	Hit   bool
}

// Validates and detects hit of player attack
func (s *Shot) PlayShot(defender *Player) *ShotResult {
	// Defaults to invalid shot
	res := &ShotResult{Type: "ShotResult", Shot: s, Valid: false, Hit: false}

	// If shot is valid, detect hit (and sink)
	if s.ValidateShot(defender) {
		res.Valid = true
		res.Hit = s.DetectHit(defender)
	}
	return res
}

// Validates the values of player attack
func (s Shot) ValidateShot(defender *Player) bool {
	if int(s.X) >= len(defender.PlayerBoard.Tiles) ||
		int(s.Y) >= len(defender.PlayerBoard.Tiles[0]) {
		return false
	}
	return true
}

// TODO: HIT DETECTION
// Hit detection for player attack
func (s *Shot) DetectHit(p *Player) bool {
	// Checks for hit on player waiting
	// TODO: Add real hit detection
	if s.X%2 == 0 || s.Y%2 == 0 {
		return true
	} else {
		return false
	}
}
// -- End Attack Moves

// TODO: CLEAN UP QUIT
// Quit Move
type Quit struct{}

// Implements Move interface
func (q PlayerQuit) Play() string { return "PlayerQuit" }

// Implements Result interface
func (res PlayerQuitResult) Played() string { return "PlayerQuitResult" }

type PlayerQuit struct {
	QuitMatch bool
}

type PlayerQuitResult struct {
	Type string
	Quit bool
}
// -- End Quit Move


// Disconnection Result
type Disconnection struct {
	player *Player
	name string
}

func (d Disconnection) Played() string {return "PlayerDisconnection" }
//-- End Disconnection