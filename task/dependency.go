package task

import (
	"fmt"
)

type DependencyType string

const (
	DependsOn     DependencyType = "depends_on"
	BlockedBy     DependencyType = "blocked_by"
	RelatedTo     DependencyType = "related_to"
)

type Dependency struct {
	TaskID     int
	DependsOn  int
	Type       DependencyType
}

func (s *TaskStore) AddDependency(taskID, dependsOnID int, depType DependencyType) error {
	if _, exists := s.Tasks[taskID]; !exists {
		return fmt.Errorf("task %d not found", taskID)
	}
	if _, exists := s.Tasks[dependsOnID]; !exists {
		return fmt.Errorf("dependency task %d not found", dependsOnID)
	}
	
	if s.hasCircularDependency(taskID, dependsOnID) {
		return fmt.Errorf("circular dependency detected")
	}
	
	task := s.Tasks[taskID]
	task.Dependencies = append(task.Dependencies, Dependency{
		TaskID:    taskID,
		DependsOn: dependsOnID,
		Type:      depType,
	})
	
	return nil
}

func (s *TaskStore) hasCircularDependency(taskID, dependsOnID int) bool {
	visited := make(map[int]bool)
	return s.detectCycle(dependsOnID, taskID, visited)
}

func (s *TaskStore) detectCycle(current, target int, visited map[int]bool) bool {
	if current == target {
		return true
	}
	
	if visited[current] {
		return false
	}
	visited[current] = true
	
	task, exists := s.Tasks[current]
	if !exists {
		return false
	}
	
	for _, dep := range task.Dependencies {
		if s.detectCycle(dep.DependsOn, target, visited) {
			return true
		}
	}
	
	return false
}

func (s *TaskStore) GetBlockedTasks() []Task {
	var blocked []Task
	
	for _, task := range s.Tasks {
		if s.isBlocked(task.ID) {
			blocked = append(blocked, *task)
		}
	}
	
	return blocked
}

func (s *TaskStore) isBlocked(taskID int) bool {
	task, exists := s.Tasks[taskID]
	if !exists {
		return false
	}
	
	for _, dep := range task.Dependencies {
		if dep.Type == BlockedBy || dep.Type == DependsOn {
			depTask, exists := s.Tasks[dep.DependsOn]
			if exists && !depTask.Completed {
				return true
			}
		}
	}
	
	return false
}

func (s *TaskStore) GetDependentTasks(taskID int) []Task {
	var dependent []Task
	
	for _, task := range s.Tasks {
		for _, dep := range task.Dependencies {
			if dep.DependsOn == taskID {
				dependent = append(dependent, *task)
				break
			}
		}
	}
	
	return dependent
}