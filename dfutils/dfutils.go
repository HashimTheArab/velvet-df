package dfutils

import (
	"fmt"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	worldmanager "github.com/emperials/df-worldmanager"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"velvet/handlers"
	"velvet/session"
	"velvet/utils"
)

func StartServer() {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.DebugLevel

	chat.Global.Subscribe(chat.StdoutSubscriber{})

	config, err := utils.ReadDragonflyConfig()
	if err != nil {
		log.Fatalln(err)
	}

	srv := server.New(&config, log)
	srv.CloseOnProgramEnd()
	if err := srv.Start(); err != nil {
		log.Fatalln(err)
	}

	utils.Srv = srv
	utils.WorldMG = worldmanager.New(srv, "worlds/", log)

	// todo: make this auto load all gamemode worlds
	if files, err := ioutil.ReadDir("worlds"); err == nil {
		for _, f := range files {
			if f.Name() != "world" {
				err = utils.WorldMG.LoadWorld(f.Name(), f.Name(), world.Overworld, nil)
				if err != nil {
					fmt.Println("Error loading world " + f.Name() + ": " + err.Error())
				} else {
					fmt.Println("Loaded world: " + f.Name())
				}
			}
		}
	}

	w := srv.World()
	w.SetDefaultGameMode(world.GameModeSurvival)
	w.SetRandomTickSpeed(0)
	w.SetTime(0)
	w.StopTime()

	for {
		p, err := srv.Accept()
		if err != nil {
			return
		}

		s := session.New(p)
		p.Handle(&handlers.PlayerHandler{
			Session: s,
			//PaletteHandler: palette.NewHandler(p),
			//BrushHandler:   brush.NewHandler(p),
		})
		p.EnableInstantRespawn()
		utils.OnlineCount.Add(1)
		for _, s := range session.All() {
			s.UpdateScoreboard(true, false)
		}
	}
}
