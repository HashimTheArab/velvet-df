package handlers

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/justtaldevelops/oomph/check"
	"github.com/justtaldevelops/oomph/player"
	"github.com/justtaldevelops/oomph/utils"
	"time"
	"velvet/db"
	"velvet/discord/webhook"
	"velvet/session"
	vu "velvet/utils"
)

type AntiCheatHandler struct {
	player.NopHandler
	p *player.Player
}

func NewACHandler(p *player.Player) *AntiCheatHandler {
	return &AntiCheatHandler{p: p}
}

func (a AntiCheatHandler) HandlePunishment(ctx *event.Context, c check.Check, m *string) {
	if pl, ok := vu.Srv.PlayerByName(a.p.Name()); ok {
		if session.Get(pl).Staff() {
			return
		}
		punishmentString := "Kick"
		name, sub := c.Name()
		reason := name + "(" + sub + ")"
		playerName := pl.Name()
		//if c.BaseSettings().Punishment == punishment.Ban() {
		//	pl.Disconnect(fmt.Sprintf("§6[§bOomph§6] Caught yo ass lackin!\n§6Reason: §b%v", reason))
		//	db.BanPlayer(pl.Name(), "Oomph", reason, time.Hour*24*14)
		//} else if c.BaseSettings().Punishment == punishment.Kick() {
		//	_, _ = fmt.Fprintf(chat.Global, vu.Config.Kick.Broadcast, pl.Name(), "Oomph", reason)
		//	pl.Disconnect(fmt.Sprintf("§6[§bOomph§6] Caught yo ass lackin!\n§6Reason: §b%v", reason))
		//} else {
		//	return
		//}
		*m = fmt.Sprintf("§6[§bOomph§6] Caught yo ass lackin!\n§6Reason: §b%v", reason)
		ctx.After(func(cancelled bool) {
			if cancelled {
				return
			}
			db.BanPlayer(playerName, "Oomph", reason, time.Hour*24*14)
			webhook.Send(vu.Config.Discord.Webhook.AntiCheatLogger, webhook.Message{
				Embeds: []webhook.Embed{{
					Title:       "**Oomph Punishment**",
					Description: fmt.Sprintf("Player: %v\nPunishment: %v\nCheck:%v\nViolations: %v", playerName, punishmentString, reason, c.Violations()),
					Color:       0xFF009F,
					Footer:      webhook.Footer{Text: time.Now().Format("01/02/06 @ 03:04:05 PM")},
				}},
			})
		})
	}
}

func (a AntiCheatHandler) HandleFlag(ctx *event.Context, c check.Check, params map[string]interface{}) {
	ctx.Cancel()
	name, sub := c.Name()
	session.AllStaff().Messagef("§7[§cOomph§7] §b%v §6flagged §b%v (%v) §6(§cx%v§6) %v", a.p.Name(), name, sub, c.Violations(), utils.PrettyParameters(params))
}

func (a AntiCheatHandler) HandleDebug(ctx *event.Context, _ check.Check, _ map[string]interface{}) {
	ctx.Cancel()
}
