# Cub Text Editor

<p align="center">
  <img src="https://i.ibb.co/88MWThZ/d9f9d3aa-2ae6-47d9-96bc-2735eda584f9.webp" alt="Cub Logo" width="300" height="300">
</p>

![GitHub release (latest by date)](https://img.shields.io/github/v/release/arthurlch/cub)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/arthurlch/cub/ci.yml)
![GitHub issues](https://img.shields.io/github/issues/arthurlch/cub)

## Overview

**Cub** is a **lightweight, BLAZINGLY-FAST, terminal-based text editor** built for speed, simplicity, and precision. Inspired by the best features of Kakoune and Vim, Cub offers **intuitive keyboard navigation** that strikes a balance between **powerful functionality** and **ease of use**. Whether you're editing configuration files, writing code, or working on documentation, Cub's minimalist design helps you stay productive without unnecessary distractions.

Cub operates with two streamlined modes:  
- **View Mode**: Navigate, select, and manipulate text with precision.  
- **Insert Mode**: Seamlessly edit and insert text where needed.  

With **advanced navigation tools** such as **line jumps, word motions, and bracket matching**, Cub offers a superior text editing experience compared to traditional editors like Nano. It’s designed to fit into your workflow, offering **faster navigation**, **fewer keystrokes**, and **better text management**.

Cub is perfect for developers, system administrators, and power users who need a **fast, no-frills editor** that maximizes efficiency without sacrificing simplicity.

---

### **Navigation and Modes**

The navigation system in Cub follows a **Kakoune-like** style with some elements inspired by Vim motions. This provides powerful yet simple navigation, aiming to be more intuitive than Nano without introducing unnecessary complexity. Navigation is separated into **simple** and **complex** actions to align with the two modes (View and Insert).

---

### **Navigation Keys (View Mode)**

- **Arrow Keys (`↑`, `↓`, `←`, `→`) or `o`, `p`, `k`, `l`**: Move the cursor in all directions.
- **`Home` / `End`**: Move to the start or end of the current line.
- **`PgUp` / `PgDn`**: Scroll up or down by a quarter of the screen.

---

### **Complex Navigation (View Mode)**  
These advanced motions allow for more efficient movement throughout the document:

- **Numbers (`0-9`)**: Build a line number to jump to.  
- **`G`**: Jump to the line specified by the accumulated number, or the end of the document if no number is provided.  
- **`w`**: Move to the **next word** boundary.  
- **`b`**: Move to the **previous word** boundary.  
- **`(` / `)`**: Move to the **matching bracket**.  
- **`e`**: Jump to the **next empty line**.  
- **`E`**: Jump to the **previous empty line**.  
- **`^`**: Move to the **first non-whitespace character** of the current line.  
- **`$`**: Move to the **end of the current line**.  
- **`z`**: **Center** the cursor on the screen.  
- **`g`**: Jump to the **top of the document** (line 0, column 0).  

---

### **Insert Mode Navigation**

While in **Insert Mode**, basic cursor navigation is still available:

- **Arrow Keys**: Move the cursor in any direction.
- **`Home` / `End`**: Jump to the beginning or end of the line.
- **`PgUp` / `PgDn`**: Scroll the view up or down by a quarter of the screen.

---

### **Jump to Line Functionality**

- **G + [number]**: Jump to the specified line. If no number is provided, it jumps to the **end of the document**.  
  _(Example: Pressing `4G` moves the cursor to line 4.)_

---

### **Bracket Matching**

Cub offers automatic navigation between matching brackets:

- **`(` / `)`**: Jump between matching parentheses.  
- **Supports both** round brackets (`()`), **and curly brackets** (`{}`).  

---

### **Page and Word Navigation**  

- **Page Up / Down**: Move the view by a quarter of the visible screen.  
- **`w` / `b`**: Jump forward to the next word boundary or backward to the previous one.  

---

### **Empty Line Detection**  
These shortcuts help jump between empty lines, improving navigation within long documents:

- **`e`**: Move to the next empty line.  
- **`E`**: Move to the previous empty line.  

---

### **Centering the Cursor**

- **`z`**: Centers the cursor vertically on the screen.  

---

### **Syntax Highlighting**

Cub now supports syntax highlighting for the following languages:

**A**: ABAP, ABNF, ActionScript, ActionScript 3, Ada, Agda, AL, Alloy, Angular2, ANTLR, ApacheConf, APL, AppleScript, ArangoDB AQL, Arduino, ArmAsm, AutoHotkey, AutoIt, Awk  
**B**: Ballerina, Bash, Bash Session, Batchfile, BibTeX, Bicep, BlitzBasic, BNF, BQN, Brainfuck  
**C**: C, C#, C++, Caddyfile, Caddyfile Directives, Cap'n Proto, Cassandra CQL, Ceylon, CFEngine3, cfstatement, ChaiScript, Chapel, Cheetah, Clojure, CMake, COBOL, CoffeeScript, Common Lisp, Coq, Crystal, CSS, Cython  
**D**: D, Dart, Dax, Desktop Entry, Diff, Django/Jinja, dns, Docker, DTD, Dylan  
**E**: EBNF, Elixir, Elm, EmacsLisp, Erlang  
**F**: Factor, Fennel, Fish, Forth, Fortran, FortranFixed, FSharp  
**G**: GAS, GDScript, Genshi, Genshi HTML, Genshi Text, Gherkin, Gleam, GLSL, Gnuplot, Go, Go HTML Template, Go Text Template, GraphQL, Groff, Groovy  
**H**: Handlebars, Hare, Haskell, Haxe, HCL, Hexdump, HLB, HLSL, HolyC, HTML, HTTP, Hy  
**I**: Idris, Igor, INI, Io, ISCdhcpd  
**J**: J, Java, JavaScript, JSON, Julia, Jungle  
**K**: Kotlin  
**L**: Lighttpd configuration file, LLVM, Lua  
**M**: Makefile, Mako, Markdown, Mason, Materialize SQL dialect, Mathematica, Matlab, MCFunction, Meson, Metal, MiniZinc, MLIR, Modula-2, MonkeyC, MorrowindScript, Myghty, MySQL  
**N**: NASM, Natural, Newspeak, Nginx configuration file, Nim, Nix  
**O**: Objective-C, OCaml, Octave, Odin, OnesEnterprise, OpenEdge ABL, OpenSCAD, Org Mode  
**P**: PacmanConf, Perl, PHP, PHTML, Pig, PkgConfig, PL/pgSQL, plaintext, Plutus Core, Pony, PostgreSQL SQL dialect, PostScript, POVRay, PowerQuery, PowerShell, Prolog, PromQL, Promela, properties, Protocol Buffer, PRQL, PSL, Puppet, Python, Python 2  
**Q**: QBasic, QML  
**R**: R, Racket, Ragel, Raku, React, ReasonML, reg, Rego, reStructuredText, Rexx, RPMSpec, Ruby, Rust  
**S**: SAS, Sass, Scala, Scheme, Scilab, SCSS, Sed, Sieve, Smali, Smalltalk, Smarty, SNBT, Snobol, Solidity, SourcePawn, SPARQL, SQL, SquidConf, Standard ML, Stas, Stylus, Svelte, Swift, SYSTEMD, SystemVerilog  
**T**: TableGen, Tal, TASM, Tcl, Tcsh, Termcap, Terminfo, Terraform, TeX, Thrift, TOML, TradingView, Transact-SQL, Turing, Turtle, Twig, TypeScript, TypoScript, TypoScriptCssData, TypoScriptHtmlData, Typst  
**V**: V, V Shell, Vala, VB.net, Verilog, VHDL, VHS, VimL, Vue  
**W**: WDTE, WebGPU Shading Language, Whiley  
**X**: XML, Xorg  
**Y**: YAML, YANG  
**Z**: Z80 Assembly, Zed, Zig

---

### **Example Workflow**

1. **Open Cub** and **navigate** using arrow keys or `o`, `p`, `k`, `l`.  
2. **Switch to Insert Mode** with `i` and type your text.  
3. **Save** the document with `Ctrl+S`.  
4. **Exit** the editor with `Ctrl+Q`.  

---

## Installation

To install and build Cub, follow these steps:

```bash
git clone https://github.com/yourusername/cub.git
cd cub
make build
```

---

## Usage

After installation, start the editor by running:

```bash
./cub <filename>
```

---

### **Common Issues and Tips**

- If you accidentally delete a line, **undo** it using `Ctrl+U`.  
- Use `Home` or `End` to jump to the beginning or end of a line.  
- Use **selection (`s` to start, `z` to end)** for bulk operations like cut or copy.

---
