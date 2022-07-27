package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"io/ioutil"
	"velvet/utils"
)

// WorldTeleport lets you teleport to another loaded world.
type WorldTeleport struct {
	Sub     cmd.SubCommand             `cmd:"teleport"`
	Name    string                     `cmd:"name"`
	Targets cmd.Optional[[]cmd.Target] `cmd:"target"`
}

// WorldList outputs all loaded worlds.
type WorldList struct {
	Sub cmd.SubCommand `cmd:"list"`
}

func (t WorldTeleport) Run(source cmd.Source, output *cmd.Output) {
	p := source.(*player.Player)
	targets := t.Targets.LoadOr(nil)
	if w, ok := utils.WorldMG.World(t.Name); ok {
		if len(targets) > 0 {
			if checkAdmin(source) {
				for _, v := range targets {
					if pl, ok := v.(*player.Player); ok {
						w.AddEntity(pl)
					}
				}
			} else {
				if len(targets) == 0 {
					if pl, ok := targets[0].(*player.Player); ok {
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

func (WorldList) Allow(s cmd.Source) bool {
	return checkStaff(s) || checkConsole(s)
}
