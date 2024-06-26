package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"velvet/db"
	"velvet/form"
	"velvet/session"
)

type StatsOnline struct {
	Targets cmd.Optional[[]cmd.Target] `cmd:"target"`
}

type StatsOffline struct {
	Target string `cmd:"target"`
}

func (t StatsOnline) Run(source cmd.Source, output *cmd.Output) {
	p := source.(*player.Player)
	targets := t.Targets.LoadOr(nil)
	if len(targets) > 0 {
		if len(targets) > 1 {
			output.Error("You can only check the stats of one player at a time.")
			return
		}
		if pl, ok := targets[0].(*player.Player); ok {
			p.SendForm(form.StatsOnline(session.Get(pl)))
		} else {
			output.Error("Player not found.")
		}
	} else {
		p.SendForm(form.StatsOnline(session.Get(p)))
	}
}

func (t StatsOffline) Run(source cmd.Source, o *cmd.Output) {
	u, err := db.LoadOfflinePlayer(t.Target)
	if err != nil {
		o.Error("That player has not joined before.")
		return
	}

	source.(*player.Player).SendForm(form.StatsOffline(u))
}

func (StatsOnline) Allow(s cmd.Source) bool  { return !checkConsole(s) }
func (StatsOffline) Allow(s cmd.Source) bool { return !checkConsole(s) }
