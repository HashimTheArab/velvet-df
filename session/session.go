package session

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/df-mc/dragonfly/server/player/title"
	"github.com/df-mc/dragonfly/server/world"
	"strconv"
	"strings"
	"time"
	"velvet/db"
	"velvet/game"
	"velvet/utils"
)

type Session struct {
	Player   *player.Player
	Flags    uint32
	clicks   []time.Time
	scoreTag struct {
		healthText string
		cpsText    string
		osText     string
	}
	stats db.PlayerStats
}

func (s *Session) OnJoin() {
	s.DefaultFlags()
	s.Player.SendTitle(title.New("§l§dVelvet").WithSubtitle("§bSeason 3 - Reformed"))
	s.UpdateScoreTag(true, true)
	s.UpdateScoreboard(false, true)
	game.DefaultKit(s.Player)
}

func (s *Session) OnQuit() {

}

func (s *Session) SetFlag(flag uint32) {
	s.Flags ^= 1 << flag
}

func (s *Session) HasFlag(flag uint32) bool {
	return s.Flags&(1<<flag) > 0
}

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
	for _, v := range utils.Config.Staff.Staff {
		if v == xuid {
			return true
		}
	}
	return false
}

func (s *Session) DefaultFlags() {
	if s.IsStaff(true) {
		s.SetFlag(Admin)
		s.SetFlag(Staff)
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

func (s *Session) TeleportToSpawn() {
	s.ChangeWorld(utils.Srv.World())
}

func (s *Session) ChangeWorld(w *world.World) {
	w.AddEntity(s.Player)
	s.Player.Teleport(w.Spawn().Vec3())
	g := game.FromWorld(s.Player.World().Name())
	if g != nil {
		g.Kit(s.Player)
	} else if s.Player.World().Name() == utils.Srv.World().Name() {
		game.DefaultKit(s.Player)
	}
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
		s.scoreTag.cpsText += " §bCPS " + strconv.Itoa(s.CPS())
	}
	tag += s.scoreTag.cpsText
	s.Player.SetScoreTag(tag)
}

func (s *Session) UpdateScoreboard(online, kd bool) { // todo: actually use these parameters
	sb := scoreboard.New("§l§dVelvet")
	sb.WriteString("Name " + s.Player.Name() + "\n" + "Online " + utils.OnlineCount.String() + "\n" + "K " + strconv.Itoa(int(s.stats.Kills)) + " D " + strconv.Itoa(int(s.stats.Deaths)) + "\n§dvelvetpractice.tk")
	//_ = sb.Set(0, "Name "+s.Player.Name())
	//_ = sb.Set(1, "Online "+utils.OnlineCount.String())
	//_ = sb.Set(2, "K "+strconv.Itoa(int(s.stats.Kills))+"D "+strconv.Itoa(int(s.stats.Deaths)))
	//_ = sb.Set(3, "§dvelvetpractice.tk")
	s.Player.SendScoreboard(sb)
}

func (s *Session) AddKills(kills uint) {
	s.stats.Kills += kills
	s.UpdateScoreboard(false, true)
}

func (s *Session) AddDeaths(deaths uint) {
	s.stats.Deaths += deaths
	s.UpdateScoreboard(false, true)
}

func (s *Session) KDR() string {
	if s.stats.Deaths == 0 || s.stats.Kills == 0 {
		return "0.0"
	}
	return strconv.FormatFloat(float64(s.stats.Kills/s.stats.Deaths), 'f', 2, 32)
}
