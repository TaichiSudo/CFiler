package pane

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View(active bool) string {
	if m.width <= 2 || m.height <= 2 {
		return ""
	}

	innerWidth := m.width - 2 // account for border
	if innerWidth < 10 {
		innerWidth = 10
	}

	var lines []string

	// Header: directory path
	header := padOrTruncate(m.dir, innerWidth)
	headerSt := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7aa2f7")).
		Bold(true)
	lines = append(lines, headerSt.Render(header))

	entries := m.Entries()
	vis := m.visibleLines()

	// Column widths
	const sizeCol = 8
	const timeCol = 13 // space + "Jan 02 15:04" (12 chars)
	showDetails := innerWidth >= 35
	nameWidth := innerWidth
	if showDetails {
		nameWidth = innerWidth - sizeCol - timeCol
		if nameWidth < 8 {
			nameWidth = 8
		}
	}

	for i := 0; i < vis; i++ {
		idx := m.offset + i
		if idx >= len(entries) {
			lines = append(lines, strings.Repeat(" ", innerWidth))
			continue
		}

		entry := entries[idx]
		isCursor := idx == m.cursor
		isMarked := m.IsMarked(entry.Name)

		// Build name part
		name := entry.Name
		if entry.IsLink {
			name += " @"
		}
		namePadded := padOrTruncate(name, nameWidth)

		// Build detail parts separately
		var sizeStr, timeStr string
		hasDetail := showDetails && entry.Name != ".."
		if hasDetail {
			sizeStr = padLeftStr(formatSize(entry.Size, entry.IsDir), sizeCol)
			timeStr = " " + padOrTruncate(formatTime(entry.ModTime), timeCol-1)
		}

		if isCursor && isMarked {
			// Cursor + marked: visual background + yellow text
			st := lipgloss.NewStyle().
				Background(lipgloss.Color("#283457")).
				Foreground(lipgloss.Color("#e0af68"))
			plain := namePadded
			if hasDetail {
				plain += sizeStr + timeStr
			} else {
				plain = padOrTruncate(plain, innerWidth)
			}
			lines = append(lines, st.Render(plain))
		} else if isCursor {
			// Cursor: uniform background, single color
			curSt := lipgloss.NewStyle().
				Background(lipgloss.Color("#283457")).
				Foreground(lipgloss.Color("#c0caf5"))
			plain := namePadded
			if hasDetail {
				plain += sizeStr + timeStr
			} else {
				plain = padOrTruncate(plain, innerWidth)
			}
			lines = append(lines, curSt.Render(plain))
		} else if isMarked {
			// Marked: yellow text + bold
			markSt := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#e0af68")).
				Bold(true)
			if hasDetail {
				line := markSt.Render(namePadded + sizeStr + timeStr)
				lines = append(lines, line)
			} else {
				lines = append(lines, markSt.Render(padOrTruncate(namePadded, innerWidth)))
			}
		} else {
			// Normal: name colored, details gray
			var nameSt lipgloss.Style
			if entry.IsDir {
				nameSt = lipgloss.NewStyle().Foreground(lipgloss.Color("#7aa2f7")).Bold(true)
			} else if entry.IsLink {
				nameSt = lipgloss.NewStyle().Foreground(lipgloss.Color("#bb9af7"))
			} else {
				nameSt = lipgloss.NewStyle().Foreground(lipgloss.Color("#c0caf5"))
			}

			if hasDetail {
				detailSt := lipgloss.NewStyle().Foreground(lipgloss.Color("#565f89"))
				line := nameSt.Render(namePadded) + detailSt.Render(sizeStr+timeStr)
				lines = append(lines, line)
			} else {
				lines = append(lines, nameSt.Render(padOrTruncate(namePadded, innerWidth)))
			}
		}
	}

	content := strings.Join(lines, "\n")

	var borderColor lipgloss.Color
	if active {
		borderColor = lipgloss.Color("#7aa2f7")
	} else {
		borderColor = lipgloss.Color("#545c7e")
	}

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(innerWidth).
		Height(m.height - 2)

	return borderStyle.Render(content)
}

// padOrTruncate pads with spaces or truncates to exactly maxLen runes
func padOrTruncate(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) > maxLen {
		if maxLen > 1 {
			return string(r[:maxLen-1]) + "~"
		}
		return string(r[:maxLen])
	}
	if len(r) < maxLen {
		return s + strings.Repeat(" ", maxLen-len(r))
	}
	return s
}

// padLeftStr right-aligns s within maxLen (rune-based)
func padLeftStr(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) > maxLen {
		return string(r[:maxLen])
	}
	if len(r) < maxLen {
		return strings.Repeat(" ", maxLen-len(r)) + s
	}
	return s
}

func formatSize(size int64, isDir bool) string {
	if isDir {
		return "<DIR>"
	}
	switch {
	case size < 1024:
		return fmt.Sprintf("%dB", size)
	case size < 1024*1024:
		return fmt.Sprintf("%.1fK", float64(size)/1024)
	case size < 1024*1024*1024:
		return fmt.Sprintf("%.1fM", float64(size)/(1024*1024))
	default:
		return fmt.Sprintf("%.1fG", float64(size)/(1024*1024*1024))
	}
}

func formatTime(t time.Time) string {
	now := time.Now()
	if t.Year() == now.Year() {
		return t.Format("Jan 02 15:04")
	}
	return t.Format("Jan 02  2006")
}
