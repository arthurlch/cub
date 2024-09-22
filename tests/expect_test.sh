#!/usr/bin/expect -f

spawn ./bin/cub

send "i"                      ;# Switch to insert mode
send "Hello, World!"           ;# Type "Hello, World!"
send "\r"                      ;# Insert new line
send "This is a test file.\r"  ;# Type on a new line
send "\033"                    ;# Exit insert mode (press Esc)

send "\033[A"                  ;# Move up (arrow key up)
send "\033[B"                  ;# Move down (arrow key down)
send "\033[D"                  ;# Move left (arrow key left)
send "\033[C"                  ;# Move right (arrow key right)

send "^S"                      ;# Press Ctrl+S to save

send "dd"                      ;# Delete the current line (dd)

send "^U"                      ;# Press Ctrl+U to undo the deletion
send "^R"                      ;# Press Ctrl+R to redo the deletion

send "s"                       ;# Start selection
send "\033[C\033[C\033[C"      ;# Move right (selecting "Hello")
send "z"                       ;# End selection
send "c"                       ;# Copy selection to clipboard
send "\033[C\033[C"            ;# Move cursor right (move to the next word)
send "v"                       ;# Paste from clipboard

send "\033OH"                  ;# Home key (move to beginning of line)
send "\033OF"                  ;# End key (move to end of line)

send "\033[5~"                 ;# PgUp
send "\033[6~"                 ;# PgDn

send "^Q"                      ;# Press Ctrl+Q to quit

expect eof
