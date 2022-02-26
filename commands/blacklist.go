package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"velvet/console"
	"velvet/db"
	"velvet/session"
	"velvet/utils"
)

type Blacklist struct {
	Player []cmd.Target `name:"victim"`
	Reason cmd.Varargs  `name:"reason"`
}

type BlacklistOffline struct {
	Player string      `name:"victim"`
	Reason cmd.Varargs `name:"reason"`
}

func (t Blacklist) Run(source cmd.Source, output *cmd.Output) {
	if len(t.Player) > 1 {
		output.Print(utils.Config.Ban.CanOnlyBanOne)
		return
	}
	if target, ok := t.Player[0].(*player.Player); ok {
		if _, ok := source.(*console.CommandSender); !ok {
			if target.Name() == source.Name() || (source.Name() != utils.Config.Staff.Owner.Name && session.Get(target).HasFlag(session.FlagStaff)) {
				output.Print(utils.Config.Message.CannotPunishPlayer)
				return
			}
		}
		ban(source, output, target.Name(), string(t.Reason), -1)
	}
}

func (t BlacklistOffline) Run(source cmd.Source, output *cmd.Output) {
	p, _ := source.(*player.Player)
	if _, ok := source.(*console.CommandSender); !ok {
		if t.Player == source.Name() || (db.IsStaff(t.Player) && p.XUID() != utils.Config.Staff.Owner.XUID) {
			output.Print(utils.Config.Message.CannotPunishPlayer)
			return
		}
	}
	ban(source, output, t.Player, string(t.Reason), -1)
}

func (Blacklist) Allow(s cmd.Source) bool        { return checkStaff(s) || checkConsole(s) }
func (BlacklistOffline) Allow(s cmd.Source) bool { return checkStaff(s) || checkConsole(s) }
