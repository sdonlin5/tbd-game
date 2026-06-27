package types

import (
	"encoding/json"
)

// payload.go
type Payload json.RawMessage

// envelope.go
type Envelope struct {
	Type    string  `json:"type"`
	Payload Payload `json:"payload"`
}
