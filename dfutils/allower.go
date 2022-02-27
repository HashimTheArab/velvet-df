package dfutils

import (
	"fmt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"net"
	"strconv"
	"time"
	"velvet/db"
	"velvet/discord/webhook"
	"velvet/session"
	"velvet/utils"
)

type allower struct{}
type oomphConnectionHandler struct{}

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
	if ban := db.GetBan(d.DisplayName); ban != nil {
		if ban.XUID == "" {
			ban.Update(d.XUID)
		}
		if ban.Blacklist() {
			log.Infof("%v tried joining but is BLACKLISTED.", d.DisplayName)
			return fmt.Sprintf(utils.Config.Ban.BlacklistScreen, ""), false
		}
		log.Infof("%v tried joining but is banned.", d.DisplayName)
		return fmt.Sprintf(utils.Config.Ban.LoginScreen, ban.Reason, ban.FormattedExpiration()), false
	}
	if ban := db.DeviceBan(c.DeviceID); ban != nil {
		if ban.Blacklist() {
			log.Infof("%v tried joining but is BLACKLISTED on another account.", d.DisplayName)
			return fmt.Sprintf(utils.Config.Ban.BlacklistScreen, ""), false
		}
		log.Infof("%v tried joining but is banned on another account.", d.DisplayName)
		return fmt.Sprintf(utils.Config.Ban.LoginScreen, ban.Reason, ban.FormattedExpiration()), false
	}
	if utils.Whitelist.Enabled && !utils.Whitelist.Contains(d.DisplayName) {
		log.Infof("%v tried joining but the server is whitelisted.", d.DisplayName)
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
	name := d.DisplayName
	time.AfterFunc(time.Second*35, func() {
		if _, ok := utils.Srv.PlayerByName(name); !ok {
			if s := session.FromName(name); s != nil {
				s.CloseWithoutSaving(name)
			}
		}
	})
	session.New(d.DisplayName)
	return "", true
}
