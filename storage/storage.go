package storage

import (
	"task-tracker/task"
	"fmt" 
)

// Storage interface defines methods for task persistence
type Storage interface {
	Save(tasks []task.Task) error
	Load() ([]task.Task, error)
	Backup() (string, error)
	Restore(backupFile string) error
	ListBackups() ([]string, error)
}

// StorageError represents storage-specific errors
type StorageError struct {
	Op      string
	Path    string
	Message string
	Err     error
}

func (e *StorageError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("storage %s failed for %s: %s (%v)", 
			e.Op, e.Path, e.Message, e.Err)
	}
	return fmt.Sprintf("storage %s failed for %s: %s", e.Op, e.Path, e.Message)
}

// Import fmt for the StorageError