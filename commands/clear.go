package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

type Clear struct {
	Targets []cmd.Target `optional:"" name:"victim"`
	Armour  bool         `optional:"" name:"armor""`
}

func (t Clear) Run(source cmd.Source, output *cmd.Output) {
	p, ok := source.(*player.Player)
	if len(t.Targets) > 0 {
		if len(t.Targets) > 1 {
			output.Print("§cYou can only clear the inventory of one player at a time.")
			return
		}
		tg, ok := t.Targets[0].(*player.Player)
		if ok {
			tg.Inventory().Clear()
			if t.Armour {
				tg.Armour().Clear()
			}
		} else {
			output.Print("§cUnknown Player.")
		}
		return
	}
	if ok {
		p.Inventory().Clear()
		if t.Armour {
			p.Armour().Clear()
		}
	}
}

func (Clear) Allow(s cmd.Source) bool { return checkAdmin(s) }
