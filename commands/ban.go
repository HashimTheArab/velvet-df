package commands

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"time"
	"velvet/db"
	"velvet/discord/webhook"
	"velvet/session"
	"velvet/utils"
)

type Ban struct {
	Player []cmd.Target `name:"victim"`
	Length string       `name:"length"`
	Reason cmd.Varargs  `name:"reason"`
}

type BanOffline struct {
	Player string      `name:"victim"`
	Length string      `name:"length"`
	Reason cmd.Varargs `name:"reason"`
}

type BanLift struct {
	Player string `name:"target"`
}

type BanInfo struct {
	Player string `name:"target"`
}

func (t Ban) Run(source cmd.Source, output *cmd.Output) {
	p := source.(*player.Player)

	if len(t.Player) > 1 {
		output.Print(utils.Config.Ban.CanOnlyBanOne)
		return
	}

	if target, ok := t.Player[0].(*player.Player); ok {
		if target.Name() == source.Name() || (source.Name() != utils.Config.Staff.Owner.Name && session.Get(target).HasFlag(session.Staff)) {
			output.Print(utils.Config.Message.CannotPunishPlayer)
			return
		}
		if t.Reason == "" {
			output.Print("§cProvide a reason.")
			return
		}
		duration := utils.DurationFromString(t.Length)
		if duration == -1 {
			output.Print(utils.Config.Message.InvalidPunishmentTime)
			return
		}
		ban(p, target.Name(), string(t.Reason), duration)
	}
}

func (t BanOffline) Run(source cmd.Source, output *cmd.Output) {
	p := source.(*player.Player)

	_, mod := utils.Config.Staff.Mods[t.Player]
	_, admin := utils.Config.Staff.Admins[t.Player]
	if t.Player == source.Name() || ((mod || admin) && p.XUID() != utils.Config.Staff.Owner.XUID) {
		output.Print(utils.Config.Message.CannotPunishPlayer)
		return
	}
	if t.Reason == "" {
		output.Print("§cProvide a reason.")
		return
	}
	duration := utils.DurationFromString(t.Length)
	if duration == -1 {
		output.Print(utils.Config.Message.InvalidPunishmentTime)
		return
	}
	ban(p, t.Player, string(t.Reason), duration)
}

func (t BanLift) Run(_ cmd.Source, output *cmd.Output) {
	ban := db.GetBan(t.Player)
	if ban != nil {
		db.UnbanPlayer(t.Player)
		output.Printf(utils.Config.Ban.PlayerUnbanned, ban.IGN)
		webhook.Send(utils.Config.Discord.Webhook.UnbanLogger, webhook.Message{
			Embeds: []webhook.Embed{{
				Title:       "Player Pardoned",
				Description: fmt.Sprintf("**Player:** %v\n**Staff:** %v\n", ban.IGN, ban.Mod),
				Color:       0xFF8900,
			}},
		})
	} else {
		output.Print(utils.Config.Ban.PlayerNotBanned)
	}
}

func (t BanInfo) Run(_ cmd.Source, output *cmd.Output) {
	ban := db.GetBan(t.Player)
	if ban != nil {
		output.Printf(utils.Config.Ban.Info, ban.IGN, ban.Mod, ban.Reason, ban.FormattedExpiration())
	} else {
		output.Print(utils.Config.Ban.PlayerNotBanned)
	}
}

func ban(p *player.Player, target, reason string, length time.Duration) {
	if db.GetBan(target) != nil {
		p.Message(utils.Config.Ban.PlayerAlreadyBanned)
		return
	}
	db.BanPlayer(target, p.Name(), reason, length)
}

func (Ban) Allow(s cmd.Source) bool        { return checkStaff(s) }
func (BanOffline) Allow(s cmd.Source) bool { return checkStaff(s) }
func (BanLift) Allow(s cmd.Source) bool    { return checkStaff(s) }
func (BanInfo) Allow(s cmd.Source) bool    { return checkStaff(s) }
