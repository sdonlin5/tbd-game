package game

import "github.com/google/uuid"

type Reply interface {}

type Result struct {
	Shot 	Shot	`json:"shot"`
	Hit 	bool	`json:"hit"`
	Sink 	bool	`json:"sink"`
}

type QuitConfirmed struct {
	PlayerConfirmed bool
	
}

type Response struct {
	ReceiverID 	uuid.UUID	`json:"receiverID"`
	Reply		Reply		`json:"result"`
}