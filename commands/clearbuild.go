package commands

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"time"
	"velvet/utils"
)

type ClearBuild struct{}

func (t ClearBuild) Run(s cmd.Source, output *cmd.Output) {
	p, pOk := s.(*player.Player)
	go func() {
		start := time.Now()
		if w, ok := utils.WorldMG.World(utils.Config.World.Build); ok {
			utils.BuildBlocks.Mutex.Lock()
			for pos := range utils.BuildBlocks.Blocks {
				w.SetBlock(pos, block.Air{})
				delete(utils.BuildBlocks.Blocks, pos)
			}
			utils.BuildBlocks.Mutex.Unlock()
			if pOk && p != nil {
				p.Messagef("§dBuild was cleared in §e%v.", time.Now().Sub(start).Round(time.Millisecond*10).String())
			}
		} else {
			if pOk && p != nil {
				p.Messagef("§cBuild is currently offline.")
			}
		}
	}()
}

func (ClearBuild) Allow(s cmd.Source) bool { return checkStaff(s) || checkConsole(s) }
