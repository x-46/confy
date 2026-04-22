package configloader

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	SourceDir string `yaml:"sourceDir"`

	ConfigFilePath string

	DBPath string `yaml:"dbPath"`

	FileExtensions []string `yaml:"fileExtensions"`

	PrimaryCommandModule string `yaml:"primaryCommandModule"`
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
	switch args.Key {
	case "sourceDir":
		config.SourceDir = args.Value
	case "configFilePath":
		config.ConfigFilePath = args.Value
	case "dbPath":
		config.DBPath = args.Value
	case "fileExtensions":
		config.FileExtensions = append(config.FileExtensions, args.Value)
	default:
		return fmt.Errorf("unknown argument: %s", args.Key)
	}
	return nil
}
