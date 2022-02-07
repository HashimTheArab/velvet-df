package game

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/player"
)

const (
	ArenaItemNbt = iota
)

func nodebuff_kit(p *player.Player) {
	p.Inventory().Clear()
	p.Armour().Clear()
	name := "§l§9Nodebuff"

	unbreaking := enchantment.Unbreaking{}.WithLevel(3)
	_ = p.Inventory().SetItem(0, item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithCustomName(name).WithEnchantments(unbreaking, enchantment.Sharpness{}))
	_ = p.Inventory().SetItem(1, item.NewStack(item.EnderPearl{}, 16))
	_, _ = p.Inventory().AddItem(item.NewStack(item.SplashPotion{Type: potion.StrongHealing()}, 34))
	p.Armour().SetBoots(item.NewStack(item.Boots{Tier: item.ArmourTierDiamond}, 1).WithCustomName(name).WithEnchantments(unbreaking))
	p.Armour().SetLeggings(item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond}, 1).WithCustomName(name).WithEnchantments(unbreaking))
	p.Armour().SetChestplate(item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond}, 1).WithCustomName(name).WithEnchantments(unbreaking))
	p.Armour().SetHelmet(item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond}, 1).WithCustomName(name).WithEnchantments(unbreaking))
	//p.AddEffect(effect.New(effect.Speed{}, 1, time.Hour*3))
}

func diamond_kit(p *player.Player) {
	p.Inventory().Clear()
	p.Armour().Clear()
	name := "§l§3Diamond"
	_ = p.Inventory().SetItem(0, item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithCustomName(name))
	_ = p.Inventory().SetItem(1, item.NewStack(item.Bow{}, 1))
	_ = p.Inventory().SetItem(9, item.NewStack(item.Arrow{}, 1))
	_ = p.Inventory().SetItem(8, item.NewStack(item.Spyglass{}, 1))
	p.Armour().SetBoots(item.NewStack(item.Boots{Tier: item.ArmourTierDiamond}, 1).WithCustomName(name))
	p.Armour().SetLeggings(item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond}, 1).WithCustomName(name))
	p.Armour().SetChestplate(item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond}, 1).WithCustomName(name))
	p.Armour().SetHelmet(item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond}, 1).WithCustomName(name))
}

func build_kit(p *player.Player) {
	p.Inventory().Clear()
	p.Armour().Clear()
	name := "§l§6Build"
	unbreaking := enchantment.Unbreaking{}.WithLevel(3)

	_ = p.Inventory().SetItem(0, item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithCustomName(name).WithEnchantments(unbreaking, enchantment.Sharpness{}))
	_ = p.Inventory().SetItem(1, item.NewStack(item.Bow{}, 1).WithCustomName(name).WithEnchantments(unbreaking, enchantment.Unbreaking{}))
	_ = p.Inventory().SetItem(2, item.NewStack(item.GoldenApple{}, 10))
	_ = p.Inventory().SetItem(3, item.NewStack(item.EnderPearl{}, 10))
	_ = p.Inventory().SetItem(4, item.NewStack(item.Axe{Tier: item.ToolTierDiamond}, 1).WithCustomName(name).WithEnchantments(unbreaking, enchantment.Efficiency{}.WithLevel(3)))
	_ = p.Inventory().SetItem(5, item.NewStack(item.Pickaxe{Tier: item.ToolTierDiamond}, 1).WithCustomName(name).WithEnchantments(unbreaking, enchantment.Efficiency{}.WithLevel(3)))
	_ = p.Inventory().SetItem(6, item.NewStack(block.Cobblestone{}, 64))
	_ = p.Inventory().SetItem(7, item.NewStack(block.Planks{}, 64))
	_ = p.Inventory().SetItem(8, item.NewStack(item.Spyglass{}, 1))
	_ = p.Inventory().SetItem(9, item.NewStack(item.Arrow{}, 16))

	p.Armour().SetBoots(item.NewStack(item.Boots{Tier: item.ArmourTierDiamond}, 1).WithCustomName(name).WithEnchantments(unbreaking))
	p.Armour().SetLeggings(item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond}, 1).WithCustomName(name).WithEnchantments(unbreaking))
	p.Armour().SetChestplate(item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond}, 1).WithCustomName(name).WithEnchantments(unbreaking))
	p.Armour().SetHelmet(item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond}, 1).WithCustomName(name).WithEnchantments(unbreaking))
}

func DefaultKit(p *player.Player) {
	_ = p.Inventory().SetItem(0, item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithCustomName("§aArenas!").WithValue("tool", 0))
}
