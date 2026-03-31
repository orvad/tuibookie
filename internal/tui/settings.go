package tui

import (
	"os"
	"path/filepath"
	"strings"

	"charm.land/huh/v2"
	tea "charm.land/bubbletea/v2"

	"example/tuibookie/internal/bookmark"
)

var settingsOptions = []string{"Export bookmarks", "Import bookmarks"}

func (m Model) updateSettings(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		m.statusMsg = ""
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc", "left", "h":
			m.currentView = categoryView
		case "up", "k":
			if m.settingsCursor > 0 {
				m.settingsCursor--
			}
		case "down", "j":
			if m.settingsCursor < len(settingsOptions)-1 {
				m.settingsCursor++
			}
		case "enter", "right", "l":
			switch m.settingsCursor {
			case 0: // Export
				filename, err := bookmark.Export(m.bookmarks)
				if err != nil {
					m.statusMsg = "Export failed: " + err.Error()
				} else {
					m.statusMsg = "Exported to " + filename
				}
			case 1: // Import
				m.formAction = formImport
				jsonFiles := findJSONFiles()
				options := make([]huh.Option[string], 0, len(jsonFiles)+1)
				for _, f := range jsonFiles {
					options = append(options, huh.NewOption(f, f))
				}
				options = append(options, huh.NewOption("Enter path manually...", ""))
				m.form = huh.NewForm(
					huh.NewGroup(
						huh.NewSelect[string]().
							Title("Select file to import").
							Key("path").
							Options(options...),
					),
				)
				m.currentView = formView
				return m, m.form.Init()
			}
		}
	}
	return m, nil
}

func (m Model) viewSettings() string {
	var b strings.Builder

	b.WriteString(m.title())
	b.WriteString("\n\n")
	b.WriteString(normalStyle.Render("Settings"))
	b.WriteString("\n\n")

	for i, opt := range settingsOptions {
		if i == m.settingsCursor {
			b.WriteString(selectedStyle.Render("> " + opt))
		} else {
			b.WriteString(normalStyle.Render("  " + opt))
		}
		b.WriteString("\n")
	}

	if m.statusMsg != "" {
		b.WriteString("\n")
		b.WriteString(selectedStyle.Render(m.statusMsg))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("[enter/→] select  [←/esc] back  [q] quit"))

	return b.String()
}

func findJSONFiles() []string {
	entries, err := os.ReadDir(".")
	if err != nil {
		return nil
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			files = append(files, e.Name())
		}
	}
	return files
}
