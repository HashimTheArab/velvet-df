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

// All returns every loaded sessions.
func All() *Sessions {
	return &sessions
}

// Get gets a session based on the player.
func Get(p *player.Player) *Session {
	if p == nil {
		return nil
	}
	sessions.mutex.Lock()
	defer sessions.mutex.Unlock()
	return sessions.list[p.Name()]
}

// FromName gets a session based on a name.
func FromName(name string) *Session {
	sessions.mutex.Lock()
	defer sessions.mutex.Unlock()
	return sessions.list[name]
}

// Close closes and saves the session.
func (s *Session) Close() {
	s.closed.Store(true)
	if s.HasFlag(FlagStaff) {
		RemoveStaff(s)
	}

	sessions.mutex.Lock()
	defer sessions.mutex.Unlock()
	delete(sessions.list, s.Player.Name())
}

// CloseWithoutSaving closes the session without saving data.
func (s *Session) CloseWithoutSaving(name string) {
	sessions.mutex.Lock()
	defer sessions.mutex.Unlock()
	if s.HasFlag(FlagStaff) {
		delete(sessions.list, name)
	}
	delete(sessions.list, name)
}

// UpdateScoreboards updates the scoreboard of every session.
func (s *Sessions) UpdateScoreboards(online, kd bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, v := range s.list {
		if !v.Offline() {
			v.UpdateScoreboard(online, kd)
		}
	}
}
