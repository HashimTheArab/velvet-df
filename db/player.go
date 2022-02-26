package db

import (
	"velvet/perm"
)

// Register registers a player into the database.
func Register(xuid, ign, deviceId string) {
	_, _ = db.Exec("INSERT INTO Players(XUID, IGN, DeviceID) VALUES (?, ?, ?) ON CONFLICT DO UPDATE SET IGN=?, DeviceID=?", xuid, ign, deviceId, ign, deviceId)
}

// Registered returns whether a player is registered
func Registered(id string) bool {
	var r bool
	_ = db.QueryRowx("SELECT EXISTS(SELECT IGN FROM Players WHERE IGN=? OR XUID=?)", id, id).Scan(&r)
	return r
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

// SetRank will set the rank for a player.
func SetRank(id, rank string) {
	_, _ = db.Exec("UPDATE Players set PlayerRank=? WHERE IGN=? OR XUID=?", rank, id, id)
}

// GetRank will return the rank of a player or nil.
func GetRank(id string) *perm.Rank {
	var rank string
	_ = db.QueryRow("SELECT PlayerRank FROM Players WHERE IGN=? OR XUID=?", id, id).Scan(&rank)
	return perm.GetRank(rank)
}

// HasRank will return true if the given player has the given rank.
func HasRank(id string, rank string) bool {
	var found string
	_ = db.QueryRow("SELECT PlayerRank FROM Players WHERE IGN=? OR XUID=?", id, id).Scan(&found)
	return found == rank
}

// IsStaff will return whether a player has a staff rank.
func IsStaff(id string) bool {
	var found string
	_ = db.QueryRow("SELECT PlayerRank FROM Players WHERE IGN=? OR XUID=?", id, id).Scan(&found)
	if found == "" {
		return false
	}
	return perm.StaffRanks.Contains(found)
}
