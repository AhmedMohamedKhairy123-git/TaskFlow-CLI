package task

type Task struct {
	ID        int
	Title     string
	Completed bool
}

type TaskStore struct {
	Tasks   map[int]Task
	NextID  int
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		Tasks:  make(map[int]Task),
		NextID: 1,
	}
}

func (s *TaskStore) AddTask(title string) Task {
	task := Task{
		ID:        s.NextID,
		Title:     title,
		Completed: false,
	}
	s.Tasks[task.ID] = task
	s.NextID++
	return task
}

func (s *TaskStore) GetTask(id int) (Task, bool) {
	task, exists := s.Tasks[id]
	return task, exists
}

func (s *TaskStore) GetAllTasks() []Task {
	var taskList []Task
	for _, task := range s.Tasks {
		taskList = append(taskList, task)
	}
	return taskList
}

func (s *TaskStore) MarkComplete(id int) bool {
	task, exists := s.Tasks[id]
	if exists {
		task.Completed = true
		s.Tasks[id] = task
		return true
	}
	return false
}

func (s *TaskStore) DeleteTask(id int) bool {
	_, exists := s.Tasks[id]
	if exists {
		delete(s.Tasks, id)
		return true
	}
	return false
}