package session

import (
	"fmt"
	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/damage"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/df-mc/dragonfly/server/player/title"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"strconv"
	"strings"
	"sync"
	"time"
	"velvet/game"
	"velvet/perm"
	"velvet/utils"
)

type Session struct {
	// Player is the player of the session.
	Player *player.Player
	// Flags are the bit flags the player has applied.
	Flags uint32

	// clicks are the timestamps the player has clicked.
	clicks   []time.Time
	clicksMu sync.Mutex

	// NetworkSession is the network session of the player.
	NetworkSession *session.Session
	scoreTag       struct {
		healthText string
		cpsText    string
		osText     string
	}
	combat   combat
	combatMu sync.Mutex

	cooldowns   cooldownMap
	cooldownsMu sync.Mutex

	scoreboard *scoreboard.Scoreboard
	rank       atomic.Value[*perm.Rank]

	kills  atomic.Uint32
	deaths atomic.Uint32

	deviceID string

	wandPos1, wandPos2 atomic.Value[mgl64.Vec3]
	bleeding           atomic.Bool

	closed atomic.Bool
}

// New creates a new session.
func New(p *player.Player, rank *perm.Rank, kills, deaths uint32, deviceID string) *Session {
	s := &Session{
		Player:         p,
		NetworkSession: player_session(p),
		cooldowns: map[cooldownType]*Cooldown{
			CooldownTypePearl: {length: time.Second * 15},
			CooldownTypeChat:  {length: time.Second * 3},
		},
		kills:    *atomic.NewUint32(kills),
		deaths:   *atomic.NewUint32(deaths),
		rank:     *atomic.NewValue[*perm.Rank](rank),
		deviceID: deviceID,
	}

	sessions.mutex.Lock()
	sessions.list[p.Name()] = s
	sessions.mutex.Unlock()

	s.DefaultFlags()

	s.UpdateScoreTag(true, true)
	utils.OnlineCount.Add(1)
	All().UpdateScoreboards(true, false)

	s.Player.SendTitle(title.New("§l§dVelvet").WithSubtitle("§bSeason 3 - Reformed"))
	s.Player.EnableInstantRespawn()
	game.DefaultKit(s.Player)

	go func() {
		t := time.NewTicker(time.Second)
		defer t.Stop()
		for range t.C {
			if s.Offline() {
				return
			}
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
		}
	}()
	go func() {
		t := time.NewTicker(time.Millisecond * 50)
		defer t.Stop()
		for range t.C {
			if s.Offline() {
				return
			}
			held, _ := s.Player.HeldItems()
			if _, ok := held.Item().(item.Stick); ok && strings.EqualFold(s.Player.World().Name(), utils.Config.World.Build) {
				s.Player.Move(entity.DirectionVector(s.Player).Mul(2), 0, 0)
			}
		}
	}()

	return s
}

// DeviceID returns the device id of the player.
func (s *Session) DeviceID() string {
	return s.deviceID
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
	return s.HasFlag(FlagStaff)
}

// Mod returns true if a player is a moderator.
func (s *Session) Mod() bool {
	return s.RankName() == perm.Mod
}

// DefaultFlags will set the default bitflags for the session.
func (s *Session) DefaultFlags() {
	if s.Rank() != nil {
		rankName := s.Rank().Name
		if perm.StaffRanks.Contains(rankName) {
			s.SetFlag(FlagStaff)
			AddStaff(s)
			if rankName == perm.Admin || rankName == perm.Owner {
				s.SetFlag(FlagAdmin)
			}
		}
		if rankName == perm.Builder {
			s.SetFlag(FlagBuilder)
		}
	}
}

// Click adds a click to the user's click history.
// Thanks, Tal, for cps!
func (s *Session) Click() {
	s.clicksMu.Lock()
	s.clicks = append(s.clicks, time.Now())
	if len(s.clicks) > 49 {
		s.clicks = s.clicks[1:]
	}
	s.clicksMu.Unlock()
	s.Player.SendTip("§6CPS §b" + strconv.Itoa(s.CPS()))
	s.UpdateScoreTag(false, true)
}

// CPS returns the user's current click per second.
func (s *Session) CPS() int {
	s.clicksMu.Lock()
	defer s.clicksMu.Unlock()
	var clicks int
	for _, past := range s.clicks {
		if time.Since(past) <= time.Second {
			clicks++
		}
	}
	return clicks
}

// UpdateScoreTag updates a users score tag showing health and stuff like that.
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
		fmt.Sprintf("§6K: §b%d §6D: §b%d", s.Kills(), s.Deaths()),
		"§bvelvetpractice.tk",
	}
	_, _ = s.scoreboard.WriteString(strings.Join(lines, "\n"))
	s.Player.SendScoreboard(s.scoreboard)
}

func (s *Session) UpdateScoreboard(online, kd bool) {
	if s.scoreboard == nil {
		s.SaveScoreboard()
		return
	}
	if online {
		s.scoreboard.Set(1, "§6Online: §b"+utils.OnlineCount.String())
	}
	if kd {
		s.scoreboard.Set(2, fmt.Sprintf("§6K: §b%d §6D: §b%d", s.Kills(), s.Deaths()))
	}
	s.Player.SendScoreboard(s.scoreboard)
}

// Vanish handles vanishing or un-vanishing a player.
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

// WandPos returns the selected wand positions for world edit.
func (s *Session) WandPos() (pos1, pos2 mgl64.Vec3) {
	return s.wandPos1.Load(), s.wandPos2.Load()
}

// SetWandPos1 sets the first wand position for world edit.
func (s *Session) SetWandPos1(pos mgl64.Vec3) {
	s.Player.Message(text.Colourf("<green>Pos1 has been set to %v", pos))
	s.wandPos1.Store(pos)
}

// SetWandPos2 sets the second wand position for world edit.
func (s *Session) SetWandPos2(pos mgl64.Vec3) {
	s.Player.Message(text.Colourf("<green>Pos2 has been set to %v", pos))
	s.wandPos2.Store(pos)
}

// StartBleeding starts bleeding for the bleed custom enchant.
func (s *Session) StartBleeding() {
	if s.bleeding.Load() {
		return
	}
	go func() {
		t := time.NewTicker(time.Second * 3)
		defer func() {
			t.Stop()
			s.bleeding.Store(false)
		}()
		runs := 10
		for range t.C {
			runs--
			if runs <= 0 || s.Player.Dead() || s.Offline() || !strings.EqualFold(s.Player.Name(), utils.Config.World.God) {
				return
			}
			s.Player.World().AddParticle(s.Player.Position(), particle.BlockBreak{Block: block.Concrete{Colour: item.ColourRed()}})
			s.Player.Hurt(1, damage.SourceInstantDamageEffect{})
		}
	}()
}

// Kills returns the kills of the player.
func (s *Session) Kills() uint32 {
	return s.kills.Load()
}

// Deaths returns the deaths of the player.
func (s *Session) Deaths() uint32 {
	return s.deaths.Load()
}

// AddKills adds to the player's kills.
func (s *Session) AddKills(kills uint32) {
	s.kills.Add(kills)
	s.UpdateScoreboard(false, true)
}

// AddDeaths adds to the player's deaths.
func (s *Session) AddDeaths(deaths uint32) {
	s.deaths.Add(deaths)
	s.UpdateScoreboard(false, true)
}

// KDR returns the formatted kill-death ratio of the player.
func (s *Session) KDR() string {
	kills, deaths := s.Kills(), s.Deaths()
	if kills == 0 || deaths == 0 {
		return "0.0"
	}
	return strconv.FormatFloat(float64(kills)/float64(deaths), 'f', 2, 32)
}

// Combat returns the players combat handler.
func (s *Session) Combat() *combat {
	s.combatMu.Lock()
	defer s.combatMu.Unlock()
	return &s.combat
}

// Cooldowns returns the players cooldown handler.
func (s *Session) Cooldowns() *cooldownMap {
	s.cooldownsMu.Lock()
	defer s.cooldownsMu.Unlock()
	return &s.cooldowns
}

// Rank will return the rank of the session.
func (s *Session) Rank() *perm.Rank {
	return s.rank.Load()
}

// RankName returns the players rank name or "None" if the player has no rank.
func (s *Session) RankName() string {
	if s.Rank() == nil {
		return "None"
	}
	return s.Rank().Name
}

// SetRank will set the rank for a session and return the new rank.
func (s *Session) SetRank(rank *perm.Rank) {
	s.rank.Store(rank)
}

// Offline checks if the session is fully online.
func (s *Session) Offline() bool {
	return s.closed.Load()
}
