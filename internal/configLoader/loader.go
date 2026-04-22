package configloader

import (
	"fmt"
	"os"
	"reflect"

	"gopkg.in/yaml.v3"
)

type Config struct {
	SourceDir string `yaml:"sourceDir" cli:"sourceDir" cliDescription:"Directory to scan for files to process"`

	ConfigFilePath string

	DBPath string `yaml:"dbPath" cli:"dbPath" cliDescription:"Path to the database file"`

	FileExtensions []string `yaml:"fileExtensions" cli:"fileExtensions" cliDescription:"List of file extensions to process"`

	PrimaryCommandModule string `yaml:"primaryCommandModule" cli:"primaryCommandModule" cliDescription:"Primary command module to execute"`
}

func InitConfig() (*Config, error) {
	args := os.Args[1:]

	parsedArgs, err := argParse(args)

	if err != nil {
		return nil, fmt.Errorf("error parsing command line arguments: %w", err)
	}

	var configFilePath string
	for _, arg := range parsedArgs.Args {
		if arg.Key == "configFilePath" {
			configFilePath = arg.Value
			break
		}
	}

	var config Config
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

	for _, arg := range parsedArgs.Args {
		if arg.Key == "configFilePath" {
			continue
		}
		if err := applyArg(&config, arg); err != nil {
			return nil, err
		}
	}

	config.PrimaryCommandModule = parsedArgs.BaseModule

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

func applyArg(config *Config, args CommandLineArg) error {
	configValue := reflect.ValueOf(config).Elem()
	for i := 0; i < configValue.NumField(); i++ {
		field := configValue.Type().Field(i)
		cliTag := field.Tag.Get("cli")
		if cliTag == args.Key {
			fieldValue := configValue.Field(i)
			if fieldValue.Kind() == reflect.Slice {
				fieldValue.Set(reflect.Append(fieldValue, reflect.ValueOf(args.Value)))
			} else {
				fieldValue.SetString(args.Value)
			}
			return nil
		}
	}

	return fmt.Errorf("unknown argument: %s", args.Key)
}
