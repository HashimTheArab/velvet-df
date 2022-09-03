package game

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"math"
	"math/rand"
	"strings"
	"velvet/utils"
)

type Game struct {
	Name         string
	DisplayName  string
	WorldName    string
	Knockback    kb
	DeathMessage string
	FormData     form
	Kit          func(p *player.Player)
}

var Games = map[string]*Game{
	NoDebuff: {
		Name:         "NoDebuff",
		DisplayName:  "§l§9Nodebuff",
		WorldName:    utils.Config.World.NoDebuff,
		Knockback:    kb{0.398, 0.405},
		DeathMessage: utils.Config.DeathMessage.NoDebuff,
		FormData: form{
			ResourcePath: "textures/items/potion_bottle_splash_heal",
		},
		Kit: nodebuff_kit,
	},
	Diamond: {
		Name:        "Diamond",
		DisplayName: "§l§3Diamond",
		WorldName:   utils.Config.World.Diamond,
		FormData: form{
			ResourcePath: "textures/items/diamond_sword",
		},
		Kit: diamond_kit,
	},
	Build: {
		Name:        "Build",
		DisplayName: "§l§6Build",
		WorldName:   utils.Config.World.Build,
		FormData: form{
			ResourcePath: "textures/items/diamond_pickaxe",
		},
		Kit: build_kit,
	},
	God: {
		Name:        "God",
		DisplayName: "§l§4God",
		WorldName:   utils.Config.World.God,
		FormData: form{
			ResourcePath: "textures/items/enchanted_apple",
		},
		Kit: gfight_kit,
	},
}

func init() {
	for _, v := range Games {
		if v.DeathMessage == "" {
			v.DeathMessage = utils.Config.DeathMessage.Default
		}
		if v.Knockback == (kb{}) {
			v.Knockback = defaultKB
		}
	}
}

func Get(name string) *Game {
	return Games[name]
}

func FromWorld(name string) *Game {
	switch strings.ToLower(name) {
	case utils.Config.World.NoDebuff:
		return Get(NoDebuff)
	case utils.Config.World.Diamond:
		return Get(Diamond)
	case utils.Config.World.Build:
		return Get(Build)
	case utils.Config.World.God:
		return Get(God)
	}
	return nil
}

func (g *Game) BroadcastDeathMessage(p, t *player.Player) {
	randomMessage := utils.Config.DeathMessage.List[rand.Intn(len(utils.Config.DeathMessage.List))]
	switch g.DeathMessage {
	case utils.Config.DeathMessage.NoDebuff:
		_, _ = fmt.Fprintf(chat.Global, utils.Config.DeathMessage.NoDebuff, p.Name(), utils.CountPots(p), randomMessage, t.Name(), math.Round(t.Health()*100)/100, utils.CountPots(t))
	case utils.Config.DeathMessage.Default:
		_, _ = fmt.Fprintf(chat.Global, utils.Config.DeathMessage.Default, p.Name(), randomMessage, t.Name(), math.Round(t.Health()*100)/100)
	}
}
