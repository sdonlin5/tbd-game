package interfaces

import (
	"encoding/json"
)

type Payload json.RawMessage

type Envelope struct {
	Type    string
	Payload Payload
}

type Handler interface {
	// Handlers update envelope passed via pointer
	InputHandler(*Envelope)
	OutputHandler(*Envelope)
}

type Response interface {
	//
	Responder() *Envelope
}


