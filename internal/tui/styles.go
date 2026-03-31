package tui

import (
	"strings"

	"charm.land/lipgloss/v2"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F8F8F2")).
			Background(lipgloss.Color("#F92672")).
			Padding(0, 2).
			MarginTop(1).
			MarginLeft(2)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#75715E"))

	keyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E6DB74"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A6E22E")).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F2"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A59F85"))
)

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
