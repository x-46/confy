package commands

import configloader "x46/confy/internal/configLoader"

type CommandModule interface {
	Execute(config *configloader.Config) error
	ValidateConfig(config *configloader.Config) error
	PrintHelp()
	GetName() string
}
