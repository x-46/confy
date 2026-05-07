package commands

import (
	"fmt"
	configloader "x46/confy/internal/configLoader"
	fileprocessing "x46/confy/internal/fileProcessing"
	"x46/confy/internal/vault"
)

type CheckoutCommand struct{}

func (c *CheckoutCommand) GetName() string {
	return "checkout"
}

func (c *CheckoutCommand) ValidateConfig(config *configloader.Config) error {
	if config.DBPath == "" {
		return fmt.Errorf("DBPath is required for checkout command")
	}
	if config.SourceDir == "" {
		return fmt.Errorf("SourceDir is required for checkout command")
	}
	if config.Password == "" {
		return fmt.Errorf("Password is required for checkout command")
	}

	return nil
}

func (c *CheckoutCommand) Execute(config *configloader.Config) error {
	openVault, err := vault.OpenKeepassVault(config.DBPath, config.Password)
	if err != nil {
		return fmt.Errorf("failed to open vault: %w", err)
	}
	defer openVault.Close()

	patterns := []fileprocessing.ReplacementPattern{}

	keys, err := openVault.GetKeys()
	if err != nil {
		return fmt.Errorf("failed to get keys from vault: %w", err)
	}

	for _, key := range keys {
		value, err := openVault.GetEntry(key)
		if err != nil {
			return fmt.Errorf("failed to get entry for key '%s': %w", key, err)
		}
		pattern := fileprocessing.ReplacementPattern{
			Pattern:     config.CreateReplacementKeyFromKey(key),
			Replacement: value,
		}
		patterns = append(patterns, pattern)
	}

	err = fileprocessing.RecursiveFileProcessing(config.SourceDir, config.FileExtensions, func(path string) error {
		content, err := fileprocessing.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file '%s': %w", path, err)
		}

		newContent, err := fileprocessing.MultiReplaceAll(content, patterns)
		if err != nil {
			return fmt.Errorf("failed to replace content in file '%s': %w", path, err)
		}

		err = fileprocessing.WriteFile(path, newContent)
		if err != nil {
			return fmt.Errorf("failed to write file '%s': %w", path, err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to process files: %w", err)
	}

	return nil
}

func (c *CheckoutCommand) PrintShortHelp() {
	fmt.Println("  checkout  Replace placeholders in files with secret values")
}

func (c *CheckoutCommand) PrintLongHelp() {
	fmt.Println("Command: checkout")
	fmt.Println("Description: Replace vault placeholders in files with secret values.")
	fmt.Println("Options:")
	printConfigOptions("sourceDir", "dbPath", "password", "fileExtensions", "configFilePath")
	fmt.Println("Examples:")
	fmt.Println("  confy checkout --sourceDir . --dbPath confy.kdbx")
	fmt.Println("  confy checkout --sourceDir . --fileExtensions .env --fileExtensions .yaml")
}
