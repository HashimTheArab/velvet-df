package session

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	_ "unsafe"
)

func (s *Session) WritePacket(pk packet.Packet) {
	session_writePacket(s.NetworkSession, pk)
}

//go:linkname player_session github.com/df-mc/dragonfly/server/player.(*Player).session
//noinspection ALL
func player_session(*player.Player) *session.Session

//go:linkname session_writePacket github.com/df-mc/dragonfly/server/session.(*Session).writePacket
//noinspection ALL
func session_writePacket(*session.Session, packet.Packet)
