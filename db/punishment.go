package db

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/player/chat"
	"time"
	"velvet/discord/webhook"
	"velvet/utils"
)

// GetBan returns the ban of a player, or nil if the player is not banned.
func GetBan(id string) *BanEntry {
	var entry BanEntry
	if err := db.QueryRowx("SELECT * FROM Bans WHERE XUID=? OR IGN=?", id, id).StructScan(&entry); err != nil {
		return nil
	}
	if !entry.Blacklist() && time.Now().Unix() >= entry.Expires {
		UnbanPlayer(id)
		return nil
	}
	return &entry
}

// BanPlayer bans a player and handles everything such as the disconnection, broadcasting, and webhook.
func BanPlayer(target, targetXUID, mod, reason string, length time.Duration) {
	p, ok := utils.Srv.PlayerByName(target)
	blacklist := length == -1
	lengthString := utils.DurationToString(length)
	if ok {
		target = p.Name()
		if blacklist {
			p.Disconnect(utils.Config.Ban.BlacklistScreen)
		} else {
			p.Disconnect(fmt.Sprintf(utils.Config.Ban.Screen, mod, reason, lengthString))
		}
	}
	if blacklist {
		_, _ = fmt.Fprintf(chat.Global, utils.Config.Ban.BlacklistBroadcast, target, mod, reason)
	} else {
		_, _ = fmt.Fprintf(chat.Global, utils.Config.Ban.Broadcast, target, mod, lengthString, reason)
	}
	go func() {
		var expires int64 = -1
		if length != 0 {
			expires = time.Now().Add(length).Unix()
		}
		_, _ = db.Exec("INSERT INTO Bans(XUID, IGN, Mod, Reason, Expires) VALUES(?, ?, ?, ?, ?)", targetXUID, target, mod, reason, expires)
		var msg webhook.Message
		if blacklist {
			msg = webhook.Message{
				Embeds: []webhook.Embed{{
					Title:       "Player Blacklisted",
					Description: fmt.Sprintf("**Player:** %v\n**Staff:** %v\n**Reason:** %v", target, mod, reason),
					Color:       0xc000ff,
				}},
			}
		} else {
			msg = webhook.Message{
				Embeds: []webhook.Embed{{
					Title:       "Player Banned",
					Description: fmt.Sprintf("**Player:** %v\n**Staff:** %v\n**Reason:** %v\n**Length:** %v", target, mod, reason, lengthString),
					Color:       0xc000ff,
				}},
			}
		}
		webhook.Send(utils.Config.Discord.Webhook.BanLogger, msg)
	}()
}

// UnbanPlayer unbans a player.
func UnbanPlayer(id string) {
	_, _ = db.Exec("DELETE FROM Bans WHERE XUID=? OR IGN=?", id, id)
}

// DeviceBan will return an existing ban with the given device id or nil.
func DeviceBan(deviceID string) *BanEntry {
	if rows, err := db.Query("SELECT XUID, IGN FROM Players WHERE DeviceID=?", deviceID); err == nil { // get all players with that device id
		var xuids, names []string
		for rows.Next() {
			var xuid, name string
			if err := rows.Scan(&xuid, &name); err == nil {
				xuids, names = append(xuids, xuid), append(names, name)
			}
		}
		var entry *BanEntry
		if err := db.QueryRowx("SELECT * FROM Bans WHERE XUID IN (?) OR IGN IN (?)", xuids, names).StructScan(entry); err != nil { // get the first ban with any of the names or xuids
			return nil
		}
		if !entry.Blacklist() && time.Now().Unix() >= entry.Expires { // unban if the ban has expired
			UnbanPlayer(entry.IGN)
			UnbanPlayer(entry.XUID)
			return nil
		}
		return entry
	}
	return nil
}
