package form

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"strconv"
	"strings"
	"velvet/game"
	"velvet/utils"
)

// ffa is the ffa form.
type ffa struct {
	p *player.Player
}

// FFA creates a new FFA form.
func FFA(p *player.Player) form.Menu {
	var buttons []form.Button
	var games = []string{game.NoDebuff, game.Diamond, game.Build, game.God}
	for _, name := range games {
		g := game.Games[name]
		w, ok := utils.WorldMG.World(g.Name)
		var name string

		if ok {
			count := 0
			for _, v := range w.Entities() {
				if _, ok := v.(*player.Player); ok {
					count++
				}
			}
			name = g.DisplayName + "\n§l§3» §r§bCurrently playing: §9" + strconv.Itoa(count)
		} else {
			name = g.DisplayName + "\n§r§cOffline"
		}

		buttons = append(buttons, form.NewButton(name, g.FormData.ResourcePath))
	}
	return form.NewMenu(ffa{p}, "§l§aFree For All!").WithButtons(buttons...)
}

// Submit ...
func (f ffa) Submit(_ form.Submitter, pressed form.Button) {
	var g *game.Game
	name := strings.Split(pressed.Text, "\n")[0]
	for _, v := range game.Games {
		if v.DisplayName == name {
			g = v
			break
		}
	}

	if g == nil {
		f.p.Message(utils.Config.Message.ModeUnavailable)
		return
	}
	w, ok := utils.WorldMG.World(g.Name)
	if !ok {
		f.p.Message(utils.Config.Message.ModeUnavailable)
		return
	}

	w.AddEntity(f.p)
	f.p.SendTip("§7Welcome to " + g.Name + ".")
}
