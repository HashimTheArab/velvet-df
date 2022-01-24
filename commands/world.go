package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"velvet/utils"
)

type WorldTeleport struct {
	Name string `name:"name"`
}

func (t WorldTeleport) Run(source cmd.Source, _ *cmd.Output) {
	p := source.(*player.Player)

	if w, ok := utils.WorldMG.World(t.Name); ok {
		w.AddEntity(p)
		p.Message("§dYou have been teleported to the world §e" + t.Name)
		return
	}

	p.Message("§cThat world does not exist or is not loaded.")
}

func (WorldTeleport) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok && checkStaff(s)
}
