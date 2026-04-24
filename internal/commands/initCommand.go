package commands

import (
	"fmt"
	"slices"
	configloader "x46/confy/internal/configLoader"
	"x46/confy/internal/vault"
)

type InitCommand struct{}

func (c *InitCommand) GetName() string {
	return "init"
}

func (c *InitCommand) ValidateConfig(config *configloader.Config) error {
	if !slices.Contains(config.SetParameters, "configFilePath") {
		config.ConfigFilePath = "confy.yaml"
	}

	if !slices.Contains(config.SetParameters, "fileExtensions") {
		config.FileExtensions = []string{".md", ".yml", ".yaml", ".json"}
	}

	if !slices.Contains(config.SetParameters, "dbPath") {
		config.DBPath = "confy.kdbx"
	}

	if !slices.Contains(config.SetParameters, "sourceDir") {
		config.SourceDir = "."
	}

	if config.DBPath == "" {
		return fmt.Errorf("DBPath is required for init command")
	}
	if config.Password == "" {
		return fmt.Errorf("Password is required for init command")
	}
	return nil
}

func (c *InitCommand) Execute(config *configloader.Config) error {
	err := vault.NewKeepassVault(config.DBPath, config.Password)
	if err != nil {
		return err
	}

	if config.ConfigFilePath != "" {
		err = configloader.SaveConfigToFile(config)
		if err != nil {
			return fmt.Errorf("error saving config to file: %w", err)
		}
	}

	fmt.Println("Keepass database initialized successfully.")

	return nil
}

func (c *InitCommand) PrintHelp() {
	fmt.Println("  init: Initializes a new Keepass Database.")
	fmt.Println("        If an configuration path is provided but no configuration file exists, it will be created with the provided values.")
}
