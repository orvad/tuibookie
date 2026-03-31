package tui

import (
	"charm.land/huh/v2"
	tea "charm.land/bubbletea/v2"

	"example/go-ssh/internal/bookmark"
)

func (m Model) updateForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Allow cancelling the form
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "esc" {
			switch m.formAction {
			case formAddCategory:
				m.currentView = categoryView
			case formAddBookmark, formEditBookmark:
				m.currentView = bookmarkView
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
		switch m.formAction {
		case formAddCategory:
			if m.formName != "" {
				bookmark.AddCategory(m.bookmarks, m.formName)
				m.refreshCategories()
				m.save()
			}
			m.currentView = categoryView

		case formAddBookmark:
			if m.formName != "" && m.formCmd != "" {
				bookmark.AddBookmark(m.bookmarks, m.selectedCat, bookmark.Bookmark{
					Name: m.formName,
					Cmd:  m.formCmd,
				})
				m.save()
			}
			m.currentView = bookmarkView

		case formEditBookmark:
			if m.formName != "" && m.formCmd != "" {
				bookmark.UpdateBookmark(m.bookmarks, m.selectedCat, m.editIndex, bookmark.Bookmark{
					Name: m.formName,
					Cmd:  m.formCmd,
				})
				m.save()
			}
			m.currentView = bookmarkView
		}
		m.form = nil
	}

	return m, cmd
}

func (m Model) viewForm() string {
	if m.form == nil {
		return ""
	}
	return m.form.View()
}
