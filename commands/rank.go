package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"velvet/console"
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

func (t SetRank) Run(source cmd.Source, output *cmd.Output) {
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
		if canSetRank(p.Name(), source, output) {
			s.SetRank(rank)
			db.SetRank(p.XUID(), rank.Name)
			setRankFlags(s, rank.Name)
			output.Printf(utils.Config.Rank.Set, p.Name(), rank.Name)
		}
	}
}

func (t RemoveRank) Run(source cmd.Source, output *cmd.Output) {
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
		if canSetRank(p.Name(), source, output) {
			s.SetRank(nil)
			db.SetRank(p.XUID(), "")
			setRankFlags(s, "")
			output.Printf(utils.Config.Rank.Removed, p.Name())
		}
	}
}

func (t SetRankOffline) Run(source cmd.Source, output *cmd.Output) {
	rank := perm.GetRank(string(t.Rank))
	if rank == nil {
		output.Error("§cRank not found. Contact the owner immediately.")
		return
	}
	if !db.Registered(t.Target) {
		output.Error("§cThat player has never joined the server.")
		return
	}
	if canSetRank(t.Target, source, output) {
		db.SetRank(t.Target, rank.Name)
		output.Printf(utils.Config.Rank.Set, t.Target, rank.Name)
	}
}

func (t RemoveRankOffline) Run(source cmd.Source, output *cmd.Output) {
	if !db.Registered(t.Target) {
		output.Error("§cThat player has never joined the server.")
		return
	}
	if canSetRank(t.Target, source, output) {
		db.SetRank(t.Target, "")
		output.Printf(utils.Config.Rank.Removed, t.Target)
	}
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

func canSetRank(target string, s cmd.Source, output *cmd.Output) bool {
	if _, ok := s.(*console.CommandSender); !ok {
		p := s.(*player.Player)
		if db.IsStaff(target) && p.XUID() != utils.Config.Staff.Owner.XUID {
			output.Print(NoPermission)
			return false
		}
	}
	return true
}

func setRankFlags(s *session.Session, newRank string) {
	if s.HasFlag(session.FlagStaff) {
		if !perm.StaffRanks.Contains(newRank) {
			s.SetFlag(session.FlagStaff)
			session.RemoveStaff(s)
		}
		if s.HasFlag(session.FlagAdmin) && newRank != perm.Admin {
			s.SetFlag(session.FlagAdmin)
		}
		if s.HasFlag(session.FlagBuilder) && newRank != perm.Builder {
			s.SetFlag(session.FlagBuilder)
		}
	} else {
		if perm.StaffRanks.Contains(newRank) {
			s.SetFlag(session.FlagStaff)
		}
		switch newRank {
		case perm.Admin:
			if !s.HasFlag(session.FlagAdmin) {
				s.SetFlag(session.FlagAdmin)
			}
			session.AddStaff(s)
		case perm.Builder:
			if !s.HasFlag(session.FlagBuilder) {
				s.SetFlag(session.FlagBuilder)
			}
		case perm.Mod:
			session.AddStaff(s)
		case perm.Owner:
			if !s.HasFlag(session.FlagAdmin) {
				s.SetFlag(session.FlagAdmin)
			}
			session.AddStaff(s)
		}
	}
}
