package handlers

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/damage"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"strings"
	"velvet/form"
	"velvet/game"
	"velvet/session"
	"velvet/utils"
)

type PlayerHandler struct {
	player.NopHandler
	Session *session.Session
	//PaletteHandler *palette.Handler
	//BrushHandler   *brush.Handler
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
	if !p.Session.HasFlag(session.Building) {
		ctx.Cancel()
	}
	//p.PaletteHandler.HandleBlockBreak(ctx, pos)
}

func (p *PlayerHandler) HandleBlockPlace(ctx *event.Context, _ cube.Pos, _ world.Block) {
	if !p.Session.HasFlag(session.Building) {
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
	} else if /* p.Session.Player.World().Name() == utils.Srv.World().Name() ||*/ source == (damage.SourceFall{}) {
		ctx.Cancel()
	}
	if !ctx.Cancelled() {
		p.Session.UpdateScoreTag(true, false)
	}
}

func (*PlayerHandler) HandleFoodLoss(ctx *event.Context, _ int, _ int) {
	ctx.Cancel()
}

func (p *PlayerHandler) HandleRespawn(*mgl64.Vec3) {
	game.DefaultKit(p.Session.Player)
}

func (p *PlayerHandler) HandleDeath(src damage.Source) {
	if s, ok := src.(damage.SourceEntityAttack); ok {
		_, _ = fmt.Fprintf(chat.Global, "§c%v §ewas killed by §a%v\n", p.Session.Player.Name(), s.Attacker.Name())
	}
}

func (p *PlayerHandler) HandleChat(ctx *event.Context, message *string) {
	ctx.Cancel()
	if strings.Contains(strings.ToLower(*message), "kkkkkkkk") {
		p.Session.Player.Message("por favor, não spam")
		return
	}
	_, _ = fmt.Fprintf(chat.Global, utils.Config.Chat.Basic+"\n", p.Session.Player.Name(), *message)
}

func (p *PlayerHandler) HandleItemUse(_ *event.Context) {
	pl := p.Session.Player
	itm, _ := pl.HeldItems()
	if pl.World().Name() == utils.Srv.World().Name() {
		if t, ok := itm.Value("tool"); ok {
			switch t {
			case game.ArenaItemNbt:
				pl.SendForm(form.FFA(p.Session.Player))
			}
		}
	}
}

//
//func (p *PlayerHandler) HandleItemUseOnBlock(ctx *event.Context, pos cube.Pos, face cube.Face, vec mgl64.Vec3) {
//	p.PaletteHandler.HandleItemUseOnBlock(ctx, pos, face, vec)
//}

func (p *PlayerHandler) HandleQuit() {
	p.Session.Close()
	utils.OnlineCount.Store(utils.OnlineCount.Load() - 1)
	for _, s := range session.All() {
		s.UpdateScoreboard(true, false)
	}
	//p.PaletteHandler.HandleQuit()
	//p.BrushHandler.HandleQuit()
}
