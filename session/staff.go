package session

import "sync"

// Staff contains the staff list and all its methods.
type Staff struct {
	list  map[string]*Session
	mutex sync.Mutex
}

// staff contains the sessions of all online staff.
var staff = Staff{
	list: make(map[string]*Session),
}

// AllStaff will return all the online staff.
func AllStaff() *Staff {
	return &staff
}

// Message will send a message to all online staff.
func (s *Staff) Message(a ...interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, v := range s.list {
		v.Player.Message(a)
	}
}

// Messagef will send a formatted message to all online staff.
func (s *Staff) Messagef(f string, a ...interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, v := range s.list {
		v.Player.Messagef(f, a...)
	}
}

// Whisper will send a gray italic message to all online staff.
func (s *Staff) Whisper(f string, a ...interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, v := range s.list {
		v.Player.Messagef("§7§o["+f+"]", a...)
	}
}

// AddStaff adds a session to the staff list.
func AddStaff(s *Session) {
	staff.mutex.Lock()
	defer staff.mutex.Unlock()
	if _, ok := staff.list[s.Player.Name()]; !ok {
		staff.list[s.Player.Name()] = s
	}
}

// RemoveStaff removes a session from the staff list.
func RemoveStaff(s *Session) {
	if s != nil {
		staff.mutex.Lock()
		defer staff.mutex.Unlock()
		delete(staff.list, s.Player.Name())
	}
}
