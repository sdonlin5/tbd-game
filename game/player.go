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
	return &Player{ClientID: id, Name: name}
}

type Vessel interface {
	Sail() bool
}

type Ship struct {
	player *Player
	color  string
	// position []*Tile
	sunk bool
}

type HitShip struct {
	ship Ship
}

func (s *HitShip) Sail() bool { return true }

type Destroyer struct {
	ship   Ship
	health uint8 // 2
	hits   uint8
}

func (s *Destroyer) Sail() bool { return true }

type Carrier struct {
	ship   Ship
	health uint8 // 5
	hits   uint8
}

func (s *Carrier) Sail() bool { return true }

type Battleship struct {
	ship   Ship
	health uint8 // 4
	hits   uint8
}

func (s *Battleship) Sail() bool { return true }

type Cruiser struct {
	ship   Ship
	health uint8 // 3
	hits   uint8
}

func (s *Cruiser) Sail() bool { return true }

type Submarine struct {
	ship   Ship
	health uint8 // 3
	hits   uint8
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
	p.shotBoard.coords[r.shot.X][r.shot.Y].ship = HitShip{}
}

func (p *Player) updateDefender(r *ShotResult) {
	// marks that a ship has been hit
	p.playerBoard.coords[r.shot.X][r.shot.Y].ship = HitShip{}
}
