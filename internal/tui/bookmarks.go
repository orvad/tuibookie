package tui

import (
	"os/exec"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"

	"github.com/orvad/tuibookie/internal/bookmark"
)

// currentBookmarkItems returns the bookmark list for the currently selected category,
// routing to either local or shared bookmarks based on context.
func (m Model) currentBookmarkItems() []bookmark.Bookmark {
	if m.isSharedContext {
		return m.sharedBookmarks[m.selectedCat]
	}
	return m.bookmarks[m.selectedCat]
}

func (m Model) updateBookmark(msg tea.Msg) (tea.Model, tea.Cmd) {
	items := m.currentBookmarkItems()

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc", "left", "h":
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
				bm := items[m.bmCursor]
				cmd := bm.Cmd
				params := bookmark.ParseParams(cmd)
				if len(params) > 0 {
					m.pendingCmd = cmd
					m.pendingParams = params
					m.formAction = formRunParam
					m.paramValues = make(map[string]*string)
					groups := make([]huh.Field, len(params))
					for i, p := range params {
						val := p.Default
						m.paramValues[p.Name] = &val
						groups[i] = huh.NewInput().
							Title(p.Name).
							Key(p.Name).
							Value(&val)
					}
					m.form = huh.NewForm(
						huh.NewGroup(groups...),
					).WithTheme(formTheme)
					m.currentView = formView
					return m, m.form.Init()
				}
				if bm.Confirm {
					m.pendingCmd = cmd
					m.confirmMsg = "Execute: " + cmd + "?"
					m.confirmAction = formConfirmExec
					m.confirmCursor = 0
					m.currentView = confirmView
					return m, nil
				}
				parts := strings.Fields(cmd)
				if len(parts) > 0 {
					m.executedCmd = cmd
					c := exec.Command(parts[0], parts[1:]...)
					return m, tea.ExecProcess(c, func(err error) tea.Msg {
						return execDoneMsg{err: err}
					})
				}
			}
		case "a":
			if m.isSharedContext && m.sharedReadOnly {
				m.statusMsg = "Shared bookmarks are read-only"
				return m, nil
			}
			addConfirm := false
			m.formAction = formAddBookmark
			m.form = huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Bookmark name").
						Key("name"),
					huh.NewInput().
						Title("Command").
						Key("cmd"),
					huh.NewConfirm().
						Title("Confirm before execute?").
						Key("confirm").
						Value(&addConfirm),
				),
			).WithTheme(formTheme)
			m.currentView = formView
			return m, m.form.Init()
		case "e":
			if len(items) > 0 {
				if m.isSharedContext && m.sharedReadOnly {
					m.statusMsg = "Shared bookmarks are read-only"
					return m, nil
				}
				bm := items[m.bmCursor]
				editName := bm.Name
				editCmd := bm.Cmd
				editConfirm := bm.Confirm
				m.editIndex = m.bmCursor
				m.formAction = formEditBookmark
				m.form = huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("Bookmark name").
							Key("name").
							Value(&editName),
						huh.NewInput().
							Title("Command").
							Key("cmd").
							Value(&editCmd),
						huh.NewConfirm().
							Title("Confirm before execute?").
							Key("confirm").
							Value(&editConfirm),
					),
				).WithTheme(formTheme)
				m.currentView = formView
				return m, m.form.Init()
			}
		case "d":
			if len(items) > 0 {
				if m.isSharedContext && m.sharedReadOnly {
					m.statusMsg = "Shared bookmarks are read-only"
					return m, nil
				}
				if m.isSharedContext {
					bookmark.DeleteBookmark(m.sharedBookmarks, m.selectedCat, m.bmCursor)
					m.saveShared()
					cmd := m.pushSharedCmd("delete bookmark from " + m.selectedCat)
					if m.bmCursor >= len(m.sharedBookmarks[m.selectedCat]) && m.bmCursor > 0 {
						m.bmCursor--
					}
					return m, cmd
				}
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

func renderLiveCmd(cmd string, values map[string]*string) string {
	parts := bookmark.ParamRegex.Split(cmd, -1)
	matches := bookmark.ParamRegex.FindAllStringSubmatch(cmd, -1)
	var b strings.Builder
	for i, part := range parts {
		b.WriteString(dimStyle.Render(part))
		if i < len(matches) {
			name := matches[i][1]
			if ptr, ok := values[name]; ok && *ptr != "" {
				b.WriteString(paramStyle.Render(*ptr))
			} else {
				b.WriteString(paramStyle.Render(name))
			}
		}
	}
	return b.String()
}

func renderCmd(cmd string) string {
	parts := bookmark.ParamRegex.Split(cmd, -1)
	matches := bookmark.ParamRegex.FindAllStringSubmatch(cmd, -1)
	var b strings.Builder
	for i, part := range parts {
		b.WriteString(dimStyle.Render(part))
		if i < len(matches) {
			label := matches[i][1]
			if matches[i][2] != "" {
				label = matches[i][2]
			}
			b.WriteString(paramStyle.Render(label))
		}
	}
	return b.String()
}

func (m Model) viewBookmark() string {
	var b strings.Builder

	b.WriteString(m.title())
	b.WriteString("\n\n")

	if m.hasBothSections() {
		prefix := "LOCAL"
		if m.isSharedContext {
			prefix = "SHARED"
		}
		b.WriteString(headingStyle.Render("  "+prefix) + dimStyle.Render(" › ") + headingStyle.Render(strings.ToUpper(m.selectedCat)))
	} else {
		b.WriteString(headingStyle.Render("  " + strings.ToUpper(m.selectedCat)))
	}
	b.WriteString("\n\n")

	items := m.currentBookmarkItems()

	if len(items) == 0 {
		if m.isSharedContext && m.sharedReadOnly {
			b.WriteString(dimStyle.Render("No bookmarks in this shared category."))
		} else {
			b.WriteString(dimStyle.Render("No bookmarks yet. Press 'a' to add one."))
		}
		b.WriteString("\n")
	} else {
		for i, bm := range items {
			indicator := ""
			if bm.Confirm {
				indicator = " " + confirmIndicatorStyle.Render("!")
			}
			if i == m.bmCursor {
				b.WriteString(selectedStyle.Render("> "+bm.Name) + indicator + "  " + renderCmd(bm.Cmd))
			} else {
				b.WriteString(normalStyle.Render("  "+bm.Name) + indicator + "  " + renderCmd(bm.Cmd))
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("\n\n")
	b.WriteString(renderHelp("[a]dd  [e]dit  [d]elete  [enter] run  [←/esc] back  [q]uit"))

	return b.String()
}
