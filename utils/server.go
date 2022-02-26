package utils

import (
	"github.com/df-mc/dragonfly/server"
	"go.uber.org/atomic"
	"velvet/utils/worldmanager"
)

var (
	Srv     *server.Server
	WorldMG *worldmanager.WorldManager
	// OnlineCount is the amount of players online
	OnlineCount atomic.Uint32
	// Started is a timestamp of when the server was turned on
	Started int64
)
