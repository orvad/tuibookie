package tui

import (
	"fmt"
	"strings"

	"charm.land/huh/v2"
	tea "charm.land/bubbletea/v2"

	"example/tuibookie/internal/bookmark"
)

func (m Model) updateCategory(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		m.statusMsg = ""
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.catCursor > 0 {
				m.catCursor--
			}
		case "down", "j":
			if m.catCursor < len(m.categories)-1 {
				m.catCursor++
			}
		case "enter", "right", "l":
			if len(m.categories) > 0 {
				m.selectedCat = m.categories[m.catCursor]
				m.bmCursor = 0
				m.currentView = bookmarkView
			}
		case "a":
			m.formAction = formAddCategory
			m.form = huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Category name").
						Key("name"),
				),
			).WithTheme(formTheme)
			m.currentView = formView
			return m, m.form.Init()
		case "e":
			if len(m.categories) > 0 {
				editCatName := m.categories[m.catCursor]
				m.formAction = formEditCategory
				m.form = huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("Rename category").
							Key("name").
							Value(&editCatName),
					),
				).WithTheme(formTheme)
				m.currentView = formView
				return m, m.form.Init()
			}
		case "d":
			if len(m.categories) > 0 {
				cat := m.categories[m.catCursor]
				bookmark.DeleteCategory(m.bookmarks, cat)
				m.refreshCategories()
				m.save()
				if m.catCursor >= len(m.categories) && m.catCursor > 0 {
					m.catCursor--
				}
			}
		case "s":
			m.settingsCursor = 0
			m.statusMsg = ""
			m.currentView = settingsView
		}
	}
	return m, nil
}

func (m Model) viewCategory() string {
	var b strings.Builder

	b.WriteString(m.title())
	b.WriteString("\n\n")

	if len(m.categories) == 0 {
		b.WriteString(dimStyle.Render("No categories yet. Press 'a' to add one."))
		b.WriteString("\n")
	} else {
		for i, cat := range m.categories {
			count := len(m.bookmarks[cat])
			label := fmt.Sprintf("%s (%d)", cat, count)

			if i == m.catCursor {
				b.WriteString(selectedStyle.Render("> " + label))
			} else {
				b.WriteString(normalStyle.Render("  " + label))
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("\n\n")
	b.WriteString(renderHelp("[a]dd  [e]dit  [d]elete  [s]ettings  [enter/→] open  [q/esc] quit"))

	return b.String()
}
