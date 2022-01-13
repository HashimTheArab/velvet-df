package commands

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"velvet/session"
	"velvet/utils"
)

type kickReason string

type Kick struct {
	Player []cmd.Target
	Reason kickReason
}

func (t Kick) Run(source cmd.Source, output *cmd.Output) {
	p, ok := source.(*player.Player)
	if ok && !session.Get(p).HasFlag(session.Staff) {
		p.Message(NoPermission)
		return
	}

	if target, ok := t.Player[0].(*player.Player); ok {
		if (session.Get(target).HasFlag(session.Staff) || target.Name() == p.Name()) && p.Name() != utils.Config.Staff.Owner.Name {
			p.Message(utils.Config.Message.CannotPunishPlayer)
			return
		}
		p.Message(utils.Config.Message.KickedPlayer, target.Name(), string(t.Reason))
		_, _ = fmt.Fprintf(chat.Global, utils.Config.Message.KickedPlayerBroadcast, target.Name(), string(t.Reason))
		return
	}

	output.Printf(PlayerNotOnline)
}

func (kickReason) Type() string {
	return "Kick"
}

func (kickReason) Options(cmd.Source) []string {
	return []string{"spam", "interfering", "2v1"}
}
