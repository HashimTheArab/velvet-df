package item

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	_ "unsafe"
	"velvet/session"
)

var items = map[world.Item]item.Usable{
	item.EnderPearl{}:   Pearl{},
	item.SplashPotion{}: Potion{},
}

// thank you so much tal
func Override(s *session.Session, ctx *event.Context) {
	p := s.Player

	held, left := p.HeldItems()

	// Clear meta that may prevent the item from being overridden.
	name, _ := held.Item().EncodeItem()
	itemToCheck, _ := world.ItemByName(name, 0)
	if replacement, ok := items[itemToCheck]; ok && !ctx.Cancelled() {
		it := held.Item()
		w := p.World()
		if p.HasCooldown(it) {
			return
		}

		ctx.Cancel()
		if cooldown, ok := it.(item.Cooldown); ok {
			p.SetCooldown(it, cooldown.Cooldown())
		}

		ctx := player_useContext(p)
		if replacement.Use(w, p, ctx) {
			// We only swing the player's arm if the item held actually does something. If it doesn't, there is no
			// reason to swing the arm.
			p.SwingArm()

			p.SetHeldItems(player_subtractItem(p, player_damageItem(p, held, ctx.Damage), ctx.CountSub), left)
			player_addNewItem(p, ctx)
		}
	}
}

//go:linkname player_subtractItem github.com/df-mc/dragonfly/server/player.(*Player).subtractItem
//noinspection ALL
func player_subtractItem(*player.Player, item.Stack, int) item.Stack

//go:linkname player_damageItem github.com/df-mc/dragonfly/server/player.(*Player).damageItem
//noinspection ALL
func player_damageItem(*player.Player, item.Stack, int) item.Stack

//go:linkname player_addNewItem github.com/df-mc/dragonfly/server/player.(*Player).addNewItem
//noinspection ALL
func player_addNewItem(*player.Player, *item.UseContext)

//go:linkname player_useContext github.com/df-mc/dragonfly/server/player.(*Player).useContext
//noinspection ALL
func player_useContext(*player.Player) *item.UseContext
