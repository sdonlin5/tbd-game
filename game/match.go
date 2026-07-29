// Implements the match game thread
package game

import (
	"log"
	_ "math/rand/v2"
	"time"
	"github.com/google/uuid"
)

// Match
//
//	PlayerOne & PlayerTwo:	Player pointers
//	ActionReceiver:			read/write chan		reads Acton from Client
//	ResponseSender:			write chan			writes Resonse to Client
//	Attacker & Defender: 	Pointers to manage turn
//	ADD HERE
//	History:				Slice of shots played
type Match struct {
	PlayerOne              *Player
	PlayerTwo              *Player
	ActionReceiver         chan *Action
	Current                *Player
	Waiting                *Player

	p1Notifier				chan <- *Response
	p2Notifier				chan <- *Response
	// History             []*ShotResult		TBD use
}

// Data type to write to ResponseSender
type Response struct {
	ReceiverID uuid.UUID `json:"receiverID"`
	Type       string    `json:"type"`
	Result     Result    `json:"result"`
}


// Interface to send updates to Client domain
// Notify implemented client.go
type EventNotifier interface {
	Notify (resp *Response)
}


// TODO: Function to select first turn
//nolint:unused
// func (m *Match) selectFirstTurn() *Player {
// 	return &Player{}
// }

// swapTurns: Swaps current and waiting player
// -
func (m *Match) swapTurns() {
	temp := m.Current
	m.Current = m.Waiting
	m.Waiting = temp
}

//func (m *Match) broadcastResponse(response *Response) {
//	m.ResponseSender <- response
//}

// Run()
// Runs the game loop
// -
func (m *Match) Run() {

	// TODO: Create function to randomize
	m.Current = m.PlayerOne
	m.Waiting = m.PlayerTwo

	//m.History = make([]*ShotResult, 0)
	timer := time.NewTimer(60 * time.Second)

	// Main game loop
	for {
		select {
		case action, ok := <-m.ActionReceiver:
			if !ok {
				log.Printf("Error: ActionReceiver channel closed and contains no remaining values.")
				return
			}
			// Current player input a move
			switch mv := action.Move.(type) {

				// Move is a shot
				case Shot:
					result := mv.PlayShot(m.Waiting)
					if result.Valid {					// valid shot played => Both players get updated
						timer.Reset(60 * time.Second)
						m.swapTurns()

						resp := Response{
							ReceiverID: m.Current.ClientID,
							Type: "Shot",
							Result: result,
						}
						m.p1Notifier <- &resp
						m.p2Notifier <- &resp
					} else {
						switch {
						case m.Current.ClientID == m.PlayerOne.ClientID:
							m.p1Notifier <- &Response{
								ReceiverID: m.Current.ClientID,
								Type: "Shot",
								Result: result,
							}
						case m.Current.ClientID == m.PlayerTwo.ClientID:
							m.p2Notifier <- &Response{
								ReceiverID: m.Current.ClientID,
								Type: "Shot",
								Result: result,
							}
						}
					}

				// Move is Quit
				case Quit:
					result := mv.PlayQuit(m.Current)
					if result.IsConfirmed {
						// TODO: function to end match
						break
					}
			}
		// Time expired for player to make a move
		case <- timer.C:
			log.Println("Time Expired!")
			m.swapTurns()
		}
	}
}

