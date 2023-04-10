package task

// TaskManager defines the interface for task operations
type TaskManager interface {
	Add(title string) (Task, error)
	Get(id int) (Task, error)
	GetAll() []Task
	Update(id int, task Task) error
	Delete(id int) error
	FindByTitle(title string) []Task  // Add this line
}

// Completer defines interface for completion operations
type Completer interface {
	MarkComplete(id int) error
	GetCompleted() []Task
	GetPending() []Task
}

// Prioritizer defines interface for priority operations
type Prioritizer interface {
	SetPriority(id int, priority Priority) error
	GetByPriority(priority Priority) []Task
	GetHighestPriority() []Task
}

// Tagger defines interface for tag operations
type Tagger interface {
	AddTag(id int, tag string) error
	RemoveTag(id int, tag string) error
	GetByTag(tag string) []Task
}

// StatsReporter defines interface for statistics
type StatsReporter interface {
	GetStats() map[string]interface{}
	GetCompletionRate() float64
	DisplayStats()  // Add this line
}

// FullTaskManager combines all interfaces
type FullTaskManager interface {
	TaskManager
	Completer
	Prioritizer
	Tagger
	StatsReporter
}