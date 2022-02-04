package session

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"time"
)

type Cooldown struct {
	length time.Duration
	last   time.Time
}

type cooldownMap map[cooldownType]*Cooldown
type cooldownType uint8

const (
	CooldownTypePearl cooldownType = iota
	CooldownTypeChat
)

// Handle ...
func (c cooldownMap) Handle(ctx *event.Context, p *player.Player, cdType cooldownType) {
	cd := c[cdType]
	diff := time.Now().Sub(cd.last)
	if diff < cd.length {
		ctx.Cancel()
		p.Messagef("§aYou are on cooldown for §c%v!", (cd.length - diff).Round(time.Millisecond*10).String())
	} else {
		c.Set(cd)
	}
}

func (c cooldownMap) Set(cd *Cooldown) {
	cd.last = time.Now()
}
