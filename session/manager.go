package session

import (
	"github.com/df-mc/dragonfly/server/player"
)

var sessions = map[string]*Session{}

func New(name string) *Session {
	session := &Session{}
	sessions[name] = session
	return session
}

func Get(p *player.Player) *Session {
	return sessions[p.Name()]
}

func FromName(name string) *Session {
	return sessions[name]
}

func (s *Session) Close() {
	s.Save()
	if s.HasFlag(FlagStaff) {
		delete(sessions, s.Player.Name())
	}
	delete(sessions, s.Player.Name())
}

func (s *Session) CloseWithoutSaving() {
	if s.HasFlag(FlagStaff) {
		delete(sessions, s.Player.Name())
	}
	delete(sessions, s.Player.Name())
}

func All() map[string]*Session {
	return sessions
}
