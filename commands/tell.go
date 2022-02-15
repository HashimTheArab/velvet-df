package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

type Tell struct {
	Target  []cmd.Target `name:"player"`
	Message cmd.Varargs  `name:"message"`
}

func (t Tell) Run(source cmd.Source, output *cmd.Output) {
	if len(t.Target) > 1 {
		output.Print("§cYou can only message one player at a time.")
		return
	}
	p, ok := t.Target[0].(*player.Player)
	if !ok {
		output.Printf(PlayerNotFound)
		return
	}
	p.Messagef("§7[§d%v §7-> §dYou§7]: §e%v", source.Name(), string(t.Message))
	output.Printf("§7[§dYou §7-> §d%v§7]: §e%v", p.Name(), string(t.Message))
}
