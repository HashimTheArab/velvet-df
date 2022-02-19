package utils

import (
	"encoding/json"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"os"
)

type whitelist struct {
	Enabled bool     `json:"enabled"`
	Players []string `json:"players"`
}

var Whitelist whitelist

func initWhitelist() {
	if _, err := os.Stat("whitelist.json"); os.IsNotExist(err) {
		if data, err := json.Marshal(Whitelist); err != nil {
			panic(err)
		} else if err := ioutil.WriteFile("whitelist.json", data, 0644); err != nil {
			panic(err)
		}
	} else {
		if data, err := ioutil.ReadFile("whitelist.json"); err != nil {
			panic(err)
		} else if err := toml.Unmarshal(data, &Whitelist); err != nil {
			panic(err)
		}
	}
}

func (w whitelist) Contains(target string) bool {
	for _, v := range w.Players {
		if v == target {
			return true
		}
	}
	return false
}

func (w whitelist) Add(target string) {
	w.Players = append(w.Players, target)
	w.Save()
}

func (w whitelist) Remove(target string) {
	var players []string
	for _, v := range w.Players {
		if v != target {
			players = append(players, v)
		}
	}
	w.Players = players
	w.Save()
}

func (w whitelist) Save() {
	if d, err := json.MarshalIndent(Whitelist, "", "    "); err == nil {
		// error
	} else if err := os.WriteFile("whitelist.json", d, 0644); err != nil {
		// error
	}
}
