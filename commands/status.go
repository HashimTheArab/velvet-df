package commands

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
	"velvet/utils"
)

type Status struct{}

type statusData struct {
	memory struct {
		Free  string
		Used  string
		Total string
	}
	cpu struct {
		Model string
	}
}

func (Status) Run(_ cmd.Source, output *cmd.Output) {
	status := "§e--§dServer Status§e--\n"
	add := func(name, value string) {
		status += "§d" + name + ": §e" + value + "\n"
	}
	s := getStatusData()
	s.format()

	add("Uptime", (time.Second * time.Duration(time.Now().Unix()-utils.Started)).String())
	add("CPU Model", s.cpu.Model)
	add("Memory", fmt.Sprintf("§d(§c%v§d/§a%v", s.memory.Used, s.memory.Total))
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
}

func init() {
	spew.Dump(getStatusData())
}

func getStatusData() statusData {
	var s statusData
	switch runtime.GOOS {
	case "linux", "android":
		if d, err := os.ReadFile("/proc/meminfo"); err == nil {
			parseFields(d, map[string]*string{"MemAvailable": &s.memory.Free, "MemTotal": &s.memory.Total}, true)
		}
		if d, err := os.ReadFile("/proc/cpuinfo"); err == nil {
			parseFields(d, map[string]*string{"model name": &s.cpu.Model}, false)
		}
	}
	return s
}

func parseFields(data []byte, fields map[string]*string, memory bool) {
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		index := strings.IndexRune(line, ':')
		if index == -1 {
			continue
		}
		if val, ok := fields[line[:index]]; ok {
			if memory {
				*val = strings.TrimSpace(strings.TrimRight(line[index+1:], "kB"))
			} else {
				// todo
			}
		}
	}
}

func (d *statusData) format() {
	var memory []int
	formatStorage := func(s *string) {
		if n, err := strconv.Atoi(*s); err == nil {
			*s = strconv.Itoa(n/1000) + "MB"
			memory = append(memory, n/1000)
		} else {
			*s = "Unavailable"
		}
	}
	formatStorage(&d.memory.Free)
	formatStorage(&d.memory.Total)
	if len(memory) >= 2 {
		d.memory.Used = strconv.Itoa(memory[1]-memory[0]) + " MB"
	} else {
		d.memory.Used = "Unavailable"
	}
}

func (Status) Allow(s cmd.Source) bool { return checkStaff(s) }
