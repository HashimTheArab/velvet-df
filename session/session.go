package session

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/df-mc/dragonfly/server/player/title"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"go.uber.org/atomic"
	"strconv"
	"strings"
	"sync"
	"time"
	"velvet/db"
	"velvet/game"
	"velvet/perm"
	"velvet/utils"
)

type Session struct {
	Player         *player.Player
	Flags          uint32
	clicks         []time.Time
	NetworkSession *session.Session
	Stats          *db.PlayerStats
	scoreTag       struct {
		healthText string
		cpsText    string
		osText     string
	}
	combat      combat
	combatMu    sync.Mutex
	cooldowns   cooldownMap
	cooldownsMu sync.Mutex
	scoreboard  *scoreboard.Scoreboard
	rank        *perm.Rank
	perms       atomic.Uint32
}

// OnJoin is called when the player joins.
func (s *Session) OnJoin() {
	s.cooldowns = map[cooldownType]*Cooldown{
		CooldownTypePearl: {length: time.Second * 15},
		CooldownTypeChat:  {length: time.Second * 3},
	}
	s.Stats = &db.PlayerStats{}
	s.NetworkSession = player_session(s.Player)
	s.Load()
	s.DefaultFlags()
	s.Player.SendTitle(title.New("§l§dVelvet").WithSubtitle("§bSeason 3 - Reformed"))
	s.UpdateScoreTag(true, true)
	s.SaveScoreboard()
	s.Player.EnableInstantRespawn()
	utils.OnlineCount.Add(1)
	for _, ses := range All() {
		ses.UpdateScoreboard(true, false)
	}
	game.DefaultKit(s.Player)
	go func() {
		for {
			s.UpdateScoreTag(true, true)
			if s.Combat().Tagged() {
				s.Combat().time--
				if !s.Combat().Tagged() {
					s.Player.Message("§aYou are no longer in combat!")
				}
			}
			if s.HasFlag(FlagVanished) {
				for _, e := range s.Player.World().Entities() {
					if pl, ok := e.(*player.Player); ok {
						if !Get(pl).HasFlag(FlagStaff) {
							pl.HideEntity(s.Player)
						}
					}
				}
				s.Player.SendTitle(title.New("").WithActionText("§aYou are vanished."))
			}
			time.Sleep(time.Second)
		}
	}()
}

// SetFlag sets a bit flag for the session, or unsets if the session already has the flag. A list of flags can be seen in flags.go
func (s *Session) SetFlag(flag uint32) {
	s.Flags ^= flag
}

// HasFlag returns whether the session has a specified bitflag.
func (s *Session) HasFlag(flag uint32) bool {
	return s.Flags&flag > 0
}

// Staff returns true if a player is a staff member.
func (s *Session) Staff() bool {
	if s.Rank() != nil && s.Rank().Name == perm.Admin || s.Rank().Name == perm.Mod {
		return true
	}
	xuid := s.Player.XUID()
	for _, v := range utils.Config.Staff.Admins {
		if v == xuid {
			return true
		}
	}
	for _, v := range utils.Config.Staff.Mods {
		if v == xuid {
			return true
		}
	}
	return false
}

// Mod returns true if a player is a moderator.
func (s *Session) Mod() bool {
	if s.Rank() != nil && s.Rank().Name == perm.Mod {
		return true
	}
	xuid := s.Player.XUID()
	for _, v := range utils.Config.Staff.Mods {
		if v == xuid {
			return true
		}
	}
	return false
}

// Admin returns true if the player is an Admin.
func (s *Session) Admin() bool {
	if s.Rank() != nil && s.Rank().Name == perm.Admin {
		return true
	}
	xuid := s.Player.XUID()
	for _, v := range utils.Config.Staff.Admins {
		if v == xuid {
			return true
		}
	}
	return false
}

// DefaultFlags will set the default bitflags for the session.
func (s *Session) DefaultFlags() {
	if s.Staff() {
		s.SetFlag(FlagStaff)
		staff[s.Player.Name()] = s
		if s.Admin() {
			s.SetFlag(FlagAdmin)
		}
	} else {
		if s.Rank() == nil {
			s.SetFlag(FlagHasChatCD)
		}
	}
}

// Click adds a click to the user's click history.
// Thanks, Tal, for cps!
func (s *Session) Click() {
	s.clicks = append(s.clicks, time.Now())
	if len(s.clicks) > 49 {
		s.clicks = s.clicks[1:]
	}
	s.Player.SendTip("§6CPS §b" + strconv.Itoa(s.CPS()))
	s.UpdateScoreTag(false, true)
}

// CPS returns the user's current click per second.
func (s *Session) CPS() int {
	var clicks int
	for _, past := range s.clicks {
		if time.Since(past) <= time.Second {
			clicks++
		}
	}
	return clicks
}

// TeleportToSpawn will teleport the player to the server spawn.
func (s *Session) TeleportToSpawn() {
	utils.Srv.World().AddEntity(s.Player)
}

func (s *Session) UpdateScoreTag(health, cps bool) {
	var tag string
	if health {
		hp := int(s.Player.Health())
		s.scoreTag.healthText = "§a" + strings.Repeat("|", hp)
		if s.Player.Health() < s.Player.MaxHealth() {
			s.scoreTag.healthText += "§c" + strings.Repeat("|", int(s.Player.MaxHealth())-hp)
		}
	}
	tag = s.scoreTag.healthText
	if cps {
		s.scoreTag.cpsText = " §bCPS " + strconv.Itoa(s.CPS())
	}
	tag += s.scoreTag.cpsText
	s.Player.SetScoreTag(tag)
}

func (s *Session) SaveScoreboard() {
	s.scoreboard = scoreboard.New("§l§6Velvet")
	s.Stats.Mutex.Lock()
	lines := []string{
		"§6Name: §b" + s.Player.Name(),
		"§6Online: §b" + utils.OnlineCount.String(),
		"§6K: §b" + strconv.Itoa(int(s.Stats.Kills)) + " §6D: §b" + strconv.Itoa(int(s.Stats.Deaths)),
		"§bvelvetpractice.tk",
	}
	s.Stats.Mutex.Unlock()
	_, _ = s.scoreboard.WriteString(strings.Join(lines, "\n"))
	s.Player.SendScoreboard(s.scoreboard)
}

func (s *Session) UpdateScoreboard(online, kd bool) {
	if online {
		_ = s.scoreboard.Set(1, "§6Online: §b"+utils.OnlineCount.String())
	}
	if kd {
		_ = s.scoreboard.Set(2, "§6K: §b"+strconv.Itoa(int(s.Stats.Kills))+" §6D: §b"+strconv.Itoa(int(s.Stats.Deaths)))
	}
	s.Player.SendScoreboard(s.scoreboard)
}

func (s *Session) Vanish() {
	p := s.Player
	s.SetFlag(FlagVanished)
	if s.HasFlag(FlagVanished) {
		p.SetNameTag(vanishPrefix + p.NameTag())
		p.Message("§aYou are now vanished!")
		AllStaff().Whisper("%v vanished", p.Name())
		p.SetGameMode(vanish{})
		for _, e := range p.World().Entities() {
			if pl, ok := e.(*player.Player); ok {
				if !Get(pl).HasFlag(FlagStaff) {
					pl.HideEntity(p)
				}
			}
		}
	} else {
		p.SetNameTag(strings.TrimPrefix(p.NameTag(), vanishPrefix))
		p.Message("§cYou are no longer vanished.")
		AllStaff().Whisper("%v unvanished", p.Name())
		p.SetGameMode(world.GameModeSurvival)
		for _, e := range p.World().Entities() {
			if pl, ok := e.(*player.Player); ok {
				pl.ShowEntity(p)
			}
		}
	}
}

func (s *Session) AddKills(kills uint32) {
	s.Stats.AddKills(kills)
	s.UpdateScoreboard(false, true)
}

func (s *Session) AddDeaths(deaths uint32) {
	s.Stats.AddDeaths(deaths)
	s.UpdateScoreboard(false, true)
}

// KDR returns the formatted kill-death ratio of the player.
func (s *Session) KDR() string {
	if s.Stats.GetDeaths() == 0 || s.Stats.GetKills() == 0 {
		return "0.0"
	}
	return strconv.FormatFloat(float64(s.Stats.GetKills()/s.Stats.GetDeaths()), 'f', 2, 32)
}

func (s *Session) Combat() *combat {
	s.combatMu.Lock()
	defer s.combatMu.Unlock()
	return &s.combat
}

func (s *Session) Cooldowns() *cooldownMap {
	s.cooldownsMu.Lock()
	defer s.cooldownsMu.Unlock()
	return &s.cooldowns
}

// Rank will return the rank of the session.
func (s *Session) Rank() *perm.Rank {
	return s.rank
}

// SetRank will set the rank for a session and return the new rank.
func (s *Session) SetRank(rank *perm.Rank) *perm.Rank {
	s.rank = rank
	return s.rank
}

// SetPerms will set the permissions of a session.
func (s *Session) SetPerms(perms uint32) {
	s.perms.Store(perms)
}

// SetPerm will set or remove a specific permission for a session
func (s *Session) SetPerm(perm perm.Permission) {
	s.perms.Store(s.Perms() ^ uint32(perm))
}

// HasPerm will return if a session has a permission.
func (s *Session) HasPerm(perm perm.Permission) bool {
	return s.Perms()&uint32(perm) > 0
}

// Perms will return the permissions of a session.
func (s *Session) Perms() uint32 {
	return s.perms.Load()
}
