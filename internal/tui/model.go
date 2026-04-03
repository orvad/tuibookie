package tui

import (
	"fmt"
	"os"
	"strings"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	tea "charm.land/bubbletea/v2"

	"github.com/orvad/tuibookie/internal/bookmark"
	"github.com/orvad/tuibookie/internal/config"
)

type view int

const (
	categoryView view = iota
	bookmarkView
	settingsView
	formView
	confirmView
)

type formAction int

const (
	formAddCategory formAction = iota
	formEditCategory
	formAddBookmark
	formEditBookmark
	formImport
	formImportManual
	formChangeBookmarksPath
	formConfirmBookmarksPath
	formSetGistToken
	formConfirmPull
	formRunParam
	formConfirmExec
)

type PathSource int

const (
	PathSourceDefault PathSource = iota
	PathSourceConfig
	PathSourceEnv
	PathSourceFlag
)

type Model struct {
	bookmarks   bookmark.Bookmarks
	configPath  string
	configDir   string
	pathSource  PathSource
	version     string
	currentView view
	categories  []string
	catCursor   int
	selectedCat string
	bmCursor    int
	form        *huh.Form
	formAction  formAction
	editIndex         int
	settingsCursor    int
	err               error
	statusMsg         string
	pendingConfigPath string
	gistToken         string
	gistID            string
	pendingGistToken  string
	confirmMsg        string
	confirmAction     formAction
	confirmCursor     int // 0=Yes, 1=No
	pendingCmd        string
	pendingParams     []bookmark.Param
	paramValues       map[string]*string
	themeSetting      string // "auto", "dark", "light"
	isDark            bool
	autoDetectedDark  bool   // cached terminal detection from startup
	width             int
	height            int
}

func NewModel(bm bookmark.Bookmarks, configPath string, configDir string, pathSource PathSource, version string) Model {
	cats := bookmark.Categories(bm)
	var gistToken, gistID, themeSetting string
	if appCfg, err := config.LoadAppConfig(configDir); err == nil {
		gistToken = appCfg.GistToken
		gistID = appCfg.GistID
		themeSetting = appCfg.Theme
	}
	if themeSetting == "" {
		themeSetting = "auto"
	}
	// Detect terminal background once before Bubble Tea takes over stdin.
	autoDetectedDark := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	isDark := resolveTheme(themeSetting, autoDetectedDark)
	ApplyTheme(isDark)
	return Model{
		bookmarks:        bm,
		configPath:       configPath,
		configDir:        configDir,
		pathSource:       pathSource,
		version:          version,
		currentView:      categoryView,
		categories:       cats,
		gistToken:        gistToken,
		gistID:           gistID,
		themeSetting:     themeSetting,
		isDark:           isDark,
		autoDetectedDark: autoDetectedDark,
	}
}

func resolveTheme(setting string, autoDetectedDark bool) bool {
	switch setting {
	case "dark":
		return true
	case "light":
		return false
	default:
		return autoDetectedDark
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case execDoneMsg:
		if msg.err != nil {
			m.err = msg.err
		}
		return m, tea.Quit
	}

	switch m.currentView {
	case categoryView:
		return m.updateCategory(msg)
	case bookmarkView:
		return m.updateBookmark(msg)
	case settingsView:
		return m.updateSettings(msg)
	case formView:
		return m.updateForm(msg)
	case confirmView:
		return m.updateConfirm(msg)
	}

	return m, nil
}

func (m Model) View() tea.View {
	var v tea.View
	switch m.currentView {
	case categoryView:
		v = tea.NewView(m.viewCategory())
	case bookmarkView:
		v = tea.NewView(m.viewBookmark())
	case settingsView:
		v = tea.NewView(m.viewSettings())
	case confirmView:
		v = tea.NewView(m.viewConfirm())
	case formView:
		v = tea.NewView(m.viewForm())
	default:
		v = tea.NewView("")
	}
	v.AltScreen = true
	return v
}

func (m Model) title() string {
	name := titlePrefixStyle.Render("Tui") + titleAccentStyle.Render("Bookie")
	ver := titleVersionStyle.Render(" " + m.version)
	left := "TuiBookie " + m.version
	sep := titleSepStyle.Render(strings.Repeat("━", max(0, m.width-4)))

	if m.statusMsg == "" {
		return fmt.Sprintf("\n  %s%s\n  %s", name, ver, sep)
	}

	styledMsg := statusMsgStyle.Render(m.statusMsg)
	gap := max(1, m.width-4-len(left)-len(m.statusMsg))
	return fmt.Sprintf("\n  %s%s%s%s\n  %s", name, ver, strings.Repeat(" ", gap), styledMsg, sep)
}

func (m *Model) refreshCategories() {
	m.categories = bookmark.Categories(m.bookmarks)
}

func (m *Model) save() {
	m.err = bookmark.Save(m.configPath, m.bookmarks)
}

