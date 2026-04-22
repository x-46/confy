package main

import (
	"fmt"
	"x46/confy/internal/commands"
	configloader "x46/confy/internal/configLoader"
)

func main() {
	fmt.Println("Startring confy...")

	config, err := configloader.InitConfig()
	if err != nil {
		fmt.Println("Error initializing config:", err)
		return
	}

	fmt.Println("Config initialized:", config)

	commandModule, err := commands.GetCommandByName(config.PrimaryCommandModule)
	if err != nil {
		fmt.Println("Error getting command module:", config.PrimaryCommandModule)
		return
	}

	if config.HelpOnly {
		commandModule.PrintHelp()
		return
	}

	err = commandModule.ValidateConfig(config)
	if err != nil {
		fmt.Println("Configuration validation failed:", err)
		return
	}

	err = commandModule.Execute(config)
	if err != nil {
		fmt.Println("Error executing command module:", err)
		return
	}
}
