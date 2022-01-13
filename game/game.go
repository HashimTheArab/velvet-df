package game

import (
	"github.com/df-mc/dragonfly/server/player"
	"velvet/utils"
)

type Game struct {
	Name        string
	DisplayName string
	WorldName   string
	Knockback   kb
	FormData    form
	Kit         func(p *player.Player)
}

var Games = map[string]*Game{
	NoDebuff: {
		Name:        "NoDebuff",
		DisplayName: "§l§9Nodebuff",
		WorldName:   utils.Config.World.NoDebuff,
		Knockback:   kb{0.398, 0.405},
		FormData: form{
			ResourcePath: "textures/items/potion_bottle_splash_heal",
		},
		Kit: nodebuff_kit,
	},
	Diamond: {
		Name:        "Diamond",
		DisplayName: "§l§3Diamond",
		WorldName:   utils.Config.World.Diamond,
		Knockback:   defaultKB,
		FormData: form{
			ResourcePath: "textures/items/diamond_sword",
		},
		Kit: diamond_kit,
	},
	Build: {
		Name:        "Build",
		DisplayName: "§l§6Build",
		WorldName:   utils.Config.World.Diamond,
		Knockback:   defaultKB,
		FormData: form{
			ResourcePath: "textures/items/diamond_pickaxe",
		},
		Kit: build_kit,
	},
}

func Get(name string) *Game {
	return Games[name]
}

func FromWorld(name string) *Game {
	switch name {
	case utils.Config.World.NoDebuff:
		return Get(NoDebuff)
	}
	return nil
}
