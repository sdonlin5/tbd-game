// Implements the player moves.
package game

import "github.com/google/uuid"

// Generic interface for player moves
type Move interface{}

//
// Shot
//

// Data type for player shots
type Shot struct {
	X uint8 `json:"x"`
	Y uint8 `json:"y"`
}

//
// Quit
//

// Data type for Quit input
type Quit struct {
	Confirmed bool `json:"confirmed"`
}

//func (q Quit) isMove() {}

//
// Action
//

type Action struct {
	SenderID uuid.UUID	`json:"senderID"`
	Move     Move		`json:"move"`
}