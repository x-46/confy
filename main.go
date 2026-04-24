package main

import (
	"fmt"
	"x46/confy/internal/commands"
	configloader "x46/confy/internal/configLoader"
)

func main() {
	/*err := vault.NewKeepassVault("vault2.kdbx", "password")
	if err != nil {
		fmt.Println("Error creating KeePass vault:", err)
		return
	}

	res, err := vault.OpenKeepassVault("vault2.kdbx", "password")
	if err != nil {
		fmt.Println("Error opening KeePass vault:", err)
		return
	}
	fmt.Println("KeePass vault opened successfully:", res)

	res.SetEntry("test", "asd")

	err = res.Close()
	if err != nil {
		fmt.Println("Error closing KeePass vault:", err)
		return
	}
	*/
	config, err := configloader.InitConfig()
	if err != nil {
		fmt.Println("Error initializing config", err)
		fmt.Println("Run confy help for a list of available options.")
		return
	}

	fmt.Println("Config initialized:", config)

	commandModule, err := commands.GetCommandByName(config.PrimaryCommandModule)
	if err != nil {
		fmt.Println("Command not found:", config.PrimaryCommandModule)
		fmt.Println("Run confy help for a list of available commands.")
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
