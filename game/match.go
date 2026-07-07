package game

type Player struct {
	color string
}

type Match struct {
	p1 chan Player
	p2 chan Player
	shots []Shot
}

func (m *Match) RecordShot()