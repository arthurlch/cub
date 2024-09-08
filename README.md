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

---

### **Modes**
The text editor operates in two primary modes:
1. **View Mode**: Allows you to navigate through the text, select/copy/cut, and delete lines.
2. **Insert Mode**: Allows you to edit the text, inserting characters, spaces, tabs, and new lines.

You can switch between these modes using the appropriate keys.

---

### **View Mode**
View mode is used to navigate and perform line operations (e.g., delete, copy, cut). You can move around the document using navigation keys and perform basic text operations.

#### **Shortcuts in View Mode:**

- **Navigation**:
  - **Arrow Keys (`↑`, `↓`, `←`, `→`) or `o`, `p`, `k`, `l`**: Move the cursor up, down, left, or right.
  - **`Home`**: Move the cursor to the beginning of the current line.
  - **`End`**: Move the cursor to the end of the current line.
  - **`PgUp`**: Move the cursor up by a quarter of the screen.
  - **`PgDn`**: Move the cursor down by a quarter of the screen.

- **Selection and Clipboard**:
  - **`s`**: Start selecting text.
  - **`z`**: End text selection.
  - **`c`**: Copy the selected text to the clipboard.
  - **`x`**: Cut the selected text (remove it from the buffer and copy it to the clipboard).
  - **`v`**: Paste the text from the clipboard at the current cursor position.

- **Deleting Lines**:
  - **`d d`**: Press `d` twice to delete the current line.

#### **Undo/Redo in View Mode**:
- **Undo (`Ctrl+U`)**: Undo the last change.
- **Redo (`Ctrl+R`)**: Redo the last undone change.

---

### **Insert Mode**
Insert mode allows you to modify the document by adding characters, spaces, tabs, or new lines. You can enter Insert mode from View mode by pressing `i`.

#### **Shortcuts in Insert Mode:**

- **Insert Characters**: Type any character to insert it at the current cursor position.
- **Insert New Line** (`Enter`): Splits the current line at the cursor position and moves the part after the cursor to a new line.
- **Insert Space** (`Space`): Inserts a space at the current cursor position.
- **Insert Tab** (`Tab`): Inserts a tab (four spaces) at the current cursor position.
- **Delete Character** (`Backspace`): Deletes the character before the cursor. If the cursor is at the beginning of a line, it merges the current line with the previous one.

#### **Exiting Insert Mode**:
- **Esc**: Exit Insert mode and return to View mode.

---

### **Saving Files**
You can save the document at any time with the following shortcut:
- **Save (`Ctrl+S`)**: Saves the current document.

---

### **Exiting the Editor**
To exit the editor, use the following shortcut:
- **Exit (`Ctrl+Q`)**: Closes the editor.

---

### **Example Workflow**
1. **Open the editor**.
2. **Navigate** through the text using the arrow keys or navigation shortcuts (e.g., `k` to move left, `l` to move right).
3. **Switch to Insert Mode** by pressing `i`.
4. **Edit** the document by typing text, pressing `Enter` to add new lines, or using `Space` to insert spaces.
5. **Switch back to View Mode** by pressing `Esc`.
6. **Delete a line** by pressing `d` twice in quick succession (`dd`).
7. **Undo/Redo** changes using `Ctrl+U` and `Ctrl+R`.
8. **Save** the document using `Ctrl+S`.
9. **Exit** the editor using `Ctrl+Q`.

---

### **Common Issues and Tips**
- If you accidentally delete a line, you can undo it using `Ctrl+U`.
- If you want to quickly move to the beginning or end of a line, use the `Home` or `End` keys.
- If a large block of text needs to be deleted or moved, use the selection (`s` to start, `z` to end) and cut (`x`) or copy (`c`).

---

## Installation

To install and build Cub, follow these steps:

```bash
git clone https://github.com/yourusername/cub.git
cd cub
make build


### Usage 

