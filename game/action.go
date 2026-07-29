// Implements the player moves.
package game

import (
	_ "log"

	"github.com/google/uuid"
)

// Move
// Interface type for player moves
type Move interface {
	Play() string
}

// Shot
type Shot struct {
	X uint8 `json:"x"`
	Y uint8 `json:"y"`
}

// Satisfies Move interface
func (s Shot) Play() string {
	return "Shot"
}

// Plays the shot
// defender: player being shot at
func (s *Shot) PlayShot(defender *Player) *ShotResult {
	// Defaults to invalid shot
	res := &ShotResult{Shot: s, Valid: false, Hit: false}

	// If shot is valid, detect hit (and sink)
	if s.ValidateShot(defender) {
		res.Valid = true
		res.Hit = s.DetectHit(defender)
	}
	return res
}

// Validates shot values
func (s Shot) ValidateShot(defender *Player) bool {
	if int(s.X) >= len(defender.PlayerBoard.Tiles) ||
		int(s.Y) >= len(defender.PlayerBoard.Tiles[0]) {
		return false
	}
	return true
}

func (s *Shot) DetectHit(p *Player) bool {
	// Checks for hit on player waiting
	// TODO: Add real hit detection
	if s.X%2 == 0 || s.Y%2 == 0 {
		return true
	} else {
		return false
	}
}

// Quit Type
type Quit struct {
	Confirmed bool `json:"confirmed"`
}

func (q Quit) Play() string {
	return "Quit"
}

func (q *Quit) PlayQuit(curr *Player) *QuitResult {
	switch {
	case q.ConfirmQuit():
		return &QuitResult{IsConfirmed: true}
	case !q.ConfirmQuit():
		return &QuitResult{IsConfirmed: false}
	default:
		return nil
	}
}

// Promps user for confirmation to quit the game
func (q Quit) ConfirmQuit() bool {
	// TODO: implement method to prompt user to confirm quit
	// Prompt yes/no set q.Confirmed
	//
	//q.Confirmed = true
	return true
}

type QuitResult struct {
	IsConfirmed bool
}

// Result
// Interface for move outcomes
type Result interface {
	Played() string
}

// ShotResult
type ShotResult struct {
	Shot  *Shot
	Valid bool
	Hit   bool
}

// Satisfies Result interface
func (res *ShotResult) Played() string {
	return "Shot"
}

// Quit

// Receives data from match

type Action struct {
	// TODO: IS SenderID NEEDED HERE?
	SenderID uuid.UUID `json:"senderID"`
	Type     string
	Move     Move `json:"move"`
}
