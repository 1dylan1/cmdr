package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Config struct {
	Servers map[string]Server `yaml:",inline"`
}

type Server struct {
	Address string `yaml:"address"`
	Password string `yaml:"password"`
}

func getConfigPath() (string, error) {
	const label = "cmdr"
	const configFileName = "config.yaml"

	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("could not get user home directory: %v", err)
		}
		configHome = filepath.Join(homeDir, ".config")

	}
	return filepath.Join(configHome, label, configFileName), nil
}

func LoadConfig() (*Config, error) {
	configFilePath, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("could not determine config file path: %v", err)
	}
	
	fileContent, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file from '%s': %v", configFilePath, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(fileContent, &cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config file '%s': %v", configFilePath, err)
	}
	
	return &cfg, nil
}
