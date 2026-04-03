package tui

import (
	"fmt"
	"strings"

	"charm.land/huh/v2"
	tea "charm.land/bubbletea/v2"

	"github.com/orvad/tuibookie/internal/bookmark"
)

func (m Model) updateCategory(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		m.statusMsg = ""
		m.statusIsError = false
		total := m.totalCategoryItems()

		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.catCursor > 0 {
				m.catCursor--
			}
		case "down", "j":
			if m.catCursor < total-1 {
				m.catCursor++
			}
		case "enter", "right", "l":
			if total > 0 {
				cat, shared := m.categoryAtCursor()
				m.selectedCat = cat
				m.isSharedContext = shared
				m.bmCursor = 0
				m.currentView = bookmarkView
			}
		case "a":
			if m.isSharedContext && m.sharedReadOnly {
				m.statusMsg = "Shared bookmarks are read-only"
				return m, nil
			}
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
			if total > 0 {
				if m.isSharedContext && m.sharedReadOnly {
					m.statusMsg = "Shared bookmarks are read-only"
					return m, nil
				}
				cat, _ := m.categoryAtCursor()
				editCatName := cat
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
			if total > 0 {
				if m.isSharedContext && m.sharedReadOnly {
					m.statusMsg = "Shared bookmarks are read-only"
					return m, nil
				}
				cat, shared := m.categoryAtCursor()
				if shared {
					bookmark.DeleteCategory(m.sharedBookmarks, cat)
					m.refreshSharedCategories()
					m.saveShared()
					if !m.sharedReadOnly {
						return m, m.pushSharedCmd("delete category: " + cat)
					}
				} else {
					bookmark.DeleteCategory(m.bookmarks, cat)
					m.refreshCategories()
					m.save()
				}
				if m.catCursor >= m.totalCategoryItems() && m.catCursor > 0 {
					m.catCursor--
				}
			}
		case "s":
			m.settingsCursor = 0
			m.statusMsg = ""
			m.currentView = settingsView
		case "S":
			if m.sharedRepoURL != "" {
				m.syncing = true
				m.statusMsg = "Syncing..."
				return m, m.syncSharedCmd()
			}
		}

		// Update isSharedContext based on current cursor position
		if total > 0 {
			_, shared := m.categoryAtCursor()
			m.isSharedContext = shared
		}
	}
	return m, nil
}

func (m Model) viewCategory() string {
	var b strings.Builder

	b.WriteString(m.title())
	b.WriteString("\n\n")

	totalItems := m.totalCategoryItems()

	if totalItems == 0 {
		b.WriteString(dimStyle.Render("No categories yet. Press 'a' to add one."))
		b.WriteString("\n")
	} else {
		idx := 0

		// Show local section
		if len(m.categories) > 0 {
			if m.hasBothSections() {
				b.WriteString(m.renderSectionHeader("Local", false))
				b.WriteString("\n")
			}
			for _, cat := range m.categories {
				count := len(m.bookmarks[cat])
				label := fmt.Sprintf("%s (%d)", cat, count)
				if idx == m.catCursor {
					b.WriteString(selectedStyle.Render("> " + label))
				} else {
					b.WriteString(normalStyle.Render("  " + label))
				}
				b.WriteString("\n")
				idx++
			}
		}

		// Show shared section
		if m.hasSharedBookmarks() {
			if m.hasBothSections() {
				b.WriteString("\n")
				b.WriteString(m.renderSectionHeader("Shared", m.sharedReadOnly))
				b.WriteString("\n")
			}
			for _, cat := range m.sharedCategories {
				count := len(m.sharedBookmarks[cat])
				label := fmt.Sprintf("%s (%d)", cat, count)
				if idx == m.catCursor {
					b.WriteString(selectedStyle.Render("> " + label))
				} else {
					b.WriteString(normalStyle.Render("  " + label))
				}
				b.WriteString("\n")
				idx++
			}
		}
	}

	b.WriteString("\n\n")

	helpText := "[a]dd  [e]dit  [d]elete  [s]ettings  [enter/→] open  [q/esc] quit"
	if m.sharedRepoURL != "" {
		helpText = "[a]dd  [e]dit  [d]elete  [s]ettings  [S]ync  [enter/→] open  [q/esc] quit"
	}
	b.WriteString(renderHelp(helpText))

	return b.String()
}

func (m Model) totalCategoryItems() int {
	return len(m.categories) + len(m.sharedCategories)
}

// categoryAtCursor returns the category name and whether it's a shared category
// based on the current cursor position.
func (m Model) categoryAtCursor() (string, bool) {
	if m.catCursor < len(m.categories) {
		return m.categories[m.catCursor], false
	}
	sharedIdx := m.catCursor - len(m.categories)
	if sharedIdx < len(m.sharedCategories) {
		return m.sharedCategories[sharedIdx], true
	}
	return "", false
}

func (m Model) renderSectionHeader(label string, readOnly bool) string {
	if readOnly {
		label += " (read-only)"
	}
	return "  " + helpStyle.Render("◆") + sectionHeaderStyle.Render(" "+strings.ToUpper(label))
}
