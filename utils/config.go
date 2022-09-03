package utils

import (
	"fmt"
	"github.com/pelletier/go-toml"
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
		Invite string
	}
	World struct {
		NoDebuff string
		Build    string
		Diamond  string
		God      string
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

// ReadVelvetConfig reads the configuration from the config/velvet.toml file and sets the proper values.
func ReadVelvetConfig() {
	Config = &config{}
	data, err := os.ReadFile("config/velvet.toml")
	if err != nil {
		fmt.Printf("error reading velvet config: %v", err)
		return
	}
	if err := toml.Unmarshal(data, Config); err != nil {
		fmt.Printf("error decoding velvet config: %v", err)
		return
	}
}
