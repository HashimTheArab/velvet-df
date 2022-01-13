package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"velvet/session"
	"velvet/utils"
)

type Build struct {
	Player []cmd.Target `optional:"" name:"target"`
}

func (t Build) Run(source cmd.Source, output *cmd.Output) {
	p, ok := source.(*player.Player)
	var s *session.Session
	if ok {
		s = session.Get(p)
		if !s.HasFlag(session.Admin) && !s.HasFlag(session.Building) {
			p.Message(NoPermission)
			return
		}
	}
	if len(t.Player) > 0 {
		if len(t.Player) > 1 {
			output.Error(utils.Config.Message.BuildTooManyPlayers)
			return
		}
		if target, ok := t.Player[0].(*player.Player); ok {
			s = session.Get(target)
			if s.HasFlag(session.Building) {
				target.Messagef(utils.Config.Message.UnsetBuilderModeByPlayer, p.Name())
				output.Printf(utils.Config.Message.UnsetPlayerInBuilderMode, target.Name())
			} else {
				target.Messagef(utils.Config.Message.SetBuilderModeByPlayer, p.Name())
				output.Printf(utils.Config.Message.SetPlayerInBuilderMode, target.Name())
			}
			s.SetFlag(session.Building)
			return
		}
		output.Errorf(PlayerNotOnline, t.Player[0].Name())
		return
	}
	if ok {
		if s.HasFlag(session.Building) {
			output.Print(utils.Config.Message.SelfNotInBuilderMode)
		} else {
			output.Print(utils.Config.Message.SelfInBuilderMode)
		}
		s.SetFlag(session.Building)
	}
}
