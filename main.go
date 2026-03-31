package main

import (
	"flag"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"example/tuibookie/internal/bookmark"
	"example/tuibookie/internal/config"
	"example/tuibookie/internal/tui"
)

func main() {
	configFlag := flag.String("config", "", "path to bookmarks JSON file")
	flag.Parse()

	configPath := config.ResolvePath(*configFlag, os.Getenv("TUIBOOKIE_CONFIG"))

	if err := config.EnsureConfigDir(configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating config directory: %v\n", err)
		os.Exit(1)
	}

	bm, err := bookmark.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading bookmarks: %v\n", err)
		os.Exit(1)
	}

	model := tui.NewModel(bm, configPath)
	p := tea.NewProgram(model)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
