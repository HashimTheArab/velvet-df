package enchants

import (
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/damage"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/entity/healing"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"math/rand"
	"time"
	"velvet/session"
)

func Kaboom(a, t *player.Player, force, height *float64) {
	if rand.Intn(50) > 3 {
		return
	}
	t.World().AddParticle(t.Position(), particle.HugeExplosion{})
	t.World().PlaySound(t.Position(), sound.Explosion{})

	*force += 0.88
	*height += 0.88

	if t.Health() < 10 {
		t.Hurt(t.Health()-3.99, damage.SourceExplosion{}) // $enchantmentLevel * 1.33
	} else {
		t.Hurt(6, damage.SourceExplosion{}) // $enchantmentLevel * 2
	}
}

func Zeus(_, t *player.Player) {
	if rand.Intn(45) > 3 {
		return
	}
	t.World().AddEntity(entity.NewLightningWithDamage(t.Position(), 0, false, false))
	t.Hurt(6, damage.SourceLightning{})
}

func Bleed(_, t *player.Player) {
	if rand.Intn(40) > 3 {
		return
	}
	if s := session.Get(t); s != nil {
		s.StartBleeding()
	}
}

func Hades(_, t *player.Player) {
	if rand.Intn(50) > 3 {
		return
	}
	t.SetOnFire(time.Millisecond * 750)
	t.Hurt(3, damage.SourceFire{})
	t.World().AddParticle(t.Position().Add(mgl64.Vec3{float64(rand.Intn(11)-10) / 10, float64(rand.Intn(21)) / 10, float64(rand.Intn(11)-10) / 10}), particle.Flame{})
}

func Poison(_, t *player.Player) {
	if rand.Intn(40) > 2 {
		return
	}
	t.AddEffect(effect.New(effect.Poison{}, 1, time.Second*12))
}

func Lifesteal(a, t *player.Player) {
	if rand.Intn(50) > 2 {
		return
	}
	if t.Health() > 10 && a.Health() < 19 {
		t.Hurt(2, damage.SourceInstantDamageEffect{})
		a.Heal(2, healing.SourceInstantHealthEffect{})
	}
}
