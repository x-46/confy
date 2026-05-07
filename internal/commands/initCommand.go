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

func (c *InitCommand) PrintShortHelp() {
	fmt.Println("  init      Initialize a new KeePass database and optional config file")
}

func (c *InitCommand) PrintLongHelp() {
	fmt.Println("Command: init")
	fmt.Println("Description: Initialize a new KeePass database and optionally create a config file.")
	fmt.Println("Options:")
	printConfigOptions("dbPath", "password", "configFilePath", "sourceDir", "fileExtensions")
	fmt.Println("Examples:")
	fmt.Println("  confy init --dbPath confy.kdbx")
	fmt.Println("  confy init --dbPath vault.kdbx --configFilePath confy.yaml --sourceDir .")
}
