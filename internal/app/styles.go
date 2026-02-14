package app

import "github.com/charmbracelet/lipgloss"

var (
	activeBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7aa2f7"))

	inactiveBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#545c7e"))

	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7aa2f7")).
		Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#c0caf5")).
		Background(lipgloss.Color("#16161e")).
		Padding(0, 1)

	statusErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#db4b4b")).
		Background(lipgloss.Color("#16161e")).
		Padding(0, 1)

	dirStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7aa2f7")).
		Bold(true)

	fileStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#c0caf5"))

	selectedStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#292e42"))

	cursorStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#283457")).
		Foreground(lipgloss.Color("#c0caf5"))

	sizeStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565f89")).
		Width(8).
		Align(lipgloss.Right)

	timeStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565f89"))

	headerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7aa2f7")).
		Bold(true)

	dialogStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#bb9af7")).
		Padding(1, 2)

	dialogTitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#bb9af7")).
		Bold(true)

	helpKeyStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7dcfff")).
		Bold(true)

	helpDescStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a9b1d6"))

	linkStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#bb9af7"))

	previewBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#545c7e"))
)
