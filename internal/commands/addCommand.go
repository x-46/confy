package commands

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"syscall"
	configloader "x46/confy/internal/configLoader"
	"x46/confy/internal/vault"

	"golang.org/x/term"
)

type AddCommand struct{}

func (c *AddCommand) GetName() string {
	return "add"
}

func (c *AddCommand) ValidateConfig(config *configloader.Config) error {
	if config.NewEntryName == "" && !slices.Contains(config.SetParameters, "entryName") {
		// read entry name from user input
		fmt.Print("Enter the name of the new entry (leave empty for random): ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			config.NewEntryName = scanner.Text()
		}
	}

	if config.NewEntryValue == "" {
		// read entry value from user input
		fmt.Print("Enter the value of the new entry: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))

		if err != nil {
			return fmt.Errorf("error reading entry value: %w", err)
		}
		fmt.Println()

		config.NewEntryValue = string(bytePassword)
	}

	if config.NewEntryName == "" {
		// random 10 character string
		newEntryName := make([]byte, 10)
		for i := range newEntryName {
			newEntryName[i] = byte(rand.Intn(26) + 97)
		}
		config.NewEntryName = string(newEntryName)
	}

	if config.NewEntryValue == "" {
		return fmt.Errorf("entry value cannot be empty")
	}

	if config.DBPath == "" {
		return fmt.Errorf("DBPath is required for add command")
	}

	if _, err := os.Stat(config.DBPath); os.IsNotExist(err) {
		return fmt.Errorf("database file does not exist at path: %s", config.DBPath)
	}

	if config.Password == "" {
		return fmt.Errorf("Password is required for add command")
	}

	return nil
}

func (c *AddCommand) Execute(config *configloader.Config) error {
	openVault, err := vault.OpenKeepassVault(config.DBPath, config.Password)
	if err != nil {
		return fmt.Errorf("failed to open vault: %w", err)
	}
	defer openVault.Close()

	if _, err := openVault.GetEntry(config.NewEntryName); err == nil {
		return fmt.Errorf("entry with name '%s' already exists", config.NewEntryName)
	}

	err = openVault.SetEntry(config.NewEntryName, config.NewEntryValue)
	if err != nil {
		return fmt.Errorf("failed to set entry: %w", err)
	}

	fmt.Printf("Entry '%s' added successfully.\n", config.NewEntryName)

	return nil
}

func (c *AddCommand) PrintShortHelp() {
	fmt.Println("  add       Add a new secret entry to the vault")
}

func (c *AddCommand) PrintLongHelp() {
	fmt.Println("Command: add")
	fmt.Println("Description: Add a new secret entry to the vault.")
	fmt.Println("Options:")
	printConfigOptions("dbPath", "password", "entryName", "entryValue")
	fmt.Println("  Note: --entryName defaults to a random name, --entryValue triggers an interactive prompt if omitted")
	fmt.Println("Examples:")
	fmt.Println("  confy add --entryName api_key --entryValue supersecret")
	fmt.Println("  confy add --dbPath confy.kdbx")
}
