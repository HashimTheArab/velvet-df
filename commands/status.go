package commands

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"os"
	"runtime"
	"strconv"
	"time"
	"velvet/utils"
)

type Status struct{}

type statusData struct {
	memory struct {
		Total string `json:"MemTotal"`
		Free  string `json:"MemAvailable"`
	}
	cpu struct {
		Model string `json:"model name"`
	}
}

func (Status) Run(_ cmd.Source, output *cmd.Output) {
	go func() {
		status := "§e--§dServer Status§e--\n"
		add := func(name, value string) {
			status += "§d" + name + ": §e" + value + "\n"
		}

		s := getStatusData()

		add("Uptime", (time.Second * time.Duration(time.Now().Unix()-utils.Started)).String())
		add("CPU Usage", "")
		add("Memory", fmt.Sprintf("§d(§c%v§d/§a%v", s.memory.Free, s.memory.Total))
		for _, w := range utils.WorldMG.Worlds() {
			players := 0
			for _, v := range w.Entities() {
				if _, ok := v.(*player.Player); ok {
					players++
				}
			}
			add("World ("+w.Name()+")", strconv.Itoa(players)+" §dPlayers, "+strconv.Itoa(len(w.Entities()))+" §dTotal Entities")
		}
		output.Printf(status)
	}()
}

func init() {
	spew.Dump(getStatusData())
}

func getStatusData() statusData {
	var s statusData
	switch runtime.GOOS {
	case "linux", "android":
		if d, err := os.ReadFile("/proc/meminfo"); err != nil {
			_ = json.Unmarshal(d, &s.memory)
		}
		if d, err := os.ReadFile("/proc/cpuinfo"); err != nil {
			_ = json.Unmarshal(d, &s.cpu)
		}
	}
	return s
}

func (Status) Allow(s cmd.Source) bool { return checkStaff(s) }
