// editor/settings.go
package editor

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Settings struct {
	TabSize int `json:"tab_size"`
}

func LoadSettings() (*Settings, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".gocodeeditor", "settings.json")

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultSettings := &Settings{
			TabSize: 4,
		}
		return defaultSettings, SaveSettings(defaultSettings)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, err
	}

	return &settings, nil
}

func SaveSettings(settings *Settings) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, ".gocodeeditor", "settings.json")

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
