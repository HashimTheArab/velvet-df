package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/player"
	"strings"
	"time"
)

type EffectType string

type Effect struct {
	Target          []cmd.Target `cmd:"target"`
	EffectName      EffectType   `cmd:"effect"`
	EffectLength    uint32       `cmd:"length"`
	EffectAmplifier int          `cmd:"amplifier"`
}

var effectIdMap = map[string]effect.Type{
	// good effects
	"speed":           effect.Speed{},
	"jump_boost":      effect.JumpBoost{},
	"strength":        effect.Strength{},
	"night_vision":    effect.NightVision{},
	"regeneration":    effect.Regeneration{},
	"resistance":      effect.Resistance{},
	"fire_resistance": effect.FireResistance{},
	"health_boost":    effect.HealthBoost{},
	"absorption":      effect.Absorption{},
	"haste":           effect.Haste{},
	"levitation":      effect.Levitation{},
	"invisibility":    effect.Invisibility{},
	"water_breathing": effect.WaterBreathing{},
	"saturation":      effect.Saturation{},
	"conduit_powder":  effect.ConduitPower{},
	// bad effects
	"slowness":       effect.Slowness{},
	"blindness":      effect.Blindness{},
	"fatal_poison":   effect.FatalPoison{},
	"poison":         effect.Poison{},
	"nausea":         effect.Nausea{},
	"mining_fatigue": effect.MiningFatigue{},
	"hunger":         effect.Hunger{},
	"wither":         effect.Wither{},
	"weakness":       effect.Weakness{},
	"slow_falling":   effect.SlowFalling{},
	// instant effects
	"instant_health": effect.InstantHealth{},
	"instant_damage": effect.InstantDamage{},
	"healing":        effect.InstantHealth{},
}

func (t Effect) Run(_ cmd.Source, output *cmd.Output) {
	if t.EffectAmplifier <= 0 {
		output.Print("§cEffect amplifier must be greater than 0.")
		return
	}

	eff, ok := effectIdMap[strings.ToLower(string(t.EffectName))]
	if ok {
		for _, e := range t.Target {
			pl, ok := e.(*player.Player)
			if ok {
				switch eff := eff.(type) {
				case effect.LastingType:
					pl.AddEffect(effect.New(eff, t.EffectAmplifier, time.Duration(t.EffectLength)*time.Second))
				case effect.Type:
					pl.AddEffect(effect.NewInstant(eff, t.EffectAmplifier))
				}
			}
		}
		return
	}
	output.Print("§cThat effect was not found.")
}

func (EffectType) Type() string {
	return "Effect"
}

func (EffectType) Options(cmd.Source) []string {
	var e []string
	for name := range effectIdMap {
		e = append(e, name)
	}
	return e
}

func (Effect) Allow(s cmd.Source) bool { return checkAdmin(s) }
