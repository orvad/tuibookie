package tui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"charm.land/huh/v2"
	tea "charm.land/bubbletea/v2"

	"github.com/orvad/tuibookie/internal/bookmark"
	"github.com/orvad/tuibookie/internal/config"
	"github.com/orvad/tuibookie/internal/gist"
)

// settingsItems are the selectable items in order. Section labels are rendered
// separately in viewSettings and are not part of this list.
var settingsItems = []string{
	"Bookmarks file",
	"Theme",
	"Export bookmarks",
	"Import bookmarks",
	"Push to Gist",
	"Pull from Gist",
	"GitHub token",
}

// Section boundaries: CONFIG [0,dataBreak), DATA [dataBreak,syncBreak), SYNC [syncBreak,len)
const (
	dataBreak = 2
	syncBreak = 4
)

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
					).WithTheme(formTheme)
					m.currentView = formView
					return m, m.form.Init()
				}
			case 1: // Theme
				m.cycleTheme()
			case 2: // Export
				filename, err := bookmark.Export(m.bookmarks)
				if err != nil {
					m.statusMsg = "Export failed: " + err.Error()
				} else {
					m.statusMsg = "Exported to " + filename
				}
			case 3: // Import
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
				).WithTheme(formTheme)
				m.currentView = formView
				return m, m.form.Init()
			case 4: // Push to Gist
				return m.pushToGist()
			case 5: // Pull from Gist
				return m.pullFromGist()
			case 6: // GitHub token
				m.formAction = formSetGistToken
				m.pendingGistToken = m.gistToken
				m.form = huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("GitHub Personal Access Token").
							Key("token").
							CharLimit(200).
							Value(&m.pendingGistToken),
					),
				).WithTheme(formTheme).WithWidth(m.width)
				m.currentView = formView
				return m, m.form.Init()
			}
		}
	}
	return m, nil
}

func (m Model) pushToGist() (tea.Model, tea.Cmd) {
	if m.gistToken == "" {
		m.statusMsg = "Set GitHub token first"
		return m, nil
	}

	data, err := json.MarshalIndent(m.bookmarks, "", "  ")
	if err != nil {
		m.statusMsg = "Failed to serialize: " + err.Error()
		return m, nil
	}

	c := &gist.Client{Token: m.gistToken}

	if m.gistID == "" {
		id, err := c.Create(data)
		if err != nil {
			m.statusMsg = "Push failed: " + err.Error()
			return m, nil
		}
		m.gistID = id
		m.saveGistConfig()
		m.statusMsg = fmt.Sprintf("Created gist %s", id[:8])
	} else {
		if err := c.Update(m.gistID, data); err != nil {
			m.statusMsg = "Push failed: " + err.Error()
			return m, nil
		}
		m.statusMsg = fmt.Sprintf("Pushed to gist %s", m.gistID[:min(8, len(m.gistID))])
	}
	return m, nil
}

func (m Model) pullFromGist() (tea.Model, tea.Cmd) {
	if m.gistToken == "" {
		m.statusMsg = "Set GitHub token first"
		return m, nil
	}
	if m.gistID == "" {
		m.statusMsg = "Push bookmarks first to create a gist"
		return m, nil
	}

	c := &gist.Client{Token: m.gistToken}
	data, err := c.Fetch(m.gistID)
	if err != nil {
		m.statusMsg = "Pull failed: " + err.Error()
		return m, nil
	}

	var bm bookmark.Bookmarks
	if err := json.Unmarshal(data, &bm); err != nil {
		m.statusMsg = "Gist contains invalid bookmark data"
		return m, nil
	}

	cats := bookmark.Categories(bm)
	totalBookmarks := 0
	for _, items := range bm {
		totalBookmarks += len(items)
	}

	m.confirmMsg = fmt.Sprintf("Replace local bookmarks? (gist has %d categories, %d bookmarks)", len(cats), totalBookmarks)
	m.confirmAction = formConfirmPull
	m.confirmCursor = 0
	m.currentView = confirmView
	return m, nil
}

func (m *Model) saveGistConfig() {
	appCfg, _ := config.LoadAppConfig(m.configDir)
	appCfg.GistToken = m.gistToken
	appCfg.GistID = m.gistID
	config.SaveAppConfig(m.configDir, appCfg)
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
	for i := 0; i < dataBreak; i++ {
		label := settingsItems[i]
		if i == 0 {
			label = "Bookmarks file: " + truncatePath(m.configPath, 40)
		}
		if i == 1 {
			label = "Theme: " + m.themeSetting
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
	for i := dataBreak; i < syncBreak; i++ {
		if i == m.settingsCursor {
			b.WriteString(selectedStyle.Render("> " + settingsItems[i]))
		} else {
			b.WriteString(normalStyle.Render("  " + settingsItems[i]))
		}
		b.WriteString("\n")
	}

	// SYNC section
	b.WriteString("\n")
	b.WriteString(dimStyle.Render("  SYNC"))
	b.WriteString("\n")
	for i := syncBreak; i < len(settingsItems); i++ {
		label := settingsItems[i]
		if i == 6 {
			label = "GitHub token: " + maskToken(m.gistToken)
		}
		if i == m.settingsCursor {
			b.WriteString(selectedStyle.Render("> " + label))
		} else {
			b.WriteString(normalStyle.Render("  " + label))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n\n")
	b.WriteString(renderHelp("[enter/→] select  [←/esc] back  [q] quit"))

	return b.String()
}

func (m *Model) cycleTheme() {
	switch m.themeSetting {
	case "auto":
		m.themeSetting = "dark"
	case "dark":
		m.themeSetting = "light"
	default:
		m.themeSetting = "auto"
	}
	m.isDark = resolveTheme(m.themeSetting, m.autoDetectedDark)
	ApplyTheme(m.isDark)
	appCfg, _ := config.LoadAppConfig(m.configDir)
	appCfg.Theme = m.themeSetting
	config.SaveAppConfig(m.configDir, appCfg)
	m.statusMsg = "Theme: " + m.themeSetting
}

func maskToken(token string) string {
	if token == "" {
		return "(not set)"
	}
	if len(token) <= 4 {
		return "****"
	}
	return "****" + token[len(token)-4:]
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
