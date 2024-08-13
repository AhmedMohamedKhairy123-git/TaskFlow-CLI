package task

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type SharePermission string

const (
	PermView   SharePermission = "view"
	PermEdit   SharePermission = "edit"
	PermAdmin  SharePermission = "admin"
)

type SharedTask struct {
	TaskID      int
	SharedBy    string
	SharedWith  string
	Permission  SharePermission
	SharedAt    time.Time
	ExpiresAt   *time.Time
	ShareLink   string
}

type ShareStore struct {
	shares map[string]SharedTask
}

func NewShareStore() *ShareStore {
	return &ShareStore{
		shares: make(map[string]SharedTask),
	}
}

func generateShareLink() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (s *TaskStore) ShareTask(taskID int, sharedWith string, perm SharePermission, expiresIn *time.Duration) (string, error) {
	task, exists := s.Tasks[taskID]
	if !exists {
		return "", fmt.Errorf("task %d not found", taskID)
	}
	
	shareLink := generateShareLink()
	
	var expiresAt *time.Time
	if expiresIn != nil {
		t := time.Now().Add(*expiresIn)
		expiresAt = &t
	}
	
	shared := SharedTask{
		TaskID:     taskID,
		SharedBy:   "current-user",
		SharedWith: sharedWith,
		Permission: perm,
		SharedAt:   time.Now(),
		ExpiresAt:  expiresAt,
		ShareLink:  shareLink,
	}
	
	if s.shareStore == nil {
		s.shareStore = NewShareStore()
	}
	s.shareStore.shares[shareLink] = shared
	
	return shareLink, nil
}

func (s *TaskStore) GetSharedTask(shareLink string) (*Task, SharePermission, error) {
	if s.shareStore == nil {
		return nil, "", fmt.Errorf("no shares found")
	}
	
	shared, exists := s.shareStore.shares[shareLink]
	if !exists {
		return nil, "", fmt.Errorf("invalid share link")
	}
	
	if shared.ExpiresAt != nil && shared.ExpiresAt.Before(time.Now()) {
		delete(s.shareStore.shares, shareLink)
		return nil, "", fmt.Errorf("share link expired")
	}
	
	task, exists := s.Tasks[shared.TaskID]
	if !exists {
		return nil, "", fmt.Errorf("shared task no longer exists")
	}
	
	return task, shared.Permission, nil
}

func (s *TaskStore) RevokeShare(shareLink string) error {
	if s.shareStore == nil {
		return nil
	}
	delete(s.shareStore.shares, shareLink)
	return nil
}

func (s *TaskStore) ListSharedTasks() []SharedTask {
	if s.shareStore == nil {
		return nil
	}
	
	var list []SharedTask
	for _, share := range s.shareStore.shares {
		list = append(list, share)
	}
	return list
}