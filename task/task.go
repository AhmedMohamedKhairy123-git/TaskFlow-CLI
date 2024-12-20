package task

import (
	"fmt"
	"runtime"  // Add this
	"sync"     // Add this
	"time"
)
type Priority int

const (
	Low Priority = iota
	Medium
	High
	Critical
)

func (p Priority) String() string {
	switch p {
	case Low:
		return "LOW"
	case Medium:
		return "MEDIUM"
	case High:
		return "HIGH"
	case Critical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

func (p Priority) Color() string {
	switch p {
	case Low:
		return "\033[32m" // Green
	case Medium:
		return "\033[33m" // Yellow
	case High:
		return "\033[31m" // Red
	case Critical:
		return "\033[35m" // Magenta
	default:
		return "\033[0m"
	}
}

type Task struct {
	ID           int
	Title        string
	Completed    bool
	Priority     Priority
	CreatedAt    time.Time
	Tags         []string
	Reminder     *Reminder
	Dependencies []Dependency
	Notes        []Note        `json:"notes,omitempty"`       // Add this
	Attachments  []Attachment  `json:"attachments,omitempty"` // Add this
}

func NewTask(id int, title string) *Task {
	return &Task{
		ID:        id,
		Title:     title,
		Completed: false,
		Priority:  Low,
		CreatedAt: time.Now(),
		Tags:      make([]string, 0),
	}
}

func (t *Task) MarkComplete() {
	t.Completed = true
}

func (t *Task) SetPriority(p Priority) {
	t.Priority = p
}

func (t *Task) AddTag(tag string) {
	for _, existing := range t.Tags {
		if existing == tag {
			return
		}
	}
	t.Tags = append(t.Tags, tag)
}

func (t *Task) RemoveTag(tag string) {
	for i, existing := range t.Tags {
		if existing == tag {
			t.Tags = append(t.Tags[:i], t.Tags[i+1:]...)
			return
		}
	}
}

func (t *Task) Display() string {
	status := " "
	if t.Completed {
		status = "✓"
	}
	
	reset := "\033[0m"
	color := t.Priority.Color()
	
	return fmt.Sprintf("%s[%s]%s %d: %s (Priority: %s, Tags: %v, Created: %s)",
		color, status, reset, t.ID, t.Title, t.Priority, t.Tags, t.CreatedAt.Format("2006-01-02"))
}

func (t *Task) DisplaySimple() string {
	status := " "
	if t.Completed {
		status = "✓"
	}
	return fmt.Sprintf("[%s] %d: %s", status, t.ID, t.Title)
}
// Add to task/task.go

// Use sync.Pool for temporary objects
var taskPool = sync.Pool{
	New: func() interface{} {
		return &Task{
			Tags:         make([]string, 0, 5),
			Dependencies: make([]Dependency, 0, 2),
			Notes:        make([]Note, 0, 3),
		}
	},
}

// Optimized Task creation with pooling
func NewTaskOptimized(id int, title string) *Task {
	task := taskPool.Get().(*Task)
	task.ID = id
	task.Title = title
	task.Completed = false
	task.Priority = Low
	task.CreatedAt = time.Now()
	task.Tags = task.Tags[:0]
	task.Dependencies = task.Dependencies[:0]
	task.Notes = task.Notes[:0]
	task.Attachments = task.Attachments[:0]
	task.Reminder = nil
	return task
}

// Return task to pool
func ReleaseTask(task *Task) {
	task.Title = ""
	task.Tags = nil
	task.Dependencies = nil
	task.Notes = nil
	task.Attachments = nil
	task.Reminder = nil
	taskPool.Put(task)
}

// Add memory stats
type MemoryStats struct {
	HeapAlloc    uint64
	HeapObjects  uint64
	NumGC        uint32
	LastGC       time.Time
}

func GetMemoryStats() MemoryStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	return MemoryStats{
		HeapAlloc:   m.HeapAlloc / 1024 / 1024, // MB
		HeapObjects: m.HeapObjects,
		NumGC:       m.NumGC,
		LastGC:      time.Unix(0, int64(m.LastGC)),
	}
}

// Add to TaskStore
func (s *TaskStore) OptimizeMemory() {
	// Force GC if needed
	if len(s.Tasks) > 10000 {
		runtime.GC()
	}
	
	// Clear expired cache
	s.cache.cleanupLoop()
	
	// Compact maps if needed
	if len(s.Tasks) < cap(s.Tasks) {
		newMap := make(map[int]*Task, len(s.Tasks))
		for k, v := range s.Tasks {
			newMap[k] = v
		}
		s.Tasks = newMap
	}
}