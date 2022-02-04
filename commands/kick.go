package commands

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"velvet/session"
	"velvet/utils"
)

type Kick struct {
	Player []cmd.Target `name:"victim"`
	Reason cmd.Varargs  `name:"reason"`
}

func (t Kick) Run(source cmd.Source, output *cmd.Output) {
	if target, ok := t.Player[0].(*player.Player); ok {
		if target.Name() == source.Name() || (source.Name() != utils.Config.Staff.Owner.Name && session.Get(target).HasFlag(session.Staff)) {
			output.Print(utils.Config.Message.CannotPunishPlayer)
			return
		}
		target.Disconnect(fmt.Sprintf(utils.Config.Kick.Screen, target.Name(), string(t.Reason)))
		_, _ = fmt.Fprintf(chat.Global, utils.Config.Kick.Broadcast+"\n", source.Name(), target.Name(), string(t.Reason))
		return
	}
	output.Printf(PlayerNotOnline)
}

func (Kick) Allow(s cmd.Source) bool { return checkStaff(s) || checkConsole(s) }
