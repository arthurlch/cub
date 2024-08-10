# Cub Text Editor

<p align="center">
  <img src="https://i.ibb.co/qpZMHns/logo-color.png" alt="Cub Logo" width="200" height="200">
</p>

![GitHub release (latest by date)](https://img.shields.io/github/v/release/arthurlch/cub)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/arthurlch/cub/ci.yml)
![GitHub issues](https://img.shields.io/github/issues/arthurlch/cub)
![GitHub license](https://img.shields.io/github/license/arthurlch/cub)
![GitHub stars](https://img.shields.io/github/stars/arthurlch/cub)

## Overview

Cub is a terminal-based text editor designed for efficiency and simplicity. It offers two primary modes of operation: **View Mode** and **Insert Mode**. Users can navigate, edit, and manage text files with a set of intuitive keyboard shortcuts.

## Modes

### View Mode

View Mode is the default mode when you open a file. In this mode, you can navigate through the text without making any changes.

- **Key `i`**: Switch to Insert Mode.
- **Key `Esc`**: If in Insert Mode, switches back to View Mode.
- **Navigation**: Use the arrow keys to move around the text.

### Insert Mode

Insert Mode allows you to edit the text. Once in this mode, any keypresses will input text into the document.

- **Key `Esc`**: Return to View Mode.
- **Editing**: Any alphanumeric keys will insert text at the cursor location.

## Key Shortcuts

- **`Esc`**: Switch from Insert Mode to View Mode.
- **`Ctrl + Q`**: Quit the editor.
- **`Ctrl + U`**: Undo the last action.
- **`Ctrl + R`**: Redo the last undone action.
- **`i`**: Switch from View Mode to Insert Mode.
- **`Ctrl + S`**: Save the current file.

## Installation

To install and build Cub, follow these steps:

```bash
git clone https://github.com/yourusername/cub.git
cd cub
make build


### Usage 

