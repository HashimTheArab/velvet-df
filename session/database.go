package session

import (
	"velvet/db"
)

// Load loads the sessions' data from the database.
func (s *Session) Load() {
	go func() {
		s.Stats = db.Stats(s.Player.XUID())
	}()
}

// Save saves the sessions data to the database.
func (s *Session) Save() {
	db.SaveStats(s.Player.XUID(), s.Stats)
}
