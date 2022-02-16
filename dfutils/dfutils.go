package dfutils

import (
	"fmt"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	worldmanager "github.com/emperials/df-worldmanager"
	"github.com/justtaldevelops/oomph"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
	"time"
	"velvet/handlers"
	"velvet/session"
	"velvet/utils"
)

var log = logrus.New()

func StartServer() {
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.InfoLevel

	chat.Global.Subscribe(chat.StdoutSubscriber{})

	config, err := utils.ReadDragonflyConfig()
	if err != nil {
		log.Fatalln(err)
	}

	srv := server.New(&config, log)
	srv.SetName("Velvet")
	srv.Allow(allower{})
	srv.CloseOnProgramEnd()
	if err := srv.Start(); err != nil {
		log.Fatalln(err)
	}

	utils.Srv = srv
	utils.WorldMG = worldmanager.New(srv, "worlds/", log)
	utils.Started = time.Now().Unix()

	if files, err := ioutil.ReadDir("worlds"); err == nil {
		for _, f := range files {
			if f.Name() != "world" {
				err = utils.WorldMG.LoadWorld(f.Name(), f.Name(), world.Overworld, nil)
				if err != nil {
					fmt.Println("Error loading world " + f.Name() + ": " + err.Error())
				} else {
					fmt.Println("Loaded world: " + f.Name())
					if f.Name() == utils.Config.World.Build {
						if w, ok := utils.WorldMG.World(f.Name()); ok {
							w.ReadOnly()
						}
					}
				}
			}
		}
	}

	w := srv.World()
	w.SetDefaultGameMode(world.GameModeSurvival)
	w.SetRandomTickSpeed(0)
	w.SetTime(0)
	w.StopTime()

	// AntiCheat start
	go func() {
		ac := oomph.New()
		go func() {
			if err := ac.Start(":"+strings.Split(config.Network.Address, ":")[1], ":19132"); err != nil {
				panic(err)
			}
		}()
		for {
			p, err := ac.Accept()
			if err != nil {
				panic(err)
			}
			p.Handle(handlers.NewACHandler(p))
		}
	}()
	// AntiCheat end

	for {
		p, err := srv.Accept()
		if err != nil {
			return
		}

		s := session.Get(p)
		s.Player = p
		s.OnJoin()
		p.Handle(&handlers.PlayerHandler{Session: s})
	}
}
