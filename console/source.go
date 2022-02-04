package console

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/sirupsen/logrus"
)

// CommandSender is the command source.
type CommandSender struct {
	log *logrus.Logger
}

func (CommandSender) Name() string         { return "CONSOLE" }
func (CommandSender) Position() mgl64.Vec3 { return mgl64.Vec3{} }
func (CommandSender) World() *world.World  { return nil }

func (c CommandSender) SendCommandOutput(output *cmd.Output) {
	for _, s := range output.Messages() {
		c.log.Info(text.ANSI(s))
	}
	for _, s := range output.Errors() {
		c.log.Error(text.ANSI(s))
	}
}
