package commands

import (
	"fmt"
	"x46/confy/internal/buildinfo"
	configloader "x46/confy/internal/configLoader"
)

type VersionCommand struct{}

func (c *VersionCommand) GetName() string {
	return "version"
}

func (c *VersionCommand) ValidateConfig(config *configloader.Config) error {
	return nil
}

func (c *VersionCommand) Execute(config *configloader.Config) error {
	fmt.Println(buildinfo.String())
	return nil
}

func (c *VersionCommand) PrintShortHelp() {
	fmt.Println("  version   Show build and release version information")
}

func (c *VersionCommand) PrintLongHelp() {
	fmt.Println("Command: version")
	fmt.Println("Description: Show the build version and build date.")
	fmt.Println("Usage: confy version")
	fmt.Println("Example:")
	fmt.Println("  confy version")
}
