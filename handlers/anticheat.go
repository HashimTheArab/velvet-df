package handlers

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/oomph-ac/oomph/check"
	"github.com/oomph-ac/oomph/player"
	"time"
	"velvet/discord/webhook"
	"velvet/session"
	"velvet/utils"
)

// AntiCheatHandler is the handler for Oomph.
type AntiCheatHandler struct {
	player.NopHandler
	p *player.Player
}

// NewACHandler creates a new Oomph handler.
func NewACHandler(p *player.Player) *AntiCheatHandler {
	return &AntiCheatHandler{p: p}
}

// HandlePunishment ...
func (a AntiCheatHandler) HandlePunishment(ctx *event.Context, ch check.Check, _ *string) {
	ctx.Cancel()

	_, autoClickerC := ch.(*check.AutoClickerC)
	_, autoClickerD := ch.(*check.AutoClickerD)
	_, timer := ch.(*check.TimerA)
	if autoClickerC || autoClickerD || timer {
		// These checks are not always accurate, and shouldn't be punished for.
		return
	}

	pl, ok := utils.Srv.PlayerByName(a.p.Name())
	if !ok {
		return
	}
	if s := session.Get(pl); s == nil || s.Staff() {
		return
	}

	name, sub := ch.Name()
	reason := fmt.Sprintf("%s (%s)", name, sub)
	punishmentString := "Kick"
	//if checkBan(ch) {
	//	punishmentString = "Ban"
	//	db.BanPlayer(pl.Name(), "Oomph", reason, time.Hour*24*14)
	//} else {
	//	_, _ = fmt.Fprintf(chat.Global, utils.Config.Kick.Broadcast, pl.Name(), "Oomph", reason)
	//	pl.Disconnect(fmt.Sprintf("§6[§bOomph§6] Caught yo ass lackin!\n§6Reason: §b%v", reason))
	//}
	_, _ = fmt.Fprintf(chat.Global, utils.Config.Kick.Broadcast, pl.Name(), "Oomph", reason)
	pl.Disconnect(fmt.Sprintf("§6[§bOomph§6] Caught yo ass lackin!\n§6Reason: §b%v", reason))

	webhook.Send(utils.Config.Discord.Webhook.AntiCheatLogger, webhook.Message{
		Embeds: []webhook.Embed{{
			Title: "**Oomph Punishment**",
			Description: fmt.Sprintf(
				"Player: %v\nPunishment: %v\nCheck:%v\nViolations: %v",
				pl.Name(), punishmentString, reason, mgl64.Round(ch.Violations(), 2),
			),
			Color:  0xFF009F,
			Footer: webhook.Footer{Text: time.Now().Format("01/02/06 @ 03:04:05 PM")},
		}},
	})
}

// HandleFlag ...
func (a AntiCheatHandler) HandleFlag(_ *event.Context, c check.Check, params map[string]any, log *bool) {
	*log = false
	name, sub := c.Name()
	//_, _ = fmt.Fprintf(chat.Global, "§7[§cOomph§7] §b%v §6flagged §b%v (%v) §6(§cx%v§6) %v", a.p.Name(), name, sub, mgl64.Round(c.Violations(), 2), utils.PrettyParams(params))
	session.AllStaff().Messagef("§7[§cOomph§7] §b%v §6flagged §b%v (%v) §6(§cx%v§6) %v", a.p.Name(), name, sub, mgl64.Round(c.Violations(), 2), utils.PrettyParams(params))
}

// HandleDebug ...
func (a AntiCheatHandler) HandleDebug(ctx *event.Context, _ check.Check, _ map[string]any) {
	ctx.Cancel()
}

// checkban checks if a check should ban.
func checkBan(ch check.Check) (ban bool) {
	ban = true
	switch ch.(type) {
	case *check.AutoClickerA, *check.AutoClickerB, *check.AutoClickerC, *check.AutoClickerD:
		ban = false
	case *check.OSSpoofer:
		ban = false
	case *check.TimerA:
		ban = false
	}
	return
}
