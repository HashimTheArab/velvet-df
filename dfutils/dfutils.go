package dfutils

import (
	"fmt"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/oomph-ac/oomph"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"time"
	"velvet/handlers"
	"velvet/session"
	"velvet/utils"
	"velvet/utils/worldmanager"
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
					if w, ok := utils.WorldMG.World(f.Name()); ok {
						switch f.Name() {
						case utils.Config.World.Build:
							w.ReadOnly()
						case utils.Config.World.NoDebuff:
							w.StopWeatherCycle()
						}
					}
				}
			}
		}
	}

	w := srv.World()
	w.SetDefaultGameMode(world.GameModeSurvival)
	w.SetRandomTickSpeed(0)
	w.SetTickRange(0)
	w.SetTime(0)
	w.StopTime()

	// AntiCheat start
	go func() {
		ac := oomph.New(log, ":19132")
		if err := ac.Listen(srv, config.Server.Name); err != nil {
			panic(err)
		}
		for {
			p, err := ac.Accept()
			if err != nil {
				return
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
		if s == nil {
			s = session.New(p.Name())
		}
		s.Player = p
		s.OnJoin()
		p.Handle(handlers.NewPlayerHandler(p, s))
	}
}
