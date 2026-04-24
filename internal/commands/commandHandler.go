package commands

import "fmt"

var commandList []CommandModule = []CommandModule{
	&HelpCommand{},
	&InitCommand{},
	&AddCommand{},
}

func GetCommandByName(name string) (CommandModule, error) {
	for _, cmd := range commandList {
		if cmd.GetName() == name {
			return cmd, nil
		}
	}
	return nil, fmt.Errorf("command not found: %s", name)
}

func GetValidCommands() []string {
	var validCommands []string
	for _, cmd := range commandList {
		validCommands = append(validCommands, cmd.GetName())
	}
	return validCommands
}
