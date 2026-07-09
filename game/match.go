package game

import (
	"game_project/interfaces"
)

type Player struct {
	PlayerID int
	Name string
	color string
}



type Match struct {
	p1 chan interfaces.Envelope
	p2 chan interfaces.Envelope
	//	shots []Shot // slice to hold shots
}



