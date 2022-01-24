package session

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/df-mc/dragonfly/server/player/title"
	"github.com/df-mc/dragonfly/server/session"
	"strconv"
	"strings"
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
	CombatTime uint8

	scoreboard *scoreboard.Scoreboard
}

// OnJoin is called when the player joins.
func (s *Session) OnJoin() {
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
}

// OnQuit is called when the session leaves.
func (s *Session) OnQuit() {
	s.Save()
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
	if s.IsStaff(true) {
		s.SetFlag(Admin)
		s.SetFlag(Staff)
		staff[s.Player.Name()] = s
	} else if s.IsStaff(false) {
		s.SetFlag(Staff)
	}
}

// Click adds a click to the user's click history.
// Thanks, Tal, for cps!
func (s *Session) Click() {
	s.clicks = append(s.clicks, time.Now())
	if len(s.clicks) > 49 {
		s.clicks = s.clicks[1:]
	}
	s.Player.SendTip("§dCPS §b" + strconv.Itoa(s.CPS()))
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
	s.scoreboard = scoreboard.New("§l§dVelvet")
	lines := []string{
		"§dName:§a " + s.Player.Name(),
		"§dOnline: §a" + utils.OnlineCount.String(),
		"§dK: §a" + strconv.Itoa(int(s.Stats.Kills.Load())) + " §dD: §a" + strconv.Itoa(int(s.Stats.Deaths.Load())),
		"§avelvetpractice.tk",
	}
	_, _ = s.scoreboard.WriteString(strings.Join(lines, "\n"))
	s.Player.SendScoreboard(s.scoreboard)
}

func (s *Session) UpdateScoreboard(online, kd bool) {
	if online {
		_ = s.scoreboard.Set(1, "§dOnline: §a"+utils.OnlineCount.String())
	}
	if kd {
		_ = s.scoreboard.Set(2, "§dK: §a"+strconv.Itoa(int(s.Stats.Kills.Load()))+" §dD: §a"+strconv.Itoa(int(s.Stats.Deaths.Load())))
	}
	s.Player.SendScoreboard(s.scoreboard)
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
