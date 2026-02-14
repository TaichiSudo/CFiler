package dialog

import tea "github.com/charmbracelet/bubbletea"

type Dialog interface {
	Update(msg tea.Msg) (Dialog, tea.Cmd)
	View() string
}

type ResultMsg struct {
	Confirmed bool
	Text      string
	Action    string
}
