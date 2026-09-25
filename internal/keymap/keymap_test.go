package keymap

import (
	"testing"

	"github.com/arthurlch/cub/internal/editor"
	"github.com/arthurlch/cub/internal/input"
)

func TestParseChords(t *testing.T) {
	cases := []struct {
		chord string
		want  input.Event
	}{
		{"ctrl+q", input.Event{Key: input.KeyRune, Rune: 'q', Ctrl: true}},
		{"ctrl+right", input.Event{Key: input.KeyRight, Ctrl: true}},
		{"alt+x", input.Event{Key: input.KeyRune, Rune: 'x', Alt: true}},
		{"i", input.Event{Key: input.KeyRune, Rune: 'i'}},
		{"$", input.Event{Key: input.KeyRune, Rune: '$'}},
		{"space", input.Event{Key: input.KeyRune, Rune: ' '}},
		{"esc", input.Event{Key: input.KeyEsc}},
	}
	for _, c := range cases {
		got, ok := Parse(c.chord)
		if !ok {
			t.Fatalf("Parse(%q) failed", c.chord)
		}
		if got != c.want {
			t.Errorf("Parse(%q) = %+v, want %+v", c.chord, got, c.want)
		}
	}
}

func TestGlobalDispatch(t *testing.T) {
	k := Default()
	ev, _ := Parse("ctrl+s")
	if got := k.Global(ev); got != Save {
		t.Errorf("ctrl+s = %q, want %q", got, Save)
	}
	if got := k.Global(input.Event{Key: input.KeyRune, Rune: 'i'}); got != "" {
		t.Errorf("view key leaked into global: %q", got)
	}
}

func TestViewInsertEntersInsertMode(t *testing.T) {
	k := Default()
	e := editor.New()
	k.HandleView(e, input.Event{Key: input.KeyRune, Rune: 'i'})
	if e.Mode() != editor.InsertMode {
		t.Fatalf("expected insert mode after 'i'")
	}
}

func TestSequenceDeleteLine(t *testing.T) {
	k := Default()
	e := editor.New()
	e.SetText("one\ntwo\nthree")
	k.HandleView(e, input.Event{Key: input.KeyRune, Rune: 'd'})
	k.HandleView(e, input.Event{Key: input.KeyRune, Rune: 'd'})
	if e.LineCount() != 2 {
		t.Fatalf("dd should delete a line, got %d lines", e.LineCount())
	}
}

func TestParseRejectsUnknownModifier(t *testing.T) {
	if _, ok := Parse("hyper+x"); ok {
		t.Fatal("unknown modifier should not parse")
	}
	if _, ok := Parse("notakey"); ok {
		t.Fatal("multi-rune bare key should not parse")
	}
}

func TestCursorAddIsGlobal(t *testing.T) {
	k := Default()
	ev, _ := Parse("ctrl+down")
	if got := k.Global(ev); got != CursorAddBelow {
		t.Errorf("ctrl+down = %q, want %q", got, CursorAddBelow)
	}
}

func TestJumpClearsCarets(t *testing.T) {
	k := Default()
	e := editor.New()
	e.SetText("aaa\nbbb\nccc")
	e.AddCaretBelow()
	if e.CaretCount() != 2 {
		t.Fatalf("setup: expected 2 carets")
	}
	k.HandleView(e, input.Event{Key: input.KeyRune, Rune: 'w'})
	if e.CaretCount() != 1 {
		t.Fatalf("word motion should collapse carets, got %d", e.CaretCount())
	}
}

func TestInsertActionKeepsCarets(t *testing.T) {
	k := Default()
	e := editor.New()
	e.SetText("aaa\nbbb")
	e.AddCaretBelow()
	k.HandleView(e, input.Event{Key: input.KeyRune, Rune: 'i'})
	if e.CaretCount() != 2 {
		t.Fatalf("entering insert must preserve carets, got %d", e.CaretCount())
	}
}

func TestEscClearsCarets(t *testing.T) {
	k := Default()
	e := editor.New()
	e.SetText("aaa\nbbb")
	e.AddCaretBelow()
	k.HandleView(e, input.Event{Key: input.KeyEsc})
	if e.CaretCount() != 1 {
		t.Fatalf("esc should collapse carets, got %d", e.CaretCount())
	}
}

func rune_(r rune) input.Event { return input.Event{Key: input.KeyRune, Rune: r} }

func TestVisualModeYankAndDelete(t *testing.T) {
	k := Default()
	e := editor.New()
	e.SetText("hello world")
	k.HandleView(e, rune_('v')) // start visual
	if !e.SelectionActive() {
		t.Fatal("v should start a visual selection")
	}
	for i := 0; i < 5; i++ {
		k.HandleView(e, rune_('l')) // extend over "hello"
	}
	k.HandleView(e, rune_('d')) // delete selection
	if e.Buffer().String() != " world" {
		t.Fatalf("visual d should delete the selection, got %q", e.Buffer().String())
	}
	if e.SelectionActive() {
		t.Fatal("selection should end after delete")
	}
}

func TestUndoIsLowercaseU(t *testing.T) {
	k := Default()
	e := editor.New()
	e.SetText("abc")
	e.InsertRune('X')
	k.HandleView(e, rune_('u'))
	if e.Buffer().String() != "abc" {
		t.Fatalf("u should undo, got %q", e.Buffer().String())
	}
}

func TestGotoBottomAndLine(t *testing.T) {
	k := Default()
	e := editor.New()
	e.SetText("l1\nl2\nl3\nl4")
	k.HandleView(e, rune_('G'))
	if r, _ := e.Cursor(); r != 3 {
		t.Fatalf("G should jump to the last line, got row %d", r)
	}
	k.HandleView(e, rune_('2'))
	k.HandleView(e, rune_('G'))
	if r, _ := e.Cursor(); r != 1 {
		t.Fatalf("2G should jump to line 2, got row %d", r)
	}
}

func TestOverrideBinding(t *testing.T) {
	chords := DefaultChords()
	chords[Insert] = "o"
	k := New(chords)
	e := editor.New()
	k.HandleView(e, input.Event{Key: input.KeyRune, Rune: 'o'})
	if e.Mode() != editor.InsertMode {
		t.Fatalf("custom insert binding 'o' did not enter insert mode")
	}

	other := editor.New()
	k.HandleView(other, input.Event{Key: input.KeyRune, Rune: 'i'})
	if other.Mode() != editor.ViewMode {
		t.Fatalf("old 'i' binding should no longer enter insert mode")
	}
}
