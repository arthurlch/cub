package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/arthurlch/cub/internal/keymap"
	"github.com/arthurlch/cub/internal/theme"
)

const currentVersion = 2

type Config struct {
	Version      int               `json:"version"`
	Theme        string            `json:"theme"`
	SidebarRight bool              `json:"sidebar_right"`
	Keys         map[string]string `json:"keys"`
}

func Path() string {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "cub-config.json"
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "cub", "config.json")
}

func Load() Config {
	data, err := os.ReadFile(Path())
	if err != nil {
		cfg := Config{Theme: theme.Default().Name, Keys: keymap.DefaultChords()}
		_ = Save(cfg)
		return cfg
	}
	cfg := Config{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{Version: currentVersion, Keys: keymap.DefaultChords()}
	}
	if cfg.Version < currentVersion {
		cfg.Keys = keymap.DefaultChords()
		_ = Save(cfg)
	}
	if cfg.Keys == nil {
		cfg.Keys = keymap.DefaultChords()
	}
	return cfg
}

func Save(cfg Config) error {
	cfg.Version = currentVersion
	path := Path()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func (c Config) Chords() map[string]string {
	chords := keymap.DefaultChords()
	for action, chord := range c.Keys {
		chords[action] = chord
	}
	return chords
}
