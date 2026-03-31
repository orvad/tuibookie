package tui

import (
	"os/exec"
	"strings"

	"charm.land/huh/v2"
	tea "charm.land/bubbletea/v2"

	"example/go-ssh/internal/bookmark"
)

func (m Model) updateBookmark(msg tea.Msg) (tea.Model, tea.Cmd) {
	items := m.bookmarks[m.selectedCat]

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			m.currentView = categoryView
			m.bmCursor = 0
		case "up", "k":
			if m.bmCursor > 0 {
				m.bmCursor--
			}
		case "down", "j":
			if m.bmCursor < len(items)-1 {
				m.bmCursor++
			}
		case "enter":
			if len(items) > 0 {
				cmd := items[m.bmCursor].Cmd
				parts := strings.Fields(cmd)
				if len(parts) > 0 {
					c := exec.Command(parts[0], parts[1:]...)
					return m, tea.ExecProcess(c, func(err error) tea.Msg {
						return execDoneMsg{err: err}
					})
				}
			}
		case "a":
			m.formName = ""
			m.formCmd = ""
			m.formAction = formAddBookmark
			m.form = huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Bookmark name").
						Value(&m.formName),
					huh.NewInput().
						Title("SSH command").
						Placeholder("ssh user@host").
						Value(&m.formCmd),
				),
			)
			m.currentView = formView
			return m, m.form.Init()
		case "e":
			if len(items) > 0 {
				bm := items[m.bmCursor]
				m.formName = bm.Name
				m.formCmd = bm.Cmd
				m.editIndex = m.bmCursor
				m.formAction = formEditBookmark
				m.form = huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("Bookmark name").
							Value(&m.formName),
						huh.NewInput().
							Title("SSH command").
							Value(&m.formCmd),
					),
				)
				m.currentView = formView
				return m, m.form.Init()
			}
		case "d":
			if len(items) > 0 {
				bookmark.DeleteBookmark(m.bookmarks, m.selectedCat, m.bmCursor)
				m.save()
				if m.bmCursor >= len(m.bookmarks[m.selectedCat]) && m.bmCursor > 0 {
					m.bmCursor--
				}
			}
		}
	}
	return m, nil
}

func (m Model) viewBookmark() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(m.selectedCat))
	b.WriteString("\n\n")

	items := m.bookmarks[m.selectedCat]

	if len(items) == 0 {
		b.WriteString(dimStyle.Render("No bookmarks yet. Press 'a' to add one."))
		b.WriteString("\n")
	} else {
		for i, bm := range items {
			if i == m.bmCursor {
				b.WriteString(selectedStyle.Render("> "+bm.Name) + "  " + dimStyle.Render(bm.Cmd))
			} else {
				b.WriteString(normalStyle.Render("  "+bm.Name) + "  " + dimStyle.Render(bm.Cmd))
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("[a]dd  [e]dit  [d]elete  [enter] connect  [esc] back  [q]uit"))

	return b.String()
}
