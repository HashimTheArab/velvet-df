package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"time"
)

type Ping struct {
	Targets cmd.Optional[[]cmd.Target] `cmd:"target"`
}

func (t Ping) Run(source cmd.Source, output *cmd.Output) {
	targets := t.Targets.LoadOr(nil)
	if len(targets) > 1 {
		output.Printf("§cYou can only specify one player at once.")
		return
	}
	if len(targets) > 0 {
		if p, ok := targets[0].(*player.Player); ok {
			output.Printf("§e%v's §dping is §e%v.", p.Name(), p.Latency().Round(time.Millisecond*10).String())
		} else {
			output.Printf(PlayerNotFound)
		}
	} else {
		if p, ok := source.(*player.Player); ok {
			output.Printf("§dYour ping is §e%v.", p.Latency().Round(time.Millisecond*10).String())
		} else {
			c, _ := cmd.ByAlias("ping")
			output.Printf("§cUsage: §7" + c.Usage())
		}
	}
}
