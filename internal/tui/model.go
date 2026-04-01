package tui

import (
	"fmt"
	"strings"

	"charm.land/huh/v2"
	tea "charm.land/bubbletea/v2"

	"example/tuibookie/internal/bookmark"
)

type view int

const (
	categoryView view = iota
	bookmarkView
	settingsView
	formView
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
	width             int
	height            int
}

func NewModel(bm bookmark.Bookmarks, configPath string, configDir string, pathSource PathSource, version string) Model {
	cats := bookmark.Categories(bm)
	return Model{
		bookmarks:   bm,
		configPath:  configPath,
		configDir:   configDir,
		pathSource:  pathSource,
		version:     version,
		currentView: categoryView,
		categories:  cats,
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
	case formView:
		v = tea.NewView(m.viewForm())
	default:
		v = tea.NewView("")
	}
	v.AltScreen = true
	return v
}

func (m Model) title() string {
	name := titleNameStyle.Render("◆ TuiBookie")
	ver := titleVersionStyle.Render(" " + m.version)
	sep := titleSepStyle.Render(strings.Repeat("━", max(0, m.width-4)))
	return fmt.Sprintf("\n  %s%s\n  %s", name, ver, sep)
}

func (m *Model) refreshCategories() {
	m.categories = bookmark.Categories(m.bookmarks)
}

func (m *Model) save() {
	m.err = bookmark.Save(m.configPath, m.bookmarks)
}

