package db

// Register registers a player into the database.
func Register(xuid, ign, deviceId string) {
	_, _ = db.Exec("REPLACE INTO Players (XUID, IGN, DeviceID) VALUES (?, ?, ?)", xuid, ign, deviceId)
}

// GetDeviceID will return the device id for the given ign or an empty string if that player has never joined before.
func GetDeviceID(ign string) string {
	var deviceID string
	_ = db.QueryRowx("SELECT DeviceID FROM Players WHERE IGN=?", ign).Scan(&deviceID)
	return deviceID
}

// GetAlias will return all the names that have the same deviceID as the given ign.
// Zero values will be returned if the player has never joined before.
func GetAlias(ign string) (deviceID string, names []string) {
	if deviceID = GetDeviceID(ign); deviceID != "" {
		if rows, err := db.Query("SELECT IGN FROM Players WHERE DeviceID=?", deviceID); err == nil { // get all players with the deviceID
			for rows.Next() {
				var name string
				if err := rows.Scan(&name); err == nil {
					names = append(names, name)
				}
			}
			if len(names) > 0 {
				var bans []string
				if rows, err := db.Query("SELECT IGN FROM Bans WHERE IGN IN (?)", names); err == nil { // get all bans from the names
					for rows.Next() {
						var ban string
						if err := rows.Scan(&ban); err == nil {
							bans = append(bans, ban)
						}
					}
				}
				for k, v := range names {
					for _, bannedPlayer := range bans {
						if bannedPlayer == v {
							names[k] = v + " §l§cBANNED§r"
						}
					}
					names[k] = "§e" + v
				}
			}
		}
	}
	return
}
