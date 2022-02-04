package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"io/ioutil"
	"velvet/utils"
)

type WorldTeleport struct {
	Sub     teleport
	Name    string       `name:"name"`
	Targets []cmd.Target `optional:"" name:"target"`
}

type WorldList struct {
	Sub list
}

func (t WorldTeleport) Run(source cmd.Source, output *cmd.Output) {
	p := source.(*player.Player)
	if w, ok := utils.WorldMG.World(t.Name); ok {
		if len(t.Targets) > 0 {
			if checkAdmin(source) {
				for _, v := range t.Targets {
					if pl, ok := v.(*player.Player); ok {
						w.AddEntity(pl)
					}
				}
			} else {
				if len(t.Targets) == 0 {
					if pl, ok := t.Targets[0].(*player.Player); ok {
						w.AddEntity(pl)
					}
				}
			}
		} else {
			w.AddEntity(p)
			output.Printf("§dYou have been teleported to the world §e%v!", t.Name)
		}
		return
	}

	output.Print("§cThat world does not exist or is not loaded.")
}

func (t WorldList) Run(_ cmd.Source, output *cmd.Output) {
	msg := "§d--World List (%v)--\n"
	var worlds uint32
	if files, err := ioutil.ReadDir("worlds"); err == nil {
		def := utils.WorldMG.DefaultWorld().Name()
		for _, f := range files {
			msg += "§e" + f.Name() + " §b" + utils.CoolAssArrow + " "
			if _, ok := utils.WorldMG.World(f.Name()); ok || f.Name() == def {
				msg += "§aOnline\n"
			} else {
				msg += "§cOffline\n"
			}
			worlds++
		}
	}
	output.Printf(msg, worlds)
}

func (WorldTeleport) Allow(s cmd.Source) bool { return checkStaff(s) }

type teleport string

func (teleport) SubName() string { return "teleport" }

func (WorldList) Allow(s cmd.Source) bool {
	return checkStaff(s) || checkConsole(s)
}

type list string

func (list) SubName() string { return "list" }
