package preview

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

const maxPreviewBytes = 64 * 1024

type Model struct {
	viewport viewport.Model
	content  string
	filePath string
	isBinary bool
	width    int
	height   int
	visible  bool
}

func New() Model {
	return Model{
		viewport: viewport.New(0, 0),
		visible:  false,
	}
}

func (m Model) Visible() bool  { return m.visible }
func (m *Model) Toggle()       { m.visible = !m.visible }
func (m *Model) SetVisible(v bool) { m.visible = v }

func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
	innerW := w - 2
	innerH := h - 2
	if innerW < 1 {
		innerW = 1
	}
	if innerH < 1 {
		innerH = 1
	}
	m.viewport.Width = innerW
	m.viewport.Height = innerH - 1 // ヘッダー行分を引く
}

func (m *Model) SetContent(path, content string, isBinary bool) {
	m.filePath = path
	m.isBinary = isBinary
	if isBinary {
		m.content = "[Binary file]"
	} else {
		m.content = content
	}
	m.viewport.SetContent(m.content)
	m.viewport.GotoTop()
}

func (m *Model) Clear() {
	m.filePath = ""
	m.content = ""
	m.isBinary = false
	m.viewport.SetContent("")
}

func (m *Model) ScrollUp() {
	m.viewport.LineUp(1)
}

func (m *Model) ScrollDown() {
	m.viewport.LineDown(1)
}

func (m Model) View() string {
	if !m.visible || m.width <= 2 || m.height <= 2 {
		return ""
	}

	innerW := m.width - 2
	title := truncatePreview(m.filePath, innerW)
	titleSt := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565f89")).
		Bold(true).
		Width(innerW)

	header := titleSt.Render(title)
	body := m.viewport.View()

	content := header + "\n" + body

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#545c7e")).
		Width(innerW).
		Height(m.height - 2)

	return borderStyle.Render(content)
}

type LoadMsg struct {
	Content  string
	IsBinary bool
	Path     string
}

func LoadFile(path string) func() (string, bool, error) {
	return func() (string, bool, error) {
		info, err := os.Stat(path)
		if err != nil {
			return "", false, err
		}
		if info.IsDir() {
			entries, err := os.ReadDir(path)
			if err != nil {
				return "", false, err
			}
			var b strings.Builder
			b.WriteString(fmt.Sprintf("Directory: %s\n", path))
			b.WriteString(fmt.Sprintf("%d items\n\n", len(entries)))
			for _, e := range entries {
				if e.IsDir() {
					b.WriteString(fmt.Sprintf("  [DIR] %s\n", e.Name()))
				} else {
					info, _ := e.Info()
					if info != nil {
						b.WriteString(fmt.Sprintf("  %s (%d bytes)\n", e.Name(), info.Size()))
					}
				}
			}
			return b.String(), false, nil
		}

		readSize := maxPreviewBytes
		if info.Size() < int64(readSize) {
			readSize = int(info.Size())
		}

		f, err := os.Open(path)
		if err != nil {
			return "", false, err
		}
		defer f.Close()

		buf := make([]byte, readSize)
		n, err := f.Read(buf)
		if err != nil && n == 0 {
			return "", false, err
		}
		buf = buf[:n]

		if isBinaryData(buf) {
			return "", true, nil
		}

		if !utf8.Valid(buf) {
			return "", true, nil
		}

		return string(buf), false, nil
	}
}

func isBinaryData(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	checkLen := 512
	if len(data) < checkLen {
		checkLen = len(data)
	}
	nullCount := 0
	for _, b := range data[:checkLen] {
		if b == 0 {
			nullCount++
		}
	}
	return nullCount > 0
}

func truncatePreview(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return string(runes[:maxLen])
	}
	return string(runes[:maxLen-1]) + "…"
}
