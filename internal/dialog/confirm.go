package dialog

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ConfirmDialog struct {
	title   string
	message string
	action  string
	width   int
}

func NewConfirm(title, message, action string, width int) *ConfirmDialog {
	return &ConfirmDialog{
		title:   title,
		message: message,
		action:  action,
		width:   width,
	}
}

func (d *ConfirmDialog) Update(msg tea.Msg) (Dialog, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y", "enter":
			return d, func() tea.Msg {
				return ResultMsg{Confirmed: true, Action: d.action}
			}
		case "n", "N", "esc":
			return d, func() tea.Msg {
				return ResultMsg{Confirmed: false, Action: d.action}
			}
		}
	}
	return d, nil
}

func (d *ConfirmDialog) View() string {
	dialogW := d.width / 2
	if dialogW < 40 {
		dialogW = 40
	}
	if dialogW > d.width-4 {
		dialogW = d.width - 4
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#bb9af7")).
		Bold(true)

	msgStyle := lipgloss.NewStyle().
		Width(dialogW - 4)

	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565f89"))

	content := fmt.Sprintf("%s\n\n%s\n\n%s",
		titleStyle.Render(d.title),
		msgStyle.Render(d.message),
		promptStyle.Render("[Y]es / [N]o"),
	)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#bb9af7")).
		Padding(1, 2).
		Width(dialogW)

	return boxStyle.Render(content)
}
