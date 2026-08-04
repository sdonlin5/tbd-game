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
	ActionReceiver chan *PlayerTurn
	Current        *Player
	Waiting        *Player
	P1Notifier     EventNotifier
	P2Notifier     EventNotifier
	History        []*ShotResult
}

// Swaps current and waiting players using a temporary variable
func (m *Match) swapTurns() {
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

// Main game loop
func (m *Match) Run() {
	m.Current = m.PlayerOne
	m.Waiting = m.PlayerTwo
	timer := time.NewTimer(60 * time.Second)
	log.Printf("Start Match")

	for {
		select {
		case action, ok := <-m.ActionReceiver:
			if !ok {
				log.Printf("ActionReceiver Channel Closed")
				return
			}
			if action.Disconnected {
				res := Disconnection{
					player: m.Current,
					name:   m.Current.Name,
				}
				switch m.Current {
				case m.PlayerOne:
					m.P1Notifier.Notify(&Response{
						Type:   "Disconnection",
						Result: res,
					})
				case m.PlayerTwo:
					m.P2Notifier.Notify(&Response{
						Type:   "Disconnection",
						Result: res,
					})
				}
			}

			// Wrong player
			if action.SenderID != m.Current.ClientID {
				log.Printf("Out of turn input received")
				continue // until correct
			}

			switch mv := action.Move.(type) {
			case Shot:
				result := mv.PlayShot(m.Waiting)
				switch {
				// invalid input - shouldn't happen from app
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

				case result.Valid:
					m.History = append(m.History, result)
					//resp := Response{Type: result.Type, Result: result}
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
				log.Printf("Match Ended By: %v", m.Current.Name)
				return
			}

		case <-timer.C:
			log.Println("Time Expired!")
			m.swapTurns()
			timer.Reset(60 * time.Second)
		}
	}
}
