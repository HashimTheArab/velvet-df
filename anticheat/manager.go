package anticheat

var checks map[byte]*Check

// todo: move ac to gophertunnel proxy
func init() {
	register(&Check{
		Name:        "Autoclicker",
		Subtype:     "A",
		Description: "Checks for high cps",
	})
	register(&Check{
		Name:        "Timer",
		Subtype:     "A",
		Description: "Checks if the player is sending packets too fast",
	})
}

func register(check *Check) {
	id := byte(len(checks) + 1)
	check.ID = id
	checks[id] = check
}
