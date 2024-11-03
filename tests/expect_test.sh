#!/usr/bin/expect -f

if {[info exists env(GITHUB_ACTIONS)]} {
    set github_actions $env(GITHUB_ACTIONS)
} else {
    set github_actions ""
}

if {[info exists env(RUNNER_OS)]} {
    set runner_os $env(RUNNER_OS)
} else {
    set runner_os ""
}

set binary_path "./bin/cub_linux-amd64"  ;# default for Linux
if {$github_actions eq "true"} {
    if {[string match {*darwin*} $runner_os]} {
        set binary_path "./bin/cub_darwin-amd64"
    } elseif {[string match {*Windows*} $runner_os]} {
        set binary_path "./bin/cub_windows-amd64"
    }
} else {
    if {[info exists env(GOOS)] && [string equal $env(GOOS) "darwin"]} {
        set binary_path "./bin/cub_darwin-amd64"
    } elseif {[info exists env(GOOS)] && [string equal $env(GOOS) "windows"]} {
        set binary_path "./bin/cub_windows-amd64"
    }
}

if {![file exists $binary_path]} {
    puts "Error: $binary_path not found!"
    exit 1
}

spawn $binary_path

send "i"                      ;# Enter insert mode
send "Hello, Cub!"            ;# Type "Hello, Cub!"
send "\r"                     ;# Insert new line
send "Testing basic insertions.\r" ;# Another line
send "\033"                   ;# Exit insert mode (press Esc)

send "\033\[A"                ;# Arrow key up
send "\033\[B"                ;# Arrow key down
send "\033\[D"                ;# Arrow key left
send "\033\[C"                ;# Arrow key right

send "\x13"                   ;# Press Ctrl+S to save

send "dd"                     ;# Delete the current line (dd)
send "\x15"                   ;# Ctrl+U to undo deletion
send "\x12"                   ;# Ctrl+R to redo deletion

send "s"                      ;# Start selection
send "\033\[C\033\[C\033\[C"  ;# Move right to select "Cub"
send "z"                      ;# End selection
send "c"                      ;# Copy selection
send "\033\[C\033\[C"         ;# Move right
send "v"                      ;# Paste clipboard content

send "\033OH"                 ;# Home key
send "\033OF"                 ;# End key

send "\033\[5~"               ;# Page Up
send "\033\[6~"               ;# Page Down

send "a"                      ;# Start "Select All" action
send "a"                      ;# Confirm "Select All"

send "\x11"                   ;# Press Ctrl+Q to quit

expect eof
