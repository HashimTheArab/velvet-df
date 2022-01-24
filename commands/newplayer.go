package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

type NewPlayer struct {
	Name string `name:"name"`
}

func (t NewPlayer) Run(source cmd.Source, _ *cmd.Output) {
	p, _ := source.(*player.Player)
	p.World().AddEntity(player.New(t.Name, p.Skin(), p.Position()))
}

func (NewPlayer) Allow(s cmd.Source) bool { return checkAdmin(s) }
