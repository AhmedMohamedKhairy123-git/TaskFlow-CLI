package task

import (
	"errors"
	"fmt"
	"strings"  
)

type TaskStore struct {
	Tasks   map[int]*Task
	NextID  int
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		Tasks:  make(map[int]*Task),
		NextID: 1,
	}
}

// TaskManager implementation
func (s *TaskStore) Add(title string) (Task, error) {
	if title == "" {
		return Task{}, errors.New("task title cannot be empty")
	}
	
	task := NewTask(s.NextID, title)
	s.Tasks[task.ID] = task
	s.NextID++
	return *task, nil
}

func (s *TaskStore) Get(id int) (Task, error) {
	task, exists := s.Tasks[id]
	if !exists {
		return Task{}, fmt.Errorf("task with ID %d not found", id)
	}
	return *task, nil
}

func (s *TaskStore) GetAll() []Task {
	var taskList []Task
	for _, task := range s.Tasks {
		taskList = append(taskList, *task)
	}
	return taskList
}

func (s *TaskStore) Update(id int, updatedTask Task) error {
	task, exists := s.Tasks[id]
	if !exists {
		return fmt.Errorf("task with ID %d not found", id)
	}
	
	if updatedTask.Title != "" {
		task.Title = updatedTask.Title
	}
	task.Completed = updatedTask.Completed
	task.Priority = updatedTask.Priority
	task.Tags = updatedTask.Tags
	
	return nil
}

func (s *TaskStore) Delete(id int) error {
	_, exists := s.Tasks[id]
	if !exists {
		return fmt.Errorf("task with ID %d not found", id)
	}
	delete(s.Tasks, id)
	return nil
}
// FindByTitle implements TaskManager interface
func (s *TaskStore) FindByTitle(title string) []Task {
	var results []Task
	titleLower := strings.ToLower(title)
	
	for _, task := range s.Tasks {
		if strings.Contains(strings.ToLower(task.Title), titleLower) {
			results = append(results, *task)
		}
	}
	return results
}
// Completer implementation
func (s *TaskStore) MarkComplete(id int) error {
	task, exists := s.Tasks[id]
	if !exists {
		return fmt.Errorf("task with ID %d not found", id)
	}
	task.MarkComplete()
	return nil
}

func (s *TaskStore) GetCompleted() []Task {
	var completed []Task
	for _, task := range s.Tasks {
		if task.Completed {
			completed = append(completed, *task)
		}
	}
	return completed
}

func (s *TaskStore) GetPending() []Task {
	var pending []Task
	for _, task := range s.Tasks {
		if !task.Completed {
			pending = append(pending, *task)
		}
	}
	return pending
}

// Prioritizer implementation
func (s *TaskStore) SetPriority(id int, priority Priority) error {
	task, exists := s.Tasks[id]
	if !exists {
		return fmt.Errorf("task with ID %d not found", id)
	}
	task.SetPriority(priority)
	return nil
}

func (s *TaskStore) GetByPriority(priority Priority) []Task {
	var tasks []Task
	for _, task := range s.Tasks {
		if task.Priority == priority {
			tasks = append(tasks, *task)
		}
	}
	return tasks
}

func (s *TaskStore) GetHighestPriority() []Task {
	if len(s.Tasks) == 0 {
		return []Task{}
	}
	
	highest := Low
	for _, task := range s.Tasks {
		if task.Priority > highest {
			highest = task.Priority
		}
	}
	
	return s.GetByPriority(highest)
}

// Tagger implementation
func (s *TaskStore) AddTag(id int, tag string) error {
	task, exists := s.Tasks[id]
	if !exists {
		return fmt.Errorf("task with ID %d not found", id)
	}
	task.AddTag(tag)
	return nil
}

func (s *TaskStore) RemoveTag(id int, tag string) error {
	task, exists := s.Tasks[id]
	if !exists {
		return fmt.Errorf("task with ID %d not found", id)
	}
	task.RemoveTag(tag)
	return nil
}

func (s *TaskStore) GetByTag(tag string) []Task {
	var tagged []Task
	for _, task := range s.Tasks {
		for _, t := range task.Tags {
			if t == tag {
				tagged = append(tagged, *task)
				break
			}
		}
	}
	return tagged
}

// StatsReporter implementation
func (s *TaskStore) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})
	total := len(s.Tasks)
	completed := len(s.GetCompleted())
	
	priorityCount := make(map[Priority]int)
	for _, task := range s.Tasks {
		priorityCount[task.Priority]++
	}
	
	stats["total"] = total
	stats["completed"] = completed
	stats["pending"] = total - completed
	stats["by_priority"] = priorityCount
	
	return stats
}

func (s *TaskStore) GetCompletionRate() float64 {
	if len(s.Tasks) == 0 {
		return 0
	}
	return float64(len(s.GetCompleted())) / float64(len(s.Tasks)) * 100
}

func (s *TaskStore) DisplayStats() {
	stats := s.GetStats()
	rate := s.GetCompletionRate()
	
	fmt.Println("\n📊 --- STATISTICS ---")
	fmt.Printf("Total Tasks: %d\n", stats["total"])
	fmt.Printf("Completed: %d\n", stats["completed"])
	fmt.Printf("Pending: %d\n", stats["pending"])
	fmt.Printf("Completion Rate: %.1f%%\n", rate)
	
	fmt.Println("\nBy Priority:")
	for p, count := range stats["by_priority"].(map[Priority]int) {
		fmt.Printf("  %s: %d\n", p, count)
	}
	fmt.Println("-------------------")
}
// Add to task/store.go

type LazyLoader struct {
	loaded    bool
	loadFn    func() (interface{}, error)
	value     interface{}
	mu        sync.RWMutex
}

func NewLazyLoader(loadFn func() (interface{}, error)) *LazyLoader {
	return &LazyLoader{
		loadFn: loadFn,
	}
}

func (l *LazyLoader) Get() (interface{}, error) {
	l.mu.RLock()
	if l.loaded {
		defer l.mu.RUnlock()
		return l.value, nil
	}
	l.mu.RUnlock()
	
	l.mu.Lock()
	defer l.mu.Unlock()
	
	if l.loaded {
		return l.value, nil
	}
	
	value, err := l.loadFn()
	if err != nil {
		return nil, err
	}
	
	l.value = value
	l.loaded = true
	return value, nil
}

// Add lazy loading to TaskStore
type TaskStore struct {
	Tasks      map[int]*Task
	NextID     int
	history    *HistoryStore
	shareStore *ShareStore
	cache      *TaskCache
	lazyStats  *LazyLoader  // Add this
}

func (s *TaskStore) GetAnalyticsLazy() Analytics {
	if s.lazyStats == nil {
		s.lazyStats = NewLazyLoader(func() (interface{}, error) {
			return s.GetAnalytics(), nil
		})
	}
	
	result, _ := s.lazyStats.Get()
	return result.(Analytics)
}

// Add pagination with lazy loading
func (s *TaskStore) GetTasksPage(offset, limit int) []Task {
	if offset >= len(s.Tasks) {
		return []Task{}
	}
	
	var page []Task
	count := 0
	for _, task := range s.Tasks {
		if count >= offset && count < offset+limit {
			page = append(page, *task)
		}
		count++
		if count >= offset+limit {
			break
		}
	}
	
	return page
}

// Add lazy loading for attachments
func (s *TaskStore) LoadAttachmentsLazy(taskID int) *LazyLoader {
	return NewLazyLoader(func() (interface{}, error) {
		task, exists := s.Tasks[taskID]
		if !exists {
			return nil, fmt.Errorf("task not found")
		}
		return task.Attachments, nil
	})
}
// Add to task/store.go

type IndexType string

const (
	IndexByID       IndexType = "id"
	IndexByPriority IndexType = "priority"
	IndexByTag      IndexType = "tag"
	IndexByComplete IndexType = "complete"
)

type Index struct {
	values map[interface{}][]int
	mu     sync.RWMutex
}

func NewIndex() *Index {
	return &Index{
		values: make(map[interface{}][]int),
	}
}

func (i *Index) Add(key interface{}, id int) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.values[key] = append(i.values[key], id)
}

func (i *Index) Remove(key interface{}, id int) {
	i.mu.Lock()
	defer i.mu.Unlock()
	
	ids := i.values[key]
	for j, existingID := range ids {
		if existingID == id {
			i.values[key] = append(ids[:j], ids[j+1:]...)
			break
		}
	}
}

func (i *Index) Get(key interface{}) []int {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.values[key]
}

// Add indexes to TaskStore
type TaskStore struct {
	Tasks        map[int]*Task
	NextID       int
	history      *HistoryStore
	shareStore   *ShareStore
	cache        *TaskCache
	lazyStats    *LazyLoader
	
	// Indexes
	priorityIndex *Index
	tagIndex      *Index
	completeIndex *Index
}

func NewTaskStore() *TaskStore {
	store := &TaskStore{
		Tasks:         make(map[int]*Task),
		NextID:        1,
		cache:         NewTaskCache(5 * time.Minute),
		priorityIndex: NewIndex(),
		tagIndex:      NewIndex(),
		completeIndex: NewIndex(),
	}
	return store
}

// Update Add method to maintain indexes
func (s *TaskStore) Add(title string) (Task, error) {
	if title == "" {
		return Task{}, errors.New("task title cannot be empty")
	}
	
	task := NewTask(s.NextID, title)
	s.Tasks[task.ID] = task
	
	// Update indexes
	s.priorityIndex.Add(task.Priority, task.ID)
	s.completeIndex.Add(task.Completed, task.ID)
	for _, tag := range task.Tags {
		s.tagIndex.Add(tag, task.ID)
	}
	
	s.NextID++
	s.cache.Delete("all_tasks")
	return *task, nil
}

// Fast lookup using indexes
func (s *TaskStore) GetByPriorityIndex(priority Priority) []Task {
	ids := s.priorityIndex.Get(priority)
	var tasks []Task
	for _, id := range ids {
		if task, exists := s.Tasks[id]; exists {
			tasks = append(tasks, *task)
		}
	}
	return tasks
}

func (s *TaskStore) GetByTagIndex(tag string) []Task {
	ids := s.tagIndex.Get(tag)
	var tasks []Task
	for _, id := range ids {
		if task, exists := s.Tasks[id]; exists {
			tasks = append(tasks, *task)
		}
	}
	return tasks
}