#!/usr/bin/expect -f

set binary_path "./bin/cub_linux-amd64"  ;# default for Linux

if {[string equal $::env.GITHUB_ACTIONS "true"]} {
    if {[string match {*darwin*} $::env.RUNNER_OS]} {
        set binary_path "./bin/cub_darwin-amd64"
    } elseif {[string match {*Windows*} $::env.RUNNER_OS]} {
        set binary_path "./bin/cub_windows-amd64"
    }
} else {
    if {[info exists ::env.GOOS] && [string equal $::env.GOOS "darwin"]} {
        set binary_path "./bin/cub_darwin-amd64"
    } elseif {[info exists ::env.GOOS] && [string equal $::env.GOOS "windows"]} {
        set binary_path "./bin/cub_windows-amd64"
    }
}

if {![file exists $binary_path]} {
    puts "Error: $binary_path not found!"
    exit 1
}

spawn $binary_path

send "i"                      ;# Switch to insert mode
send "Hello, World!"           ;# Type "Hello, World!"
send "\r"                      ;# Insert new line
send "This is a test file.\r"  ;# Type on a new line
send "\033"                    ;# Exit insert mode (press Esc)

send "\033[A"                  ;# Move up (arrow key up)
send "\033[B"                  ;# Move down (arrow key down)
send "\033[D"                  ;# Move left (arrow key left)
send "\033[C"                  ;# Move right (arrow key right)

send "\x13"                    ;# Press Ctrl+S to save

send "dd"                      ;# Delete the current line (dd)

send "\x15"                    ;# Press Ctrl+U to undo the deletion
send "\x12"                    ;# Press Ctrl+R to redo the deletion

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

send "\x11"                    ;# Press Ctrl+Q to quit

expect eof
