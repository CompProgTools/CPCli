package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/CompProgTools/CPCli/internal/models"
)

var (
	configPath string
	configFile string
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	configPath = filepath.Join(homeDir, ".cpcli")
	configFile = filepath.Join(configPath, "config.json")
}

func LoadConfig() (*models.Config, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.MkdirAll(configPath, 0755); err != nil {
			return nil, err
		}
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		config := &models.Config{}
		if err := SaveConfig(config); err != nil {
			return nil, err
		}
		return config, nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config models.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func SaveConfig(config *models.Config) error {
	data, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(configFile, data, 0644)
}

func SetAccount(platform string, handle string) error {
	config, err := LoadConfig()

	if err != nil {
		return err
	}

	switch platform {
	case "codeforces":
		config.Codeforces = handle
	case "leetcode":
		config.LeetCode = handle
	}

	return SaveConfig(config)
}

// checks whether all the account are linked
func IsAllLinked() (bool, error) {
	config, err := LoadConfig()

	if err != nil {
		return false, err
	}

	return config.Codeforces != "" && config.LeetCode != "", nil
}

func GetConfigPath() string {
	return configPath
}