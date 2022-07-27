package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"strconv"
	"velvet/utils"
)

type Transfer struct {
	Address string                     `cmd:"address"`
	Port    uint16                     `cmd:"port"`
	Targets cmd.Optional[[]cmd.Target] `cmd:"victim"`
}

func (t Transfer) Run(source cmd.Source, output *cmd.Output) {
	p, ok := source.(*player.Player)
	targets := t.Targets.LoadOr(nil)
	if len(targets) > 0 {
		for _, v := range targets {
			vp, ok := v.(*player.Player)
			if ok {
				if err := vp.Transfer(t.Address + strconv.Itoa(int(t.Port))); err != nil {
					output.Print(utils.Config.Message.ServerNotAvailable)
				}
			}
		}
	} else if ok {
		if err := p.Transfer(t.Address + ":" + strconv.Itoa(int(t.Port))); err != nil {
			p.Message(utils.Config.Message.ServerNotAvailable)
		}
	}
}

func (Transfer) Allow(s cmd.Source) bool { return checkAdmin(s) }
