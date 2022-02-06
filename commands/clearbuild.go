package commands

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/cmd"
	"time"
	"velvet/utils"
)

type ClearBuild struct{}

func (t ClearBuild) Run(_ cmd.Source, output *cmd.Output) {
	go func() {
		start := time.Now()
		if w, ok := utils.WorldMG.World(utils.Config.World.Build); ok {
			utils.BuildBlocks.Mutex.Lock()
			for pos := range utils.BuildBlocks.Blocks {
				w.SetBlock(pos, block.Air{})
				delete(utils.BuildBlocks.Blocks, pos)
			}
			utils.BuildBlocks.Mutex.Unlock()
			output.Printf("§dBuild was cleared in §e%v.", time.Now().Sub(start).Round(time.Millisecond*10).String())
		} else {
			output.Printf("§cBuild is currently offline.")
		}
	}()
}

func (ClearBuild) Allow(s cmd.Source) bool { return checkStaff(s) || checkConsole(s) }
