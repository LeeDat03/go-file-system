package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type IAction string

const (
	TOUCH IAction = "touch"
	MKDIR IAction = "mkdir"
	CD    IAction = "cd"
	TREE  IAction = "tree"
	EXIT  IAction = "exit"
)

type Node struct {
	Name     string
	IsFolder bool
	Child    []*Node
	Parent   *Node
}

type Tree struct {
	Root *Node
}

func NewTree(rootName string) *Tree {
	return &Tree{
		Root: &Node{
			Name:     rootName,
			IsFolder: true,
			Child:    []*Node{},
		},
	}
}

func (t *Tree) AddNode(parentPath, name string, isFolder bool) error {
	// pre-validate
	parentPath = strings.Trim(parentPath, " ")
	name = strings.Trim(name, " ")

	parentNode := t.FindNode(t.Root, parentPath)

	if parentNode == nil {
		return fmt.Errorf("path not found: %s", parentPath)
	}
	if !parentNode.IsFolder {
		return fmt.Errorf("the path '%s' is not a folder", parentPath)
	}

	for _, child := range parentNode.Child {
		if child.Name == name {
			return fmt.Errorf("node '%s' already exists in '%s'", name, parentPath)
		}
	}

	newNode := &Node{
		Name:     name,
		IsFolder: isFolder,
		Child:    []*Node{},
		Parent:   t.Root,
	}
	parentNode.Child = append(parentNode.Child, newNode)

	return nil
}

func (t *Tree) FindNode(curr *Node, path string) *Node {
	if curr == nil {
		return nil
	}

	splitedPath := strings.Split(path, "/")
	if len(splitedPath) == 1 && splitedPath[0] == curr.Name {
		return curr
	}
	for _, child := range curr.Child {
		if child.Name == splitedPath[1] {
			if len(splitedPath) == 1 {
				return child
			}
			return t.FindNode(child, strings.Join(splitedPath[1:],
				"/"))
		}
	}
	return nil
}

func (t *Tree) Display() {
	fmt.Println(t.Root.Name)
	t.displayNode(t.Root, 1)
}

func (t *Tree) displayNode(node *Node, level int) {
	for _, child := range node.Child {
		prefix := strings.Repeat("  ", level)
		if child.IsFolder {
			fmt.Printf("%s📁 %s\n", prefix, child.Name)
			t.displayNode(child, level+1)
		} else {
			fmt.Printf("%s📄 %s\n", prefix, child.Name)
		}
	}
}

func SplitPath(path string) (dir, file string) {
	i := strings.LastIndex(path, "/")
	if i == -1 {
		return "", path
	}

	return path[:i], path[i+1:]
}

func (t *Tree) Touch(paths ...string) error {
	for _, path := range paths {
		dir, file := SplitPath(path)
		if file == "" {
			return fmt.Errorf("invalid file name in path: %s", path)
		}
		if err := t.AddNode(dir, file, false); err != nil {
			return fmt.Errorf("failed to touch '%s': %w", path, err)
		}
	}
	return nil
}

func (t *Tree) Mkdir(paths ...string) error {
	for _, path := range paths {
		dir, folder := SplitPath(path)
		if folder == "" {
			return fmt.Errorf("invalid folder name in path: %s", path)
		}
		if err := t.AddNode(dir, folder, true); err != nil {
			return fmt.Errorf("failed to mkdir '%s': %w", path, err)
		}
	}
	return nil
}

func main() {
	fmt.Println(`
	 _______  _______    _______  ___   ___      _______   
	|       ||       |  |       ||   | |   |    |       |  
	|    ___||   _   |  |    ___||   | |   |    |    ___|  
	|   | __ |  | |  |  |   |___ |   | |   |    |   |___   
	|   ||  ||  |_|  |  |    ___||   | |   |___ |    ___|  
	|   |_| ||       |  |   |    |   | |       ||   |___   
	|_______||_______|  |___|    |___| |_______||_______|  
	`)

	// File System
	fs := NewTree("root")
	scanner := bufio.NewScanner(os.Stdin)

	// REPL = READ - EVEL - PRINT - LOOP
	for {
		fmt.Print(" > ")
		scanner.Scan()

		userInput := scanner.Text()
		userInputParts := strings.Fields(userInput)

		if len(userInputParts) == 0 {
			continue
		}

		action := IAction(userInputParts[0])
		args := userInputParts[1:]

		switch action {
		// TODO: relative path vs absolute path
		case TOUCH:
			err := fs.Touch(args...)
			if err != nil {
				fmt.Println("Error: ", err)
			}
		case MKDIR:
			err := fs.Mkdir(args...)
			if err != nil {
				fmt.Println("Error: ", err)
			}
		case TREE:
			fs.Display()
		case EXIT:
			fmt.Println("EXIT!!!")
			return
		default:
			fmt.Println("Unknow command")
		}
	}
}
