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
	MaxViolations uint8
}

type Punishment byte

const (
	None Punishment = iota
	Kick
	Ban
)

func (c *Check) Punish(s *session.Session) {
	if c.Punishment == Kick {
		name := s.Player.Name()
		s.Player.Disconnect(utils.Config.Message.AntiCheatKick, c.Name, c.Subtype)
		_, _ = fmt.Fprintf(chat.Global, utils.Config.Message.AntiCheatKickBroadcast, name, c.Name, c.Subtype)
		return
	}
	// todo: bans
}
