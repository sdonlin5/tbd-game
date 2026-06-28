// game package IO handling
//

package game

import (
	"encoding/json"
	_ "fmt"
	"log"

	"game_project/types"
)

type IOHandler struct{}

func (io IOHandler) InputHandler(env *types.Envelope) {
	// Implements the InputHandler interface declared in gameHandler.go
	log.Printf("Input Handler Called: %v", env)

	p := env.Payload
	switch {
	case env.Type == "shot":
		// call handle shot method
		// call validator
		// call hit detector
		log.Printf("Payload: %v", p)
	case env.Type == "quit":
		// call quit confirm
		log.Printf("Quit Called")
	}
}


func (io IOHandler) OutputHandler(r Response) (*types.Envelope, error) {
	//
	log.Printf("Output Handler Called: %v", r)
	env, err := r.Respond()
	if err != nil {
		return nil, err
	}
	return env, nil
}


type Response interface {
	Respond() (*types.Envelope, error)
}

// (receiver) InterfaceMethod() Return {}
func (s ShotResponse) Respond() (*types.Envelope, error) {
	var env types.Envelope
	env.Type = "shot"
	p, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	env.Payload = p
	return &env, nil
}

func (q QuitResponse) Respond() (*types.Envelope, error) {
	var env types.Envelope
	env.Type = "quit"
	p, err := json.Marshal(q)
	if err != nil {
		return nil, err
	}
	env.Payload = p
	return &env, nil
}



type TurnResponse struct {
	Player string	`json:"player"`
}

type ShotResponse struct {
	TurnResponse
	X int		`json:"x"`
	Y int		`json:"y"`
	Hit int 	`json:"hit"`
	Sink bool	`json:"sink"`
}


type QuitResponse struct {
	TurnResponse
	Confirmed bool	`json:"confirmed"`
}


func (io IOHandler) OutputHandler(r Response) *types.Envelope {
	//
	log.Printf("Output Handler Called: %v", r)

}
