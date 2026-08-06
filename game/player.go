// Implements the player structure
package game

import (
	_ "encoding/json"
	"fmt"
	_ "log"
	"sync"

	"github.com/google/uuid"
	_ "github.com/gorilla/websocket"
)

const (
	rows uint8 = 10
	cols uint8 = 10
)

// Type to hold player data
type Player struct {
	ClientID uuid.UUID
	Name     string
	// Color		string    `json:"color"`
	playerBoard [rows][cols]*Tile
	shotBoard   [rows][cols]*Tile
	//nolint:unused
	turn     bool
	Sender   EventNotifier
	Receiver chan *PlayerTurn
	mu       sync.RWMutex
}

type Board [10][10]*Tile

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
	Player *Player
	Color  string
	// position []*Tile
	Sunk bool
}

type HitShip struct {
	Ship Ship
}

func (s *HitShip) Sail() bool { return true }

type Destroyer struct {
	Ship   Ship
	Health uint8
	Hits   uint8
}

func (s *Destroyer) Sail() bool { return true }

type Carrier struct {
	Ship   Ship
	Health uint8
	Hits   uint8
}

func (s *Carrier) Sail() bool { return true }

type Battleship struct {
	Ship   Ship
	Health uint8
	Hits   uint8
}

func (s *Battleship) Sail() bool { return true }

type Cruiser struct {
	Ship   Ship
	Health uint8
	Hits   uint8
}

func (s *Cruiser) Sail() bool { return true }

type Submarine struct {
	Ship   Ship
	Health uint8
	Hits   uint8
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

func newBoard() [10][10]*Tile {
	var b Board
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			b[i][j] = newTile(i, j, false)
		}
	}
	return b
}

// confirm board
// nolint: unused
func (b *Board) printBoard() {
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if j == 9 {
				fmt.Printf(" [%v, %v] \n", b[i][j].X, b[i][j].Y)
			} else {
				fmt.Printf(" [%v, %v] ", b[i][j].X, b[i][j].Y)
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
