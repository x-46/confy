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
	c.PrintHelp()
	return nil
}

func (c *HelpCommand) PrintHelp() {
	fmt.Println("Usage: confy <command> [--arg value] [--flag]")

	config := &configloader.Config{}

	fmt.Println("\nArguments:")
	for field := range reflect.ValueOf(config).Elem().Type().Fields() {
		cliTag := field.Tag.Get("cli")
		cliDescription := field.Tag.Get("cliDescription")
		yamlTag := field.Tag.Get("yaml")

		if cliTag != "" && yamlTag != "" {
			fmt.Printf("  --%s: %s [%s]\n", cliTag, cliDescription, yamlTag)
		} else if cliTag != "" {
			fmt.Printf("  --%s: %s\n", cliTag, cliDescription)
		} else if yamlTag != "" {

			fmt.Printf("  %s: %s\n", yamlTag, cliDescription)
		}
	}
	fmt.Println("  --configFilePath: Path to a YAML config file")

	fmt.Println("\nCommands:")
	for _, cmd := range commandList {
		if cmd.GetName() == "help" {
			continue
		}
		cmd.PrintHelp()
	}

	fmt.Println("\nUse 'confy <command> --help' for more information on a specific command.")
}
