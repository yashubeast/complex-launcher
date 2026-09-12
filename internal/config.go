package internal

import (
	"cl/sources"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func LoadPrefixConfig() ([]sources.Prefix, error) {
	configDir, err := os.UserConfigDir()
	if err != nil { return nil, err }

	// TODO: generate a default config, if none exist
	configPath := filepath.Join(
		configDir,
		"complex-launcher",
		"prefixes.json",
	)
	data, err := os.ReadFile(configPath)
	if err != nil { return nil, err }

	var config sources.PrefixConfig

	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf(
			"invalid prefix config: %w",
			err,
		)
	}

	return config.Prefixes, nil
}
