package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"velvet/utils"
)

type Spawn struct{}

func (t Spawn) Run(source cmd.Source, _ *cmd.Output) {
	p := source.(*player.Player)
	utils.Srv.World().AddEntity(p)
	p.Teleport(utils.Srv.World().Spawn().Vec3())
	p.Message(utils.Config.Message.WelcomeToSpawn)
}

func (Spawn) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}
