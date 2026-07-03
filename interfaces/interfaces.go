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
	InputHandler(*Envelope)
	OutputHandler(*Envelope)
}
