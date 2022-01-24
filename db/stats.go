package db

import "go.uber.org/atomic"

type PlayerStats struct {
	Kills  atomic.Uint32 `db:"Kills"`
	Deaths atomic.Uint32 `db:"Deaths"`
}

// Stats fetches the stats of a player. id can be either the name or the XUID of the player.
func Stats(id string) PlayerStats {
	var stats PlayerStats
	_ = db.QueryRow("SELECT Kills, Deaths FROM Players WHERE XUID=? OR IGN=?", id, id).Scan(&stats)
	return stats
}

// SaveStats saves the stats of a player. This is done on another thread.
func SaveStats(xuid string, s PlayerStats) {
	go func() {
		_, _ = db.Exec("UPDATE Players set Kills=?, Deaths=? WHERE XUID=?", s.Kills.Load(), s.Deaths.Load(), xuid)
	}()
}
