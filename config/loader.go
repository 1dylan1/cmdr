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

    systemPath := filepath.Join("/usr/share", label, configFileName)
    if _, err := os.Stat(systemPath); err != nil {
        return "", fmt.Errorf("config file not found at %s", systemPath)
    }

    return systemPath, nil
}


func LoadConfig(configFileArg string) (*Config, error) {
	configFilePath := configFileArg
	if configFilePath == "" {
		defaultConfigFilePath, err := getConfigPath()
		if err != nil {
			return nil, fmt.Errorf("could not determine config file path: %v", err)
		}
		configFilePath = defaultConfigFilePath
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
