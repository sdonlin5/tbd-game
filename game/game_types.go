package game

// Move Interface & Types
type Move interface {
	Play() string
}

// Interface for output from player action
type Result interface {
	Played() string
}

// Signals an error in client package to match
type ClientCriticalError struct {}

// Satisfies move interface and allows sending via ActionSender
func (ce ClientCriticalError) Play() string { return "ClientCriticalError"}



// Satisfies Result
//func (d Disconnect) Played() string { return "Disconnect" }


// Interface to send data to Client,  Notify() implemented by Client
type EventNotifier interface {
	Notify(resp *Response)
}

// Types

// Data sent back to client
type Response struct {
	Type   string `json:"type"`
	Result Result `json:"result"`
}





