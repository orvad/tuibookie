package tui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"charm.land/huh/v2"
	tea "charm.land/bubbletea/v2"

	"github.com/orvad/tuibookie/internal/bookmark"
	"github.com/orvad/tuibookie/internal/config"
)

func (m Model) updateForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Allow cancelling the form
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "esc" {
			switch m.formAction {
			case formAddCategory, formEditCategory:
				m.currentView = categoryView
			case formAddBookmark, formEditBookmark, formRunParam:
				m.pendingCmd = ""
				m.pendingParams = nil
				m.paramValues = nil
				m.currentView = bookmarkView
			case formImport, formImportManual, formChangeBookmarksPath,
				formSetGistToken, formSetSharedRepo, formSetSharedFilePath:
				m.pendingConfigPath = ""
				m.pendingGistToken = ""
				m.currentView = settingsView
			}
			m.form = nil
			return m, nil
		}
	}

	if m.form == nil {
		m.currentView = bookmarkView
		return m, nil
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	if m.form.State == huh.StateCompleted {
		name := m.form.GetString("name")
		cmd := m.form.GetString("cmd")

		switch m.formAction {
		case formAddCategory:
			if name != "" {
				if m.isSharedContext {
					bookmark.AddCategory(m.sharedBookmarks, name)
					m.refreshSharedCategories()
					m.saveShared()
					m.selectedCat = name
					m.bmCursor = 0
					m.currentView = bookmarkView
					m.form = nil
					return m, m.pushSharedCmd("add category: " + name)
				}
				bookmark.AddCategory(m.bookmarks, name)
				m.refreshCategories()
				m.save()
				m.selectedCat = name
				m.bmCursor = 0
				m.currentView = bookmarkView
			} else {
				m.currentView = categoryView
			}

		case formEditCategory:
			if name != "" {
				if m.isSharedContext {
					_, shared := m.categoryAtCursor()
					if shared {
						sharedIdx := m.catCursor - len(m.categories)
						oldName := m.sharedCategories[sharedIdx]
						bookmark.RenameCategory(m.sharedBookmarks, oldName, name)
						m.refreshSharedCategories()
						m.saveShared()
						m.currentView = categoryView
						m.form = nil
						return m, m.pushSharedCmd("rename category: " + oldName + " -> " + name)
					}
				}
				oldName := m.categories[m.catCursor]
				bookmark.RenameCategory(m.bookmarks, oldName, name)
				m.refreshCategories()
				m.save()
			}
			m.currentView = categoryView

		case formAddBookmark:
			if name != "" && cmd != "" {
				newBm := bookmark.Bookmark{
					Name:    name,
					Cmd:     cmd,
					Confirm: m.form.GetBool("confirm"),
				}
				if m.isSharedContext {
					bookmark.AddBookmark(m.sharedBookmarks, m.selectedCat, newBm)
					m.saveShared()
					m.currentView = bookmarkView
					m.form = nil
					return m, m.pushSharedCmd("add bookmark: " + name)
				}
				bookmark.AddBookmark(m.bookmarks, m.selectedCat, newBm)
				m.save()
			}
			m.currentView = bookmarkView

		case formEditBookmark:
			if name != "" && cmd != "" {
				updatedBm := bookmark.Bookmark{
					Name:    name,
					Cmd:     cmd,
					Confirm: m.form.GetBool("confirm"),
				}
				if m.isSharedContext {
					bookmark.UpdateBookmark(m.sharedBookmarks, m.selectedCat, m.editIndex, updatedBm)
					m.saveShared()
					m.currentView = bookmarkView
					m.form = nil
					return m, m.pushSharedCmd("edit bookmark: " + name)
				}
				bookmark.UpdateBookmark(m.bookmarks, m.selectedCat, m.editIndex, updatedBm)
				m.save()
			}
			m.currentView = bookmarkView

		case formImport:
			path := m.form.GetString("path")
			if path == "" {
				// User chose "Enter path manually..."
				m.form = huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("Path to JSON file").
							Key("path"),
					),
				).WithTheme(formTheme)
				m.formAction = formImportManual
				return m, m.form.Init()
			}
			if err := bookmark.Import(path, m.bookmarks); err != nil {
				m.statusMsg = "Import failed: " + err.Error()
			} else {
				m.refreshCategories()
				m.save()
				m.statusMsg = "Imported from " + path
			}
			m.currentView = settingsView

		case formImportManual:
			path := m.form.GetString("path")
			if path != "" {
				if err := bookmark.Import(path, m.bookmarks); err != nil {
					m.statusMsg = "Import failed: " + err.Error()
				} else {
					m.refreshCategories()
					m.save()
					m.statusMsg = "Imported from " + path
				}
			}
			m.currentView = settingsView

		case formChangeBookmarksPath:
			path := m.pendingConfigPath
			if path == "" || path == m.configPath {
				m.pendingConfigPath = ""
				m.currentView = settingsView
				break
			}
			// Validate the path
			bm, err := bookmark.Load(path)
			if err != nil {
				if os.IsNotExist(err) {
					m.confirmMsg = "File not found. Create a new empty bookmarks file at this path?"
					m.confirmAction = formConfirmBookmarksPath
					m.confirmCursor = 0
					m.currentView = confirmView
					m.form = nil
					return m, nil
				}
				m.pendingConfigPath = ""
				m.statusMsg = "Invalid file: " + err.Error()
				m.currentView = settingsView
				break
			}
			// File exists and is valid — show confirmation with stats
			cats := bookmark.Categories(bm)
			totalBookmarks := 0
			for _, items := range bm {
				totalBookmarks += len(items)
			}
			m.confirmMsg = fmt.Sprintf("Switch to this file? (%d categories, %d bookmarks)", len(cats), totalBookmarks)
			m.confirmAction = formConfirmBookmarksPath
			m.confirmCursor = 0
			m.currentView = confirmView
			m.form = nil
			return m, nil

		case formSetGistToken:
			token := m.form.GetString("token")
			m.gistToken = token
			m.saveGistConfig()
			if token == "" {
				m.statusMsg = "Token removed"
			} else {
				m.statusMsg = "Token saved"
			}
			m.pendingGistToken = ""
			m.currentView = settingsView

		case formSetSharedRepo:
			url := m.form.GetString("url")
			m.sharedRepoURL = url
			appCfg, _ := config.LoadAppConfig(m.configDir)
			appCfg.SharedRepo = url
			if err := config.SaveAppConfig(m.configDir, appCfg); err != nil {
				m.statusMsg = "Error saving config: " + err.Error()
			} else if url != "" {
				m.statusMsg = "Shared repo saved — use Sync to connect"
			} else {
				m.statusMsg = "Shared repo removed"
			}
			m.currentView = settingsView

		case formSetSharedFilePath:
			path := m.form.GetString("path")
			if path == "" {
				path = "bookmarks.json"
			}
			m.sharedFilePath = path
			appCfg, _ := config.LoadAppConfig(m.configDir)
			appCfg.SharedFilePath = path
			if err := config.SaveAppConfig(m.configDir, appCfg); err != nil {
				m.statusMsg = "Error saving config: " + err.Error()
			} else {
				m.statusMsg = "Shared file path saved"
			}
			m.currentView = settingsView

		case formRunParam:
			values := make(map[string]string)
			for _, p := range m.pendingParams {
				values[p.Name] = m.form.GetString(p.Name)
			}
			resolved := bookmark.ResolveParams(m.pendingCmd, values)
			m.pendingParams = nil
			m.paramValues = nil
			items := m.bookmarks[m.selectedCat]
			if len(items) > m.bmCursor && items[m.bmCursor].Confirm {
				m.pendingCmd = resolved
				m.confirmMsg = "Execute: " + resolved + "?"
				m.confirmAction = formConfirmExec
				m.confirmCursor = 0
				m.currentView = confirmView
				m.form = nil
				return m, nil
			}
			m.pendingCmd = ""
			m.currentView = bookmarkView
			parts := strings.Fields(resolved)
			if len(parts) > 0 {
				c := exec.Command(parts[0], parts[1:]...)
				return m, tea.ExecProcess(c, func(err error) tea.Msg {
					return execDoneMsg{err: err}
				})
			}

		}
		m.form = nil
	}

	return m, cmd
}

func (m Model) viewForm() string {
	if m.form == nil {
		return ""
	}
	s := m.title() + "\n\n"
	if m.formAction == formRunParam && m.pendingCmd != "" {
		label := selectedStyle.Render("Command: ")
		s += "  " + label + renderLiveCmd(m.pendingCmd, m.paramValues) + "\n\n"
	}
	s += m.form.View()
	return s
}
