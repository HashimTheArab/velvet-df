package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity/damage"
	"github.com/df-mc/dragonfly/server/player"
	"velvet/session"
	"velvet/utils"
)

type Kill struct {
	Target []cmd.Target `optional:"" name:"victim"`
}

func (t Kill) Run(source cmd.Source, _ *cmd.Output) {
	p := source.(*player.Player)
	if len(t.Target) > 0 {
		if len(t.Target) > 1 {
			if session.Get(p).XUID != utils.Config.Staff.Owner.XUID {
				p.Message(NoPermission)
				return
			}
			p.Messagef("§cYou have killed §d%v §cpeople.", len(t.Target))
			return
		}
		if tg, ok := t.Target[0].(*player.Player); ok {
			tg.Hurt(tg.MaxHealth(), damage.SourceCustom{})
			p.Messagef("§cYou have killed %v.", tg.Name())
		}
		return
	}
	p.Hurt(p.MaxHealth(), damage.SourceCustom{})
	p.Message("§cYou have killed yourself.")
}

func (Kill) Allow(s cmd.Source) bool { return checkAdmin(s) }
