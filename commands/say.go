package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player/chat"
	"velvet/utils"
)

type Say struct {
	Message cmd.Varargs `name:"message"`
}

func (t Say) Run(cmd.Source, *cmd.Output) {
	_, _ = chat.Global.WriteString("§l§b[§6CONSOLE§b] §8" + utils.CoolAssArrow + " §b" + string(t.Message))
}

func (Say) Allow(s cmd.Source) bool { return checkConsole(s) }
