package editor_test

import (
	"strings"
	"testing"

	"github.com/arthurlch/cub/internal/editor"
)

func bigEditor(lines int) *editor.Editor {
	var sb strings.Builder
	for i := 0; i < lines; i++ {
		sb.WriteString("func handleThing(x int, name string) error { return fmt.Errorf(\"bad %d\", x) }\n")
	}
	e := editor.New()
	e.SetText(sb.String())
	return e
}

func BenchmarkInsertRune20k(b *testing.B) {
	e := bigEditor(20000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.InsertRune('x')
	}
}

func BenchmarkBackspace20k(b *testing.B) {
	e := bigEditor(20000)
	for i := 0; i < 1000; i++ {
		e.InsertRune('x')
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Backspace()
	}
}

func BenchmarkUndoRedo20k(b *testing.B) {
	e := bigEditor(20000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.InsertRune('x')
		e.Undo()
		e.Redo()
		e.Undo()
	}
}

func BenchmarkDeleteLine20k(b *testing.B) {
	e := bigEditor(20000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.DeleteLine()
		e.Undo()
	}
}

func BenchmarkOpenLargeBuffer(b *testing.B) {
	var sb strings.Builder
	for i := 0; i < 20000; i++ {
		sb.WriteString("func handleThing(x int, name string) error { return fmt.Errorf(\"bad %d\", x) }\n")
	}
	text := sb.String()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e := editor.New()
		e.SetText(text)
	}
}

// A keystroke must cost the same whether the file is tiny or huge — the delta
// undo model makes edits O(change), not O(file). These lock that in: the
// number of allocations per InsertRune must not scale with file size.
func TestInsertAllocationsAreBounded(t *testing.T) {
	small := bigEditor(100)
	large := bigEditor(50000)

	smallAllocs := testing.AllocsPerRun(200, func() { small.InsertRune('x') })
	largeAllocs := testing.AllocsPerRun(200, func() { large.InsertRune('x') })

	if largeAllocs > smallAllocs+2 {
		t.Fatalf("InsertRune allocations scale with file size: small=%.0f large=%.0f", smallAllocs, largeAllocs)
	}
	if largeAllocs > 12 {
		t.Fatalf("InsertRune allocates too much per keystroke: %.0f", largeAllocs)
	}
}

func TestUndoAllocationsAreBounded(t *testing.T) {
	large := bigEditor(50000)
	allocs := testing.AllocsPerRun(200, func() {
		large.InsertRune('x')
		large.Undo()
	})
	if allocs > 20 {
		t.Fatalf("insert+undo allocates too much per cycle on a large file: %.0f", allocs)
	}
}
