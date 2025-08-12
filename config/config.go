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
// Add to config/config.go

type Profile string

const (
	ProfileDev  Profile = "development"
	ProfileTest Profile = "testing"
	ProfileProd Profile = "production"
)

type Config struct {
	DataFile      string `json:"data_file"`
	BackupDir     string `json:"backup_dir"`
	AutoSave      bool   `json:"auto_save"`
	SaveInterval  int    `json:"save_interval"`
	MaxBackups    int    `json:"max_backups"`
	PrettyPrint   bool   `json:"pretty_print"`
	Profile       Profile `json:"profile"` // Add this
	LogLevel      string  `json:"log_level"`
	CacheSize     int     `json:"cache_size"`
}

func DefaultConfig() *Config {
	return &Config{
		DataFile:     "tasks.json",
		BackupDir:    "backups",
		AutoSave:     true,
		SaveInterval: 30,
		MaxBackups:   10,
		PrettyPrint:  true,
		Profile:      ProfileDev,
		LogLevel:     "info",
		CacheSize:    100,
	}
}

func LoadProfile(profile Profile) *Config {
	switch profile {
	case ProfileDev:
		return &Config{
			DataFile:     "tasks-dev.json",
			BackupDir:    "backups-dev",
			AutoSave:     true,
			SaveInterval: 10,
			MaxBackups:   5,
			PrettyPrint:  true,
			Profile:      ProfileDev,
			LogLevel:     "debug",
			CacheSize:    50,
		}
	case ProfileTest:
		return &Config{
			DataFile:     "tasks-test.json",
			BackupDir:    "backups-test",
			AutoSave:     false,
			SaveInterval: 0,
			MaxBackups:   2,
			PrettyPrint:  false,
			Profile:      ProfileTest,
			LogLevel:     "error",
			CacheSize:    10,
		}
	case ProfileProd:
		return &Config{
			DataFile:     "tasks-prod.json",
			BackupDir:    "backups-prod",
			AutoSave:     true,
			SaveInterval: 60,
			MaxBackups:   30,
			PrettyPrint:  false,
			Profile:      ProfileProd,
			LogLevel:     "warn",
			CacheSize:    1000,
		}
	default:
		return DefaultConfig()
	}
}

func (c *Config) Validate() []string {
	var errors []string
	
	if c.SaveInterval < 0 {
		errors = append(errors, "Save interval cannot be negative")
	}
	if c.MaxBackups < 1 {
		errors = append(errors, "Max backups must be at least 1")
	}
	if c.CacheSize < 1 {
		errors = append(errors, "Cache size must be at least 1")
	}
	
	return errors
}

func (c *Config) String() string {
	return fmt.Sprintf(
		"📋 Config Profile: %s\n"+
		"   Data File: %s\n"+
		"   Auto Save: %v (every %ds)\n"+
		"   Backups: %d\n"+
		"   Log Level: %s\n"+
		"   Cache Size: %d",
		c.Profile, c.DataFile, c.AutoSave, c.SaveInterval,
		c.MaxBackups, c.LogLevel, c.CacheSize)
}