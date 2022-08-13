package main

import (
	"fmt"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/oomph-ac/oomph"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/sirupsen/logrus"
	"log"
	"math/rand"
	"os"
	"time"
	"velvet/db"
	"velvet/handlers"
	"velvet/utils"
	"velvet/utils/worldmanager"
)

var logger = logrus.New()

func startServer() {
	logger.Formatter = &logrus.TextFormatter{ForceColors: true}
	logger.Level = logrus.InfoLevel

	chat.Global.Subscribe(chat.StdoutSubscriber{})

	config, err := readConfig()
	if err != nil {
		log.Fatalln(err)
	}
	config.WorldConfig = func(def world.Config) world.Config {
		def.ReadOnly = true
		def.RandomTickSpeed = 0
		def.Generator = nil
		def.PortalDestination = nil
		return def
	}

	srv := server.New(&config.Config, logger)
	srv.SetName("Velvet")
	srv.Allow(allower{})
	srv.CloseOnProgramEnd()
	srv.World().SetSpawn(cube.Pos{273, 66, 258})
	if err := srv.Start(); err != nil {
		logger.Fatalln(err)
	}
	//
	//if p, err := resource.Compile("resources/pack.zip"); err == nil {
	//	srv.AddResourcePack(p) // todo: p.WithContentKey
	//}

	utils.Srv = srv
	utils.WorldMG = worldmanager.New(srv, "worlds/", logger)
	utils.Started = time.Now().Unix()

	if files, err := os.ReadDir("worlds"); err == nil {
		for _, f := range files {
			if f.Name() != "world" && f.Name() != "world.zip" {
				err = utils.WorldMG.LoadWorld(f.Name(), world.Overworld, nil)
				if err != nil {
					fmt.Println("Error loading world " + f.Name() + ": " + err.Error())
				} else {
					fmt.Println("Loaded world: " + f.Name())
					if w, ok := utils.WorldMG.World(f.Name()); ok {
						switch f.Name() {
						case utils.Config.World.NoDebuff:
							w.SetSpawn(cube.Pos{556, 127, 152})
							w.StopRaining()
							w.StopWeatherCycle()
						case utils.Config.World.Diamond:
							w.SetSpawn(cube.Pos{296, 88, 286})
						case utils.Config.World.Build:
							w.SetSpawn(cube.Pos{218, 115, 255})
						}
					}
				}
			}
		}
	}

	w := srv.World()
	w.SetDefaultGameMode(world.GameModeSurvival)
	w.SetTickRange(0)
	w.SetTime(0)
	w.StopTime()

	go startBroadcasts()

	// AntiCheat start
	if config.Oomph.Enabled {
		go func() {
			ac := oomph.New(logger, config.Oomph.Address)
			//if err := ac.Listen(srv, config.Server.Name, config.Resources.Required); err != nil {
			//	panic(err)
			//}
			go func() {
				if err := ac.Start(config.Network.Address, config.Resources.Folder, config.Resources.Required); err != nil {
					panic(err)
				}
			}()
			for {
				p, err := ac.Accept()
				if err != nil {
					return
				}
				p.Handle(handlers.NewACHandler(p))
			}
		}()
	}
	// AntiCheat end

	for srv.Accept(handleJoin) {

	}

	_ = utils.WorldMG.Close()
}

// handleJoin processes a players join.
func handleJoin(p *player.Player) {
	s, err := db.LoadSession(p)
	if err != nil {
		p.Disconnect(fmt.Sprintf("There was an error loading your account, create a ticket at %s\nError: %s", utils.Config.Discord.Invite, err.Error()))
	}
	p.Handle(handlers.NewPlayerHandler(p, s))
}

// startBroadcasts starts the 5 minute broadcasts.
func startBroadcasts() {
	t := time.NewTicker(time.Minute * 5)

	var messages = [...]string{
		"Join the discord at discord.gg/yYpPX3d!",
		"Make sure to read the rules! /rules.",
		"Don't 2v1 or interfere, it can lead to a kick or temporary ban if excessive.",
		"Buy a rank at velvet.tebex.io!",
		"Simply, the cps limit is 20.",
		//"Can't see in the dark? /nightvision or /nv!",
		"Have any suggestions? Use /suggest to have your suggestion sent to the Velvet discord!",
	}

	go func() {
		time.Sleep(time.Minute * 2)
		t2 := time.NewTicker(time.Minute * 5)
		for range t2.C {
			if len(utils.Srv.Players()) > 0 {
				_, _ = chat.Global.WriteString(text.Colourf("<dark-red><b>FOR SLOW KIDS</b></dark-red><red> - The server is in <b>BETA</b>, it is <b>not released</b>, there are obviously <b>bugs</b></red>"))
			}
		}
	}()

	for range t.C {
		if len(utils.Srv.Players()) > 0 {
			_, _ = chat.Global.WriteString(text.Colourf("<dark-grey>[<purple>Velvet</purple>]<dark-grey> <purple>%s</purple>", messages[rand.Intn(len(messages))]))
		}
	}
}
