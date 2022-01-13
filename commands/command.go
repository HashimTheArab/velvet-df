package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
)

const (
	NoPermission    = "§4You do not have permission to use this command!"
	PlayerNotOnline = "§c%v is not online."
)

func init() {
	for _, command := range []cmd.Command{
		cmd.New("gamemode", "§bChange your gamemode", []string{"gm"}, GameMode{}),
		cmd.New("teleport", "§bTeleport to another player", []string{"tp"}, TeleportToPos{}, TeleportToTarget{}, TeleportTargetToTarget{}, TeleportTargetToPos{}),
		cmd.New("build", "§bUse builder mode", nil, Build{}),
		cmd.New("world", "§bManage worlds", nil, WorldTeleport{}),
		cmd.New("/worldedit", "§bManage world edit", []string{"we"}, WorldEdit{}),
		cmd.New("newplayer", "§bSpawn a fake player", []string{"np"}, NewPlayer{}),
		cmd.New("effect", "§bApply an effect to yourself or another player", nil, Effect{}),
		cmd.New("kick", "§aKick a player from the server", nil, Kick{}),
		cmd.New("ban", "§aBan a player from the server", nil, Ban{}),
	} {
		cmd.Register(command)
	}
}
