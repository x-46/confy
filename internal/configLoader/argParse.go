package configloader

import (
	"fmt"
	"strings"
)

type CommandLineArgs struct {
	BaseModule string

	Args []CommandLineArg
}

type CommandLineArg struct {
	Key   string
	Value string
}

func argParse(args []string) (CommandLineArgs, error) {
	if len(args) == 0 {
		return CommandLineArgs{}, fmt.Errorf("no command provided")
	}

	baseModule := args[0]

	allArgs := []CommandLineArg{}

	var lastArg string = ""
	for i := 1; i < len(args); i++ {
		if strings.HasPrefix(args[i], "--") {
			if lastArg != "" {
				allArgs = append(allArgs, CommandLineArg{Key: lastArg, Value: ""})
			}
			lastArg = strings.TrimPrefix(args[i], "--")
		} else if lastArg != "" {
			allArgs = append(allArgs, CommandLineArg{Key: lastArg, Value: args[i]})
			lastArg = ""
		} else {
			return CommandLineArgs{}, fmt.Errorf("unexpected argument format: %s", args[i])
		}
	}

	return CommandLineArgs{BaseModule: baseModule, Args: allArgs}, nil
}
