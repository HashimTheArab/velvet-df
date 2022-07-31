package commands

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/biome"
	"github.com/df-mc/dragonfly/server/world/generator"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"io/ioutil"
	"strings"
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

// WorldCreate creates a new world.
type WorldCreate struct {
	Sub       cmd.SubCommand `cmd:"create"`
	Name      string         `cmd:"name"`
	Generator generate       `cmd:"generator"`
}

// Run ...
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

// Run ...
func (t WorldList) Run(_ cmd.Source, output *cmd.Output) {
	sb := &strings.Builder{}
	sb.WriteString("§d--World List (%v)--\n")
	var worlds uint32
	if files, err := ioutil.ReadDir("worlds"); err == nil {
		def := utils.WorldMG.DefaultWorld().Name()
		for _, f := range files {
			sb.WriteString("§e" + f.Name() + " §b" + utils.CoolAssArrow + " ")
			if _, ok := utils.WorldMG.World(f.Name()); ok || f.Name() == def {
				sb.WriteString("§aOnline\n")
			} else {
				sb.WriteString("§cOffline\n")
			}
			worlds++
		}
	}
	output.Printf(sb.String(), worlds)
}

// Run ...
func (t WorldCreate) Run(_ cmd.Source, output *cmd.Output) {
	if files, err := ioutil.ReadDir("worlds"); err == nil {
		for _, f := range files {
			if strings.EqualFold(f.Name(), t.Name) {
				output.Error(text.Colourf("<red>The world %s already exists.</red>", t.Name))
				return
			}
		}
	}
	var gen world.Generator
	switch t.Generator {
	case "flat":
		gen = generator.NewFlat(biome.Plains{}, []world.Block{block.Grass{}, block.Dirt{}, block.Dirt{}, block.Bedrock{}})
	case "void":
		gen = world.NopGenerator{}
	}
	err := utils.WorldMG.LoadWorld(t.Name, t.Name, world.Overworld, gen)
	if err != nil {
		output.Errorf(err.Error())
		return
	}
	output.Print(text.Colourf("<green>The world <purple>%s</purple> has been created!</green>", t.Name))
}

// generator ...
type generate string

// Type ...
func (generate) Type() string { return "Gennerator" }

// Options ...
func (generate) Options(cmd.Source) []string {
	return []string{"flat", "void"}
}

// Allow ...
func (WorldTeleport) Allow(s cmd.Source) bool { return checkStaff(s) }

// Allow ...
func (WorldList) Allow(s cmd.Source) bool { return checkStaff(s) || checkConsole(s) }

// Allow ...
func (WorldCreate) Allow(s cmd.Source) bool { return checkConsole(s) }
