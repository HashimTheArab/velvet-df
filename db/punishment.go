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
	err := db.QueryRowx("SELECT * FROM Bans WHERE XUID=? OR IGN=?", id, id).StructScan(&entry)
	if err != nil {
		return nil
	}
	if !entry.Blacklist() && time.Now().Unix() >= entry.Expires {
		UnbanPlayer(id)
		return nil
	}
	return &entry
}

// BanPlayer bans a player and handles everything such as the disconnection, broadcasting, and webhook.
func BanPlayer(target, mod, reason string, length time.Duration) {
	p, ok := utils.Srv.PlayerByName(target)
	blacklist := length == -1
	lengthString := utils.DurationToString(length)
	var xuid string
	if ok {
		target = p.Name()
		xuid = p.XUID()
		if blacklist {
			p.Disconnect(utils.Config.Ban.BlacklistScreen)
		} else {
			p.Disconnect(fmt.Sprintf(utils.Config.Ban.Screen, reason, lengthString))
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
		_, _ = db.Exec("INSERT INTO Bans(XUID, IGN, Mod, Reason, Expires) VALUES(?, ?, ?, ?, ?)", xuid, target, mod, reason, expires)
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

// IsDeviceBanned returns whether a player has a banned account on their current device id.
func IsDeviceBanned(deviceID string) bool { // todo
	var xuid string
	var ign string
	_ = db.QueryRowx("SELECT XUID, IGN FROM Players WHERE DeviceID=?", deviceID).Scan(&xuid, &ign)
	return true
}
