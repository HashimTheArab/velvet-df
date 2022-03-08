package handlers

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/damage"
	"github.com/df-mc/dragonfly/server/entity/healing"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/we/brush"
	"github.com/df-mc/we/palette"
	"github.com/go-gl/mathgl/mgl64"
	"strings"
	"velvet/form"
	"velvet/game"
	vitem "velvet/item"
	"velvet/session"
	"velvet/utils"
)

type PlayerHandler struct {
	player.NopHandler
	Session *session.Session
	ph      *palette.Handler
	bh      *brush.Handler
}

func NewPlayerHandler(p *player.Player, s *session.Session) *PlayerHandler {
	return &PlayerHandler{
		Session: s,
		ph:      palette.NewHandler(p),
		bh:      brush.NewHandler(p),
	}
}

func (p *PlayerHandler) HandleAttackEntity(_ *event.Context, _ world.Entity, h *float64, v *float64, critical *bool) {
	p.Session.Click()
	g := game.FromWorld(p.Session.Player.World().Name())
	if g != nil {
		*h, *v = g.Knockback.Horizontal, g.Knockback.Vertical
	} else {
		*h, *v = 0.398, 0.405
	}
	if p.Session.Player.Sprinting() {
		*critical = false
	}
}

func (p *PlayerHandler) HandlePunchAir(_ *event.Context) {
	p.Session.Click()
}

func (p *PlayerHandler) HandleBlockBreak(ctx *event.Context, pos cube.Pos, _ *[]item.Stack) {
	p.ph.HandleBlockBreak(ctx, pos)
	if p.Session.Player.World().Name() == utils.Config.World.Build {
		utils.BuildBlocks.Mutex.Lock()
		defer utils.BuildBlocks.Mutex.Unlock()
		if _, ok := utils.BuildBlocks.Blocks[pos]; ok {
			delete(utils.BuildBlocks.Blocks, pos)
			return
		}
	}
	if !p.Session.HasFlag(session.FlagBuilding) {
		ctx.Cancel()
	}
}

func (p *PlayerHandler) HandleBlockPlace(ctx *event.Context, pos cube.Pos, _ world.Block) {
	if p.Session.Player.World().Name() == utils.Config.World.Build {
		utils.BuildBlocks.Set(pos)
	} else if !p.Session.HasFlag(session.FlagBuilding) {
		ctx.Cancel()
	}
}

func (*PlayerHandler) HandleItemDrop(ctx *event.Context, _ *entity.Item) {
	ctx.Cancel()
}

func (p *PlayerHandler) HandleHurt(ctx *event.Context, _ *float64, source damage.Source) {
	if source == (damage.SourceVoid{}) {
		ctx.Cancel()
		p.Session.TeleportToSpawn()
	} else if p.Session.Player.World().Name() == utils.Srv.World().Name() || source == (damage.SourceFall{}) {
		ctx.Cancel()
	}
	if !ctx.Cancelled() {
		p.Session.UpdateScoreTag(true, false)
		if source, ok := source.(damage.SourceEntityAttack); ok {
			if pl, ok := source.Attacker.(*player.Player); ok {
				s := session.Get(pl)
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
}

func (p *PlayerHandler) HandleHeal(ctx *event.Context, _ *float64, _ healing.Source) {
	if !ctx.Cancelled() {
		p.Session.UpdateScoreTag(true, false)
	}
}

func (*PlayerHandler) HandleFoodLoss(ctx *event.Context, _ int, _ int) {
	ctx.Cancel()
}

func (p *PlayerHandler) HandleChangeWorld(_, new *world.World) {
	p.Session.Player.Teleport(new.Spawn().Vec3())

	for _, e := range p.Session.Player.Effects() {
		p.Session.Player.RemoveEffect(e.Type())
	}

	g := game.FromWorld(new.Name())
	if g != nil {
		g.Kit(p.Session.Player)
	} else if new.Name() == utils.Srv.World().Name() {
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
				pl.Heal(pl.MaxHealth(), healing.SourceCustom{})
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
			}
		}
	}
}

func (p *PlayerHandler) HandleChat(ctx *event.Context, message *string) {
	if p.Session.Rank() == nil {
		p.Session.Cooldowns().Handle(ctx, p.Session.Player, session.CooldownTypeChat)
	}
	if strings.Contains(strings.ToLower(*message), "kkkkkkkk") {
		p.Session.Player.Message("no spam pls")
		ctx.Cancel()
		return
	}
	if !ctx.Cancelled() {
		ctx.Cancel()
		rank := p.Session.Rank()
		if rank != nil {
			_, _ = fmt.Fprintf(chat.Global, rank.ChatFormat+"\n", p.Session.Player.Name(), *message)
		} else {
			_, _ = fmt.Fprintf(chat.Global, utils.Config.Chat.Basic+"\n", p.Session.Player.Name(), *message)
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
	p.bh.HandleItemUse(ctx)
}

func (p *PlayerHandler) HandleItemUseOnBlock(ctx *event.Context, pos cube.Pos, face cube.Face, vec mgl64.Vec3) {
	p.ph.HandleItemUseOnBlock(ctx, pos, face, vec)
}

func (p *PlayerHandler) HandleQuit() {
	if p.Session.Combat().Tagged() {
		p.Session.Player.Inventory().Clear()
		p.Session.Player.Armour().Clear()
		p.Session.Player.Hurt(p.Session.Player.MaxHealth(), damage.SourceCustom{})
		_, _ = fmt.Fprintf(chat.Global, "§c%v died.", p.Session.Player.Name())
	}
	p.Session.Close()
	utils.OnlineCount.Store(utils.OnlineCount.Load() - 1)
	session.All().UpdateScoreboards(true, false)
	p.bh.HandleQuit()
	p.ph.HandleQuit()
}
