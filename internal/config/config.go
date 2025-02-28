package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

// Config struct with json tags
type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// Method - Set a user
func (c *Config) SetUser(userName string) error {
	c.CurrentUserName = userName
	return write(*c)
}

// Function - Read the JSON config file and returns a Config struct
func Read() (Config, error) {
	// Get the config file path
	jsonFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	// Open the config file
	file, err := os.Open(jsonFilePath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	// Decode cfg file into results
	decoder := json.NewDecoder(file)
	cfg := Config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Helper - Get config path
func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	cfgFilePath := filepath.Join(homePath, configFileName)
	return cfgFilePath, nil
}

// Helper - Write to JSON
func write(cfg Config) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	// Create and open json cfg file
	file, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write to JSON file
	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}
