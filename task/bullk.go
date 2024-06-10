package task

import (
	"sync"
)

type BulkResult struct {
	SuccessCount int
	FailedCount  int
	Errors       map[int]error
}

func (s *TaskStore) BulkAdd(titles []string) BulkResult {
	result := BulkResult{
		Errors: make(map[int]error),
	}
	
	for i, title := range titles {
		_, err := s.Add(title)
		if err != nil {
			result.Errors[i] = err
			result.FailedCount++
		} else {
			result.SuccessCount++
		}
	}
	
	return result
}

func (s *TaskStore) BulkComplete(ids []int) BulkResult {
	result := BulkResult{
		Errors: make(map[int]error),
	}
	
	for _, id := range ids {
		if err := s.MarkComplete(id); err != nil {
			result.Errors[id] = err
			result.FailedCount++
		} else {
			result.SuccessCount++
		}
	}
	
	return result
}

func (s *TaskStore) BulkDelete(ids []int) BulkResult {
	result := BulkResult{
		Errors: make(map[int]error),
	}
	
	for _, id := range ids {
		if err := s.Delete(id); err != nil {
			result.Errors[id] = err
			result.FailedCount++
		} else {
			result.SuccessCount++
		}
	}
	
	return result
}

func (s *TaskStore) BulkUpdate(updates map[int]Task) BulkResult {
	result := BulkResult{
		Errors: make(map[int]error),
	}
	
	for id, updatedTask := range updates {
		if err := s.Update(id, updatedTask); err != nil {
			result.Errors[id] = err
			result.FailedCount++
		} else {
			result.SuccessCount++
		}
	}
	
	return result
}

func (s *TaskStore) ParallelBulkComplete(ids []int) BulkResult {
	result := BulkResult{
		Errors: make(map[int]error),
	}
	
	var wg sync.WaitGroup
	var mu sync.Mutex
	
	for _, id := range ids {
		wg.Add(1)
		go func(taskID int) {
			defer wg.Done()
			
			err := s.MarkComplete(taskID)
			
			mu.Lock()
			if err != nil {
				result.Errors[taskID] = err
				result.FailedCount++
			} else {
				result.SuccessCount++
			}
			mu.Unlock()
		}(id)
	}
	
	wg.Wait()
	return result
}