package commands

import (
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	configloader "x46/confy/internal/configLoader"
	fileprocessing "x46/confy/internal/fileProcessing"
	"x46/confy/internal/vault"
)

type FindCommand struct{}

func (c *FindCommand) GetName() string {
	return "find"
}

func (c *FindCommand) ValidateConfig(config *configloader.Config) error {
	if config.DBPath == "" {
		return fmt.Errorf("DBPath is required for check command")
	}
	if config.SourceDir == "" {
		return fmt.Errorf("SourceDir is required for check command")
	}
	if config.Password == "" {
		return fmt.Errorf("Password is required for check command")
	}

	return nil
}

func (c *FindCommand) Execute(config *configloader.Config) error {
	openVault, err := vault.OpenKeepassVault(config.DBPath, config.Password)
	if err != nil {
		return fmt.Errorf("failed to open vault: %w", err)
	}
	defer openVault.Close()

	patterns := []string{}

	keys, err := openVault.GetKeys()
	if err != nil {
		return fmt.Errorf("failed to get keys from vault: %w", err)
	}

	for _, key := range keys {
		value, err := openVault.GetEntry(key)
		if err != nil {
			return fmt.Errorf("failed to get entry for key '%s': %w", key, err)
		}
		patterns = append(patterns, value)
	}
	numSecrets := atomic.Int64{}
	err = fileprocessing.RecursiveFileProcessing(config.SourceDir, config.FileExtensions, func(path string) error {
		content, err := fileprocessing.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file '%s': %w", path, err)
		}

		patternMatches := fileprocessing.MultipleIndex(content, patterns)
		if len(patternMatches) != 0 {
			numSecrets.Add(int64(len(patternMatches)))
			var b strings.Builder
			fmt.Fprintf(&b, "Secrets found in %s:\n", path)
			for _, match := range patternMatches {
				fmt.Fprintf(&b, "line %d column %d: %s\n", match.Line, match.Column, keys[match.PatternIndex])
			}
			fmt.Print(b.String())
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to index in files: %w", err)
	}

	fmt.Printf("Found %d remaining secrets\n", numSecrets.Load())
	if numSecrets.Load() != 0 && config.ErrOnFind {
		os.Exit(1)
	}
	return nil
}

func (c *FindCommand) PrintShortHelp() {
	fmt.Println("  find      Find for and report remaining secrets")
}

func (c *FindCommand) PrintLongHelp() {
	fmt.Println("Command: check")
	fmt.Println("Description: Find for and report remaining secrets.")
	fmt.Println("Options:")
	printConfigOptions("sourceDir", "dbPath", "password", "fileExtensions", "configFilePath", "errOnFind")
	fmt.Println("Examples:")
	fmt.Println("  confy find --sourceDir . --dbPath confy.kdbx --errOnFind")
	fmt.Println("  confy find --sourceDir . --fileExtensions .env --fileExtensions .yaml")
}
