package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"velvet/db"
	"velvet/perm"
	"velvet/session"
	"velvet/utils"
)

type SetRank struct {
	Sub     setRank
	Targets []cmd.Target `name:"target"`
	Rank    ranks        `name:"rank"`
}

type RemoveRank struct {
	Sub     removeRank
	Targets []cmd.Target `name:"target"`
}

type SetRankOffline struct {
	Sub    setRank
	Target string `name:"target"`
	Rank   ranks  `name:"rank"`
}

type RemoveRankOffline struct {
	Sub    removeRank
	Target string `name:"target"`
}

func (t SetRank) Run(_ cmd.Source, output *cmd.Output) {
	if len(t.Targets) > 1 {
		output.Error("§cYou can only set the rank of one player at a time.")
		return
	}
	if p, ok := t.Targets[0].(*player.Player); ok {
		s := session.Get(p)
		rank := perm.GetRank(string(t.Rank))
		if rank == nil {
			output.Error("§cRank not found. Contact the owner immediately.")
			return
		}
		s.SetRank(rank)
		db.SetRank(s.XUID, rank.Name)
		output.Printf(utils.Config.Rank.Set, p.Name(), rank.Name)
	}
}

func (t RemoveRank) Run(_ cmd.Source, output *cmd.Output) {
	if len(t.Targets) > 1 {
		output.Error("§cYou can only remove the rank of one player at a time.")
		return
	}
	if p, ok := t.Targets[0].(*player.Player); ok {
		s := session.Get(p)
		if s.Rank() == nil {
			output.Error("§cThat player does not have a rank.")
			return
		}
		s.SetRank(nil)
		db.SetRank(s.XUID, "")
		output.Printf(utils.Config.Rank.Removed, p.Name())
	}
}

func (t SetRankOffline) Run(_ cmd.Source, output *cmd.Output) {
	rank := perm.GetRank(string(t.Rank))
	if rank == nil {
		output.Error("§cRank not found. Contact the owner immediately.")
		return
	}
	if !db.Registered(t.Target) {
		output.Error("§cThat player has never joined the server.")
		return
	}
	db.SetRank(t.Target, rank.Name)
	output.Printf(utils.Config.Rank.Set, t.Target, rank.Name)
}

func (t RemoveRankOffline) Run(_ cmd.Source, output *cmd.Output) {
	if !db.Registered(t.Target) {
		output.Error("§cThat player has never joined the server.")
		return
	}
	db.SetRank(t.Target, "")
	output.Printf(utils.Config.Rank.Removed, t.Target)
}

type setRank string
type removeRank string

func (setRank) SubName() string    { return "set" }
func (removeRank) SubName() string { return "remove" }

type ranks string

func (ranks) Type() string { return "Rank" }
func (ranks) Options(cmd.Source) []string {
	r := perm.Ranks()
	var rankList []string
	for name, _ := range r {
		rankList = append(rankList, name)
	}
	return rankList
}

func (SetRank) Allow(s cmd.Source) bool           { return checkAdmin(s) || checkConsole(s) }
func (RemoveRank) Allow(s cmd.Source) bool        { return checkAdmin(s) || checkConsole(s) }
func (SetRankOffline) Allow(s cmd.Source) bool    { return checkAdmin(s) || checkConsole(s) }
func (RemoveRankOffline) Allow(s cmd.Source) bool { return checkAdmin(s) || checkConsole(s) }
