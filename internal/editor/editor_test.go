package editor_test

import (
	"testing"

	"github.com/arthurlch/cub/internal/editor"
	"github.com/stretchr/testify/assert"
)

func at(e *editor.Editor, row, col int) {
	for r, _ := e.Cursor(); r < row; r, _ = e.Cursor() {
		e.MoveDown()
	}
	for _, c := e.Cursor(); c < col; _, c = e.Cursor() {
		e.MoveRight()
	}
}

func TestInsertUndoRedo(t *testing.T) {
	e := editor.New()
	e.SetText("abcdef")
	at(e, 0, 2)

	e.InsertRune('X')
	assert.Equal(t, "abXcdef", e.Buffer().String())
	r, c := e.Cursor()
	assert.Equal(t, 0, r)
	assert.Equal(t, 3, c)

	e.Undo()
	assert.Equal(t, "abcdef", e.Buffer().String())
	_, c = e.Cursor()
	assert.Equal(t, 2, c)

	e.Redo()
	assert.Equal(t, "abXcdef", e.Buffer().String())
}

func TestBackspaceJoinsLines(t *testing.T) {
	e := editor.New()
	e.SetText("hello\nworld")
	at(e, 1, 0)
	e.Backspace()
	assert.Equal(t, "helloworld", e.Buffer().String())
	r, c := e.Cursor()
	assert.Equal(t, 0, r)
	assert.Equal(t, 5, c)
}

func TestNewlineSplits(t *testing.T) {
	e := editor.New()
	e.SetText("abcdef")
	at(e, 0, 3)
	e.InsertNewline()
	assert.Equal(t, "abc\ndef", e.Buffer().String())
	r, c := e.Cursor()
	assert.Equal(t, 1, r)
	assert.Equal(t, 0, c)
}

func TestCopyDoesNotMutateAndPasteMidLine(t *testing.T) {
	e := editor.New()
	e.SetText("XY")
	e.SelectAll()
	e.Copy()
	assert.Equal(t, "XY", e.Buffer().String(), "copy must not modify the document")

	e.SetText("abcdef")
	at(e, 0, 2)
	e.Paste()
	assert.Equal(t, "abXYcdef", e.Buffer().String())
}

func TestDeleteLine(t *testing.T) {
	e := editor.New()
	e.SetText("one\ntwo\nthree")
	at(e, 1, 0)
	e.DeleteLine()
	assert.Equal(t, "one\nthree", e.Buffer().String())
}

func TestSelectAllThenDelete(t *testing.T) {
	e := editor.New()
	e.SetText("alpha\nbeta")
	e.SelectAll()
	e.DeleteSelection()
	assert.Equal(t, "", e.Buffer().String())
	e.Undo()
	assert.Equal(t, "alpha\nbeta", e.Buffer().String())
}

func TestYankLineAndPaste(t *testing.T) {
	e := editor.New()
	e.SetText("one\ntwo\nthree")
	e.YankLine()
	e.Paste()
	assert.Equal(t, "one\none\ntwo\nthree", e.Buffer().String(), "yy then p pastes the line below")
}

func TestDeleteLineYanksForPaste(t *testing.T) {
	e := editor.New()
	e.SetText("keep\ncut\nkeep2")
	at(e, 1, 0)
	e.DeleteLine()
	assert.Equal(t, "keep\nkeep2", e.Buffer().String())
	e.PasteBefore()
	assert.Equal(t, "keep\ncut\nkeep2", e.Buffer().String(), "dd cuts the line into the register; P pastes it back above")
}

func TestPasteBefore(t *testing.T) {
	e := editor.New()
	e.SetText("a\nb")
	e.YankLine()
	at(e, 1, 0)
	e.PasteBefore()
	assert.Equal(t, "a\na\nb", e.Buffer().String(), "P pastes the yanked line above")
}

func TestOpenBelowEntersInsert(t *testing.T) {
	e := editor.New()
	e.SetText("first\nsecond")
	e.OpenBelow()
	assert.Equal(t, editor.InsertMode, e.Mode())
	e.InsertRune('x')
	assert.Equal(t, "first\nx\nsecond", e.Buffer().String())
}

func TestOpenAbove(t *testing.T) {
	e := editor.New()
	e.SetText("first\nsecond")
	at(e, 1, 0)
	e.OpenAbove()
	e.InsertRune('x')
	assert.Equal(t, "first\nx\nsecond", e.Buffer().String())
}

func TestAppendMovesPastCursor(t *testing.T) {
	e := editor.New()
	e.SetText("ab")
	e.Append()
	e.InsertRune('X')
	assert.Equal(t, "aXb", e.Buffer().String())
}

func TestDeleteCharUnderCursor(t *testing.T) {
	e := editor.New()
	e.SetText("abc")
	e.DeleteChar()
	assert.Equal(t, "bc", e.Buffer().String())
}

func TestDeleteForward(t *testing.T) {
	e := editor.New()
	e.SetText("abc")
	e.DeleteForward()
	assert.Equal(t, "bc", e.Buffer().String())
}

func TestGotoLineClamps(t *testing.T) {
	e := editor.New()
	e.SetText("l1\nl2\nl3")
	e.GotoLine(99)
	r, _ := e.Cursor()
	assert.Equal(t, 2, r)
	e.GotoLine(0)
	r, _ = e.Cursor()
	assert.Equal(t, 0, r)
}

func TestWordNext(t *testing.T) {
	e := editor.New()
	e.SetText("foo bar baz")
	e.WordNext()
	_, c := e.Cursor()
	assert.Equal(t, 4, c)
}

func TestWordNextCrossesLineWithoutBoundary(t *testing.T) {
	e := editor.New()
	e.SetText("aaa\nbbb")
	e.WordNext()
	r, c := e.Cursor()
	assert.Equal(t, 1, r)
	assert.Equal(t, 0, c)
}

func TestEmptySelectionCutIsNoop(t *testing.T) {
	e := editor.New()
	e.SetText("hello")
	e.StartSelection()
	e.Cut()
	assert.Equal(t, "hello", e.Buffer().String())
	assert.False(t, e.Modified(), "an empty cut must not mark the buffer modified")
}

func TestCutPaste(t *testing.T) {
	e := editor.New()
	e.SetText("hello world")
	e.StartSelection()
	for i := 0; i < 5; i++ {
		e.MoveRight()
	}
	e.Cut()
	assert.Equal(t, " world", e.Buffer().String())
	e.LineEnd()
	e.Paste()
	assert.Equal(t, " worldhello", e.Buffer().String())
}
