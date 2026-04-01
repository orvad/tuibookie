package tui

import (
	"os"
	"path/filepath"
	"strings"

	"charm.land/huh/v2"
	tea "charm.land/bubbletea/v2"

	"example/tuibookie/internal/bookmark"
)

// settingsItems are the selectable items in order. Section labels are rendered
// separately in viewSettings and are not part of this list.
var settingsItems = []string{
	"Bookmarks file",
	"Export bookmarks",
	"Import bookmarks",
}

// sectionBreak is the index where the DATA section starts (after Config items).
const sectionBreak = 1

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
			if m.settingsCursor < len(settingsItems)-1 {
				m.settingsCursor++
			}
		case "enter", "right", "l":
			switch m.settingsCursor {
			case 0: // Bookmarks file
				if m.pathSource == PathSourceFlag {
					m.statusMsg = "Path set via --config flag"
				} else if m.pathSource == PathSourceEnv {
					m.statusMsg = "Path set via TUIBOOKIE_CONFIG env"
				} else {
					m.formAction = formChangeBookmarksPath
					m.pendingConfigPath = m.configPath
					m.form = huh.NewForm(
						huh.NewGroup(
							huh.NewInput().
								Title("Bookmarks file path").
								Key("path").
								Value(&m.pendingConfigPath),
						),
					)
					m.currentView = formView
					return m, m.form.Init()
				}
			case 1: // Export
				filename, err := bookmark.Export(m.bookmarks)
				if err != nil {
					m.statusMsg = "Export failed: " + err.Error()
				} else {
					m.statusMsg = "Exported to " + filename
				}
			case 2: // Import
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
	b.WriteString(headingStyle.Render("  SETTINGS"))
	b.WriteString("\n\n")

	// CONFIG section
	b.WriteString(dimStyle.Render("  CONFIG"))
	b.WriteString("\n")
	for i := 0; i < sectionBreak; i++ {
		label := settingsItems[i]
		if i == 0 {
			label = "Bookmarks file: " + truncatePath(m.configPath, 40)
		}
		if i == m.settingsCursor {
			b.WriteString(selectedStyle.Render("> " + label))
		} else {
			b.WriteString(normalStyle.Render("  " + label))
		}
		b.WriteString("\n")
	}

	// DATA section
	b.WriteString("\n")
	b.WriteString(dimStyle.Render("  DATA"))
	b.WriteString("\n")
	for i := sectionBreak; i < len(settingsItems); i++ {
		if i == m.settingsCursor {
			b.WriteString(selectedStyle.Render("> " + settingsItems[i]))
		} else {
			b.WriteString(normalStyle.Render("  " + settingsItems[i]))
		}
		b.WriteString("\n")
	}

	if m.statusMsg != "" {
		b.WriteString("\n")
		b.WriteString(selectedStyle.Render(m.statusMsg))
		b.WriteString("\n")
	}

	b.WriteString("\n\n")
	b.WriteString(renderHelp("[enter/→] select  [←/esc] back  [q] quit"))

	return b.String()
}

func truncatePath(path string, maxLen int) string {
	if len(path) <= maxLen {
		return path
	}
	return "..." + path[len(path)-maxLen+3:]
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
