package utils

import (
	"fmt"
	"github.com/df-mc/dragonfly/server"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"os"
)

type config struct {
	Staff struct {
		Owner struct {
			Name string
			XUID string
		}
	}
	Discord struct {
		Webhook struct {
			TitleIDLogger   string
			BanLogger       string
			UnbanLogger     string
			AntiCheatLogger string
		}
	}
	World struct {
		NoDebuff string
		Build    string
		Diamond  string
	}
	Chat struct {
		Basic string
	}
	Kick struct {
		Screen, Broadcast string
	}
	Ban struct {
		Screen, LoginScreen, Broadcast                                      string
		BlacklistScreen, BlacklistBroadcast                                 string
		CanOnlyBanOne, PlayerAlreadyBanned, PlayerNotBanned, PlayerUnbanned string
		Info                                                                string
	}
	Rank struct {
		Removed, Set string
	}
	Message struct {
		WelcomeToSpawn, DefaultSpawnSet                                                      string
		GameModeSetByPlayer, GameModeSetBySelf, GameModeSetOther                             string
		TeleportSelfToPos, TeleportSelfToPlayer, TeleportTargetToPos, TeleportTargetToTarget string
		BuildTooManyPlayers, SelfNotInBuilderMode, SelfInBuilderMode, SetPlayerInBuilderMode,
		SetBuilderModeByPlayer, UnsetBuilderModeByPlayer, UnsetPlayerInBuilderMode string
		ModeUnavailable, CannotPunishPlayer, InvalidPunishmentTime, ServerNotAvailable string
		SpecifyReason, Alias, NeverJoined                                              string
	}
	DeathMessage struct {
		List              []string
		Default, NoDebuff string
	}
}

var Config *config

func init() {
	ReadVelvetConfig()
}

// ReadDragonflyConfig reads the configuration from the config/dragonfly.toml file, or creates the file if it does not yet exist.
func ReadDragonflyConfig() (server.Config, error) {
	c := server.DefaultConfig()
	if _, err := os.Stat("config/dragonfly.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return c, fmt.Errorf("failed encoding default config: %v", err)
		}
		if err := ioutil.WriteFile("config/dragonfly.toml", data, 0644); err != nil {
			return c, fmt.Errorf("failed creating config: %v", err)
		}
		return c, nil
	}
	data, err := ioutil.ReadFile("config/dragonfly.toml")
	if err != nil {
		return c, fmt.Errorf("error reading config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return c, fmt.Errorf("error decoding config: %v", err)
	}
	return c, nil
}

// ReadVelvetConfig reads the configuration from the config/velvet.toml file and sets the proper values.
func ReadVelvetConfig() {
	Config = &config{}
	data, err := ioutil.ReadFile("config/velvet.toml")
	if err != nil {
		fmt.Printf("error reading velvet config: %v", err)
		return
	}
	if err := toml.Unmarshal(data, Config); err != nil {
		fmt.Printf("error decoding velvet config: %v", err)
		return
	}
}
