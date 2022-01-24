package dfutils

import (
	"fmt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"net"
	"strconv"
	"velvet/db"
	"velvet/discord/webhook"
	"velvet/session"
	"velvet/utils"
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
	"1916611344": protocol.DeviceOS(14), // windows phone
}

/*
   1651113805
   1909043648
   1835298427
*/

func (allower) Allow(_ net.Addr, d login.IdentityData, c login.ClientData) (string, bool) {
	db.Register(d.XUID, d.DisplayName, c.DeviceID)
	ban := db.GetBan(d.XUID)
	if ban != nil {
		if ban.XUID == "" {
			ban.Update(d.XUID)
		}
		if ban.Blacklist() {
			log.Infof("%v tried joining but is BLACKLISTED.", d.DisplayName)
			return fmt.Sprintf(utils.Config.Ban.BlacklistScreen, ""), false
		} else {
			log.Infof("%v tried joining but is banned.", d.DisplayName)
			return fmt.Sprintf(utils.Config.Ban.LoginScreen, ""), false
		}
	}
	_, ok := titleIds[d.TitleID]
	if !ok {
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
	session.New(d.DisplayName)
	return "", true
}
