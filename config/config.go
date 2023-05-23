package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"task-tracker/errors"
)

type Config struct {
	DataFile      string `json:"data_file"`
	BackupDir     string `json:"backup_dir"`
	AutoSave      bool   `json:"auto_save"`
	SaveInterval  int    `json:"save_interval"` // in seconds
	MaxBackups    int    `json:"max_backups"`
	PrettyPrint   bool   `json:"pretty_print"`
}

func DefaultConfig() *Config {
	return &Config{
		DataFile:     "tasks.json",
		BackupDir:    "backups",
		AutoSave:     true,
		SaveInterval: 30,
		MaxBackups:   10,
		PrettyPrint:  true,
	}
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, errors.NewAppError(
			errors.ErrStorageFailure,
			"LoadConfig",
			"failed to load config",
		).WithCause(err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, errors.NewAppError(
			errors.ErrStorageFailure,
			"LoadConfig",
			"failed to parse config",
		).WithCause(err)
	}

	return &config, nil
}

func (c *Config) Save(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.NewAppError(
			errors.ErrStorageFailure,
			"SaveConfig",
			"failed to create config directory",
		).WithCause(err)
	}

	file, err := os.Create(path)
	if err != nil {
		return errors.NewAppError(
			errors.ErrStorageFailure,
			"SaveConfig",
			"failed to create config file",
		).WithCause(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if c.PrettyPrint {
		encoder.SetIndent("", "  ")
	}
	
	if err := encoder.Encode(c); err != nil {
		return errors.NewAppError(
			errors.ErrStorageFailure,
			"SaveConfig",
			"failed to write config",
		).WithCause(err)
	}

	return nil
}