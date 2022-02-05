package session

// StaffMap is the custom type with custom methods for the staff session map.
type StaffMap map[string]*Session

// staff is a list of the sessions of all online staff.
var staff = StaffMap{}

// AllStaff will return all the online staff.
func AllStaff() StaffMap {
	return staff
}

// Message will send a message to all online staff.
func (StaffMap) Message(a ...interface{}) {
	for _, v := range staff {
		v.Player.Message(a)
	}
}

// Messagef will send a formatted message to all online staff.
func (StaffMap) Messagef(f string, a ...interface{}) {
	for _, v := range staff {
		v.Player.Messagef(f, a)
	}
}

// Whisper will send a gray italic message to all online staff.
func (StaffMap) Whisper(f string, a ...interface{}) {
	for _, v := range staff {
		v.Player.Messagef("§7§o["+f+"]", a)
	}
}
