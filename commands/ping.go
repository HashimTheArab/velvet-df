package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

type Ping struct {
	Targets []cmd.Target `optional:"" json:"target"`
}

func (t Ping) Run(source cmd.Source, output *cmd.Output) {
	if len(t.Targets) > 1 {
		output.Printf("§cYou can only specify one player at once.")
		return
	}
	if len(t.Targets) > 0 {
		if p, ok := t.Targets[0].(*player.Player); ok {
			output.Printf("§e%v's §dping is §e%v.", p.Name(), p.Latency().String())
		} else {
			output.Printf(PlayerNotFound)
		}
	} else {
		if p, ok := source.(*player.Player); ok {
			output.Printf("§dYour ping is §e%v.", p.Latency().String())
		} else {
			c, _ := cmd.ByAlias("ping")
			output.Printf("§cUsage: §7" + c.Usage())
		}
	}
}
