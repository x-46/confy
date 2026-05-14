package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"
	configloader "x46/confy/internal/configLoader"
	fileprocessing "x46/confy/internal/fileProcessing"
)

type UnhookCommand struct{}

func (c *UnhookCommand) GetName() string {
	return "unhook"
}

func (c *UnhookCommand) ValidateConfig(config *configloader.Config) error {
	if config.SourceDir == "" {
		return fmt.Errorf("SourceDir is required for unhook command")
	}
	return nil
}

func (c *UnhookCommand) Execute(config *configloader.Config) error {
	_, err := os.Stat(".git")
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("Not in a git repository: %w", err)
	}
	_, err = os.Stat(".git/hooks/pre-commit")
	if err != nil {
		return fmt.Errorf("failed to locate pre-commit hooks: %w", err)
	}

	_, err = os.Stat(".git/hooks/pre-commit")
	content, err := fileprocessing.ReadFile(".git/hooks/pre-commit")
	if err != nil {
		return fmt.Errorf("failed to remove confy hook: %w", err)
	}
	if !strings.Contains(content, preCommitScript) {
		return fmt.Errorf("failed to remove confy hook: no hook to remove")
	}
	content = strings.ReplaceAll(content, preCommitScript, "")
	err = fileprocessing.WriteFile(".git/hooks/pre-commit", content)
	if err != nil {
		return fmt.Errorf("failed to remove confy hook: %w", err)
	}

	fmt.Println("Successfully removed hook")
	return nil
}

func (c *UnhookCommand) PrintShortHelp() {
	fmt.Println("  unhook    Remove pre-commit hook")
}

func (c *UnhookCommand) PrintLongHelp() {
	fmt.Println("Command: unhook")
	fmt.Println("Description: Remove pre-commit hook")
	fmt.Println("Options:")
	printConfigOptions("sourceDir")
	fmt.Println("Examples:")
	fmt.Println("  confy unhook --sourceDir")
}
