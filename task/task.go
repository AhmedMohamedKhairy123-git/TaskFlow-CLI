package task

import (
	"fmt"
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
	ID        int
	Title     string
	Completed bool
	Priority  Priority
	CreatedAt time.Time
	Tags      []string
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