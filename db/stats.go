package db

import (
	"sync"
)

type PlayerStats struct {
	Kills  uint32 `db:"Kills"`
	Deaths uint32 `db:"Deaths"`
	Mutex  sync.Mutex
}

type PlayerData struct {
	*PlayerStats
	Rank  string `db:"PlayerRank"`
	Perms uint32 `db:"Perms"`
}

func (s *PlayerStats) GetKills() uint32 {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	return s.Kills
}

func (s *PlayerStats) GetDeaths() uint32 {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	return s.Deaths
}

func (s *PlayerStats) AddKills(k uint32) {
	s.Mutex.Lock()
	s.Kills += k
	s.Mutex.Unlock()
}

func (s *PlayerStats) AddDeaths(k uint32) {
	s.Mutex.Lock()
	s.Deaths += k
	s.Mutex.Unlock()
}

// GetData will get the data of a player.
func GetData(id string) PlayerData {
	var d = PlayerData{PlayerStats: &PlayerStats{}}
	_ = db.QueryRowx("SELECT Kills, Deaths, PlayerRank, Perms FROM Players WHERE XUID=? OR IGN=?", id, id).StructScan(&d)
	return d
}

// SaveData will save the data of a player.
func SaveData(xuid string, d *PlayerData) {
	_, _ = db.Exec("UPDATE Players set Kills=?, Deaths=?, PlayerRank=?, Perms=? WHERE XUID=?", d.Kills, d.Deaths, d.Rank, d.Perms, xuid)
}

// Stats fetches the stats of a player. id can be either the name or the XUID of the player.
func Stats(id string) *PlayerStats {
	var stats PlayerStats
	_ = db.QueryRowx("SELECT Kills, Deaths FROM Players WHERE XUID=? OR IGN=?", id, id).StructScan(&stats)
	return &stats
}

// SaveStats saves the stats of a player. This is done on another thread.
func SaveStats(xuid string, s *PlayerStats) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	_, _ = db.Exec("UPDATE Players set Kills=?, Deaths=? WHERE XUID=?", s.Kills, s.Deaths, xuid)
}
