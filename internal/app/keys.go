package app

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Enter    key.Binding
	Back     key.Binding
	Tab      key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	Home     key.Binding
	End      key.Binding
	Copy     key.Binding
	Move     key.Binding
	Paste    key.Binding
	Mkdir    key.Binding
	Delete   key.Binding
	Rename   key.Binding
	Search   key.Binding
	Escape   key.Binding
	Bookmark key.Binding
	BookAdd  key.Binding
	Toggle    key.Binding
	Help      key.Binding
	Quit      key.Binding
	MarkToggle key.Binding
	SelectAll  key.Binding
	ShiftUp    key.Binding
	ShiftDown  key.Binding
	GotoDir    key.Binding
	Explorer   key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("Enter", "open"),
	),
	Back: key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("BS", "parent dir"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("Tab", "switch pane"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("PgUp", "page up"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("PgDn", "page down"),
	),
	Home: key.NewBinding(
		key.WithKeys("home"),
		key.WithHelp("Home", "first"),
	),
	End: key.NewBinding(
		key.WithKeys("end"),
		key.WithHelp("End", "last"),
	),
	Copy: key.NewBinding(
		key.WithKeys("f5", "c"),
		key.WithHelp("F5/c", "copy"),
	),
	Move: key.NewBinding(
		key.WithKeys("f6", "m"),
		key.WithHelp("F6/m", "move"),
	),
	Paste: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "paste"),
	),
	Mkdir: key.NewBinding(
		key.WithKeys("f7", "n"),
		key.WithHelp("F7/n", "mkdir"),
	),
	Delete: key.NewBinding(
		key.WithKeys("f8", "d"),
		key.WithHelp("F8/d", "delete"),
	),
	Rename: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rename"),
	),
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "search"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("Esc", "cancel"),
	),
	Bookmark: key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "bookmarks"),
	),
	BookAdd: key.NewBinding(
		key.WithKeys("B"),
		key.WithHelp("B", "add bookmark"),
	),
	Toggle: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "toggle preview"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	MarkToggle: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("Space", "mark/unmark"),
	),
	SelectAll: key.NewBinding(
		key.WithKeys("ctrl+a"),
		key.WithHelp("Ctrl+A", "select all"),
	),
	ShiftUp: key.NewBinding(
		key.WithKeys("shift+up"),
		key.WithHelp("Shift+↑", "mark and move up"),
	),
	ShiftDown: key.NewBinding(
		key.WithKeys("shift+down"),
		key.WithHelp("Shift+↓", "mark and move down"),
	),
	GotoDir: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "go to dir"),
	),
	Explorer: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "explorer"),
	),
}
