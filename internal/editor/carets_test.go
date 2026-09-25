package editor_test

import (
	"testing"

	"github.com/arthurlch/cub/internal/editor"
	"github.com/stretchr/testify/assert"
)

func TestAddCaretAndInsertRune(t *testing.T) {
	e := editor.New()
	e.SetText("aaa\nbbb\nccc")
	e.AddCaretBelow()
	e.AddCaretBelow()
	assert.Equal(t, 3, e.CaretCount())

	e.InsertRune('X')
	assert.Equal(t, "Xaaa\nXbbb\nXccc", e.Buffer().String())
}

func TestMultiCursorGroupedUndoRedo(t *testing.T) {
	e := editor.New()
	e.SetText("aaa\nbbb")
	e.AddCaretBelow()
	e.InsertRune('Z')
	assert.Equal(t, "Zaaa\nZbbb", e.Buffer().String())

	e.Undo()
	assert.Equal(t, "aaa\nbbb", e.Buffer().String())
	assert.False(t, e.CanUndo())
	assert.Equal(t, 1, e.CaretCount(), "undo collapses to a single cursor")

	e.Redo()
	assert.Equal(t, "Zaaa\nZbbb", e.Buffer().String())
}

func TestMultiCursorNewlineShiftsLowerCarets(t *testing.T) {
	e := editor.New()
	e.SetText("aaa\nbbb\nccc")
	e.LineEnd()
	e.AddCaretBelow()
	e.InsertNewline()
	assert.Equal(t, "aaa\n\nbbb\n\nccc", e.Buffer().String())
}

func TestMultiCursorBackspace(t *testing.T) {
	e := editor.New()
	e.SetText("aXa\nbXb")
	e.MoveRight()
	e.MoveRight()
	e.AddCaretBelow()
	e.Backspace()
	assert.Equal(t, "aa\nbb", e.Buffer().String())
}

func TestMultiCursorMoveRightMovesAll(t *testing.T) {
	e := editor.New()
	e.SetText("abc\ndef")
	e.AddCaretBelow()
	e.MoveRight()
	carets := e.Carets()
	assert.Equal(t, [][2]int{{0, 1}, {1, 1}}, carets)
}

func TestAddCaretBelowClampsColumn(t *testing.T) {
	e := editor.New()
	e.SetText("longline\nab")
	e.LineEnd()
	e.AddCaretBelow()
	carets := e.Carets()
	assert.Equal(t, 2, carets[1][1], "caret column clamps to the shorter line")
}

func TestAddCaretBelowStopsAtLastLine(t *testing.T) {
	e := editor.New()
	e.SetText("only")
	e.AddCaretBelow()
	assert.Equal(t, 1, e.CaretCount(), "no line below means no caret added")
}

func TestClearCaretsCollapses(t *testing.T) {
	e := editor.New()
	e.SetText("aaa\nbbb")
	e.AddCaretBelow()
	assert.Equal(t, 2, e.CaretCount())
	e.ClearCarets()
	assert.Equal(t, 1, e.CaretCount())
}

func TestSingleCursorUndoUnaffectedByGrouping(t *testing.T) {
	e := editor.New()
	e.SetText("abc")
	e.InsertRune('1')
	e.InsertRune('2')
	e.Undo()
	assert.Equal(t, "1abc", e.Buffer().String(), "single-cursor edits stay independent undo steps")
}
