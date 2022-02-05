package session

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/df-mc/dragonfly/server/player/title"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"strconv"
	"strings"
	"sync"
	"time"
	"velvet/db"
	"velvet/game"
	"velvet/utils"
)

type Session struct {
	Player         *player.Player
	Flags          uint32
	clicks         []time.Time
	NetworkSession *session.Session
	Stats          db.PlayerStats
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
}

// OnJoin is called when the player joins.
func (s *Session) OnJoin() {
	s.cooldowns = map[cooldownType]*Cooldown{
		CooldownTypePearl: {length: time.Second * 15},
		CooldownTypeChat:  {length: time.Second * 3},
	}
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
			}
			time.Sleep(time.Second)
		}
	}()
}

// SetFlag sets a bit flag for the session, or unsets if the session already has the flag. A list of flags can be seen in flags.go
func (s *Session) SetFlag(flag uint32) {
	s.Flags ^= 1 << flag
}

// HasFlag returns whether the session has a specified bitflag.
func (s *Session) HasFlag(flag uint32) bool {
	return s.Flags&(1<<flag) > 0
}

// IsStaff returns whether a player is a mod. If CheckAdmin is true it will return if a player is an admin.
func (s *Session) IsStaff(CheckAdmin bool) bool {
	xuid := s.Player.XUID()
	if CheckAdmin {
		for _, v := range utils.Config.Staff.Admins {
			if v == xuid {
				return true
			}
		}
		return false
	}
	for _, v := range utils.Config.Staff.Mods {
		if v == xuid {
			return true
		}
	}
	return false
}

// DefaultFlags will set the default bitflags for the session.
func (s *Session) DefaultFlags() {
	if s.IsStaff(true) || s.IsStaff(false) {
		s.SetFlag(FlagStaff)
		staff[s.Player.Name()] = s
		if s.IsStaff(true) {
			s.SetFlag(FlagAdmin)
		}
	} else {
		s.SetFlag(FlagHasChatCD)
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
	lines := []string{
		"§6Name: §b" + s.Player.Name(),
		"§6Online: §b" + utils.OnlineCount.String(),
		"§6K: §b" + strconv.Itoa(int(s.Stats.Kills.Load())) + " §6D: §b" + strconv.Itoa(int(s.Stats.Deaths.Load())),
		"§bvelvetpractice.tk",
	}
	_, _ = s.scoreboard.WriteString(strings.Join(lines, "\n"))
	s.Player.SendScoreboard(s.scoreboard)
}

func (s *Session) UpdateScoreboard(online, kd bool) {
	if online {
		_ = s.scoreboard.Set(1, "§6Online: §b"+utils.OnlineCount.String())
	}
	if kd {
		_ = s.scoreboard.Set(2, "§6K: §b"+strconv.Itoa(int(s.Stats.Kills.Load()))+" §6D: §b"+strconv.Itoa(int(s.Stats.Deaths.Load())))
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
	s.Stats.Kills.Add(kills)
	s.UpdateScoreboard(false, true)
}

func (s *Session) AddDeaths(deaths uint32) {
	s.Stats.Deaths.Add(deaths)
	s.UpdateScoreboard(false, true)
}

// KDR returns the formatted kill-death ratio of the player.
func (s *Session) KDR() string {
	if s.Stats.Deaths.Load() == 0 || s.Stats.Kills.Load() == 0 {
		return "0.0"
	}
	return strconv.FormatFloat(float64(s.Stats.Kills.Load()/s.Stats.Deaths.Load()), 'f', 2, 32)
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
