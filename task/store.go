package task

import (
	"fmt"
	"strings"
)

func (s *TaskStore) FindByTitle(title string) []Task {
	var results []Task
	titleLower := strings.ToLower(title)
	
	for _, task := range s.Tasks {
		if strings.Contains(strings.ToLower(task.Title), titleLower) {
			results = append(results, task)
		}
	}
	return results
}

func (s *TaskStore) GetStats() map[string]int {
	stats := make(map[string]int)
	stats["total"] = len(s.Tasks)
	
	completed := 0
	for _, task := range s.Tasks {
		if task.Completed {
			completed++
		}
	}
	stats["completed"] = completed
	stats["pending"] = stats["total"] - completed
	
	return stats
}

func (s *TaskStore) DisplayStats() {
	stats := s.GetStats()
	fmt.Println("\n--- STATISTICS ---")
	fmt.Printf("Total Tasks: %d\n", stats["total"])
	fmt.Printf("Completed: %d\n", stats["completed"])
	fmt.Printf("Pending: %d\n", stats["pending"])
	fmt.Println("------------------")
}