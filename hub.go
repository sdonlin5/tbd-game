package gameproject


type Hub struct {
	players 	map[*game.Player]bool
	register 	chan *game.Player		// channel to register player
}

func newHub() *Hub {
	return &Hub {
		players: make(map[*game.Player]bool),
		register: make(chan *game.Player),
	}
}

func (hub *Hub) run() {
	for {

	}

}