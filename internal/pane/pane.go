package pane

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	id         int
	dir        string
	entries    []FileEntry
	cursor     int
	offset     int
	width      int
	height     int
	search     string
	searching  bool
	filtered   []FileEntry
	err        error
	marked     map[string]bool
}

func New(id int, dir string) Model {
	return Model{
		id:  id,
		dir: dir,
	}
}

func (m Model) ID() int        { return m.id }
func (m Model) Dir() string    { return m.dir }
func (m Model) Cursor() int    { return m.cursor }
func (m Model) Width() int     { return m.width }
func (m Model) Height() int    { return m.height }
func (m Model) Search() string { return m.search }
func (m Model) Searching() bool { return m.searching }
func (m Model) Err() error     { return m.err }

func (m Model) Entries() []FileEntry {
	if m.searching && m.search != "" {
		return m.filtered
	}
	return m.entries
}

func (m Model) SelectedEntry() (FileEntry, bool) {
	entries := m.Entries()
	if m.cursor >= 0 && m.cursor < len(entries) {
		return entries[m.cursor], true
	}
	return FileEntry{}, false
}

func (m Model) SelectedPath() string {
	if entry, ok := m.SelectedEntry(); ok {
		return filepath.Join(m.dir, entry.Name)
	}
	return ""
}

func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m *Model) SetDir(dir string) {
	m.dir = dir
}

func (m *Model) SetEntries(entries []FileEntry) {
	m.entries = entries
	m.err = nil
	m.marked = nil
	m.clampCursor()
}

func (m *Model) SetError(err error) {
	m.err = err
}

func (m *Model) MoveUp() {
	if m.cursor > 0 {
		m.cursor--
	}
	m.adjustOffset()
}

func (m *Model) MoveDown() {
	entries := m.Entries()
	if m.cursor < len(entries)-1 {
		m.cursor++
	}
	m.adjustOffset()
}

func (m *Model) PageUp() {
	m.cursor -= m.visibleLines()
	if m.cursor < 0 {
		m.cursor = 0
	}
	m.adjustOffset()
}

func (m *Model) PageDown() {
	entries := m.Entries()
	m.cursor += m.visibleLines()
	if m.cursor >= len(entries) {
		m.cursor = len(entries) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	m.adjustOffset()
}

func (m *Model) GoTop() {
	m.cursor = 0
	m.adjustOffset()
}

func (m *Model) GoBottom() {
	entries := m.Entries()
	if len(entries) > 0 {
		m.cursor = len(entries) - 1
	}
	m.adjustOffset()
}

func (m *Model) StartSearch() {
	m.searching = true
	m.search = ""
	m.filtered = nil
}

func (m *Model) UpdateSearch(s string) {
	m.search = s
	m.filterEntries()
	m.cursor = 0
	m.offset = 0
}

func (m *Model) EndSearch(confirm bool) {
	if confirm && m.searching && m.search != "" {
		if entry, ok := m.SelectedEntry(); ok {
			for i, e := range m.entries {
				if e.Name == entry.Name {
					m.cursor = i
					break
				}
			}
		}
	}
	m.searching = false
	m.search = ""
	m.filtered = nil
	m.clampCursor()
}

func (m *Model) filterEntries() {
	if m.search == "" {
		m.filtered = nil
		return
	}
	lower := strings.ToLower(m.search)
	m.filtered = nil
	for _, e := range m.entries {
		if strings.Contains(strings.ToLower(e.Name), lower) {
			m.filtered = append(m.filtered, e)
		}
	}
}

func (m *Model) clampCursor() {
	entries := m.Entries()
	if m.cursor >= len(entries) {
		m.cursor = len(entries) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	m.adjustOffset()
}

func (m *Model) adjustOffset() {
	vis := m.visibleLines()
	if vis <= 0 {
		return
	}
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+vis {
		m.offset = m.cursor - vis + 1
	}
}

func (m Model) visibleLines() int {
	// height minus border (2) and header line (1)
	h := m.height - 3
	if h < 1 {
		h = 1
	}
	return h
}

// --- Mark (multi-select) methods ---

func (m *Model) ToggleMark() {
	if entry, ok := m.SelectedEntry(); ok && entry.Name != ".." {
		if m.marked == nil {
			m.marked = make(map[string]bool)
		}
		if m.marked[entry.Name] {
			delete(m.marked, entry.Name)
		} else {
			m.marked[entry.Name] = true
		}
	}
	m.MoveDown()
}

func (m *Model) SetMark(name string, on bool) {
	if name == ".." {
		return
	}
	if on {
		if m.marked == nil {
			m.marked = make(map[string]bool)
		}
		m.marked[name] = true
	} else {
		delete(m.marked, name)
	}
}

func (m Model) IsMarked(name string) bool {
	return m.marked[name]
}

func (m Model) MarkedNames() []string {
	var names []string
	for name := range m.marked {
		names = append(names, name)
	}
	return names
}

func (m Model) MarkedPaths() []string {
	var paths []string
	for name := range m.marked {
		paths = append(paths, filepath.Join(m.dir, name))
	}
	return paths
}

func (m Model) MarkedCount() int {
	return len(m.marked)
}

func (m *Model) ClearMarks() {
	m.marked = nil
}

func (m *Model) ToggleAllMarks() {
	entries := m.Entries()
	// If any are marked, clear all; otherwise mark all (except "..")
	if m.MarkedCount() > 0 {
		m.marked = nil
		return
	}
	m.marked = make(map[string]bool)
	for _, e := range entries {
		if e.Name != ".." {
			m.marked[e.Name] = true
		}
	}
}

func (m *Model) MoveUpWithMark() {
	if m.cursor > 0 {
		m.cursor--
	}
	m.adjustOffset()
	if entry, ok := m.SelectedEntry(); ok && entry.Name != ".." {
		m.SetMark(entry.Name, true)
	}
}

func (m *Model) MoveDownWithMark() {
	entries := m.Entries()
	if m.cursor < len(entries)-1 {
		m.cursor++
	}
	m.adjustOffset()
	if entry, ok := m.SelectedEntry(); ok && entry.Name != ".." {
		m.SetMark(entry.Name, true)
	}
}

type DirLoadedMsg struct {
	Entries []FileEntry
	Path    string
	PaneID  int
}

type DirLoadErrorMsg struct {
	Err    error
	PaneID int
}

func LoadDir(id int, dir string) tea.Cmd {
	return func() tea.Msg {
		dirEntries, err := os.ReadDir(dir)
		if err != nil {
			return DirLoadErrorMsg{Err: err, PaneID: id}
		}

		var entries []FileEntry
		// Add parent directory entry
		absDir, _ := filepath.Abs(dir)
		parent := filepath.Dir(absDir)
		if parent != absDir {
			entries = append(entries, FileEntry{
				Name:  "..",
				IsDir: true,
			})
		}

		var dirs, files []FileEntry
		for _, de := range dirEntries {
			info, err := de.Info()
			if err != nil {
				continue
			}
			entry := FileEntry{
				Name:    de.Name(),
				Size:    info.Size(),
				ModTime: info.ModTime(),
				IsDir:   de.IsDir(),
				Mode:    info.Mode(),
				IsLink:  de.Type()&os.ModeSymlink != 0,
			}
			if de.IsDir() {
				dirs = append(dirs, entry)
			} else {
				files = append(files, entry)
			}
		}

		sort.Slice(dirs, func(i, j int) bool {
			return strings.ToLower(dirs[i].Name) < strings.ToLower(dirs[j].Name)
		})
		sort.Slice(files, func(i, j int) bool {
			return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
		})

		entries = append(entries, dirs...)
		entries = append(entries, files...)

		return DirLoadedMsg{Entries: entries, Path: absDir, PaneID: id}
	}
}
