package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"strings"
	"velvet/utils"
)

type GameModeType string

type GameMode struct {
	GameMode GameModeType `name:"mode"`
	Player   []cmd.Target `optional:"" name:"target"`
}

func (t GameMode) Run(source cmd.Source, output *cmd.Output) {
	p := source.(*player.Player)

	var gm world.GameMode
	switch strings.ToLower(string(t.GameMode)) {
	case "survival", "0", "s":
		gm = world.GameModeSurvival
	case "creative", "1", "c":
		gm = world.GameModeCreative
	case "adventure", "2", "a":
		gm = world.GameModeAdventure
	case "spectator", "3", "sp":
		gm = world.GameModeSpectator
	default:
		return
	}

	if len(t.Player) > 0 {
		if target, ok := t.Player[0].(*player.Player); ok {
			target.SetGameMode(gm)
			target.Messagef(utils.Config.Message.GameModeSetByPlayer, p.Name(), t.GameMode)
			output.Printf(utils.Config.Message.GameModeSetOther, target.Name(), t.GameMode)
			return
		}
	}

	p.SetGameMode(gm)
	output.Printf(utils.Config.Message.GameModeSetBySelf, t.GameMode)
}

func (GameModeType) Type() string {
	return "GameMode"
}

func (GameModeType) Options(cmd.Source) []string {
	return []string{
		"survival", "0", "s",
		"creative", "1", "c",
		"adventure", "2", "a",
		"spectator", "3", "sp",
	}
}

func (GameMode) Allow(s cmd.Source) bool { return checkAdmin(s) }
