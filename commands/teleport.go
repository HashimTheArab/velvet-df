package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"strconv"
	"velvet/session"
	"velvet/utils"
)

type TeleportToPos struct {
	Destination mgl64.Vec3 `name:"destination"`
}

type TeleportToTarget struct {
	Player []cmd.Target `name:"destination"`
}

type TeleportTargetToTarget struct {
	Players []cmd.Target `name:"victim"`
	Targets []cmd.Target `name:"destination"`
}

type TeleportTargetToPos struct {
	Players     []cmd.Target `name:"victim"`
	Destination mgl64.Vec3   `name:"destination"`
}

func (t TeleportToPos) Run(source cmd.Source, output *cmd.Output) {
	if p, ok := source.(*player.Player); ok {
		if !session.Get(p).HasFlag(session.Staff) {
			output.Error(NoPermission)
			return
		}
		p.Teleport(t.Destination)
		output.Printf(utils.Config.Message.TeleportSelfToPos, t.Destination)
	}
}

func (t TeleportToTarget) Run(source cmd.Source, output *cmd.Output) {
	if p, ok := source.(*player.Player); ok {
		if !session.Get(p).HasFlag(session.Staff) {
			output.Error(NoPermission)
			return
		}
		tg, ok := t.Player[0].(*player.Player)
		if !ok {
			output.Error(PlayerNotOnline)
			return
		}
		if p.World().Name() != tg.World().Name() {
			tg.World().AddEntity(p)
		}
		p.Teleport(tg.Position())
		output.Printf(utils.Config.Message.TeleportSelfToPlayer, t.Player[0].Name())
	}
}

func (t TeleportTargetToTarget) Run(source cmd.Source, output *cmd.Output) {
	p := source.(*player.Player)
	if !session.Get(p).HasFlag(session.Staff) {
		output.Error(NoPermission)
		return
	}
	if len(t.Players) > 1 && !session.Get(p).HasFlag(session.Admin) {
		output.Error(NoPermission)
		return
	}
	tg, ok := t.Targets[0].(*player.Player)
	if !ok {
		output.Error(PlayerNotOnline)
		return
	}
	output.Printf(utils.Config.Message.TeleportTargetToTarget, teleportTargets(t.Targets, tg.Position(), tg.World()), t.Targets[0].Name())
}

func (t TeleportTargetToPos) Run(source cmd.Source, output *cmd.Output) {
	if p, ok := source.(*player.Player); ok {
		s := session.Get(p)
		if !s.HasFlag(session.Staff) {
			output.Error(NoPermission)
			return
		}
		if len(t.Players) > 1 && !s.HasFlag(session.Admin) {
			output.Error(NoPermission)
			return
		}
	}
	output.Printf(utils.Config.Message.TeleportTargetToPos, teleportTargets(t.Players, t.Destination, nil), t.Destination)
}

// teleportTargets teleports a list of targets to a specified position
func teleportTargets(targets []cmd.Target, destination mgl64.Vec3, w *world.World) string {
	for _, target := range targets {
		if p, ok := target.(*player.Player); ok {
			if w != nil {
				w.AddEntity(p)
			}
			if destination != (mgl64.Vec3{}) {
				p.Teleport(destination)
			}
		}
	}
	if len(targets) > 1 {
		return strconv.Itoa(len(targets)) + " players"
	}
	return targets[0].Name()
}
