package statusbar

import (
	"fmt"
	"strings"

	"cfiler/internal/pane"

	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	width   int
	message string
	isError bool
}

func New() Model {
	return Model{}
}

func (m *Model) SetWidth(w int) {
	m.width = w
}

func (m *Model) SetMessage(msg string, isError bool) {
	m.message = msg
	m.isError = isError
}

func (m Model) View(activePane *pane.Model, searchMode bool, searchText string) string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#c0caf5")).
		Background(lipgloss.Color("#16161e")).
		Width(m.width)

	var parts []string

	if activePane != nil {
		parts = append(parts, activePane.Dir())

		if entry, ok := activePane.SelectedEntry(); ok && entry.Name != ".." {
			info := fmt.Sprintf(" | %s", entry.Mode.String())
			if !entry.IsDir {
				info += fmt.Sprintf(" | %d bytes", entry.Size)
			}
			parts = append(parts, info)
		}
	}

	if activePane != nil && activePane.MarkedCount() > 0 {
		parts = append(parts, fmt.Sprintf(" | %d selected", activePane.MarkedCount()))
	}

	if searchMode {
		parts = append(parts, fmt.Sprintf(" | Search: %s", searchText))
	}

	if m.message != "" {
		parts = append(parts, " | "+m.message)
	}

	text := strings.Join(parts, "")
	text = truncateStatus(text, m.width)

	if m.isError {
		style = style.Foreground(lipgloss.Color("#db4b4b"))
	}

	return style.Render(text)
}

func truncateStatus(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-1]) + "…"
}
