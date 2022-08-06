package utils

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/player"
)

// CountPots returns the amount of strong healing splash potions the player has in their inventory.
func CountPots(p *player.Player) uint {
	var pots uint = 0
	for _, i := range p.Inventory().Items() {
		if p, ok := i.Item().(item.SplashPotion); ok && p.Type == potion.StrongHealing() {
			pots++
		}
	}
	return pots
}
