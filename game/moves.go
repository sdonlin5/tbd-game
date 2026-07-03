package game

import (
	"log"

	_ "game_project/interfaces"
)

// Player input actions
// ]
type Move interface {
	PlayMove()
}

type Shot struct {
	// Location of player shot
	// Receives PlayMove()
	X uint8 `json:"x"`
	Y uint8 `json:"y"`
}

type ShotResult struct {
	Shot Shot  `json:"shot"`
	Hit  uint8 `json:"hit"`
	Sink bool  `json:"sink"`
}

func (s Shot) PlayMove() *ShotResult {
	// Plays Shot
	// Testing:
	// 		IF X * Y is even --> Hit
	//      IF Y > & AND X < 5 --> SINK
	log.Printf("Shot Fired: (%d, %d)", s.X, s.Y)
	if !s.ValidateShot() {
		log.Printf("Invalid Shot Input: (%d, %d)", s.X, s.Y)
		return nil
	}
	log.Printf("Valid Shot: (%d, %d)", s.X, s.Y)
	res := ShotResult{Shot: s, Hit: 0, Sink: false}
	if (s.X*s.Y)%2 != 0 {
		res.Hit = 1
	}
	if s.Y > 7 || s.X < 5 {
		res.Sink = true
	}
	return &res
}

func (s Shot) ValidateShot() bool {
	// Validates shot coordinates are valid
	// FOR TESTING - if X or Y are >= 10
	log.Printf("ValidateShot() Called")
	if s.X >= 10 || s.Y >= 10 {
		return false
	}
	return true
}
