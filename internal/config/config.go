package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = "jyggalag.json"

type Config struct {
	NotesDir string `json:"notes_dir"`
	Editor   string `json:"editor"`
}

func SetEditor(editor string) error {
	return setConfig(func(c *Config) {
		c.Editor = editor
	})
}

func SetNotesDir(path string) error {
	return setConfig(func(c *Config) {
		c.NotesDir = path
	})
}

func setConfig(updateFunc func(*Config)) error {
	initConfig()
	c, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("Could not load config %w", err)
	}
	updateFunc(c)
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
	filename := filepath.Join(dir, configFileName)

	// don't re-init the config if it already exists
	_, err = os.Stat(filename)
	if err == nil {
		return nil
	}

	// make sure the config file always exists
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}

	c, err := json.Marshal(Config{
		NotesDir: "/home/notes",
		Editor:   "vim",
	})
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(dir, configFileName), c, 0644)
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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".config")

	return filepath.Join(configDir, "jyggalag"), nil
}
