package filetree

import (
	"os"
	"path/filepath"
	"sort"
)

type Node struct {
	Path     string
	Name     string
	IsDir    bool
	Expanded bool
	Depth    int
	Children []*Node
	loaded   bool
}

func New(root string) (*Node, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return nil, err
	}
	n := &Node{Path: abs, Name: filepath.Base(abs), IsDir: info.IsDir(), Expanded: true}
	n.load()
	return n, nil
}

func (n *Node) load() {
	if n.loaded || !n.IsDir {
		return
	}
	n.loaded = true

	entries, err := os.ReadDir(n.Path)
	if err != nil {
		return
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir()
		}
		return entries[i].Name() < entries[j].Name()
	})
	for _, e := range entries {
		n.Children = append(n.Children, &Node{
			Path:  filepath.Join(n.Path, e.Name()),
			Name:  e.Name(),
			IsDir: e.IsDir(),
			Depth: n.Depth + 1,
		})
	}
}

func (n *Node) Toggle() {
	if !n.IsDir {
		return
	}
	n.Expanded = !n.Expanded
	if n.Expanded {
		n.load()
	}
}

func Flatten(root *Node) []*Node {
	var out []*Node
	var walk func(n *Node)
	walk = func(n *Node) {
		for _, c := range n.Children {
			out = append(out, c)
			if c.IsDir && c.Expanded {
				walk(c)
			}
		}
	}
	if root != nil {
		walk(root)
	}
	return out
}
