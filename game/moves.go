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

// Game history record
type ShotRecord struct {
	Shot   Shot `json:"shot"`
	Result ShotResult
}

// Result of shot
type ShotResult struct {
	Hit  int
	Sink bool
}

// Communicated back to client
type ShotResponse struct {
	X    uint8
	Y    uint8
	Hit  bool
	Sink bool
}

func (s Shot) PlayMove() *ShotResponse {
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

	response := ShotResponse {		// construct response
		X: s.X,
		Y: s.Y,
		Hit: false,
		Sink: false,
	}
	if (s.X*s.Y)%2 != 0 {						// hit if even
		response.Hit = true
	}
	if s.Y > 7 && s.X < 5 {						// sink if Y > 7 AND x < 5
		response.Sink = true
	}
	return &response
}

func (s Shot) ValidateShot() bool {
	// Validates shot coordinates are valid
	// FOR TESTING - if X or Y are >= 10
	log.Printf("ValidateShot() Called")
	if s.X >= 10 || s.Y >= 10 {					// Invalid if coordinate >= 10
		return false
	}
	return true
}
