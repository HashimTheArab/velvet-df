package whitelist

import (
	"encoding/json"
	"golang.org/x/exp/slices"
	"io/ioutil"
	"os"
	"strings"
	"velvet/session"
)

var whitelist = struct {
	Enabled bool     `json:"enabled"`
	Players []string `json:"players"`
}{}

func init() {
	if _, err := os.Stat("whitelist.json"); os.IsNotExist(err) {
		if data, err := json.Marshal(whitelist); err != nil {
			panic(err)
		} else if err := ioutil.WriteFile("whitelist.json", data, 0644); err != nil {
			panic(err)
		}
	} else {
		if data, err := ioutil.ReadFile("whitelist.json"); err != nil {
			panic(err)
		} else if err := json.Unmarshal(data, &whitelist); err != nil {
			panic(err)
		}
	}
}

func Enabled() bool {
	return whitelist.Enabled
}

func Toggle() bool {
	whitelist.Enabled = !whitelist.Enabled
	Save()
	return whitelist.Enabled
}

func Contains(target string) bool {
	return slices.Contains(whitelist.Players, strings.ToLower(target))
}

func Add(target string) {
	whitelist.Players = append(whitelist.Players, strings.ToLower(target))
	Save()
}

func Remove(target string) {
	var players []string
	for _, v := range whitelist.Players {
		if v != target {
			players = append(players, v)
		}
	}
	whitelist.Players = players
	Save()
}

func Save() {
	if d, err := json.MarshalIndent(whitelist, "", "    "); err != nil {
		session.AllStaff().Message("§cWhitelist failed to marshal json: " + err.Error())
	} else if err := os.WriteFile("whitelist.json", d, 0644); err != nil {
		session.AllStaff().Message("§cWhitelist failed to save: " + err.Error())
	}
}
