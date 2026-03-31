package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func ResolvePath(flagPath, envPath string) string {
	if flagPath != "" {
		return flagPath
	}
	if envPath != "" {
		return envPath
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "tuibookie", "bookmarks.json")
}

func EnsureConfigDir(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0755)
}

type AppConfig struct {
	BookmarksPath string `json:"bookmarks_path,omitempty"`
}

func ConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "tuibookie")
}

func SaveAppConfig(configDir string, cfg AppConfig) error {
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(configDir, "config.json"), data, 0644)
}

func LoadAppConfig(configDir string) (AppConfig, error) {
	path := filepath.Join(configDir, "config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return AppConfig{}, nil
		}
		return AppConfig{}, err
	}
	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return AppConfig{}, err
	}
	return cfg, nil
}
