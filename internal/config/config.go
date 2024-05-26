package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

const configFileName = "jyggalag.json"

type Config struct {
	NotesDir string `json:"notes_dir"`
}

func SetNotesDir(path string) error {
	initConfig()
	c, err := LoadConfig()
	if err != nil {
		return err
	}
	c.NotesDir = path
	return saveConfig(c)
}

func saveConfig(config *Config) error {
	dir, err := getConfigFileDir()
	filename := filepath.Join(dir, configFileName)
	if err != nil {
		return err
	}

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func initConfig() error {
	dir, err := getConfigFileDir()
	if err != nil {
		return err
	}
	// make sure the config file always exists
	err = os.MkdirAll(dir, 0777)
	emptyJson, err := json.Marshal(Config{})
	if err != nil {
		return err
	}
	os.WriteFile(filepath.Join(dir, configFileName), emptyJson, 0644)
	if err != nil {
		return err
	}
	return nil
}

func LoadConfig() (*Config, error) {
	dir, err := getConfigFileDir()
	if err != nil {
		return nil, err
	}

	filename := filepath.Join(dir, configFileName)

	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func getConfigFileDir() (string, error) {
	var configDir string
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "windows":
		configDir = os.Getenv("APPDATA")
	case "darwin":
		configDir = filepath.Join(homeDir, "Library", "Application Support")
	default:
		configDir = filepath.Join(homeDir, ".config")
	}

	return filepath.Join(configDir, "jyggalag"), nil
}
