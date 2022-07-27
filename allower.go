package main

import (
	"fmt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"net"
	"strconv"
	"velvet/db"
	"velvet/discord/webhook"
	"velvet/utils"
	"velvet/utils/whitelist"
)

type allower struct{}

var titleIds = map[string]protocol.DeviceOS{
	"1739947436": protocol.DeviceAndroid,
	"1810924247": protocol.DeviceIOS,
	"1944307183": protocol.DeviceFireOS,
	"896928775":  protocol.DeviceWin10,
	"2044456598": protocol.DeviceOrbis,
	"2047319603": protocol.DeviceNX,
	"1828326430": protocol.DeviceXBOX,
	"1916611344": protocol.DeviceWP,
}

/*
   1651113805
   1909043648
   1835298427
*/

func (allower) Allow(_ net.Addr, d login.IdentityData, c login.ClientData) (string, bool) {
	db.Register(d.XUID, d.DisplayName, c.DeviceID)
	if user, ban, ok := db.GetBan(d.DisplayName); ok {
		if user.XUID == "" {
			ban.Update(user.Name, d.XUID)
		}
		if ban.Permanent {
			logger.Infof("%v tried joining but is BLACKLISTED.", d.DisplayName)
			return fmt.Sprintf(utils.Config.Ban.BlacklistScreen, ""), false
		}
		logger.Infof("%v tried joining but is banned.", d.DisplayName)
		return fmt.Sprintf(utils.Config.Ban.LoginScreen, ban.Reason, ban.FormattedExpiration()), false
	}
	if whitelist.Enabled() && !whitelist.Contains(d.DisplayName) {
		logger.Infof("%v tried joining but the server is whitelisted.", d.DisplayName)
		return fmt.Sprintf("§cThis server is whitelisted."), false
	}
	if _, ok := titleIds[d.TitleID]; !ok {
		webhook.Send(utils.Config.Discord.Webhook.TitleIDLogger, webhook.Message{
			Embeds: []webhook.Embed{{
				Fields: []webhook.Field{
					{Name: "Player", Value: d.DisplayName},
					{Name: "Title ID", Value: d.TitleID},
					{Name: "Device OS", Value: strconv.Itoa(int(c.DeviceOS))},
				},
				Color: 0x1F85DE,
			}},
		})
	}
	return "", true
}
