package entity

import (
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/damage"
	"github.com/df-mc/dragonfly/server/entity/physics"
	"github.com/df-mc/dragonfly/server/entity/physics/trace"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"image/color"
	vs "velvet/session"
)

// Pearl is a copy of an ender pearl with some edits.
type Pearl struct {
	transform
	yaw, pitch float64

	age   int
	close bool

	owner world.Entity

	c *entity.ProjectileComputer
}

// Thank you Tal!

// NewEnderPearl ...
func NewEnderPearl(pos, vel mgl64.Vec3, yaw, pitch float64, owner world.Entity) *Pearl {
	e := &Pearl{
		yaw:   yaw,
		pitch: pitch,
		c: &entity.ProjectileComputer{MovementComputer: &entity.MovementComputer{
			Gravity:           0.085, // 0.085
			Drag:              0.01,
			DragBeforeGravity: true,
		}},
		owner: owner,
	}
	e.transform = newTransform(e, pos)
	e.vel = vel
	return e
}

// Name ...
func (e *Pearl) Name() string {
	return "Ender Pearl"
}

// EncodeEntity ...
func (e *Pearl) EncodeEntity() string {
	return "minecraft:ender_pearl"
}

// Scale ...
func (e *Pearl) Scale() float64 {
	return 0.6
}

// AABB ...
func (e *Pearl) AABB() physics.AABB {
	return physics.NewAABB(mgl64.Vec3{-0.125, 0, -0.125}, mgl64.Vec3{0.125, 0.25, 0.125})
}

// Rotation ...
func (e *Pearl) Rotation() (float64, float64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.yaw, e.pitch
}

// Tick ...
func (e *Pearl) Tick(w *world.World, current int64) {
	if e.close {
		_ = e.Close()
		return
	}
	e.mu.Lock()
	m, result := e.c.TickMovement(e, e.pos, e.vel, e.yaw, e.pitch, e.ignores)
	e.yaw, e.pitch = m.Rotation()
	e.pos, e.vel = m.Position(), m.Velocity()
	e.mu.Unlock()
	w.AddParticle(m.Position(), particle.Flame{Colour: color.RGBA{G: 255, B: 255}})

	e.age++
	m.Send()

	if m.Position()[1] < float64(w.Range()[0]) && current%10 == 0 {
		e.close = true
		return
	}

	if result != nil {
		if r, ok := result.(trace.EntityResult); ok {
			if l, ok := r.Entity().(entity.Living); ok {
				if _, vulnerable := l.Hurt(0.0, damage.SourceEntityAttack{Attacker: e}); vulnerable {
					l.KnockBack(m.Position(), 0.45, 0.3608)
				}
			}
		}

		if owner := e.Owner(); owner != nil {
			if user, ok := owner.(*player.Player); ok {
				s := vs.Get(user)
				w.PlaySound(user.Position(), sound.EndermanTeleport{})

				yaw, pitch := user.Rotation()
				onGround := user.OnGround()
				translatedPos := vec64To32(m.Position()).Add(mgl32.Vec3{0, 1.621, 0})
				s.WritePacket(&packet.MovePlayer{
					EntityRuntimeID: 1,
					Position:        translatedPos,
					OnGround:        onGround,
					Yaw:             float32(yaw),
					HeadYaw:         float32(yaw),
					Pitch:           float32(pitch),
					Mode:            packet.MoveModeTeleport,
				})
				for _, v := range w.Viewers(m.Position()) {
					session_writePacket(s.NetworkSession, &packet.MovePlayer{
						EntityRuntimeID: session_entityRuntimeID(v.(*session.Session), user),
						Position:        translatedPos,
						OnGround:        onGround,
						Yaw:             float32(yaw),
						HeadYaw:         float32(yaw),
						Pitch:           float32(pitch),
						Mode:            packet.MoveModeNormal,
					})
				}

				w.AddParticle(m.Position(), particle.EndermanTeleportParticle{})
				w.PlaySound(m.Position(), sound.EndermanTeleport{})
			}
		}

		e.close = true
	}
}

// ignores returns whether the ender pearl should ignore collision with the entity passed.
func (e *Pearl) ignores(otherEntity world.Entity) bool {
	_, ok := otherEntity.(entity.Living)
	return !ok || otherEntity == e || (e.age < 5 && otherEntity == e.owner)
}

// Owner ...
func (e *Pearl) Owner() world.Entity {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.owner
}

// Own ...
func (e *Pearl) Own(owner world.Entity) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.owner = owner
}

// vec64To32 converts a mgl64.Vec3 to a mgl32.Vec3.
func vec64To32(vec3 mgl64.Vec3) mgl32.Vec3 {
	return mgl32.Vec3{float32(vec3[0]), float32(vec3[1]), float32(vec3[2])}
}
