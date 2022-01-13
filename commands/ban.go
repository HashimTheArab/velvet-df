package commands

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"velvet/session"
	"velvet/utils"
)

type banReason string
type banLength string

type Ban struct {
	Player []cmd.Target
	Reason banReason
	Length banLength `optional:""`
	Silent bool      `optional:""`
}

func (t Ban) Run(source cmd.Source, output *cmd.Output) {
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

func (banReason) Type() string { return "BanReason" }
func (banLength) Type() string { return "BanLength" }

func (banReason) Options(cmd.Source) []string {
	return []string{"cheats", "autoclicker", "debounce", "cps", "reach"}
}
func (banLength) Options(cmd.Source) []string { return []string{"14d", "30d"} }
