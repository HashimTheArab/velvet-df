package utils

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"sync"
)

type buildBlocks struct {
	Blocks map[cube.Pos]uint8
	Mutex  sync.Mutex
}

var BuildBlocks buildBlocks

func (b *buildBlocks) Set(pos cube.Pos) {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()
	b.Blocks[pos] = 0
}
