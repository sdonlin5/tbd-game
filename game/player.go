// Implements the player structure
package game

import (
	_ "encoding/json"

	"github.com/google/uuid"
	_ "github.com/gorilla/websocket"
)

// Type to hold player data
type Player struct {
	ClientID uuid.UUID `json:"client_id"`
	Name     string    `json:"name"`
	// Color		string    `json:"color"`
	playerBoard Board
	shotBoard   Board
	turn        bool
	Sender      EventNotifier
	Receiver    chan *PlayerTurn
}

func NewPlayer(id uuid.UUID, name string) *Player {
	return &Player{
		ClientID: id,
		Name: name,
		Receiver: make(chan *PlayerTurn),

	}
}



type Vessel interface {
	Sail() bool
}

type Ship struct {
	Player *Player	`json:"Player"`
	Color  string	`json:"Color"`
	//position []*Tile
	Sunk bool 		`json:"Sunk"`
}

type HitShip struct {
	Ship Ship	`json:"Ship"`
}

func (s *HitShip) Sail() bool { return true }

type Destroyer struct {
	Ship   Ship			`json:"Ship"`
	Health uint8 		`json:"Health"` // 2
	Hits   uint8 		`json:"Hits"`
}

func (s *Destroyer) Sail() bool { return true }

type Carrier struct {
	Ship   Ship			`json:"Ship"`
	Health uint8 		`json:"Health"` // 5
	Hits   uint8 		`json:"Hits"`
}

func (s *Carrier) Sail() bool { return true }

type Battleship struct {
	Ship   Ship			`json:"Ship"`
	Health uint8 		`json:"Health"` // 4
	Hits   uint8 		`json:"Hits"`
}

func (s *Battleship) Sail() bool { return true }

type Cruiser struct {
	Ship   Ship			`json:"Ship"`
	Health uint8 		`json:"Health"` // 3
	Hits   uint8 		`json:"Hits"`
}

func (s *Cruiser) Sail() bool { return true }

type Submarine struct {
	Ship   Ship			`json:"Ship"`
	Health uint8 		`json:"Health"` // 3
	Hits   uint8 		`json:"Hits"`
}

func (s *Submarine) Sail() bool { return true }

type Tile struct {
	x        uint8
	y        uint8
	occupied bool
	hit      bool
	ship     any
}

type Board struct {
	coords [10][10]*Tile
	// tiles [20][20]*Ship production board size
}

func (p *Player) updateAttacker(r *ShotResult) {
	// place hit ship indcator at shot coordinates
	p.shotBoard.coords[r.Shot.X][r.Shot.Y].ship = HitShip{}
}

func (p *Player) updateDefender(r *ShotResult) {
	// marks that a ship has been hit
	p.playerBoard.coords[r.Shot.X][r.Shot.Y].ship = HitShip{}
}
