package entity

import (
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	_ "unsafe"
)

//go:linkname session_writePacket github.com/df-mc/dragonfly/server/session.(*Session).writePacket
//noinspection ALL
func session_writePacket(*session.Session, packet.Packet)

//go:linkname session_entityRuntimeID github.com/df-mc/dragonfly/server/session.(*Session).entityRuntimeID
//noinspection ALL
func session_entityRuntimeID(*session.Session, world.Entity) uint64
