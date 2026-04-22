package commands

var commandList []*MoadlCommand = []*MoadlCommand{}

func GetCommandByName(name string) *MoadlCommand {
	for _, cmd := range commandList {
		if (*cmd).GetName() == name {
			return cmd
		}
	}
	return nil
}

func GetValidCommands() []string {
	var validCommands []string
	for _, cmd := range commandList {
		validCommands = append(validCommands, (*cmd).GetName())
	}
	return validCommands
}
