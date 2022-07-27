package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

type TimeSet struct {
	Sub  cmd.SubCommand `cmd:"set"`
	Time int            `cmd:"time"`
}

func (t TimeSet) Run(source cmd.Source, output *cmd.Output) {
	if p, ok := source.(*player.Player); ok {
		p.World().SetTime(t.Time)
		output.Printf("§dTime has been set to §e%v.", t.Time)
	}
}

func (TimeSet) Allow(s cmd.Source) bool { return checkStaff(s) }
