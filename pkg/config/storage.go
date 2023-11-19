package config

import (
	"os"
	"path/filepath"
)

const configDir = ".kpfm/configs"

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Unable to find user home directory: " + err.Error())
	}

	configPath := filepath.Join(homeDir, configDir)
	if err := os.MkdirAll(configPath, 0755); err != nil {
		panic("Unable to create config directory: " + err.Error())
	}
}

func ReadConfig(alias string) ([]byte, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(homeDir, configDir, alias+".yaml")
	return os.ReadFile(filePath)
}

func WriteConfig(alias string, data []byte) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	filePath := filepath.Join(homeDir, configDir, alias+".yaml")
	return os.WriteFile(filePath, data, 0644)
}

func DeleteConfig(alias string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	filePath := filepath.Join(homeDir, configDir, alias+".yaml")
	return os.Remove(filePath)
}
