package tui

import (
	"strings"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

var (
	titleNameStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F92672"))

	titleVersionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#75715E"))

	titleSepStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#75715E"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#75715E"))

	keyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E6DB74"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A6E22E")).
			Bold(true)

	headingStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#A6E22E"))

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F2"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A59F85"))

	statusMsgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A6E22E"))
)

func monokaiTheme(_ bool) *huh.Styles {
	t := huh.ThemeBase(true)

	accent := lipgloss.Color("#F92672")
	green := lipgloss.Color("#A6E22E")
	yellow := lipgloss.Color("#E6DB74")
	text := lipgloss.Color("#F8F8F2")
	muted := lipgloss.Color("#A59F85")
	help := lipgloss.Color("#75715E")
	red := lipgloss.Color("#F92672")

	t.Focused.Base = t.Focused.Base.BorderForeground(help)
	t.Focused.Title = t.Focused.Title.Foreground(accent).Bold(true)
	t.Focused.Description = t.Focused.Description.Foreground(muted)
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(red)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(red)
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(green)
	t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(green)
	t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(green)
	t.Focused.Option = t.Focused.Option.Foreground(text)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(green)
	t.Focused.SelectedPrefix = t.Focused.SelectedPrefix.Foreground(green)
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(text)
	t.Focused.UnselectedPrefix = t.Focused.UnselectedPrefix.Foreground(muted)
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(text).Background(accent)
	t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(muted).Background(lipgloss.Color("#1A1A1A"))

	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(green)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(help)
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(accent)
	t.Focused.TextInput.Text = t.Focused.TextInput.Text.Foreground(text)

	t.Blurred = t.Focused
	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())

	t.Help.ShortKey = t.Help.ShortKey.Foreground(yellow)
	t.Help.ShortDesc = t.Help.ShortDesc.Foreground(help)
	t.Help.ShortSeparator = t.Help.ShortSeparator.Foreground(help)
	t.Help.FullKey = t.Help.FullKey.Foreground(yellow)
	t.Help.FullDesc = t.Help.FullDesc.Foreground(help)
	t.Help.FullSeparator = t.Help.FullSeparator.Foreground(help)
	t.Help.Ellipsis = t.Help.Ellipsis.Foreground(help)

	t.Group.Title = t.Focused.Title
	t.Group.Description = t.Focused.Description
	return t
}

var formTheme = huh.ThemeFunc(monokaiTheme)

func renderHelp(text string) string {
	var b strings.Builder
	inBracket := false
	for i := 0; i < len(text); i++ {
		if text[i] == '[' {
			inBracket = true
			j := strings.IndexByte(text[i:], ']')
			if j >= 0 {
				b.WriteString(keyStyle.Render(text[i : i+j+1]))
				i += j
				inBracket = false
				continue
			}
		}
		if !inBracket {
			b.WriteString(helpStyle.Render(string(text[i])))
		}
	}
	return b.String()
}
