package session

import (
	"github.com/df-mc/dragonfly/server/player"
)

var sessions = map[string]*Session{}
var staff = map[string]*Session{}

func New(name string) *Session {
	session := &Session{}
	sessions[name] = session
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

func AllStaff() map[string]*Session {
	return staff
}
