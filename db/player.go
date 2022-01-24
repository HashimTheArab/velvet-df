package db

// Register registers a player into the database.
func Register(xuid, ign, deviceId string) {
	_, _ = db.Exec("REPLACE INTO Players (XUID, IGN, DeviceID) VALUES (?, ?, ?)", xuid, ign, deviceId)
}
