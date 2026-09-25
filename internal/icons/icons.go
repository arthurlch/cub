package icons

import (
	"os"
	"path/filepath"
	"strings"
)

var Nerd = os.Getenv("CUB_ASCII") == ""

func Powerline() string {
	if Nerd {
		return ""
	}
	return " "
}

func PowerlineLeft() string {
	if Nerd {
		return ""
	}
	return " "
}

const (
	FolderClosed = ""
	FolderOpen   = ""
	DefaultFile  = ""
	Branch       = ""
)

var byExt = map[string]string{
	"go":        "",
	"mod":       "",
	"sum":       "",
	"py":        "",
	"js":        "",
	"jsx":       "",
	"ts":        "",
	"tsx":       "",
	"json":      "",
	"md":        "",
	"yml":       "",
	"yaml":      "",
	"toml":      "",
	"sh":        "",
	"bash":      "",
	"html":      "",
	"css":       "",
	"rs":        "",
	"c":         "",
	"h":         "",
	"cpp":       "",
	"rb":        "",
	"java":      "",
	"lua":       "",
	"txt":       "",
	"lock":      "",
	"gitignore": "",
}

func File(name string) string {
	if !Nerd {
		return " "
	}
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(name)), ".")
	if g, ok := byExt[ext]; ok {
		return g
	}
	if strings.HasPrefix(name, ".git") {
		return ""
	}
	return DefaultFile
}

func Folder(open bool) string {
	if !Nerd {
		if open {
			return "▾"
		}
		return "▸"
	}
	if open {
		return FolderOpen
	}
	return FolderClosed
}
