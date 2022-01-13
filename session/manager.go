package session

import (
	"github.com/df-mc/dragonfly/server/player"
)

var sessions = map[string]*Session{}

func New(p *player.Player) *Session {
	session := &Session{Player: p}
	sessions[p.Name()] = session
	session.OnJoin()
	return session
}

func Get(p *player.Player) *Session {
	return sessions[p.Name()]
}

func (s *Session) Close() {
	s.OnQuit()
	delete(sessions, s.Player.Name())
}

func All() map[string]*Session {
	return sessions
}
