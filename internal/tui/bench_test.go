package tui

import (
	"strings"
	"testing"

	"github.com/arthurlch/cub/internal/editor"
	"github.com/arthurlch/cub/internal/filetree"
	tea "github.com/charmbracelet/bubbletea"
)

func benchModel(b *testing.B, lines int) Model {
	b.Helper()
	var sb strings.Builder
	for i := 0; i < lines; i++ {
		sb.WriteString("func handleThing(x int, name string) error { return fmt.Errorf(\"bad %d\", x) }\n")
	}
	e := editor.New()
	e.SetText(sb.String())
	m := New(e, &filetree.Node{Name: "root"}, 26)
	tm, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 50})
	return tm.(Model)
}

func BenchmarkViewLargeFile(b *testing.B) {
	m := benchModel(b, 20000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.View()
	}
}

func BenchmarkViewWithSidebar(b *testing.B) {
	m := benchModel(b, 20000)
	m.sbVisible, m.sbFocused, m.sbReveal = true, true, 1
	m.applySize()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.View()
	}
}

// The visible viewport is a fixed size, so a full render must cost the same on
// a huge file as on a small one — rendering only touches on-screen lines.
func TestViewCostIsViewportBound(t *testing.T) {
	small := New(editor.New(), &filetree.Node{Name: "root"}, 26)
	small.e.SetText(strings.Repeat("x := 1\n", 40))
	tm, _ := small.Update(tea.WindowSizeMsg{Width: 120, Height: 50})
	small = tm.(Model)

	big := New(editor.New(), &filetree.Node{Name: "root"}, 26)
	big.e.SetText(strings.Repeat("x := 1\n", 100000))
	tm, _ = big.Update(tea.WindowSizeMsg{Width: 120, Height: 50})
	big = tm.(Model)

	smallAllocs := testing.AllocsPerRun(20, func() { _ = small.View() })
	bigAllocs := testing.AllocsPerRun(20, func() { _ = big.View() })

	if bigAllocs > smallAllocs*2 {
		t.Fatalf("View cost scales with file size: small=%.0f big=%.0f allocs", smallAllocs, bigAllocs)
	}
}
