package db

type PlayerStats struct {
	Kills  uint
	Deaths uint
}

type PlayerData struct {
	PlayerStats
	Banned bool
}
