package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func parseHex(c lipgloss.Color) (r, g, b int) {
	s := strings.TrimPrefix(string(c), "#")
	if len(s) != 6 {
		return 255, 255, 255
	}
	ri, _ := strconv.ParseInt(s[0:2], 16, 0)
	gi, _ := strconv.ParseInt(s[2:4], 16, 0)
	bi, _ := strconv.ParseInt(s[4:6], 16, 0)
	return int(ri), int(gi), int(bi)
}

func lerpColor(a, b lipgloss.Color, t float64) lipgloss.Color {
	ar, ag, ab := parseHex(a)
	br, bg, bb := parseHex(b)
	r := int(float64(ar) + (float64(br)-float64(ar))*t)
	g := int(float64(ag) + (float64(bg)-float64(ag))*t)
	bl := int(float64(ab) + (float64(bb)-float64(ab))*t)
	return lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r, g, bl))
}

func (m Model) titleBar(width int) string {
	if width < 1 {
		return ""
	}
	t := m.theme
	barBg := m.st.sbBg
	active := lipgloss.NewStyle().Foreground(t.Fg).Background(m.st.selBg).Bold(true)
	marker := lipgloss.NewStyle().Foreground(t.Accent).Background(m.st.selBg).Bold(true)
	inactive := lipgloss.NewStyle().Foreground(m.st.dim).Background(barBg)

	var b strings.Builder
	used := 0
	for _, bx := range m.tabLayout(width) {
		label := m.tabLabel(m.docs[bx.doc])
		if bx.doc == m.active {
			b.WriteString(marker.Render("▎"))
			b.WriteString(active.Render(label[1:]))
		} else {
			b.WriteString(inactive.Render(label))
		}
		used += bx.w
	}
	if used < width {
		b.WriteString(lipgloss.NewStyle().Background(barBg).Render(strings.Repeat(" ", width-used)))
	}
	return b.String()
}

func (m Model) scrollbar(height, offRow, lineCount int) string {
	t := m.theme
	track := lipgloss.NewStyle().Foreground(t.Gutter).Background(t.Bg)
	thumb := lipgloss.NewStyle().Foreground(t.Accent).Background(t.Bg)

	thumbStart, thumbSize := 0, height
	if lineCount > height && height > 0 {
		thumbSize = height * height / lineCount
		if thumbSize < 1 {
			thumbSize = 1
		}
		maxOff := lineCount - height
		if maxOff > 0 {
			thumbStart = (height - thumbSize) * offRow / maxOff
		}
	}

	rows := make([]string, height)
	for i := 0; i < height; i++ {
		if lineCount > height && i >= thumbStart && i < thumbStart+thumbSize {
			rows[i] = thumb.Render("┃")
		} else {
			rows[i] = track.Render("│")
		}
	}
	return strings.Join(rows, "\n")
}
