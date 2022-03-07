package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"velvet/db"
	"velvet/form"
	"velvet/session"
)

type StatsOnline struct {
	Target []cmd.Target `name:"target" optional:""`
}

type StatsOffline struct {
	Target string `name:"target"`
}

func (t StatsOnline) Run(source cmd.Source, output *cmd.Output) {
	p, ok := source.(*player.Player)
	if !ok {
		return
	}

	var data db.PlayerData
	var name string
	if len(t.Target) > 0 {
		if len(t.Target) > 1 {
			output.Error("You can only check the stats of one player at a time.")
			return
		}
		if pl, ok := t.Target[0].(*player.Player); ok {
			s := session.Get(pl)
			data = db.PlayerData{
				PlayerStats: s.Stats,
				Rank:        s.Rank().Name,
			}
			name = pl.Name()
		} else {
			output.Error("Player not found.")
		}
	} else {
		s := session.Get(p)
		data = db.PlayerData{
			PlayerStats: s.Stats,
			Rank:        s.Rank().Name,
		}
		name = p.NameTag()
	}
	sendDataForm(p, name, data)
}

func (t StatsOffline) Run(source cmd.Source, _ *cmd.Output) {
	if p, ok := source.(*player.Player); ok {
		sendDataForm(p, t.Target, db.GetData(t.Target))
	}
}

func sendDataForm(p *player.Player, target string, data db.PlayerData) {
	p.SendForm(form.Stats(target, data))
}

func (StatsOnline) Allow(s cmd.Source) bool  { return !checkConsole(s) }
func (StatsOffline) Allow(s cmd.Source) bool { return !checkConsole(s) }
