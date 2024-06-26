package handlers

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/damage"
	heal "github.com/df-mc/dragonfly/server/entity/healing"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/we/palette"
	"github.com/go-gl/mathgl/mgl64"
	"strings"
	"time"
	"velvet/db"
	"velvet/form"
	"velvet/game"
	"velvet/healing"
	vitem "velvet/item"
	"velvet/session"
	"velvet/utils"
	"velvet/utils/enchants"
)

type PlayerHandler struct {
	player.NopHandler
	Session *session.Session
	ph      *palette.Handler
}

func NewPlayerHandler(p *player.Player, s *session.Session) *PlayerHandler {
	return &PlayerHandler{
		Session: s,
		ph:      palette.NewHandler(p),
	}
}

func (p *PlayerHandler) HandlePunchAir(_ *event.Context) {
	p.Session.Click()
}

func (p *PlayerHandler) HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack) {
	p.ph.HandleBlockBreak(ctx, pos, drops)
	held, _ := p.Session.Player.HeldItems()
	if _, ok := held.Value("wand"); ok {
		ctx.Cancel()
		p.Session.SetWandPos1(pos.Vec3())
		return
	}

	if p.Session.Player.World().Name() == utils.Config.World.Build && utils.BuildBlocks.Remove(pos) {
		return
	}

	if !p.Session.HasFlag(session.FlagBuilding) {
		ctx.Cancel()
		return
	}
}

func (p *PlayerHandler) HandleBlockPlace(ctx *event.Context, pos cube.Pos, _ world.Block) {
	if p.Session.Player.World().Name() == utils.Config.World.Build {
		utils.BuildBlocks.Set(pos)
		return
	}
	if !p.Session.HasFlag(session.FlagBuilding) {
		ctx.Cancel()
		return
	}
}

func (*PlayerHandler) HandleItemDrop(ctx *event.Context, _ *entity.Item) {
	ctx.Cancel()
}

func (p *PlayerHandler) HandleAttackEntity(_ *event.Context, e world.Entity, h *float64, v *float64, _ *bool) {
	p.Session.Click()
	g := game.FromWorld(p.Session.Player.World().Name())
	if g != nil {
		*h, *v = g.Knockback.Horizontal, g.Knockback.Vertical
	} else {
		*h, *v = 0.398, 0.405
	}
	if strings.EqualFold(g.Name, utils.Config.World.God) {
		held, _ := p.Session.Player.HeldItems()
		_, ok := held.Value("gsword")
		if !ok {
			return
		}
		t, ok := e.(*player.Player)
		if !ok || t.AttackImmune() {
			return
		}
		enchants.Kaboom(p.Session.Player, t, h, v)
		enchants.Zeus(p.Session.Player, t)
		enchants.Bleed(p.Session.Player, t)
		enchants.Hades(p.Session.Player, t)
		enchants.Poison(p.Session.Player, t)
		enchants.Lifesteal(p.Session.Player, t)
		enchants.Scorch(p.Session.Player, t)
		//switch rand.Intn(35) {
		//case 1:
		//	// todo: kaboom
		//case 2:
		//	// $enchants->lightning($player);
		//	//							$pp = mt_rand(0, 150);
		//	//							if ($pp == 3) {
		//	//								$enchants->kaboom($player, $damager);
		//	//							}
		//case 3:
		//	// if(!$player->hasFlag(Flags::BLEEDING)){
		//	//								$player->setFlag(Flags::BLEEDING);
		//	//								$this->main->getScheduler()->scheduleRepeatingTask(new BleedTask($this->main, $player), 60);
		//	//							}
		//case 4:
		//	// todo: poison
		//}
	}
}

func (p *PlayerHandler) HandleHurt(ctx *event.Context, _ *float64, attackImmunity *time.Duration, src damage.Source) {
	if _, ok := src.(damage.SourceVoid); ok {
		ctx.Cancel()
		utils.Srv.World().AddEntity(p.Session.Player)
		p.Session.Player.Teleport(utils.Srv.World().Spawn().Vec3())
		return
	}

	if p.Session.Player.World() == utils.Srv.World() || src == (damage.SourceFall{}) {
		ctx.Cancel()
		return
	}

	*attackImmunity = time.Millisecond * 475
	p.Session.UpdateScoreTag(true, false)
	if source, ok := src.(damage.SourceEntityAttack); ok {
		if t, ok := source.Attacker.(*player.Player); ok {
			s := session.Get(t)
			if !s.Combat().Tagged() {
				s.Player.Message("§cYou are now in combat.")
			}
			if !p.Session.Combat().Tagged() {
				p.Session.Player.Message("§cYou are now in combat.")
			}
			s.Combat().Tag(true)
			p.Session.Combat().Tag(true)
		}
	}
}

func (p *PlayerHandler) HandleHeal(_ *event.Context, _ *float64, _ heal.Source) {
	p.Session.UpdateScoreTag(true, false)
}

func (*PlayerHandler) HandleFoodLoss(ctx *event.Context, _ int, _ int) {
	ctx.Cancel()
}

func (p *PlayerHandler) HandleChangeWorld(_, after *world.World) {
	p.Session.Player.Teleport(after.Spawn().Vec3())

	for _, e := range p.Session.Player.Effects() {
		p.Session.Player.RemoveEffect(e.Type())
	}

	g := game.FromWorld(after.Name())
	if g != nil {
		g.Kit(p.Session.Player)
	} else if after == utils.Srv.World() {
		p.Session.Player.Armour().Clear()
		p.Session.Player.Inventory().Clear()
		game.DefaultKit(p.Session.Player)
	}
}

func (p *PlayerHandler) HandleRespawn(pos *mgl64.Vec3, w **world.World) {
	*w, *pos = utils.Srv.World(), utils.Srv.World().Spawn().Vec3()
	game.DefaultKit(p.Session.Player)
}

func (p *PlayerHandler) HandleDeath(source damage.Source) {
	g := game.FromWorld(p.Session.Player.World().Name())
	if source, ok := source.(damage.SourceEntityAttack); ok {
		if g == nil {
			_, _ = fmt.Fprintf(chat.Global, "§c%v was killed by %v", p.Session.Player.Name(), source.Attacker.Name())
		} else {
			g.BroadcastDeathMessage(p.Session.Player, source.Attacker.(*player.Player))
			if pl, ok := source.Attacker.(*player.Player); ok {
				g.Kit(pl)
				pl.Heal(pl.MaxHealth(), healing.SourceKill{})
			}
		}
		if pl, ok := source.Attacker.(*player.Player); ok {
			session.Get(pl).AddKills(1)
		}
		p.Session.AddDeaths(1)
	} else {
		_, _ = fmt.Fprintf(chat.Global, "§c%v died.", p.Session.Player.Name())
	}
	p.Session.Player.Armour().Clear()
	p.Session.Player.Inventory().Clear()
	p.Session.UpdateScoreTag(true, true)
	if p.Session.Combat().Tagged() {
		p.Session.Combat().Tag(false)
	}
}

func (p *PlayerHandler) HandleCommandExecution(ctx *event.Context, command cmd.Command, _ []string) {
	if p.Session.Combat().Tagged() {
		for _, v := range session.CombatBannedCommands {
			if command.Name() == v {
				ctx.Cancel()
				p.Session.Player.Message("§cYou cannot use this command while in combat.")
				break
			}
		}
	}
}

func (p *PlayerHandler) HandleChat(ctx *event.Context, message *string) {
	rank := p.Session.Rank()
	defer ctx.Cancel()
	if rank == nil {
		p.Session.Cooldowns().Handle(ctx, p.Session.Player, session.CooldownTypeChat)
	}
	if strings.Contains(strings.ToLower(*message), "kkkkkkkk") {
		p.Session.Player.Message("stop")
		return
	}
	if !ctx.Cancelled() {
		if rank == nil {
			_, _ = fmt.Fprintf(chat.Global, utils.Config.Chat.Basic+"\n", p.Session.Player.Name(), *message)
		} else {
			_, _ = fmt.Fprintf(chat.Global, rank.ChatFormat+"\n", p.Session.Player.Name(), *message)
		}
	}
}

func (p *PlayerHandler) HandleItemUse(ctx *event.Context) {
	pl := p.Session.Player
	itm, _ := pl.HeldItems()
	if itm.Item() == (item.EnderPearl{}) {
		p.Session.Cooldowns().Handle(ctx, pl, session.CooldownTypePearl)
	}
	vitem.Override(p.Session, ctx)
	if pl.World().Name() == utils.Srv.World().Name() {
		if t, ok := itm.Value("tool"); ok {
			switch t {
			case game.ArenaItemNbt:
				pl.SendForm(form.FFA(pl))
			}
		}
	}
}

func (p *PlayerHandler) HandleItemUseOnBlock(ctx *event.Context, pos cube.Pos, face cube.Face, vec mgl64.Vec3) {
	p.ph.HandleItemUseOnBlock(ctx, pos, face, vec)
	held, _ := p.Session.Player.HeldItems()
	_, pos2 := p.Session.WandPos()
	if _, ok := held.Value("wand"); ok && pos != cube.PosFromVec3(pos2) {
		ctx.Cancel()
		p.Session.SetWandPos2(pos.Vec3())
		return
	}
	if p.Session.Player.World() == utils.Srv.World() {
		ctx.Cancel()
	}
}

func (p *PlayerHandler) HandleQuit() {
	if p.Session.Combat().Tagged() {
		p.Session.Player.Inventory().Clear()
		p.Session.Player.Armour().Clear()
		p.Session.Player.Hurt(p.Session.Player.MaxHealth(), damage.SourceVoid{})
		_, _ = fmt.Fprintf(chat.Global, "§c%v died.", p.Session.Player.Name())
	}

	_ = db.SaveSession(p.Session)
	p.Session.Close()

	utils.OnlineCount.Store(utils.OnlineCount.Load() - 1)
	session.All().UpdateScoreboards(true, false)
	p.ph.HandleQuit()
}
