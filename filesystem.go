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
	Root       *Node
	CurrentDir *Node
}

func NewFileSystem(rootName string) *FileSystem {
	root := &Node{
		Name:     rootName,
		IsFolder: true,
	}

	return &FileSystem{
		Root:       root,
		CurrentDir: root,
	}
}

func (fs *FileSystem) Display(paths ...string) error {
	start := fs.CurrentDir
	fmt.Println(fs.CurrentPath())
	fs.displayNode(start, 1)
	return nil
}

func (fs *FileSystem) displayNode(node *Node, level int) {
	for _, child := range node.Child {
		prefix := strings.Repeat("  ", level)
		if child.IsFolder {
			fmt.Printf("%s📁 %s\n", prefix, child.Name)
			fs.displayNode(child, level+1)
		} else {
			fmt.Printf("%s📄 %s\n", prefix, child.Name)
		}
	}
}

func (fs *FileSystem) CurrentPath() string {
	parts := []string{}
	node := fs.CurrentDir
	for node != nil {
		parts = append([]string{node.Name}, parts...)
		node = node.Parent
	}
	return strings.Join(parts, "/")
}

func (fs *FileSystem) resolvePath(path string) []string {
	path = strings.TrimSpace(path)

	// absolute path
	if IsAbsolutePath(path) {
		parts := strings.Split(path, "/")
		if parts[0] == "" {
			parts = parts[1:]
		}
		return parts
	}

	// relative path
	currPath := fs.CurrentPath()
	if currPath == "" {
		currPath = "root"
	}

	parts := append(strings.Split(currPath, "/"), strings.Split(path, "/")...)
	return parts
}

func (fs *FileSystem) makePath(path string, makeDir bool) error {
	parts := strings.Split(path, "/")

	var node *Node

	if IsAbsolutePath(path) {
		node = fs.Root
		if len(parts) > 0 && parts[0] == fs.Root.Name {
			parts = parts[1:]
		}
	} else {
		node = fs.CurrentDir
	}

	for i, part := range parts {
		isLast := i == len(parts)-1
		var child *Node
		for _, c := range node.Child {
			if part == c.Name && c.IsFolder {
				child = c
				break
			}
		}
		if child != nil {
			if i == len(parts)-1 && !child.IsFolder {
				return fmt.Errorf("touch: file already exists: %s", path)
			}
			node = child
			continue
		}
		// create new node
		newNode := &Node{
			Name:     part,
			IsFolder: !isLast || makeDir,
			Parent:   node,
		}
		node.Child = append(node.Child, newNode)
		node = newNode
	}
	return nil
}

func (fs *FileSystem) Touch(paths ...string) error {
	for _, path := range paths {
		if err := fs.makePath(path, false); err != nil {
			return err
		}
	}
	return nil
}

func (fs *FileSystem) Mkdir(paths ...string) error {
	for _, path := range paths {
		if err := fs.makePath(path, true); err != nil {
			return err
		}
	}
	return nil
}

func (fs *FileSystem) Pwd(paths ...string) error {
	fmt.Println(fs.CurrentPath())
	return nil
}

func (fs *FileSystem) Cd(paths ...string) error {
	path := paths[0]
	parts := strings.Split(paths[0], "/")

	var node *Node

	if IsAbsolutePath(path) {
		node = fs.Root
		if len(parts) > 0 && parts[0] == fs.Root.Name {
			parts = parts[1:]
		}
	} else {
		node = fs.CurrentDir
	}

	for _, part := range parts {
		switch part {
		case "", ".":
			continue
		case "..":
			if node.Parent != nil {
				node = node.Parent

			}
		default:
			found := false
			for _, child := range node.Child {
				if child.Name == part && child.IsFolder {
					node = child
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("cd: no such directory: %s", path)
			}
		}
	}
	fs.CurrentDir = node
	return nil
}
