package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"
	configloader "x46/confy/internal/configLoader"
	fileprocessing "x46/confy/internal/fileProcessing"
)

const preCommitScript = `
echo "[confy pre-commit hook]"
# assign stdin to keyboard
exec < /dev/tty
# call confy find and redirect stderr to /dev/null
go run . find --errOnFind 2> /dev/null
# check exit code and exit if confy found remaining secrets
exitcode=$?
if [ $exitcode == 1 ]; then
	echo "Can not commit: Secrets are exposed"
    exit 1
fi
`

type HookCommand struct{}

func (c *HookCommand) GetName() string {
	return "hook"
}

func (c *HookCommand) ValidateConfig(config *configloader.Config) error {
	if config.SourceDir == "" {
		return fmt.Errorf("SourceDir is required for hook command")
	}
	return nil
}

func (c *HookCommand) Execute(config *configloader.Config) error {
	_, err := os.Stat(".git")
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("Not in a git repository: %w", err)
	}
	_, err = os.Stat(".git/hooks")
	if err != nil {
		return fmt.Errorf("failed to add confy hook: %w", err)
	}

	_, err = os.Stat(".git/hooks/pre-commit")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// create file with correct header if no pre-commit hook exists
			err := fileprocessing.WriteFile(".git/hooks/pre-commit", "#!/bin/sh")
			if err != nil {
				return fmt.Errorf("failed to add confy hook: %w", err)
			}
		} else {
			return fmt.Errorf("failed to add confy hook: %w", err)
		}
	} else {
		// check whether confy is already hooked
		content, err := fileprocessing.ReadFile(".git/hooks/pre-commit")
		if err != nil {
			return fmt.Errorf("failed to add confy hook: %w", err)
		}
		if strings.Contains(content, preCommitScript) {
			return fmt.Errorf("confy is already hooked")
		}
	}
	// append preCommitScript to pre-commit file
	f, err := os.OpenFile(".git/hooks/pre-commit", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to add confy hook: %w", err)
	}
	_, err = f.WriteString(preCommitScript)
	if err != nil {
		return fmt.Errorf("failed to add confy hook: %w", err)
	}

	fmt.Println("Successfully added hook")
	return nil
}

func (c *HookCommand) PrintShortHelp() {
	fmt.Println("  hook      Add pre-commit hook to git which prevents commiting when secrets are exposed")
}

func (c *HookCommand) PrintLongHelp() {
	fmt.Println("Command: hook")
	fmt.Println("Description: Add pre-commit hook to git which prevents commiting when secrets are exposed")
	fmt.Println("Options:")
	printConfigOptions("sourceDir")
	fmt.Println("Examples:")
	fmt.Println("  confy hook --sourceDir")
}
