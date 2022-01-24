package item

import (
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	ve "velvet/entity"
)

// Potion is an edited usable for splash potions.
type Potion struct{}

// Use ...
func (v Potion) Use(w *world.World, user item.User, ctx *item.UseContext) bool {
	held, _ := user.HeldItems()
	it := held.Item().(item.SplashPotion)

	yaw, pitch := user.Rotation()
	e := ve.NewSplashPotion(entity.EyePosition(user), entity.DirectionVector(user).Mul(0.5), yaw, pitch, it.Type, user)
	w.AddEntity(e)

	w.PlaySound(user.Position(), sound.ItemThrow{})
	ctx.SubtractFromCount(1)
	return true
}
