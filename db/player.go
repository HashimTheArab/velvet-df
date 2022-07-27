package db

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/upper/db/v4"
	"time"
	"velvet/discord/webhook"
	"velvet/utils"
)

var (
	PunishmentTypeBan  = "ban"
	PunishmentTypeMute = "mute"
)

// GetBan gets the ban for a player, if the ban does not exist the boolean will be false.
func GetBan(id string) (*Entry, Punishment, bool) {
	p, err := LoadOfflinePlayer(id)
	if err != nil {
		return nil, Punishment{}, false
	}
	if p.Punishments.Ban.Expired() {
		UnbanPlayer(id)
		return nil, Punishment{}, false
	}
	return p, p.Punishments.Ban, true
}

// BanPlayer bans a player and handles everything such as the disconnection, broadcasting, and webhook.
func BanPlayer(target, mod, reason string, length time.Duration) {
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
	var expires time.Time
	if length != 0 {
		expires = time.Now().Add(length)
	}
	_, _ = sess.Collection("punishments").Insert(Punishment{
		Mod:        mod,
		Reason:     reason,
		Permanent:  expires == time.Time{},
		Expiration: expires,
	})
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
}

// UnbanPlayer unbans a player.
func UnbanPlayer(id string) {
	_ = findBan(id).Delete()
}

// findBan is internally used to find a ban entry.
func findBan(id string) db.Result {
	return sess.Collection("punishments").Find(
		db.And(
			db.Or(db.Cond{"xuid": id}, db.Cond{"ign": id}),
			db.Cond{"type": PunishmentTypeBan},
		),
	)
}
