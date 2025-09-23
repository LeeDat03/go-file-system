package main

import (
	"bufio"
	"fmt"
	"os"
)

type CommandHandler func(args ...string) error

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
	fs := NewFileSystem("root")
	scanner := bufio.NewScanner(os.Stdin)

	// REPL = READ - EVEL - PRINT - LOOP
	for {
		fmt.Printf("[%s] > ", fs.CurrentPath())
		scanner.Scan()

		cmd := ParseCommand(scanner.Text())
		if cmd.Action == "" {
			continue
		}

		handlers := map[IAction]CommandHandler{
			TOUCH: fs.Touch,
			MKDIR: fs.Mkdir,
			LS:    fs.Display,
			CD:    fs.Cd,
			PWD:   fs.Pwd,
		}

		if handler, ok := handlers[cmd.Action]; ok {
			if err := handler(cmd.Args...); err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command:", cmd.Action)
		}
	}
}
