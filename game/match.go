// Implements the match game thread
package game

import "github.com/google/uuid"

type Match struct {
	PlayerOne uuid.UUID
	PlayerTwo uuid.UUID
	Actions   chan Action
}
