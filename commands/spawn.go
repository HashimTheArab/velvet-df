package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"velvet/session"
	"velvet/utils"
)

type Spawn struct{}

func (t Spawn) Run(source cmd.Source, _ *cmd.Output) {
	p := source.(*player.Player)
	session.Get(p).TeleportToSpawn()
	p.Message(utils.Config.Message.WelcomeToSpawn)
}

func (Spawn) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}
