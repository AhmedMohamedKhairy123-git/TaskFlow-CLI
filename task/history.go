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
// Add to task/history.go

type Command interface {
	Execute() error
	Undo() error
	Redo() error
}

type AddTaskCommand struct {
	store *TaskStore
	title string
	taskID int
}

func (c *AddTaskCommand) Execute() error {
	task, err := c.store.Add(c.title)
	if err == nil {
		c.taskID = task.ID
	}
	return err
}

func (c *AddTaskCommand) Undo() error {
	return c.store.Delete(c.taskID)
}

func (c *AddTaskCommand) Redo() error {
	_, err := c.store.Add(c.title)
	return err
}

type CompleteTaskCommand struct {
	store *TaskStore
	taskID int
	wasCompleted bool
}

func (c *CompleteTaskCommand) Execute() error {
	task, _ := c.store.Get(c.taskID)
	c.wasCompleted = task.Completed
	return c.store.MarkComplete(c.taskID)
}

func (c *CompleteTaskCommand) Undo() error {
	task, _ := c.store.Get(c.taskID)
	task.Completed = c.wasCompleted
	return nil
}

// Add UndoRedoManager
type UndoRedoManager struct {
	undoStack []Command
	redoStack []Command
	maxSize   int
}

func NewUndoRedoManager(maxSize int) *UndoRedoManager {
	return &UndoRedoManager{
		undoStack: make([]Command, 0, maxSize),
		redoStack: make([]Command, 0, maxSize),
		maxSize:   maxSize,
	}
}

func (u *UndoRedoManager) Execute(cmd Command) error {
	if err := cmd.Execute(); err != nil {
		return err
	}
	
	u.undoStack = append(u.undoStack, cmd)
	u.redoStack = make([]Command, 0) // Clear redo stack
	
	// Limit stack size
	if len(u.undoStack) > u.maxSize {
		u.undoStack = u.undoStack[1:]
	}
	
	return nil
}

func (u *UndoRedoManager) Undo() error {
	if len(u.undoStack) == 0 {
		return fmt.Errorf("nothing to undo")
	}
	
	cmd := u.undoStack[len(u.undoStack)-1]
	u.undoStack = u.undoStack[:len(u.undoStack)-1]
	
	if err := cmd.Undo(); err != nil {
		return err
	}
	
	u.redoStack = append(u.redoStack, cmd)
	return nil
}

func (u *UndoRedoManager) Redo() error {
	if len(u.redoStack) == 0 {
		return fmt.Errorf("nothing to redo")
	}
	
	cmd := u.redoStack[len(u.redoStack)-1]
	u.redoStack = u.redoStack[:len(u.redoStack)-1]
	
	if err := cmd.Redo(); err != nil {
		return err
	}
	
	u.undoStack = append(u.undoStack, cmd)
	return nil
}