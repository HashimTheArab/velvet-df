package main

import (
	"fmt"
	"github.com/df-mc/dragonfly/server"
	"github.com/pelletier/go-toml"
	"os"
)

// Config is an extension of the Dragonfly server config to include fields specific to Velvet.
type Config struct {
	server.Config
	// Pack contains fields related to the pack.
	Pack struct {
		// Key is the pack encryption key.
		Key string
		// Path is the path to the pack.
		Path string
	}
	// Oomph contains fields specific to Oomph.
	Oomph struct {
		Enabled bool
		// Enabled specifies if Oomph should be enabled.
		// Address is the address to run Oomph on.
		Address string
	}
}

// DefaultConfig returns a default config for the server.
func DefaultConfig() Config {
	c := Config{Config: server.DefaultConfig()}
	return c
}

// readConfig reads the configuration from the config/dragonfly.toml file, or creates the file if it does not yet exist.
func readConfig() (Config, error) {
	c := DefaultConfig()
	if _, err := os.Stat("config/dragonfly.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return c, fmt.Errorf("failed encoding default config: %v", err)
		}
		if err := os.WriteFile("config/dragonfly.toml", data, 0644); err != nil {
			return c, fmt.Errorf("failed creating config: %v", err)
		}
		return c, nil
	}
	data, err := os.ReadFile("config/dragonfly.toml")
	if err != nil {
		return c, fmt.Errorf("error reading config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return c, fmt.Errorf("error decoding config: %v", err)
	}
	return c, nil
}
