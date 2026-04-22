package commands

import configloader "x46/confy/internal/configLoader"

type MoadlCommand interface {
	Execute(config *configloader.Config) error
	ValidateConfig(config *configloader.Config) error
	GetName() string
}
