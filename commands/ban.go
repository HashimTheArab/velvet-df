package commands

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"time"
	"velvet/console"
	"velvet/db"
	"velvet/discord/webhook"
	"velvet/session"
	"velvet/utils"
)

type Ban struct {
	Player []cmd.Target `cmd:"victim"`
	Length string       `cmd:"length"`
	Reason cmd.Varargs  `cmd:"reason"`
}

type BanOffline struct {
	Player string      `cmd:"victim"`
	Length string      `cmd:"length"`
	Reason cmd.Varargs `cmd:"reason"`
}

type BanLift struct {
	Player string `cmd:"target"`
}

type BanInfo struct {
	Player string `cmd:"target"`
}

func (t Ban) Run(source cmd.Source, output *cmd.Output) {
	if len(t.Player) > 1 {
		output.Print(utils.Config.Ban.CanOnlyBanOne)
		return
	}
	if target, ok := t.Player[0].(*player.Player); ok {
		if _, ok := source.(*console.CommandSender); !ok {
			if target == source || (source.Name() != utils.Config.Staff.Owner.Name && session.Get(target).HasFlag(session.FlagStaff)) {
				output.Print(utils.Config.Message.CannotPunishPlayer)
				return
			}
		}
		if t.Reason == "" {
			output.Print("§cProvide a reason.")
			return
		}
		duration, err := utils.ParseDuration(t.Length)
		if err != nil {
			output.Print(utils.Config.Message.InvalidPunishmentTime)
			return
		}
		ban(source, output, target.Name(), string(t.Reason), duration)
	}
}

func (t BanOffline) Run(source cmd.Source, output *cmd.Output) {
	p := source.(*player.Player)
	if _, ok := source.(*console.CommandSender); !ok {
		if t.Player == source.Name() || (db.IsStaff(t.Player) && p.XUID() != utils.Config.Staff.Owner.XUID) {
			output.Print(utils.Config.Message.CannotPunishPlayer)
			return
		}
	}
	if t.Reason == "" {
		output.Print("§cProvide a reason.")
		return
	}
	duration, err := utils.ParseDuration(t.Length)
	if err != nil {
		output.Print(utils.Config.Message.InvalidPunishmentTime)
		return
	}
	ban(source, output, t.Player, string(t.Reason), duration)
}

func (t BanLift) Run(_ cmd.Source, output *cmd.Output) {
	p, ban, ok := db.GetBan(t.Player)
	if !ok {
		output.Error(utils.Config.Ban.PlayerNotBanned)
		return
	}
	db.UnbanPlayer(t.Player)
	output.Printf(utils.Config.Ban.PlayerUnbanned, p.DisplayName)
	webhook.Send(utils.Config.Discord.Webhook.UnbanLogger, webhook.Message{
		Embeds: []webhook.Embed{{
			Title:       "Player Pardoned",
			Description: fmt.Sprintf("**Player:** %v\n**Staff:** %v\n", p.DisplayName, ban.Mod),
			Color:       0xFF8900,
		}},
	})
}

func (t BanInfo) Run(_ cmd.Source, output *cmd.Output) {
	p, ban, ok := db.GetBan(t.Player)
	if !ok {
		output.Error(utils.Config.Ban.PlayerNotBanned)
		return
	}
	output.Printf(utils.Config.Ban.Info, p.DisplayName, ban.Mod, ban.Reason, ban.FormattedExpiration())
}

func ban(p cmd.Source, output *cmd.Output, target, reason string, length time.Duration) {
	if _, _, ok := db.GetBan(target); ok {
		output.Print(utils.Config.Ban.PlayerAlreadyBanned)
		return
	}
	if reason == "" {
		output.Print(utils.Config.Message.SpecifyReason)
		return
	}
	db.BanPlayer(target, p.Name(), reason, length)
}

func (Ban) Allow(s cmd.Source) bool        { return checkStaff(s) || checkConsole(s) }
func (BanOffline) Allow(s cmd.Source) bool { return checkStaff(s) || checkConsole(s) }
func (BanLift) Allow(s cmd.Source) bool    { return checkStaff(s) || checkConsole(s) }
func (BanInfo) Allow(s cmd.Source) bool    { return checkStaff(s) || checkConsole(s) }
