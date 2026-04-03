package tui

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"

	tea "charm.land/bubbletea/v2"

	"example/tuibookie/internal/bookmark"
	"example/tuibookie/internal/config"
	"example/tuibookie/internal/gist"
)

var confirmOptions = []string{"Yes", "No"}

func (m Model) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc", "left", "h":
			if m.confirmAction == formConfirmExec {
				m.pendingCmd = ""
				m.currentView = bookmarkView
			} else {
				m.currentView = settingsView
			}
		case "up", "k":
			if m.confirmCursor > 0 {
				m.confirmCursor--
			}
		case "down", "j":
			if m.confirmCursor < len(confirmOptions)-1 {
				m.confirmCursor++
			}
		case "y":
			m.confirmCursor = 0
			return m.resolveConfirm()
		case "n":
			m.confirmCursor = 1
			return m.resolveConfirm()
		case "enter":
			return m.resolveConfirm()
		}
	}
	return m, nil
}

func (m Model) resolveConfirm() (tea.Model, tea.Cmd) {
	confirmed := m.confirmCursor == 0
	if !confirmed {
		if m.confirmAction == formConfirmExec {
			m.pendingCmd = ""
			m.currentView = bookmarkView
		} else {
			m.pendingConfigPath = ""
			m.currentView = settingsView
		}
		return m, nil
	}
	return m.onConfirm()
}

func (m Model) onConfirm() (tea.Model, tea.Cmd) {
	switch m.confirmAction {
	case formConfirmBookmarksPath:
		path := m.pendingConfigPath
		m.pendingConfigPath = ""

		// Create file if it doesn't exist
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := config.EnsureConfigDir(path); err != nil {
				m.statusMsg = "Error creating directory: " + err.Error()
				m.currentView = settingsView
				return m, nil
			}
			if err := bookmark.Save(path, bookmark.Bookmarks{}); err != nil {
				m.statusMsg = "Error creating file: " + err.Error()
				m.currentView = settingsView
				return m, nil
			}
		}

		bm, err := bookmark.Load(path)
		if err != nil {
			m.statusMsg = "Error loading bookmarks: " + err.Error()
			m.currentView = settingsView
			return m, nil
		}

		appCfg, _ := config.LoadAppConfig(m.configDir)
		appCfg.BookmarksPath = path
		if err := config.SaveAppConfig(m.configDir, appCfg); err != nil {
			m.statusMsg = "Error saving config: " + err.Error()
			m.currentView = settingsView
			return m, nil
		}

		m.configPath = path
		m.bookmarks = bm
		m.pathSource = PathSourceConfig
		m.refreshCategories()
		m.catCursor = 0
		m.bmCursor = 0
		m.statusMsg = "Switched to " + path
		m.currentView = settingsView

	case formConfirmPull:
		c := &gist.Client{Token: m.gistToken}
		data, err := c.Fetch(m.gistID)
		if err != nil {
			m.statusMsg = "Pull failed: " + err.Error()
			m.currentView = settingsView
			return m, nil
		}
		var pulledBm bookmark.Bookmarks
		if err := json.Unmarshal(data, &pulledBm); err != nil {
			m.statusMsg = "Gist contains invalid bookmark data"
			m.currentView = settingsView
			return m, nil
		}
		if err := bookmark.Save(m.configPath, pulledBm); err != nil {
			m.statusMsg = "Failed to save: " + err.Error()
			m.currentView = settingsView
			return m, nil
		}
		m.bookmarks = pulledBm
		m.refreshCategories()
		m.catCursor = 0
		m.bmCursor = 0
		m.statusMsg = "Pulled from gist"
		m.currentView = settingsView

	case formConfirmExec:
		cmd := m.pendingCmd
		m.pendingCmd = ""
		parts := strings.Fields(cmd)
		if len(parts) > 0 {
			c := exec.Command(parts[0], parts[1:]...)
			return m, tea.ExecProcess(c, func(err error) tea.Msg {
				return execDoneMsg{err: err}
			})
		}
		m.currentView = bookmarkView
	}

	return m, nil
}

func (m Model) viewConfirm() string {
	var b strings.Builder

	b.WriteString(m.title())
	b.WriteString("\n\n")
	b.WriteString(normalStyle.Render("  " + m.confirmMsg))
	b.WriteString("\n\n")

	for i, opt := range confirmOptions {
		if i == m.confirmCursor {
			b.WriteString(selectedStyle.Render("> " + opt))
		} else {
			b.WriteString(normalStyle.Render("  " + opt))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n\n")
	b.WriteString(renderHelp("[enter] select  [y] yes  [n] no  [←/esc] back  [q] quit"))

	return b.String()
}
