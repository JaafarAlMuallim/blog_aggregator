package config

import (
	"encoding/json"
	"os"
)

const configFileName = "/.gatorconfig.json"

func Read() (*Config, error) {
	var config Config
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Config) SetUser(user string) error {
	c.User = user
	if err := write(*c); err != nil {
		return err
	}
	return nil
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + configFileName, nil

}

func write(cfg Config) error {
	homeDir, err := getConfigPath()
	if err != nil {
		return err
	}
	body, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(homeDir, body, 0644)
}
