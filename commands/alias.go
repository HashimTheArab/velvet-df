package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"strings"
	"velvet/db"
	"velvet/utils"
)

type Alias struct {
	Target []cmd.Target `name:"target"`
}

type AliasOffline struct {
	Target string `name:"target"`
}

func (t Alias) Run(_ cmd.Source, output *cmd.Output) {
	if p, ok := t.Target[0].(*player.Player); ok {
		handleAlias(p.Name(), output)
	} else {
		output.Error(PlayerNotFound)
	}
}

func (t AliasOffline) Run(_ cmd.Source, output *cmd.Output) {
	handleAlias(t.Target, output)
}

func handleAlias(name string, output *cmd.Output) {
	if did, accounts := db.GetAlias(name); did != "" {
		output.Printf(utils.Config.Message.Alias, name, did, strings.Join(accounts, " §d|§r "))
	} else {
		output.Printf(utils.Config.Message.NeverJoined, name)
	}
}

func (Alias) Allow(s cmd.Source) bool        { return checkStaff(s) || checkConsole(s) }
func (AliasOffline) Allow(s cmd.Source) bool { return checkStaff(s) || checkConsole(s) }
