package bookmark

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SelectMsg struct {
	Path string
}

type CloseMsg struct{}

type Model struct {
	entries []Entry
	cursor  int
	width   int
	height  int
}

func NewModel(width, height int) Model {
	entries, _ := Load()
	return Model{
		entries: entries,
		width:   width,
		height:  height,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.entries)-1 {
				m.cursor++
			}
		case "enter":
			if m.cursor < len(m.entries) {
				path := m.entries[m.cursor].Path
				return m, func() tea.Msg { return SelectMsg{Path: path} }
			}
		case "d", "delete":
			if m.cursor < len(m.entries) {
				path := m.entries[m.cursor].Path
				_ = Remove(path)
				m.entries, _ = Load()
				if m.cursor >= len(m.entries) && m.cursor > 0 {
					m.cursor--
				}
			}
		case "esc", "b", "q":
			return m, func() tea.Msg { return CloseMsg{} }
		}
	}
	return m, nil
}

func (m Model) View() string {
	dialogW := m.width / 2
	if dialogW < 50 {
		dialogW = 50
	}
	if dialogW > m.width-4 {
		dialogW = m.width - 4
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#bb9af7")).
		Bold(true)

	var b strings.Builder
	b.WriteString(titleStyle.Render("Bookmarks"))
	b.WriteString("\n\n")

	if len(m.entries) == 0 {
		b.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("#565f89")).
			Render("No bookmarks. Press B to add one."))
	} else {
		for i, entry := range m.entries {
			cursor := "  "
			if i == m.cursor {
				cursor = "▸ "
			}

			nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#7aa2f7")).Bold(true)
			pathStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#565f89"))

			line := fmt.Sprintf("%s%s  %s",
				cursor,
				nameStyle.Render(entry.Name),
				pathStyle.Render(entry.Path),
			)
			b.WriteString(line)
			if i < len(m.entries)-1 {
				b.WriteString("\n")
			}
		}
	}

	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565f89")).
		Render("Enter: open  d: delete  Esc: close"))

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#bb9af7")).
		Padding(1, 2).
		Width(dialogW)

	return boxStyle.Render(b.String())
}
