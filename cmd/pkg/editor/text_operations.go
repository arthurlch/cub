package editor

import (
	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/arthurlch/cub/cmd/pkg/utils"
	"github.com/nsf/termbox-go"
)

func (es *EditorState) InsertRunes(keyEvent termbox.Event) {
    st := es.State

    if len(st.TextBuffer) == 0 {
        st.TextBuffer = append(st.TextBuffer, []rune{})
    }

    if keyEvent.Ch != 0 {
        st.UndoBuffer = append(st.UndoBuffer, state.UndoState{
            TextBuffer: utils.DeepCopyTextBuffer(st.TextBuffer),
            CurrentRow: st.CurrentRow,
            CurrentCol: st.CurrentCol,
        })
        st.RedoBuffer = nil

        line := st.TextBuffer[st.CurrentRow]
        newLine := append(line[:st.CurrentCol], append([]rune{keyEvent.Ch}, line[st.CurrentCol:]...)...)
        st.TextBuffer[st.CurrentRow] = newLine

        st.CurrentCol++
        st.Modified = true
    }
}

func (es *EditorState) DeleteRune() {
    st := es.State

    if st.CurrentCol > 0 {
        st.UndoBuffer = append(st.UndoBuffer, state.UndoState{
            TextBuffer: utils.DeepCopyTextBuffer(st.TextBuffer),
            CurrentRow: st.CurrentRow,
            CurrentCol: st.CurrentCol,
        })
        st.RedoBuffer = nil 

        line := st.TextBuffer[st.CurrentRow]
        newLine := append(line[:st.CurrentCol-1], line[st.CurrentCol:]...)
        st.TextBuffer[st.CurrentRow] = newLine

        st.CurrentCol--
        st.Modified = true
    } else if st.CurrentRow > 0 {
        st.UndoBuffer = append(st.UndoBuffer, state.UndoState{
            TextBuffer: utils.DeepCopyTextBuffer(st.TextBuffer),
            CurrentRow: st.CurrentRow,
            CurrentCol: st.CurrentCol,
        })
        st.RedoBuffer = nil 

        prevLine := st.TextBuffer[st.CurrentRow-1]
        currentLine := st.TextBuffer[st.CurrentRow]
        st.TextBuffer[st.CurrentRow-1] = append(prevLine, currentLine...)
        st.TextBuffer = append(st.TextBuffer[:st.CurrentRow], st.TextBuffer[st.CurrentRow+1:]...)

        st.CurrentRow--
        st.CurrentCol = len(prevLine)
        st.Modified = true
    }
}

func (es *EditorState) InsertNewLine() {
    st := es.State

    st.UndoBuffer = append(st.UndoBuffer, state.UndoState{
        TextBuffer: utils.DeepCopyTextBuffer(st.TextBuffer),
        CurrentRow: st.CurrentRow,
        CurrentCol: st.CurrentCol,
    })
    st.RedoBuffer = nil

    line := st.TextBuffer[st.CurrentRow]
    beforeCursor := line[:st.CurrentCol]
    afterCursor := line[st.CurrentCol:]

    st.TextBuffer[st.CurrentRow] = beforeCursor

    st.TextBuffer = append(st.TextBuffer[:st.CurrentRow+1], append([][]rune{afterCursor}, st.TextBuffer[st.CurrentRow+1:]...)...)

    st.CurrentRow++
    st.CurrentCol = 0
    st.Modified = true
}

func deleteCurrentLine(st *state.State) {
	if st.CurrentRow >= len(st.TextBuffer) {
			return
	}

	st.UndoBuffer = append(st.UndoBuffer, state.UndoState{
			TextBuffer: utils.DeepCopyTextBuffer(st.TextBuffer),
			CurrentRow: st.CurrentRow,
			CurrentCol: st.CurrentCol,
	})
	st.RedoBuffer = nil 

	st.TextBuffer = append(st.TextBuffer[:st.CurrentRow], st.TextBuffer[st.CurrentRow+1:]...)

	if len(st.TextBuffer) == 0 {
			st.TextBuffer = append(st.TextBuffer, []rune{})
			st.CurrentRow = 0
			st.CurrentCol = 0
	} else {
			if st.CurrentRow >= len(st.TextBuffer) {
					st.CurrentRow = len(st.TextBuffer) - 1
			}
			if st.CurrentCol > len(st.TextBuffer[st.CurrentRow]) {
					st.CurrentCol = len(st.TextBuffer[st.CurrentRow])
			}
	}

	st.Modified = true
}
