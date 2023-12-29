package main

import (
	"fmt"
	"os"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

var ROWS, COLS int
var offsetX, offsetY int 

var text_buffer = [][]rune {
	{'t', 'e', 's','t'},
	{'1'},
}

func display_text_buffer() {
	 var row, col int 
	 for row = 0; row < ROWS; row++ {
		text_buffer_row := row + offsetY
		for col = 0; col < COLS; col++ {
			text_buffer_col := col + offsetX
			if text_buffer_row >= 0 && text_buffer_row < len(text_buffer) && text_buffer_col < len(text_buffer[text_buffer_row]) {
				if text_buffer[text_buffer_row][text_buffer_col] != '\t' {
					termbox.SetChar(col, row, text_buffer[text_buffer_row][text_buffer_col])
				} else { 
					termbox.SetCell(col, row, rune(' '), termbox.ColorDefault, termbox.ColorLightMagenta) 
				} 
			} else if row+offsetY > len(text_buffer) {
				termbox.SetCell(0, row, rune('*'), termbox.ColorBlue, termbox.ColorDefault)
				termbox.SetChar(col, row, rune('\n'))
			}
		}
	 }
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
	for {
		print_message(25, 11, termbox.ColorDefault, termbox.ColorDefault, "Cub -- The simple text editor." )
		termbox.Flush()
		event := termbox.PollEvent()
		if event.Type == termbox.EventKey && event.Key == termbox.KeyEsc {
			termbox.Close()
		}
	}
}

func main () {
	run_text_editor()
}

