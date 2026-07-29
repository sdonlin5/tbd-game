// Implements the match game thread

// TODO: Refactor for neatness

package game

import (
	"log"
	_ "math/rand/v2"
	"time"

	_ "github.com/google/uuid"
)


type Match struct {
	// Data type to control an instance of a match
	PlayerOne      *Player
	PlayerTwo      *Player
	ActionReceiver chan *Action			// Bidirectional channel, reads Action written by instance of Client

	Current        *Player				// State pointers
	Waiting        *Player

	P1Notifier		EventNotifier
	P2Notifier 		EventNotifier
	// History             []*ShotResult		TBD use
}

type Response struct {
	// Response written back to client
	Type       string    `json:"type"`
	Result     Result    `json:"result"`
}

type EventNotifier interface {
	// Interface to send data to Client,  Notify() implemented by Client
	Notify(resp *Response)
}


func (m *Match) swapTurns() {
	// Swaps current and waiting players using a temporary variable
	temp := m.Current
	m.Current = m.Waiting
	m.Waiting = temp
}


func (m *Match) Run() {
	// Game loop for a match between two players.
	// TODO: Create function to randomize
	m.Current = m.PlayerOne
	m.Waiting = m.PlayerTwo

	//m.History = make([]*ShotResult, 0)
	timer := time.NewTimer(60 * time.Second)

	// Main game loop
	for {
		select {

		case action, ok := <-m.ActionReceiver:
			// Ignore input in channel from player waiting
			if action.SenderID != m.Current.ClientID{
				log.Printf("Ignored out-of-turn move: %v from player: %v", action.Move.Play(), action.SenderID)
				continue
			}

			if !ok {
				log.Printf("Error: ActionReceiver channel closed and contains no remaining values.")
				return
			}
			// Current player input a move
			switch mv := action.Move.(type) {

			// Move is a shot, process
			case Shot:
				result := mv.PlayShot(m.Waiting)
				if result.Valid {
					// Both clients receive the same response until peer application is available for customization
					m.P1Notifier.Notify(&Response{
						Type:       "ShotResult",
						Result:     result,
					})

					m.P2Notifier.Notify(&Response{
						Type:       "ShotResult",
						Result:     result,
					})

					/*					Uncomment for ID
					m.P1Notifier.Notify(&Response{
						ReceiverID: m.Current.ClientID,
						Type:       "ShotResult",
						Result:     result,
					})

					m.P2Notifier.Notify(&Response{
						ReceiverID: m.Current.ClientID,
						Type:       "ShotResult",
						Result:     result,
					})
*/

					m.swapTurns()
					timer.Reset(60 * time.Second)

				} else {
					// Invalid shot played
					// Can play another move if time permits
					switch  m.Current.ClientID{
						case m.PlayerOne.ClientID:
							m.P1Notifier.Notify(
								&Response{
								Type:       "ShotResult",
								Result:     result,
							})

						default:
							m.P2Notifier.Notify(
								&Response{
								Type:       "ShotResult",
								Result:     result,
							})

/*							Uncomment for ID
							m.P1Notifier.Notify(
								&Response{
								ReceiverID: m.Current.ClientID,
								Type:       "ShotResult",
								Result:     result,
							})

						default:
							m.P2Notifier.Notify(
								&Response{
								ReceiverID: m.Waiting.ClientID,
								Type:       "ShotResult",
								Result:     result,
							})

*/


						}
					}

			case Quit:
				// Player selected to quit the fame
				result := mv.PlayQuit(m.Current)
				if result.IsConfirmed {
					// TODO: function to end match
					break
				}
			}
			case <-timer.C:
				// Time expired for player to make a move
				log.Println("Time Expired!")
				m.swapTurns()
				timer.Reset(60 * time.Second)

			// IF time remains on timer, fallthrough to top of for loop
			}
		}
	}

