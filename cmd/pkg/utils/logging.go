package utils

import (
	"log"
	"os"

	"github.com/arthurlch/cub/cmd/pkg/state"
	"github.com/nsf/termbox-go"
)

var Logger *log.Logger

func InitLogger() {
	file, err := os.OpenFile("editor.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	Logger = log.New(file, "", log.LstdFlags)
}

func LogBufferState(st *state.State, context string) {
	Logger.Printf("%s - Cursor Position - Row: %d, Col: %d", context, st.CurrentRow, st.CurrentCol)
	Logger.Printf("Mode: %v, SelectionActive: %v", st.Mode, st.SelectionActive)
	Logger.Printf("Rows: %d, Cols: %d, OffsetRow: %d, OffsetCol: %d", st.Rows, st.Cols, st.OffsetRow, st.OffsetCol)
	Logger.Printf("Modified: %v, QuitKey: %v", st.Modified, st.QuitKey)
	if st.SelectionActive {
			Logger.Printf("Selection - StartRow: %d, StartCol: %d, EndRow: %d, EndCol: %d", 
					st.StartRow, st.StartCol, st.EndRow, st.EndCol)
	}
}

func LogKeyPress(context string, keyEvent termbox.Event) {
	Logger.Printf("%s - Key: %+v (Ch: %c, Key: %v)", context, keyEvent, keyEvent.Ch, keyEvent.Key)
}
