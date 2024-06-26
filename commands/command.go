package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"velvet/console"
	"velvet/session"
)

const (
	NoPermission   = "§4You do not have permission to use this command!"
	PlayerNotFound = "§cPlayer not found."
)

func init() {
	for _, command := range []cmd.Command{
		cmd.New("gamemode", "§bChange your gamemode", []string{"gm"}, GameMode{}),
		cmd.New("_wand", "§bGet a World Edit wand", nil, Wand{}),
		cmd.New("_palette", "§bSet a block palette for world edit", nil, PaletteSet{}, PaletteSave{}, PaletteDelete{}),
		cmd.New("_fill", "§bFill an area", nil, Fill{}),
		cmd.New("_fillblock", "§bFill an area with a specific block", nil, FillBlock{}),
		cmd.New("teleport", "§bTeleport to another player", []string{"tp"}, TeleportToPos{}, TeleportToTarget{}, TeleportTargetToTarget{}, TeleportTargetToPos{}),
		cmd.New("build", "§bUse builder mode", nil, Build{}),
		cmd.New("world", "§bManage worlds", nil, WorldTeleport{}, WorldList{}, WorldCreate{}),
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
		cmd.New("status", "§cView the status of the server", nil, Status{}),
		cmd.New("clearbuild", "§aClear the build ffa arena", []string{"cb"}, ClearBuild{}),
		cmd.New("vanish", "§bHide yourself from other players", []string{"v"}, Vanish{}),
		cmd.New("alias", "§aView the alts of a player", nil, Alias{}, AliasOffline{}),
		cmd.New("rank", "§cManage ranks", nil, SetRank{}, RemoveRank{}, SetRankOffline{}, RemoveRankOffline{}),
		cmd.New("ping", "§bView the ping of yourself or another player", []string{"ms"}, Ping{}),
		cmd.New("whitelist", "§cManage the server whitelist", nil, WhitelistToggle{}, WhitelistAdd{}, WhitelistRemove{}),
		cmd.New("say", "§0Broadcast a message to everyone", nil, Say{}),
		cmd.New("list", "§aView all online players", nil, List{}),
		cmd.New("stats", "§bView a player's stats", nil, StatsOnline{}, StatsOffline{}),
		//cmd.StartNew("kill", "§bKill another player", nil, Kill{}),
	} {
		cmd.Register(command)
	}
}

func checkStaff(s cmd.Source) bool {
	return checkPerms(s, session.FlagStaff)
}

func checkAdmin(s cmd.Source) bool {
	return checkPerms(s, session.FlagAdmin)
}

func checkBuilder(s cmd.Source) bool {
	return checkPerms(s, session.FlagBuilder)
}

func checkConsole(s cmd.Source) bool {
	_, ok := s.(*console.CommandSender)
	return ok
}

func checkPerms(s cmd.Source, flag uint32) bool {
	p, ok := s.(*player.Player)
	ses := session.Get(p)
	return ses != nil && ok && ses.HasFlag(flag)
}
