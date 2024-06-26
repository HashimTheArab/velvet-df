package game

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"time"
)

const (
	ArenaItemNbt = iota
)

func nodebuff_kit(p *player.Player) {
	clear(p)
	name := "§l§9Nodebuff"

	unbreaking := item.NewEnchantment(enchantment.Unbreaking{}, 3)
	prot := item.NewEnchantment(enchantment.Protection{}, 1)
	_ = p.Inventory().SetItem(0, item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithCustomName(name).WithEnchantments(unbreaking, item.NewEnchantment(enchantment.Sharpness{}, 1)))
	_ = p.Inventory().SetItem(1, item.NewStack(item.EnderPearl{}, 16))
	_, _ = p.Inventory().AddItem(item.NewStack(item.SplashPotion{Type: potion.StrongHealing()}, 34))
	p.Armour().SetBoots(item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithCustomName(name).WithEnchantments(unbreaking, prot))
	p.Armour().SetLeggings(item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1).WithCustomName(name).WithEnchantments(unbreaking, prot))
	p.Armour().SetChestplate(item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1).WithCustomName(name).WithEnchantments(unbreaking, prot))
	p.Armour().SetHelmet(item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithCustomName(name).WithEnchantments(unbreaking, prot))
	p.AddEffect(effect.New(effect.Speed{}, 1, time.Hour*10).WithoutParticles())
}

func diamond_kit(p *player.Player) {
	clear(p)
	name := "§l§3Diamond"
	_ = p.Inventory().SetItem(0, item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithCustomName(name))
	_ = p.Inventory().SetItem(1, item.NewStack(item.Bow{}, 1))
	_ = p.Inventory().SetItem(9, item.NewStack(item.Arrow{}, 1))
	_ = p.Inventory().SetItem(8, item.NewStack(item.Spyglass{}, 1))
	p.Armour().SetBoots(item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithCustomName(name))
	p.Armour().SetLeggings(item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1).WithCustomName(name))
	p.Armour().SetChestplate(item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1).WithCustomName(name))
	p.Armour().SetHelmet(item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithCustomName(name))
}

func build_kit(p *player.Player) {
	clear(p)
	name := "§l§6Build"
	unbreaking := item.NewEnchantment(enchantment.Unbreaking{}, 3)
	prot := item.NewEnchantment(enchantment.Protection{}, 1)

	_ = p.Inventory().SetItem(0, item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithCustomName(name).WithEnchantments(unbreaking, item.NewEnchantment(enchantment.Sharpness{}, 1)))
	_ = p.Inventory().SetItem(1, item.NewStack(item.Bow{}, 1).WithCustomName(name).WithEnchantments(unbreaking))
	_ = p.Inventory().SetItem(2, item.NewStack(item.GoldenApple{}, 10))
	_ = p.Inventory().SetItem(3, item.NewStack(item.EnderPearl{}, 10))
	_ = p.Inventory().SetItem(4, item.NewStack(item.Axe{Tier: item.ToolTierDiamond}, 1).WithCustomName(name).WithEnchantments(unbreaking, item.NewEnchantment(enchantment.Efficiency{}, 3)))
	_ = p.Inventory().SetItem(5, item.NewStack(item.Pickaxe{Tier: item.ToolTierDiamond}, 1).WithCustomName(name).WithEnchantments(unbreaking, item.NewEnchantment(enchantment.Efficiency{}, 3)))
	_ = p.Inventory().SetItem(6, item.NewStack(block.Cobblestone{}, 64))
	//_ = p.Inventory().SetItem(7, item.NewStack(block.Planks{}, 64))
	_ = p.Inventory().SetItem(7, item.NewStack(item.Stick{}, 1).WithCustomName(text.Colourf("<purple>Magic Stick</purple>")))
	_ = p.Inventory().SetItem(8, item.NewStack(item.Spyglass{}, 1))
	_ = p.Inventory().SetItem(9, item.NewStack(item.Arrow{}, 16))

	p.Armour().SetBoots(item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithCustomName(name).WithEnchantments(unbreaking, prot))
	p.Armour().SetLeggings(item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1).WithCustomName(name).WithEnchantments(unbreaking, prot))
	p.Armour().SetChestplate(item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1).WithCustomName(name).WithEnchantments(unbreaking, prot))
	p.Armour().SetHelmet(item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithCustomName(name).WithEnchantments(unbreaking, prot))
}

func gfight_kit(p *player.Player) {
	clear(p)
	name := "§l§4GKit"
	unbreaking := item.NewEnchantment(enchantment.Unbreaking{}, 3)

	_ = p.Inventory().SetItem(0, item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithValue("gsword", struct{}{}).WithCustomName(name).WithEnchantments(unbreaking).WithLore(
		"§6Kaboom III", "§eZeus III", "§eBleed III", "§2Hades III", "§2Poison II", "§2Lifesteal II", "§aOOF II",
	))
	_ = p.Inventory().SetItem(1, item.NewStack(item.EnchantedApple{}, 10))

	armour := [4]item.Stack{
		item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1),
		item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1),
		item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1),
		item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1),
	}
	for k, v := range armour {
		armour[k] = v.WithCustomName(name).WithLore("§6Overlord II", "§2Adrenaline I", "§aScorch V").WithEnchantments(unbreaking)
	}
	armour[3] = armour[3].WithLore("§6Overlord II", "§eGears I", "§2Adrenaline I", "§aScorch V")
	p.Armour().Set(armour[0], armour[1], armour[2], armour[3])
	p.AddEffect(effect.New(effect.Speed{}, 1, time.Hour*24))
}

func DefaultKit(p *player.Player) {
	_ = p.Inventory().SetItem(0, item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithCustomName("§aArenas!").WithValue("tool", 0))
}

func clear(p *player.Player) {
	p.Inventory().Clear()
	p.Armour().Clear()
	for _, e := range p.Effects() {
		p.RemoveEffect(e.Type())
	}
}
