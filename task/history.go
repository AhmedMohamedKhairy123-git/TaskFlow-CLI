package task

import (
	"encoding/json"
	"time"
)

type ActionType string

const (
	ActionCreated  ActionType = "created"
	ActionUpdated  ActionType = "updated"
	ActionDeleted  ActionType = "deleted"
	ActionCompleted ActionType = "completed"
	ActionPriority ActionType = "priority_changed"
	ActionTagAdded ActionType = "tag_added"
	ActionTagRemoved ActionType = "tag_removed"
)

type HistoryEntry struct {
	ID         string     `json:"id"`
	TaskID     int        `json:"task_id"`
	Action     ActionType `json:"action"`
	Timestamp  time.Time  `json:"timestamp"`
	Changes    []Change   `json:"changes,omitempty"`
	User       string     `json:"user,omitempty"`
}

type Change struct {
	Field string      `json:"field"`
	Old   interface{} `json:"old"`
	New   interface{} `json:"new"`
}

type HistoryStore struct {
	entries []HistoryEntry
}

func NewHistoryStore() *HistoryStore {
	return &HistoryStore{
		entries: make([]HistoryEntry, 0),
	}
}

func (h *HistoryStore) Log(entry HistoryEntry) {
	if entry.ID == "" {
		entry.ID = fmt.Sprintf("hist_%d_%d", entry.TaskID, time.Now().UnixNano())
	}
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}
	h.entries = append(h.entries, entry)
}

func (h *HistoryStore) GetTaskHistory(taskID int) []HistoryEntry {
	var taskHistory []HistoryEntry
	for _, entry := range h.entries {
		if entry.TaskID == taskID {
			taskHistory = append(taskHistory, entry)
		}
	}
	return taskHistory
}

func (h *HistoryStore) GetRecent(limit int) []HistoryEntry {
	if len(h.entries) <= limit {
		return h.entries
	}
	return h.entries[len(h.entries)-limit:]
}

func (h *HistoryStore) ExportToJSON() ([]byte, error) {
	return json.MarshalIndent(h.entries, "", "  ")
}

func (s *TaskStore) WithHistory(history *HistoryStore) *TaskStore {
	s.history = history
	return s
}

func (s *TaskStore) logHistory(taskID int, action ActionType, changes ...Change) {
	if s.history == nil {
		return
	}
	
	s.history.Log(HistoryEntry{
		TaskID:    taskID,
		Action:    action,
		Timestamp: time.Now(),
		Changes:   changes,
	})
}