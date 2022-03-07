package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"strings"
	"velvet/session"
	"velvet/utils"
)

type List struct{}

func (List) Run(_ cmd.Source, output *cmd.Output) {
	var players []string
	for _, v := range utils.Srv.Players() {
		s := session.Get(v)
		if s.Rank() != nil {
			players = append(players, s.Rank().Color+v.Name())
		} else {
			players = append(players, "§7"+v.Name())
		}
	}
	output.Printf("§6Players (§b%v§6): §b%v", len(utils.Srv.Players()), strings.Join(players, ", "))
}

func (List) Allow(s cmd.Source) bool { return checkStaff(s) || checkConsole(s) }
