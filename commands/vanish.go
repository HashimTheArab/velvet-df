package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"velvet/session"
)

type Vanish struct {
	// Targets []cmd.Target `name:"target" optional:""` todo
}

func (t Vanish) Run(source cmd.Source, output *cmd.Output) {
	p, ok := source.(*player.Player)
	if !ok {
		return
	}
	session.Get(p).Vanish()
}

func (Vanish) Allow(s cmd.Source) bool { return checkStaff(s) }
