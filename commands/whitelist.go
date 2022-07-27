package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"velvet/utils/whitelist"
)

type WhitelistToggle struct {
	Sub cmd.SubCommand `cmd:"toggle"`
}

type WhitelistAdd struct {
	Sub    cmd.SubCommand `cmd:"add"`
	Target string         `cmd:"target"`
}

type WhitelistRemove struct {
	Sub    cmd.SubCommand `cmd:"remove"`
	Target string         `cmd:"target"`
}

func (t WhitelistToggle) Run(_ cmd.Source, output *cmd.Output) {
	if whitelist.Toggle() {
		output.Printf("§aWhitelist has been enabled!")
	} else {
		output.Printf("§cWhitelist has been disabled.")
	}
}

func (t WhitelistAdd) Run(_ cmd.Source, output *cmd.Output) {
	if !whitelist.Enabled() {
		output.Error("§cWhitelist is currently disabled. Use /whitelist toggle to enable it.")
		return
	}
	if whitelist.Contains(t.Target) {
		output.Error("§cThat player is already in the whitelist.")
		return
	}
	whitelist.Add(t.Target)
	output.Printf("§e%v §dhas been added to the whitelist.", t.Target)
}

func (t WhitelistRemove) Run(_ cmd.Source, output *cmd.Output) {
	if !whitelist.Enabled() {
		output.Error("§cWhitelist is currently disabled. Use /whitelist toggle to enable it.")
		return
	}
	if !whitelist.Contains(t.Target) {
		output.Error("§cThat player is not in the whitelist.")
		return
	}
	whitelist.Remove(t.Target)
	output.Printf("§e%v §dhas been removed from the whitelist", t.Target)
}

func (WhitelistToggle) Allow(s cmd.Source) bool { return checkAdmin(s) || checkConsole(s) }
func (WhitelistAdd) Allow(s cmd.Source) bool    { return checkAdmin(s) || checkConsole(s) }
func (WhitelistRemove) Allow(s cmd.Source) bool { return checkAdmin(s) || checkConsole(s) }
