package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"  // Add this
	"time"
	"task-tracker/errors"
	"task-tracker/task"
)

type JSONStorage struct {
	filePath    string
	backupDir   string
	maxBackups  int
	autoBackup  bool
}

type StorageData struct {
	Version     string      `json:"version"`
	LastUpdated time.Time   `json:"last_updated"`
	Tasks       []task.Task `json:"tasks"`
	Metadata    struct {
		TotalTasks   int `json:"total_tasks"`
		Completed    int `json:"completed"`
		Pending      int `json:"pending"`
	} `json:"metadata"`
}

func NewJSONStorage(filePath, backupDir string) *JSONStorage {
	return &JSONStorage{
		filePath:   filePath,
		backupDir:  backupDir,
		maxBackups: 10,
		autoBackup: true,
	}
}

func (js *JSONStorage) Save(tasks []task.Task) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(js.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.NewAppError(
			errors.ErrStorageFailure,
			"JSONStorage.Save",
			"failed to create directory",
		).WithCause(err).WithContext("dir", dir)
	}

	// Calculate metadata
	completed := 0
	for _, t := range tasks {
		if t.Completed {
			completed++
		}
	}

	data := StorageData{
		Version:     "1.0.0",
		LastUpdated: time.Now(),
		Tasks:       tasks,
	}
	data.Metadata.TotalTasks = len(tasks)
	data.Metadata.Completed = completed
	data.Metadata.Pending = len(tasks) - completed

	// Marshal with proper formatting
	var jsonData []byte
	var err error
	
	if true { // pretty print
		jsonData, err = json.MarshalIndent(data, "", "  ")
	} else {
		jsonData, err = json.Marshal(data)
	}
	
	if err != nil {
		return errors.NewAppError(
			errors.ErrStorageFailure,
			"JSONStorage.Save",
			"failed to marshal tasks",
		).WithCause(err)
	}

	// Write to temp file first
	tempFile := js.filePath + ".tmp"
	if err := os.WriteFile(tempFile, jsonData, 0644); err != nil {
		return errors.NewAppError(
			errors.ErrStorageFailure,
			"JSONStorage.Save",
			"failed to write temp file",
		).WithCause(err)
	}

	// Atomic rename
	if err := os.Rename(tempFile, js.filePath); err != nil {
		os.Remove(tempFile) // Clean up temp file
		return errors.NewAppError(
			errors.ErrStorageFailure,
			"JSONStorage.Save",
			"failed to rename temp file",
		).WithCause(err)
	}

	// Create backup if enabled
	if js.autoBackup {
		go js.createBackup(tasks)
	}

	return nil
}

func (js *JSONStorage) Load() ([]task.Task, error) {
	// Check if file exists
	if _, err := os.Stat(js.filePath); os.IsNotExist(err) {
		return []task.Task{}, nil // Return empty slice for new file
	}

	// Read file
	data, err := os.ReadFile(js.filePath)
	if err != nil {
		// Try to recover from backup if main file is corrupted
		if recovered, recErr := js.recoverFromBackup(); recErr == nil {
			return recovered, nil
		}
		return nil, errors.NewAppError(
			errors.ErrStorageFailure,
			"JSONStorage.Load",
			"failed to read file",
		).WithCause(err)
	}

	// Parse JSON
	var storageData StorageData
	if err := json.Unmarshal(data, &storageData); err != nil {
		// Try older format (direct task array)
		var tasks []task.Task
		if json.Unmarshal(data, &tasks) == nil {
			return tasks, nil
		}
		return nil, errors.NewAppError(
			errors.ErrStorageFailure,
			"JSONStorage.Load",
			"failed to parse JSON",
		).WithCause(err)
	}

	return storageData.Tasks, nil
}

func (js *JSONStorage) Backup() (string, error) {
	// Load current tasks
	tasks, err := js.Load()
	if err != nil {
		return "", err
	}

	return js.createBackup(tasks)
}

func (js *JSONStorage) createBackup(tasks []task.Task) (string, error) {
	// Create backup directory
	if err := os.MkdirAll(js.backupDir, 0755); err != nil {
		return "", errors.NewAppError(
			errors.ErrStorageFailure,
			"JSONStorage.createBackup",
			"failed to create backup directory",
		).WithCause(err)
	}

	// Generate backup filename
	timestamp := time.Now().Format("20060102-150405")
	backupFile := filepath.Join(js.backupDir, fmt.Sprintf("tasks-backup-%s.json", timestamp))

	// Prepare backup data
	data := StorageData{
		Version:     "1.0.0",
		LastUpdated: time.Now(),
		Tasks:       tasks,
	}

	// Marshal
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", errors.NewAppError(
			errors.ErrStorageFailure,
			"JSONStorage.createBackup",
			"failed to marshal backup data",
		).WithCause(err)
	}

	// Write backup
	if err := os.WriteFile(backupFile, jsonData, 0644); err != nil {
		return "", errors.NewAppError(
			errors.ErrStorageFailure,
			"JSONStorage.createBackup",
			"failed to write backup file",
		).WithCause(err)
	}

	// Clean old backups
	go js.cleanOldBackups()

	return backupFile, nil
}

func (js *JSONStorage) Restore(backupFile string) error {
	// Read backup file
	data, err := os.ReadFile(backupFile)
	if err != nil {
		return errors.NewAppError(
			errors.ErrStorageFailure,
			"JSONStorage.Restore",
			"failed to read backup file",
		).WithCause(err)
	}

	// Parse backup
	var storageData StorageData
	if err := json.Unmarshal(data, &storageData); err != nil {
		return errors.NewAppError(
			errors.ErrStorageFailure,
			"JSONStorage.Restore",
			"failed to parse backup file",
		).WithCause(err)
	}

	// Save to main file
	return js.Save(storageData.Tasks)
}

func (js *JSONStorage) recoverFromBackup() ([]task.Task, error) {
	backups, err := js.ListBackups()
	if err != nil || len(backups) == 0 {
		return nil, fmt.Errorf("no backups available for recovery")
	}

	// Try most recent backup
	latest := backups[len(backups)-1]
	if err := js.Restore(latest); err != nil {
		return nil, err
	}

	// Return restored tasks
	return js.Load()
}

func (js *JSONStorage) ListBackups() ([]string, error) {
	files, err := os.ReadDir(js.backupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var backups []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" && 
		   strings.HasPrefix(file.Name(), "tasks-backup-") {
			backups = append(backups, filepath.Join(js.backupDir, file.Name()))
		}
	}

	return backups, nil
}

func (js *JSONStorage) cleanOldBackups() {
	backups, err := js.ListBackups()
	if err != nil {
		return
	}

	if len(backups) <= js.maxBackups {
		return
	}

	// Remove oldest backups
	for i := 0; i < len(backups)-js.maxBackups; i++ {
		os.Remove(backups[i])
	}
}

// Import strings for ListBackups