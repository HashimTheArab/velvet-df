package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

type Clear struct {
	Targets cmd.Optional[[]cmd.Target] `cmd:"victim"`
	Armour  cmd.Optional[bool]         `cmd:"armor"`
}

func (t Clear) Run(source cmd.Source, output *cmd.Output) {
	p, ok := source.(*player.Player)
	armour := t.Armour.LoadOr(true)
	if targets := t.Targets.LoadOr(nil); len(targets) > 0 {
		if len(targets) > 1 {
			output.Print("§cYou can only clear the inventory of one player at a time.")
			return
		}
		tg, ok := targets[0].(*player.Player)
		if ok {
			tg.Inventory().Clear()
			if armour {
				tg.Armour().Clear()
			}
		} else {
			output.Print("§cUnknown Player.")
		}
		return
	}
	if ok {
		p.Inventory().Clear()
		if armour {
			p.Armour().Clear()
		}
	}
}

func (Clear) Allow(s cmd.Source) bool { return checkAdmin(s) }
