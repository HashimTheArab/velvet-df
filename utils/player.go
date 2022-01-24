package utils

import "github.com/df-mc/dragonfly/server/player"

// CountPots returns the amount of strong healing splash potions the player has in their inventory.
func CountPots(p *player.Player) uint {
	var pots uint = 0
	for _, i := range p.Inventory().Items() {
		name, meta := i.Item().EncodeItem()
		if name == "minecraft:splash_potion" && meta == 22 {
			pots++
		}
	}
	return pots
}
