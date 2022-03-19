package session

import (
	"github.com/df-mc/dragonfly/server/player"
	"sync"
)

type Sessions struct {
	list  map[string]*Session
	mutex sync.Mutex
}

var sessions = Sessions{
	list: make(map[string]*Session),
}

func New(name string) *Session {
	session := &Session{}

	sessions.mutex.Lock()
	sessions.list[name] = session
	sessions.mutex.Unlock()

	return session
}

func Get(p *player.Player) *Session {
	if p == nil {
		return nil
	}
	sessions.mutex.Lock()
	defer sessions.mutex.Unlock()
	return sessions.list[p.Name()]
}

func FromName(name string) *Session {
	sessions.mutex.Lock()
	defer sessions.mutex.Unlock()
	return sessions.list[name]
}

func (s *Session) Close() {
	s.Save()
	if s.HasFlag(FlagStaff) {
		RemoveStaff(s)
	}

	sessions.mutex.Lock()
	defer sessions.mutex.Unlock()
	delete(sessions.list, s.Player.Name())
}

func (s *Session) CloseWithoutSaving(name string) {
	sessions.mutex.Lock()
	defer sessions.mutex.Unlock()
	if s.HasFlag(FlagStaff) {
		delete(sessions.list, name)
	}
	delete(sessions.list, name)
}

func All() *Sessions {
	return &sessions
}

func (s *Sessions) UpdateScoreboards(online, kd bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, v := range s.list {
		if !v.Offline() {
			v.UpdateScoreboard(online, kd)
		}
	}
}
