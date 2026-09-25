package buffer_test

import (
	"testing"

	"github.com/arthurlch/cub/internal/buffer"
	"github.com/stretchr/testify/assert"
)

func TestNewAndString(t *testing.T) {
	b := buffer.New("one\ntwo\nthree")
	assert.Equal(t, 3, b.LineCount())
	assert.Equal(t, "one", string(b.Line(0)))
	assert.Equal(t, "three", string(b.Line(2)))
	assert.Equal(t, "one\ntwo\nthree", b.String())
}

func TestInsertMidLine(t *testing.T) {
	b := buffer.New("abcdef")
	er, ec := b.Insert(0, 2, []rune("XY"))
	assert.Equal(t, "abXYcdef", b.String())
	assert.Equal(t, 0, er)
	assert.Equal(t, 4, ec)
}

func TestInsertWithNewline(t *testing.T) {
	b := buffer.New("abcdef")
	er, ec := b.Insert(0, 3, []rune("X\nY"))
	assert.Equal(t, "abcX\nYdef", b.String())
	assert.Equal(t, 1, er)
	assert.Equal(t, 1, ec)
	assert.Equal(t, 2, b.LineCount())
}

func TestInsertDoesNotCorruptNeighbours(t *testing.T) {
	b := buffer.New("line one\nline two\nline three")
	b.Insert(1, 4, []rune("XX"))
	assert.Equal(t, "line one", string(b.Line(0)))
	assert.Equal(t, "lineXX two", string(b.Line(1)))
	assert.Equal(t, "line three", string(b.Line(2)))
}

func TestDeleteWithinLine(t *testing.T) {
	b := buffer.New("hello world")
	removed := b.Delete(0, 5, 0, 11)
	assert.Equal(t, " world", string(removed))
	assert.Equal(t, "hello", b.String())
}

func TestDeleteAcrossLinesJoins(t *testing.T) {
	b := buffer.New("aaa\nbbb\nccc")
	removed := b.Delete(0, 2, 2, 1)
	assert.Equal(t, "aacc", b.String())
	assert.Equal(t, "a\nbbb\nc", string(removed))
	assert.Equal(t, 1, b.LineCount())
}

func TestInsertDeleteRoundTrip(t *testing.T) {
	b := buffer.New("the quick brown fox")
	er, ec := b.Insert(0, 4, []rune("very "))
	assert.Equal(t, "the very quick brown fox", b.String())
	b.Delete(0, 4, er, ec)
	assert.Equal(t, "the quick brown fox", b.String())
}

func TestDeleteNormalizesReversedSpan(t *testing.T) {
	b := buffer.New("hello world")
	removed := b.Delete(0, 11, 0, 5)
	assert.Equal(t, " world", string(removed))
	assert.Equal(t, "hello", b.String())
}

func TestEmptyBufferHasOneLine(t *testing.T) {
	b := buffer.New("")
	assert.Equal(t, 1, b.LineCount())
	assert.Equal(t, "", b.String())
}

func TestInsertClampsOutOfRange(t *testing.T) {
	b := buffer.New("ab")
	er, ec := b.Insert(9, 9, []rune("Z"))
	assert.Equal(t, "abZ", b.String())
	assert.Equal(t, 0, er)
	assert.Equal(t, 3, ec)
}

func TestInsertMultipleNewlines(t *testing.T) {
	b := buffer.New("ac")
	er, ec := b.Insert(0, 1, []rune("1\n2\n3"))
	assert.Equal(t, "a1\n2\n3c", b.String())
	assert.Equal(t, 2, er)
	assert.Equal(t, 1, ec)
	assert.Equal(t, 3, b.LineCount())
}

func TestLineOutOfRangeIsNil(t *testing.T) {
	b := buffer.New("x")
	assert.Nil(t, b.Line(-1))
	assert.Nil(t, b.Line(5))
	assert.Equal(t, 0, b.LineLen(5))
}
