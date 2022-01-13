package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"velvet/session"
	"velvet/utils"
)

type WorldTeleport struct {
	Name string `name:"name"`
}

func (t WorldTeleport) Run(source cmd.Source, output *cmd.Output) {
	p, ok := source.(*player.Player)
	if !ok || !session.Get(p).HasFlag(session.Staff) {
		p.Message(NoPermission)
		return
	}

	if w, ok := utils.WorldMG.World(t.Name); ok {
		session.Get(p).ChangeWorld(w)
		p.Message("§7You have been teleported to the world §b" + t.Name)
		return
	}

	p.Message("§cThat world does not exist or is not loaded.")
}
