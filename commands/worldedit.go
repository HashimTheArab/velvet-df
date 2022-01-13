package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/we/brush"
)

type WorldEdit struct {

}

func (t WorldEdit) Run(source cmd.Source, output *cmd.Output) {
	p, ok := source.(*player.Player)
	if !ok {
		return
	}

	p.SendForm(brush.NewSelectionForm())
}