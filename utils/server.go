package utils

import (
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	worldmanager "github.com/emperials/df-worldmanager"
	"go.uber.org/atomic"
)

var (
	Srv     *server.Server
	WorldMG *worldmanager.WorldManager
	// OnlineCount is the amount of players online
	OnlineCount atomic.Uint32
	// Started is a timestamp of when the server was turned on
	Started int64
	// PlacedBuildBlocks tracks the player placed blocks in the build gamemode
	PlacedBuildBlocks = map[cube.Pos]world.Block{} // todo: no need to store the block
)
