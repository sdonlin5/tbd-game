// Implements the player moves.
package game

import (
	"log"

	"github.com/google/uuid"
)

// Generic interface for player moves
type Move interface {
	SenderID()
}

//
// Shot
//

// Data type for player shots
type Shot struct {
	X uint8 `json:"x"`
	Y uint8 `json:"y"`
}

// Plays the shot
func (s Shot) PlayShot(p *Player) {
	if !s.ValidateShot(p) {
		log.Printf("Invalid Shot Fired: [ %v, %v ]", s.X, s.Y)
	}
	log.Printf("Shot at (%v, %v) %v", s.X, s.Y, s.DetectHit(p))
}

// Validates shot values
func (s Shot) ValidateShot(p *Player) bool {
	if int(s.X) >= len(p.Board.Tiles) ||
		int(s.Y) >= len(p.Board.Tiles[0]) {
		return false
	}
	return true
}

func (s *Shot) DetectHit(p *Player) string {
	if s.X%2 == 0 || s.Y%2 == 0 {
		return "Hit"
	}
	return "Miss"
}

func (s Shot) SenderID() {}

//
// Quit
//

// Data type for Quit input
type Quit struct {
	Confirmed bool `json:"confirmed"`
}

func (q Quit) SenderID() {}

func (q Quit) ConfirmQuit() {
	// TODO: implement method to prompt user to confirm quit
	// Prompt yes/no set q.Confirmed
	//
	//q.Confirmed = true
}

//
// Action
//

type Action struct {
	SenderID uuid.UUID `json:"senderID"`
	Move     Move      `json:"move"`
}
