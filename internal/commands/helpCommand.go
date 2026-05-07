package commands

import (
	"fmt"
	"reflect"
	configloader "x46/confy/internal/configLoader"
)

type HelpCommand struct{}

func (c *HelpCommand) GetName() string {
	return "help"
}

func (c *HelpCommand) ValidateConfig(config *configloader.Config) error {
	return nil
}

func (c *HelpCommand) Execute(config *configloader.Config) error {
	c.PrintLongHelp()
	return nil
}

func (c *HelpCommand) PrintShortHelp() {
	fmt.Println("  help      Show help overview or detailed command help")
}

// printConfigOptions prints options from Config filtered by the given cli tag names.
// Pass no names to print all options.
func printConfigOptions(names ...string) {
	filter := make(map[string]bool, len(names))
	for _, n := range names {
		filter[n] = true
	}
	config := &configloader.Config{}
	for field := range reflect.ValueOf(config).Elem().Type().Fields() {
		cliTag := field.Tag.Get("cli")
		cliDescription := field.Tag.Get("cliDescription")
		if cliTag == "" {
			continue
		}
		if len(filter) == 0 || filter[cliTag] {
			fmt.Printf("  --%s: %s\n", cliTag, cliDescription)
		}
	}
}

func (c *HelpCommand) PrintLongHelp() {
	fmt.Println("Command: confy")
	fmt.Println("Description: Manage secrets in files using a KeePass vault.")
	fmt.Println("Usage: confy <command> [options]")
	fmt.Println("\nGlobal options:")
	printConfigOptions()
	fmt.Println("\nCommands:")
	for _, cmd := range commandList {
		if cmd.GetName() == "help" {
			continue
		}
		cmd.PrintShortHelp()
	}
	fmt.Println("\nUse 'confy <command> --help' for detailed usage and examples.")
}
