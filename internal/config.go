package internal

import (
	"cl/sources"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	TerminalCommand string           `json:"terminal_command"`
	WindowClassName string           `json:"window_class_name"`
	Prefixes        []sources.Prefix `json:"prefixes"`
}

func LoadConfig() (Config, error) {
	configDir, err := os.UserConfigDir()
	if err != nil { return Config{}, err }

	// TODO: generate a default config, if none exist
	configPath := filepath.Join(
		configDir,
		"complex-launcher",
		"config.json",
	)
	data, err := os.ReadFile(configPath)
	if err != nil { return Config{}, err }

	var config Config

	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, fmt.Errorf(
			"invalid config: %w",
			err,
		)
	}

	config.TerminalCommand = fmt.Sprintf(config.TerminalCommand, config.WindowClassName)
	return config, nil
}
