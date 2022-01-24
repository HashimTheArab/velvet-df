package anticheat

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/player/chat"
	"velvet/session"
	"velvet/utils"
)

type Check struct {
	Name          string
	Subtype       string
	Description   string
	Punishment    Punishment
	MaxViolations byte
	ID            byte
}

type Session struct {
}

type Punishment byte

const (
	None Punishment = iota
	Kick
	Ban
)

func (c *Check) Flag(s *session.Session) {

}

func (c *Check) Punish(s *session.Session) {
	if c.Punishment == Kick {
		name := s.Player.Name()
		s.Player.Disconnect(fmt.Sprintf(utils.Config.AntiCheat.KickScreen, c.Name, c.Subtype))
		_, _ = fmt.Fprintf(chat.Global, utils.Config.AntiCheat.KickBroadcast+"\n", name, c.Name, c.Subtype)
		return
	}
	s.Player.Disconnect(fmt.Sprintf(utils.Config.AntiCheat.BanScreen, c.Name, c.Subtype))
	_, _ = fmt.Fprintf(chat.Global, utils.Config.AntiCheat.KickBroadcast+"\n", s.Player.Name(), c.Name, c.Subtype)
}
