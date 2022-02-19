package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
)

type Wand struct { // todo
	wePerms
}

func (t Wand) Run(source cmd.Source, output *cmd.Output) {
	p := source.(*player.Player)
	axe := item.NewStack(item.Axe{Tier: item.ToolTierNetherite}, 1)
	if _, err := p.Inventory().AddItem(axe); err != nil {
		output.Error(err)
		return
	}
	// todo
}

type wePerms struct{}

func (wePerms) Allow(s cmd.Source) bool {
	return !checkConsole(s) && (checkStaff(s) || checkBuilder(s))
}
