// Implements the player structure
package game

import (
	_ "encoding/json"
	"log"

	"github.com/google/uuid"
	_ "github.com/gorilla/websocket"
)

const (
	rows uint8 = 10
	cols uint8 = 10
)

// Type to hold player data
type Player struct {
	ClientID uuid.UUID `json:"client_id"`
	Name     string    `json:"name"`
	// Color		string    `json:"color"`
	playerBoard [rows][cols]*Tile
	shotBoard   [rows][cols]*Tile
	//nolint:unused
	turn     bool
	Sender   EventNotifier
	Receiver chan *PlayerTurn
}

func NewPlayer(id uuid.UUID, name string) *Player {
	p := Player{
		ClientID: id,
		Name:     name,
		Receiver: make(chan *PlayerTurn),
	}
	p.playerBoard = newBoard()
	p.shotBoard = newBoard()
	return &p
}

type Vessel interface {
	Sail() bool
}

type Ship struct {
	Player *Player `json:"Player"`
	Color  string  `json:"Color"`
	// position []*Tile
	Sunk bool `json:"Sunk"`
}

type HitShip struct {
	Ship Ship `json:"Ship"`
}

func (s *HitShip) Sail() bool { return true }

type Destroyer struct {
	Ship   Ship  `json:"Ship"`
	Health uint8 `json:"Health"` // 2
	Hits   uint8 `json:"Hits"`
}

func (s *Destroyer) Sail() bool { return true }

type Carrier struct {
	Ship   Ship  `json:"Ship"`
	Health uint8 `json:"Health"` // 5
	Hits   uint8 `json:"Hits"`
}

func (s *Carrier) Sail() bool { return true }

type Battleship struct {
	Ship   Ship  `json:"Ship"`
	Health uint8 `json:"Health"` // 4
	Hits   uint8 `json:"Hits"`
}

func (s *Battleship) Sail() bool { return true }

type Cruiser struct {
	Ship   Ship  `json:"Ship"`
	Health uint8 `json:"Health"` // 3
	Hits   uint8 `json:"Hits"`
}

func (s *Cruiser) Sail() bool { return true }

type Submarine struct {
	Ship   Ship  `json:"Ship"`
	Health uint8 `json:"Health"` // 3
	Hits   uint8 `json:"Hits"`
}

func (s *Submarine) Sail() bool { return true }

type Tile struct {
	X        int
	Y        int
	occupied bool
	hit      bool
	//nolint:unused
	ship any
}

func newTile(x int, y int, occupied bool) *Tile {
	return &Tile{
		X:        x,
		Y:        y,
		occupied: occupied,
	}
}

type Board [10][10]*Tile

func newBoard() [10][10]*Tile {
	var b Board
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			b[i][j] = newTile(i, j, false)
		}
	}
	return b
}

// nolint: unused
func (b *Board) printBoard() {
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if i == 9 {
				log.Printf("[%v, %v]\n", b[i][j].X, b[i][j].Y)
			} else {
				log.Printf("[%v, %v]", b[i][j].X, b[i][j].Y)
			}
		}
	}
}

//nolint:unused
func (p *Player) updateAttacker(r *ShotResult) {
	// place hit ship indcator at shot coordinates
	p.shotBoard[r.Shot.X][r.Shot.Y].ship = HitShip{}
}

//nolint:unused
func (p *Player) updateDefender(r *ShotResult) {
	// marks that a ship has been hit
	p.playerBoard[r.Shot.X][r.Shot.Y].ship = HitShip{}
}
