package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

type TimeSet struct {
	Sub  set
	Time int `name:"time"`
}

type set string

func (t TimeSet) Run(source cmd.Source, output *cmd.Output) {
	if p, ok := source.(*player.Player); ok {
		p.World().SetTime(t.Time)
		output.Printf("§dTime has been set to §e%v.", t.Time)
	}
}

func (set) SubName() string {
	return "set"
}

func (TimeSet) Allow(s cmd.Source) bool { return checkStaff(s) }
