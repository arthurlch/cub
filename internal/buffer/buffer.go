package buffer

// Buffer is a simple text buffer that stores lines of text as slices of runes.
// It works great and I have no plan to work on it !

import "strings"

type Buffer struct {
	lines [][]rune
}

func New(text string) *Buffer {
	b := &Buffer{}
	b.SetText(text)
	return b
}

func (b *Buffer) SetText(text string) {
	parts := strings.Split(text, "\n")
	b.lines = make([][]rune, len(parts))
	for i, p := range parts {
		b.lines[i] = []rune(p)
	}
	if len(b.lines) == 0 {
		b.lines = [][]rune{{}}
	}
}

func (b *Buffer) LineCount() int { return len(b.lines) }

func (b *Buffer) Line(i int) []rune {
	if i < 0 || i >= len(b.lines) {
		return nil
	}
	return b.lines[i]
}

func (b *Buffer) LineLen(i int) int {
	if i < 0 || i >= len(b.lines) {
		return 0
	}
	return len(b.lines[i])
}

func (b *Buffer) String() string {
	var sb strings.Builder
	for i, ln := range b.lines {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(string(ln))
	}
	return sb.String()
}

func (b *Buffer) Insert(row, col int, text []rune) (endRow, endCol int) {
	if len(b.lines) == 0 {
		b.lines = [][]rune{{}}
	}
	row = clamp(row, 0, len(b.lines)-1)
	line := b.lines[row]
	col = clamp(col, 0, len(line))

	segments := splitRunes(text)
	if len(segments) == 1 {
		merged := make([]rune, 0, len(line)+len(segments[0]))
		merged = append(merged, line[:col]...)
		merged = append(merged, segments[0]...)
		merged = append(merged, line[col:]...)
		b.lines[row] = merged
		return row, col + len(segments[0])
	}

	head := append([]rune{}, line[:col]...)
	tail := append([]rune{}, line[col:]...)

	first := append(head, segments[0]...)
	last := append(append([]rune{}, segments[len(segments)-1]...), tail...)

	inserted := make([][]rune, 0, len(segments))
	inserted = append(inserted, first)
	for i := 1; i < len(segments)-1; i++ {
		inserted = append(inserted, append([]rune{}, segments[i]...))
	}
	inserted = append(inserted, last)

	rest := append([][]rune{}, b.lines[row+1:]...)
	b.lines = append(b.lines[:row], inserted...)
	b.lines = append(b.lines, rest...)

	endRow = row + len(segments) - 1
	endCol = len(segments[len(segments)-1])
	return endRow, endCol
}

func (b *Buffer) Delete(sr, sc, er, ec int) []rune {
	sr, sc, er, ec = normalize(sr, sc, er, ec)
	sr = clamp(sr, 0, len(b.lines)-1)
	er = clamp(er, 0, len(b.lines)-1)
	sc = clamp(sc, 0, len(b.lines[sr]))
	ec = clamp(ec, 0, len(b.lines[er]))

	if sr == er {
		line := b.lines[sr]
		removed := append([]rune{}, line[sc:ec]...)
		merged := append([]rune{}, line[:sc]...)
		merged = append(merged, line[ec:]...)
		b.lines[sr] = merged
		return removed
	}

	removed := make([]rune, 0)
	removed = append(removed, b.lines[sr][sc:]...)
	for r := sr + 1; r < er; r++ {
		removed = append(removed, '\n')
		removed = append(removed, b.lines[r]...)
	}
	removed = append(removed, '\n')
	removed = append(removed, b.lines[er][:ec]...)

	merged := append([]rune{}, b.lines[sr][:sc]...)
	merged = append(merged, b.lines[er][ec:]...)

	rest := append([][]rune{}, b.lines[er+1:]...)
	b.lines = append(b.lines[:sr], merged)
	b.lines = append(b.lines, rest...)

	return removed
}

func splitRunes(text []rune) [][]rune {
	segments := [][]rune{{}}
	for _, r := range text {
		if r == '\n' {
			segments = append(segments, []rune{})
		} else {
			last := len(segments) - 1
			segments[last] = append(segments[last], r)
		}
	}
	return segments
}

func normalize(sr, sc, er, ec int) (int, int, int, int) {
	if sr > er || (sr == er && sc > ec) {
		return er, ec, sr, sc
	}
	return sr, sc, er, ec
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
