package tui

import (
	"fmt"
	"os"

	"charm.land/huh/v2"
	tea "charm.land/bubbletea/v2"

	"example/tuibookie/internal/bookmark"
)

func (m Model) updateForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Allow cancelling the form
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "esc" {
			switch m.formAction {
			case formAddCategory, formEditCategory:
				m.currentView = categoryView
			case formAddBookmark, formEditBookmark:
				m.currentView = bookmarkView
			case formImport, formImportManual, formChangeBookmarksPath,
				formSetGistToken:
				m.pendingConfigPath = ""
				m.pendingGistToken = ""
				m.currentView = settingsView
			}
			m.form = nil
			return m, nil
		}
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
				oldName := m.categories[m.catCursor]
				bookmark.RenameCategory(m.bookmarks, oldName, name)
				m.refreshCategories()
				m.save()
			}
			m.currentView = categoryView

		case formAddBookmark:
			if name != "" && cmd != "" {
				bookmark.AddBookmark(m.bookmarks, m.selectedCat, bookmark.Bookmark{
					Name: name,
					Cmd:  cmd,
				})
				m.save()
			}
			m.currentView = bookmarkView

		case formEditBookmark:
			if name != "" && cmd != "" {
				bookmark.UpdateBookmark(m.bookmarks, m.selectedCat, m.editIndex, bookmark.Bookmark{
					Name: name,
					Cmd:  cmd,
				})
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

		}
		m.form = nil
	}

	return m, cmd
}

func (m Model) viewForm() string {
	if m.form == nil {
		return ""
	}
	return m.title() + "\n\n" + m.form.View()
}
