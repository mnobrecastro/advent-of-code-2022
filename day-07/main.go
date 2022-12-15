// Miguel Nobre Castro
// https://adventofcode.com/2022/day/7

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Node struct
type Node struct {
	name     string
	dir      bool
	size     int
	parent   *Node
	children []*Node
}

// Node constructor.
func NewNode(name string, size int, parent *Node) (node *Node) {

	dir := false
	if size == 0 {
		dir = true
	}
	node = &Node{
		name:     name,
		dir:      dir,
		size:     size,
		parent:   parent,
		children: nil,
	}
	// Verbose only
	var par string
	if node.parent != nil {
		par = node.parent.name
	} else {
		par = "nil"
	}
	fmt.Printf("Created '%s' (size=%d) inside '%s'\n", node.name, node.size, par)
	return
}

// Tree struct.
type Tree struct {
	size int
	num  int
	root *Node
	ptr  *Node
}

// Tree constructor.
func NewTree() (tree *Tree) {

	tree = &Tree{
		size: 0,
		num:  0,
		root: nil,
		ptr:  nil,
	}
	return
}

// Initializes the root node 'name' in the Tree.
func (tree *Tree) MakeRoot(name string) error {

	if tree.root == nil {
		tree.root = NewNode(name, 0, nil)
		tree.ptr = tree.root
		return nil
	} else {
		return errors.New("Terminal: root already exists.")
	}
}

// Initializes a new node 'name' in the Tree.
func (tree *Tree) MakeNode(name string, size int) error {

	if tree.ptr != nil {
		for _, child := range tree.ptr.children {
			if child.name == name {
				if size == 0 {
					return errors.New("Terminal: directory already exists.")
				} else {
					return errors.New("Terminal: file already exists.")
				}
			}
		}
		tree.ptr.children = append(tree.ptr.children, NewNode(name, size, tree.ptr))
		if size > 0 {
			ptr := tree.ptr
			for ptr != nil {
				ptr.size += size
				ptr = ptr.parent
			}
		}
		return nil
	} else {
		return errors.New("Terminal: please MakeRoot() first.")
	}
}

// Moves the Tree ptr to a new directory.
func (tree *Tree) MoveIn(name string) error {

	if tree.ptr != nil {
		for _, child := range tree.ptr.children {
			if child.name == name && child.size == 0 {
				tree.ptr = child
				return nil
			}
		}
		return errors.New("Terminal: directory not found.")
	} else {
		return errors.New("Terminal: please MakeRoot() first.")
	}
}

// Moves the Tree ptr out of the current directory.
func (tree *Tree) MoveOut() error {

	if tree.ptr != tree.root {
		tree.ptr = tree.ptr.parent
		return nil
	} else {
		return errors.New("Terminal: already in root directory.")
	}
}

// Recursively checks directories with size less than or equal to 'threshold'.
func ListDirs(node *Node, threshold int) (val int) {

	val = 0
	if node.dir {
		for _, child := range node.children {
			val += ListDirs(child, threshold)
		}
		if node.size <= threshold {
			fmt.Printf("- %s (dir, size=%d)\n", node.name, node.size)
			val += node.size
		}
	}
	return
}

// Recursively checks the smallest directory below 'threshold'.
func ListSmallest(node *Node, threshold int, storage int, used int) (val int) {

	val = storage
	if node.dir {
		for _, child := range node.children {
			new_val := ListSmallest(child, threshold, storage, used)
			if new_val < val {
				val = new_val
			}
		}
		free := storage - (used - node.size)
		if free > threshold {
			fmt.Printf("- %s (dir, size=%d)\n", node.name, node.size)
			if node.size < val {
				val = node.size
			}
		}
	}
	return
}

// Terminal struct.
type Terminal struct {
	storage int
	tree    *Tree
	cmds    chan string
}

// Terminal constructor.
func NewTerminal(storage int, filename string) (t *Terminal) {

	cmds := make(chan string, 1)
	go func() {
		f, err := os.Open(filename)
		if err != nil {
			defer close(cmds)
		}

		r := bufio.NewReader(f)
		for true {
			s, err := r.ReadString('\n')
			if !errors.Is(err, io.EOF) {
				s = s[:len(s)-1]
				if len(s) > 0 {
					cmds <- s
				} else {
					continue
				}
			} else {
				defer close(cmds)
				break
			}
		}
		f.Close()
	}()

	tree := NewTree()
	t = &Terminal{
		storage: storage,
		tree:    tree,
		cmds:    cmds,
	}
	return
}

// Builds the Terminal's tree struct.
func (t *Terminal) Build() {

	var ls bool
	for cmd := range t.cmds {
		if cmd[0] == '$' {
			if cmd[2:4] == "cd" {
				name := cmd[5:]
				if name == "/" {
					t.tree.MakeRoot("/")
				} else if name == ".." {
					t.tree.MoveOut()
				} else if name != "" {
					t.tree.MoveIn(name)
				}
				ls = false
			} else if cmd[2:4] == "ls" {
				ls = true
			}
		} else if ls {
			if cmd[0:3] == "dir" {
				t.tree.MakeNode(cmd[4:], 0)
			} else {
				str := strings.Split(cmd, " ")
				size, _ := strconv.Atoi(str[0])
				t.tree.MakeNode(str[1], size)
			}
		}
	}
}

// Lists and returns the total size of directories with size less than or equal to 'threshold'.
func (t *Terminal) ListDirs(threshold int) (val int) {

	val = ListDirs(t.tree.root, threshold)
	return
}

// Lists the directory and its size which enables at least 'threshold' unused space.
func (t *Terminal) ListSmallest(threshold int) (val int) {

	val = ListSmallest(t.tree.root, threshold, t.storage, t.tree.root.size)
	return
}

func main() {

	const filename string = "input.txt"
	term := NewTerminal(70000000, filename)
	term.Build()
	println(term.ListDirs(100000))
	println(term.ListSmallest(30000000))
}
