package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type IAction string

const (
	ADD_FILE   IAction = "touch"
	ADD_FOLDER IAction = "mkdir"
	DISPLAY    IAction = "tree"
	EXIT       IAction = "exit"
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

func (t *Tree) AddNode(parentPath, name string, isFolder bool) bool {
	// pre-validate
	parentPath = strings.Trim(parentPath, " ")
	name = strings.Trim(name, " ")

	parentNode := t.FindNode(t.Root, parentPath)
	if parentNode == nil {
		fmt.Println("ERROR: path not found")
		return false
	}
	if !parentNode.IsFolder {
		fmt.Println("The path you type is not a folder")
		return false
	}

	newNode := &Node{
		Name:     name,
		IsFolder: isFolder,
		Child:    []*Node{},
		Parent:   t.Root,
	}
	parentNode.Child = append(parentNode.Child, newNode)
	return true
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
		action := IAction(userInputParts[0])

		switch action {
		case ADD_FILE:
			parentPath, fileName := userInputParts[1], userInputParts[2]

			if fs.AddNode(parentPath, fileName, false) {
				fmt.Printf("Added %s to %s \n", fileName, parentPath)
			}
		case ADD_FOLDER:
			parentPath, folderName := userInputParts[1], userInputParts[2]

			if fs.AddNode(parentPath, folderName, true) {
				fmt.Printf("Added %s to %s \n", folderName, parentPath)
			}
		case DISPLAY:
			fs.Display()
		case EXIT:
			fmt.Println("EXIT!!!")
			return
		default:
			fmt.Println("Unknow command")
		}
	}
}
