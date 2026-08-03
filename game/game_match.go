// Implements the match game thread

// TODO: Refactor for neatness

package game

import (
	"log"
	_ "math/rand/v2"
	"time"

	"github.com/google/uuid"
)

type Match struct {
	MatchID uuid.UUID
	// Data type to control an instanatch
	PlayerOne      *Player
	PlayerTwo      *Player
	ActionReceiver chan *Action // Bidirectional channel, reads Action written by instance of Client

	Current *Player // State pointers
	Waiting *Player

	P1Notifier EventNotifier
	P2Notifier EventNotifier
	History    []*ShotResult
}

func (m *Match) swapTurns() {
	// Swaps current and waiting players using a temporary variable
	log.Printf("Swap Called\n\nBefore:\n%+v --> current\n%+v --> waiting\n",
		m.Current.ClientID,
		m.Waiting.ClientID)

	temp := m.Current
	m.Current = m.Waiting
	m.Waiting = temp

	log.Printf("After:\n%+v --> current\n%+v --> waiting\n",
		m.Current.ClientID, m.Waiting.ClientID)
}

// Send the same response to both players
func (m *Match) NotifyAll(r *Response) {
	m.P1Notifier.Notify(r)
	m.P2Notifier.Notify(r)
}

// Send sends signal to both players
func (m *Match) Disconnect() {
	disconnect := Disconnect{}
	response := &Response{Type: "Disconnect", Result: disconnect}
	m.NotifyAll(response)
}

func (m *Match) Run() {
	// Game loop for a match between two players.
	// TODO: Create function to randomize
	m.Current = m.PlayerOne
	m.Waiting = m.PlayerTwo
	m.History = make([]*ShotResult, 0)
	timer := time.NewTimer(60 * time.Second)
	log.Printf("Start Match")

	// Main Loop
	for {

		log.Printf("state:\ncurrent: %+v\nwaiting: %+v", m.Current.ClientID, m.Waiting.ClientID)
		select {
		case action, ok := <-m.ActionReceiver:
			if !ok {
				log.Printf("NO MATCH END: ActionReceiver Channel Closed")
				m.Disconnect()
				return
			}
			// Wrong player
			if action.SenderID != m.Current.ClientID {
				log.Printf("out of turn input received")
				continue // until correct
			}

			switch mv := action.Move.(type) {
			case ClientCriticalError:
				// call disconnect sequeence and return to end match
			case Shot:
				result := mv.PlayShot(m.Waiting)
				switch { // valid switch
				case !result.Valid:
					switch { // current
					case m.Current.ClientID == m.PlayerOne.ClientID:
						m.P1Notifier.Notify(
							&Response{
								Type:   result.Type,
								Result: result,
							},
						)

					case m.Current.ClientID == m.PlayerTwo.ClientID:
						m.P2Notifier.Notify(
							&Response{
								Type:   result.Type,
								Result: result,
							},
						)
					}
				// attacker can play another shot as long as time remains

				case result.Valid:
					m.History = append(m.History, result)
					// resp := Response{Type: result.Type, Result: result}
					switch m.Current.ClientID {
					// Current = PlayerOne, Valid = true, hit = true
					case m.PlayerOne.ClientID:
						// m.NotifyAll(&resp)

						m.P1Notifier.Notify(&Response{
							Type:   result.Type,
							Result: result,
						})
						m.P2Notifier.Notify(&Response{
							Type:   result.Type,
							Result: result,
						})
						m.swapTurns()
						timer.Reset(60 * time.Second)

					// Current = PlayerTwo, Valid = true, hit = true
					case m.PlayerTwo.ClientID:
						m.P2Notifier.Notify(&Response{
							Type:   result.Type,
							Result: result,
						})
						m.P1Notifier.Notify(&Response{
							Type:   result.Type,
							Result: result,
						})
						m.swapTurns()
						timer.Reset(60 * time.Second)

					}
				}

			case PlayerQuit:
				result := PlayerQuitResult{Type: "PlayerQuit", Quit: true}
				response := Response{Type: result.Type, Result: result}

				m.P1Notifier.Notify(
					&Response{
						Type:   result.Type,
						Result: result,
					},
				)
				m.P2Notifier.Notify(
					&Response{
						Type:   result.Type,
						Result: result,
					},
				)
				log.Printf("%v Ended the Match.", m.Current.Name)
				log.Println("Game Over")

				m.TerminateMatch()
				return
			}

		case <-timer.C:
			log.Println("Time Expired!")
			m.swapTurns()
			timer.Reset(60 * time.Second)
		}
	}
}
