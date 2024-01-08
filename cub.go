package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

var mode int
var ROWS, COLS int
var offsetRow, offsetCol int 
var currentX, currentY int
var currentRow, currentCol int
var source_file string
var text_buffer = [][]rune{}
var undo_buffer = [][]rune{}
var copy_buffer = [][]rune{} 
var modified bool
var quitKey termbox.Key


func read_file (filename string) {
	file, err := os.Open(filename)

	if err != nil {
		source_file = filename
		text_buffer = append(text_buffer, []rune{}); return  
	}
	defer file.Close() 
	scanner := bufio.NewScanner(file) 
	lineNumber := 0 

	for scanner.Scan() {
		line := scanner.Text() 
		text_buffer = append(text_buffer, []rune{})
		for i := 0; i < len(line); i++ {
			text_buffer[lineNumber] = append(text_buffer[lineNumber], rune(line[i])) 
		}
		lineNumber++
	}
	if lineNumber == 0 {
		text_buffer = append(text_buffer, []rune{})
	}
}


func display_text_buffer() {
	 var row, col int 
	 for row = 0; row < ROWS; row++ {
		text_buffer_row := row + offsetRow		
		for col = 0; col < COLS; col++ {
			text_buffer_col := col + offsetCol
			if text_buffer_row >= 0 && text_buffer_row < len(text_buffer) && text_buffer_col < len(text_buffer[text_buffer_row]) {
				if text_buffer[text_buffer_row][text_buffer_col] != '\t' {
					termbox.SetChar(col, row, text_buffer[text_buffer_row][text_buffer_col])
				} else { 
					termbox.SetCell(col, row, rune(' '), termbox.ColorDefault, termbox.ColorDefault) 
				} 
			} else if row+offsetCol > len(text_buffer) {
				termbox.SetCell(0, row, rune('*'), termbox.ColorLightMagenta, termbox.ColorDefault)
				termbox.SetChar(col, row, rune('\n'))
			}
		}
	}
}

func scroll_text_buffer() {
	if currentRow < offsetRow {
			offsetRow = currentRow
	}

	if currentRow >= offsetRow + ROWS {
		offsetRow = currentRow - ROWS + 1
	}

	if currentCol < offsetCol {
			offsetCol = currentCol
	}

	if currentCol >= offsetCol + COLS {
			offsetCol = currentCol - COLS + 1
			if offsetCol < 0 {
					offsetCol = 0
			}
	}
}


func display_status_bar() {
	var mode_status string 
	var file_status string
	var copy_status string 
	var undo_status string 
	var cursor_status string

	filename_len := len(source_file) 

	if filename_len > 14 {
		filename_len = 14
	} 

	file_status = source_file[:filename_len] + " - " + strconv.Itoa(len(text_buffer)) + " lines "
	if modified {
		file_status += " modified " 
	} else {
		file_status += " saved "
	}

	if mode > 0 {
		mode_status = "EDIT: "
	} else {
		mode_status = "VIEW: "
	}

	cursor_status = "Row " + strconv.Itoa(currentRow + 1) + " Col " + strconv.Itoa(currentCol) + " "

	if (len(copy_buffer ) > 0) {
		copy_status = "[copy]" 
	} 

	if len(undo_buffer) > 0 {
		undo_status = "[undo]"
	}

	used_space := len(mode_status) + len(cursor_status) + len(copy_status) + len(file_status) + len(undo_status)
	spaces := strings.Repeat(" ", COLS - used_space)
	message := mode_status + file_status + copy_status + undo_status + spaces + cursor_status
	print_message(0, ROWS, termbox.ColorBlack, termbox.ColorWhite, message)
}

func get_key() termbox.Event {
	var key_event termbox.Event 
	switch event := termbox.PollEvent(); event.Type {
		case termbox.EventKey: key_event = event
		case termbox.EventError: panic(event.Err)
	} 
	return key_event
}

func process_keypress() {
	key_event := get_key() 

	if quitKey == termbox.KeyEsc && key_event.Ch == 'q' {
		termbox.Close()
		os.Exit(0)
	}

	if key_event.Key == termbox.KeyEsc {
			quitKey = termbox.KeyEsc
			return  
	} else {
			quitKey = 0 
	}

	if key_event.Key == termbox.KeyEsc { mode = 0
	} else if key_event.Ch != 0 {

		if mode == 1 {
			insert_runes(key_event); 
			modified = true
		} else {
			switch key_event.Ch {
			case 'q': termbox.Close()
			case 'e': mode = 1
			}
		}

 	} else {
		switch key_event.Key {
		case termbox.KeyArrowUp: if currentRow != 0 { currentRow--}
		case termbox.KeyArrowDown: if currentRow < len(text_buffer) - 1 {currentRow++} 
		case termbox.KeyArrowLeft: 
			if currentCol != 0 {
				currentCol--
			} else if currentRow > 0 {
				currentRow--
				currentCol = len(text_buffer[currentRow])
			}
		case termbox.KeyArrowRight: 
			if currentCol < len(text_buffer[currentRow]) {
				currentCol++
			} else if currentRow < len(text_buffer) - 1  {
				currentRow++
				currentCol = 0
			}
		case termbox.KeyHome: currentCol = 0
		case termbox.KeyEnd: currentCol = len(text_buffer[currentRow])
		case termbox.KeyPgup: 
			if currentRow - int(ROWS / 4) > 0 {
				currentRow -= int(ROWS / 4)
			}
		case termbox.KeyPgdn: 
			if currentRow + int(ROWS / 4) < len(text_buffer) - 1 {
				currentRow += int(ROWS / 4  )
			}
		case termbox.KeyTab:
			if mode == 1 {
				for i := 0; i < 4; i++ { 
					insert_runes(key_event) 
				}
				modified = true 
			}
		case termbox.KeySpace:
			if mode == 1 {
				insert_runes(key_event)
				modified = true
			}
		case termbox.KeyBackspace:
			delete_rune()
			modified = true
		
	
		case termbox.KeyBackspace2:
			delete_rune()
			modified = true

		case termbox.KeyEnter: 
			if mode == 1 {
				insert_line()
				modified = true
			}
		}
		
		if currentCol > len(text_buffer[currentRow ]) {
			currentCol = len(text_buffer[currentRow])
		} 
	}
}

func insert_runes(event termbox.Event) {
	if currentRow >= len(text_buffer) {
			return
	}

	newRow := make([]rune, len(text_buffer[currentRow]) + 1)
	copy(newRow, text_buffer[currentRow][:currentCol])

	if event.Key == termbox.KeySpace {
			newRow[currentCol] = ' '
	} else if event.Key == termbox.KeyTab {
			newRow[currentCol] = ' '
	} else {
			newRow[currentCol] = event.Ch
	}

	copy(newRow[currentCol+1:], text_buffer[currentRow][currentCol:])
	text_buffer[currentRow] = newRow
	currentCol++
}

func delete_rune() {
	if currentCol > 0 && currentRow < len(text_buffer) {
			currentCol--
			newLine := make([]rune, len(text_buffer[currentRow])-1)
			copy(newLine, text_buffer[currentRow][:currentCol])
			copy(newLine[currentCol:], text_buffer[currentRow][currentCol+1:])
			text_buffer[currentRow] = newLine
	} else if currentRow > 0 {
		append_line := make([]rune, len(text_buffer[currentRow]))
		copy(append_line, text_buffer[currentRow][currentCol:])
		new_text_buffer := make([][]rune, len(text_buffer)-1)
		copy(new_text_buffer[:currentRow], text_buffer[:currentRow])
		copy(new_text_buffer[:currentRow], text_buffer[currentRow +1:])
		currentCol = len(text_buffer[currentRow])
		insert_line := make([]rune, len(text_buffer[currentRow]) + len(append_line))
		copy(insert_line[:len(text_buffer[currentRow])], text_buffer[currentRow])
		copy(insert_line[len(text_buffer[currentRow]): ], append_line)
		text_buffer[currentRow] = insert_line
	}
}

func insert_line() {
	right_line := make([]rune, len(text_buffer[currentRow][currentCol:]))
	copy(right_line, text_buffer[currentRow][currentCol:])
	left_line := make([]rune, len(text_buffer[currentRow][:currentCol]))
	copy(left_line, text_buffer[currentRow][:currentCol])
	text_buffer[currentRow] = left_line
	currentRow++
	currentCol = 0
	new_text_buffer := make([][]rune, len(text_buffer) + 1)
	copy(new_text_buffer, text_buffer[:currentRow])
	new_text_buffer[currentRow] = right_line 
	copy(new_text_buffer[currentRow + 1:], text_buffer[currentRow:])
	text_buffer = new_text_buffer
}

func print_message(col, row int, forground, background termbox.Attribute, message string) {
	for _, ch := range message {
		termbox.SetCell(col, row, ch, forground, background)
		col += runewidth.RuneWidth(ch)
	}
}
func run_text_editor() {
	err := termbox.Init()
	if err != nil {fmt.Println(err); os.Exit(1) } 
	if len(os.Args) > 1 {
		source_file = os.Args[1]
		read_file(source_file) 
	} else {
		text_buffer = append(text_buffer, []rune{})
	}
	for {
		COLS, ROWS = termbox.Size(); ROWS --
		if COLS < 78 { COLS = 78 }
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		scroll_text_buffer()
		display_text_buffer()
		display_status_bar()
		termbox.SetCursor(currentCol - offsetCol, currentRow - offsetRow)
		termbox.Flush()
		process_keypress()
	}
}
 
func main () {
	run_text_editor()
}

