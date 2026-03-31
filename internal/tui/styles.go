package tui

import "charm.land/lipgloss/v2"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F8F8F2")).
			Background(lipgloss.Color("#F92672")).
			Padding(0, 2).
			MarginTop(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#75715E"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A6E22E")).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F2"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A59F85"))
)
