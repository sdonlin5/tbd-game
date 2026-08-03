// Defines player offensive move struct and methods
// --
package game

// Implements Move interface
func (s Shot) Play() string { return "Shot" }

// Implements Result interface
func (res ShotResult) Played() string { return "ShotResult" }


// Shot - Holds X,Y coordinates of a player shot input
type Shot struct {
	X uint8 `json:"x"`
	Y uint8 `json:"y"`
}


// Validates shots and detects hits
// defender: player being shot at
func (s *Shot) PlayShot(defender *Player) *ShotResult {
	// Defaults to invalid shot
	res := &ShotResult{Type: "ShotResult", Shot: s, Valid: false, Hit: false}

	// If shot is valid, detect hit (and sink)
	if s.ValidateShot(defender) {
		res.Valid = true
		res.Hit = s.DetectHit(defender)
	}
	return res
}

// Validates the shot values from the player
func (s Shot) ValidateShot(defender *Player) bool {
	if int(s.X) >= len(defender.PlayerBoard.Tiles) ||
		int(s.Y) >= len(defender.PlayerBoard.Tiles[0]) {
		return false
	}
	return true
}

// Temp function to test data transit and result creation
func (s *Shot) DetectHit(p *Player) bool {
	// Checks for hit on player waiting
	// TODO: Add real hit detection
	if s.X%2 == 0 || s.Y%2 == 0 {
		return true
	} else {
		return false
	}
}

// Transport for the result of a shot.
type ShotResult struct {
	Type  string
	Shot  *Shot
	Valid bool
	Hit   bool
	// Sink  bool
}
