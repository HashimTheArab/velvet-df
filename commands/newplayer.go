package commands

import (
	"velvet/session"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

type NewPlayer struct {
	Name string
}

func (t NewPlayer) Run(source cmd.Source, _ *cmd.Output) {
	p, ok := source.(*player.Player)
	if !ok || !session.Get(p).HasFlag(session.Staff) {
		p.Message(NoPermission)
		return
	}
	p.World().AddEntity(player.New(t.Name, p.Skin(), p.Position()))
}

