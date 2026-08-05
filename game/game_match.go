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
	Player1 *Player
	Player2 *Player
	// P1Receiver chan *PlayerTurn
	// P1Sender   EventNotifier
	// P2Receiver chan *PlayerTurn
	// P2Sender   EventNotifier
	// ActionReceiver chan *PlayerTurn
	current *Player
	waiting *Player
	History []*ShotResult
}

// Spawns new match, called by Hub
func NewMatch(id1 uuid.UUID, name1 string, id2 uuid.UUID, name2 string) *Match {
	return &Match{
		MatchID: uuid.New(),
		Player1: NewPlayer(id1, name1),
		Player2: NewPlayer(id1, name2),
	}
}

// Swaps current and waiting players using a temporary variable
func (m *Match) swapTurns() {
	log.Printf("Swap Called\n\nBefore:\n%+v --> current\n%+v --> waiting\n",
		m.current.ClientID,
		m.waiting.ClientID)

	temp := m.current
	m.current = m.waiting
	m.waiting = temp

	log.Printf("After:\n%+v --> current\n%+v --> waiting\n",
		m.current.ClientID, m.waiting.ClientID)
}

// Send the same response to both players
func (m *Match) notifyAll(r *Response) {
	m.current.Sender.Notify(r)
	m.waiting.Sender.Notify(r)
}

func (m *Match) routeAction(action *PlayerTurn, sender, other *Player) Result {
	switch mv := action.Move.(type) {
	// handle quit
	case PlayerQuit:
		return &PlayerQuitResult{}
	// handle shot
	case Shot:
		// handle out of turn
		if sender.ClientID != m.current.ClientID {
			return &OutOfTurnResult{}
		}
		return mv.shotHandler(sender, other)
	// Should never get here
	default:
		return &NullResult{}
	}
}

// Main game loop
func (m *Match) Run() {
	m.current = m.Player1
	m.waiting = m.Player2
	timer := time.NewTimer(60 * time.Second)
	log.Printf("Start Match")

	for {
		select {
		case action, ok := <-m.Player1.Receiver:
			if !ok {
				m.Player2.Sender.Notify(&Response{Type: "Disconnect", Result: &Disconnect{}})
				return
			}
			switch res := m.routeAction(action, m.Player1, m.Player2).(type) {
			case *OutOfTurnResult:
				// Only sent to player who input
				m.Player1.Sender.Notify(&Response{
					Type:   "OutOfTurn",
					Result: res,
				})

			// Only sent to player who input
			case *InvalidShotResult:
				m.Player1.Sender.Notify(&Response{
					Type:   "InvalidShot",
					Result: res,
				})

			// If a player quits the match, signals to both clients that the match is over
			case *PlayerQuitResult:
				quit := &Response{Type: "PlayerQuit", Result: res}
				m.notifyAll(quit)
				return

			case *ShotResult:
				s := &Response{Type: "ShotResult", Result: res}
				// TODO: handle player, board, and ship updates
				// TODO: Account for sending updates in response
				m.notifyAll(s)
				m.swapTurns()
				timer.Reset(60 * time.Second)

			default:
				m.notifyAll(&Response{
					Type:   "NullResult",
					Result: &NullResult{},
				})
			}
		case action, ok := <-m.Player2.Receiver:
			if !ok {
				m.Player1.Sender.Notify(&Response{Type: "Disconnect", Result: &Disconnect{}})
				return
			}
			switch res := m.routeAction(action, m.Player2, m.Player1).(type) {
			case *OutOfTurnResult:
				m.Player2.Sender.Notify(&Response{
					Type:   "OutOfTurn",
					Result: res,
				})
			case *InvalidShotResult:
				m.Player2.Sender.Notify(&Response{
					Type:   "InvalidShot",
					Result: res,
				})
			case *PlayerQuitResult:
				quit := &Response{Type: "PlayerQuit", Result: res}
				m.notifyAll(quit)
				return

			case *ShotResult:
				s := &Response{Type: "ShotResult", Result: res}
				// TODO: handle player, board, and ship updates
				// TODO: Account for sending updates in response
				m.notifyAll(s)
				m.swapTurns()
				timer.Reset(60 * time.Second)
			}
		// time expired
		case <-timer.C:
			m.swapTurns()
			timer.Reset(60 * time.Second)
		}
	}
}
