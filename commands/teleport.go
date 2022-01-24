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
	source.(*player.Player).Teleport(t.Destination)
	output.Printf(utils.Config.Message.TeleportSelfToPos, t.Destination)
}

func (t TeleportToTarget) Run(source cmd.Source, output *cmd.Output) {
	p := source.(*player.Player)
	tg, ok := t.Player[0].(*player.Player)
	if !ok {
		output.Error(PlayerNotOnline)
		return
	}
	if p.World() != tg.World() {
		tg.World().AddEntity(p)
	}
	p.Teleport(tg.Position())
	output.Printf(utils.Config.Message.TeleportSelfToPlayer, t.Player[0].Name())
}

func (t TeleportTargetToTarget) Run(source cmd.Source, output *cmd.Output) {
	p, ok := source.(*player.Player)
	if len(t.Players) > 1 && ok && !session.Get(p).HasFlag(session.Admin) {
		output.Error(NoPermission)
		return
	}
	if len(t.Targets) > 1 {
		output.Error("§cYou cannot have more than one destination.")
		return
	}
	tg, ok := t.Targets[0].(*player.Player)
	if !ok {
		output.Error(PlayerNotOnline)
		return
	}
	output.Printf(utils.Config.Message.TeleportTargetToTarget, teleportTargets(t.Players, tg.Position(), tg.World()), tg.Name())
}

func (t TeleportTargetToPos) Run(source cmd.Source, output *cmd.Output) {
	p, ok := source.(*player.Player)
	s := session.Get(p)
	if len(t.Players) > 1 && ok && !s.HasFlag(session.Admin) {
		output.Error(NoPermission)
		return
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

func (TeleportToPos) Allow(s cmd.Source) bool          { return !checkConsole(s) && checkStaff(s) }
func (TeleportToTarget) Allow(s cmd.Source) bool       { return !checkConsole(s) && checkStaff(s) }
func (TeleportTargetToTarget) Allow(s cmd.Source) bool { return checkStaff(s) }
func (TeleportTargetToPos) Allow(s cmd.Source) bool    { return checkStaff(s) }
