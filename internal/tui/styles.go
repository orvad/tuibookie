package tui

import (
	"image/color"
	"strings"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

var (
	titlePrefixStyle     lipgloss.Style
	titleAccentStyle     lipgloss.Style
	titleVersionStyle    lipgloss.Style
	titleSepStyle        lipgloss.Style
	helpStyle            lipgloss.Style
	keyStyle             lipgloss.Style
	selectedStyle        lipgloss.Style
	headingStyle         lipgloss.Style
	normalStyle          lipgloss.Style
	dimStyle             lipgloss.Style
	statusMsgStyle       lipgloss.Style
	statusErrorStyle     lipgloss.Style
	sectionPillStyle     lipgloss.Style
	sectionAccentStyle   lipgloss.Style
	paramStyle           lipgloss.Style
	confirmIndicatorStyle lipgloss.Style
	sectionHeaderStyle   lipgloss.Style
	formTheme            huh.Theme
)

func init() {
	// Default to dark until ApplyTheme is called.
	applyDarkTheme()
}

// ApplyTheme sets all style variables for dark or light mode.
func ApplyTheme(isDark bool) {
	if isDark {
		applyDarkTheme()
	} else {
		applyLightTheme()
	}
}

func applyDarkTheme() {
	accent := lipgloss.Color("#F92672")
	green := lipgloss.Color("#A6E22E")
	yellow := lipgloss.Color("#E6DB74")
	text := lipgloss.Color("#F8F8F2")
	muted := lipgloss.Color("#A59F85")
	help := lipgloss.Color("#75715E")
	param := lipgloss.Color("#C87A1A")
	blurredBg := lipgloss.Color("#1A1A1A")
	titlePrefix := lipgloss.Color("#E8E8E2")

	applyPalette(accent, green, yellow, text, muted, help, param, blurredBg, titlePrefix, true)
}

func applyLightTheme() {
	accent := lipgloss.Color("#F92672")
	green := lipgloss.Color("#629755")
	yellow := lipgloss.Color("#B8860B")
	text := lipgloss.Color("#2E2E2E")
	muted := lipgloss.Color("#7A7A7A")
	help := lipgloss.Color("#9E9E9E")
	param := lipgloss.Color("#B35900")
	blurredBg := lipgloss.Color("#E8E8E8")
	titlePrefix := lipgloss.Color("#3E3D32")

	applyPalette(accent, green, yellow, text, muted, help, param, blurredBg, titlePrefix, false)
}

func applyPalette(accent, green, yellow, text, muted, help, param, blurredBg, titlePrefix color.Color, isDark bool) {
	titlePrefixStyle = lipgloss.NewStyle().Bold(true).Foreground(titlePrefix)
	titleAccentStyle = lipgloss.NewStyle().Bold(true).Foreground(accent)
	titleVersionStyle = lipgloss.NewStyle().Foreground(help)
	titleSepStyle = lipgloss.NewStyle().Foreground(help)
	helpStyle = lipgloss.NewStyle().Foreground(help)
	keyStyle = lipgloss.NewStyle().Foreground(yellow)
	selectedStyle = lipgloss.NewStyle().Foreground(green).Bold(true)
	headingStyle = lipgloss.NewStyle().Bold(true).Foreground(green)
	normalStyle = lipgloss.NewStyle().Foreground(text)
	dimStyle = lipgloss.NewStyle().Foreground(muted)
	statusMsgStyle = lipgloss.NewStyle().Foreground(green)
	statusErrorStyle = lipgloss.NewStyle().Foreground(accent).Bold(true)
	sectionPillStyle = lipgloss.NewStyle().Foreground(text).Background(blurredBg).Bold(true)
	sectionAccentStyle = lipgloss.NewStyle().Foreground(accent)
	paramStyle = lipgloss.NewStyle().Foreground(param)
	confirmIndicatorStyle = lipgloss.NewStyle().Foreground(accent).Bold(true)
	sectionHeaderStyle = lipgloss.NewStyle().Foreground(muted).Bold(true)

	formTheme = huh.ThemeFunc(func(_ bool) *huh.Styles {
		t := huh.ThemeBase(isDark)

		t.Focused.Base = t.Focused.Base.BorderForeground(help)
		t.Focused.Title = t.Focused.Title.Foreground(accent).Bold(true)
		t.Focused.Description = t.Focused.Description.Foreground(muted)
		t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(accent)
		t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(accent)
		t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(green)
		t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(green)
		t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(green)
		t.Focused.Option = t.Focused.Option.Foreground(text)
		t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(green)
		t.Focused.SelectedPrefix = t.Focused.SelectedPrefix.Foreground(green)
		t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(text)
		t.Focused.UnselectedPrefix = t.Focused.UnselectedPrefix.Foreground(muted)
		t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(text).Background(accent)
		t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(muted).Background(blurredBg)

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
	})
}

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
