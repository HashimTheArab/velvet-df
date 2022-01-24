package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"velvet/console"
	"velvet/session"
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
		//cmd.New("/worldedit", "§bManage world edit", []string{"we"}, WorldEdit{}),
		cmd.New("newplayer", "§bSpawn a fake player", []string{"np"}, NewPlayer{}),
		cmd.New("effect", "§bApply an effect to yourself or another player", nil, Effect{}),
		cmd.New("kick", "§aKick a player from the server", nil, Kick{}),
		cmd.New("ban", "§aBan a player from the server", nil, Ban{}, BanOffline{}),
		cmd.New("blacklist", "§aBlacklist a player from the server", nil, Blacklist{}, BlacklistOffline{}),
		cmd.New("unban", "§aUnban a player from the server", nil, BanLift{}),
		cmd.New("baninfo", "§aView the ban information of a player", nil, BanInfo{}),
		cmd.New("spawnpoint", "§cSet the default spawn of a world", []string{"defaultspawn"}, SpawnPoint{}),
		cmd.New("spawn", "§bTeleport to spawn", []string{"hub"}, Spawn{}),
		cmd.New("clear", "§cClear yours or another players inventory", nil, Clear{}),
		cmd.New("transfer", "§cTransfer a player to another server", nil, Transfer{}),
		cmd.New("tell", "§bSend a message to another player", []string{"w"}, Tell{}),
		cmd.New("time", "§bChange the time of the world you're in", nil, TimeSet{}),
		//cmd.New("kill", "§bKill another player", nil, Kill{}),
	} {
		cmd.Register(command)
	}
}

func checkStaff(s cmd.Source) bool {
	return checkPerms(s, session.Staff)
}

func checkAdmin(s cmd.Source) bool {
	return checkPerms(s, session.Admin)
}

func checkConsole(s cmd.Source) bool {
	_, ok := s.(console.CommandSender)
	return ok
}

func checkPerms(s cmd.Source, flag uint32) bool {
	_, ok := s.(console.CommandSender)
	if ok {
		return true
	}
	p, ok := s.(*player.Player)
	if !ok || !session.Get(p).HasFlag(flag) {
		return false
	}
	return true
}
