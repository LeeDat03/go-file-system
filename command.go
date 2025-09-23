package main

import "strings"

type IAction string

const (
	TOUCH IAction = "touch"
	MKDIR IAction = "mkdir"
	PWD   IAction = "pwd"
	CD    IAction = "cd"
	LS    IAction = "ls"
	EXIT  IAction = "exit"
)

type Command struct {
	Action IAction
	Args   []string
}

func ParseCommand(input string) Command {
	parts := strings.Fields(strings.TrimSpace(input))
	if len(parts) == 0 {
		return Command{}
	}

	return Command{
		Action: IAction(parts[0]),
		Args:   parts[1:],
	}
}
