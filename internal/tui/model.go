package tui

import (
	"charm.land/huh/v2"
	tea "charm.land/bubbletea/v2"

	"example/go-ssh/internal/bookmark"
)

type view int

const (
	categoryView view = iota
	bookmarkView
	formView
)

type formAction int

const (
	formAddCategory formAction = iota
	formAddBookmark
	formEditBookmark
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
	editIndex   int
	err         error
	width       int
	height      int

	// form field bindings
	formName string
	formCmd  string
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
	}

	switch m.currentView {
	case categoryView:
		return m.updateCategory(msg)
	case bookmarkView:
		return m.updateBookmark(msg)
	case formView:
		return m.updateForm(msg)
	}

	return m, nil
}

func (m Model) View() tea.View {
	switch m.currentView {
	case categoryView:
		return tea.NewView(m.viewCategory())
	case bookmarkView:
		return tea.NewView(m.viewBookmark())
	case formView:
		return tea.NewView(m.viewForm())
	}
	return tea.NewView("")
}

func (m *Model) refreshCategories() {
	m.categories = bookmark.Categories(m.bookmarks)
}

func (m *Model) save() {
	m.err = bookmark.Save(m.configPath, m.bookmarks)
}

// Stubs — implemented in subsequent tasks
func (m Model) updateBookmark(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m Model) updateForm(msg tea.Msg) (tea.Model, tea.Cmd)     { return m, nil }
func (m Model) viewBookmark() string                             { return "" }
func (m Model) viewForm() string                                 { return "" }
