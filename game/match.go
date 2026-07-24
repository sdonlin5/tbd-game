// Implements the match game thread
package game

import (
	"log"

	_ "github.com/google/uuid"
)

type Match struct {
	PlayerOne *Player
	PlayerTwo *Player

	// Read/write permission channel- Reads Action from client
	ActionReceiver chan Action

	// Write permission channel - writes Response for client
	ResponseSender chan Response
}

func (m *Match) Run() {

	for {
		action, ok := <-m.ActionReceiver
		if !ok {
			log.Printf("Error: ActionReceiver channel closed and contains no remaining values.")
			return
		}

		var p *Player
		switch action.SenderID {
		case m.PlayerOne.ClientID:
			p = m.PlayerTwo
		case m.PlayerTwo.ClientID:
			p = m.PlayerOne
		}
		switch mv := action.Move.(type) {
		case Shot:
			mv.PlayShot(p)
		case Quit:
			// call quit
		default:
			// default action
		}
	}
}
