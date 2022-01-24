package console

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// CommandSender is the command source.
type CommandSender struct {
}

var sender = &CommandSender{}

func (CommandSender) Name() string         { return "CONSOLE" }
func (CommandSender) Position() mgl64.Vec3 { return mgl64.Vec3{} }
func (CommandSender) World() *world.World  { return nil }

func (CommandSender) SendCommandOutput(output *cmd.Output) {
	for _, v := range output.Messages() {
		fmt.Println("[CMD OUTPUT]: " + v)
	}
	for _, v := range output.Errors() {
		fmt.Println("[CMD ERROR]: " + v.Error())
	}
}
