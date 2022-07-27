package commands

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/we/palette"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"math"
	"math/rand"
	"time"
	"velvet/session"
)

// Wand is used to give the player a wand.
type Wand struct{ wePerms }

// Fill is used to fill an area.
type Fill struct {
	wePerms
	Palette string `cmd:"palette"`
}

// PaletteSet is used to set a block palette.
type PaletteSet struct {
	wePerms
	palette.SetCommand
}

// PaletteSave is used to save a block palette.
type PaletteSave struct {
	wePerms
	palette.SaveCommand
}

// PaletteDelete is used to delete a block palette.
type PaletteDelete struct {
	wePerms
	palette.DeleteCommand
}

// Run ...
func (t Wand) Run(source cmd.Source, _ *cmd.Output) {
	_, _ = source.(*player.Player).Inventory().AddItem(item.NewStack(item.Axe{Tier: item.ToolTierNetherite}, 1).WithCustomName(text.Colourf("<gold>Wand</gold>")).WithValue("wand", true))
}

// Run ...
func (t Fill) Run(source cmd.Source, o *cmd.Output) {
	p := source.(*player.Player)
	s := session.Get(p)

	palette, ok := palette.LookupHandler(p)
	if !ok {
		return
	}
	found, ok := palette.Palette(t.Palette)
	if !ok || len(found.Blocks()) == 0 {
		o.Error("Invalid palette, create one using /palette")
		return
	}
	var names []string
	for _, b := range found.Blocks() {
		n, _ := b.EncodeBlock()
		names = append(names, n)
	}

	pos1, pos2 := s.WandPos()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	minX, maxX := int(math.Min(pos1.X(), pos2.X())), int(math.Max(pos1.X(), pos2.X()))
	minY, maxY := int(math.Min(pos1.Y(), pos2.Y())), int(math.Max(pos1.Y(), pos2.Y()))
	minZ, maxZ := int(math.Min(pos1.Z(), pos2.Z())), int(math.Max(pos1.Z(), pos2.Z()))
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			for z := minZ; z <= maxZ; z++ {
				p.World().SetBlock(cube.Pos{x, y, z}, found.Blocks()[r.Intn(len(found.Blocks()))], &world.SetOpts{
					DisableBlockUpdates:       true,
					DisableLiquidDisplacement: true,
				})
			}
		}
	}
	o.Print(text.Colourf("<green>Filled area %v to %v with palette %v (%v blocks)", pos1, pos2, names, (maxY-minY)*(maxX-minX)*(maxZ-minZ)))
}

// wePerms contains permissions for all world edit commands.
type wePerms struct{}

// Allow ...
func (wePerms) Allow(s cmd.Source) bool {
	return !checkConsole(s) && (checkStaff(s) || checkBuilder(s))
}
