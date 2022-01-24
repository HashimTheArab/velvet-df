package utils

import (
	"github.com/df-mc/dragonfly/server"
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
)
