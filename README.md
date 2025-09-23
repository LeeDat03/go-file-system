"Learn by doing": practice building data structures (tree), recursion, and simple CLI parsing in Go.

# Go File System

A tiny interactive file system simulator written in Go. It models a hierarchical file system using a tree data structure and exposes a simple REPL with basic commands (mkdir, touch, ls).

## Features

- `mkdir <path>`: create a directory; parent path must already exist; duplicate names under the same folder are rejected
- `touch <path>`: create a file; parent path must already exist; file name must be non-empty
- `ls`: print the entire tree from the root
- `cd`: move between directories

## Getting Started

### Prerequisites
- Go 1.21+ installed

### Run
```bash
go run .
```
You should see the banner and a prompt like:
```
 >
```

