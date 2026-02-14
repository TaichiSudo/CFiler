package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cfiler/internal/bookmark"
	"cfiler/internal/dialog"
	"cfiler/internal/fileops"
	"cfiler/internal/pane"
	"cfiler/internal/preview"
	"cfiler/internal/statusbar"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mode int

const (
	modeNormal mode = iota
	modeDialog
	modeSearch
	modeBookmark
	modeHelp
)

type clipAction int

const (
	clipNone clipAction = iota
	clipCopy
	clipMove
)

type App struct {
	leftPane   pane.Model
	rightPane  pane.Model
	activePane int // 0=left, 1=right
	preview    preview.Model
	statusBar  statusbar.Model
	dialog     dialog.Dialog
	bookmarks  bookmark.Model
	searchInput textinput.Model

	mode               mode
	clipboard          []string
	clipAction         clipAction
	pendingDeletePaths []string
	width              int
	height             int
	ready              bool
}

func New() App {
	startDir, err := os.Getwd()
	if err != nil {
		startDir, _ = os.UserHomeDir()
	}

	si := textinput.New()
	si.Placeholder = "search..."
	si.CharLimit = 256

	return App{
		leftPane:    pane.New(0, startDir),
		rightPane:   pane.New(1, startDir),
		activePane:  0,
		preview:     preview.New(),
		statusBar:   statusbar.New(),
		searchInput: si,
	}
}

func (a App) Init() tea.Cmd {
	return tea.Batch(
		pane.LoadDir(0, a.leftPane.Dir()),
		pane.LoadDir(1, a.rightPane.Dir()),
	)
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.ready = true
		a.updateLayout()
		return a, nil

	case pane.DirLoadedMsg:
		if msg.PaneID == 0 {
			a.leftPane.SetDir(msg.Path)
			a.leftPane.SetEntries(msg.Entries)
		} else {
			a.rightPane.SetDir(msg.Path)
			a.rightPane.SetEntries(msg.Entries)
		}
		cmds = append(cmds, a.loadPreviewCmd())
		return a, tea.Batch(cmds...)

	case pane.DirLoadErrorMsg:
		if msg.PaneID == 0 {
			a.leftPane.SetError(msg.Err)
		} else {
			a.rightPane.SetError(msg.Err)
		}
		a.statusBar.SetMessage(fmt.Sprintf("Error: %v", msg.Err), true)
		return a, nil

	case preview.LoadMsg:
		a.preview.SetContent(msg.Path, msg.Content, msg.IsBinary)
		return a, nil

	case dialog.ResultMsg:
		a.mode = modeNormal
		a.dialog = nil
		if msg.Confirmed {
			cmd := a.handleDialogResult(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		return a, tea.Batch(cmds...)

	case FileOpResultMsg:
		if msg.Err != nil {
			a.statusBar.SetMessage(fmt.Sprintf("%s failed: %v", msg.Op, msg.Err), true)
		} else {
			a.statusBar.SetMessage(fmt.Sprintf("%s completed", msg.Op), false)
		}
		// Reload both panes
		cmds = append(cmds,
			pane.LoadDir(0, a.leftPane.Dir()),
			pane.LoadDir(1, a.rightPane.Dir()),
		)
		return a, tea.Batch(cmds...)

	case bookmark.SelectMsg:
		a.mode = modeNormal
		active := a.getActivePane()
		active.SetDir(msg.Path)
		cmds = append(cmds, pane.LoadDir(active.ID(), msg.Path))
		return a, tea.Batch(cmds...)

	case bookmark.CloseMsg:
		a.mode = modeNormal
		return a, nil

	case tea.KeyMsg:
		return a.handleKey(msg)
	}

	return a, nil
}

func (a App) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch a.mode {
	case modeDialog:
		return a.handleDialogKey(msg)
	case modeSearch:
		return a.handleSearchKey(msg)
	case modeBookmark:
		return a.handleBookmarkKey(msg)
	case modeHelp:
		if msg.String() == "esc" || msg.String() == "?" || msg.String() == "q" {
			a.mode = modeNormal
		}
		return a, nil
	default:
		return a.handleNormalKey(msg)
	}
}

func (a App) handleNormalKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	active := a.getActivePane()

	switch {
	case key.Matches(msg, keys.Quit):
		return a, tea.Quit

	case key.Matches(msg, keys.Up):
		active.MoveUp()
		cmds = append(cmds, a.loadPreviewCmd())

	case key.Matches(msg, keys.Down):
		active.MoveDown()
		cmds = append(cmds, a.loadPreviewCmd())

	case key.Matches(msg, keys.PageUp):
		active.PageUp()
		cmds = append(cmds, a.loadPreviewCmd())

	case key.Matches(msg, keys.PageDown):
		active.PageDown()
		cmds = append(cmds, a.loadPreviewCmd())

	case key.Matches(msg, keys.Home):
		active.GoTop()
		cmds = append(cmds, a.loadPreviewCmd())

	case key.Matches(msg, keys.End):
		active.GoBottom()
		cmds = append(cmds, a.loadPreviewCmd())

	case key.Matches(msg, keys.Enter):
		if entry, ok := active.SelectedEntry(); ok {
			if entry.IsDir {
				var newDir string
				if entry.Name == ".." {
					newDir = filepath.Dir(active.Dir())
				} else {
					newDir = filepath.Join(active.Dir(), entry.Name)
				}
				cmds = append(cmds, pane.LoadDir(active.ID(), newDir))
			} else {
				path := active.SelectedPath()
				if err := fileops.OpenFile(path); err != nil {
					a.statusBar.SetMessage(fmt.Sprintf("Open failed: %v", err), true)
				}
			}
		}

	case key.Matches(msg, keys.Back):
		newDir := filepath.Dir(active.Dir())
		if newDir != active.Dir() {
			cmds = append(cmds, pane.LoadDir(active.ID(), newDir))
		}

	case key.Matches(msg, keys.Tab):
		if a.activePane == 0 {
			a.activePane = 1
		} else {
			a.activePane = 0
		}
		cmds = append(cmds, a.loadPreviewCmd())

	case key.Matches(msg, keys.Toggle):
		a.preview.Toggle()
		a.updateLayout()
		cmds = append(cmds, a.loadPreviewCmd())

	case key.Matches(msg, keys.Copy):
		if active.MarkedCount() > 0 {
			a.clipboard = active.MarkedPaths()
			a.clipAction = clipCopy
			a.statusBar.SetMessage(fmt.Sprintf("Copied %d files to clipboard", len(a.clipboard)), false)
		} else if path := active.SelectedPath(); path != "" {
			entry, _ := active.SelectedEntry()
			if entry.Name != ".." {
				a.clipboard = []string{path}
				a.clipAction = clipCopy
				a.statusBar.SetMessage(fmt.Sprintf("Copied to clipboard: %s", filepath.Base(path)), false)
			}
		}

	case key.Matches(msg, keys.Move):
		if active.MarkedCount() > 0 {
			a.clipboard = active.MarkedPaths()
			a.clipAction = clipMove
			a.statusBar.SetMessage(fmt.Sprintf("Cut %d files to clipboard", len(a.clipboard)), false)
		} else if path := active.SelectedPath(); path != "" {
			entry, _ := active.SelectedEntry()
			if entry.Name != ".." {
				a.clipboard = []string{path}
				a.clipAction = clipMove
				a.statusBar.SetMessage(fmt.Sprintf("Cut to clipboard: %s", filepath.Base(path)), false)
			}
		}

	case key.Matches(msg, keys.Paste):
		if len(a.clipboard) > 0 && a.clipAction != clipNone {
			other := a.getOtherPane()
			dst := other.Dir()
			srcs := a.clipboard
			action := a.clipAction
			a.clipboard = nil
			a.clipAction = clipNone
			active.ClearMarks()

			if action == clipCopy {
				cmds = append(cmds, func() tea.Msg {
					for _, src := range srcs {
						if err := fileops.Copy(src, dst); err != nil {
							return FileOpResultMsg{Err: err, Op: "Copy"}
						}
					}
					return FileOpResultMsg{Op: "Copy"}
				})
			} else {
				cmds = append(cmds, func() tea.Msg {
					for _, src := range srcs {
						if err := fileops.Move(src, dst); err != nil {
							return FileOpResultMsg{Err: err, Op: "Move"}
						}
					}
					return FileOpResultMsg{Op: "Move"}
				})
			}
		}

	case key.Matches(msg, keys.Delete):
		if active.MarkedCount() > 0 {
			a.pendingDeletePaths = active.MarkedPaths()
			a.mode = modeDialog
			a.dialog = dialog.NewConfirm(
				"Delete",
				fmt.Sprintf("Delete %d files?", len(a.pendingDeletePaths)),
				"delete-multi:",
				a.width,
			)
		} else if entry, ok := active.SelectedEntry(); ok && entry.Name != ".." {
			path := active.SelectedPath()
			a.mode = modeDialog
			a.dialog = dialog.NewConfirm(
				"Delete",
				fmt.Sprintf("Delete %q?", entry.Name),
				"delete:"+path,
				a.width,
			)
		}

	case key.Matches(msg, keys.Rename):
		if entry, ok := active.SelectedEntry(); ok && entry.Name != ".." {
			path := active.SelectedPath()
			a.mode = modeDialog
			a.dialog = dialog.NewInput(
				"Rename",
				"rename:"+path,
				"new name",
				entry.Name,
				a.width,
			)
		}

	case key.Matches(msg, keys.Mkdir):
		a.mode = modeDialog
		a.dialog = dialog.NewInput(
			"New Directory",
			"mkdir:"+active.Dir(),
			"directory name",
			"",
			a.width,
		)

	case key.Matches(msg, keys.Search):
		a.mode = modeSearch
		a.updateLayout()
		active.StartSearch()
		a.searchInput.SetValue("")
		a.searchInput.Focus()
		return a, textinput.Blink

	case key.Matches(msg, keys.Bookmark):
		a.mode = modeBookmark
		a.bookmarks = bookmark.NewModel(a.width, a.height)

	case key.Matches(msg, keys.BookAdd):
		dir := active.Dir()
		name := filepath.Base(dir)
		if err := bookmark.Add(name, dir); err != nil {
			a.statusBar.SetMessage(fmt.Sprintf("Bookmark error: %v", err), true)
		} else {
			a.statusBar.SetMessage(fmt.Sprintf("Bookmarked: %s", dir), false)
		}

	case key.Matches(msg, keys.MarkToggle):
		active.ToggleMark()
		cmds = append(cmds, a.loadPreviewCmd())

	case key.Matches(msg, keys.SelectAll):
		active.ToggleAllMarks()

	case key.Matches(msg, keys.ShiftUp):
		active.MoveUpWithMark()
		cmds = append(cmds, a.loadPreviewCmd())

	case key.Matches(msg, keys.ShiftDown):
		active.MoveDownWithMark()
		cmds = append(cmds, a.loadPreviewCmd())

	case key.Matches(msg, keys.GotoDir):
		a.mode = modeDialog
		a.dialog = dialog.NewInput(
			"Go to Directory",
			"goto:"+fmt.Sprintf("%d", active.ID()),
			"path",
			active.Dir(),
			a.width,
		)
		return a, textinput.Blink

	case key.Matches(msg, keys.Explorer):
		dir := active.Dir()
		if err := fileops.OpenFile(dir); err != nil {
			a.statusBar.SetMessage(fmt.Sprintf("Open failed: %v", err), true)
		}

	case key.Matches(msg, keys.Help):
		a.mode = modeHelp
	}

	return a, tea.Batch(cmds...)
}

func (a App) handleDialogKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if a.dialog == nil {
		a.mode = modeNormal
		return a, nil
	}
	newDialog, cmd := a.dialog.Update(msg)
	a.dialog = newDialog
	return a, cmd
}

func (a App) handleSearchKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	active := a.getActivePane()

	switch msg.String() {
	case "enter":
		active.EndSearch(true)
		a.mode = modeNormal
		a.updateLayout()
		return a, a.loadPreviewCmd()
	case "esc":
		active.EndSearch(false)
		a.mode = modeNormal
		a.updateLayout()
		return a, nil
	default:
		var cmd tea.Cmd
		a.searchInput, cmd = a.searchInput.Update(msg)
		active.UpdateSearch(a.searchInput.Value())
		return a, cmd
	}
}

func (a App) handleBookmarkKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.bookmarks, cmd = a.bookmarks.Update(msg)
	return a, cmd
}

func (a *App) handleDialogResult(msg dialog.ResultMsg) tea.Cmd {
	parts := strings.SplitN(msg.Action, ":", 2)
	if len(parts) != 2 {
		return nil
	}
	action := parts[0]
	target := parts[1]

	switch action {
	case "delete":
		return func() tea.Msg {
			err := fileops.Delete(target)
			return FileOpResultMsg{Err: err, Op: "Delete"}
		}
	case "delete-multi":
		paths := a.pendingDeletePaths
		a.pendingDeletePaths = nil
		active := a.getActivePane()
		active.ClearMarks()
		return func() tea.Msg {
			for _, p := range paths {
				if err := fileops.Delete(p); err != nil {
					return FileOpResultMsg{Err: err, Op: "Delete"}
				}
			}
			return FileOpResultMsg{Op: "Delete"}
		}
	case "rename":
		return func() tea.Msg {
			err := fileops.Rename(target, msg.Text)
			return FileOpResultMsg{Err: err, Op: "Rename"}
		}
	case "mkdir":
		return func() tea.Msg {
			err := fileops.Mkdir(target, msg.Text)
			return FileOpResultMsg{Err: err, Op: "Mkdir"}
		}
	case "goto":
		paneID := 0
		if target == "1" {
			paneID = 1
		}
		dir := msg.Text
		if !filepath.IsAbs(dir) {
			p := a.getActivePane()
			dir = filepath.Join(p.Dir(), dir)
		}
		return pane.LoadDir(paneID, dir)
	}
	return nil
}

func (a *App) getActivePane() *pane.Model {
	if a.activePane == 0 {
		return &a.leftPane
	}
	return &a.rightPane
}

func (a *App) getOtherPane() *pane.Model {
	if a.activePane == 0 {
		return &a.rightPane
	}
	return &a.leftPane
}

func (a *App) updateLayout() {
	if a.width == 0 || a.height == 0 {
		return
	}

	statusH := 1
	if a.mode == modeSearch {
		statusH = 2 // status bar + search bar
	}
	contentH := a.height - statusH

	if a.preview.Visible() {
		leftW := a.width * 35 / 100
		rightW := a.width * 35 / 100
		previewW := a.width - leftW - rightW
		a.leftPane.SetSize(leftW, contentH)
		a.rightPane.SetSize(rightW, contentH)
		a.preview.SetSize(previewW, contentH)
	} else {
		halfW := a.width / 2
		a.leftPane.SetSize(halfW, contentH)
		a.rightPane.SetSize(a.width-halfW, contentH)
	}
	a.statusBar.SetWidth(a.width)
}

func (a App) loadPreviewCmd() tea.Cmd {
	if !a.preview.Visible() {
		return nil
	}
	active := a.getActivePane()
	path := active.SelectedPath()
	if path == "" {
		return nil
	}
	return func() tea.Msg {
		loader := preview.LoadFile(path)
		content, isBinary, err := loader()
		if err != nil {
			return preview.LoadMsg{Content: fmt.Sprintf("Error: %v", err), Path: path}
		}
		return preview.LoadMsg{Content: content, IsBinary: isBinary, Path: path}
	}
}

func (a App) View() string {
	if !a.ready {
		return "Loading..."
	}

	// Render panes
	leftView := a.leftPane.View(a.activePane == 0)
	rightView := a.rightPane.View(a.activePane == 1)

	var contentView string
	if a.preview.Visible() {
		previewView := a.preview.View()
		contentView = lipgloss.JoinHorizontal(lipgloss.Top, leftView, rightView, previewView)
	} else {
		contentView = lipgloss.JoinHorizontal(lipgloss.Top, leftView, rightView)
	}

	// Search bar
	var searchBar string
	if a.mode == modeSearch {
		searchStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#c0caf5")).
			Background(lipgloss.Color("#16161e")).
			Width(a.width)
		searchBar = searchStyle.Render("/" + a.searchInput.View())
	}

	// Status bar
	active := a.getActivePane()
	statusView := a.statusBar.View(active, a.mode == modeSearch, a.searchInput.Value())

	var mainView string
	if a.mode == modeSearch {
		mainView = lipgloss.JoinVertical(lipgloss.Left, contentView, searchBar, statusView)
	} else {
		mainView = lipgloss.JoinVertical(lipgloss.Left, contentView, statusView)
	}

	// Overlay dialogs
	switch a.mode {
	case modeDialog:
		if a.dialog != nil {
			return a.overlayCenter(mainView, a.dialog.View())
		}
	case modeBookmark:
		return a.overlayCenter(mainView, a.bookmarks.View())
	case modeHelp:
		return a.overlayCenter(mainView, a.helpView())
	}

	return mainView
}

func (a App) overlayCenter(bg, fg string) string {
	return lipgloss.Place(
		a.width, a.height,
		lipgloss.Center, lipgloss.Center,
		fg,
		lipgloss.WithWhitespaceBackground(lipgloss.Color("#1a1b26")),
	)
}

func (a App) helpView() string {
	w := a.width / 2
	if w < 50 {
		w = 50
	}
	if w > a.width-4 {
		w = a.width - 4
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#bb9af7")).
		Bold(true)

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7dcfff")).
		Bold(true).
		Width(16)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a9b1d6"))

	helpItems := []struct{ key, desc string }{
		{"↑/↓", "Move cursor"},
		{"Enter", "Open dir/file"},
		{"Backspace", "Parent directory"},
		{"Tab", "Switch pane"},
		{"PgUp/PgDn", "Page scroll"},
		{"Home/End", "First/Last"},
		{"Space", "Mark/Unmark"},
		{"Shift+↑/↓", "Mark and move"},
		{"Ctrl+A", "Select all"},
		{"F5/c", "Copy to clipboard"},
		{"F6/m", "Move to clipboard"},
		{"p", "Paste to other pane"},
		{"F7/n", "New directory"},
		{"F8/d", "Delete"},
		{"r", "Rename"},
		{"/", "Search"},
		{"t", "Toggle preview"},
		{"g", "Go to directory"},
		{"e", "Open in explorer"},
		{"b", "Bookmarks"},
		{"B", "Add bookmark"},
		{"?", "This help"},
		{"q/Ctrl+C", "Quit"},
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render("CFiler Help"))
	b.WriteString("\n\n")

	for _, item := range helpItems {
		b.WriteString(keyStyle.Render(item.key))
		b.WriteString(descStyle.Render(item.desc))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565f89")).
		Render("Press Esc or ? to close"))

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#bb9af7")).
		Padding(1, 2).
		Width(w)

	return boxStyle.Render(b.String())
}
