package dialog

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InputDialog struct {
	title     string
	action    string
	textInput textinput.Model
	width     int
}

func NewInput(title, action, placeholder, initial string, width int) *InputDialog {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.SetValue(initial)
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = width/2 - 8

	return &InputDialog{
		title:     title,
		action:    action,
		textInput: ti,
		width:     width,
	}
}

func (d *InputDialog) Update(msg tea.Msg) (Dialog, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return d, func() tea.Msg {
				return ResultMsg{
					Confirmed: true,
					Text:      d.textInput.Value(),
					Action:    d.action,
				}
			}
		case "esc":
			return d, func() tea.Msg {
				return ResultMsg{Confirmed: false, Action: d.action}
			}
		}
	}

	var cmd tea.Cmd
	d.textInput, cmd = d.textInput.Update(msg)
	return d, cmd
}

func (d *InputDialog) View() string {
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

	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565f89"))

	content := fmt.Sprintf("%s\n\n%s\n\n%s",
		titleStyle.Render(d.title),
		d.textInput.View(),
		promptStyle.Render("Enter to confirm / Esc to cancel"),
	)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#bb9af7")).
		Padding(1, 2).
		Width(dialogW)

	return boxStyle.Render(content)
}
