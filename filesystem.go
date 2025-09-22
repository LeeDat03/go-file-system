package main

import (
	"fmt"
	"strings"
)

type Node struct {
	Name     string
	IsFolder bool
	Child    []*Node
	Parent   *Node
}

type FileSystem struct {
	Root *Node
}

func NewFileSystem(rootName string) *FileSystem {
	return &FileSystem{
		Root: &Node{
			Name:     rootName,
			IsFolder: true,
			Child:    []*Node{},
		},
	}
}

func (t *FileSystem) AddNode(parentPath, name string, isFolder bool) error {
	// pre-validate
	parentPath = strings.Trim(parentPath, " ")
	name = strings.Trim(name, " ")

	if parentPath == "" {
		return fmt.Errorf("path not found: %s", parentPath)
	}

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

func (t *FileSystem) FindNode(curr *Node, path string) *Node {
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

// TODO: display tree with absolute path
func (t *FileSystem) Display(paths ...string) error {
	fmt.Println(t.Root.Name)
	t.displayNode(t.Root, 1)

	return nil
}

func (t *FileSystem) displayNode(node *Node, level int) {
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

func splitPath(path string) (dir, file string) {
	i := strings.LastIndex(path, "/")
	if i == -1 {
		return "", path
	}

	return path[:i], path[i+1:]
}

func (t *FileSystem) Touch(paths ...string) error {
	for _, path := range paths {
		dir, file := splitPath(path)
		fmt.Printf("Path: %s, Dir: %s, File: %s\n", path, dir, file)
		if file == "" {
			return fmt.Errorf("invalid file name in path: %s", path)
		}

		if err := t.AddNode(dir, file, false); err != nil {
			return fmt.Errorf("failed to touch '%s': %w", path, err)
		}
	}
	return nil
}

func (t *FileSystem) Mkdir(paths ...string) error {
	for _, path := range paths {
		dir, folder := splitPath(path)
		if folder == "" {
			return fmt.Errorf("invalid folder name in path: %s", path)
		}
		if err := t.AddNode(dir, folder, true); err != nil {
			return fmt.Errorf("failed to mkdir '%s': %w", path, err)
		}
	}
	return nil
}
