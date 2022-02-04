package session

type combat struct {
	time uint8
}

var CombatBannedCommands = []string{"spawn"}

func (c *combat) Tag(tag bool) {
	if tag {
		c.time = 15
	} else {
		c.time = 0
	}
}

func (c *combat) Tagged() bool {
	return c.time > 0
}

func (c *combat) Time() uint8 {
	return c.time
}
