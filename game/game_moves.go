// Defines interfaces for player moves and their results
// Handles types of moves a player can make
// --
package game

import (
	"github.com/google/uuid"
)

// Interface type for player moves
type Move interface {
	Play() string
}

// Packge to receive data from match
type Action struct {
	SenderID uuid.UUID `json:"senderID"`
	// Type     string		`json:"type"`
	Move Move `json:"move"`
}

// Interface for results of player moves
type Result interface {
	Played() string
}

// TODO: Future Ideas
/*
	- Move ship?
	- Repair ship?
	- Look at other games for ideas to make
*/
