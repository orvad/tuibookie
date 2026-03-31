package tui

import (
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
)

type Model struct {
	bookmarks   bookmark.Bookmarks
	configPath  string
	currentView view
	categories  []string
	catCursor   int
	selectedCat string
	bmCursor    int
	form        *huh.Form
	formAction  formAction
	editIndex      int
	settingsCursor int
	err            error
	statusMsg      string
	width       int
	height      int

}

func NewModel(bm bookmark.Bookmarks, configPath string) Model {
	cats := bookmark.Categories(bm)
	return Model{
		bookmarks:   bm,
		configPath:  configPath,
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

func (m *Model) refreshCategories() {
	m.categories = bookmark.Categories(m.bookmarks)
}

func (m *Model) save() {
	m.err = bookmark.Save(m.configPath, m.bookmarks)
}

