package utils

import (
	"github.com/df-mc/dragonfly/server"
	worldmanager "github.com/emperials/df-worldmanager"
	"go.uber.org/atomic"
)

var (
	Srv         *server.Server
	WorldMG     *worldmanager.WorldManager
	OnlineCount atomic.Uint32
)
