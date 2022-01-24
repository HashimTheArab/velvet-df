package commands

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"velvet/utils"
)

type SpawnPoint struct{}

func (t SpawnPoint) Run(source cmd.Source, _ *cmd.Output) {
	p := source.(*player.Player)
	pos := cube.PosFromVec3(p.Position())
	p.World().SetSpawn(pos)
	p.Messagef(utils.Config.Message.DefaultSpawnSet, pos)
}

func (SpawnPoint) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok && checkAdmin(s)
}
