package task

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"task-tracker/errors"
	"time"
)

type BackupManager struct {
	store        *TaskStore
	backupDir    string
	mu           sync.RWMutex
	lastBackup   time.Time
	backupCount  int
	autoBackup   bool
	backupChan   chan struct{}
	errorHandler func(error)
}

type BackupMetadata struct {
	ID          string    `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
	TaskCount   int       `json:"task_count"`
	File        string    `json:"file"`
	Size        int64     `json:"size"`
	Checksum    string    `json:"checksum"`
}

func NewBackupManager(store *TaskStore, backupDir string) *BackupManager {
	bm := &BackupManager{
		store:        store,
		backupDir:    backupDir,
		backupChan:   make(chan struct{}, 1),
		autoBackup:   true,
		errorHandler: defaultErrorHandler,
	}
	
	// Create backup directory if it doesn't exist
	os.MkdirAll(backupDir, 0755)
	
	// Start auto-backup goroutine
	go bm.autoBackupRoutine()
	
	return bm
}

func defaultErrorHandler(err error) {
	fmt.Printf("Backup error: %v\n", err)
}

func (bm *BackupManager) SetErrorHandler(handler func(error)) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	bm.errorHandler = handler
}

func (bm *BackupManager) CreateBackup() (*BackupMetadata, error) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	
	defer func() {
		if r := recover(); r != nil {
			bm.errorHandler(fmt.Errorf("panic during backup: %v", r))
		}
	}()
	
	// Generate backup filename
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("backup-%s.json", timestamp)
	backupPath := filepath.Join(bm.backupDir, filename)
	
	// Get all tasks
	tasks := bm.store.GetAll()
	
	// Create backup data
	data := struct {
		Timestamp time.Time `json:"timestamp"`
		Tasks     []Task    `json:"tasks"`
		Version   string    `json:"version"`
	}{
		Timestamp: time.Now(),
		Tasks:     tasks,
		Version:   "1.0.0",
	}
	
	// Marshal with indentation
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, errors.NewAppError(
			errors.ErrBackupFailed,
			"CreateBackup",
			"failed to marshal backup data",
		).WithCause(err)
	}
	
	// Write to file
	err = os.WriteFile(backupPath, jsonData, 0644)
	if err != nil {
		return nil, errors.NewAppError(
			errors.ErrBackupFailed,
			"CreateBackup",
			"failed to write backup file",
		).WithCause(err)
	}
	
	// Get file info
	info, err := os.Stat(backupPath)
	if err != nil {
		return nil, err
	}
	
	// Create metadata
	metadata := &BackupMetadata{
		ID:        fmt.Sprintf("backup-%d", time.Now().UnixNano()),
		Timestamp: time.Now(),
		TaskCount: len(tasks),
		File:      filename,
		Size:      info.Size(),
	}
	
	bm.lastBackup = metadata.Timestamp
	bm.backupCount++
	
	// Clean old backups (keep last 10)
	go bm.cleanOldBackups(10)
	
	return metadata, nil
}

func (bm *BackupManager) RestoreLatest() error {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	
	backups, err := bm.ListBackups()
	if err != nil {
		return err
	}
	
	if len(backups) == 0 {
		return errors.NewAppError(
			errors.ErrRecoveryFailed,
			"RestoreLatest",
			"no backups found",
		)
	}
	
	// Get latest backup
	latest := backups[0]
	return bm.RestoreBackup(latest.File)
}

func (bm *BackupManager) RestoreBackup(filename string) error {
	backupPath := filepath.Join(bm.backupDir, filename)
	
	// Read backup file
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return errors.NewAppError(
			errors.ErrRecoveryFailed,
			"RestoreBackup",
			"failed to read backup file",
		).WithCause(err)
	}
	
	// Parse backup data
	var backupData struct {
		Tasks []Task `json:"tasks"`
	}
	
	err = json.Unmarshal(data, &backupData)
	if err != nil {
		return errors.NewAppError(
			errors.ErrRecoveryFailed,
			"RestoreBackup",
			"failed to parse backup data",
		).WithCause(err)
	}
	
	// Clear current store and restore
	bm.store.Tasks = make(map[int]*Task)
	bm.store.NextID = 1
	
	for _, t := range backupData.Tasks {
		taskCopy := t
		bm.store.Tasks[t.ID] = &taskCopy
		if t.ID >= bm.store.NextID {
			bm.store.NextID = t.ID + 1
		}
	}
	
	return nil
}

func (bm *BackupManager) ListBackups() ([]BackupMetadata, error) {
	files, err := os.ReadDir(bm.backupDir)
	if err != nil {
		return nil, err
	}
	
	var backups []BackupMetadata
	
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}
		
		info, err := file.Info()
		if err != nil {
			continue
		}
		
		backups = append(backups, BackupMetadata{
			File:      file.Name(),
			Timestamp: info.ModTime(),
			Size:      info.Size(),
		})
	}
	
	// Sort by timestamp descending (newest first)
	for i := 0; i < len(backups)-1; i++ {
		for j := i + 1; j < len(backups); j++ {
			if backups[i].Timestamp.Before(backups[j].Timestamp) {
				backups[i], backups[j] = backups[j], backups[i]
			}
		}
	}
	
	return backups, nil
}

func (bm *BackupManager) cleanOldBackups(keep int) {
	backups, err := bm.ListBackups()
	if err != nil {
		bm.errorHandler(err)
		return
	}
	
	if len(backups) <= keep {
		return
	}
	
	for i := keep; i < len(backups); i++ {
		path := filepath.Join(bm.backupDir, backups[i].File)
		os.Remove(path)
	}
}

func (bm *BackupManager) autoBackupRoutine() {
	defer func() {
		if r := recover(); r != nil {
			bm.errorHandler(fmt.Errorf("auto-backup panic: %v", r))
			go bm.autoBackupRoutine() // Restart
		}
	}()
	
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		if !bm.autoBackup {
			continue
		}
		
		_, err := bm.CreateBackup()
		if err != nil {
			bm.errorHandler(err)
		}
	}
}

func (bm *BackupManager) EnableAutoBackup() {
	bm.autoBackup = true
}

func (bm *BackupManager) DisableAutoBackup() {
	bm.autoBackup = false
}