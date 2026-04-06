package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"

	"github.com/orvad/tuibookie/internal/bookmark"
	"github.com/orvad/tuibookie/internal/config"
	"github.com/orvad/tuibookie/internal/tui"
)

var version = "dev"

func main() {
	configFlag := flag.String("config", "", "path to bookmarks JSON file")
	flag.Parse()

	configDir := config.ConfigDir()

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0700); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating config directory: %v\n", err)
		os.Exit(1)
	}

	// Ensure config.json exists
	configJsonPath := filepath.Join(configDir, "config.json")
	if _, err := os.Stat(configJsonPath); os.IsNotExist(err) {
		if err := config.SaveAppConfig(configDir, config.AppConfig{}); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating config file: %v\n", err)
			os.Exit(1)
		}
	}

	flagVal := *configFlag
	envVal := os.Getenv("TUIBOOKIE_CONFIG")
	configPath := config.ResolvePath(flagVal, envVal, configDir)

	var pathSource tui.PathSource
	switch {
	case flagVal != "":
		pathSource = tui.PathSourceFlag
	case envVal != "":
		pathSource = tui.PathSourceEnv
	default:
		appCfg, _ := config.LoadAppConfig(configDir)
		if appCfg.BookmarksPath != "" {
			pathSource = tui.PathSourceConfig
		} else {
			pathSource = tui.PathSourceDefault
		}
	}

	if err := config.EnsureConfigDir(configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating config directory: %v\n", err)
		os.Exit(1)
	}

	bm, err := bookmark.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading bookmarks: %v\n", err)
		os.Exit(1)
	}

	model := tui.NewModel(bm, configPath, configDir, pathSource, version)
	p := tea.NewProgram(model)

	m, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Write last executed command to file for shell history integration
	lastCmdPath := filepath.Join(configDir, "lastcmd")
	if final, ok := m.(tui.Model); ok && final.ExecutedCmd() != "" {
		os.WriteFile(lastCmdPath, []byte(final.ExecutedCmd()), 0600)
	} else {
		os.Remove(lastCmdPath)
	}
}
