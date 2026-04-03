package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	tea "charm.land/bubbletea/v2"

	"github.com/orvad/tuibookie/internal/bookmark"
	"github.com/orvad/tuibookie/internal/config"
	"github.com/orvad/tuibookie/internal/gitrepo"
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

// sharedSyncMsg is sent when an async shared bookmark sync completes.
type sharedSyncMsg struct {
	bookmarks bookmark.Bookmarks
	readOnly  bool
	err       error
}

// sharedPushMsg is sent when an async shared bookmark push completes.
type sharedPushMsg struct {
	err error
}

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
	sharedBookmarks  bookmark.Bookmarks
	sharedCategories []string
	sharedReadOnly   bool
	sharedRepoURL    string
	sharedFilePath   string
	sharedCloneDir   string
	isSharedContext  bool
	syncing          bool
}

func NewModel(bm bookmark.Bookmarks, configPath string, configDir string, pathSource PathSource, version string) Model {
	cats := bookmark.Categories(bm)
	var gistToken, gistID, themeSetting string
	var sharedRepoURL, sharedFilePath string
	var sharedReadOnly bool
	if appCfg, err := config.LoadAppConfig(configDir); err == nil {
		gistToken = appCfg.GistToken
		gistID = appCfg.GistID
		themeSetting = appCfg.Theme
		sharedRepoURL = appCfg.SharedRepo
		sharedFilePath = appCfg.SharedFilePath
		sharedReadOnly = appCfg.SharedReadOnly
	}
	if themeSetting == "" {
		themeSetting = "auto"
	}
	if sharedFilePath == "" {
		sharedFilePath = "bookmarks.json"
	}

	sharedCloneDir := filepath.Join(configDir, "shared-repo")

	// Load cached shared bookmarks from disk (instant, no git)
	var sharedBm bookmark.Bookmarks
	if sharedRepoURL != "" {
		sharedPath := filepath.Join(sharedCloneDir, sharedFilePath)
		if loaded, err := bookmark.Load(sharedPath); err == nil {
			sharedBm = loaded
		}
	}

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
		sharedBookmarks:  sharedBm,
		sharedCategories: bookmark.Categories(sharedBm),
		sharedReadOnly:   sharedReadOnly,
		sharedRepoURL:    sharedRepoURL,
		sharedFilePath:   sharedFilePath,
		sharedCloneDir:   sharedCloneDir,
		syncing:          sharedRepoURL != "",
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
	if m.sharedRepoURL == "" {
		return nil
	}
	// Note: syncing=true is set in NewModel when sharedRepoURL is configured
	return m.syncSharedCmd()
}

func (m Model) syncSharedCmd() tea.Cmd {
	repoURL := m.sharedRepoURL
	cloneDir := m.sharedCloneDir
	filePath := m.sharedFilePath
	return func() tea.Msg {
		if !gitrepo.GitInstalled() {
			return sharedSyncMsg{err: fmt.Errorf("git not found — install git to use shared bookmarks")}
		}

		if !gitrepo.IsCloned(cloneDir) {
			if err := gitrepo.Clone(repoURL, cloneDir); err != nil {
				return sharedSyncMsg{err: fmt.Errorf("clone failed: %w", err)}
			}
		} else {
			if err := gitrepo.Pull(cloneDir); err != nil {
				// Pull failed — try reset if diverged, otherwise report error
				if resetErr := gitrepo.ResetToRemote(cloneDir); resetErr != nil {
					return sharedSyncMsg{err: fmt.Errorf("sync failed: %w", err)}
				}
			}
		}

		bmPath := filepath.Join(cloneDir, filePath)
		bm, err := bookmark.Load(bmPath)
		if err != nil {
			return sharedSyncMsg{err: fmt.Errorf("shared bookmarks file not found at %s", filePath)}
		}

		canPush, _ := gitrepo.CanPush(cloneDir)
		return sharedSyncMsg{bookmarks: bm, readOnly: !canPush}
	}
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
	case sharedSyncMsg:
		m.syncing = false
		if msg.err != nil {
			m.statusMsg = msg.err.Error()
			return m, nil
		}
		m.sharedBookmarks = msg.bookmarks
		m.sharedCategories = bookmark.Categories(msg.bookmarks)
		m.sharedReadOnly = msg.readOnly
		// Persist read-only status to config
		appCfg, _ := config.LoadAppConfig(m.configDir)
		appCfg.SharedReadOnly = msg.readOnly
		config.SaveAppConfig(m.configDir, appCfg)
		m.statusMsg = "Synced"
		return m, nil
	case sharedPushMsg:
		if msg.err != nil {
			m.statusMsg = "Push failed — will retry on next sync"
		}
		return m, nil
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

func (m *Model) refreshSharedCategories() {
	m.sharedCategories = bookmark.Categories(m.sharedBookmarks)
}

func (m *Model) saveShared() {
	bmPath := filepath.Join(m.sharedCloneDir, m.sharedFilePath)
	m.err = bookmark.Save(bmPath, m.sharedBookmarks)
}

func (m Model) pushSharedCmd(commitMsg string) tea.Cmd {
	cloneDir := m.sharedCloneDir
	filePath := m.sharedFilePath
	return func() tea.Msg {
		err := gitrepo.CommitAndPush(cloneDir, filePath, commitMsg)
		return sharedPushMsg{err: err}
	}
}

func (m *Model) hasSharedBookmarks() bool {
	return len(m.sharedCategories) > 0
}

func (m *Model) hasBothSections() bool {
	return len(m.categories) > 0 && m.hasSharedBookmarks()
}

