package item

import (
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	ve "velvet/entity"
)

// Pearl is an edited usable for ender pearls.
type Pearl struct{}

// Use ...
func (v Pearl) Use(w *world.World, user item.User, ctx *item.UseContext) bool {
	yaw, pitch := user.Rotation()
	e := ve.NewEnderPearl(entity.EyePosition(user), entity.DirectionVector(user).Mul(2.3), yaw, pitch, user)
	w.AddEntity(e)

	w.PlaySound(user.Position(), sound.ItemThrow{})
	ctx.SubtractFromCount(1)
	return true
}
