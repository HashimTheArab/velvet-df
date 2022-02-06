package session

import (
	"velvet/db"
	"velvet/perm"
)

// Load loads the sessions' data from the database.
func (s *Session) Load() {
	data := db.GetData(s.Player.XUID())
	s.Stats.Kills = data.Kills
	s.Stats.Deaths = data.Deaths
	s.SetRank(perm.GetRank(data.Rank))
	s.SetPerms(data.Perms)
}

// Save saves the sessions data to the database.
func (s *Session) Save() {
	data := db.PlayerData{
		PlayerStats: s.Stats,
		Perms:       s.perms.Load(),
	}
	if s.Rank() != nil {
		data.Rank = s.Rank().Name
	}
	db.SaveData(s.Player.XUID(), &data)
}
