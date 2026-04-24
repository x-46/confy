package configloader

import (
	"fmt"
	"os"
	"reflect"
	"syscall"

	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

type Config struct {
	SourceDir string `yaml:"sourceDir" cli:"sourceDir" cliDescription:"Directory to scan for files to process"`

	ConfigFilePath string `cli:"configFilePath" cliDescription:"Path to the configuration file"`

	DBPath string `yaml:"dbPath" cli:"dbPath" cliDescription:"Path to the database file"`

	FileExtensions []string `yaml:"fileExtensions" cli:"fileExtensions" cliDescription:"List of file extensions to process"`

	PrimaryCommandModule string

	Password string `cli:"password" cliDescription:"Password for encrypting/decrypting data"`

	HelpOnly bool `cli:"help" cliDescription:"If set, only the help command will be executed"`

	NewEntryName string `cli:"entryName" cliDescription:"Name of the new entry to add"`

	NewEntryValue string `cli:"entryValue" cliDescription:"Value of the new entry to add"`

	SetParameters []string
}

func InitConfig() (*Config, error) {
	args := os.Args[1:]

	parsedArgs, err := argParse(args)

	if err != nil {
		return nil, fmt.Errorf("error parsing command line arguments: %w", err)
	}

	var configFilePath string
	var configFilePathSet bool = false
	for _, arg := range parsedArgs.Args {
		if arg.Key == "configFilePath" {
			configFilePath = arg.Value
			configFilePathSet = true
			break
		}
	}

	var config Config

	// if no config file path is provided, check if a default config file exists in the current directory
	if !configFilePathSet {
		if _, err := os.Stat("confy.yaml"); err == nil {
			configFilePath = "confy.yaml"
		}
	}

	// if a config file path is provided, attempt to load the config from the file
	if configFilePath != "" {
		loadedConfig, err := loadConfigFromFile(configFilePath)
		if err != nil {
			return nil, err
		}
		if loadedConfig != nil {
			config = *loadedConfig
		}
	} else {
		config = Config{}
	}

	// applys all command line arguments to the config struct, overriding any values from the config file
	for _, arg := range parsedArgs.Args {
		config.SetParameters = append(config.SetParameters, arg.Key)
		if err := applyArg(&config, arg); err != nil {
			return nil, err
		}
	}

	config.PrimaryCommandModule = parsedArgs.BaseModule

	if config.Password == "" && !(config.HelpOnly || config.PrimaryCommandModule == "help") {
		fmt.Print("Enter Password: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return nil, fmt.Errorf("error reading password: %w", err)
		}
		config.Password = string(bytePassword)
		fmt.Println()
	}

	return &config, nil
}

func loadConfigFromFile(filePath string) (*Config, error) {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	var c Config
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling config file: %w", err)
	}

	return &c, nil
}

func SaveConfigToFile(config *Config) error {
	if config.ConfigFilePath == "" {
		return fmt.Errorf("config file path is not set")
	}

	// Create a map to store only fields with yaml tags
	configMap := make(map[string]interface{})
	configValue := reflect.ValueOf(config).Elem()

	for i := 0; i < configValue.NumField(); i++ {
		field := configValue.Type().Field(i)
		yamlTag := field.Tag.Get("yaml")
		if yamlTag != "" && yamlTag != "-" {
			configMap[yamlTag] = configValue.Field(i).Interface()
		}
	}

	yamlData, err := yaml.Marshal(configMap)
	if err != nil {
		return fmt.Errorf("error marshalling config to yaml: %w", err)
	}

	err = os.WriteFile(config.ConfigFilePath, yamlData, 0644)
	if err != nil {
		return fmt.Errorf("error writing config to file: %w", err)
	}

	return nil
}

func applyArg(config *Config, args CommandLineArg) error {
	configValue := reflect.ValueOf(config).Elem()
	for i := 0; i < configValue.NumField(); i++ {
		field := configValue.Type().Field(i)
		cliTag := field.Tag.Get("cli")
		if cliTag == args.Key {
			fieldValue := configValue.Field(i)
			if fieldValue.Kind() == reflect.Slice {
				fieldValue.Set(reflect.Append(fieldValue, reflect.ValueOf(args.Value)))
			} else if fieldValue.Kind() == reflect.Bool {
				if args.Value == "false" {
					fieldValue.SetBool(false)
				} else {
					fieldValue.SetBool(true)
				}
			} else {
				fieldValue.SetString(args.Value)
			}
			return nil
		}
	}

	return fmt.Errorf("unknown argument: %s", args.Key)
}
